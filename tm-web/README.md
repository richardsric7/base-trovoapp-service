# tm-web — Trovo Manager (admin dashboard)

## What this project does

`tm-web` is the website of **Trovo Manager**, the admin dashboard Trovo
staff use to run the platform. From it, staff:

- look up and manage users (KYC, suspensions, fee exemptions);
- curate assets and manage tokenized assets, their payouts (dividends and
  interest) and Public Markets (tokenized Nigerian stocks and bonds);
- resolve P2P trading disputes and review compliance and AML cases;
- manage partner businesses ("service links"), fees, platform signing keys
  (vault manager), admins and their permissions, and see audit trails and
  system health.

It also hosts the **organization portal** (`/organisation`), where partner
organizations (custodians, trustees, issuing houses, advisers) manage their
own tokenized assets, with a login of their own.

Admins log in without a password: they approve the login in their Trovo
mobile app. tm-web has no data of its own; everything comes from
[`tm-api`](../tm-api/README.md).

Next: [DEPLOYMENT.md](DEPLOYMENT.md) to run it,
[CONFIGURATION.md](CONFIGURATION.md) for its settings and
[INTEGRATION.md](INTEGRATION.md) for what it connects to.

## Tech stack

| Concern | Choice |
|---|---|
| Framework | [Next.js 16](https://nextjs.org/), App Router (`src/app/`) |
| Language | TypeScript 7, React 19 |
| UI library | [Ant Design 6](https://ant.design/) (`antd`, `@ant-design/icons`, `@ant-design/nextjs-registry`) |
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
└── middleware.ts                # Route guard: sends logged-out visitors to /sign-in or /organizations/login
```

At the repo root, `sentry.client.config.ts`, `sentry.server.config.ts`, and
`sentry.edge.config.ts` configure the three Sentry/GlitchTip runtimes Next.js
uses (browser, Node.js server, and Edge middleware respectively).

## Running locally

The short version (the full guide is [DEPLOYMENT.md](DEPLOYMENT.md)), with
Node.js 20.9 or newer, Yarn 1 and tm-api running at `http://localhost:8082`:

```bash
git clone https://github.com/richardsric7/base-trovoapp-service.git
cd base-trovoapp-service/tm-web
yarn install --frozen-lockfile
cp .env.example .env.local     # NEXT_PUBLIC_API_BASE_URL=http://localhost:8082/api/v1
yarn dev                       # http://localhost:3000
```

Sign in with the Trovo Wallet username of one of tm-api's
`DEFAULT_SUPER_ADMINS` and approve the request in the Trovo app (see
[INTEGRATION.md](INTEGRATION.md#2-the-trovo-app-for-logins)).

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
