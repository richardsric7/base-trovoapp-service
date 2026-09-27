# Integration Guide

This document is for anyone integrating a new client, a new admin feature,
or a new white-label partner against `app-backend`. It explains how the
three kinds of callers this API already has — the end-user apps, the admin
backend (`tm-api`), and white-label "service link" partners — each
authenticate and talk to it.

For the exact request/response shape of any individual endpoint, the
canonical reference is always the generated **Swagger UI**, not this
document — see [the bottom of this file](#swagger-ui-the-per-endpoint-reference).

## Authentication

`app-backend` has **three** distinct authentication schemes, chosen per
route based on who's expected to call it. There is no cookie-based session
and no classic OAuth bearer flow for end users — every scheme is stateless
and checked on every request.

### 1. Request-signing (`app-web` and `app-mobile`)

This is how the two end-user client apps authenticate almost every
request. There is no login token to attach — instead, **every request is
individually signed** with the calling wallet's private key. The client
sends four custom headers:

| Header | Meaning |
| --- | --- |
| `X-TW-PUBLIC-KEY` | The wallet address making the call |
| `X-TW-SIGNER` | The address whose private key actually signed this request (usually the same as `X-TW-PUBLIC-KEY`, but can differ for shared-access wallets where an approver signs on behalf of the wallet) |
| `X-TW-TIMESTAMP` | A current timestamp, included in what gets signed so a captured request can't be replayed later |
| `X-TW-SIGNATURE` | A base64-encoded EIP-191 `personal_sign` signature |

The signature covers `<full request path + query string> + <X-TW-SIGNER> +
<X-TW-TIMESTAMP>` (see `internal/middleware/authentication_middleware.go`'s
`authenticationChecks` and `internal/middleware/security_checks.go`'s
`SignHttp`/`VerifyHttpSignature` for the exact construction). The server
recomputes the expected signature from the request it received and rejects
the call with `401` if it doesn't match, or if `X-TW-SIGNER` fails basic
address-format validation.

Practically: whichever wallet SDK/library your client uses needs to be able
to sign an arbitrary string with the user's private key (never send the
private key itself to the server). The middleware enforcing this is
`AuthenticationMiddlewareUsingTimestamp()`, applied per-route in each
component's `controllers/main.go`.

In the Swagger UI and generated docs, routes that require this are marked
`@Security SignatureAuth`.

### 2. API key (white-label "service links", and a couple of server-to-server admin calls)

Routes under `/v1/servicelinks/...` and `/v1/trovo-api/...`, plus a small
number of admin-only endpoints elsewhere (e.g. P2P dispute resolution),
authenticate with a single static header instead:

```
X-TW-SERVICE-LINK-API-KEY: <the partner's API key>
```

The server looks this key up against the `service_links` table and rejects
the call with `401` if it's missing, unknown, or the service link is
suspended. See [Service links](#service-links-white-label-partner-integration)
below for how a partner gets one of these keys. The middleware is
`AuthenticationMiddlewareUsingAPIKey()`; these routes are marked
`@Security ServiceLinkApiKey` in the generated docs.

### 3. JWT bearer token (`tm-api`'s own admin routes)

A small set of routes under `/v1/trovo-manager/...` — internal admin
read/write endpoints `tm-api`'s staff UI calls — use a conventional JWT
bearer token instead:

```
Authorization: Bearer <jwt>
```

checked by `JwtTokenAuthMiddleware()`, signed with the `JWT_ACCESS_SECRET`
described in [CONFIGURATION.md](CONFIGURATION.md). These are marked
`@Security BearerAuth` in the generated docs. End-user clients and
white-label partners never use this scheme.

---

## How `tm-api` relates to `app-backend`

`tm-api` is Trovo's internal admin/catalog backend, and it integrates with
`app-backend` in **two different ways** at once:

1. **Direct database access, for admin/catalog data.** `tm-api` reads and
   writes tables like the curated-asset catalog, fee configuration, and
   service-link records directly against the **same physical Postgres
   database** `app-backend` uses. There's no API call involved for this —
   `app-backend`'s `assets`/`rates`/`servicelinks` components largely just
   *read* what `tm-api` wrote.
2. **HTTP calls into `app-backend`'s own API, for P2P transactional
   writes.** For creating orders, resolving disputes, and other P2P
   state-changing actions, `tm-api` calls this service's HTTP API (e.g.
   `POST /v1/p2p/disputes/{disputeID}/admin-resolve`, authenticated with
   `ServiceLinkApiKey`) rather than writing P2P tables directly.

