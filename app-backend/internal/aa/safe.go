// Package aa builds ERC-4337 (EntryPoint v0.7) UserOperations for Trovo
// wallets.
//
// Every Trovo wallet is a Safe v1.4.1 proxy with the Safe4337Module
// enabled (also its fallback handler), owned by the user's mnemonic key
// (or, for shared wallets, by its approvers). A wallet's address is fixed
// by CREATE2 before it is deployed, so it can receive funds from
// registration on; its first UserOperation deploys it (initCode) and acts
// in one step, paying its own gas - in ETH, or in a stablecoin through
// TrovoTokenPaymaster.
//
// The owners sign the Safe4337Module's EIP-712 SafeOp hash with the apps'
// existing personal_sign; Safe accepts those as eth_sign signatures (v+4).
// The address formula is shared with wallet-core (src/safe.rs), and both
// are tested against real Safe deployments (testdata/safe_fixture.json).
package aa

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"

	"trovo-wallet-api/internal/gnosissafe"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Canonical deployments (identical on Base and Base Sepolia).
const (
	DefaultEntryPoint        = "0x0000000071727De22E5E9d8BAf0edAc6f37da032" // ERC-4337 v0.7
	DefaultSafeProxyFactory  = "0x4e1DCf7AD4e460CfD30791CCC4F9c8a4f820ec67" // v1.4.1
	DefaultSafeSingleton     = "0x29fcB43b46531BcA003ddC8FCB67FFE91900C762" // SafeL2 v1.4.1
	DefaultSafeModuleSetup   = "0x2dd68b007B46fBe91B9A7c3EDa5A7a1063cB5b47" // safe-modules v0.3.0
	DefaultSafe4337Module    = "0x75cf11467937ce3F2f357CE24ffc3DBF8fD5c226" // safe-modules v0.3.0 (EntryPoint v0.7)
	DefaultMultiSendCallOnly = "0x9641d764fc13c8B624c04430C7356C1C7C8102e2" // v1.4.1

	// SafeVersion is recorded on every wallet so a future layout change
	// can tell existing wallets apart.
	SafeVersion = "safe-1.4.1+4337-0.3.0"
)

// SafeProxyCreationCode is `type(SafeProxy).creationCode`, as returned by
// the canonical v1.4.1 SafeProxyFactory's proxyCreationCode().
const SafeProxyCreationCode = "608060405234801561001057600080fd5b506040516101e63803806101e68339818101604052602081101561003357600080fd5b8101908080519060200190929190505050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156100ca576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001806101c46022913960400191505060405180910390fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505060ab806101196000396000f3fe608060405273ffffffffffffffffffffffffffffffffffffffff600054167fa619486e0000000000000000000000000000000000000000000000000000000060003514156050578060005260206000f35b3660008037600080366000845af43d6000803e60008114156070573d6000fd5b3d6000f3fea264697066735822122003d1488ee65e08fa41e58e888a9865554c535f2c77126a82cb4c0f917f31441364736f6c63430007060033496e76616c69642073696e676c65746f6e20616464726573732070726f7669646564"

// Config is the set of contracts wallets are built from.
type Config struct {
	ChainID           *big.Int
	EntryPoint        common.Address
	ProxyFactory      common.Address
	Singleton         common.Address
	ModuleSetup       common.Address
	Safe4337Module    common.Address
	MultiSendCallOnly common.Address
	ProxyCreationCode []byte
}

func envAddress(key, def string) common.Address {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return common.HexToAddress(v)
	}
	return common.HexToAddress(def)
}

// ConfigFromEnv returns the canonical deployments, overridable per
// contract through SAFE_PROXY_FACTORY_ADDRESS, SAFE_SINGLETON_ADDRESS,
// SAFE_MODULE_SETUP_ADDRESS, SAFE_4337_MODULE_ADDRESS,
// SAFE_MULTISEND_CALL_ONLY_ADDRESS and ENTRYPOINT_ADDRESS (for local
// chains; production uses the defaults).
func ConfigFromEnv(chainID *big.Int) Config {
	code, _ := hex.DecodeString(SafeProxyCreationCode)
	return Config{
		ChainID:           chainID,
		EntryPoint:        envAddress("ENTRYPOINT_ADDRESS", DefaultEntryPoint),
		ProxyFactory:      envAddress("SAFE_PROXY_FACTORY_ADDRESS", DefaultSafeProxyFactory),
		Singleton:         envAddress("SAFE_SINGLETON_ADDRESS", DefaultSafeSingleton),
		ModuleSetup:       envAddress("SAFE_MODULE_SETUP_ADDRESS", DefaultSafeModuleSetup),
		Safe4337Module:    envAddress("SAFE_4337_MODULE_ADDRESS", DefaultSafe4337Module),
		MultiSendCallOnly: envAddress("SAFE_MULTISEND_CALL_ONLY_ADDRESS", DefaultMultiSendCallOnly),
		ProxyCreationCode: code,
	}
}

