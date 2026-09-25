package vaultclient

import (
	"context"
	"errors"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// ErrAlreadyExists means Create was called for a personal secret that's
// already set — personal secrets are create/delete only (Section 5b), never
// an upsert.
var ErrAlreadyExists = errors.New("personal secret already exists")

// ErrNotFound means Delete was called for a personal secret that doesn't
// exist (404 — Section 6).
var ErrNotFound = errors.New("personal secret not found")

// Exists checks metadata only — never pulls the value into memory just to
// answer "does this exist."
func Exists(ctx context.Context, client *vaultapi.Client, mount, path string) (bool, error) {
	_, err := client.KVv2(mount).GetMetadata(ctx, path)
	if err != nil {
		if errors.Is(err, vaultapi.ErrSecretNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("checking existence of %s/%s: %w", mount, path, err)
	}
	return true, nil
}

// Create fails if the secret already exists — this is not an upsert. The
// handler is expected to call Exists first and return 409 before ever
// reaching here, but this is checked again at this layer too, defensively.
func Create(ctx context.Context, client *vaultapi.Client, mount, path, value string) (version int, err error) {
	exists, err := Exists(ctx, client, mount, path)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, ErrAlreadyExists
	}
	secret, err := client.KVv2(mount).Put(ctx, path, map[string]interface{}{"value": value})
	if err != nil {
		return 0, fmt.Errorf("creating %s/%s: %w", mount, path, err)
	}
	return secret.VersionMetadata.Version, nil
}

// Delete fully purges the secret and all of its version history — not a
// soft/versioned delete. A soft delete (client.KVv2(mount).Delete) would
// leave the metadata path resolvable, so Exists would keep reporting
// exists:true after a "delete." DeleteMetadata removes the metadata itself,
// so a subsequent Exists call correctly reports false and a later Create for
// the same prefix starts clean.
func Delete(ctx context.Context, client *vaultapi.Client, mount, path string) error {
	if err := client.KVv2(mount).DeleteMetadata(ctx, path); err != nil {
		return fmt.Errorf("deleting %s/%s: %w", mount, path, err)
	}
	return nil
}

// There is deliberately no Get/Update function in this file — the API
// surface intentionally has no code path capable of reading a personal-env
// value back out of Vault once it's stored (Section 0, Section 5b).
