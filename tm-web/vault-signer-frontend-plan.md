# Vault Signer Frontend — Plan

## Auth pattern this plan preserves (confirmed from real usage)

Backend endpoints gated by `AllowOrgOrTrovoAdminNormalized` already exist and have real frontend
callers to model against — the `/organizations/*` self-service-style endpoints
(`GET /organizations/members`, `GET /organizations/details`, `POST /organizations/members/link-wallet/*`).
Each is called from **two completely separate places**, one per portal:

| Endpoint | Dashboard (Trovo Admin) caller | Org portal caller |
|---|---|---|
| `GET /organizations/members` | `src/redux/api/organizations/api.ts` → `getOrganizationMembers`, injected into **`baseApi`** | `src/redux/api/org/api.ts` → `getOrganizationMembers`, injected into **`orgApi`** |
| `GET /organizations/details` | `organizations/api.ts` → `getOrganizationDetails`, `baseApi` | `org/api.ts` → `getOrganizationDetailsData`, `orgApi` |
| `POST /organizations/members/link-wallet/*` | *(admin never calls this)* | `sharedstakeholders/api.ts` → injected into `orgApi` |

This is forced by how auth is actually wired: `baseApi` (`src/redux/baseApi/index.ts`) uses
`axiosBaseQuery`, whose global axios interceptor reads the Trovo Admin JWT from
`getPreloadedState().auth.token` (localStorage key `TOKEN`) and sends
`Authorization: Bearer <token>`; `orgApi` (`src/redux/baseApi/orgApi.ts`) uses
`orgAxiosBaseQuery`, which reads the org member's token from localStorage key `org_auth_token`
(`TOKEN_ORG`) and sends it raw, unprefixed, as `Authorization: <token>`. These are two different
axios instances with two different interceptors — there is no shared query layer between the two
portals. Both interceptors already redirect on auth failure (`baseApi`'s 401 handler
refreshes-then-redirects to `/sign-in`; `orgApi`'s redirects straight to `/organizations/login` on
401/403) — this is the entire route-protection mechanism; neither `(dashboard)/layout.tsx` nor
`organisation/layout.tsx` does any auth check of its own, so the new pages don't need one either.

**Conclusion:** every `/me/vault-signer/*` self-service endpoint is injected **twice** — once into
`baseApi` for the dashboard, once into `orgApi` for the org portal — exactly like
`getOrganizationMembers` is today.

