package services

import (
	"context"
	"fmt"
	"strings"

	vaultsignermodels "admin-panel-dashboard/internal/components/vaultsigner/models"
	"admin-panel-dashboard/internal/components/vaultsigner/vaultclient"
	"admin-panel-dashboard/internal/evmkeypair"
	"admin-panel-dashboard/internal/gnosissafe"
	"admin-panel-dashboard/internal/network"

	"github.com/ethereum/go-ethereum/common"
	vaultapi "github.com/hashicorp/vault/api"
)

// SwapResult is what the handler needs to build the PUT /me/vault-signer/secrets/:secretId
// response (Section 6) and write the audit log row.
type SwapResult struct {
	VaultVersionBefore int
	VaultVersionAfter  int
	Swapped            bool // true if an on-chain swap was attempted
	BaseTxHash         string
	BaseTxStatus       string // "success" | "failed", only meaningful if Swapped
	SwapError          error  // non-nil if Swapped and submission failed
}

// SwapSigner is the full synchronous flow for PUT /me/vault-signer/secrets/:secretId
// on an active index (Section 5d). Ownership/guard checks (is targetIndex
// assigned to the caller) are the handler's responsibility — this function
// assumes the caller is already authorized and focuses purely on the
// Vault/on-chain mechanics.
//
// The managed secret's WalletAddress is expected to be a deployed Safe
// (Safe{Wallet}, formerly Gnosis Safe) contract's address: Base EOAs have
// no native multisig/weight concept the way a Stellar account did, so
// "swap this wallet's signer" now means calling the Safe's own
// swapOwner(prevOwner, oldOwner, newOwner) through execTransaction, with
// enough of the Safe's other current owners co-signing to meet its
// threshold - see internal/gnosissafe's package doc.
func SwapSigner(ctx context.Context, vc *vaultapi.Client, secret vaultsignermodels.VaultSignerManagedSecret, targetIndex int, newValueRaw string) (SwapResult, error) {
	// Step 2 (Section 5d/original 4d): validate the new key through the full
	// Section 5c pipeline.
	newAddressStr, err := ValidateSignerValue(newValueRaw)
	if err != nil {
		return SwapResult{}, err
	}
	newValueTrimmed := strings.TrimSpace(newValueRaw)

	// Step 3: validate the old key currently at targetIndex — structural
	// check only, no on-chain re-check.
	currentValue, _, err := vaultclient.ReadValue(ctx, vc, secret, targetIndex)
	if err != nil {
		return SwapResult{}, err
	}
	oldValueTrimmed := strings.TrimSpace(currentValue)
	oldKP, oldErr := evmkeypair.ParseFull(oldValueTrimmed)

	// Step 4: if the old value isn't a valid keypair, or new == old, skip the
	// on-chain swap entirely and fall through to a plain Vault write.
	if oldErr != nil || oldValueTrimmed == newValueTrimmed {
		versionBefore, versionAfter, writeErr := vaultclient.WriteAtIndex(ctx, vc, secret, targetIndex, newValueTrimmed)
		if writeErr != nil {
			return SwapResult{}, writeErr
		}
		return SwapResult{VaultVersionBefore: versionBefore, VaultVersionAfter: versionAfter, Swapped: false}, nil
	}

	oldAddress := common.HexToAddress(oldKP.Address())
	newAddress := common.HexToAddress(newAddressStr)

	// Step 5: fetch the Safe's current owner set and confirm the old key is
	// actually still one of them - the Base equivalent of reading the
	// wallet's current signer weight off the Stellar account.
	client := network.GetBlockchainClient()
	safe, err := gnosissafe.New(client, secret.WalletAddress)
	if err != nil {
		return SwapResult{}, fmt.Errorf("resolving safe address: %w", err)
	}
	owners, err := safe.Owners(ctx)
	if err != nil {
		return SwapResult{}, fmt.Errorf("fetching safe owners: %w", err)
	}
	prevOwner, err := gnosissafe.PrevOwner(owners, oldAddress)
	if err != nil {
		return SwapResult{}, fmt.Errorf("old signer %s is not currently registered on safe %s: %w", oldAddress.Hex(), secret.WalletAddress, err)
	}

	// Step 6: update Vault first, per the established (accepted-risk) order.
	versionBefore, versionAfter, err := vaultclient.WriteAtIndex(ctx, vc, secret, targetIndex, newValueTrimmed)
	if err != nil {
		return SwapResult{}, err
	}

	// Step 7: gather signing keys — every active index except targetIndex,
	// read from the now-updated CSV. The new value is deliberately not used
	// to sign — it isn't a registered Safe owner yet.
	signers, err := GatherActiveSigners(ctx, vc, secret, targetIndex)
	if err != nil {
		return SwapResult{VaultVersionBefore: versionBefore, VaultVersionAfter: versionAfter}, err
	}

	// Step 8-9: build the swapOwner call, collect one signature per
	// gathered signer over the Safe's own getTransactionHash, and submit.
	innerCalldata, err := gnosissafe.SwapOwnerCalldata(prevOwner, oldAddress, newAddress)
	if err != nil {
		return SwapResult{VaultVersionBefore: versionBefore, VaultVersionAfter: versionAfter}, fmt.Errorf("encoding swapOwner call: %w", err)
	}
	result := SwapResult{VaultVersionBefore: versionBefore, VaultVersionAfter: versionAfter, Swapped: true}
	execResult, submitErr := submitSafeOwnerChange(ctx, safe, innerCalldata, signers)
	if submitErr != nil {
		result.BaseTxStatus = "failed"
		result.SwapError = submitErr
		return result, nil // Vault already updated — the mismatch is a documented, accepted risk (Section 5d)
	}
	result.BaseTxHash = execResult.TxHash
	if execResult.Success {
		result.BaseTxStatus = "success"
	} else {
		result.BaseTxStatus = "failed"
		result.SwapError = fmt.Errorf("execTransaction for swapOwner reverted on-chain (tx %s)", execResult.TxHash)
	}
	return result, nil
}

