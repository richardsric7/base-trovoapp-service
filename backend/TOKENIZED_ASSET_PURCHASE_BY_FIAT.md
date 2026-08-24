# Tokenized Asset Purchase by Fiat — Implementation Plan

Status: **PLANNING ONLY — not implemented.**
Author: Claude (research grounded in current state of `backend`, branch `ric-asset-tokenization-purchase`)
Related prior work: `backend/issue-internal-balance-on-demand.md`, commits `35d34192`, `771059a5`, `08691be0`, `3c53afd0`, `e784093c`.

## 1. Goal

Let a user buy a tokenized asset with fiat instead of an on-chain stablecoin, without the backend needing any further call from the client after payment starts. The buyer never touches Stellar directly: the server builds and signs (with every server-held key) a transaction that mints the country's InternalBalance token to the buyer's wallet and swaps it into the tokenized asset — the buyer only adds their own signature. Submission to the blockchain is deferred until Flutterwave confirms the fiat payment via webhook.

**One identifier drives the whole flow: a client-supplied `ID`.** It's required on the first call, becomes the primary key of both the `FiatPaymentInvoice` row and the `TokenizedAssetSubscription` row, and is reused as Flutterwave's `tx_ref` — one id, no separate linking columns needed to go from a webhook payload back to the right transaction.

Two-call, three-phase flow:

