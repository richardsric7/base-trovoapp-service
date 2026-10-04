const { expect } = require("chai");
const { ethers } = require("hardhat");
const { loadFixture, time } = require("@nomicfoundation/hardhat-toolbox/network-helpers");

const { AbiCoder, concat, getBytes, keccak256, toBeHex, zeroPadValue, MaxUint256 } = ethers;

const USDC = (n) => BigInt(Math.round(n * 1e6));
const RATE = USDC(3000); // 3000 USDC (6 decimals) per ETH, spread included
const MIN_RATE = USDC(1000);
const MAX_RATE = USDC(10000);
const POST_OP_OVERHEAD = 40000n;

const GAS = {
  verification: 300000n,
  call: 200000n,
  preVerification: 60000n,
  pmVerification: 200000n,
  postOp: 80000n,
  maxPriority: ethers.parseUnits("1", "gwei"),
  maxFee: ethers.parseUnits("5", "gwei"),
};

const pack128 = (hi, lo) => zeroPadValue(toBeHex((BigInt(hi) << 128n) | BigInt(lo)), 32);

async function deployFixture() {
  const [deployer, owner, quoteSigner, pauser, walletOwner, beneficiary, recipient, stranger] = await ethers.getSigners();

  const entryPoint = await ethers.deployContract("EntryPoint");
  const factory = await ethers.deployContract("SimpleAccountFactory", [entryPoint]);
  const usdc = await ethers.deployContract("MockERC20", ["USD Coin", "USDC", 6]);
  const other = await ethers.deployContract("MockERC20", ["Other", "OTH", 6]);

  const paymaster = await ethers.deployContract("TrovoTokenPaymaster", [
    entryPoint,
    deployer,
    [quoteSigner.address],
    pauser.address,
    POST_OP_OVERHEAD,
  ]);
  await paymaster.addStake(86400, { value: ethers.parseEther("1") });
  await paymaster.deposit({ value: ethers.parseEther("10") });
  await paymaster.setToken(usdc, true, MIN_RATE, MAX_RATE);
  await paymaster.transferOwnership(owner.address);

  // a deployed wallet that approved the paymaster once
  await factory.createAccount(walletOwner.address, 0);
  const account = await ethers.getContractAt("SimpleAccount", await factory["getAddress(address,uint256)"](walletOwner.address, 0));
  await usdc.mint(account, USDC(100));
  await account
    .connect(walletOwner)
    .execute(usdc, 0, usdc.interface.encodeFunctionData("approve", [await paymaster.getAddress(), MaxUint256]));

  return { deployer, owner, quoteSigner, pauser, walletOwner, beneficiary, recipient, stranger, entryPoint, factory, usdc, other, paymaster, account };
}

function transferCall(account, token, to, amount) {
  return account.interface.encodeFunctionData("execute", [
    token.target,
    0,
    token.interface.encodeFunctionData("transfer", [to, amount]),
  ]);
}

async function buildOp(f, { sender, initCode = "0x", callData }) {
  return {
    sender,
    nonce: await f.entryPoint.getNonce(sender, 0),
    initCode,
    callData,
    accountGasLimits: pack128(GAS.verification, GAS.call),
    preVerificationGas: GAS.preVerification,
    gasFees: pack128(GAS.maxPriority, GAS.maxFee),
    paymasterAndData: "0x",
    signature: "0x",
  };
}

function paymasterAndData(paymaster, { token, validUntil = 0, validAfter = 0, rate = RATE, signature }) {
  return concat([
    paymaster.target,
    zeroPadValue(toBeHex(GAS.pmVerification), 16),
    zeroPadValue(toBeHex(GAS.postOp), 16),
    token,
    zeroPadValue(toBeHex(validUntil), 6),
    zeroPadValue(toBeHex(validAfter), 6),
    zeroPadValue(toBeHex(rate), 32),
    signature ?? "0x" + "00".repeat(65),
  ]);
}

// quote: what the quote service does - hash the op and eth_sign it
async function quote(f, op, { token = f.usdc.target, validUntil = 0, validAfter = 0, rate = RATE, signer = f.quoteSigner } = {}) {
  op.paymasterAndData = paymasterAndData(f.paymaster, { token, validUntil, validAfter, rate });
  const hash = await f.paymaster.getHash(op, token, validUntil, validAfter, rate);
  const signature = await signer.signMessage(getBytes(hash));
  op.paymasterAndData = paymasterAndData(f.paymaster, { token, validUntil, validAfter, rate, signature });
  return op;
}

