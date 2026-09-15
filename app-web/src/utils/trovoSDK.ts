// Base (EVM/secp256k1) wallet operations, backed by the shared wallet-core
// Rust crate (compiled to WASM - see ../walletCore/) instead of this
// module's own Stellar-SDK-based key generation/signing. This is the same
// core the app-mobile Flutter app calls through Dart FFI, so key
// derivation and signing are implemented once and reviewed once, not
// duplicated per platform - see /wallet-core's crate doc.
//
// wasm-bindgen's init() is async (it fetches/instantiates the .wasm file),
// but several call sites need these functions synchronously (e.g.
// `useState<Account>(createAccount())` as a hook initializer). Top-level
// await below blocks this module's own evaluation - and therefore every
// module that imports it, transitively up to the app's entry point -
// until the wasm module is ready, so by the time any component using
// these functions can mount, they're safe to call synchronously.
import init, {
  generateKeypair,
  keypairFromPrivateKey,
  addressFromPrivateKey,
  signPersonal,
  signPersonalBytes,
  generateMnemonic as wcGenerateMnemonic,
  keypairFromMnemonic,
} from '../walletCore/wallet_core.js';

await init();

const createAccount = (): Account => {
  // create a completely new and unique pair of keys.
  const kp = generateKeypair();
  return { publicKey: kp.address, secretKey: kp.privateKeyHex };
};

const signHTTP = (toSign: string, secretKey: string): string => {
  try {
    return signPersonal(secretKey, toSign);
  } catch (error: any) {
    console.log('error signing request', error);
    return '';
  }
};

const importAccount = (secretKey: string): string => {
  try {
    return addressFromPrivateKey(secretKey);
  } catch (error: any) {
    console.log('error signing request', error);
    return '';
  }
};

// signBase64Txn signs a transaction digest. On Base there is no XDR
// envelope to parse and re-hash: transactionDigest is a base64-encoded
// digest computed upstream (by the backend - see
// internal/middleware/security_checks.go's SignBase64Txn), and this
// function's job is only to sign it with EIP-191 personal_sign.
// networkPassphrase is vestigial (kept only so call sites don't need to
// change their argument count).
const signBase64Txn = (
  secretKey: string,
  transactionDigest: string,
  _networkPassphrase: string,
): string => {
  try {
    const raw = atob(transactionDigest);
    const bytes = Uint8Array.from(raw, (c) => c.charCodeAt(0));
    return signPersonalBytes(secretKey, bytes);
  } catch (error: any) {
    console.log('error signing request', error);
    return '';
  }
};

const parseSecretKey = (secretKey: string): Account => {
  const kp = keypairFromPrivateKey(secretKey);
  return { publicKey: kp.address, secretKey: kp.privateKeyHex };
};

const getCredsFromPassPhrase = (passphrase: string): Account | null => {
  try {
    const kp = keypairFromMnemonic(passphrase, 0);
    return { publicKey: kp.address, secretKey: kp.privateKeyHex };
  } catch (error: any) {
    console.log(error);
    return null;
  }
};

const generateMnemonic = (): string => wcGenerateMnemonic();

export {
  createAccount,
  signHTTP,
  importAccount,
  signBase64Txn,
  parseSecretKey,
  getCredsFromPassPhrase,
  generateMnemonic,
};
