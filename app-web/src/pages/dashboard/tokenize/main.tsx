import { useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { FaCirclePlus } from 'react-icons/fa6';
import Header from '../../../components/header';
import ButtonSecondary from '../../../components/buttonSecondary';
import { RootState } from '../../../store/reduxStore';
import { useFetchTokenizationDataQuery } from '../../../store/api/walletApis';
import { useFetchTokenizationAssetsQuery } from '../../../store/api/tokenizationApis';
import { Encryptor } from '../../../utils/encryptor';
import { showNotification } from '../../../utils/showToaster';
import {
  setActiveTokenizedAsset,
  setTokenizationData,
} from '../../../store/appStateSlice';
import { TokenizedAsset } from '../../../types/tokenizedAsset';
import { TokenizationData } from '../../../types/tokenizationData';
import TokenizationAssetListItem from '../../../components/tokenizationAssetListItem';

export const Tokenize = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const appUser = useSelector((state: RootState) => state.auth.user as any);
  const [secretKey, setSecretKey] = useState('');
  const walletList = appUser?.userWallets ?? [];
  const primaryWallet = useMemo(() => {
    return (
      walletList.find((wallet: any) => wallet.primaryWallet) ?? walletList[0]
    );
  }, [walletList]);

  useEffect(() => {
    if (!appUser) {
      return;
    }

    (async () => {
      try {
        const encryptor = new Encryptor();
        const currentSecretKey = await encryptor.getSecretKey(appUser);
        setSecretKey(currentSecretKey);
      } catch (error) {
        showNotification('error', 'Unable to load wallet credentials.');
      }
    })();
  }, [appUser]);

  const {
    data: tokenizationDataResponse,
    isLoading: isLoadingTokenizationData,
    refetch: refetchTokenizationData,
  } = useFetchTokenizationDataQuery(
    {
      signer: primaryWallet?.signer ?? '',
      address: primaryWallet?.address ?? '',
      secretKey,
      body: {},
    },
    {
      skip: !secretKey || !primaryWallet,
    },
  );

  useEffect(() => {
    if (tokenizationDataResponse) {
      dispatch(
        setTokenizationData(tokenizationDataResponse as TokenizationData),
      );
    }
  }, [dispatch, tokenizationDataResponse]);

  const {
    data: tokenizedAssetsResponse,
    isLoading: isLoadingAssets,
    isFetching: isFetchingAssets,
    refetch: refetchTokenizedAssets,
    error: assetsError,
  } = useFetchTokenizationAssetsQuery(
    {
      signer: primaryWallet?.signer ?? '',
      address: primaryWallet?.address ?? '',
      secretKey,
      body: {
        status: 0,
      },
    },
    {
      skip: !secretKey || !primaryWallet,
    },
  );

  const records = useMemo(() => {
    const rawRecords = (tokenizedAssetsResponse as any)?.records ?? [];
    return rawRecords as TokenizedAsset[];
  }, [tokenizedAssetsResponse]);

  const tokenizationConfig = useMemo(() => {
    return (tokenizationDataResponse as any)?.countryConfigs?.find(
      (config: any) => {
        return config.countryCode?.toLowerCase() === 'ng';
      },
    );
  }, [tokenizationDataResponse]);

  const tokenizationApplicationFee = Number(
    tokenizationConfig?.tokenizationApplicationFee ?? 0,
  );
  const tokenizationApplicationFeeAsset = (
    tokenizationConfig?.tokenizationApplicationFeeAsset ?? ''
  )
    .split(':')[0]
    .toUpperCase();

  const hasEnoughTrov = useMemo(() => {
    if (
      !tokenizationApplicationFeeAsset ||
      !primaryWallet?.claimedAssets?.length
    ) {
      return false;
    }

    const matchingAsset = primaryWallet.claimedAssets.find((asset: any) => {
      return asset.assetCode?.toUpperCase() === tokenizationApplicationFeeAsset;
    });

    if (!matchingAsset) {
      return false;
    }

    const balance = Number(matchingAsset.amount ?? 0);
    const usdPrice = Number(matchingAsset.usdPrice ?? 0);

    if (usdPrice > 0) {
      return balance >= tokenizationApplicationFee / usdPrice;
    }

    return balance >= tokenizationApplicationFee;
  }, [
    primaryWallet,
    tokenizationApplicationFee,
    tokenizationApplicationFeeAsset,
  ]);

  const canCreateNewTokenization = useMemo(() => {
    return !records.some((asset) => asset.assetTokenizationStatus === 0);
  }, [records]);

  const handleStartTokenization = () => {
    dispatch(setActiveTokenizedAsset({}));

    if (!tokenizationDataResponse) {
      showNotification(
        'error',
        'Cannot initiate this process at the moment. Please check your network, refresh this view and try again.',
      );
      return;
    }

    if (!canCreateNewTokenization) {
      showNotification(
        'error',
        'You must complete the active tokenization process before starting a new one.',
      );
      return;
    }

    if (appUser?.kycVerified === 0) {
      showNotification(
        'error',
        'Your account must be KYC verified before you can start a new tokenization application.',
      );
      return;
    }

    if (!hasEnoughTrov) {
      showNotification(
        'error',
        `You must have at least ${tokenizationApplicationFee} ${tokenizationApplicationFeeAsset} tokens in your wallet to begin a new tokenization process.`,
      );
      return;
    }

    navigate('/dashboard/tokenize/apply');
  };

  const handleAssetSelect = (asset: TokenizedAsset) => {
    dispatch(setActiveTokenizedAsset(asset));

    if (asset.assetTokenizationStatus === 0) {
      navigate(`/dashboard/tokenize/apply/${asset.id}`);
      return;
    }

    if (
      asset.assetTokenizationStatus === 1 ||
      asset.assetTokenizationStatus === 2 ||
      asset.assetTokenizationStatus === 3
    ) {
      navigate(`/dashboard/tokenize/confirm-details/${asset.id}`);
      return;
    }

    navigate(`/dashboard/tokenize/asset-dashboard/${asset.id}`);
  };

  return (
    <div>
      <Header />
      <div className="max-w-full p-5">
        <div className="rounded-2xl space-y-8 bg-primary-100 py-6 px-10">
          <div className="rounded-2xl bg-[#004988] bg-[url('/images/topographicIcon.png')] bg-cover bg-center p-6 text-center text-white shadow-md">
            <h2 className="mb-2 text-xl font-semibold">
              Unlock Liquidity From Fixed Assets
            </h2>
            <p className="mb-4 text-sm font-light">
              Digitally fragment any asset of value to unlock capital from
              multiple contributors. Use the power of tokenization to give more
              people access to investment opportunities previously unavailable
              to them. Tokenize land, real estate, infrastructure projects,
              commodities, natural resources and other qualifying assets.
            </p>
            <button
              type="button"
              onClick={handleStartTokenization}
              className="mt-2 inline-flex items-center gap-2 rounded-lg bg-[#FFFFFF33] px-4 py-2 font-semibold text-[#FFFFFF] transition hover:bg-primary-600"
            >
              <FaCirclePlus />
              Tokenize Asset
            </button>
          </div>
          <div className="bg-primary-100 rounded-2xl p-4 shadow-sm">
            {isLoadingTokenizationData ||
            (isLoadingAssets && !tokenizedAssetsResponse) ? (
              <div className="py-8 text-center text-primary-800">
                Loading tokenization data…
              </div>
            ) : assetsError ? (
              <div className="space-y-3 py-8 text-center">
                <p className="text-primary-800">
                  We could not load your tokenized assets right now.
                </p>
                <ButtonSecondary
                  label="Retry"
                  onclick={() => {
                    void refetchTokenizedAssets();
                    void refetchTokenizationData();
                  }}
                />
              </div>
            ) : records.length === 0 ? (
              <div className="rounded-xl bg-[#EDF3F9] p-8 text-center">
                <div className="mb-6 flex justify-center">
                  <img
                    src="/images/emptyBox.png"
                    alt="Tokenize Graphic"
                    className="h-185 w-397 object-contain"
                  />
                </div>

                <p className="mx-auto mb-4 max-w-2xl text-sm text-[#0E1F51]">
                  Tokenized assets that are for public offering are securities
                  and MUST be approved by the Securities and Exchange Commission
                  (SEC).
                </p>
                <p className="mx-auto mb-4 max-w-2xl text-sm text-[#0E1F51]">
                  Only assets approved by the SEC can be tokenized and offered
                  to the public on the Trovo tokenization platform.
                </p>
                <p className="mx-auto mb-6 max-w-2xl text-sm text-[#0E1F51]">
                  Asset issued by private entities and are for private offering
                  are not under the purview of the SEC.
                </p>

                <p className="mb-4 text-sm font-medium text-[#0E1F51]">
                  <span className="font-semibold">
                    Are you an asset owner and want to tokenize your asset to
                    attract capital from around the world?
                  </span>
                </p>

                <div className="flex flex-col flex-wrap items-center justify-center gap-3 md:flex-row">
                  <button
                    type="button"
                    onClick={handleStartTokenization}
                    className="inline-flex items-center gap-2 rounded-lg bg-[#004988] px-5 py-2 text-white transition hover:bg-[#00285B]"
                  >
                    <FaCirclePlus size={20} />
                    Proceed to Tokenize Asset
                  </button>
                </div>
              </div>
            ) : (
              <div className="bg-primary-100">
                <div className="bg-primary-100 mb-4 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                  <div>
                    <p className="font-montserratSemiBold text-lg text-primary-800">
                      My Tokenized Assets
                    </p>
                  </div>
                  {/* <div className="flex flex-wrap gap-3">
                      <Dropdown
                        key={filterMode}
                        label="Filter by"
                        defaultValue={
                          filterModeOptions.find(
                            (option) => option.value === filterMode,
                          ) ?? filterModeOptions[0]
                        }
                        options={filterModeOptions}
                        onSelect={(item) => {
                          setFilterMode(item.value as FilterMode);
                          setFilterLabel(item.text);
                          setShowFilterModal(false);
                        }}
                      />
                      <Button
                        label={filterQuery ? 'Edit filter' : 'Add filter'}
                        onclick={() => openFilterModal(filterMode)}
                      />
                      {filterQuery ? (
                        <ButtonSecondary label="Clear" onclick={resetFilters} />
                      ) : null}
                  </div> */}
                </div>
                <div className="space-y-4">
                  {isFetchingAssets ? (
                    <p className="text-sm text-gray-500">Updating results…</p>
                  ) : null}
                  {records.map((asset, index) => (
                    <TokenizationAssetListItem
                      key={`${asset.id ?? asset.assetCode ?? index}`}
                      asset={asset}
                      onclick={() => handleAssetSelect(asset)}
                    />
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
