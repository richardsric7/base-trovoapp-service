// Package safesigner assembles and submits Gnosis Safe multisig release
// transactions for P2P escrow settlement (Plan Section 59). This is
// app-backend's own code - per your decision it does not port or depend on
// tm-api's vaultsigner package, though that package's Safe transaction
// conventions were used as a reference for the CSV signer format.
package safesigner

import (
	"fmt"
	"os"
	"strings"
	"trovo-wallet-api/internal/evmkeypair"
)

// signerSeparator matches VaultSignerManagedSecret's existing CSV separator
// in tm-api (tm-api/internal/components/vaultsigner/models/db_models.go).
const signerSeparator = ";"

// legacySignerSeparator is accepted for backward compatibility, matching
// VaultSignerManagedSecret's own legacy fallback.
const legacySignerSeparator = ","

// requiredActiveSigners is the fixed number of signers P2P's own release
// code uses to sign every settlement transaction - the escrow Safe's
// on-chain policy may be 3-of-M with M >= 4 (matching
// VaultSignerManagedSecret's minimum of 4 entries), but P2P deterministically
// signs with only the first 3 configured signers, never an arbitrary
// combination.
const requiredActiveSigners = 3

// ActiveSigners parses P2P_ESCROW_SIGNERS and returns the first 3 signer
// keypairs, in the order they appear in the env var.
func ActiveSigners() ([]*evmkeypair.Full, error) {
	raw := os.Getenv("P2P_ESCROW_SIGNERS")
	if raw == "" {
		return nil, fmt.Errorf("P2P_ESCROW_SIGNERS is not configured")
	}
	entries := splitCSV(raw)
	if len(entries) < requiredActiveSigners {
		return nil, fmt.Errorf("P2P_ESCROW_SIGNERS must contain at least %d signers, found %d", requiredActiveSigners, len(entries))
	}
	signers := make([]*evmkeypair.Full, 0, requiredActiveSigners)
	for i := 0; i < requiredActiveSigners; i++ {
		kp, err := evmkeypair.ParseFull(entries[i])
		if err != nil {
			return nil, fmt.Errorf("P2P_ESCROW_SIGNERS entry %d is invalid: %w", i, err)
		}
		signers = append(signers, kp)
	}
	return signers, nil
}

// splitCSV mirrors VaultSignerManagedSecret's SplitCSV helper: prefer the
// standard ";" separator, fall back to the legacy ",".
func splitCSV(raw string) []string {
	sep := signerSeparator
	if !strings.Contains(raw, signerSeparator) && strings.Contains(raw, legacySignerSeparator) {
		sep = legacySignerSeparator
	}
	parts := strings.Split(raw, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
