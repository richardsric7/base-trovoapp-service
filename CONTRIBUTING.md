# Contributing to this monorepo

This file exists mainly for one rule: **documentation is not optional, and it
is not a follow-up task.** It ships in the same commit/PR as the change that
made it necessary. This applies to a human contributor and to Claude (or any
other AI assistant) working in this repo identically — there is no separate,
looser bar for automated changes.

Read [ARCHITECTURE.md](ARCHITECTURE.md) first if you haven't already — it
explains how the projects in this monorepo relate to each other and the
one design decision (direct shared-DB write vs. calling app-backend's API)
that most new admin features need to make correctly.

## Every project keeps the same four documents, at its own root

| File | What it covers |
|---|---|
| `README.md` | Starts with **"What this project does"**: a few plain sentences a non-developer could follow (what it is for, who or what uses it, what would stop working without it). Then its place in the monorepo, tech stack, directory layout and how to run it locally. Links to the other three rather than repeating them. |
| `DEPLOYMENT.md` | A step-by-step guide a first-time deployer can follow without help (format below). |
| `CONFIGURATION.md` | **Every** parameter the project reads (environment variables, build-time values, Vault keys, config files), each in the format below. |
| `INTEGRATION.md` | Every connection to another project or outside service, each in the format below. |

The reader to write for is a novice: someone deploying this for the first
time who has never seen the code. Do not assume they know what a term means,
where a value comes from, or which other service it must match. Never put a
real secret in an example.

### Parameter format (CONFIGURATION.md, and deploy-only parameters in DEPLOYMENT.md)

Every parameter gets its own heading and all five fields:

```markdown
### `DB_CONNECTION_STRING`

- **What it is:** The address and login of the Postgres database this service stores its data in.
- **Why it's needed:** Without it the service has nowhere to keep users or payments and refuses to start.
- **Required:** Yes.  *(or: No, default `8080`, meaning ...)*
- **Example:** `postgres://trovo:change-me@db.internal:5432/trovo?sslmode=require`
- **How to get it:** 1. Create a database (...). 2. Create a user with a password (...). 3. Put them together as shown in the example.
```

- **What it is** says what the value *is*, in plain words.
- **Why it's needed** says what it is used for and what happens if it is
  missing or wrong.
- **Required** says Yes or No; for No, give the default and what it means.
- **Example** is realistic and never a real secret.
- **How to get it** gives concrete steps: a command to generate a secret, the
  provider dashboard page for an API key, or, for a value that must match
  another project in this repo, exactly which parameter of which project it
  must equal (with a link).

Group parameters under headings by area (Server, Database, Blockchain, ...)
and start the file with how to set them (`.env` file, Docker, hosting
platform) and a short table of the ones needed just to start.

### DEPLOYMENT.md format

1. **What you are deploying:** one paragraph in plain words, and what must
   already be running first (database, other services).
2. **Before you start:** every account and tool needed, with the install
   link or command and the command that checks it is installed.
3. **Steps:** numbered, one action each, with the exact command to run and
   what you should see when it worked.
4. **Parameters:** the runtime parameters needed to start, listed with links
   to their `CONFIGURATION.md` entries, and any deploy-only parameters
   (build arguments, hosting settings) written out in the parameter format
   above.
5. **Check it works:** how to tell the deployment is healthy.
6. **Updating and rolling back.**
7. **Troubleshooting:** common symptoms, their cause and the fix.

Only describe infrastructure you can confirm from the project's own files;
say clearly what you could not confirm.

### INTEGRATION.md format

One section per connection (another project in this repo, or an outside
service such as a bank API, Vault or a blockchain node):

- **What it is and why:** what the other side is and what this project uses
  it for.
- **Direction:** who calls whom (or which database or queue is shared).
- **How they connect:** URL or address, protocol, and how requests are
  authenticated, described from the real code.
- **Settings on this side** and **settings on the other side:** the
  parameters that must be set (and must match), each with an example value
  and how to get it, or a link to its `CONFIGURATION.md` entry.
- **How to check it works**, and **what happens when it is down**.

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
4. Add it to the project list and deploy order in the root
   [README.md](README.md).
5. There is no CI pipeline yet (see ARCHITECTURE.md); document the
   project's build/lint/test commands in its `DEPLOYMENT.md` so they can be
   wired into one later.

## When you modify an existing project

Ask, for every change:

- **Did I add, remove, or change an environment variable / config value?**
  Update that project's `CONFIGURATION.md` in the same commit, with all five
  fields of the parameter format. Don't leave a
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
