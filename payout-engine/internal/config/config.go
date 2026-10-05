// Package config is payout-engine's configuration: environment variables
// (see CONFIGURATION.md), and the payout Safe's signers.
package config

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	vaultapi "github.com/hashicorp/vault/api"
)

// Config is read once at start.
type Config struct {
	RPCURL            string
	ChainID           *big.Int // optional; read from the RPC when unset
	PayoutSafe        common.Address
	MultiSendCallOnly common.Address
	OfferBook         common.Address // optional
	IssuingProfile    string         // its wallets (issuing / distribution Safes) are never paid
	ExcludedAddresses []common.Address
	SweepAddress      common.Address // where an admin's sweep sends leftover funds (optional)

	Confirmations      uint64 // blocks behind head treated as final
	LogChunk           uint64 // blocks per eth_getLogs
	CatchUpWithin      uint64 // preparation locks once within this many blocks of the head
	BatchSize          int    // transfers per Safe transaction
	MaxBatchGas        uint64 // a batch is split if its estimate exceeds this
	MinPayoutUnits     *big.Int
	GasPriceMultiplier float64
	PollInterval       time.Duration
	ReceiptTimeout     time.Duration
	DefaultStartBlock  uint64 // fallback start of a token's first scan

	RedisAddr     string
	RedisPassword string
	RedisTLS      bool

	Instance string
	Version  string
}

func env(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}

func envUint(k string, def uint64) uint64 {
	if v, err := strconv.ParseUint(env(k, ""), 10, 64); err == nil {
		return v
	}
	return def
}

func envAddress(k string, required bool) (common.Address, error) {
	v := env(k, "")
	if v == "" {
		if required {
			return common.Address{}, fmt.Errorf("%s is not set", k)
		}
		return common.Address{}, nil
	}
	if !common.IsHexAddress(v) {
		return common.Address{}, fmt.Errorf("%s is not an address: %q", k, v)
	}
	return common.HexToAddress(v), nil
}

// FromEnv reads the configuration.
func FromEnv() (*Config, error) {
	c := &Config{
		RPCURL:             env("BASE_RPC_URL", ""),
		IssuingProfile:     env("TOKENIZATION_ISSUING_PROFILE", "atprofile"),
		Confirmations:      envUint("PROCEED_PAYOUT_CONFIRMATIONS", 3),
		LogChunk:           envUint("PROCEED_PAYOUT_LOG_CHUNK", 5000),
		CatchUpWithin:      envUint("PROCEED_PAYOUT_CATCH_UP_WITHIN", 2),
		BatchSize:          int(envUint("PROCEED_PAYOUT_BATCH_SIZE", 150)),
		MaxBatchGas:        envUint("PROCEED_PAYOUT_MAX_BATCH_GAS", 12_000_000),
		MinPayoutUnits:     new(big.Int).SetUint64(envUint("PROCEED_PAYOUT_MIN_UNITS", 1)),
		GasPriceMultiplier: 1.25,
		PollInterval:       time.Duration(envUint("PROCEED_PAYOUT_POLL_SECONDS", 10)) * time.Second,
		ReceiptTimeout:     time.Duration(envUint("PROCEED_PAYOUT_RECEIPT_TIMEOUT_SECONDS", 180)) * time.Second,
		DefaultStartBlock:  envUint("PROCEED_PAYOUT_DEFAULT_START_BLOCK", 0),
		RedisAddr:          "",
		RedisPassword:      env("REDIS_PASSWORD", ""),
		RedisTLS:           env("REDIS_TLS", "0") == "1",
		Version:            env("APP_VERSION", "dev"),
	}
	if c.RPCURL == "" {
		return nil, fmt.Errorf("BASE_RPC_URL is not set")
	}
	if v := env("BASE_CHAIN_ID", ""); v != "" {
		id, ok := new(big.Int).SetString(v, 10)
		if !ok {
			return nil, fmt.Errorf("BASE_CHAIN_ID is not a number: %q", v)
		}
		c.ChainID = id
	}
	if m, err := strconv.ParseFloat(env("PROCEED_PAYOUT_GAS_PRICE_MULTIPLIER", ""), 64); err == nil && m >= 1 {
		c.GasPriceMultiplier = m
	}
	if c.BatchSize < 1 {
		c.BatchSize = 1
	}
	var err error
	if c.PayoutSafe, err = envAddress("PROCEED_PAYOUT_SAFE_ADDRESS", true); err != nil {
		return nil, err
	}
	// Safe v1.4.1's MultiSendCallOnly (the one app-backend's wallets use)
	if c.MultiSendCallOnly, err = envAddress("SAFE_MULTISEND_CALL_ONLY_ADDRESS", false); err != nil {
		return nil, err
	}
	if c.MultiSendCallOnly == (common.Address{}) {
		c.MultiSendCallOnly = common.HexToAddress("0x9641d764fc13c8B624c04430C7356C1C7C8102e2")
	}
	if c.OfferBook, err = envAddress("OFFER_BOOK_ADDRESS", false); err != nil {
		return nil, err
	}
	if c.SweepAddress, err = envAddress("PROCEED_PAYOUT_SWEEP_ADDRESS", false); err != nil {
		return nil, err
	}
	for _, a := range strings.FieldsFunc(env("PROCEED_PAYOUT_EXCLUDED_ADDRESSES", ""), func(r rune) bool { return r == ',' || r == ';' || r == ' ' }) {
		if !common.IsHexAddress(a) {
			return nil, fmt.Errorf("PROCEED_PAYOUT_EXCLUDED_ADDRESSES has an invalid address %q", a)
		}
		c.ExcludedAddresses = append(c.ExcludedAddresses, common.HexToAddress(a))
	}
	if host := env("REDIS_HOST", ""); host != "" {
		c.RedisAddr = host + ":" + env("REDIS_PORT", "6379")
	}
	host, _ := os.Hostname()
	c.Instance = env("INSTANCE_ID", host)
	return c, nil
}

