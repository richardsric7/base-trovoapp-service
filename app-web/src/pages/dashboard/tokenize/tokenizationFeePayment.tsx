import { useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import { TokenizedAsset } from '../../../types/tokenizedAsset';
import { TokenizationData } from '../../../types/tokenizationData';
import Dropdown from '../../../components/dropdown';
import Button from '../../../components/button';
import ButtonSecondary from '../../../components/buttonSecondary';
import Modal from '../../../components/modal';
import {
  hideLoader,
  showLoader,
  showNotification,
} from '../../../utils/showToaster';
import { Encryptor } from '../../../utils/encryptor';
import { formatNumber, formatNumberShort } from '../../../utils/fomatNumber';
import {
  useConfirmTokenizationFeeMutation,
  useDeleteTokenizationFeeProofMutation,
  useLazyGetTokenizationDetailQuery,
  useUploadTokenizationFeeProofMutation,
} from '../../../store/api/tokenizationApis';
import { setActiveTokenizedAsset } from '../../../store/appStateSlice';
import { useActiveTokenizedAsset } from './useActiveTokenizedAsset';
import { useTokenizationData } from './useTokenizationData';

type PaymentMethodOption = {
  text: string;
  value: string;
};

type FeeDocument = {
  id?: number | string;
  documentUrl?: string;
  documentTitle?: string;
  transactionReference?: string;
};

function getStatusCode(responseData: unknown): number | undefined {
  const record = responseData as
    | { statusCode?: number; data?: { statusCode?: number } }
    | undefined;
  return record?.statusCode ?? record?.data?.statusCode;
}

function getMessage(responseData: unknown): string {
  const record = responseData as
    | {
        data?: {
          message?: string;
          error?: string;
          data?: { message?: string; error?: string };
        };
        message?: string;
      }
    | undefined;

  return (
    record?.data?.data?.message ??
    record?.data?.data?.error ??
    record?.data?.message ??
    record?.data?.error ??
    record?.message ??
    'Sorry, something went wrong. Please try again.'
  );
}

function truncateString(value: string, max: number = 48): string {
  if (!value) {
    return '';
  }

  if (value.length <= max) {
    return value;
  }

  return `${value.slice(0, max)}...`;
}

export function TokenizationFeePayment() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const appUser = useSelector((state: RootState) => state.auth.user);
  const appState = useSelector((state: RootState) => state.appState);
  const {
    activeTokenizedAsset,
    isLoading: isAssetLoading,
    fetchError: assetFetchError,
  } = useActiveTokenizedAsset();
  const {
    tokenizationData: tokenizationDataFromHook,
    isLoading: isTokenizationDataLoading,
    fetchError: tokenizationDataFetchError,
  } = useTokenizationData();

  const tokenizedAsset = (activeTokenizedAsset ??
    appState.activeTokenizedAsset) as TokenizedAsset | undefined;
  const tokenizationData = (tokenizationDataFromHook ??
    appState.tokenizationData) as TokenizationData | undefined;

  const [uploadFeeProof] = useUploadTokenizationFeeProofMutation();
  const [deleteFeeProof] = useDeleteTokenizationFeeProofMutation();
  const [confirmTokenizationFee] = useConfirmTokenizationFeeMutation();
  const [getTokenizationDetail] = useLazyGetTokenizationDetailQuery();
  const [isLoading, setIsLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);

  const paymentMethodOptions = useMemo<PaymentMethodOption[]>(() => {
    const methods = tokenizationData?.feePaymentMethods ?? [];
    return methods.map((method) => ({
      text: method.id.replaceAll(' ', ''),
      value: method.id,
    }));
  }, [tokenizationData]);

  const defaultPaymentMethod = paymentMethodOptions[0]?.value ?? 'FIAT';

  const [preferredPaymentMethod, setPreferredPaymentMethod] =
    useState(defaultPaymentMethod);
  const [hasMadePayment, setHasMadePayment] = useState(
    (tokenizedAsset?.ProofOfPaymentDocuments?.length ?? 0) > 0,
  );
  const [showUploadModal, setShowUploadModal] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [uploadError, setUploadError] = useState('');
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [transactionReference, setTransactionReference] = useState('');
  const [isRefreshing, setIsRefreshing] = useState(false);

  useEffect(() => {
    if (isAssetLoading || isTokenizationDataLoading) {
      setIsLoading(true);
      return;
    }

    const loadData = () => {
      setIsLoading(true);
      setLoadError(null);

      if (!tokenizedAsset) {
        setLoadError(assetFetchError ?? 'No tokenized asset data found.');
        setIsLoading(false);
        return;
      }

      if (!tokenizationData) {
        setLoadError(
          tokenizationDataFetchError ??
            'No tokenization metadata found for this application.',
        );
        setIsLoading(false);
        return;
      }

      setIsLoading(false);
    };

    void loadData();
  }, [
    isAssetLoading,
    isTokenizationDataLoading,
    tokenizedAsset,
    tokenizationData,
    assetFetchError,
    tokenizationDataFetchError,
  ]);

  useEffect(() => {
    if (
      paymentMethodOptions.length > 0 &&
      !paymentMethodOptions.some(
        (item) => item.value === preferredPaymentMethod,
      )
    ) {
      setPreferredPaymentMethod(paymentMethodOptions[0].value);
    }
  }, [paymentMethodOptions, preferredPaymentMethod]);

  const countryConfig = useMemo(() => {
    if (!tokenizedAsset || !tokenizationData?.countryConfigs) {
      return undefined;
    }

    return tokenizationData.countryConfigs.find((config) => {
      return (
        config.countryCode.toLowerCase() ===
        String(tokenizedAsset.assetCountryLocation ?? '').toLowerCase()
      );
    });
  }, [tokenizationData, tokenizedAsset]);

  const fiatCurrency = useMemo(() => {
    if (!tokenizedAsset || !tokenizationData) {
      return '';
    }

    const currency = tokenizationData.tokenizationCurrencies.find((entry) => {
      return (
        entry.assetCode.toLowerCase() ===
        String(countryConfig?.quoteCurrencyCode ?? '').toLowerCase()
      );
    });

    return (
      currency?.label?.toUpperCase() ??
      tokenizedAsset.assetQuoteCurrency?.toUpperCase() ??
      ''
    );
  }, [countryConfig, tokenizationData, tokenizedAsset]);

  const getFeeInfo = (feeIndex: number): number => {
    if (!tokenizedAsset || !tokenizationData?.tokenizationFees) {
      return 0;
    }

    const fee = tokenizationData.tokenizationFees.find(
      (entry, index) => index === feeIndex || entry.id === feeIndex,
    );

    if (!fee) {
      return 0;
    }

    const fiatPercentage = Number(fee.feeFiatPercentage ?? 0);
    const fiatCap = Number(fee.feeFiatCap ?? 0);
    const fiatFee =
      (Number(tokenizedAsset.assetCurrentValue ?? 0) * fiatPercentage) / 100;

    return fiatCap > fiatFee ? fiatCap : fiatFee;
  };

  const getVat = (): number => {
    if (!tokenizedAsset) {
      return 0;
    }

    const vatPercent = Number(tokenizedAsset.vatPercent ?? 0);
    const baseAmount =
      Number(tokenizedAsset.SECTokenizationFeeValue ?? 0) +
      Number(tokenizedAsset.issuingHouseFeeValue ?? 0) +
      Number(tokenizedAsset.legalAndProfessionalFeeValue ?? 0) +
      Number(tokenizedAsset.ratingAgencyFeeValue ?? 0) +
      getFeeInfo(tokenizedAsset.tokenizationFeeId ?? 0);

    return (baseAmount * vatPercent) / 100;
  };

  const getTotalFeeValue = (): number => {
    if (!tokenizedAsset) {
      return 0;
    }

    return (
      Number(tokenizedAsset.SECTokenizationFeeValue ?? 0) +
      Number(tokenizedAsset.issuingHouseFeeValue ?? 0) +
      Number(tokenizedAsset.legalAndProfessionalFeeValue ?? 0) +
      Number(tokenizedAsset.ratingAgencyFeeValue ?? 0) +
      getVat() +
      getFeeInfo(tokenizedAsset.tokenizationFeeId ?? 0)
    );
  };

  const tokenizationFee = formatNumber(
    getFeeInfo(tokenizedAsset?.tokenizationFeeId ?? 0),
  );
  const totalFee = formatNumberShort(getTotalFeeValue());

  const paymentCurrency =
    preferredPaymentMethod === 'STABLE COIN'
      ? (tokenizedAsset?.proceedPayoutCurrency?.replaceAll(' ', '') ?? '')
      : fiatCurrency;

  const proofDocuments =
    (tokenizedAsset?.ProofOfPaymentDocuments as FeeDocument[] | undefined) ??
    [];

  const copyToClipboard = async (value: string, label: string) => {
    try {
      await navigator.clipboard.writeText(value);
      showNotification('success', `${label} copied`);
    } catch {
      showNotification('error', `Unable to copy ${label.toLowerCase()}`);
    }
  };

  const refreshTokenizationInfo = async () => {
    if (!appUser || !tokenizedAsset) {
      return;
    }

    setIsRefreshing(true);
    try {
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);

      const res = await getTokenizationDetail({
        signer: appUser.primarySigner,
        address: appUser.address,
        secretKey: currentSecretKey,
        tokenizedAssetId: tokenizedAsset.id,
      });

      if ('data' in res) {
        const responseData = res.data;
        const statusCode = getStatusCode(responseData);
        if (statusCode === 200) {
          const detailData =
            (responseData as { data?: { data?: TokenizedAsset } })?.data
              ?.data ??
            (responseData as { data?: TokenizedAsset })?.data ??
            tokenizedAsset;
          dispatch(setActiveTokenizedAsset(detailData));
        }
      }
    } catch {
      // ignore refresh errors after a successful upload/delete action
    } finally {
      setIsRefreshing(false);
    }
  };

  const handleUploadProof = async () => {
    if (!appUser || !tokenizedAsset) {
      return;
    }

    if (!transactionReference.trim()) {
      setUploadError('Enter transaction reference.');
      return;
    }

    if (preferredPaymentMethod !== 'STABLE COIN' && !selectedFile) {
      setUploadError('Select proof of payment file.');
      return;
    }

    if (selectedFile && selectedFile.size > 900000) {
      setUploadError('File size must be less than 900KB.');
      return;
    }

    setUploading(true);
    setUploadError('');
    showLoader();

    try {
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);

      const res = await uploadFeeProof({
        signer: appUser.primarySigner,
        address: appUser.address,
        secretKey: currentSecretKey,
        tokenizedAssetId: tokenizedAsset.id,
        transactionReference: transactionReference.trim(),
        tokenizationFeePaymentMethodID: preferredPaymentMethod,
        file: selectedFile ?? undefined,
      });

      hideLoader();

      if ('data' in res) {
        showNotification('success', 'Proof of payment uploaded successfully.');
        setShowUploadModal(false);
        setSelectedFile(null);
        setTransactionReference('');
        await refreshTokenizationInfo();
        return;
      }

      if ('data' in res) {
        setUploadError(getMessage(res.data));
      } else {
        setUploadError('Sorry, something went wrong. Please try again.');
      }
    } catch {
      hideLoader();
      setUploadError('Sorry, something went wrong. Please try again.');
    } finally {
      setUploading(false);
    }
  };

  const handleDeleteProof = async (documentId: string) => {
    if (!appUser || !documentId) {
      return;
    }

    const confirmed = window.confirm(
      'Are you sure you want to delete this proof document?',
    );
    if (!confirmed) {
      return;
    }

    showLoader();
    try {
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);

      const res = await deleteFeeProof({
        signer: appUser.primarySigner,
        address: appUser.address,
        secretKey: currentSecretKey,
        documentId,
      });

      hideLoader();

      if ('data' in res) {
        showNotification('success', 'Proof of payment deleted.');
        await refreshTokenizationInfo();
        return;
      }

      if ('data' in res) {
        showNotification('error', getMessage(res.data));
      } else {
        showNotification('error', 'Failed to delete proof document.');
      }
    } catch {
      hideLoader();
      showNotification(
        'error',
        'Sorry, something went wrong. Please try again.',
      );
    }
  };

  const confirmPayment = async () => {
    if (!appUser || !tokenizedAsset) {
      return;
    }

    showLoader();
    try {
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);

      const res = await confirmTokenizationFee({
        signer: appUser.primarySigner,
        address: appUser.address,
        secretKey: currentSecretKey,
        tokenizedAssetId: tokenizedAsset.id,
        body: '',
      });

      hideLoader();

      if ('data' in res) {
        showNotification(
          'success',
          'Your proof of payment has been submitted successfully and is awaiting confirmation.',
        );
        setShowSuccessModal(true);
        return;
      }

      if ('error' in res) {
        console.error('error', getMessage(res.error));
        showNotification('error', 'Failed to confirm payment.');
      }
    } catch {
      hideLoader();
      showNotification(
        'error',
        'Sorry, something went wrong. Please try again.',
      );
    }
  };

  if (isLoading) {
    return (
      <div className="rounded-lg border border-gray-200 bg-primary-100 p-6 text-sm text-gray-600">
        Loading payment details...
      </div>
    );
  }

  if (loadError && !tokenizedAsset) {
    return (
      <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-600">
        {loadError}
      </div>
    );
  }

  if (!tokenizedAsset) {
    return (
      <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-600">
        No tokenized asset data found.
      </div>
    );
  }

  const dropdownDefault = paymentMethodOptions.find(
    (item) => item.value === preferredPaymentMethod,
  );

  return (
    <div className="w-full max-w-5xl mr-auto bg-white p-4 md:p-6 lg:p-8 space-y-6">
      <h2 className="text-lg md:text-xl font-montserratSemiBold text-[#0E1F51]">
        Payment
      </h2>

      <div className="rounded-xl border border-gray-200 bg-primary-100 p-4 md:p-5 space-y-2">
        <p className="text-sm font-medium text-[#0E1F51]">
          Select payment method
        </p>
        <Dropdown
          label="Select payment method"
          defaultValue={dropdownDefault}
          options={paymentMethodOptions}
          onSelect={(item) => {
            setPreferredPaymentMethod(String(item.value));
          }}
        />
      </div>

      <div className="rounded-xl border border-gray-200 bg-primary-100 p-4 md:p-5">
        <div className="flex flex-wrap items-center gap-2 text-[#0E1F51]">
          <p className="text-lg font-montserratSemiBold">Pay</p>
          <p className="text-lg font-montserratSemiBold">{totalFee}</p>
          <p className="text-lg font-montserratSemiBold">{paymentCurrency}</p>
          <button
            type="button"
            onClick={() => {
              void copyToClipboard(totalFee.replaceAll(',', ''), 'Amount');
            }}
            className="rounded-md border border-primary-800 px-3 py-1 text-xs font-semibold text-primary-800"
          >
            Copy Amount
          </button>
        </div>
      </div>

      {preferredPaymentMethod === 'STABLE COIN' ? (
        <div className="rounded-xl border border-gray-200 bg-primary-100 p-4 md:p-5 space-y-3">
          <p className="font-montserratSemiBold text-[#0E1F51]">
            Wallet address
          </p>
          <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <p className="text-sm text-[#0E1F51] break-all md:max-w-[78%]">
              {tokenizedAsset.walletToHoldAssetsNotForSale?.toLowerCase() ??
                'N/A'}
            </p>
            <button
              type="button"
              onClick={() => {
                void copyToClipboard(
                  tokenizedAsset.walletToHoldAssetsNotForSale ?? '',
                  'Wallet address',
                );
              }}
              className="rounded-md border border-primary-800 px-3 py-1 text-xs font-semibold text-primary-800"
            >
              Copy Address
            </button>
          </div>
        </div>
      ) : (
        <div className="rounded-xl border border-gray-200 bg-primary-100 p-4 md:p-5 space-y-3">
          <p className="font-montserratSemiBold text-[#0E1F51]">
            Bank account details
          </p>
          <div className="space-y-3">
            <div className="flex items-center justify-between text-sm">
              <span className="text-[#0E1F51]">Account Name</span>
              <span className="font-montserratSemiBold text-[#0E1F51] text-right">
                Trovotech Limited
              </span>
            </div>
            <div className="flex items-center justify-between text-sm gap-3">
              <span className="text-[#0E1F51]">Account Number</span>
              <div className="flex items-center gap-3">
                <span className="font-montserratSemiBold text-[#0E1F51]">
                  0088066577
                </span>
                <button
                  type="button"
                  onClick={() => {
                    void copyToClipboard('0088066577', 'Account number');
                  }}
                  className="rounded-md border border-primary-800 px-3 py-1 text-xs font-semibold text-primary-800"
                >
                  Copy
                </button>
              </div>
            </div>
            <div className="flex items-center justify-between text-sm">
              <span className="text-[#0E1F51]">Bank</span>
              <span className="font-montserratSemiBold text-[#0E1F51] text-right">
                Sterling Bank
              </span>
            </div>
          </div>
        </div>
      )}

      <div className="rounded-xl border border-gray-200 bg-primary-100 p-4 md:p-5">
        <p className="text-red-500 font-montserratSemiBold">IMPORTANT!</p>
        <p className="mt-2 text-sm text-[#0E1F51]">
          Your outstanding fee details is as follows:
        </p>

        <div className="mt-3 space-y-2 text-sm">
          <div className="flex items-center justify-between gap-3">
            <span className="text-[#0E1F51]">Tokenization Fee</span>
            <span className="font-montserratSemiBold text-[#0E1F51] text-right">
              {tokenizationFee} {fiatCurrency}
            </span>
          </div>
          <div className="flex items-center justify-between gap-3">
            <span className="text-[#0E1F51]">SEC Fee</span>
            <span className="font-montserratSemiBold text-[#0E1F51] text-right">
              {formatNumberShort(tokenizedAsset.SECTokenizationFeeValue ?? 0)}{' '}
              {fiatCurrency}
            </span>
          </div>

          {(tokenizedAsset.issuingHouseFeeValue ?? 0) > 0 && (
            <div className="flex items-center justify-between gap-3">
              <span className="text-[#0E1F51]">Issuing House Fee</span>
              <span className="font-montserratSemiBold text-[#0E1F51] text-right">
                {formatNumberShort(tokenizedAsset.issuingHouseFeeValue ?? 0)}{' '}
                {fiatCurrency}
              </span>
            </div>
          )}

          {(tokenizedAsset.legalAndProfessionalFeeValue ?? 0) > 0 && (
            <div className="flex items-center justify-between gap-3">
              <span className="text-[#0E1F51]">Legal/Professional Fee</span>
              <span className="font-montserratSemiBold text-[#0E1F51] text-right">
                {formatNumberShort(
                  tokenizedAsset.legalAndProfessionalFeeValue ?? 0,
                )}{' '}
                {fiatCurrency}
              </span>
            </div>
          )}

          {(tokenizedAsset.ratingAgencyFeeValue ?? 0) > 0 && (
            <div className="flex items-center justify-between gap-3">
              <span className="text-[#0E1F51]">Rating Agency Fee</span>
              <span className="font-montserratSemiBold text-[#0E1F51] text-right">
                {formatNumberShort(tokenizedAsset.ratingAgencyFeeValue ?? 0)}{' '}
                {fiatCurrency}
              </span>
            </div>
          )}

          <div className="flex items-center justify-between gap-3">
            <span className="text-[#0E1F51]">VAT</span>
            <span className="font-montserratSemiBold text-[#0E1F51] text-right">
              {formatNumberShort(getVat())} {fiatCurrency}
            </span>
          </div>

          <div className="flex items-center justify-between gap-3 border-t border-gray-200 pt-2">
            <span className="text-[#0E1F51]">Total Fee</span>
            <span className="font-montserratSemiBold text-[#0E1F51] text-right">
              {totalFee} {fiatCurrency}
            </span>
          </div>
        </div>
      </div>

      <div className="rounded-xl border border-gray-200 bg-primary-100 p-4 md:p-5 space-y-4">
        <label className="flex items-start gap-3">
          <input
            type="checkbox"
            className="mt-1 h-4 w-4 accent-primary-800"
            checked={hasMadePayment}
            onChange={(event) => {
              setHasMadePayment(event.target.checked);
            }}
          />
          <span className="text-sm text-[#0E1F51]">I have made payment</span>
        </label>
      </div>

      {proofDocuments.length > 0 && (
        <div className="space-y-3">
          {proofDocuments.map((item, index) => {
            const documentId =
              item.id != null
                ? String(item.id)
                : `${index}-${item.documentUrl}`;
            const label =
              item.documentTitle ?? item.documentUrl ?? 'Proof file';
            return (
              <div
                key={documentId}
                className="rounded-lg border border-gray-200 bg-primary-100 px-4 py-3"
              >
                <button
                  type="button"
                  onClick={() => {
                    const targetUrl = item.documentUrl ?? '';
                    if (targetUrl.trim().length > 0) {
                      window.open(targetUrl, '_blank', 'noopener,noreferrer');
                    }
                  }}
                  className="text-left underline text-primary-800 break-all"
                >
                  {truncateString(label)}
                </button>

                {item.transactionReference ? (
                  <p className="mt-1 text-sm text-[#0E1F51] break-all">
                    <span className="font-montserratSemiBold">Ref:</span>{' '}
                    {item.transactionReference}
                  </p>
                ) : null}

                {item.id != null ? (
                  <button
                    type="button"
                    onClick={() => {
                      void handleDeleteProof(String(item.id));
                    }}
                    className="mt-2 text-xs font-semibold text-red-500"
                  >
                    Delete
                  </button>
                ) : null}
              </div>
            );
          })}
        </div>
      )}

      {hasMadePayment && (
        <button
          type="button"
          onClick={() => {
            setUploadError('');
            setShowUploadModal(true);
          }}
          className="text-primary-700 underline font-montserratSemiBold"
        >
          Upload proof of payment
        </button>
      )}

      <div className="flex items-end w-4/6 space-x-3">
        <ButtonSecondary
          label="Go to homepage"
          onclick={() => {
            navigate('/dashboard');
          }}
        />
        <Button
          label={isRefreshing ? 'Please wait...' : 'Confirm payment'}
          disabled={isRefreshing}
          onclick={() => {
            void confirmPayment();
          }}
        />
      </div>

      <Modal
        showModal={showUploadModal}
        onClose={() => {
          if (!uploading) {
            setShowUploadModal(false);
          }
        }}
      >
        <div className="px-5 pb-6 space-y-4">
          <h3 className="font-montserratSemiBold text-lg text-[#0E1F51]">
            Upload Proof of Payment
          </h3>

          <div>
            <p className="mb-1 text-sm text-[#0E1F51]">Transaction Reference</p>
            <input
              type="text"
              value={transactionReference}
              onChange={(event) => {
                setTransactionReference(event.target.value);
              }}
              placeholder="Enter transaction reference"
              className="h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 outline-none focus:border-primary-600"
            />
          </div>

          {preferredPaymentMethod !== 'STABLE COIN' && (
            <div>
              <p className="mb-1 text-sm text-[#0E1F51]">Proof Document</p>
              {/* <input
                type="file"
                onChange={(event) => {
                  const file = event.target.files?.[0] ?? null;
                  setSelectedFile(file);
                }}
                className="w-full text-sm text-[#0E1F51] file:mr-4 file:rounded-md file:border-0 file:bg-primary-700 file:px-4 file:py-2 file:text-sm file:font-semibold file:text-white"
              />
              <p className="mt-1 text-xs text-gray-500">
                Maximum file size: 900KB
              </p> */}

              {/* Hidden file input for logo upload */}
              <input
                ref={fileInputRef}
                type="file"
                className="hidden"
                onChange={(event) => {
                  const file = event.target.files?.[0] ?? null;
                  setSelectedFile(file);
                }}
                accept=".jpg,.jpeg,.gif,.png,.pdf,.svg"
              />

              <div
                className="cursor-pointer mt-2"
                onClick={() => fileInputRef.current?.click()}
              >
                <div className="border-2 border-dashed border-primary-300 rounded-lg p-8 flex flex-col items-center justify-center gap-2 hover:bg-primary-50 transition-colors">
                  <svg
                    className="w-8 h-8 text-gray-400"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                    />
                  </svg>
                  <span className="text-sm font-semibold text-gray-600">
                    {uploading
                      ? 'Uploading...'
                      : 'Browse files to upload asset logo'}
                  </span>
                  <span className="text-xs text-gray-400">
                    Supported: JPG, JPEG, GIF, PNG, PDF (Max 900KB)
                  </span>
                </div>
              </div>
              <p className="my-2 text-sm text-[#0E1F51]">
                {selectedFile?.name}
              </p>
            </div>
          )}

          {uploadError && <p className="text-sm text-red-500">{uploadError}</p>}

          <div className="flex gap-3">
            <div className="w-1/2">
              <ButtonSecondary
                label="Cancel"
                onclick={() => {
                  if (!uploading) {
                    setShowUploadModal(false);
                  }
                }}
              />
            </div>
            <div className="w-1/2">
              <Button
                label={uploading ? 'Uploading...' : 'Submit'}
                disabled={uploading}
                onclick={() => {
                  void handleUploadProof();
                }}
              />
            </div>
          </div>
        </div>
      </Modal>
      <Modal
        showModal={showSuccessModal}
        onClose={() => {
          setShowSuccessModal(false);
          navigate('/dashboard/tokenize');
        }}
      >
        <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
          <img src="/images/launch.png" alt="success" />
          <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
            <p className="text-primary-800 text-md xl:text-lg font-semibold">
              Your Proof of Payment has been submitted successfully and is
              awaiting confirmation. Your asset tokenization application will be
              processed once payment has been confirmed.
            </p>
            <div></div>
          </div>
          <div className="w-3/4">
            <Button
              label="Done"
              onclick={() => {
                navigate('/dashboard/tokenize');
              }}
            />
          </div>
        </div>
      </Modal>
    </div>
  );
}

export default TokenizationFeePayment;
