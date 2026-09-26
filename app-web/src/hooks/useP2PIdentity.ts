import { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../store/reduxStore';
import { Encryptor } from '../utils/encryptor';
import { P2PCreds } from '../store/api/p2pApis';

// useP2PIdentity resolves the signed-request credentials every P2P call
// needs (Plan Section 92's finding: any new P2P endpoint must go through
// app-web's existing signed-request axiosBaseQuery, same as walletApis.ts),
// plus the caller's own username - the client has no raw user-id, only
// username (mirrors app-mobile's identical finding), so every P2P
// role-check compares by username.
export function useP2PIdentity() {
  const appUser = useSelector((state: RootState) => state.auth.user);
  const [secretKey, setSecretKey] = useState('');
  const [ready, setReady] = useState(false);

  useEffect(() => {
    let cancelled = false;
    if (!appUser) return;
    const encryptor = new Encryptor();
    encryptor.getSecretKey(appUser).then((key: string) => {
      if (!cancelled) {
        setSecretKey(key);
        setReady(true);
      }
    });
    return () => {
      cancelled = true;
    };
  }, [appUser]);

  const creds: P2PCreds = {
    signer: appUser?.primarySigner ?? '',
    address: appUser?.address ?? '',
    secretKey,
  };

  return {
    creds,
    username: appUser?.username ?? '',
    ready: ready && !!appUser,
  };
}
