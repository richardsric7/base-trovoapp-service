import Button from '../../components/button';
import Header from '../../components/header';
import Tabs from '../../components/tabs';
import { useNavigate } from 'react-router-dom';
import { useEffect, useRef, useState } from 'react';
import WalletCard from '../../components/walletCard';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../store/reduxStore';
import {
  setActiveWallet as setActiveWalletToStore,
  setActiveAsset as setActiveAssetToStore,
} from '../../store/appStateSlice';
import React from 'react';
import {
  calculateFiatValue,
  formatToDecimal,
  getAssetCode,
  getBytesLength,
  getExplorerBaseUrl,
  canInitiate,
  totalWalletBalanceInCurrency,
} from '../../utils/utilities';
import { Wallet } from '../../types/wallet';
import {
  TEDropdown,
  TEDropdownItem,
  TEDropdownMenu,
  TEDropdownToggle,
  TERipple,
} from 'tw-elements-react';
import Dropdown from '../../components/dropdown';
import TextInput from '../../components/textInput';
import { Asset } from '../../types/asset';
import { showNotification, toggleLoader } from '../../utils/showToaster';
import {
  useLazyReceiveAssetQuery,
  useSendAssetMutation,
  useSwapAssetMutation,
} from '../../store/api/walletApis';
import { ErrorResponse } from '../../store/api/baseapi/axiosBaseQuery';
import Modal from '../../components/modal';
import { Encryptor } from '../../utils/encryptor';
import { signBase64Txn } from '../../utils/trovoSDK';
import ButtonSecondary from '../../components/buttonSecondary';
import { TransactionDirection } from '../../types/transactionInfo';
import WalletDropdown from '../../components/walletDropdown';
import AssetDropdown from '../../components/assettDropdown';
import AssetItem from '../../components/assetItem';

