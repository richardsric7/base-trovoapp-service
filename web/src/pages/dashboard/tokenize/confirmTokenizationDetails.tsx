import { useEffect, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import { Encryptor } from '../../../utils/encryptor';
import { signBase64Txn } from '../../../utils/trovoSDK';
import { formatNumber, formatNumberShort } from '../../../utils/fomatNumber';
import {
  showNotification,
  showLoader,
  hideLoader,
} from '../../../utils/showToaster';
import { useConfirmTokenizationMutation } from '../../../store/api/tokenizationApis';
import { ErrorResponse } from '../../../store/api/baseapi/axiosBaseQuery';
import Button from '../../../components/button';
import ButtonSecondary from '../../../components/buttonSecondary';
import Modal from '../../../components/modal';
import { TokenizedAsset } from '../../../types/tokenizedAsset';
import { TokenizationData } from '../../../types/tokenizationData';
import { useActiveTokenizedAsset } from './useActiveTokenizedAsset';
import { useTokenizationData } from './useTokenizationData';

function truncateToDecimalPlaces(
  value: number,
  decimalPlaces: number = 0,
): string {
  return value.toFixed(decimalPlaces);
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr);
  return date.toLocaleDateString('en-US', {
    month: 'long',
    day: '2-digit',
    year: 'numeric',
  });
}

