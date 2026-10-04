const { expect } = require("chai");
const { ethers } = require("hardhat");
const { loadFixture, time } = require("@nomicfoundation/hardhat-toolbox/network-helpers");

// 1500 NGN per whole asset token: the asset has 7 decimals, cNGN 6, the
// internal balance token 18. Price = payment base units per asset base units.
const ASSET = (n) => BigInt(n) * 10n ** 7n;
const CNGN = (n) => BigInt(n) * 10n ** 6n;
const PRICE_CNGN = { num: 1500n * 10n ** 6n, den: 10n ** 7n };
const PRICE_INTERNAL = { num: 1500n * 10n ** 18n, den: 10n ** 7n };

const FILL_TYPES = {
  Fill: [
    { name: "offerId", type: "uint256" },
    { name: "taker", type: "address" },
    { name: "recipient", type: "address" },
    { name: "paymentToken", type: "address" },
    { name: "amount", type: "uint256" },
    { name: "maxPayment", type: "uint256" },
    { name: "nonce", type: "uint256" },
    { name: "deadline", type: "uint256" },
  ],
};

async function deployFixture() {
  const [owner, authorizer, seller, buyer, recipient, proceeds, stranger, issuer] = await ethers.getSigners();

  const book = await ethers.deployContract("TrovoOfferBook", [owner.address, [authorizer.address]]);
  const asset = await ethers.deployContract("TokenizedAsset", ["Farm Fund", "FARM", 7, issuer.address]);
  const cngn = await ethers.deployContract("MockERC20", ["cNGN", "CNGN", 6]);
  const internal = await ethers.deployContract("MockERC20", ["Internal NGN", "NGNI", 18]);
  const unlisted = await ethers.deployContract("MockERC20", ["Unlisted", "UNL", 6]);

  for (const t of [asset, cngn, internal]) {
    await book.connect(owner).setTradable(t, true);
  }
  await asset.connect(issuer).mint(seller.address, ASSET(1000));
  await asset.connect(seller).approve(book, ethers.MaxUint256);
  await cngn.mint(buyer.address, CNGN(1_000_000));
  await cngn.connect(buyer).approve(book, ethers.MaxUint256);

  const domain = {
    name: "TrovoOfferBook",
    version: "1",
    chainId: (await ethers.provider.getNetwork()).chainId,
    verifyingContract: await book.getAddress(),
  };
  let nextNonce = 1n;
  // an authorization for `taker` to make fill `f`
  const authorize = async (f, taker, signer = authorizer) =>
    signer.signTypedData(domain, FILL_TYPES, { ...f, taker: taker.address ?? taker });
  const request = async (overrides = {}) => ({
    offerId: 1n,
    paymentToken: await cngn.getAddress(),
    amount: ASSET(10),
    maxPayment: CNGN(15_000),
    recipient: recipient.address,
    nonce: nextNonce++,
    deadline: BigInt(await time.latest()) + 600n,
    ...overrides,
  });

  return { owner, authorizer, seller, buyer, recipient, proceeds, stranger, issuer, book, asset, cngn, internal, unlisted, authorize, request };
}

async function withOfferFixture() {
  const f = await deployFixture();
  await f.book
    .connect(f.seller)
    .createOffer(f.asset, ASSET(100), f.proceeds.address, [f.cngn, f.internal], [PRICE_CNGN, PRICE_INTERNAL]);
  return f;
}

