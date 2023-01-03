package users

import (
	"fmt"
	"log"
	"os"
	"strings"
	bc "trovo-wallet-api/internal/blockchainalgofuncs"
	userBc "trovo-wallet-api/internal/components/users/blockchain"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ecnepsnai/discord"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

func CreateNewSubWallet(accountOwner *userModels.User, subWalletInfo *userModels.SubWalletInfo, gc *sharedconfig.GlobalConfig) (*userModels.SubWalletInfo, error) {
	client := network.GetBlockchainClient()
	var err error
	var xdrBase64 string
	var subWalletObj userModels.UserWallet

	subWalletInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	subWalletInfo.SubWalletMustSign = 1
	if len(subWalletInfo.ChannelAccount) == 56 {
		//generate xdr for channel account
		xdrBase64, subWalletObj, err = generateSubWalletXdrWithChannelAccount(accountOwner, subWalletInfo, gc, client)
		if err != nil {
			log.Printf("[CreateNewSubWallet] create sub [%v] for [%v] generateSubWalletXdrWithChannelAccount error:[%v] \n", subWalletInfo.PublicKey, accountOwner.Username, err)
			return subWalletInfo, err
		}
	} else {
		xdrBase64, subWalletObj, err = generateSubWalletXdr(accountOwner, subWalletInfo, gc, client)
		if err != nil {
			log.Printf("[CreateNewSubWallet] create sub [%v] for [%v] generateSubWalletXdr error:[%v] \n", subWalletInfo.PublicKey, accountOwner.Username, err)
			return subWalletInfo, err
		}
	}

	oldTransaction := subWalletInfo.Transaction

	subWalletInfo.Transaction = xdrBase64

	if len(subWalletInfo.PrimarySignature) == 0 || len(subWalletInfo.SubWalletSignature) == 0 {
		//it has not been signed before
		return subWalletInfo, nil
	}

	if oldTransaction != subWalletInfo.Transaction {
		err = &tErrors.CustomError{
			Param:      "transaction",
			Err:        "error-transaction-mismatch",
			ErrMessage: "Transaction mismatch",
		}
		return subWalletInfo, err

	}

	if len(subWalletInfo.ChannelAccountSignature) == 0 && len(subWalletInfo.ChannelAccount) == 56 {
		err = &tErrors.CustomError{
			Param:      "transaction",
			Err:        "error-transaction-mismatch",
			ErrMessage: "signature for channel account does not validate",
		}
		return subWalletInfo, err

	}
	//begin a database transaction here
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()

	//create the data to be sure it goes through
	errDBTX := dbTX.Create(&subWalletObj).Error
	if errDBTX != nil {
		//unable to save sub wallet. abort
		log.Printf("[CreateNewSubWallet] by [%v] for [%v] Error saving subwallet error:[%v] \n", accountOwner.Username, subWalletInfo.PublicKey, errDBTX)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-saving-subwallet",
			ErrMessage: "There is an error saving sub-wallet. Please, try again later.",
		}
		return subWalletInfo, err
	}

	if len(subWalletInfo.ChannelAccountSignature) > 0 && len(subWalletInfo.ChannelAccount) == 56 {
		txnHash, err := SubmitSubWalletXdrForChannelAccountWithSignature(client, accountOwner.PublicKey, subWalletInfo.PublicKey, subWalletInfo.ChannelAccount, xdrBase64, subWalletInfo.PrimarySignature, subWalletInfo.SubWalletSignature, subWalletInfo.ChannelAccountSignature, subWalletInfo.SubWalletMustSign)
		if err != nil {
			log.Printf("[CreateNewSubWallet] by [%v] for [%v] SubmitSubWalletXdrForChannelAccountWithSignature error:[%v] \n", accountOwner.Username, subWalletInfo.PublicKey, err)
			return subWalletInfo, err
		}
		subWalletInfo.TransactionID = txnHash
		dbTX.Commit()

		{
			//send to monitoring service
			trackPublicKey := userModels.TrackedPublicKey{
				PublicKey: subWalletInfo.PublicKey,
			}
			errTrack := gc.RoachDB.Create(&trackPublicKey).Error
			if errTrack != nil {
				//if tracking of public key fails, then payment history generation service will pick it up and do justice to it
				discord.Say(fmt.Sprintf("[CreateNewSubWallet] tracking public key for payment history failed for user:%v, with DB Error:%v\n\n\nFailedData:%+v", accountOwner.Username, errTrack, subWalletInfo))

			}
		}
	} else {
		txnHash, err := SubmitSubWalletXdrWithSignature(client, accountOwner.PublicKey, subWalletInfo.PublicKey, xdrBase64, subWalletInfo.PrimarySignature, subWalletInfo.SubWalletSignature, subWalletInfo.SubWalletMustSign)
		if err != nil {
			log.Printf("[CreateNewSubWallet] by [%v] for [%v] SubmitSubwalletXdrWithSignature error:[%v] \n", accountOwner.Username, subWalletInfo.PublicKey, err)
			return subWalletInfo, err
		}
		subWalletInfo.TransactionID = txnHash
		dbTX.Commit()
		{
			//send to monitoring service
			trackPublicKey := userModels.TrackedPublicKey{
				PublicKey: subWalletInfo.PublicKey,
			}
			errTrack := gc.RoachDB.Create(&trackPublicKey).Error
			if errTrack != nil {
				//if tracking of public key fails, then payment history generation service will pick it up and do justice to it
				discord.Say(fmt.Sprintf("[CreateNewSubWallet] tracking public key for payment history failed for user:%v, with DB Error:%v\n\n\nFailedData:%+v", accountOwner.Username, errTrack, subWalletInfo))

			}
		}

	}
	accountOwner.InvalidateUserCache(gc)

	return subWalletInfo, nil
}

