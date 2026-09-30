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

## Payment history: source vs. destination

`GET /v1/users/payments/:targetAddressForHistory` and
`GET /v1/trovo-api/users/payment-history/:walletAddress` (the service-link
equivalent) both return each row split into a **source** side (what left
`fromAddress`) and a **destination** side (what arrived at `toAddress`),
replacing the old flat `assetCode`/`contractAddress`/`amount` fields:

```json
{
  "sourceNetwork": "base",
  "sourceContractAddress": "0x...",
  "sourceAssetCode": "USDC",
  "sourceAmount": "100.0000000",
  "destinationNetwork": "base",
  "destinationContractAddress": "0x...",
  "destinationAssetCode": "USDC",
  "destinationAmount": "100.0000000"
}
```

For a plain payment (everything today, since Base has no live on-chain
swap integration yet) the two sides are identical — that's not a
placeholder, it's the correct representation of "the same asset left and
arrived." A row where `sourceAssetCode`/`sourceNetwork` differ from
`destinationAssetCode`/`destinationNetwork` is, by definition, a swap.
`network` is currently always `"base"` on both sides; it's a real column
(not inferred/hardcoded client-side) so a future second network or
Base-native bridge doesn't require another breaking response-shape
change — just a different value here.

The read endpoints' query filters follow the same split:
`destinationAssetCode`/`destinationContractAddress`/`destinationAmount`
and `sourceAssetCode`/`sourceContractAddress`/`sourceAmount` (a
`min%max` range for the amount filters), replacing the old flat
`assetCode`/`contractAddress`/`amount` filter names.

Both of these endpoints, along with every other API-key-authenticated
service-link route (`/v1/servicelinks/...`, `/v1/trovo-api/...`),
`POST /v1/users/payment`, `POST /v1/shared-access/payment`,
`POST /v1/users/swap`, `POST /v1/shared-access/swap`, and
`GET /v1/users/:targetUser`, are rate-limited (see
[Rate limiting](CONFIGURATION.md#rate-limiting) in `CONFIGURATION.md`) —
a request over the limit gets `429 Too Many Requests` with a
`Retry-After` header. A service-link partner with unusually high (or low)
traffic needs can get a per-partner limit override instead of the
route's shared default — see CONFIGURATION.md's "Per-service-link
override", set from tm-api's Service Links admin page, not an env var.

---

## Live P2P updates over the websocket

`GET /v1/users/websocket/:identifier` is a websocket endpoint every
authenticated user can connect to for live, in-app updates. Once
connected, it pushes JSON messages shaped like:

```json
{ "stream": <payload>, "streamType": "<type>" }
```

The `streamType` you care about for live P2P updates is **`p2pEvent`** —
you'll get one of these whenever something happens on an order/offer/
escrow/dispute you're a party to (accepted, payment confirmed, dispute
opened, etc.), with the same title/body/data your push notification for
that event carries:

```json
{
  "streamType": "p2pEvent",
  "stream": {
    "title": "Order accepted",
    "body": "Your order was accepted by the merchant",
    "data": { "route": "orderDetail", "orderId": "..." }
  }
}
```

**How to connect:** open the websocket connection the same way you'd open
any other (`wss://.../v1/users/websocket/<your-wallet-address>`), send the
initial subscription handshake the endpoint's Swagger entry describes,
and then just listen — `p2pEvent` messages arrive automatically for as
long as the connection stays open. There is nothing else to subscribe to
per-order; every event for orders/offers you're involved in reaches this
one connection.

**Important: this is a live nudge, not a source of truth.** If your app
has no open connection at the moment an event happens (app closed,
connection dropped, etc.), that specific message is simply never
delivered — there's no catch-up/replay. Treat a `p2pEvent` message as a
hint to re-fetch the relevant order/offer from its REST endpoint, and
always fall back to the existing push notification (which you already
handle) plus normal polling/pull-to-refresh for anything you can't afford
to miss. This is the same tradeoff the admin panel's (`tm-api`) login
notification stream makes, and it's deliberate: the REST API and the
database are always the source of truth, this socket is purely a "check
now" signal.

If you're the first client integrating this (`app-web`/`app-mobile`
haven't wired it up yet as of this writing), the backend side is fully
built and tested — you just need to open the connection and handle
`streamType: "p2pEvent"` messages as described above.

---

## Wallets and sending: Safe UserOperations

Every user wallet is a **Safe** (v1.4.1 with the ERC-4337 Safe4337Module)
owned by the user's key. Its address is fixed before it is deployed, so
it receives funds from registration on; its first send deploys it
("activates" it) in the same operation.

**Registration.** The app computes the user's address with wallet-core's
`primarySafeAddress(signer, "0")` and sends it as `X-TW-PUBLIC-KEY` (the
signer as `X-TW-SIGNER`). `POST /v1/users` refuses an address that is not
that Safe. Service-link onboarding may omit `publicKey`; it is derived.

**Sending (two calls, unchanged shape).**

1. The first call (e.g. `POST /v1/users/payment` without
   `transactionSignature`) validates the request, builds the operation
   and returns it as `transaction`: **base64 of the 32-byte SafeOp hash**,
   plus `messages` (including the network fee, and a note when the send
   also activates the wallet).
2. The app signs it exactly as before - `signBase64Txn` base64-decodes
   `transaction` and `personal_sign`s the bytes - and calls again with the
   same `transaction` and the base64 `transactionSignature`. The backend
   submits the stored operation (it is never rebuilt) to the bundler and
   returns its `userOpHash` as `transactionId`.

Operations expire (`WALLET_OPERATION_VALIDITY`, 10 minutes by default);
an expired or already-submitted `transaction` is refused. Shared wallets
with approvers get an approval request instead, and their operations
stay signable for `SHARED_WALLET_OPERATION_VALIDITY`.

**Gas.** The wallet pays its own gas: in the stablecoin its owner picked
(`User.gasFeeAsset`, one of the curated assets with `gasFeeEligible`)
through the paymaster when the wallet holds some, otherwise in ETH. A
payment that would leave too little to pay the fee in the same asset is
refused with `error-insufficient-for-network-fee`.

**Sub-wallets** (`POST /v1/users/subwallet`, and the service-link
`POST /v1/trovo-api/users/subwallet`). Each sub-wallet is its own Safe,
deployed by the **primary** wallet, so only the primary wallet signs:

1. First call with `walletType` (0 normal, 1 issuing, 2 market making, 3
   bulk payment), `walletTag`, `walletDescription`. The backend picks the
   new Safe's address (the user's signer as owner, a random salt; types 2
   and 3 also get the platform's co-signer as a second owner, threshold
   1). An issuing wallet - or a request with any `linkedWalletAddress` -
   also gets a linked distribution Safe. The response carries `publicKey`
   (the new wallet), `linkedWalletAddress`, `transaction` (the primary
   wallet's operation that deploys them, seeds them with ETH and pays the
   creation fee - set in USD, charged in its stablecoin), `messages` and
   `feeAmount`/`feeCode` (the fee in that stablecoin). A `publicKey`
   sent by the client is ignored.
