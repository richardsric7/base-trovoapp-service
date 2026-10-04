// Local development deployment of the market contracts, on the same local
// node as paymaster/contracts/scripts/local-stack.js: TrovoOfferBook (owned
// by hardhat account #0, authorizer hardhat #2), a mock cNGN stablecoin and
// a mock internal balance token, all listed as tradable.
//
//   npx hardhat run scripts/local-market.js --network localhost
//
// It prints the settings and writes them to deployments/local-market.json.
// Tokenized asset tokens are deployed per asset (TokenizedAsset, owned by
// the asset's issuing Safe) and listed with setTradable by the owner.

const fs = require("fs");
const path = require("path");
const { ethers } = require("hardhat");

const AUTHORIZER = process.env.OFFER_AUTHORIZER_ADDRESS || "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"; // hardhat #2

async function main() {
  const [deployer] = await ethers.getSigners();
  const book = await ethers.deployContract("TrovoOfferBook", [deployer.address, [AUTHORIZER]]);
  const cngn = await ethers.deployContract("MockERC20", ["cNGN", "CNGN", 6]);
  const internal = await ethers.deployContract("MockERC20", ["Internal NGN", "NGNI", 18]);
  for (const t of [cngn, internal]) await (await book.setTradable(t, true)).wait();

  const settings = {
    OFFER_BOOK_ADDRESS: book.target,
    OFFER_BOOK_OWNER: deployer.address,
    OFFER_AUTHORIZER_ADDRESS: AUTHORIZER,
    MOCK_CNGN_ADDRESS: cngn.target,
    MOCK_INTERNAL_BALANCE_ADDRESS: internal.target,
    TOKENIZED_ASSET_ARTIFACT: path.join(__dirname, "..", "artifacts", "src", "TokenizedAsset.sol", "TokenizedAsset.json"),
  };
  const dir = path.join(__dirname, "..", "deployments");
  fs.mkdirSync(dir, { recursive: true });
  fs.writeFileSync(path.join(dir, "local-market.json"), JSON.stringify(settings, null, 2) + "\n");
  for (const [k, v] of Object.entries(settings)) console.log(`${k}=${v}`);
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
