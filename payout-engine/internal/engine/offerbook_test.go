package engine

import (
	"math/big"
	"testing"

	"trovo-payout-engine/internal/config"
	"trovo-payout-engine/internal/store"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Tokens on sale in the offer book are credited to their sellers, pro rata
// to what each open offer still sells, and the platform's wallets are found.
func TestOfferBookCreditAndPlatformWallets(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(store.AllModels()...)
	book := common.HexToAddress("0x00000000000000000000000000000000000b00c0")
	token := "0x00000000000000000000000000000000000070c0"
	seller1, seller2 := "0x0000000000000000000000000000000000000051", "0x0000000000000000000000000000000000000052"
	db.Create(&store.OfferBookOffer{ID: "1", Seller: seller1, SellToken: token, Remaining: "30", Open: true})
	db.Create(&store.OfferBookOffer{ID: "2", Seller: seller2, SellToken: token, Remaining: "10", Open: true})
	db.Create(&store.OfferBookOffer{ID: "3", Seller: seller2, SellToken: token, Remaining: "99", Open: false}) // closed
	mm := "0x00000000000000000000000000000000000000aa"
	db.Create(&store.TokenizedAsset{ID: "ta", MarketMakingWallet: &mm})
	db.Create(&store.User{ID: "7", Username: "atprofile"})
	db.Create(&store.UserWallet{ID: "0x00000000000000000000000000000000000000bb", UserID: "7"})

	e := &Engine{DB: db, Cfg: &config.Config{OfferBook: book, IssuingProfile: "atprofile", PayoutSafe: common.HexToAddress("0xcc")}}
	balances := map[string]*big.Int{key(book): big.NewInt(41), seller1: big.NewInt(5)}
	if err := e.creditOfferBookSellers(balances, token, nil); err != nil {
		t.Fatal(err)
	}
	// 41 held for 40 on sale: 30.75 -> 30, 10.25 -> 10, 1 stays with the book
	if balances[seller1].Int64() != 35 || balances[seller2].Int64() != 10 || balances[key(book)].Int64() != 1 {
		t.Fatalf("balances %v", balances)
	}
	ex, err := e.platformWallets(&store.ProceedPayout{TokenizedAssetID: "ta"})
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range []string{key(book), mm, "0x00000000000000000000000000000000000000bb", key(common.HexToAddress("0xcc"))} {
		if ex[a] == "" {
			t.Fatalf("%s not excluded: %v", a, ex)
		}
	}
}
