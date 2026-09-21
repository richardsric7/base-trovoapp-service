import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import TextInput from '../../components/textInput';
import WidgetCard from '../../components/widgetCard';
import AssetListItem from '../../components/assetListItem';
import TransactionItem from '../../components/transactionItem';
import Button from '../../components/button';
import ButtonSecondary from '../../components/buttonSecondary';
import Header from '../../components/header';
import Modal from '../../components/modal';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../store/reduxStore';
import capitalizeFirstLetter from '../../utils/capitalizeFirst';
import {
  totalAccountBalanceInCurrency,
  formatToDecimal,
  getBytesLength,
} from '../../utils/utilities';
import { Encryptor } from '../../utils/encryptor';
import { showNotification, toggleLoader } from '../../utils/showToaster';
import { ErrorResponse } from '../../store/api/baseapi/axiosBaseQuery';
import {
  useLazyReceiveAssetQuery,
  useLazyFetchFiatAmountForActivationQuery,
  useFetchFiatPaymentsQuery,
  useFetchTokenizedAssetsQuery,
  useFetchExpressedInterestsQuery,
  useFetchTokenizationDataQuery,
} from '../../store/api/walletApis';
import { useFlutterwave, closePaymentModal } from 'flutterwave-react-v3';
import {
  formatAmount,
  formatRelativeTime,
  getTransactionAssetCode,
  normalizeTransactionType,
  parseAmountValue,
  parseDate,
} from '../../utils/transactionUtils';
import { truncateAddress } from '../../utils/truncateValues';
import { TokenizedAsset } from '../../types/tokenizedAsset';
import {
  setActiveTokenizedAsset,
  setTokenizationData,
} from '../../store/appStateSlice';
import { TokenizationData } from '../../types/tokenizationData';
import { ExpressedInterest } from '../../types/expressedInterest';
import {
  BuyTokenModal,
  ExpressInterestModal,
  KycUnverifiedModal,
  P2pComingSoonModal,
  SoldOutModal,
} from '../../components/tokenizedAssetActionModals';

