import WalletDropdown from '../../../components/walletDropdown';
import { useEffect, useState } from 'react';
import { Wallet } from '../../../types/wallet';
import { RootState } from '../../../store/reduxStore';
import { useDispatch, useSelector } from 'react-redux';
import Header from '../../../components/header';
import Modal from '../../../components/modal';
import Button from '../../../components/button';
import { CuratedAsset } from '../../../types/curatedAsset';
import { truncatePublicKey } from '../../../utils/truncateValues';
import { showNotification, toggleLoader } from '../../../utils/showToaster';
import { Asset } from '../../../types/asset';
import { formatAmount } from '../../../utils/transactionUtils';
import ButtonSecondary from '../../../components/buttonSecondary';
import { Encryptor } from '../../../utils/encryptor';
import { signBase64Txn } from '../../../utils/trovoSDK';
import { ErrorResponse } from '../../../store/api/baseapi/axiosBaseQuery';
import {
  useAddAssetMutation,
  useRemoveAssetMutation,
} from '../../../store/api/walletApis';
import { useLazyGetUserQuery } from '../../../store/api/authApi';
import { deserializeUserData } from '../../../utils/deserializeAndStoreUserData';
import { User } from '../../../types/user';
import { setTempData, setUser } from '../../../store/authSlice';

interface AssetTile {
  isRemovable: boolean;
  balance: number;
  asset: CuratedAsset;
}
const filteredAssets: Record<string, AssetTile> = {};

