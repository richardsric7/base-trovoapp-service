# Deployment

A step-by-step guide to getting `tm-web` running, written for someone doing
it for the first time. Commands run from the `tm-web/` folder unless a step
says otherwise.

## 1. What you are deploying

`tm-web` is **Trovo Manager's website**: the admin dashboard Trovo staff
use to manage users, assets, fees, payouts, P2P disputes, Public Markets,
partners and platform keys. It also holds the **organization portal**
(`/organisation`) where partner organizations (custodians, trustees,
issuing houses) manage their tokenized assets. It is a Next.js app; all
data comes from tm-api.

**It needs these to be running first:**

| What | Why | Where |
|---|---|---|
| tm-api | every page's data and every login | [tm-api/DEPLOYMENT.md](../tm-api/DEPLOYMENT.md) |
| app-backend, and a Trovo Wallet account for each admin | admins log in by approving a request in the Trovo app | [app-backend/DEPLOYMENT.md](../app-backend/DEPLOYMENT.md) |

## 2. Before you start

| Tool | Install | Check |
|---|---|---|
| Git | <https://git-scm.com/downloads> | `git --version` |
| Node.js 20.9 or newer (22 also works) | <https://nodejs.org/en/download> | `node --version` |
| Yarn 1 | `corepack enable` (comes with Node), or <https://classic.yarnpkg.com/en/docs/install> | `yarn --version` (shows `1.22.x`) |
| Docker (for production) | <https://docs.docker.com/get-docker/> | `docker --version` |

## 3. Get the code and install

```bash
git clone https://github.com/richardsric7/base-trovoapp-service.git
cd base-trovoapp-service/tm-web
yarn install --frozen-lockfile
```

**You should see:** `Done in ...s`. `--frozen-lockfile` installs exactly
what `yarn.lock` lists. (Use Yarn, not npm: `yarn.lock` is the lockfile.)

## 4. Set the parameters

```bash
cp .env.example .env.local
```

Set at least:

| Parameter | Example | How to get it |
|---|---|---|
| [`NEXT_PUBLIC_API_BASE_URL`](CONFIGURATION.md#next_public_api_base_url) | `http://localhost:8082/api/v1` | tm-api's address plus `/api/v1` |

Everything else is optional error reporting ([CONFIGURATION.md](CONFIGURATION.md)).
These values are copied into the website when it is built, so **change
them before building**, and never put a secret in a `NEXT_PUBLIC_` value.

## 5. Run it on your machine

```bash
yarn dev
```

**You should see:** `▲ Next.js 16...` and `Local: http://localhost:3000`.
Open <http://localhost:3000>: you are sent to **Sign in**.

**Log in:** enter the Trovo Wallet username of one of tm-api's
`DEFAULT_SUPER_ADMINS`. A QR code appears and an approval request arrives
in that user's Trovo app; approve it, and the dashboard opens. (The page
waits for the approval live and also checks every 5 seconds.)

## 6. Build and run for production

**A. With Docker (recommended):**

```bash
docker build \
  --build-arg NEXT_PUBLIC_API_BASE_URL=https://manager-api.trovo.example.com/api/v1 \
  --build-arg NEXT_PUBLIC_APP_ENV=production \
  --build-arg NEXT_PUBLIC_APP_VERSION=$(git rev-parse --short HEAD) \
  -t tm-web .
docker run -d --name tm-web --restart unless-stopped -p 3000:3000 tm-web
```

Add `--build-arg NEXT_PUBLIC_SENTRY_DSN=...` to turn on error reporting.
The build needs about 4 GB of memory (the image sets
`NODE_OPTIONS=--max-old-space-size=4096`).

**B. Without Docker:**

```bash
yarn build      # reads .env.local / the environment
yarn start      # serves on PORT, default 3000
```

**You should see** from the build: a list of routes ending with
`○ (Static)` / `ƒ (Dynamic)` and no errors.

**How the image works:** three stages on `node:20-alpine`: install
packages, build (where the build arguments are read), and a small image that
runs `yarn start` as an unprivileged user on port 3000.

Put HTTPS in front of it (your load balancer or a reverse proxy). If tm-api
sets [`CORS_ALLOWED_ORIGINS`](../tm-api/CONFIGURATION.md#cors_allowed_origins),
add tm-web's address to it.

## 7. Check it works

```bash
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:3000/sign-in
```

**You should see:** `200`. Then log in as in step 5. Point your hosting
platform's health check at `/sign-in` (other pages redirect there when
logged out).

## 8. Updating and rolling back

- **Updating:**
  ```bash
  git pull
  docker build --build-arg NEXT_PUBLIC_API_BASE_URL=... -t tm-web .
  docker rm -f tm-web
  docker run -d --name tm-web --restart unless-stopped -p 3000:3000 tm-web
  ```
  Deploy the matching tm-api first when a release adds API routes.
- **Rolling back:** `git checkout <previous-tag>`, then build and replace
  the container the same way; or keep tagged images (`tm-web:2026-10-06`)
  and run the previous one.
- More than one copy can run behind a load balancer; sessions are kept in
  the admin's browser.
- There is no CI/CD pipeline in this repository: run `yarn ci` (lint and
  build) before merging, and build and deploy the image by hand.

## 9. Troubleshooting

| You see | Cause | Fix |
|---|---|---|
| pages show network errors, or calls go to `localhost:3000/...` | `NEXT_PUBLIC_API_BASE_URL` missing or wrong when building | set it and **build again** |
| a new `NEXT_PUBLIC_...` value has no effect | they are copied in at build time | rebuild the image |
| browser console shows CORS errors | tm-api's `CORS_ALLOWED_ORIGINS` doesn't list tm-web's address | add it in tm-api |
| "not an admin" at sign-in | the username is not a tm-api admin | add it to tm-api's `DEFAULT_SUPER_ADMINS`, or have an admin invite them |
| the QR page waits and never logs in | the approval never reached the phone, or tm-api cannot reach app-backend | see tm-api's troubleshooting table |
| sent back to sign-in every 15 minutes | the login token expired and could not be refreshed | log in again; check tm-api is reachable |
| `JavaScript heap out of memory` during the build | not enough memory | build with `NODE_OPTIONS=--max-old-space-size=4096` (the image does) |
| the build stops at "sentry-cli" | source-map upload with a wrong organization or project | fix `SENTRY_ORG` / `SENTRY_PROJECT`, or unset `SENTRY_UPLOAD_SOURCEMAPS` |
| `yarn install` fails with a lockfile error | `package.json` changed without `yarn.lock` | run `yarn install` and commit both |
