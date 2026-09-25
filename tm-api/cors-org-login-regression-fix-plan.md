# CORS Org-Login Regression — Hotfix

## Symptom

After the general bug sweep's CORS fix (PR #45, finding 1.1) merged to `dev`, organisation
members got a CORS error at login.

## Root cause

The fix replaced `Access-Control-Allow-Origin: *` with a strict allow-list: the backend only
echoed back the request's `Origin` when it exactly matched `CORS_ALLOWED_ORIGINS`, falling back to
the single `TROVO_MANAGER_FRONTEND_URL` value when that env var was unset. Any origin outside that
one configured value got **no** `Access-Control-Allow-Origin` header at all, and the browser
blocked the response outright — not a "credentials rejected" failure, a full block. If the
organisation portal is reached from any origin not equal to that one fallback value, this is what
happens.

Neither `axiosBasedQuery.ts` nor `orgBasedQuery.ts` actually sets `withCredentials: true`, which
means the original CORS-spec violation (`*` + `Allow-Credentials: true` being an invalid
combination) wasn't itself causing failures in production — the fix over-corrected by also adding
origin *restriction*, which needs every real deployed frontend origin known and configured ahead
of time, and evidently wasn't.

## Fix

`internal/middleware/cors.go` — still resolves the actual spec defect (never emit the literal
`"*"` alongside `Access-Control-Allow-Credentials: true`; always echo a concrete origin instead),
but stops restricting to an allow-list by default:

- `CORS_ALLOWED_ORIGINS` unset (the default) → every origin is echoed back, exactly matching the
  pre-fix "any origin" behavior, just never as the literal `"*"` string.
- `CORS_ALLOWED_ORIGINS` set → only origins in that list are echoed (opt-in restriction, for when
  every legitimate frontend origin is actually known and worth enumerating).

Removed the `TROVO_MANAGER_FRONTEND_URL` fallback entirely — a single-value fallback is exactly
what caused this regression, since it can never cover every real deployed frontend origin (the
organisation portal being the immediate example). `.env-sample` updated to describe the new
opt-in semantics and commented out by default.

## Verification

- `go build ./...`, `go vet ./...`, `gofmt -l internal/middleware/cors.go` — clean.
- Not exercised against a live org-portal login in this environment (no backend/browser
  available) — worth a manual pass confirming org member login succeeds before/after this
  deploys.

## Status

Implemented on branch `fix/cors-org-login-regression`, off `dev`.
