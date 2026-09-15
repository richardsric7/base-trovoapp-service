//! WASM bindings (wasm-bindgen) over [`crate::core`] for app-web.
//! Build with `wasm-pack build --target web` - see README.md.

use crate::core;
use wasm_bindgen::prelude::*;

fn to_js_err(e: core::WalletCoreError) -> JsValue {
    JsValue::from_str(&e.to_string())
}

/// A generated/imported keypair, returned to JS as a plain object
/// ({ privateKeyHex, address }).
#[wasm_bindgen]
pub struct WcKeypair {
    #[wasm_bindgen(getter_with_clone, js_name = privateKeyHex)]
    pub private_key_hex: String,
    #[wasm_bindgen(getter_with_clone)]
    pub address: String,
}

impl From<core::Keypair> for WcKeypair {
    fn from(kp: core::Keypair) -> Self {
        WcKeypair {
            private_key_hex: kp.private_key_hex,
            address: kp.address,
        }
    }
}

/// Generates a fresh keypair - the Base equivalent of trovoSDK.ts's
/// `createAccount()`.
#[wasm_bindgen(js_name = generateKeypair)]
pub fn generate_keypair() -> WcKeypair {
    core::generate_keypair().into()
}

/// Parses a hex private key and returns its (privateKeyHex, address) pair
/// - the Base equivalent of trovoSDK.ts's `parseSecretKey`.
#[wasm_bindgen(js_name = keypairFromPrivateKey)]
pub fn keypair_from_private_key(private_key_hex: &str) -> Result<WcKeypair, JsValue> {
    core::keypair_from_private_key(private_key_hex)
        .map(Into::into)
        .map_err(to_js_err)
}

/// Returns just the address for a hex private key - the Base equivalent
/// of trovoSDK.ts's `importAccount`.
#[wasm_bindgen(js_name = addressFromPrivateKey)]
pub fn address_from_private_key(private_key_hex: &str) -> Result<String, JsValue> {
    core::address_from_private_key(private_key_hex).map_err(to_js_err)
}

/// Signs `message` (UTF-8 text) with EIP-191 personal_sign, returning a
/// base64 signature - the Base equivalent of trovoSDK.ts's `signHTTP` and
/// `signBase64Txn` (both now sign a message/digest the same way; see the
/// package doc and internal/middleware/security_checks.go on the backend).
#[wasm_bindgen(js_name = signPersonal)]
pub fn sign_personal(private_key_hex: &str, message: &str) -> Result<String, JsValue> {
    core::sign_personal(private_key_hex, message.as_bytes()).map_err(to_js_err)
}

/// Signs raw bytes (e.g. a pre-computed transaction digest) with EIP-191
/// personal_sign, returning a base64 signature.
#[wasm_bindgen(js_name = signPersonalBytes)]
pub fn sign_personal_bytes(private_key_hex: &str, message: &[u8]) -> Result<String, JsValue> {
    core::sign_personal(private_key_hex, message).map_err(to_js_err)
}

/// Verifies a base64 EIP-191 personal_sign signature against `address`.
#[wasm_bindgen(js_name = verifyPersonal)]
pub fn verify_personal(address: &str, message: &str, signature_b64: &str) -> Result<bool, JsValue> {
    core::verify_personal(address, message.as_bytes(), signature_b64).map_err(to_js_err)
}

/// Recovers the signer address from a base64 EIP-191 personal_sign
/// signature, without needing a claimed address to check against.
#[wasm_bindgen(js_name = recoverPersonalSigner)]
pub fn recover_personal_signer(message: &str, signature_b64: &str) -> Result<String, JsValue> {
    core::recover_personal_signer(message.as_bytes(), signature_b64).map_err(to_js_err)
}

/// Generates a fresh 12-word BIP39 mnemonic.
#[wasm_bindgen(js_name = generateMnemonic)]
pub fn generate_mnemonic() -> Result<String, JsValue> {
    core::generate_mnemonic().map_err(to_js_err)
}

/// Derives a keypair from a BIP39 mnemonic at BIP44 path
/// m/44'/60'/0'/0/{index} - the Base equivalent of trovoSDK.ts's
/// `getCredsFromPassPhrase`.
#[wasm_bindgen(js_name = keypairFromMnemonic)]
pub fn keypair_from_mnemonic(mnemonic: &str, index: u32) -> Result<WcKeypair, JsValue> {
    core::keypair_from_mnemonic(mnemonic, index)
        .map(Into::into)
        .map_err(to_js_err)
}
