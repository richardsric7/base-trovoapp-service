import { setUser, setTempUser } from "../store/authSlice";
import { USER_DETAILS } from "../store/constants";
import { store } from "../store/reduxStore";
import { User } from "../types/user";
import { getStorage } from "./storage";

export class Encryptor {  
    async deriveKeyFromPassword(password: string, mySalt: string): Promise<CryptoKey> {
      const salt = new TextEncoder().encode(mySalt);
      const encodedPassword = new TextEncoder().encode(password);
      const key = await window.crypto.subtle.importKey(
        'raw',
        encodedPassword,
        { name: 'PBKDF2' },
        false,
        ['deriveBits', 'deriveKey']
      );
      return window.crypto.subtle.deriveKey(
        {
          name: 'PBKDF2',
          salt,
          iterations: 10000,
          hash: 'SHA-256'
        },
        key,
        { name: 'AES-GCM', length: 256 },
        true,
        ['encrypt', 'decrypt']
      );
    }
  
    async encryptData(data: string, password: string, primaryWalletPublicKey: string): Promise<string> {
      const key = await this.deriveKeyFromPassword(password, primaryWalletPublicKey);
      const encodedData = new TextEncoder().encode(data);
      const iv = window.crypto.getRandomValues(new Uint8Array(12));
      const encryptedData =  await window.crypto.subtle.encrypt(
        { name: 'AES-GCM', iv },
        key,
        encodedData
      );

      const ivBase64 = btoa(String.fromCharCode.apply(null, Array.from(iv)));
      const encryptDataB64 = btoa(
        String.fromCharCode.apply(
          null,
          Array.from(new Uint8Array(encryptedData)),
        ),
      );

      return `${ivBase64}|${encryptDataB64}`;
    }
  
    async decryptData(base64EncryptedData: string, password: string, primaryWalletPublicKey: string): Promise<string> {
      const splitB64String = base64EncryptedData.split('|');
      const encryptedDataArray = Uint8Array.from(
        atob(splitB64String[1]),
        (c) => c.charCodeAt(0),
      );
      const iv = Uint8Array.from(atob(splitB64String[0]), (c) =>
        c.charCodeAt(0),
      );
      const encryptedDataUint = new Uint8Array(encryptedDataArray);

      const key = await this.deriveKeyFromPassword(password, primaryWalletPublicKey);
      const decrypted = await window.crypto.subtle.decrypt(
        { name: 'AES-GCM', iv },
        key,
        encryptedDataUint
      );
      return new TextDecoder().decode(decrypted);
    }

    async createHash(plainText: string): Promise<string> {
      const utf8 = new TextEncoder().encode(plainText);
      const hashBuffer = await crypto.subtle.digest('SHA-256', utf8);
      const hashArray = Array.from(new Uint8Array(hashBuffer));
      const hashHex = hashArray
        .map((bytes) => bytes.toString(16).padStart(2, '0'))
        .join('');
      return hashHex;      
    }

    async encryptUserData(user: User): Promise<void> {
      const hash = await this.createHash(user.username);
      const base64EncryptedUserData = await this.encryptData(
        JSON.stringify(user),
        hash,
        user.publicKey,
      );
      store.dispatch(
        setUser({
          key: hash,
          user,
          encryptedUser: base64EncryptedUserData,
        }),
      );
    }

    async decryptUserData(): Promise<User|undefined> {
      const storedInfo = getStorage(USER_DETAILS);
      if (!storedInfo) {
        return undefined;
      }
      const result = await this.decryptData(
        storedInfo.__slw31H408,
        storedInfo.__39deR7sx4,
        storedInfo.__i34dcY9Mn,
      );
      const user = JSON.parse(result) as User;
      store.dispatch(
        setTempUser({
          ...user,
        }),
      );
      return user;
    }
  }
  