//! Counterfactual Safe addresses.
//!
//! Every Trovo wallet is a Safe (v1.4.1) with the Safe4337Module enabled,
//! owned by the user's mnemonic key. Its address is fixed by CREATE2 before
//! it is ever deployed, so the app can show it (and receive funds on it)
//! from registration on. This module computes that address exactly as
//! `SafeProxyFactory.createProxyWithNonce` will deploy it, and exactly as
//! app-backend's `internal/aa` package does - the backend re-derives it at
//! registration and refuses a mismatch.

use crate::core::{checksum_address, WalletCoreError};
use sha3::{Digest, Keccak256};

/// Canonical Base (and Base Sepolia) deployments this module defaults to.
pub const SAFE_PROXY_FACTORY: &str = "0x4e1DCf7AD4e460CfD30791CCC4F9c8a4f820ec67";
pub const SAFE_L2_SINGLETON: &str = "0x29fcB43b46531BcA003ddC8FCB67FFE91900C762";
pub const SAFE_MODULE_SETUP: &str = "0x2dd68b007B46fBe91B9A7c3EDa5A7a1063cB5b47";
pub const SAFE_4337_MODULE: &str = "0x75cf11467937ce3F2f357CE24ffc3DBF8fD5c226";

/// `type(SafeProxy).creationCode` as returned by the canonical v1.4.1
/// SafeProxyFactory's `proxyCreationCode()`.
pub const SAFE_PROXY_CREATION_CODE: &str = "608060405234801561001057600080fd5b506040516101e63803806101e68339818101604052602081101561003357600080fd5b8101908080519060200190929190505050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156100ca576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001806101c46022913960400191505060405180910390fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505060ab806101196000396000f3fe608060405273ffffffffffffffffffffffffffffffffffffffff600054167fa619486e0000000000000000000000000000000000000000000000000000000060003514156050578060005260206000f35b3660008037600080366000845af43d6000803e60008114156070573d6000fd5b3d6000f3fea264697066735822122003d1488ee65e08fa41e58e888a9865554c535f2c77126a82cb4c0f917f31441364736f6c63430007060033496e76616c69642073696e676c65746f6e20616464726573732070726f7669646564";

/// The contracts a Safe address depends on.
#[derive(Debug, Clone)]
pub struct SafeConfig {
    pub proxy_factory: [u8; 20],
    pub singleton: [u8; 20],
    pub module_setup: [u8; 20],
    pub safe_4337_module: [u8; 20],
    pub proxy_creation_code: Vec<u8>,
}

impl SafeConfig {
    /// The canonical deployments (identical on Base and Base Sepolia).
    pub fn base() -> SafeConfig {
        SafeConfig {
            proxy_factory: parse_address(SAFE_PROXY_FACTORY).expect("valid constant"),
            singleton: parse_address(SAFE_L2_SINGLETON).expect("valid constant"),
            module_setup: parse_address(SAFE_MODULE_SETUP).expect("valid constant"),
            safe_4337_module: parse_address(SAFE_4337_MODULE).expect("valid constant"),
            proxy_creation_code: hex::decode(SAFE_PROXY_CREATION_CODE).expect("valid constant"),
        }
    }
}

pub fn parse_address(s: &str) -> Result<[u8; 20], WalletCoreError> {
    let t = s.trim();
    let t = t.strip_prefix("0x").or_else(|| t.strip_prefix("0X")).unwrap_or(t);
    let bytes = hex::decode(t).map_err(|_| WalletCoreError::InvalidAddress)?;
    bytes.try_into().map_err(|_| WalletCoreError::InvalidAddress)
}

/// Parses a salt nonce given as a decimal or 0x-hex string into a uint256.
pub fn parse_uint256(s: &str) -> Result<[u8; 32], WalletCoreError> {
    let t = s.trim();
    let mut out = [0u8; 32];
    if let Some(h) = t.strip_prefix("0x").or_else(|| t.strip_prefix("0X")) {
        if h.is_empty() || h.len() > 64 {
            return Err(WalletCoreError::InvalidNumber);
        }
        let padded = format!("{:0>64}", h);
        let bytes = hex::decode(padded).map_err(|_| WalletCoreError::InvalidNumber)?;
        out.copy_from_slice(&bytes);
        return Ok(out);
    }
    if t.is_empty() {
        return Err(WalletCoreError::InvalidNumber);
    }
    for c in t.chars() {
        let d = c.to_digit(10).ok_or(WalletCoreError::InvalidNumber)?;
        // out = out * 10 + d, big-endian
        let mut carry = d;
        for byte in out.iter_mut().rev() {
            let v = (*byte as u32) * 10 + carry;
            *byte = (v & 0xff) as u8;
            carry = v >> 8;
        }
        if carry != 0 {
            return Err(WalletCoreError::InvalidNumber);
        }
    }
    Ok(out)
}

