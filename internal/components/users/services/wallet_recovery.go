package users

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	bc "trovo-wallet-api/internal/blockchainalgofuncs"
	userBc "trovo-wallet-api/internal/components/users/blockchain"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ecnepsnai/discord"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

func EnableAccountRecovery(user *userModels.User, payload *userModels.UserAccountRecoveryPayload, gc *sharedconfig.GlobalConfig) (err error) {
	var e error
	client := gc.BantuExpansionClient
	multiAccessWallets := make([]userModels.UserWallet, 0)
	ops := make([]txnbuild.Operation, 0)
	messages := make([]string, 0)
	payload.Messages = make([]string, 0)
	if user.HasSecurityQuestions == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error security answers not set", ErrMessage: "security answers has not been set for this account."}
	}
	if user.AccountRecoveryEnabled == 1 {
		return &tErrors.CustomError{Param: "username", Err: "error account recovery already enabled.", ErrMessage: "Account recovery already enabled."}
	}
	wallet, _ := userModels.UserWalletID(user.PublicKey).GetWallet(gc.DB)
	if wallet.SharedAccessEnabled == 1 && !WalletHasViewOnlyAccess(&wallet, gc) {
		return &tErrors.CustomError{Param: "username", Err: "error primary wallet has shared access.", ErrMessage: "Primary wallet has shared access enabled! Only primary wallets without shared access or with view only shared access can participate in account recovery at this time."}
	}

	dbtx := gc.DB.Begin()
	defer dbtx.Rollback()
	user.AccountRecoveryEnabled = 1
	exp := time.Now().AddDate(1, 0, 0)
	user.AccountRecoveryExpiresOn = &exp
	dbErr := dbtx.Save(user).Error
	if dbErr != nil {
		log.Printf("[EnableAccountRecovery] Error saving account recovery state: %v\n", dbErr)
		return &tErrors.ErrorTemporaryServerError{}
	}

	var userAccount horizon.Account
	if userAccount, e = userBc.GetBlockchainAccountDetail(user.PublicKey); e != nil {
		if e.Error() == "error-blockchain-account-not-activated" {
			return &tErrors.CustomError{Param: "username", Err: "error primary account not yet activated", ErrMessage: fmt.Sprintf("Primary account is not yet activated. Please send upto 50 %v to the primary wallet to continue.", os.Getenv("NATIVE_ASSET_CODE"))}
		}
		log.Printf("[EnableAccountRecovery] Error on blockchain validating primary account %v\n", e)
		return &tErrors.ErrorTemporaryServerError{}

	}
	for _, v := range userAccount.Balances {
		if v.Code == "" && decimal.RequireFromString(v.Balance).LessThan(decimal.RequireFromString(os.Getenv("ACCOUNT_RECOVERY_MINIMUM_BALANCE"))) {
			return &tErrors.CustomError{Param: "username", Err: "error primary wallet needs funding", ErrMessage: fmt.Sprintf("Primary wallet needs minimum of 50 %v to proceed.", os.Getenv("NATIVE_ASSET_CODE"))}

		}
	}

	//get the recovery keypair
	recoveryAddress := bc.GetRecoveryAccountAddress(user.Username, user.PublicKey)
	if len(recoveryAddress) == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error generating recovery address", ErrMessage: "Could not generate valid recovery address for account."}

	}
	//activate recovery Address and add signer to primary key
	_, e = userBc.GetBlockchainAccountDetail(recoveryAddress)
	if e != nil {
		if e.Error() == "error-blockchain-account-not-activated" {
			//activate account
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   recoveryAddress,
				Amount:        os.Getenv("RECOVERY_SIGNER_ACTIVATION_AMOUNT"),
				SourceAccount: user.PublicKey,
			})
			messages = append(messages, fmt.Sprintf("%v %v will be deducted from your wallet [%v] to activate your unique recovery key on the blockchain.", os.Getenv("RECOVERY_SIGNER_ACTIVATION_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE"), user.Username))

		}
	}
	//  else {
	// 	//account already active
	// 	ops = append(ops, &txnbuild.Payment{
	// 		Destination:   recoveryAddress,
	// 		Amount:        "6",
	// 		Asset:         txnbuild.NativeAsset{},
	// 		SourceAccount: user.PublicKey,
	// 	})
	// }
	if !userBc.SignerIsValid(user.PublicKey, recoveryAddress) {
		//recovery not a signer to the primary wallet.
		ops = append(ops, &txnbuild.SetOptions{
			Signer: &txnbuild.Signer{
				Address: recoveryAddress,
				Weight:  1,
			},
			SourceAccount: user.PublicKey,
		})
	}
	//check for subwallets
	wallets := user.GetAllWallets(gc)
	bulkPaymentSignerKeyPairs, marketMakingSignerKeyPairs := make([]*keypair.Full, 0), make([]*keypair.Full, 0)

	if len(wallets) > 1 {
		//has subwallets other than the primary wallet, which has already be added to the ops
		for _, w := range wallets {
			if w.PrimaryWallet == 1 {
				// ensures the primary wallet is not added to the ops
				continue
			}
			if w.PrimaryWallet == 0 {
				if w.SharedAccessEnabled == 1 {
					if !w.HasViewOnlyAccess(gc) {
						// wallet does not have only view-only access, so we need to add it to the multiaccessWallets list
						multiAccessWallets = append(multiAccessWallets, w)
						continue
					}
				}
				//a subwallet, not primary wallet
				//activate subwallet Address and add signer to subwallet key
				_, e = userBc.GetBlockchainAccountDetail(w.ID)
				if e != nil {
					if e.Error() == "error-blockchain-account-not-activated" {
						//activate account
						ops = append(ops, &txnbuild.CreateAccount{
							Destination:   w.ID,
							Amount:        os.Getenv("RECOVERY_SIGNER_ACTIVATION_AMOUNT"),
							SourceAccount: user.PublicKey,
						})
						messages = append(messages, fmt.Sprintf("%v %v will be deducted from your wallet [%v] to activate your subwallet [%v] on the blockchain.", os.Getenv("RECOVERY_SIGNER_ACTIVATION_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE"), user.Username, w.Alias))

					}
				}
				if w.WalletType == 0 || w.WalletType == 1 {
					if !userBc.SignerIsValid(w.ID, recoveryAddress) {
						//recovery not a signer to the sub wallet. add it
						ops = append(ops, &txnbuild.SetOptions{
							Signer: &txnbuild.Signer{
								Address: recoveryAddress,
								Weight:  1,
							},
							SourceAccount: w.ID,
						})
					}
					// {
					// 	//adjust account threshold

					// 	ops = append(ops, &txnbuild.SetOptions{
					// 		LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(1)),
					// 		MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(1)),
					// 		HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(1)),
					// 		SourceAccount:   wallet.ID,
					// 	})

					// }
				}
				if w.WalletType == 2 {
					mm, _ := bc.MarketMakingSignerKeypair(user.Username, w.ID)
					marketMakingSignerKeyPairs = append(marketMakingSignerKeyPairs, mm)
					if !userBc.SignerIsValid(w.ID, recoveryAddress) {
						//recovery not a signer to the sub wallet. add it
						ops = append(ops, &txnbuild.SetOptions{
							Signer: &txnbuild.Signer{
								Address: recoveryAddress,
								Weight:  3,
							},
							SourceAccount: w.ID,
						})
					}

				}
				if w.WalletType == 3 {
					bp, _ := bc.BulkPaymentSignerKeypair(user.Username, w.ID)
					bulkPaymentSignerKeyPairs = append(bulkPaymentSignerKeyPairs, bp)
					if !userBc.SignerIsValid(w.ID, recoveryAddress) {
						//recovery not a signer to the sub wallet. add it
						ops = append(ops, &txnbuild.SetOptions{
							Signer: &txnbuild.Signer{
								Address: recoveryAddress,
								Weight:  3,
							},
							SourceAccount: w.ID,
						})
					}
				}

			}

		}
	}

	{ //add fee for transaction

	}

	var xdrBase64 string
	payload.Messages = messages
	payload.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	// oldTrx := payload.Transaction

	// generate new transaction

	memo := "Enabling Account Recovery"
	log.Println("[EnableAccountRecovery] Memo:", memo)

	if len(ops) == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error no operations to perform", ErrMessage: "Could not find any operations to perform for this action."}

	}
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &userAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              2000,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
			Memo:                 txnbuild.MemoText(memo),
		},
	)

	if err != nil {
		log.Println("[EnableAccountRecovery]error constructing transaction ", err)
		return &tErrors.ErrorTemporaryServerError{}
	}
	if len(bulkPaymentSignerKeyPairs) > 0 {
		tx, err = tx.Sign(gc.BantuNetworkPassphrase, bulkPaymentSignerKeyPairs...)
		if err != nil {
			log.Println("[EnableAccountRecovery] error signning with bulkPaymentSignerKeyPairs", err)
			return &tErrors.ErrorTemporaryServerError{}
		}
	}
	if len(marketMakingSignerKeyPairs) > 0 {
		tx, err = tx.Sign(gc.BantuNetworkPassphrase, marketMakingSignerKeyPairs...)
		if err != nil {
			log.Println("[EnableAccountRecovery] error signning with marketMakingSignerKeyPairs", err)
			return &tErrors.ErrorTemporaryServerError{}
		}
	}
	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[EnableAccountRecovery] error getting txn base64", err)
		return &tErrors.ErrorTemporaryServerError{}
	}
	if len(payload.Transaction) == 0 || len(payload.TransactionSignature) == 0 {
		payload.Transaction = xdrBase64

		return nil
	}

	// if oldTrx != payload.Transaction {
	// 	return &tErrors.ErrorInvalidTransaction{}
	// }

	if len(payload.TransactionSignature) == 0 {
		return &tErrors.CustomError{Param: "TransactionSignature", Err: "error transaction signature is required", ErrMessage: "Transaction signature is required."}
	}
	//submit to blockchain

	txnHash, err := network.SubmitXdrWithSignature(client, user.PrimarySigner, xdrBase64, payload.TransactionSignature)
	if err != nil {
		logDiscordFailedRecovery(fmt.Sprintf("Error submitting account recovery enable [%+v] transaction: %s", payload, err.Error()))
		return
	}
	payload.TransactionID = txnHash
	dbtx.Commit()
	if err == nil && len(multiAccessWallets) > 0 {
		log.Printf("Skipped wallets %+v\n", multiAccessWallets)
	}
	return nil

}
func DisableAccountRecovery(user *userModels.User, payload *userModels.UserAccountRecoveryPayload, gc *sharedconfig.GlobalConfig) (err error) {
	var e error
	client := gc.BantuExpansionClient
	multiAccessWallets := make([]userModels.UserWallet, 0)
	ops := make([]txnbuild.Operation, 0)
	messages := make([]string, 0)
	payload.Messages = make([]string, 0)
	if user.HasSecurityQuestions == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error security answers not set", ErrMessage: "Security answers has not been set for this account."}
	}
	if user.AccountRecoveryEnabled == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error account recovery not enabled.", ErrMessage: "Account recovery not enabled."}
	}
	dbtx := gc.DB.Begin()
	defer dbtx.Rollback()
	user.AccountRecoveryEnabled = 0
	// exp := time.Now().AddDate(1, 0, 0)
	user.AccountRecoveryExpiresOn = nil
	dbErr := dbtx.Save(user).Error
	if dbErr != nil {
		log.Printf("[DisableAccountRecovery] Error saving account recovery state: %v\n", dbErr)
		return &tErrors.ErrorTemporaryServerError{}
	}

	if !ValidateSecurityAnswers(user, payload.SecurityAnswers, gc) {
		return &tErrors.CustomError{Param: "username", Err: "error invalid security answers", ErrMessage: "Answers to the security questions are invalid."}
	}
	var userAccount horizon.Account
	if userAccount, e = userBc.GetBlockchainAccountDetail(user.PublicKey); e != nil {
		if e.Error() == "error-blockchain-account-not-activated" {
			return &tErrors.CustomError{Param: "username", Err: "error primary account not yet activated", ErrMessage: fmt.Sprintf("Primary account is not yet activated. Please send upto 50 %v to the primary wallet to continue.", os.Getenv("NATIVE_ASSET_CODE"))}
		}
		log.Printf("[DisableAccountRecovery] Error on blockchain validating primary account %v\n", e)
		return &tErrors.ErrorTemporaryServerError{}

	}
	for _, v := range userAccount.Balances {
		if v.Code == "" && decimal.RequireFromString(v.Balance).LessThan(decimal.RequireFromString(os.Getenv("ACCOUNT_RECOVERY_MINIMUM_BALANCE"))) {
			return &tErrors.CustomError{Param: "username", Err: "error primary wallet needs funding", ErrMessage: fmt.Sprintf("Primary wallet needs minimum of 50 %v to proceed.", os.Getenv("NATIVE_ASSET_CODE"))}

		}
	}

	//get the recovery keypair
	recoveryAddress := bc.GetRecoveryAccountAddress(user.Username, user.PublicKey)
	if len(recoveryAddress) == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error generating recovery address", ErrMessage: "Could not generate valid recovery address for account."}

	}
	//activate recovery Address and add signer to primary key
	_, e = userBc.GetBlockchainAccountDetail(recoveryAddress)
	if e != nil {
		// if e.Error() == "error-blockchain-account-not-activated" {
		// 	//activate account
		// 	ops = append(ops, &txnbuild.CreateAccount{
		// 		Destination:   recoveryAddress,
		// 		Amount:        os.Getenv("RECOVERY_SIGNER_ACTIVATION_AMOUNT"),
		// 		SourceAccount: user.PublicKey,
		// 	})
		// 	messages = append(messages, fmt.Sprintf("%v %v will be deducted from your wallet [%v] to activate your unique recovery key on the blockchain.", os.Getenv("RECOVERY_SIGNER_ACTIVATION_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE"), user.Username))

		// }
		return &tErrors.CustomError{Param: "username", Err: "error feature not enabled previously", ErrMessage: fmt.Sprintf("Feature not enabled previously on account [%v].", user.Username)}

	}
	//  else {
	// 	//account already active
	// 	ops = append(ops, &txnbuild.Payment{
	// 		Destination:   recoveryAddress,
	// 		Amount:        "6",
	// 		Asset:         txnbuild.NativeAsset{},
	// 		SourceAccount: user.PublicKey,
	// 	})
	// }

	//check for all wallets
	wallets := user.GetAllWallets(gc)
	bulkPaymentSignerKeyPairs, marketMakingSignerKeyPairs := make([]*keypair.Full, 0), make([]*keypair.Full, 0)

	if len(wallets) > 0 {
		//has subwallets other than the primary wallet, which has already be added to the ops
		for _, w := range wallets {

			if w.SharedAccessEnabled == 1 {
				if !w.HasViewOnlyAccess(gc) {
					// wallet does not have only view-only access, so we need to add it to the multiaccessWallets list
					multiAccessWallets = append(multiAccessWallets, w)
					continue
				}
			}

			_, e = userBc.GetBlockchainAccountDetail(w.ID)
			if e != nil {

				continue

			}

			// if userBc.SignerIsValid(w.ID, recoveryAddress) {
			// 	//recovery a signer to the wallet. remove it
			// 	ops = append(ops, &txnbuild.SetOptions{
			// 		Signer: &txnbuild.Signer{
			// 			Address: recoveryAddress,
			// 			Weight:  0,
			// 		},
			// 		SourceAccount: w.ID,
			// 	})
			// }
			if w.WalletType == 0 || w.WalletType == 1 {
				if userBc.SignerIsValid(w.ID, recoveryAddress) {
					//recovery a signer to the wallet. remove it
					ops = append(ops, &txnbuild.SetOptions{
						Signer: &txnbuild.Signer{
							Address: recoveryAddress,
							Weight:  0,
						},
						SourceAccount: w.ID,
					})
				}
			}
			if w.WalletType == 2 {
				mm, _ := bc.MarketMakingSignerKeypair(user.Username, w.ID)
				marketMakingSignerKeyPairs = append(marketMakingSignerKeyPairs, mm)
				if userBc.SignerIsValid(w.ID, recoveryAddress) {
					//recovery a signer to the wallet. remove it
					ops = append(ops, &txnbuild.SetOptions{
						Signer: &txnbuild.Signer{
							Address: recoveryAddress,
							Weight:  0,
						},
						SourceAccount: w.ID,
					})
				}
			}
			if w.WalletType == 3 {
				bp, _ := bc.BulkPaymentSignerKeypair(user.Username, w.ID)
				bulkPaymentSignerKeyPairs = append(bulkPaymentSignerKeyPairs, bp)
				if userBc.SignerIsValid(w.ID, recoveryAddress) {
					//recovery a signer to the wallet. remove it
					ops = append(ops, &txnbuild.SetOptions{
						Signer: &txnbuild.Signer{
							Address: recoveryAddress,
							Weight:  0,
						},
						SourceAccount: w.ID,
					})
				}
			}
		}

		{ //add fee for transaction

		}
	}
	var xdrBase64 string
	payload.Messages = messages
	payload.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	// oldTrx := payload.Transaction

	// generate new transaction

	memo := "Disabling Account Recovery"
	log.Println("[DisableAccountRecovery] Memo:", memo)
	if len(ops) == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error no operations to perform", ErrMessage: "Could not find any operations to perform for this action."}

	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &userAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              2000,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
			Memo:                 txnbuild.MemoText(memo),
		},
	)

	if err != nil {
		log.Println("[DisableAccountRecovery]error constructing transaction ", err)
		return &tErrors.ErrorTemporaryServerError{}
	}
	if len(bulkPaymentSignerKeyPairs) > 0 {
		tx, err = tx.Sign(gc.BantuNetworkPassphrase, bulkPaymentSignerKeyPairs...)
		if err != nil {
			log.Println("[DisableAccountRecovery] error signning with bulkPaymentSignerKeyPairs", err)
			return &tErrors.ErrorTemporaryServerError{}
		}
	}
	if len(marketMakingSignerKeyPairs) > 0 {
		tx, err = tx.Sign(gc.BantuNetworkPassphrase, marketMakingSignerKeyPairs...)
		if err != nil {
			log.Println("[DisableAccountRecovery] error signning with marketMakingSignerKeyPairs", err)
			return &tErrors.ErrorTemporaryServerError{}
		}
	}
	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[DisableAccountRecovery] error getting txn base64", err)
		return &tErrors.ErrorTemporaryServerError{}
	}
	if len(payload.Transaction) == 0 || len(payload.TransactionSignature) == 0 {
		payload.Transaction = xdrBase64

		return nil
	}

	if len(payload.TransactionSignature) == 0 {
		return &tErrors.CustomError{Param: "TransactionSignature", Err: "error transaction signature is required", ErrMessage: "Transaction signature is required."}
	}
	//submit to blockchain

	txnHash, err := network.SubmitXdrWithSignature(client, user.PrimarySigner, xdrBase64, payload.TransactionSignature)
	if err != nil {
		logDiscordFailedRecovery(fmt.Sprintf("Error submitting disable account recovery [%+v] transaction: %s", payload, err.Error()))
		return
	}
	payload.TransactionID = txnHash
	dbtx.Commit()
	if err == nil && len(multiAccessWallets) > 0 {
		log.Printf("Skipped wallets %+v\n", multiAccessWallets)
	}
	return nil

}