func generateSubWalletXdr(accountOwner *userModels.User, subWalletInfo *userModels.SubWalletInfo, gc *sharedconfig.GlobalConfig, client *horizonclient.Client) (xdrbase64 string, subWalletObj userModels.UserWallet, err error) {
	ops := make([]txnbuild.Operation, 0)
	subWalletInfo.Messages = make([]string, 0)
	var activationAmount = decimal.NewFromFloat(6)
	var minBalance = decimal.NewFromFloat(3.0)
	dab := strings.Split(os.Getenv("DOLLAR_ASSET"), ":")
	dollarAsset := txnbuild.CreditAsset{Code: dab[0], Issuer: dab[1]}
	if len(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT")) > 0 {
		activationAmount = decimal.RequireFromString(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("ISSUING_SUB_WALLET_ACTIVATION_AMOUNT")) > 0 && subWalletInfo.WalletType == 1 {
		activationAmount = decimal.RequireFromString(os.Getenv("ISSUING_SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("MM_SUB_WALLET_ACTIVATION_AMOUNT")) > 0 && subWalletInfo.WalletType == 2 {
		activationAmount = decimal.RequireFromString(os.Getenv("MM_SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("BULKPAYMENT_SUB_WALLET_ACTIVATION_AMOUNT")) > 0 && subWalletInfo.WalletType == 3 {
		activationAmount = decimal.RequireFromString(os.Getenv("BULKPAYMENT_SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("WALLET_MINIMUM_BALANCE")) > 0 {
		minBalance = decimal.RequireFromString(os.Getenv("WALLET_MINIMUM_BALANCE"))
	}
	{
		//check if the sub-wallet passes the validation
		subWalletObj, err = accountOwner.BuildNewSubWallet(subWalletInfo.PublicKey, subWalletInfo.WalletTag, subWalletInfo.WalletDescription, subWalletInfo.WalletType, gc)
		if err != nil {
			log.Printf("[generateSubWalletXdr] by [%v] for [%v] BuildNewSubWallet error:[%v] \n", accountOwner.Username, subWalletInfo.PublicKey, err)
			return "", subWalletObj, err
		}

		//set the subwallet suggested alias
		subWalletInfo.Alias = subWalletObj.Alias

	}

	//check if it is first call to create sub-wallet

	//populate the subwallet Info and generate the transaction

	//check if primary account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	primaryAccountExists, _, primaryAccountNativeBalance, _, primarySourceAccount, errAct := network.BlockchainAccountProperties(client, accountOwner.PublicKey, nativeAsset)
	if errAct != nil {
		log.Printf("[generateSubWalletXdr] by [%v] for [%v] Primary Account Properties error:[%v] \n", accountOwner.Username, subWalletInfo.PublicKey, errAct)

		return "", subWalletObj, errAct
	}
	if !primaryAccountExists || (primaryAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
		log.Printf("[generateSubWalletXdr] by [%v] for [%v] PrimaryAccount underfunded \n", accountOwner.Username, subWalletInfo.PublicKey)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-primary-account-underfunded",
			ErrMessage: fmt.Sprintf("Primary account does not have enough %v balance to create sub-wallet", os.Getenv("NATIVE_ASSET_CODE")),
			Code:       404,
		}
		return "", subWalletObj, err
	}
	var walletSigner *keypair.Full
	if subWalletInfo.WalletType == 2 {
		walletSigner, _ = bc.MarketMakingSignerKeypair(accountOwner.Username, subWalletInfo.PublicKey)

	}
	if subWalletInfo.WalletType == 3 {
		walletSigner, _ = bc.BulkPaymentSignerKeypair(accountOwner.Username, subWalletInfo.PublicKey)

	}

	subWalletAccountExists, _, subWalletAccountNativeBalance, _, subWalletAccountObject, _ := network.BlockchainAccountProperties(client, subWalletInfo.PublicKey, nativeAsset)
	if !subWalletAccountExists {
		//if subwallet is not activated
		ops = append(ops, &txnbuild.CreateAccount{
			Destination:   subWalletInfo.PublicKey,
			Amount:        activationAmount.String(),
			SourceAccount: accountOwner.PublicKey,
		})
		//if a minting wallet do not create trustline
		if subWalletInfo.WalletType != 1 {
			//enable dollar asset if not minting wallet

			ops = append(ops, &txnbuild.ChangeTrust{
				Line:          txnbuild.ChangeTrustAssetWrapper{Asset: dollarAsset},
				Limit:         "900000000000",
				SourceAccount: subWalletInfo.PublicKey,
			})
		}

		//build transaction that will activate the subwallet from the primary wallet
		if subWalletInfo.WalletType == 0 || subWalletInfo.WalletType == 1 {

			//after creation, it now exists with enough balance to add primary wallet as signer
			ops = append(ops, &txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: accountOwner.PrimarySigner,
					Weight:  1,
				},
				SourceAccount: subWalletInfo.PublicKey,
			})
		}
		//make it custodial
		if subWalletInfo.WalletType == 2 || subWalletInfo.WalletType == 3 {

			signerExists, _, _, _, _, _ := network.BlockchainAccountProperties(client, walletSigner.Address(), nativeAsset)
			if !signerExists {
				//activate signer
				ops = append(ops, &txnbuild.CreateAccount{
					Destination:   walletSigner.Address(),
					Amount:        os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT"),
					SourceAccount: accountOwner.PublicKey,
				})

			} else {
				//topup signer
				ops = append(ops, &txnbuild.Payment{
					Destination:   walletSigner.Address(),
					Amount:        os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT"),
					Asset:         nativeAsset,
					SourceAccount: accountOwner.PublicKey,
				})

			}

			//after creation, it now exists with enough balance to add primary wallet as signer
			ops = append(ops, &txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: accountOwner.PrimarySigner,
					Weight:  1,
				},
				SourceAccount: subWalletInfo.PublicKey,
			})
			ops = append(ops, &txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: walletSigner.Address(),
					Weight:  3,
				},
				SourceAccount: subWalletInfo.PublicKey,
			})
		}

	}

	if subWalletAccountExists {
		if subWalletAccountNativeBalance.LessThan(minBalance) {
			ops = append(ops, &txnbuild.Payment{
				Destination:   subWalletInfo.PublicKey,
				Amount:        activationAmount.String(),
				Asset:         nativeAsset,
				SourceAccount: accountOwner.PublicKey,
			})
		}
		if subWalletInfo.WalletType != 1 {
			//enable dollar asset if not minting wallet
			_, trusted, _, _, _, _ := network.BlockchainAccountProperties(client, subWalletInfo.PublicKey, dollarAsset)
			if !trusted {
				ops = append(ops, &txnbuild.ChangeTrust{
					Line:          txnbuild.ChangeTrustAssetWrapper{Asset: dollarAsset},
					Limit:         "900000000000",
					SourceAccount: subWalletInfo.PublicKey,
				})
			}

		}

		//account exists and native balance is less than needed. add 3 native token to the wallet
		if subWalletInfo.WalletType == 0 || subWalletInfo.WalletType == 1 {

			//after topping up, it now has enough balance to add primary wallet as signer if it is not already a signer
			if !accountOwner.SignerIsValidWA(accountOwner.PrimarySigner, subWalletAccountObject) {

				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: accountOwner.PrimarySigner,
						Weight:  1,
					},
					SourceAccount: subWalletInfo.PublicKey,
				})
			} else {
				subWalletInfo.SubWalletMustSign = 0
			}

		}
		//make it custodial
		if subWalletInfo.WalletType == 2 || subWalletInfo.WalletType == 3 {

			signerExists, _, _, _, _, _ := network.BlockchainAccountProperties(client, walletSigner.Address(), nativeAsset)
			if !signerExists {
				ops = append(ops, &txnbuild.CreateAccount{
					Destination:   walletSigner.Address(),
					Amount:        os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT"),
					SourceAccount: accountOwner.PublicKey,
				})

			} else {
				ops = append(ops, &txnbuild.Payment{
					Destination:   walletSigner.Address(),
					Amount:        os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT"),
					Asset:         nativeAsset,
					SourceAccount: accountOwner.PublicKey,
				})

			}

			//after creation, it now exists with enough balance to add primary wallet as signer
			if !userBc.SignerIsValid(subWalletInfo.PublicKey, accountOwner.PrimarySigner) {
				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: accountOwner.PrimarySigner,
						Weight:  1,
					},
					SourceAccount: subWalletInfo.PublicKey,
				})
			}
			if !userBc.SignerIsValid(subWalletInfo.PublicKey, walletSigner.Address()) {
				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: walletSigner.Address(),
						Weight:  3,
					},
					SourceAccount: subWalletInfo.PublicKey,
				})
			}
		}

	}

	//add recovery key if account recovery is enabled
	if accountOwner.AccountRecoveryEnabled == 1 {
		recoveryKeyAddress := bc.GetRecoveryAccountAddress(accountOwner.Username, accountOwner.PublicKey)

		if len(recoveryKeyAddress) == 56 {
			if subWalletInfo.WalletType == 0 || subWalletInfo.WalletType == 1 {
				if !userBc.SignerIsValid(subWalletInfo.PublicKey, recoveryKeyAddress) {
					ops = append(ops, &txnbuild.SetOptions{
						Signer: &txnbuild.Signer{
							Address: recoveryKeyAddress,
							Weight:  1,
						},
						SourceAccount: subWalletInfo.PublicKey,
					})
				}
			}

			if subWalletInfo.WalletType == 2 || subWalletInfo.WalletType == 3 {
				if !userBc.SignerIsValid(subWalletInfo.PublicKey, recoveryKeyAddress) {
					ops = append(ops, &txnbuild.SetOptions{
						Signer: &txnbuild.Signer{
							Address: recoveryKeyAddress,
							Weight:  3,
						},
						SourceAccount: subWalletInfo.PublicKey,
					})
				}
			}

		}

	}
	serviceFee, e := decimal.NewFromString(os.Getenv("SUBWALLET_FEE_AMOUNT_USD"))
	if e != nil {
		serviceFee = decimal.Zero
	}
	if serviceFee.IsPositive() {
		//add fees if enabled.
		//process service fee
		subWalletInfo.FeeAmount = serviceFee.String()
		subWalletInfo.FeeCode = os.Getenv("SUBWALLET_FEE_ASSET_CODE")

		if len(os.Getenv("SUBWALLET_FEE_ASSET_ISSUER")) == 56 {
			ops = append(ops, &txnbuild.Payment{
				Destination:   os.Getenv("SUBWALLET_FEE_ADDRESS"),
				Amount:        os.Getenv("SUBWALLET_FEE_AMOUNT_USD"),
				SourceAccount: accountOwner.PublicKey,
				Asset:         txnbuild.CreditAsset{Code: os.Getenv("SUBWALLET_FEE_ASSET_CODE"), Issuer: os.Getenv("SUBWALLET_FEE_ASSET_ISSUER")},
			})
			subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("%v %v will be deducted as service fee.", os.Getenv("SUBWALLET_FEE_AMOUNT_USD"), os.Getenv("SUBWALLET_FEE_ASSET_CODE")))

		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   os.Getenv("SUBWALLET_FEE_ADDRESS"),
				Amount:        os.Getenv("SUBWALLET_FEE_AMOUNT_USD"),
				SourceAccount: accountOwner.PublicKey,
				Asset:         txnbuild.NativeAsset{},
			})
			subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("%v %v will be deducted as service fee.", os.Getenv("SUBWALLET_FEE_AMOUNT_USD"), os.Getenv("NATIVE_ASSET_CODE")))

		}
	}

	if subWalletInfo.WalletType == 0 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("%v %v will be sent from your primary wallet to this new sub-wallet for wallet activation.", activationAmount.String(), os.Getenv("NATIVE_ASSET_CODE")))

	}
	if subWalletInfo.WalletType == 1 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Because this subwallet is designated to be a token minting wallet, %v %v will be deducted from your primary wallet and be used to activate it. Please note that token minting wallets cannot be used to send payments.", activationAmount.String(), os.Getenv("NATIVE_ASSET_CODE")))

	}
	if subWalletInfo.WalletType == 2 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Because this subwallet is designated to be an market making wallet, %v %v will be deducted from your primary wallet and be used to activate it and the custodial signer. Please note that MM wallets cannot be used to send normal payments, but only used for market making.", (activationAmount.Add(decimal.RequireFromString(os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT")))).String(), os.Getenv("NATIVE_ASSET_CODE")))

	}

	if subWalletInfo.WalletType == 3 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Because this subwallet is designated to be an bulk-payment wallet, %v %v will be deducted from your primary wallet and be used to activate it and the custodial signer. Please note that bulk-payment wallets cannot be used to send normal payments, but only be used by internal system to disburse bulk payments on your behalf.", (activationAmount.Add(decimal.RequireFromString(os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT")))).String(), os.Getenv("NATIVE_ASSET_CODE")))

	}
	if subWalletInfo.WalletType == 2 || subWalletInfo.WalletType == 3 {
		ops = append(ops, &txnbuild.SetOptions{
			LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(3)),
			MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(3)),
			HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(3)),
			SourceAccount:   subWalletInfo.PublicKey,
		})
	}
	// if subWalletInfo.WalletType == 0 || subWalletInfo.WalletType == 1 {
	// 	ops = append(ops, &txnbuild.SetOptions{
	// 		LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(1)),
	// 		MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(1)),
	// 		HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(1)),
	// 		SourceAccount:   subWalletInfo.PublicKey,
	// 	})
	// }
	// Construct the transaction that holds the operations to execute on the network
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        primarySourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText("Create Sub Wallet"),
		},
	)
	if err != nil {
		log.Println("[generateSubWalletXdr] error constructing transaction ", err)
		return "", subWalletObj, err
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateSubWalletXdr] error getting txn base64", err)
		return "", subWalletObj, err
	}

	return xdrBase64, subWalletObj, nil

}

