package users

import (
	"fmt"
	"log"
	"os"
	"strings"
	"trovo-wallet-api/internal/basetxn"
	bc "trovo-wallet-api/internal/blockchainalgofuncs"
	blockchain "trovo-wallet-api/internal/components/assets/blockchain"
	userBc "trovo-wallet-api/internal/components/users/blockchain"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ecnepsnai/discord"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

func CreateNewSubWallet(accountOwner *userModels.User, subWalletInfo *userModels.SubWalletInfo, gc *sharedconfig.GlobalConfig) (*userModels.SubWalletInfo, error) {
	client := gc.BantuExpansionClient
	var err error
	var xdrBase64 string
	var subWalletObj userModels.UserWallet
	var linkedWallet userModels.UserWallet

	subWalletInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	subWalletInfo.SubWalletMustSign = 1
	if len(subWalletInfo.LinkedWalletAddress) == 42 {
		subWalletInfo.LinkedWalletMustSign = 1
	}
	if len(subWalletInfo.ChannelAccount) == 42 {
		//generate xdr for channel account
		xdrBase64, subWalletObj, linkedWallet, err = generateSubWalletXdrWithChannelAccount(accountOwner, subWalletInfo, gc, client)
		if err != nil {
			log.Printf("[CreateNewSubWallet] create sub [%v] for [%v] generateSubWalletXdrWithChannelAccount error:[%v] \n", subWalletInfo.Address, accountOwner.Username, err)
			return subWalletInfo, err
		}
	} else {
		xdrBase64, subWalletObj, linkedWallet, err = generateSubWalletXdr(accountOwner, subWalletInfo, gc, client)
		if err != nil {
			log.Printf("[CreateNewSubWallet] create sub [%v] for [%v] generateSubWalletXdr error:[%v] \n", subWalletInfo.Address, accountOwner.Username, err)
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

	if len(subWalletInfo.ChannelAccountSignature) == 0 && len(subWalletInfo.ChannelAccount) == 42 {
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
	errDBTX := dbTX.Omit(clause.Associations).Create(&subWalletObj).Error
	if errDBTX != nil {
		//unable to save sub wallet. abort
		log.Printf("[CreateNewSubWallet] by [%v] for [%v] Error saving subwallet error:[%v] \n", accountOwner.Username, subWalletInfo.Address, errDBTX)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-saving-subwallet",
			ErrMessage: "There is an error saving sub-wallet. Please, try again later.",
		}
		return subWalletInfo, err
	}
	if len(linkedWallet.ID) > 0 && len(subWalletInfo.LinkedWalletAddress) > 0 {
		//create the data  of linked walletto be sure it goes through
		errDBTX := dbTX.Omit(clause.Associations).Create(&linkedWallet).Error
		if errDBTX != nil {
			//unable to save linked wallet. abort
			log.Printf("[CreateNewSubWallet] by [%v] for [%v] Error saving linked subwallet error:[%v] \n", accountOwner.Username, subWalletInfo.LinkedWalletAddress, errDBTX)

			err = &tErrors.CustomError{
				Param:      "LinkedWalletAddress",
				Err:        "error-saving-subwallet-linked-wallet",
				ErrMessage: "There is an error saving sub-wallet from the linked wallet. Please, try again later.",
			}
			return subWalletInfo, err
		}
	}

	if len(subWalletInfo.ChannelAccountSignature) > 0 && len(subWalletInfo.ChannelAccount) == 42 {
		txnHash, err := SubmitSubWalletXdrForChannelAccountWithSignature(client, accountOwner.Address, subWalletInfo.Address, subWalletInfo.ChannelAccount, xdrBase64, subWalletInfo.PrimarySignature, subWalletInfo.SubWalletSignature, subWalletInfo.ChannelAccountSignature, subWalletInfo.SubWalletMustSign)
		if err != nil {
			log.Printf("[CreateNewSubWallet] by [%v] for [%v] SubmitSubWalletXdrForChannelAccountWithSignature error:[%v] \n", accountOwner.Username, subWalletInfo.Address, err)
			return subWalletInfo, err
		}
		subWalletInfo.TransactionID = txnHash
		dbTX.Commit()

		{
			//send to monitoring service
			trackAddress := userModels.TrackedAddress{
				Address: subWalletInfo.Address,
			}
			errTrack := gc.RoachDB.Create(&trackAddress).Error
			if errTrack != nil {
				//if tracking of public key fails, then payment history generation service will pick it up and do justice to it
				discord.Say(fmt.Sprintf("[CreateNewSubWallet] tracking public key for payment history failed for user:%v, with DB Error:%v\n\n\nFailedData:%+v", accountOwner.Username, errTrack, subWalletInfo))

			}

			if subWalletInfo.LinkedWalletMustSign == 1 {

				{
					//send to monitoring service
					trackAddress := userModels.TrackedAddress{
						Address: linkedWallet.ID,
					}
					errTrack := gc.RoachDB.Create(&trackAddress).Error
					if errTrack != nil {
						//if tracking of public key fails, then payment history generation service will pick it up and do justice to it
						discord.Say(fmt.Sprintf("[CreateNewSubWallet] tracking linked wallet public key for payment history failed for user:%v, with DB Error:%v\n\n\nFailedData:%+v", accountOwner.Username, errTrack, subWalletInfo))

					}
				}
			}
		}
	} else {
		signatures := make(map[string]string, 0)
		signatures[accountOwner.PrimarySigner] = subWalletInfo.PrimarySignature

		if subWalletInfo.SubWalletMustSign == 1 {
			signatures[subWalletInfo.Address] = subWalletInfo.SubWalletSignature
		}

		if subWalletInfo.LinkedWalletMustSign == 1 {
			//it is a linked wallet operation. build a map of signers
			signatures[subWalletInfo.LinkedWalletAddress] = subWalletInfo.LinkedWalletSignature
		}

		// txnHash, err := SubmitSubWalletXdrWithSignature(client, accountOwner.Address, accountOwner.PrimarySigner, subWalletInfo.Address, xdrBase64, subWalletInfo.PrimarySignature, subWalletInfo.SubWalletSignature, subWalletInfo.SubWalletMustSign)
		txnHash, err := SubmitSubWalletXdrWithSignatures(client, signatures, xdrBase64)
		if err != nil {
			log.Printf("[CreateNewSubWallet] by [%v] for [%v] SubmitSubWalletXdrWithSignatures error:[%v] \n", accountOwner.Username, subWalletInfo.Address, err)
			return subWalletInfo, err
		}
		subWalletInfo.TransactionID = txnHash
		dbTX.Commit()
		{
			//send to monitoring service
			trackAddress := userModels.TrackedAddress{
				Address: subWalletInfo.Address,
			}
			errTrack := gc.RoachDB.Create(&trackAddress).Error
			if errTrack != nil {
				//if tracking of public key fails, then payment history generation service will pick it up and do justice to it
				discord.Say(fmt.Sprintf("[CreateNewSubWallet] tracking public key for payment history failed for user:%v, with DB Error:%v\n\n\nFailedData:%+v", accountOwner.Username, errTrack, subWalletInfo))

			}
			if len(linkedWallet.ID) > 0 && len(subWalletInfo.LinkedWalletAddress) > 0 {
				//send to monitoring service
				trackAddress := userModels.TrackedAddress{
					Address: linkedWallet.ID,
				}
				errTrack := gc.RoachDB.Create(&trackAddress).Error
				if errTrack != nil {
					//if tracking of public key fails, then payment history generation service will pick it up and do justice to it
					discord.Say(fmt.Sprintf("[CreateNewSubWallet] tracking linked public key for payment history failed for user:%v, with DB Error:%v\n\n\nFailedData:%+v", accountOwner.Username, errTrack, subWalletInfo))

				}
			}
		}

	}
	accountOwner.InvalidateUserCache(gc)

	return subWalletInfo, nil
}

func generateSubWalletXdr(accountOwner *userModels.User, subWalletInfo *userModels.SubWalletInfo, gc *sharedconfig.GlobalConfig, client *ethclient.Client) (xdrbase64 string, subWalletObj, linkedWallet userModels.UserWallet, err error) {
	// var linkedWallet userModels.UserWallet
	ops := make([]basetxn.Operation, 0)
	subWalletInfo.Messages = make([]string, 0)
	var activationAmount = decimal.NewFromFloat(6)
	WALLET_SIGNER_ACTIVATION_AMOUNT := accountOwner.UserWallets[0].GetActivationFee("WALLET_SIGNER_ACTIVATION_AMOUNT", gc)
	SUB_WALLET_ACTIVATION_AMOUNT := accountOwner.UserWallets[0].GetActivationFee("SUB_WALLET_ACTIVATION_AMOUNT", gc)
	ISSUING_SUB_WALLET_ACTIVATION_AMOUNT := accountOwner.UserWallets[0].GetActivationFee("ISSUING_SUB_WALLET_ACTIVATION_AMOUNT", gc)
	MM_SUB_WALLET_ACTIVATION_AMOUNT := accountOwner.UserWallets[0].GetActivationFee("MM_SUB_WALLET_ACTIVATION_AMOUNT", gc)
	BULKPAYMENT_SUB_WALLET_ACTIVATION_AMOUNT := accountOwner.UserWallets[0].GetActivationFee("BULKPAYMENT_SUB_WALLET_ACTIVATION_AMOUNT", gc)
	SUBWALLET_CREATION_FEE := accountOwner.UserWallets[0].GetSubwalletCreationFee(gc)
	feeKeypair, e := evmkeypair.ParseFull(SUBWALLET_CREATION_FEE.FeeWalletSecretKey)
	if e != nil {
		gc.LogDiscordFailedRequest("SUBWALLET CREATION SECRET KEY IS INVALID")
		return "", subWalletObj, linkedWallet, &tErrors.ErrorTemporaryServerError{}
	}
	var minBalance = decimal.NewFromFloat(3.0)
	dab := strings.Split(os.Getenv("DOLLAR_ASSET"), ":")
	dollarAsset := basetxn.CreditAsset{Code: dab[0], Issuer: dab[1]}
	if SUB_WALLET_ACTIVATION_AMOUNT.Amount > 0 {
		activationAmount = decimal.NewFromFloat(SUB_WALLET_ACTIVATION_AMOUNT.Amount)
	}
	if ISSUING_SUB_WALLET_ACTIVATION_AMOUNT.Amount > 0 && subWalletInfo.WalletType == 1 {
		activationAmount = decimal.NewFromFloat(ISSUING_SUB_WALLET_ACTIVATION_AMOUNT.Amount)
	}
	if MM_SUB_WALLET_ACTIVATION_AMOUNT.Amount > 0 && subWalletInfo.WalletType == 2 {
		activationAmount = decimal.NewFromFloat(MM_SUB_WALLET_ACTIVATION_AMOUNT.Amount)
	}
	if BULKPAYMENT_SUB_WALLET_ACTIVATION_AMOUNT.Amount > 0 && subWalletInfo.WalletType == 3 {
		activationAmount = decimal.NewFromFloat(BULKPAYMENT_SUB_WALLET_ACTIVATION_AMOUNT.Amount)
	}
	if len(os.Getenv("WALLET_MINIMUM_BALANCE")) > 0 {
		minBalance = decimal.RequireFromString(os.Getenv("WALLET_MINIMUM_BALANCE"))
	}
	{
		//check if the sub-wallet passes the validation
		subWalletObj, err = accountOwner.BuildNewSubWallet(subWalletInfo.Address, subWalletInfo.WalletTag, subWalletInfo.WalletDescription, subWalletInfo.WalletType, subWalletInfo.LinkedWalletAddress, gc)
		if err != nil {
			log.Printf("[generateSubWalletXdr] by [%v] for [%v] BuildNewSubWallet error:[%v] \n", accountOwner.Username, subWalletInfo.Address, err)
			return "", subWalletObj, linkedWallet, err
		}

		//set the subwallet suggested alias
		subWalletInfo.Alias = subWalletObj.Alias

		//build linked wallet. Linked wallet public key already validated in buildnewsubwallet function. so if it is not valid it won't get here. and if it is valid, then below procedure will execute.
		if len(subWalletInfo.LinkedWalletAddress) > 0 {
			linkedWallet, err = subWalletObj.BuildNewLinkedSubWallet(accountOwner, gc)
			if err != nil {
				log.Printf("[generateSubWalletXdr] by [%v] for [%v] BuildNewLinkedSubWallet error:[%v] \n", accountOwner.Username, subWalletInfo.Address, err)
				return "", subWalletObj, linkedWallet, err
			}
		}

	}

	//check if it is first call to create sub-wallet

	//populate the subwallet Info and generate the transaction

	//check if primary account has native enough native balance
	var nativeAsset basetxn.Asset = basetxn.NativeAsset{}
	_, _, primaryAccountNativeBalance, _, primarySourceAccount, errAct := network.BlockchainAccountProperties(client, accountOwner.Address, nativeAsset)
	if errAct != nil {
		log.Printf("[generateSubWalletXdr] by [%v] for [%v] Primary Account Properties error:[%v] \n", accountOwner.Username, subWalletInfo.Address, errAct)

		return "", subWalletObj, linkedWallet, errAct
	}
	if subWalletInfo.LinkedWalletAddress == "" {
		if (primaryAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
			log.Printf("[generateSubWalletXdr] by [%v] for [%v] PrimaryAccount underfunded \n", accountOwner.Username, subWalletInfo.Address)

			err = &tErrors.CustomError{
				Param:      "publicKey",
				Err:        "error-primary-account-underfunded",
				ErrMessage: fmt.Sprintf("Primary account does not have enough %v balance to create sub-wallet", os.Getenv("NATIVE_ASSET_CODE")),
				Code:       404,
			}
			return "", subWalletObj, linkedWallet, err
		}
	} else {
		//since linked wallet is present, two wallets would be activated. check that balance is double at least
		if (primaryAccountNativeBalance.Sub(activationAmount.Mul(decimal.NewFromInt(2)))).LessThan(minBalance) {
			log.Printf("[generateSubWalletXdr] by [%v] for [%v] PrimaryAccount underfunded \n", accountOwner.Username, subWalletInfo.Address)

			err = &tErrors.CustomError{
				Param:      "publicKey",
				Err:        "error-primary-account-underfunded",
				ErrMessage: fmt.Sprintf("Primary account does not have enough %v balance (%v) to create sub-wallet", os.Getenv("NATIVE_ASSET_CODE"), activationAmount.Mul(decimal.NewFromInt(2)).String()),
				Code:       404,
			}
			return "", subWalletObj, linkedWallet, err
		}
	}

	var walletSigner *evmkeypair.Full
	if subWalletInfo.WalletType == 2 {
		walletSigner, _ = bc.MarketMakingSignerKeypair(accountOwner.Username, subWalletInfo.Address)

	}
	if subWalletInfo.WalletType == 3 {
		walletSigner, _ = bc.BulkPaymentSignerKeypair(accountOwner.Username, subWalletInfo.Address)

	}

	subWalletAccountExists, _, subWalletAccountNativeBalance, _, subWalletAccountObject, _ := network.BlockchainAccountProperties(client, subWalletInfo.Address, nativeAsset)

	if subWalletAccountExists {
		if subWalletAccountNativeBalance.LessThan(minBalance) {
			ops = append(ops, &basetxn.Payment{
				Destination:   subWalletInfo.Address,
				Amount:        activationAmount.String(),
				Asset:         nativeAsset,
				SourceAccount: accountOwner.Address,
			})
		}
		if subWalletInfo.WalletType != 1 {
			if os.Getenv("ENABLE_DOLLAR_ASSET_BY_DEFAULT") != "0" {
				//enable dollar asset if not minting wallet
				_, trusted, _, _, _, _ := network.BlockchainAccountProperties(client, subWalletInfo.Address, dollarAsset)
				if !trusted {
					ops = append(ops, &basetxn.ChangeTrust{
						Line:          dollarAsset,
						Limit:         "900000000000",
						SourceAccount: subWalletInfo.Address,
					})
				}
			}
			if os.Getenv("ENABLE_NAIRA_ASSET_BY_DEFAULT") != "0" {
				//enable NAIRA asset if not minting wallet
				ndab := strings.Split(os.Getenv("NAIRA_ASSET"), ":")
				nairaAsset := basetxn.CreditAsset{Code: ndab[0], Issuer: ndab[1]}
				_, ntrusted, _, _, _, _ := network.BlockchainAccountProperties(client, subWalletInfo.Address, nairaAsset)
				if !ntrusted {
					ops = append(ops, &basetxn.ChangeTrust{
						Line:          nairaAsset,
						Limit:         "900000000000",
						SourceAccount: subWalletInfo.Address,
					})
				}
			}
			if os.Getenv("ENABLE_TROV_ASSET_BY_DEFAULT") != "0" {

				issuer := os.Getenv("TROV_ASSET_ISSUER")
				if issuer != "" {
					trovAsset := basetxn.CreditAsset{Code: "TROV", Issuer: issuer}
					_, ntrusted, _, _, _, _ := network.BlockchainAccountProperties(client, subWalletInfo.Address, trovAsset)
					if !ntrusted {
						ops = append(ops, &basetxn.ChangeTrust{
							Line:          trovAsset,
							Limit:         "900000000000",
							SourceAccount: subWalletInfo.Address,
						})
					}
				}
			}

		}

		//account exists and native balance is less than needed. add 3 native token to the wallet
		if subWalletInfo.WalletType == 0 {

			//after topping up, it now has enough balance to add primary wallet as signer if it is not already a signer
			if !network.IsAccountSigner(subWalletAccountObject.Address, accountOwner.PrimarySigner) {

				ops = append(ops, &basetxn.SetOptions{
					Signer: &basetxn.Signer{
						Address: accountOwner.PrimarySigner,
						Weight:  1,
					},
					SourceAccount: subWalletInfo.Address,
				})
				// //prevent any future changes to the auth flag of the wallet.
				// ops = append(ops, &basetxn.SetOptions{
				// 	SetFlags:      []basetxn.AccountFlag{basetxn.AuthImmutable},
				// 	SourceAccount: subWalletInfo.Address,
				// })
			} else {
				subWalletInfo.SubWalletMustSign = 0
			}

		}
		if subWalletInfo.WalletType == 1 {

			//after topping up, it now has enough balance to add primary wallet as signer if it is not already a signer
			if !network.IsAccountSigner(subWalletAccountObject.Address, accountOwner.PrimarySigner) {

				ops = append(ops, &basetxn.SetOptions{
					Signer: &basetxn.Signer{
						Address: accountOwner.PrimarySigner,
						Weight:  1,
					},
					SetFlags:      []basetxn.AccountFlag{basetxn.AuthRequired, basetxn.AuthClawbackEnabled},
					SourceAccount: subWalletInfo.Address,
				})
			} else {
				subWalletInfo.SubWalletMustSign = 0
			}

		}
		//make it custodial
		if subWalletInfo.WalletType == 2 || subWalletInfo.WalletType == 3 {

			signerExists, _, _, _, _, _ := network.BlockchainAccountProperties(client, walletSigner.Address(), nativeAsset)
			if !signerExists {
				ops = append(ops, &basetxn.CreateAccount{
					Destination:   walletSigner.Address(),
					Amount:        fmt.Sprintf("%v", WALLET_SIGNER_ACTIVATION_AMOUNT.Amount),
					SourceAccount: accountOwner.Address,
				})

			} else {
				ops = append(ops, &basetxn.Payment{
					Destination:   walletSigner.Address(),
					Amount:        fmt.Sprintf("%v", WALLET_SIGNER_ACTIVATION_AMOUNT.Amount),
					Asset:         nativeAsset,
					SourceAccount: accountOwner.Address,
				})

			}

			//after creation, it now exists with enough balance to add primary wallet as signer
			if !userBc.SignerIsValid(subWalletInfo.Address, accountOwner.PrimarySigner) {
				ops = append(ops, &basetxn.SetOptions{
					Signer: &basetxn.Signer{
						Address: accountOwner.PrimarySigner,
						Weight:  1,
					},
					SourceAccount: subWalletInfo.Address,
				})
			}
			if !userBc.SignerIsValid(subWalletInfo.Address, walletSigner.Address()) {
				ops = append(ops, &basetxn.SetOptions{
					Signer: &basetxn.Signer{
						Address: walletSigner.Address(),
						Weight:  3,
					},
					SourceAccount: subWalletInfo.Address,
				})
			}
		}

	}

	//add recovery key if account recovery is enabled
	if accountOwner.AccountRecoveryEnabled == 1 {
		recoveryKeyAddress := bc.GetRecoveryAccountAddress(accountOwner.Username, accountOwner.Address)

		if len(recoveryKeyAddress) == 42 {
			if subWalletInfo.WalletType == 0 || subWalletInfo.WalletType == 1 {
				if !userBc.SignerIsValid(subWalletInfo.Address, recoveryKeyAddress) {
					ops = append(ops, &basetxn.SetOptions{
						Signer: &basetxn.Signer{
							Address: recoveryKeyAddress,
							Weight:  1,
						},
						SourceAccount: subWalletInfo.Address,
					})
				}
			}

			if subWalletInfo.WalletType == 2 || subWalletInfo.WalletType == 3 {
				if !userBc.SignerIsValid(subWalletInfo.Address, recoveryKeyAddress) {
					ops = append(ops, &basetxn.SetOptions{
						Signer: &basetxn.Signer{
							Address: recoveryKeyAddress,
							Weight:  3,
						},
						SourceAccount: subWalletInfo.Address,
					})
				}
			}

		}

	}

	//perform routine for linked wallet if available
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	{
		if len(linkedWallet.ID) == 42 {

			linkedSubWalletAccountExists, _, linkedSubWalletAccountNativeBalance, _, linkedSubWalletAccountObject, _ := network.BlockchainAccountProperties(client, linkedWallet.ID, nativeAsset)

			if linkedSubWalletAccountExists {
				if linkedSubWalletAccountNativeBalance.LessThan(minBalance) {
					ops = append(ops, &basetxn.Payment{
						Destination:   linkedWallet.ID,
						Amount:        activationAmount.String(),
						Asset:         nativeAsset,
						SourceAccount: accountOwner.Address,
					})
				}
				if os.Getenv("ENABLE_DOLLAR_ASSET_BY_DEFAULT") != "0" {
					//enable default assets
					_, trusted, _, _, _, _ := network.BlockchainAccountProperties(client, linkedWallet.ID, dollarAsset)
					if !trusted {
						ops = append(ops, &basetxn.ChangeTrust{
							Line:          dollarAsset,
							Limit:         "900000000000",
							SourceAccount: linkedWallet.ID,
						})
					}
				}
				if os.Getenv("ENABLE_NAIRA_ASSET_BY_DEFAULT") != "0" {
					//enable NAIRA asset if not minting wallet
					ndab := strings.Split(os.Getenv("NAIRA_ASSET"), ":")
					nairaAsset := basetxn.CreditAsset{Code: ndab[0], Issuer: ndab[1]}
					_, ntrusted, _, _, _, _ := network.BlockchainAccountProperties(client, linkedWallet.ID, nairaAsset)
					if !ntrusted {
						ops = append(ops, &basetxn.ChangeTrust{
							Line:          nairaAsset,
							Limit:         "900000000000",
							SourceAccount: linkedWallet.ID,
						})
					}
				}

				if os.Getenv("ENABLE_TROV_ASSET_BY_DEFAULT") != "0" {

					issuer := os.Getenv("TROV_ASSET_ISSUER")
					if issuer != "" {
						trovAsset := basetxn.CreditAsset{Code: "TROV", Issuer: issuer}
						_, ntrusted, _, _, _, _ := network.BlockchainAccountProperties(client, linkedWallet.ID, trovAsset)
						if !ntrusted {
							ops = append(ops, &basetxn.ChangeTrust{
								Line:          trovAsset,
								Limit:         "900000000000",
								SourceAccount: linkedWallet.ID,
							})
						}
					}
				}

				if !network.IsAccountSigner(linkedSubWalletAccountObject.Address, accountOwner.PrimarySigner) {

					subWalletInfo.SubWalletMustSign = 0

					ops = append(ops, &basetxn.SetOptions{
						Signer: &basetxn.Signer{
							Address: accountOwner.PrimarySigner,
							Weight:  1,
						},
						SourceAccount: linkedWallet.ID,
					})
				}

			}

			//add recovery key if account recovery is enabled
			if accountOwner.AccountRecoveryEnabled == 1 {
				recoveryKeyAddress := bc.GetRecoveryAccountAddress(accountOwner.Username, accountOwner.Address)

				if len(recoveryKeyAddress) == 42 {

					if !userBc.SignerIsValid(linkedWallet.ID, recoveryKeyAddress) {
						ops = append(ops, &basetxn.SetOptions{
							Signer: &basetxn.Signer{
								Address: recoveryKeyAddress,
								Weight:  1,
							},
							SourceAccount: linkedWallet.ID,
						})
					}
				}

			}

		}
	}
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	signForFeeTrustLine := 0
	if accountOwner.Username != "atprofile" {
		//add fee for transaction
		usdPrice, _, _ := blockchain.GetDollarPrice(SUBWALLET_CREATION_FEE.FeeAssetCode, SUBWALLET_CREATION_FEE.FeeAssetIssuer, gc, true)

		serviceFee := decimal.NewFromFloat(SUBWALLET_CREATION_FEE.FeeFixed)
		if e != nil {
			serviceFee = decimal.Zero
		}
		if serviceFee.IsPositive() {
			serviceFee = decimal.RequireFromString(usdPrice).Div(serviceFee).Truncate(7)
		} else {
			serviceFee = decimal.Zero
		}

		if serviceFee.IsPositive() {

			feeAddress := feeKeypair.Address()

			feeAsset := basetxn.CreditAsset{Code: SUBWALLET_CREATION_FEE.FeeAssetCode, Issuer: SUBWALLET_CREATION_FEE.FeeAssetIssuer}
			_, _, _, assetBalance, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, accountOwner.Address, feeAsset)
			if assetBalance.LessThan(serviceFee) {
				return "", subWalletObj, linkedWallet, &tErrors.CustomError{Param: "username", Err: "error-primary-wallet-underfunded", ErrMessage: fmt.Sprintf("%v %v is required on wallet %v to pay for fees for this service. Please first fund the wallet with at least %v %v.", serviceFee.String(), SUBWALLET_CREATION_FEE.FeeAssetCode, accountOwner.Username, serviceFee.Sub(assetBalance), SUBWALLET_CREATION_FEE.FeeAssetCode)}

			}

			_, feeAccountTrustsAsset, _, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, feeAddress, feeAsset)
			if !feeAccountTrustsAsset {
				signForFeeTrustLine = 1
				//establish trustline automatically
				ops = append(ops, &basetxn.ChangeTrust{
					Line:          feeAsset,
					Limit:         "900000000000",
					SourceAccount: feeAddress,
				})

			}

			ops = append(ops, &basetxn.Payment{
				Destination:   feeAddress,
				Amount:        serviceFee.String(),
				SourceAccount: accountOwner.Address,
				Asset:         feeAsset,
			})
			// paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("%v %v will be added from wallet %v as service fee (%v).", serviceFee.String(), assetCode, sourceWallet.Alias, feeLabel))
			subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("%v %v ($%v USD) will be deducted from wallet %v as service fee.", serviceFee.String(), SUBWALLET_CREATION_FEE.FeeAssetCode, SUBWALLET_CREATION_FEE.FeeFixed, accountOwner.Username))

		}

	}

	if subWalletInfo.WalletType == 0 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("%v %v will be sent from your primary wallet to this new sub-wallet for wallet activation. It will become the new balance of the subwallet.", activationAmount.String(), os.Getenv("NATIVE_ASSET_CODE")))

	}
	if subWalletInfo.WalletType == 1 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Because this subwallet is designated to be a token minting wallet, %v %v will be deducted from your primary wallet and be used to activate it alongside the distriution wallet. Please note that token minting wallets cannot be used to send payments.", (activationAmount.Mul(decimal.NewFromInt(2))).String(), os.Getenv("NATIVE_ASSET_CODE")))

	}
	if subWalletInfo.WalletType == 2 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Because this subwallet is designated to be a market making wallet, %v %v will be deducted from your primary wallet and be used to activate it and the custodial signer. It will become the new balance of the subwallet and custodial signer. Please note that MM wallets cannot be used to send normal payments, but only used for market making.", (activationAmount.Add(decimal.NewFromFloat(WALLET_SIGNER_ACTIVATION_AMOUNT.Amount))).String(), os.Getenv("NATIVE_ASSET_CODE")))

	}

	if subWalletInfo.WalletType == 3 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Because this subwallet is designated to be a bulk-payment wallet, %v %v will be deducted from your primary wallet and be used to activate it and the custodial signer. It will become the new balance of the subwallet and custodial signer. Please note that bulk-payment wallets cannot be used to send normal payments, but only be used by internal system to disburse bulk payments on your behalf.", (activationAmount.Add(decimal.NewFromFloat(WALLET_SIGNER_ACTIVATION_AMOUNT.Amount))).String(), os.Getenv("NATIVE_ASSET_CODE")))

	}
	if subWalletInfo.WalletType == 2 || subWalletInfo.WalletType == 3 {
		ops = append(ops, &basetxn.SetOptions{
			LowThreshold:    basetxn.NewThreshold(uint32(3)),
			MediumThreshold: basetxn.NewThreshold(uint32(3)),
			HighThreshold:   basetxn.NewThreshold(uint32(3)),
			SourceAccount:   subWalletInfo.Address,
		})
	}
	// if subWalletInfo.WalletType == 0 || subWalletInfo.WalletType == 1 {
	// 	ops = append(ops, &basetxn.SetOptions{
	// 		LowThreshold:    basetxn.NewThreshold(uint32(1)),
	// 		MediumThreshold: basetxn.NewThreshold(uint32(1)),
	// 		HighThreshold:   basetxn.NewThreshold(uint32(1)),
	// 		SourceAccount:   subWalletInfo.Address,
	// 	})
	// }

	// Construct the transaction that holds the operations to execute on the network
	tx, err := basetxn.NewTransaction(
		basetxn.TransactionParams{
			SourceAccount:        primarySourceAccount.Address,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              2000,
			Memo:                 "Create Sub Wallet",
		},
	)
	if err != nil {
		log.Println("[generateSubWalletXdr] error constructing transaction ", err)
		return "", subWalletObj, linkedWallet, err
	}

	if signForFeeTrustLine == 1 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), feeKeypair)

		if err != nil {
			log.Println("[generateSubWalletXdr] error signing transaction with fee wallet key ", err)
			return "", subWalletObj, linkedWallet, &tErrors.ErrorTemporaryServerError{}
		}
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateSubWalletXdr] error getting txn base64", err)
		return "", subWalletObj, linkedWallet, err
	}

	return xdrBase64, subWalletObj, linkedWallet, nil

}