1. **Call 1 — quote/generate.** Client calls the fiat subscription endpoint with `FiatTokenizedAssetSubscriptionInput` (new merged model, §3), including a client-generated `ID` (**required**) and no `TransactionSignature`. Server resolves the country's InternalBalance asset, builds the transaction XDR using a **channel account** as the source (signed by the channel account, the InternalBalance issuer's multisig signers, and the tokenization issuing profile wallet if a new trustline is opened), and persists it — in one DB transaction — as a `FiatPaymentInvoice` row (`ID` = the supplied id, `Status: "PENDING"`, `PaymentType: "ASSET PURCHASE"`, `Transaction` = the built XDR, plus a new channel-account field) **and** a `TokenizedAssetSubscription` row (same `ID`). Returns the XDR to the client. A retried call 1 (same `ID`, still unsigned) returns the already-stored XDR instead of regenerating and reserving a second channel account.
2. **Call 2 — sign & finalize the invoice.** Client signs the XDR locally (buyer's own key — the operations sourced from the buyer's wallet still need the buyer's signature) and calls the endpoint again with the same `ID` + `TransactionSignature`. Server loads the existing `FiatPaymentInvoice` row by `ID`, stores the signature on it, and commits. **Nothing is submitted to the blockchain here.** The client then takes that same `ID` as `tx_ref` into its Flutterwave checkout — payment happens entirely client-side from here, no further backend involvement until the webhook fires.
3. **Async — webhook settles it.** Flutterwave calls `postCallbacksFlutterwaveWebhookHandler`. Its existing `Product == "ASSET PURCHASE"` branch (currently a TODO stub) loads the `FiatPaymentInvoice` by `event.Data.TxRef` (== the same `ID`), retrieves `Transaction` + `TransactionSignature` from it, submits via `network.SubmitXdrWithSignature` (same call shape as the cloned function's non-multiparty submission path), and on success saves the `FiatPayment` audit record, flips the invoice to `Status: "COMPLETED"`, updates the linked `TokenizedAssetSubscription.TransactionID`, releases the channel account, and sends a push notification.

A background goroutine (started at service boot, alongside the other `main.go` background loops) sweeps `FiatPaymentInvoice` rows stuck `PENDING` for more than 2 days and marks them `EXPIRED`, releasing any channel account still reserved for them.

## 2. Reference: existing on-chain flow being cloned

`SubscribeToTokenizedAsset` (`internal/components/users/services/tokenized_assets.go:2631`):
- Validates KYC, sale status, per-wallet purchase cap.
- Resolves `PaymentAssetCode`/`PaymentAssetIssuer` (defaults to CNGN, rejects internal-balance/tokenized-asset codes).
- Validates swap info + liquidity.
- If `input.TransactionSignature == ""`: calls `generateAssetSubscriptionXdr(...)` (line 2920) to build the XDR, returns it unsubmitted.
- If signature present (non-multiparty): `network.SubmitXdrWithSignature(client, subscriber.PrimarySigner, input.Transaction, input.TransactionSignature)`, saves `TransactionID`.
- If multiparty: creates a `PendingAuth` row instead of submitting.

`generateAssetSubscriptionXdr` (line 2920) stacks 7 operations — **the fiat clone keeps 6 of them, dropping op 1**:
1. `Payment` — buyer wallet → `ta.FundsHoldingWalletPublicKey`, stablecoin (**removed in the fiat clone** — fiat payment substitutes for this on-chain leg).
2. (conditional) `ChangeTrust` — buyer trusts InternalBalance asset.
3. (conditional) `SetTrustLineFlags` — issuer authorizes that trustline.
4. `Payment` — InternalBalance issuer → buyer, InternalBalance asset, 1:1 mint of the amount paid.
5. (conditional) `ChangeTrust` + `SetTrustLineFlags` — buyer trusts the tokenized asset if needed.
6. `PathPaymentStrictSend` — InternalBalance asset → tokenized asset (the actual swap).
7. `ChangeTrust` — InternalBalance asset, `Limit: "0"` (drops the trustline in the same atomic tx).

Signing today: channel account keypair (multiparty only), `getInternalBalanceIssuingSigners()` (always — parses `INTERNAL_BALANCE_ISSUING_SIGNERS`, a CSV of secret seeds, trimmed and parsed to `keypair.Full`; already implemented at `tokenized_assets.go:3183`), and the `TOKENIZATION_ISSUING_PROFILE_WALLET` keypair when a new tokenized-asset trustline is opened. The fiat clone reuses `getInternalBalanceIssuingSigners()` as-is — no new env var needed.

Channel accounts live only in memory: `sharedconfig.GlobalConfig.ChannelAccounts` (buffered channel, acts as a pool) and `InUseChannelAccounts` (map). `ReleaseInUseChannelAccount(pk)` / `StoreInUseChannelAccount(kp)` are at `internal/sharedconfig/models.go:264-284`. At startup, `main.go:359-385` walks the configured channel accounts and, for each, checks `PendingAuth` for `transaction_status = 'PENDING' AND transaction_source = <address>` — if found, marks it in-use (`StoreInUseChannelAccount`) instead of adding it to the available pool. **This scan needs a second check added for the fiat flow** (§7), since this flow never creates a `PendingAuth` row.

## 3. New merged input model: `FiatTokenizedAssetSubscriptionInput`

Location: `internal/components/users/models/tokenization.go`, alongside `TokenizedAssetSubscriptionInput`.

This is a genuine merge — the two source structs already share several concepts (`ID`/amount/wallet/tokenized-asset), so the merged struct doesn't duplicate them:

```go
type FiatTokenizedAssetSubscriptionInput struct {
	// identity — required on call 1, resent unchanged on call 2
	ID string `json:"id"`

	// from TokenizedAssetSubscriptionInput
	TokenizedAssetID     string   `json:"tokenizedAssetId"`
	WalletPublicKey      string   `json:"walletPublicKey"`
	Amount               float64  `json:"amount"` // fiat amount in tokenized asset quote currency
	SwappedEstimate      string   `json:"swappedEstimate"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
	Messages             []string `json:"messages"`
	Memo                 string   `json:"memo"`
	SignatureRequired    int      `json:"signatureRequired"`

	// from FiatPaymentInvoice — accepted for JSON-decoding convenience only; ServiceProvider,
	// Status, PaymentType are always overridden server-side (see §4), matching the existing
	// postUsersFiatFlutterwaveHandler convention of force-setting these fields after decode
	ServiceProvider string `json:"serviceProvider"`

	// server-populated, not client input
	SubscriberUsername string `json:"-"`
	TransactionSource  string `json:"-"`
}
```

`PaymentAssetCode`/`PaymentAssetIssuer` from the original `TokenizedAssetSubscriptionInput` are **not** carried into this merged struct at all — for the fiat flow there's no client-supplied payment asset; it's always resolved server-side (§4.1). `WalletAlias` and `TokenizedAssetID` don't need separate client-supplied variants either — they're resolved server-side from `TokenizedAssetID`/`WalletPublicKey` lookups the same way `postUsersFiatFlutterwaveHandler`-adjacent code already resolves a wallet/asset from its public key/id, and then copied onto both persisted rows.

Call 1 vs. call 2 is distinguished by whether `TransactionSignature` is present — the same discriminator the original `SubscribeToTokenizedAsset` already uses at line 2773 — not by a separate flag.

## 4. Data model changes

### 4.1 `TokenizedAssetSubscription` (`internal/components/users/models/tokenization.go:1956`)

Add fields:

```go
type TokenizedAssetSubscription struct {
	// ...existing fields unchanged, including the existing ID primary key...
	PaymentAssetCode   string `json:"paymentAssetCode"`   // NEW
	PaymentAssetIssuer string `json:"paymentAssetIssuer"` // NEW
}
```

For the fiat path: `PaymentAssetCode = *countryConfig.InternalBalanceTokenCode`, `PaymentAssetIssuer = "FIAT"` — a sentinel, not a real Stellar issuer, signaling "this subscription's InternalBalance leg was fiat-funded, not swap-validated against a real payment asset."

The row's `ID` is set explicitly to the client-supplied `FiatTokenizedAssetSubscriptionInput.ID` at creation time (call 1), not auto-generated — check however `ID` is normally populated for this model (a `BeforeCreate` hook, if one exists) and make sure the fiat creation path sets `ID` explicitly before `Create` so it isn't overwritten.

No `Transaction`/`TransactionSignature`/`TransactionSource` fields are added here — per §4.2, those already live on `FiatPaymentInvoice`. `TokenizedAssetSubscription.TransactionID` (existing field) is still set, but only later, by the webhook, once the transaction actually lands on-chain (§6.5) — mirroring how the original function uses it.

### 4.2 `FiatPaymentInvoice` (`internal/components/users/models/fiatPayment.go:39`)

**Current state, already in place:**

```go
type FiatPaymentInvoice struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `gorm:"default:now()" json:"createdAt"`
	ServiceProvider      string    `json:"serviceProvider"`
	Username             string    `json:"username"`
	Amount               float64   `json:"amount"`
	PaymentType          string    `json:"paymentType"`
	Status               string    `gorm:"default:'PENDING'" json:"status"`
	Refunded             int       `gorm:"default:0" json:"refunded"`
	TokenizedAssetID     *string   `gorm:"null;size:100" json:"tokenizedAssetId"`
	WalletAlias          *string   `gorm:"null;size:100" json:"walletAlias"`
	WalletPublicKey      *string   `gorm:"null;size:100" json:"walletPublicKey"`
	Transaction          *string   `json:"transaction"`
	TransactionSignature *string   `json:"transactionSignature"`
}
```

`Transaction`/`TransactionSignature`/`TokenizedAssetID`/`WalletAlias`/`WalletPublicKey` already exist (all nullable, since most other invoice types — e.g. `"ACTIVATION"` — don't use them). **One field is still missing and needs to be added** — the channel account reservation, which nothing on this model currently tracks:

```go
TransactionSource *string `gorm:"null;size:100" json:"-"` // NEW — channel account public key used as this transaction's source, for pool bookkeeping (§7)
```

`PaymentType` stays a free-form string (no enum today — the only other branched-on value in the codebase is `"ACTIVATION"`, set server-side). `"ASSET PURCHASE"` is this feature's convention, matching the string the webhook TODO already checks (`event.MetaData.Product`).

**Assumption to verify during implementation**: `FiatPaymentInvoice.ID` is client-supplied today (`postUsersFiatFlutterwaveHandler` unmarshals the whole invoice from the request body including `ID`, at `internal/components/users/controllers/handlers_impl.go:874`), and the existing "activation" webhook branch matches by `event.MetaData.UserID`, not by invoice id — there's no precedent in this codebase for looking an invoice up by `event.Data.TxRef`. This plan assumes the client sets its generated `ID` as both `FiatPaymentInvoice.ID` *and* the Flutterwave `tx_ref`, so the webhook can do `WHERE id = event.Data.TxRef`. Confirm `tx_ref` is client-settable in the Flutterwave SDK integration before implementing.

## 5. New service function: `SubscribeToTokenizedAssetByFiat`

Location: `internal/components/users/services/tokenized_assets.go`, alongside `SubscribeToTokenizedAsset`.

```go
func SubscribeToTokenizedAssetByFiat(subscriber *userModels.User, subscriberWallet *userModels.UserWallet, ta *userModels.TokenizedAsset, input *userModels.FiatTokenizedAssetSubscriptionInput, gc *sharedconfig.GlobalConfig) (invoice userModels.FiatPaymentInvoice, err error)
```

Cloned from `SubscribeToTokenizedAsset` with these differences:

1. **No client-supplied payment asset.** Resolve `countryConfig := userModels.CountryCode(*ta.AssetCountryLocation).GetConfig(gc)`, require `InternalBalanceTokenCode`/`InternalTokenIssuer` non-nil (same guard as `generateAssetSubscriptionXdr` today). This is what feeds `TokenizedAssetSubscription.PaymentAssetCode`/`PaymentAssetIssuer` (§4.1). Skip the stablecoin-resolution/rejection block entirely (lines 2709-2734 of the original) — there is no payment asset to validate against a real issuer.
2. **Require `input.ID != ""` unconditionally** (reject with 400 on either call if missing) — it's the row identity for both `FiatPaymentInvoice` and `TokenizedAssetSubscription` on both calls, not just call 1.
3. **Call 1** (`input.TransactionSignature == ""`): first check for an existing invoice, `gc.DB.Where("id = ?", input.ID).First(&existing)`. If found and unsigned/unsubmitted (`TransactionSignature == nil && Status == "PENDING"`), return it as-is — **do not regenerate the XDR or reserve a second channel account** on a retried call 1. If found but already signed, reject (call 1 shouldn't be re-invoked past that point). If not found:
   - Run the existing validations (KYC, sale status, purchase cap, swap/liquidity checks).
   - Call `generateAssetSubscriptionFiatXdr(...)` (new function, §6) to build the XDR using a channel account, returning `(xdrBase64, channelAccountPublicKey, err)`.
   - Open `dbTX := gc.DB.Begin()`. Create the `TokenizedAssetSubscription` row with `ID: input.ID` (set explicitly), `PaymentAssetCode`/`PaymentAssetIssuer` from step 1, and the usual fields (`TokenizedAssetID`, `AssetCode`/`AssetIssuer` from `ta`, `WalletAlias`/`WalletPublicKey` from `subscriberWallet`, `Amount`, `Price`, `SubscriberUsername`). Create the `FiatPaymentInvoice` row with `ID: input.ID`, `ServiceProvider: "flutterwave"`, `Username: subscriber.Username`, `Amount: input.Amount`, `PaymentType: "ASSET PURCHASE"`, `Status: "PENDING"`, `TokenizedAssetID: &ta.ID`, `WalletAlias: &subscriberWallet.WalletAlias`, `WalletPublicKey: &subscriberWallet.PublicKey`, `Transaction: &xdrBase64`, `TransactionSource: &channelAccountPublicKey`. Commit both in the same transaction.
   - Return the new invoice — caller returns `Transaction` to the client.
4. **Call 2** (`input.TransactionSignature != ""`): load the existing `FiatPaymentInvoice` row by `Where("id = ? AND username = ?", input.ID, subscriber.Username)` (ownership check; 404 if missing — call 2 must follow a call 1), require `TransactionSignature == nil && Status == "PENDING"` (reject as already-processed otherwise — idempotency guard against double submission). Set `TransactionSignature = &input.TransactionSignature` **as-is, with no signature verification** (confirmed — do not check it's well-formed or matches the wallet's signer set; any validity problem surfaces later at submission time in the webhook), save. **No blockchain submission happens in this branch** — this is the "commits the DB operation" step, not a chain write.
5. **Never calls `network.SubmitXdrWithSignature` itself** — unlike the original function, submission is entirely the webhook's job (§6). This function only ever generates or persists; it never reaches the original's "submission path" or "multiparty path" branches.
6. **No `PendingAuth` row.** This isn't the multiparty co-signer approval flow — there's exactly one buyer signature expected, collected synchronously in call 2. Channel-account "in use" tracking rides on `FiatPaymentInvoice.TransactionSource` instead (§7), not on `PendingAuth`.

## 6. New XDR builder: `generateAssetSubscriptionFiatXdr`

Location: same file, cloned from `generateAssetSubscriptionXdr` (line 2920).

Differences from the original:

- **Omit operation 1** (buyer → `FundsHoldingWalletPublicKey` stablecoin payment) entirely — fiat settlement (confirmed later by the webhook) is what authorizes the mint; there's no on-chain payment leg.
- **Always use a channel account as the transaction source**, not conditionally on `swapInfo.Multiparty` — pop one via `chanAccount := <-gc.ChannelAccounts` (same pattern as the existing conditional branch), but don't `defer` it back to the pool: call `gc.StoreInUseChannelAccount(chanAccount)` and keep its public key to return to the caller. It gets released explicitly later — on successful blockchain submission (§8) or on invoice expiry (§9) — not when this function returns, since the transaction isn't submitted until the async webhook fires, potentially minutes or hours later.
- Keep the remaining 6 operations (InternalBalance trustline + authorize, 1:1 mint payment, tokenized-asset trustline if needed, `PathPaymentStrictSend` swap, InternalBalance trustline removal) unchanged, still sourced from `wallet.ID` per-operation and the InternalBalance issuer per-operation as today.
- **Signing**: channel account keypair (always, since it's always the tx source now), `getInternalBalanceIssuingSigners()` (unchanged, reused as-is), and the `TOKENIZATION_ISSUING_PROFILE_WALLET` keypair if a new tokenized-asset trustline was opened (unchanged). **Do not** sign with the buyer's key — that signature arrives later, out of band, as `input.TransactionSignature` in call 2, and gets appended at submission time via `network.SubmitXdrWithSignature`, exactly like the original's non-multiparty submission path.
- Signature: `func generateAssetSubscriptionFiatXdr(wallet *userModels.UserWallet, ta *userModels.TokenizedAsset, swapInfo *swapModels.SwapSendInfo, gc *sharedconfig.GlobalConfig) (xdrBase64 string, channelAccountPublicKey string, err error)`.

## 7. Channel account pool bookkeeping

`main.go:359-385` reconstructs `InUseChannelAccounts` at startup by checking `PendingAuth` for each configured channel account. Since the fiat subscription flow reserves a channel account without ever creating a `PendingAuth` row, add a second check right next to the existing one (inside the same `for _, v := range scas` loop, alongside the `PendingAuth` lookup at line 372-383). This check must look at **both `PENDING` and `COMPLETED`** invoices referencing the account, not just `PENDING`:

```go
// existing: check PendingAuth for transaction_status = 'PENDING' AND transaction_source = k.Address()
// NEW: also check FiatPaymentInvoice for any row that reserved this channel account
var fiatInvoice userModels.FiatPaymentInvoice
errFetchFiat := database.Where("transaction_source = ?", k.Address()).First(&fiatInvoice).Error
if errFetchFiat == nil {
	if fiatInvoice.Status == "PENDING" {
		// still mid-flight in the fiat purchase flow — genuinely in use
		log.Printf("[ADDING KEY TO IN-USE CHANNEL ACCOUNT LIST] %v\n", k.Address())
		globalConfig.StoreInUseChannelAccount(k)
		continue
	}
	if fiatInvoice.Status == "COMPLETED" {
		// the webhook's success path (§8.6) should have already released this account —
		// its presence here means that release step didn't run for some reason (e.g. a
		// crash between the status flip and the release call). Self-heal: don't mark it
		// in-use, let it fall through to the available pool below, and clear the stale
		// reference so this fallback doesn't have to re-detect it on every future restart.
		log.Printf("[CHANNEL ACCOUNT RELEASE NOT PERSISTED — SELF-HEALING] %v (invoice %v was COMPLETED but still held a channel account)\n", k.Address(), fiatInvoice.ID)
		database.Model(&userModels.FiatPaymentInvoice{}).Where("id = ?", fiatInvoice.ID).Update("transaction_source", nil)
	}
	// any other status (e.g. EXPIRED) already had its channel account released by
	// ExpireStalePaymentInvoices (§9) — nothing to do, falls through to the pool.
}
```

This ensures a service restart doesn't hand out a channel account still reserved for an in-flight fiat purchase awaiting its webhook, while also self-healing the case where a completed purchase's release step (§8.6) never ran — rather than leaking that channel account out of the pool forever.

## 8. Webhook completion: `postCallbacksFlutterwaveWebhookHandler`

File: `internal/components/callbacks/controllers/handlers_impl.go:771-773`. Replace:

```go
//Condition to process asset purchase
if strings.EqualFold(event.MetaData.Product, "ASSET PURCHASE") && strings.EqualFold(event.Data.Status, "successful") {
	//TODO: perform asset purchase logic here
}
```

with logic mirroring the "activation" branch's shape above it:

1. Load the invoice: `gc.DB.Where("id = ? AND status = ? AND payment_type = ?", event.Data.TxRef, "PENDING", "ASSET PURCHASE").First(&invoice)`. If not found or already processed, log and skip (webhook may retry/duplicate — this guard makes it idempotent).
2. Guard `invoice.Transaction != nil && invoice.TransactionSignature != nil` (signature was actually collected in call 2 — if not, this invoice never got signed and shouldn't be submitted).
3. Load the linked subscription by the same id: `gc.DB.Where("id = ?", invoice.ID).First(&subscription)` — `TokenizedAssetSubscription.ID` and `FiatPaymentInvoice.ID` are the same client-supplied `ID` by construction (§4.1/§4.2), so no extra linking column is needed. Guard `subscription.TransactionID == ""` (not already submitted).
4. Load the subscriber (`userModels.Username(subscription.SubscriberUsername).GetSimpleUser(gc.DB, gc)`) for `PrimarySigner` and for sending the push notification. `invoice.WalletAlias`/`invoice.TokenizedAssetID` are already on hand for the notification copy without a further join.
5. Submit: `txResult, err := network.SubmitXdrWithSignature(gc.GetBlockchainClient(), subscriber.PrimarySigner, *invoice.Transaction, *invoice.TransactionSignature)` — same call shape as `SubscribeToTokenizedAsset`'s existing non-multiparty submission path (line ~2815).
6. On success:
   - `subscription.TransactionID = txResult`, save.
   - `userServices.SaveUserPaymentData(user.Username, invoice.ServiceProvider, invoice.PaymentType, txResult, invoice.Amount, gc)` — the audit-trail `FiatPayment` write, matching the existing "activation" branch's pattern.
   - Flip the invoice to `"COMPLETED"` via a targeted update — `gc.DB.Model(&userModels.FiatPaymentInvoice{}).Where("id = ?", invoice.ID).Update("status", "COMPLETED")` — **not** a blind `Save` of a fresh struct; the existing "activation" branch's use of `SaveUserPaymentInvoiceData` doesn't transition an existing row by id today, which would be a bug here since the specific `PENDING` row must flip to `COMPLETED`, not get upserted as a new one.
   - Release the channel account: `gc.ReleaseInUseChannelAccount(*invoice.TransactionSource)`.
   - Invalidate caches, then send a push notification via `user.SendPushMessage(title, msg, "", dataPayload, gc)`, following the existing convention (`internal/components/users/models/user_methods.go:2791`).
7. On failure: log the error, leave the invoice `PENDING` (**confirmed policy** — no explicit `"FAILED"` status) and do **not** release the channel account — the same invoice may still be resubmitted (e.g. a duplicated webhook delivery, or a manual reconciliation job). Recovery relies entirely on the 2-day expiry sweep (§9), which will eventually flip a permanently-failing invoice to `"EXPIRED"` and release its channel account.

## 9. Background goroutine: expire stale invoices

New service function, `internal/components/users/services/fiatPayment.go`:

```go
func ExpireStalePaymentInvoices(gc *sharedconfig.GlobalConfig) error {
	cutoff := time.Now().Add(-48 * time.Hour)
	var stale []userModels.FiatPaymentInvoice
	if err := gc.DB.Where("status = ? AND created_at < ?", "PENDING", cutoff).Find(&stale).Error; err != nil {
		return err
	}
	for _, inv := range stale {
		if inv.TransactionSource != nil {
			gc.ReleaseInUseChannelAccount(*inv.TransactionSource)
		}
		gc.DB.Model(&userModels.FiatPaymentInvoice{}).Where("id = ?", inv.ID).Update("status", "EXPIRED")
	}
	return nil
}
```

Registered in `main.go`, in the same block as the other background loops (~line 820, right before `//setup router`), following the codebase's uniform `for { <work>; time.Sleep(duration) }` idiom (no `time.Ticker` used anywhere in this file):

