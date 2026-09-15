//! C ABI bindings over [`crate::core`] for app-mobile (consumed through
//! Dart's `dart:ffi`). Every fallible function returns a heap-allocated,
//! NUL-terminated JSON C string - `{"...fields":...}` on success, or
//! `{"error":"..."}` on failure - so Dart-side call sites can uniformly
//! `jsonDecode` the result and branch on an `error` key, without needing
//! a second FFI call to fetch error text. Every returned `*mut c_char`
//! must be released with [`wc_free_string`] once decoded, to avoid
//! leaking the Rust-allocated buffer across the FFI boundary.

use crate::core;
use serde::Serialize;
use std::ffi::{CStr, CString};
use std::os::raw::c_char;

fn ok_json<T: Serialize>(value: &T) -> *mut c_char {
    let s = serde_json::to_string(value)
        .unwrap_or_else(|_| "{\"error\":\"encode failed\"}".to_string());
    CString::new(s).unwrap_or_default().into_raw()
}

fn err_json(e: impl std::fmt::Display) -> *mut c_char {
    let s = serde_json::json!({ "error": e.to_string() }).to_string();
    CString::new(s).unwrap_or_default().into_raw()
}

/// Reads a NUL-terminated C string from `ptr` as UTF-8. Returns an error
/// string (not a panic) on a null/invalid pointer, since these functions
/// are called across the FFI boundary where Rust cannot rely on Dart
/// having upheld its side of the contract.
unsafe fn read_cstr(ptr: *const c_char) -> Result<String, ()> {
    if ptr.is_null() {
        return Err(());
    }
    CStr::from_ptr(ptr)
        .to_str()
        .map(str::to_owned)
        .map_err(|_| ())
}

#[derive(Serialize)]
struct KeypairJson {
    #[serde(rename = "privateKeyHex")]
    private_key_hex: String,
    address: String,
}

impl From<core::Keypair> for KeypairJson {
    fn from(kp: core::Keypair) -> Self {
        KeypairJson {
            private_key_hex: kp.private_key_hex,
            address: kp.address,
        }
    }
}

/// Frees a string previously returned by any `wc_*` function in this
/// module. Calling this on any other pointer (or twice on the same
/// pointer) is undefined behavior, same as `free()`.
#[no_mangle]
pub extern "C" fn wc_free_string(ptr: *mut c_char) {
    if ptr.is_null() {
        return;
    }
    unsafe {
        drop(CString::from_raw(ptr));
    }
}

/// Generates a fresh keypair - the Base equivalent of trovo-sdk.dart's
/// `createAccount()`. Returns `{"privateKeyHex","address"}`.
#[no_mangle]
pub extern "C" fn wc_generate_keypair() -> *mut c_char {
    ok_json(&KeypairJson::from(core::generate_keypair()))
}

/// Parses a hex private key and returns its (privateKeyHex, address) pair
/// - the Base equivalent of trovo-sdk.dart's `parseSecretKey`.
#[no_mangle]
pub extern "C" fn wc_keypair_from_private_key(private_key_hex: *const c_char) -> *mut c_char {
    let Ok(private_key_hex) = (unsafe { read_cstr(private_key_hex) }) else {
        return err_json("invalid private key argument");
    };
    match core::keypair_from_private_key(&private_key_hex) {
        Ok(kp) => ok_json(&KeypairJson::from(kp)),
        Err(e) => err_json(e),
    }
}

/// Returns just `{"address":...}` for a hex private key - the Base
/// equivalent of trovo-sdk.dart's `importAccount`.
#[no_mangle]
pub extern "C" fn wc_address_from_private_key(private_key_hex: *const c_char) -> *mut c_char {
    let Ok(private_key_hex) = (unsafe { read_cstr(private_key_hex) }) else {
        return err_json("invalid private key argument");
    };
    match core::address_from_private_key(&private_key_hex) {
        Ok(address) => ok_json(&serde_json::json!({ "address": address })),
        Err(e) => err_json(e),
    }
}

