package users

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	// bc "trovo-wallet-api/internal/blockchainalgofuncs"
	"trovo-wallet-api/internal/basetxn"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

func MakeOffer(signerUser, walletOwner *userModels.User, sourceWallet *userModels.UserWallet, offerRequest *userModels.MarketOfferRequest, gc *sharedconfig.GlobalConfig) (err error) {
	offerRequest.Messages = make([]string, 0)

	if sourceWallet.SharedAccessEnabled == 1 && sourceWallet.NumberOfApprovalsNeeded > 0 {
		offerRequest.Multiparty = 1
	}
	if sourceWallet.HasViewOnlyAccess(gc) {
		offerRequest.SignatureRequired = 1
	}
	// mmWallet, err := walletOwner.GetMartketMakingWallet(gc.DB)
	// if err != nil {
	// 	log.Printf("[MakeOffer]Error validating Market making wallet: %v", err)
	// 	return err
	// }
	if !strings.EqualFold(offerRequest.OfferType, "BUY") && !strings.EqualFold(offerRequest.OfferType, "SELL") {
		return &tErrors.CustomError{
			Param:      "offerType",
			Err:        "error-invalid-paramter",
			ErrMessage: "Offer type can only be either BUY or SELL.",
		}
	}
	if len(offerRequest.Quantity) == 0 {
		return &tErrors.CustomError{
			Param:      "quantity",
			Err:        "error-missing-quantity",
			ErrMessage: "Missing quantity.",
		}
	}
	if !decimal.RequireFromString(offerRequest.Quantity).IsPositive() {
		return &tErrors.CustomError{
			Param:      "quantity",
			Err:        "error-missing-quantity",
			ErrMessage: "Quantity must be greater than zero.",
		}
	}
	if len(offerRequest.PricePerUnit) == 0 {
		return &tErrors.CustomError{
			Param:      "pricePerUnit",
			Err:        "error-missing-price",
			ErrMessage: "Missing price.",
		}
	}

	if !decimal.RequireFromString(offerRequest.PricePerUnit).IsPositive() {
		return &tErrors.CustomError{
			Param:      "pricePerUnit",
			Err:        "error-missing-price",
			ErrMessage: "Price must be greater than zero.",
		}
	}
	// mmSignerKeyPair, err := bc.MarketMakingSignerKeypair(walletOwner.Username, mmWallet.ID)
	// if err != nil {

	// 	return &tErrors.ErrorTemporaryServerError{}
	// }
	assetOfMarket := os.Getenv("NATIVE_ASSET_CODE")
	if len(offerRequest.AssetIssuer) == 56 {
		assetOfMarket = fmt.Sprintf("%v:%v...%v", offerRequest.AssetCode, offerRequest.AssetIssuer[0:4], offerRequest.AssetIssuer[51:55])
	}
	currencyOfMarket := os.Getenv("NATIVE_ASSET_CODE")
	if len(offerRequest.CurrencyIssuer) == 56 {
		currencyOfMarket = fmt.Sprintf("%v:%v...%v", offerRequest.CurrencyCode, offerRequest.CurrencyIssuer[0:4], offerRequest.CurrencyIssuer[51:55])
	}

	if strings.EqualFold(offerRequest.AssetCode, os.Getenv("NATIVE_ASSET_CODE")) {
		offerRequest.AssetCode = ""
		offerRequest.AssetIssuer = ""
	}

	if strings.EqualFold(offerRequest.CurrencyCode, os.Getenv("NATIVE_ASSET_CODE")) {
		offerRequest.CurrencyCode = ""
		offerRequest.CurrencyIssuer = ""
	}
	serviceFee, e := decimal.NewFromString(os.Getenv("MARKET_MAKING_FEE_AMOUNT"))
	if e != nil || os.Getenv("MARKET_MAKING_FEE_ENABLED") == "0" {
		serviceFee = decimal.Zero
	}
	totalQty := decimal.RequireFromString(offerRequest.Quantity).Truncate(7)
	// feeQuantity := (totalQty.Mul(serviceFee.Div(decimal.NewFromInt(100)))).Truncate(7)
	feeQuantity := serviceFee.Div(decimal.NewFromInt(0)) //set to zero since fees r now removed.
	netQuantity := totalQty.Sub(feeQuantity)
	offerRequest.NetQuantity = netQuantity.String()
	offerRequest.FeeChargedOnAsset = serviceFee.String()
	offerRequest.FeeValue = feeQuantity.String()

	var marketOffer *userModels.MarketOffer
	xdrBase64, err := generateMakeMarketXdr(sourceWallet, offerRequest, gc)
	if err != nil {
		return err
	}

	//do not overwrite transaction
	if len(offerRequest.Transaction) == 0 {

		offerRequest.Transaction = xdrBase64
	}

	offerRequest.NetworkPassPhrase = gc.BantuNetworkPassphrase

	oldTxn := offerRequest.Transaction

	if len(offerRequest.TransactionSignature) == 0 && offerRequest.Commit == 0 {
		//no signature
		return nil

	}

	//there was a signature... let's submit

	if oldTxn != xdrBase64 && offerRequest.Commit == 0 {
		return &tErrors.CustomError{
			Param:      "transaction",
			Err:        "transaction mismatch",
			ErrMessage: "transaction mismatch, please try again",
			Code:       404,
		}
	}
	var assetIssuer, currencyIssuer *string
	if len(offerRequest.AssetIssuer) > 0 {
		assetIssuer = &offerRequest.AssetIssuer
	}
	if len(offerRequest.CurrencyIssuer) > 0 {
		currencyIssuer = &offerRequest.CurrencyIssuer
	}
	{
		marketOffer = &userModels.MarketOffer{
			ID:                          uuid.NewString(),
			SourceWalletAlias:           sourceWallet.Alias,
			SourceWalletPublicKey:       sourceWallet.ID,
			MarketMakingWalletPublicKey: sourceWallet.ID,
			OfferType:                   offerRequest.OfferType,
			AssetCode:                   offerRequest.AssetCode,
			AssetIssuer:                 assetIssuer,
			CurrencyCode:                offerRequest.CurrencyCode,
			CurrencyIssuer:              currencyIssuer,
			PricePerUnit:                offerRequest.PricePerUnit,
			Quantity:                    offerRequest.Quantity,
			FeeChargedOnAsset:           offerRequest.FeeChargedOnAsset,
			FeeValue:                    offerRequest.FeeValue,
			NetQuantity:                 offerRequest.NetQuantity,
			// RemainingQuantity:           offerRequest.NetQuantity,
			// RemainingFeeValue:           offerRequest.FeeValue,
		}
	}
	var msgs string
	for i, m := range offerRequest.Messages {
		msgs = m
		if i < len(offerRequest.Messages)-1 {
			msgs = fmt.Sprintf("%s\n", msgs)
		}
	}
	fp := offerRequest.FeeChargedOnAsset + "%"
	offerRequest.ReturnedDescription = fmt.Sprintf("%v %v %v @ %v %v with %v fee of %v %v", offerRequest.OfferType, offerRequest.NetQuantity, assetOfMarket, offerRequest.PricePerUnit, currencyOfMarket, fp, offerRequest.FeeValue, assetOfMarket)
	if len(msgs) > 0 {
		offerRequest.ReturnedDescription = fmt.Sprintf("%s\nMessages: %v", offerRequest.ReturnedDescription, msgs)
	}
	if (offerRequest.Commit == 0 && sourceWallet.SharedAccessEnabled == 1 && sourceWallet.NumberOfApprovalsNeeded == 0) || (sourceWallet.SharedAccessEnabled == 0 && offerRequest.Commit == 0) {
		dbTX := gc.DB.Begin()
		defer dbTX.Rollback()
		e = dbTX.Omit(clause.Associations).Create(&marketOffer).Error
		if e != nil {
			log.Printf("[MakeOffer]Error creating market offer: %+v\nError: %v\n", marketOffer, e)
			return &tErrors.ErrorTemporaryServerError{}
		}
		txnResult, err := network.SubmitXdrWithSignatureReturnsTrx(gc.BantuExpansionClient, sourceWallet.Signer, xdrBase64, offerRequest.TransactionSignature)

		if err != nil {
			log.Printf("[MakeOffer]error submitting txn: %v\n", err)
			return &tErrors.ErrorTemporaryServerError{}
		}

		offerRequest.TransactionID = txnResult.Hash
		marketOffer.TransactionID = &txnResult.Hash
		// Market-making DEX offers have no Base equivalent (see basetxn.ManageSellOffer's
		// doc) - there is no on-chain offer ID to extract from the submitted
		// transaction's result, so marketOffer.BlockchainOfferID stays unset
		// pending a real Base market-making design.
		e = dbTX.Omit(clause.Associations).Save(&marketOffer).Error
		if e != nil {
			log.Printf("[MakeOffer]Error saving transactionID on market offer: %+v\nError: %v\n", marketOffer, e)
			return &tErrors.ErrorTemporaryServerError{}
		}
		dbTX.Commit()
		walletOwner.InvalidateUserCache(gc)
		return nil
	}

	if offerRequest.Multiparty == 1 {
		offerRequest.TransactionID = "PENDING_AUTH"
		log.Printf("[MakeOffer]shared access with approver permission enabled for %v \n", sourceWallet.Alias)
		id := uuid.NewString()

		description := offerRequest.ReturnedDescription
		//use market offer to be inserted into the table after authorization, instead of constructing again
		transactionByte, _ := json.Marshal(*marketOffer)
		transactionStr := string(transactionByte)
		pendingAuth := userModels.PendingAuth{
			ID:                       id,
			Initiator:                signerUser.Username,
			InitiatorSignerPublicKey: signerUser.PrimarySigner,
			WalletPublicKey:          sourceWallet.ID,
			TransactionType:          "MAKE MARKET OFFER",
			Description:              description,
			TransactionSource:        offerRequest.TransactionSource,
			ApprovalsNeeded:          sourceWallet.NumberOfApprovalsNeeded,
			TransactionXdr:           xdrBase64,
			TransactionInfoStr:       &transactionStr,
		}
		//save and commit this to database
		e := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[MakeOffer] Error saving MAKE MARKET OFFER txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}
			return err
		}
		return nil

	}
	log.Println("[MakeOffer]unknown conditions for makr market")
	return &tErrors.ErrorTemporaryServerError{}
}