export default function WalletView() {
  const ref = useRef<HTMLDivElement>(null);
  const dispatch = useDispatch();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const appState = useSelector((state: RootState) => state.appState!);
  const navigate = useNavigate();
  let mutableWalletArray = [...appUser.userWallets];
  const [wallets, setWallets] = useState(
    mutableWalletArray.sort((w) => (w.primaryWallet ? 0 : 1)),
  );
  const [activeWalletIndex, setActiveWalletIndex] = useState(0);
  const [activeWallet, setActiveWallet] = useState<Wallet>(wallets[0]);
  const [assetFilterMode, setAssetFilterMode] = useState(0);
  const [walletActionMode, setWalletActionMode] = useState(1);
  const itemRefs = useRef<HTMLDivElement[]>([]);
  const fiatRates = useSelector((state: RootState) => state.cache.fiatRates);
  const [secretKey, setSecretKey] = useState('');
  const [password, setPassword] = useState('');
  const [passwordErr, setPasswordErr] = useState('');
  const [showConfirmSendModal, setShowConfirmSendModal] = useState(false);
  const [showSendSuccessModal, setShowSendSuccessModal] = useState(false);
  const [showConfirmSwapModal, setShowConfirmSwapModal] = useState(false);
  const [showSwapSuccess, setShowSwapSuccess] = useState(false);
  const [qrCodeLink, setQrCodeLink] = useState('');
  const [requestSpecificAmount, setRequestSpecificAmount] = useState(false);
  const [receiptQuery, setReceiptQuery] = useState('');
  const [selectedAsset, setSelectedAsset] = useState<Asset>(
    activeWallet?.claimedAssets.find(
      (a) => a.assetCode === '' && a.assetIssuer === '',
    )!,
  );
  const [currentTabIndex, setCurrentTabIndex] = useState(1);
  const [gas, setGas] = useState(
    activeWallet?.claimedAssets.find(
      (a) => a.assetCode === '' && a.assetIssuer === '',
    )?.amount!,
  );
  const [sendAsset] = useSendAssetMutation();
  const [swapAsset] = useSwapAssetMutation();
  const [receiveAsset] = useLazyReceiveAssetQuery();

  type WalletFormData = {
    sendTo?: string;
    amount?: string;
    memo?: string;
    swapFrom?: Asset;
    swapTo?: Asset;
    transactionData?: any;
  };

  const [formData, setFormData] = useState<WalletFormData>({
    sendTo: '',
    amount: '',
    memo: '',
    swapFrom: undefined,
    swapTo: undefined,
    transactionData: undefined,
  });

  const [errorObj, setErrorObj] = useState({
    sendTo: '',
    amount: '',
    memo: '',
    swapFrom: '',
    swapTo: '',
  });

  const resetForm = () => {
    setErrorObj({
      sendTo: '',
      amount: '',
      memo: '',
      swapFrom: '',
      swapTo: '',
    });
    setFormData({
      sendTo: '',
      amount: '',
      memo: '',
      swapFrom: undefined,
      swapTo: undefined,
      transactionData: undefined,
    });
  };

  const validateSendForm = (): boolean => {
    let isValid = true;
    let newObj = errorObj;

    if (!formData.sendTo) {
      newObj = {
        ...newObj,
        sendTo: 'Please enter receiver username or public key',
      };
      isValid = false;
    } else if (formData.sendTo.length < 3) {
      newObj = {
        ...newObj,
        sendTo: 'Invalid username or public key',
      };
      isValid = false;
    } else {
      newObj = {
        ...newObj,
        sendTo: '',
      };
    }

    if (!formData.amount) {
      newObj = {
        ...newObj,
        amount: 'Please enter amount to send',
      };
      isValid = false;
    } else if (isNaN(Number(formData.amount))) {
      newObj = {
        ...newObj,
        amount: 'Please enter a valid amount',
      };
      isValid = false;
    } else if (Number(formData.amount) <= 0) {
      newObj = {
        ...newObj,
        amount: 'Please enter a valid amount',
      };
      isValid = false;
    } else if (Number(formData.amount) > selectedAsset?.amount!) {
      newObj = {
        ...newObj,
        amount: "You don't have sufficient balance",
      };
      isValid = false;
    } else if (
      getAssetCode(selectedAsset!.assetCode) == 'XBN' &&
      Number(formData.amount) > selectedAsset!.amount! - 6
    ) {
      newObj = {
        ...newObj,
        amount: "You don't have sufficient balance",
      };
      isValid = false;
    } else {
      newObj = {
        ...newObj,
        amount: '',
      };
    }

    if (formData.memo && getBytesLength(formData.memo) > 28) {
      newObj = {
        ...newObj,
        memo: 'Memo length cannot be more than 28 bytes',
      };
      isValid = false;
    } else {
      newObj = {
        ...newObj,
        memo: '',
      };
    }

    setErrorObj({ ...newObj });
    return isValid;
  };

  const validateRecieveForm = (): boolean => {
    let isValid = true;
    let newObj = errorObj;

    if (!formData.amount) {
      newObj = {
        ...newObj,
        amount: 'Please enter amount to recieve',
      };
      isValid = false;
    } else if (isNaN(Number(formData.amount))) {
      newObj = {
        ...newObj,
        amount: 'Please enter a valid amount',
      };
    } else {
      newObj = {
        ...newObj,
        amount: '',
      };
    }

    if (formData.memo && getBytesLength(formData.memo) > 28) {
      newObj = {
        ...newObj,
        memo: 'Memo length cannot be more than 28 bytes',
      };
    } else {
      newObj = {
        ...newObj,
        memo: '',
      };
    }

    setErrorObj({ ...newObj });
    return isValid;
  };

  const validateSwapForm = (): boolean => {
    let isValid = true;
    let newObj = errorObj;

    if (!formData.amount) {
      newObj = {
        ...newObj,
        amount: 'Please enter amount to swap',
      };
      isValid = false;
    } else if (isNaN(Number(formData.amount))) {
      newObj = {
        ...newObj,
        amount: 'Please enter a valid amount',
      };
    } else if (Number(formData.amount) > formData.swapFrom?.amount!) {
      newObj = {
        ...newObj,
        amount: "You don't have sufficient balance",
      };
      isValid = false;
    } else if (
      getAssetCode(selectedAsset!.assetCode) == 'XBN' &&
      Number(formData.amount) > selectedAsset!.amount! - 6
    ) {
      newObj = {
        ...newObj,
        amount: "You don't have sufficient balance",
      };
      isValid = false;
    } else {
      newObj = {
        ...newObj,
        amount: '',
      };
    }

    if (!formData.swapFrom) {
      newObj = {
        ...newObj,
        swapFrom: 'Please select source asset',
      };
    } else {
      newObj = {
        ...newObj,
        swapFrom: '',
      };
    }

    if (!formData.swapTo) {
      newObj = {
        ...newObj,
        swapTo: 'Please select destination asset',
      };
    } else {
      newObj = {
        ...newObj,
        swapTo: '',
      };
    }

    setErrorObj({ ...newObj });
    return isValid;
  };

  useEffect(() => {
    (async () => {
      const encryptor = new Encryptor();
      setSecretKey(await encryptor.getSecretKey(appUser));
    })();
  }, []);

  useEffect(() => {
    setWallets(
      mutableWalletArray
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
        .sort((w) => (w.primaryWallet ? 0 : 1)),
    );

    setActiveWalletIndex(0);
    const ref = itemRefs.current[0];
    if (ref) {
      // Change background color
      ref.scrollIntoView({
        behavior: 'smooth',
        block: 'nearest',
        inline: 'center',
      });
    }
  }, [currentTabIndex]);

  useEffect(() => {
    setGas(
      activeWallet?.claimedAssets.find(
        (a) => a.assetCode === '' && a.assetIssuer === '',
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

  const handleSend = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!validateSendForm()) {
      return;
    }

    const payload = {
      signer: activeWallet.signer,
      publicKey: activeWallet.publicKey,
      secretKey: secretKey,
      body: {
        isSharedWallet: activeWallet.sharedAccessEnabled,
        destination: formData.sendTo,
        memo: formData.memo,
        amount: formData.amount,
        assetCode: getAssetCode(selectedAsset.assetCode),
        assetIssuer: selectedAsset.assetIssuer,
      },
    };

    toggleLoader();
    const res = await sendAsset(payload);
    console.log('response', res);

    if ('data' in res) {
      console.log('response', res);
      setFormData({
        ...formData,
        transactionData: res.data,
      });
      setShowConfirmSendModal(true);
    } else if ('error' in res) {
      const errorResponse = res.error as ErrorResponse;
      showNotification(
        'error',
        errorResponse.data.message ?? 'Something went wrong. Please try again.',
      );
    }
    toggleLoader();
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

  const handleRecieve = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!validateRecieveForm()) {
      return;
    }

    const payload = {
      signer: activeWallet.signer,
      publicKey: activeWallet.publicKey,
      secretKey: secretKey,
      body: {
        destination: formData.sendTo,
        memo: formData.memo,
        amount: formData.amount,
        publicKey: activeWallet.publicKey,
        alias: activeWallet.alias,
        assetCode: selectedAsset.assetCode,
        assetIssuer: selectedAsset.assetIssuer,
      },
    };

    toggleLoader();
    const res = await receiveAsset(payload);
    console.log('response', res);

    if ('data' in res) {
      console.log('response', res);
      setQrCodeLink(res.data.qrCode);
      setRequestSpecificAmount(false);
    } else if ('error' in res) {
      const errorResponse = res.error as ErrorResponse;
      showNotification(
        'error',
        errorResponse.data.message ?? 'Something went wrong. Please try again.',
      );
    }
    toggleLoader();
  };

  const handleSwap = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!validateSwapForm()) {
      return;
    }

    const payload = {
      signer: activeWallet.signer,
      publicKey: activeWallet.publicKey,
      secretKey: secretKey,
      body: {
        isSharedWallet: activeWallet.sharedAccessEnabled,
        destinationAssetCode: formData.swapTo?.assetCode,
        destinationAssetIssuer: formData.swapTo?.assetIssuer,
        sourceAssetCode: formData.swapFrom?.assetCode,
        sourceAssetIssuer: formData.swapFrom?.assetIssuer,
        sourceAmount: formData.amount,
      },
    };

    toggleLoader();
    const res = await swapAsset(payload);
    console.log('response', res);

    if ('data' in res) {
      console.log('response', res);
      setFormData({
        ...formData,
        transactionData: res.data,
      });
      setShowConfirmSwapModal(true);
    } else if ('error' in res) {
      const errorResponse = res.error as ErrorResponse;
      showNotification(
        'error',
        errorResponse.data.message ?? 'Something went wrong. Please try again.',
      );
    }
    toggleLoader();
  };

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

  return (
    <div className="flex text-primary-800 text-sm overflow-x-hidden md:text-md flex-col space-y-5 p-5">
      <Header isHomeView />
      <div className="flex md:h-screen w-full items-center justify-center">
        <div className="h-full w-full md:p-3">
          <div className="flex md:pt-5 p-5 rounded-lg flex-col items-center bg-primary-100 max-w-5xl min-h-full">
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
                  onclick={() => {
                    /* */
                  }}
                />
              </div>
            </div>
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
                        (a) => a.assetCode === '' && a.assetIssuer === '',
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
                        (a) => a.assetCode === '' && a.assetIssuer === '',
                      )!,
                    );

                    if (!canInitiate(wallets[activeWalletIndex + 1])) {
                      setWalletActionMode(1);
                    }
                    resetForm();

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
                  <img src="/images/gas.png" alt="gas" width="30" height="30" />
                  <span>Gas</span>
                </div>
                <img
                  src="/images/export.png"
                  alt="export"
                  width="30"
                  height="30"
                />
              </div>
              <p>{formatToDecimal(gas)} XBN</p>
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
              }).length > 1 ? ( // since every wallet must have an XBN asset check length greater than 1
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
                            dispatch(setActiveWalletToStore(activeWallet));
                            dispatch(
                              setActiveAssetToStore(
                                `${asset.assetCode}|${asset.assetIssuer}`,
                              ),
                            );
                            navigate('/dashboard/asset-details');
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
                            dispatch(setActiveWalletToStore(activeWallet));
                            dispatch(
                              setActiveAssetToStore(
                                `${asset.assetCode}|${asset.assetIssuer}`,
                              ),
                            );
                            console.log(appState);
                            navigate('/dashboard/asset-details');
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
        <div className="hidden md:block w-3/6 flex flex-col space-y-5 h-full py-3 lg:px-3">
          <div className="flex space-y-6 rounded-lg py-5 px-6 flex-col items-center bg-primary-100">
            <div className="flex justify-between items-center space-x-5 mt-3">
              <button
                className={`flex space-y-3 ring-1 rounded-full px-6 py-2 flex-col items-center ${
                  walletActionMode === 1
                    ? 'ring-primary bg-primary-200 font-bold font-montserratSemiBold'
                    : 'ring-gray-400 text-gray-400'
                }`}
                onClick={() => {
                  resetForm();
                  setWalletActionMode(1);
                }}
              >
                Receive
              </button>
              <button
                disabled={!canInitiate(activeWallet)}
                className={`flex space-y-3 ring-1 rounded-full px-6 py-2 flex-col items-center disabled:bg-gray-200 ${
                  walletActionMode === 0
                    ? 'ring-primary bg-primary-200 font-bold font-montserratSemiBold'
                    : 'ring-gray-400 text-gray-400'
                }`}
                onClick={() => {
                  resetForm();
                  setWalletActionMode(0);
                }}
              >
                Send
              </button>
              <button
                disabled={!canInitiate(activeWallet)}
                className={`flex space-y-3 ring-1 rounded-full px-6 py-2 flex-col items-center disabled:bg-gray-200 ${
                  walletActionMode === 2
                    ? 'ring-primary bg-primary-200 font-bold font-montserratSemiBold'
                    : 'ring-gray-400 text-gray-400'
                }`}
                onClick={() => {
                  resetForm();
                  setWalletActionMode(2);
                }}
              >
                Swap
              </button>
            </div>
            {(() => {
              switch (walletActionMode) {
                case 0:
                  return (
                    <form
                      id="send-asset-form"
                      onSubmit={handleSend}
                      className="w-full space-y-5"
                    >
                      <div className="space-y-3 w-full">
                        <p>Select wallet</p>
                        <WalletDropdown
                          label={activeWallet.alias}
                          defaultValue={{
                            text: activeWallet.alias,
                            value: activeWallet,
                            index: activeWalletIndex,
                          }}
                          options={[
                            ...wallets
                              .filter((w) => canInitiate(w))
                              .map((w, index) => ({
                                text: w.alias,
                                value: w,
                                index,
                              })),
                          ]}
                          onSelect={(selectedItem) => {
                            console.log(selectedItem);
                            setActiveWallet(selectedItem.value);
                            setActiveWalletIndex(selectedItem.index);
                            const ref = itemRefs.current[selectedItem.index];
                            if (ref) {
                              // Change background color
                              ref.scrollIntoView({
                                behavior: 'smooth',
                                block: 'nearest',
                                inline: 'center',
                              });
                            }
                          }}
                          key={`${activeWallet.publicKey}-${activeWalletIndex}`}
                        />
                      </div>
                      <div className="space-y-3 w-full">
                        <TextInput
                          inputType="text"
                          label="Send to"
                          defaultValue={formData.sendTo}
                          onInputChange={(value) => {
                            setFormData({ ...formData, sendTo: value });
                          }}
                          error={errorObj.sendTo}
                        />
                      </div>
                      <div className="w-full space-y-2">
                        <div className="flex items-center space-x-3 space-y-3 w-full">
                          <div className="w-4/6">
                            <TextInput
                              inputType="number"
                              label="Token quantity/Amount"
                              onInputChange={(value) => {
                                setFormData({ ...formData, amount: value });
                              }}
                            />
                          </div>
                          <div className="w-2/6 self-end">
                            <Dropdown
                              label={getAssetCode(selectedAsset.assetCode)}
                              options={[
                                ...activeWallet.claimedAssets.map((asset) => ({
                                  text:
                                    asset.assetCode === ''
                                      ? 'XBN'
                                      : asset.assetCode,
                                  value: asset,
                                })),
                              ]}
                              onSelect={(asset) => {
                                console.log(asset);
                                setSelectedAsset(asset.value);
                              }}
                            />
                          </div>
                        </div>
                        {errorObj.amount && (
                          <p className="text-red-500 text-sm mt-1">
                            {errorObj.amount}
                          </p>
                        )}
                        {!appState.hideBalances && (
                          <div className="w-full flex justify-between items-center">
                            <span>
                              {formData.amount &&
                              !isNaN(Number(formData.amount))
                                ? formatToDecimal(Number(formData.amount))
                                : 0}{' '}
                              {getAssetCode(selectedAsset.assetCode)}
                            </span>
                            <span>
                              {formatToDecimal(selectedAsset.amount)}{' '}
                              {getAssetCode(selectedAsset.assetCode)}
                            </span>
                          </div>
                        )}
                      </div>
                      <div className="space-y-3 w-full">
                        <TextInput
                          inputType="text"
                          label="Add memo (optional)"
                          onInputChange={(value) => {
                            setFormData({ ...formData, memo: value });
                          }}
                        />
                        <div
                          className={`w-full flex items-center ${
                            errorObj.memo ? 'justify-between' : 'justify-end'
                          }`}
                        >
                          {errorObj.memo && (
                            <span className="text-red-500 text-sm mt-1">
                              {errorObj.memo}
                            </span>
                          )}
                          <span
                            className={
                              formData.memo &&
                              getBytesLength(formData.memo) > 28
                                ? 'text-red-500'
                                : ''
                            }
                          >
                            {formData.memo ? getBytesLength(formData.memo) : 0}
                            /28
                          </span>
                        </div>
                      </div>
                      <Button
                        type="submit"
                        label="Proceed"
                        onclick={() => {}}
                      />
                    </form>
                  );
                case 1:
                  return (
                    <form
                      id="recieve-asset-form"
                      onSubmit={handleRecieve}
                      className="w-full space-y-5"
                    >
                      <div className="w-full space-y-4">
                        {!qrCodeLink ? (
                          <div className="flex justify-between space-y-5 w-full">
                            <div className="w-4/6 mr-2 space-y-2">
                              <p>Receiving wallet</p>
                              <WalletDropdown
                                label={activeWallet.alias}
                                defaultValue={{
                                  text: activeWallet.alias,
                                  value: activeWallet,
                                  index: activeWalletIndex,
                                }}
                                options={[
                                  ...wallets
                                    .filter((w) => canInitiate(w))
                                    .map((w, index) => ({
                                      text: w.alias,
                                      value: w,
                                      index,
                                    })),
                                ]}
                                onSelect={(selectedItem) => {
                                  console.log(selectedItem);
                                  setActiveWallet(selectedItem.value);
                                  setActiveWalletIndex(selectedItem.index);
                                  setSelectedAsset(
                                    activeWallet?.claimedAssets.find(
                                      (a) =>
                                        a.assetCode === '' &&
                                        a.assetIssuer === '',
                                    )!,
                                  );
                                  const ref =
                                    itemRefs.current[selectedItem.index];
                                  if (ref) {
                                    // Change background color
                                    ref.scrollIntoView({
                                      behavior: 'smooth',
                                      block: 'nearest',
                                      inline: 'center',
                                    });
                                  }
                                }}
                                key={`${activeWallet.publicKey}-${activeWalletIndex}`}
                              />
                            </div>
                            <div className="w-2/6 self-end">
                              <Dropdown
                                label={getAssetCode(selectedAsset.assetCode)}
                                options={[
                                  ...activeWallet.claimedAssets.map((w) => ({
                                    text:
                                      w.assetCode === '' ? 'XBN' : w.assetCode,
                                    value: w,
                                  })),
                                ]}
                                onSelect={(item) => {
                                  setSelectedAsset(item.value);
                                }}
                              />
                            </div>
                          </div>
                        ) : (
                          <>
                            <p className="text-lg font-montserratSemiBold">
                              Receive {formData.amount}{' '}
                              {getAssetCode(selectedAsset.assetCode)}
                            </p>
                            <div>
                              <p className="font-montserratSemiBold">
                                Receiving wallet
                              </p>
                              <div className="flex justify-between break-all">
                                <p className="max-w-sm">{activeWallet.alias}</p>
                                <button
                                  type="button"
                                  onClick={() =>
                                    navigator.clipboard
                                      .writeText(activeWallet.alias)
                                      .then(() => {
                                        showNotification(
                                          'info',
                                          'Wallet alias copied!',
                                        );
                                      })
                                  }
                                >
                                  <img src="/images/copy.png" alt="copy" />
                                </button>
                              </div>
                            </div>
                            <div>
                              <p className="font-montserratSemiBold">Memo</p>
                              <p className="max-w-sm">{formData.memo}</p>
                            </div>
                          </>
                        )}
                        {requestSpecificAmount ? (
                          <>
                            <div className="w-full space-y-2">
                              <div className="flex items-center space-x-3 space-y-3 w-full">
                                <div className="w-full">
                                  <TextInput
                                    inputType="text"
                                    label="Token quantity/Amount"
                                    onInputChange={(value) => {
                                      setFormData({
                                        ...formData,
                                        amount: value,
                                      });
                                    }}
                                  />
                                </div>
                              </div>
                              {errorObj.amount && (
                                <p className="text-red-500 text-sm mt-1">
                                  {errorObj.amount}
                                </p>
                              )}
                            </div>
                            <div className="space-y-3 w-full">
                              <TextInput
                                inputType="text"
                                label="Add memo (optional)"
                                onInputChange={(value) => {
                                  setFormData({ ...formData, memo: value });
                                }}
                              />
                              <div
                                className={`w-full flex items-center ${
                                  errorObj.memo
                                    ? 'justify-between'
                                    : 'justify-end'
                                }`}
                              >
                                {errorObj.memo && (
                                  <span className="text-red-500 text-sm mt-1">
                                    {errorObj.memo}
                                  </span>
                                )}
                                <span
                                  className={
                                    formData.memo &&
                                    getBytesLength(formData.memo) > 28
                                      ? 'text-red-500'
                                      : ''
                                  }
                                >
                                  {formData.memo
                                    ? getBytesLength(formData.memo)
                                    : 0}
                                  /28
                                </span>
                              </div>
                            </div>
                            <Button
                              type="submit"
                              label="Proceed"
                              onclick={() => {}}
                            />
                            <ButtonSecondary
                              label="Cancel"
                              onclick={() => setRequestSpecificAmount(false)}
                            />
                          </>
                        ) : (
                          <div className="space-y-4">
                            <div>
                              <p className="font-montserratSemiBold">
                                Receive from non Trovo App
                              </p>
                              <div className="flex justify-between break-all">
                                <p className="max-w-sm">
                                  {activeWallet.publicKey}
                                </p>
                                <button
                                  type="button"
                                  onClick={() =>
                                    navigator.clipboard
                                      .writeText(activeWallet.publicKey)
                                      .then(() => {
                                        showNotification(
                                          'info',
                                          'Public key copied!',
                                        );
                                      })
                                  }
                                >
                                  <img src="/images/copy.png" alt="copy" />
                                </button>
                              </div>
                            </div>
                            <img
                              src={
                                qrCodeLink ? qrCodeLink : selectedAsset.qrCode
                              }
                              alt="qrcode"
                            />
                            {!qrCodeLink ? (
                              <Button
                                label="Request specific amount"
                                onclick={() => setRequestSpecificAmount(true)}
                              />
                            ) : (
                              <Button
                                label="Reset"
                                onclick={() => setQrCodeLink('')}
                              />
                            )}
                          </div>
                        )}
                      </div>
                    </form>
                  );
                case 2:
                  return (
                    <form
                      id="swap-asset-form"
                      onSubmit={handleSwap}
                      className="w-full space-y-5"
                    >
                      <div className="w-full space-y-5">
                        <div className="space-y-3 w-full">
                          <p>Select wallet</p>
                          <WalletDropdown
                            label={activeWallet.alias}
                            defaultValue={{
                              text: activeWallet.alias,
                              value: activeWallet,
                              index: activeWalletIndex,
                            }}
                            options={[
                              ...wallets
                                .filter((w) => canInitiate(w))
                                .map((w, index) => ({
                                  text: w.alias,
                                  value: w,
                                  index,
                                })),
                            ]}
                            onSelect={(selectedItem) => {
                              console.log(selectedItem);
                              setActiveWallet(selectedItem.value);
                              setActiveWalletIndex(selectedItem.index);
                              const ref = itemRefs.current[selectedItem.index];
                              if (ref) {
                                // Change background color
                                ref.scrollIntoView({
                                  behavior: 'smooth',
                                  block: 'nearest',
                                  inline: 'center',
                                });
                              }
                            }}
                            key={`${activeWallet.publicKey}-${activeWalletIndex}`}
                          />
                        </div>
                        <div className="w-full space-y-2">
                          <p>Swap from</p>
                          <AssetDropdown
                            label={
                              formData.swapFrom
                                ? getAssetCode(formData.swapFrom.assetCode)
                                : 'Choose asset'
                            }
                            options={[
                              ...activeWallet.claimedAssets
                                .filter(
                                  (a) =>
                                    a.assetCode !== formData.swapTo?.assetCode,
                                )
                                .map((w) => ({
                                  text: getAssetCode(w.assetCode),
                                  value: w,
                                })),
                            ]}
                            onSelect={(asset) => {
                              var data = {
                                ...formData,
                                swapFrom: asset.value,
                              };

                              setFormData(data);
                              console.log(formData);
                            }}
                          />
                          {errorObj.amount && (
                            <p className="text-red-500 text-sm mt-1">
                              {errorObj.swapFrom}
                            </p>
                          )}
                        </div>
                        <div className="w-full space-y-2">
                          <p>Swap to</p>
                          <Dropdown
                            label={
                              formData.swapTo
                                ? getAssetCode(formData.swapTo.assetCode)
                                : 'Choose asset'
                            }
                            options={[
                              ...appUser.curatedSwapList
                                .filter(
                                  (a) =>
                                    a.assetCode !==
                                    formData.swapFrom?.assetCode,
                                )
                                .map((w) => ({
                                  text: getAssetCode(w.assetCode),
                                  value: w,
                                })),
                            ]}
                            onSelect={(asset) => {
                              var data = {
                                ...formData,
                                swapTo: asset.value,
                              };

                              setFormData(data);
                              console.log(
                                formData.swapTo
                                  ? getAssetCode(formData.swapTo.assetCode)
                                  : 'Choose asset',
                              );
                            }}
                          />
                          {errorObj.amount && (
                            <p className="text-red-500 text-sm mt-1">
                              {errorObj.swapTo}
                            </p>
                          )}
                        </div>
                        <div className="w-full space-y-2">
                          <div className="flex items-center space-x-3 space-y-3 w-full">
                            <div className="w-full">
                              <TextInput
                                inputType="text"
                                label="Token quantity/Amount"
                                onInputChange={(value) => {
                                  setFormData({ ...formData, amount: value });
                                }}
                                error={errorObj.amount}
                              />
                            </div>
                          </div>
                          {!appState.hideBalances && formData.swapFrom && (
                            <div className="w-full flex justify-between items-center">
                              <span>
                                {formData.amount &&
                                !isNaN(Number(formData.amount))
                                  ? formatToDecimal(Number(formData.amount))
                                  : 0}{' '}
                                {getAssetCode(formData.swapFrom?.assetCode)}
                              </span>
                              <span>
                                {formatToDecimal(formData.swapFrom?.amount)}{' '}
                                {getAssetCode(formData.swapFrom?.assetCode)}
                              </span>
                            </div>
                          )}
                        </div>
                        <Button
                          type="submit"
                          label="Proceed"
                          onclick={() => {}}
                        />
                      </div>
                    </form>
                  );
              }
            })()}
          </div>
          <div className="flex w-full space-y-3 py-7 px-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
            <img src="/images/bankNotes.png" alt="bank notes" />
            <p className="text-center font-semibold">
              Make Deposits and Withdrawals on your assets with ease on Trovo
              Wallet
            </p>
            <div className="w-3/4">
              <Button
                label="Deposit/Withdraw"
                onclick={() => {
                  /* */
                }}
              />
            </div>
          </div>
          <div />
          <div />
          <div />
          <div />
        </div>
      </div>
      <Modal
        showModal={showConfirmSendModal}
        onClose={() => {
          setShowConfirmSendModal(false);
        }}
      >
        <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
          <p className="text-primary-800 w-full mb-5 text-lg md:text-xl font-montserratSemiBold">
            Confirm transaction
          </p>

          {formData.transactionData?.messages > 0 && (
            <div className="flex flex-col items-center">
              <span className="font-montserratSemiBold">Note:</span>
              {formData.transactionData?.messages.map((message: string) => (
                <p className="text-red-500">{message}</p>
              ))}
            </div>
          )}
          <p>You are about to send</p>
          <div className="w-full bg-primary-100 rounded-xl py-5 space-y-2 text-center">
            <p className="font-montserratSemiBold">
              {formData.amount} {getAssetCode(selectedAsset.assetCode)}
            </p>
            <p>
              {`${calculateFiatValue(
                Number(formData.amount || 0),
                (fiatRates as Record<string, number>)[
                  appUser.currency.toUpperCase()
                ] || 0,
                selectedAsset.usdPrice,
              )} ${appUser.currency.toUpperCase()}`}
            </p>
          </div>
          <p>To</p>
          <div className="flex space-x-3 justify-center items-center px-10 w-full bg-primary-100 rounded-xl py-5 space-y-2">
            <img
              src={formData.transactionData?.destinationThumbnail}
              className="rounded-full w-16 h-16"
            />
            <p>
              {formData.transactionData?.destinationFirstName}{' '}
              {formData.transactionData?.destinationLastName}
            </p>
          </div>
          {formData.memo && (
            <>
              <p>Description/Memo</p>
              <div className="flex space-x-3 justify-center items-center px-10 w-full bg-primary-100 rounded-xl py-5 space-y-2">
                {formData.memo}
              </div>
            </>
          )}
          {formData.transactionData?.fee && (
            <>
              <p>Service Fee</p>
              <div className="flex flex-col space-y-3 justify-center items-center px-10 w-full bg-primary-100 rounded-xl py-5 space-y-2">
                <p className="space-x-1">
                  <span className="font-montserratSemiBold">Fee:</span>
                  <span>{formData.transactionData?.fee}%</span>
                </p>
                <p className="space-x-1">
                  <span className="font-montserratSemiBold">
                    Amount (Calculated):
                  </span>
                  <span>
                    {formData.transactionData?.feeAmount}{' '}
                    {getAssetCode(selectedAsset.assetCode)}
                  </span>
                </p>
                {formData.transactionData.vatAmount > 0 && (
                  <>
                    <p className="space-x-1">
                      <span className="font-montserratSemiBold">VAT:</span>
                      <span>{formData.transactionData?.vat}%</span>
                    </p>
                    <p className="space-x-1">
                      <span className="font-montserratSemiBold">
                        VAT Amount:
                      </span>
                      <span>
                        {formData.transactionData?.vatAmount}{' '}
                        {getAssetCode(selectedAsset.assetCode)}
                      </span>
                    </p>
                  </>
                )}
              </div>
            </>
          )}
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
                  isSharedWallet: activeWallet.sharedAccessEnabled,
                  ...formData.transactionData,
                  commit: 1,
                  transactionSignature: signBase64Txn(
                    secretKey,
                    formData.transactionData?.transaction,
                    formData.transactionData?.networkPassPhrase,
                  ),
                };

                const payload = {
                  signer: activeWallet.signer,
                  publicKey: activeWallet.publicKey,
                  secretKey: secretKey,
                  body,
                };
                console.log('body', payload);

                toggleLoader();
                const res = await sendAsset(payload);
                console.log('response', res);

                if ('data' in res) {
                  console.log('response', res);
                  setFormData({
                    ...formData,
                    transactionData: res.data,
                  });

                  const queryString = btoa(
                    JSON.stringify({
                      to: `${res.data.destinationFirstName} ${res.data.destinationLastName} [${res.data.destination}]`,
                      from: `${appUser.firstName} ${appUser.lastName} [${appUser.username}]`,
                      fromPublicKey: activeWallet.publicKey,
                      toPublicKey: '',
                      memo: res.data.memo,
                      amount: res.data.amount,
                      transactionId: res.data.transactionId,
                      assetCode: res.data.assetCode,
                      assetIssuer: res.data.assetIssuer,
                      transactionDate: new Date(),
                      transactionDirection: TransactionDirection.Send,
                    }),
                  );
                  setReceiptQuery(queryString);

                  setShowConfirmSendModal(false);
                  setShowSendSuccessModal(true);
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
          <div />
        </div>
      </Modal>
      <Modal
        showModal={showSendSuccessModal}
        onClose={() => {
          setShowSendSuccessModal(false);
        }}
      >
        <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
          <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
            <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
              {!activeWallet.sharedAccessEnabled
                ? 'Your transaction was successful!'
                : 'Your request was successful'}
            </p>
          </div>
          <img src="/images/success.png" alt="success" />
          {!activeWallet.sharedAccessEnabled ? (
            <>
              <div className="px-5 w-full bg-primary-100 rounded-xl py-5 space-y-2">
                <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                  Sent to
                </p>
                <div className="flex space-x-3 justify-start items-center">
                  <img
                    src={formData.transactionData?.destinationThumbnail}
                    className="rounded-full w-16 h-16"
                  />
                  <p>
                    {formData.transactionData?.destinationFirstName}{' '}
                    {formData.transactionData?.destinationLastName}
                  </p>
                </div>
                {formData.memo && (
                  <>
                    <hr className="border-1" />
                    <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                      For
                    </p>
                    <p>{formData.memo}</p>
                  </>
                )}
                <hr className="border-1" />
                <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                  Blockchain Proof (Transaction ID)
                </p>
                <div className="flex space-x-4">
                  <a
                    className="underline"
                    href={`${getExplorerBaseUrl(appState.walletMode)}${
                      formData.transactionData?.transactionId
                    }`}
                    target="_blank"
                  >
                    {formData.transactionData?.transactionId}
                  </a>
                  <button
                    type="button"
                    onClick={() =>
                      navigator.clipboard
                        .writeText(formData.transactionData?.transactionId)
                        .then(() => {
                          showNotification('info', 'Username copied!');
                        })
                    }
                  >
                    <img src="/images/copy.png" alt="copy" />
                  </button>
                </div>
              </div>
              <div className="w-full space-y-3">
                <a
                  href={`/#/send-asset-receipt?q=${receiptQuery}`} // since we use hash router
                  target="_blank"
                  className="bg-primary-800 block rounded-lg w-full text-center text-white h-12 py-4 px-5"
                >
                  Generate receipt
                </a>
                <ButtonSecondary
                  label="Close"
                  additionalClasses="font-montserratSemiBold"
                  onclick={async () => {
                    resetForm();
                    setShowSendSuccessModal(false);
                  }}
                />
              </div>
            </>
          ) : (
            <>
              <div className="px-5 w-full bg-primary-100 rounded-xl py-5 space-y-2">
                <p className="text-center font-montserratSemiBold">
                  Payment request submitted
                </p>
                <p className="text-center">
                  You have successfully requested payment of [{formData.amount}{' '}
                  {getAssetCode(selectedAsset.assetCode)}] from [
                  {activeWallet.alias}] to [{formData.sendTo}]. This transaction
                  will be completed when it gets the required number of
                  approvals by those who have approver access on this wallet.
                </p>
              </div>
              <div className="w-full space-y-3">
                <Button
                  label="Done"
                  additionalClasses="font-montserratSemiBold"
                  onclick={async () => {
                    resetForm();
                    setShowSendSuccessModal(false);
                  }}
                />
              </div>
            </>
          )}
          <div />
        </div>
      </Modal>
      <Modal
        showModal={showConfirmSwapModal}
        onClose={() => {
          setShowConfirmSwapModal(false);
          setShowSwapSuccess(false);
        }}
      >
        {!showSwapSuccess ? (
          <div className="flex flex-col items-center overflow-y-auto w-full px-5 md:px-20 space-y-5 py-5 justify-center rounded-md">
            <p className="text-primary-800 w-full text-lg md:text-xl font-montserratSemiBold">
              Confirm swap
            </p>
            {formData.transactionData?.messages.length > 0 && (
              <div className="flex flex-col text-center">
                <p className="font-montserratSemiBold">Note:</p>
                {formData.transactionData?.messages.map((message: string) => (
                  <p className="text-red-500">{message}</p>
                ))}
              </div>
            )}
            <p>You are swapping</p>
            <div className="w-full bg-primary-100 rounded-xl py-20 space-y-2 text-center">
              <p className="font-montserratSemiBold">
                {formData.amount} {getAssetCode(selectedAsset.assetCode)}
              </p>
              <p>
                {formData.swapFrom &&
                  `${calculateFiatValue(
                    Number(formData.amount || 0),
                    (fiatRates as Record<string, number>)[
                      appUser.currency.toUpperCase()
                    ] || 0,
                    formData.swapFrom?.usdPrice,
                  )} ${appUser.currency.toUpperCase()}`}
              </p>
            </div>
            <p>To</p>
            <div className="flex flex-col space-y-3 justify-center items-center px-10 w-full bg-primary-100 rounded-xl py-5 space-y-2">
              <p className="font-montserratSemiBold">
                {formData.transactionData?.swappedEstimate}{' '}
                {getAssetCode(formData.transactionData?.destinationAssetCode)}
              </p>
              <p>
                {formData.swapTo &&
                  `${formatToDecimal(
                    calculateFiatValue(
                      Number(formData.amount || 0),
                      (fiatRates as Record<string, number>)[
                        appUser.currency.toUpperCase()
                      ] || 0,
                      formData.swapTo?.usdPrice,
                    ),
                  )} ${appUser.currency.toUpperCase()}`}
              </p>
            </div>
            {formData.transactionData?.fee && (
              <>
                <p>Service Fee</p>
                <div className="flex flex-col space-y-3 justify-center items-center px-10 w-full bg-primary-100 rounded-xl py-5 space-y-2">
                  <p className="space-x-1">
                    <span className="font-montserratSemiBold">Fee:</span>
                    <span>{formData.transactionData?.fee}%</span>
                  </p>
                  <p className="space-x-1">
                    <span className="font-montserratSemiBold">
                      Amount (Calculated):
                    </span>
                    <span>
                      {formData.transactionData?.feeAmount}{' '}
                      {getAssetCode(selectedAsset.assetCode)}
                    </span>
                  </p>
                  {formData.transactionData.vatAmount > 0 && (
                    <>
                      <p className="space-x-1">
                        <span className="font-montserratSemiBold">VAT:</span>
                        <span>{formData.transactionData?.vat}%</span>
                      </p>
                      <p className="space-x-1">
                        <span className="font-montserratSemiBold">
                          VAT Amount:
                        </span>
                        <span>
                          {formData.transactionData?.vatAmount}{' '}
                          {getAssetCode(selectedAsset.assetCode)}
                        </span>
                      </p>
                    </>
                  )}
                </div>
              </>
            )}
            <p>Wallet</p>
            <div className="flex space-x-3 justify-center items-center px-10 w-full bg-primary-100 rounded-xl py-5 space-y-2">
              <span className="font-montserratSemiBold">
                {activeWallet.alias}
              </span>
            </div>
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
                    isSharedWallet: activeWallet.sharedAccessEnabled,
                    ...formData.transactionData,
                    commit: activeWallet.sharedAccessEnabled ? 1 : 0,
                    transactionSignature: signBase64Txn(
                      secretKey,
                      formData.transactionData?.transaction,
                      formData.transactionData?.networkPassPhrase,
                    ),
                  };

                  const payload = {
                    signer: activeWallet.signer,
                    publicKey: activeWallet.publicKey,
                    secretKey: secretKey,
                    body,
                  };
                  console.log('body', payload);

                  toggleLoader();
                  const res = await swapAsset(payload);
                  console.log('response', res);

                  if ('data' in res) {
                    console.log('response', res);
                    setFormData({
                      ...formData,
                      transactionData: res.data,
                    });

                    setShowSwapSuccess(true);
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
            <div />
          </div>
        ) : (
          <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
            <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
              <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                {!activeWallet.sharedAccessEnabled
                  ? 'Your transaction was successful!'
                  : 'Your request was successful'}
              </p>
            </div>
            <img src="/images/success.png" alt="success" />
            {!activeWallet.sharedAccessEnabled ? (
              <>
                <div className="px-5 w-full bg-primary-100 rounded-xl py-5 space-y-2">
                  <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                    Swapped
                  </p>
                  {formData.swapFrom && formData.amount && (
                    <div className="flex flex-col">
                      <p>
                        {formatToDecimal(Number(formData.amount))}{' '}
                        {getAssetCode(formData.swapFrom?.assetCode)}
                      </p>
                      <p>
                        -{' '}
                        {`${calculateFiatValue(
                          Number(formData.amount || 0),
                          (fiatRates as Record<string, number>)[
                            appUser.currency.toUpperCase()
                          ] || 0,
                          formData.swapFrom.usdPrice,
                        )} ${appUser.currency.toUpperCase()}`}
                      </p>
                    </div>
                  )}
                  {formData.swapTo && (
                    <>
                      <hr className="border-1" />
                      <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                        To
                      </p>
                      <p>
                        {formatToDecimal(Number(formData.amount))}{' '}
                        {getAssetCode(formData.swapTo?.assetCode)}
                      </p>
                      <p>
                        +{' '}
                        {`${calculateFiatValue(
                          Number(formData.amount || 0),
                          (fiatRates as Record<string, number>)[
                            appUser.currency.toUpperCase()
                          ] || 0,
                          formData.swapTo.usdPrice,
                        )} ${appUser.currency.toUpperCase()}`}
                      </p>
                      <hr className="border-1" />
                    </>
                  )}
                  {formData.transactionData?.fee && (
                    <>
                      <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                        Fee
                      </p>
                      <p>
                        {formData.transactionData?.feeAmount}{' '}
                        {getAssetCode(selectedAsset.assetCode)} (
                        {formData.transactionData?.fee}%)
                      </p>
                      <hr className="border-1" />
                      {formData.transactionData.vatAmount > 0 && (
                        <>
                          <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                            VAT
                          </p>
                          <p>
                            {formData.transactionData?.vatAmount}{' '}
                            {getAssetCode(selectedAsset.assetCode)} (
                            {formData.transactionData?.vat}%)
                          </p>
                          <hr className="border-1" />
                        </>
                      )}
                    </>
                  )}
                  <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                    Blockchain Proof (Transaction ID)
                  </p>
                  <div className="flex space-x-4">
                    <a
                      className="underline"
                      href={`${getExplorerBaseUrl(appState.walletMode)}${
                        formData.transactionData?.transactionId
                      }`}
                      target="_blank"
                    >
                      {formData.transactionData?.transactionId}
                    </a>
                    <button
                      type="button"
                      onClick={() =>
                        navigator.clipboard
                          .writeText(formData.transactionData?.transactionId)
                          .then(() => {
                            showNotification('info', 'Username copied!');
                          })
                      }
                    >
                      <img src="/images/copy.png" alt="copy" />
                    </button>
                  </div>
                </div>
                <div className="w-full space-y-3">
                  <Button
                    label="Done"
                    additionalClasses="font-montserratSemiBold"
                    onclick={async () => {
                      resetForm();
                      setShowConfirmSwapModal(false);
                      setShowSwapSuccess(false);
                    }}
                  />
                </div>
              </>
            ) : (
              <>
                <div className="px-5 w-full bg-primary-100 rounded-xl py-5 space-y-2">
                  <p className="text-center font-montserratSemiBold">
                    Swap request submitted
                  </p>
                  <p className="text-center">
                    You have successfully requested swap of [{formData.amount}{' '}
                    {getAssetCode(selectedAsset.assetCode)}] to [
                    {formData.transactionData?.swappedEstimate}
                    {formData.transactionData?.destinationAssetCode}] on wallet
                    [{activeWallet.alias}]. This transaction will be completed
                    when it gets the required number of approvals by those who
                    have approver access on this wallet.
                  </p>
                </div>
                <div className="w-full space-y-3">
                  <Button
                    label="Done"
                    additionalClasses="font-montserratSemiBold"
                    onclick={async () => {
                      resetForm();
                      setShowConfirmSwapModal(false);
                      setShowSwapSuccess(false);
                    }}
                  />
                </div>
              </>
            )}
            <div />
          </div>
        )}
      </Modal>
    </div>
  );
}
