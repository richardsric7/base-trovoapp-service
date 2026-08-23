# Tokenized Asset Purchase by Fiat — Implementation Plan

Status: **PLANNING ONLY — not implemented.**
Author: Claude (research grounded in current `main` branch of `backend`, branch `ric-asset-tokenization-purchase`)
Related prior work: `backend/issue-internal-balance-on-demand.md`, commits `35d34192`, `771059a5`, `08691be0`, `3c53afd0`, `e784093c`.

## 1. Goal

Let a user buy a tokenized asset with fiat instead of an on-chain stablecoin, without the backend ever needing a second callback from the client after payment. Today's on-chain flow (`SubscribeToTokenizedAsset`) requires the buyer to pay a stablecoin leg on-chain themselves. The fiat flow replaces that leg with an off-chain fiat payment (Flutterwave) and mints the equivalent InternalBalance token directly to the buyer (this "on-demand mint" mechanism already exists and is reused as-is).

Two-call, three-phase flow:

**One identifier drives all of it: `PaymentInvoiceID`.** It's client-generated, required on call 1, resent unchanged on call 2, and is used directly as the primary key of both the `TokenizedAssetSubscription` row and the `FiatPaymentInvoice` row (rather than each row getting its own separately-generated id) — and again as Flutterwave's `tx_ref`. One id, three roles, no separate `TokenizedAssetSubscriptionID` needed. Call 1 vs. call 2 is distinguished the same way the original `SubscribeToTokenizedAsset` already distinguishes its two calls: by whether `TransactionSignature` is present.

