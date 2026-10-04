# Vendored: Candide Social Recovery Module v0.2.0

Unmodified copies of these files from
[candidelabs/candide-contracts](https://github.com/candidelabs/candide-contracts)
at commit `d0959d28ae084a7d549bb2c12d04482653456c2c` - the commit at which
Certora reviewed the fixes to its audit findings (released as v0.2.0):

- `modules/social_recovery/SocialRecoveryModule.sol`
- `modules/social_recovery/storage/GuardianStorage.sol`
- `modules/social_recovery/storage/IGuardianStorage.sol`
- `interfaces/ISafe.sol`, `IModuleManager.sol`, `IOwnerManager.sol`,
  `IFallbackManager.sol`, `IGuardManager.sol`

Audits (see the upstream `audit/` folder): Ackee Blockchain (v0.0.1),
Nethermind Security and Certora (safe-modules commit `8076191f`; Certora's
fixes reviewed at the commit above).

Dependencies as upstream: OpenZeppelin Contracts 4.9 (`SignatureChecker`)
and Safe Contracts 1.4.1 (`Enum`).

Licence: GPL-3.0 (see `LICENSE`). These files are deployed as their own
contract; nothing else in this repository is compiled into them.

Do not edit these files. To update, copy the files from a newer audited
commit and record it here.
