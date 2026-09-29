# Deployment Guide

This is a from-scratch guide to building, running and deploying `tm-api`,
written for someone who has never touched this codebase before. Everything
below is verified against this repo's own `Dockerfile` and `Makefile` —
nothing is assumed from memory. See
[CONFIGURATION.md](./CONFIGURATION.md) for what every environment variable
does, and [README.md](./README.md) for the project overview and the
two-database architecture.

## 1. Prerequisites

| Tool | Why you need it | Install |
|---|---|---|
| **Go 1.23+** (toolchain 1.24.3, per `go.mod`) | Build/run the API | https://go.dev/doc/install |
| **Docker** | Build/run the production container image | https://docs.docker.com/get-docker/ |
| **Access to a Postgres database** for `ADMIN_CONNECTION_STRING` | tm-api's own schema (`AdminDB`) — see README's two-database section | Any reachable Postgres 13+ instance. For local dev, the easiest path is `docker run -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:16` |
| **Access to app-backend's shared Postgres database**, if you need working wallet/P2P features locally | `TrovoWalletDB`/`P2P` — as of this snapshot this is *the same connection string* as `ADMIN_CONNECTION_STRING` (see README), so you need a database that already has app-backend's schema (`users`, curated assets, etc) applied, not just an empty one. Ask whoever owns `app-backend` locally for a dump/seed, or run `app-backend`'s own migrations against the same database first. | — |
| **`swag` CLI** (only if you're changing Swagger annotations) | Regenerates `docs/` | `go install github.com/swaggo/swag/cmd/swag@latest` |
| **`golangci-lint`** (only if you're running `make lint`/`make ci` locally) | Lints the code with this project's `.golangci.yml` | https://golangci-lint.run/usage/install/ |

## 2. Building

### 2a. Plain `go build`

```bash
go mod download
go build -o main .
./main
```

### 2b. Via `make`

```bash
make build   # go build ./...
make vet     # go vet ./...
make lint    # golangci-lint run -v --fix
make test    # go test ./... -race -coverprofile=coverage.out -covermode=atomic
make ci      # tidy-check + build + vet + lint + test — the same checks CI runs
```

### 2c. Walking through the `Dockerfile`

The `Dockerfile` at the repo root is a two-stage build:

```dockerfile
FROM golang:alpine3.20 AS builder
```
**Stage 1 (`builder`)** — compiles the binary:

1. Sets `CGO_ENABLED=0`, `GOOS=linux`, `GOARCH=amd64` for a static, portable binary.
2. Installs `ca-certificates build-base runc curl bash` (bash is needed for the swagger-fix script below; the others support the Go toolchain and outbound TLS calls the app itself makes at runtime, like `app-backend`/Discord/Mailgun/Vault).
3. `COPY go.mod go.sum` then `go mod download` — dependencies are cached in their own Docker layer, so a code-only change doesn't re-download the module graph.
4. `COPY . .` — copies the full build context in.
5. `go install github.com/swaggo/swag/cmd/swag@latest` — installs the swag CLI *inside the image*, so Swagger docs are always regenerated fresh at build time rather than trusting whatever's committed in `docs/`.
6. `bash scripts/docker-build-swagger.sh` — locates the just-installed `swag` binary, deletes `docs/`, runs `swag init`, then runs `scripts/fix-swagger-yaml.sh` (see the note below on why that fix-up exists).
7. `go build -ldflags "-X admin-panel-dashboard/internal/components/health/services.buildVersion=${APP_VERSION}" -o main .` — builds the binary, stamping the `APP_VERSION` build arg into the health package so `/health` reports exactly which build is running. CI passes the git SHA as `APP_VERSION`; for a local `docker build` it defaults to `unknown`.
8. Copies the resulting `main` binary into `/dist`.

```dockerfile
FROM alpine:latest
```
**Stage 2 (final image)** — a minimal runtime image: just `ca-certificates` and `curl` (for any outbound HTTPS calls and for a container healthcheck script, if one is added) plus the compiled `main` binary, run directly as the container's `ENTRYPOINT`. No Go toolchain, no source, no swag — the final image is small.

**Why the Swagger YAML fix-up exists:** the installed version of `swag` has a
known quirk where `swag init`'s generated `docs/swagger.yaml` sometimes puts
the `swagger:` and `info:` keys after other content instead of at the top,
which some YAML/OpenAPI parsers reject. `scripts/fix-swagger-yaml.sh` checks
`docs/swagger.yaml`, and if `swagger:` doesn't appear before `info:`,
rewrites the file with those two blocks moved to the front. `make swagger`
and `make swagger-clean` both run this automatically after `swag init` — see
the `Makefile`. **You should never need to run `swag init` directly without
also running this script (or `make swagger`/`make swagger-clean`).**

### 2d. Building the Docker image directly

```bash
docker build -t tm-api:local .
# with a real version stamp:
docker build --build-arg APP_VERSION=$(git rev-parse --short HEAD) -t tm-api:local .
```

## 3. Running locally (no Docker)

```bash
cp .env-sample .env
# edit .env — at minimum set ADMIN_CONNECTION_STRING to a real Postgres DSN
# (see CONFIGURATION.md for every variable)

go mod download
go run main.go
# or: make run
```

On startup, `main.go`:

1. Loads `.env` (skipped in `GIN_MODE=release`, where real env vars are expected instead).
2. Opens `AdminDB` and runs GORM `AutoMigrate` against it (creates/updates tm-api's own tables — see `internal/db/main.go`).
3. Seeds the default super-admins (`DEFAULT_SUPER_ADMINS`, or `obi,toluwase` if unset — see CONFIGURATION.md), the role/permission config, and suspension reasons.
4. Wires up every route (see the list of components in README.md).
5. Starts listening on `PORT` (default `8082`).

Confirm it's up:

```bash
curl http://localhost:8082/health
curl http://localhost:8082/ping
open http://localhost:8082/swagger/index.html   # Swagger UI
```

### Running two or more instances

`migrateAdminSchema` (`internal/db/main.go`) already guards `AutoMigrate`
with a Postgres advisory lock, so two instances booting at the same time
against the same database serialize their migration correctly instead of
racing each other — nothing extra to configure there.

Separately, a real bug was found and fixed while working on this: the
pinned `gorm.io/driver/postgres` version (`v1.5.2`) was out of sync with
this project's `gorm.io/gorm` version in a way that broke `AutoMigrate`
against any table that already existed — meaning every restart after the
very first boot would have fatally failed migration against a real
Postgres database (reproduced and confirmed against real Postgres 13, 14,
and 16). This is fixed by bumping the driver to `v1.5.11`, verified
end-to-end via `AdminDB()` run three times in a row against the same
database (simulating first boot plus two redeploys) and the existing
SQLite migration test suite. If you're working from this repo, you
already have the fix.

## 4. Running via Docker

```bash
docker build -t tm-api:local .

docker run --rm -p 8082:8082 \
  --env-file .env \
  -e PORT=8082 \
  tm-api:local
```

Notes:

- `.env` must contain a Postgres DSN reachable **from inside the container**
  — `localhost` in `.env` will not resolve to your host machine's Postgres
  from inside Docker. Use `host.docker.internal` (Docker Desktop) or run
  Postgres as a linked container / on a shared Docker network instead.
- The container has no `ENV PORT` default baked in beyond what `main.go`
  falls back to (`8082`), so make sure the `-p` mapping matches whatever
  `PORT` you set.

## 5. CI/CD

There is **no CI/CD pipeline** in this repository (the inherited GitHub
Actions workflows were removed - see the root `ARCHITECTURE.md`). Before
merging, run `make ci` (tidy check, build, vet, lint, test). To deploy,
regenerate the Swagger docs (`make swagger`), build the image from this
directory (`docker build -t tm-api .`), push it to your registry and
redeploy it.

## 6. Swagger docs in deployment

The Docker build regenerates `docs/`
(`docs.go`, `swagger.json`, `swagger.yaml`) from source at build time via
`swag init` + the YAML fix-up script — it does not simply trust whatever is
committed. If you change a route or a `@Summary`/`@Router`/... annotation,
regenerate locally before committing so the committed copy in `docs/` stays
in sync with what reviewers see on GitHub:

```bash
make swagger          # swag init + fix-swagger-yaml.sh
# or, to force a completely clean regeneration:
make swagger-clean     # rm -rf docs/ && swag init + fix-swagger-yaml.sh

head -5 docs/swagger.yaml   # should start with "swagger:"
```