```go
go func() {
	for {
		if err := userServices.ExpireStalePaymentInvoices(&globalConfig); err != nil {
			log.Printf("[MAIN] error expiring stale payment invoices: %v\n", err)
		}
		time.Sleep(30 * time.Minute)
	}
}()
```

This sweep is generic to all `FiatPaymentInvoice` rows, not just asset-purchase ones — any stuck `PENDING` invoice of any payment type should expire the same way; the channel-account release step naturally only fires for rows that have one (`TransactionSource != nil`), i.e. asset-purchase invoices.

## 10. HTTP layer

New endpoint (mirrors `postUsersFiatFlutterwaveHandler`'s auth/registration shape, but drives the two-call subscription flow instead of a single invoice write):

```
POST /v1/users/fiat/tokenized-assets/subscribe
```

Handler `postUsersTokenizedAssetSubscriptionFiatHandler(gc)`, registered in `internal/components/users/controllers/main.go` next to the existing `postUsersFiatFlutterwaveHandler` route:

1. Auth via `middleware.AuthenticationMiddlewareUsingTimestamp()` + `usersDB.GetUserFromPrimarySigner`, same as `postUsersFiatFlutterwaveHandler`.
2. Decode `userModels.FiatTokenizedAssetSubscriptionInput` from the body; reject with 400 if `ID` is empty.
3. Load `ta` (`TokenizedAssetID`) and `subscriberWallet` (`WalletPublicKey` + user), same lookups the existing `SubscribeToTokenizedAsset` HTTP call sites perform (`internal/components/users/controllers/handlers_impl.go:6224`, `:6385` — mirror their pre-fetch pattern).
4. Set `input.SubscriberUsername = user.Username`, `input.ServiceProvider = "flutterwave"` (server overrides, same convention as `postUsersFiatFlutterwaveHandler` force-setting `tInput.ServiceProvider`/`tInput.Status` after decode).
5. Call `userServices.SubscribeToTokenizedAssetByFiat(user, subscriberWallet, ta, &input, gc)`.
6. Respond:
   - Call 1: `{"id": ..., "transaction": ...}` (unsigned-by-buyer XDR for the client to sign).
   - Call 2: `{"id": ..., "status": "awaiting-payment"}` — no XDR needed back, submission happens later via webhook.

## 11. File-by-file checklist

| File | Change |
|---|---|
| `internal/components/users/models/tokenization.go` | Add `FiatTokenizedAssetSubscriptionInput`; add `PaymentAssetCode`/`PaymentAssetIssuer` to `TokenizedAssetSubscription`. |
| `internal/components/users/models/fiatPayment.go` | Add `TransactionSource *string` to `FiatPaymentInvoice` (the other new fields — `Transaction`, `TransactionSignature`, `TokenizedAssetID`, `WalletAlias`, `WalletPublicKey` — already exist). |
| `internal/components/users/services/tokenized_assets.go` | Add `SubscribeToTokenizedAssetByFiat` and `generateAssetSubscriptionFiatXdr`. |
| `internal/components/users/services/fiatPayment.go` | Add `ExpireStalePaymentInvoices`. |
| `internal/components/users/controllers/handlers_impl.go` | Add `postUsersTokenizedAssetSubscriptionFiatHandler`. |
| `internal/components/users/controllers/main.go` | Register the new route. |
| `internal/components/callbacks/controllers/handlers_impl.go` | Complete the `ASSET PURCHASE` TODO block (line ~771-773). |
| `main.go` | Extend the channel-account startup scan (~line 372-385, inside the existing `scas` loop) with the `FiatPaymentInvoice` check; register the new expiry goroutine (~line 820). |

## 12. Resolved decisions

1. **Invoice/tx_ref correlation** — **confirmed**: Flutterwave's `tx_ref` is client-settable, and the client sets it to the same `ID` used for the invoice/subscription. The webhook matches on `event.Data.TxRef` as designed (§8.1).
2. **Failure/retry policy** for a webhook-side submission failure (§8.7) — **confirmed**: leave the invoice `PENDING` indefinitely; no explicit `"FAILED"` status. Recovery relies on the 2-day expiry sweep (§9).
3. **Signature verification** — **confirmed**: do not sanity-check or verify the buyer's signature at call 2. It's persisted as-is; any validity problem surfaces later at submission time in the webhook via `network.SubmitXdrWithSignature`.
