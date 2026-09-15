import { useEffect, useMemo, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import { TokenizationData } from '../../../types/tokenizationData';
import { Encryptor } from '../../../utils/encryptor';
import { useLazyFetchTokenizationDataQuery } from '../../../store/api/walletApis';
import { setTokenizationData } from '../../../store/appStateSlice';

/**
 * Ensures tokenization metadata is available in the Redux store.
 *
 * It reads the user's primary wallet and, when the store has no
 * `tokenizationData`, lazily fetches it from the API and dispatches it to the
 * store. This keeps pages resilient to a full page reload (where the in-memory
 * store is reset) and mirrors the behaviour of `useActiveTokenizedAsset`.
 */
export function useTokenizationData() {
  const dispatch = useDispatch();
  const appUser = useSelector((state: RootState) => state.auth.user as any);
  const tokenizationData = useSelector(
    (state: RootState) => state.appState.tokenizationData,
  ) as TokenizationData | undefined;

  const [fetchTokenizationData] = useLazyFetchTokenizationDataQuery();
  const [isLoading, setIsLoading] = useState(false);
  const [fetchError, setFetchError] = useState<string | null>(null);

  const walletList = appUser?.userWallets ?? [];
  const primaryWallet = useMemo(() => {
    return (
      walletList.find((wallet: any) => wallet.primaryWallet) ?? walletList[0]
    );
  }, [walletList]);

  useEffect(() => {
    if (tokenizationData) {
      return;
    }

    if (!appUser || !primaryWallet) {
      return;
    }

    let cancelled = false;

    const fetchData = async () => {
      setIsLoading(true);
      setFetchError(null);

      try {
        const encryptor = new Encryptor();
        const currentSecretKey = await encryptor.getSecretKey(appUser);

        const res = await fetchTokenizationData({
          signer: primaryWallet?.signer ?? '',
          publicKey: primaryWallet?.publicKey ?? '',
          secretKey: currentSecretKey,
          body: {},
        });

        if (cancelled) {
          return;
        }

        if ('data' in res && res.data) {
          dispatch(setTokenizationData(res.data as TokenizationData));
        } else {
          setFetchError('Unable to load tokenization metadata.');
        }
      } catch {
        if (!cancelled) {
          setFetchError('Unable to load tokenization metadata.');
        }
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    };

    void fetchData();

    return () => {
      cancelled = true;
    };
  }, [
    appUser,
    primaryWallet,
    tokenizationData,
    fetchTokenizationData,
    dispatch,
  ]);

  return { tokenizationData, isLoading, fetchError };
}
