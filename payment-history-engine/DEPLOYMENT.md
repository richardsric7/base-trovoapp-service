# Deployment

This document covers building and running `payment-history-engine` on its
own. **There is no CI/CD wiring for this service found anywhere in this
monorepo** — see [No CI/CD found](#no-cicd-found-in-this-repo) below. Treat
everything here as manual/local build-and-run steps you can verify yourself
from the Dockerfile, not a description of an existing pipeline.

## 1. Prerequisites

| Tool | Why you need it | Install |
|---|---|---|
| Go 1.22+ (repo pins toolchain `go1.24.7`) | Build/run the service | https://go.dev/doc/install |
| Docker | Build/run the container image | https://docs.docker.com/get-docker/ |
| A Postgres-compatible primary "wallet" DB | Stores `user_wallets` (read) and `payment_history` (written) — see [INTEGRATION.md](./INTEGRATION.md) | https://www.postgresql.org/download/ (or use `DB_TYPE=sqlite` locally, no install needed) |
| A CockroachDB instance ("RoachDB") | Stores this engine's own tracking/cursor state | https://www.cockroachlabs.com/docs/stable/install-cockroachdb (or use `ROACH_DB_TYPE=sqlite` locally, no install needed) |
| (Optional) Redis | Caching layer, only if `ENABLE_CACHING=1` | https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/ |
| (Optional) A Base RPC endpoint | Required in practice — `BASE_RPC_URL` is a hard-required env var at boot | Base Sepolia (testnet, free): `https://sepolia.base.org`. Base mainnet: run your own node or use a provider (Alchemy, Infura, QuickNode, etc.) |

### Standing up a local CockroachDB, if you want the real thing instead of SQLite

```bash
# macOS
brew install cockroachdb/tap/cockroach
# Linux — see https://www.cockroachlabs.com/docs/stable/install-cockroachdb-linux
cockroach start-single-node --insecure --listen-addr=localhost:26257 --http-addr=localhost:8081 --background
cockroach sql --insecure -e "CREATE DATABASE payment_history;"
```

Then:

```bash
CDB_CONNECTION_STRING="postgresql://root@localhost:26257/payment_history?sslmode=disable"
```

## 2. Building

### Plain `go build`

```bash
cd payment-history-engine
go mod download
go build -o payment-history-engine .
./payment-history-engine
```

### Via the Dockerfile (stage by stage)

The `Dockerfile` at the project root is a two-stage build:

**Stage 1 — `builder` (`golang:alpine`)**

```dockerfile
FROM golang:alpine AS builder
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN apk add --no-cache ca-certificates build-base runc curl
WORKDIR /build
COPY go.mod . 
COPY go.sum .
RUN go mod download
COPY . .
ARG APP_VERSION=unknown
RUN go build -ldflags "-X trovo-wallet-payment-history-engine/internal/components/health.buildVersion=${APP_VERSION}" -o main .
WORKDIR /dist
RUN cp /build/main .
```

1. Starts from `golang:alpine`, sets `CGO_ENABLED=0` (a fully static binary —
   needed since the final stage has no glibc/musl dev headers) and targets
   `linux/amd64`.
2. Installs `ca-certificates` (for outbound TLS to the Base RPC endpoint and
   Postgres/CockroachDB over TLS), `build-base` (cgo/gcc toolchain, kept even
   with `CGO_ENABLED=0` since some transitive deps' `go generate`/build
   tooling can want it), `runc` and `curl`.
3. Downloads modules (`go mod download`) with only `go.mod`/`go.sum` copied
   first, so this layer is cached across builds that don't touch dependencies.
4. Copies the rest of the source and builds, stamping `APP_VERSION` (pass
   this as `--build-arg APP_VERSION=$(git rev-parse --short HEAD)` — CI would
   normally pass the git SHA here, but see the CI note below, nothing in this
   repo does that today) into
   `internal/components/health.buildVersion` via `-ldflags -X`, which is what
   `GET /health` and `GET /ready` report as `version`.

**Stage 2 — final image (`alpine:latest`)**

```dockerfile
FROM alpine:latest
RUN apk add --no-cache ca-certificates runc curl
COPY --from=builder /dist/main /
ENTRYPOINT ["/main"]
```

Copies only the compiled binary into a fresh, minimal Alpine image (no Go
toolchain, no source) and runs it as PID 1's entrypoint.

Build and run it:

```bash
docker build -t payment-history-engine:local --build-arg APP_VERSION=$(git rev-parse --short HEAD 2>/dev/null || echo dev) .
docker run --rm -p 8080:8080 --env-file .env payment-history-engine:local
```

`HEALTH_PORT` defaults to `8080` inside the container; the `-p 8080:8080`
maps that out. If you point `--env-file` at a `.env` with `DB_TYPE=sqlite`
and a relative `DB_CONNECTION_STRING`, remember the SQLite file is written
**inside the container's filesystem** and disappears when the container is
removed — mount a volume (`-v $(pwd)/data:/data` and point
`DB_CONNECTION_STRING=/data/wallet.db`) if you need it to persist.

## 3. Running

### Locally (no Docker)

```bash
cp .env.example .env
# edit .env — see CONFIGURATION.md
go run .
```

### Via Docker Compose (minimal example, not present in this repo — write your own)

```yaml
services:
  payment-history-engine:
    build: .
    ports:
      - "8080:8080"
    env_file: .env
    restart: unless-stopped
```

```bash
docker compose up --build
```

### Verifying it's up

```bash
curl -s localhost:8080/health | jq
curl -s localhost:8080/ready  | jq
open http://localhost:8080/swagger/index.html   # or just visit it in a browser
```

### A dependency fix worth knowing about

A real bug was found and fixed while auditing this monorepo for
multi-instance readiness: the pinned `gorm.io/driver/postgres` version
(`v1.3.7`) was out of sync with this project's `gorm.io/gorm` version in a
way that broke `AutoMigrate` (used by `internal/db/main.go`'s
`OpenRoachDB`) against any table that already existed — meaning every
restart after the very first boot would have fatally failed migration
against a real Postgres/CockroachDB database (reproduced and confirmed
against real Postgres 16). This is fixed by bumping the driver to
`v1.5.11`, verified against both real Postgres and the `ROACH_DB_TYPE=sqlite`
local escape hatch. If you're working from this repo, you already have
the fix — nothing to do.

## 4. No CI/CD found in this repo

This monorepo's two GitHub Actions workflows were checked directly:

- **`.github/workflows/deploy.yml`** — triggers only on pushes touching
  `backend/**`, `web/**`, `trovotech-io/**`, `trovo-app-website/**`, or the
  workflow file itself. It builds and pushes `app-backend` and the web
  frontend to DigitalOcean Container Registry. It explicitly documents (in a
  comment on its migration step) that it does **not** configure
  `CDB_CONNECTION_STRING`, because "RoachDB belongs to the payment-history
  service" and the migrator doesn't touch it.
- **`.github/workflows/pr-checks.yml`** — its `paths-filter` only defines
  `backend`, `web` and `mobile` filters. There is no `payment-history-engine`
  filter, job, or path reference anywhere in either file.

**Conclusion**: nothing in `.github/workflows/` builds, tests, lints, or
deploys this service. If it is deployed anywhere today, that pipeline lives
outside this repository (e.g. triggered manually, or from a separate
CI system/registry not checked into `.github/workflows/` here). Treat the
build/run steps above as the verified, manual path until CI wiring for this
service is added.
