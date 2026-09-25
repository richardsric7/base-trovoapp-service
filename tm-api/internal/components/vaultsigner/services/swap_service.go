package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	vaultsignermodels "admin-panel-dashboard/internal/components/vaultsigner/models"
	"admin-panel-dashboard/internal/components/vaultsigner/vaultclient"
	tErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/network"

	"github.com/ecnepsnai/discord"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

// SwapResult is what the handler needs to build the PUT /me/vault-signer/secrets/:secretId
// response (Section 6) and write the audit log row.
type SwapResult struct {
	VaultVersionBefore int
	VaultVersionAfter  int
	Swapped            bool // true if an on-chain swap was attempted
	StellarTxHash      string
	StellarTxStatus    string // "success" | "failed", only meaningful if Swapped
	SwapError          error  // non-nil if Swapped and submission failed
}

// SwapSigner is the full synchronous flow for PUT /me/vault-signer/secrets/:secretId
// on an active index (Section 5d). Ownership/guard checks (is targetIndex
// assigned to the caller) are the handler's responsibility — this function
// assumes the caller is already authorized and focuses purely on the
// Vault/on-chain mechanics.
func SwapSigner(ctx context.Context, vc *vaultapi.Client, secret vaultsignermodels.VaultSignerManagedSecret, targetIndex int, newValueRaw string) (SwapResult, error) {
	// Step 2 (Section 5d/original 4d): validate the new key through the full
	// Section 5c pipeline.
	newPublicKey, err := ValidateSignerValue(newValueRaw)
	if err != nil {
		return SwapResult{}, err
	}
	newValueTrimmed := strings.TrimSpace(newValueRaw)

	// Step 3: validate the old key currently at targetIndex — structural
	// check only, no Horizon re-check.
	currentValue, _, err := vaultclient.ReadValue(ctx, vc, secret, targetIndex)
	if err != nil {
		return SwapResult{}, err
	}
	oldValueTrimmed := strings.TrimSpace(currentValue)
	oldKP, oldErr := keypair.ParseFull(oldValueTrimmed)

	// Step 4: if the old value isn't a valid keypair, or new == old, skip the
	// on-chain swap entirely and fall through to a plain Vault write.
	if oldErr != nil || oldValueTrimmed == newValueTrimmed {
		versionBefore, versionAfter, writeErr := vaultclient.WriteAtIndex(ctx, vc, secret, targetIndex, newValueTrimmed)
		if writeErr != nil {
			return SwapResult{}, writeErr
		}
		return SwapResult{VaultVersionBefore: versionBefore, VaultVersionAfter: versionAfter, Swapped: false}, nil
	}

	// Step 5: fetch the wallet's current on-chain state — the old key's
	// current signer weight (preserving threshold math) and sequence number.
	client := network.GetBlockchainClient()
	account, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: secret.WalletPublicKey})
	if err != nil {
		return SwapResult{}, fmt.Errorf("fetching wallet account detail: %w", err)
	}
	oldPublicKey := oldKP.Address()
	var oldWeight int32
	found := false
	for _, s := range account.Signers {
		if s.Key == oldPublicKey {
			oldWeight = int32(s.Weight)
			found = true
			break
		}
	}
	if !found {
		return SwapResult{}, fmt.Errorf("old signer %s is not currently registered on wallet %s", oldPublicKey, secret.WalletPublicKey)
	}

	// Step 6: update Vault first, per the established (accepted-risk) order.
	versionBefore, versionAfter, err := vaultclient.WriteAtIndex(ctx, vc, secret, targetIndex, newValueTrimmed)
	if err != nil {
		return SwapResult{}, err
	}

	// Step 7: gather signing keys — every active index except targetIndex,
	// read from the now-updated CSV. The new value is deliberately not used
	// to sign — it isn't a registered on-chain signer yet.
	signers, err := GatherActiveSigners(ctx, vc, secret, targetIndex)
	if err != nil {
		return SwapResult{VaultVersionBefore: versionBefore, VaultVersionAfter: versionAfter}, err
	}

	// Step 8: build the two-operation SetOptions transaction.
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		BaseFee:              txnbuild.MinBaseFee,
		Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewTimeout(300)},
		Operations: []txnbuild.Operation{
			&txnbuild.SetOptions{Signer: &txnbuild.Signer{Address: oldPublicKey, Weight: 0}},
			&txnbuild.SetOptions{Signer: &txnbuild.Signer{Address: newPublicKey, Weight: txnbuild.Threshold(oldWeight)}},
		},
	})
	if err != nil {
		return SwapResult{VaultVersionBefore: versionBefore, VaultVersionAfter: versionAfter}, fmt.Errorf("building swap transaction: %w", err)
	}

	// Step 9: sign with every gathered active signer except targetIndex.
	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), signers...)
	if err != nil {
		return SwapResult{VaultVersionBefore: versionBefore, VaultVersionAfter: versionAfter}, fmt.Errorf("signing swap transaction: %w", err)
	}

	// Step 10: submit, following this repo's established Horizon submission
	// error-handling shape (internal/network/main.go's SubmitXdrWithSignature*).
	result := SwapResult{VaultVersionBefore: versionBefore, VaultVersionAfter: versionAfter, Swapped: true}
	txResp, submitErr := SubmitSetOptionsTransaction(client, tx)
	if submitErr != nil {
		result.StellarTxStatus = "failed"
		result.SwapError = submitErr
		return result, nil // Vault already updated — the mismatch is a documented, accepted risk (Section 5d)
	}
	result.StellarTxHash = txResp.Hash
	result.StellarTxStatus = "success"
	return result, nil
}

