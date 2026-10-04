// Local development stack for Trovo wallets: deploys the ERC-4337
// EntryPoint v0.7, Safe v1.4.1 + Safe4337Module, a mock USDC and the
// TrovoTokenPaymaster to a local node, then serves a minimal ERC-4337
// bundler so app-backend and the quote service can run end to end without
// a real bundler.
//
//   npx hardhat node                                  # terminal 1
//   npx hardhat run scripts/local-stack.js --network localhost   # terminal 2
//
// It prints the settings to use (ENTRYPOINT_ADDRESS, SAFE_* overrides,
// BUNDLER_URL, PAYMASTER_ADDRESS, ...) and writes them to
// deployments/local-stack.json. The bundler is for development only: it
// submits each operation on its own with handleOps and returns fixed,
// generous gas estimates.

const fs = require("fs");
const http = require("http");
const path = require("path");
const { ethers } = require("hardhat");

const safeArtifact = (name) => require(path.join(__dirname, "..", "test", "safe-artifacts", `${name}.json`));
const BUNDLER_PORT = Number(process.env.DEV_BUNDLER_PORT || 4337);
const QUOTE_SIGNER = process.env.QUOTE_SIGNER_ADDRESS || "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"; // hardhat #1

async function deploySafe(signer, name, ...args) {
  const a = safeArtifact(name);
  const c = await new ethers.ContractFactory(a.abi, a.bytecode, signer).deploy(...args);
  await c.waitForDeployment();
  return c;
}

function unpackRpc(op) {
  const initCode = op.factory ? ethers.concat([op.factory, op.factoryData || "0x"]) : "0x";
  const pack = (hi, lo) => ethers.zeroPadValue(ethers.toBeHex((BigInt(hi) << 128n) | BigInt(lo)), 32);
  const paymasterAndData = op.paymaster
    ? ethers.concat([op.paymaster, ethers.zeroPadValue(ethers.toBeHex(op.paymasterVerificationGasLimit), 16), ethers.zeroPadValue(ethers.toBeHex(op.paymasterPostOpGasLimit), 16), op.paymasterData || "0x"])
    : "0x";
  return {
    sender: op.sender,
    nonce: op.nonce,
    initCode,
    callData: op.callData,
    accountGasLimits: pack(op.verificationGasLimit, op.callGasLimit),
    preVerificationGas: op.preVerificationGas,
    gasFees: pack(op.maxPriorityFeePerGas, op.maxFeePerGas),
    paymasterAndData,
    signature: op.signature,
  };
}

