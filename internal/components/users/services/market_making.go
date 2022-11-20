package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	bc "trovo-wallet-api/internal/blockchainalgofuncs"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func MakeOffer(signerUser, walletOwner *userModels.User, sourceWallet *userModels.UserWallet, offerRequest *userModels.MarketOfferRequest, gc *sharedconfig.GlobalConfig) (err error) {
	offerRequest.Messages = make([]string, 0)

	if sourceWallet.SharedAccessEnabled == 1 && sourceWallet.NumberOfApprovalsNeeded > 0 {
		offerRequest.Multiparty = 1
	}
	if sourceWallet.HasViewOnlyAccess(gc) {
		offerRequest.SignatureRequired = 1
	}
	mmWallet, err := walletOwner.GetMartketMakingWallet(gc.DB)
	if err != nil {
		log.Printf("[MakeOffer]Error validating Market making wallet: %v", err)
		return err
	}
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
	mmSignerKeyPair, err := bc.MarketMakingSignerKeypair(walletOwner.Username, mmWallet.ID)
	if err != nil {

		return &tErrors.ErrorTemporaryServerError{}
	}
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
	if e != nil {
		serviceFee = decimal.Zero
	}
	totalQty := decimal.RequireFromString(offerRequest.Quantity).Truncate(7)
	feeQuantity := (totalQty.Mul(serviceFee.Div(decimal.NewFromInt(100)))).Truncate(7)
	netQuantity := totalQty.Sub(feeQuantity)
	offerRequest.NetQuantity = netQuantity.String()
	offerRequest.FeeChargedOnAsset = serviceFee.String()
	offerRequest.FeeValue = feeQuantity.String()

	var marketOffer *userModels.MarketOffer
	xdrBase64, err := generateMakeMarketXdr(sourceWallet, &mmWallet, offerRequest, mmSignerKeyPair, gc)
	if err != nil {
		return err
	}
	offerRequest.Transaction = xdrBase64
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
	{
		marketOffer = &userModels.MarketOffer{
			ID:                          uuid.NewString(),
			SourceWalletAlias:           sourceWallet.Alias,
			SourcewalletPublicKey:       sourceWallet.ID,
			MarketMakingWalletPublicKey: mmWallet.ID,
			OfferType:                   offerRequest.OfferType,
			AssetCode:                   offerRequest.AssetCode,
			AssetIssuer:                 &offerRequest.AssetIssuer,
			CurrencyCode:                offerRequest.CurrencyCode,
			CurrencyIssuer:              &offerRequest.CurrencyIssuer,
			PricePerUnit:                offerRequest.PricePerUnit,
			Quantity:                    offerRequest.Quantity,
			FeeChargedOnAsset:           offerRequest.FeeChargedOnAsset,
			FeeValue:                    offerRequest.FeeValue,
			NetQuantity:                 offerRequest.NetQuantity,
			RemainingQuantity:           offerRequest.NetQuantity,
			RemainingFeeValue:           offerRequest.FeeValue,
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
		e = dbTX.Create(&marketOffer).Error
		if e != nil {
			log.Printf("[MakeOffer]Error saving market offer: %+v\nError: %v\n", marketOffer, err)
			return &tErrors.ErrorTemporaryServerError{}
		}
		txnResult, err := network.SubmitXdrWithSignatureReturnsTrx(gc.BantuExpansionClient, sourceWallet.Signer, xdrBase64, offerRequest.TransactionSignature)

		if err != nil {
			log.Printf("[MakeOffer]error submitting txn: %v\n", err)
			return &tErrors.ErrorTemporaryServerError{}
		}

		offerRequest.TransactionID = txnResult.Hash
		marketOffer.TransactionID = &txnResult.Hash
		//get and set the offerID
		{
			var re xdr.TransactionResult
			e := xdr.SafeUnmarshalBase64(txnResult.ResultXdr, &re)
			if e != nil {
				fmt.Println(e)
			}
			log.Println(re)
			or, _ := re.OperationResults()

			// log.Printf("OK: %v\nOperationResult: %#v", ok, or)
			for _, r := range or {

				ms, ok := r.Tr.GetManageSellOfferResult()
				if !ok {
					continue
				}
				// log.Printf("Index: %v\nManageSellOffer Offer ID: %#v", i, ms.Success.Offer.Offer.OfferId)
				offerID := fmt.Sprintf("%v", ms.Success.Offer.Offer.OfferId)

				marketOffer.BlockchainOfferID = &offerID
			}
		}
		e = dbTX.Save(&marketOffer).Error
		if e != nil {
			log.Printf("[MakeOffer]Error saving transactionID on market offer: %+v\nError: %v\n", marketOffer, err)
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
			ApprovalsNeeded:          sourceWallet.NumberOfApprovalsNeeded,
			TransactionXdr:           xdrBase64,
			TransactionInfoStr:       &transactionStr,
		}
		//save and commit this to database
		e := gc.DB.Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[MakeOffer] Error saving MAKE MARKET OFFER txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}
			return err
		}
		return nil

	}

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

	if decimal.RequireFromString(marketOffer.RemainingQuantity).IsZero() {
		return &tErrors.CustomError{
			Param:      "id",
			Err:        "error-offer-has been filled",
			ErrMessage: "Offer cannot be canceled because it has been filled.",
		}
	}
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

	mmWallet, err := walletOwner.GetMartketMakingWallet(gc.DB)
	if err != nil {
		log.Printf("[MakeOffer]Error validating Market making wallet: %v", err)
		return err
	}

	mmSignerKeyPair, err := bc.MarketMakingSignerKeypair(walletOwner.Username, mmWallet.ID)
	if err != nil {

		return &tErrors.ErrorTemporaryServerError{}
	}

	xdrBase64, err := generateDeleteMarketXdr(sourceWallet, &mmWallet, &marketOffer, mmSignerKeyPair, bOffer, gc)
	if err != nil {
		return err
	}
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
		e := dbTX.Create(&marketOffer).Error
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
		//get and set the offerID
		{
			var re xdr.TransactionResult
			e := xdr.SafeUnmarshalBase64(txnResult.ResultXdr, &re)
			if e != nil {
				fmt.Println(e)
			}
			log.Println(re)
			or, _ := re.OperationResults()

			// log.Printf("OK: %v\nOperationResult: %#v", ok, or)
			for _, r := range or {

				ms, ok := r.Tr.GetManageSellOfferResult()
				if !ok {
					continue
				}
				// log.Printf("Index: %v\nManageSellOffer Offer ID: %#v", i, ms.Success.Offer.Offer.OfferId)
				offerID := fmt.Sprintf("%v", ms.Success.Offer.Offer.OfferId)

				marketOffer.BlockchainOfferID = &offerID
			}
		}
		e = dbTX.Save(&marketOffer).Error
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
			ApprovalsNeeded:          sourceWallet.NumberOfApprovalsNeeded,
			TransactionXdr:           xdrBase64,
			TransactionInfoStr:       &transactionStr,
		}
		//save and commit this to database
		e := gc.DB.Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[MakeOffer] Error saving MAKE MARKET OFFER txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}
			return err
		}
		return nil

	}

	return &tErrors.ErrorTemporaryServerError{}

}

