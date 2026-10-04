# recovery

Opt-in account recovery for Trovo wallets, on
[Candide's Social Recovery Module](contracts/src/candide/README.md)
**v0.2.0**, vendored unmodified (GPL-3.0) and deployed by Trovo.

Trovo is non-custodial: users keep their secret key and the platform holds
no key that controls their wallets. A user who opts in to recovery adds one
narrow exception: the platform's **recovery guardian** may *replace the
key* on their wallets, after a waiting period, never move funds.

| | |
|---|---|
| Contract | [`SocialRecoveryModule`](contracts/src/candide/modules/social_recovery/SocialRecoveryModule.sol) - a Safe module; each wallet enables it and chooses its guardians. |
| Guardian | The platform's recovery guardian: a Safe owned by dedicated recovery keys (recommended) or one key (app-backend `ACCOUNT_RECOVERY_GUARDIAN_*`). The only guardian of each covered wallet, threshold 1. |
| Covered wallets | A user's primary wallet and their own sub-wallets without approvers. Wallets with co-signers are recovered by their co-signers instead. |
| Recovery period | Fixed when the module is deployed (e.g. 7 days). |

## How it works

1. **Opt in** - the user signs one operation per wallet: enable the module
   and add the guardian (the primary wallet's operation also pays the
   recovery fee). The user pays the gas, like any of their operations.
2. **Recover** - on a new device, after the email OTP and security
   answers, app-backend has the guardian call `confirmRecovery` on each
   covered wallet with the new key in place of the old one. The user is
   notified (push and email).
3. **Waiting period** - during the recovery period nothing changes, and the
   wallet's current key can `cancelRecovery`. The apps offer this on the
   account recovery page.
4. **Finalize** - after the period anyone may call `finalizeRecovery`;
   app-backend does, and moves the account to the new key. Wallets shared
   with approvers get an approval request to swap the old key for the new
   one, which their other co-signers approve.

## Safety properties

- The guardian is not an owner of the wallet: it cannot sign transactions,
  transfer tokens or ETH, or change modules. The module only lets
  guardians replace the owner set, through the delayed request above.
- The guardian cannot make itself an owner (`SM: new owner cannot be
  guardian`), cannot cancel or change guardians (only the wallet can), and
  every request needs the module's current nonce for the wallet, so a
  cancelled request cannot be replayed.
- A wallet's owners can cancel any request during the period, and remove
  the guardian at any time (turning recovery off).
- Adding approvers to a covered wallet removes the guardian in the same
  operation, so the platform never replaces keys on wallets approvers
  control.
- app-backend watches the module's `RecoveryExecuted` events and alerts
  (Discord, and an urgent push and email to the owner) on any recovery it
  did not start.

**Residual risk**: someone holding the guardian's keys can *start* a
recovery of a covered wallet. It only succeeds if the owner does not cancel
it within the recovery period, which is why the period should be days, the
guardian should be a multi-key Safe, and the alerts above must reach
someone. The guardian keys should be held as carefully as any treasury key.

## Why v0.2.0, deployed by Trovo

Candide's publicly deployed recovery modules (abstractionkit 0.4.x
addresses) predate the v0.2.0 fixes to Certora's findings (one medium,
five low) - their `confirmRecovery` has no nonce argument. Trovo deploys
the audited v0.2.0 source itself. See [DEPLOYMENT.md](DEPLOYMENT.md).

## Documents

- [DEPLOYMENT.md](DEPLOYMENT.md) - build, test, deploy, local deployment.
- [CONFIGURATION.md](CONFIGURATION.md) - the deploy script's settings and
  what app-backend needs.
- [INTEGRATION.md](INTEGRATION.md) - how app-backend and the apps use it.
