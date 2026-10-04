package safesigner

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/network"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// multiSendCallOnlyAddress returns Safe's MultiSendCallOnly deployment
// (gnosissafe.DefaultMultiSendCallOnlyAddress unless overridden via
// P2P_ESCROW_MULTISEND_ADDRESS, e.g. for a newer Safe version).
func multiSendCallOnlyAddress() string {
	if addr := strings.TrimSpace(os.Getenv("P2P_ESCROW_MULTISEND_ADDRESS")); addr != "" {
		return addr
	}
	return gnosissafe.DefaultMultiSendCallOnlyAddress
}

// Transfer is one outgoing transfer the Safe should execute as part of a
// settlement release (Plan Section 60).
type Transfer struct {
	Token     string // ERC20 contract address, or "" for the native asset
	Recipient string
	Amount    *big.Int
}

// SafeAddress returns the escrow Safe's own address. It is the same address
// deposits are paid into (P2P_ESCROW_WALLET_ADDRESS) - a Gnosis Safe is
// itself a contract address that can receive tokens directly, so one env
// var serves both purposes.
func SafeAddress() string {
	return os.Getenv("P2P_ESCROW_WALLET_ADDRESS")
}

var erc20ABI abi.ABI

func init() {
	var err error
	erc20ABI, err = abi.JSON(strings.NewReader(`[
		{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"}
	]`))
	if err != nil {
		panic("safesigner: invalid erc20 ABI: " + err.Error())
	}
}

// ExecuteTransfers assembles, signs (with the first 3 P2P_ESCROW_SIGNERS),
// and submits transfers as Gnosis Safe execTransaction(s).
//
// When every non-zero transfer moves the same ERC-20 token (Plan Section
// 60's usual case: buyer net amount + up to 3 fee wallets, all in
// order.AssetContractAddress), they are batched into ONE Safe
// execTransaction via the canonical MultiSendCallOnly contract - one Safe
// nonce, one signature round, one Base transaction/gas bill, and the legs
// become atomic with each other (all succeed or the whole batch reverts).
// Otherwise (a mix of tokens, or a native-asset leg, or just one transfer)
// each transfer is submitted as its own sequential Safe transaction (own
// nonce, own signatures), as before - simpler for the cases batching
// doesn't apply to, at the cost of those legs not being atomic with each
// other. Returns the transaction hash of the single batched transaction,
// or of the last transfer submitted in the sequential fallback - which
// Order.AssetReleaseTransactionHash stores as the canonical release hash.
func ExecuteTransfers(transfers []Transfer) (string, error) {
	signers, err := ActiveSigners()
	if err != nil {
		return "", err
	}
	safeAddress := SafeAddress()
	if safeAddress == "" {
		return "", fmt.Errorf("P2P_ESCROW_WALLET_ADDRESS is not configured")
	}
	client := network.GetBlockchainClient()
	chainID := network.GetBlockchainChainID()

	nonZero := make([]Transfer, 0, len(transfers))
	for _, t := range transfers {
		if t.Amount != nil && t.Amount.Sign() > 0 {
			nonZero = append(nonZero, t)
		}
	}
	if len(nonZero) == 0 {
		return "", fmt.Errorf("no non-zero transfers to execute")
	}

	safe := common.HexToAddress(safeAddress)
	multiSend := common.HexToAddress(multiSendCallOnlyAddress())
	ctx := context.Background()

	if sameERC20Token(nonZero) {
		calls := make([]gnosissafe.Call, 0, len(nonZero))
		for _, t := range nonZero {
			c, err := transferCall(t)
			if err != nil {
				return "", err
			}
			calls = append(calls, c)
		}
		return gnosissafe.ExecCalls(ctx, client, chainID, safe, signers, calls, multiSend)
	}

	var lastTxHash string
	for _, t := range nonZero {
		c, err := transferCall(t)
		if err != nil {
			return lastTxHash, err
		}
		txHash, err := gnosissafe.ExecCalls(ctx, client, chainID, safe, signers, []gnosissafe.Call{c}, multiSend)
		if err != nil {
			return lastTxHash, fmt.Errorf("safe transfer to %v failed: %w", t.Recipient, err)
		}
		lastTxHash = txHash
	}
	return lastTxHash, nil
}

// sameERC20Token reports whether every transfer moves the same non-native
// ERC-20 token - the only shape ExecuteTransfers batches via MultiSend. A
// single transfer is left to the sequential path (nothing to gain from
// batching one leg), and a native-asset transfer (Token == "") is never
// batched: MultiSendCallOnly could technically carry it too, but Plan
// Section 60 only ever needs to combine the ERC-20 legs, so mixing in a
// native leg here would just be unexercised complexity.
func sameERC20Token(transfers []Transfer) bool {
	if len(transfers) < 2 {
		return false
	}
	token := transfers[0].Token
	if token == "" {
		return false
	}
	for _, t := range transfers[1:] {
		if !strings.EqualFold(t.Token, token) {
			return false
		}
	}
	return true
}

// transferCall turns a Transfer into the Safe call that performs it: a
// plain value transfer for the native asset, or an ERC-20 transfer() call
// on the token contract.
func transferCall(t Transfer) (gnosissafe.Call, error) {
	if t.Token == "" {
		return gnosissafe.Call{To: common.HexToAddress(t.Recipient), Value: t.Amount}, nil
	}
	data, err := erc20ABI.Pack("transfer", common.HexToAddress(t.Recipient), t.Amount)
	if err != nil {
		return gnosissafe.Call{}, fmt.Errorf("packing ERC20 transfer to %v: %w", t.Recipient, err)
	}
	return gnosissafe.Call{To: common.HexToAddress(t.Token), Value: big.NewInt(0), Data: data}, nil
}
