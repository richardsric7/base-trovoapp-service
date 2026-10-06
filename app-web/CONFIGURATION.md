# Configuration

Every parameter `app-web` reads. Each one says what it is, why it is
needed, whether you must set it, an example and how to get a real value.

## How to set them

app-web is a website made of static files. Its settings come in two kinds,
and the difference matters:

| Kind | Read when | Where you set it | Changing it needs |
|---|---|---|---|
| **Build-time** (`VITE_...`) | when you run `npm run build` or `npm start`; the value is copied into the website's files | `.env.local` in `app-web/`, or the shell before building | a rebuild |
| **Run-time** (`BACKEND_URL`) | when the Docker container starts | `docker run -e ...` or your hosting platform | a container restart |

Locally: `cp .env.example .env.local`, then edit it (`.env.local` is never
committed). Build-time values are visible to anyone who opens the website,
so **never put a secret in a `VITE_` variable**.

**You usually need only one of them:**

- **Docker (production):** set only `BACKEND_URL`. The website calls its
  own address (`/api`, `/ws`) and the container forwards those calls to
  app-backend.
- **Dev server (`npm start`):** set `VITE_API_URL`, because there is no
  forwarding.

---

## 1. Connection to app-backend

### `VITE_API_URL`

- **What it is:** The web address of app-backend that the website calls.
  Read in `src/store/config.ts`.
- **Why it's needed:** Every screen loads its data from app-backend.
- **Required:** For the dev server, yes. With Docker, no: when unset the
  website calls `/api` on its own address, which the container forwards
  (see `BACKEND_URL`).
- **Example:** `http://localhost:8080` (app-backend on your machine) or
  `https://api.trovo.example.com`
- **How to get it:** app-backend's address in that environment: `http://localhost:8080`
  if you run it yourself ([app-backend/DEPLOYMENT.md](../app-backend/DEPLOYMENT.md)),
  otherwise ask whoever runs the shared server. No trailing `/`.

### `VITE_SOCKET_URL`

- **What it is:** A websocket address for live updates. `src/store/config.ts`
  reads it into `SOCKET_URL`, but **nothing uses `SOCKET_URL` today**: the
  website opens no websocket and refreshes by asking again (for example an
  order page every 5 seconds).
- **Why it's needed:** It isn't, yet. It is kept for a future live-update
  feature.
- **Required:** No.
- **Example:** `ws://localhost:8080/v1` (local) or
  `wss://api.trovo.example.com/v1`
- **How to get it:** If you set it, use the address of `VITE_API_URL` with
  `ws://` instead of `http://` (or `wss://` instead of `https://`), plus
  `/v1`.

### `BACKEND_URL`

- **What it is:** The address of app-backend that the **Docker container**
  forwards `/api/...` and `/ws/...` to. It is written into the web
  server's (nginx) settings when the container starts
  (`nginx.conf.template`); the website's code never reads it.
- **Why it's needed:** It lets one built image serve any environment: the
  website calls its own address and the container passes calls on.
- **Required:** Yes, when you run the Docker image (nginx refuses to start
  without it). Not used by the dev server.
- **Example:** `http://app-backend.internal:8080` or
  `https://api.trovo.example.com`
- **How to get it:** app-backend's address as seen **from the container**
  (inside Docker, `localhost` is the container itself; use the service
  name or a real host name). No trailing `/`.

## 2. Dev server only

These affect only `npm start` on your own machine.

### `DISPLAY`, `WAYLAND_DISPLAY`

- **What they are:** Set by Linux desktops to say a screen is available.
  `vite.config.mts` opens a browser tab on `npm start` only on macOS,
  Windows, or Linux when one of these is set.
- **Why they're needed:** On a Linux machine without a screen (a
  container, a server over SSH) there is no browser to open, and trying
  printed an `xdg-open` error.
- **Required:** No. Your desktop sets them; don't set them yourself.
- **Example:** `:0`
- **How to get them:** Nothing to do.

### `BROWSER`

- **What it is:** Which browser the dev server opens; `none` opens none.
- **Why it's needed:** To stop the dev server opening a tab you don't want.
- **Required:** No.
- **Example:** `none`
- **How to get it:** `BROWSER=none npm start`.

The dev server's port is fixed at `3000` in `vite.config.mts`.

## 3. Values fixed in the code

These are not settings today; changing them means editing the code.

### Flutterwave public key

- **What it is:** The public key of the Flutterwave account used by the
  "buy with card" pop-up (`src/pages/dashboard/home.tsx`,
  `public_key: 'FLWPUBK-**************************-X'`).
- **Why it's needed:** Flutterwave identifies the merchant by it. The
  committed value is a placeholder, so card purchases cannot work until it
  is replaced.
- **Required:** Only if you use the Flutterwave purchase.
- **Example:** `FLWPUBK_TEST-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx-X` (placeholder)
- **How to get it:** [Flutterwave dashboard](https://dashboard.flutterwave.com/)
  → **Settings → API Keys** → **Public Key** (test and live keys differ).
  Only the public key belongs in a website, never the secret key.

### Block explorer links

- **What they are:** `https://basescan.org/tx/` and
  `https://sepolia.basescan.org/tx/` (`src/store/constants.ts`), used to
  link a transaction to the block explorer.
- **Why they're needed:** So users can look a transaction up.
- **Required:** Already set.
- **Example:** `https://basescan.org/tx/`
- **How to get them:** Change them only for another network.

Nothing else is read from the environment: the code reads only
`VITE_API_URL` and `VITE_SOCKET_URL` (unused), and `vite.config.mts` only `DISPLAY`
and `WAYLAND_DISPLAY`. There is no error-reporting service to configure.