// requiredSigners is the established standard for platform Safes: they are
// 3-of-N, the managed secret lists at least 3 keys (the vault manager keeps
// at least 4) and the first 3 sign; the first also broadcasts and pays gas.
const requiredSigners = 3

// Signer is one of the payout Safe's signing keys.
type Signer struct {
	Key     *ecdsa.PrivateKey
	Address common.Address
}

// SignerSource yields the payout Safe's active signers.
type SignerSource interface {
	Signers(ctx context.Context) ([]Signer, error)
}

// ParseSigners parses a managed signer secret: keys separated by ";"
// (legacy ","), of which the first 3 are returned.
func ParseSigners(raw, name string) ([]Signer, error) {
	sep := ";"
	if !strings.Contains(raw, ";") && strings.Contains(raw, ",") {
		sep = ","
	}
	var entries []string
	for _, p := range strings.Split(raw, sep) {
		if p = strings.TrimSpace(p); p != "" {
			entries = append(entries, p)
		}
	}
	if len(entries) < requiredSigners {
		return nil, fmt.Errorf("%s must contain at least %d signers, found %d", name, requiredSigners, len(entries))
	}
	out := make([]Signer, 0, requiredSigners)
	for i := 0; i < requiredSigners; i++ {
		k, err := crypto.HexToECDSA(strings.TrimPrefix(entries[i], "0x"))
		if err != nil {
			return nil, fmt.Errorf("%s entry %d is not a private key: %w", name, i+1, err)
		}
		out = append(out, Signer{Key: k, Address: crypto.PubkeyToAddress(k.PublicKey)})
	}
	return out, nil
}

// EnvSigners reads PROCEED_PAYOUT_SIGNERS (injected from the vault manager's
// managed secret at deploy time).
type EnvSigners struct{}

func (EnvSigners) Signers(context.Context) ([]Signer, error) {
	raw := os.Getenv("PROCEED_PAYOUT_SIGNERS")
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("PROCEED_PAYOUT_SIGNERS is not configured")
	}
	return ParseSigners(raw, "PROCEED_PAYOUT_SIGNERS")
}

// VaultSigners reads the managed secret from Vault on every use, so a
// signer the vault manager rotates (it swaps the Safe owner on-chain at the
// same time) is used without restarting the engine.
type VaultSigners struct {
	Client *vaultapi.Client
	Mount  string
	Path   string
	Field  string
}

func (v VaultSigners) Signers(ctx context.Context) ([]Signer, error) {
	secret, err := v.Client.KVv2(v.Mount).Get(ctx, v.Path)
	if err != nil {
		return nil, fmt.Errorf("reading vault secret %s/%s: %w", v.Mount, v.Path, err)
	}
	raw, _ := secret.Data[v.Field].(string)
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("vault secret %s/%s has no %q", v.Mount, v.Path, v.Field)
	}
	return ParseSigners(raw, "PROCEED_PAYOUT_SIGNERS")
}

// SignerSourceFromEnv picks Vault when PROCEED_PAYOUT_SIGNERS_VAULT_PATH
// ("mount/path#field", field defaults to PROCEED_PAYOUT_SIGNERS) and
// VAULT_ADDR / VAULT_TOKEN are set, else the environment.
func SignerSourceFromEnv() (SignerSource, string, error) {
	ref := env("PROCEED_PAYOUT_SIGNERS_VAULT_PATH", "")
	if ref == "" {
		return EnvSigners{}, "PROCEED_PAYOUT_SIGNERS (environment)", nil
	}
	path, field, _ := strings.Cut(ref, "#")
	if field == "" {
		field = "PROCEED_PAYOUT_SIGNERS"
	}
	mount, p, ok := strings.Cut(strings.Trim(path, "/"), "/")
	if !ok {
		return nil, "", fmt.Errorf("PROCEED_PAYOUT_SIGNERS_VAULT_PATH must be mount/path[#field], got %q", ref)
	}
	if env("VAULT_ADDR", "") == "" || env("VAULT_TOKEN", "") == "" {
		return nil, "", fmt.Errorf("PROCEED_PAYOUT_SIGNERS_VAULT_PATH is set but VAULT_ADDR / VAULT_TOKEN are not")
	}
	cfg := vaultapi.DefaultConfig()
	cfg.Address = env("VAULT_ADDR", "")
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, "", err
	}
	client.SetToken(env("VAULT_TOKEN", ""))
	if ns := env("VAULT_NAMESPACE", ""); ns != "" {
		client.SetNamespace(ns)
	}
	return VaultSigners{Client: client, Mount: mount, Path: p, Field: field}, fmt.Sprintf("vault %s/%s#%s", mount, p, field), nil
}
