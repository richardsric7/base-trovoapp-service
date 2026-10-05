# recovery - integration

How the recovery module connects to the rest of Trovo. The module lives on
the Base blockchain; only app-backend talks to it, and the apps reach it
through app-backend.

| Connects to | Direction | What for |
|---|---|---|
| [app-backend](#1-app-backend) | app-backend → module | turn recovery on and off, start, cancel and finalize recoveries, watch for unexpected ones |
| [The guardian Safe](#2-the-guardian-safe) | guardian → module | the platform's recovery guardian sends starts and finalizations |
| [app-mobile and app-web](#3-app-mobile-and-app-web) | apps → app-backend | the user-facing recovery screens |
| [Alerts: Discord, push, email](#4-alerts) | app-backend → people | warnings about recoveries Trovo did not start |

---

## 1. app-backend

- **What it is and why:** Trovo's main backend. It is the only component
  that builds calls to the module: the users' wallet operations that turn
  recovery on or off or cancel a recovery, and the guardian's calls that
  start and finish one. Code: `internal/aa/recovery.go` (the calls) and
  `internal/components/users/services/recovery.go` (the flows). The full
  API is in [app-backend/INTEGRATION.md](../app-backend/INTEGRATION.md#account-recovery-opt-in-guardian).
- **Direction:** app-backend → module (through users' wallets and the
  guardian), and app-backend reads the module's events.
- **How they connect:** JSON-RPC to a Base node; users' operations go
  through the bundler like any other wallet operation.
- **What is sent, by whom:**

| Step | Who sends it | Calls |
|---|---|---|
| Turn on | The user's wallet (one operation per covered wallet, user pays gas) | `Safe.enableModule(module)` (unless enabled), `module.addGuardianWithThreshold(guardian, 1)`; the primary wallet's also pays the fee |
| Turn off | The user's wallet (one per covered wallet) | `module.revokeGuardianWithThreshold(0x1, guardian, 0)` |
| Add approvers to a covered wallet | The wallet, in the shared access operation | the owner changes, then `revokeGuardianWithThreshold` |
| Start a recovery | The guardian (pays gas) | `module.confirmRecovery(wallet, owners with the new key, threshold, nonce, true)` per covered wallet |
| Cancel | The user's wallet, current key (one per wallet being recovered) | `module.cancelRecovery()` |
| Finalize | The guardian (pays gas), after the period | `module.finalizeRecovery(wallet)` |

  When the guardian is a Safe, its calls for all of a user's wallets are
  batched into one Safe transaction.
- **What app-backend records:** `user_account_recovery_logs` (one per recovery: old and new key,
wallets, status `PENDING` / `CANCELED` / `COMPLETED` / `FAILED`, when it
takes effect, the guardian's transactions) and `recovery_watch_cursors`
(the last block the watcher scanned).
- **Settings on the other side (app-backend):**

  | Parameter | Example | How to get it |
  |---|---|---|
  | `RECOVERY_MODULE_ADDRESS` | `0xB7f8BC63BbcaD18155201308C8f3540b07f84F5e` | printed by the deploy script; `recoveryModule` in `contracts/deployments/<chainId>.json` |
  | `RECOVERY_PERIOD_SECONDS` | `604800` | `recoveryPeriodSeconds` in the same file |
  | `ACCOUNT_RECOVERY_GUARDIAN_SAFE` | `0xA11CE00000000000000000000000000000000002` | the guardian Safe ([DEPLOYMENT.md step 5](DEPLOYMENT.md#5-create-the-guardian)) |
  | `ACCOUNT_RECOVERY_GUARDIAN_SIGNERS` | `0x59c6...;0x5de4...` | the guardian's online keys, `;` separated |

  Each is described in [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md).
- **How to check it works:** turn recovery on for a test account in the
  app; basescan shows the wallet's operation calling
  `addGuardianWithThreshold` on the module.
- **When it is not configured:** the recovery endpoints answer
  `error-account-recovery-not-configured` (503); nothing else is affected.

## 2. The guardian Safe

- **What it is and why:** a Safe wallet owned by dedicated recovery keys.
  It is the only guardian of each covered wallet, so it is what starts and
  finishes recoveries. It is not an owner of any user wallet and cannot
  move funds.
- **Direction:** guardian → module, with transactions app-backend builds
  and signs with enough of `ACCOUNT_RECOVERY_GUARDIAN_SIGNERS` to meet the
  Safe's threshold (the first key pays the network fee).
- **Settings:** `ACCOUNT_RECOVERY_GUARDIAN_SAFE` and
  `ACCOUNT_RECOVERY_GUARDIAN_SIGNERS` in app-backend.
- **How to check it works:** the guardian's online keys hold a little ETH,
  and a test recovery on the test network starts and finishes.
- **When its keys run out of ETH:** recoveries cannot start or finish
  until they are topped up.

## 3. app-mobile and app-web

- **app-mobile:** turning recovery on and off signs every transaction
  app-backend returns (`transactions` → `transactionSignatures`). The
  account recovery page has *Cancel a recovery in progress*, which recovery
  push notifications (`route: accountRecovery`) open. After recovering, the
  new device shows the recovery as in progress and checks its status until
  it completes.
- **app-web:** recovering on the web shows the recovery as in progress with
  the time it takes effect and checks its status, then sends the user to
  import the account with the new key.
- **Settings:** none specific to recovery; the apps only need app-backend's
  address (see their CONFIGURATION.md).

## 4. Alerts

- **What it is and why:** someone holding the guardian's keys could start a
  recovery on their own. app-backend watches the module's
  `RecoveryExecuted` events, and for any recovery it did not start it
  alerts the team on Discord and sends the wallet's owner an urgent push
  notification and email, so they can cancel within the recovery period.
- **Settings:** app-backend's Discord webhook, Firebase (`GC`) and email
  settings (see [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md)).

## Contract details (v0.2.0)

- `confirmRecovery` takes the module's current nonce for the wallet
  (`nonce(wallet)`); `cancelRecovery` increments it.
- `finalizeRecovery` adds the new owners before removing the old ones, so
  a wallet's owner order can change; compare owner sets, not lists.
- `executeRecovery` (signatures from several guardians) is not used: the
  platform guardian is a wallet's only guardian.
