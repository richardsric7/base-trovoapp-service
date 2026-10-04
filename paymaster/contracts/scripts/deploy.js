// Deploys and configures TrovoTokenPaymaster from a Vault KV v2 secret.
//
//   npm run build && node scripts/deploy.js [--dry-run]
//
// Bootstrap environment: VAULT_ADDR, VAULT_TOKEN, VAULT_KV_MOUNT (default
// "secret"), VAULT_NAMESPACE (optional), PAYMASTER_DEPLOY_SECRET_PATH
// (default "trovo/paymaster/deploy"). Every other setting is read from that
// secret - see ../CONFIGURATION.md for each key.
//
// Steps: deploy (deployer as temporary owner) -> addStake -> deposit ->
// setToken for each gas token -> transferOwnership to the owner Safe.
// The result is written to deployments/<chainId>.json.

const fs = require("fs");
const path = require("path");
const { ethers } = require("ethers");
const { loadConfig } = require("./lib/vault");

const DEFAULT_ENTRYPOINT = "0x0000000071727De22E5E9d8BAf0edAc6f37da032"; // ERC-4337 v0.7
const IENTRYPOINT_INTERFACE_ID = "0x915074d8"; // type(IEntryPoint).interfaceId, v0.7

function required(cfg, key) {
  const v = cfg.get(key);
  if (v === undefined || v.trim() === "") throw new Error(`${key} is required`);
  return v.trim();
}

function optional(cfg, key, def) {
  const v = cfg.get(key);
  return v === undefined || v.trim() === "" ? def : v.trim();
}

function address(key, v) {
  if (!ethers.isAddress(v)) throw new Error(`${key} is not an address: ${v}`);
  return ethers.getAddress(v);
}

function parseTokens(raw) {
  let list;
  try {
    list = JSON.parse(raw);
  } catch (e) {
    throw new Error(`GAS_TOKENS is not valid JSON: ${e.message}`);
  }
  if (!Array.isArray(list) || list.length === 0) throw new Error("GAS_TOKENS must be a non-empty JSON array");
  return list.map((t, i) => {
    const minRate = BigInt(t.minRate);
    const maxRate = BigInt(t.maxRate);
    if (minRate <= 0n || minRate > maxRate) throw new Error(`GAS_TOKENS[${i}]: need 0 < minRate <= maxRate`);
    return { symbol: String(t.symbol || `token${i}`), address: address(`GAS_TOKENS[${i}].address`, t.address), minRate, maxRate };
  });
}

function readConfig(cfg) {
  const signers = required(cfg, "QUOTE_SIGNER_ADDRESSES")
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean)
    .map((s, i) => address(`QUOTE_SIGNER_ADDRESSES[${i}]`, s));
  if (signers.length === 0) throw new Error("QUOTE_SIGNER_ADDRESSES has no addresses");
  return {
    rpcUrl: required(cfg, "RPC_URL"),
    chainId: BigInt(required(cfg, "CHAIN_ID")),
    deployerKey: required(cfg, "DEPLOYER_PRIVATE_KEY"),
    entryPoint: address("ENTRYPOINT_ADDRESS", optional(cfg, "ENTRYPOINT_ADDRESS", DEFAULT_ENTRYPOINT)),
    owner: address("OWNER_ADDRESS", required(cfg, "OWNER_ADDRESS")),
    allowEoaOwner: optional(cfg, "ALLOW_EOA_OWNER", "false") === "true",
    quoteSigners: signers,
    pauser: address("PAUSER_ADDRESS", required(cfg, "PAUSER_ADDRESS")),
    postOpGasOverhead: BigInt(optional(cfg, "POST_OP_GAS_OVERHEAD", "45000")),
    stake: ethers.parseEther(optional(cfg, "STAKE_AMOUNT_ETH", "0.1")),
    unstakeDelaySec: Number(optional(cfg, "UNSTAKE_DELAY_SEC", "86400")),
    deposit: ethers.parseEther(optional(cfg, "INITIAL_DEPOSIT_ETH", "0.05")),
    tokens: parseTokens(required(cfg, "GAS_TOKENS")),
  };
}

function artifact(name) {
  const file = path.join(__dirname, "..", "artifacts", "src", `${name}.sol`, `${name}.json`);
  if (!fs.existsSync(file)) throw new Error(`${file} not found - run "npm run build" first`);
  return JSON.parse(fs.readFileSync(file, "utf8"));
}

