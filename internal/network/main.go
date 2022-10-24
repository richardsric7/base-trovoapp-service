package network

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	tErrors "trovo-wallet-api/internal/errors"

	"github.com/ecnepsnai/discord"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var kTempAccountSalt string = "j4rkTZQ2mLk3NAhK"

func GetBlockchainNetworkPassPhrase() string {
	return os.Getenv("BLOCKCHAIN_NETWORK_PASSPHRASE")
}

func GetBlockchainBaseReserve() decimal.Decimal {
	val, err := decimal.NewFromString(os.Getenv("BLOCKCHAIN_BASE_RESERVE"))

	if err != nil {
		return decimal.NewFromInt(1)
	}

	return val

}

func GetBlockchainSwapDestinationMin() decimal.Decimal {
	val, err := decimal.NewFromString(os.Getenv("BLOCKCHAIN_SWAP_DESTINATION_MIN"))

	if err != nil {

		return decimal.RequireFromString("0.0000500")
	}

	return val

}

func GetBlockchainClient() *horizonclient.Client {
	var client *horizonclient.Client = horizonclient.DefaultPublicNetClient
	client.HorizonURL = os.Getenv("EXPANSION_URL")
	return client
}

func TempAccountKeypair(publicKey string) (*keypair.Full, error) {

	mnemonic := os.Getenv("MNEMONIC_TEMP_ACCOUNTS")

	h := sha256.New()
	h.Write([]byte(kTempAccountSalt))
	h.Write([]byte(mnemonic))
	h.Write([]byte(publicKey))

	hashed := h.Sum(nil)

	var rawSeed [32]byte
	copy(rawSeed[:], hashed[0:32])

	return keypair.FromRawSeed([32]byte(rawSeed))

}

// BlockchainAccountProperties returns whether the account exists , whether the account trusts the asset, the native balance, the asset balance, the account object, and whether there was an error
func BlockchainAccountProperties(client *horizonclient.Client, destinationPublicKey string, asset txnbuild.Asset) (bool, bool, decimal.Decimal, decimal.Decimal, *horizon.Account, error) {

	log.Printf("[BlockchainAccountProperties] obtaining blockchain account properties for  %v \n", destinationPublicKey)

	destinationAccountExists := true
	destinationAccountTrustsAsset := false

	if asset.GetIssuer() == destinationPublicKey {
		destinationAccountTrustsAsset = true

	}

	destinationKeyPair, _ := keypair.ParseAddress(destinationPublicKey)

	destinationAccountRequest := horizonclient.AccountRequest{AccountID: destinationKeyPair.Address()}

	destinationAccountDetail, destinationHorizonError := client.AccountDetail(destinationAccountRequest)

	if destinationHorizonError != nil {
		log.Print(destinationHorizonError)
		userNotFoundError := false

		horizonException, ok := destinationHorizonError.(*horizonclient.Error)

		if ok {

			log.Print("[BlockchainAccountProperties] error is known ", horizonException.Problem.Status)

			if horizonException.Problem.Status == http.StatusNotFound {
				userNotFoundError = true
				destinationAccountExists = false
				return destinationAccountExists, destinationAccountTrustsAsset, decimal.Zero, decimal.Zero, &destinationAccountDetail, nil
			}
		}

		if !userNotFoundError {
			return destinationAccountExists, destinationAccountTrustsAsset, decimal.Zero, decimal.Zero, &destinationAccountDetail, &tErrors.ErrorTemporaryServerError{}
		}
	}

	//check if destinationAccount trusts asset

	balances := destinationAccountDetail.Balances

	nativeAccountBalance := decimal.Zero
	customAssetAccountBalance := decimal.Zero

	subEntryCount := decimal.NewFromInt32(destinationAccountDetail.SubentryCount)
	subEntryCountMultiplier := GetBlockchainBaseReserve()

	amountToSubtractFromNativeAccountBalance := (subEntryCount).Mul(subEntryCountMultiplier)

	log.Printf("[BlockchainAccountProperties] amount to subtract is %v\n", amountToSubtractFromNativeAccountBalance.String())

	for _, balance := range balances {
		_asset := balance.Asset

		if _asset.Code == "" {

			nativeAccountBalance = decimal.RequireFromString(balance.Balance).Sub(amountToSubtractFromNativeAccountBalance).Truncate(7)
		}

		if _asset.Code == asset.GetCode() && _asset.Issuer == asset.GetIssuer() {
			destinationAccountTrustsAsset = true

			if _asset.Code != "" {
				customAssetAccountBalance = decimal.RequireFromString(balance.Balance)
			}
		}
	}

	return destinationAccountExists, destinationAccountTrustsAsset, nativeAccountBalance, customAssetAccountBalance, &destinationAccountDetail, nil
}

