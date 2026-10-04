// Writes the quote-hash fixture the Go quote service tests against, so the
// Go implementation of getHash (and its eth_sign signature) is checked
// against the contract itself:
//
//   npx hardhat run scripts/fixtures.js
//
// Output: ../quote-service/internal/quote/testdata/hash_fixture.json

const fs = require("fs");
const path = require("path");
const { ethers } = require("hardhat");

const pack128 = (hi, lo) => ethers.zeroPadValue(ethers.toBeHex((BigInt(hi) << 128n) | BigInt(lo)), 32);

async function main() {
  const [deployer, quoteSigner] = await ethers.getSigners();
  const entryPoint = await ethers.deployContract("EntryPoint");
  const paymaster = await ethers.deployContract("TrovoTokenPaymaster", [entryPoint, deployer, [quoteSigner.address], deployer.address, 45000]);
  const { chainId } = await ethers.provider.getNetwork();

  const cases = [
    {
      name: "deployed wallet",
      sender: "0x5a6b0c1c6f1c0f4b1d7aa2d1c3e0bee8a9f2d4c1",
      nonce: "7",
      initCode: "0x",
      callData: "0x7bb374280000000000000000000000000000000000000000000000000000000000000001",
      verificationGasLimit: "300000",
      callGasLimit: "200000",
      preVerificationGas: "60000",
      maxFeePerGas: "5000000000",
      maxPriorityFeePerGas: "1000000",
      paymasterVerificationGasLimit: "150000",
      paymasterPostOpGasLimit: "80000",
      token: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
      validUntil: 1790000000,
      validAfter: 1789999000,
      exchangeRate: "3000000000",
    },
    {
      name: "activation with initCode, large values",
      sender: "0x00000000000000000000000000000000000000ff",
      nonce: "340282366920938463463374607431768211456", // key 1 << 64 style nonce
      initCode: "0x4e1DCf7AD4e460CfD30791CCC4F9c8a4f820ec671688f0b9000000000000000000000000000000000000000000000000000000000000abcd",
      callData: "0x",
      verificationGasLimit: "1000000",
      callGasLimit: "0",
      preVerificationGas: "123456789",
      maxFeePerGas: "340282366920938463463374607431768211455",
      maxPriorityFeePerGas: "1",
      paymasterVerificationGasLimit: "1",
      paymasterPostOpGasLimit: "340282366920938463463374607431768211455",
      token: "0x1000000000000000000000000000000000000001",
      validUntil: 281474976710655,
      validAfter: 0,
      exchangeRate: "115792089237316195423570985008687907853269984665640564039457584007913129639935",
    },
  ];

  const out = { chainId: chainId.toString(), paymaster: await paymaster.getAddress(), quoteSignerPrivateKey: "", cases: [] };
  // hardhat's well-known second test account; never used outside tests
  out.quoteSignerPrivateKey = "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d";
  if (new ethers.Wallet(out.quoteSignerPrivateKey).address !== quoteSigner.address) throw new Error("unexpected signer");

  for (const c of cases) {
    const paymasterAndData = ethers.concat([
      await paymaster.getAddress(),
      ethers.zeroPadValue(ethers.toBeHex(c.paymasterVerificationGasLimit), 16),
      ethers.zeroPadValue(ethers.toBeHex(c.paymasterPostOpGasLimit), 16),
    ]);
    const op = {
      sender: c.sender,
      nonce: c.nonce,
      initCode: c.initCode,
      callData: c.callData,
      accountGasLimits: pack128(c.verificationGasLimit, c.callGasLimit),
      preVerificationGas: c.preVerificationGas,
      gasFees: pack128(c.maxPriorityFeePerGas, c.maxFeePerGas),
      paymasterAndData,
      signature: "0x",
    };
    const hash = await paymaster.getHash(op, c.token, c.validUntil, c.validAfter, c.exchangeRate);
    const signature = await quoteSigner.signMessage(ethers.getBytes(hash));
    out.cases.push({ ...c, hash, signature });
  }

  const file = path.join(__dirname, "..", "..", "quote-service", "internal", "quote", "testdata", "hash_fixture.json");
  fs.mkdirSync(path.dirname(file), { recursive: true });
  fs.writeFileSync(file, JSON.stringify(out, null, 2) + "\n");
  console.log(`wrote ${file}`);
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
