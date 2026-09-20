package network

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
	"trovo-wallet-payment-history-engine/internal/basetxn"
	tErrors "trovo-wallet-payment-history-engine/internal/errors"
	"trovo-wallet-payment-history-engine/internal/evmkeypair"

	"github.com/ecnepsnai/discord"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

var kTempAccountSalt string = "j4rkTZQ2mLk3NAhK"

// GetBlockchainNetworkPassPhrase is vestigial on Base (Stellar used a
// network passphrase for signature domain separation; Base's chain ID,
// baked into every EIP-1559 transaction, serves that role instead). Kept so
// existing call sites that pass its result through to tx.Sign(passphrase,
// ...) don't need to change.
func GetBlockchainNetworkPassPhrase() string {
	return os.Getenv("BLOCKCHAIN_NETWORK_PASSPHRASE")
}

// GetBlockchainChainID returns the Base chain ID transactions are signed
// and submitted against (BASE_CHAIN_ID, e.g. 8453 for Base mainnet, 84532
// for Base Sepolia).
func GetBlockchainChainID() *big.Int {
	raw := strings.TrimSpace(os.Getenv("BASE_CHAIN_ID"))
	id, ok := new(big.Int).SetString(raw, 10)
	if !ok {
		return big.NewInt(84532)
	}
	return id
}

