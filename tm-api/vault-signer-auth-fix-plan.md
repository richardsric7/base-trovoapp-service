# Vault Signer Auth Fix — Plan

## 0. Symptom

Users clicking into the Vault Signer self-service pages (`/vault-signer` on the dashboard,
specifically "My Signer Keys" and "Personal Secrets") get logged out — a hard redirect back to
sign-in. This happens for every Trovo Admin regardless of role, including SuperAdmins. Other
dashboard modules (Admin Management, Settings, etc.) are unaffected. Organisation-portal users
are unaffected.

## 1. Root cause

`internal/middleware/authentication_middleware.go` has two different ways of pulling the JWT out
of the `Authorization` header, and they are inconsistently applied:

- `AuthenticateSuperAdmin` and `JwtTokenAuthMiddleware` both call `ExtractToken(c.Request)`
  first (`internal/middleware/extractToken.go`), which strips a `"Bearer "` prefix if one is
  present, before handing the raw token to `sl.JwtTokenVerify(...)`.
- `tryTrovoAdminAuth` — the function both `AllowOrgOrTrovoAdmin` and
  `AllowOrgOrTrovoAdminNormalized` call to check the Trovo-Admin branch — does **not**:
  ```go
  // authentication_middleware.go:232
  jwtResponse, err := sl.JwtTokenVerify(authHeader)   // authHeader passed straight through, unstripped
  ```

`JwtTokenVerify` itself always re-adds its own `"Bearer "` prefix before forwarding to the
external Trovo verification service:
```go
// internal/trovosdk/methods.go:26
Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken)).
```

So `tryTrovoAdminAuth` only works correctly when the caller sends the **raw token with no
`Bearer ` prefix at all**. If the incoming header is already `Bearer <token>`, it's forwarded
verbatim and `JwtTokenVerify` wraps it a second time — the external call ends up sending
`Authorization: Bearer Bearer <token>`, which fails verification. This fails at token
*verification*, before any role check runs, which is why it fails identically for every admin
role, including SuperAdmin.

## 2. Why this is scoped to vault-signer specifically

The frontend already has an undocumented workaround for this exact bug, but only for three URL
patterns. In `src/redux/baseApi/axiosBasedQuery.ts`'s request interceptor:
```js
const isTokenizationEndpoint = config.url?.includes("/tokenization");
const isWalletBalanceEndpoint = config.url?.includes("/wallet-balances/");
const isOrganizations = config.url?.includes("/organizations/");
Authorization:
  isTokenizationEndpoint || isOrganizations || isWalletBalanceEndpoint
    ? token.accessToken                // raw, no prefix — the workaround
    : `Bearer ${token.accessToken}`,   // every other module
```
`/organizations/*` routes through the same `AllowOrgOrTrovoAdminNormalized` middleware, so it was
special-cased at some point to send the raw token and dodge this bug. Vault-signer's self-service
endpoints (`/me/vault-signer/secrets*`, `/me/vault-signer/personal-envs*`) go through the same
middleware but were never added to that special-cased list, so they send the normal
`Bearer <token>` format like most of the app — and hit the bug.

| Area | Middleware | Bug present? |
|---|---|---|
| `/me/vault-signer/secrets*`, `/me/vault-signer/personal-envs*` (My Signer Keys **and** Personal Secrets, dashboard) | `AllowOrgOrTrovoAdminNormalized` → `tryTrovoAdminAuth` | **Yes** |
| `/admin/vault-signer/*` (Settings → Vault Signer: managed secrets/assignments/audit log) | `AuthenticateSuperAdmin` | No — already uses `ExtractToken` correctly |
| Organisation portal (`orgApi`) calls, any endpoint | `orgApi` always sends the raw token, never `Bearer`-prefixed | No |

This is not limited to Personal Secrets — My Signer Keys hits the identical failure, since it's
the same middleware and the same header format. Both dashboard self-service panels are affected
equally.

## 3. Fix options

### Option 1 — Fix the backend middleware (recommended)

Change `tryTrovoAdminAuth` (`authentication_middleware.go:222-249`) to call
`sl.JwtTokenVerify(ExtractToken(c.Request))` instead of `sl.JwtTokenVerify(authHeader)`.
`ExtractToken` already handles both `"Bearer X"` and bare `"X"`, so this is backward-compatible
with every existing caller — including `/organizations/*`, which currently sends the bare form —
and fixes the defect at its source. No frontend module ever needs to know about this quirk
again.

For full consistency, apply the same change to `tryOrganizationAuth`'s two `authHeader` uses
(`jwt.ParseWithClaims(authHeader, ...)` and `jwt.Parse(authHeader, ...)`) — not the trigger for
this specific report (org tokens are never sent with a `Bearer ` prefix today), but the same
latent bug in principle, and cheap to close alongside the admin-branch fix while this file is
already open.

**Blast radius:** shared authentication middleware used by `/organizations/*` in addition to
vault-signer. Low risk given `ExtractToken`'s backward compatibility with the raw-token format
those callers already send, but it's auth code — worth a deliberate look/test pass, not a
drive-by change.

**Files touched:** `internal/middleware/authentication_middleware.go` only.

### Option 2 — Frontend-only workaround

