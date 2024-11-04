import * as StellarSdk from '@stellar/stellar-sdk';
import { mnemonicToSeedSync } from 'bip39';
import {Buffer} from "buffer";

const createAccount = (): Account => {
    // create a completely new and unique pair of keys.
    const keyPair = StellarSdk.Keypair.random();   

    /// Returns the human readable account ID of this key pair.
    //return {"publicKey": keyPair.accountId, "secretKey": keyPair.secretSeed};
    return {publicKey:keyPair.publicKey(), secretKey:keyPair.secret()};
}

const signHTTP = (toSign: string, secretKey: string): string => {
    try{
        const keypair = StellarSdk.Keypair.fromSecret(secretKey);
        const encoder = new TextEncoder();
        const list = encoder.encode(toSign);
        const buffer = Buffer.from(list);

        const signedData = keypair.sign(buffer);
        const signedBase64Str = btoa(Array.from(signedData).toString());

        return signedBase64Str;
    }catch(error: any){
        console.log('error signing request', error);
        return '';
    }
}

const importAccount = (secretKey: string): string => {
    try{
        const keypair = StellarSdk.Keypair.fromSecret(secretKey);
        return keypair.publicKey();        
    }catch(error: any){
        console.log('error signing request', error);
        return '';
    }
}

const signBase64Txn = (secretKey: string, transactionXDR: string, networkPassphrase: string): string => {
    try{
        const keypair = StellarSdk.Keypair.fromSecret(secretKey);
        console.log('dkslfsd', keypair, secretKey);
        const txn = new StellarSdk.Transaction(transactionXDR, networkPassphrase);
        console.log('txn', txn);
        const bytes = txn.hash();
        console.log('bytes', bytes);
        
        const signedData = keypair.sign(bytes);
        console.log('bytes', bytes);
        const signedBase64Str = Buffer.from(signedData).toString('base64');
console.log('btoa result', signedBase64Str);
        return signedBase64Str;
    }catch(error: any){
        console.log('error signing request', error);
        return '';
    }
}

const parseSecretKey = (secretKey: string): Account => {
    const keypair = StellarSdk.Keypair.fromSecret(secretKey);

    return {publicKey: keypair.publicKey(), secretKey: keypair.secret()};
};

const getCredsFromPassPhrase = (passphrase: string): Account | null => {
    try {    
        const seed = mnemonicToSeedSync(passphrase);
        const keypair = StellarSdk.Keypair.fromRawEd25519Seed(Buffer.alloc(32, seed));

        return {publicKey: keypair.publicKey(), secretKey: keypair.secret()};
    } catch (error: any) {
        console.log(error);
        return null;        
    }
};

export {createAccount, signHTTP, importAccount, signBase64Txn, parseSecretKey, getCredsFromPassPhrase};