var contractsABI = mustABI(`[
 {"name":"setup","type":"function","inputs":[
   {"name":"_owners","type":"address[]"},{"name":"_threshold","type":"uint256"},{"name":"to","type":"address"},
   {"name":"data","type":"bytes"},{"name":"fallbackHandler","type":"address"},{"name":"paymentToken","type":"address"},
   {"name":"payment","type":"uint256"},{"name":"paymentReceiver","type":"address"}],"outputs":[]},
 {"name":"enableModules","type":"function","inputs":[{"name":"modules","type":"address[]"}],"outputs":[]},
 {"name":"createProxyWithNonce","type":"function","inputs":[
   {"name":"_singleton","type":"address"},{"name":"initializer","type":"bytes"},{"name":"saltNonce","type":"uint256"}],
   "outputs":[{"name":"proxy","type":"address"}]},
 {"name":"executeUserOp","type":"function","inputs":[
   {"name":"to","type":"address"},{"name":"value","type":"uint256"},{"name":"data","type":"bytes"},{"name":"operation","type":"uint8"}],"outputs":[]},
 {"name":"multiSend","type":"function","inputs":[{"name":"transactions","type":"bytes"}],"outputs":[]},
 {"name":"addOwnerWithThreshold","type":"function","inputs":[{"name":"owner","type":"address"},{"name":"_threshold","type":"uint256"}],"outputs":[]},
 {"name":"removeOwner","type":"function","inputs":[{"name":"prevOwner","type":"address"},{"name":"owner","type":"address"},{"name":"_threshold","type":"uint256"}],"outputs":[]},
 {"name":"swapOwner","type":"function","inputs":[{"name":"prevOwner","type":"address"},{"name":"oldOwner","type":"address"},{"name":"newOwner","type":"address"}],"outputs":[]},
 {"name":"changeThreshold","type":"function","inputs":[{"name":"_threshold","type":"uint256"}],"outputs":[]},
 {"name":"getOwners","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"address[]"}]},
 {"name":"getThreshold","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"uint256"}]},
 {"name":"getNonce","type":"function","stateMutability":"view","inputs":[{"name":"sender","type":"address"},{"name":"key","type":"uint192"}],"outputs":[{"type":"uint256"}]},
 {"name":"transfer","type":"function","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"type":"bool"}]},
 {"name":"burn","type":"function","inputs":[{"name":"value","type":"uint256"}],"outputs":[]},
 {"name":"approve","type":"function","inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"type":"bool"}]},
 {"name":"allowance","type":"function","stateMutability":"view","inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"outputs":[{"type":"uint256"}]},
 {"name":"balanceOf","type":"function","stateMutability":"view","inputs":[{"name":"account","type":"address"}],"outputs":[{"type":"uint256"}]}
]`)

