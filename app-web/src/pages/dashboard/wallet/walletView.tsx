import Button from '../../../components/button';
import Header from '../../../components/header';
import Tabs from '../../../components/tabs';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useEffect, useRef, useState } from 'react';
import WalletCard from '../../../components/walletCard';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import React from 'react';
import {
  formatToDecimal,
  canInitiate,
  totalWalletBalanceInCurrency,
  getBytesLength,
} from '../../../utils/utilities';
import { Wallet } from '../../../types/wallet';
import {
  TEDropdown,
  TEDropdownItem,
  TEDropdownMenu,
  TEDropdownToggle,
  TERipple,
} from 'tw-elements-react';
import { Asset } from '../../../types/asset';
import AssetItem from '../../../components/assetItem';
import WalletOperations from '../../../components/walletOperations';
import Modal from '../../../components/modal';
import TextInput from '../../../components/textInput';
import ButtonSecondary from '../../../components/buttonSecondary';
import { showNotification, toggleLoader } from '../../../utils/showToaster';
import { Encryptor } from '../../../utils/encryptor';
import {
  createAccount,
  importAccount,
  signBase64Txn,
} from '../../../utils/trovoSDK';
import { ErrorResponse } from '../../../store/api/baseapi/axiosBaseQuery';
import { useAddSubwalletMutation } from '../../../store/api/walletApis';
import { setTempData, setUser } from '../../../store/authSlice';
import { User } from '../../../types/user';
import { useLazyGetUserQuery } from '../../../store/api/authApi';
import { deserializeUserData } from '../../../utils/deserializeAndStoreUserData';

