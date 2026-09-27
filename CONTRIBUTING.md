# Contributing to this monorepo

This file exists mainly for one rule: **documentation is not optional, and it
is not a follow-up task.** It ships in the same commit/PR as the change that
made it necessary. This applies to a human contributor and to Claude (or any
other AI assistant) working in this repo identically — there is no separate,
looser bar for automated changes.

Read [ARCHITECTURE.md](ARCHITECTURE.md) first if you haven't already — it
explains how the seven projects in this monorepo relate to each other and the
one design decision (direct shared-DB write vs. calling app-backend's API)
that most new admin features need to make correctly.

## Every project keeps the same four documents, at its own root

| File | What it covers |
|---|---|
| `README.md` | What the project is, its role in the monorepo, tech stack, directory layout, how to run it locally. The front door — link out to the other three rather than duplicating them. |
| `DEPLOYMENT.md` | Every step from a clean checkout to a running instance: prerequisites (with install links), build, run locally, run via Docker, and — only where verifiable from the project's own files — how it's actually deployed. Never invent infrastructure you can't confirm exists; say what you can verify and flag what you can't. |
| `CONFIGURATION.md` | Every environment variable / build-time config value the project reads. Each entry needs a name, a realistic (never real) example value, a plain-English explanation of its effect, and concrete instructions for obtaining a real value — a shell command for a generated secret, the actual third-party provider's dashboard for an API key, a pointer to a sibling project's own `DEPLOYMENT.md` for a value that should point at another service in this monorepo. |
| `INTEGRATION.md` | What calls this project and what it calls, with the actual auth scheme, base-URL config, and any shared-database relationship — described from reading the real client/auth code, not assumed from convention. |

A REST API project (currently `app-backend`, `tm-api`, and — once its own
Swagger setup is added — `payment-history-engine`) additionally maintains
Swagger/OpenAPI documentation generated from `swaggo/swag` annotations on
every handler. A frontend or library project does not need Swagger — it has
no HTTP API of its own to document that way.

## When you add a new project to this monorepo

1. Give it the same four documents from day one — don't defer them to "later."
2. Add a row for it to the table in [ARCHITECTURE.md](ARCHITECTURE.md),
   including what it talks to and what talks to it.
3. If it exposes a REST API, wire up `swaggo/swag` from the start (see
   "Adding or changing an endpoint" below) rather than retrofitting it once
   the API has grown — retrofitting is exactly the situation this
   documentation initiative had to dig out of for `app-backend`, `tm-api`,
   and `payment-history-engine`.
4. Add it to `.github/workflows/pr-checks.yml`'s path-filtered build/lint/test
   matrix if CI wiring for this repo's actual directory layout gets fixed
   (see the note in ARCHITECTURE.md — as of this writing those workflows
   reference a different directory layout and don't run against this repo's
   projects at all).

## When you modify an existing project

Ask, for every change:

- **Did I add, remove, or change an environment variable / config value?**
  Update that project's `CONFIGURATION.md` in the same commit. Don't leave a
  variable undocumented because "it's obvious" — the whole point of this doc
  is that a first-time deployer has never seen this codebase before.
- **Did I add, remove, or change an HTTP endpoint** (on `app-backend`,
  `tm-api`, or `payment-history-engine`)? Add/update the `swaggo/swag`
  annotation on the handler (see below) and regenerate the docs.
- **Did I change how one project calls another**, or add a new
  cross-project dependency (a new caller of an existing API, a new shared
  table, a new library consumer)? Update `INTEGRATION.md` on both the
  calling and the called project, and update the flow diagram in
  `ARCHITECTURE.md` if the shape of the relationship changed.
- **Did I change how the project is built, run, or deployed** (a new
  Dockerfile stage, a new build prerequisite, a new required tool version)?
  Update `DEPLOYMENT.md`.

If none of these apply — a pure refactor, a bug fix with no behavior change a
deployer or integrator would need to know about — no doc update is needed.
Don't pad a PR with busywork doc edits that don't carry real information.

## Adding or changing an endpoint (app-backend / tm-api / payment-history-engine)

Every handler function backing a registered route carries a `swaggo/swag`
doc comment block directly above it. Minimum bar for every endpoint:

```go
// @Summary Suspend a user account
// @Description Suspends a normal user, logging the mandatory reason. Blocks the user's wallet(s) from any transaction until lifted.
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param data body models.SuspendOrLiftUserPayload true "Email of the user to suspend, and the mandatory reason"
// @Success 200 {object} response.Data
// @Failure 400,401,404,500 {object} object
// @Router /admin/users/suspend [patch]
func SuspendUser(s *serverModels.Server) gin.HandlerFunc {
```

- `@Summary` is one plain-English sentence a novice could understand — not a
  restatement of the function name.
- `@Tags` groups the endpoint with its feature area in the generated UI.
- `@Param` documents every path, query, and body parameter — not just "there
  is a body," but what's in it.
- `@Success`/`@Failure` list the actual HTTP codes the handler returns, not a
  generic guess.
- `@Router` must exactly match the path and method it's registered under.

After adding/editing annotations, regenerate the docs from that project's own
root directory:

```bash
go install github.com/swaggo/swag/cmd/swag@latest   # once, if not already installed
swag init                                             # from app-backend/, tm-api/, or payment-history-engine/
```

Check that project's own `docs/docs.go` header (or its `Makefile`, if it has a
`swag`/`docs` target) for any extra flags it was originally generated with
(e.g. `--parseDependency --parseInternal`) and match them — inconsistent flags
between runs can make `swag init` silently drop coverage.

Confirm `go build ./...` and `go vet ./...` still pass, and confirm the
Swagger UI is reachable (each project's `README.md` states the exact path,
typically `/swagger/index.html`) before opening the PR.

## Git workflow

Nothing here changes the repo's normal git conventions — small, reviewable
commits; a PR per logically-independent change rather than one giant "docs +
features" PR; don't force-push over someone else's work. See the individual
project READMEs for anything project-specific (test commands, lint commands)
that should pass before you open a PR.