fn keccak(data: &[u8]) -> [u8; 32] {
    Keccak256::digest(data).into()
}

fn selector(signature: &str) -> [u8; 4] {
    let h = keccak(signature.as_bytes());
    [h[0], h[1], h[2], h[3]]
}

fn word_address(a: &[u8; 20]) -> [u8; 32] {
    let mut w = [0u8; 32];
    w[12..].copy_from_slice(a);
    w
}

fn word_u64(v: u64) -> [u8; 32] {
    let mut w = [0u8; 32];
    w[24..].copy_from_slice(&v.to_be_bytes());
    w
}

/// `SafeModuleSetup.enableModules([safe_4337_module])` calldata.
fn enable_modules_data(module: &[u8; 20]) -> Vec<u8> {
    let mut out = selector("enableModules(address[])").to_vec();
    out.extend_from_slice(&word_u64(0x20)); // offset of the array
    out.extend_from_slice(&word_u64(1)); // length
    out.extend_from_slice(&word_address(module));
    out
}

/// `Safe.setup(...)` calldata for a Safe owned by `owners` with
/// `threshold`, enabling the Safe4337Module and using it as the fallback
/// handler (so the EntryPoint can call `validateUserOp` on the Safe).
pub fn safe_initializer(owners: &[[u8; 20]], threshold: u64, cfg: &SafeConfig) -> Vec<u8> {
    let data = enable_modules_data(&cfg.safe_4337_module);
    let head_words = 8u64;
    let owners_offset = head_words * 32;
    let data_offset = owners_offset + 32 * (1 + owners.len() as u64);

    let mut out = selector("setup(address[],uint256,address,bytes,address,address,uint256,address)").to_vec();
    out.extend_from_slice(&word_u64(owners_offset));
    out.extend_from_slice(&word_u64(threshold));
    out.extend_from_slice(&word_address(&cfg.module_setup));
    out.extend_from_slice(&word_u64(data_offset));
    out.extend_from_slice(&word_address(&cfg.safe_4337_module)); // fallback handler
    out.extend_from_slice(&[0u8; 32]); // payment token
    out.extend_from_slice(&[0u8; 32]); // payment
    out.extend_from_slice(&[0u8; 32]); // payment receiver
    out.extend_from_slice(&word_u64(owners.len() as u64));
    for o in owners {
        out.extend_from_slice(&word_address(o));
    }
    out.extend_from_slice(&word_u64(data.len() as u64));
    out.extend_from_slice(&data);
    let pad = (32 - data.len() % 32) % 32;
    out.extend(std::iter::repeat(0u8).take(pad));
    out
}

/// The CREATE2 address `createProxyWithNonce(singleton, initializer,
/// salt_nonce)` deploys to.
pub fn safe_address_bytes(owners: &[[u8; 20]], threshold: u64, salt_nonce: &[u8; 32], cfg: &SafeConfig) -> [u8; 20] {
    let initializer = safe_initializer(owners, threshold, cfg);
    let mut salt_input = keccak(&initializer).to_vec();
    salt_input.extend_from_slice(salt_nonce);
    let salt = keccak(&salt_input);

    let mut init_code = cfg.proxy_creation_code.clone();
    init_code.extend_from_slice(&word_address(&cfg.singleton));
    let init_code_hash = keccak(&init_code);

    let mut buf = Vec::with_capacity(85);
    buf.push(0xff);
    buf.extend_from_slice(&cfg.proxy_factory);
    buf.extend_from_slice(&salt);
    buf.extend_from_slice(&init_code_hash);
    let h = keccak(&buf);
    let mut addr = [0u8; 20];
    addr.copy_from_slice(&h[12..]);
    addr
}