export default function WalletView() {
  const ref = useRef<HTMLDivElement>(null);
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const isSharedRel = searchParams.get('rel') === 'shared';
  const walletQueryParam = searchParams.get('wallet');

  let mutableWalletArray = [...appUser.userWallets];
  const [wallets, setWallets] = useState(
    mutableWalletArray.sort((w) => (w.primaryWallet ? 0 : 1)),
  );
  const [getUser, {}] = useLazyGetUserQuery();
  const dispatch = useDispatch();
  const tempData = useSelector((state: RootState) => state.auth.tempData);
  const [addSubwallet] = useAddSubwalletMutation();
  const [secretKey, setSecretKey] = useState('');
  const [activeWalletIndex, setActiveWalletIndex] = useState(0);
  const [activeWallet, setActiveWallet] = useState<Wallet>(wallets[0]);
  const [assetFilterMode, setAssetFilterMode] = useState(0);
  const [walletActionMode, setWalletActionMode] = useState(1);
  const itemRefs = useRef<HTMLDivElement[]>([]);
  const fiatRates = useSelector((state: RootState) => state.cache.fiatRates);
  const [selectedAsset, setSelectedAsset] = useState<Asset>(
    activeWallet?.claimedAssets.find(
      (a) => a.assetCode === '' && a.contractAddress === '',
    )!,
  );
  const [passwordErr, setPasswordErr] = useState('');
  const [password, setPassword] = useState('');
  const [showAddSubwalletModal, setShowAddSubwalletModal] = useState(false);
  const [showConfirmView, setShowConfirmView] = useState(false);
  const [showAuthorizeView, setShowAuthorizeView] = useState(false);
  const [showBackupModal, setShowBackupModal] = useState(false);
  const [currentTabIndex, setCurrentTabIndex] = useState(1);
  const [gas, setGas] = useState(
    activeWallet?.claimedAssets.find(
      (a) => a.assetCode === '' && a.contractAddress === '',
    )?.amount!,
  );

  type AddSubwalletType = {
    tag: string;
    description: string;
    newAddress: string;
    newSecretKey: string;
    isImport: boolean;
    transactionData: any;
  };

  const [formData, setFormData] = useState<AddSubwalletType>({
    tag: '',
    description: '',
    newAddress: '',
    newSecretKey: '',
    isImport: false,
    transactionData: undefined,
  });
  const [errorObj, setErrorObj] = useState({
    tag: '',
    description: '',
    newSecretKey: '',
  });

  const resetForm = () => {
    setErrorObj({
      tag: '',
      description: '',
      newSecretKey: '',
    });
    setFormData({
      tag: '',
      description: '',
      newAddress: '',
      newSecretKey: '',
      isImport: false,
      transactionData: undefined,
    });
  };

  const isValidPassword = async () => {
    try {
      const encryptor = new Encryptor();
      const key = await encryptor.decryptData(
        appUser?.secretKeys[0],
        password,
        appUser?.primarySigner,
      );
      setSecretKey(key);
      return true;
    } catch (error: any) {
      return false;
    }
  };

  useEffect(() => {
    (async () => {
      const encryptor = new Encryptor();
      setSecretKey(await encryptor.getSecretKey(appUser));
    })();
  }, []);

  useEffect(() => {
    console.log('using effect here!');
    let filtered = mutableWalletArray
      .filter((w) => {
        if (currentTabIndex === 1) {
          return w;
        } else if (
          currentTabIndex === 2 &&
          (w.owner === appUser.username ||
            !w.sharedAccessEnabled ||
            w.walletThreshold === 2)
        ) {
          return w;
        } else if (currentTabIndex === 3 && w.sharedAccessEnabled) {
          return w;
        }
      })
      .sort((w) => (w.primaryWallet ? 0 : 1));

    if (isSharedRel && walletQueryParam) {
      filtered = filtered.filter((w) => w.address === walletQueryParam);
    }

    setWallets(filtered);

    let targetIndex = 0;
    if (isSharedRel && walletQueryParam) {
      const idx = filtered.findIndex((w) => w.address === walletQueryParam);
      if (idx !== -1) {
        targetIndex = idx;
      }
    }

    setActiveWalletIndex(targetIndex);
    if (filtered[targetIndex]) {
      setActiveWallet(filtered[targetIndex]);
      const nativeAsset = filtered[targetIndex].claimedAssets.find(
        (a) => a.assetCode === '' && a.contractAddress === '',
      );
      if (nativeAsset) {
        setSelectedAsset(nativeAsset);
        setGas(nativeAsset.amount);
      }
    }

    setTimeout(() => {
      const ref = itemRefs.current[targetIndex];
      if (ref) {
        // Change background color
        ref.scrollIntoView({
          behavior: 'smooth',
          block: 'nearest',
          inline: 'center',
        });
      }
    }, 100);
  }, [currentTabIndex, isSharedRel, walletQueryParam]);

  useEffect(() => {
    setGas(
      activeWallet?.claimedAssets.find(
        (a) => a.assetCode === '' && a.contractAddress === '',
      )?.amount!,
    );
  }, [activeWallet, walletActionMode]);

  // Ensure `itemRefs` has refs for each item on each render
  if (itemRefs.current.length !== wallets.length) {
    // Adjust the array size based on items length
    itemRefs.current = Array(wallets.length)
      .fill(null)
      .map(
        (_, i) =>
          itemRefs.current[i] || React.createRef<HTMLDivElement>().current,
      );
  }

  const userWallets = wallets.map((wallet, index) => (
    // const walletRef = useRef<HTMLDivElement>(null);
    <WalletCard
      localCurrencyBalance={formatToDecimal(
        totalWalletBalanceInCurrency(
          wallet,
          (fiatRates as Record<string, number>)[
            appUser.currency.toUpperCase()
          ] || 0,
        ),
      )}
      usdBalance={formatToDecimal(
        totalWalletBalanceInCurrency(
          wallet,
          (fiatRates as { USD: number }).USD,
        ),
      )}
      currency={appUser.currency.toUpperCase()}
      alias={wallet.alias}
      isSharedAccess={wallet.sharedAccessEnabled}
      walletType={wallet.walletType!}
      key={index}
      ref={(el) => (itemRefs.current[index] = el!)}
    />
  ));

  const assetClassOptions = [
    {
      text: 'Asset Tokens',
      value: 0,
    },
    {
      text: 'Other Tokens',
      value: 1,
    },
  ];

  const dropdownItems = assetClassOptions.map((item) => (
    <TEDropdownItem
      key={`${new Date().getTime()}${item.text.replace(' ', '')}`}
      id={`${new Date().getTime()}${item.text.replace(' ', '')}`}
      className="bg-primary-200"
    >
      <button
        type="button"
        className="block w-full cursor-pointer hover:bg-primary-200 bg-primary-100 whitespace-nowrap px-8 py-2 text-sm text-left font-normal pointer-events-auto active:text-primary-800 focus:hover:bg-primary-200 focus:text-primary-800 focus:outline-none active:no-underline"
        onClick={(e: React.MouseEvent<HTMLButtonElement>) => {
          e.preventDefault;
          setAssetFilterMode(item.value);
        }}
      >
        {item.text}
      </button>
    </TEDropdownItem>
  ));

  const isValidForm = (value?: typeof formData): boolean => {
    let newObj = errorObj;
    var isValid = true;
    const pattern = /^[a-zA-Z0-9_]*$/;

    if (!value?.tag || value.tag.trim().length === 0) {
      newObj = { ...newObj, tag: 'Enter wallet tag' };
      isValid = false;
    } else if (
      value?.tag == null ||
      !pattern.test(value.tag.trim().replaceAll(' ', ''))
    ) {
      newObj = { ...newObj, tag: 'Invalid tag name' };
      isValid = false;
    } else {
      newObj = { ...newObj, tag: '' };
    }

    if (!value?.description || value.description.trim().length === 0) {
      newObj = { ...newObj, description: 'Enter wallet description' };
      isValid = false;
    } else if (getBytesLength(value.description.trim()) > 100) {
      newObj = {
        ...newObj,
        description:
          'Wallet description must be less than 100 characters in length',
      };
      isValid = false;
    } else {
      newObj = { ...newObj, description: '' };
    }

    if (
      value?.isImport &&
      (!value?.newSecretKey || value.newSecretKey.trim().length === 0)
    ) {
      newObj = { ...newObj, newSecretKey: 'Enter wallet secret key' };
      isValid = false;
    } else {
      newObj = { ...newObj, newSecretKey: '' };
    }

    setErrorObj(newObj);
    console.log('dalksd', errorObj);
    return isValid;
  };

  const validateAndConfirm = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!isValidForm(formData)) {
      return;
    }

    if (formData.isImport) {
      var address = importAccount(formData.newSecretKey);
      if (address.length == 0) {
        setErrorObj({ ...errorObj, newSecretKey: 'Invalid secret key' });
        return;
      }

      setFormData({
        ...formData,
        newAddress: address,
      });
    } else {
      // generate keypair for the new subwallet
      var ac = createAccount();
      setFormData({
        ...formData,
        newAddress: ac.address,
        newSecretKey: ac.secretKey,
      });
    }

    setShowConfirmView(true);
  };

  const handleSend = async () => {
    const payload = {
      signer: activeWallet.signer,
      address: activeWallet.address,
      secretKey: secretKey,
      body: {
        publickey: formData.newAddress,
        walletTag: formData.tag,
        WalletDescription: formData.description,
        walletType: 0,
        linkedWalletAddress: '',
      },
    };
    toggleLoader();
    const res = await addSubwallet(payload);
    if ('data' in res) {
      console.log('response', res);
      setFormData({
        ...formData,
        transactionData: res.data,
      });
      setShowAuthorizeView(true);
    } else if ('error' in res) {
      const errorResponse = res.error as ErrorResponse;
      showNotification(
        'error',
        errorResponse.data.message ?? 'Something went wrong. Please try again.',
      );
    }
    toggleLoader();
  };

  return (
    <>
      <div className="flex text-primary-800 text-sm overflow-x-hidden md:text-md flex-col space-y-5 p-5">
        <Header isHomeView />
        <div className="flex md:h-screen w-full items-center justify-center">
          <div className="h-full w-full md:p-3">
            <div className="flex md:pt-5 p-5 rounded-lg flex-col items-center bg-primary-100 max-w-5xl min-h-full">
              {isSharedRel ? (
                <div className="w-full flex justify-start mb-4">
                  <button
                    onClick={() => navigate('/dashboard/shared-access')}
                    className="flex items-center gap-2 text-primary-800 font-montserratSemiBold hover:text-primary-700 text-lg"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentColor">
                      <path d="m313-440 224 224-57 57-320-320 320-320 57 57-224 224h527v80H313Z"/>
                    </svg>
                    <span>Back</span>
                  </button>
                </div>
              ) : (
                <div className="w-full flex justify-between space-x-5">
                  <Tabs
                    tabList={['All Wallets', 'My Wallets', 'Shared Wallets']}
                    onTabChanged={(index) => {
                      setCurrentTabIndex(index);
                    }}
                  />
                  <div className="w-1/4">
                    <Button
                      label="Add Subwallet"
                      onclick={() => setShowAddSubwalletModal(true)}
                    />
                  </div>
                </div>
              )}
              <div className="relative max-w-full">
                <div
                  className="flex overflow-x-hidden whitespace-nowrap w-full space-x-4 "
                  ref={ref}
                >
                  {userWallets.length > 0 && userWallets}
                </div>
                {activeWalletIndex > 0 && (
                  <button
                    onClick={() => {
                      activeWalletIndex > 0 &&
                        setActiveWalletIndex(activeWalletIndex - 1);
                      setActiveWallet(wallets[activeWalletIndex - 1]);
                      if (!canInitiate(wallets[activeWalletIndex - 1])) {
                        setWalletActionMode(1);
                      }
                      setSelectedAsset(
                        activeWallet?.claimedAssets.find(
                          (a) => a.assetCode === '' && a.contractAddress === '',
                        )!,
                      );
                      const ref = itemRefs.current[activeWalletIndex - 1];
                      if (ref) {
                        // Change background color
                        ref.scrollIntoView({
                          behavior: 'smooth',
                          block: 'nearest',
                          inline: 'center',
                        });
                      }
                    }}
                    className="absolute ml-2 text-3xl top-1/2 left-0 transform -translate-y-1/2 bg-gray-200 w-10  h-10 rounded-full hover:bg-gray-300"
                  >
                    ←
                  </button>
                )}

                {activeWalletIndex + 1 <= itemRefs.current.length - 1 && (
                  <button
                    onClick={() => {
                      activeWalletIndex <= itemRefs.current.length &&
                        setActiveWalletIndex(activeWalletIndex + 1);
                      setActiveWallet(wallets[activeWalletIndex + 1]);
                      setSelectedAsset(
                        activeWallet?.claimedAssets.find(
                          (a) => a.assetCode === '' && a.contractAddress === '',
                        )!,
                      );

                      if (!canInitiate(wallets[activeWalletIndex + 1])) {
                        setWalletActionMode(1);
                      }
                      // resetForm();

                      const ref = itemRefs.current[activeWalletIndex + 1];
                      if (ref) {
                        // Change background color
                        ref.scrollIntoView({
                          behavior: 'smooth',
                          block: 'nearest',
                          inline: 'center',
                        });
                      }
                    }}
                    className="absolute mr-2 text-3xl top-1/2 right-0 transform -translate-y-1/2 bg-gray-200 w-10  h-10 rounded-full hover:bg-gray-300"
                  >
                    →
                  </button>
                )}
              </div>
              <div className="w-full ring-2 ring-primary-600 rounded-lg p-3">
                <div className="flex justify-between w-full">
                  <div className="flex space-x-2 items-center">
                    <img
                      src="/images/gas.png"
                      alt="gas"
                      width="30"
                      height="30"
                    />
                    <span>Gas</span>
                  </div>
                  <img
                    src="/images/export.png"
                    alt="export"
                    width="30"
                    height="30"
                  />
                </div>
                <p>{formatToDecimal(gas)} ETH</p>
              </div>
              <div className="w-full p-3 mt-7">
                <div className="flex justify-between w-full">
                  <p className="font-bold font-montserratSemiBold text-lg">
                    My Assets
                  </p>
                  <TEDropdown className="flex justify-center">
                    <TERipple className="w-full" rippleColor="light">
                      <TEDropdownToggle className="flex items-center whitespace-nowrap rounded bg-primary-100 hover:bg-primary-200 text-primary-800 p-2 justify-between rounded-xl w-full">
                        <img src="/images/filter.png" alt="filter" />
                      </TEDropdownToggle>
                    </TERipple>
                    <TEDropdownMenu className="w-full bg-primary-100">
                      {dropdownItems}
                    </TEDropdownMenu>
                  </TEDropdown>
                </div>
                <p>
                  <span>Sort by:</span>{' '}
                  <span className="font-bold font-montserratSemiBold">
                    {assetClassOptions[assetFilterMode].text}
                  </span>
                </p>
              </div>
              <div className="w-full mt-4 space-y-4">
                {activeWallet.claimedAssets.filter((a) => {
                  if (assetFilterMode === 0) return a.tokenizedAsset;
                  else return !a.tokenizedAsset;
                }).length > 1 ? ( // since every wallet must have an ETH asset check length greater than 1
                  // if filter is set on other tokens
                  assetFilterMode === 0 ? (
                    activeWallet.claimedAssets.map(
                      (asset, index) =>
                        asset.tokenizedAsset == true && (
                          <AssetItem
                            key={index}
                            image={asset.imageUrl}
                            assetCode={asset.assetCode}
                            amount={asset.amount.toString()}
                            usdPrice={asset.usdPrice.toString()}
                            nativePrice={asset.nativePrice.toString()}
                            currency={appUser.currency}
                            onclick={() => {
                              navigate(
                                `/dashboard/asset-details?wallet=${activeWallet.address}&assetCode=${asset.assetCode}&contractAddress=${asset.contractAddress}`,
                              );
                            }}
                          />
                        ),
                    )
                  ) : (
                    activeWallet.claimedAssets.map(
                      (asset, index) =>
                        !asset.tokenizedAsset && (
                          <AssetItem
                            key={index}
                            image={asset.imageUrl}
                            assetCode={asset.assetCode}
                            amount={asset.amount.toString()}
                            usdPrice={asset.usdPrice.toString()}
                            nativePrice={asset.nativePrice.toString()}
                            currency={appUser.currency}
                            onclick={() => {
                              navigate(
                                `/dashboard/asset-details?wallet=${activeWallet.address}&assetCode=${asset.assetCode}&contractAddress=${asset.contractAddress}`,
                              );
                            }}
                          />
                        ),
                    )
                  )
                ) : (
                  <div className="flex items-center justify-center w-full h-60">
                    <p className="text-lg">Nothing to show here</p>
                  </div>
                )}
              </div>
            </div>
          </div>
          <WalletOperations
            wallets={wallets}
            selectedAsset={selectedAsset}
            activeWallet={activeWallet}
            onWalletChanged={(wallet, index) => {
              setActiveWallet(wallet);
              setActiveWalletIndex(index);
              setSelectedAsset(
                activeWallet?.claimedAssets.find(
                  (a) => a.assetCode === '' && a.contractAddress === '',
                )!,
              );
              const ref = itemRefs.current[index];
              if (ref) {
                // Change background color
                ref.scrollIntoView({
                  behavior: 'smooth',
                  block: 'nearest',
                  inline: 'center',
                });
              }
            }}
            onAssetChanged={(asset) => {
              setSelectedAsset(asset);
            }}
            isWalletDetailsPage={false}
          />
        </div>
      </div>
      {/* Address copied modal */}
      <Modal
        showModal={showAddSubwalletModal}
        onClose={() => {
          resetForm();
          setShowAddSubwalletModal(false);
        }}
      >
        <>
          {!showConfirmView ? (
            <form
              id="send-asset-form"
              onSubmit={validateAndConfirm}
              className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center"
            >
              <p className="text-primary-800 w-full mb-5 text-lg text-center font-montserratSemiBold">
                Add Subwallet
              </p>
              <div className="flex w-full py-3 px-2 xl:px-5 xl:space-y-5 rounded-xl flex-col items-center bg-primary-100">
                <div className="flex flex-col space-y-2 items-center w-full py-3 justify-between space-y-5">
                  <p className="font-montserratSemiBold text-primary-800">
                    You are about to add a subwallet to your Trovo App account
                  </p>
                  <p className="text-sm text-primary-800">Choose a method</p>
                  <div className="flex flex-col items-start space-y-3 mt-6">
                    <div className="flex justify-between items-between space-x-2">
                      <input
                        type="radio"
                        id="addSubwalletMethodImport"
                        name="subwalletMethod"
                        className="w-4 h-4 accent-primary-800 cursor-pointer"
                        checked={formData.isImport}
                        onChange={() =>
                          setFormData({
                            ...formData,
                            newSecretKey: '',
                            isImport: true,
                          })
                        }
                      />
                      <label
                        htmlFor="addSubwalletMethodImport"
                        className="text-sm font-montserratSemiBold text-primary-800 cursor-pointer"
                      >
                        Import existing wallet
                      </label>
                    </div>
                    <div className="flex justify-between items-between space-x-2">
                      <input
                        type="radio"
                        id="addSubwalletMethodCreate"
                        name="subwalletMethod"
                        className="w-4 h-4 accent-primary-800 cursor-pointer"
                        checked={!formData.isImport}
                        onChange={() =>
                          setFormData({
                            ...formData,
                            newSecretKey: '',
                            isImport: false,
                          })
                        }
                      />
                      <label
                        htmlFor="addSubwalletMethodCreate"
                        className="text-sm font-montserratSemiBold text-primary-800 cursor-pointer"
                      >
                        Create new wallet
                      </label>
                    </div>
                  </div>
                  <p className="text-xs text-primary-800">
                    Please note that completing this process will attract some
                    charges
                  </p>
                </div>
              </div>
              <div className="w-full space-y-2">
                <TextInput
                  inputType="text"
                  label="Tag"
                  defaultValue={formData.tag}
                  placeholder="Tag"
                  onInputChange={(value) =>
                    setFormData({ ...formData, tag: value })
                  }
                />
                <div className={`w-full flex items-center justify-between`}>
                  <span className="text-sm mt-1">
                    {appUser.username}_{formData.tag}
                  </span>
                  <span
                    className={
                      formData.tag && getBytesLength(formData.tag) > 12
                        ? 'text-sm text-red-500'
                        : 'text-sm '
                    }
                  >
                    {formData.tag ? getBytesLength(formData.tag) : 0}
                    /12
                  </span>
                </div>
                <div className="w-full flex items-center">
                  {errorObj.tag && (
                    <span className="text-red-500 text-sm mt-1">
                      {errorObj.tag}
                    </span>
                  )}
                </div>
                {/* <p className="text-red-500">{error}</p> */}
              </div>
              <div className="w-full space-y-2">
                <TextInput
                  inputType="text"
                  label="Description"
                  placeholder="Savings wallet..."
                  defaultValue={formData.description}
                  onInputChange={(value) =>
                    setFormData({ ...formData, description: value })
                  }
                />
                <div
                  className={`w-full flex items-center ${
                    errorObj.description ? 'justify-between' : 'justify-end'
                  }`}
                >
                  {errorObj.description && (
                    <span className="text-red-500 text-sm mt-1">
                      {errorObj.description}
                    </span>
                  )}
                  <span
                    className={
                      formData.description &&
                      getBytesLength(formData.description) > 100
                        ? 'text-sm text-red-500'
                        : 'text-sm '
                    }
                  >
                    {formData.description
                      ? getBytesLength(formData.description)
                      : 0}
                    /100
                  </span>
                </div>
                {/* <p className="text-red-500">{error}</p> */}
              </div>
              {formData.isImport && (
                <div className="w-full space-y-2">
                  <TextInput
                    inputType="text"
                    label="Secret Key"
                    defaultValue={formData.isImport && formData.newSecretKey}
                    placeholder="Secret key..."
                    onInputChange={(value) =>
                      setFormData({ ...formData, newSecretKey: value })
                    }
                  />
                  <div className="w-full flex items-center">
                    {errorObj.newSecretKey && (
                      <span className="text-red-500 text-sm mt-1">
                        {errorObj.newSecretKey}
                      </span>
                    )}
                  </div>
                  {/* <p className="text-red-500">{error}</p> */}
                </div>
              )}
              <div></div>
              <div className="w-full self-center space-y-3">
                <Button
                  type="submit"
                  label="Continue"
                  additionalClasses="font-montserratSemiBold"
                  onclick={() => {}}
                />
              </div>
              <div className="w-full self-center space-y-3">
                <ButtonSecondary
                  label="Close"
                  additionalClasses="font-montserratSemiBold"
                  onclick={() => {
                    resetForm();
                    setShowAddSubwalletModal(false);
                  }}
                />
              </div>
              <div></div>
            </form>
          ) : (
            <>
              {' '}
              {!showAuthorizeView ? (
                <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
                  <p className="text-primary-800 w-full mb-5 text-lg text-center font-montserratSemiBold">
                    Add Subwallet
                  </p>
                  <div className="flex w-full py-3 px-2 xl:px-5 xl:space-y-5 rounded-xl flex-col items-center bg-primary-100">
                    <div className="flex flex-col space-y-2 items-center w-full py-3 justify-between space-y-5">
                      <p className="font-montserratSemiBold text-primary-800">
                        You have requested to create a subwallet with these
                        details:
                      </p>
                      <p className="text-sm text-primary-800 font-montserratSemiBold">
                        Tag
                      </p>
                      <p className="text-sm text-primary-800 ">
                        {formData.tag}
                      </p>
                      <p className="text-sm text-primary-800 font-montserratSemiBold">
                        Description
                      </p>
                      <p className="text-sm text-primary-800 ">
                        {formData.description}
                      </p>
                      <p className="text-sm text-primary-800 font-montserratSemiBold">
                        Address
                      </p>
                      <p className="text-sm text-primary-800 ">
                        {formData.newAddress}
                      </p>
                      <p className="text-xs text-primary-800">
                        Please note that completing this process will attract
                        some charges
                      </p>
                    </div>
                  </div>
                  <div></div>
                  <div className="w-full self-center space-y-3">
                    <Button
                      type="submit"
                      label="Create Subwallet"
                      additionalClasses="font-montserratSemiBold"
                      onclick={handleSend}
                    />
                  </div>
                  <div className="w-full self-center space-y-3">
                    <ButtonSecondary
                      label="Back"
                      additionalClasses="font-montserratSemiBold"
                      onclick={() => {
                        setShowConfirmView(false);
                      }}
                    />
                  </div>
                  <div></div>
                </div>
              ) : (
                <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
                  <p className="text-primary-800 w-full mb-5 text-lg text-center font-montserratSemiBold">
                    Authorize
                  </p>
                  <div className="flex w-full py-3 px-2 xl:px-5 xl:space-y-5 rounded-xl flex-col items-center bg-primary-100">
                    <div className="flex flex-col space-y-2 items-center w-full py-3 justify-between space-y-5">
                      {formData.transactionData.messages.map((m: any) => (
                        <p className="text-sm text-red-800 ">{m}</p>
                      ))}
                    </div>
                  </div>
                  <div></div>
                  <div className="w-full space-y-1">
                    <TextInput
                      label=""
                      leadingIcon="/images/lock.png"
                      inputType="password"
                      placeholder="Enter answer"
                      onInputChange={(newValue: string) => {
                        setPassword(newValue);
                      }}
                    />
                    {passwordErr && (
                      <p className="text-red-500 text-sm">{passwordErr}</p>
                    )}
                  </div>
                  <div className="w-full space-y-3">
                    <Button
                      label="Authorize with password"
                      additionalClasses="font-montserratSemiBold"
                      onclick={async () => {
                        setPasswordErr('');

                        if (!password) {
                          setPasswordErr('Please enter a password!');
                          return;
                        }

                        if (!(await isValidPassword())) {
                          setPasswordErr('Password is invalid!');
                          return;
                        }

                        const body = {
                          ...formData.transactionData,
                          primarySignature: signBase64Txn(
                            secretKey,
                            formData.transactionData?.transaction,
                            formData.transactionData?.networkPassPhrase,
                          ),
                          subWalletSignature: signBase64Txn(
                            formData.newSecretKey,
                            formData.transactionData?.transaction,
                            formData.transactionData?.networkPassPhrase,
                          ),
                        };

                        const payload = {
                          signer: activeWallet.signer,
                          address: activeWallet.address,
                          secretKey: secretKey,
                          body,
                        };

                        toggleLoader();
                        const res = await addSubwallet(payload);

                        const importedPayload = {
                          signer: appUser.address,
                          address: appUser.address,
                          secretKey: appUser.secretKeys[0],
                          body: { userId: appUser.username, import: 1 },
                        };
                        const { data: importData, error: importError } =
                          await getUser(importedPayload);

                        var userData = {} as User;
                        if (importData) {
                          userData = deserializeUserData(importData);
                        } else if (importError) {
                          const err = importError as ErrorResponse;
                          showNotification(
                            'error',
                            err.data.message ??
                              'Sorry we could not complete the request. Please try again.',
                          );
                        }
                        console.log('response', res);

                        if ('data' in res) {
                          const encryptor = new Encryptor();
                          const base64EncryptedSecretKey =
                            await encryptor.encryptData(
                              formData.newSecretKey,
                              password,
                              appUser.address,
                            );

                          const passwordHash = await encryptor.createHash(
                            userData.username,
                          );
                          const encryptedPassword = await encryptor.encryptData(
                            password,
                            passwordHash,
                            appUser.address,
                          );

                          const user = {
                            ...userData,
                            password: encryptedPassword,
                            isLoggedIn: true,
                            currency: 'USD',
                            secretKeys: [
                              ...appUser.secretKeys,
                              base64EncryptedSecretKey,
                            ],
                          };

                          const hash = await encryptor.createHash(
                            user.username,
                          );
                          const base64EncryptedUserData =
                            await encryptor.encryptData(
                              JSON.stringify(user),
                              hash,
                              user.address,
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
                              username: formData.tag,
                              address: formData.newAddress,
                              secretKey: base64EncryptedSecretKey,
                            }),
                          );
                          setShowAddSubwalletModal(false);
                          setShowBackupModal(true);
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
                  <div></div>
                </div>
              )}
            </>
          )}
        </>
      </Modal>
      <Modal
        showModal={showBackupModal}
        onClose={() => {
          setShowBackupModal(false);
        }}
      >
        <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
          <img src="/images/launch.png" alt="success" />
          <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
            <p className="text-primary-800 text-md xl:text-lg font-semibold">
              Congratulations!
            </p>
            <p className="text-primary-800 text-md xl:text-lg">
              Your wallet has been successfully created
            </p>
            <p className="text-primary-800 text-md xl:text-lg">
              We strongly recommend that you backup your wallet before
              proceeding
            </p>
            <p className="text-primary-800 text-md xl:text-lg">
              Backing up your wallet is a way to restore your wallet if you
              loose your device
            </p>
          </div>
          <div className="w-3/4">
            <Button
              label="Backup"
              onclick={() => {
                navigate('/backup');
              }}
            />
          </div>
          <div className="w-3/4">
            <ButtonSecondary
              label="Skip"
              onclick={() => {
                setShowBackupModal(false);
                navigate('/dashboard');
              }}
            />
          </div>
        </div>
      </Modal>
    </>
  );
}
