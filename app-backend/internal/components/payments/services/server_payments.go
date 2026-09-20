package payments

import (
	"fmt"
	"trovo-wallet-api/internal/basetxn"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/network"

	"os"

	"github.com/ecnepsnai/discord"
	"github.com/shopspring/decimal"
)

// //DoFaucetBlockchainPayment makes payment from faucet on server
// func DoFaucetBlockchainPayment(faucetSecret, receiver, assetCode, assetIssuer, amount, memo string, db *gorm.DB) (returnedPaymentInfo *paymentModels.PaymentInfo, err error) {
// 	var paymentInfo paymentModels.PaymentInfo
// 	var asset basetxn.Asset
// 	kp := evmkeypair.MustParseFull(faucetSecret)

// 	// faucetPK := kp.Address()
// 	paymentInfo = paymentModels.PaymentInfo{
// 		Destination: receiver, Memo: memo, AssetIssuer: assetIssuer,
// 		AssetCode: assetCode, Amount: amount,
// 	}
// 	if len(assetIssuer) == 42 {
// 		asset = basetxn.CreditAsset{Code: assetCode, Issuer: assetIssuer}
// 	} else {
// 		asset = basetxn.NativeAsset{}
// 	}

// 	receiverAccount, err := usersdb.GetUserInfo(receiver, db)
// 	if err != nil {
// 		log.Println("[DoFaucetBlockchainPayment]", err)
// 		return nil, err
// 	}
// 	client := network.GetBlockchainClient()
// 	var receiverStatus bool
// 	sourceAccount, _ := bc.GetBlockchainAccountDetail(kp.Address())
// 	_, err = bc.GetBlockchainAccountDetail(receiverAccount.Address)
// 	if err != nil {
// 		if err.Error() == "error-blockchain-account-not-activated" {
// 			receiverStatus = true

// 		} else {
// 			log.Println("[DoFaucetBlockchainPayment]", err)
// 			return nil, err
// 		}
// 	}

// 	if assetCode == "BNR" {
// 		asset = basetxn.NativeAsset{}
// 		amountFloat, _ := decimal.NewFromString(amount)
// 		xbnEquivalent := amountFloat.Div(decimal.NewFromInt(4))
// 		amount = xbnEquivalent.Truncate(7).String()
// 	}

// 	var txnParams basetxn.TransactionParams

// 	// var txnParams basetxn.TransactionParams
// 	if receiverStatus {
// 		creatAccountOpRequest := basetxn.CreateAccount{
// 			Destination: receiverAccount.Address,
// 			Amount:      amount,
// 		}

// 		txnParams = basetxn.TransactionParams{
// 			SourceAccount:        &sourceAccount,
// 			IncrementSequenceNum: true,
// 			Operations:           []basetxn.Operation{&creatAccountOpRequest},
// 			BaseFee:              1000,
// 			Memo:                 txnbuild.MemoText(memo),
// 			Timebounds:           txnbuild.NewTimeout(30),
// 		}
// 	} else {

// 		paymentOpRequest := basetxn.Payment{
// 			Destination: receiverAccount.Address,
// 			Amount:      amount,
// 			Asset:       asset,
// 		}

// 		txnParams = basetxn.TransactionParams{
// 			SourceAccount:        &sourceAccount,
// 			IncrementSequenceNum: true,
// 			Operations:           []basetxn.Operation{&paymentOpRequest},
// 			BaseFee:              1000,
// 			Memo:                 txnbuild.MemoText(memo),
// 			Timebounds:           txnbuild.NewTimeout(30),
// 		}
// 	}
// 	tx, err := basetxn.NewTransaction(txnParams)
// 	if err != nil {
// 		log.Println("[DoFaucetBlockchainPayment]error building transaction:", err)

// 		return nil, &bantuErrors.ErrorTemporaryServerError{}
// 	}
// 	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), kp)

// 	if err != nil {
// 		log.Println("[DoFaucetBlockchainPayment]error signing tramsaction:", err)
// 		return nil, &bantuErrors.ErrorTemporaryServerError{}
// 	}

// 	resp, err := client.SubmitTransaction(tx)
// 	if err != nil {
// 		log.Println("[DoFaucetBlockchainPayment]error submitting PayProfileCompletionReward txn:", err)