func generateSubWalletXdrWithChannelAccount(user *userModels.User, subWalletInfo *userModels.SubWalletInfo, gc *sharedconfig.GlobalConfig, client *ethclient.Client) (xdrbase64 string, subWalletObj, linkedWallet userModels.UserWallet, err error) {
	// var linkedWallet userModels.UserWallet
	ops := make([]basetxn.Operation, 0)
	subWalletInfo.Messages = make([]string, 0)
	var activationAmount = decimal.NewFromFloat(6.0)
	var minBalance = decimal.NewFromFloat(3.0)
	dab := strings.Split(os.Getenv("DOLLAR_ASSET"), ":")
	dollarAsset := basetxn.CreditAsset{Code: dab[0], Issuer: dab[1]}
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
		subWalletObj, err = user.BuildNewSubWallet(subWalletInfo.Address, subWalletInfo.WalletTag, subWalletInfo.WalletDescription, subWalletInfo.WalletType, subWalletInfo.LinkedWalletAddress, gc)
		if err != nil {
			log.Printf("[generateSubWalletXdrWithChannelAccount] by [%v] for [%v] BuildNewSubWallet error:[%v] \n", user.Username, subWalletInfo.Address, err)

			return "", subWalletObj, linkedWallet, err
		}
		//set the subwallet suggested alias
		subWalletInfo.Alias = subWalletObj.Alias

		//build linked wallet. Linked wallet public key already validated in buildnewsubwallet function. so if it is not valid it won't get here. and if it is valid, then below procedure will execute.
		if len(subWalletInfo.LinkedWalletAddress) > 0 {
			linkedWallet, err = subWalletObj.BuildNewLinkedSubWallet(user, gc)
			if err != nil {
				log.Printf("[generateSubWalletXdr] by [%v] for [%v] BuildNewLinkedSubWallet error:[%v] \n", user.Username, subWalletInfo.Address, err)
				return "", subWalletObj, linkedWallet, err
			}
		}
	}

	//check if it is first call to create sub-wallet

	//populate the subwallet Info and generate the transaction

	//check if primary account has native enough native balance
	var nativeAsset basetxn.Asset = basetxn.NativeAsset{}
	_, _, primaryAccountNativeBalance, _, _, errAct := network.BlockchainAccountProperties(client, user.Address, nativeAsset)
	if errAct != nil {
		log.Printf("[generateSubWalletXdrWithChannelAccount] by [%v] for [%v] Primary Account Error error:[%v] \n", user.Username, subWalletInfo.Address, errAct)

		return "", subWalletObj, linkedWallet, errAct
	}

	if subWalletInfo.LinkedWalletAddress == "" {
		if (primaryAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
			log.Printf("[generateSubWalletXdrWithChannelAccount] by [%v] for [%v] Primary Account Underfunded\n", user.Username, subWalletInfo.Address)

			err = &tErrors.CustomError{
				Param:      "publicKey",
				Err:        "error-primary-account-underfunded",
				ErrMessage: fmt.Sprintf("Primary account does not have enough %s balance to create sub-wallet", os.Getenv("NATIVE_ASSET_CODE")),
				Code:       404,
			}
			return "", subWalletObj, linkedWallet, err
		}
	} else {
		//since linked wallet is present, two wallets would be activated. check that balance is double at least

		if (primaryAccountNativeBalance.Sub(activationAmount.Mul(decimal.NewFromInt(2)))).LessThan(minBalance) {
			log.Printf("[generateSubWalletXdrWithChannelAccount] by [%v] for [%v] Primary Account Underfunded\n", user.Username, subWalletInfo.Address)

			err = &tErrors.CustomError{
				Param:      "publicKey",
				Err:        "error-primary-account-underfunded",
				ErrMessage: fmt.Sprintf("Primary account does not have enough %s balance to create sub-wallet", os.Getenv("NATIVE_ASSET_CODE")),
				Code:       404,
			}
			return "", subWalletObj, linkedWallet, err
		}
	}

	var walletSigner *evmkeypair.Full
	if subWalletInfo.WalletType == 2 {
		walletSigner, _ = bc.MarketMakingSignerKeypair(user.Username, subWalletInfo.Address)

	}
	if subWalletInfo.WalletType == 3 {
		walletSigner, _ = bc.BulkPaymentSignerKeypair(user.Username, subWalletInfo.Address)

	}

	subWalletAccountExists, _, subWalletAccountNativeBalance, _, subWalletAccountObject, _ := network.BlockchainAccountProperties(client, subWalletInfo.Address, nativeAsset)

	if subWalletAccountExists {
		if subWalletAccountNativeBalance.LessThan(minBalance) {
			ops = append(ops, &basetxn.Payment{
				Destination:   subWalletInfo.Address,
				Amount:        activationAmount.String(),
				Asset:         nativeAsset,
				SourceAccount: user.Address,
			})
		}
		//account exists and native balance is less than needed. add 3 native token to the wallet
		if subWalletInfo.WalletType == 0 || subWalletInfo.WalletType == 1 {

			//after topping up, it now has enough balance to add primary wallet as signer if it is not already a signer
			if !network.IsAccountSigner(subWalletAccountObject.Address, user.PrimarySigner) {

				ops = append(ops, &basetxn.SetOptions{
					Signer: &basetxn.Signer{
						Address: user.PrimarySigner,
						Weight:  1,
					},
					SourceAccount: subWalletInfo.Address,
				})
			} else {
				subWalletInfo.SubWalletMustSign = 0
			}

		}
		//make it custodial
		if subWalletInfo.WalletType == 2 || subWalletInfo.WalletType == 3 {

			signerExists, _, _, _, _, _ := network.BlockchainAccountProperties(client, walletSigner.Address(), nativeAsset)
			if !signerExists {
				ops = append(ops, &basetxn.CreateAccount{
					Destination:   walletSigner.Address(),
					Amount:        os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT"),
					SourceAccount: user.Address,
				})

			} else {
				ops = append(ops, &basetxn.Payment{
					Destination:   walletSigner.Address(),
					Amount:        os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT"),
					Asset:         nativeAsset,
					SourceAccount: user.Address,
				})

			}

			//after creation, it now exists with enough balance to add primary wallet as signer
			if !userBc.SignerIsValid(subWalletInfo.Address, user.PrimarySigner) {
				ops = append(ops, &basetxn.SetOptions{
					Signer: &basetxn.Signer{
						Address: user.PrimarySigner,
						Weight:  1,
					},
					SourceAccount: subWalletInfo.Address,
				})
			}
			if !userBc.SignerIsValid(subWalletInfo.Address, walletSigner.Address()) {
				ops = append(ops, &basetxn.SetOptions{
					Signer: &basetxn.Signer{
						Address: walletSigner.Address(),
						Weight:  3,
					},
					SourceAccount: subWalletInfo.Address,
				})
			}
		}

	}

	//add recovery key if account recovery is enabled
	if user.AccountRecoveryEnabled == 1 {
		recoveryKeyAddress := bc.GetRecoveryAccountAddress(user.Username, user.Address)

		if len(recoveryKeyAddress) == 42 {
			if !userBc.SignerIsValid(subWalletInfo.Address, recoveryKeyAddress) {
				ops = append(ops, &basetxn.SetOptions{
					Signer: &basetxn.Signer{
						Address: recoveryKeyAddress,
						Weight:  3,
					},
					SourceAccount: subWalletInfo.Address,
				})
			}

		}

	}

	//perform routine for linked wallet if available
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	{
		if len(linkedWallet.ID) == 42 {

			linkedSubWalletAccountExists, _, linkedSubWalletAccountNativeBalance, _, linkedSubWalletAccountObject, _ := network.BlockchainAccountProperties(client, linkedWallet.ID, nativeAsset)

			if linkedSubWalletAccountExists {
				if linkedSubWalletAccountNativeBalance.LessThan(minBalance) {
					ops = append(ops, &basetxn.Payment{
						Destination:   linkedWallet.ID,
						Amount:        activationAmount.String(),
						Asset:         nativeAsset,
						SourceAccount: user.Address,
					})
				}

				//enable default assets
				_, trusted, _, _, _, _ := network.BlockchainAccountProperties(client, linkedWallet.ID, dollarAsset)
				if !trusted {
					ops = append(ops, &basetxn.ChangeTrust{
						Line:          dollarAsset,
						Limit:         "900000000000",
						SourceAccount: linkedWallet.ID,
					})
				}

				if os.Getenv("ENABLE_NAIRA_ASSET_BY_DEFAULT") == "1" {
					//enable NAIRA asset if not minting wallet
					ndab := strings.Split(os.Getenv("NAIRA_ASSET"), ":")
					nairaAsset := basetxn.CreditAsset{Code: ndab[0], Issuer: ndab[1]}
					_, ntrusted, _, _, _, _ := network.BlockchainAccountProperties(client, linkedWallet.ID, nairaAsset)
					if !ntrusted {
						ops = append(ops, &basetxn.ChangeTrust{
							Line:          nairaAsset,
							Limit:         "900000000000",
							SourceAccount: linkedWallet.ID,
						})
					}
				}

				if !network.IsAccountSigner(linkedSubWalletAccountObject.Address, user.PrimarySigner) {

					subWalletInfo.SubWalletMustSign = 0

					ops = append(ops, &basetxn.SetOptions{
						Signer: &basetxn.Signer{
							Address: user.PrimarySigner,
							Weight:  1,
						},
						SourceAccount: linkedWallet.ID,
					})
				}

			}

			//add recovery key if account recovery is enabled
			if user.AccountRecoveryEnabled == 1 {
				recoveryKeyAddress := bc.GetRecoveryAccountAddress(user.Username, user.Address)

				if len(recoveryKeyAddress) == 42 {

					if !userBc.SignerIsValid(linkedWallet.ID, recoveryKeyAddress) {
						ops = append(ops, &basetxn.SetOptions{
							Signer: &basetxn.Signer{
								Address: recoveryKeyAddress,
								Weight:  1,
							},
							SourceAccount: linkedWallet.ID,
						})
					}
				}

			}

		}
	}
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// TODO: FEE
	////////

	if subWalletInfo.WalletType == 0 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("%v %v will be deducted from your primary wallet and be used to activate the sub-wallet.", activationAmount.String(), os.Getenv("NATIVE_ASSET_CODE")))

	}
	if subWalletInfo.WalletType == 1 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Because this subwallet is designated to be an asset issuing wallet, %v %v will be deducted from your primary wallet and be used to activate it. Please note that asset issuing wallets cannot be used to send payments.", (activationAmount.Mul(decimal.NewFromInt(2))).String(), os.Getenv("NATIVE_ASSET_CODE")))

	}
	if subWalletInfo.WalletType == 2 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Because this subwallet is designated to be an market making wallet, %v %v will be deducted from your primary wallet and be used to activate it and the custodial signer. Please note that MM wallets cannot be used to send normal payments, but only used for market making.", (activationAmount.Add(decimal.RequireFromString(os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT")))).String(), os.Getenv("NATIVE_ASSET_CODE")))

	}

	if subWalletInfo.WalletType == 3 {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Because this subwallet is designated to be an bulk-payment wallet, %v %v will be deducted from your primary wallet and be used to activate it and the custodial signer. Please note that bulk-payment wallets cannot be used to send normal payments, but only be used by internal system to disburse bulk payments on your behalf.", (activationAmount.Add(decimal.RequireFromString(os.Getenv("WALLET_SIGNER_ACTIVATION_AMOUNT")))).String(), os.Getenv("NATIVE_ASSET_CODE")))

	}
	if subWalletInfo.WalletType == 2 || subWalletInfo.WalletType == 3 {
		ops = append(ops, &basetxn.SetOptions{
			LowThreshold:    basetxn.NewThreshold(uint32(3)),
			MediumThreshold: basetxn.NewThreshold(uint32(3)),
			HighThreshold:   basetxn.NewThreshold(uint32(3)),
			SourceAccount:   subWalletInfo.Address,
		})
	}

	_, _, channelSourceAccountNativeBalance, _, channelSourceAccount, channelSourceAccountErr := network.BlockchainAccountProperties(client, subWalletInfo.ChannelAccount, basetxn.NativeAsset{})
	if channelSourceAccountErr != nil || (channelSourceAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
		log.Printf("[generateSubWalletXdrWithChannelAccount] by [%v] for [%v] Channel Account underfunded.\n", user.Username, subWalletInfo.Address)

		err = &tErrors.CustomError{
			Param:      "channelAccount",
			Err:        "error-channel-account-underfunded",
			ErrMessage: fmt.Sprintf("Channel account does not have minimum %v balance required to complete this operation", os.Getenv("NATIVE_ASSET_CODE")),
			Code:       400,
		}
		return "", subWalletObj, linkedWallet, err
	}
	// Construct the transaction that holds the operations to execute on the network
	tx, err := basetxn.NewTransaction(
		basetxn.TransactionParams{
			SourceAccount:        channelSourceAccount.Address,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              2000,
			Memo:                 "Create Sub Wallet",
		},
	)
	if err != nil {
		log.Println("[generateSubWalletXdrWithChannelAccount] error constructing transaction ", err)
		return "", subWalletObj, linkedWallet, err
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateSubWalletXdrWithChannelAccount] error getting txn base64", err)
		return "", subWalletObj, linkedWallet, err
	}

	return xdrBase64, subWalletObj, linkedWallet, nil

}

