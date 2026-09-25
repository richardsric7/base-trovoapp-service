# Login Verify Cache Fix — Plan

## 0. Symptom

The frontend login/QR flow polls `GET /users/verify/:targetUser/:loginID` every 5 seconds
(`useVerifyLoginQuery`, `pollingInterval: 5000` in `frontend/src/app/(auth)/qr-scan/page.tsx`)
while waiting for the user to approve the login in their Trovo Wallet app. Even after the user
approves and the backend has everything it needs to return the success payload (access token,
refresh token, user info), the endpoint kept returning as if nothing had happened for well over
30 seconds — sometimes up to a minute — before finally reflecting the successful login.

## 1. Root cause

`VerifyLoginID` (`internal/components/auth/services/auth.go`) wraps its entire response in a
Redis cache, keyed only on `(targetUser, loginID)` — a key that never changes across polls for
the same login attempt:

```go
cacheKey := fmt.Sprintf("[GET] /v1/users/verify/%v/%v", targetUser, loginID)
ok, status, response := s.GC.Cache.CachedHttpResponse(cacheKey)
if ok {
    c.JSON(status, response)   // returns immediately — never re-checks with Trovo
    return
}
```

Before the wallet approves, the external Trovo verification call fails, and that
**pending/failure response was cached for 60 seconds**:
```go
cacheDurationInSeconds := 60 // 1 minute
s.GC.Cache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)
```

So the first poll before approval caches "pending" for 60s, and every poll in that window — even
ones that land *after* the user has approved on their phone — hits the cache check at the top of
the handler and replays the same stale "pending" answer without ever calling Trovo again. Only
once that cache entry naturally expires does a poll actually re-check and see the real state.
"Over 30s" is the expected average wait, since approval timing is essentially random relative to
when that 60-second window started.

This caching is gated behind `ENABLE_CACHING=1` (`main.go:129-133`) — off by default, on wherever
the symptom was observed.

A commented-out cache-invalidation call sits right before the success path caches its own
response (`auth.go`, just above where `tokenResponse` is cached) — evidence someone had already
run into this and started a fix that never landed. Even if re-enabled, it would only clear the
cache *after* a poll happens to get through to a fresh check — it doesn't stop the up-to-60s of
stale replay beforehand, so it wasn't a sufficient fix on its own.

### Two secondary bugs found while reading this code

1. **Status-code mismatch.** The live success response is sent as `201 Created`
   (`serverResponse.JSON(c, http.StatusCreated, ...)`), but the *cached* copy of that same
   response was stored under `http.StatusOK` (200) — a cache hit within the 2-minute success
   cache window would have replayed the wrong status code.
2. **Body-shape mismatch.** The live path wraps the token payload in `serverResponse.JSON`'s
   envelope (`internal/server/response/json.go`'s `Data{Message, Data, Timestamp, Errors,
   Status}`), but the cached value was the raw, unwrapped `tokenResponse` struct. A cache hit
   would have sent the frontend a body without the `data` wrapper it expects — the frontend reads
   `loginDetails.data.accessToken`, which would have been `undefined` against the unwrapped
   shape.

## 2. Fix implemented (Option B + status code fix)

Per the fix options discussed (A: stop caching pending responses entirely, B: shorten the
pending-cache TTL to at/below the polling interval) — **Option B** was chosen, plus fixing the
status-code mismatch (which, on inspection, turned out to include the body-shape mismatch too,
fixed alongside for the same reason: a cached response must match its live counterpart).

- Added `pendingLoginVerifyCacheSeconds = 3` — a named constant, replacing the two active `60`
  (1-minute) pending/error cache-duration literals in `VerifyLoginID`. 3 seconds is at/below the
  frontend's 5-second polling interval, so no poll can ever replay an answer more than one polling
  cycle stale, while still deduping near-simultaneous requests (e.g. multiple browser tabs).
- The success-path cache write now stores `http.StatusCreated` (matching the live status) instead
  of `http.StatusOK`.
- The success-path cache write now stores the same `serverResponse.Data{Message, Data, Status,
  Timestamp}` envelope the live response actually sends, instead of the raw unwrapped
  `tokenResponse` struct — so a cache hit and a cache miss are byte-shape-identical from the
  frontend's point of view.

**Not done:** re-enabling the commented-out `InvalidateCachedHttpResponse` call, and Option A
(dropping pending-state caching entirely). With the TTL now bounded at 3 seconds, the residual
staleness window is small enough that neither was judged necessary on top of the TTL fix — worth
revisiting only if 3 seconds of caching still causes visible issues in practice.

## 3. Files touched

- `internal/components/auth/services/auth.go` — the three changes above, plus a `time` import
  (needed to stamp the cached envelope's `Timestamp` field the same way the live path does).

## 4. Verification performed

- `go build ./...` — clean, whole repo.
- `go vet ./...` — clean, whole repo.
- `gofmt -l internal/components/auth/services/auth.go` — clean.
- `go test ./internal/components/auth/...` — no test files exist for this package (`[no test
  files]` for controllers/models/services), so nothing to run; this fix does not add new test
  coverage.

## 5. Verification not performed, and why

No live backend, Redis instance, or reachable Trovo service-link in this sandbox, so the actual
caching behavior (a poll returning a fresh "pending" check every ≤3s instead of every 60s, and a
cache-hit response matching a cache-miss response byte-for-byte) was traced by reading the code
rather than exercised live. **Recommended before merging:** a manual smoke test with
`ENABLE_CACHING=1` and a real Redis instance — start a QR login, wait past the old 60s pending-
cache window without approving (confirm the poll still eventually says "pending," not stuck on an
error), then approve and confirm the frontend reflects success within a few seconds rather than
up to a minute.

## 6. Status

**Implemented, on branch `fix/login-verify-cache-staleness`.**
