// End-to-end with real Safe v1.4.1 wallets and the Safe4337Module: the
// wallet layout every Trovo wallet uses. Proves that
//   - the counterfactual address formula (shared by wallet-core and
//     app-backend) matches what the factory deploys,
//   - a personal_sign (eth_sign, v+4) owner signature over the SafeOp hash
//     is accepted by the Safe,
//   - an undeployed Safe that received funds activates and acts in one
//     UserOperation, paying its own gas through TrovoTokenPaymaster.
// It also writes the fixture wallet-core and app-backend test against.

const fs = require("fs");
const path = require("path");
const { expect } = require("chai");
const { ethers } = require("hardhat");
const { loadFixture } = require("@nomicfoundation/hardhat-toolbox/network-helpers");

const { concat, getBytes, keccak256, solidityPacked, toBeHex, zeroPadValue, MaxUint256 } = ethers;

const safeArtifact = (name) => require(`./safe-artifacts/${name}.json`);
const pack128 = (hi, lo) => zeroPadValue(toBeHex((BigInt(hi) << 128n) | BigInt(lo)), 32);
const USDC = (n) => BigInt(Math.round(n * 1e6));
const RATE = USDC(3000);

async function deploy(signer, name, ...args) {
  const a = safeArtifact(name);
  const c = await new ethers.ContractFactory(a.abi, a.bytecode, signer).deploy(...args);
  await c.waitForDeployment();
  return c;
}

async function fixture() {
  const [deployer, quoteSigner, owner, recipient, beneficiary, secondOwner] = await ethers.getSigners();
  const entryPoint = await ethers.deployContract("EntryPoint");
  const singleton = await deploy(deployer, "SafeL2");
  const factory = await deploy(deployer, "SafeProxyFactory");
  const moduleSetup = await deploy(deployer, "SafeModuleSetup");
  const module = await deploy(deployer, "Safe4337Module", entryPoint.target);
  const multiSend = await deploy(deployer, "MultiSendCallOnly");
  const usdc = await ethers.deployContract("MockERC20", ["USD Coin", "USDC", 6]);
  const paymaster = await ethers.deployContract("TrovoTokenPaymaster", [entryPoint, deployer, [quoteSigner.address], deployer.address, 45000]);
  await paymaster.addStake(86400, { value: ethers.parseEther("1") });
  await paymaster.deposit({ value: ethers.parseEther("5") });
  await paymaster.setToken(usdc, true, USDC(1000), USDC(10000));
  return { deployer, quoteSigner, owner, recipient, beneficiary, secondOwner, entryPoint, singleton, factory, moduleSetup, module, multiSend, usdc, paymaster };
}

function initializer(f, owners, threshold) {
  const enable = f.moduleSetup.interface.encodeFunctionData("enableModules", [[f.module.target]]);
  return f.singleton.interface.encodeFunctionData("setup", [owners, threshold, f.moduleSetup.target, enable, f.module.target, ethers.ZeroAddress, 0, ethers.ZeroAddress]);
}

// the formula wallet-core and app-backend implement
function predictAddress(f, init, saltNonce, creationCode) {
  const salt = keccak256(solidityPacked(["bytes32", "uint256"], [keccak256(init), saltNonce]));
  const initCodeHash = keccak256(concat([creationCode, zeroPadValue(f.singleton.target, 32)]));
  return ethers.getCreate2Address(f.factory.target, salt, initCodeHash);
}

function multiSendData(f, calls) {
  const packed = concat(calls.map((c) => solidityPacked(["uint8", "address", "uint256", "uint256", "bytes"], [0, c.to, 0, getBytes(c.data).length, c.data])));
  return f.multiSend.interface.encodeFunctionData("multiSend", [packed]);
}

function executeUserOp(f, to, value, data, operation) {
  return f.module.interface.encodeFunctionData("executeUserOp", [to, value, data, operation]);
}

async function quoted(f, op, token) {
  const pmStatic = concat([f.paymaster.target, zeroPadValue(toBeHex(200000), 16), zeroPadValue(toBeHex(80000), 16)]);
  op.paymasterAndData = concat([pmStatic, token, zeroPadValue("0x00", 6), zeroPadValue("0x00", 6), zeroPadValue(toBeHex(RATE), 32), "0x" + "00".repeat(65)]);
  const hash = await f.paymaster.getHash(op, token, 0, 0, RATE);
  const sig = await f.quoteSigner.signMessage(getBytes(hash));
  op.paymasterAndData = concat([pmStatic, token, zeroPadValue("0x00", 6), zeroPadValue("0x00", 6), zeroPadValue(toBeHex(RATE), 32), sig]);
  return op;
}

