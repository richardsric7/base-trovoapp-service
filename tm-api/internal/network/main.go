// Package network is tm-api's Base (EVM) blockchain access point - the
// equivalent of the old Stellar Horizon client this service used to talk
// to. Only GetBlockchainClient and GetBlockchainNetworkPassPhrase have any
// real callers in this service (the vault-signer feature); everything
// else this package used to hold (TempAccountKeypair,
// BlockchainAccountProperties, the SubmitXdrWithSignature* family,
// GetBlockchainBaseReserve, GetBlockchainSwapDestinationMin) had zero
// callers anywhere in tm-api and was removed rather than ported.
package network

import (
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
)

// GetBlockchainNetworkPassPhrase is vestigial on Base (Stellar used a
// network passphrase for signature domain separation; Base's chain ID,
// baked into every EIP-1559 transaction, serves that role instead). Kept,
// matching app-backend's own basetxn.Transaction.Sign, so call sites that
// still pass its result through don't need to change shape.
func GetBlockchainNetworkPassPhrase() string {
	return os.Getenv("BLOCKCHAIN_NETWORK_PASSPHRASE")
}

// GetBlockchainClient returns the Base JSON-RPC client, the Base
// equivalent of Stellar's Horizon client.
func GetBlockchainClient() *ethclient.Client {
	url := os.Getenv("BASE_RPC_URL")
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Panicf("[GetBlockchainClient] invalid BASE_RPC_URL %q: %v", url, err)
	}
	return client
}

// GetBlockchainChainID returns the Base chain ID transactions are signed
// and submitted against (BASE_CHAIN_ID, e.g. 8453 for Base mainnet, 84532
// for Base Sepolia) - matches app-backend's own network.GetBlockchainChainID.
func GetBlockchainChainID() *big.Int {
	raw := strings.TrimSpace(os.Getenv("BASE_CHAIN_ID"))
	id, ok := new(big.Int).SetString(raw, 10)
	if !ok {
		// Base Sepolia testnet - a safe default so a misconfigured
		// deployment fails against a public testnet, not mainnet.
		return big.NewInt(84532)
	}
	return id
}
