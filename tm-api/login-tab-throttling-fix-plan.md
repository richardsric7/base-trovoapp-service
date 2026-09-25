# Login Tab-Throttling Fix — Plan

## 0. Symptom

When a user starts the QR login flow and then switches to another browser tab before scanning
or before the wallet approval lands, the web app appears to "pause" — it doesn't pick up the
successful login until the user switches back to the tab.

## 1. Investigation — confirmed, two contributing causes

`useVerifyLoginQuery` (`frontend/src/app/(auth)/qr-scan/page.tsx`) polls
`GET /users/verify/:targetUser/:loginID` every 5 seconds (`pollingInterval: 5000`) while waiting
for wallet approval. That polling is timer-based (`setInterval`/`setTimeout` under RTK Query),
and two things compound to produce the reported pause:

1. **Browser background-tab timer throttling — not app-specific, not fixable by disabling
   anything in this codebase.** Every modern browser throttles JS timers in backgrounded/inactive
   tabs for battery and performance reasons, increasingly aggressively the longer a tab stays
   hidden (Chrome, for example, can clamp a background tab's timers to firing about once a
   minute). Checked for app-side causes first — grepped for `visibilitychange`, `document.hidden`,
   `requestAnimationFrame` across the frontend: **nothing**. This isn't a deliberate pause
   feature; it's the browser engine's own policy acting on a plain `setInterval`.
2. **`setupListeners(store.dispatch)` is never called** (`frontend/src/redux/store.ts`). This is
   RTK Query's own built-in mitigation for exactly this scenario: wired up, it listens for the
   window regaining focus/visibility and, combined with `refetchOnFocus: true` on a query,
   triggers an immediate refetch the moment the tab is revisited. Without it, nothing proactively
   catches up when the user returns — the app just waits for the browser's already-throttled timer
   to fire again on its own schedule, which can add real time on top of the throttling itself.

### A related discovery: half of a push-based fix already exists, unused

While tracing the approval path, found `GlobalConfig.LoginUsers` and
`BroadcastToLoginID(loginID, streamType, payload)` (`internal/models/streams.go`) — and
confirmed this is **live, not dead, code**: the real webhook Trovo Wallet calls on approval,
`LoginCallback` (`POST /v1/callbacks/login/:serviceName`, registered in
`internal/components/auth/controllers/main.go`), already calls
`s.GC.BroadcastToLoginID(callbakInput.LoginID, "loginNotice", pl)` the instant approval happens —
right next to the cache-invalidation call from the previous fix (`login-verify-cache-fix-plan.md`).

But nothing ever *subscribes*. Confirmed by grep: `LoginUsers[loginID]` is only ever read inside
`BroadcastToLoginID` itself (`if conChan, ok := n.LoginUsers[loginID]; ok`) — no code anywhere
registers a channel into that map, and there is no WebSocket library or SSE
(`text/event-stream`) usage anywhere in this backend. Every real approval today calls
`BroadcastToLoginID` and it silently does nothing, because `ok` is always `false`. The "publish"
half of a push-based login flow already exists and fires correctly on the real webhook; the
"subscribe" half (an endpoint the frontend could actually connect to) was never built.

This matters for the fix: an open server-push connection isn't subject to the same
timer-throttling that breaks polling in a backgrounded tab — the browser throttles JS *timers*,
not incoming data on an already-open connection — so completing this mechanism fixes the
tab-switch pause at its root, not just mitigates it.

## 2. Fix options discussed

- **Option 1 (small, unconditionally worth doing):** wire up `setupListeners(store.dispatch)` and
  opt the login-verify query into `refetchOnFocus: true`. Doesn't stop throttling while
  backgrounded, but removes the compounding "and then also wait for the stale timer" delay on
  return — the check reruns the instant focus returns.
- **Option 2 (not chosen):** move the polling into a dedicated Web Worker, which browsers throttle
  less aggressively than the main thread.
