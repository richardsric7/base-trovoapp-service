# Plan: Country Internal Balance Token replacing CNGN as tokenized-asset quote currency

*(quote currency resolved directly from `CountryConfig`; purchase and dividend-payout flows default to CNGN so existing clients keep working unchanged; distribution-wallet trustline auto-created/authorized at mint time via a dedicated co-signer; already-minted CNGN-quoted assets grandfathered)*

## What's already in the codebase (confirmed by reading it)

- [CountryConfig](internal/components/users/models/country_config.go) is already in the `AutoMigrate` list in [internal/db/main.go](internal/db/main.go#L159), keyed by `CountryCode`. Adding two new columns is enough — `AutoMigrate` adds them automatically. **(Done — see below.)**
- [`CountryCode.GetConfig(gc)`](internal/components/users/models/country_config.go#L46) already loads the full `CountryConfig` row by country code, and is already called at [SubmitTokenizationAssetInfoByInitiator](internal/components/users/services/tokenized_assets.go#L790) for the TROV-balance check — reused as the single source of truth for the quote currency everywhere.
- `TokenizationCurrency` ([internal/components/users/models/tokenization.go](internal/components/users/models/tokenization.go#L1708)) stays in use as the registry of acceptable **stablecoins** — for purchase payment (item 5) and dividend payout (item 2b) — it's only removed from the **`AssetQuoteCurrency`** resolution path.
- `ProceedPayoutCurrency` is the currency the **dividend engine** ([payouts.go](internal/components/users/services/payouts.go#L129), models at `ProceedPayout`/`TokenizedAssetPayoutSchedule`/`TokenizedAssetPayoutEngineTask` in [tokenization.go](internal/components/users/models/tokenization.go#L2054)) pays token holders in — distinct from `TokenizedAssetEarlyExit.PayoutCurrency` ([tokenization.go](internal/components/users/models/tokenization.go#L6735)), a separate field for early-exit/liquidation payouts. `payouts.go` looks like an unfinished/scratch script (placeholder `KGM`/`ISSUER_PUBLIC_KEY` constants, an unused `func main()`) rather than a wired-in production path — confirm its actual production equivalent before implementing item 2b.
- The rest of the tokenization pipeline is threaded through `AssetQuoteCurrency`:
  - Creation validation: [SubmitTokenizationAssetInfoByInitiator](internal/components/users/services/tokenized_assets.go#L714), [SubmitTokenizationAssetInfo](internal/components/users/services/tokenized_assets.go#L818)
  - Mint: [generateMintRegulatedTokenizedAssetXdr](internal/components/users/services/tokenized_assets.go#L3264) — resolves the quote currency at [L3296](internal/components/users/services/tokenized_assets.go#L3296), builds the `ManageSellOffer` with `Buying = quoteCurrency` at [L3530](internal/components/users/services/tokenized_assets.go#L3530), gates on [checkDistributionWalletHasQuoteCurrencyAuthorization](internal/components/users/services/tokenized_assets.go#L4030) (today hardcoded to `NAIRA_ASSET`).
  - Purchase: [SubscribeToTokenizedAsset](internal/components/users/services/tokenized_assets.go#L2620) — resolves quote currency at [L2701](internal/components/users/services/tokenized_assets.go#L2701), sets `swapInfo.SourceAssetCode = ta.AssetQuoteCurrency` today.
  - Early exit: [buildTokenizedAssetEarlyExit](internal/components/users/services/tokenized_assets.go#L4164) sets `ee.PayoutCurrency = *ta.AssetQuoteCurrency`.
  - Swap/path-payment already uses Stellar strict-send pathfinding ([generateAssetSubscriptionXdr](internal/components/users/services/tokenized_assets.go#L2888) → `GetStrictSendPaths`), which **already supports multi-hop routing** — "CNGN → NGN → AssetCode" works automatically once we stop forcing `SourceAssetCode == AssetQuoteCurrency`.
- [Pay()](internal/components/users/services/pay.go#L34) is the single choke point for wallet-to-wallet transfers, shared by payments, service-links, and callbacks controllers.
- [SwapSend/SwapReceive](internal/components/swaps/services/swap.go#L32) can deliver an asset to an arbitrary `DestinationAccount` ([L959](internal/components/swaps/services/swap.go#L959)) — a second path to guard, independent of the tokenization purchase flow (which builds its own ops and never calls these).
- `CanWithdraw` on `BantuAsset` ([internal/components/users/models/assets.go](internal/components/users/models/assets.go#L222)) exists but isn't wired into anything today — external withdrawal needs an explicit check.
- Shared-access / multi-party users aren't a separate code path anywhere here — `/v1/tokenization/subscriptions/:id` and `/v1/shared-access/tokenization/subscriptions/:id` both call `SubscribeToTokenizedAsset`, mint-approval is shared via `ApproveTransaction` in [approvals.go](internal/components/users/services/approvals.go), and `Pay()` already branches on `Multiparty` internally — no separate work item needed there.

---

## 1. Data model — `CountryConfig` ✅ done

```go
InternalBalanceTokenCode *string `gorm:"size:12;default:null" json:"internalBalanceTokenCode"`
InternalTokenIssuer      *string `gorm:"size:68;default:null" json:"internalTokenIssuer"`
```
Nullable pointer fields, `AutoMigrate`-safe against existing rows. No changes needed in [internal/db/main.go](internal/db/main.go).

Still needed: package-level `IsInternalBalanceAsset(code, issuer string, gc)` helper (reverse lookup across all countries) for the Pay/Swap guards (item 7) and for rejecting the internal token as a stablecoin choice in items 5/2b.

## 2a. `AssetQuoteCurrency` resolution — derive from `CountryConfig`, only for new tokenizations ✅ done

In [SubmitTokenizationAssetInfoByInitiator](internal/components/users/services/tokenized_assets.go#L714) and [SubmitTokenizationAssetInfo](internal/components/users/services/tokenized_assets.go#L818): replace the `GetTokenizationCurrencyByCode(input.AssetQuoteCurrency, ...)` validation block (around L725, L831) with `countryConfig := userModels.CountryCode(*ato.AssetCountryLocation).GetConfig(gc)` (reusing the call already made at L790), then:
- error if `countryConfig.InternalBalanceTokenCode == nil`,
- force `ato.AssetQuoteCurrency = countryConfig.InternalBalanceTokenCode`, ignoring any client-supplied value.

Only runs pre-mint, so already-minted/live tokenized assets keep their CNGN `AssetQuoteCurrency` untouched — no backfill needed.

## 2b. `ProceedPayoutCurrency` (dividend payout) — stays a stablecoin, defaults to CNGN (non-breaking) ✅ done (creation-time validation only — see note below)

Dividend payout uses stablecoin, not the internal balance token. Remove the current default-to-`AssetQuoteCurrency` behavior (around L735-739 / L839-846) — since `AssetQuoteCurrency` is now the internal token, blindly defaulting to it would wrongly pull dividends into the internal token too. Instead, mirror item 5's pattern:
- if `input.ProceedPayoutCurrency == ""`, default it to `"CNGN"` and resolve the issuer via `GetTokenizationCurrencyByCode("CNGN", gc.DB)` — reproduces today's practical behavior,
- if provided, validate it against `TokenizationCurrency` as today (must be a registered stablecoin), and explicitly reject it if it equals the internal balance asset (reuse `IsInternalBalanceAsset` from item 1).

Confirm the actual production dividend-payout path (vs. the `payouts.go` script, which looks unfinished) before implementing, so this change lands wherever dividends are genuinely computed/paid.

**Status**: implemented in `SubmitTokenizationAssetInfoByInitiator` and `SubmitTokenizationAssetInfo` — `ProceedPayoutCurrency` now defaults to CNGN, is validated against `TokenizationCurrency`, and is rejected if it resolves to the internal balance asset. The actual dividend-payout *engine* (`payouts.go` vs. its real production equivalent) was **not** tracked down or touched — still pending confirmation per the note above.

## 3. Mint — resolve quote currency + issuer from `CountryConfig` ✅ done

In [generateMintRegulatedTokenizedAssetXdr](internal/components/users/services/tokenized_assets.go#L3264), replace the `GetTokenizationCurrencyByCode(*t.AssetQuoteCurrency, gc.DB)` call (around L3296) with `countryConfig := userModels.CountryCode(*t.AssetCountryLocation).GetConfig(gc)`, erroring if `InternalBalanceTokenCode`/`InternalTokenIssuer` is nil. The `ManageSellOffer` (around L3530-3531) then buys `*countryConfig.InternalBalanceTokenCode`/`*countryConfig.InternalTokenIssuer` directly.

Already-minted assets are grandfathered (item 2a), so this path only runs for tokenizations created with the internal token as quote currency.

## 4. Distribution wallet must hold/trust the internal balance token — auto-created and auto-authorized at mint time ✅ done

The internal token's **issuer** is per-country: `*countryConfig.InternalTokenIssuer` (from `CountryConfig`, resolved in item 3). `INTERNAL_BALANCE_AUTHORIZER_WALLET` is **not** that issuer — it's just one signer *on* the issuer account (a multi-sig co-signer used specifically to authorize trustlines), parsed the same way `TOKENIZATION_ISSUING_PROFILE_WALLET` already is (see [tokenized_assets.go](internal/components/users/services/tokenized_assets.go#L3367)):
```go
authorizerKP := keypair.MustParseFull(strings.TrimSpace(os.Getenv("INTERNAL_BALANCE_AUTHORIZER_WALLET")))
```
No public-key-equality check against `InternalTokenIssuer` — they're expected to differ, since `authorizerKP` is a co-signer, not the account itself.

In `generateMintRegulatedTokenizedAssetXdr`, alongside the existing asset-code trustline ops (around L3446-3476), add the mirrored pair for the quote currency:
```go
ops = append(ops, &txnbuild.ChangeTrust{
    Line:          txnbuild.ChangeTrustAssetWrapper{Asset: txnbuild.CreditAsset{Code: *countryConfig.InternalBalanceTokenCode, Issuer: *countryConfig.InternalTokenIssuer}},
    Limit:         gc.TokenLimitAsString(),
    SourceAccount: distributionWallet.ID,
})
ops = append(ops, &txnbuild.SetTrustLineFlags{
    Trustor:       distributionWallet.ID,
    Asset:         txnbuild.CreditAsset{Code: *countryConfig.InternalBalanceTokenCode, Issuer: *countryConfig.InternalTokenIssuer},
    SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized},
    SourceAccount: *countryConfig.InternalTokenIssuer,
})
```
`SourceAccount` on `SetTrustLineFlags` is the actual issuer (`*countryConfig.InternalTokenIssuer`), but the transaction gets **signed** with `authorizerKP` (the co-signer) rather than the issuer account's own master key — wire that signature into whatever step signs/co-signs this mint transaction alongside the fee wallet and tokenization-issuer-profile keys. This assumes the issuer account's signing threshold accepts `authorizerKP` alone for this operation (i.e. it's already set up as a valid signer with sufficient weight on that account) — an operational/multi-sig-configuration prerequisite on the issuer account itself, not something this code sets up.

[checkDistributionWalletHasQuoteCurrencyAuthorization](internal/components/users/services/tokenized_assets.go#L4030) still runs as a pre-flight check, generalized to accept `(code, issuer string, distributionWallet, gc)`, dropping the CNGN/`NAIRA_ASSET` special-casing — catches the case where a prior mint attempt's trustline op failed/wasn't submitted.

## 5. Subscription/purchase — let the user pay in any approved stablecoin, defaulting to CNGN (non-breaking) ✅ done

Today [SubscribeToTokenizedAsset](internal/components/users/services/tokenized_assets.go#L2620) resolves the quote currency issuer at L2701 and sets `swapInfo.SourceAssetCode = ta.AssetQuoteCurrency` (L2699).

- Add `PaymentAssetCode` / `PaymentAssetIssuer` (optional) to `TokenizedAssetSubscriptionInput` ([tokenization.go](internal/components/users/models/tokenization.go#L1970)) (and the service-link mirror struct).
- Replace the L2699-2702 block:
  - if `input.PaymentAssetCode == ""`, default to `"CNGN"`, resolve via `GetTokenizationCurrencyByCode("CNGN", dbTX)` — reproduces today's exact behavior for unchanged clients,
  - if provided, use it as the code, always re-resolve/validate the issuer server-side via `GetTokenizationCurrencyByCode(input.PaymentAssetCode, dbTX)` — a client-supplied issuer is never trusted verbatim; mismatch is rejected.
  - reject if the resolved payment asset is the internal balance asset (`IsInternalBalanceAsset`) or `ta.AssetCode`.
- `swapInfo.DestinationAssetCode/Issuer` stays `ta.AssetCode`/`ta.IssuingWalletPublicKey`. For internal-token-quoted assets, NGN is only the intermediate DEX hop (requires the 1:1 CNGN/NGN offer to exist on-chain — ops prerequisite). For grandfathered CNGN-quoted assets, paying in CNGN is direct, unchanged.
- Cap/liquidity checks (around L2710-2740) read the market offer directly — no change needed.

## 6. Early exit — no core logic change, stays on the internal balance token

[buildTokenizedAssetEarlyExit](internal/components/users/services/tokenized_assets.go#L4164) already sets `PayoutCurrency = *ta.AssetQuoteCurrency` — once item 2a lands, early-exit payouts automatically denominate in the internal balance token for new tokenizations (CNGN for grandfathered ones), matching the intent that liquidation/exit proceeds sit in the internal token pending a future redemption tool. Unlike dividends (item 2b), this is unchanged.

## 7. Enforce "cannot send the internal balance asset" ✅ done

- **[Pay()](internal/components/users/services/pay.go#L34)**: guard near the top checking `paymentInfo.AssetCode`/`AssetIssuer` against `IsInternalBalanceAsset(...)`, before `Multiparty` branching — covers payments, service-links, and callbacks controllers, and shared-access payments, in one change.
- **[SwapSend/SwapReceive](internal/components/swaps/services/swap.go#L32)**: same guard on `Source`/`DestinationAssetCode/Issuer`.
- **External/crypto withdrawal**: explicit check at the withdrawal-request entry point in [stablerail_withdrawal.go](internal/components/users/services/stablerail_withdrawal.go) — uses a new `IsInternalBalanceAssetCode(code, gc)` code-only variant (added alongside `IsInternalBalanceAsset`) since the Stablerail request only carries a ticker, no issuer.

## 8. Migration/rollout (data, not code)

- Populate `InternalBalanceTokenCode`/`InternalTokenIssuer` per `CountryConfig` row.
- Set `INTERNAL_BALANCE_AUTHORIZER_WALLET` in the environment, and ensure it's configured as a valid signer (with sufficient threshold weight) on each country's `InternalTokenIssuer` account.
- `CNGN` must remain a valid `TokenizationCurrency` row (hardcoded default for both purchase and dividend payout, and still used by grandfathered assets).
- Ops: create the 1:1 DEX offers (CNGN/NGN, etc.) on-chain.
- Already-minted CNGN-quoted assets: left as-is, no migration.

---

## Implementation order

1. Data model (done)
2. Currency resolution (2a, 3) — done
3. Mint / distribution wallet trustline (4) — done
4. Subscribe, CNGN-default (5) — done
5. Dividend payout, CNGN-default (2b) — creation-time validation done; real payout *engine* still pending confirmation
6. Pay/Swap/withdrawal guards (7) — done

**Remaining work**: item 6 (early exit) needed no code change, as noted in its section — it inherits the new behavior automatically once 2a landed. The only genuinely open item is confirming the production dividend-payout engine (vs. the seemingly-unfinished `payouts.go`) and wiring `ProceedPayoutCurrency` through it if that engine does anything with the currency beyond what's already validated at creation time. Item 8 (migration/rollout) is ops work, not code, and is unstarted.
