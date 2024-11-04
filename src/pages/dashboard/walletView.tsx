import Button from '../../components/button';
import Header from '../../components/header';
import Tabs from '../../components/tabs';
import { useEffect, useRef, useState } from 'react';
import WalletCard from '../../components/walletCard';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/reduxStore';
import React from 'react';
import {
  calculateFiatValue,
  formatToDecimal,
  getAssetCode,
  getBytesLength,
  getExplorerBaseUrl,
  totalWalletBalanceInCurrency,
} from '../../utils/utilities';
import { Wallet } from '../../types/wallet';
import AssetTokenItem from '../../components/assetTokenListItem';
import {
  TEDropdown,
  TEDropdownItem,
  TEDropdownMenu,
  TEDropdownToggle,
  TERipple,
} from 'tw-elements-react';
import AssetListItem from '../../components/assetListItem';
import Dropdown from '../../components/dropdown';
import TextInput from '../../components/textInput';
import { Asset } from '../../types/asset';
import { showNotification, toggleLoader } from '../../utils/showToaster';
import { useSendAssetMutation } from '../../store/api/walletApis';
import { ErrorResponse } from '../../store/api/baseapi/axiosBaseQuery';
import Modal from '../../components/modal';
import { Encryptor } from '../../utils/encryptor';
import { signBase64Txn } from '../../utils/trovoSDK';
import ButtonSecondary from '../../components/buttonSecondary';

