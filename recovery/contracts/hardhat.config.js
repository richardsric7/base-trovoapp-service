// Hardhat 2 configuration for the Trovo recovery contracts (vendored Candide Social Recovery Module).
//
// Compiler: the Solidity compiler is taken from the pinned `solc` npm package
// (the official WebAssembly build) instead of being downloaded from
// binaries.soliditylang.org at build time. Builds are then reproducible and
// work in environments where that host is unreachable.
require("@nomicfoundation/hardhat-toolbox");
const { subtask } = require("hardhat/config");
const { TASK_COMPILE_SOLIDITY_GET_SOLC_BUILD } = require("hardhat/builtin-tasks/task-names");

const SOLC_VERSION = "0.8.28";

subtask(TASK_COMPILE_SOLIDITY_GET_SOLC_BUILD, async (args, hre, runSuper) => {
  if (args.solcVersion !== SOLC_VERSION) {
    return runSuper();
  }
  const solc = require("solc");
  return {
    compilerPath: require.resolve("solc/soljson.js"),
    isSolcJs: true,
    version: args.solcVersion,
    longVersion: solc.version(),
  };
});

module.exports = {
  solidity: {
    version: SOLC_VERSION,
    settings: {
      optimizer: { enabled: true, runs: 200 },
      evmVersion: "cancun",
    },
  },
  paths: {
    sources: "./src",
    tests: "./test",
    artifacts: "./artifacts",
    cache: "./cache",
  },
  networks: {
    // `npx hardhat node` (or any local node); LOCAL_RPC_URL overrides the URL
    localhost: { url: process.env.LOCAL_RPC_URL || "http://127.0.0.1:8545" },
  },
};
