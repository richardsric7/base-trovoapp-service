// Deploys and configures TrovoOfferBook from a Vault KV v2 secret.
//
//   npm run build && node scripts/deploy.js [--dry-run]
//
// Bootstrap environment: VAULT_ADDR, VAULT_TOKEN, VAULT_KV_MOUNT (default
// "secret"), VAULT_NAMESPACE (optional), MARKET_DEPLOY_SECRET_PATH
// (default "trovo/market/deploy"). Every other setting is read from that
// secret - see ../CONFIGURATION.md for each key.
//
// Steps: deploy (deployer as temporary owner, with the authorizers) ->
// setTradable for each listed token -> transferOwnership to the owner
// Safe. The result is written to deployments/<chainId>.json.

const fs = require("fs");
const path = require("path");
const { ethers } = require("ethers");
const { loadConfig } = require("./lib/vault");

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

function addressList(cfg, key, allowEmpty) {
  const list = optional(cfg, key, "")
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean)
    .map((s, i) => address(`${key}[${i}]`, s));
  if (!allowEmpty && list.length === 0) throw new Error(`${key} has no addresses`);
  return list;
}

function readConfig(cfg) {
  return {
    rpcUrl: required(cfg, "RPC_URL"),
    chainId: BigInt(required(cfg, "CHAIN_ID")),
    deployerKey: required(cfg, "DEPLOYER_PRIVATE_KEY"),
    owner: address("OWNER_ADDRESS", required(cfg, "OWNER_ADDRESS")),
    allowEoaOwner: optional(cfg, "ALLOW_EOA_OWNER", "false") === "true",
    authorizers: addressList(cfg, "AUTHORIZER_ADDRESSES", false),
    tradable: addressList(cfg, "TRADABLE_TOKENS", true),
  };
}

function artifact(name) {
  const file = path.join(__dirname, "..", "artifacts", "src", `${name}.sol`, `${name}.json`);
  if (!fs.existsSync(file)) throw new Error(`${file} not found - run "npm run build" first`);
  return JSON.parse(fs.readFileSync(file, "utf8"));
}

async function main() {
  const dryRun = process.argv.includes("--dry-run");
  const cfg = await loadConfig("MARKET_DEPLOY_SECRET_PATH", "trovo/market/deploy");
  const c = readConfig(cfg);
  console.log(`config: ${cfg.source}`);

  const provider = new ethers.JsonRpcProvider(c.rpcUrl);
  const network = await provider.getNetwork();
  if (network.chainId !== c.chainId) throw new Error(`RPC_URL is chain ${network.chainId}, CHAIN_ID says ${c.chainId}`);
  const wallet = new ethers.Wallet(c.deployerKey, provider);
  const deployer = new ethers.NonceManager(wallet);

  if (!c.allowEoaOwner && (await provider.getCode(c.owner)) === "0x") {
    throw new Error(`OWNER_ADDRESS ${c.owner} has no code; it should be the platform admin Safe (set ALLOW_EOA_OWNER=true only for testnets)`);
  }
  for (const t of c.tradable) {
    if ((await provider.getCode(t)) === "0x") throw new Error(`no token contract at ${t}`);
  }
  const balance = await provider.getBalance(wallet.address);
  console.log(`chain ${network.chainId}, deployer ${wallet.address} (${ethers.formatEther(balance)} ETH)`);
  console.log(`owner ${c.owner}, authorizers ${c.authorizers.join(", ")}`);
  console.log(`tradable tokens: ${c.tradable.join(", ") || "(none yet)"}`);
  if (dryRun) {
    console.log("dry run: nothing deployed");
    return;
  }

  const a = artifact("TrovoOfferBook");
  const book = await new ethers.ContractFactory(a.abi, a.bytecode, deployer).deploy(wallet.address, c.authorizers);
  const deployTx = book.deploymentTransaction();
  await book.waitForDeployment();
  const addr = await book.getAddress();
  console.log(`deployed TrovoOfferBook at ${addr} (tx ${deployTx.hash})`);

  const step = async (label, txPromise) => {
    const tx = await txPromise;
    const r = await tx.wait();
    if (r.status !== 1) throw new Error(`${label} failed (tx ${tx.hash})`);
    console.log(`${label}: ${tx.hash}`);
    return tx.hash;
  };
  const txs = { deploy: deployTx.hash, setTradable: {} };
  for (const t of c.tradable) txs.setTradable[t] = await step(`setTradable ${t}`, book.setTradable(t, true));
  txs.transferOwnership = await step("transferOwnership", book.transferOwnership(c.owner));
  if ((await book.owner()) !== c.owner) throw new Error("ownership transfer did not take effect");

  const out = {
    chainId: network.chainId.toString(),
    offerBook: addr,
    owner: c.owner,
    authorizers: c.authorizers,
    tradable: c.tradable,
    transactions: txs,
    deployedAt: new Date().toISOString(),
  };
  const dir = path.join(__dirname, "..", "deployments");
  fs.mkdirSync(dir, { recursive: true });
  const file = path.join(dir, `${network.chainId}.json`);
  fs.writeFileSync(file, JSON.stringify(out, null, 2) + "\n");
  console.log(`wrote ${file}`);
  console.log(`next: set OFFER_BOOK_ADDRESS=${addr} in app-backend`);
}

main().catch((e) => {
  console.error(e.message || e);
  process.exit(1);
});
