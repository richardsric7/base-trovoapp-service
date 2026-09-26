package safesigner

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"os"
	"sort"
	"time"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/network"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Gnosis Safe v1.3.0/v1.4.x standard operation types.
const (
	operationCall = uint8(0)
)

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

var (
	erc20ABI    abi.ABI
	safeABI     abi.ABI
	structTypes abi.Arguments
)

func init() {
	var err error
	erc20ABI, err = abi.JSON(stringsReader(`[
		{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"}
	]`))
	if err != nil {
		panic("safesigner: invalid erc20 ABI: " + err.Error())
	}
	safeABI, err = abi.JSON(stringsReader(`[
		{"constant":true,"inputs":[],"name":"nonce","outputs":[{"name":"","type":"uint256"}],"type":"function"},
		{"constant":false,"inputs":[
			{"name":"to","type":"address"},
			{"name":"value","type":"uint256"},
			{"name":"data","type":"bytes"},
			{"name":"operation","type":"uint8"},
			{"name":"safeTxGas","type":"uint256"},
			{"name":"baseGas","type":"uint256"},
			{"name":"gasPrice","type":"uint256"},
			{"name":"gasToken","type":"address"},
			{"name":"refundReceiver","type":"address"},
			{"name":"signatures","type":"bytes"}
		],"name":"execTransaction","outputs":[{"name":"","type":"bool"}],"type":"function"}
	]`))
	if err != nil {
		panic("safesigner: invalid safe ABI: " + err.Error())
	}

	addressTy, _ := abi.NewType("address", "", nil)
	uint256Ty, _ := abi.NewType("uint256", "", nil)
	bytes32Ty, _ := abi.NewType("bytes32", "", nil)
	uint8Ty, _ := abi.NewType("uint8", "", nil)
	structTypes = abi.Arguments{
		{Type: bytes32Ty}, // SAFE_TX_TYPEHASH
		{Type: addressTy}, // to
		{Type: uint256Ty}, // value
		{Type: bytes32Ty}, // keccak256(data)
		{Type: uint8Ty},   // operation
		{Type: uint256Ty}, // safeTxGas
		{Type: uint256Ty}, // baseGas
		{Type: uint256Ty}, // gasPrice
		{Type: addressTy}, // gasToken
		{Type: addressTy}, // refundReceiver
		{Type: uint256Ty}, // nonce
	}
}

// SAFE_TX_TYPEHASH = keccak256("SafeTx(address to,uint256 value,bytes data,uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice,address gasToken,address refundReceiver,uint256 nonce)")
var safeTxTypehash = crypto.Keccak256([]byte("SafeTx(address to,uint256 value,bytes data,uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice,address gasToken,address refundReceiver,uint256 nonce)"))

// DOMAIN_SEPARATOR_TYPEHASH = keccak256("EIP712Domain(uint256 chainId,address verifyingContract)")
var domainSeparatorTypehash = crypto.Keccak256([]byte("EIP712Domain(uint256 chainId,address verifyingContract)"))

// ExecuteTransfers assembles, signs (with the first 3 P2P_ESCROW_SIGNERS),
// and submits one Gnosis Safe execTransaction per non-zero transfer. Each
// transfer is its own sequential Safe transaction (own nonce, own
// signatures) rather than bundled via Safe's MultiSend contract - simpler
// and safer to get right without a live network to test against, at the
// cost of the four legs (buyer net amount + 3 fee wallets, Plan Section 60)
// not being atomic with each other. Returns the transaction hash of the
// last (largest/buyer-facing) transfer submitted, which Order.
// AssetReleaseTransactionHash stores as the canonical release hash - the
// others remain visible via each transfer's own on-chain record.
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

	var lastTxHash string
	for _, t := range transfers {
		if t.Amount == nil || t.Amount.Sign() <= 0 {
			continue // skip zero-amount legs (e.g. a fee that rounded to 0)
		}
		txHash, err := executeSingleTransfer(client, chainID, safeAddress, signers, t)
		if err != nil {
			return lastTxHash, fmt.Errorf("safe transfer to %v failed: %w", t.Recipient, err)
		}
		lastTxHash = txHash
	}
	if lastTxHash == "" {
		return "", fmt.Errorf("no non-zero transfers to execute")
	}
	return lastTxHash, nil
}

