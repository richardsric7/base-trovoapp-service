const { expect } = require("chai");
const { ethers } = require("hardhat");
const { loadFixture, time } = require("@nomicfoundation/hardhat-toolbox/network-helpers");

const PERIOD = 7 * 24 * 3600;
const SENTINEL = "0x0000000000000000000000000000000000000001";

// Trovo's use of Candide's module, against real Safe v1.4.1 wallets: the
// platform's recovery key is a wallet's only guardian (threshold 1). It can
// start replacing the wallet's owners; the owners can cancel during the
// recovery period; after it anyone can finalize. It can never move funds.
async function deployFixture() {
  const [deployer, user, newUser, guardian, cosigner, stranger, approver] = await ethers.getSigners();
  const singleton = await ethers.deployContract("SafeL2");
  const factory = await ethers.deployContract("SafeProxyFactory");
  const handler = await ethers.deployContract("CompatibilityFallbackHandler");
  const module = await ethers.deployContract("SocialRecoveryModule", [PERIOD]);
  const token = await ethers.deployContract("MockERC20");

  let salt = 0;
  const newSafe = async (owners, threshold = 1) => {
    const init = singleton.interface.encodeFunctionData("setup", [
      owners.map((o) => o.address), threshold, ethers.ZeroAddress, "0x", await handler.getAddress(), ethers.ZeroAddress, 0, ethers.ZeroAddress,
    ]);
    const tx = await factory.createProxyWithNonce(singleton, init, salt++);
    const rc = await tx.wait();
    const ev = rc.logs.map((l) => { try { return factory.interface.parseLog(l); } catch (_) { return null; } }).find((e) => e && e.name === "ProxyCreation");
    return ethers.getContractAt("SafeL2", ev.args.proxy);
  };
  // the owner `by` executes a call from the Safe (pre-validated signature)
  const exec = async (safe, by, to, data) => {
    const sig = ethers.concat([ethers.zeroPadValue(by.address, 32), ethers.ZeroHash, "0x01"]);
    return safe.connect(by).execTransaction(to, 0, data, 0, 0, 0, 0, ethers.ZeroAddress, ethers.ZeroAddress, sig);
  };
  const enableRecovery = async (safe, by) => {
    await exec(safe, by, await safe.getAddress(), safe.interface.encodeFunctionData("enableModule", [await module.getAddress()]));
    await exec(safe, by, await module.getAddress(), module.interface.encodeFunctionData("addGuardianWithThreshold", [guardian.address, 1]));
  };

  const safe = await newSafe([user]);
  await token.mint(await safe.getAddress(), 1000n);
  return { deployer, user, newUser, guardian, cosigner, stranger, approver, module, token, safe, newSafe, exec, enableRecovery };
}

