import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useParams } from 'react-router-dom';
import { RootState } from '../../../store/reduxStore';
import { Encryptor } from '../../../utils/encryptor';
import { useLazyGetTokenizationDetailQuery } from '../../../store/api/tokenizationApis';
import { setActiveTokenizedAsset } from '../../../store/appStateSlice';
import { TokenizedAsset } from '../../../types/tokenizedAsset';

function extractTokenizedAssetFromResponse(
  responseData: unknown,
): TokenizedAsset | undefined {
  const source = responseData as
    | { data?: { data?: TokenizedAsset } }
    | { data?: TokenizedAsset }
    | TokenizedAsset
    | undefined;

  if (!source) {
    return undefined;
  }

  if ('data' in source && source.data) {
    const nested = source.data as { data?: TokenizedAsset } | TokenizedAsset;
    if (
      nested &&
      typeof nested === 'object' &&
      'data' in nested &&
      (nested as { data?: TokenizedAsset }).data
    ) {
      return (nested as { data?: TokenizedAsset }).data;
    }

    return nested as TokenizedAsset;
  }

  return source as TokenizedAsset;
}

/**
 * Ensures the active tokenized asset is available in the Redux store.
 *
 * It reads the `:id` route parameter and, when the store either has no active
 * tokenized asset or contains a different asset, fetches the record from the
 * API using the id and dispatches it to the store. This keeps the pages
 * resilient to a full page reload (where the in-memory store is reset).
 */
export function useActiveTokenizedAsset() {
  const { id } = useParams<{ id: string }>();
  const dispatch = useDispatch();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const activeTokenizedAsset = useSelector(
    (state: RootState) => state.appState.activeTokenizedAsset,
  ) as TokenizedAsset | undefined;
  const [getTokenizationDetail] = useLazyGetTokenizationDetailQuery();
  const [isLoading, setIsLoading] = useState(false);
  const [fetchError, setFetchError] = useState<string | null>(null);

  const assetId = id ?? activeTokenizedAsset?.id ?? '';

  useEffect(() => {
    if (!assetId) {
      return;
    }

    if (
      activeTokenizedAsset?.id &&
      String(activeTokenizedAsset.id) === String(assetId)
    ) {
      return;
    }

    let cancelled = false;

    const fetchAsset = async () => {
      setIsLoading(true);
      setFetchError(null);

      try {
        const encryptor = new Encryptor();
        const currentSecretKey = await encryptor.getSecretKey(appUser);

        const res = await getTokenizationDetail({
          signer: appUser.primarySigner,
          publicKey: appUser.publicKey,
          secretKey: currentSecretKey,
          tokenizedAssetId: assetId,
        });

        if (cancelled) {
          return;
        }

        if ('data' in res) {
          const detailData = extractTokenizedAssetFromResponse(res.data);
          if (detailData) {
            dispatch(setActiveTokenizedAsset(detailData));
          } else {
            setFetchError('Unable to load the tokenized asset.');
          }
        } else {
          setFetchError('Unable to load the tokenized asset.');
        }
      } catch {
        if (!cancelled) {
          setFetchError('Unable to load the tokenized asset.');
        }
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    };

    void fetchAsset();

    return () => {
      cancelled = true;
    };
  }, [
    assetId,
    activeTokenizedAsset?.id,
    appUser,
    getTokenizationDetail,
    dispatch,
  ]);

  return { activeTokenizedAsset, isLoading, fetchError, assetId };
}