func DoAccountRecovery(user *userModels.User, payload *userModels.AccountRecoveryRequest, gc *sharedconfig.GlobalConfig) (multiAccessWallets []userModels.UserWallet, err error) {
	var e error
	client := gc.BantuExpansionClient
	multiAccessWallets = make([]userModels.UserWallet, 0)
	ops := make([]txnbuild.Operation, 0)
	messages := make([]string, 0)
	payload.Messages = make([]string, 0)

	if len(payload.NewSignerPublicKey) != 56 {
		return multiAccessWallets, &tErrors.CustomError{Param: "newSignerPublicKey", Err: "error invalid new signer public key.", ErrMessage: "Invalid new signer public key."}
	}
	{
		// PARSE SIGNER KEY
		_, e := keypair.ParseAddress(payload.NewSignerPublicKey)
		if e != nil {
			return multiAccessWallets, &tErrors.CustomError{Param: "newSignerPublicKey", Err: "error invalid new signer public key.", ErrMessage: "Invalid new signer public key."}
		}
	}
	if _, e := usersDB.GetUser(payload.NewSignerPublicKey, gc.DB, gc); e == nil {
		return multiAccessWallets, &tErrors.CustomError{Param: "newSignerPublicKey", Err: "error new signer public key already in use.", ErrMessage: "The new signer public key is already in use on another account."}
	}
	if user.AccountRecoveryEnabled == 0 {
		return multiAccessWallets, &tErrors.CustomError{Param: "username", Err: "error account recovery not enabled.", ErrMessage: "Account recovery not enabled."}
	}
	dbtx := gc.DB.Begin()
	defer dbtx.Rollback()

	if !ValidateSecurityAnswers(user, payload.SecurityAnswers, gc) {
		return multiAccessWallets, &tErrors.CustomError{Param: "username", Err: "error invalid security answers", ErrMessage: "Answers to the security questions are invalid."}
	}
	if CheckAccountRecoveryEmailOTP(user, payload.EmailOTP, gc.DB) != nil {
		return multiAccessWallets, &tErrors.CustomError{Param: "username", Err: "error invalid email otp", ErrMessage: "Email OTP is invalid."}
	}

	var userAccount horizon.Account
	if userAccount, e = userBc.GetBlockchainAccountDetail(user.PublicKey); e != nil {
		if e.Error() == "error-blockchain-account-not-activated" {
			return multiAccessWallets, &tErrors.CustomError{Param: "username", Err: "error primary account not yet activated", ErrMessage: fmt.Sprintf("Primary account is not yet activated. Please send upto 50 %v to the primary wallet to continue.", os.Getenv("NATIVE_ASSET_CODE"))}
		}
		log.Printf("[DisableAccountRecovery] Error on blockchain validating primary account %v\n", e)
		return multiAccessWallets, &tErrors.ErrorTemporaryServerError{}

	}
	for _, v := range userAccount.Balances {
		if v.Code == "" && decimal.RequireFromString(v.Balance).LessThan(decimal.RequireFromString(os.Getenv("ACCOUNT_RECOVERY_MINIMUM_BALANCE"))) {
			return multiAccessWallets, &tErrors.CustomError{Param: "username", Err: "error primary wallet needs funding", ErrMessage: fmt.Sprintf("Primary wallet needs minimum of %v %v to proceed.", os.Getenv("ACCOUNT_RECOVERY_MINIMUM_BALANCE"), os.Getenv("NATIVE_ASSET_CODE"))}

		}
	}

	//get the recovery keypair
	recoveryKeyPair, _ := bc.RecoveryAccountKeypair(user.Username, user.PublicKey)
	recoveryAddress := recoveryKeyPair.Address()
	if len(recoveryAddress) == 0 {
		return multiAccessWallets, &tErrors.CustomError{Param: "username", Err: "error generating recovery address", ErrMessage: "Could not generate valid recovery address for account."}

	}
	//activate new signer Address and add signer to primary key
	_, e = userBc.GetBlockchainAccountDetail(payload.NewSignerPublicKey)
	if e != nil {
		if e.Error() == "error-blockchain-account-not-activated" {
			//activate account
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   payload.NewSignerPublicKey,
				Amount:        os.Getenv("RECOVERY_SIGNER_ACTIVATION_AMOUNT"),
				SourceAccount: user.PublicKey,
			})
			messages = append(messages, fmt.Sprintf("%v %v will be deducted from your wallet [%v] to activate your new signer key on the blockchain.", os.Getenv("RECOVERY_SIGNER_ACTIVATION_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE"), user.Username))

		} else {
			return multiAccessWallets, &tErrors.ErrorTemporaryServerError{}
		}

	}

	if !userBc.SignerIsValid(user.PublicKey, payload.NewSignerPublicKey) {
		//newsigner not a signer to the primary wallet.
		ops = append(ops, &txnbuild.SetOptions{
			Signer: &txnbuild.Signer{
				Address: payload.NewSignerPublicKey,
				Weight:  1,
			},
			SourceAccount: user.PublicKey,
		})
	}

	//check for subwallets
	wallets := user.GetAllWallets(gc)
	if len(wallets) > 1 {
		//has subwallets other than the primary wallet, which has already be added to the ops

		for _, w := range wallets {
			if w.Alias == user.Username {
				//skip primary wallet
				continue
			}
			if w.Alias != user.Username {
				if w.SharedAccessEnabled == 1 {
					if !w.HasViewOnlyAccess(gc) {
						// wallet does not have only view-only access, so we need to add it to the multiaccessWallets list
						multiAccessWallets = append(multiAccessWallets, w)
						continue
					}
				}
				//a subwallet, not primary wallet
				//activate subwallet Address and add signer to subwallet key
				_, e = userBc.GetBlockchainAccountDetail(w.ID)
				if e != nil {
					if e.Error() == "error-blockchain-account-not-activated" {
						//activate account
						ops = append(ops, &txnbuild.CreateAccount{
							Destination:   w.ID,
							Amount:        os.Getenv("RECOVERY_SIGNER_ACTIVATION_AMOUNT"),
							SourceAccount: user.PublicKey,
						})
						messages = append(messages, fmt.Sprintf("%v %v will be deducted from your wallet [%v] to activate your subwallet [%v] on the blockchain.", os.Getenv("RECOVERY_SIGNER_ACTIVATION_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE"), user.Username, w.Alias))

					}
				}
				if w.WalletType == 0 || w.WalletType == 1 {
					if !userBc.SignerIsValid(w.ID, payload.NewSignerPublicKey) {
						//recovery not a signer to the sub wallet. add it
						ops = append(ops, &txnbuild.SetOptions{
							Signer: &txnbuild.Signer{
								Address: payload.NewSignerPublicKey,
								Weight:  1,
							},
							SourceAccount: w.ID,
						})
					}
				}

				if w.WalletType == 2 || w.WalletType == 3 {
					if !userBc.SignerIsValid(w.ID, payload.NewSignerPublicKey) {
						//recovery not a signer to the sub wallet. add it
						ops = append(ops, &txnbuild.SetOptions{
							Signer: &txnbuild.Signer{
								Address: payload.NewSignerPublicKey,
								Weight:  1,
							},
							SourceAccount: w.ID,
						})
					}
				}

				if userBc.SignerIsValid(w.ID, user.PrimarySigner) {
					//former signer exists. remove it
					ops = append(ops, &txnbuild.SetOptions{
						Signer: &txnbuild.Signer{
							Address: user.PrimarySigner,
							Weight:  0,
						},
						SourceAccount: w.ID,
					})
				}
			}

		}

	}

	{ //add fee for transaction

	}
	//last operation
	if payload.DisableOldSignerFromPrimaryWallet == 1 {

		if userBc.SignerIsValid(user.PublicKey, user.PrimarySigner) {
			//remove old signer directive is enabled. remove old signer
			if user.PublicKey == user.PrimarySigner {
				//it is the master key you need to disable
				masterWeight := txnbuild.Threshold(0)
				ops = append(ops, &txnbuild.SetOptions{
					MasterWeight:  &masterWeight,
					SourceAccount: user.PublicKey,
				})
			} else {
				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: user.PrimarySigner,
						Weight:  0,
					},
					SourceAccount: user.PublicKey,
				})
			}

			messages = append(messages, "You have chosen to remove old signer from your account. This will remove your previous signer and it will not be able to authorize any more transactions on your account ever again.")
		}
	}

	if payload.Commit == 0 {
		//initial request
		payload.Messages = messages
		return multiAccessWallets, nil
	}
	//commit request
	if payload.Commit == 1 {
		// generate new transaction
		user.LastRecoveredAccountOn = time.Now().UTC()
		user.PrimarySigner = payload.NewSignerPublicKey
		dbErr := dbtx.Save(user).Error
		if dbErr != nil {
			log.Println("[DoAccountRecovery]error saving user database status ", err)
			return multiAccessWallets, &tErrors.ErrorTemporaryServerError{}
		}
		for i, v := range wallets {
			v.Signer = payload.NewSignerPublicKey
			wallets[i] = v
		}
		dbErr = dbtx.Save(&wallets).Error
		if dbErr != nil {
			log.Println("[DoAccountRecovery]error saving wallet signers database status ", err)
			return multiAccessWallets, &tErrors.ErrorTemporaryServerError{}
		}
		RemoveAccountRecoveryEmailOTP(user, payload.EmailOTP, dbtx)
		memo := "Recovering Account"
		log.Println("[DoAccountRecovery] Memo:", memo)

		if len(ops) == 0 {
			return multiAccessWallets, &tErrors.CustomError{Param: "username", Err: "error no operations to perform", ErrMessage: "Could not find any operations to perform for this action."}

		}
		tx, err := txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        &userAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
				Memo:                 txnbuild.MemoText(memo),
			},
		)

		if err != nil {
			log.Println("[DoAccountRecovery]error constructing transaction ", err)
			return multiAccessWallets, &tErrors.ErrorTemporaryServerError{}
		}

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), recoveryKeyPair)
		if err != nil {
			log.Println("[DoAccountRecovery]error signing transaction ", err)
			return multiAccessWallets, &tErrors.ErrorTemporaryServerError{}
		}
		resp, err := client.SubmitTransaction(tx)
		if err != nil {
			log.Println("[DoAccountRecovery]error submitting transaction ", err)
			return multiAccessWallets, &tErrors.ErrorTemporaryServerError{}
		}
		log.Println("[DoAccountRecovery]successfully submitted transaction", resp)
		//do database operation
		dbtx.Commit()
		payload.TransactionID = resp.Hash
		return multiAccessWallets, nil
	}
	return multiAccessWallets, &tErrors.ErrorTemporaryServerError{}
}