describe("TrovoOfferBook", function () {
  describe("creating offers", function () {
    it("escrows the tokens and records the prices", async function () {
      const { book, seller, proceeds, asset, cngn, internal } = await loadFixture(deployFixture);
      await expect(book.connect(seller).createOffer(asset, ASSET(100), proceeds.address, [cngn, internal], [PRICE_CNGN, PRICE_INTERNAL]))
        .to.emit(book, "OfferCreated")
        .withArgs(1n, seller.address, await asset.getAddress(), ASSET(100), proceeds.address)
        .and.to.emit(book, "PriceSet")
        .withArgs(1n, await cngn.getAddress(), PRICE_CNGN.num, PRICE_CNGN.den);
      expect(await asset.balanceOf(book)).to.equal(ASSET(100));
      expect(await asset.balanceOf(seller)).to.equal(ASSET(900));
      const o = await book.getOffer(1n);
      expect([o.seller, o.remaining, o.open]).to.deep.equal([seller.address, ASSET(100), true]);
      expect(await book.quote(1n, cngn, ASSET(1))).to.equal(CNGN(1500));
      expect(await book.quote(1n, internal, ASSET(1))).to.equal(1500n * 10n ** 18n);
      expect(await book.nextOfferId()).to.equal(2n);
    });

    it("trades only tradable tokens, at real prices", async function () {
      const { book, seller, proceeds, asset, cngn, unlisted } = await loadFixture(deployFixture);
      const c = book.connect(seller);
      await expect(c.createOffer(unlisted, 1n, proceeds.address, [cngn], [PRICE_CNGN])).to.be.revertedWithCustomError(book, "NotTradable");
      await expect(c.createOffer(asset, 1n, proceeds.address, [unlisted], [PRICE_CNGN])).to.be.revertedWithCustomError(book, "NotTradable");
      await expect(c.createOffer(asset, 1n, proceeds.address, [cngn], [{ num: 0n, den: 1n }])).to.be.revertedWithCustomError(book, "InvalidPrice");
      await expect(c.createOffer(asset, 1n, proceeds.address, [cngn], [{ num: 1n, den: 0n }])).to.be.revertedWithCustomError(book, "InvalidPrice");
      await expect(c.createOffer(asset, 1n, proceeds.address, [asset], [PRICE_CNGN])).to.be.revertedWithCustomError(book, "InvalidPrice");
      await expect(c.createOffer(asset, 1n, proceeds.address, [cngn], [])).to.be.revertedWithCustomError(book, "InvalidAmount");
      await expect(c.createOffer(asset, 0n, proceeds.address, [cngn], [PRICE_CNGN])).to.be.revertedWithCustomError(book, "InvalidAmount");
      await expect(c.createOffer(asset, 1n, ethers.ZeroAddress, [cngn], [PRICE_CNGN])).to.be.revertedWithCustomError(book, "InvalidAddress");
    });

    it("records what actually arrives from a fee-on-transfer token", async function () {
      const { book, owner, seller, proceeds, cngn } = await loadFixture(deployFixture);
      const fee = await ethers.deployContract("FeeOnTransferERC20");
      await book.connect(owner).setTradable(fee, true);
      await fee.mint(seller.address, 10_000n);
      await fee.connect(seller).approve(book, 10_000n);
      await book.connect(seller).createOffer(fee, 10_000n, proceeds.address, [cngn], [PRICE_CNGN]);
      expect((await book.getOffer(1n)).remaining).to.equal(9_900n);
    });
  });

  describe("filling", function () {
    it("pays the proceeds recipient and delivers to the recipient in one step", async function () {
      const { book, buyer, recipient, proceeds, asset, cngn, authorize, request } = await loadFixture(withOfferFixture);
      const f = await request();
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer)))
        .to.emit(book, "Filled")
        .withArgs(1n, buyer.address, recipient.address, await cngn.getAddress(), ASSET(10), CNGN(15_000), f.nonce);
      expect(await cngn.balanceOf(proceeds)).to.equal(CNGN(15_000));
      expect(await asset.balanceOf(recipient)).to.equal(ASSET(10));
      expect((await book.getOffer(1n)).remaining).to.equal(ASSET(90));
      expect(await book.nonceUsed(f.nonce)).to.equal(true);
    });

    it("rounds the payment up, never in the buyer's favour", async function () {
      const { book, seller, buyer, proceeds, cngn, authorize, request } = await loadFixture(withOfferFixture);
      // 1.5 base units of cNGN per base unit of the asset: 1 costs 2
      await book.connect(seller).setPrice(1n, cngn, { num: 3n, den: 2n });
      let f = await request({ amount: 1n, maxPayment: 1n });
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer))).to.be.revertedWithCustomError(book, "PaymentTooHigh");
      f = await request({ amount: 1n, maxPayment: 2n });
      await book.connect(buyer).fill(f, await authorize(f, buyer));
      expect(await cngn.balanceOf(proceeds)).to.equal(2n);
    });

    it("lets the platform pay in internal balance for a buyer (fiat purchase)", async function () {
      const { book, issuer, recipient, proceeds, asset, internal, authorize, request } = await loadFixture(withOfferFixture);
      await internal.mint(issuer.address, 15_000n * 10n ** 18n);
      await internal.connect(issuer).approve(book, ethers.MaxUint256);
      const f = await request({ paymentToken: await internal.getAddress(), maxPayment: 15_000n * 10n ** 18n });
      await book.connect(issuer).fill(f, await authorize(f, issuer));
      expect(await asset.balanceOf(recipient)).to.equal(ASSET(10));
      expect(await internal.balanceOf(proceeds)).to.equal(15_000n * 10n ** 18n);
    });

    it("refuses fills the platform did not authorize exactly", async function () {
      const { book, buyer, stranger, authorize, request } = await loadFixture(withOfferFixture);
      const f = await request();
      const sig = await authorize(f, buyer);
      // another key
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer, stranger))).to.be.revertedWithCustomError(book, "Unauthorized");
      // someone else submitting the buyer's authorization
      await expect(book.connect(stranger).fill(f, sig)).to.be.revertedWithCustomError(book, "Unauthorized");
      // any field changed
      for (const change of [{ amount: ASSET(9) }, { maxPayment: CNGN(20_000) }, { recipient: stranger.address }, { deadline: f.deadline + 1n }]) {
        await expect(book.connect(buyer).fill({ ...f, ...change }, sig)).to.be.revertedWithCustomError(book, "Unauthorized");
      }
      await expect(book.connect(buyer).fill(f, "0x1234")).to.be.revertedWithCustomError(book, "Unauthorized");
    });

    it("uses each authorization once, before its deadline", async function () {
      const { book, buyer, authorize, request } = await loadFixture(withOfferFixture);
      const f = await request();
      const sig = await authorize(f, buyer);
      await book.connect(buyer).fill(f, sig);
      await expect(book.connect(buyer).fill(f, sig)).to.be.revertedWithCustomError(book, "NonceUsed");

      const late = await request();
      const lateSig = await authorize(late, buyer);
      await time.increaseTo(late.deadline + 1n);
      await expect(book.connect(buyer).fill(late, lateSig)).to.be.revertedWithCustomError(book, "AuthorizationExpired");
    });

    it("refuses a payment above the buyer's maximum, an unaccepted token or more than is left", async function () {
      const { book, buyer, unlisted, internal, owner, authorize, request } = await loadFixture(withOfferFixture);
      let f = await request({ maxPayment: CNGN(14_999) });
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer))).to.be.revertedWithCustomError(book, "PaymentTooHigh");
      f = await request({ amount: ASSET(101), maxPayment: CNGN(1_000_000) });
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer))).to.be.revertedWithCustomError(book, "InsufficientRemaining");
      f = await request({ amount: 0n });
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer))).to.be.revertedWithCustomError(book, "InvalidAmount");
      f = await request({ recipient: ethers.ZeroAddress });
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer))).to.be.revertedWithCustomError(book, "InvalidAddress");
      await book.connect(owner).setTradable(unlisted, true);
      f = await request({ paymentToken: await unlisted.getAddress() });
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer))).to.be.revertedWithCustomError(book, "NotAccepted");
      f = await request({ offerId: 7n });
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer))).to.be.revertedWithCustomError(book, "OfferClosed");
      expect(await book.quote(1n, internal, 0n)).to.equal(0n);
    });

    it("stops trading a delisted token, and while paused", async function () {
      const { book, owner, buyer, asset, cngn, authorize, request } = await loadFixture(withOfferFixture);
      await book.connect(owner).setTradable(cngn, false);
      let f = await request();
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer))).to.be.revertedWithCustomError(book, "NotTradable");
      await book.connect(owner).setTradable(cngn, true);
      await book.connect(owner).setTradable(asset, false);
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer))).to.be.revertedWithCustomError(book, "NotTradable");
      await book.connect(owner).setTradable(asset, true);

      await book.connect(owner).pause();
      f = await request();
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer))).to.be.revertedWithCustomError(book, "EnforcedPause");
      await book.connect(owner).unpause();
      await book.connect(buyer).fill(f, await authorize(f, buyer));
    });

    it("sells out exactly", async function () {
      const { book, buyer, recipient, asset, authorize, request } = await loadFixture(withOfferFixture);
      const f = await request({ amount: ASSET(100), maxPayment: CNGN(150_000) });
      await book.connect(buyer).fill(f, await authorize(f, buyer));
      expect((await book.getOffer(1n)).remaining).to.equal(0n);
      expect(await asset.balanceOf(recipient)).to.equal(ASSET(100));
      const more = await request({ amount: 1n });
      await expect(book.connect(buyer).fill(more, await authorize(more, buyer))).to.be.revertedWithCustomError(book, "InsufficientRemaining");
    });
  });

  describe("managing offers", function () {
    it("lets only the seller reprice or withdraw a price", async function () {
      const { book, seller, stranger, buyer, cngn, authorize, request } = await loadFixture(withOfferFixture);
      await expect(book.connect(stranger).setPrice(1n, cngn, PRICE_CNGN)).to.be.revertedWithCustomError(book, "NotSeller");
      await book.connect(seller).setPrice(1n, cngn, { num: 2n * PRICE_CNGN.num, den: PRICE_CNGN.den });
      expect(await book.quote(1n, cngn, ASSET(1))).to.equal(CNGN(3000));
      await book.connect(seller).setPrice(1n, cngn, { num: 0n, den: 0n });
      const f = await request();
      await expect(book.connect(buyer).fill(f, await authorize(f, buyer))).to.be.revertedWithCustomError(book, "NotAccepted");
    });

    it("returns what is left on cancel, even while paused or delisted", async function () {
      const { book, owner, seller, stranger, buyer, asset, authorize, request } = await loadFixture(withOfferFixture);
      const f = await request();
      await book.connect(buyer).fill(f, await authorize(f, buyer));
      await expect(book.connect(stranger).cancelOffer(1n)).to.be.revertedWithCustomError(book, "NotSeller");
      await book.connect(owner).pause();
      await book.connect(owner).setTradable(asset, false);
      await expect(book.connect(seller).cancelOffer(1n)).to.emit(book, "OfferCancelled").withArgs(1n, ASSET(90));
      expect(await asset.balanceOf(seller)).to.equal(ASSET(990));
      expect(await asset.balanceOf(book)).to.equal(0n);
      await expect(book.connect(seller).cancelOffer(1n)).to.be.revertedWithCustomError(book, "OfferClosed");
      await expect(book.connect(seller).setPrice(1n, asset, PRICE_CNGN)).to.be.revertedWithCustomError(book, "OfferClosed");
    });

    it("keeps offers' escrows apart", async function () {
      const { book, seller, stranger, proceeds, asset, cngn, issuer } = await loadFixture(withOfferFixture);
      await asset.connect(issuer).mint(stranger.address, ASSET(5));
      await asset.connect(stranger).approve(book, ethers.MaxUint256);
      await book.connect(stranger).createOffer(asset, ASSET(5), proceeds.address, [cngn], [PRICE_CNGN]);
      await book.connect(stranger).cancelOffer(2n);
      expect(await asset.balanceOf(stranger)).to.equal(ASSET(5));
      expect((await book.getOffer(1n)).remaining).to.equal(ASSET(100));
      expect(await asset.balanceOf(book)).to.equal(ASSET(100));
      await expect(book.connect(stranger).cancelOffer(1n)).to.be.revertedWithCustomError(book, "NotSeller");
      void seller;
    });
  });

  describe("administration", function () {
    it("is owner-only", async function () {
      const { book, owner, stranger, authorizer, cngn } = await loadFixture(deployFixture);
      for (const call of [
        () => book.connect(stranger).setTradable(cngn, false),
        () => book.connect(stranger).setAuthorizer(stranger.address, true),
        () => book.connect(stranger).pause(),
        () => book.connect(stranger).unpause(),
      ]) {
        await expect(call()).to.be.revertedWithCustomError(book, "OwnableUnauthorizedAccount");
      }
      await expect(book.connect(owner).setAuthorizer(authorizer.address, false)).to.emit(book, "AuthorizerSet").withArgs(authorizer.address, false);
      expect(await book.authorizers(authorizer.address)).to.equal(false);
      await expect(book.connect(owner).setTradable(ethers.ZeroAddress, true)).to.be.revertedWithCustomError(book, "InvalidAddress");
      await expect(book.connect(owner).setAuthorizer(ethers.ZeroAddress, true)).to.be.revertedWithCustomError(book, "InvalidAddress");

      await book.connect(owner).transferOwnership(stranger.address);
      expect(await book.owner()).to.equal(stranger.address);
    });

    it("a removed authorizer's signatures stop working", async function () {
      const { book, owner, authorizer, buyer, authorize, request } = await loadFixture(withOfferFixture);
      const f = await request();
      const sig = await authorize(f, buyer);
      await book.connect(owner).setAuthorizer(authorizer.address, false);
      await expect(book.connect(buyer).fill(f, sig)).to.be.revertedWithCustomError(book, "Unauthorized");
    });

    it("exposes the digest authorizers sign", async function () {
      const { book, buyer, authorizer, authorize, request } = await loadFixture(withOfferFixture);
      const f = await request();
      const digest = await book.fillDigest(f, buyer.address);
      const sig = await authorize(f, buyer);
      expect(ethers.recoverAddress(digest, sig)).to.equal(authorizer.address);
    });
  });

  describe("TokenizedAsset", function () {
    it("only its owner (the issuing Safe) mints; holders can burn", async function () {
      const { asset, issuer, seller, stranger } = await loadFixture(deployFixture);
      expect(await asset.decimals()).to.equal(7n);
      await expect(asset.connect(stranger).mint(stranger.address, 1n)).to.be.revertedWithCustomError(asset, "OwnableUnauthorizedAccount");
      await asset.connect(seller).burn(ASSET(1));
      expect(await asset.totalSupply()).to.equal(ASSET(999));
      void issuer;
    });
  });
});