**Existing wallet-link UI reused, not reinvented:** the org portal already has a live "is my Trovo
wallet linked" indicator and CTA — `OrgNavBar.tsx` reads `wallet.is_wallet_linked` /
`wallet.trovo_wallet_username` off `useGetStakeholderProfileQuery()`
(`sharedstakeholders` → `orgApi`, `GET /stakeholder/shared/profile`) and, when unlinked, opens the
existing `LinkWalletModal` (`useRequestWalletLinkAuthorizationMutation` +
`useVerifyWalletLinkMutation`, both already on `orgApi`). This is precisely the condition that
gates Personal Secrets (`TrovoWalletUsername` unset → backend's `403 trovo_wallet_not_linked`).
The Personal Secrets page checks `wallet?.is_wallet_linked` the same way and, if false, shows the
same "Link Wallet" affordance instead of a dead-end error message.

---

## 0. Scope

Two audiences, matching the backend's own split:

- **Self-service** (`/api/v1/me/vault-signer/*`) — usable by *either* a Trovo Admin or an
  Organization member (`AllowOrgOrTrovoAdminNormalized`). UI in **both** portals, each wired to
  its own API instance per the pattern above.
- **Admin-only** (`/api/v1/admin/vault-signer/*`) — Trovo SuperAdmin only
  (`AuthenticateSuperAdmin`). Lives **only** in the `(dashboard)` app, on `baseApi`, under
  Settings.

No new design language — reuse the exact building blocks already in `src/components`
(`CustomTable`, `Modal`, `SearchBar`, `PrimaryButton`/`SecondaryButton`, `Tab`, `EmptyState`,
`showSuccessToast`/`showErrorToast`) and the styled-components look of `settings/mintingapprovers`
and `admin-users`.

---

## 1. Redux API layer — `vaultSigner` domain, dual-injected exactly like `organizations`/`org`

```
src/redux/api/vaultSigner/
├── interface.ts     # shared TS types for both portals — one source of truth
├── api.ts           # injectEndpoints(baseApi) — self-service (6) + all admin endpoints (7)
├── orgApi.ts        # injectEndpoints(orgApi)  — self-service (6) only, mirrors org/api.ts
└── index.ts          # re-exports both
```

New tag types added to `src/redux/baseApi/tagTypes.ts`:
```ts
VAULT_SIGNER_SECRETS: "vaultSignerSecrets",
VAULT_SIGNER_PERSONAL_ENVS: "vaultSignerPersonalEnvs",
VAULT_SIGNER_MANAGED_SECRETS: "vaultSignerManagedSecrets",
VAULT_SIGNER_ASSIGNMENTS: "vaultSignerAssignments",
VAULT_SIGNER_AUDIT_LOG: "vaultSignerAuditLog",
```

Types match the actual handler code (`signer_handler.go`, `personalenv_handler.go`,
`db_models.go`), not the generic `additionalProperties: true` swagger stubs.

`api.ts` (`baseApi`) endpoints:

| Hook | Method/URL | Tag |
|---|---|---|
| `useListMySecretsQuery` | `GET /me/vault-signer/secrets` | provides `VAULT_SIGNER_SECRETS` |
| `useGetMySecretValueQuery(secretId)` | `GET /me/vault-signer/secrets/:secretId` | — |
| `usePutMySecretValueMutation` | `PUT /me/vault-signer/secrets/:secretId` | invalidates `VAULT_SIGNER_SECRETS` |
| `useListPersonalEnvsQuery` | `GET /me/vault-signer/personal-envs` | provides `VAULT_SIGNER_PERSONAL_ENVS` |
| `useCreatePersonalEnvMutation` | `POST /me/vault-signer/personal-envs/:prefix` | invalidates `VAULT_SIGNER_PERSONAL_ENVS` |
| `useDeletePersonalEnvMutation` | `DELETE /me/vault-signer/personal-envs/:prefix` | invalidates `VAULT_SIGNER_PERSONAL_ENVS` |
| `useListManagedSecretsQuery` | `GET /admin/vault-signer/managed-secrets` | provides `VAULT_SIGNER_MANAGED_SECRETS` |
| `useRegisterManagedSecretMutation` | `POST /admin/vault-signer/managed-secrets` | invalidates `VAULT_SIGNER_MANAGED_SECRETS` |
| `useListAssignmentsQuery(id)` | `GET /admin/vault-signer/managed-secrets/:id/assignments` | provides `VAULT_SIGNER_ASSIGNMENTS` |
| `useCreateAssignmentMutation` | `POST .../assignments` | invalidates `VAULT_SIGNER_MANAGED_SECRETS`, `VAULT_SIGNER_ASSIGNMENTS` |
| `useEditAssignmentMutation` | `PATCH .../assignments/:assignmentId` | invalidates same |
| `useDeleteAssignmentMutation` | `DELETE .../assignments/:assignmentId` | invalidates same |
| `useListAuditLogQuery` | `GET /admin/vault-signer/audit-log` | provides `VAULT_SIGNER_AUDIT_LOG` |

`orgApi.ts` (`orgApi`) — the same six self-service endpoints, same URLs/shapes/tags, hook names
suffixed `Org` to avoid collisions when re-exported alongside `api.ts` from `index.ts`
(`useListMySecretsOrgQuery`, `usePutMySecretValueOrgMutation`, etc.) — the same disambiguation
`organizations`/`org` already needs for `getOrganizationMembers` today.

---

## 2. Self-service pages — one shared UI, two portal-specific wrappers

- Dashboard: new sidebar item in `SideBar.tsx`, **"Vault Signer"** → `/vault-signer`, using the
  `baseApi` hooks.
- Organisation: added to `sharedNavigation` in `organizationNavigation.ts` → `/organisation/vault-signer`,
  using the `orgApi` (`...Org`) hooks.

Shared UI lives once in `src/app/vault-signer/components/` (`MySignerSecrets`, `PersonalSecrets`),
written against a small hook-set prop so each portal's thin `page.tsx` passes in its own hooks.

`MySignerSecrets`: table of assigned slots + Signer Position (`—` when unclaimed — expected, not
an error) + "Set Value" modal → `PUT`. Error messages read via the same `getErrorMessage` helper
pattern already used in `LinkWalletModal.tsx` (`error?.data?.error ?? error?.data?.message ?? ...`).

`PersonalSecrets`: table of prefixes + `exists` pill + Create/Delete (never Edit — no update path
exists). Org portal gates on `wallet?.is_wallet_linked` and reuses `LinkWalletModal` instead of a
dead-end error panel; dashboard side never needs this gate (`AdminUser.Username` is required).

---

## 3. Admin UI — `(dashboard)/settings/vault-signer` (baseApi only)

Two tabs via the existing `Tab.tsx`:

- **Managed Secrets** — table + "Register Secret" modal; row click drills into
  `/settings/vault-signer/[id]` for that secret's **Assignments** (assign/reassign/remove owner,
  no position field ever accepted from the admin — it's claimed by the owner's own first write).
