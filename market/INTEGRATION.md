# market - integration

## TrovoOfferBook interface

| Function | Who | What it does |
|---|---|---|
| `createOffer(sellToken, amount, proceedsRecipient, paymentTokens[], prices[])` | seller | Escrows `amount` of `sellToken` (what actually arrives) in a new offer; returns its id (`OfferCreated`). |
| `setPrice(offerId, paymentToken, price)` | seller | Sets a price; `num = 0` stops accepting that token. |
| `cancelOffer(offerId)` | seller | Closes the offer and returns what is left. Works while paused or delisted. |
| `fill(FillRequest f, bytes authorization)` | buyer | Pays at most `f.maxPayment` of `f.paymentToken` to the proceeds recipient and sends `f.amount` of the sell token to `f.recipient` (`Filled`). The caller must have approved the book for the payment. |
| `getOffer`, `priceOf`, `quote`, `fillDigest`, `nonceUsed`, `tradable`, `authorizers` | anyone | Views. |
| `setTradable(token, bool)`, `setAuthorizer(addr, bool)`, `pause()`, `unpause()` | owner | Administration. |

Prices are exact fractions in base units: buying `amount` base units of
the sell token costs `ceil(amount * num / den)` base units of the payment
token. app-backend converts a price per whole token in the quote currency
into one fraction per payment token, taking each token as 1:1 with the
quote currency.

`FillRequest` is `(offerId, paymentToken, amount, maxPayment, recipient,
nonce, deadline)`. The authorization is an EIP-712 signature (domain
`TrovoOfferBook`, version `1`, chain id, book address) over

```
Fill(uint256 offerId,address taker,address recipient,address paymentToken,uint256 amount,uint256 maxPayment,uint256 nonce,uint256 deadline)
```

where `taker` is the address calling `fill`. app-backend computes it in
`internal/aa/offerbook.go` (`FillDigest`, `SignFill`).

## How app-backend uses it

- **Mint**: the asset's issuing wallet (seller) escrows the amount for
  sale with prices in each tokenization payment stablecoin and the
  internal balance token, proceeds to the funds holding wallet. When the
  mint is mined, app-backend records the offer id from `OfferCreated`.
- **Stablecoin purchase**: the buyer wallet's operation is
  `approve(book, payment)` + `fill(...)`, authorized for that wallet.
- **Fiat purchase**: the internal balance minting Safe's transaction is
  `mint(self, payment)` + `approve(book, payment)` + `fill(...)` with the
  buyer's wallet as recipient, authorized for the Safe.
- **Liquidity**: purchase limits and the app's remaining amount come from
  `getOffer(offerId).remaining`.

See [app-backend/INTEGRATION.md](../app-backend/INTEGRATION.md#tokenized-assets-token-issuing-and-distribution-wallets-sale-offer).

## What the owner Safe does

- **List each tokenized asset's token** (`setTradable(token, true)`) after
  it is deployed and registered, before minting - app-backend refuses to
  build a mint for an unlisted token (`error-token-not-listed`). Payment
  tokens not listed are left out of an asset's offer.
- **Rotate authorizers**: add the new address, switch app-backend's key,
  remove the old one.
- **Pause** fills in an incident; sellers can still cancel.
