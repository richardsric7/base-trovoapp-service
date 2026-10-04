# Deployment Guide

This guide assumes you have never deployed a Go service before. It covers
building, running, and (as far as this repository's own files let us verify
accurately) how `app-backend` is deployed today.

All commands below are run from the `app-backend/` directory unless stated
otherwise.

## 1. Prerequisites

| Tool | Why you need it | Install |
| --- | --- | --- |
| **Go** (see `go.mod` for the exact version — currently `1.23`, toolchain `1.24.3`) | To build/run the service directly | [go.dev/doc/install](https://go.dev/doc/install) |
| **Docker** | To build/run the container image the way it's deployed in production | [docs.docker.com/get-docker](https://docs.docker.com/get-docker/) |
| **PostgreSQL** (optional for local dev — SQLite works too) | The production database | [postgresql.org/download](https://www.postgresql.org/download/) |

You do **not** need Postgres to get the app running locally: the
`.env.example` file ships with `DB_TYPE=sqlite`, which needs no separate
database server at all.

## 2. Building

### Plain Go build

```bash
go build ./...
```

This compiles every package and reports any compile errors. It does not
produce a single runnable binary by itself for `main.go` if you're inside a
subdirectory — to build the actual server binary:

```bash
go build -o trovo-wallet-api .
```

### Docker build

The repository's `Dockerfile` is a two-stage build:

```dockerfile
FROM golang:alpine3.20 AS builder
...
RUN go build -o main .
...
FROM alpine:latest
COPY --from=builder /dist/main /
COPY trovo-logo.png /trovo-logo.png
COPY ht2.png /ht2.png
ENTRYPOINT ["/main"]
```

Stage by stage:

1. **`builder` stage** (`golang:alpine3.20`) — installs `ca-certificates`,
   `build-base`, `runc`, and `curl`; copies `go.mod`/`go.sum` and runs
   `go mod download` first (so Docker's layer cache is reused across builds
   that don't change dependencies); then copies the rest of the source and
   runs `go build -o main .` to produce a single static-ish binary.
2. **Final stage** (`alpine:latest`) — a small runtime image. It installs
   only the runtime packages it needs (`ca-certificates`, `runc`, `curl`),
   copies the compiled `main` binary and the two image assets the app
   embeds by file path (`trovo-logo.png`, `ht2.png`) from the builder
   stage, and sets `ENTRYPOINT ["/main"]`.

To build the image yourself:

```bash
docker build -t trovo-wallet-api:local .
```

## 3. Running locally

### Directly with `go run`

```bash
go run main.go
```

At minimum, `main.go` refuses to start (see the `requiredEnvironmentVariables`
list near the top of `main()`) unless roughly 50 environment variables are
set — covering the blockchain RPC endpoint, mnemonics for several signer
roles, the RoachDB connection string, Mailgun, Redis, JWT secrets, and a
long list of fee-wallet addresses. The full, documented list — with example
values and where to get real ones — is in
**[CONFIGURATION.md](CONFIGURATION.md)**. The fastest way to get a complete,
valid set for local development is:

```bash
cp .env.example .env
# then fill in the blanks per CONFIGURATION.md
go run main.go
```

The server listens on `PORT` if set, otherwise `:8080`.

### Via Docker

```bash
docker run --rm -p 8080:8080 --env-file .env trovo-wallet-api:local
```

Same environment-variable requirements as above. If you're using SQLite in
`DB_TYPE`, remember the SQLite file path is relative to the container's
filesystem, not your host's — either bind-mount a directory in or switch to
Postgres for a containerized run.

## 4. Database migrations and `DB_AUTOMIGRATE`

This service uses **GORM `AutoMigrate`**, not hand-written SQL migration
files. On every boot (see `internal/db/main.go`'s `MigrateDB` and
`main.go`'s call to it), unless `DB_AUTOMIGRATE=0` is set, the app calls
`AutoMigrate` against roughly **107 GORM models** — one call per model —
which inspects and, if needed, alters the live schema to match each Go
struct.

This has one important consequence: **it is slow**, especially over a
network connection to the database rather than a local/same-datacenter one.
Running this AutoMigrate step from a machine far from the database (e.g.
a hosted build runner, over the public internet) has been observed to take
**over 20 minutes**, because
each of the ~107 models' AutoMigrate calls does multiple catalog round-trips
against Postgres — roughly **~800 round-trips in total** for the full set.
The practical implications for you:

- **Local development**: leave `DB_AUTOMIGRATE` unset or `1` (the
  `.env.example` default) — your DB is local, so the round-trip cost is
  negligible.
- **A real deployment**: run migrations from *inside* the same network as
  the database (e.g. on the server itself, or a same-region CI runner), not
  from a general-purpose CI runner reaching the DB over the public internet.
  This codebase supports a **`MIGRATE_ONLY=1`** mode (see the top of
  `main()` in `main.go`) specifically for this: with `MIGRATE_ONLY=1`, the
  process opens only the wallet's own Postgres database, runs
  `AutoMigrate`, and exits — without booting the HTTP server or requiring
  the ~50 other required environment variables the serving process needs.
  The serving process itself is then expected to run with
  `DB_AUTOMIGRATE=0`, so it boots straight into serving traffic without
  re-running AutoMigrate on every restart.

### Running two or more instances (and a dependency fix this needed)

If you scale this service to more than one instance — or if `DB_AUTOMIGRATE`
ends up enabled on a serving instance alongside a separate `MIGRATE_ONLY=1`
job — more than one process can end up calling `AutoMigrate` against the
same database at the same time. `MigrateDB` (`internal/db/main.go`) now
guards this with a small **distributed lock**: a `distributed_locks`
database table, claimed with a plain `UPDATE ... WHERE` statement that's
atomic on both SQLite and Postgres (deliberately *not* a Postgres advisory
lock or `SELECT ... FOR UPDATE`, since SQLite — this project's local/test
database — supports neither). A second instance booting at the same time
just waits (up to 10 minutes) for the first one's migration to finish,
instead of racing it. You don't need to configure anything for this; it's
automatic. See `internal/sharedconfig/distributed_lock.go` if you want the
details.

Alongside this, a real bug was found and fixed: `gorm.io/driver/postgres`
was pinned to a version (`v1.3.7`) badly out of sync with this project's
`gorm.io/gorm` version. The practical effect: `AutoMigrate` worked the
*first* time it ran against a fresh table, but **every subsequent run
against an already-existing table failed** with an `insufficient
arguments` error — meaning, in practice, every restart or redeploy after
the very first one would have fatally errored during migration against a
real Postgres database (this was reproduced and confirmed against real
Postgres 13, 14, and 16). This has been fixed by bumping the driver to
`v1.5.11` (verified against all three Postgres versions and SQLite; no
other dependency needed to change). If you're setting up a fresh
environment from this repo, you already have the fix — nothing to do.

This service's background sweep loops (sales activation, Stablerail
pollers, order-expiry checks, etc. — see `main.go` and
`internal/components/p2p/controllers/main.go`) are also now safe to run
on multiple instances: each tick is guarded by the same kind of
distributed lock, so only one instance actually does the work on a given
tick and the others skip it, rather than every instance doing the same
work redundantly. Again, nothing to configure — see
`internal/sharedconfig/singleton_lock.go`.

Further multi-instance fixes (all automatic):

- **Locks are renewed while their work runs.** A background job may take
  longer than its lock's timeout without another instance starting the
  same job; the timeout only decides how soon the others take over from
  an instance that died. Waits between runs happen outside the lock.
- **Platform signing keys are used by one instance at a time.** Every
  transaction from a platform Safe (P2P escrow settlements and refunds,
  fiat purchase delivery, the account-recovery guardian) holds a lock on
  that Safe from reading its nonce until the transaction is mined, and
  every transaction from a platform key holds a lock on the key from
  reading its nonce until it is broadcast. Without this, two requests at
  once (on one instance or several) signed the same nonce and one failed.
  Requests wait for the lock (up to 3 minutes) rather than fail.
- **Sale-start notifications** are flagged on the asset in the database
  instead of passed through an in-memory queue, so the instance that
  starts a sale and the one that sends its notifications can differ.
- **Partner callbacks** (service-link login/authorization/event callbacks
  and payment notifications) are recorded in `callback_deliveries` before
  they are sent and retried with backoff (up to 10 attempts) by one
  instance at a time, so a stopped instance loses none. Before this,
  payment-notification callbacks were never sent at all.
- **Graceful shutdown.** On SIGTERM/SIGINT (an autoscaler removing the
  instance) it stops starting background work, lets in-flight requests
  finish and waits for running jobs and platform-key transactions, for up
  to `SHUTDOWN_GRACE_PERIOD` (default 60s). Give the platform's stop
  timeout at least that long.
- The crypto-deposit minting loop and referral-link generator run on one
  instance at a time; channel-account funding at boot takes the funder
  key's lock, and new channel accounts are generated by at most one
  instance per hour.

## 5. How this service is deployed

There is **no CI/CD pipeline** in this repository (the inherited GitHub
Actions workflows were removed - see the root `ARCHITECTURE.md`). Build
and deploy by hand:

1. Regenerate the Swagger docs if handlers changed (`swag init` - see the
   header of `docs/docs.go` for the flags this project uses).
2. Build the image from this directory: `docker build -t trovo-wallet-api .`
3. Push it to your registry and redeploy it wherever it runs.

Run schema migrations as a separate step when deploying to a remote
database (see section 4: `MIGRATE_ONLY` / `DB_AUTOMIGRATE`).

## 6. Makefile targets

The `Makefile` in this directory has the checks to run before merging:

```bash
make build   # go build ./...
make vet     # go vet -stringintconv=false ./...  (see comment in Makefile for why)
make lint    # golangci-lint run (requires golangci-lint installed)
make test    # go test ./internal/... -race -vet=off -coverprofile=coverage.out
make ci      # tidy-check + build + vet + lint + test, in that order
```

`make test` intentionally excludes the root package's `main_test.go`, which
is an integration suite that expects a live server already running at
`localhost:8080`.
