# market - integration

How the market contracts connect to the rest of Trovo. The contracts live
on the Base blockchain; other projects talk to them by sending transactions
and reading their state.

| Connects to | Direction | What for |
|---|---|---|
| [app-backend](#1-app-backend) | app-backend → contracts | mints assets, signs purchase permissions, builds buyers' transactions, reads offers |
| [The platform admin Safe (owner)](#2-the-platform-admin-safe-owner) | owner → offer book | lists tokens, rotates authorizers, pauses |
| [payout-engine](#3-payout-engine) | engine reads | pays proceeds for tokens held in open offers to the sellers |
| [Trovo Manager (tm-api, tm-web)](#4-trovo-manager) | indirect, through app-backend | registering an asset's token, minting |

---

## 1. app-backend

- **What it is and why:** app-backend is Trovo's main backend. It is the
  only part of Trovo that creates purchase permissions, so every sale goes
  through its checks (KYC, sale window, purchase cap) first.
- **Direction:** app-backend builds transactions that users' wallets (and
  platform Safes) send to the contracts, signs purchase permissions, and
  reads offers.
- **How they connect:** JSON-RPC to a Base node. Purchase permissions are
  EIP-712 signatures, computed in app-backend's `internal/aa/offerbook.go`
  (`FillDigest`, `SignFill`).
- **What app-backend does with it:**
  - **Mint:** the asset's issuing wallet (the seller) creates the tokens and,
    in the same operation, puts the amount for sale in an offer, priced in
    each payment stablecoin and in the internal balance token, with the
    money going to the asset's funds holding wallet. app-backend records the
    offer id from the `OfferCreated` event.
  - **Buying with a stablecoin:** the buyer's wallet sends
    `approve(book, payment)` then `fill(...)`, with a permission for that
    wallet.
  - **Buying with fiat:** after the payment provider confirms, the internal
    balance minting Safe sends `mint(self, payment)`, `approve(book,
    payment)` and `fill(...)` with the buyer's wallet as recipient.
  - **How much is left:** purchase limits use `getOffer(offerId).remaining`.
  - **Market-making offers and swaps:** users' offers (`POST /v1/users/trades`)
    are offers on the same book, and swaps (`POST /v1/users/swap`) fill them,
    cheapest first.
- **Settings on the other side (app-backend):**

  | Parameter | Example | How to get it |
  |---|---|---|
  | `OFFER_BOOK_ADDRESS` | `0x5FbDB2315678afecb367f032d93F642f64180aa3` | printed by the deploy script, and in `contracts/deployments/<chainId>.json` |
  | `OFFER_AUTHORIZER_PRIVATE_KEY` | `0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a` (test key) | the key whose address is in [`AUTHORIZER_ADDRESSES`](CONFIGURATION.md#authorizer_addresses) |

  See [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md) and
  [app-backend/INTEGRATION.md](../app-backend/INTEGRATION.md#tokenized-assets-token-issuing-and-distribution-wallets-sale-offer).
- **How to check it works:** on a test network, buy a small amount of a test
  asset in the app; the purchase appears as a `Filled` event on the offer
  book in basescan.
- **When it is down or misconfigured:** nothing can be bought or minted. A
  wrong authorizer key makes every purchase fail with an authorization
  error; the contracts and everyone's tokens are unaffected.

### The offer book's functions

| Function | Who calls it | What it does |
|---|---|---|
| `createOffer(sellToken, amount, proceedsRecipient, paymentTokens[], prices[])` | seller | Escrows `amount` of `sellToken` (what actually arrives) in a new offer; returns its id (`OfferCreated`). |
| `setPrice(offerId, paymentToken, price)` | seller | Sets a price; a price of 0 stops accepting that token. |
| `cancelOffer(offerId)` | seller | Closes the offer and returns what is left. Works even while paused or delisted. |
| `fill(FillRequest f, bytes authorization)` | buyer | Pays at most `f.maxPayment` of `f.paymentToken` to the seller's proceeds recipient and sends `f.amount` of the token to `f.recipient` (`Filled`). The buyer must first approve the book to take the payment. |
| `getOffer`, `priceOf`, `quote`, `fillDigest`, `nonceUsed`, `tradable`, `authorizers` | anyone | Read-only views. |
| `setTradable(token, bool)`, `setAuthorizer(addr, bool)`, `pause()`, `unpause()` | owner | Administration (see below). |

**Prices** are exact fractions in the tokens' smallest units: buying
`amount` units costs `ceil(amount * num / den)` units of the payment token,
so the buyer never pays less than the price. app-backend turns a price per
whole token into one fraction per payment token, treating each payment
token as worth exactly one unit of the quote currency.

**Purchase permissions:** `FillRequest` is `(offerId, paymentToken, amount,
maxPayment, recipient, nonce, deadline)`. The permission is an EIP-712
signature (domain name `TrovoOfferBook`, version `1`, the chain id and the
book's address) over:

```
Fill(uint256 offerId,address taker,address recipient,address paymentToken,uint256 amount,uint256 maxPayment,uint256 nonce,uint256 deadline)
```

where `taker` is the address calling `fill`. Each nonce works once, and a
permission is only valid for that exact taker, so copying one is useless.

## 2. The platform admin Safe (owner)

- **What it is and why:** a multi-signature wallet controlled by Trovo's
  platform administrators. It owns the offer book after deployment. It can
  manage the book but can never move anyone's escrowed tokens.
- **Direction:** the owner sends administration transactions to the book.
- **How they connect:** transactions made in Safe{Wallet}
  (<https://app.safe.global>) with **New transaction → Transaction Builder**,
  signed by enough owners and then executed.
- **What it does:**
  - **Lists each tokenized asset's token** (`setTradable(token, true)`)
    after the token is deployed and registered and before minting.
    app-backend refuses to mint an unlisted token (`error-token-not-listed`).
    Payment tokens that are not listed are left out of an asset's offer.
  - **Rotates authorizers:** `setAuthorizer(new, true)`, switch app-backend's
    `OFFER_AUTHORIZER_PRIVATE_KEY`, then `setAuthorizer(old, false)`. If an
    authorizer key ever leaks, rotate it straight away: a leaked key lets
    someone skip Trovo's KYC and cap checks (though only ever at the
    seller's own price).
  - **Pauses** purchases in an incident (`pause()` / `unpause()`). Sellers
    can still cancel.
- **Settings:** [`OWNER_ADDRESS`](CONFIGURATION.md#owner_address) at
  deployment (example `0xA11CE00000000000000000000000000000000001`).
- **How to check it works:** `owner()` on basescan's **Read Contract** tab
  returns the Safe's address.

## 3. payout-engine

- **What it is and why:** pays tokenized assets' proceeds to their holders.
  Tokens in an open sell offer sit in the offer book but still belong to
  the seller, so the engine credits them to the seller.
- **Direction:** payout-engine reads; it never sends anything to the book.
- **How they connect:** the open offers in app-backend's database
  (`offer_book_offers`, kept up to date by app-backend) and the book's
  address.
- **Settings on the other side:** payout-engine's
  [`OFFER_BOOK_ADDRESS`](../payout-engine/CONFIGURATION.md#offer_book_address),
  the same value as app-backend's.
- **When it is unset:** tokens in open offers earn nothing in payouts.

## 4. Trovo Manager

- **What it is and why:** the admin dashboard. Admins register each asset's
  token contract and run minting there.
- **Direction:** Trovo Manager → app-backend → contracts.
- **How it works:** after deploying an asset's token
  ([DEPLOYMENT.md step 5](DEPLOYMENT.md#5-for-each-new-tokenized-asset)),
  the contract address is entered on the asset's tokenization page. That
  calls app-backend's `PUT /v1/trovo-manager/tokenization/contract/:tid`,
  which checks on-chain that the contract exists, its symbol is the asset
  code, `decimals()` works, its supply is 0 and the issuing wallet can mint
  (`400 error-invalid-token-contract` otherwise).