/// The EIP-55 address of the 1-of-1 Safe owned by `owner` with
/// `salt_nonce` (a user's primary wallet uses salt nonce "0"), on the
/// canonical Base deployments.
pub fn primary_safe_address(owner: &str, salt_nonce: &str) -> Result<String, WalletCoreError> {
    let owner = parse_address(owner)?;
    let salt = parse_uint256(salt_nonce)?;
    Ok(checksum_address(&safe_address_bytes(&[owner], 1, &salt, &SafeConfig::base())))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parses_salt_nonces() {
        assert_eq!(parse_uint256("0").unwrap(), [0u8; 32]);
        let mut one = [0u8; 32];
        one[31] = 1;
        assert_eq!(parse_uint256("1").unwrap(), one);
        assert_eq!(parse_uint256("0x1").unwrap(), one);
        let mut v = [0u8; 32];
        v[30] = 0x01;
        v[31] = 0x00;
        assert_eq!(parse_uint256("256").unwrap(), v);
        assert!(parse_uint256("abc").is_err());
        assert!(parse_uint256("").is_err());
        // 2^256 overflows
        assert!(parse_uint256("115792089237316195423570985008687907853269984665640564039457584007913129639936").is_err());
        assert!(parse_uint256("115792089237316195423570985008687907853269984665640564039457584007913129639935").is_ok());
    }

    #[test]
    fn known_selectors() {
        assert_eq!(hex::encode(selector("setup(address[],uint256,address,bytes,address,address,uint256,address)")), "b63e800d");
        assert_eq!(hex::encode(selector("enableModules(address[])")), "8d0dc49f");
    }

    #[test]
    fn proxy_creation_code_matches_the_canonical_factory() {
        // keccak256 of SafeProxyFactory v1.4.1 proxyCreationCode(), read
        // from a factory whose runtime code hash equals the canonical one
        let code = hex::decode(SAFE_PROXY_CREATION_CODE).unwrap();
        assert_eq!(code.len(), 486);
        assert_eq!(hex::encode(keccak(&code)), "1856e0ee08399d74e0ea0b03adca210aeade6f748969ac023cdcb4dd62dcaf5f");
    }

    #[test]
    fn address_depends_on_owner_and_salt() {
        let a = primary_safe_address("0x70997970C51812dc3A010C7d01b50e0d17dc79C8", "0").unwrap();
        let b = primary_safe_address("0x70997970C51812dc3A010C7d01b50e0d17dc79C8", "1").unwrap();
        let c = primary_safe_address("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC", "0").unwrap();
        assert_ne!(a, b);
        assert_ne!(a, c);
        assert_eq!(a, primary_safe_address("0x70997970c51812dc3a010c7d01b50e0d17dc79c8", "0x0").unwrap());
        assert!(a.starts_with("0x") && a.len() == 42);
        // pinned with app-backend's internal/aa (TestBaseAddressesMatchWalletCore)
        assert_eq!(a, "0xe7a9D4D8a9633bea8f6f891C5D98744356A9259F");
        assert_eq!(b, "0x3A741746d076eCF518186E8644a70C0982Dd584E");
    }
}

#[cfg(test)]
mod fixture_tests {
    use super::*;

    /// Addresses produced by real Safe v1.4.1 + Safe4337Module deployments
    /// (paymaster/contracts/test/SafeWallet.test.js).
    #[test]
    fn matches_real_deployments() {
        let raw = std::fs::read_to_string(concat!(env!("CARGO_MANIFEST_DIR"), "/testdata/safe_fixture.json")).unwrap();
        let f: serde_json::Value = serde_json::from_str(&raw).unwrap();
        let c = &f["config"];
        let addr = |k: &str| parse_address(c[k].as_str().unwrap()).unwrap();
        let cfg = SafeConfig {
            proxy_factory: addr("proxyFactory"),
            singleton: addr("singleton"),
            module_setup: addr("moduleSetup"),
            safe_4337_module: addr("safe4337Module"),
            proxy_creation_code: hex::decode(c["proxyCreationCode"].as_str().unwrap().trim_start_matches("0x")).unwrap(),
        };
        assert_eq!(hex::encode(&cfg.proxy_creation_code), SAFE_PROXY_CREATION_CODE);
        for case in f["addresses"].as_array().unwrap() {
            let owners: Vec<[u8; 20]> = case["owners"].as_array().unwrap().iter().map(|o| parse_address(o.as_str().unwrap()).unwrap()).collect();
            let threshold = case["threshold"].as_u64().unwrap();
            let init = safe_initializer(&owners, threshold, &cfg);
            assert_eq!(format!("0x{}", hex::encode(&init)), case["initializer"].as_str().unwrap());
            let salt = parse_uint256(case["saltNonce"].as_str().unwrap()).unwrap();
            let got = checksum_address(&safe_address_bytes(&owners, threshold, &salt, &cfg));
            assert_eq!(got, case["address"].as_str().unwrap());
        }
    }
}