2. Second call with `transaction` and `primarySignature`. The wallets are
   recorded and the operation submitted; they are marked activated when it
   is mined. `subWalletSignature` / `linkedWalletSignature` /
   `channelAccountSignature` are no longer used (`subWalletMustSign` and
   `linkedWalletMustSign` are always 0).

If the primary wallet is not activated yet, the same operation activates
it too.

**Shared wallets.** Shared access is enforced by the wallet's Safe:
its owners are the wallet owner's key plus every APPROVER's key, and its
threshold is `numberOfApprovalsNeeded` (1 - the owner alone - when there
are no approvers). VIEW-ONLY and INITIATOR permissions stay in the
database. An issuing wallet's linked distribution wallet is deployed with
the issuing wallet's Safe enabled as a module, so the same operation that
changes the issuing wallet's owners changes the distribution wallet's too
(one set of signatures mirrors both).

- *Enable / modify / disable* keep their request shapes. While the owner
  alone controls the wallet, the owner signs the returned `transaction`
  (a wallet operation when approvers change, or a statement of the change
  when only view/initiator permissions change). Once the wallet has
  approvers, the initiator previews (`commit: 0`), then submits the change
  for approval (`commit: 1` with the same `transaction`).
- *Approvals* (`POST /v1/shared-access/approval/:ID`; `DELETE` rejects):
  the first call returns the
  request's `transaction` (base64 of the operation hash or statement);
  each approver signs it with `signBase64Txn`. A signature is checked
  against the approver's key, which must be one of the Safe's owners, and
  the Safe's current threshold decides how many approvals it takes. The
  last approval submits the operation with the collected signatures.
  Payments and crypto withdrawals from shared wallets work the same way.
- Operations waiting for approvers use their own EntryPoint nonce key and
  stay signable for `SHARED_WALLET_OPERATION_VALIDITY` (24h by default), so
  several can be pending on one wallet. Rejecting a request retires its
  operation.

**Gas-fee asset.** `GET /v1/users/settings/gas-fee-assets` lists the
curated assets users can pay network fees in (`gasFeeEligible`, set from
tm-web) and the user's current choice; `PUT /v1/users/settings/gas-fee-asset`
with `{"assetCode": "USDC"}` sets it (`""` clears it - fees in ETH). The
user profile carries it as `gasFeeAsset`.