- **Audit Log** — read-only table, server-capped at 200 rows.

**Access control.** Every route under `/admin/vault-signer/*` — managed secrets, assignments,
and the audit log alike — is gated by `middleware.AuthenticateSuperAdmin(s.AdminDB)` on the
backend, confirmed against the middleware's actual implementation
(`internal/middleware/authentication_middleware.go`): it has no `OrganizationAuth` branch at
all, so an Organization member's token is rejected outright (`401`) before the SuperAdmin role
check even runs — this is not merely a routing convention, an org member is structurally
incapable of authenticating against it. That is why these seven endpoints are only ever injected
into `baseApi` (`src/redux/api/vaultSigner/api.ts`) and have no `orgApi.ts` counterpart, and why
the organisation portal has no admin UI at all — there would be nothing for it to call. On the
client, `ManagedSecretsTable`/`AssignmentsTable`/`AuditLogTable` check for the backend's own
`403` on their list queries and render `AccessDenied` instead of a broken table (see
"Corrections made during implementation" below for why this replaced the originally planned
client-side `Role` check).

---

## 4. File list

```
src/redux/api/vaultSigner/interface.ts
src/redux/api/vaultSigner/api.ts          # → baseApi
src/redux/api/vaultSigner/orgApi.ts       # → orgApi
src/redux/api/vaultSigner/index.ts

src/app/vault-signer/components/MySignerSecrets.tsx
src/app/vault-signer/components/PersonalSecrets.tsx
src/app/vault-signer/components/SetSignerValueModal.tsx
src/app/vault-signer/components/PersonalSecretModal.tsx

src/app/(dashboard)/vault-signer/page.tsx        # MySignerSecrets + PersonalSecrets, baseApi hooks
src/app/organisation/vault-signer/page.tsx        # same components, orgApi hooks + wallet-link gate

src/app/(dashboard)/settings/vault-signer/page.tsx
src/app/(dashboard)/settings/vault-signer/[id]/page.tsx
src/app/(dashboard)/settings/vault-signer/components/
  ├── ManagedSecretsTable.tsx
  ├── RegisterManagedSecretModal.tsx
  ├── AssignmentsTable.tsx
  ├── AssignOwnerModal.tsx
  ├── DeleteAssignmentConfirm.tsx
  └── AuditLogTable.tsx

# modified
src/redux/baseApi/tagTypes.ts                              # + 5 new tags
src/app/(dashboard)/components/SideBar.tsx                 # + "Vault Signer" item
src/app/organisation/components/organizationNavigation.ts  # + sharedNavigation item
```

---

## Implementation Summary

