package users

import (
	"context"
	"log"
	"time"
	"trovo-wallet-api/internal/basetxn"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	// db "trovo-wallet-api/internal/db"

	"github.com/ethereum/go-ethereum/common"
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

// fetchInitialAccounts lists every holder of token (a tokenized asset's
// contract) with its balance in whole tokens, from the token's Transfer
// events since fromBlock (see network.TokenHolders).
func fetchInitialAccounts(ctx context.Context, gc *sharedconfig.GlobalConfig, token common.Address, fromBlock uint64) ([]BasicBalance, error) {
	dec, err := network.AssetDecimals(ctx, gc.BantuExpansionClient, basetxn.CreditAsset{Code: "TOKEN", Issuer: token.Hex()})
	if err != nil {
		return nil, err
	}
	holders, err := network.TokenHolders(ctx, gc.BantuExpansionClient, token, fromBlock)
	if err != nil {
		return nil, err
	}
	out := make([]BasicBalance, 0, len(holders))
	for a, b := range holders {
		out = append(out, BasicBalance{Address: a.Hex(), Balance: decimal.NewFromBigInt(b, -int32(dec)).InexactFloat64()})
	}
	return out, nil
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
	if payout.TokenizedAsset.ContractAddress == nil {
		return // not minted: there are no holders of a token contract yet
	}

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
		PayoutContractAddress:          *payout.TokenizedAsset.ContractAddress,
		BeneficiaryAddress:             publicKey,
		ConfirmedTokenizedAssetBalance: bal,
		AmountToReceive:                amountToReceive.InexactFloat64(),
		CannotReceiveAsset:             canReceiveAsset,
	}
	gc.DB.Omit(clause.Associations).Create(&scheduleItem)
}
