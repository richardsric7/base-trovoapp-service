/* tslint:disable */
/* eslint-disable */
export const memory: WebAssembly.Memory;
export const __wbg_get_wckeypair_address: (a: number) => [number, number];
export const __wbg_get_wckeypair_privateKeyHex: (a: number) => [number, number];
export const __wbg_set_wckeypair_address: (a: number, b: number, c: number) => void;
export const __wbg_set_wckeypair_privateKeyHex: (a: number, b: number, c: number) => void;
export const __wbg_wckeypair_free: (a: number, b: number) => void;
export const addressFromPrivateKey: (a: number, b: number) => [number, number, number, number];
export const generateKeypair: () => number;
export const generateMnemonic: () => [number, number, number, number];
export const keypairFromMnemonic: (a: number, b: number, c: number) => [number, number, number];
export const keypairFromPrivateKey: (a: number, b: number) => [number, number, number];
export const recoverPersonalSigner: (a: number, b: number, c: number, d: number) => [number, number, number, number];
export const signPersonal: (a: number, b: number, c: number, d: number) => [number, number, number, number];
export const signPersonalBytes: (a: number, b: number, c: number, d: number) => [number, number, number, number];
export const verifyPersonal: (a: number, b: number, c: number, d: number, e: number, f: number) => [number, number, number];
export const __wbindgen_exn_store: (a: number) => void;
export const __externref_table_alloc: () => number;
export const __wbindgen_externrefs: WebAssembly.Table;
export const __wbindgen_malloc: (a: number, b: number) => number;
export const __wbindgen_realloc: (a: number, b: number, c: number, d: number) => number;
export const __externref_table_dealloc: (a: number) => void;
export const __wbindgen_free: (a: number, b: number, c: number) => void;
export const __wbindgen_start: () => void;