// GatherActiveSigners reads the now-updated CSV and parses every active
// index (0..ActiveSigningCount-1) except targetIndex into a signing keypair.
func GatherActiveSigners(ctx context.Context, vc *vaultapi.Client, secret vaultsignermodels.VaultSignerManagedSecret, targetIndex int) ([]*evmkeypair.Full, error) {
	var signers []*evmkeypair.Full
	for i := 0; i < secret.ActiveSigningCount; i++ {
		if i == targetIndex {
			continue
		}
		value, _, err := vaultclient.ReadValue(ctx, vc, secret, i)
		if err != nil {
			return nil, fmt.Errorf("reading active signer at index %d: %w", i, err)
		}
		kp, err := evmkeypair.ParseFull(strings.TrimSpace(value))
		if err != nil {
			return nil, fmt.Errorf("active signer at index %d is not a valid keypair: %w", i, err)
		}
		signers = append(signers, kp)
	}
	return signers, nil
}

// submitSafeOwnerChange is the shared "gather one signature per active
// signer over the Safe's current transaction hash, then submit and wait
// for the receipt" tail end SwapSigner and DeleteAssignment's full path
// both need, given an already-ABI-encoded swapOwner/removeOwner call. The
// first gathered signer pays gas and is the execTransaction's msg.sender —
// execTransaction doesn't require the caller to be an owner, only that
// signatures satisfy the Safe's threshold, and Vault-managed signer keys
// are already expected to hold enough Base ETH to act as transaction
// senders elsewhere in this system.
func submitSafeOwnerChange(ctx context.Context, safe *gnosissafe.Safe, innerCalldata []byte, signers []*evmkeypair.Full) (gnosissafe.ExecResult, error) {
	if len(signers) == 0 {
		return gnosissafe.ExecResult{}, fmt.Errorf("no active signers available to co-sign or submit this change")
	}

	nonce, err := safe.Nonce(ctx)
	if err != nil {
		return gnosissafe.ExecResult{}, fmt.Errorf("fetching safe nonce: %w", err)
	}
	safeTxHash, err := safe.TransactionHash(ctx, safe.Address, innerCalldata, nonce)
	if err != nil {
		return gnosissafe.ExecResult{}, fmt.Errorf("computing safe transaction hash: %w", err)
	}

	sigs := make(map[common.Address][]byte, len(signers))
	for _, signer := range signers {
		sig, sigErr := gnosissafe.SignTransactionHash(signer.PrivateKey(), safeTxHash)
		if sigErr != nil {
			return gnosissafe.ExecResult{}, fmt.Errorf("signing safe transaction hash: %w", sigErr)
		}
		sigs[common.HexToAddress(signer.Address())] = sig
	}
	signatures := gnosissafe.ConcatSignatures(sigs)

	return safe.SubmitOwnerChange(ctx, signers[0].PrivateKey(), innerCalldata, signatures, network.GetBlockchainChainID())
}
