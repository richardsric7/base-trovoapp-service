package vaultclient

import (
	"context"
	"errors"
	"fmt"
	"strings"

	vaultsignermodels "admin-panel-dashboard/internal/components/vaultsigner/models"

	vaultapi "github.com/hashicorp/vault/api"
)

// MinCSVEntries is the floor established for every managed secret's CSV —
// enforced both when a value is written (Section 4) and when an assignment
// is deleted (Section 5e).
const MinCSVEntries = 4

// Separator delimits entries in a managed secret's stored CSV value. ";"
// rather than "," - a Base signer value is commonly pasted in from a
// wallet's mnemonic recovery phrase, and some wallets export that phrase
// comma-separated rather than space-separated, so a comma can appear
// inside a single entry. A semicolon can't collide with anything a pasted
// mnemonic or hex private key ever contains.
const Separator = ";"

// legacySeparator is the delimiter every managed secret's CSV used before
// the Base port (Stellar secret seeds never contained a comma, so it was
// safe at the time). SplitCSV falls back to it on read so a managed secret
// registered before the cutover keeps parsing correctly; the next write
// through WriteAtIndex/CollapseCSV always rejoins with Separator, migrating
// that one managed secret's CSV in place with no separate migration step.
const legacySeparator = ","

// SplitCSV splits a managed secret's stored value into its entries. It
// tries Separator first; only when the value contains no Separator at all
// but does contain legacySeparator does it fall back to the pre-cutover
// format - a single entry (hex private key or mnemonic) never legitimately
// contains either character, so a bare comma with no semicolon anywhere in
// the value is unambiguously the old format, not part of one entry.
func SplitCSV(raw string) []string {
	if !strings.Contains(raw, Separator) && strings.Contains(raw, legacySeparator) {
		return strings.Split(raw, legacySeparator)
	}
	return strings.Split(raw, Separator)
}

var (
	// ErrIndexOutOfRange means the requested index has no corresponding CSV
	// entry — checked against the live CSV length, since capacity isn't a
	// stored column (Section 4).
	ErrIndexOutOfRange = errors.New("signer index out of range")

	// ErrCSVBelowMinimum guards the floor from Section 4/5e: a managed
	// secret's CSV must never be written down to fewer than MinCSVEntries.
	ErrCSVBelowMinimum = errors.New("csv would fall below the minimum entry count")
)

// getCSV reads the current CSV value and Vault KV v2 version for a managed
// secret.
func getCSV(ctx context.Context, client *vaultapi.Client, secret vaultsignermodels.VaultSignerManagedSecret) (csvValue string, version int, err error) {
	kv := client.KVv2(secret.VaultMount)
	result, err := kv.Get(ctx, secret.VaultPath)
	if err != nil {
		return "", 0, fmt.Errorf("reading %s/%s: %w", secret.VaultMount, secret.VaultPath, err)
	}
	raw, ok := result.Data[secret.VaultField]
	if !ok {
		return "", 0, fmt.Errorf("field %q not present at %s/%s", secret.VaultField, secret.VaultMount, secret.VaultPath)
	}
	csvValue, ok = raw.(string)
	if !ok {
		return "", 0, fmt.Errorf("field %q at %s/%s is not a string", secret.VaultField, secret.VaultMount, secret.VaultPath)
	}
	return csvValue, result.VersionMetadata.Version, nil
}

// putCSV writes a new CSV value with CAS against the version it's expected
// to still be at, returning the version the write produced.
func putCSV(ctx context.Context, client *vaultapi.Client, secret vaultsignermodels.VaultSignerManagedSecret, csvValue string, casVersion int) (newVersion int, err error) {
	kv := client.KVv2(secret.VaultMount)
	result, err := kv.Put(ctx, secret.VaultPath, map[string]interface{}{
		secret.VaultField: csvValue,
	}, vaultapi.WithCheckAndSet(casVersion))
	if err != nil {
		return 0, fmt.Errorf("writing %s/%s: %w", secret.VaultMount, secret.VaultPath, err)
	}
	return result.VersionMetadata.Version, nil
}

// ReadValue returns the current CSV value at index (Section 4a's "Read +
// version" step) — used by the GET-my-slot handler and by the swap flow's
// "read the old key" step.
func ReadValue(ctx context.Context, client *vaultapi.Client, secret vaultsignermodels.VaultSignerManagedSecret, index int) (value string, version int, err error) {
	csvValue, version, err := getCSV(ctx, client, secret)
	if err != nil {
		return "", 0, err
	}
	parts := SplitCSV(csvValue)
	if index < 0 || index >= len(parts) {
		return "", 0, ErrIndexOutOfRange
	}
	return parts[index], version, nil
}

// ReadCSV returns the raw CSV value and its current Vault version — used by
// Section 5e's assignment-delete flow, which needs the whole CSV (to count
// valid keypairs and to snapshot the pre-collapse content), not just one
// index.
func ReadCSV(ctx context.Context, client *vaultapi.Client, secret vaultsignermodels.VaultSignerManagedSecret) (csvValue string, version int, err error) {
	return getCSV(ctx, client, secret)
}

// WriteAtIndex rebuilds the CSV with newValue at index and writes it back
// with CAS (Section 4a's "Rebuild CSV" + "Write with CAS" steps). newValue
// must already have passed the Section 5c validation pipeline before this is
// called — this function performs no key validation, only CSV-shape checks.
func WriteAtIndex(ctx context.Context, client *vaultapi.Client, secret vaultsignermodels.VaultSignerManagedSecret, index int, newValue string) (versionBefore, versionAfter int, err error) {
	csvValue, version, err := getCSV(ctx, client, secret)
	if err != nil {
		return 0, 0, err
	}
	parts := SplitCSV(csvValue)
	if index < 0 || index >= len(parts) {
		return 0, 0, ErrIndexOutOfRange
	}
	if len(parts) < MinCSVEntries {
		return 0, 0, ErrCSVBelowMinimum
	}
	parts[index] = newValue
	newVersion, err := putCSV(ctx, client, secret, strings.Join(parts, Separator), version)
	if err != nil {
		return 0, 0, err
	}
	return version, newVersion, nil
}

// CollapseCSV removes the entry at `index` and shifts everything after it
// left by one (Section 5e's assignment-delete flow). Takes the
// already-read csvValue/version rather than re-reading, since the caller
// needs that same read for its own pre-collapse snapshot before this
// function ever runs. Returns the version the collapse produced.
func CollapseCSV(ctx context.Context, client *vaultapi.Client, secret vaultsignermodels.VaultSignerManagedSecret, csvValue string, version, index int) (newVersion int, err error) {
	parts := SplitCSV(csvValue)
	if index < 0 || index >= len(parts) {
		return 0, ErrIndexOutOfRange
	}
	parts = append(parts[:index], parts[index+1:]...)
	return putCSV(ctx, client, secret, strings.Join(parts, Separator), version)
}

// CountValidKeypairs reports how many entries in parts are structurally
// valid Base keypairs — the Section 5e trigger for whether an assignment
// deletion needs to touch the chain at all.
func CountValidKeypairs(parts []string, isValidKeypair func(string) bool) int {
	count := 0
	for _, raw := range parts {
		if isValidKeypair(strings.TrimSpace(raw)) {
			count++
		}
	}
	return count
}