func mustABI(s string) abi.ABI {
	a, err := abi.JSON(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	return a
}

func pack(method string, args ...interface{}) []byte {
	data, err := contractsABI.Pack(method, args...)
	if err != nil {
		panic(fmt.Sprintf("aa: packing %s: %v", method, err))
	}
	return data
}

// SortedOwners returns owners in ascending address order (the order a
// Safe's signatures must be in; the owner list order itself is kept as
// given in the initializer).
func SortedOwners(owners []common.Address) []common.Address {
	out := append([]common.Address(nil), owners...)
	sort.Slice(out, func(i, j int) bool { return bytes.Compare(out[i].Bytes(), out[j].Bytes()) < 0 })
	return out
}

// Initializer is the Safe.setup calldata for owners/threshold, enabling the
// Safe4337Module and using it as the fallback handler.
func (c Config) Initializer(owners []common.Address, threshold int64) []byte {
	enable := pack("enableModules", []common.Address{c.Safe4337Module})
	return pack("setup", owners, big.NewInt(threshold), c.ModuleSetup, enable, c.Safe4337Module,
		common.Address{}, big.NewInt(0), common.Address{})
}

// SafeAddress is the address createProxyWithNonce deploys the Safe to.
func (c Config) SafeAddress(owners []common.Address, threshold int64, saltNonce *big.Int) common.Address {
	init := c.Initializer(owners, threshold)
	salt := crypto.Keccak256(crypto.Keccak256(init), common.LeftPadBytes(saltNonce.Bytes(), 32))
	initCode := append(append([]byte{}, c.ProxyCreationCode...), common.LeftPadBytes(c.Singleton.Bytes(), 32)...)
	return crypto.CreateAddress2(c.ProxyFactory, common.BytesToHash(salt), crypto.Keccak256(initCode))
}

// InitCode is the UserOperation initCode (factory ++ factoryData) that
// deploys the Safe.
func (c Config) InitCode(owners []common.Address, threshold int64, saltNonce *big.Int) []byte {
	data := pack("createProxyWithNonce", c.Singleton, c.Initializer(owners, threshold), saltNonce)
	return append(append([]byte{}, c.ProxyFactory.Bytes()...), data...)
}

// DeploySafeCall is a call to the proxy factory deploying the Safe owned by
// owners/threshold with saltNonce - how a primary wallet deploys its
// sub-wallets inside its own operation.
func (c Config) DeploySafeCall(owners []common.Address, threshold int64, saltNonce *big.Int) Call {
	return Call{To: c.ProxyFactory, Value: big.NewInt(0), Data: pack("createProxyWithNonce", c.Singleton, c.Initializer(owners, threshold), saltNonce)}
}

// Call is one call the wallet makes.
type Call = gnosissafe.Call

// CallData is the UserOperation callData executing calls from the Safe:
// executeUserOp(call) for one, or a delegatecall to MultiSendCallOnly
// batching them atomically.
func (c Config) CallData(calls []Call) ([]byte, error) {
	switch len(calls) {
	case 0:
		return nil, fmt.Errorf("aa: no calls")
	case 1:
		v := calls[0].Value
		if v == nil {
			v = big.NewInt(0)
		}
		return pack("executeUserOp", calls[0].To, v, calls[0].Data, gnosissafe.OperationCall), nil
	}
	var packed []byte
	for _, call := range calls {
		v := call.Value
		if v == nil {
			v = big.NewInt(0)
		}
		packed = append(packed, gnosissafe.PackMultiSendTx(gnosissafe.OperationCall, call.To, v, call.Data)...)
	}
	return pack("executeUserOp", c.MultiSendCallOnly, big.NewInt(0), pack("multiSend", packed), gnosissafe.OperationDelegateCall), nil
}

// ERC20Transfer is a token transfer call.
func ERC20Transfer(token, to common.Address, amount *big.Int) Call {
	return Call{To: token, Value: big.NewInt(0), Data: pack("transfer", to, amount)}
}

// ERC20Burn burns the wallet's own tokens (OpenZeppelin ERC20Burnable's
// burn(uint256)); Trovo-issued tokens are redeemed this way.
func ERC20Burn(token common.Address, amount *big.Int) Call {
	return Call{To: token, Value: big.NewInt(0), Data: pack("burn", amount)}
}

// ERC20Approve is a token approval call.
func ERC20Approve(token, spender common.Address, amount *big.Int) Call {
	return Call{To: token, Value: big.NewInt(0), Data: pack("approve", spender, amount)}
}

// NativeTransfer sends ETH.
func NativeTransfer(to common.Address, wei *big.Int) Call {
	return Call{To: to, Value: wei, Data: nil}
}

// sentinelOwners is Safe's linked-list head for owner management.
var sentinelOwners = common.HexToAddress("0x0000000000000000000000000000000000000001")

// prevOwner finds the owner before target in the Safe's owner list (as
// returned by getOwners), as removeOwner/swapOwner need.
func prevOwner(owners []common.Address, target common.Address) (common.Address, error) {
	for i, o := range owners {
		if o == target {
			if i == 0 {
				return sentinelOwners, nil
			}
			return owners[i-1], nil
		}
	}
	return common.Address{}, fmt.Errorf("aa: %s is not an owner", target.Hex())
}

// SwapOwner replaces oldOwner with newOwner (a call the Safe makes to itself).
func SwapOwner(safe common.Address, currentOwners []common.Address, oldOwner, newOwner common.Address) (Call, error) {
	prev, err := prevOwner(currentOwners, oldOwner)
	if err != nil {
		return Call{}, err
	}
	return Call{To: safe, Value: big.NewInt(0), Data: pack("swapOwner", prev, oldOwner, newOwner)}, nil
}

// AddOwner adds owner and sets the threshold.
func AddOwner(safe, owner common.Address, threshold int64) Call {
	return Call{To: safe, Value: big.NewInt(0), Data: pack("addOwnerWithThreshold", owner, big.NewInt(threshold))}
}

// RemoveOwner removes owner and sets the threshold.
func RemoveOwner(safe common.Address, currentOwners []common.Address, owner common.Address, threshold int64) (Call, error) {
	prev, err := prevOwner(currentOwners, owner)
	if err != nil {
		return Call{}, err
	}
	return Call{To: safe, Value: big.NewInt(0), Data: pack("removeOwner", prev, owner, big.NewInt(threshold))}, nil
}

// ChangeThreshold sets the threshold.
func ChangeThreshold(safe common.Address, threshold int64) Call {
	return Call{To: safe, Value: big.NewInt(0), Data: pack("changeThreshold", big.NewInt(threshold))}
}

// OwnersHash is a stable digest of a Safe's owners and threshold, stored on
// the wallet to detect drift between the database and the chain.
func OwnersHash(owners []common.Address, threshold int64) string {
	var buf []byte
	for _, o := range SortedOwners(owners) {
		buf = append(buf, o.Bytes()...)
	}
	buf = append(buf, common.LeftPadBytes(big.NewInt(threshold).Bytes(), 32)...)
	return common.BytesToHash(crypto.Keccak256(buf)).Hex()
}