**Status: built, on branch `ric-vault-signer-frontend-feature`.** Everything in this plan has
been implemented directly against the actual backend contract (`vault-signer-service-plan.md`,
`internal/components/vaultsigner/handlers/*.go`, `models/db_models.go`) and the frontend
conventions confirmed by reading the codebase, not guessed from the swagger stubs.

### What was added

**Redux API layer** (`src/redux/api/vaultSigner/`):
- `interface.ts` — TS types matched to the handler code's actual JSON field names (e.g.
  `ownerLabel`, `vaultKey`, `friendlyLabel`, `deletedPosition`), not the generic
  `additionalProperties: true` swagger stubs.
- `api.ts` — all 13 endpoints (6 self-service + 7 admin) injected into `baseApi`, for the
  dashboard.
- `orgApi.ts` — the same 6 self-service endpoints injected into `orgApi` for the organisation
  portal, hook names suffixed `Org` to avoid collisions with `api.ts`'s exports when both are
  re-exported from `index.ts`. Payloads sent via `body:` (matching `org/api.ts` and
  `sharedstakeholders/api.ts`'s own convention on `orgApi`) rather than `data:`.
- `index.ts` re-exports both plus the shared interfaces.
- `src/redux/baseApi/tagTypes.ts` — 5 new tags (`VAULT_SIGNER_SECRETS`,
  `VAULT_SIGNER_PERSONAL_ENVS`, `VAULT_SIGNER_MANAGED_SECRETS`, `VAULT_SIGNER_ASSIGNMENTS`,
  `VAULT_SIGNER_AUDIT_LOG`).

**Shared self-service UI** (`src/app/vault-signer/components/`):
- `MySignerSecrets.tsx` — table of assigned slots + Signer Position (`Unclaimed` when `null`,
  not treated as an error state) + a "Set Value" modal (`SetSignerValueModal.tsx`) that calls
  `PUT`. Surfaces on-chain swap confirmation (`stellarTxHash`) in the success toast.
- `PersonalSecrets.tsx` — prefix list with an `exists` pill and a single contextual
  Create/Delete action (no Edit — the backend has no update-in-place path). Accepts an optional
  `walletGate` prop; when supplied and `isLinked` is false, renders a "Link Wallet" CTA instead
  of the list.
- `PersonalSecretModal.tsx` — create-only value form / delete confirmation, two modes on one
  component.
- Both components are typed against small structural hook signatures
  (`UseListMySecrets`, `UsePutMySecretValue`, etc.) rather than `typeof useXQuery` — this was a
  deliberate fix made during implementation (see "Corrections" below): the dashboard's `baseApi`
  hooks and the org portal's `orgApi` hooks are bound to two different `BaseQueryFn` generics
  and are not directly assignable to one another under `typeof`, even though both satisfy the
  same calling shape at runtime.
- `getErrorMessage` (exported from `MySignerSecrets.tsx`) reuses the exact precedence
  `LinkWalletModal.tsx` already established (`error?.data?.error ?? error?.data?.message ??
  error?.error ?? fallback`) — reused by every mutation across both the self-service and admin
  UI instead of a new convention per component.

**Dashboard self-service page:** `src/app/(dashboard)/vault-signer/page.tsx` — wires
`MySignerSecrets`/`PersonalSecrets` to the `baseApi` hooks, no wallet gate (a Trovo Admin's
`AdminUser.Username` is always present). Added to `SideBar.tsx`'s `sidebarItems` (key-square
icon, `/vault-signer`).

**Organisation self-service page:** `src/app/organisation/vault-signer/page.tsx` — wires the
same components to the `orgApi` (`...Org`) hooks, and gates `PersonalSecrets` on
`wallet?.is_wallet_linked` from the already-existing `useGetStakeholderProfileQuery()`
(`sharedstakeholders` → `orgApi`), reusing the existing `LinkWalletModal` component from
`organisation/components/` (the same one `OrgNavBar.tsx` already uses) rather than inventing a
new "connect your wallet" flow. Added to `organizationNavigation.ts`'s `sharedNavigation` (same
key-square icon, `/organisation/vault-signer`) so every stakeholder role gets it.

