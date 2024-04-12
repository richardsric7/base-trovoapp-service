import * as StellarSdk from '@stellar/stellar-sdk';
import {Buffer} from 'buffer';

const createAccount = (): Account => {
    // create a completely new and unique pair of keys.
    const keyPair = StellarSdk.Keypair.random();   

    /// Returns the human readable account ID of this key pair.
    //return {"publicKey": keyPair.accountId, "secretKey": keyPair.secretSeed};
    console.log('key pair', keyPair.publicKey(), keyPair.secret());
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
        const txn = new StellarSdk.Transaction(transactionXDR, networkPassphrase);
        const bytes = txn.hash();

        const signedData = keypair.sign(bytes);
        const signedBase64Str = btoa(Array.from(signedData).toString());

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

export {createAccount, signHTTP, importAccount, signBase64Txn, parseSecretKey};