// GetBlockchainBaseReserve is vestigial on Base: Stellar reserved a fixed
// XLM balance per ledger subentry; EVM accounts have no equivalent reserve
// concept. Kept, returning the configured value (or 0), for call-site
// compatibility.
func GetBlockchainBaseReserve() decimal.Decimal {
	val, err := decimal.NewFromString(os.Getenv("BLOCKCHAIN_BASE_RESERVE"))
	if err != nil {
		return decimal.Zero
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

// GetBlockchainClient returns the Base JSON-RPC client, the Base equivalent
// of Stellar's Horizon client. Reads BASE_RPC_URL, falling back to the
// original EXPANSION_URL env var if BASE_RPC_URL isn't set yet.
func GetBlockchainClient() *ethclient.Client {
	url := os.Getenv("BASE_RPC_URL")
	if url == "" {
		url = os.Getenv("EXPANSION_URL")
	}
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Panicf("[GetBlockchainClient] invalid BASE_RPC_URL %q: %v", url, err)
	}
	return client
}

func TempAccountKeypair(publicKey string) (*evmkeypair.Full, error) {
	mnemonic := os.Getenv("MNEMONIC_TEMP_ACCOUNTS")

	h := crypto.Keccak256(
		[]byte(kTempAccountSalt),
		[]byte(mnemonic),
		[]byte(publicKey),
	)

	var rawSeed [32]byte
	copy(rawSeed[:], h[0:32])

	return evmkeypair.FromRawSeed(rawSeed)
}

// AccountInfo is the Base equivalent of Stellar's *horizon.Account -
// deliberately minimal, since Base/EVM accounts have no ledger-native
// equivalent of Stellar's sub-entries or trustlines-as-account-state.
type AccountInfo struct {
	Address string
	Nonce   uint64
}

var erc20ABI abi.ABI

func init() {
	var err error
	erc20ABI, err = abi.JSON(strings.NewReader(`[
		{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},
		{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"}
	]`))
	if err != nil {
		log.Panicf("[network] invalid embedded B20/ERC20 ABI: %v", err)
	}
}

// BlockchainAccountProperties returns whether the account exists, whether
// it is authorized to hold/send asset, the native balance, the asset
// balance, the account object, and whether there was an error - same shape
// as the original Stellar version. This module has no app-layer
// authorization registry (that restriction is enforced by app-backend,
// which owns the payment flow this engine only records history for), so a
// plain B20/ERC20 token here has no on-chain "trust" step either - any
// well-formed address can hold/receive one, so authorized is always true.
func BlockchainAccountProperties(client *ethclient.Client, destinationAddress string, asset basetxn.Asset) (bool, bool, decimal.Decimal, decimal.Decimal, *AccountInfo, error) {
	log.Printf("[BlockchainAccountProperties] obtaining blockchain account properties for  %v \n", destinationAddress)

	if !common.IsHexAddress(destinationAddress) {
		return false, false, decimal.Zero, decimal.Zero, nil, &tErrors.ErrorInvalidAddress{Address: destinationAddress}
	}
	addr := common.HexToAddress(destinationAddress)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	destinationAccountExists := true

	nonce, err := client.PendingNonceAt(ctx, addr)
	if err != nil {
		log.Print("[BlockchainAccountProperties] error fetching nonce: ", err)
		return destinationAccountExists, false, decimal.Zero, decimal.Zero, &AccountInfo{Address: destinationAddress}, &tErrors.ErrorTemporaryServerError{}
	}
	account := &AccountInfo{Address: destinationAddress, Nonce: nonce}

	weiBalance, err := client.BalanceAt(ctx, addr, nil)
	if err != nil {
		log.Print("[BlockchainAccountProperties] error fetching native balance: ", err)
		return destinationAccountExists, false, decimal.Zero, decimal.Zero, account, &tErrors.ErrorTemporaryServerError{}
	}
	nativeAccountBalance := weiToDecimal(weiBalance)

	if asset.IsNative() {
		return destinationAccountExists, true, nativeAccountBalance, decimal.Zero, account, nil
	}

	assetBalance, err := B20BalanceOf(client, asset.GetIssuer(), destinationAddress)
	if err != nil {
		log.Print("[BlockchainAccountProperties] error fetching B20 balance: ", err)
		return destinationAccountExists, true, nativeAccountBalance, decimal.Zero, account, &tErrors.ErrorTemporaryServerError{}
	}

	return destinationAccountExists, true, nativeAccountBalance, assetBalance, account, nil
}

// B20BalanceOf calls the B20 (ERC-20-shaped) token contract's balanceOf.
func B20BalanceOf(client *ethclient.Client, tokenContract, holder string) (decimal.Decimal, error) {
	if !common.IsHexAddress(tokenContract) || !common.IsHexAddress(holder) {
		return decimal.Zero, fmt.Errorf("invalid token contract or holder address")
	}
	data, err := erc20ABI.Pack("balanceOf", common.HexToAddress(holder))
	if err != nil {
		return decimal.Zero, err
	}
	to := common.HexToAddress(tokenContract)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := client.CallContract(ctx, ethereum.CallMsg{To: &to, Data: data}, nil)
	if err != nil {
		return decimal.Zero, err
	}
	outputs, err := erc20ABI.Unpack("balanceOf", result)
	if err != nil || len(outputs) == 0 {
		return decimal.Zero, fmt.Errorf("could not decode balanceOf result")
	}
	raw, ok := outputs[0].(*big.Int)
	if !ok {
		return decimal.Zero, fmt.Errorf("unexpected balanceOf result type")
	}
	return weiToDecimal(raw), nil
}

func weiToDecimal(wei *big.Int) decimal.Decimal {
	return decimal.NewFromBigInt(wei, -18)
}

func decimalToWei(amount decimal.Decimal) *big.Int {
	return amount.Shift(18).BigInt()
}

// TxBuilder implements basetxn.Builder, turning a basetxn.Payment into an
// unsigned Base EIP-1559 transaction: a native transfer, or an ERC-20
// `transfer` call for a B20 asset.
type TxBuilder struct {
	Client *ethclient.Client
}

func NewTxBuilder(client *ethclient.Client) *TxBuilder {
	return &TxBuilder{Client: client}
}

func (b *TxBuilder) BuildPaymentTx(ctx context.Context, from common.Address, op basetxn.Payment) (*types.Transaction, error) {
	to := common.HexToAddress(op.Destination)
	nonce, err := b.Client.PendingNonceAt(ctx, from)
	if err != nil {
		return nil, err
	}
	header, err := b.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	gasTipCap, err := b.Client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}
	gasFeeCap := new(big.Int).Add(gasTipCap, new(big.Int).Mul(header.BaseFee, big.NewInt(2)))
	chainID := GetBlockchainChainID()
	amount, err := decimal.NewFromString(op.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid payment amount %q: %w", op.Amount, err)
	}

	if op.Asset.IsNative() {
		return types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        &to,
			Value:     decimalToWei(amount),
			Gas:       21000,
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
		}), nil
	}

	data, err := erc20ABI.Pack("transfer", to, decimalToWei(amount))
	if err != nil {
		return nil, err
	}
	tokenAddr := common.HexToAddress(op.Asset.GetIssuer())
	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &tokenAddr,
		Value:     big.NewInt(0),
		Gas:       120000,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Data:      data,
	}), nil
}