Add `/me/vault-signer` to the existing special-cased URL list in
`src/redux/baseApi/axiosBasedQuery.ts` (frontend), alongside `/tokenization`, `/wallet-balances/`,
`/organizations/`, so those specific calls send the raw token without a `Bearer ` prefix.

**Trade-off:** minimal, no backend deploy required, mirrors the existing (if undocumented)
pattern exactly — but leaves the real bug in place. The next module built behind
`AllowOrgOrTrovoAdminNormalized` will hit this exact same failure with no clue why, the same way
vault-signer just did.

**Files touched:** `frontend/src/redux/baseApi/axiosBasedQuery.ts` only.

### Recommendation

Option 1 alone, since it's the durable fix and is backward-compatible with everything already
working. Option 2 can be layered on top as an immediate belt-and-suspenders patch if a same-day
frontend-only unblock is wanted while the backend fix goes through review — the two are not
mutually exclusive and don't conflict.

## 4. Verification plan (once a direction is approved)

- `go build ./...`, `go vet ./...`, `gofmt -l` on the touched middleware file.
- Trace both branches of `AllowOrgOrTrovoAdminNormalized` by hand against the fix: a Trovo Admin
  token sent as `Bearer <token>` (vault-signer's current dashboard behavior) and a bare `<token>`
  (the existing `/organizations/*` behavior) should both now reach `tryTrovoAdminAuth` correctly.
- If Option 2 is included: confirm the frontend still builds/type-checks/lints clean
  (`tsc --noEmit`, `next lint`) with the URL-matching change.
- No live backend/Trovo service-link instance is available in this environment, so the actual
  external `JwtTokenVerify` round-trip cannot be exercised end-to-end here — recommend a manual
  smoke test (a real Trovo Admin login, then opening My Signer Keys / Personal Secrets) before
  merging.

## 5. Status

**Implemented: Option 1, on branch `ric-vault-signer-auth-fix`.** Option 2 (the frontend
workaround) was not applied — with the root cause fixed, it isn't needed for vault-signer, and
adding it would just be one more special case to maintain.

### What changed

`internal/middleware/authentication_middleware.go`:

- **`AllowOrgOrTrovoAdmin`** and **`AllowOrgOrTrovoAdminNormalized`** now compute
  `token := ExtractToken(c.Request)` once (after the existing empty-header guard) and pass that
  — not the raw `c.GetHeader("Authorization")` value — into both `tryOrganizationAuth` and
  `tryTrovoAdminAuth`.
- **`tryTrovoAdminAuth`**'s parameter is renamed `authHeader` → `token` (it now genuinely holds
  the bare token) and its call becomes `sl.JwtTokenVerify(token)` — this is the actual bug fix:
  previously it received the raw header, which for a `Bearer`-prefixed caller (vault-signer's
  self-service endpoints) meant `JwtTokenVerify` re-added its own `"Bearer "` and sent
  `"Bearer Bearer <token>"` to the external Trovo verification service, failing before any role
  check.
- **`tryOrganizationAuth`**'s parameter is likewise renamed `authHeader` → `token`, matching
  Section 3's note that it carries the identical latent issue in principle (not the trigger for
  this report, since `orgApi` never sends a `Bearer`-prefixed token today, but now closed
  defensively at the same time while this file was already open). Its internal parsed-JWT
  variable was renamed `token` → `parsedToken` to avoid shadowing the (now same-named) input
  parameter — a pure rename, no logic change.

`ExtractToken` was already correct and unchanged (`internal/middleware/extractToken.go`) — it
already handled both `"Bearer X"` and bare `"X"` forms, which is exactly why routing both
`AllowOrgOrTrovoAdmin` helpers through it is backward-compatible with every existing caller,
including `/organizations/*` (which sends the bare form today and continues to work identically).

### Verification performed

- `go build ./...` — clean, whole repo.
- `go vet ./...` — clean, whole repo.
- `gofmt -l internal/middleware/authentication_middleware.go` — clean.
- `go test ./internal/middleware/...` — passes (`ok`). Neither existing test file in this
  package (`organization_session_test.go`, `stakeholder_auth_middleware_test.go`) exercises
  `tryOrganizationAuth`/`tryTrovoAdminAuth`/`AllowOrgOrTrovoAdmin(Normalized)` directly, so this
  confirms no regression in what *was* covered, not new coverage for the fixed functions
  themselves — this repo's admin-DB test suite (`go test ./...` broadly) still can't build in
  this sandbox for the pre-existing, unrelated `sqlite`-cgo reason noted in
  `vault-signer-service-plan.md`'s own verification section.

### Verification not performed, and why

No live backend or reachable Trovo service-link instance in this sandbox, so the actual external
`JwtTokenVerify` round-trip was traced by hand against the code (Section 1's trace, now
confirmed correct on both branches) rather than exercised live. **Recommended before merging:** a
manual smoke test — real Trovo Admin login, then open My Signer Keys and Personal Secrets on
`/vault-signer` — to confirm the 401 is gone end to end, plus a quick check that
`/organizations/*` (Admin Management, org listing, etc.) still authenticates correctly for both
Trovo Admins and Organization members, since that's the one already-working path this change
routes through the same new code.
