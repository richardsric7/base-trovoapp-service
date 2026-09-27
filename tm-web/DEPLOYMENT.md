# Deployment Guide

This guide assumes no prior familiarity with this codebase or with Next.js.
It covers everything needed to install, build, and run `tm-web` — for local
development and for a production-like deployment.

## 1. Prerequisites

| Requirement | Version | Notes |
|---|---|---|
| Node.js | **20.x** | Inferred from the Dockerfile (`node:20-alpine`); `package.json` has no `engines` field pinning this for local installs, so match the Docker image's major version to avoid surprises. Install from [nodejs.org](https://nodejs.org/en/download) or with a version manager such as [nvm](https://github.com/nvm-sh/nvm#installing-and-updating) (`nvm install 20`). |
| Yarn | 1.x (Classic) | This project uses `yarn.lock`, not `package-lock.json` or `pnpm-lock.yaml` — confirmed by the presence of `yarn.lock` at the project root and no other lockfile. `package.json` pins `packageManager: yarn@1.22.22`. Install via [yarnpkg.com/getting-started/install](https://classic.yarnpkg.com/en/docs/install) or `corepack enable` (bundled with Node 16.9+). |
| A running `tm-api` | any | `tm-web` has no data of its own — see `tm-api/DEPLOYMENT.md` in this monorepo for how to stand one up (locally or otherwise). |
| Docker (optional) | any recent version | Only needed if you want to build/run the production container image. |

## 2. Installing dependencies

```bash
cd tm-web
yarn install --frozen-lockfile
```

`--frozen-lockfile` fails the install if `yarn.lock` is out of sync with
`package.json`, which is what CI/Docker builds should use; for day-to-day
local development plain `yarn install` is fine too.

## 3. Environment variables — the Next.js nuance

**Read this before setting anything up.** Next.js treats environment
variables very differently from a plain single-page app:

- Any variable whose name starts with **`NEXT_PUBLIC_`** is inlined
  (baked in as a literal string) into the JavaScript bundle **at build
  time** (`yarn build` / `next build`). Once built, that value is fixed —
  it ships to every browser that loads the app, is visible to anyone who
  opens devtools, and **cannot be changed by setting an environment
  variable on the running container**. To change a `NEXT_PUBLIC_*` value
  you must rebuild the app with the new value set.
- Any variable **without** the `NEXT_PUBLIC_` prefix stays server-side
  only — it's readable by code that runs on the Next.js server (e.g. the
  Sentry build/config plumbing) but is never sent to the browser. Since
  this app has no server-side data-fetching or API routes of its own,
  server-only variables here are used only for build tooling (source-map
  upload), not for anything the running app depends on at request time.

In short: for `tm-web`, **all the variables that matter to the running app
are `NEXT_PUBLIC_*` and must be correct before you run `yarn build`,**
not before `yarn start`/`next start`. See [CONFIGURATION.md](./CONFIGURATION.md)
for the full list and example values.

## 4. Building for production

```bash
yarn build
```

This runs `next build` (see the `"build"` script in `package.json`). Set the
`NEXT_PUBLIC_*` variables (at minimum `NEXT_PUBLIC_API_BASE_URL`) in your
shell or `.env.production.local` *before* running this, per the note above.

There is also a combined CI script:

```bash
yarn ci   # runs `yarn lint && yarn build`
```

`next.config.mjs` does **not** set `output: 'standalone'` — this project
produces a regular `.next/` build output that still needs `node_modules`
present at runtime (a "standalone" build bundles a minimal self-contained
server and doesn't need `node_modules` copied in; that's not the case
here — see the Dockerfile below, which does copy `node_modules`).

## 5. Running in production

```bash
yarn start   # runs `next start`, serving on PORT (default 3000)
```

This starts the regular Next.js production server, reading the already-built
`.next/` directory produced by `yarn build`. It needs `node_modules`,
`.next/`, `public/`, `package.json`, and `next.config.mjs` present (exactly
what the Dockerfile below copies into the runtime stage).

### 5.1 The Dockerfile

The repo root has a **multi-stage** `Dockerfile` with three stages:

1. **`deps`** (`node:20-alpine`) — copies only `package.json` and
   `yarn.lock`, then runs `yarn install --frozen-lockfile`. Isolating this
   step lets Docker cache the (slow) dependency install as its own layer,
   only invalidated when the lockfile/manifest actually change.

2. **`builder`** (`node:20-alpine`) — copies in `node_modules` from `deps`,
   copies the full source tree, and runs `yarn build`. This is where all the
   `NEXT_PUBLIC_*` build args get consumed (see below) — they're declared as
   Docker `ARG`s and re-exported as `ENV` so `next build` can read them via
   `process.env`. It also sets `NODE_OPTIONS=--max-old-space-size=4096`
   because the default Node heap isn't enough once the Sentry webpack plugin
   instruments the build (documented in the Dockerfile's own comments — the
   build previously died with an out-of-memory error at the default ~2 GB
   heap).