describe("SocialRecoveryModule (Trovo's use)", function () {
  it("replaces the lost key after the recovery period, keeping the funds", async function () {
    const { module, safe, user, newUser, guardian, stranger, token, exec, enableRecovery } = await loadFixture(deployFixture);
    await enableRecovery(safe, user);
    expect(await module.isGuardian(safe, guardian.address)).to.equal(true);
    expect(await module.threshold(safe)).to.equal(1n);

    const nonce = await module.nonce(safe);
    await module.connect(guardian).confirmRecovery(safe, [newUser.address], 1, nonce, true);
    const req = await module.getRecoveryRequest(safe);
    expect(req.newOwners).to.deep.equal([newUser.address]);
    expect(req.executableAt).to.equal(BigInt(await time.latest()) + BigInt(PERIOD));

    await expect(module.connect(stranger).finalizeRecovery(safe)).to.be.revertedWith("SM: recovery period still pending");
    await time.increase(PERIOD);
    await module.connect(stranger).finalizeRecovery(safe); // anyone may finalize
    expect(await safe.getOwners()).to.deep.equal([newUser.address]);
    expect(await safe.getThreshold()).to.equal(1n);

    // the new key controls the wallet and its funds; the old one no longer does
    await exec(safe, newUser, await token.getAddress(), token.interface.encodeFunctionData("transfer", [stranger.address, 10n]));
    expect(await token.balanceOf(stranger.address)).to.equal(10n);
    await expect(exec(safe, user, await token.getAddress(), token.interface.encodeFunctionData("transfer", [user.address, 1n]))).to.be.reverted;
    // the guardian stays, so the wallet remains covered
    expect(await module.isGuardian(safe, guardian.address)).to.equal(true);
  });

  it("lets the owner cancel a recovery during the period", async function () {
    const { module, safe, user, newUser, guardian, stranger, exec, enableRecovery } = await loadFixture(deployFixture);
    await enableRecovery(safe, user);
    await module.connect(guardian).confirmRecovery(safe, [stranger.address], 1, await module.nonce(safe), true);
    await exec(safe, user, await module.getAddress(), module.interface.encodeFunctionData("cancelRecovery"));
    expect((await module.getRecoveryRequest(safe)).executableAt).to.equal(0n);
    await time.increase(PERIOD);
    await expect(module.finalizeRecovery(safe)).to.be.revertedWith("SM: no ongoing recovery");
    expect(await safe.getOwners()).to.deep.equal([user.address]);
    // a new request needs the new nonce
    await expect(module.connect(guardian).confirmRecovery(safe, [newUser.address], 1, 0, true)).to.be.revertedWith("SM: invalid nonce");
  });

  it("gives the recovery key no way to move funds", async function () {
    const { module, safe, user, guardian, token, enableRecovery } = await loadFixture(deployFixture);
    await enableRecovery(safe, user);
    // not an owner: cannot execute Safe transactions
    const sig = ethers.concat([ethers.zeroPadValue(guardian.address, 32), ethers.ZeroHash, "0x01"]);
    await expect(
      safe.connect(guardian).execTransaction(token, 0, token.interface.encodeFunctionData("transfer", [guardian.address, 1n]), 0, 0, 0, 0, ethers.ZeroAddress, ethers.ZeroAddress, sig),
    ).to.be.reverted;
    // cannot make itself an owner through recovery
    await expect(module.connect(guardian).confirmRecovery(safe, [guardian.address], 1, 0, true)).to.be.revertedWith("SM: new owner cannot be guardian");
    // cannot cancel or change guardians (only the wallet can)
    await expect(module.connect(guardian).revokeGuardianWithThreshold(SENTINEL, guardian.address, 0)).to.be.reverted;
    expect(await token.balanceOf(guardian.address)).to.equal(0n);
  });

  it("keeps a co-signer when recovering a two-owner wallet", async function () {
    const { module, newSafe, user, newUser, cosigner, guardian, enableRecovery } = await loadFixture(deployFixture);
    const safe = await newSafe([user, cosigner], 1);
    await enableRecovery(safe, user);
    await module.connect(guardian).confirmRecovery(safe, [newUser.address, cosigner.address], 1, await module.nonce(safe), true);
    await time.increase(PERIOD);
    await module.finalizeRecovery(safe);
    // the same owner set (Safe puts added owners first)
    expect([...(await safe.getOwners())].sort()).to.deep.equal([newUser.address, cosigner.address].sort());
    expect(await safe.getThreshold()).to.equal(1n);
  });

  it("stops covering a wallet once its guardian is revoked (approvers added)", async function () {
    const { module, safe, user, guardian, newUser, exec, enableRecovery } = await loadFixture(deployFixture);
    await enableRecovery(safe, user);
    await exec(safe, user, await module.getAddress(), module.interface.encodeFunctionData("revokeGuardianWithThreshold", [SENTINEL, guardian.address, 0]));
    expect(await module.guardiansCount(safe)).to.equal(0n);
    await expect(module.connect(guardian).confirmRecovery(safe, [newUser.address], 1, await module.nonce(safe), true)).to.be.revertedWith("SM: sender not a guardian");
  });

  it("refuses a guardian that is an owner, and disabled modules", async function () {
    const { module, newSafe, user, guardian, exec } = await loadFixture(deployFixture);
    const safe = await newSafe([user]);
    // module not enabled
    await expect(exec(safe, user, await module.getAddress(), module.interface.encodeFunctionData("addGuardianWithThreshold", [guardian.address, 1]))).to.be.reverted;
    await exec(safe, user, await safe.getAddress(), safe.interface.encodeFunctionData("enableModule", [await module.getAddress()]));
    await expect(exec(safe, user, await module.getAddress(), module.interface.encodeFunctionData("addGuardianWithThreshold", [user.address, 1]))).to.be.reverted;
  });
});