const AssetsList = () => {
  const [secretKey, setSecretKey] = useState('');
  const appUser = useSelector((state: RootState) => state.auth.user!);
  let mutableWalletArray = [...appUser.userWallets];
  const [wallets] = useState(
    mutableWalletArray.sort((w) => (w.primaryWallet ? 0 : 1)),
  );
  const [activeWallet, setActiveWallet] = useState<Wallet>();
  const [filteredWallets, setFilteredWallets] = useState<AssetTile[]>();
  const tempData = useSelector((state: RootState) => state.auth.tempData);
  const dispatch = useDispatch();
  const [getUser, {}] = useLazyGetUserQuery();

  useEffect(() => {
    (async () => {
      const encryptor = new Encryptor();
      setSecretKey(await encryptor.getSecretKey(appUser));
    })();
  }, []);

  const importWallet = async () => {
    const importedPayload = {
      signer: appUser.publicKey,
      publicKey: appUser.publicKey,
      secretKey: secretKey,
      body: { userId: appUser.username, import: 1 },
    };
    const { data: importData, error: importError } =
      await getUser(importedPayload);

    var userData = importData as User;
    if (importData) {
      userData = deserializeUserData(importData);
    } else if (importError) {
      const err = importError as ErrorResponse;
      showNotification(
        'error',
        err.data.message ??
          'Sorry we could not complete the request. Please try again.',
      );
      return;
    }

    if (userData) {
      const wallet = userData.userWallets.find(
        (w) => w.alias == activeWallet?.alias,
      )!;

      const filteredList = (appUser?.curatedSwapList ?? []).filter(
        (asset: any) => asset.assetCode !== '' && asset.assetIssuer !== '',
      );

      filteredList.forEach((asset: any) => {
        filteredAssets[`${asset.assetCode!}|${asset.assetIssuer}`] = {
          isRemovable: false,
          balance: 0,
          asset: asset as CuratedAsset,
        };
      });

      const filteredList2 = (wallet.claimedAssets ?? []).filter(
        (asset: Asset) => asset.assetCode !== '' && asset.assetIssuer !== '',
      );

      filteredList2.forEach((asset: Asset) => {
        const assetKey = `${asset.assetCode!}|${asset.assetIssuer}`;

        if (!filteredAssets[assetKey]) {
          filteredAssets[assetKey] = {
            isRemovable: true,
            balance: asset.amount,
            asset: asset as any,
          };
        } else {
          filteredAssets[assetKey].isRemovable = true;
          filteredAssets[assetKey].balance = asset.amount;
        }
      });

      setFilteredWallets(Object.values(filteredAssets));
      setActiveWallet(wallet);

      const encryptor = new Encryptor();
      const base64EncryptedSecretKey = await encryptor.encryptData(
        secretKey,
        tempData.password,
        appUser.publicKey,
      );

      const passwordHash = await encryptor.createHash(userData.username);
      const encryptedPassword = await encryptor.encryptData(
        tempData.password,
        passwordHash,
        appUser.publicKey,
      );

      const user = {
        ...userData,
        password: encryptedPassword,
        isLoggedIn: true,
        currency: 'USD',
        secretKeys: appUser.secretKeys,
      };

      const hash = await encryptor.createHash(user.username);
      const base64EncryptedUserData = await encryptor.encryptData(
        JSON.stringify(user),
        hash,
        user.publicKey,
      );

      dispatch(
        setUser({
          key: hash,
          user,
          encryptedUser: base64EncryptedUserData,
        }),
      );

      dispatch(
        setTempData({
          ...tempData,
          secretKey: base64EncryptedSecretKey,
        }),
      );
    } else {
      showNotification('error', 'Something went wrong. Please try again.');
    }
  };

  return (
    <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
      <Header isHomeView />
      <div className="w-full px-3 md:px-5 md:pt-5 flex flex-col md:flex-row space-y-5 md:space-y-0 items-center justify-between">
        <div className="flex space-x-5 w-full md:w-auto items-center">
          <p className="font-montserratSemiBold text-lg xl:text-xl">
            Add/Remove Asset
          </p>
        </div>
      </div>
      <div className="w-full md:w-2/3 px-3 md:px-5 md:pt-5 flex flex-col items-center md:space-x-3 space-y-5">
        <form
          id="swap-asset-form"
          onSubmit={() => {}}
          className="w-full space-y-5"
        >
          <div className="w-full space-y-5">
            <div className="space-y-3 w-full">
              <WalletDropdown
                label={activeWallet?.alias ?? 'Select wallet'}
                options={[
                  ...wallets.map((w, index) => ({
                    text: w.alias,
                    value: w,
                    index,
                  })),
                ]}
                onSelect={(selectedItem) => {
                  console.log(selectedItem);
                  const wallet = selectedItem.value as Wallet;
                  setActiveWallet(wallet);

                  const filteredList = (appUser?.curatedSwapList ?? []).filter(
                    (asset: any) =>
                      asset.assetCode !== '' && asset.assetIssuer !== '',
                  );

                  filteredList.forEach((asset: any) => {
                    filteredAssets[`${asset.assetCode!}|${asset.assetIssuer}`] =
                      {
                        isRemovable: false,
                        balance: 0,
                        asset: asset as CuratedAsset,
                      };
                  });

                  const filteredList2 = (wallet.claimedAssets ?? []).filter(
                    (asset: Asset) =>
                      asset.assetCode !== '' && asset.assetIssuer !== '',
                  );

                  filteredList2.forEach((asset: Asset) => {
                    const assetKey = `${asset.assetCode!}|${asset.assetIssuer}`;

                    if (!filteredAssets[assetKey]) {
                      filteredAssets[assetKey] = {
                        isRemovable: true,
                        balance: asset.amount,
                        asset: asset as any,
                      };
                    } else {
                      filteredAssets[assetKey].isRemovable = true;
                      filteredAssets[assetKey].balance = asset.amount;
                    }
                  });

                  setFilteredWallets(Object.values(filteredAssets));
                }}
                key={`${activeWallet?.publicKey}`}
              />
            </div>
          </div>
          {!filteredWallets && (
            <div className="flex w-4/6 space-y-3 py-7 px-2 xl:px-5 xl:space-y-5 rounded-lg bg-primary-100">
              Select the wallet where you want to add or remove assets from.{' '}
              {secretKey} || {appUser.password}
            </div>
          )}
        </form>
        {filteredWallets && (
          <div className="w-full space-y-3">
            {filteredWallets.map((asset, index) => (
              <CuratedAssetTile
                key={index}
                assetTile={asset}
                wallet={activeWallet!}
                secretKey={secretKey}
                onImport={() => importWallet()}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

interface CuratedAssetTileProps {
  assetTile: AssetTile;
  secretKey: string;
  onImport: () => void;
  wallet: Wallet;
}

const CuratedAssetTile: React.FC<CuratedAssetTileProps> = ({
  assetTile,
  wallet,
  secretKey,
  onImport,
}) => {
  const [showAddAssetModal, setShowAddAssetModal] = useState(false);
  const [showRemoveAssetModal, setShowRemoveAssetModal] = useState(false);
  const [addAsset] = useAddAssetMutation();
  const [removeAsset] = useRemoveAssetMutation();
  const handleButtonClick = () => {
    // Handle asset add/remove logic here
    console.log(assetTile);
    assetTile.isRemovable
      ? setShowRemoveAssetModal(true)
      : setShowAddAssetModal(true);
  };

  return (
    <>
      <div className="rounded-lg bg-white shadow-md p-2.5 md:p-5 flex justify-between items-center space-x-4">
        <div className="flex items-center space-x-5">
          <div className="w-10 h-10 rounded-full bg-white flex items-center justify-center overflow-hidden flex-shrink-0">
            <img
              src={assetTile.asset.imageUrl || '/images/trovoLogo.png'}
              alt="asset"
              className="w-full h-full object-cover"
              onError={(e) => {
                (e.target as HTMLImageElement).src = '/images/trovoLogo.png';
              }}
            />
          </div>
          <div className="flex flex-col space-y-1">
            <p className="font-semibold text-xs md:text-sm text-gray-900">
              {assetTile.asset.assetCode}
            </p>
            <p className="text-xs text-gray-600">{assetTile.asset.assetName}</p>
            {assetTile.asset.assetClassId && (
              <p className="text-xs text-gray-600">
                {assetTile.asset.assetClass.assetClass}
              </p>
            )}
          </div>
        </div>
        <button
          onClick={handleButtonClick}
          className={`px-4 py-2 rounded-lg text-xs font-semibold text-white flex-shrink-0 ${
            assetTile.isRemovable
              ? 'bg-red-500 hover:bg-red-400'
              : 'bg-primary-800 hover:bg-blue-500'
          }`}
        >
          {assetTile.isRemovable ? 'Remove' : 'Add'}
        </button>
      </div>
      <Modal
        showModal={showRemoveAssetModal}
        onClose={() => {
          setShowRemoveAssetModal(false);
        }}
      >
        <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
          <p className="font-montserratSemiBold text-lg xl:text-xl">
            Remove {assetTile.asset.assetCode}
          </p>
          {assetTile.balance > 0 ? (
            <div className="flex w-4/6 space-y-3 py-7 px-2 xl:px-5 xl:space-y-5 rounded-lg bg-primary-100 flex flex-col items-center">
              <p>
                Your still have some [{assetTile.asset.assetCode}] balance [
                {formatAmount(assetTile.balance, '')}] available on your wallet
                [{wallet.alias}].
              </p>
              <p>
                Please transfer the remaining [{assetTile.asset.assetCode}]
                balance on your wallet to another wallet or burn it by sending
                it back to the issuer on the address below before you can
                proceed.
              </p>
            </div>
          ) : (
            <div className="flex w-4/6 space-y-3 py-7 px-2 xl:px-5 xl:space-y-5 rounded-lg bg-primary-100 flex flex-col items-center">
              {wallet.canInitiate ? (
                <>
                  <p>
                    Do you wish to remove this asset [
                    {assetTile.asset.assetCode}
                    ]?
                  </p>
                  <p>
                    This means that [{assetTile.asset.assetCode}] will no longer
                    appear on the list of assets on your wallet [{wallet.alias}
                    ].
                  </p>
                </>
              ) : (
                <p>
                  You do not have enough permission to initiate this transaction
                  on this wallet [{wallet.alias}].
                </p>
              )}
            </div>
          )}

          <div className="w-3/4">
            {assetTile.balance > 0 ? (
              <Button
                label="Close"
                onclick={() => setShowRemoveAssetModal(false)}
              />
            ) : (
              <div className="space-y-2">
                <Button
                  label="Remove"
                  disabled={!wallet.canInitiate}
                  additionalClasses="bg-red-500"
                  onclick={async () => {
                    const payload = {
                      signer: wallet.signer,
                      publicKey: wallet.publicKey,
                      secretKey: secretKey,
                      body: {
                        isSharedWallet: wallet.sharedAccessEnabled,
                        assetCode: assetTile.asset.assetCode,
                        assetIssuer: assetTile.asset.assetIssuer,
                      },
                    };

                    toggleLoader();
                    const res = await removeAsset(payload);

                    if ('data' in res) {
                      const body = {
                        isSharedWallet: wallet.sharedAccessEnabled,
                        ...res.data,
                        commit: res.data.signatureRequired == 1 ? 0 : 1,
                        transactionSignature:
                          res.data.signatureRequired == 1
                            ? signBase64Txn(
                                secretKey,
                                res.data.transaction,
                                res.data.networkPassPhrase,
                              )
                            : '',
                      };

                      const payload = {
                        signer: wallet.signer,
                        publicKey: wallet.publicKey,
                        secretKey: secretKey,
                        body,
                      };
                      const res2 = await removeAsset(payload);
                      if ('data' in res2) {
                        await onImport();
                        setShowRemoveAssetModal(false);
                        showNotification(
                          'success',
                          `[${assetTile.asset.assetCode}] removed successfully.`,
                        );
                      } else if ('error' in res2) {
                        const errorResponse = res2.error as ErrorResponse;
                        showNotification(
                          'error',
                          errorResponse.data.message ??
                            'Something went wrong. Please try again.',
                        );
                      }
                    } else if ('error' in res) {
                      const errorResponse = res.error as ErrorResponse;
                      showNotification(
                        'error',
                        errorResponse.data.message ??
                          'Something went wrong. Please try again.',
                      );
                    }
                    toggleLoader();
                  }}
                />
                <ButtonSecondary
                  label="Close"
                  onclick={() => setShowRemoveAssetModal(false)}
                />
              </div>
            )}
          </div>
        </div>
      </Modal>
      <Modal
        showModal={showAddAssetModal}
        onClose={() => {
          setShowAddAssetModal(false);
        }}
      >
        <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
          <p className="font-montserratSemiBold text-lg xl:text-xl">
            Add {assetTile.asset.assetCode}
          </p>
          <div className="flex w-4/6 space-y-3 py-7 px-2 xl:px-5 xl:space-y-5 rounded-lg bg-primary-100 flex flex-col items-center">
            <p>
              Do you wish to add this asset [{assetTile.asset.assetCode}] to
              your wallet [{wallet.alias}]
            </p>
            <p>
              This will make this asset appear on the list of assets on your
              wallet [{wallet.alias}] and will enable you to start transacting
              with it
            </p>
          </div>
          <div className="flex w-4/6 space-y-3 py-7 px-2 xl:px-5 xl:space-y-5 rounded-lg bg-primary-100 flex flex-col items-center">
            <p className="font-montserratSemiBold text-lg xl:text-xl">
              {assetTile.asset.assetCode} Token
            </p>
            <div className="flex flex-col items-center space-y-2">
              <img
                src={assetTile.asset.imageUrl}
                alt={`${assetTile.asset.assetCode} logo`}
                className="h-24 w-24 rounded-full"
              />
              <p>{assetTile.asset.website}</p>
            </div>
            <div className="flex flex-col items-center space-y-2">
              <p>{assetTile.asset.assetName}</p>
              <p className="text-center">{assetTile.asset.description}</p>
            </div>
            <div className="w-full flex flex-col items-center space-y-2">
              <p className="font-montserratSemiBold flex items-center">
                Issuer Public Key
              </p>
              <p className="flex items-center space-x-2">
                <span>
                  {truncatePublicKey(assetTile.asset.assetIssuer ?? '')}
                </span>
                <button
                  type="button"
                  onClick={() =>
                    navigator.clipboard
                      .writeText(assetTile.asset.assetIssuer ?? '')
                      .then(() => {
                        showNotification('info', 'Public key copied!');
                      })
                  }
                >
                  <img
                    src="/images/copy.png"
                    alt="copy"
                    className="text-primary-500"
                  />
                </button>
              </p>
            </div>
          </div>
          {!wallet.canInitiate && (
            <div className="flex w-4/6 space-y-3 py-7 px-2 xl:px-5 xl:space-y-5 rounded-lg bg-primary-100 flex flex-col items-center">
              <p className="text-red-500">
                You do not have enough permission to initiate this transaction
                on this wallet [{wallet.alias}].
              </p>
            </div>
          )}
          <div className="w-3/4">
            <Button
              disabled={!wallet.canInitiate}
              label="Add"
              onclick={async () => {
                const payload = {
                  signer: wallet.signer,
                  publicKey: wallet.publicKey,
                  secretKey: secretKey,
                  body: {
                    isSharedWallet: wallet.sharedAccessEnabled,
                    assetCode: assetTile.asset.assetCode,
                    assetIssuer: assetTile.asset.assetIssuer,
                  },
                };

                toggleLoader();
                const res = await addAsset(payload);

                if ('data' in res) {
                  console.log('response', res);
                  const body = {
                    isSharedWallet: wallet.sharedAccessEnabled,
                    ...res.data,
                    commit: res.data.signatureRequired == 1 ? 0 : 1,
                    transactionSignature:
                      res.data.signatureRequired == 1
                        ? signBase64Txn(
                            secretKey,
                            res.data.transaction,
                            res.data.networkPassPhrase,
                          )
                        : '',
                  };

                  const payload = {
                    signer: wallet.signer,
                    publicKey: wallet.publicKey,
                    secretKey: secretKey,
                    body,
                  };

                  const res2 = await addAsset(payload);
                  if ('data' in res2) {
                    console.log('response2', res2);
                    await onImport();
                    setShowAddAssetModal(false);
                    showNotification(
                      'success',
                      `[${assetTile.asset.assetCode}] added successfully.`,
                    );
                  } else if ('error' in res2) {
                    const errorResponse = res2.error as ErrorResponse;
                    showNotification(
                      'error',
                      errorResponse.data.message ??
                        'Something went wrong. Please try again.',
                    );
                  }
                } else if ('error' in res) {
                  const errorResponse = res.error as ErrorResponse;
                  showNotification(
                    'error',
                    errorResponse.data.message ??
                      'Something went wrong. Please try again.',
                  );
                }
                toggleLoader();
              }}
            />
          </div>
        </div>
      </Modal>
    </>
  );
};

export default AssetsList;
