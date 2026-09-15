/* tslint:disable */
/* eslint-disable */

/**
 * A generated/imported keypair, returned to JS as a plain object
 * ({ privateKeyHex, address }).
 */
export class WcKeypair {
    private constructor();
    free(): void;
    [Symbol.dispose](): void;
    address: string;
    privateKeyHex: string;
}

/**
 * Returns just the address for a hex private key - the Base equivalent
 * of trovoSDK.ts's `importAccount`.
 */
export function addressFromPrivateKey(private_key_hex: string): string;

/**
 * Generates a fresh keypair - the Base equivalent of trovoSDK.ts's
 * `createAccount()`.
 */
export function generateKeypair(): WcKeypair;

/**
 * Generates a fresh 12-word BIP39 mnemonic.
 */
export function generateMnemonic(): string;

/**
 * Derives a keypair from a BIP39 mnemonic at BIP44 path
 * m/44'/60'/0'/0/{index} - the Base equivalent of trovoSDK.ts's
 * `getCredsFromPassPhrase`.
 */
export function keypairFromMnemonic(mnemonic: string, index: number): WcKeypair;

/**
 * Parses a hex private key and returns its (privateKeyHex, address) pair
 * - the Base equivalent of trovoSDK.ts's `parseSecretKey`.
 */
export function keypairFromPrivateKey(private_key_hex: string): WcKeypair;

/**
 * Recovers the signer address from a base64 EIP-191 personal_sign
 * signature, without needing a claimed address to check against.
 */
export function recoverPersonalSigner(message: string, signature_b64: string): string;

/**
 * Signs `message` (UTF-8 text) with EIP-191 personal_sign, returning a
 * base64 signature - the Base equivalent of trovoSDK.ts's `signHTTP` and
 * `signBase64Txn` (both now sign a message/digest the same way; see the
 * package doc and internal/middleware/security_checks.go on the backend).
 */
export function signPersonal(private_key_hex: string, message: string): string;

/**
 * Signs raw bytes (e.g. a pre-computed transaction digest) with EIP-191
 * personal_sign, returning a base64 signature.
 */
export function signPersonalBytes(private_key_hex: string, message: Uint8Array): string;

/**
 * Verifies a base64 EIP-191 personal_sign signature against `address`.
 */
export function verifyPersonal(address: string, message: string, signature_b64: string): boolean;

export type InitInput = RequestInfo | URL | Response | BufferSource | WebAssembly.Module;

export interface InitOutput {
    readonly memory: WebAssembly.Memory;
    readonly __wbg_get_wckeypair_address: (a: number) => [number, number];
    readonly __wbg_get_wckeypair_privateKeyHex: (a: number) => [number, number];
    readonly __wbg_set_wckeypair_address: (a: number, b: number, c: number) => void;
    readonly __wbg_set_wckeypair_privateKeyHex: (a: number, b: number, c: number) => void;
    readonly __wbg_wckeypair_free: (a: number, b: number) => void;
    readonly addressFromPrivateKey: (a: number, b: number) => [number, number, number, number];
    readonly generateKeypair: () => number;
    readonly generateMnemonic: () => [number, number, number, number];
    readonly keypairFromMnemonic: (a: number, b: number, c: number) => [number, number, number];
    readonly keypairFromPrivateKey: (a: number, b: number) => [number, number, number];
    readonly recoverPersonalSigner: (a: number, b: number, c: number, d: number) => [number, number, number, number];
    readonly signPersonal: (a: number, b: number, c: number, d: number) => [number, number, number, number];
    readonly signPersonalBytes: (a: number, b: number, c: number, d: number) => [number, number, number, number];
    readonly verifyPersonal: (a: number, b: number, c: number, d: number, e: number, f: number) => [number, number, number];
    readonly __wbindgen_exn_store: (a: number) => void;
    readonly __externref_table_alloc: () => number;
    readonly __wbindgen_externrefs: WebAssembly.Table;
    readonly __wbindgen_malloc: (a: number, b: number) => number;
    readonly __wbindgen_realloc: (a: number, b: number, c: number, d: number) => number;
    readonly __externref_table_dealloc: (a: number) => void;
    readonly __wbindgen_free: (a: number, b: number, c: number) => void;
    readonly __wbindgen_start: () => void;
}

export type SyncInitInput = BufferSource | WebAssembly.Module;

/**
 * Instantiates the given `module`, which can either be bytes or
 * a precompiled `WebAssembly.Module`.
 *
 * @param {{ module: SyncInitInput }} module - Passing `SyncInitInput` directly is deprecated.
 *
 * @returns {InitOutput}
 */
export function initSync(module: { module: SyncInitInput } | SyncInitInput): InitOutput;

/**
 * If `module_or_path` is {RequestInfo} or {URL}, makes a request and
 * for everything else, calls `WebAssembly.instantiate` directly.
 *
 * @param {{ module_or_path: InitInput | Promise<InitInput> }} module_or_path - Passing `InitInput` directly is deprecated.
 *
 * @returns {Promise<InitOutput>}
 */
export default function __wbg_init (module_or_path?: { module_or_path: InitInput | Promise<InitInput> } | InitInput | Promise<InitInput>): Promise<InitOutput>;