3. **`runner`** (`node:20-alpine`) — the actual image that ships. Creates an
   unprivileged `nextjs` user/group, copies `public/`, `.next/` (with the
   webpack build cache deleted — it's only useful for speeding up a
   *rebuild*, and copying its hundreds of MB into the shipped image bloats
   every layer/push/pull for no runtime benefit), `node_modules`,
   `package.json`, and `next.config.mjs` from the `builder` stage, then runs
   as the `nextjs` user on port 3000 via `CMD ["yarn", "start"]`.

Because there's no `output: 'standalone'`, the runtime stage has to carry
the full `node_modules` (not just a pruned server bundle), which is why that
stage explicitly copies it.

**Build args the Dockerfile accepts** (pass with `--build-arg`, all
optional except the API URL, which the app needs to be useful):

```bash
docker build \
  --build-arg NEXT_PUBLIC_API_BASE_URL=https://admin-api.dev.trovowallet.example.com/api/v1 \
  --build-arg NEXT_PUBLIC_SENTRY_DSN=https://examplePublicKey@errors.dev.trovo.app/1 \
  --build-arg NEXT_PUBLIC_APP_ENV=staging \
  --build-arg NEXT_PUBLIC_APP_VERSION=$(git rev-parse HEAD) \
  -t tm-web:local .
```

Then run it:

```bash
docker run --rm -p 3000:3000 tm-web:local
```

Remember: none of the `--build-arg` values above can be changed by an
environment variable passed to `docker run` — they're already compiled into
the JavaScript bundle. To change any `NEXT_PUBLIC_*` value you must rebuild
the image with new `--build-arg`s.

### 5.2 CI workflows — flagged as likely stale for this project

`/home/user/src-monorepo/.github/workflows/deploy.yml` and `pr-checks.yml`
exist at the monorepo root, but their path filters reference `backend/`,
`web/`, and `mobile/` (e.g. `paths: ["backend/**", "web/**"]`,
`dockerfile: web/Dockerfile`, `cache-dependency-path: web/package-lock.json`).
**None of those directory names exist in this monorepo** — the actual
top-level projects are `tm-web/`, `tm-api/`, `app-web/`, `app-backend/`,
`app-mobile/`, `payment-history-engine/`, and `wallet-core/`. `pr-checks.yml`
also references a `web/package-lock.json`, but this project has no
`package-lock.json` at all (it's Yarn-based, per §2 above).

This strongly suggests those two workflow files were copied in from a
differently-laid-out sibling repository and never adapted to this monorepo's
actual directory names — they likely do not currently run against `tm-web`
on push/PR. Treat the Dockerfile and `package.json` scripts described above
(verified directly from this project's own files) as the source of truth for
how to build and deploy `tm-web`, not those workflow files, until/unless
someone updates their path filters to match.

## 6. Error reporting (Sentry / GlitchTip)

This project uses the `@sentry/nextjs` SDK, but points it at a **self-hosted
GlitchTip** instance (GlitchTip implements Sentry's ingestion API, so the
standard Sentry SDK works against it unmodified — see the comment in
`.env.example`). There are three config files, one per Next.js runtime:

- `sentry.client.config.ts` — runs in the browser.
- `sentry.server.config.ts` — runs in the Node.js server process.
- `sentry.edge.config.ts` — runs in Edge middleware (`src/middleware.ts`).

All three are no-ops unless `NEXT_PUBLIC_SENTRY_DSN` is set — with no DSN,
`Sentry.init` is never called and the SDK stays completely dormant (see
`sentry.client.config.ts`'s guard). When a DSN *is* set:

- Only errors are captured (`tracesSampleRate: 0`) — GlitchTip doesn't do
  performance tracing, and this project treats response-time monitoring as
  Prometheus's job instead.
  This is a comment/assumption inside `sentry.client.config.ts`; no
  Prometheus/metrics wiring was found inside this project itself, so treat
  it as a stated intent rather than something verified in this repo.
- Session Replay is explicitly disabled (`replaysSessionSampleRate: 0`,
  `replaysOnErrorSampleRate: 0`) — this is an admin dashboard showing real
  customer data, so screen recordings are deliberately never sent anywhere.
- A `beforeSend` hook strips query strings and `Authorization`/`Cookie`
  headers from outgoing error reports, since admin URLs and headers here can
  carry session tokens.

Separately, `SENTRY_UPLOAD_SOURCEMAPS=1` plus `SENTRY_AUTH_TOKEN`,
`SENTRY_ORG`, `SENTRY_PROJECT`, and `SENTRY_URL` (all build-time, non-public
variables) opt a CI build into uploading source maps to GlitchTip after
`next build`, so stack traces on error reports show real component/file
names instead of minified output. This is deliberately opt-in and
independent of whether a DSN is set — a misconfigured org/project would
otherwise abort the whole build (see the extensive comments in
`next.config.mjs`), so it's off unless explicitly turned on. See
[CONFIGURATION.md](./CONFIGURATION.md) for each variable individually.