export default function WalletView() {
  const ref = useRef<HTMLDivElement>(null);
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const appState = useSelector((state: RootState) => state.appState!);
  let mutableWalletArray = [...appUser.userWallets];
  const [wallets, setWallets] = useState(
    mutableWalletArray.sort((w) => (w.primaryWallet ? 0 : 1)),
  );
  const [activeWalletIndex, setActiveWalletIndex] = useState(0);
  const [activeWallet, setActiveWallet] = useState<Wallet>(wallets[0]);
  const [assetFilterMode, setAssetFilterMode] = useState('Asset Tokens');
  const [walletActionMode, setWalletActionMode] = useState(0);
  const itemRefs = useRef<HTMLDivElement[]>([]);
  const fiatRates = useSelector((state: RootState) => state.cache.fiatRates);
  const [secretKey, setSecretKey] = useState('');
  const [password, setPassword] = useState('');
  const [passwordErr, setPasswordErr] = useState('');
  const [showConfirmSendModal, setShowConfirmSendModal] = useState(false);
  const [showSendSuccessModal, setShowSendSuccessModal] = useState(false);
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
    console.log('fjsdlkfs', Number(formData.amount));

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
    } else if (Number(formData.amount) <= 0) {
      newObj = {
        ...newObj,
        amount: 'Please enter a valid amount',
      };
    } else if (Number(formData.amount) > selectedAsset?.amount!) {
      newObj = {
        ...newObj,
        amount: "You don't have sufficient balance",
      };
    } else if (
      getAssetCode(selectedAsset!.assetCode) == 'XBN' &&
      Number(formData.amount) > selectedAsset!.amount! - 7
    ) {
      newObj = {
        ...newObj,
        amount: "You don't have sufficient balance",
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

  const validateRecieveForm = (): boolean => {
    let isValid = true;
    let newObj = errorObj;

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
    } else {
      newObj = {
        ...newObj,
        amount: '',
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
  }, [activeWallet]);

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

    console.log('secret key here ', secretKey);

    const payload = {
      signer: activeWallet.publicKey,
      publicKey: activeWallet.publicKey,
      secretKey: secretKey,
      body: {
        // isSharedWallet: activeWallet.sharedAccessEnabled,
        destination: formData.sendTo,
        memo: formData.memo,
        amount: formData.amount,
        assetCode: getAssetCode(selectedAsset.assetCode),
        assetIssuer: selectedAsset.assetIssuer,
      },
    };

    console.log('secret key 2 here ', payload);
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
      console.log('sdafsdfsdsfdsf1111', key);
      setSecretKey(key);
      return true;
    } catch (error: any) {
      return false;
    }
  };

  const handleRecieve = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    validateRecieveForm();
  };

  const handleSwap = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    validateSwapForm();
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
      key={index}
      ref={(el) => (itemRefs.current[index] = el!)}
    />
  ));

  const dropdownItems = [
    {
      text: 'Asset Tokens',
      value: 'Asset Tokens',
    },
    {
      text: 'Other Tokens',
      value: 'Other Tokens',
    },
  ].map((item) => (
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
    <div className="flex text-primary-800 text-sm overflow-x-hidden md:text-md flex-col space-y-5 p-3">
      <Header isHomeView />
      <div className="flex md:h-screen w-full items-center justify-center">
        <div className="h-full w-full md:p-3">
          <div className="flex md:pt-5 p-5 rounded-lg flex-col items-center bg-primary-100 max-w-5xl min-h-full">
            <div className="w-full flex justify-between space-x-5">
              <Tabs
                tabList={['All Wallets', 'My Wallets', 'Shared Wallets']}
                onTabChanged={(index) => {
                  console.log('tab index', index);
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
                  {assetFilterMode}
                </span>
              </p>
            </div>
            <div className="w-full mt-4 space-y-4">
              {activeWallet.claimedAssets.length > 1 ? ( // since every wallet must have an XBN asset check length greater than 1
                // if filter is set on other tokens
                assetFilterMode === 'Other Tokens' ? (
                  activeWallet.claimedAssets.map(
                    (asset) =>
                      asset.assetCode !== '' && (
                        <AssetTokenItem
                          image={asset.imageUrl}
                          assetCode={asset.assetCode}
                          usdPrice={asset.usdPrice.toString()}
                          amount={asset.amount.toString()}
                          valueInFiat={`${calculateFiatValue(
                            asset.amount,
                            (fiatRates as Record<string, number>)[
                              appUser.currency.toUpperCase()
                            ] || 0,
                            asset.usdPrice,
                          )} ${appUser.currency.toUpperCase()}`}
                          onclick={() => {
                            // navigate('/dashboard/tokenized-asset');
                          }}
                        />
                      ),
                  )
                ) : (
                  activeWallet.claimedAssets.map(
                    (asset, index) =>
                      asset.assetCode !== '' && (
                        <AssetListItem
                          key={index}
                          image="/images/avatar.png"
                          assetName="Atlantis 1"
                          assetClass="Real Estate"
                          isSubscribed
                          onclick={() => {
                            // navigate('/dashboard/tokenized-asset');
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
                className={`flex space-y-3 ring-1 rounded-full px-6 py-2 flex-col items-center ${
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
                        <Dropdown
                          label={activeWallet.alias}
                          options={[
                            ...wallets.map((w, index) => ({
                              text: w.alias,
                              value: index,
                            })),
                          ]}
                          onSelect={(selectedItem) => {
                            console.log(selectedItem);
                            setActiveWallet(wallets[selectedItem.value]);
                            setActiveWalletIndex(selectedItem.value);
                            const ref = itemRefs.current[selectedItem.value];
                            if (ref) {
                              // Change background color
                              ref.scrollIntoView({
                                behavior: 'smooth',
                                block: 'nearest',
                                inline: 'center',
                              });
                            }
                          }}
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
                      <div className="w-full space-y-5">
                        <div className="space-y-3 w-full">
                          <p>Receiving wallet</p>
                          <Dropdown
                            label={activeWallet.alias}
                            options={[
                              ...wallets.map((w) => ({
                                text: w.alias,
                                value: w.publicKey,
                              })),
                            ]}
                            onSelect={() => {
                              //
                            }}
                          />
                        </div>
                        <div className="w-full space-y-2">
                          <div className="flex items-center space-x-3 space-y-3 w-full">
                            <div className="w-4/6">
                              <TextInput
                                inputType="text"
                                label="Token quantity/Amount"
                                onInputChange={(value) => {
                                  setFormData({ ...formData, sendTo: value });
                                }}
                              />
                            </div>
                            <div className="w-2/6 self-end">
                              <Dropdown
                                label={
                                  activeWallet.claimedAssets[0].assetCode === ''
                                    ? 'XBN'
                                    : activeWallet.claimedAssets[0].assetCode
                                }
                                options={[
                                  ...activeWallet.claimedAssets.map((w) => ({
                                    text:
                                      w.assetCode === '' ? 'XBN' : w.assetCode,
                                    value: w.assetCode,
                                  })),
                                ]}
                                onSelect={() => {
                                  //
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
                          <Dropdown
                            label={activeWallet.alias}
                            options={[
                              ...wallets.map((w) => ({
                                text: w.alias,
                                value: w.publicKey,
                              })),
                            ]}
                            onSelect={() => {
                              //
                            }}
                          />
                        </div>
                        <div className="w-full space-y-2">
                          <p>Swap from</p>
                          <Dropdown
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
                  // isSharedWallet: activeWallet.sharedAccessEnabled,
                  ...formData.transactionData,
                  commit: 1,
                  transactionSignature: signBase64Txn(
                    secretKey,
                    formData.transactionData.transaction,
                    formData.transactionData.networkPassPhrase,
                  ),
                };

                const payload = {
                  signer: activeWallet.publicKey,
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
              Your transaction was successful!
            </p>
          </div>
          <img src="/images/success.png" alt="success" />
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
            <Button
              label="Generate receipt"
              additionalClasses="font-montserratSemiBold"
              onclick={async () => {
                // setShowConfirmSendModal(false);
              }}
            />
            <ButtonSecondary
              label="Close"
              additionalClasses="font-montserratSemiBold"
              onclick={async () => {
                resetForm();
                setShowSendSuccessModal(false);
              }}
            />
          </div>
          <div />
        </div>
      </Modal>
    </div>
  );
}