**Why the split exists:** P2P orders and disputes involve real
in-flight, stateful business logic — escrow balances, order-state
transitions, notifications, on-chain settlement — the kind of thing that
needs to go through `app-backend`'s own validation and side effects (see
`internal/components/p2p/controllers/admin_handlers.go`'s own comment: *"app-backend
has no admin/staff user model of its own... this is a server-to-server
endpoint an admin-authorized tm-api call can reach"*). Writing directly to
the P2P tables from `tm-api` would risk leaving escrow state, balances, or
notifications out of sync with what `app-backend` believes is true.
Catalog/config data, by contrast, is comparatively inert — a curated
asset's metadata or a fee percentage doesn't have in-flight state that
another process could race with — so a direct, low-overhead DB write is
fine for that category.

**A consequence for you if you're adding a new admin feature:** if it's
read-only config/catalog data, writing directly to the DB from `tm-api`
(following existing patterns there) is consistent with how this already
works. If it's anything that changes P2P (or similar stateful,
in-flight) records, add a new `app-backend` endpoint instead of writing
those tables directly, and gate it behind `AuthenticationMiddlewareUsingAPIKey`
the same way `postAdminResolveDisputeHandler` does — see
`internal/components/p2p/controllers/admin_handlers.go` as the template.

---

## Service links (white-label partner integration)

A "service link" is how an external partner integrates Trovo Wallet
functionality into their own product under their own brand — a
white-label / embedded-wallet integration. The routes live in
`internal/components/servicelinks/controllers/` and are all under
`/v1/servicelinks/...` and `/v1/trovo-api/...`.

### Getting an API key

A service link is a row in the `service_links` table (owned by a specific
Trovo user account, `OwnerUsername`) carrying a generated `ApiKey`. There
is currently no self-serve signup endpoint for this — a service link
record is provisioned by Trovo (via `tm-api`'s admin tooling, direct DB
access as described above). Once provisioned, the partner receives their
`ApiKey` out of band and sends it as `X-TW-SERVICE-LINK-API-KEY` on every
request (see [Authentication](#2-api-key-white-label-service-links-and-a-couple-of-server-to-server-admin-calls) above).

### What a service link can do

Broadly, three categories of endpoint:

- **User-facing login/authorization handshakes** — `POST
  /v1/servicelinks/login/request/{targetUser}`,
  `POST /v1/servicelinks/authorize/request/{targetUser}`, and their
  `/verify/...` counterparts. These let the partner ask a Trovo user (by
  username) to approve the partner accessing their account, which the user
  confirms from inside the Trovo app itself (`POST
  /v1/users/servicelinks/login/approval/{targetUser}` and
  `.../authorize/approval/{targetUser}`, both signed by the user, not the
  partner). Approved sessions/authorizations are exchanged for short-lived
  tokens (`POST /v1/servicelinks/token/refresh`, `POST
  /v1/servicelinks/token/verify`) the partner then presents on subsequent
  calls.
- **Server-to-server "trovo-api" endpoints** — `/v1/trovo-api/users/...`,
  `/v1/trovo-api/assets/...`, `/v1/trovo-api/tokens/mint`, etc. These let
  an already-authorized partner onboard users, check balances, request
  payments, mint tokens, and manage their own asset listings, all
  authenticated purely by the API key (no per-request user approval).
- **Push/event delivery** — `POST
  /v1/servicelinks/{ownerUsername}/{targetUser}/push` and `POST
  /v1/servicelinks/events/request` let a partner push a notification to a
  user or register for event callbacks.

### The permission-flag model

Not every service link can do everything above. The `ServiceLink` row
(`internal/components/servicelinks/models/servicelink.go`) carries a set of
boolean-ish integer flags — `LoginPermission`, `PaymentPermission`,
`TokenInfoPermission`, `AuthorizationPermission`, `EventPermission`,
`AllowUserInfo`, `PushNotificationPermission`, `IncludePhoneNumbers`,
`IncludeUserBalances`, `TokenizedAssetAuthorizationPermission`,
`CreateUsersPermission`, `AllowReferralForRegisteredUsers`, `Verified`,
`Inactive`, `Suspended` — and each relevant handler checks the specific
flag(s) it needs before proceeding, returning a permission error if the
service link isn't authorized for that action. When you're orienting
yourself in this component for the first time, you don't need to memorize
every flag — just know that **each capability above is individually
gated**, so a partner integration is scoped to exactly what it was
provisioned for, and a new capability you add should follow the same
pattern (add a flag, check it in your new handler) rather than reusing an
existing, differently-scoped flag.

---

## Swagger UI: the per-endpoint reference

Once the server is running, every documented endpoint — request
parameters, body shape, response codes, and which of the three auth
schemes above it requires — is browsable at:

```
http://localhost:8080/swagger/index.html
```

(replace the host/port with wherever the service is actually running). The
underlying spec is generated by [`swaggo/swag`](https://github.com/swaggo/swag)
from `@Summary`/`@Router`/etc. comments directly above each handler
function, regenerated with `swag init --parseDependency=false` (see
`docs/docs.go`'s header and [DEPLOYMENT.md](DEPLOYMENT.md)). If you're
integrating against a specific endpoint, always check its Swagger entry
first — it reflects the current handler code, not a hand-maintained
description that can drift.
