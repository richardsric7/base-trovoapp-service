package network

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
	"trovo-wallet-api/internal/basetxn"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"

	"github.com/ecnepsnai/discord"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var kTempAccountSalt string = "j4rkTZQ2mLk3NAhK"

// GetBlockchainNetworkPassPhrase is vestigial on Base (Stellar used a
// network passphrase for signature domain separation; Base's chain ID,
// baked into every EIP-1559 transaction, serves that role instead). Kept
// so the ~75 existing call sites that pass its result through to
// tx.Sign(passphrase, ...) don't all need to change in this pass.
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
		// Base Sepolia testnet - a safe default so a misconfigured
		// deployment fails against a public testnet, not mainnet.
		return big.NewInt(84532)
	}
	return id
}

// GetBlockchainBaseReserve is vestigial on Base: Stellar reserved a fixed
// XLM balance per ledger subentry (trustlines, offers, ...); EVM accounts
// have no equivalent reserve/subentry concept. Kept, returning the
// configured value (or 0), for the few callers that haven't been
// re-audited yet.
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

// GetBlockchainClient returns the Base JSON-RPC client, the Base
// equivalent of Stellar's Horizon client.
func GetBlockchainClient() *ethclient.Client {
	url := os.Getenv("BASE_RPC_URL")
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
// equivalent of Stellar's sub-entries, trustlines-as-account-state, or
// manage_data key/value store (asset metadata now lives in this app's own
// curated-asset DB rows instead - see internal/components/assets).
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
// it is authorized to hold/send asset (see IsWalletAuthorizedForAsset),
// the native balance, the asset balance, the account object, and whether
// there was an error - same shape as the original Stellar version, so
// existing call sites keep destructuring the same 6 values.
func BlockchainAccountProperties(client *ethclient.Client, destinationAddress string, asset basetxn.Asset) (bool, bool, decimal.Decimal, decimal.Decimal, *AccountInfo, error) {

	log.Printf("[BlockchainAccountProperties] obtaining blockchain account properties for  %v \n", destinationAddress)

	if !common.IsHexAddress(destinationAddress) {
		return false, false, decimal.Zero, decimal.Zero, nil, &tErrors.ErrorInvalidAddress{Address: destinationAddress}
	}
	addr := common.HexToAddress(destinationAddress)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Any well-formed Base address is a valid account - unlike Stellar,
	// there is no separate on-chain "account creation" step.
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

	authorized := IsWalletAuthorizedForAsset(destinationAddress, asset)

	if asset.IsNative() {
		return destinationAccountExists, authorized, nativeAccountBalance, decimal.Zero, account, nil
	}

	assetBalance, err := B20BalanceOf(client, asset.GetIssuer(), destinationAddress)
	if err != nil {
		log.Print("[BlockchainAccountProperties] error fetching B20 balance: ", err)
		return destinationAccountExists, authorized, nativeAccountBalance, decimal.Zero, account, &tErrors.ErrorTemporaryServerError{}
	}

	return destinationAccountExists, authorized, nativeAccountBalance, assetBalance, account, nil
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

// --- Wallet/asset authorization (the Base equivalent of a Stellar
// trustline's authorization flag) ---
//
// A plain B20 (ERC-20-shaped) token contract has no on-chain concept of
// "only these wallets may hold/send me" the way a Stellar
// AUTH_REQUIRED/AUTH_REVOCABLE asset does. Per instructions, that
// restriction is preserved at this application's layer instead: this
// backend will not build/relay a payment or report a positive balance
// authorization for a restricted asset unless the destination wallet has
// an explicit authorization row here - the same gate the original
// enforced by checking a Stellar trustline's authorized flag before
// submitting.

var authDB *gorm.DB

// SetDB wires this package's DB handle for wallet/asset authorization.
// Called once at startup (see main.go).
func SetDB(db *gorm.DB) {
	authDB = db
	if authDB != nil {
		// Rename the legacy "asset_issuer" column (Stellar-era name) to
		// "contract_address" before AutoMigrate, which never renames an
		// existing column on its own - left alone it would add a new
		// empty contract_address column while asset_issuer's data sat
		// unused. No-op if already renamed or the table doesn't exist yet.
		migrator := authDB.Migrator()
		const table = "wallet_asset_authorizations"
		if migrator.HasTable(table) && migrator.HasColumn(table, "asset_issuer") && !migrator.HasColumn(table, "contract_address") {
			if err := migrator.RenameColumn(table, "asset_issuer", "contract_address"); err != nil {
				log.Printf("[SetDB] failed to rename %s.asset_issuer -> contract_address: %v\n", table, err)
			}
		}
		authDB.AutoMigrate(&WalletAssetAuthorization{})
		authDB.AutoMigrate(&AccountSigner{})
	}
	basetxn.SetDefaultBuilder(NewTxBuilder(GetBlockchainClient()))
}

// DB returns the DB handle SetDB wired up - a shared handle for the
// several "does this asset/issuer already exist" style lookups that,
// under Stellar, queried Horizon's global asset registry instead (see
// e.g. internal/components/users/models' GetBlockchainAssets).
func DB() *gorm.DB {
	return authDB
}

// WalletAssetAuthorization is one wallet's authorization to hold/send one
// B20 asset - the Base equivalent of a Stellar trustline's authorization
// flag. Rows are only ever written by a compliance approval (see
// internal/components/users/services' SetWalletAssetAuthorization),
// signed by the asset's own issuing wallet - ApprovedBy/Reason record who
// made that call and why, for audit purposes.
type WalletAssetAuthorization struct {
	ID              uint64 `gorm:"primaryKey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	WalletAddress   string `gorm:"size:64;not null;uniqueIndex:idx_wallet_asset_auth"`
	AssetCode       string `gorm:"size:12;not null;uniqueIndex:idx_wallet_asset_auth"`
	ContractAddress string `gorm:"size:64;not null;uniqueIndex:idx_wallet_asset_auth"`
	Authorized      bool   `gorm:"not null;default:false"`
	ApprovedBy      string `gorm:"size:64;not null;default:''" json:"approvedBy"`
	Reason          string `gorm:"size:255;not null;default:''" json:"reason"`
}

// IsWalletAuthorizedForAsset reports whether wallet may hold/send asset.
// The native asset, and any B20 asset that doesn't require authorization,
// is always authorized - only market-ready tokenized/regulated assets
// (isTokenizedAsset) are gated behind an explicit authorization row.
func IsWalletAuthorizedForAsset(wallet string, asset basetxn.Asset) bool {
	if asset.IsNative() || authDB == nil || !isTokenizedAsset(asset.GetCode()) {
		return true
	}
	var row WalletAssetAuthorization
	err := authDB.Where("wallet_address = ? AND asset_code = ? AND contract_address = ?",
		strings.ToLower(wallet), asset.GetCode(), strings.ToLower(asset.GetIssuer())).First(&row).Error
	if err != nil {
		return false
	}
	return row.Authorized
}

// isTokenizedAsset mirrors sharedconfig.GlobalConfig.IsValidTokenizedAsset's
// query (duplicated rather than imported, to avoid a sharedconfig<->network
// import cycle): an asset only requires wallet-level authorization once
// it's a market-ready tokenized/regulated asset (Asset_Tokenization_Status
// > 3). Every other B20 asset needs no opt-in on Base.
func isTokenizedAsset(assetCode string) bool {
	type Result struct {
		ID string
	}
	var result Result
	authDB.Raw("SELECT id FROM Tokenized_Assets WHERE Asset_Tokenization_Status > 3 AND Asset_Code = upper(?)", assetCode).Scan(&result)
	return len(result.ID) > 0
}

// SetWalletAssetAuthorization grants or revokes wallet's authorization to
// hold/send asset - the Base equivalent of submitting a
// SetTrustLineFlags/ChangeTrust operation on Stellar. approvedBy is the
// address the compliance approval was authenticated against (the asset's
// own issuing wallet - see the service-layer caller), reason an optional
// free-text compliance note; both are recorded for audit purposes.
func SetWalletAssetAuthorization(wallet string, asset basetxn.Asset, authorized bool, approvedBy, reason string) error {
	if authDB == nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	row := WalletAssetAuthorization{
		WalletAddress:   strings.ToLower(wallet),
		AssetCode:       asset.GetCode(),
		ContractAddress: strings.ToLower(asset.GetIssuer()),
		Authorized:      authorized,
		ApprovedBy:      strings.ToLower(approvedBy),
		Reason:          reason,
	}
	// Assign takes a map, not a struct literal: GORM's struct-to-assignment
	// conversion silently drops zero-value fields (false, ""), which would
	// make a revoke (authorized=false) never actually persist.
	return authDB.Where("wallet_address = ? AND asset_code = ? AND contract_address = ?",
		row.WalletAddress, row.AssetCode, row.ContractAddress).
		Assign(map[string]interface{}{"authorized": authorized, "approved_by": row.ApprovedBy, "reason": reason}).
		FirstOrCreate(&row).Error
}

// WalletAssetAuthorizations lists every wallet's authorization row for
// one asset (contractAddress identifies the specific issuance, matching
// IsWalletAuthorizedForAsset's own lookup) - for a compliance review of
// who currently holds/is entitled to hold it.
func WalletAssetAuthorizations(assetCode, contractAddress string) ([]WalletAssetAuthorization, error) {
	rows := make([]WalletAssetAuthorization, 0)
	if authDB == nil {
		return rows, &tErrors.ErrorTemporaryServerError{}
	}
	err := authDB.Where("asset_code = ? AND contract_address = ?", assetCode, strings.ToLower(contractAddress)).
		Order("created_at desc").Find(&rows).Error
	return rows, err
}

// --- Account signer registry (app-layer stand-in for Stellar's
// on-chain weighted multisig; see basetxn.SetOptions' doc) ---

// AccountSigner records that signerAddress is authorized to act for
// account - the Base equivalent of a Stellar SetOptions.Signer entry.
// This is bookkeeping only: unlike Stellar's on-chain multisig, it grants
// no actual on-chain signing power over account (a plain EOA has none to
// grant). Real Base-native multisig needs account to be a smart-contract
// account (e.g. a Safe) with signerAddress as one of its owners - out of
// scope for this alteration pass, per the user's own note that every
// Base wallet needing multisig (subwallets, shared access, recovery)
// will eventually need to be a smart account.
type AccountSigner struct {
	ID            uint64 `gorm:"primaryKey"`
	CreatedAt     time.Time
	Account       string `gorm:"size:64;not null;uniqueIndex:idx_account_signer"`
	SignerAddress string `gorm:"size:64;not null;uniqueIndex:idx_account_signer"`
	Weight        uint32 `gorm:"not null;default:1"`
}

// AddAccountSigner registers signerAddress as a signer for account.
func AddAccountSigner(account, signerAddress string, weight uint32) error {
	if authDB == nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	row := AccountSigner{
		Account:       strings.ToLower(account),
		SignerAddress: strings.ToLower(signerAddress),
		Weight:        weight,
	}
	return authDB.Where("account = ? AND signer_address = ?", row.Account, row.SignerAddress).
		Assign(AccountSigner{Weight: weight}).
		FirstOrCreate(&row).Error
}

// IsAccountSigner reports whether signerAddress is registered as a
// signer for account (account's own address always counts as its own
// signer - a plain EOA is always valid to act for itself).
func IsAccountSigner(account, signerAddress string) bool {
	if strings.EqualFold(account, signerAddress) {
		return true
	}
	if authDB == nil {
		return false
	}
	var count int64
	authDB.Model(&AccountSigner{}).Where("account = ? AND signer_address = ?", strings.ToLower(account), strings.ToLower(signerAddress)).Count(&count)
	return count > 0
}

// --- Transaction submission ---

// SubmitSignedRawTx broadcasts one already-signed raw transaction
// (RLP-encoded, hex or "0x"-prefixed) and waits for it to be mined - the
// Base equivalent of Stellar's client.SubmitTransactionXDR.
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
		return "", classifySubmitError(err)
	}
	return tx.Hash().Hex(), nil
}

// SubmitXdrWithSignature submits a Transaction built via internal/basetxn
// once every operation is signed, attaching signature (a raw signed tx
// hex string) as the final missing piece for signerAddress's operation
// if one is still outstanding - the Base equivalent of Stellar's
// client.SubmitTransactionXDR after AddSignatureBase64. Kept named/shaped
// like the original (xdrBase64 in, tx hash out) for call-site
// compatibility; xdrBase64 is now a comma-joined list of signed raw txs
// (basetxn.Transaction.Base64's shape) rather than an XDR envelope.
func SubmitXdrWithSignature(client *ethclient.Client, signerAddress string, xdrBase64 string, signature string) (string, error) {
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
			return lastHash, err
		}
		lastHash = hash
	}
	return lastHash, nil
}

// SubmittedTransaction is the Base equivalent of Stellar's
// horizon.Transaction, for the "...ReturnsTrx" function variants below.
type SubmittedTransaction struct {
	Hash string
}

func SubmitXdrWithSignatureReturnsTrx(client *ethclient.Client, signerAddress string, xdrBase64 string, signature string) (txnResult SubmittedTransaction, err error) {
	hash, err := SubmitXdrWithSignature(client, signerAddress, xdrBase64, signature)
	return SubmittedTransaction{Hash: hash}, err
}

// SubmitXdrWithSignatures submits a Transaction that multiple signers
// (map key = signer address, value = their signed raw tx hex) each
// contributed one operation's signature for - the Base equivalent of a
// Stellar transaction co-signed by multiple accounts. Kept for call-site
// compatibility; db is unused (kept only so callers' argument lists don't
// need to change) since authorization state is read/written through
// IsWalletAuthorizedForAsset/SetWalletAssetAuthorization instead.
func SubmitXdrWithSignatures(client *ethclient.Client, transactionXdr string, signatures map[string]string, db *gorm.DB) (string, error) {
	rawTxs := strings.Split(transactionXdr, ",")
	for _, sig := range signatures {
		rawTxs = append(rawTxs, sig)
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

// pendingAuthForSubmit mirrors the DB row internal/network reads to
// resolve a pending multi-approver approval into raw signed
// transactions - see SubmitApprovalsXdrWithSignatures's doc for why each
// approver's "signature" is now itself a signed raw tx.
type pendingAuthForSubmit struct {
	ID                           string `gorm:"size:56"`
	TransactionXdr               string
	PendingTransactionSignatures []struct {
		TransactionWithSignature string
	} `gorm:"-"`
}

// SubmitApprovalsXdrWithSignatures submits every approver's contributed
// raw signed transaction for a pending multi-party approval. Stellar's
// version let N approvers co-sign ONE multisig-account transaction; Base
// has no plain-account equivalent (that needs a deployed multisig/Safe
// contract, out of scope for this alteration pass), so each approver's
// "signature" is itself a fully signed raw transaction for their own
// operation, and this submits each in turn - application-sequenced, not
// chain-atomic multisig.
func SubmitApprovalsXdrWithSignatures(client *ethclient.Client, approvalID string, db *gorm.DB) (string, error) {
	type pendingTransactionSignature struct {
		PendingAuthID            string `gorm:"size:56"`
		TransactionWithSignature string
	}
	type pendingAuth struct {
		ID             string `gorm:"size:56"`
		TransactionXdr string
	}
	var auth pendingAuth
	if e := db.Table("pending_auths").Where("id = ?", approvalID).First(&auth).Error; e != nil {
		return "", &tErrors.ErrorTemporaryServerError{}
	}
	var sigs []pendingTransactionSignature
	db.Table("pending_transaction_signatures").Where("pending_auth_id = ?", approvalID).Find(&sigs)

	rawTxs := strings.Split(auth.TransactionXdr, ",")
	for _, s := range sigs {
		rawTxs = append(rawTxs, s.TransactionWithSignature)
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

func SubmitApprovalsXdrWithSignaturesReturnsTrx(client *ethclient.Client, approvalID string, db *gorm.DB) (txnResult SubmittedTransaction, err error) {
	hash, err := SubmitApprovalsXdrWithSignatures(client, approvalID, db)
	return SubmittedTransaction{Hash: hash}, err
}

// SubmitXdrWithSignatureChannelAccounts submits a transaction co-signed
// by a channel (gas-sponsoring) account and the sender - both already
// signed raw txs (see SubmitXdrWithSignatures's doc for why).
func SubmitXdrWithSignatureChannelAccounts(client *ethclient.Client, signerAddress, channelAccountPK string, xdrBase64 string, txSignature, chanSignature string) (string, error) {
	rawTxs := []string{}
	for _, s := range strings.Split(xdrBase64, ",") {
		if strings.TrimSpace(s) != "" {
			rawTxs = append(rawTxs, s)
		}
	}
	if chanSignature != "" {
		rawTxs = append(rawTxs, chanSignature)
	}
	if txSignature != "" {
		rawTxs = append(rawTxs, txSignature)
	}
	var lastHash string
	for _, raw := range rawTxs {
		hash, err := SubmitSignedRawTx(client, raw)
		if err != nil {
			return "", err
		}
		lastHash = hash
	}
	return lastHash, nil
}

func classifySubmitError(err error) error {
	errorString := err.Error()
	if strings.Contains(errorString, "insufficient funds") {
		return &tErrors.CustomError{
			Param:      "destinationAssetCode",
			Err:        "error-low-liquidity",
			ErrMessage: "There is not enough balance to complete this transaction. Please try again later or reduce the amount and try again.",
		}
	}
	if strings.Contains(errorString, "nonce too low") || strings.Contains(errorString, "already known") || strings.Contains(errorString, "replacement transaction underpriced") {
		return &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-please-reject-transaction",
			ErrMessage: "This transaction has become invalid and cannot be completed. Please reject this transaction and initiate a fresh one.",
		}
	}
	return &tErrors.CustomError{Param: "publicKey", Err: "error operation failed", ErrMessage: "Operation Failed", Code: 500}
}

func notifyNetworkFailure(prefix string, err error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	errStr := err.Error()
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "handshake") || strings.Contains(errStr, "read tcp") || strings.Contains(errStr, "connection reset by peer") || strings.Contains(errStr, "dial tcp") || strings.Contains(errStr, "no such host") {
		discord.Say(fmt.Sprintf("%s error connecting to Base RPC service: %v", prefix, err))
	}
}

func logDiscordFailedPayment(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}
