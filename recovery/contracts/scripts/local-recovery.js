// Local deployment of the recovery module on the paymaster local stack's
// node, with a short recovery period (RECOVERY_PERIOD_SECONDS, default 60)
// for end-to-end tests. Writes deployments/local-recovery.json.
//
//   npx hardhat run scripts/local-recovery.js --network localhost

const fs = require("fs");
const path = require("path");
const { ethers } = require("hardhat");

async function main() {
  const period = BigInt(process.env.RECOVERY_PERIOD_SECONDS || "60");
  const module = await ethers.deployContract("SocialRecoveryModule", [period]);
  await module.waitForDeployment();
  const settings = { RECOVERY_MODULE_ADDRESS: module.target, RECOVERY_PERIOD_SECONDS: period.toString() };
  const dir = path.join(__dirname, "..", "deployments");
  fs.mkdirSync(dir, { recursive: true });
  fs.writeFileSync(path.join(dir, "local-recovery.json"), JSON.stringify(settings, null, 2) + "\n");
  for (const [k, v] of Object.entries(settings)) console.log(`${k}=${v}`);
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
