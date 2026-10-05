# Integration with the rest of the monorepo

## Summary

`tm-web` has exactly one integration point with the rest of the system:
it is an HTTP client of **`tm-api`**, a sibling Go REST API in this
monorepo. It has **no direct relationship** with `app-backend`, `app-web`,
or `app-mobile` — those are a separate product line (the consumer-facing
Trovo App and its own backend). `tm-api` is, in turn, the service that talks
to (and reportedly shares a database with) `app-backend`; `tm-web` never
reaches past `tm-api` to get there. This section explains what was checked
to confirm that, and the sections below explain exactly how the `tm-api`
integration works.

Confirmation of the above, from this project's own source:

- `grep -rn "app-backend\|app-web\|app-mobile"` across `tm-web/src` and the
  config files at the project root turns up **no HTTP calls, URLs, or
  imports** referencing those projects — only a handful of source comments
  (`src/redux/api/serviceLinks/interface.ts`, `src/redux/api/curatedAssets/interface.ts`,
  `src/redux/api/p2p/interface.ts`) noting that certain `tm-api` response
  shapes mirror tables/modules that live in `app-backend`'s own database
  schema (e.g. "`service_links` table app-backend's servicelinks module
  reads from"). That's exactly the "`tm-api` shares a DB with `app-backend`"
  relationship described above, documented from `tm-api`'s side of the
  integration for the benefit of whoever is typing these interfaces — it is
  not `tm-web` talking to `app-backend`. `tm-web` itself has no configured
  base URL, client, or import for `app-backend`, `app-web`, or `app-mobile`.
- The only externally-configurable HTTP base URL in this project is
  `NEXT_PUBLIC_API_BASE_URL` (`src/config/index.ts`), used everywhere the
  RTK Query API slices and the auth/token-refresh code build request URLs.
  There is exactly one such base URL — not several per-backend — which is
  consistent with `tm-web` only ever calling one upstream service.
  `tm-api`'s default port (per the comment in `.env.example`, referencing
  `tm-api/internal/server/models/server.go`) is `8082`, distinct from
  `app-backend`'s own port.
- `tm-web` has **zero Next.js Route Handlers** of its own (no `route.ts`
  anywhere under `src/app` — confirmed by `find src/app -iname "route.ts"`
  returning nothing), so it isn't itself acting as a proxy or gateway in
  front of any other backend.

## How `tm-web` calls `tm-api`

### Base URL configuration

`src/config/index.ts`:

```ts
export const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? '';
```

This is the single source of truth for where `tm-api` is. It's consumed by:

- `src/redux/baseApi/axiosBasedQuery.ts` — the custom RTK Query `baseQuery`
  used by the **admin dashboard's** API slice (`baseApi`, in
  `src/redux/baseApi/index.ts`). Despite being wired into RTK Query, every
  request actually goes out via **`axios`**, not `fetch` — `axiosBaseQuery`
  wraps a global `axios` instance and its interceptors (see below) in an
  RTK Query-shaped function.
