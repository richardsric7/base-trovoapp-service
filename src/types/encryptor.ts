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
  }
  