export default function Home() {
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const fiatRates = useSelector((state: RootState) => state.cache.fiatRates);
  const dispatch = useDispatch();

  const primaryWallet = appUser.userWallets.find((w) => w.primaryWallet)!;
  const navigate = useNavigate();
  const [isActivated] = useState(
    Number(
      primaryWallet?.claimedAssets.find((a) => !a.assetCode && !a.assetIssuer)
        ?.amount,
    ) !== 0,
  );
  const [hasAssets, setHasAssets] = useState(
    primaryWallet.claimedAssets.length > 0,
  );

  const config = {
    public_key: 'FLWPUBK-**************************-X',
    tx_ref: Date.now().toString(),
    amount: 100,
    currency: 'NGN',
    payment_options: 'card,mobilemoney,ussd',
    customer: {
      email: 'user@gmail.com',
      phone_number: '070********',
      name: 'john doe',
    },
    customizations: {
      title: 'Activate Trovo Account',
      description: 'Buy ETH and TROV with fiat',
      logo: 'https://st2.depositphotos.com/4403291/7418/v/450/depositphotos_74189661-stock-illustration-online-shop-log.jpg',
    },
  };

  const handleFlutterPayment = useFlutterwave(config);

  const getImage = (patronPackageId: string): string => {
    switch (patronPackageId.toLowerCase()) {
      case 'platinum':
        return '/images/platinum.png';
      case 'diamond':
        return '/images/diamond.png';
      default:
        return '/images/gold.png';
    }
  };

  const [showSubscribeModal, setShowSubscribeModal] = useState(false);
  const [showSoldOutModal, setShowSoldOutModal] = useState(false);
  const [showKycUnverifiedModal, setShowKycUnverifiedModal] = useState(false);
  const [showBuyTokenModal, setShowBuyTokenModal] = useState(false);
  const [showP2pComingSoonModal, setShowP2pComingSoonModal] = useState(false);
  const [showCopiedModal, setShowCopiedModal] = useState(false);
  const [showComingSoonModal, setShowComingSoonModal] = useState(false);
  const [showRequestETHModal, setShowRequestETHModal] = useState(false);
  const [showBuyETHWithFiatModal, setShowBuyETHWithFiatModal] = useState(false);
  const [activationAmount, setActivationAmount] = useState(0);
  const [gasPercent, setGasPercent] = useState(0);
  const [confirmRequestETHModal, setConfirmRequestETHModal] = useState(false);
  const [qrCodeLink, setQrCodeLink] = useState('');
  const [receiveAsset] = useLazyReceiveAssetQuery();
  const [fetchFiatAmounts] = useLazyFetchFiatAmountForActivationQuery();
  const [formData, setFormData] = useState({
    sendTo: '',
    amount: '',
    memo: '',
  });
  const [errorObj, setErrorObj] = useState({
    amount: '',
    memo: '',
  });
  const [secretKey, setSecretKey] = useState('');
  const [asset, setAsset] = useState(
    useSelector((state: RootState) => state.appState.activeTokenizedAsset),
  );

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

    if (formData.memo && getBytesLength(formData.memo) > 60) {
      newObj = {
        ...newObj,
        memo: 'Memo length cannot be more than 60 bytes',
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

  useEffect(() => {
    (async () => {
      const encryptor = new Encryptor();
      setSecretKey(await encryptor.getSecretKey(appUser));
    })();
  }, []);

  const handleRecieve = async () => {
    if (!validateRecieveForm()) {
      return;
    }

    const payload = {
      signer: primaryWallet.signer,
      address: primaryWallet.address,
      secretKey: secretKey,
      body: {
        destination: formData.sendTo,
        memo: formData.memo,
        amount: formData.amount,
        address: primaryWallet.address,
        alias: primaryWallet.alias,
        assetCode: '',
        assetIssuer: '',
      },
    };

    toggleLoader();
    const res = await receiveAsset(payload);
    console.log('response', res);

    if ('data' in res) {
      console.log('response', res);
      setQrCodeLink(res.data.qrCode);
      setConfirmRequestETHModal(true);
      setShowRequestETHModal(false);
    } else if ('error' in res) {
      const errorResponse = res.error as ErrorResponse;
      showNotification(
        'error',
        errorResponse.data.message ?? 'Something went wrong. Please try again.',
      );
    }
    toggleLoader();
  };

  const handleBuyWithFiat = async () => {
    const payload = {
      signer: primaryWallet.signer,
      address: primaryWallet.address,
      secretKey: secretKey,
      body: {},
    };

    toggleLoader();
    const res = await fetchFiatAmounts(payload);
    if ('data' in res) {
      console.log('response', res.data);
      setActivationAmount(res.data.activationAmount);
      setGasPercent(res.data.gasPercent);
      setShowBuyETHWithFiatModal(true);
    } else if ('error' in res) {
      const errorResponse = res.error as ErrorResponse;
      showNotification(
        'error',
        errorResponse.data.message ?? 'Something went wrong. Please try again.',
      );
    }
    toggleLoader();
  };

  const getPercentageValue = (percentage: number, amount: number): number =>
    (percentage * amount) / 100;

  const { data, isLoading } = useFetchFiatPaymentsQuery(
    {
      signer: primaryWallet?.signer ?? '',
      address: primaryWallet.address,
      secretKey,
      body: { limit: 5 },
    },
    { skip: !secretKey || !primaryWallet },
  );

  const rawRecords = data?.data?.records ?? data?.records ?? [];

  const rows: HistoryRow[] = rawRecords.map(
    (item: Record<string, any>, index: number) => {
      const type = normalizeTransactionType(
        item.transactionType ?? item.type,
        item,
        primaryWallet.address,
      );
      const amountValue = parseAmountValue(item.amount);
      const assetCode = getTransactionAssetCode(item);
      const amountPrefix =
        type === 'Sent' ? '-' : type === 'Received' ? '+' : '+';
      const amountText =
        typeof item.amount === 'string' && item.amount.trim().length > 0
          ? item.amount
          : formatAmount(amountValue, assetCode);
      const createdAt = parseDate(
        item.transactionDate ?? item.createdAt ?? item.date,
      );
      const description =
        item.description ??
        item.narration ??
        item.memo ??
        (type === 'Received'
          ? `Received ${assetCode}`
          : type === 'Sent'
            ? `Sent ${assetCode}`
            : `Swapped asset to ${assetCode}`);

      const walletMatch =
        appUser.userWallets.find(
          (wallet) =>
            wallet.address === item.address ||
            wallet.address === item.walletAddress ||
            wallet.alias === item.walletAlias,
        ) ?? primaryWallet;

      return {
        id: `${item.transactionId ?? item.id ?? item._id ?? item.reference ?? index}`,
        type,
        assetCode,
        amountValue,
        price: `${amountPrefix}${amountText}`,
        description,
        dateLabel: formatRelativeTime(createdAt),
        dateValue: createdAt,
        walletKey: walletMatch?.address ?? 'unknown',
        walletLabel: walletMatch?.alias ?? 'Primary wallet',
        username: `${item.username ?? item.fullName ?? item.name ?? ''}`.trim(),
        fromUsername: item.from,
        toUsername: item.to,
        fromAddress: `${item.fromAddress ?? item.senderAddress ?? ''}`,
        toAddress: `${item.toAddress ?? item.receiverAddress ?? ''}`,
        memo: `${item.memo ?? item.narration ?? ''}`.trim(),
      };
    },
  );

  const { data: primaryOffers, isLoading: isLoadingPrimaryOffers } =
    useFetchTokenizedAssetsQuery(
      {
        signer: primaryWallet?.signer ?? '',
        address: primaryWallet.address,
        secretKey,
        body: { status: 0 },
      },
      { skip: !secretKey || !primaryWallet },
    );

  const { data: expressedInterests } = useFetchExpressedInterestsQuery(
    {
      signer: primaryWallet?.signer ?? '',
      address: primaryWallet.address,
      secretKey,
      body: {},
    },
    { skip: !secretKey || !primaryWallet },
  );

  var typedPrimaryOffers: TokenizedAsset[] = [];

  const { data: secondaryListing, isLoading: isLoadingSecondaryListing } =
    useFetchTokenizedAssetsQuery(
      {
        signer: primaryWallet?.signer ?? '',
        address: primaryWallet.address,
        secretKey,
        body: { status: 1 },
      },
      { skip: !secretKey || !primaryWallet },
    );

  var typedSecondaryListing =
    secondaryListing?.records == undefined
      ? []
      : (secondaryListing?.records as TokenizedAsset[]);

  primaryOffers?.records.map((x: TokenizedAsset) => {
    var record = expressedInterests?.records?.find(
      (a: ExpressedInterest) => a.assetCode == x.assetCode,
    );
    var y = {
      ...x,
      expressedInterestAmount: record != null ? record.amount : 0,
    };
    typedPrimaryOffers = [...typedPrimaryOffers, y];
  });

  const { data: tokenizationData } = useFetchTokenizationDataQuery(
    {
      signer: primaryWallet?.signer ?? '',
      address: primaryWallet.address,
      secretKey,
      body: { limit: 5 },
    },
    { skip: !secretKey || !primaryWallet },
  );

  dispatch(
    setTokenizationData(
      tokenizationData != undefined
        ? (tokenizationData as TokenizationData)
        : undefined,
    ),
  );

  return (
    <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
      <Header isHomeView />
      <div className="flex md:h-full w-full items-center justify-center">
        <div className="h-full w-full md:p-3">
          {!appUser.hasSecurityQuestions && (
            <div
              className="bg-trovored-light flex justify-between mb-2 rounded-md ring-1 ring-trovored-primary
            text-trovored-primary px-4 py-2 font-semibold"
            >
              <div className="flex space-x-5 items-center">
                <Link
                  className="flex space-x-5 items-center"
                  to={'setup-security-questions'}
                >
                  <img src="/images/alert.png" alt="" />
                  <p className="text-xs md:text-md">
                    You have not setup security questions yet. Tap to setup
                    security questions
                  </p>
                </Link>
              </div>
            </div>
          )}
          <div className="md:hidden items-center px-5 mb-5 justify-center flex">
            <WidgetCard
              localCurrencyBalance={formatToDecimal(
                totalAccountBalanceInCurrency(
                  appUser.userWallets,
                  (fiatRates as Record<string, number>)[
                    appUser.currency.toUpperCase()
                  ] || 0,
                ),
              )}
              usdBalance={formatToDecimal(
                totalAccountBalanceInCurrency(
                  appUser.userWallets,
                  (fiatRates as { USD: number }).USD,
                ),
              )}
              currency={appUser.currency.toUpperCase()}
            />
          </div>
          <div className="flex space-y-3 mb-20 md:pt-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
            <div className="hidden w-full items-center px-5 justify-center md:flex">
              <WidgetCard
                localCurrencyBalance={formatToDecimal(
                  totalAccountBalanceInCurrency(
                    appUser.userWallets,
                    (fiatRates as Record<string, number>)[
                      appUser.currency.toUpperCase()
                    ] || 0,
                  ),
                )}
                usdBalance={formatToDecimal(
                  totalAccountBalanceInCurrency(
                    appUser.userWallets,
                    (fiatRates as { USD: number }).USD,
                  ),
                )}
                currency={appUser.currency.toUpperCase()}
              />
            </div>
            {isActivated ? (
              <div className="w-full">
                <p
                  onClick={() => {
                    navigator.clipboard
                      .writeText(JSON.stringify(tokenizationData))
                      .then(() => {
                        showNotification('success', 'Address copied!');
                      });
                  }}
                  className="text-primary-800 text-center text-lg md:text-xl font-montserratSemiBold"
                >
                  Tokenized Assets
                </p>
                <div className="md:px-5 pb-10 w-full space-y-5">
                  <div className="flex items-center w-full justify-between">
                    <p className="font-montserratSemiBold text-md">
                      PRIMARY OFFERS
                    </p>
                    <Link
                      className="flex space-x-5 items-center"
                      to={'tokenized-asset-list?rel=primary'}
                    >
                      <p className="test-xs text-primary-400 underline">
                        View all
                      </p>
                    </Link>
                  </div>
                  {isLoadingPrimaryOffers ? (
                    <p className="text-center">Loading primary offers...</p>
                  ) : typedPrimaryOffers.length == 0 ? (
                    <div className="w-full">
                      <div className="md:px-5 w-full">
                        <div className="w-full flex flex-col justify-center space-y-5 md:mb-10">
                          <img
                            className="self-center"
                            src="/images/empty.png"
                            alt=""
                          />
                          <p className="font-montserratSemiBold text-md text-center">
                            No Assets Yet
                          </p>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div className="space-y-5">
                      {typedPrimaryOffers.map((asset, index) => {
                        return (
                          <AssetListItem
                            key={index}
                            asset={asset}
                            isSubscribed={asset.expressedInterestAmount > 0}
                            onclick={() => {
                              dispatch(setActiveTokenizedAsset(asset));
                              navigate(
                                `/dashboard/tokenized-asset?salesList=0&assetCode=${asset.assetCode}`,
                              );
                            }}
                            onBuy={() => {
                              setAsset(asset);
                              setShowBuyTokenModal(true);
                            }}
                            onExpressInterest={() => {
                              setAsset(asset);
                              setShowSubscribeModal(true);
                            }}
                          />
                        );
                      })}
                    </div>
                  )}
                </div>
                <div className="md:px-5 pb-10 w-full space-y-5">
                  <div className="flex items-center w-full justify-between">
                    <p className="font-montserratSemiBold text-md">
                      SECONDARY LISTING
                    </p>
                    <Link
                      className="flex space-x-5 items-center"
                      to={'tokenized-asset-list?rel=secondary'}
                    >
                      <p className="test-xs text-primary-400 underline">
                        View all
                      </p>
                    </Link>
                  </div>
                  {isLoadingSecondaryListing ? (
                    <p className="text-center">Loading secondary listing...</p>
                  ) : typedSecondaryListing.length == 0 ? (
                    <div className="w-full">
                      <div className="md:px-5 w-full">
                        <div className="w-full flex flex-col justify-center space-y-5 md:mb-10">
                          <img
                            className="self-center"
                            src="/images/empty.png"
                            alt=""
                          />
                          <p className="font-montserratSemiBold text-md text-center">
                            No Assets Yet
                          </p>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div className="space-y-5">
                      {typedSecondaryListing.map((asset, index) => (
                        <AssetListItem
                          key={index}
                          asset={asset}
                          isSecondary={true}
                          onclick={() => {
                            dispatch(setActiveTokenizedAsset(asset));
                            navigate(
                              `/dashboard/tokenized-asset?salesList=1&assetCode=${asset.assetCode}`,
                            );
                          }}
                          onBuy={() => {
                            setAsset(asset);
                            setShowBuyTokenModal(true);
                          }}
                          onExpressInterest={() => {
                            setAsset(asset);
                            setShowSubscribeModal(true);
                          }}
                        />
                      ))}
                    </div>
                  )}
                </div>
              </div>
            ) : (
              <div className="w-full">
                <p className="text-primary-800 text-center text-lg md:text-xl font-montserratSemiBold">
                  Wallet Activation
                </p>
                <div className="md:px-5 w-full">
                  <div className="w-full flex flex-col items-center justify-center space-y-5 md:mb-10">
                    <img
                      className="self-center"
                      src="/images/activateWallet.png"
                      alt=""
                    />
                    <p className="font-montserratSemiBold text-md text-center">
                      Your Wallet is ready!
                    </p>
                    <div className="space-y-2">
                      <p className="text-md text-center">
                        You'll need at least 10 ETH in your wallet before you
                        can start transacting
                      </p>
                      <p className="text-md text-center">
                        You can add ETH to your wallet using either of the 3
                        easy ways displayed below
                      </p>
                    </div>
                    <div className="px-3 pb-5 space-y-3">
                      <div className="flex space-x-5 ">
                        <Button
                          label="Buy ETH with Fiat"
                          additionalClasses="font-montserratSemiBold py-10 px-5"
                          onclick={handleBuyWithFiat}
                          leftIcon={
                            <img
                              src="/images/buy_with_fiat.png"
                              height={40}
                              width={40}
                              alt="filter"
                            />
                          }
                        />
                        <Button
                          label="Buy ETH on TrovoP2P"
                          additionalClasses="font-montserratSemiBold py-10 px-5"
                          onclick={() => {
                            // navigate('/import');
                            setShowComingSoonModal(true);
                          }}
                          leftIcon={
                            <img
                              src="/images/buy_from_trovop2p.png"
                              height={40}
                              width={40}
                              alt="Buy ETH on TrovoP2P"
                            />
                          }
                        />
                      </div>
                      <div className="flex space-x-5 ">
                        <Button
                          label="Request ETH from Trovo User"
                          additionalClasses="font-montserratSemiBold py-10 px-5"
                          onclick={() => {
                            setShowRequestETHModal(true);
                          }}
                          leftIcon={
                            <img
                              src="/images/request_from_user.png"
                              height={40}
                              width={40}
                              alt="Request from Trovo User"
                            />
                          }
                        />
                        <Button
                          label="Send ETH to your Wallet"
                          additionalClasses="font-montserratSemiBold py-10 px-5"
                          onclick={() => {
                            navigator.clipboard
                              .writeText(primaryWallet.alias)
                              .then(() => {
                                showNotification(
                                  'success',
                                  'Address copied!',
                                );
                              });
                          }}
                          leftIcon={
                            <img
                              src="/images/send_to_wallet.png"
                              height={40}
                              width={40}
                              alt="Send ETH to your Wallet"
                            />
                          }
                        />
                      </div>
                    </div>
                  </div>
                </div>
                {/* Address copied modal */}
                <Modal
                  showModal={showCopiedModal}
                  onClose={() => {
                    setShowCopiedModal(false);
                  }}
                >
                  <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
                    <img src="/images/success.png" alt="success" />
                    <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
                      <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                        Address Copied Successfully!
                      </p>
                    </div>
                  </div>
                </Modal>
                {/* P2P coming soon modal */}
                <Modal
                  showModal={showComingSoonModal}
                  onClose={() => {
                    setShowComingSoonModal(false);
                  }}
                >
                  <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
                    <p className="text-primary-800 text-md xl:text-xl font-montserratSemiBold">
                      Coming soon!
                    </p>
                    <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
                      <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
                        P2P will be launching soon.
                      </p>
                    </div>
                  </div>
                </Modal>
                {/* Request ETH modal */}
                <Modal
                  showModal={showRequestETHModal}
                  onClose={() => {
                    setShowRequestETHModal(false);
                  }}
                >
                  <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
                    <p className="text-primary-800 w-full mb-5 text-lg md:text-xl font-montserratSemiBold">
                      Request ETH
                    </p>
                    <div className="w-full">
                      <TextInput
                        defaultValue={primaryWallet.alias}
                        inputType="text"
                        readonly={true}
                        label="Receiving Wallet"
                        placeholder="Efizee"
                        onInputChange={() => {
                          // console.log('input has changed', newValue);
                        }}
                      />
                    </div>
                    <div className="w-full">
                      <TextInput
                        inputType="text"
                        label="Amount"
                        placeholder="100"
                        defaultValue={formData.amount}
                        onInputChange={(value) => {
                          setFormData({ ...formData, amount: value });
                        }}
                      />
                      <div
                        className={`w-full flex items-center ${
                          errorObj.amount ? 'justify-between' : 'justify-end'
                        }`}
                      >
                        {errorObj.amount && (
                          <span className="text-red-500 text-sm mt-1">
                            {errorObj.amount}
                          </span>
                        )}
                      </div>
                    </div>

                    <div className="space-y-3 w-full">
                      <TextInput
                        inputType="text"
                        label="Add memo (optional)"
                        defaultValue={formData.memo}
                        onInputChange={(value) => {
                          console.log('value', value);
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
                            formData.memo && getBytesLength(formData.memo) > 60
                              ? 'text-red-500'
                              : ''
                          }
                        >
                          {formData.memo ? getBytesLength(formData.memo) : 0}
                          /60
                        </span>
                      </div>
                    </div>
                    <div className="w-full md:w-2/4 space-y-3">
                      <Button
                        label="Proceed"
                        additionalClasses="font-montserratSemiBold"
                        onclick={handleRecieve}
                      />
                    </div>
                    <div />
                  </div>
                </Modal>
                <Modal
                  showModal={confirmRequestETHModal}
                  onClose={() => {
                    setConfirmRequestETHModal(false);
                  }}
                >
                  <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
                    <p className="text-primary-800 w-full mb-5 text-lg font-montserratSemiBold">
                      Receive {formData.amount} ETH
                    </p>
                    <div className="flex w-full py-3 px-2 xl:px-5 xl:space-y-5 rounded-xl flex-col items-center bg-primary-100">
                      <div className="flex flex-col space-y-2 w-full justify-between">
                        <p className="font-montserratSemiBold">
                          Receiving Wallet
                        </p>
                        <div className="flex w-full py-3  rounded-xl justify-start">
                          <div className="flex flex-col space-y-2 w-full justify-between">
                            <p className="font-montserratSemiBold">
                              {primaryWallet.alias}
                            </p>
                            <p className="text-xs text-primary-400">
                              {primaryWallet.address}
                            </p>
                          </div>
                          <button
                            type="button"
                            onClick={() => {
                              navigator.clipboard
                                .writeText(primaryWallet.address)
                                .then(() => {
                                  showNotification(
                                    'success',
                                    'Address copied!',
                                  );
                                });
                            }}
                          >
                            <img src="/images/copy.png" alt="copy" />
                          </button>
                        </div>
                      </div>
                    </div>
                    {formData.memo ? (
                      <div className="flex w-full py-5 px-5 rounded-xl mb-10 justify-center items-center bg-primary-100">
                        <div className="flex flex-col space-y-2 w-full justify-between">
                          <p className="font-montserratSemiBold">For</p>
                          <p className="font-montserratSemiBold">
                            {formData.memo}
                          </p>
                        </div>
                      </div>
                    ) : (
                      <></>
                    )}
                    <div className="flex w-full py-5 px-2 xl:px-10 xl:space-y-5 rounded-xl mb-10 justify-center items-center bg-primary-100">
                      <div className="flex flex-col space-y-2 items-center w-full justify-between">
                        <img className="w-3/5" src={qrCodeLink} alt="QR Code" />
                      </div>
                    </div>
                    <div className="w-full flex space-x-3 pt-5">
                      <Button
                        label="Share"
                        additionalClasses="font-montserratSemiBold"
                        onclick={() => {
                          setConfirmRequestETHModal(false);
                        }}
                      />
                      <ButtonSecondary
                        label="Back"
                        additionalClasses="font-montserratSemiBold"
                        onclick={() => {
                          setShowRequestETHModal(true);
                          setConfirmRequestETHModal(false);
                        }}
                      />
                    </div>
                    <div />
                  </div>
                </Modal>
                {/* Buy ETH with FIAT */}
                <Modal
                  showModal={showBuyETHWithFiatModal}
                  onClose={() => {
                    setShowBuyETHWithFiatModal(false);
                  }}
                >
                  <div className="flex flex-col items-center w-full px-5 md:px-20 space-y-5 py-5 justify-center">
                    <p className="text-primary-800 w-full mb-5 text-lg font-montserratSemiBold">
                      Activate Account
                    </p>
                    <div className="flex w-full py-10 px-2 xl:px-5 xl:space-y-5 rounded-xl flex-col items-center bg-primary-100">
                      <img
                        className="self-center"
                        src="/images/activateWallet.png"
                        alt=""
                      />
                      <div className="flex flex-col space-y-2 w-full justify-between">
                        <p className="text-center">
                          An easy option to quickly activate your account and
                          buy your first virtual asset on the Trovo App. You
                          will get ETH as well as TROV when you send fiat.
                        </p>
                      </div>
                    </div>
                    <p>You Pay</p>
                    <div className="flex w-full py-3 px-2 xl:px-5 xl:space-y-5 rounded-xl flex-col items-center bg-primary-100">
                      <div className="flex flex-col space-y-2 w-full justify-between">
                        <p className="font-montserratSemiBold text-center">
                          NGN {activationAmount}
                        </p>
                      </div>
                    </div>
                    <p>You Get</p>
                    <div className="flex w-full py-3 px-2 xl:px-5 xl:space-y-5 rounded-xl flex-col items-center bg-primary-100">
                      <div className="flex flex-col space-y-2 w-full justify-between">
                        <p className="font-montserratSemiBold text-center">
                          NGN {getPercentageValue(gasPercent, activationAmount)}{' '}
                          worth of Gas and NGN{' '}
                          {activationAmount -
                            getPercentageValue(
                              gasPercent,
                              activationAmount,
                            )}{' '}
                          worth of TROV
                        </p>
                      </div>
                    </div>
                    <div className="w-full flex space-x-3 pt-5">
                      <Button
                        label="Make Payment"
                        additionalClasses="font-montserratSemiBold"
                        onclick={() => {
                          // setConfirmRequestETHModal(false);
                          handleFlutterPayment({
                            callback: (response) => {
                              console.log(response);
                              closePaymentModal(); // this will close the modal programmatically
                            },
                            onClose: () => {},
                          });
                        }}
                      />
                    </div>
                    <div />
                  </div>
                </Modal>
              </div>
            )}
          </div>
        </div>
        <div className="hidden md:block w-3/6 flex flex-col space-y-5 h-full py-3 lg:px-3">
          <div className="flex space-y-3 rounded-lg py-5 flex-col items-center bg-primary-100">
            {appUser.patronMembership && (
              <>
                <p>Current Plan</p>
                <div className="flex space-x-1">
                  <img
                    src={getImage(
                      appUser?.patronMembership?.patronPackageId ?? '',
                    )}
                    alt="gold"
                    className="h-5 w-5"
                  />
                  <p className="font-bold">
                    {capitalizeFirstLetter(
                      appUser?.patronMembership?.patronPackageId ?? '',
                    )}{' '}
                    Patron (
                    {capitalizeFirstLetter(
                      appUser?.patronMembership?.patronTierId ?? '',
                    )}
                    )
                  </p>
                </div>
              </>
            )}
            <div className="w-full px-3 lg:px-10">
              <div
                className="max-w-xl w-full flex flex-col xl:flex-row text-white font-matahariRegular text-xl xl:text-2xl
              rounded-2xl bg-primary-800 items-center justify-center space-y-3 xl:space-y-0 xl:items-start xl:justify-between xl:space-x-3 mt-5 xl:mb-5 py-5 xl:py-10 xl:px-5"
              >
                <img
                  className="w-10"
                  src="/images/platinum.png"
                  alt="platinum"
                />
                <div className="flex flex-col space-y-2 text-center xl:text-start">
                  <p className="text-lg">Try the Diamond Plan</p>
                  <p className="text-sm">Upgrade to unlock more features</p>
                </div>
                <img className="w-10" src="/images/go.png" alt="platinum" />
              </div>
            </div>
          </div>
          <div className="flex w-full space-y-3 py-7 px-2 xl:px-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
            <div className="flex items-center w-full justify-between">
              <p className="font-semibold">Recent Transactions</p>
              <Link className="flex space-x-5 items-center" to={'history'}>
                <p className="test-xs text-primary-400">View all</p>
              </Link>
            </div>
            {isLoading ? (
              <div>Loading transaction history...</div>
            ) : rows.length > 0 ? (
              rows.map((item, index) => (
                <TransactionItem
                  key={index}
                  addressOrUsername={
                    item.type == 'Sent'
                      ? item.toUsername?.length == 0
                        ? truncateAddress(item.toAddress)
                        : (item.toUsername ?? '')
                      : item.fromUsername?.length == 0
                        ? truncateAddress(item.fromAddress)
                        : (item.fromUsername ?? '')
                  }
                  transactionType={item.type == 'Sent' ? 0 : 1}
                  amount={item.amountValue.toString()}
                  assetCode={item.assetCode}
                  date={item.dateLabel}
                  onclick={() => {
                    setHasAssets(!hasAssets);
                  }}
                />
              ))
            ) : (
              <div>No recent transactions.</div>
            )}
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
      {showSubscribeModal && (
        <ExpressInterestModal
          show={showSubscribeModal}
          onClose={() => {
            setShowSubscribeModal(false);
          }}
          asset={asset}
          onSuccess={() => {
            setShowSubscribeModal(false);
          }}
        />
      )}
      {/* Address copied modal */}
      <P2pComingSoonModal
        show={showP2pComingSoonModal}
        onClose={() => {
          setShowP2pComingSoonModal(false);
        }}
      />
      <KycUnverifiedModal
        show={showKycUnverifiedModal}
        onClose={() => {
          setShowKycUnverifiedModal(false);
        }}
      />
      <SoldOutModal
        show={showSoldOutModal}
        onClose={() => {
          setShowSoldOutModal(false);
        }}
      />
      {showBuyTokenModal && (
        <BuyTokenModal
          show={showBuyTokenModal}
          onClose={() => {
            setShowBuyTokenModal(false);
          }}
          onSuccess={() => {
            setShowBuyTokenModal(false);
          }}
          asset={asset}
        />
      )}
    </div>
  );
}
