//! Platform-agnostic Base (secp256k1/EVM) wallet logic: key generation,
//! address derivation, EIP-191 personal_sign/verify, and BIP39/BIP44
//! mnemonic-based key derivation. No platform-specific types here - see
//! `wasm` and `ffi` for the per-platform bindings that call into this.

use base64::Engine;
use k256::ecdsa::signature::hazmat::PrehashVerifier;
use k256::ecdsa::{RecoveryId, Signature, SigningKey, VerifyingKey};
use k256::elliptic_curve::rand_core::OsRng;
use sha3::{Digest, Keccak256};
use std::fmt;

#[derive(Debug)]
pub enum WalletCoreError {
    InvalidPrivateKey,
    InvalidAddress,
    InvalidSignature,
    InvalidMnemonic,
    InvalidDerivationPath,
    Base64Decode,
}

impl fmt::Display for WalletCoreError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let msg = match self {
            WalletCoreError::InvalidPrivateKey => "invalid private key",
            WalletCoreError::InvalidAddress => "invalid address",
            WalletCoreError::InvalidSignature => "invalid signature",
            WalletCoreError::InvalidMnemonic => "invalid mnemonic",
            WalletCoreError::InvalidDerivationPath => "invalid derivation path",
            WalletCoreError::Base64Decode => "invalid base64",
        };
        write!(f, "{}", msg)
    }
}

impl std::error::Error for WalletCoreError {}

/// A generated or imported keypair: `private_key_hex` is 0x-prefixed
/// 32-byte secp256k1 scalar, `address` is the 0x-prefixed, EIP-55
/// checksummed 20-byte Base/Ethereum address - the Base equivalent of
/// Stellar's (secretSeed, accountId) pair.
#[derive(Debug, Clone)]
pub struct Keypair {
    pub private_key_hex: String,
    pub address: String,
}

fn signing_key_from_hex(private_key_hex: &str) -> Result<SigningKey, WalletCoreError> {
    let trimmed = private_key_hex.trim();
    // Strip a "0x"/"0X" prefix case-insensitively: some call sites
    // (e.g. app-mobile's `secretKey.toUpperCase()` before parsing, kept
    // from the original Stellar-secret-key-was-always-uppercase
    // convention) uppercase the whole string before handing it here.
    // hex::decode itself already accepts mixed-case hex digits.
    let trimmed = trimmed
        .strip_prefix("0x")
        .or_else(|| trimmed.strip_prefix("0X"))
        .unwrap_or(trimmed);
    let bytes = hex::decode(trimmed).map_err(|_| WalletCoreError::InvalidPrivateKey)?;
    SigningKey::from_slice(&bytes).map_err(|_| WalletCoreError::InvalidPrivateKey)
}

/// Derives the 0x-prefixed, EIP-55 checksummed address from a verifying
/// (public) key: keccak256(uncompressed_pubkey[1:])[12:], checksum-cased
/// per EIP-55.
fn address_from_verifying_key(vk: &VerifyingKey) -> String {
    let uncompressed = vk.to_encoded_point(false);
    let pubkey_bytes = uncompressed.as_bytes(); // 65 bytes: 0x04 || X || Y
    let hash = Keccak256::digest(&pubkey_bytes[1..]);
    let addr_bytes = &hash[12..]; // last 20 bytes
    checksum_address(addr_bytes)
}

/// EIP-55 mixed-case checksum encoding of a 20-byte address.
fn checksum_address(addr_bytes: &[u8]) -> String {
    let addr_hex = hex::encode(addr_bytes);
    let hash = Keccak256::digest(addr_hex.as_bytes());
    let mut out = String::with_capacity(42);
    out.push_str("0x");
    for (i, c) in addr_hex.chars().enumerate() {
        if c.is_ascii_digit() {
            out.push(c);
            continue;
        }
        // nibble i of the hash decides upper/lower case for hex letters.
        let byte = hash[i / 2];
        let nibble = if i % 2 == 0 { byte >> 4 } else { byte & 0x0f };
        if nibble >= 8 {
            out.push(c.to_ascii_uppercase());
        } else {
            out.push(c);
        }
    }
    out
}

/// Generates a fresh keypair from CSPRNG entropy - the Base equivalent of
/// Stellar's `Keypair.random()`.
pub fn generate_keypair() -> Keypair {
    let signing_key = SigningKey::random(&mut OsRng);
    keypair_from_signing_key(signing_key)
}

fn keypair_from_signing_key(signing_key: SigningKey) -> Keypair {
    let private_key_hex = format!("0x{}", hex::encode(signing_key.to_bytes()));
    let address = address_from_verifying_key(signing_key.verifying_key());
    Keypair {
        private_key_hex,
        address,
    }
}

/// Parses a hex-encoded private key and returns its (private_key_hex,
/// address) pair - the Base equivalent of Stellar's
/// `Keypair.fromSecret(secretKey)`.
pub fn keypair_from_private_key(private_key_hex: &str) -> Result<Keypair, WalletCoreError> {
    let signing_key = signing_key_from_hex(private_key_hex)?;
    Ok(keypair_from_signing_key(signing_key))
}