**Admin UI** (`src/app/(dashboard)/settings/vault-signer/`):
- `page.tsx` — two tabs (`Tab.tsx`) for Managed Secrets and Audit Log.
- `[id]/page.tsx` — one managed secret's assignments. Reads the secret's own label/wallet key
  out of the already-fetched `useListManagedSecretsQuery()` cache rather than inventing a
  singular-fetch endpoint the backend doesn't expose (there is no
  `GET /admin/vault-signer/managed-secrets/:id`).
- `components/ManagedSecretsTable.tsx` + `RegisterManagedSecretModal.tsx` — list + register
  form, no `position`/CSV-seeding fields (the backend expects the CSV to already hold real
  signer keys before a secret is registered).
- `components/AssignmentsTable.tsx` + `AssignOwnerModal.tsx` (assign/reassign, one component
  two modes) + `DeleteAssignmentConfirm.tsx` — no position field ever collected in the UI at
  creation or reassignment time (it's claimed by the owner's own first `PUT`, never by an
  admin); the delete confirmation copy explicitly warns about the possible on-chain removal and
  the result panel surfaces `deletedPosition`/`onChainRemoval` from the response.
- `components/AuditLogTable.tsx` — read-only, no pagination controls (the backend caps this at
  200 rows server-side with no pagination params).
- Added a "Vault Signer Settings" menu group to `(dashboard)/settings/layout.tsx`'s existing
  settings sidebar, linking to `/settings/vault-signer`.

### Corrections made during implementation (not just planning)

1. **Shared-hook-prop typing had to be loosened.** The original plan's approach of typing
   `MySignerSecrets`/`PersonalSecrets` props as `typeof useListMySecretsQuery` etc. (the exact
   `baseApi`-bound hook types) failed `tsc` when the same components were wired to the
   `orgApi`-bound `...Org` hooks on the organisation page — RTK Query hook types carry their
   `BaseQueryFn` generic, and `baseApi`'s and `orgApi`'s are different concrete types, so one
   hook's type isn't assignable to a prop typed as the other's. Fixed by declaring small,
   structural function-signature types instead (what the component actually calls: `data`,
   `isLoading`, `error?`, `refetch` for queries; a tuple of `[trigger, state]` for mutations),
   which both hook families satisfy at the call-site shape that matters.
2. **Mutation hook tuples are `readonly`.** The first pass at those structural mutation types
   used a mutable tuple (`[fn, state]`); RTK Query's `useMutation` actually returns a `readonly`
   tuple, which isn't assignable to a mutable tuple type. Fixed by declaring the prop types as
   `readonly [fn, state]`.
3. **No client-side "is this admin a SuperAdmin" field actually exists.** The original plan
   said to gate the Settings entry point "the same lightweight way `AdminFilter.tsx` already
   reads the stored admin's `Role`" — on inspection, `AdminFilter.tsx` doesn't read the current
   user's own role at all; it's a filter dropdown of role *options* for filtering the list of
   *other* admins. The actually-stored current-user object (`state.auth.user`, populated from
   the wallet-connect login response, `IWalletConnectResponse`) only carries `adminLevel: 0 | 1`,
   not the `"SUPER_ADMIN"` / `"EDIT_LEVEL_ADMIN"` / `"VIEW_ONLY_ADMIN"` `Role` string other
   admins are listed with — there was no reliable field to gate on. Rather than fabricate a
   check against a field that might not correspond to the backend's actual `AuthenticateSuperAdmin`
   decision, `ManagedSecretsTable`/`AuditLogTable`/`AssignmentsTable` instead check for the
   backend's own `403` on their list queries and render a plain `AccessDenied` panel — the real
   gate stays server-side, and the UI explains the denial instead of showing a broken empty
   table.

### Verification performed

- `yarn install` — clean (one optional native-module build failure for `canvas`, an unrelated
  pre-existing optional transitive dependency of `pdfjs-dist`; explicitly non-fatal, "Done").
