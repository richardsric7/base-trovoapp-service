# Configuration

This document lists every environment variable `app-web` reads, at either
**build time** (Vite inlines `import.meta.env.VITE_*` values into the JS
bundle when you run `npm run build` / `npm start`) or **runtime** (read
by nginx inside the production Docker container, not by the JS bundle —
see [`DEPLOYMENT.md`](./DEPLOYMENT.md) for the mechanics of that
distinction).

This list was produced by searching the whole `src/` tree for
`import.meta.env` / `VITE_` usage and by reading `.env.example` and the
Dockerfile/`nginx.conf.template`. Nothing else in `src/` reads an
environment variable.

## How to set these locally

Copy the example file and edit it:

```bash
cp .env.example .env.local
```

`.env.local` is git-ignored (see `.gitignore`) — it's per-developer and
never committed. Vite automatically loads `.env.local` (and `.env`,
`.env.development`, etc.) when you run `npm start` / `npm run build`.

## Build-time variables (baked into the JS bundle by Vite)

### `VITE_API_URL`

- **Example value**: `https://api.dev.trovo.app`
- **What it configures**: The base URL `app-web` prefixes onto every REST
  API request to `app-backend`. Read in `src/store/config.ts`:
  ```ts
  export const BASE_URL = import.meta.env.VITE_API_URL ?? '/api';
  ```
  If unset, the app falls back to the relative path `/api`, which only
  resolves correctly when the app is served from behind the nginx
  container described in `DEPLOYMENT.md` (nginx proxies `/api/` to the
  real backend using its own, separate `BACKEND_URL` runtime variable).
  For the Vite **dev server** (`npm start`), there's no such proxy by
  default, so you almost always want to set this explicitly during local
  development.
- **How to get a real value**: Point it at wherever `app-backend` is
  running. If you're running `app-backend` locally, this is typically
  `http://localhost:8080` (or whatever port it listens on). If you're
  hitting a shared dev/staging deployment, ask a teammate for its URL or
  check the deployment's Portainer/infra config. See
  `app-backend/DEPLOYMENT.md` (in the sibling `app-backend` project) for
  how to stand up a backend instance to point at.

### `VITE_SOCKET_URL`

- **Example value**: `wss://api.dev.trovo.app/v1`
- **What it configures**: The WebSocket URL `app-web` connects to for
  real-time updates from `app-backend` (e.g. live balance/order/price
  updates). Read in `src/store/config.ts`:
  ```ts
  const wsHost = typeof window !== 'undefined' ? window.location.host : '';
  export const SOCKET_URL = import.meta.env.VITE_SOCKET_URL ?? `wss://${wsHost}/ws/v1`;
  ```
  If unset, it falls back to `wss://<current page's host>/ws/v1` — again,
  only correct behind the nginx proxy in production.
- **How to get a real value**: The WebSocket equivalent of `VITE_API_URL`
  — same backend instance, `ws://` instead of `http://` for a local
  backend (e.g. `ws://localhost:8080/v1`), or `wss://` for a TLS-fronted
  one. See `app-backend/DEPLOYMENT.md` for what path/port it exposes for
  WebSocket connections.

`.env.example` documents both of these with the same guidance:

```bash
# Hit the deployed dev backend directly (easiest):
VITE_API_URL=https://api.dev.trovo.app
VITE_SOCKET_URL=wss://api.dev.trovo.app/v1

# Or hit a locally running backend:
# VITE_API_URL=http://localhost:8080
# VITE_SOCKET_URL=ws://localhost:8080/v1
```

## Runtime variable (read by nginx, not by the JS bundle)

### `BACKEND_URL`

- **Example value**: `https://api.dev.trovo.app`
- **What it configures**: Set on the **production Docker container**
  (e.g. `docker run -e BACKEND_URL=... app-web`). It's substituted by
  `envsubst` into `nginx.conf.template` at container startup and used as
  the `proxy_pass` target for the `/api/` and `/ws/` nginx locations —
  i.e. it's what the *server* proxies same-origin requests to, not
  something the frontend JS reads directly. See
  [`DEPLOYMENT.md`](./DEPLOYMENT.md) for the full explanation of how this
  differs from `VITE_API_URL` above.
- **How to get a real value**: Same guidance as `VITE_API_URL` — point it
  at a running `app-backend` instance; see `app-backend/DEPLOYMENT.md`.

## Third-party keys found in the code (not currently environment-configurable)

### Flutterwave public key

`app-web` uses `flutterwave-react-v3` (see `src/pages/dashboard/home.tsx`)
for a fiat on-ramp payment flow ("Activate Trovo Account" / "Buy ETH and
TROV with fiat"). The `useFlutterwave` config in that file currently has
a **hardcoded placeholder public key**:

```ts
const config = {
  public_key: 'FLWPUBK-**************************-X',
  ...
};
```

This is **not** read from an environment variable today — it's a literal
string in the component. Flag this as a gap: a real deployment would need
this replaced with a genuine Flutterwave public key (and, ideally,
refactored to come from a `VITE_FLUTTERWAVE_PUBLIC_KEY`-style env var
instead of being hardcoded). To obtain a real key: sign up / log in at
the [Flutterwave Dashboard](https://dashboard.flutterwave.com/), go to
**Settings → API Keys**, and copy the **Public Key** for the environment
you need (test vs. live keys are separate). Flutterwave secret keys
should never be used in frontend code — only the public key belongs here.

## Things checked and found not to apply here

- **No error-tracking SDK** (e.g. Sentry) is present in `package.json` or
  imported anywhere in `src/` — there is no error-reporting env var to
  configure.
- **`web-vitals`** is used only via `src/reportWebVitals.ts`, which is the
  stock Create-React-App wiring: it measures CLS/FCP/LCP/TTFB and passes
  them to a callback the app provides. It does not send data anywhere by
  itself and has no associated configuration.
- No other `process.env.*` or `import.meta.env.*` reads exist anywhere in
  `src/` beyond `VITE_API_URL` and `VITE_SOCKET_URL`.