// 		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
// 			log.Print("[DoFaucetBlockchainPayment]error connecting to expansion service: ", err)

// 			return nil, &bantuErrors.ErrorTemporaryServerError{}
// 		}

// 		if hError, ok := err.(*horizonclient.Error); ok {
// 			//something went wrong, verify stage and check approprate action
// 			rCode, _ := hError.ResultCodes()
// 			rS, _ := hError.ResultString()
// 			log.Println("[DoFaucetBlockchainPayment] Problem in Transaction:", hError.Problem)
// 			log.Println("[DoFaucetBlockchainPayment] Result Codes in Transaction:", rCode)
// 			log.Println("[DoFaucetBlockchainPayment] Result String in Transaction:", rS)
// 			log.Printf("[DoFaucetBlockchainPayment] Problem in Transaction - RESPONSE: %+v\n", hError.Response)
// 			log.Println("[DoFaucetBlockchainPayment] Error submitting transaction:", err)
// 		}

// 		return nil, &bantuErrors.ErrorTemporaryServerError{}
// 	}
// 	paymentInfo.TransactionID = resp.Hash

// 	returnedPaymentInfo = &paymentInfo

// 	return returnedPaymentInfo, nil

// }

func AlertFaucetLowBalance(faucetKP *evmkeypair.Full) {

	//check faucet balance

	faucetPK := faucetKP.Address()
	faucetSecret := faucetKP.Seed()
	client := network.GetBlockchainClient()
	_, _, nativeBalance, _, _, err := network.BlockchainAccountProperties(client, faucetPK, basetxn.NativeAsset{})
	if err != nil {
		return
	}

	if faucetSecret == os.Getenv("XBN_FAUCET") {
		xbnFaucetMin := "200000"
		if os.Getenv("XBN_FAUCET_MIN_BALANCE") != "" && os.Getenv("XBN_FAUCET_MIN_BALANCE") != "0" {
			xbnFaucetMin = os.Getenv("XBN_FAUCET_MIN_BALANCE")
		}
		if nativeBalance.LessThan(decimal.RequireFromString(xbnFaucetMin)) {
			LogDiscordFaucetLowBalance(fmt.Sprintf("XBN FAUCET: %v has gone below minimum  warning amount %v. the balance is: %v", faucetPK, xbnFaucetMin, nativeBalance.String()))
		}
	}

	if faucetSecret == os.Getenv("REWARD_FAUCET") {
		rewardFaucetXBNMin := "20000"
		if os.Getenv("REWARD_FAUCET_XBN_MIN_BALANCE") != "" && os.Getenv("REWARD_FAUCET_XBN_MIN_BALANCE") != "0" {
			rewardFaucetXBNMin = os.Getenv("REWARD_FAUCET_XBN_MIN_BALANCE")
		}
		if nativeBalance.LessThan(decimal.RequireFromString(rewardFaucetXBNMin)) {
			LogDiscordFaucetLowBalance(fmt.Sprintf("REWARD FAUCET: %v has low XBN minimum balance %v. the balance is: %v", faucetPK, rewardFaucetXBNMin, nativeBalance.String()))
		}

		if os.Getenv("REWARD_ASSET_CODE") != "" && os.Getenv("REWARD_ASSET_ISSUER") != "" {
			rewardAsset := basetxn.CreditAsset{Code: os.Getenv("REWARD_ASSET_CODE"), Issuer: os.Getenv("REWARD_ASSET_ISSUER")}
			_, _, _, rewardBalance, _, e := network.BlockchainAccountProperties(client, faucetPK, rewardAsset)
			if e == nil {
				rewardFaucetMin := "200000"
				if os.Getenv("REWARD_FAUCET_MIN_BALANCE") != "" && os.Getenv("REWARD_FAUCET_MIN_BALANCE") != "0" {
					rewardFaucetMin = os.Getenv("REWARD_FAUCET_MIN_BALANCE")
				}
				if rewardBalance.LessThan(decimal.RequireFromString(rewardFaucetMin)) {
					LogDiscordFaucetLowBalance(fmt.Sprintf("REWARD FAUCET: %v has gone below minimum warning amount %v %v. the balance is: %v %v", faucetPK, rewardFaucetMin, os.Getenv("REWARD_ASSET_CODE"), rewardBalance.String(), os.Getenv("REWARD_ASSET_CODE")))
				}
			}
		}
	}
}

