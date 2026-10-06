# Configuration

Every parameter `tm-web` reads. Each one says what it is, why it is needed,
whether you must set it, an example and how to get a real value.

## How to set them

tm-web (Trovo Manager's website) is a Next.js app. **All its settings are
build-time settings**: they are read when you run `yarn build` (or
`yarn dev`) and copied into the website. Changing one means building again;
setting it on a running container does nothing.

- **Locally:** `cp .env.example .env.local`, then edit `.env.local` (never
  committed). `yarn dev` reads it.
- **Docker:** pass each as a build argument:
  `docker build --build-arg NEXT_PUBLIC_API_BASE_URL=https://manager-api.trovo.example.com/api/v1 -t tm-web .`

Settings starting with `NEXT_PUBLIC_` end up in the browser and are visible
to anyone who opens the website, so never put a secret in one.

The only required setting is [`NEXT_PUBLIC_API_BASE_URL`](#next_public_api_base_url).
Everything else is for error reporting and is optional.

---

## 1. Connection to tm-api

### `NEXT_PUBLIC_API_BASE_URL`

- **What it is:** The address of tm-api's API, including `/api/v1`. Read in
  `src/config/index.ts`; every page's data, the login and the token
  refresh use it.
- **Why it's needed:** tm-web has no data of its own; everything comes from
  tm-api.
- **Required:** Yes. Unset, calls go to the website's own address and fail.
- **Example:** `http://localhost:8082/api/v1` (tm-api on your machine) or
  `https://manager-api.trovo.example.com/api/v1`
- **How to get it:** tm-api's public address
  ([tm-api/DEPLOYMENT.md](../tm-api/DEPLOYMENT.md); its default port is
  `8082`) followed by `/api/v1`. It must be reachable **from the admins'
  browsers**, not only from the server. If tm-api sets
  `CORS_ALLOWED_ORIGINS`, tm-web's address must be in that list.

## 2. Error reporting (GlitchTip / Sentry), all optional

tm-web can send browser and server errors to GlitchTip (or Sentry, which
uses the same protocol). With no DSN, nothing is sent and the build is
unchanged.

### `NEXT_PUBLIC_SENTRY_DSN`

- **What it is:** The address errors are sent to (the project's "DSN").
- **Why it's needed:** To be told about crashes admins hit.
- **Required:** No. Unset turns error reporting off.
- **Example:** `https://0123abcd@errors.trovo.example.com/3` (placeholder)
- **How to get it:** In GlitchTip: your project → **Settings → Client Keys
  (DSN)** → copy the DSN.

### `NEXT_PUBLIC_APP_ENV`

- **What it is:** The environment name shown on each error report.
- **Why it's needed:** To tell production errors from test errors.
- **Required:** No, default `development`.
- **Example:** `production`
- **How to get it:** `development`, `staging` or `production`.

### `NEXT_PUBLIC_APP_VERSION`

- **What it is:** The version shown on each error report (and in
  `src/observability.ts`).
- **Why it's needed:** To know which build an error came from.
- **Required:** No, default `unknown`.
- **Example:** `a1b2c3d`
- **How to get it:** `git rev-parse --short HEAD` when building.

## 3. Source-map upload, all optional

Source maps make error reports show real file names and lines. Uploading
them is **off unless `SENTRY_UPLOAD_SOURCEMAPS=1` and all four values
below are set**, because a wrong organization or project name stops the
build. Without upload, error reporting still works.

### `SENTRY_UPLOAD_SOURCEMAPS`

- **What it is:** `1` turns source-map upload on.
- **Why it's needed:** It is a deliberate switch: having the credentials
  alone does not turn upload on.
- **Required:** No.
- **Example:** `1`
- **How to get it:** Set it only after checking the four values below.

### `SENTRY_AUTH_TOKEN`

- **What it is:** A token allowed to upload to GlitchTip. A **secret**:
  give it only to the build, never as `NEXT_PUBLIC_`.
- **Why it's needed:** To upload source maps.
- **Required:** For upload.
- **Example:** `<your-glitchtip-auth-token>`
- **How to get it:** GlitchTip → **Profile → Auth Tokens → Create**, with
  `project:releases` permission.

### `SENTRY_ORG`

- **What it is:** Your GlitchTip organization's short name ("slug").
- **Why it's needed:** Where to upload.
- **Required:** For upload.
- **Example:** `trovo`
- **How to get it:** The first part of the address after logging in:
  `https://errors.trovo.example.com/<org-slug>/<project-slug>/`.

### `SENTRY_PROJECT`

- **What it is:** The project's short name.
- **Why it's needed:** Where to upload.
- **Required:** For upload.
- **Example:** `trovo-manager-web`
- **How to get it:** The second part of that address.

### `SENTRY_URL`

- **What it is:** The address of your GlitchTip server.
- **Why it's needed:** Where to upload.
- **Required:** For upload.
- **Example:** `https://errors.trovo.example.com`
- **How to get it:** The address you log in to GlitchTip at.

## 4. Set for you

### `NODE_ENV`

- **What it is:** `production` or `development`, set by Next.js itself
  (`yarn build`/`yarn start` use `production`, `yarn dev` uses
  `development`) and by the Docker image.
- **Why it's needed:** Next.js optimizes production builds. tm-web's own
  code does not read it.
- **Required:** No; don't set it yourself.
- **Example:** `production`
- **How to get it:** Nothing to do.

### `PORT`

- **What it is:** The port `yarn start` listens on.
- **Why it's needed:** To run on a different port.
- **Required:** No, default `3000` (the Docker image sets `3000`).
- **Example:** `3000`
- **How to get it:** Any free port.

`NEXT_RUNTIME` is set by Next.js to load the right error-reporting setup;
`NEXT_TELEMETRY_DISABLED=1` (set in the image) stops Next.js sending usage
statistics.