func generateMakeMarketXdr(sourceWallet, mmWallet *userModels.UserWallet, offerRequest *userModels.MarketOfferRequest, mmSignerKeyPair *keypair.Full, gc *sharedconfig.GlobalConfig) (txnBase64 string, err error) {
	offerFeePercentage := offerRequest.FeeChargedOnAsset + "%"
	minBalance := decimal.RequireFromString(os.Getenv("STANDARD_WALLET_MINIMUM_BALANCE"))
	offerRequest.Messages = make([]string, 0)
	var memo string
	var currencyAsset, mainAsset txnbuild.Asset
	if offerRequest.AssetCode+offerRequest.AssetIssuer == offerRequest.CurrencyCode+offerRequest.CurrencyIssuer {
		return "", &tErrors.CustomError{
			Param:      "AssetCode",
			Err:        "error-asset-and-currency-are-the-same",
			ErrMessage: "You cannot make a market against same asset",
		}
	}
	offerRequest.OfferType = strings.ToUpper(offerRequest.OfferType)
	fraction := decimal.RequireFromString(offerRequest.PricePerUnit).Rat()
	d := int32(fraction.Denom().Int64())
	n := int32(fraction.Num().Int64())
	if offerRequest.AssetIssuer == "" {
		mainAsset = txnbuild.NativeAsset{}
	} else {
		mainAsset = txnbuild.CreditAsset{Code: offerRequest.AssetCode, Issuer: offerRequest.AssetIssuer}
	}

	if offerRequest.CurrencyIssuer == "" {
		currencyAsset = txnbuild.NativeAsset{}
	} else {
		currencyAsset = txnbuild.CreditAsset{Code: offerRequest.CurrencyCode, Issuer: offerRequest.CurrencyIssuer}
	}

	chanAccount := <-gc.ChannelAccounts
	defer func(c *keypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)

	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, chanAccount.Address(), txnbuild.NativeAsset{})
	var sourceAccount *horizon.Account
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
		sourceMAccountExists, sourceMAccountTrustsAsset, nativeMAccountBalance, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, sourceWallet.ID, mainAsset)

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

	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)

	if strings.EqualFold(offerRequest.OfferType, "BUY") {
		if !currencyAsset.IsNative() {
			//check if it has trustline to it and then create it.
			_, mmAccountTrustsAsset, mmnativeAccountBalance, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, mmWallet.ID, currencyAsset)
			if mmnativeAccountBalance.LessThan(minBalance) {
				ops = append(ops, &txnbuild.Payment{
					Asset:         txnbuild.NativeAsset{},
					Destination:   mmWallet.ID,
					Amount:        minBalance.Mul(decimal.NewFromInt(3)).String(),
					SourceAccount: sourceWallet.ID,
				})
			}

			if !mmAccountTrustsAsset {
				//establish trustline automatically
				ops = append(ops, &txnbuild.ChangeTrust{
					Line:          txnbuild.ChangeTrustAssetWrapper{Asset: currencyAsset},
					Limit:         "900000000000",
					SourceAccount: mmWallet.ID,
				})
			}

		}
		{
			_, mmAccountTrustsAsset, _, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, mmWallet.ID, mainAsset)
			if !mmAccountTrustsAsset {
				//establish trustline automatically
				ops = append(ops, &txnbuild.ChangeTrust{
					Line:          txnbuild.ChangeTrustAssetWrapper{Asset: mainAsset},
					Limit:         "900000000000",
					SourceAccount: mmWallet.ID,
				})
			}
		}
		// move sellinng funds to the MM wallet
		ops = append(ops, &txnbuild.Payment{
			Asset:         currencyAsset,
			Destination:   mmWallet.ID,
			Amount:        offerRequest.Quantity,
			SourceAccount: sourceWallet.ID,
		})

		//invert the price fraction
		ops = append(ops, &txnbuild.ManageSellOffer{
			Selling:       currencyAsset,
			Buying:        mainAsset,
			Amount:        offerRequest.NetQuantity,
			Price:         xdr.Price{D: xdr.Int32(n), N: xdr.Int32(d)},
			SourceAccount: mmWallet.ID,
		})
		memo = fmt.Sprintf("s%v-b%v", currencyAsset.GetCode(), mainAsset.GetCode())

	}
	if strings.EqualFold(offerRequest.OfferType, "SELL") {
		if !mainAsset.IsNative() {
			//check if it has trustline to it and then create it.
			_, mmAccountTrustsAsset, mmnativeAccountBalance, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, mmWallet.ID, mainAsset)
			if mmnativeAccountBalance.LessThan(minBalance) {
				ops = append(ops, &txnbuild.Payment{
					Asset:         txnbuild.NativeAsset{},
					Destination:   mmWallet.ID,
					Amount:        minBalance.Mul(decimal.NewFromInt(3)).String(),
					SourceAccount: sourceWallet.ID,
				})
			}

			if !mmAccountTrustsAsset {
				//establish trustline automatically
				ops = append(ops, &txnbuild.ChangeTrust{
					Line:          txnbuild.ChangeTrustAssetWrapper{Asset: mainAsset},
					Limit:         "900000000000",
					SourceAccount: mmWallet.ID,
				})
			}

		}
		{
			_, mmAccountTrustsAsset, _, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, mmWallet.ID, currencyAsset)
			if !mmAccountTrustsAsset {
				//establish trustline automatically
				ops = append(ops, &txnbuild.ChangeTrust{
					Line:          txnbuild.ChangeTrustAssetWrapper{Asset: currencyAsset},
					Limit:         "900000000000",
					SourceAccount: mmWallet.ID,
				})
			}
		}
		// move sellinng funds to the MM wallet
		ops = append(ops, &txnbuild.Payment{
			Asset:         mainAsset,
			Destination:   mmWallet.ID,
			Amount:        offerRequest.Quantity,
			SourceAccount: sourceWallet.ID,
		})

		// make offer with net quantity so that fee can be returned.
		ops = append(ops, &txnbuild.ManageSellOffer{
			Selling:       mainAsset,
			Buying:        currencyAsset,
			Amount:        offerRequest.NetQuantity,
			Price:         xdr.Price{N: xdr.Int32(n), D: xdr.Int32(d)},
			SourceAccount: mmWallet.ID,
		})
		memo = fmt.Sprintf("s%v-b%v", mainAsset.GetCode(), currencyAsset.GetCode())
	}

	//service fee
	if decimal.RequireFromString(offerRequest.FeeValue).IsPositive() {

		//process service fee

		// ops = append(ops, &txnbuild.Payment{
		// 	Destination:   os.Getenv("MARKET_MAKING_FEE_ADDRESS"),
		// 	Amount:        offerRequest.FeeValue,
		// 	SourceAccount: sourceWallet.ID,
		// 	Asset:         mainAsset,
		// })
		//no need deducting it as we deduct it as market executes
		feeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
		if !mainAsset.IsNative() {
			feeAssetCode = mainAsset.GetCode()
		}
		offerRequest.Messages = append(offerRequest.Messages, fmt.Sprintf("%v %v (%v) will be deducted from the total quantity as service fee and your offer will be placed with %v %v.", offerRequest.FeeValue, feeAssetCode, offerFeePercentage, offerRequest.NetQuantity, feeAssetCode))

	}

	// Construct the transaction that holds the operations to execute on the network

	var tx *txnbuild.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if offerRequest.Multiparty == 1 {
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        chanSourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              txnbuild.MinBaseFee,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText(memo),
			},
		)
	} else {
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        sourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              txnbuild.MinBaseFee,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText(memo),
			},
		)
	}
	if err != nil {
		log.Println("[generateMakeMarketXdr]error constructing transaction", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), mmSignerKeyPair)

	if err != nil {
		log.Println("[generateMakeMarketXdr] error signing transaction with custodial signer key ", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	if offerRequest.Multiparty == 1 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

		if err != nil {
			log.Println("[generateMakeMarketXdr] error signing transaction with channelAccount key ", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}

	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", err
	}
	offerRequest.Memo = memo
	return xdrBase64, nil

}
func generateDeleteMarketXdr(sourceWallet, mmWallet *userModels.UserWallet, offerRequest *userModels.MarketOffer, mmSignerKeyPair *keypair.Full, bOffer horizon.Offer, gc *sharedconfig.GlobalConfig) (txnBase64 string, err error) {
	// offerFeePercentage := offerRequest.FeeChargedOnAsset + "%"
	minBalance := decimal.RequireFromString(os.Getenv("STANDARD_WALLET_MINIMUM_BALANCE"))
	// offerRequest.Messages = make([]string, 0)
	var memo string
	var currencyAsset, mainAsset txnbuild.Asset

	offerRequest.OfferType = strings.ToUpper(offerRequest.OfferType)
	fraction := decimal.RequireFromString(offerRequest.PricePerUnit).Rat()
	d := int32(fraction.Denom().Int64())
	n := int32(fraction.Num().Int64())
	if offerRequest.AssetIssuer == nil {
		mainAsset = txnbuild.NativeAsset{}
	} else {
		mainAsset = txnbuild.CreditAsset{Code: offerRequest.AssetCode, Issuer: *offerRequest.AssetIssuer}
	}

	if offerRequest.CurrencyIssuer == nil {
		currencyAsset = txnbuild.NativeAsset{}
	} else {
		currencyAsset = txnbuild.CreditAsset{Code: offerRequest.CurrencyCode, Issuer: *offerRequest.CurrencyIssuer}
	}

	chanAccount := <-gc.ChannelAccounts
	defer func(c *keypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)

	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, chanAccount.Address(), txnbuild.NativeAsset{})
	var sourceAccount *horizon.Account
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
		sourceMAccountExists, sourceMAccountTrustsAsset, nativeMAccountBalance, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, sourceWallet.ID, mainAsset)

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

	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)

	if strings.EqualFold(offerRequest.OfferType, "BUY") {

		//invert the price fraction
		ops = append(ops, &txnbuild.ManageSellOffer{
			OfferID:       bOffer.ID,
			Selling:       currencyAsset,
			Buying:        mainAsset,
			Amount:        "0",
			Price:         xdr.Price{D: xdr.Int32(n), N: xdr.Int32(d)},
			SourceAccount: mmWallet.ID,
		})

		// move selling funds from the MM wallet to the wallet it was made from
		ops = append(ops, &txnbuild.Payment{
			Asset:         currencyAsset,
			Destination:   sourceWallet.ID,
			Amount:        bOffer.Amount,
			SourceAccount: mmWallet.ID,
		})
		memo = fmt.Sprintf("cancelOffer %v", bOffer.ID)

	}
	if strings.EqualFold(offerRequest.OfferType, "SELL") {

		// cancel offer with net 0 so that fee can be returned.
		ops = append(ops, &txnbuild.ManageSellOffer{
			OfferID:       bOffer.ID,
			Selling:       mainAsset,
			Buying:        currencyAsset,
			Amount:        "0",
			Price:         xdr.Price{N: xdr.Int32(n), D: xdr.Int32(d)},
			SourceAccount: mmWallet.ID,
		})

		// move sellinng funds from the MM wallet to the offer source wallet
		ops = append(ops, &txnbuild.Payment{
			Asset:         mainAsset,
			Destination:   sourceWallet.ID,
			Amount:        bOffer.Amount,
			SourceAccount: mmWallet.ID,
		})
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

	var tx *txnbuild.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if sourceWallet.NumberOfApprovalsNeeded > 0 {
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        chanSourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              txnbuild.MinBaseFee,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText(memo),
			},
		)
	} else {
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        sourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              txnbuild.MinBaseFee,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText(memo),
			},
		)
	}
	if err != nil {
		log.Println("[generateDeleteMarketXdr]error constructing transaction", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), mmSignerKeyPair)

	if err != nil {
		log.Println("[generateDeleteMarketXdr] error signing transaction with custodial signer key ", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	if sourceWallet.NumberOfApprovalsNeeded > 0 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

		if err != nil {
			log.Println("[generateDeleteMarketXdr] error signing transaction with channelAccount key ", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}

	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", err
	}
	return xdrBase64, nil

}