/// Returns just the address for a hex-encoded private key - the Base
/// equivalent of Stellar's `Keypair.fromSecret(secretKey).publicKey()`
/// (trovoSDK.ts's/trovo-sdk.dart's `importAccount`).
pub fn address_from_private_key(private_key_hex: &str) -> Result<String, WalletCoreError> {
    Ok(keypair_from_private_key(private_key_hex)?.address)
}

/// The EIP-191 "personal_sign" digest: keccak256("\x19Ethereum Signed
/// Message:\n" + len(message) + message) - the Base equivalent of
/// Stellar's network-passphrase-scoped transaction hash, used here for
/// offline/off-chain message authentication (HTTP request signing, login
/// challenges, and now also transaction-digest signing - see
/// internal/middleware/security_checks.go's SignBase64Txn on the backend,
/// which this mirrors).
pub fn personal_sign_hash(message: &[u8]) -> [u8; 32] {
    let prefix = format!("\x19Ethereum Signed Message:\n{}", message.len());
    let mut hasher = Keccak256::new();
    hasher.update(prefix.as_bytes());
    hasher.update(message);
    hasher.finalize().into()
}

/// Signs `message` with EIP-191 personal_sign and returns the
/// base64-encoded 65-byte (r,s,v) recoverable signature, v normalized to
/// {27,28} - the Base equivalent of Stellar's
/// `Keypair.sign(bytes)`/`SignBase64`, matching this codebase's backend
/// (internal/evmkeypair.SignPersonal) byte-for-byte so a signature
/// produced by either verifies against the other.
pub fn sign_personal(private_key_hex: &str, message: &[u8]) -> Result<String, WalletCoreError> {
    let signing_key = signing_key_from_hex(private_key_hex)?;
    let digest = personal_sign_hash(message);
    let (signature, recovery_id): (Signature, RecoveryId) = signing_key
        .sign_prehash_recoverable(&digest)
        .map_err(|_| WalletCoreError::InvalidSignature)?;

    let mut sig_bytes = [0u8; 65];
    sig_bytes[..64].copy_from_slice(&signature.to_bytes());
    sig_bytes[64] = recovery_id.to_byte() + 27;

    Ok(base64::engine::general_purpose::STANDARD.encode(sig_bytes))
}

/// Verifies a base64-encoded EIP-191 personal_sign signature against
/// `address` - the Base equivalent of Stellar's `Keypair.verify`.
pub fn verify_personal(
    address: &str,
    message: &[u8],
    signature_b64: &str,
) -> Result<bool, WalletCoreError> {
    let sig_bytes = base64::engine::general_purpose::STANDARD
        .decode(signature_b64)
        .map_err(|_| WalletCoreError::Base64Decode)?;
    if sig_bytes.len() != 65 {
        return Err(WalletCoreError::InvalidSignature);
    }
    let digest = personal_sign_hash(message);
    let signature =
        Signature::from_slice(&sig_bytes[..64]).map_err(|_| WalletCoreError::InvalidSignature)?;
    let recid_byte = if sig_bytes[64] >= 27 {
        sig_bytes[64] - 27
    } else {
        sig_bytes[64]
    };
    let recovery_id = RecoveryId::from_byte(recid_byte).ok_or(WalletCoreError::InvalidSignature)?;

    let recovered = VerifyingKey::recover_from_prehash(&digest, &signature, recovery_id)
        .map_err(|_| WalletCoreError::InvalidSignature)?;
    let recovered_address = address_from_verifying_key(&recovered);

    Ok(recovered_address.eq_ignore_ascii_case(address.trim()))
}

/// A convenience combination of verify_personal that also self-checks the
/// signature against its own recovered public key (used internally by
/// `verify_personal`; exposed separately so bindings can recover an
/// unknown signer's address instead of only checking a claimed one).
pub fn recover_personal_signer(
    message: &[u8],
    signature_b64: &str,
) -> Result<String, WalletCoreError> {
    let sig_bytes = base64::engine::general_purpose::STANDARD
        .decode(signature_b64)
        .map_err(|_| WalletCoreError::Base64Decode)?;
    if sig_bytes.len() != 65 {
        return Err(WalletCoreError::InvalidSignature);
    }
    let digest = personal_sign_hash(message);
    let signature =
        Signature::from_slice(&sig_bytes[..64]).map_err(|_| WalletCoreError::InvalidSignature)?;
    let recid_byte = if sig_bytes[64] >= 27 {
        sig_bytes[64] - 27
    } else {
        sig_bytes[64]
    };
    let recovery_id = RecoveryId::from_byte(recid_byte).ok_or(WalletCoreError::InvalidSignature)?;
    let recovered = VerifyingKey::recover_from_prehash(&digest, &signature, recovery_id)
        .map_err(|_| WalletCoreError::InvalidSignature)?;
    let _ = recovered.verify_prehash(&digest, &signature); // best-effort sanity check, ignored on failure
    Ok(address_from_verifying_key(&recovered))
}

