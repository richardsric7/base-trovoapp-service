# recovery - integration

## app-backend

app-backend is the only component that talks to the module
(`internal/aa/recovery.go` for the calls, `internal/components/users/services/recovery.go`
for the flows). The full API is in
[app-backend/INTEGRATION.md](../app-backend/INTEGRATION.md#account-recovery-opt-in-guardian).

| Step | Who sends it | Calls |
|---|---|---|
| Turn on | The user's wallet (one operation per covered wallet, user pays gas) | `Safe.enableModule(module)` (unless enabled), `module.addGuardianWithThreshold(guardian, 1)`; the primary wallet's also pays the fee |
| Turn off | The user's wallet (one per covered wallet) | `module.revokeGuardianWithThreshold(0x1, guardian, 0)` |
| Add approvers to a covered wallet | The wallet, in the shared access operation | the owner changes, then `revokeGuardianWithThreshold` |
| Start a recovery | The guardian (pays gas) | `module.confirmRecovery(wallet, owners with the new key, threshold, nonce, true)` per covered wallet |
| Cancel | The user's wallet, current key (one per wallet being recovered) | `module.cancelRecovery()` |
| Finalize | The guardian (pays gas), after the period | `module.finalizeRecovery(wallet)` |

The guardian's calls are batched into one Safe transaction
(MultiSendCallOnly) when the guardian is a Safe.

Records: `user_account_recovery_logs` (one per recovery: old and new key,
wallets, status `PENDING` / `CANCELED` / `COMPLETED` / `FAILED`, when it
takes effect, the guardian's transactions) and `recovery_watch_cursors`
(the last block the watcher scanned).

## Apps

- **app-mobile** - turning recovery on and off signs every returned
  transaction (`transactions` -> `transactionSignatures`); the account
  recovery page has *Cancel a recovery in progress*, which recovery push
  notifications (`route: accountRecovery`) open; after recovering, the new
  device shows the recovery as in progress and polls its status until it
  completes.
- **app-web** - recovering on the web shows the recovery as in progress
  with the time it takes effect and polls its status, then sends the user
  to import the account with the new key.

## Contract interface notes (v0.2.0)

- `confirmRecovery` takes the module's current nonce for the wallet
  (`nonce(wallet)`); `cancelRecovery` increments it.
- `finalizeRecovery` adds the new owners before removing the old ones, so
  a wallet's owner order can change; compare owner sets, not lists.
- `executeRecovery` (signatures from several guardians) is not used: the
  platform guardian is a wallet's only guardian.
