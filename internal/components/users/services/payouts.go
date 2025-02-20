package users

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
)

const (
	assetCode    = "KGM"
	assetIssuer  = "ISSUER_PUBLIC_KEY" // Replace with actual issuer public key
	batchSize    = 10000               // Batch size for saving to file/database
	saveInterval = 5 * time.Minute     // How often to save to disk
	pageLimit    = 200                 // Number of records per page for initial fetch
	rateLimit    = time.Second         // Wait time between requests
)

func main() {
	// Initialize Horizon client for Mainnet
	client := horizonclient.DefaultPublicNetClient

	// First, fetch all current accounts holding KGM using batch request
	initialAccounts, err := fetchInitialAccounts(client)
	if err != nil {
		log.Fatalf("Failed to fetch initial accounts: %v", err)
	}

	// Create a map to track accounts holding KGM
	accountsWithKGM := make(map[string]bool)
	// var mutex sync.Mutex

	// Initialize with accounts from initial fetch
	for _, accountID := range initialAccounts {
		accountsWithKGM[accountID.PublicKey] = true
	}

	// Create a channel to receive streaming events
	// events := make(chan horizon.Account)

	// // Start streaming account updates for these accounts
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// Final save before exiting
	// fmt.Printf("Streaming completed. Total unique accounts holding KGM: %d\n", len(accountsWithKGM))
}

type BasicBalance struct {
	PublicKey string
	Balance   float64
}

// fetchInitialAccounts fetches all current accounts holding the KGM asset
func fetchInitialAccounts(client *horizonclient.Client) ([]BasicBalance, error) {
	var allAccounts []BasicBalance
	var cursor string

	for {
		request := horizonclient.AccountsRequest{
			Asset:  fmt.Sprintf("%s:%s", assetCode, assetIssuer),
			Limit:  pageLimit,
			Cursor: cursor,
		}

		accounts, err := client.Accounts(request)
		if err != nil {
			return nil, fmt.Errorf("error fetching accounts: %v", err)
		}

		for _, account := range accounts.Embedded.Records {
			if ok, bal := hasPositiveBalance(account, assetCode, assetIssuer); ok {
				allAccounts = append(allAccounts, BasicBalance{
					PublicKey: account.AccountID,
					Balance:   decimal.RequireFromString(bal).InexactFloat64(),
				})
			}
		}

		if len(accounts.Links.Next.Href) == 0 {
			break
		}
		uri, _ := url.Parse(accounts.Links.Next.Href)
		cursor = uri.Query().Get("cursor")
		time.Sleep(rateLimit) // Respect rate limits
	}

	return allAccounts, nil
}

// hasPositiveBalance checks if an account holds a positive balance of the KGM asset
func hasPositiveBalance(account horizon.Account, assetCode, assetIssuer string) (bool, string) {
	for _, balance := range account.Balances {
		if balance.Asset.Code == assetCode && balance.Asset.Issuer == assetIssuer {
			amount := parseBalance(balance.Balance)
			return amount > 0, balance.Balance
		}
	}
	return false, "0"
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

// // contains checks if a slice contains a string
// func contains(slice []string, item string) bool {
// 	for _, s := range slice {
// 		if s == item {
// 			return true
// 		}
// 	}
// 	return false
// }