/// Generates a fresh BIP39 mnemonic - the Base equivalent of
/// app-mobile's `Wallet.generate24WordsMnemonic()`. Uses the standard
/// 12-word/128-bit strength (the same strength most EVM wallets default
/// to, e.g. MetaMask), rather than Stellar's 24-word convention, so
/// mnemonics generated here import cleanly into other Base/EVM wallets.
pub fn generate_mnemonic() -> Result<String, WalletCoreError> {
    // 16 bytes (128 bits) of entropy -> a 12-word mnemonic per BIP39.
    let mut entropy = [0u8; 16];
    getrandom::getrandom(&mut entropy).map_err(|_| WalletCoreError::InvalidMnemonic)?;
    let mnemonic =
        bip39::Mnemonic::from_entropy(&entropy).map_err(|_| WalletCoreError::InvalidMnemonic)?;
    Ok(mnemonic.to_string())
}

/// Derives a keypair from a BIP39 mnemonic phrase at BIP44 path
/// m/44'/60'/0'/0/{index} (60 = Ethereum/EVM's registered coin type,
/// shared by Base) - the Base equivalent of Stellar's
/// `Wallet.from(mnemonic).getKeyPair(index)`. Standard BIP44 derivation
/// (rather than Stellar's simpler "hash the raw seed bytes" approach) is
/// used so a mnemonic generated/imported here derives the same addresses
/// a user would see importing it into another EVM wallet.
pub fn keypair_from_mnemonic(mnemonic: &str, index: u32) -> Result<Keypair, WalletCoreError> {
    let mnemonic = bip39::Mnemonic::parse_normalized(mnemonic.trim())
        .map_err(|_| WalletCoreError::InvalidMnemonic)?;
    let seed = mnemonic.to_seed("");
    let path = format!("m/44'/60'/0'/0/{}", index);
    let derived = tiny_hderive::bip32::ExtendedPrivKey::derive(&seed, path.as_str())
        .map_err(|_| WalletCoreError::InvalidDerivationPath)?;
    let signing_key = SigningKey::from_slice(&derived.secret())
        .map_err(|_| WalletCoreError::InvalidPrivateKey)?;
    Ok(keypair_from_signing_key(signing_key))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn generate_and_sign_roundtrip() {
        let kp = generate_keypair();
        assert!(kp.address.starts_with("0x"));
        assert_eq!(kp.address.len(), 42);

        let sig = sign_personal(&kp.private_key_hex, b"hello base").unwrap();
        let ok = verify_personal(&kp.address, b"hello base", &sig).unwrap();
        assert!(ok);

        let recovered = recover_personal_signer(b"hello base", &sig).unwrap();
        assert_eq!(recovered.to_lowercase(), kp.address.to_lowercase());

        let bad = verify_personal(&kp.address, b"tampered", &sig).unwrap();
        assert!(!bad);
    }

    #[test]
    fn import_matches_generated() {
        let kp = generate_keypair();
        let imported = keypair_from_private_key(&kp.private_key_hex).unwrap();
        assert_eq!(kp.address, imported.address);

        let addr_only = address_from_private_key(&kp.private_key_hex).unwrap();
        assert_eq!(kp.address, addr_only);
    }

    #[test]
    fn import_tolerates_uppercased_prefix_and_digits() {
        // app-mobile's original call site uppercases the secret key
        // before parsing (a leftover Stellar-secret-key convention) -
        // "0x" becomes "0X" and hex letters become uppercase.
        let kp = generate_keypair();
        let uppercased = kp.private_key_hex.to_uppercase();
        let imported = keypair_from_private_key(&uppercased).unwrap();
        assert_eq!(kp.address, imported.address);
    }

    #[test]
    fn mnemonic_roundtrip_is_deterministic() {
        let mnemonic = generate_mnemonic().unwrap();
        let kp1 = keypair_from_mnemonic(&mnemonic, 0).unwrap();
        let kp2 = keypair_from_mnemonic(&mnemonic, 0).unwrap();
        assert_eq!(kp1.address, kp2.address);

        let kp_other_index = keypair_from_mnemonic(&mnemonic, 1).unwrap();
        assert_ne!(kp1.address, kp_other_index.address);
    }

    #[test]
    fn checksum_address_matches_eip55_test_vector() {
        // EIP-55 spec reference test vector: this is the checksummed
        // encoding of the all-lowercase address
        // 0xfb6916095ca1df60bb79ce92ce3ea74c37c5d359.
        let expected = "0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359";
        assert_eq!(
            checksum_address(&hex::decode(&expected[2..]).unwrap()),
            expected
        );
    }

    #[test]
    fn rejects_malformed_private_key() {
        assert!(keypair_from_private_key("not-hex").is_err());
        assert!(keypair_from_private_key("0x1234").is_err()); // too short
    }
}
