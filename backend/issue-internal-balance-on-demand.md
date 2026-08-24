# Replace DEX 1:1 swap with on-demand InternalBalance issuance in tokenized-asset purchase

## Context

When a user buys a tokenized asset today, `SubscribeToTokenizedAsset` builds a single `PathPaymentStrictSend` operation that sends the buyer's stablecoin (e.g. CNGN) and lets Horizon's path-finder route it through the DEX order book — first through an assumed 1:1 CNGN↔InternalBalance offer, then through the asset's market-maker offer (InternalBalance↔TokenizedAsset, created at mint time) — to land the tokenized asset in the buyer's wallet. This depends on a real 1:1 DEX offer existing for CNGN↔InternalBalance, which is fragile/indirect.

The new design removes that dependency: the stablecoin the buyer sends is redirected to a dedicated `FundsHoldingWalletPublicKey` (per tokenized asset), and the InternalBalance token is minted on-demand 1:1 directly to the buyer by the InternalBalance issuer, all within the same atomic Stellar transaction as the InternalBalance→TokenizedAsset swap. The InternalBalance trustline is opened and closed within that same transaction so the buyer's wallet never holds a lingering InternalBalance balance/trustline outside of the app's control.

Per the user's explicit decision: if a tokenized asset doesn't yet have `FundsHoldingWalletPublicKey` set, purchases against it must be rejected with a clear config error — there is no fallback to the old DEX-swap path. This keeps exactly one purchase mechanism live.

## 1. Add `FundsHoldingWalletPublicKey` field

File: `internal/components/users/models/tokenization.go`

Follow the exact pattern already used for `WalletToHoldAssetsNotForSale` (closest existing analog — nullable wallet-public-key field):

- **`TokenizedAsset` struct** (near line 150): add `FundsHoldingWalletPublicKey *string `json:"fundsHoldingWalletPublicKey"`` — nullable, no `gorm` tag needed (GORM `AutoMigrate`, already wired against `TokenizedAsset` in `internal/db/main.go:191`, will add the column on next startup).
- **`TokenizedAssetJSONInput` struct** (near line 675): add `FundsHoldingWalletPublicKey string `json:"fundsHoldingWalletPublicKey"``.
- **`TokenizedAssetJSON` struct** (near line 1226): add `FundsHoldingWalletPublicKey string `json:"fundsHoldingWalletPublicKey"``.
- **`TokenizedAsset.UpdateTokenizedAssetFromInput`** (near line 2785-2790): mirror the `WalletToHoldAssetsNotForSale` if/else-nil block:
  ```go
  if len(ti.FundsHoldingWalletPublicKey) > 0 {
      t.FundsHoldingWalletPublicKey = &ti.FundsHoldingWalletPublicKey
  } else {
      t.FundsHoldingWalletPublicKey = nil
  }
  ```
- **`TokenizedAsset.ToJSON`** (near line 5257-5259): mirror the pointer-deref block:
  ```go
  if ti.FundsHoldingWalletPublicKey != nil {
      t.FundsHoldingWalletPublicKey = *ti.FundsHoldingWalletPublicKey
  }
  ```
- **`TokenizedAsset.UpdateCalculation`** (`tokenization.go:4793-4978`): no change — this method only does financial/fee arithmetic and never touches wallet fields.

## 2. Parse `INTERNAL_BALANCE_ISSUING_SIGNERS`

New helper in `internal/components/users/services/tokenized_assets.go` (near `generateAssetSubscriptionXdr`), following the CSV-of-secret-seeds pattern already used for `CHANNEL_ACCOUNTS` in `main.go:351-380` (`strings.Split` + `keypair.ParseFull` + trim), rather than the single-secret `keypair.MustParseFull(os.Getenv(...))` pattern used elsewhere:

```go
func getInternalBalanceIssuingSigners() ([]*keypair.Full, error) {
    raw := os.Getenv("INTERNAL_BALANCE_ISSUING_SIGNERS")
    var signers []*keypair.Full
    for _, v := range strings.Split(raw, ",") {
        v = strings.TrimSpace(v)
        if len(v) == 0 {
            continue
        }
        kp, e := keypair.ParseFull(v)
        if e != nil {
            return nil, e
        }
        signers = append(signers, kp)
    }
    if len(signers) == 0 {
        return nil, errors.New("INTERNAL_BALANCE_ISSUING_SIGNERS not configured")
    }
    return signers, nil
}
```

All parsed signers sign the transaction (extra valid signatures beyond the multisig threshold are harmless on Stellar).

## 3. Rework `generateAssetSubscriptionXdr`

File: `internal/components/users/services/tokenized_assets.go:2920-3111`

**Signature change:** add `ta *userModels.TokenizedAsset` as a parameter (it already exists in the caller `SubscribeToTokenizedAsset`, just not threaded through). Update the call site at line 2774: `generateAssetSubscriptionXdr(subscriberWallet, ta, &swapInfo, gc)`.