func SubmitSubWalletXdrWithSignatures(client *ethclient.Client, signatures map[string]string, xdrBase64 string) (string, error) {
	txnHash, err := network.SubmitXdrWithSignatures(client, xdrBase64, signatures, nil)
	if err != nil {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error subwallet activation failed", ErrMessage: "SubWallet Failed", Code: 500}
	}
	return txnHash, nil
}

func SubmitSubWalletXdrWithSignature(client *ethclient.Client, accountAddress, signerAddress, subWalletAddress string, xdrBase64 string, primarySignature, subWalletSignature string, subWalletMustSign int) (string, error) {
	signatures := map[string]string{signerAddress: primarySignature}
	if subWalletMustSign == 1 {
		signatures[subWalletAddress] = subWalletSignature
	}
	txnHash, err := network.SubmitXdrWithSignatures(client, xdrBase64, signatures, nil)
	if err != nil {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error subwallet activation failed", ErrMessage: "SubWallet Failed", Code: 500}
	}
	return txnHash, nil
}

func SubmitSubWalletXdrForChannelAccountWithSignature(client *ethclient.Client, ownerAddress, subWalletAddress, channelPK string, xdrBase64 string, primarySignature, subWalletSignature, channelAccountSignature string, subWalletMustSign int) (string, error) {
	signatures := map[string]string{
		ownerAddress:     primarySignature,
		subWalletAddress: subWalletSignature,
	}
	if subWalletMustSign == 1 {
		signatures[channelPK] = channelAccountSignature
	}
	txnHash, err := network.SubmitXdrWithSignatures(client, xdrBase64, signatures, nil)
	if err != nil {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error subwallet activation failed", ErrMessage: "SubWallet Failed", Code: 500}
	}
	return txnHash, nil
}

func GetUserWallets(user userModels.User, gc *sharedconfig.GlobalConfig) (wallets []userModels.UserWallet) {
	wallets = make([]userModels.UserWallet, 0)
	gc.DB.Where("user_id = ?", user.ID).Find(&wallets)
	return

}
