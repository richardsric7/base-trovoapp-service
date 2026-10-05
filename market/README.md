# market

The contracts Trovo's tokenized assets are issued and sold with.

| Contract | What it is |
|---|---|
| [`TrovoOfferBook`](contracts/src/TrovoOfferBook.sol) | Fixed-price offers between curated tokens. A seller escrows a token in an offer and prices it in each payment token it accepts; a buyer fills any part of it, paying the seller's proceeds recipient and receiving the tokens in one transaction, or nothing happens. |
| [`TokenizedAsset`](contracts/src/TokenizedAsset.sol) | The token of one tokenized asset: a plain mintable, burnable ERC-20 with its own decimals, owned by the asset's issuing Safe - the only account that can mint it. |

Hardhat project in [`contracts/`](contracts) with tests, a Vault-configured
deploy script and a local deployment for end-to-end tests.

## How a tokenized asset is sold

1. The asset's **issuing wallet** (a Safe owned by the tokenization issuing
   profile's key and the minting approvers) mints the supply in one
   operation the approvers sign, and in the same operation escrows the
   amount for sale in an offer on the book, priced in each tokenization
   payment stablecoin (e.g. cNGN) and in the country's internal balance
   token, with proceeds going to the asset's funds holding wallet.
2. A **buyer** asks app-backend to buy. After its checks (KYC, sale window,
   purchase cap), app-backend's authorizer key signs an authorization for
   exactly that purchase - buyer, recipient, amount, maximum payment,
   single-use nonce, deadline - and builds the buyer wallet's operation:
   approve the book for the payment, then `fill`. The buyer signs it as any
   other send.
3. A **fiat** buyer pays the payment provider; once it confirms, the
   internal balance token's minting Safe mints the amount to itself and
   fills the offer with the buyer's wallet as recipient.

The book has no order matching and no price discovery: every fill is at the
seller's price. Users' market-making offers (app-backend
`POST /v1/users/trades`) are offers on the same book, and swaps
(`POST /v1/users/swap`) fill them, cheapest first, each fill authorized
by the platform like a purchase - see app-backend's INTEGRATION.md
"Swaps and market-making offers".

## Safety properties

- Fills need a platform authorization (EIP-712) for exactly that taker,
  recipient, offer, payment token, amount, maximum payment, nonce and
  deadline; each nonce works once. Copying someone's authorization does
  not let anyone else use it.
- An authorizer key can only let someone buy at the seller's own price; it
  can never move a seller's escrow or a buyer's tokens. A leaked key
  bypasses the platform's KYC and cap checks, so rotate it
  (`setAuthorizer`) if that happens.
- Only tokens the owner lists (`setTradable`) can be sold or accepted, so
  only curated regulated assets trade. Delisting a token stops its trading;
  sellers can always cancel and take back what is left, even while the
  book is paused.
- The buyer always pays the price rounded up, never less.
- The owner (the platform admin Safe) can pause fills and list tokens; it
  cannot touch escrows.
- Not upgradeable. **Not externally audited yet - get it audited before
  mainnet.**

## Docs

- [DEPLOYMENT.md](DEPLOYMENT.md) - building, testing, deploying, and the
  local deployment for end-to-end tests.
- [CONFIGURATION.md](CONFIGURATION.md) - the deploy script's settings
  (Vault).
- [INTEGRATION.md](INTEGRATION.md) - the contract interface, how
  app-backend uses it, and what the owner Safe does.