func CancelOffer(signerUser, walletOwner *userModels.User, sourceWallet *userModels.UserWallet, deleteOfferRequest *userModels.DeleteOfferRequest, gc *sharedconfig.GlobalConfig) (err error) {

	if len(deleteOfferRequest.ID) == 0 {
		return &tErrors.ErrorMissingParameter{Parameter: "Id"}
	}

	marketOffer, err := sourceWallet.GetMarketOfferByID(deleteOfferRequest.ID, gc.DB, gc)
	if err != nil {
		return err
	}

	if marketOffer.Canceled == 1 {
		return &tErrors.CustomError{
			Param:      "id",
			Err:        "error-offer-has been canceled",
			ErrMessage: "Offer cannot be canceled because it has been already been canceled before.",
		}
	}

	// if decimal.RequireFromString(marketOffer.RemainingQuantity).IsZero() {
	// 	return &tErrors.CustomError{
	// 		Param:      "id",
	// 		Err:        "error-offer-has been filled",
	// 		ErrMessage: "Offer cannot be canceled because it has been filled.",
	// 	}
	// }
	bOffer, err := marketOffer.GetBlockchainOfferDetail(gc)
	if err != nil {
		return err
	}
	if decimal.RequireFromString(bOffer.Amount).IsZero() {
		return &tErrors.CustomError{
			Param:      "id",
			Err:        "error-offer-has been filled",
			ErrMessage: "Offer cannot be canceled because it has been filled.",
		}
	}

	// mmWallet, err := walletOwner.GetMartketMakingWallet(gc.DB)
	// if err != nil {
	// 	log.Printf("[MakeOffer]Error validating Market making wallet: %v", err)
	// 	return err
	// }

	// mmSignerKeyPair, err := bc.MarketMakingSignerKeypair(walletOwner.Username, mmWallet.ID)
	// if err != nil {

	// 	return &tErrors.ErrorTemporaryServerError{}
	// }

	xdrBase64, transactionSource, err := generateDeleteMarketXdr(sourceWallet, &marketOffer, bOffer, gc)
	if err != nil {
		return err
	}
	if sourceWallet.HasViewOnlyAccess(gc) {
		deleteOfferRequest.SignatureRequired = 1
	}
	deleteOfferRequest.TransactionSource = transactionSource
	deleteOfferRequest.Transaction = xdrBase64
	deleteOfferRequest.NetworkPassPhrase = gc.BantuNetworkPassphrase

	oldTxn := deleteOfferRequest.Transaction

	if len(deleteOfferRequest.TransactionSignature) == 0 && deleteOfferRequest.Commit == 0 {
		//no signature
		return nil

	}

	//there was a signature... let's submit

	if oldTxn != xdrBase64 && deleteOfferRequest.Commit == 0 {
		return &tErrors.CustomError{
			Param:      "transaction",
			Err:        "transaction mismatch",
			ErrMessage: "transaction mismatch, please try again",
			Code:       404,
		}
	}

	// deleteOfferRequest.ReturnedDescription = fmt.Sprintf("%v %v %v @ %v %v with %v fee of %v %v", offerRequest.OfferType, offerRequest.NetQuantity, assetOfMarket, offerRequest.PricePerUnit, currencyOfMarket, fp, offerRequest.FeeValue, assetOfMarket)
	// if len(msgs) > 0 {
	// 	deleteOfferRequest.ReturnedDescription = fmt.Sprintf("%s\nMessages: %v", deleteOfferRequest.ReturnedDescription, msgs)
	// }
	if (deleteOfferRequest.Commit == 0 && sourceWallet.SharedAccessEnabled == 1 && sourceWallet.NumberOfApprovalsNeeded == 0) || (sourceWallet.SharedAccessEnabled == 0 && deleteOfferRequest.Commit == 0) {
		dbTX := gc.DB.Begin()
		defer dbTX.Rollback()
		e := dbTX.Omit(clause.Associations).Create(&marketOffer).Error
		if e != nil {
			log.Printf("[DeleteOffer]Error saving market offer: %+v\nError: %v\n", marketOffer, err)
			return &tErrors.ErrorTemporaryServerError{}
		}
		txnResult, err := network.SubmitXdrWithSignatureReturnsTrx(gc.BantuExpansionClient, sourceWallet.Signer, xdrBase64, deleteOfferRequest.TransactionSignature)

		if err != nil {
			log.Printf("[MakeOffer]error submitting txn: %v\n", err)
			return &tErrors.ErrorTemporaryServerError{}
		}

		deleteOfferRequest.TransactionID = txnResult.Hash
		marketOffer.TransactionID = &txnResult.Hash
		// Market-making DEX offers have no Base equivalent (see
		// basetxn.ManageSellOffer's doc) - marketOffer.BlockchainOfferID
		// stays unset pending a real Base market-making design.
		e = dbTX.Omit(clause.Associations).Save(&marketOffer).Error
		if e != nil {
			log.Printf("[MakeOffer]Error saving transactionID on market offer: %+v\nError: %v\n", marketOffer, err)
			return &tErrors.ErrorTemporaryServerError{}
		}
		dbTX.Commit()
		walletOwner.InvalidateUserCache(gc)
		return nil
	}

	if deleteOfferRequest.Multiparty == 1 {
		deleteOfferRequest.TransactionID = "PENDING_AUTH"
		log.Printf("[MakeOffer]shared access with approver permission enabled for %v \n", sourceWallet.Alias)
		id := uuid.NewString()

		description := deleteOfferRequest.ReturnedDescription
		//use market offer to be inserted into the table after authorization, instead of constructing again
		transactionByte, _ := json.Marshal(*deleteOfferRequest)
		transactionStr := string(transactionByte)
		pendingAuth := userModels.PendingAuth{
			ID:                       id,
			Initiator:                signerUser.Username,
			InitiatorSignerPublicKey: signerUser.PrimarySigner,
			WalletPublicKey:          sourceWallet.ID,
			TransactionType:          "DELETE MARKET OFFER",
			Description:              description,
			TransactionSource:        deleteOfferRequest.TransactionSource,
			ApprovalsNeeded:          sourceWallet.NumberOfApprovalsNeeded,
			TransactionXdr:           xdrBase64,
			TransactionInfoStr:       &transactionStr,
		}
		//save and commit this to database
		e := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[MakeOffer] Error saving MAKE MARKET OFFER txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}
			return err
		}
		return nil

	}

	return &tErrors.ErrorTemporaryServerError{}

}