- **Option 3 (chosen, the real fix):** finish the push mechanism already half-built —  add an SSE
  endpoint that subscribes a connection into `LoginUsers[loginID]`, and have the frontend open
  that stream alongside its existing (now-fixed, 3-second-cache) polling. The instant
  `BroadcastToLoginID` fires, the open connection delivers it regardless of tab focus, and the
  frontend immediately re-checks for the real token payload rather than waiting on any timer.

**Both Option 1 and Option 3 are implemented below** — Option 1 as a small, unconditional safety
net; Option 3 as the actual fix; polling is kept as a fallback in both cases rather than removed,
so nothing regresses if a client's network path doesn't support a long-lived SSE connection (e.g.
an overly aggressive corporate proxy).

## 3. Implementation

### Backend — new SSE subscription endpoint

- `internal/components/auth/services/auth.go` — new `LoginNotificationStream(c *gin.Context, s
  *serverModels.Server)` handler. Registers a **buffered** channel (capacity 1) into
  `s.GC.LoginUsers[loginID]` under the existing `s.GC.Mutex`, matching the locking convention
  `BroadcastToLoginID` already uses. The buffer of 1 matters: `BroadcastToLoginID` sends
  (`conChan <- message`) while holding that same global mutex, and a login notice fires at most
  once per login attempt — an unbuffered channel would risk that send blocking (and, since it
  holds the global mutex, stalling every other broadcast map in `GlobalConfig` — `GeneralUsers`,
  `UsersOnline`, `AuthUsers` — while it waits) if the handler weren't in the middle of a `select`
  at that exact instant. A buffer of 1 makes the send always succeed immediately, no timing
  dependency on the reader.
  - Streams as `text/event-stream`, one heartbeat comment every 20s to keep intermediary proxies
    from closing an idle connection, and a 5-minute hard cap so an abandoned tab doesn't hold a
    server-side connection open forever.
  - Cleans up via `defer`, deleting its own map entry only if it's still the current registration
    for that `loginID` (guards against a reconnect's newer registration being deleted by an older
    handler's deferred cleanup).
  - Unauthenticated, exactly like the existing `GET /users/verify/:targetUser/:loginID` it
    complements — the user has no JWT yet at this point in the flow, and `EventSource` cannot send
    custom headers anyway. `loginID` (an unguessable, single-use identifier minted when the login
    attempt starts) is the same de facto capability boundary the polling endpoint already relies
    on.
- `internal/components/auth/controllers/main.go` — registers `GET /login/stream/:loginID`. Placed
  under a distinct path (not nested under `/users/verify/...`) to avoid any Gin route-tree
  ambiguity between a static segment and the existing `:targetUser` wildcard at the same depth.

### Frontend — Option 1 (focus-triggered refetch)

- `redux/store.ts` — calls `setupListeners(store.dispatch)` once, globally required for RTK Query
  to observe `visibilitychange`/`focus`/`online` at all.
- `qr-scan/page.tsx` — `useVerifyLoginQuery(..., { pollingInterval: 5000, refetchOnFocus: true })`.
  Scoped to this one call site, not set globally on `createApi`, so no other table/list in the
  dashboard starts refetching on every tab-focus as a side effect.

### Frontend — Option 3 (SSE fast path)

- `qr-scan/page.tsx` — once `loginProps.loginId` is available (right after the QR code is shown),
  opens `new EventSource(\`${BASE_URL}/login/stream/${loginProps.loginId}\`)`. On any message,
  calls `refetch()` from `useVerifyLoginQuery` to immediately re-check with the backend for the
  real token payload — the SSE message itself only carries `{username}` (`LoginCallback`'s
  payload), it's a "check now" signal, not the login result itself. Closes the connection on
  unmount, on a new QR code being generated, and once `loginDetails.data` is populated (login
  complete). Polling keeps running unchanged alongside it — this is additive, not a replacement.

## 4. Files touched

- `internal/components/auth/services/auth.go` — `LoginNotificationStream`.
- `internal/components/auth/controllers/main.go` — route registration.
- `frontend/src/redux/store.ts` — `setupListeners`.
- `frontend/src/app/(auth)/qr-scan/page.tsx` — `refetchOnFocus`, `EventSource` subscription.

## 5. Verification performed

- `go build ./...` — clean, whole repo.
- `go vet ./...` — clean, whole repo.
- `gofmt -l internal/components/auth/services/auth.go internal/components/auth/controllers/main.go`
  — clean.
- `make swagger` (`swag init` pinned `v1.16.2` + `scripts/fix-swagger-yaml.sh`) — regenerated
  `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`; diff is additive-only (97 lines across
  the three files), adding exactly the new `GET /login/stream/{loginID}` route and nothing else.
  The pre-existing "route declared multiple times" warnings `swag` prints are unrelated duplicate
  route registrations elsewhere in the codebase, not caused by this change.
- `npx tsc --noEmit` — this branch is based on `origin/dev` (not the earlier vault-signer feature
  lineage), and on `dev` the whole frontend already type-checks with **zero** errors — no
  pre-existing baseline to compare against here, unlike prior fixes in this series. `tsc` exits 0
  after these changes too, including the `SetSignerValueModal.tsx`, `PersonalSecretModal.tsx`, and
  `qr-scan/page.tsx` copy/error-handling changes below.
- `yarn lint` (`next lint`) — clean; the only warnings reported are pre-existing, in files this
  change never touches.
- `yarn build` (`next build`) — succeeds; the full route manifest builds, including the touched
  `qr-scan` and `vault-signer` routes.

### Additional frontend changes folded into this branch

While implementing Options 1 and 3, also addressed a related usability request on the same
signer-secret entry flow:

- `frontend/src/app/vault-signer/components/SetSignerValueModal.tsx` — the dialog's caption is now
  the managed secret's own friendly `label` directly (was `Set Value — ${label}`), and the field
  label reads **Signer** (was "Stellar signer key").
- `frontend/src/app/vault-signer/components/PersonalSecretModal.tsx` (create mode) — the dialog's
  caption is now the personal-env's `friendlyLabel` directly (was `Create ${friendlyLabel}`), and
  the field label reads **Signer** (was "Value") — matching `SetSignerValueModal.tsx`. Delete-mode
  caption is untouched (`Delete ${friendlyLabel}`) since it's a confirmation dialog, not a value
  entry dialog.
- `frontend/src/app/(auth)/qr-scan/page.tsx` — `handleLogin`'s failure toast now reads the
  backend's human-readable `message` field first (falling back to the `error` slug, then a generic
  string) instead of passing the raw error object into `showErrorToast`, which expects a string.
  The vault-signer secrets flow (`MySignerSecrets.tsx`'s `getErrorMessage`, shared by
  `PersonalSecrets.tsx`) already preferred `message` over `error` from an earlier fix in this
  series and needed no change here.

## 6. Verification not performed, and why

No live backend, Redis, or browser available in this sandbox to actually open an SSE connection,
background a tab, and confirm the push arrives — this was built and reasoned through by reading
the code and the SSE/EventSource/Gin-streaming mechanics directly, not exercised live.
**Recommended before relying on this in production:**
- A real QR login with the tab backgrounded immediately after the QR code appears — confirm
  login completes within a second or two of approval rather than waiting for the tab to regain
  focus.
- Confirm whatever reverse proxy sits in front of this backend doesn't buffer or time out
  long-lived `text/event-stream` responses (a common nginx gotcha — this handler sends
  `X-Accel-Buffering: no`, but the proxy's own idle-timeout configuration is outside this
  codebase's control and worth checking).
- Confirm the SSE connection is torn down correctly on tab close (no leaked server-side
  goroutines/connections building up over time in production).

## 7. Status

**Implemented, on branch `fix/login-tab-throttling`.**