**New up-front validation** (mirrors the nil-checks already used in `generateMintRegulatedTokenizedAssetXdr` for `AssetQuoteCurrency`/`AssetCountryLocation`, lines 3316-3350):
- If `ta.FundsHoldingWalletPublicKey == nil` → return `&tErrors.CustomError{Param: "fundsHoldingWalletPublicKey", Err: "error-funds-holding-wallet-not-set", ErrMessage: "This tokenized asset is not yet configured to accept purchases."}`.
- Resolve `countryConfig := userModels.CountryCode(*ta.AssetCountryLocation).GetConfig(gc)`; if `InternalBalanceTokenCode`/`InternalTokenIssuer` are nil, return the same kind of config error used at mint time.
- Build `internalBalanceAsset := txnbuild.CreditAsset{Code: *countryConfig.InternalBalanceTokenCode, Issuer: *countryConfig.InternalTokenIssuer}`.

**Replace the operation-building block (current lines 2966-3025)** with the 7-op sequence, in this order (mirrors the user's spec exactly):

1. `txnbuild.Payment` — `SourceAccount: wallet.ID`, `Destination: *ta.FundsHoldingWalletPublicKey`, `Asset: sourceAsset` (the stablecoin), `Amount: swapInfo.SourceAmount`.
2. `txnbuild.ChangeTrust` — internal balance asset, `SourceAccount: wallet.ID`, `Limit: gc.TokenLimitAsString()` — **only if** the buyer's wallet doesn't already trust it (check via `network.BlockchainAccountProperties(client, wallet.ID, internalBalanceAsset)`, same pattern as the existing `sourceAccountTrustsDestinationAsset` check).
3. `txnbuild.SetTrustLineFlags` — `Trustor: wallet.ID`, internal balance asset, `SetFlags: [TrustLineAuthorized]`, `SourceAccount: *countryConfig.InternalTokenIssuer` — same conditional as step 2.
4. `txnbuild.Payment` — `SourceAccount: *countryConfig.InternalTokenIssuer`, `Destination: wallet.ID`, `Asset: internalBalanceAsset`, `Amount: swapInfo.SwapAmount` (1:1 with the stablecoin amount sent in step 1).
5. Existing tokenized-asset trustline block, unchanged logic, just moved after step 4: `ChangeTrust` + `SetTrustLineFlags` for the destination (tokenized) asset, still gated on `!sourceAccountTrustsDestinationAsset`, still signed later by `TOKENIZATION_ISSUING_PROFILE_WALLET`.
6. `txnbuild.PathPaymentStrictSend` — `SourceAccount: wallet.ID`, `SendAsset: internalBalanceAsset`, `SendAmount: swapInfo.SwapAmount`, `DestAsset: destinationAsset`, `Destination: wallet.ID`, `DestMin: swapDestMin.String()`, `Path` from `GetStrictSendPaths(...)` called with `SourceAssetCode/Issuer` = the internal balance token instead of the stablecoin (reuses the existing `GetStrictSendPaths` helper and its liquidity-error handling — this is now the only leg that still needs DEX path-finding, since step 1→4 replaced the CNGN↔InternalBalance leg with direct issuance).
7. `txnbuild.ChangeTrust` — internal balance asset, `SourceAccount: wallet.ID`, `Limit: "0"` (removes the trustline; safe because step 6 spends the entire internal-balance balance within the same atomic transaction, so the balance is zero by the time this op executes).

**Signing (after `tx, err = txnbuild.NewTransaction(...)`):**
- Channel account signs, same as today, only `if swapInfo.Multiparty == 1`.
- `TOKENIZATION_ISSUING_PROFILE_WALLET` signs, same as today, only `if !sourceAccountTrustsDestinationAsset` (covers step 5's `SetTrustLineFlags`).
- **New:** always call `getInternalBalanceIssuingSigners()` and sign with all returned keypairs — `tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), signers...)` — to satisfy steps 3 and 4's `SourceAccount` (the InternalBalance issuer, a multisig account). Return `&tErrors.ErrorTemporaryServerError{}` (logged) if parsing fails, consistent with how other missing-signer config is handled.

No other part of `SubscribeToTokenizedAsset` needs to change: KYC checks, cap checks, market-offer liquidity checks (still against the same tokenized-asset `ManageSellOffer` created at mint time, unaffected by this change), the multiparty/`PendingAuth` branch, and submission via `network.SubmitXdrWithSignature` all stay as-is.

## Verification

- `go build ./...` (or `go vet ./...`) from `/Users/ric/trovo-wallet-monorepo/backend` to confirm the new field, signature change, and new operations compile cleanly across all three call sites of `SubscribeToTokenizedAsset` (`servicelinks/controllers/handlers_impl.go:4945`, `users/controllers/handlers_impl.go:6224` and `:6385`).
- Manually inspect the generated XDR (e.g. via `stellar-xdr` decode or Horizon's transaction-preview) for a test purchase to confirm the 7 operations appear in order with the expected source accounts/assets/amounts, and that `swapInfo.SwapAmount` matches between steps 1, 4, and 6.
- Exercise the purchase flow end-to-end against testnet: buy a tokenized asset (a) whose `FundsHoldingWalletPublicKey` is unset — expect the new config error — and (b) whose `FundsHoldingWalletPublicKey` is set — expect a successful purchase, confirm the buyer's wallet ends with no InternalBalance trustline and the correct tokenized-asset balance, and confirm the stablecoin landed in the holding wallet.
- Confirm `INTERNAL_BALANCE_ISSUING_SIGNERS` missing/malformed produces a clear server error rather than a panic (the CSV parser returns an error rather than using `MustParseFull`).
