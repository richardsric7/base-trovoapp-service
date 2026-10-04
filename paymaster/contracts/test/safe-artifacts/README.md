# Vendored Safe artifacts (tests only)

ABI + creation bytecode of the Safe contracts the Trovo wallets use, so the
tests can run real Safes with the ERC-4337 module on a local chain. Only
`abi` and `bytecode` are kept from each build artifact.

| File | Source package | License |
|---|---|---|
| `SafeL2.json`, `SafeProxyFactory.json`, `MultiSendCallOnly.json` | `@safe-global/safe-contracts@1.4.1-build.0` (`build/artifacts`) | LGPL-3.0 |
| `Safe4337Module.json`, `SafeModuleSetup.json` | `@safe-global/safe-4337@0.3.0-1` (`build/artifacts`) | GPL-3.0 |

The SafeL2 and SafeProxyFactory runtime code deployed from these artifacts
hash to the canonical v1.4.1 code hashes published in
`@safe-global/safe-deployments` (`0xb1f92697…81ff` and `0x50c3cdc4…3317`).

They are vendored rather than installed because those packages' peer
dependencies conflict with this project's Hardhat toolbox.
