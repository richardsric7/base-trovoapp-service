# tm-web — Trovo Manager (Admin Dashboard)

`tm-web` is the internal admin dashboard for the Trovo Wallet platform. It is
a [Next.js](https://nextjs.org/) + TypeScript single-page-style application
used by Trovo staff to manage platform data: users, curated/tokenized assets,
service links, P2P trades and appeals, fee/payment configuration, compliance
and AML review, KYC, vault-signer secrets, audit trails, and more.

It also hosts a smaller, separate **organisation/stakeholder portal** (under
`/organisation`) used by external partners (asset custodians, trustees,
issuing houses, legal/financial advisers, rating agencies) to manage their
own tokenized assets. Both portals live in this one app but authenticate and
call the backend independently of each other (see [INTEGRATION.md](./INTEGRATION.md)).

`npm run test`-style unit tests are not part of this project's current setup;
correctness is enforced via `yarn lint` and `next build` type-checking (see
`yarn ci` in `package.json`).

## Role in the monorepo

This is a **pure client application with no database and (practically) no
server-side business logic of its own**. Every piece of real data — users,
assets, trades, fees, KYC records, etc. — is fetched from and written to the
sibling Go project **`tm-api`** over HTTP/REST. `tm-web` renders UI, manages
client-side auth state, and calls `tm-api`; it does not talk to any other
service in the monorepo directly. See [INTEGRATION.md](./INTEGRATION.md) for
the full picture, including why `tm-web` has no relationship with
`app-backend`, `app-web`, or `app-mobile`.

A check of `src/app` confirms there are no Next.js Route Handlers
(`route.ts`) anywhere in this project — `tm-web` exposes no API of its own.
The `external-api-clients` folder under `src/app` is a *page* (for managing
external API client credentials issued by `tm-api`), not an API route.

## Tech stack

| Concern | Choice |
|---|---|
| Framework | [Next.js 14](https://nextjs.org/) (`^14.2.35`), **App Router** (`src/app/`) |
| Language | TypeScript |
| UI library | [Ant Design 5](https://ant.design/) (`antd`, `@ant-design/icons`, `@ant-design/nextjs-registry`) as the component kit |
| Styling | [styled-components](https://styled-components.com/) (`^6.1.11`), enabled via the `compiler.styledComponents` option in `next.config.mjs` — most custom/bespoke UI (forms, layout, one-off pages) is hand-styled with `styled-components`, while `antd` supplies tables, modals, date pickers, etc. |
| State / data fetching | [Redux Toolkit](https://redux-toolkit.js.org/) + **RTK Query**, via `react-redux`. Two separate RTK Query API slices exist: `baseApi` (admin dashboard, `src/redux/baseApi/index.ts`) and `orgApi` (organisation portal, `src/redux/baseApi/orgApi.ts`) — both use a **custom Axios-based `baseQuery`** (`axiosBasedQuery.ts` / `orgBasedQuery.ts`) rather than `fetchBaseQuery`, so all HTTP calls actually go through `axios` with its own interceptors. |
| Forms | [Formik](https://formik.org/) + [Yup](https://github.com/jquense/yup) validation |
| Charts | [ApexCharts](https://apexcharts.com/) via `react-apexcharts` |
| Drag & drop | `@dnd-kit/*` |
| Error reporting | [Sentry SDK](https://docs.sentry.io/platforms/javascript/guides/nextjs/) (`@sentry/nextjs`) pointed at a self-hosted **GlitchTip** instance (Sentry-API-compatible) — see [DEPLOYMENT.md](./DEPLOYMENT.md) |
| Package manager | **Yarn 1** (`yarn.lock` is the lockfile; `packageManager` in `package.json` pins `yarn@1.22.22`) |

## Directory structure (`src/`)

```
src/
├── app/                        # Next.js App Router — route-group based
│   ├── (auth)/                 # Public admin auth routes: /sign-in, /qr-scan, /invite-expired, ...
│   ├── (dashboard)/            # The admin dashboard proper — one folder per feature area, e.g.:
│   │   ├── users/              #   /users
│   │   ├── p2pusers/           #   /p2pusers, p2pappeal/ etc. — P2P trading & disputes
│   │   ├── assetcuration/      #   curated assets
│   │   ├── complaince-and-aml/ #   compliance / AML review workflows
│   │   ├── vault-signer/       #   vault-signer secret management UI
│   │   ├── servicelinks/       #   service link config
│   │   ├── feeexemptions/      #   accounts exempt from platform service fees
│   │   ├── dividendandyield/   #   proceeds payouts: stages, approvals, schedule, fee & engine controls, reports (tm-api /proceed-payouts)
│   │   ├── publicmarkets/      #   Public Markets (tokenized NGX/FMDQ instruments): overview, assets, orders & approvals, escalations,
│   │   │                       #   reconciliation, corporate actions, exchange partners, wallet provisioning, prices, settings (tm-api /public-markets)
│   │   ├── payment/, revenue/, earlyexit/, ...
│   │   └── settings/, audit-trail/, system-health/, notification/, ...
│   ├── organisation/           # SEPARATE stakeholder portal (external orgs), own auth (see INTEGRATION.md)
│   ├── external-api-clients/   # Admin UI for managing tm-api's external API client credentials
│   ├── middleware.ts           # Route guard — redirects unauthenticated requests to /sign-in or /organizations/login
│   └── layout.tsx / page.tsx   # Root layout and landing page
├── redux/
│   ├── baseApi/                # RTK Query setup: axios-based baseQuery, admin (baseApi) + org (orgApi) API slices, tag types
│   ├── api/                    # One folder per domain, each `injectEndpoints`-ing into baseApi, e.g.:
│   │   ├── auth/                #   login, verifyLogin, logout
│   │   ├── users/, admin/, organizations/, p2p/, payment/, assettokenization/, ...
│   │   ├── publicMarkets/       #   Public Markets (tm-api /public-markets)
│   │   └── vaultSigner/, compliance/, auditTrail/, jobs/, observability/, ...
│   ├── slices/                 # Plain Redux slices: authSlice (admin session), orgSlice (org session)
│   ├── store.ts                # configureStore wiring both API slices + both auth slices
│   └── refreshAccessToken.ts   # Single-flight access-token refresh used by the axios interceptor
├── components/                 # Shared/reusable UI components
├── config/                     # Runtime config, chiefly `BASE_URL` (tm-api base URL)
├── hooks/, utils/, lib/, constants/  # Shared helpers, localStorage/cookie helpers, constants
├── observability.ts            # Sentry helpers shared across client/server/edge configs
└── middleware.ts                # (see above)
```

At the repo root, `sentry.client.config.ts`, `sentry.server.config.ts`, and
`sentry.edge.config.ts` configure the three Sentry/GlitchTip runtimes Next.js
uses (browser, Node.js server, and Edge middleware respectively).

## Running locally

Prerequisites: Node.js 20.x and Yarn (see [DEPLOYMENT.md](./DEPLOYMENT.md#prerequisites)
for install links). You'll also need `tm-api` running somewhere reachable —
see `tm-api/DEPLOYMENT.md` in this monorepo for how to stand it up locally.

```bash
cd tm-web
yarn install
cp .env.example .env.local
# edit .env.local: set NEXT_PUBLIC_API_BASE_URL to your local tm-api,
# e.g. http://localhost:8082/api/v1
yarn dev
```

Then open [http://localhost:3000](http://localhost:3000). Sign-in uses a
QR-code flow tied to the Trovo mobile app (see [INTEGRATION.md](./INTEGRATION.md#authentication-flow))
rather than a username/password form, so you'll need a way to complete that
against whichever `tm-api`/backend environment you pointed at.

## Further reading

- **[DEPLOYMENT.md](./DEPLOYMENT.md)** — installing, building, running in
  production, the Dockerfile, and Next.js build-time vs runtime env var
  handling.
- **[CONFIGURATION.md](./CONFIGURATION.md)** — every environment variable
  this app reads, with example values and how to obtain real ones.
- **[INTEGRATION.md](./INTEGRATION.md)** — how `tm-web` authenticates against
  and calls `tm-api`, and how it relates (or doesn't) to the rest of the
  monorepo.
- [compliance-aml-flow.md](./compliance-aml-flow.md) — feature-specific notes
  on the compliance/AML review flow.
- [vault-signer-frontend-plan.md](./vault-signer-frontend-plan.md) —
  feature-specific notes on the vault-signer UI.