// //DoFaucetPaymentWithChannelAccount makes payment from faucet with channel account on server
// func DoFaucetPaymentWithChannelAccount(faucetSecret, receiver, assetCode, assetIssuer, amount, memo string, db *gorm.DB) (returnedPaymentInfo *paymentModels.PaymentInfo, err error) {
// 	var paymentInfo paymentModels.PaymentInfo
// 	kp := evmkeypair.MustParseFull(faucetSecret)

// 	faucetPK := kp.Address()
// 	paymentInfo = paymentModels.PaymentInfo{
// 		Destination: receiver, Memo: memo, AssetIssuer: assetIssuer,
// 		AssetCode: assetCode, Amount: amount,
// 	}
// 	AlertFaucetLowBalance(kp)
// 	FundChannelAccount()
// 	retPaymentInfo, err := PayWithChannelAccount(faucetPK, &paymentInfo, db)
// 	if err != nil {
// 		log.Printf("[DoFaucetPayment] ###### payment transaction request failed:%v\n%+v\n", err, retPaymentInfo)
// 		return
// 	}
// 	// log.Printf("########..... success: ReturnedPayment Transaction to sign: %+v\n", *retPaymentInfo)
// 	// log.Printf("########..... success: PaymentInfo Transaction to sign: %+v\n", paymentInfo)

// 	signedBase64, err := bantupaysdk.SignBase64Txn(kp.Seed(), paymentInfo.Transaction, paymentInfo.NetworkPassPhrase)
// 	if err != nil {
// 		log.Println("[DoFaucetPayment] Transaction signing failed:", err)
// 		return nil, err
// 	}

// 	paymentInfo.TransactionSignature = signedBase64

// 	returnedPaymentInfo, err = PayWithChannelAccount(faucetPK, &paymentInfo, db)
// 	if err != nil {
// 		log.Println("[DoFaucetPayment] error making Payment:", err)
// 		return nil, err
// 	}
// 	return

// }

// //DoFaucetPayment makes payment from faucet on server
// func DoFaucetPayment(faucetSecret, receiver, assetCode, assetIssuer, amount, memo string, db *gorm.DB) (returnedPaymentInfo *paymentModels.PaymentInfo, err error) {
// 	var paymentInfo paymentModels.PaymentInfo
// 	kp := evmkeypair.MustParseFull(faucetSecret)

// 	faucetPK := kp.Address()
// 	paymentInfo = paymentModels.PaymentInfo{
// 		Destination: receiver, Memo: memo, AssetIssuer: assetIssuer,
// 		AssetCode: assetCode, Amount: amount,
// 	}

// 	alertFaucetLowBalance(kp)

// 	retPaymentInfo, _, err := Pay(faucetPK, &paymentInfo, db)
// 	if err != nil {
// 		log.Printf("[DoFaucetPayment] ###### payment transaction request failed:%v\n%+v\n", err, retPaymentInfo)
// 		return
// 	}
// 	// log.Printf("########..... success: ReturnedPayment Transaction to sign: %+v\n", *retPaymentInfo)
// 	// log.Printf("########..... success: PaymentInfo Transaction to sign: %+v\n", paymentInfo)

// 	signedBase64, err := bantupaysdk.SignBase64Txn(kp.Seed(), paymentInfo.Transaction, paymentInfo.NetworkPassPhrase)
// 	if err != nil {
// 		log.Println("[DoFaucetPayment] Transaction signing failed:", err)
// 		return nil, err
// 	}
// 	paymentInfo.TransactionSignature = signedBase64

// 	returnedPaymentInfo, _, err = Pay(faucetPK, &paymentInfo, db)
// 	if err != nil {
// 		log.Println("[DoFaucetPayment] error making Payment:", err)
// 		return nil, err
// 	}
// 	return

// }

func LogDiscordFaucetLowBalance(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/828199026263982120/phflpVArKrTPCT1wfqzc5QNu4jlmWI3dc3n5xGpy1dV7cIHGul_V9IC1tnwggdUa2GvQ"
	if len(os.Getenv("FAUCET_LOW_BALANCE_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("FAUCET_LOW_BALANCE_WEBHOOK")
	}
	discord.Say(msg)
}