func generateSubWalletXdrWithChannelAccount(user *userModels.User, subWalletInfo *userModels.SubWalletInfo, gc *sharedconfig.GlobalConfig, client *horizonclient.Client) (xdrbase64 string, subWalletObj userModels.UserWallet, err error) {
	ops := make([]txnbuild.Operation, 0)
	subWalletInfo.Messages = make([]string, 0)
	var activationAmount = decimal.NewFromFloat(6.0)
	var minBalance = decimal.NewFromFloat(3.0)
	if len(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT")) > 0 {
		activationAmount = decimal.RequireFromString(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("ISSUING_SUB_WALLET_ACTIVATION_AMOUNT")) > 0 && subWalletInfo.WalletType == 1 {
		activationAmount = decimal.RequireFromString(os.Getenv("ISSUING_SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("MM_SUB_WALLET_ACTIVATION_AMOUNT")) > 0 && subWalletInfo.WalletType == 2 {
		activationAmount = decimal.RequireFromString(os.Getenv("MM_SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("BULKPAYMENT_SUB_WALLET_ACTIVATION_AMOUNT")) > 0 && subWalletInfo.WalletType == 3 {
		activationAmount = decimal.RequireFromString(os.Getenv("BULKPAYMENT_SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("WALLET_MINIMUM_BALANCE")) > 0 {
		minBalance = decimal.RequireFromString(os.Getenv("WALLET_MINIMUM_BALANCE"))
	}
	{
		//check if the sub-wallet passes the validation
		subWalletObj, err = user.BuildNewSubWallet(subWalletInfo.PublicKey, subWalletInfo.WalletTag, subWalletInfo.WalletDescription, subWalletInfo.WalletType, gc)
		if err != nil {
			log.Printf("[generateSubWalletXdrWithChannelAccount] by [%v] for [%v] BuildNewSubWallet error:[%v] \n", user.Username, subWalletInfo.PublicKey, err)

			return "", subWalletObj, err
		}
		//set the subwallet suggested alias
		subWalletInfo.Alias = subWalletObj.Alias

	}

	//check if it is first call to create sub-wallet

	//populate the subwallet Info and generate the transaction

	//check if primary account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	primaryAccountExists, _, primaryAccountNativeBalance, _, _, errAct := network.BlockchainAccountProperties(client, user.PublicKey, nativeAsset)
	if errAct != nil {
		log.Printf("[generateSubWalletXdrWithChannelAccount] by [%v] for [%v] Primary Account Error error:[%v] \n", user.Username, subWalletInfo.PublicKey, errAct)

		return "", subWalletObj, errAct
	}
	if !primaryAccountExists || (primaryAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
		log.Printf("[generateSubWalletXdrWithChannelAccount] by [%v] for [%v] Primary Account Underfunded\n", user.Username, subWalletInfo.PublicKey)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-primary-account-underfunded",
			ErrMessage: fmt.Sprintf("Primary account does not have enough %s balance to create sub-wallet", os.Getenv("NATIVE_ASSET_CODE")),
			Code:       404,
		}
		return "", subWalletObj, err
	}
	var walletSigner *keypair.Full
	if subWalletInfo.WalletType == 2 {
		walletSigner, _ = bc.MarketMakingSignerKeypair(user.Username, subWalletInfo.PublicKey)

	}
	if subWalletInfo.WalletType == 3 {
		walletSigner, _ = bc.BulkPaymentSignerKeypair(user.Username, subWalletInfo.PublicKey)

	}

	subWalletAccountExists, _, subWalletAccountNativeBalance, _, subWalletAccountObject, _ := network.BlockchainAccountProperties(client, subWalletInfo.PublicKey, nativeAsset)
	if !subWalletAccountExists {
		//if subwallet is not activated
		//build transaction that will activate the subwallet from the primary wallet
		if subWalletInfo.WalletType == 0 || subWalletInfo.WalletType == 1 {

			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   subWalletInfo.PublicKey,
				Amount:        activationAmount.String(),
				SourceAccount: user.PublicKey,
			})

			//after creation, it now exists with enough balance to add primary wallet as signer
			ops = append(ops, &txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: user.PrimarySigner,
					Weight:  1,
				},
				SourceAccount: subWalletInfo.PublicKey,
			})
		}
		//make it custodial
		if subWalletInfo.WalletType == 2 || subWalletInfo.WalletType == 3 {

			signerExists, _, _, _, _, _ := network.BlockchainAccountProperties(client, walletSigner.Address(), nativeAsset)
			if !signerExists {
				ops = append(ops, &txnbuild.CreateAccount{
					Destination:   walletSigner.Address(),
					Amount:        os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT"),
					SourceAccount: user.PublicKey,
				})

			} else {
				ops = append(ops, &txnbuild.Payment{
					Destination:   walletSigner.Address(),
					Amount:        os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT"),
					Asset:         nativeAsset,
					SourceAccount: user.PublicKey,
				})

			}

			//after creation, it now exists with enough balance to add primary wallet as signer
			ops = append(ops, &txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: user.PrimarySigner,
					Weight:  1,
				},
				SourceAccount: subWalletInfo.PublicKey,
			})
			ops = append(ops, &txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: walletSigner.Address(),
					Weight:  3,
				},
				SourceAccount: subWalletInfo.PublicKey,
			})
		}

	}

	if subWalletAccountExists {
		if subWalletAccountNativeBalance.LessThan(minBalance) {
			ops = append(ops, &txnbuild.Payment{
				Destination:   subWalletInfo.PublicKey,
				Amount:        activationAmount.String(),
				Asset:         nativeAsset,
				SourceAccount: user.PublicKey,
			})
		}
		//account exists and native balance is less than needed. add 3 native token to the wallet
		if subWalletInfo.WalletType == 0 || subWalletInfo.WalletType == 1 {

			//after topping up, it now has enough balance to add primary wallet as signer if it is not already a signer
			if !user.SignerIsValidWA(user.PrimarySigner, subWalletAccountObject) {

				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: user.PrimarySigner,
						Weight:  1,
					},
					SourceAccount: subWalletInfo.PublicKey,
				})
			} else {
				subWalletInfo.SubWalletMustSign = 0
			}

		}
		//make it custodial
		if subWalletInfo.WalletType == 2 || subWalletInfo.WalletType == 3 {

			signerExists, _, _, _, _, _ := network.BlockchainAccountProperties(client, walletSigner.Address(), nativeAsset)
			if !signerExists {
				ops = append(ops, &txnbuild.CreateAccount{
					Destination:   walletSigner.Address(),
					Amount:        os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT"),
					SourceAccount: user.PublicKey,
				})

			} else {
				ops = append(ops, &txnbuild.Payment{
					Destination:   walletSigner.Address(),
					Amount:        os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT"),
					Asset:         nativeAsset,
					SourceAccount: user.PublicKey,
				})

			}

			//after creation, it now exists with enough balance to add primary wallet as signer
			if !userBc.SignerIsValid(subWalletInfo.PublicKey, user.PrimarySigner) {
				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: user.PrimarySigner,
						Weight:  1,
					},
					SourceAccount: subWalletInfo.PublicKey,
				})
			}
			if !userBc.SignerIsValid(subWalletInfo.PublicKey, walletSigner.Address()) {
				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: walletSigner.Address(),
						Weight:  3,
					},
					SourceAccount: subWalletInfo.PublicKey,
				})
			}
		}

	}

	//add recovery key if account recovery is enabled
	if user.AccountRecoveryEnabled == 1 {
		recoveryKeyAddress := bc.GetRecoveryAccountAddress(user.Username, user.PublicKey)

		if len(recoveryKeyAddress) == 56 {
			if !userBc.SignerIsValid(subWalletInfo.PublicKey, recoveryKeyAddress) {
				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: recoveryKeyAddress,
						Weight:  3,
					},
					SourceAccount: subWalletInfo.PublicKey,
				})
			}

		}

	}
	fee := decimal.RequireFromString(os.Getenv("SHARED_ACCESS_FEE_AMOUNT"))
	if !fee.IsZero() {
		//add fees if enabled.
		//process service fee
		if len(os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")) == 56 {
			ops = append(ops, &txnbuild.Payment{
				Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
				Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
				SourceAccount: user.PublicKey,
				Asset:         txnbuild.CreditAsset{Code: os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), Issuer: os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")},
			})
			subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("%v %v will be deducted as service fee.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE")))

		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
				Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
				SourceAccount: user.PublicKey,
				Asset:         txnbuild.NativeAsset{},
			})
			subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("%v %v will be deducted as service fee for creating view only access.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE")))

		}
	}

	if subWalletInfo.WalletType == 0 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("%v %v will be deducted from your primary wallet and be used to activate the sub-wallet.", activationAmount.String(), os.Getenv("NATIVE_ASSET_CODE")))

	}
	if subWalletInfo.WalletType == 1 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Because this subwallet is designated to be an asset issuing wallet, %v %v will be deducted from your primary wallet and be used to activate it. Please note that asset issuing wallets cannot be used to send payments.", activationAmount.String(), os.Getenv("NATIVE_ASSET_CODE")))

	}
	if subWalletInfo.WalletType == 2 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Because this subwallet is designated to be an market making wallet, %v %v will be deducted from your primary wallet and be used to activate it and the custodial signer. Please note that MM wallets cannot be used to send normal payments, but only used for market making.", (activationAmount.Add(decimal.RequireFromString(os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT")))).String(), os.Getenv("NATIVE_ASSET_CODE")))

	}

	if subWalletInfo.WalletType == 3 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Because this subwallet is designated to be an bulk-payment wallet, %v %v will be deducted from your primary wallet and be used to activate it and the custodial signer. Please note that bulk-payment wallets cannot be used to send normal payments, but only be used by internal system to disburse bulk payments on your behalf.", (activationAmount.Add(decimal.RequireFromString(os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT")))).String(), os.Getenv("NATIVE_ASSET_CODE")))

	}
	if subWalletInfo.WalletType == 2 || subWalletInfo.WalletType == 3 {
		ops = append(ops, &txnbuild.SetOptions{
			LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(3)),
			MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(3)),
			HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(3)),
			SourceAccount:   subWalletInfo.PublicKey,
		})
	}

	channelSourceAccountExists, _, channelSourceAccountNativeBalance, _, channelSourceAccount, channelSourceAccountErr := network.BlockchainAccountProperties(client, subWalletInfo.ChannelAccount, txnbuild.NativeAsset{})
	if !channelSourceAccountExists || channelSourceAccountErr != nil || (channelSourceAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
		log.Printf("[generateSubWalletXdrWithChannelAccount] by [%v] for [%v] Channel Account underfunded.\n", user.Username, subWalletInfo.PublicKey)

		err = &tErrors.CustomError{
			Param:      "channelAccount",
			Err:        "error-channel-account-underfunded",
			ErrMessage: fmt.Sprintf("Channel account does not have minimum %v balance required to complete this operation", os.Getenv("NATIVE_ASSET_CODE")),
			Code:       400,
		}
		return "", subWalletObj, err
	}
	// Construct the transaction that holds the operations to execute on the network
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        channelSourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText("Create Sub Wallet"),
		},
	)
	if err != nil {
		log.Println("[generateSubWalletXdrWithChannelAccount] error constructing transaction ", err)
		return "", subWalletObj, err
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateSubWalletXdrWithChannelAccount] error getting txn base64", err)
		return "", subWalletObj, err
	}

	return xdrBase64, subWalletObj, nil

}