func SubmitXdrWithSignature(client *horizonclient.Client, signerPublicKey string, xdrBase64 string, signature string) (string, error) {
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

	txn, err = txn.AddSignatureBase64(GetBlockchainNetworkPassPhrase(), signerPublicKey, signature)

	if err != nil {
		log.Println("[SubmitXdrWithSignature] Failed to verify signature on [", GetBlockchainNetworkPassPhrase(), "] and [", signature, "] for [", xdrBase64, "] and public key ", signerPublicKey, ", error [", err, "]")

		return "", err
	}

	xdrBase64, err = txn.Base64()

	if err != nil {
		log.Println(err)
		return "", err
	}

	// log.Println("signed xdr is " + xdrBase64)

	txnResult, err := client.SubmitTransactionXDR(xdrBase64)

	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			discord.Say(fmt.Sprintf("[SubmitXdrWithSignature] error connecting to expansion service: %v\nXDR: %v", err, xdrBase64))
		}

		horizonException, ok := err.(*horizonclient.Error)

		if ok {

			extraErrors := horizonException.Problem.Extras

			for key, val := range extraErrors {
				log.Printf("[SubmitXdrWithSignature] Extras: %v is %v\nOwner publicKey: %v\n", key, val, signerPublicKey)
				logDiscordFailedPayment(fmt.Sprintf("[SubmitXdrWithSignature] Extras: %v is %v\nSigner publicKey: %v\n", key, val, signerPublicKey))

			}

			resultCodes, errRes := horizonException.ResultCodes()
			if errRes == nil {
				for key, val := range resultCodes.OperationCodes {
					log.Printf("[SubmitXdrWithSignature] Result code: %v is %v\nSigner publicKey: %v\n", key, val, signerPublicKey)
					logDiscordFailedPayment(fmt.Sprintf("[SubmitXdrWithSignature] Result code: %v is %v\nSigner publicKey: %v\n", key, val, signerPublicKey))

				}
			} else {
				log.Printf("[SubmitXdrWithSignature] Error getting result codes: %v\n", errRes)
			}

		} else {
			log.Printf("[SubmitXdrWithSignature] not horizon error: %v\n", err)

		}

		return "", &tErrors.CustomError{Param: "publicKey", Err: "error operation failed", ErrMessage: "Operation Failed", Code: 500}

	}

	return txnResult.Hash, nil

}
func SubmitApprovalsXdrWithSignatures(client *horizonclient.Client, approvalID string, db *gorm.DB) (string, error) {
	type PendingTransactionSignature struct {
		CreatedAt                time.Time `json:"createdAt"`
		ID                       string    `gorm:"size:56"`
		PendingAuthID            string    `gorm:"size:56;not null;index:idx_pending_trxsig_pending_auth,unique"`
		Approver                 string    `gorm:"size:20;not null;index:idx_pending_trxsig_pending_auth,unique;index:idx_pending_trxsig_approver" json:"approver"`
		ApproverSignerPublicKey  string    `gorm:"size:56;not null;index:idx_pending_trxsig_approver_signer" json:"approverSignerPublicKey"`
		TransactionWithSignature string    `gorm:"not null" json:"transactionWithSignature"`
	}
	type PendingAuth struct {
		CreatedAt                    time.Time                     `json:"createdAt"`
		UpdatedAt                    time.Time                     `gorm:"default:now()" json:"updatedAt"`
		ID                           string                        `gorm:"size:56" json:"id"`
		Initiator                    string                        `gorm:"size:20;not null;index:idx_pending_auth_initiator" json:"initiator"`
		InitiatorSignerPublicKey     string                        `gorm:"size:56;not null;index:idx_pending_auth_signer_public_key" json:"initiatorSignerPublicKey"`
		WalletPublicKey              string                        `gorm:"size:56;not null;index:idx_pending_auth_wallet_public_key" json:"walletPublicKey"`
		TransactionType              string                        `gorm:"size:28;not null;index:idx_pending_auth_transaction_type" json:"transactionType"`
		Description                  string                        `gorm:"not null;" json:"description"`
		ApprovalsNeeded              int                           `gorm:"not null;" json:"approvalsNeeded"`
		ApprovalsGotten              int                           `gorm:"not null;default:0" json:"approvalsGotten"`
		TransactionStatus            string                        `gorm:"size:20;not null;default:'PENDING';index:idx_pending_auth_transaction_status" json:"transactionStatus"`
		RejectedBy                   *string                       `gorm:"size:20;null;index:idx_pending_auth_rejected_by" json:"rejectedBy"`
		ReasonForRejection           *string                       `gorm:"size:200;null;" json:"reasonForRejection"`
		TransactionXdr               string                        `gorm:"not null;" json:"transactionXdr"`
		TransactionInfoStr           *string                       `gorm:"null;" json:"transactionInfoStr"`
		TransactionID                *string                       `gorm:"size:70;null;index:idx_pending_auth_transaction_id" json:"transactionID"`
		PendingTransactionSignatures []PendingTransactionSignature `json:"pendingTransactionSignatures"`
	}
	var pendingAuth PendingAuth
	e := db.Preload(clause.Associations).Where("id = ?", approvalID).First(&pendingAuth).Error
	if e != nil {

		return "", &tErrors.ErrorTemporaryServerError{}
	}

	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFb"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}

	gTxn, err := txnbuild.TransactionFromXDR(pendingAuth.TransactionXdr)

	if err != nil {
		return "", err
	}

	txn, ok := gTxn.Transaction()

	if !ok {
		return "", &tErrors.ErrorInvalidTransaction{}
	}
	for _, sig := range pendingAuth.PendingTransactionSignatures {

		txn, err = txn.AddSignatureBase64(GetBlockchainNetworkPassPhrase(), sig.ApproverSignerPublicKey, sig.TransactionWithSignature)

		if err != nil {
			log.Printf("[SubmitApprovalXdrWithSignature] Failed to verify signature of [%v] on [%v] and [%v] for [%v] on signerPublicKey [%v], error: [%v]\n", sig.Approver, GetBlockchainNetworkPassPhrase(), sig.TransactionWithSignature, pendingAuth.TransactionXdr, sig.ApproverSignerPublicKey, err)
			return "", err
		}
	}

	xdrBase64, err := txn.Base64()

	if err != nil {
		log.Println(err)
		return "", err
	}

	// log.Println("signed xdr is " + xdrBase64)

	txnResult, err := client.SubmitTransactionXDR(xdrBase64)

	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			discord.Say(fmt.Sprintf("[SubmitApprovalXdrWithSignature] error connecting to expansion service: %v\nXDR: %v", err, xdrBase64))
		}

		horizonException, ok := err.(*horizonclient.Error)

		if ok {

			extraErrors := horizonException.Problem.Extras

			for key, val := range extraErrors {
				log.Printf("[SubmitApprovalXdrWithSignature] Extras: %v is %v\nWallet publicKey: %v\n", key, val, pendingAuth.WalletPublicKey)
				logDiscordFailedPayment(fmt.Sprintf("[SubmitApprovalXdrWithSignature] Extras: %v is %v\nWallet publicKey: %v\n", key, val, pendingAuth.WalletPublicKey))

			}

			resultCodes, errRes := horizonException.ResultCodes()
			if errRes == nil {
				for key, val := range resultCodes.OperationCodes {
					log.Printf("[v] Result code: %v is %v\nWallet publicKey: %v\n", key, val, pendingAuth.WalletPublicKey)
					logDiscordFailedPayment(fmt.Sprintf("[SubmitApprovalXdrWithSignature] Result code: %v is %v\nWallet publicKey: %v\n", key, val, pendingAuth.WalletPublicKey))

				}
			} else {
				log.Printf("[SubmitApprovalXdrWithSignature] Error getting result codes: %vzz\n", errRes)
			}

		} else {
			log.Printf("[SubmitApprovalXdrWithSignature] not horizon error: %v\n", err)

		}

		return "", &tErrors.CustomError{Param: "publicKey", Err: "error operation failed", ErrMessage: "Operation Failed", Code: 500}

	}

	return txnResult.Hash, nil

}