async function main() {
  const dryRun = process.argv.includes("--dry-run");
  const cfg = await loadConfig("PAYMASTER_DEPLOY_SECRET_PATH", "trovo/paymaster/deploy");
  const c = readConfig(cfg);
  console.log(`config: ${cfg.source}`);

  const provider = new ethers.JsonRpcProvider(c.rpcUrl);
  const network = await provider.getNetwork();
  if (network.chainId !== c.chainId) throw new Error(`RPC_URL is chain ${network.chainId}, CHAIN_ID says ${c.chainId}`);
  const wallet = new ethers.Wallet(c.deployerKey, provider);
  // tracks nonces locally: the provider's cached nonce can lag between steps
  const deployer = new ethers.NonceManager(wallet);

  if ((await provider.getCode(c.entryPoint)) === "0x") throw new Error(`no EntryPoint deployed at ${c.entryPoint}`);
  const ep = new ethers.Contract(c.entryPoint, ["function supportsInterface(bytes4) view returns (bool)"], provider);
  if (!(await ep.supportsInterface(IENTRYPOINT_INTERFACE_ID))) throw new Error(`${c.entryPoint} is not an EntryPoint v0.7`);
  if (!c.allowEoaOwner && (await provider.getCode(c.owner)) === "0x") {
    throw new Error(`OWNER_ADDRESS ${c.owner} has no code; it should be the treasury Safe (set ALLOW_EOA_OWNER=true only for testnets)`);
  }
  for (const t of c.tokens) {
    if ((await provider.getCode(t.address)) === "0x") throw new Error(`no token contract at ${t.address} (${t.symbol})`);
  }

  const needed = c.stake + c.deposit;
  const balance = await provider.getBalance(wallet.address);
  console.log(`chain ${network.chainId}, deployer ${wallet.address} (${ethers.formatEther(balance)} ETH), EntryPoint ${c.entryPoint}`);
  console.log(`owner ${c.owner}, pauser ${c.pauser}, quote signers ${c.quoteSigners.join(", ")}`);
  console.log(`stake ${ethers.formatEther(c.stake)} ETH (unstake delay ${c.unstakeDelaySec}s), deposit ${ethers.formatEther(c.deposit)} ETH, postOp overhead ${c.postOpGasOverhead} gas`);
  for (const t of c.tokens) console.log(`gas token ${t.symbol} ${t.address}: rate bounds [${t.minRate}, ${t.maxRate}] base units per ETH`);
  if (balance < needed) throw new Error(`deployer needs at least ${ethers.formatEther(needed)} ETH plus gas`);
  if (dryRun) {
    console.log("dry run: nothing deployed");
    return;
  }

  const a = artifact("TrovoTokenPaymaster");
  const factory = new ethers.ContractFactory(a.abi, a.bytecode, deployer);
  const paymaster = await factory.deploy(c.entryPoint, wallet.address, c.quoteSigners, c.pauser, c.postOpGasOverhead);
  const deployTx = paymaster.deploymentTransaction();
  await paymaster.waitForDeployment();
  const addr = await paymaster.getAddress();
  console.log(`deployed TrovoTokenPaymaster at ${addr} (tx ${deployTx.hash})`);

  const step = async (label, txPromise) => {
    const tx = await txPromise;
    const r = await tx.wait();
    if (r.status !== 1) throw new Error(`${label} failed (tx ${tx.hash})`);
    console.log(`${label}: ${tx.hash}`);
    return tx.hash;
  };
  const txs = { deploy: deployTx.hash };
  txs.addStake = await step("addStake", paymaster.addStake(c.unstakeDelaySec, { value: c.stake }));
  txs.deposit = await step("deposit", paymaster.deposit({ value: c.deposit }));
  txs.setToken = {};
  for (const t of c.tokens) {
    txs.setToken[t.symbol] = await step(`setToken ${t.symbol}`, paymaster.setToken(t.address, true, t.minRate, t.maxRate));
  }
  txs.transferOwnership = await step("transferOwnership", paymaster.transferOwnership(c.owner));
  if ((await paymaster.owner()) !== c.owner) throw new Error("ownership transfer did not take effect");

  const out = {
    chainId: network.chainId.toString(),
    paymaster: addr,
    entryPoint: c.entryPoint,
    owner: c.owner,
    pauser: c.pauser,
    quoteSigners: c.quoteSigners,
    postOpGasOverhead: c.postOpGasOverhead.toString(),
    tokens: c.tokens.map((t) => ({ symbol: t.symbol, address: t.address, minRate: t.minRate.toString(), maxRate: t.maxRate.toString() })),
    transactions: txs,
    deployedAt: new Date().toISOString(),
  };
  const dir = path.join(__dirname, "..", "deployments");
  fs.mkdirSync(dir, { recursive: true });
  const file = path.join(dir, `${network.chainId}.json`);
  fs.writeFileSync(file, JSON.stringify(out, null, 2) + "\n");
  console.log(`wrote ${file}`);
  console.log(`next: set PAYMASTER_ADDRESS=${addr} in the quote service's Vault secret`);
}

main().catch((e) => {
  console.error(e.message || e);
  process.exit(1);
});