func generateMakeMarketXdr(sourceWallet *userModels.UserWallet, offerRequest *userModels.MarketOfferRequest, gc *sharedconfig.GlobalConfig) (txnBase64 string, err error) {
	// offerFeePercentage := offerRequest.FeeChargedOnAsset + "%"
	minBalance := decimal.RequireFromString(os.Getenv("STANDARD_WALLET_MINIMUM_BALANCE"))
	offerRequest.Messages = make([]string, 0)
	var memo string
	var currencyAsset, mainAsset basetxn.Asset
	if offerRequest.AssetCode+offerRequest.AssetIssuer == offerRequest.CurrencyCode+offerRequest.CurrencyIssuer {
		return "", &tErrors.CustomError{
			Param:      "AssetCode",
			Err:        "error-asset-and-currency-are-the-same",
			ErrMessage: "You cannot make a market against same asset",
		}
	}
	offerRequest.OfferType = strings.ToUpper(offerRequest.OfferType)
	n, d := ToFractionInt32(decimal.RequireFromString(offerRequest.PricePerUnit).InexactFloat64())

	if offerRequest.AssetIssuer == "" {
		mainAsset = basetxn.NativeAsset{}
	} else {
		mainAsset = basetxn.CreditAsset{Code: offerRequest.AssetCode, Issuer: offerRequest.AssetIssuer}
	}

	if offerRequest.CurrencyIssuer == "" {
		currencyAsset = basetxn.NativeAsset{}
	} else {
		currencyAsset = basetxn.CreditAsset{Code: offerRequest.CurrencyCode, Issuer: offerRequest.CurrencyIssuer}
	}

	chanAccount := <-gc.ChannelAccounts
	defer func(c *evmkeypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)

	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, chanAccount.Address(), basetxn.NativeAsset{})
	var sourceAccount *network.AccountInfo
	var sourceMAccountExists, sourceMAccountTrustsAsset bool
	var nativeMAccountBalance decimal.Decimal
	if strings.EqualFold(offerRequest.OfferType, "BUY") {
		sourceMAccountExists, sourceMAccountTrustsAsset, nativeMAccountBalance, _, sourceAccount, _ = network.BlockchainAccountProperties(gc.BantuExpansionClient, sourceWallet.ID, currencyAsset)

		if !sourceMAccountExists {
			return "", &tErrors.CustomError{Param: "publicKey", Err: "error-account-not-activated-on-blockchain", ErrMessage: "The Wallet public key is currently underfunded. Please send about 3XBN to it to activate it before you can perform this task", Code: http.StatusBadRequest}

		}
		if nativeMAccountBalance.LessThan(minBalance) {
			return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", minBalance.String(), os.Getenv("NATIVE_ASSET_CODE")), Code: http.StatusBadRequest}

		}

		if !mainAsset.IsNative() && !sourceMAccountTrustsAsset {

			return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", offerRequest.Quantity, currencyAsset.GetCode()), Code: http.StatusBadRequest}

		}
	}
	if strings.EqualFold(offerRequest.OfferType, "SELL") {
		sourceMAccountExists, sourceMAccountTrustsAsset, nativeMAccountBalance, _, sourceAccount, _ = network.BlockchainAccountProperties(gc.BantuExpansionClient, sourceWallet.ID, mainAsset)

		if !sourceMAccountExists {
			return "", &tErrors.CustomError{Param: "publicKey", Err: "error-account-not-activated-on-blockchain", ErrMessage: "The Wallet public key is currently underfunded. Please send about 3XBN to it to activate it before you can perform this task", Code: http.StatusBadRequest}

		}
		if nativeMAccountBalance.LessThan(minBalance) {
			return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", minBalance.String(), os.Getenv("NATIVE_ASSET_CODE")), Code: http.StatusBadRequest}

		}

		if !mainAsset.IsNative() && !sourceMAccountTrustsAsset {

			return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", offerRequest.Quantity, mainAsset.GetCode()), Code: http.StatusBadRequest}

		}
	}

	var ops []basetxn.Operation = make([]basetxn.Operation, 0)

	if strings.EqualFold(offerRequest.OfferType, "BUY") {
		// if !currencyAsset.IsNative() {
		// 	//check if it has trustline to it and then create it.
		// 	_, _, _, _, _, _ = network.BlockchainAccountProperties(gc.BantuExpansionClient, sourceWallet.ID, currencyAsset)
		// 	_, mmAccountTrustsAsset, mmnativeAccountBalance, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, sourceAccount.ID, currencyAsset)
		// 	if mmnativeAccountBalance.LessThan(minBalance) {
		// 		ops = append(ops, &basetxn.Payment{
		// 			Asset:         basetxn.NativeAsset{},
		// 			Destination:   mmWallet.ID,
		// 			Amount:        minBalance.Mul(decimal.NewFromInt(3)).String(),
		// 			SourceAccount: sourceWallet.ID,
		// 		})
		// 	}

		// 	if !mmAccountTrustsAsset {
		// 		//establish trustline automatically
		// 		ops = append(ops, &basetxn.ChangeTrust{
		// 			Line:          txnbuild.ChangeTrustAssetWrapper{Asset: currencyAsset},
		// 			Limit:         "900000000000",
		// 			SourceAccount: mmWallet.ID,
		// 		})
		// 	}

		// }
		// {
		// 	_, _, _, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, sourceAccount.ID, mainAsset)
		// 	if !mmAccountTrustsAsset {
		// 		//establish trustline automatically
		// 		ops = append(ops, &basetxn.ChangeTrust{
		// 			Line:          txnbuild.ChangeTrustAssetWrapper{Asset: mainAsset},
		// 			Limit:         "900000000000",
		// 			SourceAccount: mmWallet.ID,
		// 		})
		// 	}
		// }
		// // move sellinng funds to the MM wallet
		// ops = append(ops, &basetxn.Payment{
		// 	Asset:         currencyAsset,
		// 	Destination:   sourceWallet.ID,
		// 	Amount:        offerRequest.Quantity,
		// 	SourceAccount: sourceWallet.ID,
		// })

		//invert the price fraction
		ops = append(ops, &basetxn.ManageSellOffer{
			Selling:       currencyAsset,
			Buying:        mainAsset,
			Amount:        offerRequest.NetQuantity,
			Price:         fmt.Sprintf("%v/%v", n, d),
			SourceAccount: sourceWallet.ID,
		})
		bcode, scode := mainAsset.GetCode(), currencyAsset.GetCode()
		if mainAsset.IsNative() {
			bcode = os.Getenv("NATIVE_ASSET_CODE")
		}
		if currencyAsset.IsNative() {
			scode = os.Getenv("NATIVE_ASSET_CODE")
		}
		memo = fmt.Sprintf("s%v-b%v", scode, bcode)

	}
	if strings.EqualFold(offerRequest.OfferType, "SELL") {
		// if !mainAsset.IsNative() {
		// 	//check if it has trustline to it and then create it.
		// 	_, mmAccountTrustsAsset, mmnativeAccountBalance, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, mmWallet.ID, mainAsset)
		// 	if mmnativeAccountBalance.LessThan(minBalance) {
		// 		ops = append(ops, &basetxn.Payment{
		// 			Asset:         basetxn.NativeAsset{},
		// 			Destination:   mmWallet.ID,
		// 			Amount:        minBalance.Mul(decimal.NewFromInt(3)).String(),
		// 			SourceAccount: sourceWallet.ID,
		// 		})
		// 	}

		// 	if !mmAccountTrustsAsset {
		// 		//establish trustline automatically
		// 		ops = append(ops, &basetxn.ChangeTrust{
		// 			Line:          txnbuild.ChangeTrustAssetWrapper{Asset: mainAsset},
		// 			Limit:         "900000000000",
		// 			SourceAccount: mmWallet.ID,
		// 		})
		// 	}

		// }
		// {
		// 	_, mmAccountTrustsAsset, _, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, mmWallet.ID, currencyAsset)
		// 	if !mmAccountTrustsAsset {
		// 		//establish trustline automatically
		// 		ops = append(ops, &basetxn.ChangeTrust{
		// 			Line:          txnbuild.ChangeTrustAssetWrapper{Asset: currencyAsset},
		// 			Limit:         "900000000000",
		// 			SourceAccount: mmWallet.ID,
		// 		})
		// 	}
		// }
		// // move sellinng funds to the MM wallet
		// ops = append(ops, &basetxn.Payment{
		// 	Asset:         mainAsset,
		// 	Destination:   mmWallet.ID,
		// 	Amount:        offerRequest.Quantity,
		// 	SourceAccount: sourceWallet.ID,
		// })

		// make offer with net quantity so that fee can be returned.
		ops = append(ops, &basetxn.ManageSellOffer{
			Selling:       mainAsset,
			Buying:        currencyAsset,
			Amount:        offerRequest.NetQuantity,
			Price:         fmt.Sprintf("%v/%v", n, d),
			SourceAccount: sourceWallet.ID,
		})
		scode, bcode := mainAsset.GetCode(), currencyAsset.GetCode()
		if mainAsset.IsNative() {
			scode = os.Getenv("NATIVE_ASSET_CODE")
		}
		if currencyAsset.IsNative() {
			bcode = os.Getenv("NATIVE_ASSET_CODE")
		}
		memo = fmt.Sprintf("s%v-b%v", scode, bcode)
	}

	//service fee
	// signForFeeTrustLine := 0
	// if decimal.RequireFromString(offerRequest.FeeValue).IsPositive() && os.Getenv("MARKET_MAKING_FEE_ENABLED") == "1" {

	// 	//process service fee

	// 	//no need deducting it as we deduct it as market executes
	// 	feeAssetCode := mainAsset.GetCode()
	// 	if mainAsset.IsNative() {
	// 		feeAssetCode = os.Getenv("NATIVE_ASSET_CODE")
	// 	}

	// 	if !mainAsset.IsNative() {
	// 		//ensure that the fee address is can accept the asset.
	// 		// but bcos  fee address needs to sign, it cannot be done here
	// 		mmfeeKeypair := evmkeypair.MustParseFull(os.Getenv("MARKET_MAKING_FEE_WALLET"))
	// 		mmFeeAddress := mmfeeKeypair.Address()

	// 		{
	// 			_, feeAccountTrustsAsset, _, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, mmFeeAddress, mainAsset)
	// 			if !feeAccountTrustsAsset {
	// 				signForFeeTrustLine = 1
	// 				//establish trustline automatically
	// 				ops = append(ops, &basetxn.ChangeTrust{
	// 					Line:          txnbuild.ChangeTrustAssetWrapper{Asset: mainAsset},
	// 					Limit:         "900000000000",
	// 					SourceAccount: mmFeeAddress,
	// 				})

	// 				//throw error
	// 				// return "", &tErrors.CustomError{Param: "publicKey", Err: "error-asset-not-configured-for-fee-address", ErrMessage: fmt.Sprintf("Please contact support to configure %v fee for before you can perform this task.", mainAsset.GetCode()), Code: http.StatusBadRequest}

	// 			}
	// 		}
	// 		feeAssetCode = mainAsset.GetCode()
	// 	}
	// 	offerRequest.Messages = append(offerRequest.Messages, fmt.Sprintf("%v %v (%v) will be deducted from the total quantity as service fee and your offer will be placed with %v %v.", offerRequest.FeeValue, feeAssetCode, offerFeePercentage, offerRequest.NetQuantity, feeAssetCode))

	// }

	// Construct the transaction that holds the operations to execute on the network

	var tx *basetxn.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if offerRequest.Multiparty == 1 {

		offerRequest.TransactionSource = chanSourceAccount.Address

		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        chanSourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 memo,
			},
		)
	} else {
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        sourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 memo,
			},
		)
	}
	if err != nil {
		log.Println("[generateMakeMarketXdr]error constructing transaction", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	// tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), mmSignerKeyPair)

	// if err != nil {
	// 	log.Println("[generateMakeMarketXdr] error signing transaction with custodial signer key ", err)
	// 	return "", &tErrors.ErrorTemporaryServerError{}
	// }

	if offerRequest.Multiparty == 1 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

		if err != nil {
			log.Println("[generateMakeMarketXdr] error signing transaction with channelAccount key ", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}

	// if signForFeeTrustLine == 1 {
	// 	mmfeeKeypair := evmkeypair.MustParseFull(os.Getenv("MARKET_MAKING_FEE_WALLET"))
	// 	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), mmfeeKeypair)

	// 	if err != nil {
	// 		log.Println("[generateMakeMarketXdr] error signing transaction with market making fee key ", err)
	// 		return "", &tErrors.ErrorTemporaryServerError{}
	// 	}
	// }

	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", err
	}
	offerRequest.Memo = memo
	return xdrBase64, nil

}
func generateDeleteMarketXdr(sourceWallet *userModels.UserWallet, offerRequest *userModels.MarketOffer, bOffer userModels.MarketOfferDetail, gc *sharedconfig.GlobalConfig) (txnBase64, transactionSource string, err error) {
	// offerFeePercentage := offerRequest.FeeChargedOnAsset + "%"
	minBalance := decimal.RequireFromString(os.Getenv("STANDARD_WALLET_MINIMUM_BALANCE"))
	// offerRequest.Messages = make([]string, 0)
	var memo string
	var currencyAsset, mainAsset basetxn.Asset

	offerRequest.OfferType = strings.ToUpper(offerRequest.OfferType)
	fraction := decimal.RequireFromString(offerRequest.PricePerUnit).Rat()
	d := int32(fraction.Denom().Int64())
	n := int32(fraction.Num().Int64())
	offerIDInt, _ := strconv.ParseInt(bOffer.ID, 10, 64)
	if offerRequest.AssetIssuer == nil {
		mainAsset = basetxn.NativeAsset{}
	} else {
		mainAsset = basetxn.CreditAsset{Code: offerRequest.AssetCode, Issuer: *offerRequest.AssetIssuer}
	}

	if offerRequest.CurrencyIssuer == nil {
		currencyAsset = basetxn.NativeAsset{}
	} else {
		currencyAsset = basetxn.CreditAsset{Code: offerRequest.CurrencyCode, Issuer: *offerRequest.CurrencyIssuer}
	}

	chanAccount := <-gc.ChannelAccounts
	defer func(c *evmkeypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)

	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, chanAccount.Address(), basetxn.NativeAsset{})
	var sourceAccount *network.AccountInfo
	var sourceMAccountExists, sourceMAccountTrustsAsset bool
	var nativeMAccountBalance decimal.Decimal
	if strings.EqualFold(offerRequest.OfferType, "BUY") {
		sourceMAccountExists, sourceMAccountTrustsAsset, nativeMAccountBalance, _, sourceAccount, _ = network.BlockchainAccountProperties(gc.BantuExpansionClient, sourceWallet.ID, currencyAsset)

		if !sourceMAccountExists {
			return "", "", &tErrors.CustomError{Param: "publicKey", Err: "error-account-not-activated-on-blockchain", ErrMessage: "The Wallet public key is currently underfunded. Please send about 3XBN to it to activate it before you can perform this task", Code: http.StatusBadRequest}

		}
		if nativeMAccountBalance.LessThan(minBalance) {
			return "", "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", minBalance.String(), os.Getenv("NATIVE_ASSET_CODE")), Code: http.StatusBadRequest}

		}

		if !mainAsset.IsNative() && !sourceMAccountTrustsAsset {

			return "", "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", offerRequest.Quantity, currencyAsset.GetCode()), Code: http.StatusBadRequest}

		}
	}
	if strings.EqualFold(offerRequest.OfferType, "SELL") {
		sourceMAccountExists, sourceMAccountTrustsAsset, nativeMAccountBalance, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, sourceWallet.ID, mainAsset)

		if !sourceMAccountExists {
			return "", "", &tErrors.CustomError{Param: "publicKey", Err: "error-account-not-activated-on-blockchain", ErrMessage: "The Wallet public key is currently underfunded. Please send about 3XBN to it to activate it before you can perform this task", Code: http.StatusBadRequest}

		}
		if nativeMAccountBalance.LessThan(minBalance) {
			return "", "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", minBalance.String(), os.Getenv("NATIVE_ASSET_CODE")), Code: http.StatusBadRequest}

		}

		if !mainAsset.IsNative() && !sourceMAccountTrustsAsset {

			return "", "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", offerRequest.Quantity, mainAsset.GetCode()), Code: http.StatusBadRequest}

		}
	}

	var ops []basetxn.Operation = make([]basetxn.Operation, 0)

	if strings.EqualFold(offerRequest.OfferType, "BUY") {

		//invert the price fraction
		ops = append(ops, &basetxn.ManageSellOffer{
			OfferID:       offerIDInt,
			Selling:       currencyAsset,
			Buying:        mainAsset,
			Amount:        "0",
			Price:         fmt.Sprintf("%v/%v", n, d),
			SourceAccount: sourceWallet.ID,
		})

		// // move selling funds from the MM wallet to the wallet it was made from
		// ops = append(ops, &basetxn.Payment{
		// 	Asset:         currencyAsset,
		// 	Destination:   sourceWallet.ID,
		// 	Amount:        bOffer.Amount,
		// 	SourceAccount: mmWallet.ID,
		// })
		memo = fmt.Sprintf("cancelOffer %v", bOffer.ID)

	}
	if strings.EqualFold(offerRequest.OfferType, "SELL") {

		// cancel offer with net 0 so that fee can be returned.
		ops = append(ops, &basetxn.ManageSellOffer{
			OfferID:       offerIDInt,
			Selling:       mainAsset,
			Buying:        currencyAsset,
			Amount:        "0",
			Price:         fmt.Sprintf("%v/%v", n, d),
			SourceAccount: sourceWallet.ID,
		})

		// // move sellinng funds from the MM wallet to the offer source wallet
		// ops = append(ops, &basetxn.Payment{
		// 	Asset:         mainAsset,
		// 	Destination:   sourceWallet.ID,
		// 	Amount:        bOffer.Amount,
		// 	SourceAccount: mmWallet.ID,
		// })
		memo = fmt.Sprintf("cancelOffer %v", bOffer.ID)
	}

	// //service fee
	// if decimal.RequireFromString(offerRequest.FeeValue).IsPositive() {

	// 	//process service fee

	// 	//no need deducting it as we deduct it as market executes
	// 	feeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
	// 	if !mainAsset.IsNative() {
	// 		feeAssetCode = mainAsset.GetCode()
	// 	}
	// 	offerRequest.Messages = append(offerRequest.Messages, fmt.Sprintf("%v %v (%v) will be deducted from the total quantity as service fee and your offer will be placed with %v %v.", offerRequest.FeeValue, feeAssetCode, offerFeePercentage, offerRequest.NetQuantity, feeAssetCode))

	// }

	// Construct the transaction that holds the operations to execute on the network

	var tx *basetxn.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if sourceWallet.NumberOfApprovalsNeeded > 0 {
		transactionSource = chanSourceAccount.Address
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        chanSourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 memo,
			},
		)
	} else {
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        sourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 memo,
			},
		)
	}
	if err != nil {
		log.Println("[generateDeleteMarketXdr]error constructing transaction", err)
		return "", "", &tErrors.ErrorTemporaryServerError{}
	}

	// tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), mmSignerKeyPair)

	// if err != nil {
	// 	log.Println("[generateDeleteMarketXdr] error signing transaction with custodial signer key ", err)
	// 	return "", "", &tErrors.ErrorTemporaryServerError{}
	// }

	if sourceWallet.NumberOfApprovalsNeeded > 0 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

		if err != nil {
			log.Println("[generateDeleteMarketXdr] error signing transaction with channelAccount key ", err)
			return "", "", &tErrors.ErrorTemporaryServerError{}
		}
	}

	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", "", err
	}
	return xdrBase64, transactionSource, nil

}

// ToFractionInt32 approximates a float64 as a fraction with int32 numerator and denominator.
// It uses a continued fraction algorithm for a good approximation.
func ToFractionInt32(x float64) (num, den int32) {
	if x < 0 {
		return 0, 0 // Handle negative numbers as per requirement
	}

	const (
		maxDen = math.MaxInt32
		maxNum = math.MaxInt32
	)

	// Initial values for the continued fraction
	n1, d1 := int64(math.Floor(x)), int64(1)
	n2, d2 := int64(1), int64(0)
	x -= math.Floor(x)

	for x != 0 {
		a := math.Floor(1 / x)
		n3 := n1*int64(a) + n2
		d3 := d1*int64(a) + d2

		if d3 > maxDen || n3 > maxNum {
			break // Stop if numerator or denominator exceeds int32 limits
		}

		// Update values for the next iteration
		n2, d2 = n1, d1
		n1, d1 = n3, d3
		x = 1/x - a
	}

	return int32(n1), int32(d1)
}
