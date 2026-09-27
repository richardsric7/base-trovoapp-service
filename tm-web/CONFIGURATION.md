# Configuration Reference

Every environment variable `tm-web` reads, found via:

```bash
grep -rn "process\.env\." src/ next.config.mjs sentry.client.config.ts sentry.server.config.ts sentry.edge.config.ts
```

cross-referenced against the committed `.env.example` at the project root.

All the variables below except the `SENTRY_*` build/upload ones are
**`NEXT_PUBLIC_*`**, meaning they are compiled into the browser bundle at
**build time**. See [DEPLOYMENT.md §3](./DEPLOYMENT.md#3-environment-variables--the-nextjs-nuance)
for why that matters. Set them before `yarn build` (or pass them as Docker
`--build-arg`s, per the Dockerfile) — setting them only at container runtime
has no effect on any of the `NEXT_PUBLIC_*` values.

For local development, copy the template and fill it in:

```bash
cp .env.example .env.local
```

## Core

### `NEXT_PUBLIC_API_BASE_URL`

- **Example**: `https://admin-api.dev.trovowallet.example.com/api/v1`
- **What it is**: The base URL of the `tm-api` REST API that this dashboard
  calls for literally all of its data (users, assets, trades, KYC, etc.).
  Read in `src/config/index.ts` as `BASE_URL`, and used everywhere RTK
  Query endpoints and the auth/token-refresh code build request URLs (e.g.
  `src/redux/baseApi/axiosBasedQuery.ts`, `src/redux/refreshAccessToken.ts`).
- **How to get a real value**: Point it at wherever `tm-api` is actually
  running for your environment. For local development that's typically
  `http://localhost:8082/api/v1` (`tm-api`'s default port, per the comment in
  `.env.example`, which cites `tm-api/internal/server/models/server.go`).
  See `tm-api/DEPLOYMENT.md` in this monorepo for how to stand up a `tm-api`
  instance and confirm its actual base path/port. For a shared
  staging/production deployment, use whatever hostname ops/infra has that
  environment's `tm-api` reachable at, including the `/api/v1` path prefix.

## Error reporting (GlitchTip / Sentry-compatible) — all optional

`tm-web` uses the `@sentry/nextjs` SDK pointed at a self-hosted **GlitchTip**
instance (a Sentry-API-compatible error tracker — not sentry.io). Every
variable in this section is optional: leave `NEXT_PUBLIC_SENTRY_DSN` empty
and the SDK never initializes, so none of the others matter either.

### `NEXT_PUBLIC_SENTRY_DSN`

- **Example**: `https://a1b2c3d4e5f6@errors.dev.trovo.app/7`
- **What it is**: The DSN (Data Source Name) that tells the Sentry SDK where
  to send error events, and identifies which GlitchTip project they belong
  to. Read in `src/observability.ts`, `sentry.client.config.ts`,
  `sentry.server.config.ts`, and `sentry.edge.config.ts`. Also checked in
  `next.config.mjs` to decide whether to wrap the Next.js config with the
  Sentry build plugin at all.
- **How to get a real value**: Log into the project's GlitchTip instance
  (per `.env.example`, this project's is at `https://errors.dev.trovo.app/`),
  open the relevant project's **Settings → Client Keys (DSN)** page, and copy
  the DSN shown there. If no such project exists yet, create one in
  GlitchTip first.

### `NEXT_PUBLIC_APP_ENV`

- **Example**: `staging` (other realistic values: `development`, `production`)
- **What it is**: A free-text label attached to every error report so you can
  tell which environment (dev/staging/prod) an error came from. Defaults to
  `"development"` if unset. Read in `src/observability.ts` and all three
  Sentry config files as `environment`.
- **How to get a real value**: Not fetched from anywhere — just set it to
  match whatever environment you're deploying (`development`, `staging`,
  `production`, etc.), consistently with how other services in this
  environment are labeled.

### `NEXT_PUBLIC_APP_VERSION`

- **Example**: `a1b2c3d4e5f67890abcdef1234567890abcdef12` (a git commit SHA)
- **What it is**: A build/release identifier attached to error reports (as
  `release` in the Sentry config), so a GlitchTip error can be traced back to
  the exact build that produced it. Per `.env.example`, "CI passes the git
  sha."
- **How to get a real value**: In CI, use the commit SHA being built, e.g.
  `$(git rev-parse HEAD)`. Locally it can be left empty (defaults to
  `"unknown"` per `src/observability.ts`).

### `SENTRY_UPLOAD_SOURCEMAPS`

- **Example**: `1`
- **What it is**: Explicit opt-in switch (checked in `next.config.mjs`) for
  uploading source maps to GlitchTip after a production build, so stack
  traces in error reports show real file/function names instead of minified
  code. Deliberately **not** automatic just because credentials are present —
  a wrong org/project slug makes the underlying `sentry-cli` abort the whole
  build, so this must be turned on knowingly. Only relevant in CI/production
  builds, not local dev.
- **How to get a real value**: Set to `1` only when you also have all four
  of `SENTRY_AUTH_TOKEN`, `SENTRY_ORG`, `SENTRY_PROJECT`, and `SENTRY_URL`
  below correctly configured for a real GlitchTip project.

### `SENTRY_AUTH_TOKEN`

- **Example**: `glitchtip_sample_token_abcdef1234567890` *(placeholder — real
  tokens are opaque secrets, never guess or reuse an example value)*
- **What it is**: An auth token used by the Sentry build plugin's
  `sentry-cli` to create a release and upload source maps to GlitchTip.
  Server-side/build-time only — never shipped to the browser.
- **How to get a real value**: In GlitchTip, go to your user or
  organization's **Settings → API Tokens** (or the project's own settings,
  depending on GlitchTip version) and generate a token with permission to
  create releases and upload artifacts for the relevant project.

### `SENTRY_ORG` / `SENTRY_PROJECT`

- **Example**: `SENTRY_ORG=trovo`, `SENTRY_PROJECT=tm-web`
- **What it is**: The GlitchTip organization and project slugs that source
  maps get uploaded to.
- **How to get a real value**: Per the comment in `.env.example`, read them
  directly out of the GlitchTip URL after logging in:
  `https://errors.dev.trovo.app/<org-slug>/<project-slug>/`.

### `SENTRY_URL`

- **Example**: `https://errors.dev.trovo.app`
- **What it is**: The base URL of the GlitchTip instance itself (defaults to
  the value shown in `.env.example`), so `sentry-cli` knows which server to
  talk to instead of the default `sentry.io`.
- **How to get a real value**: Use the URL of whichever GlitchTip deployment
  this environment's errors should go to — ask whoever administers the
  error-reporting stack if you're unsure it's still `https://errors.dev.trovo.app`.

## Not application config (context only)

These aren't variables `tm-web`'s own code reads for its behavior, but show
up around the build/runtime and are worth knowing about:

- **`NODE_OPTIONS`** — set to `--max-old-space-size=4096` in the Dockerfile's
  build stage only, to give Node enough heap to survive the Sentry webpack
  plugin instrumenting every module during `yarn build`. Not something you
  need to set yourself unless reproducing an out-of-memory build failure.
- **`PORT`** — set to `3000` in the Dockerfile's runtime stage; `next start`
  reads this to decide which port to listen on.
- **`NEXT_RUNTIME`** — set automatically by Next.js itself (not by you) to
  distinguish the `nodejs` vs `edge` runtime in `src/instrumentation.ts`,
  which loads the matching `sentry.server.config.ts` or
  `sentry.edge.config.ts`.
