package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Bootstrap is the only configuration read from the process environment:
// where to find the service's Vault secret. Everything else lives in that
// secret (see CONFIGURATION.md).
type Bootstrap struct {
	Source     string // "vault" (default) or "env" (local development only)
	VaultAddr  string
	VaultToken string
	Namespace  string
	Mount      string
	SecretPath string
}

func BootstrapFromEnv() (Bootstrap, error) {
	b := Bootstrap{
		Source:     strings.ToLower(envOr("PAYMASTER_CONFIG_SOURCE", "vault")),
		VaultAddr:  os.Getenv("VAULT_ADDR"),
		VaultToken: os.Getenv("VAULT_TOKEN"),
		Namespace:  os.Getenv("VAULT_NAMESPACE"),
		Mount:      envOr("VAULT_KV_MOUNT", "secret"),
		SecretPath: envOr("QUOTE_SERVICE_SECRET_PATH", "trovo/paymaster/quote-service"),
	}
	switch b.Source {
	case "env":
	case "vault":
		if b.VaultAddr == "" {
			return b, fmt.Errorf("VAULT_ADDR is not set")
		}
		if b.VaultToken == "" {
			return b, fmt.Errorf("VAULT_TOKEN is not set")
		}
	default:
		return b, fmt.Errorf(`PAYMASTER_CONFIG_SOURCE must be "vault" or "env", got %q`, b.Source)
	}
	return b, nil
}

// Values is the flat key/value view of the service's configuration secret.
type Values map[string]string

// Load reads the configuration values: the Vault KV v2 secret, or the
// process environment when Source is "env".
func (b Bootstrap) Load(ctx context.Context) (Values, error) {
	if b.Source == "env" {
		v := Values{}
		for _, kv := range os.Environ() {
			if k, val, ok := strings.Cut(kv, "="); ok {
				v[k] = val
			}
		}
		return v, nil
	}

	cfg := vaultapi.DefaultConfig()
	cfg.Address = b.VaultAddr
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}
	client.SetToken(b.VaultToken)
	if b.Namespace != "" {
		client.SetNamespace(b.Namespace)
	}
	secret, err := client.KVv2(b.Mount).Get(ctx, b.SecretPath)
	if err != nil {
		return nil, fmt.Errorf("reading vault secret %s/%s: %w", b.Mount, b.SecretPath, err)
	}
	v := Values{}
	for k, val := range secret.Data {
		if val == nil {
			continue
		}
		if s, ok := val.(string); ok {
			v[k] = s
		} else {
			v[k] = fmt.Sprint(val)
		}
	}
	return v, nil
}

// Describe names where the values come from, for logs.
func (b Bootstrap) Describe() string {
	if b.Source == "env" {
		return "process environment"
	}
	return fmt.Sprintf("vault %s/%s", b.Mount, b.SecretPath)
}

func envOr(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}
