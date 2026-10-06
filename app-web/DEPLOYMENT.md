# Deployment

A step-by-step guide to getting `app-web` running, written for someone
doing it for the first time. Commands run from the `app-web/` folder unless
a step says otherwise.

## 1. What you are deploying

`app-web` is the **Trovo web wallet**: the website where users sign up,
see balances, send and receive money, trade P2P, buy tokenized assets and
Public Markets stocks and bonds. It is a set of static files (HTML,
JavaScript, CSS); every piece of data comes from app-backend. Keys are made
and used only in the user's browser.

In production it runs as a small Docker container: a web server (nginx)
that serves the files and forwards `/api/...` calls to app-backend.

**It needs these to be running first:**

| What | Why | Where |
|---|---|---|
| app-backend | every screen loads its data from it | [app-backend/DEPLOYMENT.md](../app-backend/DEPLOYMENT.md) |

## 2. Before you start

| Tool | Install | Check |
|---|---|---|
| Git | <https://git-scm.com/downloads> | `git --version` |
| Node.js 20.19 or newer (22 also works) | <https://nodejs.org/en/download> (or `nvm install 22`) | `node --version` |
| npm (comes with Node) | — | `npm --version` |
| Docker (for production) | <https://docs.docker.com/get-docker/> | `docker --version` |

## 3. Get the code and install

```bash
git clone https://github.com/richardsric7/base-trovoapp-service.git
cd base-trovoapp-service/app-web
npm ci
```

**You should see:** `added NNN packages`. `npm ci` installs the exact
versions in `package-lock.json`.

## 4. Run it on your machine (dev server)

1. Tell it where app-backend is:
   ```bash
   cp .env.example .env.local
   ```
   The template points at `http://localhost:8080` (app-backend on your
   machine). For a shared server, change
   [`VITE_API_URL`](CONFIGURATION.md#vite_api_url).
2. Start it:
   ```bash
   npm start
   ```
   **You should see:** `VITE v8... ready` and `Local: http://localhost:3000/`.
   A browser tab opens if your computer has a screen.
3. Open <http://localhost:3000>, choose **Create account**, and finish
   sign-up. The account appears in app-backend's `users` table.

Changes to the code reload the page by themselves. Changes to `.env.local`
need `npm start` again.

## 5. Build for production

```bash
npm run build
```

This checks the TypeScript code (`tsc`) and then builds (`vite build`).
**You should see:** `✓ built in ...`, and a `dist/` folder with
`index.html` and an `assets/` folder. A warning that some chunks are
larger than 500 kB is expected. Any type error stops the build.

To look at the built website: `npm run preview`, then open the address it
prints.

## 6. Run it with Docker (production)

```bash
docker build -t app-web .
docker run -d --name app-web --restart unless-stopped -p 8081:80 \
  -e BACKEND_URL=https://api.trovo.example.com app-web
```

| Parameter | Example | What it is |
|---|---|---|
| [`BACKEND_URL`](CONFIGURATION.md#backend_url) | `https://api.trovo.example.com` | app-backend's address, as seen from the container. Required. |
| `-p 8081:80` | `8081` | the port on your server; the container listens on `80` |

Leave `VITE_API_URL` **unset** when building the image: the website then
calls `/api` on its own address, and the container forwards those calls to
`BACKEND_URL`. One image works for every environment; only `BACKEND_URL`
changes. (If you do set `VITE_API_URL` before building, that address is
copied into the files and `BACKEND_URL` is ignored for API calls.)

Put HTTPS in front of the container (your load balancer or a reverse
proxy such as Caddy or Traefik): browsers only allow some features, such as
the clipboard, on HTTPS pages.

**How the image works:** it builds the website with Node 20, then copies
`dist/` into an `nginx:1.27-alpine` image. When the container starts,
nginx writes `BACKEND_URL` into its settings (`nginx.conf.template`):
`/api/...` and `/ws/...` go to app-backend, and every other address
returns `index.html`, so links like `/dashboard/bank` work after a page
reload.

## 7. Check it works

```bash
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8081/
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8081/api/swagger/index.html
```

**You should see:** `200` twice. The second proves the container reaches
app-backend. Then open <http://localhost:8081> and sign in.

Point your hosting platform's health check at `/`.

## 8. Updating and rolling back

- **Updating:**
  ```bash
  git pull
  docker build -t app-web .
  docker rm -f app-web
  docker run -d --name app-web --restart unless-stopped -p 8081:80 \
    -e BACKEND_URL=https://api.trovo.example.com app-web
  ```
  Deploy the matching app-backend first when a release needs new API
  routes.
- **Rolling back:** `git checkout <previous-tag>`, then build and replace
  the container the same way. Tag images (`app-web:2026-10-06`) to roll
  back without rebuilding.
- More than one copy can run behind a load balancer; they keep no state.
- There is no CI/CD pipeline in this repository: run `npm run build` before
  merging, and build and deploy the image by hand.

## 9. Troubleshooting

| You see | Cause | Fix |
|---|---|---|
| the dev server's pages show network errors | `VITE_API_URL` missing or wrong, or app-backend is not running | set it in `.env.local` and run `npm start` again |
| `xdg-open ENOENT` on `npm start` | an older version on a Linux machine without a screen | update; or `BROWSER=none npm start` |
| `npm ci` fails with "lock file out of sync" | `package.json` changed without `package-lock.json` | run `npm install` and commit both files |
| the build fails with TypeScript errors | a type error in the code | fix the file and line it names |
| the container exits at start | `BACKEND_URL` not set | add `-e BACKEND_URL=...` |
| `502 Bad Gateway` on `/api/...` | the container cannot reach `BACKEND_URL` | use an address reachable from the container (not `localhost`) |
| a page shows `404` after reload on another host | that web server doesn't send unknown paths to `index.html` | serve with the provided image, or configure that fallback |
| a change to `VITE_API_URL` has no effect in production | it is copied into the files at build time | rebuild the image |
| "buy with card" does nothing | the Flutterwave key is a placeholder | see [CONFIGURATION.md](CONFIGURATION.md#flutterwave-public-key) |
