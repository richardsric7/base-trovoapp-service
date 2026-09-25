package users

import (
	"fmt"
	"log"
	"time"
	"trovo-wallet-api/internal/basetxn"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	// db "trovo-wallet-api/internal/db"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

const (
	assetCode       = "KGM"
	contractAddress = "ISSUER_PUBLIC_KEY" // Replace with actual issuer public key
	batchSize       = 10000               // Batch size for saving to file/database
	saveInterval    = 5 * time.Minute     // How often to save to disk
	pageLimit       = 200                 // Number of records per page for initial fetch
	rateLimit       = time.Second         // Wait time between requests
)

func main() {
	// Initialize Horizon client for Mainnet
	// client := horizonclient.DefaultPublicNetClient

	// First, fetch all current accounts holding KGM using batch request
	// initialAccounts, err := fetchInitialAccounts(client)
	// if err != nil {
	// 	log.Fatalf("Failed to fetch initial accounts: %v", err)
	// }

	// Create a map to track accounts holding KGM
	// accountsWithKGM := make(map[string]bool)
	// var mutex sync.Mutex

	// Initialize with accounts from initial fetch
	// for _, accountID := range initialAccounts {
	// 	accountsWithKGM[accountID.Address] = true
	// }

	// Create a channel to receive streaming events
	// events := make(chan horizon.Account)

	// // Start streaming account updates for these accounts
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// Final save before exiting
	// fmt.Printf("Streaming completed. Total unique accounts holding KGM: %d\n", len(accountsWithKGM))
}

type BasicBalance struct {
	Address string
	Balance float64
}

// fetchInitialAccounts enumerated every Stellar account holding the KGM
// asset via Horizon's global "accounts holding this asset" query. Base
// has no equivalent enumeration API for a B20 (ERC-20-shaped) token -
// finding every holder needs indexing Transfer event logs (e.g. via a
// subgraph or a log-scanning service), a follow-up out of scope for this
// alteration pass. Stubbed to return no accounts rather than crash;
// currently unreferenced (the one call site is commented out).
func fetchInitialAccounts(gc *sharedconfig.GlobalConfig) ([]BasicBalance, error) {
	return nil, fmt.Errorf("holder enumeration is not available on Base yet - needs an event-log indexer")
}

// parseBalance converts a balance string to float64
func parseBalance(balanceStr string) float64 {

	balDec, err := decimal.NewFromString(balanceStr)
	if err != nil {
		log.Printf("Warning: Failed to parse balance '%s': %v", balanceStr, err)
		return 0
	}

	return balDec.Truncate(7).InexactFloat64()
}

func processData(balance, publicKey string, payout *userModels.ProceedPayout, gc *sharedconfig.GlobalConfig) {
	var scheduleItem userModels.TokenizedAssetPayoutSchedule

	// var payout userModels.ProceedPayout
	bal := decimal.RequireFromString(balance).InexactFloat64()
	amountToReceive := decimal.NewFromFloat(payout.AmountPerTokenizedAssetHeld * bal).Truncate(7)
	if amountToReceive.IsZero() {
		return
	}
	ca, _ := userModels.Currency(*payout.TokenizedAsset.ProceedPayoutCurrency).GetCurratedAsset(gc)

	payoutAsset := basetxn.CreditAsset{Code: ca.AssetCode, Issuer: ca.ContractAddress}

	_, trustsAsset, _, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, publicKey, payoutAsset)

	var canReceiveAsset int
	if trustsAsset {
		canReceiveAsset = 1
	}
	scheduleItem = userModels.TokenizedAssetPayoutSchedule{
		ID:                             uuid.NewString(),
		TokenizedAssetID:               payout.TokenizedAssetID,
		Batch:                          payout.Batch,
		PayoutAssetCode:                *payout.TokenizedAsset.AssetCode,
		PayoutContractAddress:          *payout.TokenizedAsset.IssuingWalletAddress,
		BeneficiaryAddress:             publicKey,
		ConfirmedTokenizedAssetBalance: bal,
		AmountToReceive:                amountToReceive.InexactFloat64(),
		CannotReceiveAsset:             canReceiveAsset,
	}
	gc.DB.Omit(clause.Associations).Create(&scheduleItem)
}
