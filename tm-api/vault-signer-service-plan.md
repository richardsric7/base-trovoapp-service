# Vault Signer-Key & Personal-Env Manager — Plan (backend-only, `admin-panel-dashboard`)

## 0. Scope

This is a Go-only backend feature built directly inside this existing repo (`admin-panel-dashboard`, the Gin/GORM backend at `trovo-admin-monorepo/backend`), as one more component alongside `organizations`, `stakeholder`, `careers`, etc. No separate service, no separate `go.mod`. Everything below is a Go-only addition: new routes under the existing `apiV1` group, new GORM models added to the existing AutoMigrate list, and a new Vault/Stellar integration reusing what this codebase already has for Horizon. There is no client of any kind in scope here — this document ends at the HTTP API contract.

There are two secret models:

**A. Shared signer-slot secrets** — a Vault secret holds a CSV of Stellar signer keys (typically 3). **Every entry in that CSV is assumed to already exist as a real, registered signer on the wallet's public key from the start** — this service never writes a placeholder or empty entry into it (Section 4/5e). An owner's slot has a 1-based **Signer Position** — never a raw 0-based index — but assigning an owner to a managed secret does *not* hand them a position: `Position` starts `NULL` and is only claimed, automatically, the first time that owner submits their own real signer value (Section 5e). Reassigning an assignment to a different owner (Section 6) always retains whatever `Position` currently holds, null or claimed — the new owner picks it up the normal way, by logging in later and submitting their own value, which is what actually triggers the swap/rotation. An owner can edit only their own position; the service rebuilds the CSV and writes it back, but only after the submitted key passes format validation *and* is confirmed to be an activated account on the ledger (Section 5c).

**Stellar protocol note:** the network currently runs **Protocol 17**, with forward compatibility required through **Protocol 19**. This repo already pins `github.com/stellar/go` (see `go.mod`) and uses it in `internal/network`, `internal/trovosdk`, and elsewhere — no new SDK dependency here, just confirm the pinned version's protocol support before this feature ships.

**B. Personal per-owner env secrets — called "personal secrets" in any user-facing text.** `PERSONAL_ENVS`, the env var name, and internal identifiers (`personalenv_handler.go`, `personal_env_create`, etc.) keep the "personal env" name since that's just internal plumbing — but every error message, friendly copy, or Swagger summary a caller might see refers to this feature as **"personal secrets"**, never "personal env(s)."

A configurable list of env-name *prefixes*, each paired with a friendly label, is injected via `PERSONAL_ENVS` — a CSV of `PREFIX:FriendlyLabel` entries (e.g. `AUTO_APPROVE:Auto Approval,MARKET_MAKER:Market Maker`). Every owner automatically "owns" one secret per prefix. Unlike signer slots, there's no admin assignment step, and unlike a typical secrets UI, **this is create/delete only — there is no read-back and no update-in-place**: an owner can list their prefixes and see whether each `exists`, **create** the secret (which writes the value into Vault), and **delete** it (which removes the Vault secret entirely). There is no endpoint that ever returns a personal-env value to any client, including the owner who set it. To change a value, the owner deletes and re-creates it — there is no `PUT`/update path. There is no database table for the prefix/label list either — both live entirely in the `PERSONAL_ENVS` env var, parsed once at startup into an in-memory map.

**Personal secrets belong to an individual, never to an organization.** An organization member's personal secrets are keyed by their `TrovoWalletUsername` (on `OrganizationMember`) — the org itself has no notion of "its" personal secrets. If an org member hasn't linked a Trovo wallet (`TrovoWalletUsername` is unset), they cannot use personal secrets at all, and are shown a friendly message to connect their profile to the Trovo app instead of a raw error. See Section 5b for the exact resolution/gating logic.

Both flows sit behind this repo's real auth middleware — see Section 1.

---

## 1. Auth: Reuse `AllowOrgOrTrovoAdminNormalized`, Not a New JWT Stub

The original plan assumed a single stubbed `ValidateJWT()`. This repo already has three real auth populations (Trovo Admin, Organization member, Stakeholder portal), each with its own middleware in `internal/middleware`. Per decision: **this feature is shared between Trovo Admins and Organization members**, using the existing composed middleware:

```go
middleware.AllowOrgOrTrovoAdminNormalized(s.AdminDB)
```

This already:
- Accepts either a Trovo Admin JWT (verified via `trovosdk.ServiceLink.JwtTokenVerify`, checked against the `AdminUser` table) or an Organization member JWT (`JWT_SECRET`-signed, checked via `validateOrganizationSession`).
- Normalizes context regardless of which one authenticated:
  - `c.GetString("user_type")` → `"trovo_admin"` or `"org_member"`
  - `c.GetString("current_user_id")` → `AdminUser.ID` (as string) or `OrganizationMember.ID`
  - `c.GetString("current_user_email")` → admin email or member email
  - `c.GetString("user_organization_role")` → org role, empty for admins