// SubmitSignedRawTx broadcasts one already-signed raw transaction
// (RLP-encoded, hex or "0x"-prefixed) and waits for it to be accepted -
// the Base equivalent of Stellar's client.SubmitTransactionXDR.
func SubmitSignedRawTx(client *ethclient.Client, rawTxHex string) (txHash string, err error) {
	rawTxHex = strings.TrimPrefix(strings.TrimSpace(rawTxHex), "0x")
	raw := common.FromHex("0x" + rawTxHex)
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(raw); err != nil {
		return "", &tErrors.ErrorInvalidTransaction{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := client.SendTransaction(ctx, tx); err != nil {
		notifyNetworkFailure("[SubmitSignedRawTx]", err)
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// SubmitXdrWithSignature submits a Transaction built via internal/basetxn
// once its Payment operation is signed - the Base equivalent of Stellar's
// client.SubmitTransactionXDR. Kept named/shaped like the original
// (xdrBase64 in, tx hash out) for call-site compatibility; xdrBase64 is now
// a comma-joined list of signed raw txs (basetxn.Transaction.Base64's
// shape) rather than an XDR envelope.
func SubmitXdrWithSignature(client *ethclient.Client, ownerAddress string, xdrBase64 string, signature string) (string, error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if v := os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK"); len(v) > 50 {
		discord.WebhookURL = v
	}
	rawTxs := strings.Split(xdrBase64, ",")
	if signature != "" {
		rawTxs = append(rawTxs, signature)
	}
	var lastHash string
	for _, raw := range rawTxs {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		hash, err := SubmitSignedRawTx(client, raw)
		if err != nil {
			logDiscordFailedPayment(fmt.Sprintf("[SubmitXdrWithSignature] error submitting tx: %v\nOwner publickey: %v\n", err, ownerAddress))
			return lastHash, &tErrors.CustomError{Param: "destination", Err: "error payment failed", ErrMessage: "Payment Failed", Code: 500}
		}
		lastHash = hash
	}
	return lastHash, nil
}

// SubmitXdrWithSignatureChannelAccounts submits a Transaction co-signed by
// a channel account and the sending signer - the Base equivalent of a
// Stellar transaction with a fee-sponsoring channel-account source.
func SubmitXdrWithSignatureChannelAccounts(client *ethclient.Client, signerAddress, channelAccountPK string, xdrBase64 string, txSignature, chanSignature string) (string, error) {
	rawTxs := strings.Split(xdrBase64, ",")
	for _, sig := range []string{chanSignature, txSignature} {
		sig = strings.TrimSpace(sig)
		if sig != "" {
			rawTxs = append(rawTxs, sig)
		}
	}
	var lastHash string
	for _, raw := range rawTxs {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		hash, err := SubmitSignedRawTx(client, raw)
		if err != nil {
			return lastHash, err
		}
		lastHash = hash
	}
	return lastHash, nil
}

func notifyNetworkFailure(prefix string, err error) {
	msg := err.Error()
	if strings.Contains(msg, "timeout") || strings.Contains(msg, "handshake") || strings.Contains(msg, "read tcp") || strings.Contains(msg, "connection reset by peer") || strings.Contains(msg, "dial tcp") || strings.Contains(msg, "no such host") {
		discord.Say(fmt.Sprintf("%s error connecting to blockchain RPC: %v", prefix, err))
	}
}

func logDiscordFailedPayment(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if v := os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK"); len(v) > 50 {
		discord.WebhookURL = v
	}
	discord.Say(msg)
}