/// Signs UTF-8 `message` with EIP-191 personal_sign, returning
/// `{"signature":...}` (base64) - the Base equivalent of trovo-sdk.dart's
/// `signHTTP`/`signBase64Txn` (both now sign a message/digest the same
/// way; see the crate doc and internal/middleware/security_checks.go on
/// the backend).
#[no_mangle]
pub extern "C" fn wc_sign_personal(
    private_key_hex: *const c_char,
    message: *const c_char,
) -> *mut c_char {
    let (Ok(private_key_hex), Ok(message)) = (unsafe { read_cstr(private_key_hex) }, unsafe {
        read_cstr(message)
    }) else {
        return err_json("invalid argument");
    };
    match core::sign_personal(&private_key_hex, message.as_bytes()) {
        Ok(signature) => ok_json(&serde_json::json!({ "signature": signature })),
        Err(e) => err_json(e),
    }
}

/// Signs raw bytes (e.g. a pre-computed transaction digest, which is
/// arbitrary binary data - not valid UTF-8 in general, so it cannot go
/// through [`wc_sign_personal`]'s NUL-terminated C string parameter) with
/// EIP-191 personal_sign, returning `{"signature":...}` (base64). Callers
/// own `message_ptr`'s buffer and may free it immediately after this call
/// returns - it is only read, never retained.
#[no_mangle]
pub extern "C" fn wc_sign_personal_bytes(
    private_key_hex: *const c_char,
    message_ptr: *const u8,
    message_len: usize,
) -> *mut c_char {
    let Ok(private_key_hex) = (unsafe { read_cstr(private_key_hex) }) else {
        return err_json("invalid private key argument");
    };
    if message_ptr.is_null() && message_len != 0 {
        return err_json("invalid message argument");
    }
    let message = if message_len == 0 {
        &[][..]
    } else {
        unsafe { std::slice::from_raw_parts(message_ptr, message_len) }
    };
    match core::sign_personal(&private_key_hex, message) {
        Ok(signature) => ok_json(&serde_json::json!({ "signature": signature })),
        Err(e) => err_json(e),
    }
}

/// Verifies a base64 EIP-191 personal_sign signature against `address`.
/// Returns 1 (valid), 0 (invalid), or -1 (a malformed argument/signature
/// - distinct from "invalid" so callers don't mistake a decode failure
/// for a confirmed-bad signature).
#[no_mangle]
pub extern "C" fn wc_verify_personal(
    address: *const c_char,
    message: *const c_char,
    signature_b64: *const c_char,
) -> i32 {
    let (Ok(address), Ok(message), Ok(signature_b64)) = (
        unsafe { read_cstr(address) },
        unsafe { read_cstr(message) },
        unsafe { read_cstr(signature_b64) },
    ) else {
        return -1;
    };
    match core::verify_personal(&address, message.as_bytes(), &signature_b64) {
        Ok(true) => 1,
        Ok(false) => 0,
        Err(_) => -1,
    }
}

/// Recovers the signer address from a base64 EIP-191 personal_sign
/// signature. Returns `{"address":...}`.
#[no_mangle]
pub extern "C" fn wc_recover_personal_signer(
    message: *const c_char,
    signature_b64: *const c_char,
) -> *mut c_char {
    let (Ok(message), Ok(signature_b64)) = (unsafe { read_cstr(message) }, unsafe {
        read_cstr(signature_b64)
    }) else {
        return err_json("invalid argument");
    };
    match core::recover_personal_signer(message.as_bytes(), &signature_b64) {
        Ok(address) => ok_json(&serde_json::json!({ "address": address })),
        Err(e) => err_json(e),
    }
}

/// Generates a fresh 12-word BIP39 mnemonic. Returns `{"mnemonic":...}`.
#[no_mangle]
pub extern "C" fn wc_generate_mnemonic() -> *mut c_char {
    match core::generate_mnemonic() {
        Ok(mnemonic) => ok_json(&serde_json::json!({ "mnemonic": mnemonic })),
        Err(e) => err_json(e),
    }
}

/// Derives a keypair from a BIP39 mnemonic at BIP44 path
/// m/44'/60'/0'/0/{index} - the Base equivalent of trovo-sdk.dart's
/// `generateCredentialsFromPassPhrase`/`retrieveCredentialsFromPassPhrase`.
#[no_mangle]
pub extern "C" fn wc_keypair_from_mnemonic(mnemonic: *const c_char, index: u32) -> *mut c_char {
    let Ok(mnemonic) = (unsafe { read_cstr(mnemonic) }) else {
        return err_json("invalid mnemonic argument");
    };
    match core::keypair_from_mnemonic(&mnemonic, index) {
        Ok(kp) => ok_json(&KeypairJson::from(kp)),
        Err(e) => err_json(e),
    }
}