export function ConfirmTokenizationDetails() {
  const navigate = useNavigate();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const appState = useSelector((state: RootState) => state.appState);
  const [confirmTokenization] = useConfirmTokenizationMutation();
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

  const [isLoading, setIsLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [password, setPassword] = useState('');
  const [passwordError, setPasswordError] = useState('');
  const [pendingTransaction, setPendingTransaction] = useState<any>(null);
  const [showTxnModal, setShowTxnModal] = useState(false);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [tokenizedAsset, setTokenizedAsset] = useState<TokenizedAsset | null>(
    null,
  );
  const [fiatCurrency, setFiatCurrency] = useState('');
  const [tokenizationApplicationFee, setTokenizationApplicationFee] =
    useState(0);
  const [tokenizationApplicationFeeAsset, setTokenizationApplicationFeeAsset] =
    useState('');
  const [feeInfo, setFeeInfo] = useState('');
  const [isFinancialAssetType, setIsFinancialAssetType] = useState(false);

  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      setLoadError(null);
      try {
        const asset = activeTokenizedAsset;
        // const asset = {} as TokenizedAsset;
        if (!asset) {
          setLoadError(assetFetchError ?? 'No tokenized asset data found');
          return;
        }

        setTokenizedAsset(asset);
        const isFinancial =
          asset.assetSector?.toLowerCase() === 'finance and investment markets';
        setIsFinancialAssetType(isFinancial);

        // Find country config
        const tokenizationData = (tokenizationDataFromHook ??
          appState.tokenizationData) as TokenizationData | undefined;
        if (!tokenizationData) {
          setLoadError(
            tokenizationDataFetchError ??
              'No tokenization metadata found for this application.',
          );
          return;
        }

        let quoteCurrencyCode = '';
        let appFee = 0;
        let appFeeAsset = '';

        if (tokenizationData?.countryConfigs) {
          for (const config of tokenizationData.countryConfigs) {
            if (
              config.countryCode.toLowerCase() ===
              asset.assetCountryLocation?.toLowerCase()
            ) {
              quoteCurrencyCode = config.quoteCurrencyCode;
              appFee = config.tokenizationApplicationFee ?? 0;
              appFeeAsset =
                config.tokenizationApplicationFeeAsset?.split(':')[0] ?? '';
            }
          }
        }

        setTokenizationApplicationFee(appFee);
        setTokenizationApplicationFeeAsset(appFeeAsset);

        // Find fiat currency label
        if (tokenizationData?.tokenizationCurrencies) {
          for (const currency of tokenizationData.tokenizationCurrencies) {
            if (
              currency.assetCode?.toLowerCase() ===
              quoteCurrencyCode.toLowerCase()
            ) {
              setFiatCurrency(
                currency.label?.toUpperCase() ?? quoteCurrencyCode,
              );
            }
          }
        }

        // Calculate fee info
        if (
          asset.tokenizationFeeId != null &&
          tokenizationData?.tokenizationFees
        ) {
          const feeIndex = asset.tokenizationFeeId;
          // Find the fee by index or id
          const fee = tokenizationData.tokenizationFees.find(
            (f, i) => i === feeIndex || f.id === feeIndex,
          );
          if (fee) {
            const fiatPercentage = fee.feeFiatPercentage ?? 0;
            const fiatCap = fee.feeFiatCap ?? 0;
            const calculatedFee =
              ((asset.assetCurrentValue ?? 0) * fiatPercentage) / 100;
            const finalFee = fiatCap > calculatedFee ? fiatCap : calculatedFee;
            setFeeInfo(formatNumber(finalFee));
          }
        }
      } catch (err) {
        setLoadError('Failed to load tokenization details');
        showNotification('error', 'Failed to load tokenization details');
      } finally {
        setIsLoading(false);
      }
    };

    if (isAssetLoading || isTokenizationDataLoading) {
      setIsLoading(true);
      return;
    }

    void loadData();
  }, [
    activeTokenizedAsset,
    isAssetLoading,
    isTokenizationDataLoading,
    assetFetchError,
    tokenizationDataFetchError,
    tokenizationDataFromHook,
    appUser,
    appState,
    navigate,
  ]);

  const getFeeInfo = useCallback(
    (feeIndex: number): number => {
      const tokenizationData = appState.tokenizationData as
        | TokenizationData
        | undefined;
      if (!tokenizationData?.tokenizationFees) return 0;

      const fee = tokenizationData.tokenizationFees.find(
        (f, i) => i === feeIndex || f.id === feeIndex,
      );
      if (!fee) return 0;

      const fiatPercentage = fee.feeFiatPercentage ?? 0;
      const fiatCap = fee.feeFiatCap ?? 0;
      const calculatedFee =
        ((tokenizedAsset?.assetCurrentValue ?? 0) * fiatPercentage) / 100;
      return fiatCap > calculatedFee ? fiatCap : calculatedFee;
    },
    [appState.tokenizationData, tokenizedAsset],
  );

  const getTotalFee = useCallback((): string => {
    if (!tokenizedAsset) return '0';
    const total =
      (tokenizedAsset.SECTokenizationFeeValue ?? 0) +
      (tokenizedAsset.custodianFeeValue ?? 0) +
      (tokenizedAsset.assetManagerFeeValue ?? 0) +
      (tokenizedAsset.issuingHouseFeeValue ?? 0) +
      (tokenizedAsset.legalAndProfessionalFeeValue ?? 0) +
      (tokenizedAsset.ratingAgencyFeeValue ?? 0) +
      (tokenizedAsset.trusteeFeeValue ?? 0) +
      (tokenizedAsset.vatValue ?? 0) +
      getFeeInfo(tokenizedAsset.tokenizationFeeId ?? 0);

    return formatNumberShort(total);
  }, [tokenizedAsset, getFeeInfo]);

  const isValidPassword = async (): Promise<boolean> => {
    try {
      const encryptor = new Encryptor();
      await encryptor.decryptData(
        appUser.secretKeys[0],
        password,
        appUser.primarySigner,
      );
      return true;
    } catch (error: any) {
      return false;
    }
  };

  const handlePasswordSubmit = async () => {
    if (!password.trim()) {
      setPasswordError('Enter your password');
      return;
    }
    if (password.trim().length < 6) {
      setPasswordError('Use 6 characters or more for your password');
      return;
    }

    setPasswordError('');

    if (!(await isValidPassword())) {
      setPasswordError('Password is invalid!');
      return;
    }

    await submitForm();
  };

  const submitForm = async () => {
    if (!tokenizedAsset) return;

    setIsSubmitting(true);
    showLoader();

    try {
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);

      const requestBody = JSON.stringify({});

      const res = await confirmTokenization({
        signer: appUser.primarySigner,
        publicKey: appUser.publicKey,
        secretKey: currentSecretKey,
        tokenizedAssetId: tokenizedAsset.id,
        body: requestBody,
      });

      hideLoader();

      if ('data' in res) {
        const data = res.data;
        // The confirm endpoint returns the unsigned transaction directly
        // (with messages, transaction, networkPassPhrase, etc.) when a
        // signature is required. It must be shown to the user for approval
        // before signing and resubmitting.
        if (data?.transaction) {
          setPendingTransaction(data);
          setShowTxnModal(true);
        }
      } else if ('error' in res) {
        const errorData = (res as any).error as ErrorResponse;
        showNotification(
          'error',
          errorData?.data?.message ?? 'Failed to confirm tokenization',
        );
      }
    } catch (err: any) {
      hideLoader();
      showNotification('error', err?.message ?? 'An error occurred');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleApproveTransaction = async () => {
    if (!tokenizedAsset || !pendingTransaction) return;

    setIsSubmitting(true);
    setShowTxnModal(false);
    showLoader();

    try {
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);

      const signature = signBase64Txn(
        currentSecretKey,
        pendingTransaction?.transaction,
        pendingTransaction?.networkPassPhrase,
      );

      const requestBody = JSON.stringify({
        ...pendingTransaction,
        transactionSignature: signature,
      });

      const res = await confirmTokenization({
        signer: appUser.primarySigner,
        publicKey: appUser.publicKey,
        secretKey: currentSecretKey,
        tokenizedAssetId: tokenizedAsset.id,
        body: requestBody,
      });

      hideLoader();

      if ('data' in res) {
        setShowSuccessModal(true);
      } else if ('error' in res) {
        const errorData = (res as any).error as ErrorResponse;
        showNotification(
          'error',
          errorData?.data?.message ?? 'Failed to submit tokenization',
        );
      }
    } catch (err: any) {
      hideLoader();
      showNotification('error', err?.message ?? 'An error occurred');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeclineTransaction = () => {
    setPendingTransaction(null);
    setShowTxnModal(false);
  };

  const renderItem = (
    key: string,
    value: string,
    isNotUpfront: boolean = false,
  ) => {
    return (
      <div className="flex justify-between items-center py-1 px-4">
        <div className="max-w-[45%]">
          <span className="text-sm font-medium text-gray-700">
            {key}
            {isNotUpfront && <span className="text-red-500 ml-0.5"> *</span>}
          </span>
        </div>
        <div className="w-[45%]">
          <span className="text-sm font-semibold text-gray-700 block text-right">
            {value}
          </span>
        </div>
      </div>
    );
  };

  if (isLoading) {
    return (
      <div className="bg-white p-6 w-full mx-auto">
        <div className="rounded-lg border border-gray-200 bg-primary-100 p-6 text-sm text-gray-600">
          Loading confirmation details...
        </div>
      </div>
    );
  }

  if (loadError && !tokenizedAsset) {
    return (
      <div className="bg-white p-6 w-full mx-auto">
        <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-600">
          {loadError}
        </div>
      </div>
    );
  }

  if (!tokenizedAsset) {
    return (
      <div className="bg-white p-6 w-full mx-auto">
        <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-600">
          No tokenized asset data found.
        </div>
      </div>
    );
  }

  const isAlreadySubmitted =
    (tokenizedAsset.vettingStatus ?? 0) >= 1 ||
    (tokenizedAsset.assetTokenizationStatus ?? 0) >= 1;
  const isVetted = tokenizedAsset.vettingStatus === 1;

  return (
    <div className="bg-white p-6 w-full mx-auto">
      {/* Page Title */}
      <h2 className="text-lg md:text-xl font-semibold text-[#0E1F51] mb-6">
        {isAlreadySubmitted && isVetted
          ? 'Vetted Summary'
          : 'Confirm Your Information'}
      </h2>

      <div className="space-y-4">
        {/* Vetted info message */}
        {tokenizedAsset.assetTokenizationStatus === 2 ? (
          <div className="px-5">
            <p className="text-sm font-medium text-green-600">
              Your proof of payment has been submitted successfully and is
              awaiting confirmation.
            </p>
          </div>
        ) : (
          isAlreadySubmitted &&
          isVetted && (
            <div className="px-5">
              <p className="text-sm font-medium text-gray-600">
                Your application has been vetted, please confirm the information
                below and proceed to pay for tokenization.
              </p>
            </div>
          )
        )}

        {/* Asset Information Section */}
        <div className="rounded-xl border border-gray-200 bg-primary-100 p-4">
          <h3 className="text-base font-semibold text-[#0E1F51] mb-3 px-4">
            Asset Information
          </h3>
          {renderItem('Asset Name', `${tokenizedAsset.assetName ?? ''}`)}
          {renderItem('Asset Code', `${tokenizedAsset.assetCode ?? ''}`)}

          {(tokenizedAsset.assetAlreadyExists ?? 0) === 1 ? (
            <>
              {renderItem(
                isFinancialAssetType
                  ? 'Net Asset Value'
                  : 'Original Asset Value',
                `${truncateToDecimalPlaces(tokenizedAsset.assetCurrentValue ?? 0, 2)} ${fiatCurrency}`,
              )}
              {!isFinancialAssetType && (
                <>
                  {renderItem(
                    'Percentage Retained',
                    `${formatNumber(((tokenizedAsset.assetOwnerRetainedOrContributedValue ?? 0) / (tokenizedAsset.assetCurrentValue ?? 1)) * 100)}%`,
                  )}
                  {renderItem(
                    'Value Retained',
                    `${truncateToDecimalPlaces(tokenizedAsset.assetOwnerRetainedOrContributedValue ?? 0, 2)} ${fiatCurrency}`,
                  )}
                  {renderItem(
                    'Additional Cost',
                    `${truncateToDecimalPlaces(tokenizedAsset.assetMscCostOutisdeOfValuation ?? 0, 2)} ${fiatCurrency}`,
                  )}
                </>
              )}
            </>
          ) : (
            <>
              {renderItem(
                'Project Budget',
                `${truncateToDecimalPlaces(tokenizedAsset.assetCurrentValue ?? 0, 2)} ${fiatCurrency}`,
              )}
              {renderItem(
                'Equity Contribution',
                `${formatNumber(((tokenizedAsset.assetOwnerRetainedOrContributedValue ?? 0) / (tokenizedAsset.assetCurrentValue ?? 1)) * 100)}%`,
              )}
              {renderItem(
                'Value of Equity',
                `${truncateToDecimalPlaces(tokenizedAsset.assetOwnerRetainedOrContributedValue ?? 0, 2)} ${fiatCurrency}`,
              )}
            </>
          )}
        </div>

        {/* Asset Token Information Section */}
        <div className="rounded-xl border border-gray-200 bg-primary-100 p-4">
          <h3 className="text-base font-semibold text-[#0E1F51] mb-3 px-4">
            Asset Token Information
          </h3>
          {isVetted ? (
            <>
              {renderItem(
                'Total Token',
                `${truncateToDecimalPlaces(tokenizedAsset.numberOfTokenToBeIssued ?? 0)} ${tokenizedAsset.assetCode ?? ''}`,
              )}
              {renderItem(
                'Final Value',
                `${truncateToDecimalPlaces(tokenizedAsset.valueOfTokenizedAsset ?? 0, 2)} ${fiatCurrency}`,
              )}
              {renderItem(
                'Price per Token',
                `${truncateToDecimalPlaces(tokenizedAsset.pricePerToken ?? 0)} ${fiatCurrency}`,
              )}
              {renderItem(
                'Tokens not for Sale',
                `${truncateToDecimalPlaces(
                  (tokenizedAsset.numberOfTokenToBeIssued ?? 0) -
                    (tokenizedAsset.numberOfTokenToBeSold ?? 0) -
                    (tokenizedAsset.feeInAsset ?? 0) -
                    (tokenizedAsset.feeInAsset ?? 0) * 0.075,
                )} ${tokenizedAsset.assetCode ?? ''}`,
              )}
              {renderItem(
                'Tokens for Sale',
                `${truncateToDecimalPlaces(tokenizedAsset.numberOfTokenToBeSold ?? 0)} ${tokenizedAsset.assetCode ?? ''}`,
              )}
              {renderItem(
                'Amount to be Raised',
                `${truncateToDecimalPlaces(
                  (tokenizedAsset.numberOfTokenToBeSold ?? 0) *
                    (tokenizedAsset.pricePerToken ?? 0),
                  2,
                )} ${fiatCurrency}`,
              )}
            </>
          ) : (
            <>
              {renderItem(
                'Proposed Total Tokens to be Issued',
                `${formatNumber(tokenizedAsset.numberOfTokenToBeIssued ?? 0)} ${tokenizedAsset.assetCode ?? ''}`,
              )}
            </>
          )}
        </div>

        {/* Primary Offering Section */}
        <div className="rounded-xl border border-gray-200 bg-primary-100 p-4">
          <h3 className="text-base font-semibold text-[#0E1F51] mb-3 px-4">
            Primary Offering
          </h3>
          {isVetted ? (
            <>
              {renderItem(
                'Start Date',
                `${formatDate(tokenizedAsset.salesStart ?? '')}`,
              )}
              {renderItem(
                'End Date',
                `${formatDate(tokenizedAsset.salesEnd ?? '')}`,
              )}
            </>
          ) : (
            <>
              {renderItem(
                'Proposed Start Date',
                `${formatDate(tokenizedAsset.salesStart ?? '')}`,
              )}
              {renderItem(
                'Proposed End Date',
                `${formatDate(tokenizedAsset.salesEnd ?? '')}`,
              )}
            </>
          )}
        </div>

        {/* Fees Section */}
        <div className="rounded-xl border border-gray-200 bg-primary-100 p-4">
          <h3 className="text-base font-semibold text-[#0E1F51] mb-3 px-4">
            {isVetted ? 'Fees in Fiat' : 'Application Fee'}
          </h3>
          {isVetted ? (
            <>
              {renderItem('Tokenization Fee', `${feeInfo} ${fiatCurrency}`)}
              {renderItem(
                'SEC Fee',
                `${formatNumberShort(tokenizedAsset.SECTokenizationFeeValue ?? 0)} ${fiatCurrency}`,
              )}
              {renderItem(
                'Custody Fee',
                `${formatNumberShort(tokenizedAsset.custodianFeeValue ?? 0)} ${fiatCurrency}`,
                true,
              )}
              {renderItem(
                'Management Fee',
                `${formatNumberShort(tokenizedAsset.assetManagerFeeValue ?? 0)} ${fiatCurrency}`,
                true,
              )}
              {(tokenizedAsset.issuingHouseFeeValue ?? 0) > 0 &&
                renderItem(
                  'Issuing House Fee',
                  `${formatNumberShort(tokenizedAsset.issuingHouseFeeValue ?? 0)} ${fiatCurrency}`,
                )}
              {(tokenizedAsset.legalAndProfessionalFeeValue ?? 0) > 0 &&
                renderItem(
                  'Legal/Professional Fee',
                  `${formatNumberShort(tokenizedAsset.legalAndProfessionalFeeValue ?? 0)} ${fiatCurrency}`,
                )}
              {(tokenizedAsset.ratingAgencyFeeValue ?? 0) > 0 &&
                renderItem(
                  'Rating Agency Fee',
                  `${formatNumberShort(tokenizedAsset.ratingAgencyFeeValue ?? 0)} ${fiatCurrency}`,
                )}
              {(tokenizedAsset.trusteeFeeValue ?? 0) > 0 &&
                renderItem(
                  'Trustee Fee',
                  `${formatNumberShort(tokenizedAsset.trusteeFeeValue ?? 0)} ${fiatCurrency}`,
                )}
              {renderItem(
                'VAT (Fiat)',
                `${formatNumberShort(tokenizedAsset.vatValue ?? 0)} ${fiatCurrency}`,
              )}
              {renderItem('Total', `${getTotalFee()} ${fiatCurrency}`)}
            </>
          ) : (
            <>
              {renderItem(
                'Application Fee',
                `${formatNumber(tokenizationApplicationFee)} ${tokenizationApplicationFeeAsset}${tokenizedAsset.assetTokenizationStatus === 1 ? ' (Paid)' : ''}`,
              )}
            </>
          )}
        </div>

        {/* Fees in Asset Section (only when vetted and feeInAsset > 0) */}
        {isVetted && (tokenizedAsset.feeInAsset ?? 0) > 0 && (
          <div className="rounded-xl border border-gray-200 bg-primary-100 p-4">
            <h3 className="text-base font-semibold text-[#0E1F51] mb-3 px-4">
              Fees in Asset
            </h3>
            {renderItem(
              'Tokenization Fee',
              `${formatNumberShort(tokenizedAsset.feeInAsset ?? 0)} ${tokenizedAsset.assetCode?.toUpperCase() ?? ''}`,
              true,
            )}
            {renderItem(
              'VAT (Asset)',
              `${formatNumberShort((tokenizedAsset.feeInAsset ?? 0) * 0.075)} ${tokenizedAsset.assetCode?.toUpperCase() ?? ''}`,
              true,
            )}
            {renderItem(
              'Total',
              `${formatNumberShort((tokenizedAsset.feeInAsset ?? 0) * 0.075 + (tokenizedAsset.feeInAsset ?? 0))} ${tokenizedAsset.assetCode?.toUpperCase() ?? ''}`,
            )}
          </div>
        )}

        {/* Info Notice */}
        {isVetted && isAlreadySubmitted && (
          <div className="px-5">
            <div className="flex items-start gap-2 p-3 rounded-lg border border-gray-200 bg-gray-50">
              <div className="flex-shrink-0 w-6 h-6 rounded-full border border-gray-500 flex items-center justify-center">
                <span className="text-xs font-bold text-gray-600">i</span>
              </div>
              <p className="text-xs text-gray-600">
                Items marked in asterisks (*) are not to be paid upfront.
              </p>
            </div>
          </div>
        )}

        {!isVetted && (
          <div className="px-5">
            <div className="flex items-start gap-2 p-3 rounded-lg border border-gray-200 bg-gray-50">
              <div className="flex-shrink-0 w-6 h-6 rounded-full border border-gray-500 flex items-center justify-center">
                <span className="text-xs font-bold text-gray-600">i</span>
              </div>
              <p className="text-xs text-gray-600">
                Please note that tokenization fee and other statutory fees will
                be displayed after vetting.
              </p>
            </div>
          </div>
        )}

        {/* Application Fee Notice */}
        {!isAlreadySubmitted && (
          <div className="px-5">
            <div className="p-3 rounded-lg border border-gray-200 bg-gray-50">
              <p className="text-xs text-gray-600">
                Please authorize the deduction of application fee to submit
                application.
              </p>
            </div>
          </div>
        )}

        {/* Password Form */}
        {!isAlreadySubmitted && (
          <div className="px-5 space-y-4">
            <div>
              <label className="block font-medium text-[#0E1F51] mb-1">
                Password
              </label>
              <input
                type="password"
                className="h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 outline-none focus:border-primary-600"
                placeholder="Enter your password"
                value={password}
                onChange={(e) => {
                  setPassword(e.target.value);
                  setPasswordError('');
                }}
              />
              {passwordError && (
                <p className="mt-1 text-xs text-red-500">{passwordError}</p>
              )}
            </div>
          </div>
        )}
      </div>

      {/* Buttons */}
      <div className="w-full flex justify-end mt-10">
        <div className="flex self-end justify-end w-full max-w-2xl space-x-6">
          {!isAlreadySubmitted && (
            <div className="w-2/4">
              <ButtonSecondary label="Back" onclick={() => navigate(-1)} />
            </div>
          )}
          <div className="w-3/4">
            {isAlreadySubmitted &&
            isVetted &&
            tokenizedAsset.assetTokenizationStatus === 1 ? (
              <Button
                label="Proceed to Pay"
                onclick={() => {
                  // Navigate to fee payment view
                  navigate(
                    `/dashboard/tokenize/fee-payment/${tokenizedAsset.id}`,
                  );
                }}
              />
            ) : (isAlreadySubmitted && !isVetted) ||
              tokenizedAsset.assetTokenizationStatus === 2 ? (
              <ButtonSecondary label="Back" onclick={() => navigate(-1)} />
            ) : (
              <Button
                label={isSubmitting ? 'Authorizing...' : 'Authorize'}
                onclick={() => void handlePasswordSubmit()}
                disabled={isSubmitting}
              />
            )}
          </div>
        </div>
      </div>

      {/* View Asset Details Link */}
      <div className="flex justify-center mt-4">
        <button
          type="button"
          className="text-sm font-semibold text-gray-700 underline hover:text-gray-900"
          onClick={() => {
            // Navigate to asset details view
            navigate(
              `/dashboard/tokenize/asset-dashboard/${tokenizedAsset.id}`,
            );
          }}
        >
          View Asset Details
        </button>
      </div>

      <Modal
        showModal={showSuccessModal}
        onClose={() => {
          setShowSuccessModal(false);
        }}
      >
        <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
          <img src="/images/launch.png" alt="success" />
          <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
            <p className="text-primary-800 text-md xl:text-lg font-semibold">
              Your Asset Tokenization application has been submitted
              successfully, please wait for vetting to be done, you will be
              notified once it is vetted so that you can proceed to pay the
              asset tokenization fee and other statutory fees.
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

      {/* Transaction Confirmation Modal */}
      <Modal
        showModal={showTxnModal}
        onClose={() => {
          handleDeclineTransaction();
        }}
      >
        <div className="flex flex-col items-center w-full px-5 md:px-10 space-y-5 py-6 justify-center">
          <p className="text-primary-800 w-full mb-2 text-lg text-center font-montserratSemiBold">
            Authorize
          </p>
          <div className="flex w-full py-3 px-2 xl:px-5 rounded-xl flex-col items-center bg-primary-100">
            <div className="flex flex-col space-y-2 items-center w-full py-3 justify-between">
              {isSubmitting ? (
                <p className="text-sm text-primary-800">
                  Authorizing transaction, please wait...
                </p>
              ) : (
                <>
                  {(pendingTransaction?.messages ?? []).length > 0 ? (
                    (pendingTransaction?.messages ?? []).map(
                      (msg: string, index: number) => (
                        <p
                          key={index}
                          className="text-sm text-primary-800 text-center"
                        >
                          {msg}
                        </p>
                      ),
                    )
                  ) : (
                    <p className="text-sm text-primary-800 text-center">
                      Do you want to authorize this transaction?
                    </p>
                  )}
                </>
              )}
            </div>
          </div>
          <div className="w-full space-y-3">
            <Button
              label={isSubmitting ? 'Authorizing...' : 'Continue'}
              disabled={isSubmitting}
              onclick={() => {
                void handleApproveTransaction();
              }}
            />
          </div>
          <div className="w-full space-y-3">
            <ButtonSecondary
              label="Decline"
              onclick={() => {
                handleDeclineTransaction();
              }}
            />
          </div>
        </div>
      </Modal>

      <div className="h-20" />
    </div>
  );
}