// what the apps do: personal_sign the SafeOp hash; the backend adds 4 to v
// (Safe's eth_sign signature type) and prefixes validAfter/validUntil
async function ownerSignature(f, op, signers, validAfter = 0, validUntil = 0) {
  const withTimes = { ...op, signature: solidityPacked(["uint48", "uint48"], [validAfter, validUntil]) };
  const safeOpHash = await f.module.getOperationHash(withTimes);
  const sigs = [];
  for (const s of [...signers].sort((a, b) => (BigInt(a.address) < BigInt(b.address) ? -1 : 1))) {
    const sig = getBytes(await s.signMessage(getBytes(safeOpHash)));
    sig[64] += 4;
    sigs.push(sig);
  }
  return { safeOpHash, signature: concat([solidityPacked(["uint48", "uint48"], [validAfter, validUntil]), ...sigs]) };
}

function baseOp(sender, nonce, initCode, callData) {
  return {
    sender,
    nonce,
    initCode,
    callData,
    accountGasLimits: pack128(600000, 300000),
    preVerificationGas: 80000,
    gasFees: pack128(ethers.parseUnits("1", "gwei"), ethers.parseUnits("3", "gwei")),
    paymasterAndData: "0x",
    signature: "0x",
  };
}

describe("Safe wallets (Safe v1.4.1 + Safe4337Module)", function () {
  it("derives the counterfactual address the factory deploys", async function () {
    const f = await loadFixture(fixture);
    const creationCode = await f.factory.proxyCreationCode();
    for (const [owners, threshold, salt] of [[[f.owner.address], 1, 0n], [[f.owner.address], 1, 7n], [[f.owner.address, f.secondOwner.address], 2, 12345n]]) {
      const init = initializer(f, owners, threshold);
      const predicted = predictAddress(f, init, salt, creationCode);
      expect(await f.factory.createProxyWithNonce.staticCall(f.singleton, init, salt)).to.equal(predicted);
    }
  });

  it("activates a funded, undeployed Safe and pays gas in USDC in one UserOperation", async function () {
    const f = await loadFixture(fixture);
    const creationCode = await f.factory.proxyCreationCode();
    const init = initializer(f, [f.owner.address], 1);
    const safe = predictAddress(f, init, 0n, creationCode);

    // funds arrive before activation (card funding, another user, P2P)
    await f.usdc.mint(safe, USDC(50));
    expect(await ethers.provider.getCode(safe)).to.equal("0x");

    const initCode = concat([f.factory.target, f.factory.interface.encodeFunctionData("createProxyWithNonce", [f.singleton.target, init, 0])]);
    const callData = executeUserOp(f, f.multiSend.target, 0, multiSendData(f, [
      { to: f.usdc.target, data: f.usdc.interface.encodeFunctionData("approve", [f.paymaster.target, MaxUint256]) },
      { to: f.usdc.target, data: f.usdc.interface.encodeFunctionData("transfer", [f.recipient.address, USDC(10)]) },
    ]), 1);
    const op = await quoted(f, baseOp(safe, await f.entryPoint.getNonce(safe, 0), initCode, callData), f.usdc.target);
    op.signature = (await ownerSignature(f, op, [f.owner])).signature;

    const tx = await f.entryPoint.handleOps([op], f.beneficiary.address);
    await expect(tx).to.emit(f.paymaster, "GasPaidInToken");
    const charged = await f.usdc.balanceOf(f.paymaster);
    expect(await f.usdc.balanceOf(f.recipient)).to.equal(USDC(10));
    expect(await f.usdc.balanceOf(safe)).to.equal(USDC(40) - charged);

    const deployed = new ethers.Contract(safe, safeArtifact("SafeL2").abi, ethers.provider);
    expect(await deployed.getOwners()).to.deep.equal([f.owner.address]);
    expect(await deployed.isModuleEnabled(f.module.target)).to.equal(true);

    // a second, pre-charged operation from the now-deployed Safe
    const op2 = await quoted(f, baseOp(safe, await f.entryPoint.getNonce(safe, 0), "0x", executeUserOp(f, f.usdc.target, 0, f.usdc.interface.encodeFunctionData("transfer", [f.recipient.address, USDC(1)]), 0)), f.usdc.target);
    op2.signature = (await ownerSignature(f, op2, [f.owner])).signature;
    await expect(f.entryPoint.handleOps([op2], f.beneficiary.address)).to.emit(f.paymaster, "GasPaidInToken");
    expect(await f.usdc.balanceOf(f.recipient)).to.equal(USDC(11));
  });

  it("requires every owner's signature on a 2-of-2 Safe and rejects a wrong signer", async function () {
    const f = await loadFixture(fixture);
    const creationCode = await f.factory.proxyCreationCode();
    const init = initializer(f, [f.owner.address, f.secondOwner.address], 2);
    const safe = predictAddress(f, init, 99n, creationCode);
    await f.usdc.mint(safe, USDC(50));
    const initCode = concat([f.factory.target, f.factory.interface.encodeFunctionData("createProxyWithNonce", [f.singleton.target, init, 99])]);
    const callData = executeUserOp(f, f.multiSend.target, 0, multiSendData(f, [
      { to: f.usdc.target, data: f.usdc.interface.encodeFunctionData("approve", [f.paymaster.target, MaxUint256]) },
    ]), 1);

    const op = await quoted(f, baseOp(safe, 0, initCode, callData), f.usdc.target);
    op.signature = (await ownerSignature(f, op, [f.owner, f.recipient])).signature;
    await expect(f.entryPoint.handleOps([op], f.beneficiary.address)).to.be.revertedWithCustomError(f.entryPoint, "FailedOp").withArgs(0, "AA24 signature error");

    op.signature = (await ownerSignature(f, op, [f.owner, f.secondOwner])).signature;
    await expect(f.entryPoint.handleOps([op], f.beneficiary.address)).to.emit(f.paymaster, "GasPaidInToken");
  });

  it("writes the cross-language fixture", async function () {
    const f = await loadFixture(fixture);
    const creationCode = await f.factory.proxyCreationCode();
    const cases = [];
    for (const [owners, threshold, salt] of [[[f.owner.address], 1, "0"], [[f.owner.address], 1, "115792089237316195423570985008687907853269984665640564039457584007913129639935"], [[f.owner.address, f.secondOwner.address, f.recipient.address], 2, "4242"]]) {
      const init = initializer(f, owners, threshold);
      const address = await f.factory.createProxyWithNonce.staticCall(f.singleton, init, salt);
      cases.push({ owners, threshold, saltNonce: salt, initializer: init, address });
    }
    // a SafeOp hash case, from the module itself
    const op = baseOp(cases[0].address, 5, "0x", executeUserOp(f, f.usdc.target, 0, f.usdc.interface.encodeFunctionData("transfer", [f.recipient.address, 1]), 0));
    op.paymasterAndData = "0x1234";
    const { safeOpHash, signature } = await ownerSignature(f, op, [f.owner], 100, 200);
    const userOpHash = await f.entryPoint.getUserOpHash({ ...op, signature });
    const { chainId } = await ethers.provider.getNetwork();
    const out = {
      comment: "Generated by paymaster/contracts/test/SafeWallet.test.js from real Safe v1.4.1 + Safe4337Module deployments on a local chain.",
      config: { proxyFactory: f.factory.target, singleton: f.singleton.target, moduleSetup: f.moduleSetup.target, safe4337Module: f.module.target, entryPoint: f.entryPoint.target, proxyCreationCode: creationCode },
      addresses: cases,
      safeOp: {
        chainId: chainId.toString(),
        sender: op.sender, nonce: "5", initCode: op.initCode, callData: op.callData,
        verificationGasLimit: "600000", callGasLimit: "300000", preVerificationGas: "80000",
        maxPriorityFeePerGas: ethers.parseUnits("1", "gwei").toString(), maxFeePerGas: ethers.parseUnits("3", "gwei").toString(),
        paymasterAndData: op.paymasterAndData, validAfter: 100, validUntil: 200, hash: safeOpHash,
        userOpHash, ownerSignatureWithV4: signature,
      },
    };
    for (const file of [
      path.join(__dirname, "..", "..", "..", "wallet-core", "testdata", "safe_fixture.json"),
      path.join(__dirname, "..", "..", "..", "app-backend", "internal", "aa", "testdata", "safe_fixture.json"),
    ]) {
      fs.mkdirSync(path.dirname(file), { recursive: true });
      fs.writeFileSync(file, JSON.stringify(out, null, 2) + "\n");
    }
  });
});