async function main() {
  const [deployer, , , , , , , , , bundlerKey] = await ethers.getSigners();
  const entryPoint = await ethers.deployContract("EntryPoint");
  const singleton = await deploySafe(deployer, "SafeL2");
  const factory = await deploySafe(deployer, "SafeProxyFactory");
  const moduleSetup = await deploySafe(deployer, "SafeModuleSetup");
  const module = await deploySafe(deployer, "Safe4337Module", entryPoint.target);
  const multiSend = await deploySafe(deployer, "MultiSendCallOnly");
  const usdc = await ethers.deployContract("MockERC20", ["USD Coin", "USDC", 6]);
  const paymaster = await ethers.deployContract("TrovoTokenPaymaster", [entryPoint, deployer, [QUOTE_SIGNER], deployer.address, 45000]);
  await (await paymaster.addStake(86400, { value: ethers.parseEther("1") })).wait();
  await (await paymaster.deposit({ value: ethers.parseEther("10") })).wait();
  await (await paymaster.setToken(usdc, true, 1_000_000_000n, 20_000_000_000n)).wait();
  const { chainId } = await ethers.provider.getNetwork();

  const settings = {
    BASE_CHAIN_ID: chainId.toString(),
    ENTRYPOINT_ADDRESS: entryPoint.target,
    SAFE_PROXY_FACTORY_ADDRESS: factory.target,
    SAFE_SINGLETON_ADDRESS: singleton.target,
    SAFE_MODULE_SETUP_ADDRESS: moduleSetup.target,
    SAFE_4337_MODULE_ADDRESS: module.target,
    SAFE_MULTISEND_CALL_ONLY_ADDRESS: multiSend.target,
    SAFE_PROXY_CREATION_CODE: await factory.proxyCreationCode(),
    PAYMASTER_ADDRESS: paymaster.target,
    MOCK_USDC_ADDRESS: usdc.target,
    BUNDLER_URL: `http://127.0.0.1:${BUNDLER_PORT}`,
  };
  const dir = path.join(__dirname, "..", "deployments");
  fs.mkdirSync(dir, { recursive: true });
  fs.writeFileSync(path.join(dir, "local-stack.json"), JSON.stringify(settings, null, 2) + "\n");
  for (const [k, v] of Object.entries(settings)) if (k !== "SAFE_PROXY_CREATION_CODE") console.log(`${k}=${v}`);

  const receipts = new Map(); // userOpHash -> receipt
  const ep = entryPoint.connect(bundlerKey);
  const handlers = {
    eth_chainId: async () => ethers.toQuantity(chainId),
    eth_supportedEntryPoints: async () => [entryPoint.target],
    eth_estimateUserOperationGas: async ([op]) => ({
      preVerificationGas: "0x186a0", // 100k
      verificationGasLimit: op.factory ? "0xc3500" : "0x61a80", // 800k / 400k
      callGasLimit: "0x7a120", // 500k
      paymasterVerificationGasLimit: "0x30d40", // 200k
      paymasterPostOpGasLimit: "0x186a0", // 100k
    }),
    eth_sendUserOperation: async ([op]) => {
      const packed = unpackRpc(op);
      const hash = await entryPoint.getUserOpHash(packed);
      // gas for the whole bundle: the operation's limits plus overhead
      const limits = [op.verificationGasLimit, op.callGasLimit, op.preVerificationGas, op.paymasterVerificationGasLimit || 0, op.paymasterPostOpGasLimit || 0];
      const gasLimit = limits.reduce((a, b) => a + BigInt(b), 250000n);
      const tx = await ep.handleOps([packed], bundlerKey.address, { gasLimit: gasLimit > 15000000n ? 15000000n : gasLimit });
      const rc = await tx.wait();
      for (const log of rc.logs) {
        try {
          const ev = entryPoint.interface.parseLog(log);
          if (ev && ev.name === "UserOperationEvent" && ev.args.userOpHash === hash) {
            receipts.set(hash, {
              userOpHash: hash, sender: ev.args.sender, success: ev.args.success, reason: "",
              actualGasCost: ethers.toQuantity(ev.args.actualGasCost), actualGasUsed: ethers.toQuantity(ev.args.actualGasUsed),
              receipt: { transactionHash: rc.hash, blockNumber: ethers.toQuantity(rc.blockNumber) },
            });
          }
        } catch (_) {}
      }
      return hash;
    },
    eth_getUserOperationReceipt: async ([hash]) => receipts.get(hash) || null,
  };

  http
    .createServer((req, res) => {
      let body = "";
      req.on("data", (d) => (body += d));
      req.on("end", async () => {
        let msg;
        try {
          msg = JSON.parse(body);
        } catch (_) {
          res.statusCode = 400;
          return res.end();
        }
        const reply = { jsonrpc: "2.0", id: msg.id };
        try {
          const h = handlers[msg.method];
          if (!h) throw Object.assign(new Error(`method ${msg.method} not supported`), { code: -32601 });
          reply.result = await h(msg.params || []);
        } catch (e) {
          const reason = e.revert ? `${e.revert.name}(${e.revert.args.join(", ")})` : e.shortMessage || e.message;
          reply.error = { code: e.code && Number.isInteger(e.code) ? e.code : -32500, message: reason };
        }
        res.setHeader("content-type", "application/json");
        res.end(JSON.stringify(reply));
      });
    })
    .listen(BUNDLER_PORT, () => console.log(`dev bundler listening on ${settings.BUNDLER_URL}`));
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