- `npx tsc --noEmit` — compared against a true clean baseline (`git stash -u`): baseline 174
  pre-existing errors (all in files this change never touches, e.g. the repo-wide `FC<SvgProps>`
  vs `StaticImageData` mismatch on `next/image` across dozens of `.svg` icon imports), branch
  175 — the one delta is the same pre-existing icon-typing pattern firing once more for the new
  `key-square.svg` import, not a new category of error. Zero errors anywhere under
  `vault-signer`/`vaultSigner` paths.
- `yarn lint` (`next lint`) — clean; all reported warnings are pre-existing and outside every
  file this change touches.
- `yarn build` (`next build`) — succeeds. All four new routes compile and appear in the route
  manifest: `/vault-signer`, `/organisation/vault-signer`, `/settings/vault-signer`,
  `/settings/vault-signer/[id]`.

### Verification **not** performed, and why

- **No live backend.** Nothing here has been exercised against a running
  `admin-panel-dashboard` instance, Vault, or Horizon — request/response shapes are written to
  match the handler code read directly (Section-by-section against
  `signer_handler.go`/`personalenv_handler.go`), not guessed, but that's not a substitute for an
  actual round trip.
- **No visual QA.** No dev server was run in a browser; layout/spacing follows the existing
  styled-components conventions from `mintingapprovers`/`admin-users` by inspection, not by
  screenshot comparison.

### Recommended before merging

1. Run this against a real backend (or a mocked one matching the handler shapes above) and
   exercise all three flows end-to-end: signer slot claim-and-swap, personal secret
   create/delete, and the assignment delete/on-chain-removal path.
2. Confirm the `403`-based `AccessDenied` gating is an acceptable interim approach, or wire in
   whatever real "is this admin a SuperAdmin" signal the team decides on (see Correction 3
   above) if a friendlier pre-emptive hide is wanted.
3. Visual QA in a browser against both portals, including the wallet-link gate on the
   organisation portal's Personal Secrets panel.

---

## Reflected backend change: one assignment per owner per managed secret

The backend added a constraint (backend plan Section 12) that `(ManagedSecretID, OwnerRefID)` is
unique — an owner can hold at most one assignment on a given managed secret, guarded at both the
DB and app level, surfaced as `409 owner_already_assigned` on both
`POST .../assignments` (create) and `PATCH .../assignments/:assignmentId` (reassign).

**No request/response shape changed**, so no `interface.ts` types needed updating — this is
purely a new error case on two endpoints already wired up. Two things did need fixing on the
frontend side to actually surface it correctly:

1. **`getErrorMessage`'s precedence was backwards.** The helper (exported from
   `MySignerSecrets.tsx`, reused by every mutation in this feature) read
   `error?.data?.error ?? error?.data?.message ?? ...` — `error` (the machine-readable slug)
   was checked *before* `message` (the sentence meant for a person). For the vault-signer
   backend's structured error responses, which carry both (e.g. the new
   `{"error": "owner_already_assigned", "message": "This owner already has an assignment on
   this managed secret."}`), this meant the toast showed the raw slug, not the sentence — a
   latent bug that predates this backend change but was only surfaced by tracing through what
   the new 409 would actually render. Fixed by flipping the precedence to
   `error?.data?.message ?? error?.data?.error ?? error?.error ?? fallback`, which now also
   fixes the same problem for the pre-existing `no_position_available` 409
   (`PUT /me/vault-signer/secrets/:secretId`) and both `502` assignment-delete responses — none
   of those were showing their friendly `message` correctly either, for the same reason.
2. **`AssignOwnerModal.tsx`** already routed every create/reassign error through
   `getErrorMessage`, so once the precedence was fixed, the new `409` needed no other code
   change — added a comment explaining why this form has no client-side pre-check for the
   conflict: the admin types a username/email, and only the backend can resolve that to the
   stable `OwnerRefID` the constraint actually keys on (`services.ResolveOwnerRefID` has no
   frontend equivalent), so the form always submits and lets the server be the judge.

**Verification:** `tsc --noEmit` and `eslint`, scoped to every vault-signer file, both clean
after the change (no new errors, no new warnings).