**Other sends built the same way** (two calls, base64 SafeOp hash to sign):

- *Crypto withdrawal* (`POST /v1/crypto/withdrawals` and the
  shared-access variant): crypto deposits are minted as Trovo tokens, so
  a withdrawal **burns** `amountSubmitted` from the wallet (`burn(uint256)`
  of an OpenZeppelin `ERC20Burnable` token contract) and records the
  withdrawal request for payout. The first response carries `messages`.
- *Patron subscription*: paid from the primary wallet to the patron fee
  wallet (price + VAT). Payable in the dollar (`DOLLAR_ASSET`) or naira
  (`NAIRA_ASSET`, converted with the USD/cNGN rate) stablecoin; other
  payment assets need a DEX and are refused with
  `error-payment-asset-not-supported`.
- *P2P escrow deposit from the depositor's own wallet* goes through the
  payment flow; the deposit is recorded as pending under the operation's
  `userOpHash` and applied to the order only once it is mined, under the
  mined transaction's hash (the reconciliation sweep finishes deposits
  that take longer than the request).

Swaps still need a DEX route and are not available yet. (Closed-group
creation had no endpoint; its unused transaction builder was removed.)

**What app-backend calls.** The bundler (`BUNDLER_URL`:
`eth_estimateUserOperationGas`, `eth_sendUserOperation`,
`eth_getUserOperationReceipt`) and the paymaster quote service
(`PAYMASTER_QUOTE_SERVICE_URL`, `POST /v1/quote` with `X-API-Key`; see
[`paymaster/INTEGRATION.md`](../paymaster/INTEGRATION.md)). A background
loop follows submitted operations to inclusion (`wallet_operations`
table) and marks wallets activated.

## Tokenized assets: token contract vs issuing Safe

On Base a tokenized asset has two different addresses, and nothing in the
API treats one as the other:

| Field | What it is | Used for |
|---|---|---|
| `contractAddress` | The asset's deployed **B20 token contract** - the asset itself | Every balance, transfer, purchase, wallet authorization, curated-asset listing, subscription/interest record, payout and deep link (`?action=tokenizedAsset&assetCode=...&contractAddress=...`) |
| `issuingWalletAddress` | The asset's **issuing Safe**: a Safe multisig that owns the token contract, is its only minter, and holds unsold supply (treasury). `marketMakingWallet` / `walletToHoldAssetsNotForSale` are the same Safe | Minting, and as the issuer that approves wallet authorizations |

Lifecycle, from the admin panel's point of view:

1. **Issuing Safe** - assigned automatically when the asset's details are
   submitted/vetted: a Safe is deployed through the canonical SafeProxyFactory
   with `TOKENIZATION_ISSUING_SAFE_SIGNERS` as owners (see
   [CONFIGURATION.md](CONFIGURATION.md#tokenization)) and recorded as a
   wallet of the `TOKENIZATION_ISSUING_PROFILE` user.
2. **Token contract** - operations deploy the asset's B20 token with the
   issuing Safe as owner (Ownable) or `MINTER_ROLE` holder (AccessControl),
   the asset code as `symbol()` and zero supply, exposing
   `mint(address,uint256)`. Then register it:

   ```
   PUT /v1/trovo-manager/tokenization/contract/:tid
   {"contractAddress": "0x..."}
   ```

   Allowed for minting initiators/approvers until the asset is minted. The
   backend verifies on-chain that the contract exists, its symbol equals the
   asset code, `decimals()` works, `totalSupply()` is 0, and the issuing Safe
   can mint - otherwise it answers `400 error-invalid-token-contract` with
   the reason. (tm-api proxies this as `PUT /tokenization/contract/:id`;
   tm-web's tokenization Wallets panel has the form.)
3. **Mint** - `POST /v1/trovo-manager/tokenization/mint/:tid` creates a
   `TOKENIZE ASSET` approval on the issuing Safe's wallet. Its
   `transaction` is the base64 of a plain-text statement of the exact Safe
   call (`Safe 0x.. on chain N executes 0xToken.mint(0xSafe, amount) and
   0xToken.mint(0xFeeWallet, fee)`), which approvers sign in the app. On the
   final approval the backend re-derives that statement - it must match
   what was signed - checks the token still has zero supply, and has the
   Safe execute both mints atomically. A retried approval can never mint
   twice. Rejecting the approval keeps the Safe and contract and resets
   the approvers for a new request.

Purchases, early exits and payouts refuse assets whose token contract is
not registered (`error-token-contract-not-registered`).

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
