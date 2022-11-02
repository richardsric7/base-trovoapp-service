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

func MakeOffer(signerUser *userModels.User, sourceWallet *userModels.UserWallet, offerRequest *userModels.MakeOfferRequest, gc *sharedconfig.GlobalConfig) (err error) {
	if sourceWallet.SharedAccessEnabled == 1 && sourceWallet.NumberOfApprovalsNeeded > 0 {
		offerRequest.Multiparty = 1
	}
	walletOwner, err := sourceWallet.GetWalletOwner(gc.DB)
	if err != nil {
		log.Printf("[MakeOffer]Error validating Market making wallet: %v", err)
		return err
	}

	mmWallet, err := walletOwner.GetMartketMakingWallet(gc.DB)
	if err != nil {
		log.Printf("[MakeOffer]Error validating Market making wallet: %v", err)
		return err
	}
	mmSignerKeyPair, err := bc.MarketMakingSignerKeypair(walletOwner.Username, mmWallet.ID)
	if err != nil {
		return &tErrors.CustomError{
			Param:      "marketMakingWalletPK",
			Err:        "error-retrieving-market-making-wallet",
			ErrMessage: "Market Making wallet could not be validated.",
		}
	}
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
	if (offerRequest.Commit == 0 && sourceWallet.SharedAccessEnabled == 1 && sourceWallet.NumberOfApprovalsNeeded == 0) || (sourceWallet.SharedAccessEnabled == 0 && offerRequest.Commit == 0) {
		txnID, err := network.SubmitXdrWithSignature(gc.BantuExpansionClient, sourceWallet.Signer, xdrBase64, offerRequest.TransactionSignature)

		if err != nil {
			log.Printf("[MakeOffer]error submitting txn: %v\n", err)
			return &tErrors.ErrorTemporaryServerError{}
		}

		offerRequest.TransactionID = txnID

		return nil
	}

	if offerRequest.Multiparty == 1 {
		offerRequest.TransactionID = "PENDING_AUTH"
		log.Printf("[MakeOffer]shared access with approver permission enabled for %v \n", sourceWallet.Alias)
		id := uuid.NewString()

		description := offerRequest.ReturnedDescription
		transactionByte, _ := json.Marshal(*offerRequest)
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
			log.Printf("[MakeOffer] Error saving payment txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}
			return err
		}
		return nil

	}

	return &tErrors.ErrorTemporaryServerError{}
}

func generateMakeMarketXdr(sourceWallet, mmWallet *userModels.UserWallet, offerRequest *userModels.MakeOfferRequest, mmSignerKeyPair *keypair.Full, gc *sharedconfig.GlobalConfig) (txnBase64 string, err error) {

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
	fraction := decimal.RequireFromString(offerRequest.PricePerAsset).Rat()
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
	if strings.EqualFold(offerRequest.OfferType, "Buy") {
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
	if strings.EqualFold(offerRequest.OfferType, "Sell") {
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

	if strings.EqualFold(offerRequest.OfferType, "Buy") {
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
			Amount:        offerRequest.Quantity,
			Price:         xdr.Price{D: xdr.Int32(n), N: xdr.Int32(d)},
			SourceAccount: mmWallet.ID,
		})
		memo = "sell-" + currencyAsset.GetCode()

	}
	if strings.EqualFold(offerRequest.OfferType, "Sell") {
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

		ops = append(ops, &txnbuild.ManageSellOffer{
			Selling:       mainAsset,
			Buying:        currencyAsset,
			Amount:        offerRequest.Quantity,
			Price:         xdr.Price{N: xdr.Int32(n), D: xdr.Int32(d)},
			SourceAccount: mmWallet.ID,
		})
		memo = "sell-" + mainAsset.GetCode()
	}

	//service fee

	serviceFee, e := decimal.NewFromString(os.Getenv("MARKET_MAKING_FEE_AMOUNT"))
	if e != nil {
		serviceFee = decimal.Zero
	}
	if serviceFee.IsPositive() {
		if sourceWallet.FeeDisabled == 0 {
			//process service fee
			if len(os.Getenv("MARKET_MAKING_FEE_ASSET_ISSUER")) == 0 {
				ops = append(ops, &txnbuild.Payment{
					Destination:   os.Getenv("MARKET_MAKING_FEE_ADDRESS"),
					Amount:        os.Getenv("MARKET_MAKING_FEE_AMOUNT"),
					SourceAccount: sourceWallet.ID,
					Asset:         txnbuild.NativeAsset{},
				})
				offerRequest.Messages = append(offerRequest.Messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee.", os.Getenv("MARKET_MAKING_FEE_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE"), sourceWallet.Alias))

			} else {
				ops = append(ops, &txnbuild.Payment{
					Destination:   os.Getenv("MARKET_MAKING_FEE_ADDRESS"),
					Amount:        os.Getenv("MARKET_MAKING_FEE_AMOUNT"),
					SourceAccount: sourceWallet.ID,
					Asset:         txnbuild.CreditAsset{Code: os.Getenv("MARKET_MAKING_FEE_ASSET_CODE"), Issuer: os.Getenv("MARKET_MAKING_FEE_ASSET_ISSUER")},
				})
				offerRequest.Messages = append(offerRequest.Messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee.", os.Getenv("MARKET_MAKING_FEE_AMOUNT"), os.Getenv("MARKET_MAKING_FEE_ASSET_CODE"), sourceWallet.Alias))

			}

		}
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