- `src/redux/baseApi/orgBasedQuery.ts` — the equivalent for the
  **organisation/stakeholder portal's** API slice (`orgApi`). Also
  axios-based, against the same `BASE_URL`, but with an entirely separate
  token/interceptor setup (see [Two independent auth flows](#two-independent-auth-flows-admin-vs-organisation-portal) below).
- `src/redux/refreshAccessToken.ts` — the token-refresh call, made with a
  bare `axios` instance deliberately kept outside the interceptor chain (see
  below).
- `src/app/(auth)/qr-scan/page.tsx` — opens a raw `EventSource` (SSE) against
  `${BASE_URL}/login/stream/:loginId` to get a push notification the instant
  a login is approved on the mobile app (used alongside, not instead of, a
  5-second poll).

Every RTK Query endpoint definition under `src/redux/api/**` (`auth`,
`users`, `admin`, `p2p`, `payment`, `assettokenization`, `vaultSigner`,
`compliance`, `auditTrail`, `jobs`, `organizations`, etc.) is `injectEndpoints`
onto `baseApi` (or, for the organisation portal's endpoints, `orgApi`), so
all of them inherit this same base URL and the same request/response
interceptor behavior automatically.

### Authentication flow (admin dashboard)

Admin login in `tm-web` is **not** a username/password form — it's a
QR-code flow that authenticates the admin's session against their approval
on the Trovo mobile app. The relevant files are
`src/app/(auth)/sign-in/page.tsx`, `src/app/(auth)/qr-scan/page.tsx`,
`src/redux/api/auth/api.ts`, `src/redux/slices/authSlice.ts`, and
`src/middleware.ts`.

1. **`/sign-in`** — the admin enters their Trovo App username. This calls
   `POST /login` (`useLoginMutation`, `src/redux/api/auth/api.ts`), which
   returns a `dynamicLink`, `loginId`, and a `qrCode` image. That response is
   stashed in a short-lived cookie (`CookieType.LoginResults`) and the user
   is routed to `/qr-scan`.
2. **`/qr-scan`** — renders the QR code and both polls and listens for the
   login to be approved:
   - Every 5 seconds it calls `GET /users/verify/:targetUser/:loginId`
     (`useVerifyLoginQuery`) to check whether the mobile app has approved the
     pending login.
   - In parallel, it opens an `EventSource` against
     `${BASE_URL}/login/stream/:loginId` so an approval is picked up
     instantly rather than waiting out the poll interval (particularly
     useful since background-tab throttling can delay the poll).
3. Once `verifyLogin` returns data, the response contains
   `accessToken`, `refreshToken`, and `userInfo`. The page:
   - `dispatch(setToken({ accessToken, refreshToken, expiresAt }))` —
     stores the token pair in the `auth` Redux slice **and** mirrors it into
     `localStorage` (`authSlice.ts`'s `setToken` reducer calls
     `localStorage.setItem(TOKEN, ...)` as a side effect). `expiresAt` is
     computed client-side as `Date.now() + 15 minutes` — the token itself is
     never decoded client-side to read a real expiry.
   - `setCookie(TOKEN, accessToken)` — **also** writes the raw access token
     into a browser cookie (via `cookies-next`), separately from the
     Redux/localStorage copy. This cookie is what `src/middleware.ts` (which
     runs server-side/at the edge, and cannot read `localStorage`) checks to
     decide whether a request is authenticated, before it even reaches a
     page component.
   - `dispatch(setUser(userInfo))` — stores basic user info, also mirrored
     to `localStorage`.
   - Routes to `/` and clears the temporary `LoginResults` cookie.

**Attaching the token to requests**: a global `axios.interceptors.request.use`
in `src/redux/baseApi/axiosBasedQuery.ts` reads the access token out of
Redux state (`getPreloadedState().auth.token`) on every outgoing request and
sets `Authorization: Bearer <accessToken>` — except for a few endpoint
families (tokenization, wallet-balances, organizations) that the backend
expects to receive the **raw token without the `Bearer ` prefix**, which the
interceptor special-cases by URL substring match.

**Refreshing an expired token**: a global `axios.interceptors.response.use`
watches for `401` responses. On the first 401 for a given request, it calls
`refreshAccessToken()` (`src/redux/refreshAccessToken.ts`), which:

- Is single-flight — concurrent 401s during the same expiry window share one
  in-flight refresh call rather than firing several.
- Uses a **bare `axios` instance with no interceptors attached**, calling
  `POST ${BASE_URL}/login/token/refresh` with the stored `refreshToken` —
  deliberately bypassing the global interceptors so a refresh call can never
  itself trigger another refresh.
- On success, updates `localStorage` with the new token pair and retries the
  original request (and any others queued behind it) with the new access
  token.
- On failure, clears the stored session (`localStorage` and the `TOKEN`
  cookie) and redirects to `/sign-in?reason=unauthorized`.

**Route protection**: `src/middleware.ts` runs on (almost) every request and
redirects to `/sign-in` if the `TOKEN` cookie is absent for a non-public
route (and to `/organizations/login` for `/organisation/**` routes missing
the separate `TOKEN_ORG` cookie — see below). This is why the QR-scan
success handler sets *both* the Redux/localStorage copy of the token *and*
a plain cookie: the cookie is what the edge middleware can actually see.

### Two independent auth flows: admin vs. organisation portal

`tm-web` actually contains two separately-authenticated areas that happen to
share one Next.js app and one `tm-api` base URL:

- **Admin dashboard** (`src/app/(dashboard)/**`) — internal Trovo staff,
  authenticated via the QR/wallet flow above, token stored under the `TOKEN`
  key, requests sent with a standard `Authorization: Bearer <token>` header
  via `baseApi`/`axiosBasedQuery.ts`.
- **Organisation portal** (`src/app/organisation/**`) — external partner
  users (asset custodians, trustees, issuing houses, etc.), authenticated
  separately, token stored under a distinct `TOKEN_ORG` key/cookie, and
  requests sent via `orgApi`/`orgBasedQuery.ts` with
  `Authorization: <token>` **without** the `Bearer ` prefix (per a comment
  in that file: "Organization endpoints use raw token based on your Swagger
  test"). A 401 from an org endpoint redirects to `/organizations/login`; a
  401 from an admin endpoint deliberately does **not** trigger that
  redirect while the user is anywhere under `/organisation` (see the guard
  comment in `handleNavigation`, `axiosBasedQuery.ts`), so the two sessions
  can't eject each other.

Both flows still ultimately call the same `tm-api`, just as different
authenticated principals hitting different endpoint groups.

## Public Markets pages

`/publicmarkets/*` covers tokenized NGX equities and FMDQ bonds. It has a
sidebar group of its own and the section's tabs at the top of each page.
Every page reads tm-api's `/public-markets` endpoints through
`src/redux/api/publicMarkets`.

- **Writes need permissions.** Reads need any Trovo admin. Writes need
  `MANAGE_PUBLIC_MARKETS`, and settings and partners need
  `MANAGE_SETTINGS`. `GET /public-markets/me` tells the pages what to show:
  they hide the buttons the admin cannot use, and only listed approvers see
  the approve buttons.
- **Some work is done by app-backend's engine.** Resuming an asset, running
  reconciliation, a mock position feed and paying a balance back are
  requests that app-backend's engine picks up. Recorded partner outcomes
  and declared corporate actions are applied by the engine within seconds.
  The pages show what is still "waiting for the engine". Corporate
  Actions › "Recorded in Trovo Manager" lists each entry with its result.
- **Secrets are shown once.** A new exchange's signing secret, and a
  rotated one, appear once in a dialog and are never shown again.
- **Personal data is masked.** Exchange customers' legal names and tax IDs
  arrive masked unless the admin holds `VIEW_PUBLIC_MARKETS_PII`.
- **System Health** also lists the engine's background jobs.

## Canonical per-endpoint reference

This document describes the *shape* of the integration (base URL, auth
headers, token lifecycle) but is **not** an endpoint-by-endpoint API
reference. For the authoritative, current list of every route, request/
response schema, and status code `tm-web` can call, use **`tm-api`'s own
Swagger/OpenAPI UI**, expected at:

```
<tm-api base URL>/swagger/index.html
```

e.g. `http://localhost:8082/swagger/index.html` for a local `tm-api`
started per `tm-api/DEPLOYMENT.md`. If that path isn't live yet in a given
`tm-api` checkout, treat this note as forward-looking — the Swagger UI is
expected to be the single source of truth for endpoint contracts once set
up there, superseding any endpoint URL strings that happen to be quoted in
this document or in `tm-web`'s own source.