1. **Call 1 — quote/generate.** Client calls the fiat subscription endpoint with `TokenizedAssetSubscriptionInput`, including a client-generated `PaymentInvoiceID` (**required**, fixed for the lifetime of this purchase attempt — generated once, reused verbatim on call 2 and as Flutterwave's `tx_ref`) and no `TransactionSignature`. Server builds the transaction XDR (signed by every server-held signer: channel account, InternalBalance issuer multisig, and tokenization issuing profile wallet if needed), persists it against a new `TokenizedAssetSubscription` row **whose `ID` is set explicitly to `PaymentInvoiceID`** (not auto-generated), and returns the XDR to the client. A retried call 1 (same `PaymentInvoiceID`, still unsigned) returns the already-stored XDR instead of regenerating and reserving a second channel account.
2. **Call 2 — sign & register invoice.** Client signs the XDR locally (buyer's own key, since ops sourced from the buyer's wallet still need the buyer's signature) and calls the same endpoint again with the same `PaymentInvoiceID` + `TransactionSignature`. Server loads the subscription row by that id, stores the signature on it, and, in the same DB transaction, creates a `FiatPaymentInvoice` with **`ID` also set to `PaymentInvoiceID`** (`PaymentType: "ASSET PURCHASE"`, `Status: "PENDING"`). **Nothing is submitted to the blockchain here.** The client then takes that same id as `tx_ref` into its Flutterwave checkout — payment happens entirely client-side, no further backend involvement until the webhook fires.
3. **Async — webhook settles it.** Flutterwave calls `postCallbacksFlutterwaveWebhookHandler`. Its existing `Product == "ASSET PURCHASE"` branch (currently a TODO stub) loads the invoice by `tx_ref`, then loads the subscription **by that same id** (no separate link column needed — the two rows share a primary key by construction), submits the already-built, already-server-signed XDR with the stored buyer signature via `network.SubmitXdrWithSignature`, marks the invoice `COMPLETED`, releases the channel account, and pushes a notification to the user.

A background goroutine sweeps `FiatPaymentInvoice` rows stuck `PENDING` for more than 2 days and marks them `EXPIRED`, releasing any channel account still reserved for their subscription so the pool doesn't leak.

## 2. Reference: existing on-chain flow being cloned

`internal/components/users/services/tokenized_assets.go:2631` — `SubscribeToTokenizedAsset(subscriber, subscriberWallet, ta, input, gc)`:
- Validates KYC, sale status, per-wallet purchase cap.
- Resolves `PaymentAssetCode`/`PaymentAssetIssuer` (defaults to CNGN, rejects internal-balance/tokenized-asset codes).
- Validates swap info + liquidity.
- If `input.TransactionSignature == ""`: calls `generateAssetSubscriptionXdr(...)` (line 2920) to build XDR, returns it unsaved-to-chain.
- If signature present (non-multiparty): `network.SubmitXdrWithSignature(client, subscriber.PrimarySigner, input.Transaction, input.TransactionSignature)`, then saves `TransactionID`.
- If multiparty: creates a `PendingAuth` row instead of submitting.

`generateAssetSubscriptionXdr` (line 2920) stacks 7 operations:
1. `Payment` — buyer wallet → `ta.FundsHoldingWalletPublicKey`, stablecoin (**this is the on-chain payment leg the fiat flow removes**).
2. (conditional) `ChangeTrust` — buyer trusts InternalBalance asset.
3. (conditional) `SetTrustLineFlags` — issuer authorizes that trustline.
4. `Payment` — InternalBalance issuer → buyer, InternalBalance asset, 1:1 mint of the swap amount.
5. (conditional) `ChangeTrust` + `SetTrustLineFlags` — buyer trusts the tokenized asset if needed.
6. `PathPaymentStrictSend` — InternalBalance asset → tokenized asset (the actual swap).
7. `ChangeTrust` — InternalBalance asset, `Limit: "0"` (drops the trustline in the same atomic tx).

Signing: channel account keypair (multiparty only, today), `getInternalBalanceIssuingSigners()` (always — parses `INTERNAL_BALANCE_ISSUING_SIGNERS` CSV env var of secret seeds, already implemented at `tokenized_assets.go:3179`), and the `TOKENIZATION_ISSUING_PROFILE_WALLET` keypair when a new tokenized-asset trustline is opened.

**This is exactly the operation list the user specified for the fiat flow, minus operation 1** — fiat payment substitutes for the on-chain stablecoin leg.

Channel accounts today live only in memory: `sharedconfig.GlobalConfig.ChannelAccounts` (buffered channel, acts as a pool/lock) and `InUseChannelAccounts` (map). `ReleaseInUseChannelAccount(pk)` / `StoreInUseChannelAccount(kp)` at `internal/sharedconfig/models.go:264-284`. On startup, `main.go:370-383` repopulates `InUseChannelAccounts` by scanning `PendingAuth` rows with `transaction_status = 'PENDING'`. **There is no dedicated `channel_accounts` DB table** — "in use" state is reconstructed from whichever business table recorded the channel account's public key.

## 3. Data model changes

### 3.1 `TokenizedAssetSubscription` (`internal/components/users/models/tokenization.go:1956`)

Add fields:

```go
type TokenizedAssetSubscription struct {
	// ...existing fields unchanged, including the existing `ID` primary key...
	PaymentAssetCode     string `json:"paymentAssetCode"`     // NEW — persisted, currently only lives on the Input struct
	PaymentAssetIssuer   string `json:"paymentAssetIssuer"`   // NEW — persisted, currently only lives on the Input struct
	Transaction           string `gorm:"type:text" json:"-"`           // NEW — the built, server-signed XDR, so call 2 doesn't need to rebuild it
	TransactionSignature  string `gorm:"type:text" json:"-"`           // NEW — buyer's signature, captured on call 2, consumed by the webhook
	TransactionSource     string `json:"-"`                             // NEW — channel account public key used as the tx source, for pool bookkeeping
}
```

No new `PaymentInvoiceID` field is added — for the fiat flow, the existing `ID` field on this row is set explicitly to the client-supplied `PaymentInvoiceID` at creation time (call 1), instead of being auto-generated. This is a deliberate deviation from however `ID` is normally populated for this model (check the existing `BeforeCreate` hook / id-generation path if one exists, and make sure the fiat creation path explicitly sets `ID` before `Create` so it isn't overwritten). Doing it this way means the same id threads through `TokenizedAssetSubscription.ID` → `FiatPaymentInvoice.ID` → Flutterwave `tx_ref` with no extra linking column anywhere.

For the fiat path specifically: `PaymentAssetCode = *countryConfig.InternalBalanceTokenCode`, `PaymentAssetIssuer = "FIAT"` (a sentinel, not a real Stellar issuer — signals "this subscription's InternalBalance leg was fiat-funded, not swap-validated against a real payment asset").

`Transaction`/`TransactionSignature` should not round-trip in JSON responses (raw XDR/signature aren't useful to callers beyond the two calls that need them) — mark `json:"-"` and return them explicitly in the handler's response payload only when needed (call 1 needs to return `Transaction`; call 2 doesn't need to return anything but the subscription id and status).

### 3.2 `TokenizedAssetSubscriptionInput` (`internal/components/users/models/tokenization.go:1984`)

Add one field used only by the fiat path:

```go
PaymentInvoiceID string `json:"paymentInvoiceId"` // NEW — REQUIRED on call 1, resent unchanged on call 2: doubles as the id of both the TokenizedAssetSubscription row and the FiatPaymentInvoice row
```

No separate `TokenizedAssetSubscriptionID` field — `PaymentInvoiceID` alone identifies the row on both calls, since it *is* `TokenizedAssetSubscription.ID` (§3.1). It's **required on call 1** and must be resent identically on call 2 (the server looks the row up by it); consistent with the existing convention that `FiatPaymentInvoice.ID` is always client-supplied (`postUsersFiatFlutterwaveHandler` unmarshals it straight from the request body today, at `internal/components/users/controllers/handlers_impl.go:845`). The handler must reject either call with a 400 if it's empty — never fall back to server-generating one, and never accept a call 2 whose `PaymentInvoiceID` doesn't match an existing row.

`PaymentAssetCode`/`PaymentAssetIssuer` already exist on this struct — for the fiat handler these are ignored if the client sends them; the server always overrides them with the InternalBalance code / `"FIAT"`.

### 3.3 `FiatPaymentInvoice` (`internal/components/users/models/fiatPayment.go:39`)

**Fields already added** (by the user, ahead of this plan):

```go
TokenizedAssetID *string `gorm:"null;size:100" json:"tokenizedAssetId"`
WalletAlias      string  `gorm:"not null;size:100" json:"walletAlias"`
WalletPublicKey  string  `gorm:"not null;size:100" json:"walletPublicKey"`
```

`TokenizedAssetID` is nullable because most existing invoice types (e.g. `"ACTIVATION"`) have no associated tokenized asset — it's only populated for `"ASSET PURCHASE"` invoices. `WalletAlias`/`WalletPublicKey` are `not null`, so every invoice (asset-purchase or not) is expected to carry the wallet it belongs to going forward.

`FiatPaymentInvoice.ID` is still set to the same `PaymentInvoiceID` value on call 2 (§3.2), which is also `TokenizedAssetSubscription.ID` — so the webhook can still go straight from the invoice's `ID` to `TokenizedAssetSubscription.ID` with a plain `Where("id = ?")` lookup, no separate linking column required to find the transaction/signature. The new fields don't replace that lookup; they let the invoice **self-describe** which wallet and asset it's for, without joining to the subscription row — useful for admin/reporting queries over invoices (e.g. "all pending asset-purchase invoices for this wallet") and for building the webhook's push-notification content and sanity checks without an extra load.

At invoice-creation time (call 2, §4), populate them from the already-loaded `taSubscription` row (which already carries `WalletAlias`/`WalletPublicKey`/`TokenizedAssetID` from call 1 via `UpdateTokenizedAssetSubscriptionFromInput` — no need to re-derive from `ta`/`subscriberWallet` a second time):
```go
invoice := userModels.FiatPaymentInvoice{
	ID:               input.PaymentInvoiceID,
	ServiceProvider:  "flutterwave",
	Username:         subscriber.Username,
	Amount:           taSubscription.Amount,
	PaymentType:      "ASSET PURCHASE",
	Status:           "PENDING",
	TokenizedAssetID: &taSubscription.TokenizedAssetID,
	WalletAlias:      taSubscription.WalletAlias,
	WalletPublicKey:  taSubscription.WalletPublicKey,
}
```

`PaymentType` stays a free-form string (there's no enum today — the only other branched-on value in the codebase is `"ACTIVATION"`, set server-side). Document `"ASSET PURCHASE"` as the convention for this feature, matching the string the webhook TODO already checks (`event.MetaData.Product`).

**Assumption to verify during implementation**: `FiatPaymentInvoice.ID` is client-supplied today (`postUsersFiatFlutterwaveHandler` unmarshals the whole invoice from the request body including `ID`), and the existing "activation" webhook branch matches by `event.MetaData.UserID`, not by invoice id — there's no precedent in this codebase for looking an invoice up by `event.Data.TxRef`. This plan assumes the client sets its generated `PaymentInvoiceID` as both `FiatPaymentInvoice.ID` *and* the Flutterwave `tx_ref`, so the webhook can do `WHERE id = event.Data.TxRef`. Confirm this is an acceptable convention (or find the actual field Flutterwave webhooks use to correlate to a `tx_ref` the client set) before implementing.

## 4. New service function: `SubscribeToTokenizedAssetByFiat`

Location: `internal/components/users/services/tokenized_assets.go`, alongside `SubscribeToTokenizedAsset`.

```go
func SubscribeToTokenizedAssetByFiat(subscriber *userModels.User, subscriberWallet *userModels.UserWallet, ta *userModels.TokenizedAsset, input *userModels.TokenizedAssetSubscriptionInput, gc *sharedconfig.GlobalConfig) (taSubscription userModels.TokenizedAssetSubscription, err error)
```

Cloned from `SubscribeToTokenizedAsset` with these differences:

1. **No client-supplied payment asset.** Resolve `countryConfig := userModels.CountryCode(*ta.AssetCountryLocation).GetConfig(gc)`, require `InternalBalanceTokenCode`/`InternalTokenIssuer` non-nil (same guard as `generateAssetSubscriptionXdr` today). Set `input.PaymentAssetCode = *countryConfig.InternalBalanceTokenCode`, `input.PaymentAssetIssuer = "FIAT"`. Skip the stablecoin-resolution/rejection block entirely (lines 2709-2734 in the original) — there is no payment asset to validate against a real issuer.
2. **Require `input.PaymentInvoiceID != ""` unconditionally** (reject with 400 on either call if missing) — it's the row's identity on both calls, not just call 1.
3. **Idempotent two-call handling, discriminated by `TransactionSignature` presence** (same discriminator the original `SubscribeToTokenizedAsset` already uses at line 2773) **rather than by a separate subscription-id field:**
   - If `input.TransactionSignature == ""` (call 1): first check for an existing row, `gc.DB.Where("id = ?", input.PaymentInvoiceID).First(&existing)`. If found and unsigned/unsubmitted (`TransactionSignature == "" && TransactionID == ""`), return it as-is — **do not regenerate the XDR or reserve a second channel account** on a retried call 1. If found but already signed or submitted, reject (call 1 shouldn't be re-invoked past that point). If not found: run the existing validations (KYC, sale status, purchase cap, swap/liquidity checks), call `generateAssetSubscriptionFiatXdr(...)` (new function, §5), and persist a **new** `TokenizedAssetSubscription` row with `ID: input.PaymentInvoiceID` (set explicitly, not auto-generated — see §3.1), `Transaction`, `TransactionSource` populated, `TransactionSignature` empty. Return it — caller returns `Transaction` to the client.
   - If `input.TransactionSignature != ""` (call 2): load the existing subscription row by `Where("id = ? AND subscriber_username = ?", input.PaymentInvoiceID, subscriber.Username)` (ownership check; 404 if missing — call 2 must follow a call 1), require the row not already signed/submitted (`TransactionSignature == "" && TransactionID == ""` — reject as already-processed otherwise, idempotency guard against double submission). Open `dbTX := gc.DB.Begin()`, set `TransactionSignature`, save the row, then create the `FiatPaymentInvoice` inside the same transaction — `ID: input.PaymentInvoiceID` (same value, now also the invoice's id), `PaymentType: "ASSET PURCHASE"`, `Status: "PENDING"`, plus `TokenizedAssetID`/`WalletAlias`/`WalletPublicKey` copied from `taSubscription` (full field list in §3.3) — then commit. **No blockchain submission happens in this branch.**
4. **Never calls `network.SubmitXdrWithSignature` itself** — unlike the original function, submission is entirely the webhook's job (§6). This function only ever generates or persists; it never reaches the original's "submission path" or "multiparty path" branches.
5. **No `PendingAuth` row.** This isn't the multiparty co-signer approval flow — there's exactly one buyer signature expected, collected synchronously in call 2. Channel-account "in use" tracking rides on the subscription row itself instead (see §7), not on `PendingAuth`.

## 5. New XDR builder: `generateAssetSubscriptionFiatXdr`

Location: same file, cloned from `generateAssetSubscriptionXdr` (line 2920).

Differences from the original:

- **Omit operation 1** (buyer → `FundsHoldingWalletPublicKey` stablecoin payment) entirely — there is no on-chain payment leg; fiat settlement (confirmed later by the webhook) is what authorizes the mint.
- **Always use a channel account as the transaction source**, not conditionally on `swapInfo.Multiparty` — pop one via `chanAccount := <-gc.ChannelAccounts` (same pattern as the existing conditional branch), but don't `defer` it back to the pool: instead call `gc.StoreInUseChannelAccount(chanAccount)` and keep its public key to return to the caller as `TransactionSource`. It gets released explicitly later — on successful blockchain submission (§6) or on invoice expiry (§8) — not when this function returns, since the transaction isn't submitted until the async webhook fires, potentially minutes or hours later.
- Keep the remaining 6 operations (InternalBalance trustline + authorize, 1:1 mint payment, tokenized-asset trustline if needed, `PathPaymentStrictSend` swap, InternalBalance trustline removal) unchanged from the original, still sourced from `wallet.ID` per-operation and the InternalBalance issuer per-operation as today.
- **Signing**: sign with the channel account keypair (always, since it's always the tx source now), `getInternalBalanceIssuingSigners()` (unchanged, reuse the existing helper as-is), and the `TOKENIZATION_ISSUING_PROFILE_WALLET` keypair if a new tokenized-asset trustline was opened (unchanged). **Do not** sign with the buyer's key — that signature arrives later, out of band, as `input.TransactionSignature` in call 2, and gets appended at submission time via `network.SubmitXdrWithSignature`, exactly like the original's non-multiparty submission path.
- Return `(xdrBase64 string, channelAccountPublicKey string, err error)` — the extra return value is new, needed so the caller can persist `TransactionSource`.

## 6. Webhook completion: `postCallbacksFlutterwaveWebhookHandler`

File: `internal/components/callbacks/controllers/handlers_impl.go:771-773`. Replace:

```go
//Condition to process asset purchase
if strings.EqualFold(event.MetaData.Product, "ASSET PURCHASE") && strings.EqualFold(event.Data.Status, "successful") {
	//TODO: perform asset purchase logic here
}
```

with logic mirroring the "activation" branch's shape above it:

1. Load the invoice: `gc.DB.Where("id = ? AND status = ? AND payment_type = ?", event.Data.TxRef, "PENDING", "ASSET PURCHASE").First(&invoice)`. If not found or already processed, log and skip (webhook may retry/duplicate — this guard makes it idempotent).
2. Load the subscription by the same id: `gc.DB.Where("id = ?", invoice.ID).First(&subscription)` — `TokenizedAssetSubscription.ID` and `FiatPaymentInvoice.ID` are the same `PaymentInvoiceID` value by construction (§3.1/§3.2), so no separate link column is needed. Guard `subscription.TransactionID == ""` (not already submitted) and `subscription.TransactionSignature != ""` (signature was actually collected in call 2). As a cheap integrity check (belt-and-braces, not strictly required for the happy path), assert `invoice.TokenizedAssetID != nil` (it's always set for `"ASSET PURCHASE"` invoices — nil here means the row is corrupt or mistyped) and then `subscription.WalletPublicKey == invoice.WalletPublicKey` and `subscription.TokenizedAssetID == *invoice.TokenizedAssetID` before proceeding — both rows were written from the same `taSubscription` in call 2 so they should always agree; a mismatch means the wrong pair got loaded and submission should be aborted rather than risk minting into the wrong wallet.
3. Load the subscriber (`userModels.Username(subscription.SubscriberUsername).GetSimpleUser(gc.DB, gc)`) for `PrimarySigner` and for sending the push notification. Note `invoice.WalletAlias`/`invoice.WalletPublicKey`/`invoice.TokenizedAssetID` are already on hand at this point for the notification copy (§6.5) without a further join.
4. Submit: `txResult, err := network.SubmitXdrWithSignature(gc.GetBlockchainClient(), subscriber.PrimarySigner, subscription.Transaction, subscription.TransactionSignature)` — same call shape as `SubscribeToTokenizedAsset`'s existing non-multiparty submission path (line ~2815).
5. On success: `subscription.TransactionID = txResult`, save; `invoice.Status = "COMPLETED"`, save (`SaveUserPaymentInvoiceData`-style update, but must be a targeted `Where("id = ?").Updates(...)` rather than a blind `Save` of a fresh struct — the existing "activation" branch's use of `SaveUserPaymentInvoiceData` doesn't transition an existing row by id today, which would be a bug here since we need to flip the specific `PENDING` row to `COMPLETED`, not upsert a new one); release the channel account: `gc.ReleaseInUseChannelAccount(subscription.TransactionSource)`; invalidate caches; send push notification via `user.SendPushMessage(title, msg, "", dataPayload, gc)` — using `invoice.WalletAlias` and the tokenized asset's name/code (loadable via `invoice.TokenizedAssetID`) in the message copy, following the existing convention (`internal/components/users/models/user_methods.go:2791`).
6. On failure: log the error, leave the invoice `PENDING` (so it's picked up by a retry — a manual reconciliation job or a future webhook retry) and do **not** release the channel account yet — the same subscription may still be resubmitted (e.g. a duplicated webhook delivery, or a manual retry endpoint added later). If a maximum-retry/definitive-failure policy is wanted, mark the invoice `"FAILED"` and release the channel account then — leaving that decision for implementation-time judgment call rather than baking retry-count tracking into this plan.

## 7. Channel account pool bookkeeping

There's no dedicated channel-accounts DB table today — `main.go:370-383` reconstructs `InUseChannelAccounts` at startup purely from `PendingAuth` rows. Since the fiat subscription flow reserves a channel account without ever creating a `PendingAuth` row, that startup scan needs a second query added right next to the existing one:

```go
// existing: scan PendingAuth for transaction_status = 'PENDING'
// NEW: also scan TokenizedAssetSubscription for rows that reserved a channel account
// but haven't submitted yet (i.e. still mid-flight in the fiat purchase flow)
var pendingFiatSubs []userModels.TokenizedAssetSubscription
database.Where("transaction_source <> '' AND transaction_id = ''").Find(&pendingFiatSubs)
for _, s := range pendingFiatSubs {
	// mark s.TransactionSource as in-use in InUseChannelAccounts / remove from the pool channel,
	// same bookkeeping the existing PendingAuth loop performs
}
```

This ensures a service restart doesn't hand out a channel account that's still reserved for an in-flight fiat purchase awaiting its webhook.

## 8. Background goroutine: expire stale invoices

New service function, `internal/components/users/services/fiatPayment.go`:

```go
func ExpireStalePaymentInvoices(gc *sharedconfig.GlobalConfig) error {
	cutoff := time.Now().Add(-48 * time.Hour)
	var stale []userModels.FiatPaymentInvoice
	if err := gc.DB.Where("status = ? AND created_at < ?", "PENDING", cutoff).Find(&stale).Error; err != nil {
		return err
	}
	for _, inv := range stale {
		// release any channel account still reserved for this invoice's subscription. Since
		// TokenizedAssetSubscription.ID and FiatPaymentInvoice.ID are the same PaymentInvoiceID
		// value by construction (§3.1/§3.2), this is a direct id lookup — no link column needed.
		// Only asset-purchase invoices will match a subscription row; other payment types simply won't.
		var sub userModels.TokenizedAssetSubscription
		if err := gc.DB.Where("id = ? AND transaction_id = ''", inv.ID).First(&sub).Error; err == nil && sub.TransactionSource != "" {
			gc.ReleaseInUseChannelAccount(sub.TransactionSource)
		}
		gc.DB.Model(&userModels.FiatPaymentInvoice{}).Where("id = ?", inv.ID).Update("status", "EXPIRED")
	}
	return nil
}
```

Registered in `main.go`, in the same block as the other background loops (~line 820, before `//setup router`), following the codebase's uniform `for { <work>; time.Sleep(duration) }` idiom (no `time.Ticker` used anywhere in this file):

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

(This sweep is generic to all `FiatPaymentInvoice` rows, not just asset-purchase ones — reasonable, since any stuck `PENDING` invoice of any payment type should expire the same way; only the channel-account release step is specific to asset-purchase invoices, gated by `TokenizedAssetSubscriptionID != ""`.)

## 9. HTTP layer

New endpoint (mirrors `postUsersFiatFlutterwaveHandler`'s auth/registration shape, but drives the two-call subscription flow instead of a single invoice write):

```
POST /v1/users/fiat/tokenized-assets/subscribe
```

Handler `postUsersTokenizedAssetSubscriptionFiatHandler(gc)`:
1. Auth via `middleware.AuthenticationMiddlewareUsingTimestamp()` + `usersDB.GetUserFromPrimarySigner`, same as `postUsersFiatFlutterwaveHandler`.
2. Decode `userModels.TokenizedAssetSubscriptionInput` from the body.
3. Load `ta` (`TokenizedAssetID`) and `subscriberWallet` (`WalletPublicKey` + user), same lookups the existing `SubscribeToTokenizedAsset` HTTP call sites perform (`internal/components/users/controllers/handlers_impl.go:6224`, `:6385` — mirror their pre-fetch pattern).
4. Call `userServices.SubscribeToTokenizedAssetByFiat(user, subscriberWallet, ta, &input, gc)`.
5. Respond:
   - Call 1: `{"paymentInvoiceId": ..., "transaction": ...}` (unsigned-by-buyer XDR for the client to sign; `paymentInvoiceId` just echoes what the client sent, for convenience).
   - Call 2: `{"paymentInvoiceId": ..., "status": "awaiting-payment"}` — no XDR needed back, submission happens later via webhook.

## 10. Env vars

`INTERNAL_BALANCE_ISSUING_SIGNERS` already exists and is already exactly what's needed (`tokenized_assets.go:3179`, `getInternalBalanceIssuingSigners()` — CSV of secret seeds, trimmed, parsed to `keypair.Full`). Reuse it as-is in `generateAssetSubscriptionFiatXdr`; no new env var required. (Aside, not required for this feature: it's also not currently in `main.go`'s `requiredEnvironmentVariables` startup validation list, which is a pre-existing gap from the prior on-demand-mint work, not something this plan needs to fix.)

## 11. File-by-file checklist

| File | Change |
|---|---|
| `internal/components/users/models/tokenization.go` | Add `PaymentAssetCode`, `PaymentAssetIssuer`, `Transaction`, `TransactionSignature`, `TransactionSource` to `TokenizedAssetSubscription`; add `PaymentInvoiceID` to `TokenizedAssetSubscriptionInput`. |
| `internal/components/users/models/fiatPayment.go` | Model fields already added (`TokenizedAssetID`, `WalletAlias`, `WalletPublicKey`); no further schema change — `FiatPaymentInvoice.ID` is set to `PaymentInvoiceID` at creation time and the three new fields are populated from `taSubscription` (§3.3). |
| `internal/components/users/services/tokenized_assets.go` | Add `SubscribeToTokenizedAssetByFiat` and `generateAssetSubscriptionFiatXdr`. |
| `internal/components/users/services/fiatPayment.go` | Add `ExpireStalePaymentInvoices`. |
| `internal/components/users/controllers/handlers_impl.go` | Add `postUsersTokenizedAssetSubscriptionFiatHandler`. |
| `internal/components/users/controllers/main.go` | Register the new route. |
| `internal/components/callbacks/controllers/handlers_impl.go` | Complete the `ASSET PURCHASE` TODO block (~line 771). |
| `main.go` | Extend the `InUseChannelAccounts` startup scan (~line 370-383) to include mid-flight `TokenizedAssetSubscription` rows; register the new expiry goroutine (~line 820). |

## 12. Open questions to resolve before implementation

1. **Invoice/tx_ref correlation** — confirmed: `FiatPaymentInvoice.ID` is always client-supplied (matches the existing `postUsersFiatFlutterwaveHandler` convention), so `PaymentInvoiceID` is **required on both calls** — the handler rejects the request if it's missing, and the same value is used as `TokenizedAssetSubscription.ID`, `FiatPaymentInvoice.ID`, and Flutterwave's `tx_ref` throughout. Still worth confirming `tx_ref` is client-settable in the Flutterwave SDK integration before implementing.
2. **Failure/retry policy** for a webhook-side submission failure (§6.6) — leave `PENDING` indefinitely (relying only on the 2-day expiry sweep) or introduce an explicit `FAILED` status with immediate channel-account release? Needs a product decision.
3. **Signature verification** — should call 2 verify the buyer's signature is well-formed / matches the wallet's signer set before persisting it (saving a bad signature that will only fail much later, async, in the webhook), or is that deferred entirely to submission time as today's `SubmitXdrWithSignature` does? Recommend an early sanity check (e.g. `keypair.ParseAddress` on the recovered signer) but not a full verify, to fail fast on obviously malformed input.
