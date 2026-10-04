// Deploys Candide's Social Recovery Module (vendored v0.2.0) from a Vault
// KV v2 secret.
//
//   npm run build && node scripts/deploy.js [--dry-run]
//
// Bootstrap environment: VAULT_ADDR, VAULT_TOKEN, VAULT_KV_MOUNT (default
// "secret"), VAULT_NAMESPACE (optional), RECOVERY_DEPLOY_SECRET_PATH
// (default "trovo/recovery/deploy"). See ../CONFIGURATION.md.
//
// The module has no owner and no settings besides its recovery period,
// fixed at deployment. The result is written to deployments/<chainId>.json.

const fs = require("fs");
const path = require("path");
const { ethers } = require("ethers");
const { loadConfig } = require("./lib/vault");

function required(cfg, key) {
  const v = cfg.get(key);
  if (v === undefined || v.trim() === "") throw new Error(`${key} is required`);
  return v.trim();
}

async function main() {
  const dryRun = process.argv.includes("--dry-run");
  const cfg = await loadConfig("RECOVERY_DEPLOY_SECRET_PATH", "trovo/recovery/deploy");
  const rpcUrl = required(cfg, "RPC_URL");
  const chainId = BigInt(required(cfg, "CHAIN_ID"));
  const period = BigInt(required(cfg, "RECOVERY_PERIOD_SECONDS"));
  if (period < 86400n && (cfg.get("ALLOW_SHORT_PERIOD") || "") !== "true") {
    throw new Error("RECOVERY_PERIOD_SECONDS is under a day; set ALLOW_SHORT_PERIOD=true only for testnets");
  }
  console.log(`config: ${cfg.source}`);

  const provider = new ethers.JsonRpcProvider(rpcUrl);
  const network = await provider.getNetwork();
  if (network.chainId !== chainId) throw new Error(`RPC_URL is chain ${network.chainId}, CHAIN_ID says ${chainId}`);
  const wallet = new ethers.Wallet(required(cfg, "DEPLOYER_PRIVATE_KEY"), provider);
  console.log(`chain ${chainId}, deployer ${wallet.address}, recovery period ${period}s (${Number(period) / 86400} days)`);
  if (dryRun) {
    console.log("dry run: nothing deployed");
    return;
  }

  const file = path.join(__dirname, "..", "artifacts", "src", "candide", "modules", "social_recovery", "SocialRecoveryModule.sol", "SocialRecoveryModule.json");
  if (!fs.existsSync(file)) throw new Error(`${file} not found - run "npm run build" first`);
  const a = JSON.parse(fs.readFileSync(file, "utf8"));
  const module = await new ethers.ContractFactory(a.abi, a.bytecode, wallet).deploy(period);
  const tx = module.deploymentTransaction();
  await module.waitForDeployment();
  const addr = await module.getAddress();
  if ((await module.VERSION()) !== "0.2.0") throw new Error("unexpected module version");
  console.log(`deployed SocialRecoveryModule v0.2.0 at ${addr} (tx ${tx.hash})`);

  const out = { chainId: chainId.toString(), recoveryModule: addr, recoveryPeriodSeconds: period.toString(), transaction: tx.hash, deployedAt: new Date().toISOString() };
  const dir = path.join(__dirname, "..", "deployments");
  fs.mkdirSync(dir, { recursive: true });
  fs.writeFileSync(path.join(dir, `${chainId}.json`), JSON.stringify(out, null, 2) + "\n");
  console.log(`next: set RECOVERY_MODULE_ADDRESS=${addr} in app-backend`);
}

main().catch((e) => {
  console.error(e.message || e);
  process.exit(1);
});