// the wallet owner signs the userOpHash, which covers the quote
async function signOp(f, op, signer = f.walletOwner) {
  op.signature = await signer.signMessage(getBytes(await f.entryPoint.getUserOpHash(op)));
  return op;
}

async function send(f, op) {
  return f.entryPoint.handleOps([op], f.beneficiary.address);
}

function ceilDiv(a, b) {
  return (a + b - 1n) / b;
}

async function userOpEvent(f, tx) {
  const receipt = await tx.wait();
  for (const log of receipt.logs) {
    try {
      const parsed = f.entryPoint.interface.parseLog(log);
      if (parsed && parsed.name === "UserOperationEvent") return parsed.args;
    } catch (_) {}
  }
  throw new Error("no UserOperationEvent");
}

function paymasterError(f, name) {
  return (inner) => {
    const parsed = f.paymaster.interface.parseError(inner);
    return parsed !== null && parsed.name === name;
  };
}

describe("TrovoTokenPaymaster", function () {
  it("charges the actual gas in the token and refunds the pre-charge", async function () {
    const f = await loadFixture(deployFixture);
    const op = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, USDC(10)) })));

    const depositBefore = await f.paymaster.getDeposit();
    const tx = await send(f, op);
    const ev = await userOpEvent(f, tx);
    expect(ev.success).to.equal(true);
    expect(ev.paymaster).to.equal(f.paymaster.target);

    const feePerGas = ev.actualGasCost / ev.actualGasUsed;
    // postOp's cost is estimated with the overhead; EntryPoint's actualGasCost
    // passed to postOp excludes postOp itself, so recompute from the event.
    const charged = await f.usdc.balanceOf(f.paymaster);
    expect(charged).to.be.greaterThan(0n);
    expect(await f.usdc.balanceOf(f.recipient)).to.equal(USDC(10));
    expect(await f.usdc.balanceOf(f.account)).to.equal(USDC(100) - USDC(10) - charged);

    // charge is at most the actual gas (incl. overhead) at the quoted rate
    const upperBound = ceilDiv((ev.actualGasCost + POST_OP_OVERHEAD * feePerGas) * RATE, 10n ** 18n);
    expect(charged).to.be.lessThanOrEqual(upperBound);
    await expect(tx).to.emit(f.paymaster, "GasPaidInToken").withArgs(f.account.target, f.usdc.target, ev.userOpHash, RATE, charged, (v) => v > 0n);

    // the paymaster's ETH deposit paid for it
    expect(await f.paymaster.getDeposit()).to.equal(depositBefore - ev.actualGasCost);
  });

  it("deploys a wallet that approves in its first op and charges in postOp", async function () {
    const f = await loadFixture(deployFixture);
    const sender = await f.factory["getAddress(address,uint256)"](f.walletOwner.address, 7);
    await f.usdc.mint(sender, USDC(50)); // funded before activation (card, P2P, another user)

    const initCode = concat([f.factory.target, f.factory.interface.encodeFunctionData("createAccount", [f.walletOwner.address, 7])]);
    const account = await ethers.getContractAt("SimpleAccount", sender);
    const callData = account.interface.encodeFunctionData("executeBatch", [
      [f.usdc.target, f.usdc.target],
      [],
      [
        f.usdc.interface.encodeFunctionData("approve", [f.paymaster.target, MaxUint256]),
        f.usdc.interface.encodeFunctionData("transfer", [f.recipient.address, USDC(5)]),
      ],
    ]);
    const op = await signOp(f, await quote(f, await buildOp(f, { sender, initCode, callData })));
    const tx = await send(f, op);
    const ev = await userOpEvent(f, tx);
    expect(ev.success).to.equal(true);
    expect(await ethers.provider.getCode(sender)).to.not.equal("0x");

    const charged = await f.usdc.balanceOf(f.paymaster);
    expect(charged).to.be.greaterThan(0n);
    expect(await f.usdc.balanceOf(sender)).to.equal(USDC(50) - USDC(5) - charged);
    await expect(tx).to.emit(f.paymaster, "GasPaidInToken");
    expect(await f.paymaster.debtTokenCount(sender)).to.equal(0n);
  });

  it("records debt when the postOp charge fails and blocks the wallet until settled", async function () {
    const f = await loadFixture(deployFixture);
    const sender = await f.factory["getAddress(address,uint256)"](f.walletOwner.address, 9);
    await f.usdc.mint(sender, USDC(50));
    const initCode = concat([f.factory.target, f.factory.interface.encodeFunctionData("createAccount", [f.walletOwner.address, 9])]);
    const account = await ethers.getContractAt("SimpleAccount", sender);

    // no approve in the first op's calldata
    const op = await signOp(f, await quote(f, await buildOp(f, { sender, initCode, callData: transferCall(account, f.usdc, f.recipient.address, USDC(1)) })));
    const tx = await send(f, op);
    await expect(tx).to.emit(f.paymaster, "GasChargeFailed");
    const owed = await f.paymaster.debt(sender, f.usdc);
    expect(owed).to.be.greaterThan(0n);
    expect(await f.paymaster.debtTokenCount(sender)).to.equal(1n);
    // the wallet's action itself still ran
    expect(await f.usdc.balanceOf(f.recipient)).to.equal(USDC(1));

    const next = await signOp(f, await quote(f, await buildOp(f, { sender, callData: transferCall(account, f.usdc, f.recipient.address, USDC(1)) })));
    await expect(send(f, next))
      .to.be.revertedWithCustomError(f.entryPoint, "FailedOpWithRevert")
      .withArgs(0, "AA33 reverted", paymasterError(f, "OutstandingDebt"));

    // settle: approve (e.g. in an ETH-paid op) then anyone settles
    await account.connect(f.walletOwner).execute(f.usdc, 0, f.usdc.interface.encodeFunctionData("approve", [f.paymaster.target, MaxUint256]));
    await expect(f.paymaster.connect(f.stranger).settleDebt(sender, f.usdc))
      .to.emit(f.paymaster, "DebtSettled")
      .withArgs(sender, f.usdc.target, owed, false);
    expect(await f.paymaster.debtTokenCount(sender)).to.equal(0n);
    await expect(f.paymaster.settleDebt(sender, f.usdc)).to.be.revertedWithCustomError(f.paymaster, "NoDebt");

    const retry = await signOp(f, await quote(f, await buildOp(f, { sender, callData: transferCall(account, f.usdc, f.recipient.address, USDC(1)) })));
    await expect(send(f, retry)).to.emit(f.paymaster, "GasPaidInToken");
  });

  it("lets the owner forgive debt", async function () {
    const f = await loadFixture(deployFixture);
    const sender = await f.factory["getAddress(address,uint256)"](f.walletOwner.address, 11);
    await f.usdc.mint(sender, USDC(50));
    const initCode = concat([f.factory.target, f.factory.interface.encodeFunctionData("createAccount", [f.walletOwner.address, 11])]);
    const account = await ethers.getContractAt("SimpleAccount", sender);
    await send(f, await signOp(f, await quote(f, await buildOp(f, { sender, initCode, callData: transferCall(account, f.usdc, f.recipient.address, 1n) }))));
    const owed = await f.paymaster.debt(sender, f.usdc);

    await expect(f.paymaster.connect(f.stranger).forgiveDebt(sender, f.usdc)).to.be.revertedWithCustomError(f.paymaster, "OwnableUnauthorizedAccount");
    await expect(f.paymaster.connect(f.owner).forgiveDebt(sender, f.usdc)).to.emit(f.paymaster, "DebtSettled").withArgs(sender, f.usdc.target, owed, true);
    expect(await f.paymaster.debt(sender, f.usdc)).to.equal(0n);
  });

  it("still charges when the wallet's call reverts", async function () {
    const f = await loadFixture(deployFixture);
    // transferring more than the balance makes the call revert
    const op = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, USDC(1000)) })));
    const tx = await send(f, op);
    const ev = await userOpEvent(f, tx);
    expect(ev.success).to.equal(false);
    const charged = await f.usdc.balanceOf(f.paymaster);
    expect(charged).to.be.greaterThan(0n);
    expect(await f.usdc.balanceOf(f.account)).to.equal(USDC(100) - charged);
    await expect(tx).to.emit(f.paymaster, "GasPaidInToken");
  });

  it("rejects a quote not signed by a quote signer", async function () {
    const f = await loadFixture(deployFixture);
    const op = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) }), { signer: f.stranger }));
    await expect(send(f, op)).to.be.revertedWithCustomError(f.entryPoint, "FailedOp").withArgs(0, "AA34 signature error");
  });

  it("rejects an op whose calldata changed after the quote", async function () {
    const f = await loadFixture(deployFixture);
    const op = await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) }));
    op.callData = transferCall(f.account, f.usdc, f.stranger.address, USDC(90));
    await expect(send(f, await signOp(f, op))).to.be.revertedWithCustomError(f.entryPoint, "FailedOp").withArgs(0, "AA34 signature error");
  });

  it("rejects a rate changed after the quote", async function () {
    const f = await loadFixture(deployFixture);
    const op = await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) }));
    const sig = "0x" + op.paymasterAndData.slice(2 + 116 * 2);
    op.paymasterAndData = paymasterAndData(f.paymaster, { token: f.usdc.target, rate: MIN_RATE, signature: sig });
    await expect(send(f, await signOp(f, op))).to.be.revertedWithCustomError(f.entryPoint, "FailedOp").withArgs(0, "AA34 signature error");
  });

  it("rejects an expired or not-yet-valid quote", async function () {
    const f = await loadFixture(deployFixture);
    const now = await time.latest();
    const expired = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) }), { validUntil: now - 10 }));
    await expect(send(f, expired)).to.be.revertedWithCustomError(f.entryPoint, "FailedOp").withArgs(0, "AA32 paymaster expired or not due");

    const early = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) }), { validAfter: now + 3600 }));
    await expect(send(f, early)).to.be.revertedWithCustomError(f.entryPoint, "FailedOp").withArgs(0, "AA32 paymaster expired or not due");
  });

  it("accepts a long-lived quote within its window", async function () {
    const f = await loadFixture(deployFixture);
    const now = await time.latest();
    const op = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) }), { validAfter: now - 60, validUntil: now + 24 * 3600 }));
    await time.increase(23 * 3600); // e.g. a shared wallet collecting approvals
    await expect(send(f, op)).to.emit(f.paymaster, "GasPaidInToken");
  });

  it("rejects tokens that are not enabled", async function () {
    const f = await loadFixture(deployFixture);
    await f.other.mint(f.account, USDC(100));
    const op = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) }), { token: f.other.target }));
    await expect(send(f, op))
      .to.be.revertedWithCustomError(f.entryPoint, "FailedOpWithRevert")
      .withArgs(0, "AA33 reverted", paymasterError(f, "TokenNotEnabled"));

    await f.paymaster.connect(f.owner).setToken(f.usdc, false, 0, 0);
    const op2 = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) })));
    await expect(send(f, op2))
      .to.be.revertedWithCustomError(f.entryPoint, "FailedOpWithRevert")
      .withArgs(0, "AA33 reverted", paymasterError(f, "TokenNotEnabled"));
  });

  it("rejects rates outside the owner's bounds, even when signed", async function () {
    const f = await loadFixture(deployFixture);
    for (const rate of [MIN_RATE - 1n, MAX_RATE + 1n]) {
      const op = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) }), { rate }));
      await expect(send(f, op))
        .to.be.revertedWithCustomError(f.entryPoint, "FailedOpWithRevert")
        .withArgs(0, "AA33 reverted", paymasterError(f, "RateOutOfBounds"));
    }
  });

  it("rejects a wallet that cannot cover the maximum cost", async function () {
    const f = await loadFixture(deployFixture);
    await f.account.connect(f.walletOwner).execute(f.usdc, 0, f.usdc.interface.encodeFunctionData("transfer", [f.recipient.address, USDC(99.99)]));
    const op = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) })));
    await expect(send(f, op))
      .to.be.revertedWithCustomError(f.entryPoint, "FailedOpWithRevert")
      .withArgs(0, "AA33 reverted", paymasterError(f, "InsufficientTokenBalance"));
  });

  it("rejects malformed paymaster data", async function () {
    const f = await loadFixture(deployFixture);
    const op = await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) });
    op.paymasterAndData = concat([paymasterAndData(f.paymaster, { token: f.usdc.target }), "0x00"]);
    await expect(send(f, await signOp(f, op)))
      .to.be.revertedWithCustomError(f.entryPoint, "FailedOpWithRevert")
      .withArgs(0, "AA33 reverted", paymasterError(f, "InvalidPaymasterData"));
  });

  it("can be paused by the pauser and unpaused only by the owner", async function () {
    const f = await loadFixture(deployFixture);
    await expect(f.paymaster.connect(f.stranger).pause()).to.be.revertedWithCustomError(f.paymaster, "NotPauser");
    await f.paymaster.connect(f.pauser).pause();

    const op = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) })));
    await expect(send(f, op))
      .to.be.revertedWithCustomError(f.entryPoint, "FailedOpWithRevert")
      .withArgs(0, "AA33 reverted", paymasterError(f, "EnforcedPause"));

    await expect(f.paymaster.connect(f.pauser).unpause()).to.be.revertedWithCustomError(f.paymaster, "OwnableUnauthorizedAccount");
    await f.paymaster.connect(f.owner).unpause();
    await expect(send(f, op)).to.emit(f.paymaster, "GasPaidInToken");
  });

  it("rotates quote signers", async function () {
    const f = await loadFixture(deployFixture);
    await f.paymaster.connect(f.owner).setQuoteSigner(f.stranger.address, true);
    await f.paymaster.connect(f.owner).setQuoteSigner(f.quoteSigner.address, false);
    const oldSigner = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) })));
    await expect(send(f, oldSigner)).to.be.revertedWithCustomError(f.entryPoint, "FailedOp").withArgs(0, "AA34 signature error");
    const newSigner = await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) }), { signer: f.stranger }));
    await expect(send(f, newSigner)).to.emit(f.paymaster, "GasPaidInToken");
  });

  it("restricts administration to the owner", async function () {
    const f = await loadFixture(deployFixture);
    const p = f.paymaster.connect(f.stranger);
    for (const call of [
      () => p.setToken(f.usdc, true, 1, 2),
      () => p.setQuoteSigner(f.stranger.address, true),
      () => p.setPauser(f.stranger.address),
      () => p.setPostOpGasOverhead(1),
      () => p.withdrawToken(f.usdc, f.stranger.address, 1),
      () => p.withdrawTo(f.stranger.address, 1),
      () => p.unlockStake(),
    ]) {
      await expect(call()).to.be.revertedWithCustomError(f.paymaster, "OwnableUnauthorizedAccount");
    }
    await expect(f.paymaster.connect(f.owner).setToken(f.usdc, true, 0, 1)).to.be.revertedWithCustomError(f.paymaster, "InvalidRateBounds");
    await expect(f.paymaster.connect(f.owner).setToken(f.usdc, true, 2, 1)).to.be.revertedWithCustomError(f.paymaster, "InvalidRateBounds");
  });

  it("lets the owner withdraw collected tokens to the treasury", async function () {
    const f = await loadFixture(deployFixture);
    await send(f, await signOp(f, await quote(f, await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) }))));
    const collected = await f.usdc.balanceOf(f.paymaster);
    await expect(f.paymaster.connect(f.owner).withdrawToken(f.usdc, f.owner.address, collected))
      .to.emit(f.paymaster, "TokenWithdrawn")
      .withArgs(f.usdc.target, f.owner.address, collected);
    expect(await f.usdc.balanceOf(f.owner)).to.equal(collected);
  });

  it("rounds token cost up", async function () {
    const f = await loadFixture(deployFixture);
    expect(await f.paymaster.tokenCost(1n, RATE)).to.equal(1n);
    expect(await f.paymaster.tokenCost(10n ** 18n, RATE)).to.equal(RATE);
    expect(await f.paymaster.tokenCost(0n, RATE)).to.equal(0n);
  });

  it("matches the quote hash encoding the quote service uses", async function () {
    const f = await loadFixture(deployFixture);
    const op = await buildOp(f, { sender: f.account.target, callData: transferCall(f.account, f.usdc, f.recipient.address, 1n) });
    op.paymasterAndData = paymasterAndData(f.paymaster, { token: f.usdc.target, validUntil: 100, validAfter: 5 });
    const coder = AbiCoder.defaultAbiCoder();
    const opHash = keccak256(
      coder.encode(
        ["address", "uint256", "bytes32", "bytes32", "bytes32", "uint256", "uint256", "bytes32"],
        [op.sender, op.nonce, keccak256(op.initCode), keccak256(op.callData), op.accountGasLimits, (GAS.pmVerification << 128n) | GAS.postOp, op.preVerificationGas, op.gasFees],
      ),
    );
    const { chainId } = await ethers.provider.getNetwork();
    const expected = keccak256(
      coder.encode(
        ["bytes32", "uint256", "address", "address", "uint48", "uint48", "uint256"],
        [opHash, chainId, f.paymaster.target, f.usdc.target, 100, 5, RATE],
      ),
    );
    expect(await f.paymaster.getHash(op, f.usdc, 100, 5, RATE)).to.equal(expected);
  });
});
