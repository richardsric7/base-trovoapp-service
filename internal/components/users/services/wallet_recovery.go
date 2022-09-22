package users

import (
	"fmt"
	"log"
	"os"
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
	if user.HasSecretQuestions == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error secret answers not set", ErrMessage: "Secret answers has not been set for this account."}
	}
	if user.AccountRecoveryEnabled == 1 {
		return &tErrors.CustomError{Param: "username", Err: "error account recovery already enabled.", ErrMessage: "Account recovery already enabled."}
	}
	wallet, _ := userModels.UserWalletID(user.PublicKey).GetWallet(gc.DB)
	if wallet.ManagedAccessEnabled == 1 && !WalletHasViewOnlyAccess(&wallet, gc) {
		return &tErrors.CustomError{Param: "username", Err: "error primary wallet has shared access.", ErrMessage: "Primary wallet has shared access enabled! Only primary wallets without shared access or with view only shared access can participate in account recovery at this time."}
	}
	// if !ValidateSecretAnswers(user, answer, gc) {
	// 	return &tErrors.CustomError{Param: "username", Err: "error invalid secret answers", ErrMessage: "Answers to the secret questions are invalid."}
	// }

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
	if len(wallets) > 1 {
		//has subwallets other than the primary wallet, which has already be added to the ops
		for _, w := range wallets {
			if w.Alias != user.Username {
				if w.ManagedAccessEnabled == 1 {
					if !WalletHasViewOnlyAccess(&w, gc) {
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
			}
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

			{ //add fee for transaction

			}

		}
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
	if user.HasSecretQuestions == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error secret answers not set", ErrMessage: "Secret answers has not been set for this account."}
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

	if !ValidateSecretAnswers(user, payload.SecretAnswers, gc) {
		return &tErrors.CustomError{Param: "username", Err: "error invalid secret answers", ErrMessage: "Answers to the secret questions are invalid."}
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
	if len(wallets) > 0 {
		//has subwallets other than the primary wallet, which has already be added to the ops
		for _, w := range wallets {

			if w.ManagedAccessEnabled == 1 {
				if !WalletHasViewOnlyAccess(&w, gc) {
					// wallet does not have only view-only access, so we need to add it to the multiaccessWallets list
					multiAccessWallets = append(multiAccessWallets, w)
					continue
				}
			}
			//a subwallet, not primary wallet
			//activate subwallet Address and add signer to subwallet key
			_, e = userBc.GetBlockchainAccountDetail(w.ID)
			if e != nil {
				// if e.Error() == "error-blockchain-account-not-activated" {
				// 	//activate account
				// 	ops = append(ops, &txnbuild.CreateAccount{
				// 		Destination:   w.ID,
				// 		Amount:        os.Getenv("RECOVERY_SIGNER_ACTIVATION_AMOUNT"),
				// 		SourceAccount: user.PublicKey,
				// 	})
				// 	messages = append(messages, fmt.Sprintf("%v %v will be deducted from your wallet [%v] to activate your subwallet [%v] on the blockchain.", os.Getenv("RECOVERY_SIGNER_ACTIVATION_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE"), user.Username, w.Alias))

				// }
				continue

			}

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

			{ //add fee for transaction

			}

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

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[DisableAccountRecovery] error getting txn base64", err)
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
	if _, e := usersDB.GetUser(payload.NewSignerPublicKey, gc.DB); e == nil {
		return multiAccessWallets, &tErrors.CustomError{Param: "newSignerPublicKey", Err: "error new signer public key already in use.", ErrMessage: "The new signer public key is already in use on another account."}
	}
	if user.AccountRecoveryEnabled == 0 {
		return multiAccessWallets, &tErrors.CustomError{Param: "username", Err: "error account recovery not enabled.", ErrMessage: "Account recovery not enabled."}
	}
	dbtx := gc.DB.Begin()
	defer dbtx.Rollback()

	if !ValidateSecretAnswers(user, payload.SecretAnswers, gc) {
		return multiAccessWallets, &tErrors.CustomError{Param: "username", Err: "error invalid secret answers", ErrMessage: "Answers to the secret questions are invalid."}
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
			if w.Alias != user.Username {
				if w.ManagedAccessEnabled == 1 {
					if !WalletHasViewOnlyAccess(&w, gc) {
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
			}
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

			{ //add fee for transaction

			}

		}
	}

	//last operation

	if userBc.SignerIsValid(user.PublicKey, user.PrimarySigner) && payload.DisableOldSignerFromPrimaryWallet == 1 {
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

func logDiscordFailedRecovery(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}
