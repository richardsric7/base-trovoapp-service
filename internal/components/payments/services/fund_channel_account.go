package payments

import (
	"log"
	"os"
	bc "trovo-wallet-api/internal/components/users/blockchain"

	bantupaysdk "github.com/bantublockchain/bantupaysdk-go/payments"

	"github.com/shopspring/decimal"
	"github.com/stellar/go/keypair"
)

//FundChannelAccount makes payment from faucet on server
func FundChannelAccount() (paymentHash string, err error) {
	chanAccount := os.Getenv("CHANNEL_ACCOUNT")
	kp := keypair.MustParseFull(os.Getenv("XBN_FAUCET"))

	channelkp := keypair.MustParseFull(chanAccount)
	log.Println("faucet public key:", kp.Address())
	log.Println("channel account public key:", channelkp.Address())
	amount := "12.6"
	var pay bool
	sourceAccount, sError := bc.GetBlockchainAccountDetail(kp.Address())
	channelAccount, err := bc.GetBlockchainAccountDetail(channelkp.Address())
	if sError == nil {

		for _, b := range sourceAccount.Balances {
			if b.Issuer == "" {
				log.Printf("source account balance:%v\n", b.Balance)
			}
		}
	}

	if err != nil {
		if err.Error() == "error-blockchain-account-not-activated" {
			pay = true

		} else {
			log.Println("[FundChannelAccount]", err)
			return "", err
		}
	} else {
		for _, v := range channelAccount.Balances {
			if len(v.Issuer) == 0 {
				log.Printf("channel account balance:%v\n", v.Balance)

				//balance has dropped below 12.5XBN then top up
				if decimal.RequireFromString(v.Balance).LessThan(decimal.NewFromFloat(12.5)) {
					amount = decimal.NewFromFloat(13.6).Sub(decimal.RequireFromString(v.Balance)).Truncate(7).String()
					pay = true
					log.Println("channel account is below minimum balance...topping up...")

				}
			}
		}
	}
	if !pay {
		log.Println("channel account does not need funding")
		return "good", nil
	}
	p := bantupaysdk.PaymentInstance()
	p.Amount = amount
	p.Destination = channelkp.Address()
	p.Memo = "ChannelAccountFunding"
	err = p.ExpressPay(os.Getenv("BANTUPAY_BASE_URL"), "bantupay", kp.Seed(), "", "")
	if err != nil {
		log.Println("error funding channel account at this time:", err)
		return "", err
	}
	return p.TransactionID, nil

}