func SubmitSubWalletXdrWithSignature(client *horizonclient.Client, ownerPublicKey, subWalletPublicKey string, xdrBase64 string, primarySignature, subWalletSignature string, subWalletMustSign int) (string, error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	gTxn, err := txnbuild.TransactionFromXDR(xdrBase64)

	if err != nil {
		return "", err
	}

	txn, ok := gTxn.Transaction()

	if !ok {
		return "", &tErrors.ErrorInvalidTransaction{}
	}

	{
		//add signature of the primaryWallet to the new transaction Instance
		txn, err = txn.AddSignatureBase64(network.GetBlockchainNetworkPassPhrase(), ownerPublicKey, primarySignature)
		if err != nil {
			log.Println("[SubmitSubwalletXdrWithSignature] Failed to verify primary signature of primaryWallet on [", network.GetBlockchainNetworkPassPhrase(), "] and [", primarySignature, "] for [", xdrBase64, "] and public key ", ownerPublicKey, ", error [", err, "]")

			return "", err
		}

		//add signature of the subWallet to the new transaction Instance
		if subWalletMustSign == 1 {
			txn, err = txn.AddSignatureBase64(network.GetBlockchainNetworkPassPhrase(), subWalletPublicKey, subWalletSignature)
			if err != nil {
				log.Println("[SubmitSubwalletXdrWithSignature] Failed to verify signature of subwallet on [", network.GetBlockchainNetworkPassPhrase(), "] and [", subWalletSignature, "] for [", xdrBase64, "] and public key ", subWalletPublicKey, ", error [", err, "]")

				return "", err
			}
		}
	}

	xdrBase64, err = txn.Base64()

	if err != nil {
		log.Printf("[SubmitSubwalletXdrWithSignature] error converting transaction to base64: %v\n", err)
		return "", err
	}

	// log.Println("signed xdr is " + xdrBase64)

	txnResult, err := client.SubmitTransactionXDR(xdrBase64)

	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			discord.Say(fmt.Sprintf("[SubmitSubwalletXdrWithSignature] error connecting to expansion service: %v\nXDR: %v", err, xdrBase64))
		}

		horizonException, ok := err.(*horizonclient.Error)

		if ok {

			extraErrors := horizonException.Problem.Extras

			for key, val := range extraErrors {
				log.Printf("[SubmitSubwalletXdrWithSignature] Extras: %v is %v\nOwner publicKey: %v, subwallet: %v\n", key, val, ownerPublicKey, subWalletPublicKey)

			}

			resultCodes, errRes := horizonException.ResultCodes()
			if errRes == nil {
				for key, val := range resultCodes.OperationCodes {
					log.Printf("[SubmitSubwalletXdrWithSignature] Result code: %v is %v\nOwner publicKey: %v\nSubWalletPublicKey: %v\n", key, val, ownerPublicKey, subWalletPublicKey)
					// logDiscordFailedPayment(fmt.Sprintf("[SubmitSubwalletXdrWithSignature]Result code: %v is %v\nOwner publicKey: %v\nSubWalletPublicKey: %v\n", key, val, ownerPublicKey, subWalletPublicKey))

				}
			} else {
				log.Printf("[SubmitSubwalletXdrWithSignature] Error getting result codes: %v\n", errRes)
			}

		} else {
			log.Printf("[SubmitSubwalletXdrWithSignature] not horizon error: %v\n", err)

		}

		return "", &tErrors.CustomError{Param: "publicKey", Err: "error subwallet activation failed", ErrMessage: "SubWallet Failed", Code: 500}

	}

	return txnResult.Hash, nil

}