**Important departure from the original plan:** there is no single `username` field shared by both populations — `OrganizationMember` has no username. Two different fields end up mattering for two different flows, but (per Section 4's optimized `VaultSignerAssignment` design below) only **one** of them is ever *stored*:

| Flow | `trovo_admin` field | `org_member` field | Where it's used |
|---|---|---|---|
| **Managed secrets** (signer slots, `VaultSignerAssignment`) | `AdminUser.Username` | `OrganizationMember.Email` | Only as a **lookup key at assignment-creation time** (Section 4) — an admin identifies who to assign by username/email, since neither admin should have to know an opaque ID. Never stored, never compared at request time. |
| **Personal secrets** (Section 5b) | `AdminUser.Username` | `OrganizationMember.TrovoWalletUsername` | Used directly to derive the Vault key on every create/delete — this one genuinely has to be resolved on every request, since it *is* the key. |

For managed secrets specifically: **an org member never needs a linked Trovo wallet to be assigned a signer slot or manage it** — that requirement is scoped exclusively to personal secrets.

**What's actually stored and compared for managed secrets is `current_user_id` — the real, stable primary key** (`AdminUser.ID` or `OrganizationMember.ID`), not a username or email. See Section 4 for why, and for how the API still lets an admin assign by email/username without ever exposing the underlying ID as something either population's identity check depends on. This means ownership checks on `me/vault-signer/secrets*` (is this caller allowed to touch this slot?) are a direct field comparison against context — `current_user_id` and `user_type` are already provided by `AllowOrgOrTrovoAdminNormalized` with **zero additional DB lookups**, unlike the personal-secrets gate in Section 5b which does need one.

Admin-only actions (registering/listing a managed secret, creating/reassigning/removing an assignment, and reading the audit log) are **not** appropriate for any org member, so they're gated more strictly, using the existing role check:

```go
middleware.AuthenticateSuperAdmin(s.AdminDB) // Trovo SuperAdmin only
```

**Confirmed, re-verified against the actual middleware implementation (`internal/middleware/authentication_middleware.go`):** `AuthenticateSuperAdmin` has no `OrganizationAuth` branch at all — unlike `AllowOrgOrTrovoAdmin`/`AllowOrgOrTrovoAdminNormalized`, it only ever verifies a Trovo Admin JWT (via `trovosdk.ServiceLink.JwtTokenVerify` against the `AdminUser` table) and then further requires `existingAdmin.Role == SuperAdmin`. An Organization member's JWT is structurally incapable of satisfying this middleware — there is no code path here that inspects org auth at all — so this isn't just a routing convention, it's enforced at the type/data level: an org member is rejected with `401` before the SuperAdmin role check is even reached, and a Trovo Admin who isn't a SuperAdmin is rejected with `403` after it. This is the same middleware `internal/components/vaultsigner/controllers/main.go` already wires onto every route under the `admin/vault-signer` group (Section 6) — `managed-secrets*`, `managed-secrets/:id/assignments*`, and `audit-log` alike, with no exceptions and no separate, looser gate for any of the three. The self-service group (`me/vault-signer/*`) is the only part of this feature reachable by both populations, via `AllowOrgOrTrovoAdminNormalized` as described above.

---

## 2. High-Level Architecture (in-repo)

```
┌─────────────┐      ┌───────────────────────────────────┐      ┌─────────────┐
│  Any client │ ───► │  admin-panel-dashboard (this repo) │ ───► │   Vault     │
│  (existing  │ ◄─── │  internal/components/vaultsigner   │ ◄─── │  (KV v2)    │
│  admin app) │      └──────────────┬──────────────────────┘      └─────────────┘
└─────────────┘                     │
                                     ├──► s.AdminDB (existing GORM Postgres connection)
                                     │      vault_signer_managed_secrets
                                     │      vault_signer_assignments
                                     │      vault_signer_audit_logs
                                     │
                                     └──► internal/network (existing Horizon/Stellar helpers)
                                            network.GetBlockchainClient()
                                            network.GetBlockchainNetworkPassPhrase()
```

No new database, no new connection string. This reuses `s.AdminDB` exactly like `stakeholder` and `organizations` do.

**New config needed (Vault only — Stellar config already exists):**
```
VAULT_ADDR
VAULT_TOKEN
PERSONAL_ENVS = "AUTO_APPROVE:Auto Approval,MARKET_MAKER:Market Maker"
PERSONAL_ENV_VAULT_MOUNT              # default: secret
PERSONAL_ENV_VAULT_PATH_PREFIX        # default: personal-envs
```
Loaded via `os.Getenv(...)`, matching how `internal/network` already reads `EXPANSION_URL` and `BLOCKCHAIN_NETWORK_PASSPHRASE` — not through `config/config.go`'s `envconfig` struct, which is scoped to DB/port/env fields only. Add these four to `.env-sample` alongside the existing entries.

**`PERSONAL_ENVS` format:** each comma-separated entry is `PREFIX:FriendlyLabel`, not a bare prefix. Parsed once at startup into an in-memory `map[string]string` (prefix → label) — there is no `VaultSignerPersonalEnvDefinition` table or any other DB-backed source for the label; the env var is the single source of truth for both the prefix and its display name. Fail-fast at startup if any entry doesn't contain the `:` separator, rather than silently falling back to an auto-derived label — consistent with this feature's "fail closed" posture elsewhere (Section 5c, 9).

**Horizon/Stellar config — reuse, don't duplicate:** the original plan introduced its own `HORIZON_URL`. This repo already has that concept:
- `network.GetBlockchainClient()` → `*horizonclient.Client`, pointed at `os.Getenv("EXPANSION_URL")`
- `network.GetBlockchainNetworkPassPhrase()` → `os.Getenv("BLOCKCHAIN_NETWORK_PASSPHRASE")`

Both are already used by `internal/models/streams.go`, `internal/trovosdk`, and payment flows. The vault-signer feature calls these exact functions rather than standing up a second, possibly-divergent Horizon client pointed at a different network.

`github.com/hashicorp/vault/api` is **not** currently a dependency (`go.mod`/`go.sum` checked — absent). It needs to be added via `go get github.com/hashicorp/vault/api && go mod tidy`.

---

## 3. Project Structure — Mirrors `stakeholder`, Not a Standalone Service

The original Section 8/11 laid out `/cmd/server`, its own `go.mod`, its own `README.md`. None of that applies now. New code goes under one new component, following the same five-folder shape `stakeholder` and `organizations` already use:

```
internal/components/vaultsigner/
├── controllers/
│   └── main.go            # Init(router *gin.Engine, s *serverModels.Server) — registers routes
├── handlers/
│   ├── signer_handler.go       # me/secrets*, admin/managed-secrets*
│   └── personalenv_handler.go  # me/personal-envs* (list, create, delete — no get/update)
├── services/
│   ├── validation_service.go   # trim → keypair.ParseFull() → Horizon activation check (Section 5c)
│   ├── swap_service.go         # SwapSigner(): old/new keypair gate + build/sign/submit (Section 5d)
│   ├── signer_service.go       # csvIndex(position) conversion (Section 4); assignment create (auto-position), edit (reassign owner), delete (renumber + CSV collapse + conditional on-chain removal, Section 5e); resolves ownerIdentifier -> OwnerRefID; joins OwnerRefID -> display label for listings
│   └── personalenv_service.go  # ResolveTrovoUsername (identity gate, Section 5b), key derivation, create/delete
├── models/
│   └── db_models.go        # VaultSignerManagedSecret, VaultSignerAssignment, VaultSignerAuditLog
└── vaultclient/
    ├── signer.go            # CSV rebuild + CAS logic against Vault KV v2 (0-based index only); splice/collapse for delete (Section 5e)
    └── personalenv.go       # existence check (metadata-only) + create/update
```

Wired into `main.go` exactly like the existing components:
```go
vaultsigner "admin-panel-dashboard/internal/components/vaultsigner/controllers"
...
vaultsigner.Init(router, s)
log.Println("##vaultsigner services initialized##")
```
placed after `stakeholder.Init(router, s)`.

---

## 4. Data Model — GORM-Owned, Not a Raw SQL Migration

This repo's AdminDB schema is **GORM-owned**: `internal/db/main.go`'s `migrateAdminSchemaTransaction` calls `gormDB.AutoMigrate(...)` with the full list of admin/org/stakeholder models (this was the point of the recent "make stakeholder schema GORM-owned" fix on `main`). The flat files under `/migrations/*.sql` are legacy incremental patches from before that change, not the current path for adding new tables.

**Directive: every new model in this feature MUST be registered in `migrateAdminSchemaTransaction`'s `AutoMigrate(...)` call, in the same commit that introduces the model.** A model that only exists as a Go struct and is never added to that call will never get a table — there is no separate migration step, no CLI command, and no other place in this codebase that creates AdminDB tables. All three new tables are plain GORM model structs added to the same `AutoMigrate(...)` call, the same way `stakeholderModels.StakeholderAssetAssignment` etc. are today:

```go
// internal/db/main.go, migrateAdminSchemaTransaction
vaultSignerModels "admin-panel-dashboard/internal/components/vaultsigner/models"

if err := gormDB.AutoMigrate(&models.AdminUser{}, /* ...existing... */
    &vaultSignerModels.VaultSignerManagedSecret{},
    &vaultSignerModels.VaultSignerAssignment{},
    &vaultSignerModels.VaultSignerAuditLog{},
); err != nil {
    return err
}
```

No new raw SQL file needed unless a future column needs the kind of guarded backfill `prepareAdminSchemaForAutoMigrate` does for `FundReleaseRequest` — not expected here since these are brand-new tables with no legacy data.

### `VaultSignerManagedSecret` (shared signer-slot secrets)

| field | type | notes |
|---|---|---|
| ID | uuid pk | |
| Label | text | e.g. "Mainnet Payment Signers" |
| VaultMount | text | e.g. `secret` (KV v2 mount) |
| VaultPath | text | e.g. `stellar/mainnet-signers` |
| VaultField | text | key inside the KV data map, e.g. `signers_csv` |
| WalletPublicKey | text | the Stellar `G...` account these signers control on-chain |
| ActiveSigningCount | int | entries actively used to sign at once (e.g. `3`) |
| CreatedAt | timestamp | |

**No `CSVCapacity` column.** Capacity isn't tracked in the DB — the CSV's actual length in Vault *is* the capacity, and it now grows and shrinks in lockstep with assignments (Section 5e), rather than independently. `Position` validity is checked against `len(strings.Split(csvValue, ","))` at write time (Section 5a), converted via `csvIndex(position)`, not against a stored capacity value. The only fixed rule is a **minimum of 4 entries**, enforced wherever a managed secret is registered/edited (admin endpoint, Section 6) and re-checked defensively in the CSV rebuild path before any write — reject with a validation error if the live CSV has fewer than 4 entries. This same floor also **blocks a `DELETE`** that would bring the assignment count below 4 (Section 5e) — the rule always meant "this managed secret must have at least 4 CSV entries," and now that assignment count and CSV length are the same number, deleting past that floor is exactly the case it was meant to catch.

### `VaultSignerAssignment` — optimized: store the stable ID, resolve the friendly identifier only at write time

| field | type | notes |
|---|---|---|
| ID | uuid pk | |
| ManagedSecretID | uuid fk | |
| Position | \*int, **nullable** | **1-based**, user-facing — "Signer Position," never "index," anywhere an admin sees it (Section 6). **Starts `NULL`** — a correction from an earlier draft that auto-assigned it serially at creation time. It is claimed exactly once, automatically, the first time the owner submits their own real signer value (`services.ClaimPosition`, Section 5e), and is never touched again except by that same claim happening for a *different* owner if reassignment leaves it unclaimed. `NULL` means "no one — this owner or a previous one — has ever set a real key for this assignment yet," not "invalid" or "deleted." |
| OwnerRefID | text | The real, stable primary key of the owner: `AdminUser.ID` (stringified) for `trovo_admin`, `OrganizationMember.ID` for `org_member`. **Not** `Username`/`Email` — see rationale below. |
| OwnerType | text | `trovo_admin` \| `org_member` — this repo has no cross-table foreign key, so `(OwnerType, OwnerRefID)` together is the effective polymorphic reference, the same idiom GORM's own polymorphic-association support uses |
| AssignedAt | timestamp | |

Composite index on `(OwnerType, OwnerRefID)` for "list my managed secrets" lookups. Unique constraint on `(ManagedSecretID, Position)` — Postgres treats every `NULL` in a unique index as distinct from every other `NULL`, so any number of not-yet-claimed assignments coexist on the same managed secret without conflict; the constraint only ever actually engages once a position is claimed.

**`(ManagedSecretID, OwnerRefID)` is also unique** (`idx_vault_signer_assignment_secret_owner`) — the same owner can never hold two assignment rows on the same managed secret. This is guarded at **both** levels, not just documented as a convention:
- **DB level:** a real unique index on the table, so it holds under concurrent writes, not just whenever the app-level check happens to run first.
- **App level:** `services.CreateAssignment` pre-checks for an existing `(ManagedSecretID, OwnerRefID)` row before inserting, and `services.EditAssignment` does the same before reassigning to a new owner (excluding the assignment's own row, so keeping the same owner is always a no-op save). Both also catch the underlying Postgres unique-violation (`SQLSTATE 23505`, via `services.isUniqueViolation`) on the rare concurrent race that slips past the pre-check, translating it to the same `services.ErrOwnerAlreadyAssigned` the handlers turn into a `409 owner_already_assigned` — a raw DB constraint error never reaches the caller.

Without this, `handlers.loadOwnedSecret` — which resolves the caller's own assignment via `Where("managed_secret_id = ? AND owner_ref_id = ? AND owner_type = ?", ...).First(...)` for every `GET`/`PUT /me/vault-signer/secrets/:secretId` — would have no reliable way to pick between two rows for the same owner on the same secret; `First()` with no `ORDER BY` picks one arbitrarily, silently making the other row's slot unreachable through self-service even if it had already claimed a live position. This constraint is what makes that lookup well-defined.

This also keeps `Position` well-defined per `(ManagedSecretID, OwnerRefID)`: since an owner can now only ever hold one assignment on a given managed secret, that pair maps to at most one `Position`. Position numbering itself doesn't change — `ClaimPosition` (Section 5e) already scopes its lowest-unclaimed-position search to one `ManagedSecretID` at a time, so positions run **1, 2, 3, ... independently per managed secret**, never globally across secrets; this constraint doesn't touch that numbering, it just guarantees no single owner can occupy more than one of those numbers on the same secret.

**`Position` vs. the Vault CSV index — one deliberate translation point, not two parallel numbering schemes:** everything in the DB, the API, and any user-facing text uses 1-based `Position` ("Signer Position 1", "Signer Position 2", ...). The Vault CSV itself is a plain 0-based array, because that's what `strings.Split` naturally produces — there's no reasonable way to make a Go slice 1-based, and no reason to fight that. So there is exactly **one** conversion point, in `services/signer_service.go`, and every Vault/on-chain touchpoint (Sections 5a, 5d, 5e) goes through it rather than re-deriving the conversion locally:
```go
// services/signer_service.go
func csvIndex(position int) int { return position - 1 }
```
Handlers and the DB layer never see or reason about a 0-based value; `vaultclient/` and the CSV-splicing logic in Section 5e never see or reason about a 1-based one. If a bug ever mixes the two up, the visible symptom would be an off-by-one on which owner's key gets read/written/removed — worth a unit test that specifically asserts `csvIndex(1) == 0`.

**Why `OwnerRefID` instead of storing `Username`/`Email` directly (a correction from the previous draft):** an earlier version of this plan stored the resolved `Username`/`Email` string itself as the row's identity. That has a real correctness problem — `OrganizationMember.Email` is a mutable field (an org member's email can change), so a stored copy of it silently goes stale: the assignment row would keep pointing at an email nobody owns anymore, and every future "is this my slot?" check (which compares against the *current* `current_user_email`) would start failing for that member with no error, no audit trail, nothing — the slot just quietly stops being theirs. `AdminUser.Username` is less likely to change but has the identical failure mode in principle. A primary key doesn't have this problem: `AdminUser.ID` and `OrganizationMember.ID` are the one thing on each row guaranteed not to change for the life of the record.

**This doesn't give up the ergonomics the email/username design was for** — an admin creating an assignment still never needs to know or paste an opaque ID (Section 6): the request still accepts a human identifier (`ownerIdentifier`: username for `trovo_admin`, email for `org_member`), and it's resolved to the real `OwnerRefID` exactly once, at creation time, by `services/signer_service.go` (Section 1/8). From that point on:
- **Authorization** (does this request own this slot?) is a direct comparison against `current_user_id`/`user_type` from `AllowOrgOrTrovoAdminNormalized` — **zero DB lookups**, and immune to the member later changing their email.
- **Display** (an admin viewing who owns each slot, Section 6) does a fresh join/lookup from `OwnerRefID` back to `AdminUser.Username`/`OrganizationMember.Email` at read time — so the label shown is always current, never a frozen snapshot.

**Created via `POST`, edited via `PATCH`, removed via `DELETE`, all at `/api/v1/admin/vault-signer/managed-secrets/:id/assignments[/:assignmentId]`** (Section 6) — these are the only endpoints that write to this table; there is no other path (no seed script, no admin UI form outside this API) that touches a `VaultSignerAssignment` row. `PATCH`/`DELETE` address an assignment by its own ID, not by `Position` — `Position` can be null, so it can't reliably identify a row on its own. `DELETE` is the complex one — see Section 5e for the full transactional flow (position renumbering, CSV collapse, and a conditional on-chain removal) — and `Position` is also claimed from a completely different endpoint entirely: `PUT /me/vault-signer/secrets/:secretId` (Section 5e/6), the owner's own first write.

**No `VaultSignerPersonalEnvDefinition` table.** Personal-env prefixes and their friendly labels come entirely from parsing `PERSONAL_ENVS` (`PREFIX:FriendlyLabel` CSV, Section 2) into an in-memory map at startup — there is nothing to persist here, so this model from the original draft is dropped.

### `VaultSignerAuditLog` (record only — both flows)

| field | type | notes |
|---|---|---|
| ID | uuid pk | |
| Kind | text | `signer_slot`, `signer_swap`, `assignment_create`, `assignment_edit`, `assignment_delete`, `personal_env_create`, or `personal_env_delete`. The three `assignment_*` kinds (Section 5e) are administrative lifecycle events on `VaultSignerAssignment` itself — distinct from `signer_slot`, which is about an owner writing their *value* into an already-assigned slot. |
| ManagedSecretID | uuid, nullable | `signer_slot`/`signer_swap` only |
| Position | int, nullable | `signer_slot`/`signer_swap`/`assignment_*` only — 1-based, same convention as `VaultSignerAssignment.Position` |
| VaultKey | text, nullable | full personal secret name, `personal_env_create`/`personal_env_delete` only |
| ActorID | text | For `signer_slot`/`signer_swap`: `current_user_id` (same real-ID convention as `VaultSignerAssignment.OwnerRefID`, Section 4 — not `Username`/`Email`, for the identical staleness reasons). For `personal_env_create`/`personal_env_delete`: the resolved Trovo username (`AdminUser.Username` or `OrganizationMember.TrovoWalletUsername`, Section 5b) — this one *is* a point-in-time display value, which is fine here since it's the same string already embedded in `VaultKey`. |
| ActorType | text | `trovo_admin` \| `org_member` |
| VaultVersionBefore | int, nullable | null if this was a create; on `personal_env_delete`, the version that existed immediately before the purge |
| VaultVersionAfter | int, nullable | **now nullable** — a `personal_env_delete` purges the secret entirely, so there is no "after" version to record |
| StellarTxHash | text, nullable | `signer_swap` only |
| StellarTxStatus | text, nullable | `success` or `failed`, `signer_swap` only |
| ChangedAt | timestamp | |
| IPAddress | text, optional | |

This table is a historical record only, written after the fact — nothing polls it or reads it back to make decisions. Note the field-nullability change from the earlier draft: `VaultVersionAfter` was previously non-nullable, which only worked while every write flow produced a new version. `personal_env_delete` doesn't, so it's now `*int` like `VaultVersionBefore`.

---

## 5. Vault & Stellar Interaction

`github.com/hashicorp/vault/api` (to be added), implemented in `internal/components/vaultsigner/vaultclient/`.

### 5a. Signer slots — unchanged mechanics, relocated

`vaultclient/signer.go` does the KVv2 get/rebuild/CAS-write for signer slots. Same Go snippets as the original plan's 4a — no repo-specific change needed beyond the package path.

Since capacity isn't a stored column (Section 4), `signer.go`'s rebuild step does two live checks against the CSV pulled from Vault on every write: the converted `index := csvIndex(position)` must be `< len(parts)`, and `len(parts) >= 4` — reject with a validation error otherwise rather than trusting a number that could drift from what's actually in Vault. `signer.go` itself only ever deals in `index` — the 1-based `Position` never crosses into this file; the conversion happens once in `signer_service.go` before calling in (Section 4).

`signer.go` also owns the CSV-shape operation the assignment-delete lifecycle (Section 5e) needs — a plain CAS write with no validation pipeline involved (there's no signer *value* being written, just the CSV's shape). **There is deliberately no "grow the CSV" function anymore** — a correction from an earlier draft that grew the CSV with an empty placeholder entry whenever an assignment needed a slot that didn't exist yet. That contradicts Section 0/4's premise that every CSV entry is already a real, pre-existing on-chain signer from the start: this service never has a legitimate reason to append a placeholder. If every entry is already claimed by some assignment, there's nothing left to claim (`ErrNoPositionAvailable`, Section 5e) — the CSV simply doesn't grow to make room.
```go
// CollapseCSV removes the entry at `index` and shifts everything after it
// left by one — the exact splice from Section 5e's worked example. Takes the
// already-read csvValue/version rather than re-reading, because Section 5e's
// delete flow needs that same read for its own pre-deletion snapshot before
// this function ever runs — one read, shared. Returns the version the
// collapse produced.
func CollapseCSV(ctx context.Context, secret VaultSignerManagedSecret, csvValue string, version, index int) (newVersion int, err error) {
    parts := strings.Split(csvValue, ",")
    if index >= len(parts) {
        return 0, ErrIndexOutOfRange
    }
    parts = append(parts[:index], parts[index+1:]...)
    return putCSV(ctx, secret, strings.Join(parts, ","), version) // putCSV returns the version it just wrote
}
```

### 5b. Personal envs — create + delete only, no update, no read-back

**Identity for personal secrets is not `owner_id` — it's the caller's actual Trovo username, resolved from a different field per population.** This is a correction from an earlier draft of this plan, which derived the Vault key from `current_user_id`/`owner_id`. That's wrong for this flow specifically:

- **Trovo Admin** → `AdminUser.Username` (`gorm:"unique;not null"` — always present, no gate needed).
- **Organization member** → `OrganizationMember.TrovoWalletUsername` (`*string`, nullable) — the org member's linked username in the Trovo wallet app itself. Vault secrets are never assigned to an organization as a whole; they're personal to the individual member, identified by *this* field, not the member's `ID` or `Email`.

`AllowOrgOrTrovoAdminNormalized` doesn't put either of these in context — it only sets `current_user_id`/`current_user_email`/`user_type`. So every personal-secrets handler (list, create, delete — Section 6) starts by resolving the real Trovo username with a small helper, before touching Vault at all:

```go
// services/personalenv_service.go
//
// ResolveTrovoUsername looks up the caller's actual Trovo app username —
// AdminUser.Username for a Trovo admin, OrganizationMember.TrovoWalletUsername
// for an org member — using current_user_id/user_type from
// AllowOrgOrTrovoAdminNormalized. This is the identity used to derive the
// Vault key for PERSONAL SECRETS ONLY. It is deliberately NOT the identity
// used for managed secrets (Section 4's VaultSignerAssignment.OwnerRefID,
// which stores OrganizationMember.ID, resolved from Email only once at
// assignment-creation time — never TrovoWalletUsername) — the two
// resolvers exist because these two flows use two different, unrelated
// identity fields for an org member, and conflating them is a correctness
// bug: an org member with no linked wallet can still own a signer slot, but
// cannot use personal secrets.
func ResolveTrovoUsername(c *gin.Context, db *gorm.DB) (username string, ok bool) {
    userType := c.GetString("user_type")
    currentUserID := c.GetString("current_user_id")

    if userType == "trovo_admin" {
        var admin middleware.AdminUser
        if err := db.First(&admin, "id = ?", currentUserID).Error; err != nil {
            return "", false
        }
        return admin.Username, true // always non-empty — not-null in schema
    }

    var member models.OrganizationMember
    if err := db.First(&member, "id = ?", currentUserID).Error; err != nil {
        return "", false
    }
    if member.TrovoWalletUsername == nil || *member.TrovoWalletUsername == "" {
        return "", false // not linked — caller must show the friendly message below
    }
    return *member.TrovoWalletUsername, true
}
```

**The friendly-message gate:** if `ResolveTrovoUsername` returns `ok == false` for an org member (i.e. `TrovoWalletUsername` isn't set), the handler stops before any Vault call and responds with a message that never mentions internal jargon like "PERSONAL_ENVS" or "personal env" — the feature is always **"personal secrets"** in anything user-facing:

```go
c.JSON(http.StatusForbidden, gin.H{
    "error":   "trovo_wallet_not_linked",
    "message": "Connect your profile to the Trovo app to enable personal secrets.",
})
```

This gate applies uniformly to **all three** personal-secrets endpoints (`GET` list, `POST` create, `DELETE`) — there's no partial-access mode where an unlinked org member can list prefixes but not create, since without a resolvable username there is no way to even compute the Vault key to check `exists` against. Trovo Admins never hit this gate, since `AdminUser.Username` is a required column.

**Key derivation**, once a username is resolved:
```go
fullKey := strings.ToUpper(prefix) + "_" + strings.ToUpper(username)
// e.g. AUTO_APPROVE + tunde -> AUTO_APPROVE_TUNDE
// "tunde" here is AdminUser.Username OR OrganizationMember.TrovoWalletUsername,
// resolved above — never owner_id, never Email.
```

With the username resolved and the key derived, `vaultclient/personalenv.go` exposes three operations, none of which ever return a value:

```go
// Exists checks metadata only — never pulls the value into memory just to
// answer "does this exist."
func Exists(ctx context.Context, client *vaultapi.Client, mount, path string) (bool, error) {
    _, err := client.KVv2(mount).GetMetadata(ctx, path)
    if err != nil {
        if errors.Is(err, vaultapi.ErrSecretNotFound) {
            return false, nil
        }
        return false, err
    }
    return true, nil
}

// Create fails if the secret already exists — this is not an upsert.
// The handler calls Exists first and returns 409 before ever reaching here,
// but this is checked again at this layer too, defensively.
func Create(ctx context.Context, client *vaultapi.Client, mount, path, value string) (version int, err error) {
    if exists, err := Exists(ctx, client, mount, path); err != nil {
        return 0, err
    } else if exists {
        return 0, ErrAlreadyExists
    }
    secret, err := client.KVv2(mount).Put(ctx, path, map[string]interface{}{"value": value})
    if err != nil {
        return 0, err
    }
    return secret.VersionMetadata.Version, nil
}

// Delete fully purges the secret and all of its version history — not a
// soft/versioned delete. A soft delete (client.KVv2(mount).Delete) would
// leave the metadata path resolvable, so Exists would keep reporting
// exists:true after a "delete." DeleteMetadata removes the metadata itself,
// so a subsequent Exists call correctly reports false and a later Create
// for the same prefix starts clean.
func Delete(ctx context.Context, client *vaultapi.Client, mount, path string) error {
    return client.KVv2(mount).DeleteMetadata(ctx, path)
}
```

There is no `Get`/`Update` function in this file at all — the API surface intentionally has no code path capable of reading a personal-env value back out of Vault once it's stored.

### 5c. Signer Key Validation Pipeline — reuse existing error types

Same three steps as before (trim → `keypair.ParseFull()` → Horizon activation check via `network.GetBlockchainClient()`), but **reuse this repo's existing typed errors instead of inventing new ones**:

- `internal/errors.ErrorInvalidPublicKey` — already implements `HTTPCode() → 400` and a `JSONError()` shape matching `{"error": "...", "data": "...", "message": "..."}`. Use this for a `ParseFull()` failure instead of a bespoke error.
- `internal/errors.ErrorBlockchainAccountNotActivated{PublicKey: ...}` — **already exists in this exact codebase for this exact purpose** (`Message() → "Account not funded"`, `HTTPCode() → 404`). Use this for the "not found on Horizon" case rather than a new `422`.

```go
import (
    tErrors "admin-panel-dashboard/internal/errors"
    "admin-panel-dashboard/internal/network"
    "github.com/stellar/go/keypair"
    "github.com/stellar/go/clients/horizonclient"
)

func ValidateSignerValue(raw string) (publicKey string, err error) {
    trimmed := strings.TrimSpace(raw)

    kp, parseErr := keypair.ParseFull(trimmed)
    if parseErr != nil {
        return "", &tErrors.ErrorInvalidPublicKey{PublicKey: trimmed}
    }
    publicKey = kp.Address()

    client := network.GetBlockchainClient()
    if _, horizonErr := client.AccountDetail(horizonclient.AccountRequest{AccountID: publicKey}); horizonErr != nil {
        // fail closed on any Horizon error, not just 404 — matches the original plan's intent
        return "", &tErrors.ErrorBlockchainAccountNotActivated{PublicKey: publicKey}
    }
    return publicKey, nil
}
```

This keeps the response shape consistent with how the rest of this codebase already reports Stellar validation failures (see `internal/network/main.go`'s use of the same `internal/errors` package), rather than introducing a second error envelope just for this feature.

**Applies to both flows**, same as the original plan — signer slots and personal-env **creates** both run this pipeline before any Vault write, and personal envs still stop at the Vault write (never reach 5d). Personal-env **deletes** submit no value, so this pipeline doesn't run at all on `DELETE` — there's nothing to validate when removing a secret.

### 5d. On-Chain Signer Swap — build on `internal/network`, don't duplicate it

**Scope unchanged:** signer slots only, active indices only, gated by the same three conditions (old value is a valid keypair, new value is valid+activated, and they differ).

The original plan's `SwapSigner()` steps 1–11 are unchanged in logic. What changes is *where the transaction-building/submission code lives* — this codebase already has hardened helpers for exactly this in `internal/network/main.go`:

- `network.GetBlockchainClient()` for the Horizon client
- `network.GetBlockchainNetworkPassPhrase()` for the network passphrase used to sign
- The existing `SubmitXdrWithSignature*` family shows this repo's established pattern for Horizon submission error handling (checking `*horizonclient.Error`, logging `Problem.Extras` and `ResultCodes()`, alerting via `discord.Say` on connectivity failures). `swap_service.go` should follow the same error-handling shape for consistency, even though it submits a freshly-built `SetOptions` transaction rather than a pre-signed XDR from a caller.

Everything else — fetching the wallet's current signer weight and sequence number from Horizon, building the two-operation `SetOptions` transaction (remove old signer weight 0, add new signer at old weight), signing with every active signer except the target index, submitting, and writing the `VaultSignerAuditLog` row with `StellarTxHash`/`StellarTxStatus` — is unchanged from the original plan's Section 4d, including:
- The **Vault-then-ledger ordering risk** (Vault is written before submission; a submission failure leaves them mismatched by design, no rollback/compensation).
- The **threshold assumption** that `active_signing_count - 1` signers meet the wallet's high threshold — worth confirming against real account config before relying on this in production.
- The **fallback path** (old value invalid, or new == old) skipping the on-chain step entirely and doing a plain Vault write.

### 5e. Assignment Lifecycle: Claim-on-First-Write, and the Delete/Collapse Transaction

This is new relative to the earlier draft, which treated `VaultSignerAssignment` as create-only, and relative to an even later revision of *this* section, which auto-assigned `Position` serially at creation time. Neither is correct: **`Position` starts `NULL` and is claimed only when the owner submits their own real signer value for the first time** — creation itself never touches `Position`, Vault, or the chain at all.

**Creation** (`POST .../assignments`, `signer_service.go`) — no position in the request, and none computed either. A plain row insert:
```go
func CreateAssignment(db *gorm.DB, managedSecretID, ownerRefID, ownerType string) (*VaultSignerAssignment, error) {
    assignment := &VaultSignerAssignment{
        ID: uuid.NewString(), ManagedSecretID: managedSecretID, Position: nil,
        OwnerRefID: ownerRefID, OwnerType: ownerType, AssignedAt: time.Now(),
    }
    if err := db.Create(assignment).Error; err != nil {
        return nil, err
    }
    return assignment, nil
}
```
No locking is needed here — there's no position to race over yet.

**Claiming a position** (`services.ClaimPosition`, called from `PUT /me/vault-signer/secrets/:secretId` — Section 6 — the *first* time this assignment's owner submits a value): finds the lowest 1-based position not already claimed by any other assignment on the same managed secret, and assigns it. Every CSV entry is assumed to already be a real, pre-existing on-chain signer (Section 0/4) — claiming a position never grows the CSV or writes a placeholder into it; it only decides *which already-existing* entry this assignment now controls. If every entry is already claimed by some other assignment, there's nothing left (`ErrNoPositionAvailable`, surfaced as `409`).
```go
func ClaimPosition(ctx context.Context, db *gorm.DB, vc *vaultapi.Client, assignmentID string) (position int, err error) {
    err = db.Transaction(func(tx *gorm.DB) error {
        var assignment VaultSignerAssignment
        if err := tx.First(&assignment, "id = ?", assignmentID).Error; err != nil {
            return err
        }
        if assignment.Position != nil {
            position = *assignment.Position // already claimed — no-op
            return nil
        }

        // SELECT ... FOR UPDATE serializes concurrent first-time claims on
        // the same managed secret, so two owners' first PUT can't race for
        // the same lowest-available position.
        var secret VaultSignerManagedSecret
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&secret, "id = ?", assignment.ManagedSecretID).Error; err != nil {
            return err
        }

        csvValue, _, err := getCSV(ctx, secret) // Section 5a
        if err != nil {
            return err
        }
        totalPositions := len(strings.Split(csvValue, ","))

        var claimedPositions []int
        tx.Model(&VaultSignerAssignment{}).
            Where("managed_secret_id = ? AND position IS NOT NULL", assignment.ManagedSecretID).
            Pluck("position", &claimedPositions)
        claimed := make(map[int]bool, len(claimedPositions))
        for _, p := range claimedPositions {
            claimed[p] = true
        }

        for p := 1; p <= totalPositions; p++ {
            if !claimed[p] {
                position = p
                break
            }
        }
        if position == 0 {
            return ErrNoPositionAvailable
        }
        return tx.Model(&VaultSignerAssignment{}).Where("id = ?", assignmentID).
            UpdateColumn("position", position).Error
    })
    return position, err
}
```
A no-op if the assignment already has a position — covering both "already claimed by this owner" and "reassigned, but the previous owner had already claimed it," since `EditAssignment` (Section 6) always retains `Position` exactly as-is on reassignment. Once claimed, `PUT`'s handler converts to `index := position - 1` (Section 4's `csvIndex`) and runs the normal swap/write logic (Section 5d) — since the claimed CSV entry is always a real pre-existing key by construction, this is always a genuine swap opportunity, not a first-time-provisioning fallback.

**Deletion** (`DELETE .../assignments/:assignmentId`, looked up by the assignment's own ID — not by position, since position may be null) is the involved one, but **only for a claimed position**. An assignment that never claimed a position deletes with nothing more than a row delete:
```go
if assignment.Position == nil {
    return db.Delete(&assignment).Error // no CSV entry was ever assigned to this row
}
```
No renumbering (nothing was ever part of the position sequence), no floor check (no CSV entry to collapse against), no Vault or chain interaction at all.

For a **claimed** position, there are two separate `keypair.ParseFull()` questions in this flow, and it matters not to conflate them:
1. Does the *deleted position itself* hold a real key? (Section 5d's existing structural check, reused) — decides light path vs. full path, exactly as before.
2. **How many entries in the *whole CSV* are currently valid keypairs?** — decides whether the on-chain removal happens at all. This is the corrected trigger (a correction from an earlier draft of this section, which tied it to the deleted position's relationship to `ActiveSigningCount` — wrong): **the request specifically means the count of valid signer keypairs currently in the CSV, not the raw entry count.** This count is a capacity/redundancy check ("is it safe to shrink the active set by one"), not a literal tally of live on-chain signers, so it's still counted across every entry, active and spare alike.

   **A further, independent gate added during implementation, on top of the correction above:** on-chain removal additionally requires the *deleted position itself* to be within the active range (`index < ActiveSigningCount`). This is not a restriction on the count — it's a separate condition that must also hold. Restricting the *count* to the active range was tried and rejected during implementation: the active range only has `ActiveSigningCount` slots, so a count capped to that range can never exceed `ActiveSigningCount`, which would make the `> ActiveSigningCount` condition below permanently false and silently disable on-chain removal entirely. The position gate exists for a different, narrower reason: Section 4a/5d already establishes that a spare position's key, however valid, was never registered on-chain to begin with — so attempting to remove one doesn't hold up regardless of the count.

   Putting both together, on-chain removal is attempted only if `index < ActiveSigningCount` **and** `validKeypairCount > ActiveSigningCount` *(using `ActiveSigningCount` as the generalization of the request's illustrative `3`, consistent with every other active/spare boundary in this plan)* — i.e., only when the position being removed is itself active, and removing one still leaves at least `ActiveSigningCount` valid signers. If either condition fails — a spare position, or an active position when the whole CSV currently has `ActiveSigningCount` or fewer valid keys — the assignment is still deleted at the DB/CSV level, but the on-chain removal is skipped. For the spare case, there was simply never anything on-chain to remove. For the "at the floor" active case, the key stays registered on the wallet even though this service no longer tracks it, because dropping the live signer set below the safe minimum is worse than leaving one untracked-but-harmless extra signer on-chain.

Both checks come from the same single CSV read, done once up front:

1. **Read the CSV, then floor-check its live length:** `len(parts) - 1 < 4` → reject with `409` — the same minimum-4 rule from Section 4, now also enforced on the way down. This is checked against the CSV's actual length, not the assignment count — the two are no longer 1:1 now that positions are claimed rather than auto-assigned at creation (some assignments may still be unclaimed), so an unclaimed assignment must never count toward "how many CSV entries would remain."
2. **The same read becomes `originalCSV`/`originalVersion`**, reused by Step 5's collapse and, if needed, Step 5's on-chain reversal:
   ```go
   originalCSV, originalVersion, err := getCSV(ctx, secret) // Section 5a
   parts := strings.Split(originalCSV, ",")
   index := csvIndex(position)
   targetIsValidKey := keypair.ParseFull(strings.TrimSpace(parts[index])) == nil // question 1

   // Counted across the whole CSV, not just the active range — this is a
   // capacity/redundancy check, not a literal tally of live on-chain
   // signers. Restricting it to the active range would cap it at
   // ActiveSigningCount and make "> ActiveSigningCount" below unsatisfiable.
   validKeypairCount := countValidKeypairs(parts) // question 2
   ```
   ```go
   func countValidKeypairs(parts []string) int {
       count := 0
       for _, raw := range parts {
           if _, err := keypair.ParseFull(strings.TrimSpace(raw)); err == nil {
               count++
           }
       }
       return count
   }
   ```
   `validKeypairCount` is deliberately computed **before** this deletion touches anything, across the whole CSV — it answers "how much redundancy currently exists," and the on-chain trigger in Step 4 reads directly off this pre-deletion number (combined with the separate `index < ActiveSigningCount` position gate), not a recomputed post-collapse one.
   - **If `targetIsValidKey` is false:** defensive fallback only — light path. Every CSV entry is assumed to already be a real, pre-existing signer (Section 0), so a *claimed* position's entry should always be a valid keypair by construction; this branch exists in case that assumption is ever violated (e.g. Vault drift outside this service's control), not because it's an expected outcome of the normal claim/swap flow. Delete the row, renumber every `Position > deleted` down by one, splice the entry out of the CSV, done. **No on-chain call, regardless of `validKeypairCount`** — there was never a registered signer at *this* position to remove. Audit as `assignment_delete` with no `StellarTxHash`.
   - **If `targetIsValidKey` is true:** the full flow below.
3. **Begin a DB transaction.** Within it (nothing committed yet):
   - Delete the `VaultSignerAssignment` row at the target position.
   - Renumber: `UPDATE vault_signer_assignments SET position = position - 1 WHERE managed_secret_id = ? AND position > ?` (the deleted position).
4. **On-chain removal now happens *before* the Vault write — a correction from the previous draft, which collapsed Vault first.** Attempted only if `index < ActiveSigningCount` **and** `validKeypairCount > ActiveSigningCount` (Step 2, from the pre-deletion, active-range-only CSV read). Before building the transaction, capture the removed key's **current signer weight** from Horizon (same pattern as Section 5d's swap, Step 5) — needed if a reversal turns out to be required in Step 6. This is a single-operation `SetOptions` (remove only — `Signer{Address: removedPublicKey, Weight: 0}`), signed by every *remaining* active signer.
   - **If not required** (`validKeypairCount <= ActiveSigningCount`): skip straight to Step 5.
   - **If attempted and it fails** (Horizon rejects the submission): **nothing has touched Vault yet** — just roll back the DB transaction. This is now the simple case: no partial state, no compensation needed, because the write ordering means a failure here is caught before it can leave anything collapsed. Audit as `assignment_delete` with `StellarTxStatus = failed` and no Vault change.
   - **If it succeeds:** proceed to Step 5.
5. **Collapse the Vault CSV** — `_, err := vaultclient.CollapseCSV(ctx, secret, originalCSV, originalVersion, csvIndex(deletedPosition))` (Section 5a), CAS'd against `originalVersion`.
   - **If it succeeds:** commit the DB transaction. Fully consistent across all three systems — this is the common case.
   - **If it fails** (e.g. Vault is unreachable) **and Step 4's on-chain removal already succeeded:** this is the situation the correction targets. Rather than leaving the signer removed on-chain while the DB and Vault still show the old assignment, **attempt to add the signer back on-chain** — a second `SetOptions` (`Signer{Address: removedPublicKey, Weight: capturedWeight}`, the weight captured in Step 4), signed by the same active-signer set, submitted against the wallet's now-current sequence number (it moved forward after Step 4's transaction):
     ```go
     if collapseErr != nil && onChainRemovalHappened {
         if _, reAddErr := network.SubmitXdrWithSignature /* or equivalent SetOptions add */ (...); reAddErr == nil {
             // Reversal succeeded — the wallet's on-chain signer set is back
             // to where it started. Nothing was ever persisted to Vault, so
             // there's nothing to discard there; roll back the DB transaction
             // so it matches too.
             tx.Rollback()
             return &tErrors.CustomError{
                 Err: "vault-unreachable-changes-reverted",
                 ErrMessage: "Vault service could not be reached; all changes have been reverted.",
                 Code: http.StatusBadGateway,
             }
         }
         // Reversal ALSO failed — see the residual-risk note below.
     }
     ```
6. **Write the audit log row** (`assignment_delete`) after the outcome is known, including `StellarTxHash`/`StellarTxStatus` for both the removal and, if it happened, the reversal.

**The one case this can't fully resolve:** Step 5's on-chain reversal *also* failing — on-chain removal succeeded, the Vault write failed, and the attempt to add the signer back also failed. At that point the wallet's real signer set has diverged from what the DB and Vault both still say, and no automatic action here can safely reconcile it — this is the genuine last-resort case, and it's *narrower* than what the previous draft had to accept (that version could reach an inconsistent state from a single Vault failure; this one needs on-chain removal to succeed, then Vault to fail, then the on-chain reversal to *also* fail). Handle it the same way as any other double-failure in this plan: bounded retries on the reversal, then alert loudly (this codebase's `discord.Say` pattern from `internal/network/main.go`) rather than letting it surface only as an audit-log row — and leave the DB assignment uncommitted/unresolved rather than guessing which side to trust, since a human now needs to look at the wallet's actual on-chain signers before anything here can be corrected safely.

**Worked example — claim, reassign, delete-before-claim:** a managed secret's CSV already holds 5 real signers `[k1,k2,k3,k4,k5]` (seeded at registration, per Section 0 — this service never put any of them there). An admin creates assignment A for Alice: `Position = NULL`, nothing else happens. Weeks later Alice logs in and calls `PUT` for the first time: `ClaimPosition` finds no positions claimed yet, so she claims position `1`; her write then goes through the normal swap flow against `k1` (a real, already-valid key), rotating it to her own. Separately, the admin reassigns assignment A to Bob (`PATCH`) before Alice ever logs in: `Position` stays `NULL` through the reassignment; when Bob eventually logs in and calls `PUT`, *he's* the one who claims position `1` and rotates `k1` — Alice, despite being the original assignee, never touched Vault or the chain at all. If the admin instead deletes assignment A while it's still unclaimed (`DELETE`), it's a plain row delete — `k1` was never associated with a DB row's CSV entry, so nothing renumbers and nothing collapses.

**Worked examples for the claimed-position delete flow, both against `ActiveSigningCount = 3`:**
- **5 assignments, all 5 claimed and holding real keys, position 2 (index 1) is deleted.** `index 1 < 3` (active) and `validKeypairCount = 5 > 3` → both conditions hold → on-chain removal attempted.
  - *Full success:* `k2` is removed on-chain; the CSV collapses from `[k1,k2,k3,k4,k5]` to `[k1,k3,k4,k5]`; DB commits. Positions 3, 4, 5 become 2, 3, 4.
  - *On-chain removal fails outright:* nothing else has happened yet — DB rolls back immediately, Vault was never touched, still `[k1,k2,k3,k4,k5]`.
  - *On-chain removal succeeds, then Vault is unreachable:* the service re-adds `k2` back on-chain at its original weight. If that succeeds, the DB rolls back and the caller gets "Vault service could not be reached; all changes have been reverted" — the wallet, Vault, and DB all end up exactly as they were before the request.
- **Exactly 3 assignments, all active, all 3 hold real keys, position 1 (index 0) is deleted.** `index 0 < 3` (active) but `validKeypairCount = 3`, and `3 > 3` is false → **on-chain removal is skipped** even though the deleted position is active. The row deletes, positions renumber, the CSV collapses — but `k1` stays registered on the wallet, since dropping the live signer set below the floor is worse than one untracked-but-harmless extra signer.
- **5 assignments (3 active + 2 spare), all 5 hold real keys, position 4 (index 3, a spare position) is deleted.** `validKeypairCount = 5 > 3` holds, but `index 3 < 3` is false → **on-chain removal is skipped** regardless of the count, because a spare position's key was never registered on-chain to begin with (Section 4a/5d) — there's nothing there to remove. The row deletes, positions renumber, the CSV collapses as normal; no chain interaction at all.

---

## 6. API Design

All routes behind `apiV1 := router.Group("/api/v1")`, matching every other component's `controllers/main.go`.

### Self-service endpoints (Trovo Admin or Organization member)

| Method | Path | Middleware | Purpose |
|---|---|---|---|
| GET | `/api/v1/me/vault-signer/secrets` | `AllowOrgOrTrovoAdminNormalized` | List signer secrets + my Signer Position, `null` if I haven't claimed one yet |
| GET | `/api/v1/me/vault-signer/secrets/:secretId` | `AllowOrgOrTrovoAdminNormalized` | Get current value of *my* slot — `{"position": null}` with no value if I haven't claimed one yet (Section 5e); there's nothing in Vault tied to this assignment until I do |
| PUT | `/api/v1/me/vault-signer/secrets/:secretId` | `AllowOrgOrTrovoAdminNormalized` | Submit new value for *my* slot. If this is the first time I've ever set a value here, claims my Signer Position first (Section 5e), then branches into on-chain swap (active) vs plain write (spare) |
| GET | `/api/v1/me/vault-signer/personal-envs` | `AllowOrgOrTrovoAdminNormalized` | List all `PERSONAL_ENVS` prefixes (with the label parsed from `PREFIX:FriendlyLabel`), my derived key name, `exists` flag — **never a value** |
| POST | `/api/v1/me/vault-signer/personal-envs/:prefix` | `AllowOrgOrTrovoAdminNormalized` | **Create** my personal secret (writes the value into Vault). `409` if it already exists — this is not an upsert |
| DELETE | `/api/v1/me/vault-signer/personal-envs/:prefix` | `AllowOrgOrTrovoAdminNormalized` | **Delete** my personal secret (fully purges it from Vault). `404` if it doesn't exist |

There is deliberately **no `GET /me/vault-signer/personal-envs/:prefix`** and **no `PUT`**. Once a personal-env secret is created, its value is never readable again through this API — by anyone, including its owner — and it cannot be edited in place. Changing a value means `DELETE` then `POST` again.

**All three personal-secrets endpoints require a resolvable Trovo username first** (Section 5b). An org member without a linked wallet gets this on every one of them, before any Vault interaction:
```json
// 403 — org member with no OrganizationMember.TrovoWalletUsername
{ "error": "trovo_wallet_not_linked", "message": "Connect your profile to the Trovo app to enable personal secrets." }
```
Trovo Admins never see this response (`AdminUser.Username` is always set). Note the message says "personal secrets," matching the user-facing naming convention from Section 0 — never "personal env(s)."

### Admin-only endpoints (Trovo SuperAdmin)

| Method | Path | Middleware | Purpose |
|---|---|---|---|
| GET | `/api/v1/admin/vault-signer/managed-secrets` | `AuthenticateSuperAdmin` | List managed signer secrets, each with its current `VaultSignerAssignment` rows (owner per position) |
| POST | `/api/v1/admin/vault-signer/managed-secrets` | `AuthenticateSuperAdmin` | Register a new managed signer secret (incl. `walletPublicKey`) |
| GET | `/api/v1/admin/vault-signer/managed-secrets/:id/assignments` | `AuthenticateSuperAdmin` | List current `VaultSignerAssignment` rows for one managed secret, ordered by `Position` |
| POST | `/api/v1/admin/vault-signer/managed-secrets/:id/assignments` | `AuthenticateSuperAdmin` | **Creates a `VaultSignerAssignment` row with `Position = NULL`.** The owner claims a position themselves the first time they call `PUT` on their own slot (Section 5e) — an admin never supplies or sees a position at creation time. |
| PATCH | `/api/v1/admin/vault-signer/managed-secrets/:id/assignments/:assignmentId` | `AuthenticateSuperAdmin` | **Reassigns the owner** of an assignment, addressed by its own ID (not by position, which may be null). `Position` — null or claimed — is always retained exactly as-is. Never touches Vault or the chain — see the note below. |
| DELETE | `/api/v1/admin/vault-signer/managed-secrets/:id/assignments/:assignmentId` | `AuthenticateSuperAdmin` | **Removes an assignment**, addressed by its own ID. If a position was ever claimed, renumbers every later position down by one so there's never a gap. Full mechanics in Section 5e — this is the one non-trivial endpoint in the whole feature. |
| GET | `/api/v1/admin/vault-signer/audit-log` | `AuthenticateSuperAdmin` | View change history (all kinds) |

Route prefix is `vault-signer/...` rather than bare `secrets`/`personal-envs` to avoid collisions with any future unrelated `/me/secrets` route in this shared router — this codebase has many components sharing one `apiV1` group, unlike the original single-purpose service where a bare path was safe.

**Assignment-creation request/response.** No `position` field in the request, and none in the response either — `position` comes back `null`, because it starts that way and only gets claimed by the owner's own first `PUT` (Section 5e).
```json
// POST /api/v1/admin/vault-signer/managed-secrets/:id/assignments
// assigning a Trovo admin
{ "ownerIdentifier": "obi", "ownerType": "trovo_admin" }
```
```json
// assigning an org member — no wallet linkage required for this
{ "ownerIdentifier": "member@acme-trust.com", "ownerType": "org_member" }
```
```json
// 201 response — position is null until this owner (or a future reassigned
// owner) logs in and submits their own value for the first time
{ "data": { "id": "...", "managedSecretId": "...", "position": null, "ownerRefId": "8f1a...", "ownerType": "org_member", "ownerLabel": "member@acme-trust.com", "assignedAt": "..." } }
```
`ownerRefId` is the real stored FK value; `ownerLabel` is resolved fresh via a join at read time (Section 4), so it stays current even if the member's email later changes — it is never what's compared or stored for authorization. The handler rejects with `404`/`422` if `ownerIdentifier` doesn't resolve to a real `AdminUser.Username`/`OrganizationMember.Email`, and with `409` if this owner already has an assignment on this managed secret (Section 4's `(ManagedSecretID, OwnerRefID)` uniqueness):
```json
// 409 response — this owner already has a row on this managed secret
{ "error": "owner_already_assigned", "message": "This owner already has an assignment on this managed secret." }
```

**Assignment-edit request/response** — addressed by the assignment's own ID (`:assignmentId`), not by position, since position may still be null: swap who owns an assignment without touching its `Position` (null or claimed — always retained exactly as-is), its Vault CSV entry, or the chain at all.
```json
// PATCH /api/v1/admin/vault-signer/managed-secrets/:id/assignments/7f3e...
{ "ownerIdentifier": "different-member@acme-trust.com", "ownerType": "org_member" }
```
```json
// 200 response — this assignment had already claimed position 3 before
// being reassigned; the new owner inherits that same position
{ "data": { "id": "7f3e...", "managedSecretId": "...", "position": 3, "ownerRefId": "9c2d...", "ownerType": "org_member", "ownerLabel": "different-member@acme-trust.com", "assignedAt": "..." } }
```
**Why this never touches Vault:** whatever Stellar key currently sits at this assignment's CSV entry (if any position has been claimed at all) doesn't change just because the DB row's `OwnerRefID` changed — the *new* owner simply doesn't have a valid key registered there yet. If `Position` was still `null` at the time of reassignment, it stays `null` — the new owner is the one who ends up claiming it, on their own first `PUT`. Either way, the very next time the new owner calls `PUT /me/vault-signer/secrets/:secretId`, that request already runs claim-if-needed (Section 5e) followed by the full Section 5c/5d pipeline, and triggers a normal on-chain swap exactly as if any owner were rotating their own key. Reassignment is purely "who is allowed to call `PUT` for this assignment from now on" — it doesn't need to know or care what value (if any) is currently sitting there.

**Reassignment also respects `(ManagedSecretID, OwnerRefID)` uniqueness (Section 4):** if `ownerIdentifier` resolves to someone who already owns a *different* assignment on this same managed secret, the request is rejected with the same `409 owner_already_assigned` shape as create — reassigning is not a way around the constraint. Reassigning an assignment back to the owner it already has is always a no-op save, never a conflict with itself.

**Assignment-deletion response** — also addressed by `:assignmentId`; see Section 5e for the full request lifecycle. An assignment that never claimed a position deletes trivially:
```json
// DELETE /api/v1/admin/vault-signer/managed-secrets/:id/assignments/7f3e...
// 200 response — this assignment's Position was still null
{ "data": { "claimed": false } }
```
For a claimed position, the response carries the same fields as before:
```json
// 200 response — success, claimed position 2
{ "data": { "claimed": true, "deletedPosition": 2, "renumberedCount": 3, "vaultCollapsed": true, "onChainRemoval": { "attempted": true, "stellarTxHash": "a1b2...", "stellarTxStatus": "success" } } }
```
`onChainRemoval` is `null` when the deleted position's entry wasn't a valid keypair (the defensive fallback, Section 5e), or when `validKeypairCount <= ActiveSigningCount` or the position is a spare one — nothing to remove on-chain in any of those cases, and the row/renumber still commit normally.
```json
// 502 response — on-chain removal itself failed outright (Section 5e)
// Nothing was ever written to Vault; the assignment was not deleted.
{ "error": "assignment_delete_failed", "message": "Could not remove the signer on-chain; the assignment was not deleted.", "stellarTxStatus": "failed" }
```
```json
// 502 response — on-chain removal succeeded, but Vault couldn't be reached
// afterward (Section 5e); the removal was reversed on-chain and the
// assignment was not deleted.
{ "error": "vault-unreachable-changes-reverted", "message": "Vault service could not be reached; all changes have been reverted." }
```
Both failure responses leave the DB transaction rolled back (the row and everyone's `Position` are exactly as they were) — the caller always sees a clean failure, never a partial one, because Vault is only ever written *after* the on-chain step succeeds (Section 5e). The one case neither response can promise to have resolved is the second one's own reversal call also failing — that residual, double-failure case is surfaced operationally (alerting), not through this response, since by definition the API has no reliable way to guarantee a clean outcome once two independent systems have both failed on the same request.

### Request/response shapes

**Signer slots** (`PUT /me/vault-signer/secrets/:secretId`) keep the same request/response shape as the original plan: `{"newValue": "..."}` in, `{"vaultVersion": N}` / `{"vaultVersion": N, "stellarTxHash": "...", "stellarTxStatus": "..."}` out. The one addition: if this is the caller's first-ever write for this assignment, the handler calls `ClaimPosition` first (Section 5e) — this can fail with `409 {"error": "no_position_available", "message": "Every signer position on this managed secret is already claimed."}` if every CSV entry is already claimed by some other assignment, before the request ever reaches the validation/swap pipeline.

**Personal envs** are asymmetric, matching the create/delete-only lifecycle:
```json
// POST /api/v1/me/vault-signer/personal-envs/:prefix
{ "newValue": "SABC...TUNDESEED" }
```
```json
// 201 response — no value echoed back, ever
{ "data": { "prefix": "AUTO_APPROVE", "vaultKey": "AUTO_APPROVE_...", "vaultVersion": 1 } }
```
```json
// DELETE /api/v1/me/vault-signer/personal-envs/:prefix
// 204 No Content — nothing to return; the secret no longer exists
```
The `POST` response includes `vaultVersion` for audit/debugging purposes (it's metadata, not the secret), but under no circumstances does any personal-env response body include the submitted or stored value — the value only ever flows in on `POST`, and only ever out of scope (deleted) on `DELETE`.

Anywhere a response needs to display *who* owns a managed secret (e.g. the admin-only assignment listing in Section 6), the label is resolved fresh from `OwnerRefID` via a join — `AdminUser.Username` for `trovo_admin`, `OrganizationMember.Email` for `org_member` (Section 4) — never a stored copy, so it can't go stale if the underlying field changes. No fallback logic is needed here, since both source fields are mandatory; this is a different, simpler situation than personal secrets, where an unresolved `TrovoWalletUsername` blocks the action entirely rather than displaying anything.

### Response/error conventions — match `stakeholder`, not the original plan's proposed envelope

The original plan proposed a service-wide `{"error": "string", "code": "string"}` envelope (Section 12, point 6), reasonable for a standalone service with outside consumers. This codebase's newest component (`stakeholder`) instead uses the simplest idiom already established here:
```go
c.JSON(http.StatusOK, gin.H{"data": result})
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
```
with typed `internal/errors` values (like `ErrorInvalidPublicKey`) used where they already exist, falling back to `gin.H{"error": err.Error()}` otherwise. New vault-signer handlers should follow this same shape rather than introducing a third response convention into an already-mixed codebase (which currently has `response.Data`, `models.ErrorResponse`, and bare `gin.H` all in use across different components — `stakeholder`'s plain `gin.H` is the most recent and simplest, so that's the one to extend).

---

## 7. Swagger — Already Wired, Just Add Annotations

This repo already runs `swaggo/swag` (`internal/components/swagger/controllers`, generated `docs/swagger.json` + `docs/swagger.yaml`, `make swagger` / `make swagger-clean` in the `Makefile`). No new setup from the original plan's Section 7 is needed — just annotate the new handlers the same way every other component does, then run `make swagger` to regenerate:

```go
// PutSignerSecret godoc
// @Summary      Submit a new value for the caller's signer slot
// @Tags         vault-signer
// @Security     JwtTokenAuth
// @Security     OrganizationAuth
// @Produce      json
// @Param        secretId path string true "Managed secret ID"
// @Success      200 {object} vaultsignermodels.SignerSecretResponse
// @Failure      422 {object} gin.H
// @Router       /api/v1/me/vault-signer/secrets/{secretId} [put]
func (h *Handler) PutSignerSecret(c *gin.Context) { ... }
```

Note both security schemes (`JwtTokenAuth` and `OrganizationAuth`) are already defined in `main.go`'s top-level doc block — this feature is the first to legitimately need *both* on the same endpoint, since `AllowOrgOrTrovoAdminNormalized` accepts either.

---

## 8. Build Order

1. `go get github.com/hashicorp/vault/api && go mod tidy`
2. New GORM models in `internal/components/vaultsigner/models/db_models.go`, **and in the same step**, register all three in `internal/db/main.go`'s `migrateAdminSchemaTransaction` `AutoMigrate(...)` call — do not defer this to a later step or a separate PR; an unregistered model never gets a table
3. `vaultclient/signer.go` — CSV rebuild/CAS logic (unit-testable standalone, no Gin/DB dependency)
4. `services/validation_service.go` — trim → `ParseFull()` → Horizon activation check, using `network.GetBlockchainClient()` and reusing `internal/errors.ErrorInvalidPublicKey` / `ErrorBlockchainAccountNotActivated`
5. `services/swap_service.go` — `SwapSigner()`: old/new keypair gate, account/weight lookup via `network.GetBlockchainClient()`, build/sign/submit the 2-op `SetOptions` transaction, following the error-handling shape already established in `internal/network/main.go`'s `SubmitXdrWithSignature*` functions
6. `services/personalenv_service.go` — `ResolveTrovoUsername` (Section 5b: `AdminUser.Username` or `OrganizationMember.TrovoWalletUsername`, friendly-message gate if unresolved) + key derivation (`{PREFIX}_{username}`) on top of `vaultclient/personalenv.go`'s `Exists`/`Create`/`Delete`; also parses `PERSONAL_ENVS` (`PREFIX:FriendlyLabel` CSV) into the in-memory prefix→label map at startup, failing fast on any malformed entry
7. `handlers/signer_handler.go` + `controllers/main.go` route registration — `me/vault-signer/secrets*`. `PUT` claims a position first (`services.ClaimPosition`, Section 5e) if this assignment's `Position` is still null, then branches into plain write vs on-chain swap by `Position` (converted to `index` via `csvIndex`, Section 4) exactly as before
8. `handlers/personalenv_handler.go` — `me/vault-signer/personal-envs*` (list, **create**, **delete** — no get, no update)
9. `services/signer_service.go` — assignment **create** (a plain insert with `Position = nil`, no locking, no Vault touch — `ownerIdentifier` → `OwnerRefID` resolution, `404`/`422` if unresolvable) and assignment **list** (joins `OwnerRefID` → display label); then admin handlers for both, gated by `AuthenticateSuperAdmin`
10. `services/signer_service.go` — `ClaimPosition` (the lowest-unclaimed-position lookup + lock, Section 5e — build this before step 7's handler needs it); assignment **edit** (`PATCH`, addressed by assignment ID, reassign `OwnerRefID` only, `Position` always retained, no Vault/chain touch, Section 6); and assignment **delete** (`DELETE`, addressed by assignment ID — a plain row delete if `Position` is still null, otherwise the full renumber + CSV collapse + conditional on-chain removal transaction, Section 5e) — build and test the delete path last among these since it's the highest-risk piece, ideally against a Horizon testnet sandbox like `swap_service.go`
11. Swagger annotations on all new handlers + `make swagger`
12. Audit log wiring (all kinds, including `signer_swap`, `assignment_create`, `assignment_edit`, `assignment_delete` with tx hash/status where applicable)
13. Wire `vaultsigner.Init(router, s)` into `main.go` after `stakeholder.Init(router, s)`

---

## 9. Security Notes (updated for this repo's context)

- Least-privilege Vault token scoped to both `managed_secrets` paths **and** `personal-envs/*`.
- CAS on every signer-slot write (required); CAS on personal-env writes (defensive — single owner/writer per secret).
- Never log or return the full signer CSV to any client. Personal-env values go further: **no endpoint ever returns one, to anyone, including the owner** — enforced by the API surface itself (Section 6: no `GET`/`PUT` on a single personal-env, only `POST` create and `DELETE` purge).
- Signer values pass the same three validation layers (trim → `keypair.ParseFull()` → Horizon activation) for **both** flows — enforced server-side via `internal/errors.ErrorBlockchainAccountNotActivated`, never left to client-side validation.
- Horizon errors (timeouts, 5xx) fail closed on the activation check, matching how `internal/network` already treats Horizon failures elsewhere in this codebase.
- `VaultSignerAssignment.OwnerRefID`/`OwnerType` (Section 4) stores the real `AdminUser.ID`/`OrganizationMember.ID`, **not** `Username`/`Email` — a bug that compares `OwnerRefID` across the two `OwnerType`s without checking type first could let an admin's ID collide with an org member's ID string (they aren't drawn from the same namespace, and `AdminUser.ID` is numeric while `OrganizationMember.ID` is a UUID-style string). Always compare `(OwnerRefID, OwnerType)` as a pair, never `OwnerRefID` alone.
- **Storing `OwnerRefID` instead of a resolved `Username`/`Email` string is a deliberate anti-staleness fix, not just a style preference:** `OrganizationMember.Email` can change; a stored copy of it on `VaultSignerAssignment` would silently drift from reality the moment it does, breaking "is this my slot?" for that member with no error and no audit trail. The real ID can't do that. `ownerIdentifier` (email/username) is still how an admin *specifies* who to assign (Section 6) — it's just never what's persisted.
- **Managed secrets and personal secrets deliberately use two different identity fields for an org member, and this is intentional, not an oversight:** managed secrets resolve `ownerIdentifier` (an org member's `Email`) to `OwnerRefID` once, at assignment-creation time, with no ongoing gate; personal secrets resolve `TrovoWalletUsername` on every request, because the Vault key itself is derived from it (Section 5b) — falling back to `Email` there would tie a secret's actual storage path to a field this feature was specifically told not to use for that purpose. Do not "simplify" by unifying these two into one identity field.
- Personal secrets can only ever be created/deleted for the *caller's own* resolved Trovo username — there is no admin override and no endpoint that accepts a username or member ID as a parameter for this flow. An org member with no `TrovoWalletUsername` is blocked outright (Section 5b's friendly-message gate) rather than falling back to `Email`. Managed secrets have no such gate — an org member can be assigned and manage a signer slot with zero Trovo wallet linkage (Section 1).
- **`VaultSignerAssignment.Position` is nullable by design, not by omission** (Section 4/5e) — it starts `NULL` at creation and is claimed only by the owner's own first `PUT`. `ClaimPosition`'s `SELECT ... FOR UPDATE` on the managed secret row is required, not defensive polish: without it, two different owners' first-ever writes on the same managed secret could race for the same lowest-available position and both succeed, silently corrupting the `(ManagedSecretID, Position)` uniqueness guarantee at the exact moment it matters most. Every CSV entry is assumed to already be a real, pre-existing on-chain signer (Section 0) — this service never writes a placeholder into the CSV to "reserve" a position; if every entry is already claimed, there is nothing left (`ErrNoPositionAvailable`), not a reason to grow the CSV.
- Admin actions (`managed-secrets*`, assignments) require `AuthenticateSuperAdmin` specifically — not just "any authenticated admin" — since this grants control over live signing infrastructure. Confirm this matches the sensitivity bar the team wants; `EditLevelAdmin` is deliberately excluded here unless decided otherwise.
- Vault-then-ledger ordering is not atomic (Section 5d) — monitor via `VaultSignerAuditLog.StellarTxStatus` rather than assuming self-correction.
- The active-signer-minus-target-index signing set must meet the wallet account's **high threshold** — confirm against real account config before relying on this flow.
- The Vault token used for the swap flow can move signing authority on a live wallet — treat it operationally like a hot signing key, not app config storage.
- Test all four combinations of the swap gate explicitly (old valid/invalid × new same/different) before shipping — get this wrong and either a real rotation is silently skipped, or a signer that was never really registered gets "removed."
- **Assignment deletion (Section 5e) orders on-chain removal *before* the Vault write specifically to avoid the shape mismatch the swap flow accepts.** Because Vault is only ever touched after the on-chain step has already succeeded, an outright on-chain failure leaves Vault untouched — nothing to compensate. The one case that does need compensation is on-chain succeeding and the *following* Vault write then failing (e.g. Vault unreachable): the service re-adds the removed signer on-chain (at its captured original weight) and, if that reversal succeeds, rolls back the DB transaction and reports the failure plainly. The genuine residual risk is now narrower and requires three independent things to go wrong in sequence: on-chain removal succeeds, then Vault fails, then the on-chain reversal *also* fails. Handle that combination the same way as any other double-failure in this plan — bounded retries, then alert loudly (this codebase's `discord.Say` pattern from `internal/network/main.go`) rather than letting it surface only as an audit-log row — and leave the assignment uncommitted/unresolved for manual reconciliation rather than guessing which of the three systems to trust.
- The minimum-4 CSV floor (Section 4) is checked against the CSV's own live length, not the assignment count — the two are no longer 1:1 now that `Position` is claimed rather than auto-assigned at creation (Section 5e), so an unclaimed assignment must never be allowed to count toward "how many CSV entries would remain" on either side of that check. `DELETE` on a claimed position is rejected with `409` before any transaction begins if collapsing would bring the CSV below 4; `DELETE` on an unclaimed one never reaches this check at all, since there's no CSV entry involved.
- **The on-chain-removal trigger for assignment deletion is `index < ActiveSigningCount && validKeypairCount > ActiveSigningCount` (Section 5e)** — `validKeypairCount` is the count of *every* CSV entry (active and spare alike) that currently passes `keypair.ParseFull()`, computed *before* the deletion; it's a capacity/redundancy check, not a literal tally of live on-chain signers, so it deliberately isn't scoped to the active range (an active-range-only count could never exceed `ActiveSigningCount`, which would make the `>` comparison permanently false and silently disable on-chain removal — this was caught and reverted during implementation). The `index < ActiveSigningCount` half is a separate, independent gate added on top of the original count-only correction: a spare position's key is never registered on-chain regardless of the count (Section 4a/5d), so attempting to remove one doesn't hold up. Deleting a real, active key when the whole CSV currently has `ActiveSigningCount` or fewer valid keys skips the on-chain step entirely (that key stays registered on the wallet, untracked but harmless) rather than letting the live signer count drop below the safe minimum; deleting a spare position skips it because there was never anything on-chain to remove in the first place.

---

## 10. Decisions (previously open questions)

- **First-time provisioning:** confirmed as-is, and now further reinforced by the claim-on-first-write redesign (Section 5e) — an active slot with no valid old key only writes Vault and never adds an on-chain signer; no "add-only" swap variant is being built. Under the current model this case should be rare in practice, not the common path the original decision anticipated: every CSV entry is assumed to already be a real, pre-existing on-chain signer from the start (Section 0), and a position is only ever claimed once (`ClaimPosition`, Section 5e) — so by the time any `PUT` reaches the swap logic, the "old value" at that position should always be a genuine valid keypair. The plain-write fallback (Section 5d) still exists purely as a defensive path for if that assumption is ever violated (e.g. Vault drift outside this service's control), not as an expected first-time-provisioning route.
- **`EditLevelAdmin` access:** confirmed — no admin-only vault-signer action is reachable below `SuperAdmin`. `AuthenticateSuperAdmin` stays the gate for all of `admin/vault-signer/*` (managed-secrets, assignments, audit-log), with no `EditLevelAdmin` carve-out.

---

## 11. Implementation Summary

**Status: built.** Everything described above (Sections 1–8) has been implemented directly in this repo, on branch `ric-vault-signer-manager`, and matches the code exactly as of this writing — this plan was kept in sync with the implementation through several corrections found while actually writing it, not written up front and left to drift.

### What was added

**New dependency:** `github.com/hashicorp/vault/api v1.16.0` — deliberately *not* the latest (`v1.23.0`), which pulls in a `go 1.24` requirement and would have forced this repo's `go.mod` `go` directive up from `1.21`, plus unrelated transitive upgrades. `v1.16.0` is the newest version that still declares `go 1.21`, so it adds cleanly with no toolchain change and no unrelated dependency churn beyond `vault/api`'s own transitive requirements.

**New component**, `internal/components/vaultsigner/` (14 files, ~2,050 lines):
```
controllers/main.go          # Init(router, s) — builds the Vault client + parses PERSONAL_ENVS once, registers every route
handlers/handler.go          # Handler struct, shared httpError-based error response helper
handlers/signer_handler.go   # me/vault-signer/secrets*, all admin/vault-signer/* endpoints
handlers/personalenv_handler.go  # me/vault-signer/personal-envs*
handlers/audit.go            # writeAuditLog + small pointer helpers
handlers/swagger_docs.go     # swag annotations for all 11 routes (doc-stub pattern, matching stakeholder's convention)
models/db_models.go          # VaultSignerManagedSecret, VaultSignerAssignment, VaultSignerAuditLog
services/validation_service.go   # ValidateSignerValue / IsValidKeypair (Section 5c)
services/swap_service.go     # SwapSigner + SubmitSetOptionsTransaction + GatherActiveSigners (Section 5d)
services/signer_service.go   # CreateAssignment, ClaimPosition, EditAssignment, DeleteAssignment, owner identity resolution (Sections 1, 5e)
services/personalenv_service.go  # ResolveTrovoUsername, PERSONAL_ENVS parsing, key derivation, personal-secret orchestration (Section 5b)
vaultclient/client.go        # Vault API client bootstrap from VAULT_ADDR/VAULT_TOKEN
vaultclient/signer.go        # CSV get/rebuild/CAS-write, CollapseCSV, CountValidKeypairs (Section 5a)
vaultclient/personalenv.go   # Exists/Create/Delete (Section 5b) — no Get/Update, by design
```

**Modified existing files:**
- `internal/db/main.go` — the three GORM models registered in `migrateAdminSchemaTransaction`'s `AutoMigrate(...)` call, per Section 4's directive.
- `main.go` — `vaultsigner.Init(router, s)` wired in after `stakeholder.Init(router, s)`.
- `.env-sample` — `VAULT_ADDR`, `VAULT_TOKEN`, `PERSONAL_ENVS`, `PERSONAL_ENV_VAULT_MOUNT`, `PERSONAL_ENV_VAULT_PATH_PREFIX` added.
- `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml` — regenerated via `swag init` + `scripts/fix-swagger-yaml.sh`; contain all 11 vault-signer routes.
- `go.mod` / `go.sum` — the new dependency and its transitive requirements.

### Verification performed

- `go build ./...` — passes cleanly, whole repo, no errors.
- `go vet ./...` — passes cleanly (the only output is a pre-existing `xcode-select`/cgo notice from this sandbox lacking Xcode command-line tools, confirmed present on `main` before any of this work by stashing and re-running).
- `swag init` — succeeds, generates correct entries for all 11 routes and their DTOs (`putSecretRequest`, `createPersonalEnvRequest`, `registerManagedSecretRequest`, `ownerRequest`); verified by grepping the generated `swagger.json` for every expected path.
- `gofmt -l` — clean on every new/modified file.

### Verification **not** performed, and why

- **`go test ./...`** — fails to even build in this sandbox. Confirmed pre-existing on `main` (stashed this branch's changes and reproduced the identical failure): the repo's admin-DB tests use the `gorm.io/driver/sqlite` cgo driver, and this sandbox has no Xcode command-line tools installed. Not caused by this work, but it does mean **no automated test coverage has actually been run** for any of the new logic — only manual build/vet verification and careful tracing against the plan.
- **No live Vault, Postgres, or Horizon instance** in this sandbox — nothing here has been exercised end-to-end against real infrastructure. The AutoMigrate registration, the Vault KV v2 calls, and the Stellar transaction building/submission are all written to match their respective SDKs' real signatures (confirmed by reading the actual `vault/api`/`stellar/go` source in the module cache while writing this code, not guessed), but that's not a substitute for running them.
- **`golangci-lint`** — not installed in this sandbox; the repo's own `.golangci.yml`/`make lint` target was not run against this code.

### Two real bugs caught and fixed while implementing (not just planning)

1. **Spare-position on-chain removal.** The delete flow's on-chain trigger needed an explicit `index < ActiveSigningCount` gate — without it, deleting a spare position that happened to hold a valid key would attempt to remove a signer that was never registered on-chain to begin with, contradicting Section 4a/5d's own established active/spare distinction.
2. **Self-defeating count scoping.** An earlier attempt to restrict `validKeypairCount` to only the active range would have capped it at `ActiveSigningCount`, making the `> ActiveSigningCount` trigger condition permanently false and silently disabling on-chain removal for good. Caught by working through the arithmetic on a concrete example before shipping it, reverted to whole-CSV counting with the position gate kept as a separate, independent condition.

### Recommended before merging

1. Run this repo's real test suite and `make lint` on a machine with Xcode command-line tools installed (or in CI, which presumably has them).
2. Exercise the delete-then-reversal path (Section 5e) against a Horizon **testnet** sandbox specifically — it's the highest-risk, least-common code path (three independent systems have to interact in a specific order) and the one this plan repeatedly flags as needing real verification.
3. Confirm the wallet's actual on-chain threshold configuration matches the "active-signers-minus-one" assumption Section 9 flags — this was never verifiable from inside this sandbox.
4. Code review focused on `services/signer_service.go`'s `DeleteAssignment`/`reverseOnChainRemoval` — it's the most intricate control flow in the feature and the one place where getting an edge case wrong has real on-chain consequences.

---

## 12. Post-Launch Fix: One Assignment Per Owner Per Managed Secret

**Gap found:** as originally shipped (Sections 1–11), nothing prevented the same owner from holding more than one `VaultSignerAssignment` row on the same managed secret. `VaultSignerAssignment`'s only unique index was `(ManagedSecretID, Position)`; there was no constraint at all on `(ManagedSecretID, OwnerRefID)`, and neither `services.CreateAssignment` nor `services.EditAssignment` checked for an existing row before writing. An admin could `POST` the same `ownerIdentifier` onto a managed secret twice and get two separate assignment rows for that person.

**Why this mattered beyond a data-hygiene concern:** `handlers.loadOwnedSecret` — the lookup behind every `GET`/`PUT /me/vault-signer/secrets/:secretId` — resolves the caller's assignment with `Where("managed_secret_id = ? AND owner_ref_id = ? AND owner_type = ?", ...).First(...)`. With two rows for the same owner on the same secret, `First()` (no `ORDER BY`) picks one arbitrarily; the other becomes permanently unreachable through self-service, even if it had already claimed a live, on-chain-registered position no `PUT` could ever rotate again.

**Fix — guarded at both levels, per Section 4's own established pattern of "DB constraint is the real guarantee, app-level check is the friendly error":**

- **DB level:** added `idx_vault_signer_assignment_secret_owner`, a unique index on `(ManagedSecretID, OwnerRefID)`, via `gorm` tags on `models.VaultSignerAssignment` — picked up automatically by the existing `AutoMigrate(...)` registration (Section 4), no new migration mechanism needed. `OwnerType` was deliberately left out of this index: `AdminUser.ID` (a small `uint`, stringified) and `OrganizationMember.ID` (a UUID-format string) are different-enough formats that a collision between the two namespaces is not a practical concern, and the *stricter* two-column constraint is what actually closes the "same owner ID string, same secret" hole — a three-column version would let that exact case through if it ever occurred.
- **App level:** `services.CreateAssignment` now pre-checks for an existing `(ManagedSecretID, OwnerRefID)` row and returns `ErrOwnerAlreadyAssigned` before inserting; `services.EditAssignment` runs the same pre-check before reassigning to a new owner (excluding the assignment's own row — keeping the current owner is always a no-op save, never a self-conflict). Both also catch a Postgres unique-violation (`SQLSTATE 23505`) via a small `isUniqueViolation` helper (`github.com/jackc/pgx/v5/pgconn`) on the rare race that slips past the pre-check, translating it to the same sentinel error rather than letting a raw DB error reach the caller.
- **API level:** `handlers.CreateAssignment`/`handlers.EditAssignment` translate `ErrOwnerAlreadyAssigned` into `409 {"error": "owner_already_assigned", "message": "This owner already has an assignment on this managed secret."}` (Section 6).

**Signer Position numbering — confirmed unchanged, not touched by this fix:** `ClaimPosition` (Section 5e) already scoped its lowest-unclaimed-position search to a single `managed_secret_id` at a time, so positions have always run 1, 2, 3, ... independently per managed secret, never globally. This fix doesn't alter that numbering — it only guarantees a single owner can't occupy more than one of those numbers on the same secret, which also makes `Position` well-defined per `(ManagedSecretID, OwnerRefID)` now that the pair maps to at most one assignment row.

**Files touched:**
- `internal/components/vaultsigner/models/db_models.go` — new unique index tags + doc comment on `VaultSignerAssignment`.
- `internal/components/vaultsigner/services/signer_service.go` — `ErrOwnerAlreadyAssigned`, `isUniqueViolation`, and the pre-checks in `CreateAssignment`/`EditAssignment`.
- `internal/components/vaultsigner/handlers/signer_handler.go` — `409` translation in both handlers.
- `internal/components/vaultsigner/handlers/swagger_docs.go` — `@Failure 409` added to both `createAssignmentDocs`/`editAssignmentDocs`.
- `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml` — regenerated via `swag init` (pinned `v1.16.2`, matching `go.mod`) + `scripts/fix-swagger-yaml.sh`.
- `go.mod` — `github.com/jackc/pgx/v5` moved from indirect to direct (already an indirect dependency of `gorm.io/driver/postgres`; now imported directly for `pgconn.PgError`). `go.sum` unchanged.

**Verification performed:** `go build ./...` clean; `go vet ./internal/components/vaultsigner/...` clean; `gofmt -l` clean on every touched file; `swag init` regenerated only the expected `409` additions (checked via diff — 3 files, ~18 lines each, no unrelated churn).

**Verification not performed, and why:** same as Section 11 — no live Postgres in this sandbox, so the new unique index has not actually been exercised against a concurrent-write race, and `go test ./...` still fails to build here for the pre-existing, unrelated `sqlite` cgo reason. Before merging: run this against a real Postgres instance and confirm (a) the `AutoMigrate` successfully adds the new index (it will fail loudly instead of silently if any pre-existing duplicate `(ManagedSecretID, OwnerRefID)` rows exist — worth checking for those first if this deploys against a database that already has assignment data), and (b) a concurrent double-`POST` for the same owner actually resolves to one `201` and one `409`, not two `201`s.