func SubmitXdrWithSignatureChannelAccounts(client *horizonclient.Client, signerPublicKey, channelAccountPK string, xdrBase64 string, txSignature, chanSignature string) (string, error) {

	gTxn, err := txnbuild.TransactionFromXDR(xdrBase64)

	if err != nil {
		return "", err
	}

	txn, ok := gTxn.Transaction()

	if !ok {
		return "", &tErrors.ErrorInvalidTransaction{}
	}
	txn, err = txn.AddSignatureBase64(GetBlockchainNetworkPassPhrase(), channelAccountPK, chanSignature)

	if err != nil {
		log.Println("Failed to verify signature with channelAccount Public Key on [", GetBlockchainNetworkPassPhrase(), "] and [", chanSignature, "] for [", xdrBase64, "] and public key ", channelAccountPK)
		// chantx, _ := txnbuild.TransactionFromXDR(chanSignature)
		// chtxn, _ := chantx.Transaction()
		// log.Printf("<<<<<Decorated signatures>>>>>>>>>>>>>>>>>>>>>>>>\n\n%+v\n\n", chtxn.Signatures())
		return "", err
	}
	txn, err = txn.AddSignatureBase64(GetBlockchainNetworkPassPhrase(), signerPublicKey, txSignature)

	if err != nil {
		log.Println("Failed to verify signature with sender public key on [", GetBlockchainNetworkPassPhrase(), "] and [", txSignature, "] for [", xdrBase64, "] and singer public key ", signerPublicKey)

		return "", err
	}

	xdrBase64, err = txn.Base64()

	if err != nil {
		log.Println(err)
		return "", err
	}

	// log.Println("signed xdr is " + txSignature)

	// txnResult, err := client.SubmitTransactionXDR(txSignature)
	log.Println("signed xdr is " + xdrBase64)

	txnResult, err := client.SubmitTransactionXDR(xdrBase64)

	if err != nil {

		horizonException, ok := err.(*horizonclient.Error)

		if ok {

			extraErrors := horizonException.Problem.Extras

			for key, val := range extraErrors {
				log.Printf("Extras: %v is %v\n", key, val)
			}

			resultCodes, _ := horizonException.ResultCodes()

			for key, val := range resultCodes.OperationCodes {
				log.Printf("Result code: %v is %v\n", key, val)
			}

		} else {
			log.Printf("[SubmitXdrWithSignatureChannelAccounts] not horizon error: %v\n", err)

		}

		return "", err

	}

	return txnResult.Hash, nil

}
func logDiscordFailedPayment(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}