func SubmitSubWalletXdrForChannelAccountWithSignature(client *horizonclient.Client, ownerPublicKey, subWalletPublicKey, channelPK string, xdrBase64 string, primarySignature, subWalletSignature, channelAccountSignature string, subWalletMustSign int) (string, error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	gTxn, err := txnbuild.TransactionFromXDR(xdrBase64)

	if err != nil {
		return "", err
	}

	txn, ok := gTxn.Transaction()

	if !ok {
		return "", &tErrors.ErrorInvalidTransaction{}
	}

	{
		//add signature of the primaryWallet to the new transaction Instance
		txn, err = txn.AddSignatureBase64(network.GetBlockchainNetworkPassPhrase(), ownerPublicKey, primarySignature)
		if err != nil {
			log.Println("[SubmitSubwalletXdrWithSignature] Failed to verify primary signature of primaryWallet on [", network.GetBlockchainNetworkPassPhrase(), "] and [", primarySignature, "] for [", xdrBase64, "] and public key ", ownerPublicKey, ", error [", err, "]")

			return "", err
		}

		//add signature of the subWallet to the new transaction Instance
		txn, err = txn.AddSignatureBase64(network.GetBlockchainNetworkPassPhrase(), subWalletPublicKey, subWalletSignature)
		if err != nil {
			log.Println("[SubmitSubwalletXdrWithSignature] Failed to verify signature of subwallet on [", network.GetBlockchainNetworkPassPhrase(), "] and [", subWalletSignature, "] for [", xdrBase64, "] and public key ", subWalletPublicKey, ", error [", err, "]")

			return "", err
		}
		if subWalletMustSign == 1 {

			//add signature of the channelAccount to the new transaction Instance
			txn, err = txn.AddSignatureBase64(network.GetBlockchainNetworkPassPhrase(), subWalletPublicKey, channelAccountSignature)
			if err != nil {
				log.Println("[SubmitSubwalletXdrWithSignature] Failed to verify signature of channelAccount on [", network.GetBlockchainNetworkPassPhrase(), "] and [", channelAccountSignature, "] for [", xdrBase64, "] and public key ", channelPK, ", error [", err, "]")

				return "", err
			}
		}
	}

	xdrBase64, err = txn.Base64()

	if err != nil {
		log.Printf("[SubmitSubwalletXdrWithSignature] error converting transaction to base64: %v\n", err)
		return "", err
	}

	// log.Println("signed xdr is " + xdrBase64)

	txnResult, err := client.SubmitTransactionXDR(xdrBase64)

	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			discord.Say(fmt.Sprintf("[SubmitSubwalletXdrWithSignature] error connecting to expansion service: %v\nXDR: %v", err, xdrBase64))
		}

		horizonException, ok := err.(*horizonclient.Error)

		if ok {

			extraErrors := horizonException.Problem.Extras

			for key, val := range extraErrors {
				log.Printf("[SubmitSubwalletXdrWithSignature] Extras: %v is %v\nOwner publicKey: %v, subwallet: %v\n", key, val, ownerPublicKey, subWalletPublicKey)

			}

			resultCodes, errRes := horizonException.ResultCodes()
			if errRes == nil {
				for key, val := range resultCodes.OperationCodes {
					log.Printf("[SubmitSubwalletXdrWithSignature] Result code: %v is %v\nOwner publicKey: %v\nSubWalletPublicKey: %v\n", key, val, ownerPublicKey, subWalletPublicKey)
					// logDiscordFailedPayment(fmt.Sprintf("[SubmitSubwalletXdrWithSignature]Result code: %v is %v\nOwner publicKey: %v\nSubWalletPublicKey: %v\n", key, val, ownerPublicKey, subWalletPublicKey))

				}
			} else {
				log.Printf("[SubmitSubwalletXdrWithSignature] Error getting result codes: %v\n", errRes)
			}

		} else {
			log.Printf("[SubmitSubwalletXdrWithSignature] not horizon error: %v\n", err)

		}

		return "", &tErrors.CustomError{Param: "publicKey", Err: "error subwallet activation failed", ErrMessage: "SubWallet Failed", Code: 500}

	}

	return txnResult.Hash, nil

}

func GetUserWallets(user userModels.User, gc *sharedconfig.GlobalConfig) (wallets []userModels.UserWallet) {
	wallets = make([]userModels.UserWallet, 0)
	gc.DB.Where("user_id = ?", user.ID).Find(&wallets)
	return

}
