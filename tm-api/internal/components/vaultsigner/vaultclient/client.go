package vaultclient

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// NewClient builds a Vault API client from VAULT_ADDR/VAULT_TOKEN, matching
// how internal/network reads EXPANSION_URL/BLOCKCHAIN_NETWORK_PASSPHRASE via
// plain os.Getenv rather than the envconfig-based config/config.go, which is
// scoped to DB/port/env fields only.
func NewClient() (*vaultapi.Client, error) {
	addr := os.Getenv("VAULT_ADDR")
	if addr == "" {
		return nil, fmt.Errorf("VAULT_ADDR is not set")
	}
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("VAULT_TOKEN is not set")
	}

	config := vaultapi.DefaultConfig()
	config.Address = addr
	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("could not create vault client: %w", err)
	}
	client.SetToken(token)
	return client, nil
}