func executeSingleTransfer(client *ethclient.Client, chainID *big.Int, safeAddress string, signers []*evmkeypair.Full, t Transfer) (string, error) {
	to := common.HexToAddress(t.Recipient)
	value := big.NewInt(0)
	var data []byte
	if t.Token == "" {
		// native asset transfer: Safe sends value directly to `to`, no data
		value = t.Amount
	} else {
		packed, err := erc20ABI.Pack("transfer", to, t.Amount)
		if err != nil {
			return "", fmt.Errorf("packing ERC20 transfer: %w", err)
		}
		data = packed
		to = common.HexToAddress(t.Token)
	}

	nonce, err := fetchSafeNonce(client, safeAddress)
	if err != nil {
		return "", fmt.Errorf("fetching Safe nonce: %w", err)
	}

	safeTxHash, err := computeSafeTxHash(chainID, safeAddress, to, value, data, operationCall, big.NewInt(0), big.NewInt(0), big.NewInt(0), common.Address{}, common.Address{}, nonce)
	if err != nil {
		return "", fmt.Errorf("computing safeTxHash: %w", err)
	}

	signatures, err := signSafeTxHash(safeTxHash, signers)
	if err != nil {
		return "", fmt.Errorf("signing safeTxHash: %w", err)
	}

	execData, err := safeABI.Pack("execTransaction", to, value, data, operationCall,
		big.NewInt(0), big.NewInt(0), big.NewInt(0), common.Address{}, common.Address{}, signatures)
	if err != nil {
		return "", fmt.Errorf("packing execTransaction: %w", err)
	}

	// The first configured signer broadcasts and pays gas for the outer
	// transaction; its own signature is also one of the 3 Safe signatures.
	broadcaster := signers[0]
	return submitRawTransaction(client, chainID, broadcaster, safeAddress, execData)
}

func fetchSafeNonce(client *ethclient.Client, safeAddress string) (*big.Int, error) {
	data, err := safeABI.Pack("nonce")
	if err != nil {
		return nil, err
	}
	to := common.HexToAddress(safeAddress)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := client.CallContract(ctx, callMsg(to, data), nil)
	if err != nil {
		return nil, err
	}
	out, err := safeABI.Unpack("nonce", result)
	if err != nil || len(out) == 0 {
		return nil, fmt.Errorf("unexpected nonce() response")
	}
	nonce, ok := out[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected nonce() type")
	}
	return nonce, nil
}

// computeSafeTxHash implements Safe's EIP-712 SafeTx typed-data hash exactly
// as the Safe contract itself verifies it (getTransactionHash), so
// signatures produced here validate on-chain without modification.
func computeSafeTxHash(chainID *big.Int, safeAddress string, to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas, baseGas, gasPrice *big.Int, gasToken, refundReceiver common.Address, nonce *big.Int) ([]byte, error) {
	addressTy, _ := abi.NewType("address", "", nil)
	uint256Ty, _ := abi.NewType("uint256", "", nil)
	domainArgs := abi.Arguments{
		{Type: mustBytes32Type()}, {Type: uint256Ty}, {Type: addressTy},
	}
	domainPacked, err := domainArgs.Pack(toBytes32(domainSeparatorTypehash), chainID, common.HexToAddress(safeAddress))
	if err != nil {
		return nil, err
	}
	domainSeparator := crypto.Keccak256(domainPacked)

	dataHash := crypto.Keccak256(data)
	structPacked, err := structTypes.Pack(
		toBytes32(safeTxTypehash),
		to,
		value,
		toBytes32(dataHash),
		operation,
		safeTxGas,
		baseGas,
		gasPrice,
		gasToken,
		refundReceiver,
		nonce,
	)
	if err != nil {
		return nil, err
	}
	structHash := crypto.Keccak256(structPacked)

	// EIP-712: keccak256("\x19\x01" || domainSeparator || structHash)
	preimage := append([]byte{0x19, 0x01}, domainSeparator...)
	preimage = append(preimage, structHash...)
	return crypto.Keccak256(preimage), nil
}

// signSafeTxHash produces Safe's "type 0" (regular EOA-over-raw-hash)
// signature format for each signer, concatenated in ascending signer-
// address order as Safe's checkSignatures requires.
func signSafeTxHash(safeTxHash []byte, signers []*evmkeypair.Full) ([]byte, error) {
	type sig struct {
		addr common.Address
		raw  []byte
	}
	sigs := make([]sig, 0, len(signers))
	for _, s := range signers {
		privKey := s.PrivateKey()
		signature, err := crypto.Sign(safeTxHash, privKey)
		if err != nil {
			return nil, err
		}
		// crypto.Sign returns [R || S || V] with V in {0,1}; Safe expects
		// V in {27,28} for a direct (non eth_sign-prefixed) hash signature.
		signature[64] += 27
		sigs = append(sigs, sig{addr: common.HexToAddress(s.Address()), raw: signature})
	}
	sort.Slice(sigs, func(i, j int) bool {
		return bytes.Compare(sigs[i].addr.Bytes(), sigs[j].addr.Bytes()) < 0
	})
	out := make([]byte, 0, 65*len(sigs))
	for _, s := range sigs {
		out = append(out, s.raw...)
	}
	return out, nil
}

func submitRawTransaction(client *ethclient.Client, chainID *big.Int, broadcaster *evmkeypair.Full, to string, data []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fromAddr := common.HexToAddress(broadcaster.Address())
	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return "", err
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}
	toAddr := common.HexToAddress(to)
	callMessage := callMsg(toAddr, data)
	callMessage.From = fromAddr
	gasLimit, err := client.EstimateGas(ctx, callMessage)
	if err != nil {
		// fall back to a generous fixed limit if estimation fails (e.g. the
		// node doesn't support eth_estimateGas for pending state)
		gasLimit = 500000
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddr,
		Value:    big.NewInt(0),
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), broadcaster.PrivateKey())
	if err != nil {
		return "", err
	}
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}