func DoInactiveAccountRecover(subjectUser *userModels.User, payload *userModels.InactiveAccountRecoveryRequest, gc *sharedconfig.GlobalConfig) (userInfo userModels.UserInfo, err error) {
	var e error
	payload.NewSignerPublicKey = strings.ToUpper(payload.NewSignerPublicKey)
	answers := payload.SecurityAnswers
	if len(payload.NewSignerPublicKey) != 56 {
		return userInfo, &tErrors.CustomError{Param: "newSignerPublicKey", Err: "error invalid new signer public key.", ErrMessage: "Invalid new signer public key."}
	}
	{
		// PARSE SIGNER KEY
		_, e := keypair.ParseAddress(payload.NewSignerPublicKey)
		if e != nil {
			return userInfo, &tErrors.CustomError{Param: "newSignerPublicKey", Err: "error invalid new signer public key.", ErrMessage: "Invalid new signer public key."}
		}
	}

	if _, e := usersDB.GetUser(payload.NewSignerPublicKey, gc.DB, gc); e == nil {
		return userInfo, &tErrors.CustomError{Param: "newSignerPublicKey", Err: "error new signer public key already in use.", ErrMessage: "The new signer public key is already in use on another account."}
	}

	if _, e := usersDB.GetUserFromPrimarySigner(payload.NewSignerPublicKey, gc.DB, gc); e == nil {
		return userInfo, &tErrors.CustomError{Param: "newSignerPublicKey", Err: "error new signer public key already in use.", ErrMessage: "The new signer public key is already in use on another account."}
	}

	if subjectUser.HasSecurityQuestions == 1 {
		log.Println("[DoInactiveAccountRecover] error account has security question enabled")

		return userInfo, &tErrors.CustomError{Param: "username", Err: "error option not allowed", ErrMessage: "Account not qualified to use this option. This user is not qualified to use this option of recovery. Please use wallet recovery option."}
	}

	if subjectUser.AccountRecoveryEnabled == 1 {

		log.Println("[DoInactiveAccountRecover] error account recovery enabled")

		return userInfo, &tErrors.CustomError{Param: "username", Err: "error option not allowed", ErrMessage: "Account not qualified to use this option. This user is not qualified to use this option of recovery. Please use wallet recovery option."}
	}

	if subjectUser.PublicKey != subjectUser.PrimarySigner {
		log.Println("[DoInactiveAccountRecover] error primary signer and main public key does not match")

		return userInfo, &tErrors.CustomError{Param: "username", Err: "error option not allowed", ErrMessage: "Account not qualified to use this option. This user is not qualified to use this option of recovery. Please use wallet recovery option."}
	}
	wallets := subjectUser.UserWallets
	_, err = userBc.GetBlockchainAccountDetail(subjectUser.PublicKey)
	if err == nil {
		//account already active
		log.Println("[DoInactiveAccountRecover] error account is already activated")
		return userInfo, &tErrors.CustomError{Param: "username", Err: "error option not allowed", ErrMessage: "Account not qualified to use this option. This user is not qualified to use this option of recovery. Please use wallet recovery option."}
	} else {
		if err.Error() != "error-blockchain-account-not-activated" {
			//other blockchain error
			log.Println("[DoInactiveAccountRecover] error other blockchain error: ", err)

			return userInfo, err
		}
	}

	if CheckAccountRecoveryEmailOTP(subjectUser, payload.EmailOTP, gc.DB) != nil {
		return userInfo, &tErrors.CustomError{Param: "username", Err: "error invalid email otp", ErrMessage: "Email OTP is invalid."}
	}

	if len(answers.A1) == 0 || len(answers.A2) == 0 || len(answers.A3) == 0 || answers.Q1 == 0 || answers.Q2 == 0 || answers.Q3 == 0 {
		return userInfo, &tErrors.CustomError{Param: "securityAnswers", Err: "Questions-or-Answers must be 3", ErrMessage: "Questions/Answers must be 3"}
	}

	dbtx := gc.DB.Begin()
	defer dbtx.Rollback()

	// if !ValidateSecurityAnswers(&subjectUser, payload.SecurityAnswers, gc) {
	// 	return userInfo, &tErrors.CustomError{Param: "username", Err: "error invalid security answers", ErrMessage: "Answers to the security questions are invalid."}
	// }
	e = dbtx.Delete(&wallets).Error
	if e != nil {
		// error saving security questions
		log.Printf("[DoInactiveAccountRecover] error removing existing user wallet data for %v. error: %v\n", subjectUser.Username, e)
		return userInfo, &tErrors.ErrorTemporaryServerError{}
	}
	subjectUser.PrimarySigner = payload.NewSignerPublicKey
	subjectUser.PublicKey = payload.NewSignerPublicKey
	subjectUser.HasSecurityQuestions = 1
	subjectUser.UserWallets = make([]userModels.UserWallet, 0)
	subjectUser.BuildPrimaryWallet()

	err = SaveUserSecurityQuestions(subjectUser, answers, dbtx)
	if err != nil {
		// error saving security questions
		return userInfo, err
	}

	e = dbtx.Save(subjectUser).Error
	if e != nil {
		// error saving security questions
		log.Printf("[DoInactiveAccountRecover] error saving user data for %v. error: %v\n", subjectUser.Username, e)
		return userInfo, &tErrors.ErrorTemporaryServerError{}
	}
	dbtx.Commit()
	RemoveAccountRecoveryEmailOTP(subjectUser, payload.EmailOTP, dbtx)

	return GetUserInfo(subjectUser.ID, payload.NewSignerPublicKey, gc)

}

func logDiscordFailedRecovery(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}