// GatherActiveSigners reads the now-updated CSV and parses every active
// index (0..ActiveSigningCount-1) except targetIndex into a signing keypair.
func GatherActiveSigners(ctx context.Context, vc *vaultapi.Client, secret vaultsignermodels.VaultSignerManagedSecret, targetIndex int) ([]*keypair.Full, error) {
	var signers []*keypair.Full
	for i := 0; i < secret.ActiveSigningCount; i++ {
		if i == targetIndex {
			continue
		}
		value, _, err := vaultclient.ReadValue(ctx, vc, secret, i)
		if err != nil {
			return nil, fmt.Errorf("reading active signer at index %d: %w", i, err)
		}
		kp, err := keypair.ParseFull(strings.TrimSpace(value))
		if err != nil {
			return nil, fmt.Errorf("active signer at index %d is not a valid keypair: %w", i, err)
		}
		signers = append(signers, kp)
	}
	return signers, nil
}

// SubmitSetOptionsTransaction submits a freshly-built SetOptions transaction,
// mirroring internal/network/main.go's SubmitXdrWithSignature* error
// handling (horizonclient.Error inspection, Discord alert on connectivity
// failure) even though this submits a transaction this service itself built,
// not a pre-signed XDR from a caller.
func SubmitSetOptionsTransaction(client *horizonclient.Client, tx *txnbuild.Transaction) (result horizon.Transaction, err error) {
	xdrBase64, err := tx.Base64()
	if err != nil {
		return result, fmt.Errorf("encoding swap transaction: %w", err)
	}

	txnResult, submitErr := client.SubmitTransactionXDR(xdrBase64)
	if submitErr != nil {
		if strings.Contains(submitErr.Error(), "timeout") || strings.Contains(submitErr.Error(), "handshake") ||
			strings.Contains(submitErr.Error(), "read tcp") || strings.Contains(submitErr.Error(), "connection reset by peer") ||
			strings.Contains(submitErr.Error(), "dial tcp") || strings.Contains(submitErr.Error(), "no such host") {
			discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
			if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
				discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
			}
			if sayErr := discord.Say(fmt.Sprintf("[vaultsigner.SwapSigner] error connecting to expansion service: %v\nXDR: %v", submitErr, xdrBase64)); sayErr != nil {
				log.Println("[vaultsigner.SwapSigner] discord alert failed:", sayErr)
			}
		}

		if horizonException, ok := submitErr.(*horizonclient.Error); ok {
			for key, val := range horizonException.Problem.Extras {
				log.Printf("[vaultsigner.SwapSigner] Extras: %v is %v\n", key, val)
			}
			if resultCodes, codesErr := horizonException.ResultCodes(); codesErr == nil {
				for key, val := range resultCodes.OperationCodes {
					log.Printf("[vaultsigner.SwapSigner] Result code: %v is %v\n", key, val)
				}
			}
		} else {
			log.Printf("[vaultsigner.SwapSigner] not horizon error: %v\n", submitErr)
		}

		return result, &tErrors.CustomError{Param: "publicKey", Err: "error operation failed", ErrMessage: "Operation Failed", Code: 500}
	}

	return txnResult, nil
}
