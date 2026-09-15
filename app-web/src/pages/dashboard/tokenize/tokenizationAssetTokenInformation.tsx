import { useState, useRef, useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../../components/button';
import ButtonSecondary from '../../../components/buttonSecondary';
import TextInput from '../../../components/textInput';
import Dropdown from '../../../components/dropdown';
import {
  hideLoader,
  showLoader,
  showNotification,
} from '../../../utils/showToaster';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import { Encryptor } from '../../../utils/encryptor';
import { setActiveTokenizedAsset } from '../../../store/appStateSlice';
import {
  TokenizedAsset,
  excludedSubmissionFields,
} from '../../../types/tokenizedAsset';
import { useSubmitTokenizationMutation } from '../../../store/api/tokenizationApis';
import { useUploadAssetLogoMutation } from '../../../store/api/tokenizationApis';
import { useActiveTokenizedAsset } from './useActiveTokenizedAsset';
import { useTokenizationData } from './useTokenizationData';
import { formatNumber } from '../../../utils/fomatNumber';

type FormData = {
  assetName: string;
  assetCode: string;
  assetLogo: string | null;
  numberOfTokenToBeIssued: string;
  walletToHoldAssetsNotForSale: string;
  tokenizationFeeId: string;
  assetQuoteCurrency: string;
  salesStart: string;
  salesEnd: string;
  capOnPurchase: boolean;
  capAmountInFiat: string;
  capDurationInDays: string;
  proceedCycle: string;
  proceedPayoutCurrency: string;
  withholdingTaxDisclosure: boolean;
  bankId: string;
  accountNumber: string;
  beneficiaryName: string;
  exemptedCountries: string[];
  hasAdditionalKYCRequirements: boolean;
  additionalKYCRequirements: string;
  investorCategory: string;
  minimumKycTier: string;
  acknowledgedSuitabilityCriteria: boolean;
  authorizedRepresentativeName: string;
  authorizedRepresentativeTitleOrPosition: string;
  authorizedRepresentativeEmail: string;
  attestInformationAccurateAndVerifiable: boolean;
};

const INITIAL_FORM_DATA: FormData = {
  assetName: '',
  assetCode: '',
  assetLogo: null,
  numberOfTokenToBeIssued: '',
  walletToHoldAssetsNotForSale: '',
  tokenizationFeeId: '',
  assetQuoteCurrency: '',
  salesStart: '',
  salesEnd: '',
  capOnPurchase: false,
  capAmountInFiat: '',
  capDurationInDays: '',
  proceedCycle: '',
  proceedPayoutCurrency: '',
  withholdingTaxDisclosure: false,
  bankId: '',
  accountNumber: '',
  beneficiaryName: '',
  exemptedCountries: [],
  hasAdditionalKYCRequirements: false,
  additionalKYCRequirements: '',
  investorCategory: '',
  minimumKycTier: '',
  acknowledgedSuitabilityCriteria: false,
  authorizedRepresentativeName: '',
  authorizedRepresentativeTitleOrPosition: '',
  authorizedRepresentativeEmail: '',
  attestInformationAccurateAndVerifiable: false,
};

type TokenizedAssetFormPayload = Partial<TokenizedAsset>;

const INVESTOR_CATEGORY_OPTIONS = [
  { text: 'Retail', value: 'Retail' },
  { text: 'Qualified', value: 'Qualified' },
  { text: 'Institutional', value: 'Institutional' },
];

const KYC_TIER_OPTIONS = [
  { text: 'Tier 1', value: 'Tier 1' },
  { text: 'Tier 2', value: 'Tier 2' },
  { text: 'Tier 3', value: 'Tier 3' },
];

const EXEMPTED_COUNTRIES = [
  { code: 'US', name: 'United States' },
  { code: 'GB', name: 'United Kingdom' },
  { code: 'NG', name: 'Nigeria' },
  { code: 'CA', name: 'Canada' },
  { code: 'AU', name: 'Australia' },
  { code: 'DE', name: 'Germany' },
  { code: 'FR', name: 'France' },
  { code: 'JP', name: 'Japan' },
  { code: 'CH', name: 'Switzerland' },
  { code: 'SG', name: 'Singapore' },
  { code: 'AE', name: 'United Arab Emirates' },
  { code: 'ZA', name: 'South Africa' },
  { code: 'KE', name: 'Kenya' },
  { code: 'GH', name: 'Ghana' },
  { code: 'EG', name: 'Egypt' },
  { code: 'CN', name: 'China' },
  { code: 'IN', name: 'India' },
  { code: 'BR', name: 'Brazil' },
];

const BANKS = [
  { id: 0, name: 'Select Bank' },
  { id: 1, name: 'Access Bank Plc' },
  { id: 2, name: 'Zenith Bank Plc' },
  { id: 3, name: 'First Bank of Nigeria Plc' },
  { id: 4, name: 'GTBank Plc' },
  { id: 5, name: 'UBA Plc' },
  { id: 6, name: 'Stanbic IBTC Bank Plc' },
  { id: 7, name: 'Fidelity Bank Plc' },
  { id: 8, name: 'Polaris Bank Plc' },
  { id: 9, name: 'Union Bank of Nigeria Plc' },
  { id: 10, name: 'Wema Bank Plc' },
];

export function TokenizationAssetTokenInformation() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const {
    activeTokenizedAsset,
    assetId,
    isLoading: isAssetLoading,
    fetchError,
  } = useActiveTokenizedAsset();
  const { tokenizationData } = useTokenizationData();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [submitTokenization] = useSubmitTokenizationMutation();
  const [uploadAssetLogo] = useUploadAssetLogoMutation();

  const [formData, setFormData] = useState<FormData>(INITIAL_FORM_DATA);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [showCountryPicker, setShowCountryPicker] = useState(false);
  const fieldRefs = useRef<Map<string, HTMLElement>>(new Map());

  useEffect(() => {
    if (isAssetLoading && !activeTokenizedAsset) {
      setIsLoading(true);
      return;
    }

    const loadData = async () => {
      setIsLoading(true);
      setLoadError(null);
      try {
        if (activeTokenizedAsset) {
          setFormData((prev) => ({
            ...prev,
            assetName: activeTokenizedAsset.assetName ?? '',
            assetCode: activeTokenizedAsset.assetCode ?? '',
            assetLogo: activeTokenizedAsset.assetLogo || null,
            assetQuoteCurrency: activeTokenizedAsset.assetQuoteCurrency ?? '',
            proceedPayoutCurrency:
              activeTokenizedAsset.proceedPayoutCurrency ?? '',
            proceedCycle: activeTokenizedAsset.proceedCycle ?? '',
            minimumKycTier: activeTokenizedAsset.minimumKycTier ?? '',
            investorCategory: activeTokenizedAsset.investorCategory ?? '',
            bankId: activeTokenizedAsset.bankId
              ? String(activeTokenizedAsset.bankId)
              : '',
            accountNumber: activeTokenizedAsset.accountNumber ?? '',
            beneficiaryName: activeTokenizedAsset.beneficiaryName ?? '',
            exemptedCountries: String(
              activeTokenizedAsset.exemptedCountries ?? '',
            )
              .replaceAll(' ', '')
              .split(',')
              .filter((entry) => entry.trim().length > 0),
            hasAdditionalKYCRequirements:
              (activeTokenizedAsset.hasAdditionalKYCRequirements ?? 0) === 1,
            additionalKYCRequirements:
              activeTokenizedAsset.additionalKYCRequirements ?? '',
            authorizedRepresentativeName:
              activeTokenizedAsset.authorizedRepresentativeName ?? '',
            authorizedRepresentativeTitleOrPosition:
              activeTokenizedAsset.authorizedRepresentativeTitleOrPosition ??
              '',
            authorizedRepresentativeEmail:
              activeTokenizedAsset.authorizedRepresentativeEmail ?? '',
            withholdingTaxDisclosure:
              (activeTokenizedAsset.withholdingTaxDisclosure ?? 0) === 1,
            acknowledgedSuitabilityCriteria:
              (activeTokenizedAsset.acknowledgedSuitabilityCriteria ?? 0) === 1,
            attestInformationAccurateAndVerifiable:
              (activeTokenizedAsset.attestInformationAccurateAndVerifiable ??
                0) === 1,
            numberOfTokenToBeIssued:
              activeTokenizedAsset.numberOfTokenToBeIssued
                ? String(activeTokenizedAsset.numberOfTokenToBeIssued)
                : '',
            walletToHoldAssetsNotForSale:
              activeTokenizedAsset.walletToHoldAssetsNotForSale ?? '',
            tokenizationFeeId: activeTokenizedAsset.tokenizationFeeId
              ? String(activeTokenizedAsset.tokenizationFeeId)
              : '',
            salesStart: activeTokenizedAsset.salesStart
              ? new Date(activeTokenizedAsset.salesStart)
                  .toISOString()
                  .slice(0, 10)
              : '',
            salesEnd: activeTokenizedAsset.salesEnd
              ? new Date(activeTokenizedAsset.salesEnd)
                  .toISOString()
                  .slice(0, 10)
              : '',
            capOnPurchase: (activeTokenizedAsset.capOnPurchase ?? 0) === 1,
            capAmountInFiat: activeTokenizedAsset.capAmountInFiat
              ? String(activeTokenizedAsset.capAmountInFiat)
              : '',
            capDurationInDays: activeTokenizedAsset.capDurationInDays
              ? String(activeTokenizedAsset.capDurationInDays)
              : '',
          }));
          return;
        }

        setLoadError(
          fetchError ??
            'Unable to load the selected tokenized asset. Please reopen the page and try again.',
        );
      } catch {
        setLoadError('Unable to load the selected tokenized asset.');
        // silently fail
      } finally {
        setIsLoading(false);
      }
    };
    void loadData();
  }, [activeTokenizedAsset, fetchError, isAssetLoading]);

  const updateField = <K extends keyof FormData>(
    field: K,
    value: FormData[K],
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    // Clear error for this field on change
    if (errors[field]) {
      setErrors((prev) => {
        const next = { ...prev };
        delete next[field];
        return next;
      });
    }
  };

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};
    const f = formData;

    if (!f.assetCode.trim()) newErrors.assetCode = 'Asset code is required';
    if (!f.numberOfTokenToBeIssued || Number(f.numberOfTokenToBeIssued) <= 0) {
      newErrors.numberOfTokenToBeIssued =
        'Number of tokens to issue must be greater than 0';
    }
    if (!f.walletToHoldAssetsNotForSale) {
      newErrors.walletToHoldAssetsNotForSale =
        'Please select a wallet to hold assets not for sale';
    }
    if (!f.tokenizationFeeId || f.tokenizationFeeId === '0') {
      newErrors.tokenizationFeeId = 'Please select a tokenization fee';
    }
    if (!f.assetQuoteCurrency) {
      newErrors.assetQuoteCurrency = 'Please select a quote currency';
    }
    if (!f.salesStart) {
      newErrors.salesStart = 'Please select a sales start date';
    }
    if (!f.salesEnd) {
      newErrors.salesEnd = 'Please select a sales end date';
    }
    if (f.salesStart && f.salesEnd && f.salesEnd <= f.salesStart) {
      newErrors.salesEnd = 'Sales end date must be after sales start date';
    }
    if (f.capOnPurchase) {
      if (!f.capAmountInFiat || Number(f.capAmountInFiat) <= 0) {
        newErrors.capAmountInFiat = 'Cap amount must be greater than 0';
      }
      if (!f.capDurationInDays || Number(f.capDurationInDays) <= 0) {
        newErrors.capDurationInDays = 'Cap duration must be greater than 0';
      }
    }
    if (!f.proceedCycle) {
      newErrors.proceedCycle = 'Please select a payout cycle';
    }
    if (!f.proceedPayoutCurrency) {
      newErrors.proceedPayoutCurrency = 'Please select a payout currency';
    }
    if (!f.bankId || f.bankId === '0') {
      newErrors.bankId = 'Please select a bank';
    }
    if (!f.accountNumber.trim()) {
      newErrors.accountNumber = 'Account number is required';
    }
    if (!f.beneficiaryName.trim()) {
      newErrors.beneficiaryName = 'Beneficiary name is required';
    }
    if (!f.investorCategory) {
      newErrors.investorCategory = 'Please select an investor category';
    }
    if (f.hasAdditionalKYCRequirements && !f.additionalKYCRequirements.trim()) {
      newErrors.additionalKYCRequirements =
        'Please describe the additional KYC requirements';
    }
    if (!f.authorizedRepresentativeName.trim()) {
      newErrors.authorizedRepresentativeName =
        'Authorized representative name is required';
    }
    if (!f.authorizedRepresentativeTitleOrPosition.trim()) {
      newErrors.authorizedRepresentativeTitleOrPosition =
        'Title/position is required';
    }
    if (!f.authorizedRepresentativeEmail.trim()) {
      newErrors.authorizedRepresentativeEmail = 'Contact email is required';
    } else if (
      !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(f.authorizedRepresentativeEmail.trim())
    ) {
      newErrors.authorizedRepresentativeEmail =
        'Please enter a valid email address';
    }
    if (!f.attestInformationAccurateAndVerifiable) {
      newErrors.attestInformationAccurateAndVerifiable =
        'You must attest that the information is accurate';
    }
    if (!f.assetLogo) {
      newErrors.assetLogo = 'Please upload an asset logo';
    }
    if (!f.acknowledgedSuitabilityCriteria) {
      newErrors.acknowledgedSuitabilityCriteria =
        'You must acknowledge the suitability criteria';
    }

    setErrors(newErrors);
    const hasErrors = Object.keys(newErrors).length > 0;
    if (hasErrors) {
      scrollToFirstError(newErrors);
    }
    return !hasErrors;
  };

  const scrollToFirstError = (validationErrors: Record<string, string>) => {
    const firstErrorKey = Object.keys(validationErrors)[0];
    if (!firstErrorKey) {
      return;
    }

    const targetElement = fieldRefs.current.get(firstErrorKey);
    if (targetElement) {
      targetElement.scrollIntoView({
        behavior: 'smooth',
        block: 'center',
      });

      const focusable = targetElement.querySelector<HTMLElement>(
        'input, select, textarea, button',
      );
      focusable?.focus();
    }
  };

  const handleLogoUpload = async (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const file = event.target.files?.[0];
    if (!file) return;

    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }

    if (file.size > 900000) {
      setErrors((prev) => ({
        ...prev,
        assetLogo: 'File size must be less than 900KB',
      }));
      return;
    }

    const tokenizedAssetId = assetId ?? activeTokenizedAsset?.id ?? '';
    if (!tokenizedAssetId) {
      setErrors((prev) => ({
        ...prev,
        assetLogo:
          'Unable to find selected tokenized asset. Please reopen the application and try again.',
      }));
      return;
    }

    setIsUploading(true);
    setErrors((prev) => {
      const next = { ...prev };
      delete next.assetLogo;
      return next;
    });

    try {
      showLoader();
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);

      const res = await uploadAssetLogo({
        signer: appUser.primarySigner,
        publicKey: appUser.publicKey,
        secretKey: currentSecretKey,
        tokenizedAssetId,
        file,
      });

      if ('data' in res) {
        const uploadedLogoUrl = (res.data as any)?.data ?? res.data;
        const logoUrl = String(uploadedLogoUrl ?? '')
          .replaceAll('"', '')
          .trim();

        if (logoUrl) {
          updateField('assetLogo', logoUrl);
        } else {
          setErrors((prev) => ({
            ...prev,
            assetLogo: 'Failed to upload asset logo. Please try again.',
          }));
        }
      } else {
        setErrors((prev) => ({
          ...prev,
          assetLogo: 'Failed to upload asset logo. Please try again.',
        }));
      }
    } catch {
      setErrors((prev) => ({
        ...prev,
        assetLogo: 'Sorry, something went wrong. Please try again.',
      }));
    } finally {
      setIsUploading(false);
      hideLoader();
    }
  };

  const toggleExemptedCountry = (countryCode: string) => {
    setFormData((prev) => {
      const current = prev.exemptedCountries;
      const updated = current.includes(countryCode)
        ? current.filter((c) => c !== countryCode)
        : [...current, countryCode];
      return { ...prev, exemptedCountries: updated };
    });
  };

  const removeExemptedCountry = (countryCode: string) => {
    setFormData((prev) => ({
      ...prev,
      exemptedCountries: prev.exemptedCountries.filter(
        (c) => c !== countryCode,
      ),
    }));
  };

  function excludeSubmissionFields(payload: any) {
    return Object.entries(payload).reduce((accumulator: any, [key, value]) => {
      if (!excludedSubmissionFields.has(key)) {
        accumulator[key] = value as never;
      }

      return accumulator;
    }, {});
  }

  const handleSubmit = async () => {
    if (!validateForm()) return;

    setIsSubmitting(true);
    showLoader();

    try {
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);
      const existingPayload: TokenizedAssetFormPayload = {
        ...(activeTokenizedAsset ?? {}),
      };

      const formPayload = {
        id: existingPayload.id,
        assetCode: formData.assetCode,
        assetName: formData.assetName,
        assetLogo: formData.assetLogo,
        numberOfTokenToBeIssued: Number(formData.numberOfTokenToBeIssued),
        walletToHoldAssetsNotForSale: formData.walletToHoldAssetsNotForSale,
        tokenizationFeeId: Number(formData.tokenizationFeeId),
        assetQuoteCurrency: formData.assetQuoteCurrency,
        salesStart: formData.salesStart
          ? new Date(formData.salesStart).toISOString()
          : null,
        salesEnd: formData.salesEnd
          ? new Date(formData.salesEnd).toISOString()
          : null,
        capOnPurchase: formData.capOnPurchase ? 1 : 0,
        capAmountInFiat: formData.capOnPurchase
          ? Number(formData.capAmountInFiat)
          : 0,
        capDurationInDays: formData.capOnPurchase
          ? Number(formData.capDurationInDays)
          : 0,
        proceedCycle: formData.proceedCycle,
        proceedPayoutCurrency: formData.proceedPayoutCurrency,
        withholdingTaxDisclosure: formData.withholdingTaxDisclosure ? 1 : 0,
        bankId: Number(formData.bankId),
        accountNumber: formData.accountNumber,
        beneficiaryName: formData.beneficiaryName,
        exemptedCountries: formData.exemptedCountries.join(','),
        hasAdditionalKYCRequirements: formData.hasAdditionalKYCRequirements
          ? 1
          : 0,
        additionalKYCRequirements: formData.additionalKYCRequirements,
        investorCategory: formData.investorCategory,
        minimumKycTier: formData.minimumKycTier,
        acknowledgedSuitabilityCriteria:
          formData.acknowledgedSuitabilityCriteria ? 1 : 0,
        authorizedRepresentativeName: formData.authorizedRepresentativeName,
        authorizedRepresentativeTitleOrPosition:
          formData.authorizedRepresentativeTitleOrPosition,
        authorizedRepresentativeEmail: formData.authorizedRepresentativeEmail,
        attestInformationAccurateAndVerifiable:
          formData.attestInformationAccurateAndVerifiable ? 1 : 0,
      };

      const mergedPayload = {
        ...existingPayload,
        ...formPayload,
        id: existingPayload.id,
      };

      const payload = excludeSubmissionFields(mergedPayload);

      const res = await submitTokenization({
        signer: appUser.primarySigner,
        publicKey: appUser.publicKey,
        secretKey: currentSecretKey,
        body: payload,
      });

      hideLoader();

      if ('data' in res) {
        const savedAsset = res.data;

        if (savedAsset && typeof savedAsset === 'object') {
          dispatch(setActiveTokenizedAsset(savedAsset));
        }

        const savedAssetId =
          (savedAsset &&
            typeof savedAsset === 'object' &&
            'id' in savedAsset &&
            savedAsset.id) ||
          assetId ||
          '';

        showNotification('success', 'Token information saved successfully.');
        navigate(`/dashboard/tokenize/confirm-details/${savedAssetId}`);
        return;
      }

      if ('error' in res) {
        const errorData = (res as any).error;
        showNotification(
          'error',
          errorData?.data?.message ??
            'Unable to save token information. Please try again.',
        );
        return;
      }

      showNotification(
        'error',
        'Unable to save token information. Please try again.',
      );
    } catch (submitError: any) {
      hideLoader();
      showNotification(
        'error',
        submitError?.message ??
          'Unable to save token information. Please try again.',
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  const getCountryName = (code: string): string => {
    return EXEMPTED_COUNTRIES.find((c) => c.code === code)?.name ?? code;
  };

  const tokenizationFeeOptions = useMemo(() => {
    const fees = tokenizationData?.tokenizationFees ?? [];
    const numberOfTokens = Number(formData.numberOfTokenToBeIssued) || 0;
    const assetCurrentValue =
      Number(activeTokenizedAsset?.assetCurrentValue) || 0;
    const quoteCurrency =
      formData.assetQuoteCurrency ||
      activeTokenizedAsset?.assetQuoteCurrency ||
      '';
    const currentAssetCode =
      formData.assetCode || activeTokenizedAsset?.assetCode || '';

    return fees.map((fee, index) => {
      if (index === 0) {
        return {
          text: 'Select fee',
          value: String(fee.id),
        };
      }

      const fiatPercentage = Number(fee.feeFiatPercentage) || 0;
      const assetPercentage = Number(fee.feeAssetPercentage) || 0;
      const fiatFeeCap = Number(fee.feeFiatCap) || 0;
      const tokenFee = (numberOfTokens * assetPercentage) / 100;
      const fiatFee = (assetCurrentValue * fiatPercentage) / 100;
      const effectiveFiatFee = fiatFeeCap > fiatFee ? fiatFeeCap : fiatFee;

      return {
        text: `Option ${index} - ${quoteCurrency}${formatNumber(effectiveFiatFee)} + ${formatNumber(tokenFee)}${currentAssetCode ? ` ${currentAssetCode}` : ''}`,
        value: String(fee.id),
      };
    });
  }, [
    tokenizationData?.tokenizationFees,
    formData.numberOfTokenToBeIssued,
    formData.assetQuoteCurrency,
    formData.assetCode,
    activeTokenizedAsset?.assetCurrentValue,
    activeTokenizedAsset?.assetQuoteCurrency,
    activeTokenizedAsset?.assetCode,
  ]);

  const tokenizationCurrencyOptions = useMemo(() => {
    const currencies = tokenizationData?.tokenizationCurrencies ?? [];

    return currencies.map((currency) => ({
      text: String(currency.assetCode),
      value: String(currency.assetCode),
    }));
  }, [tokenizationData?.tokenizationCurrencies]);

  const payoutCycleOptions = useMemo(() => {
    const cycles = tokenizationData?.assetProceedCycle ?? [];

    return cycles.map((cycle) => ({
      text: String(cycle.id),
      value: String(cycle.id),
    }));
  }, [tokenizationData?.assetProceedCycle]);

  if (isLoading) {
    return (
      <div className="bg-white p-6 w-full mx-auto">
        <div className="rounded-lg border border-gray-200 bg-primary-100 p-6 text-sm text-gray-600">
          Loading token information...
        </div>
      </div>
    );
  }

  if (loadError && !activeTokenizedAsset) {
    return (
      <div className="bg-white p-6 w-full mx-auto">
        <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-600">
          {loadError}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white p-6 w-full mx-auto">
      <div className="space-y-6">
        {errors._form ? (
          <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-600">
            {errors._form}
          </div>
        ) : null}

        {/* Token Identity Section */}
        <div className="flex flex-col gap-6 md:flex-row md:gap-x-6 pb-8">
          <div className="w-full md:w-2/6">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Token Identity
            </label>
            <p className="text-sm text-gray-500 mb-2 max-w-sm">
              Basic identifying information for the token
            </p>
          </div>
          <div className="flex-1 w-full md:w-4/6 space-y-4">
            {/* Asset Name (read-only) */}
            <div>
              <label className="block font-medium text-[#0E1F51] mb-1">
                Asset Name
              </label>
              <TextInput
                label=""
                inputType="text"
                defaultValue={formData.assetName}
                placeholder="Enter the asset name"
                readonly
                onInputChange={(value) => updateField('assetName', value)}
              />
            </div>

            {/* Asset Code / Token Symbol */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('assetCode', node);
                } else {
                  fieldRefs.current.delete('assetCode');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Asset Code / Token Symbol{' '}
                <span className="text-red-500">*</span>
              </label>
              <TextInput
                label=""
                inputType="text"
                defaultValue={formData.assetCode}
                placeholder="e.g. AAF001"
                onInputChange={(value) => updateField('assetCode', value)}
              />
              {errors.assetCode ? (
                <p className="mt-1 text-xs text-red-500">{errors.assetCode}</p>
              ) : null}
            </div>

            {/* Asset Logo Upload */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('assetLogo', node);
                } else {
                  fieldRefs.current.delete('assetLogo');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Asset Logo <span className="text-red-500">*</span>
              </label>

              {formData.assetLogo ? (
                <div className="mt-2">
                  <img
                    src={formData.assetLogo}
                    alt="Asset logo"
                    className="max-w-xs max-h-32 rounded-lg border border-gray-200 cursor-pointer object-contain"
                    onClick={() => fileInputRef.current?.click()}
                  />
                  <button
                    type="button"
                    className="mt-1 text-xs text-red-500 hover:text-red-700 underline"
                    onClick={() => updateField('assetLogo', null)}
                  >
                    Remove
                  </button>
                </div>
              ) : (
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
                      {isUploading
                        ? 'Uploading...'
                        : 'Browse files to upload asset logo'}
                    </span>
                    <span className="text-xs text-gray-400">
                      Supported: JPG, JPEG, GIF, PNG, PDF (Max 900KB)
                    </span>
                  </div>
                </div>
              )}
              {errors.assetLogo ? (
                <p className="mt-1 text-xs text-red-500">{errors.assetLogo}</p>
              ) : null}
            </div>
          </div>
        </div>

        <hr className="border-gray-200" />

        {/* Token Supply Section */}
        <div className="flex flex-col gap-6 md:flex-row md:gap-x-6 pb-8">
          <div className="w-full md:w-2/6">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Token Supply
            </label>
            <p className="text-sm text-gray-500 mb-2 max-w-sm">
              Define the number of tokens to be issued for this asset
            </p>
          </div>
          <div className="flex-1 w-full md:w-4/6 space-y-4">
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('numberOfTokenToBeIssued', node);
                } else {
                  fieldRefs.current.delete('numberOfTokenToBeIssued');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Number of Tokens to Issue{' '}
                <span className="text-red-500">*</span>
              </label>
              <input
                type="number"
                className="h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 outline-none focus:border-primary-600"
                value={formData.numberOfTokenToBeIssued}
                placeholder="e.g. 1000000"
                onChange={(e) =>
                  updateField('numberOfTokenToBeIssued', e.target.value)
                }
              />
              {errors.numberOfTokenToBeIssued ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.numberOfTokenToBeIssued}
                </p>
              ) : null}
            </div>
          </div>
        </div>

        <hr className="border-gray-200" />

        {/* Wallet & Fees Section */}
        <div className="flex flex-col gap-6 md:flex-row md:gap-x-6 pb-8">
          <div className="w-full md:w-2/6">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Wallet & Fees
            </label>
            <p className="text-sm text-gray-500 mb-2 max-w-sm">
              Select the wallet to hold assets not for sale and choose a
              tokenization fee
            </p>
          </div>
          <div className="flex-1 w-full md:w-4/6 space-y-4">
            {/* Wallet to hold assets not for sale */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('walletToHoldAssetsNotForSale', node);
                } else {
                  fieldRefs.current.delete('walletToHoldAssetsNotForSale');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Wallet to Hold Assets Not for Sale{' '}
                <span className="text-red-500">*</span>
              </label>
              <Dropdown
                label="Select wallet"
                defaultValue={
                  appUser.userWallets.find(
                    (w) =>
                      w.publicKey === formData.walletToHoldAssetsNotForSale,
                  )
                    ? {
                        text: appUser.userWallets.find(
                          (w) =>
                            w.publicKey ===
                            formData.walletToHoldAssetsNotForSale,
                        )!.alias,
                        value: formData.walletToHoldAssetsNotForSale,
                      }
                    : null
                }
                options={appUser.userWallets.map((w) => {
                  return { text: w.alias, value: w.publicKey };
                })}
                onSelect={(option) =>
                  updateField('walletToHoldAssetsNotForSale', option.value)
                }
              />
              {errors.walletToHoldAssetsNotForSale ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.walletToHoldAssetsNotForSale}
                </p>
              ) : null}
            </div>

            {/* Tokenization Fee */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('tokenizationFeeId', node);
                } else {
                  fieldRefs.current.delete('tokenizationFeeId');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Tokenization Fee <span className="text-red-500">*</span>
              </label>
              <Dropdown
                label="Select fee"
                defaultValue={
                  formData.tokenizationFeeId
                    ? (tokenizationFeeOptions.find(
                        (o) => String(o.value) === formData.tokenizationFeeId,
                      ) ?? null)
                    : null
                }
                options={tokenizationFeeOptions}
                onSelect={(item) =>
                  updateField('tokenizationFeeId', String(item.value))
                }
              />
              {errors.tokenizationFeeId ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.tokenizationFeeId}
                </p>
              ) : null}
            </div>
          </div>
        </div>

        <hr className="border-gray-200" />

        {/* Asset Token Sale Section */}
        <div className="flex flex-col gap-6 md:flex-row md:gap-x-6 pb-8">
          <div className="w-full md:w-2/6">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Asset Token Sale
            </label>
            <p className="text-sm text-gray-500 mb-2 max-w-sm">
              Configure the token sale parameters
            </p>
          </div>
          <div className="flex-1 w-full md:w-4/6 space-y-4">
            {/* Quote Currency */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('assetQuoteCurrency', node);
                } else {
                  fieldRefs.current.delete('assetQuoteCurrency');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Asset Quote Currency <span className="text-red-500">*</span>
              </label>
              <Dropdown
                label="Select currency"
                defaultValue={
                  formData.assetQuoteCurrency
                    ? (tokenizationCurrencyOptions.find(
                        (o) => String(o.value) === formData.assetQuoteCurrency,
                      ) ?? null)
                    : null
                }
                options={tokenizationCurrencyOptions}
                onSelect={(item) =>
                  updateField('assetQuoteCurrency', String(item.value))
                }
              />
              {errors.assetQuoteCurrency ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.assetQuoteCurrency}
                </p>
              ) : null}
            </div>

            {/* Sale Dates */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div
                ref={(node) => {
                  if (node) {
                    fieldRefs.current.set('salesStart', node);
                  } else {
                    fieldRefs.current.delete('salesStart');
                  }
                }}
              >
                <label className="block font-medium text-[#0E1F51] mb-1">
                  Sales Starts From <span className="text-red-500">*</span>
                </label>
                <input
                  type="date"
                  className="h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 outline-none focus:border-primary-600"
                  value={formData.salesStart}
                  onChange={(e) => updateField('salesStart', e.target.value)}
                />
                {errors.salesStart ? (
                  <p className="mt-1 text-xs text-red-500">
                    {errors.salesStart}
                  </p>
                ) : null}
              </div>
              <div
                ref={(node) => {
                  if (node) {
                    fieldRefs.current.set('salesEnd', node);
                  } else {
                    fieldRefs.current.delete('salesEnd');
                  }
                }}
              >
                <label className="block font-medium text-[#0E1F51] mb-1">
                  Sales End On <span className="text-red-500">*</span>
                </label>
                <input
                  type="date"
                  className="h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 outline-none focus:border-primary-600"
                  value={formData.salesEnd}
                  onChange={(e) => updateField('salesEnd', e.target.value)}
                />
                {errors.salesEnd ? (
                  <p className="mt-1 text-xs text-red-500">{errors.salesEnd}</p>
                ) : null}
              </div>
            </div>

            {/* Cap on Purchase Toggle */}
            <div className="flex items-center justify-between pt-2">
              <div>
                <label className="block font-medium text-[#0E1F51] mb-1">
                  Cap on Purchase
                </label>
                <p className="text-xs text-gray-500 max-w-sm">
                  Limit the purchase amount for each investor
                </p>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  className="sr-only peer"
                  checked={formData.capOnPurchase}
                  onChange={() =>
                    updateField('capOnPurchase', !formData.capOnPurchase)
                  }
                />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
              </label>
            </div>

            {/* Conditional Cap Fields */}
            {formData.capOnPurchase ? (
              <div className="space-y-4 pl-0 md:pl-4 border-l-0 md:border-l-2 border-primary-200 md:pl-4">
                <div
                  ref={(node) => {
                    if (node) {
                      fieldRefs.current.set('capAmountInFiat', node);
                    } else {
                      fieldRefs.current.delete('capAmountInFiat');
                    }
                  }}
                >
                  <label className="block font-medium text-[#0E1F51] mb-1">
                    Cap Amount (in Fiat) <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    className="h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 outline-none focus:border-primary-600"
                    value={formData.capAmountInFiat}
                    placeholder="e.g. 100000"
                    onChange={(e) =>
                      updateField('capAmountInFiat', e.target.value)
                    }
                  />
                  {errors.capAmountInFiat ? (
                    <p className="mt-1 text-xs text-red-500">
                      {errors.capAmountInFiat}
                    </p>
                  ) : null}
                </div>
                <div
                  ref={(node) => {
                    if (node) {
                      fieldRefs.current.set('capDurationInDays', node);
                    } else {
                      fieldRefs.current.delete('capDurationInDays');
                    }
                  }}
                >
                  <label className="block font-medium text-[#0E1F51] mb-1">
                    Cap Duration (in Days){' '}
                    <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="number"
                    className="h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 outline-none focus:border-primary-600"
                    value={formData.capDurationInDays}
                    placeholder="e.g. 30"
                    onChange={(e) =>
                      updateField('capDurationInDays', e.target.value)
                    }
                  />
                  {errors.capDurationInDays ? (
                    <p className="mt-1 text-xs text-red-500">
                      {errors.capDurationInDays}
                    </p>
                  ) : null}
                </div>
              </div>
            ) : null}
          </div>
        </div>

        <hr className="border-gray-200" />

        {/* Proceeds Payout Section */}
        <div className="flex flex-col gap-6 md:flex-row md:gap-x-6 pb-8">
          <div className="w-full md:w-2/6">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Proceeds Payout
            </label>
            <p className="text-sm text-gray-500 mb-2 max-w-sm">
              Configure how sale proceeds will be paid out
            </p>
          </div>
          <div className="flex-1 w-full md:w-4/6 space-y-4">
            {/* Payout Cycle */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('proceedCycle', node);
                } else {
                  fieldRefs.current.delete('proceedCycle');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Proceeds Payout Cycle <span className="text-red-500">*</span>
              </label>
              <Dropdown
                label="Select payout cycle"
                defaultValue={
                  formData.proceedCycle
                    ? (payoutCycleOptions.find(
                        (o) => String(o.value) === formData.proceedCycle,
                      ) ?? null)
                    : null
                }
                options={payoutCycleOptions}
                onSelect={(item) =>
                  updateField('proceedCycle', String(item.value))
                }
              />
              {errors.proceedCycle ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.proceedCycle}
                </p>
              ) : null}
            </div>

            {/* Payout Currency */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('proceedPayoutCurrency', node);
                } else {
                  fieldRefs.current.delete('proceedPayoutCurrency');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Proceeds Payout Currency <span className="text-red-500">*</span>
              </label>
              <Dropdown
                label="Select payout currency"
                defaultValue={
                  formData.proceedPayoutCurrency
                    ? (tokenizationCurrencyOptions.find(
                        (o) =>
                          String(o.value) === formData.proceedPayoutCurrency,
                      ) ?? null)
                    : null
                }
                options={tokenizationCurrencyOptions}
                onSelect={(item) =>
                  updateField('proceedPayoutCurrency', String(item.value))
                }
              />
              {errors.proceedPayoutCurrency ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.proceedPayoutCurrency}
                </p>
              ) : null}
            </div>

            {/* Withholding Tax Disclosure */}
            <label className="flex items-center gap-3 cursor-pointer pt-2">
              <input
                type="checkbox"
                className="accent-primary-700 w-4 h-4"
                checked={formData.withholdingTaxDisclosure}
                onChange={() =>
                  updateField(
                    'withholdingTaxDisclosure',
                    !formData.withholdingTaxDisclosure,
                  )
                }
              />
              <span className="text-sm text-gray-700">
                Withholding Tax Disclosure
              </span>
            </label>
          </div>
        </div>

        <hr className="border-gray-200" />

        {/* Asset Sale Proceeds Section */}
        <div className="flex flex-col gap-6 md:flex-row md:gap-x-6 pb-8">
          <div className="w-full md:w-2/6">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Asset Sale Proceeds
            </label>
            <p className="text-sm text-gray-500 mb-2 max-w-sm">
              Provide bank account details for receiving sale proceeds
            </p>
          </div>
          <div className="flex-1 w-full md:w-4/6 space-y-4">
            {/* Bank Selection */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('bankId', node);
                } else {
                  fieldRefs.current.delete('bankId');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Select Bank <span className="text-red-500">*</span>
              </label>
              <Dropdown
                label="Please select bank"
                options={BANKS.map((bank) => ({
                  text: bank.name,
                  value: String(bank.id),
                }))}
                onSelect={(value) => updateField('bankId', value)}
              />
              {errors.bankId ? (
                <p className="mt-1 text-xs text-red-500">{errors.bankId}</p>
              ) : null}
            </div>

            {/* Account Number */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('accountNumber', node);
                } else {
                  fieldRefs.current.delete('accountNumber');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Account Number <span className="text-red-500">*</span>
              </label>
              <TextInput
                label=""
                inputType="text"
                defaultValue={formData.accountNumber}
                placeholder="Enter account number"
                onInputChange={(value) => updateField('accountNumber', value)}
              />
              {errors.accountNumber ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.accountNumber}
                </p>
              ) : null}
            </div>

            {/* Beneficiary Name */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('beneficiaryName', node);
                } else {
                  fieldRefs.current.delete('beneficiaryName');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Beneficiary Name <span className="text-red-500">*</span>
              </label>
              <TextInput
                label=""
                inputType="text"
                defaultValue={formData.beneficiaryName}
                placeholder="Enter beneficiary name"
                onInputChange={(value) => updateField('beneficiaryName', value)}
              />
              {errors.beneficiaryName ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.beneficiaryName}
                </p>
              ) : null}
            </div>
          </div>
        </div>

        <hr className="border-gray-200" />

        {/* Primary Buyer Requirements Section */}
        <div className="flex flex-col gap-6 md:flex-row md:gap-x-6 pb-8">
          <div className="w-full md:w-2/6">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Primary Buyer Requirements
            </label>
            <p className="text-sm text-gray-500 mb-2 max-w-sm">
              Define buyer eligibility requirements
            </p>
          </div>
          <div className="flex-1 w-full md:w-4/6 space-y-4">
            {/* Exempted Countries */}
            <div>
              <label className="block font-medium text-[#0E1F51] mb-1">
                Exempted Countries
              </label>
              <p className="text-xs text-gray-500 mb-2">
                Select countries that are exempted from purchasing this token
              </p>

              {/* Country Selection */}
              <div className="relative">
                <button
                  type="button"
                  className="h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 text-left flex items-center justify-between outline-none focus:border-primary-600"
                  onClick={() => setShowCountryPicker(!showCountryPicker)}
                >
                  <span className="text-sm">
                    {formData.exemptedCountries.length > 0
                      ? `${formData.exemptedCountries.length} country(s) selected`
                      : 'Select countries'}
                  </span>
                  <svg
                    className={`w-5 h-5 text-gray-400 transition-transform ${showCountryPicker ? 'rotate-180' : ''}`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M19 9l-7 7-7-7"
                    />
                  </svg>
                </button>

                {showCountryPicker ? (
                  <div className="absolute z-10 mt-1 w-full bg-white border border-gray-200 rounded-md shadow-lg max-h-60 overflow-y-auto">
                    {EXEMPTED_COUNTRIES.map((country) => (
                      <label
                        key={country.code}
                        className="flex items-center gap-2 px-3 py-2 hover:bg-primary-50 cursor-pointer text-sm"
                      >
                        <input
                          type="checkbox"
                          className="accent-primary-700"
                          checked={formData.exemptedCountries.includes(
                            country.code,
                          )}
                          onChange={() => toggleExemptedCountry(country.code)}
                        />
                        {country.name}
                      </label>
                    ))}
                  </div>
                ) : null}
              </div>

              {/* Selected Country Tags */}
              {formData.exemptedCountries.length > 0 ? (
                <div className="mt-2 flex flex-wrap gap-2">
                  {formData.exemptedCountries.map((code) => (
                    <span
                      key={code}
                      className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-primary-100 text-primary-700"
                    >
                      {getCountryName(code)}
                      <button
                        type="button"
                        className="text-primary-500 hover:text-primary-700"
                        onClick={() => removeExemptedCountry(code)}
                      >
                        <svg
                          className="w-3.5 h-3.5"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M6 18L18 6M6 6l12 12"
                          />
                        </svg>
                      </button>
                    </span>
                  ))}
                </div>
              ) : null}
            </div>

            {/* Additional KYC Requirements */}
            <div>
              <label className="block font-medium text-[#0E1F51] mb-1">
                Are there additional KYC requirements?
              </label>
              <div className="flex gap-6 mt-2">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="radio"
                    className="accent-primary-700"
                    checked={formData.hasAdditionalKYCRequirements === true}
                    onChange={() =>
                      updateField('hasAdditionalKYCRequirements', true)
                    }
                  />
                  <span className="text-sm text-gray-700">Yes</span>
                </label>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="radio"
                    className="accent-primary-700"
                    checked={formData.hasAdditionalKYCRequirements === false}
                    onChange={() =>
                      updateField('hasAdditionalKYCRequirements', false)
                    }
                  />
                  <span className="text-sm text-gray-700">No</span>
                </label>
              </div>
            </div>

            {/* Conditional Additional KYC Description */}
            {formData.hasAdditionalKYCRequirements ? (
              <div
                ref={(node) => {
                  if (node) {
                    fieldRefs.current.set('additionalKYCRequirements', node);
                  } else {
                    fieldRefs.current.delete('additionalKYCRequirements');
                  }
                }}
              >
                <label className="block font-medium text-[#0E1F51] mb-1">
                  Describe additional KYC requirements{' '}
                  <span className="text-red-500">*</span>
                </label>
                <textarea
                  className="w-full rounded-md border border-gray-200 bg-primary-100 px-3 py-2 text-gray-700 outline-none focus:border-primary-600 min-h-[80px] resize-y"
                  value={formData.additionalKYCRequirements}
                  placeholder="Describe specific documents or identity checks required"
                  onChange={(e) =>
                    updateField('additionalKYCRequirements', e.target.value)
                  }
                  rows={3}
                />
                {errors.additionalKYCRequirements ? (
                  <p className="mt-1 text-xs text-red-500">
                    {errors.additionalKYCRequirements}
                  </p>
                ) : null}
              </div>
            ) : null}

            {/* Investor Category */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('investorCategory', node);
                } else {
                  fieldRefs.current.delete('investorCategory');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Investor Category Eligibility{' '}
                <span className="text-red-500">*</span>
              </label>
              <p className="text-xs text-gray-500 mb-2">
                Define categories of eligible investors
              </p>
              <Dropdown
                label="Select option"
                defaultValue={
                  formData.investorCategory
                    ? (INVESTOR_CATEGORY_OPTIONS.find(
                        (o) => o.value === formData.investorCategory,
                      ) ?? null)
                    : null
                }
                options={INVESTOR_CATEGORY_OPTIONS}
                onSelect={(option) =>
                  updateField('investorCategory', option.value)
                }
              />
              {errors.investorCategory ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.investorCategory}
                </p>
              ) : null}
            </div>

            {/* Minimum KYC Tier */}
            <div>
              <label className="block font-medium text-[#0E1F51] mb-1">
                Minimum KYC Tier Required
              </label>
              <p className="text-xs text-gray-500 mb-2">
                Select platform verification tier required to invest
              </p>
              <Dropdown
                label="Select option"
                defaultValue={
                  formData.minimumKycTier
                    ? (KYC_TIER_OPTIONS.find(
                        (o) => o.value === formData.minimumKycTier,
                      ) ?? null)
                    : null
                }
                options={KYC_TIER_OPTIONS}
                onSelect={(option) =>
                  updateField('minimumKycTier', option.value)
                }
              />
            </div>
          </div>
        </div>

        <hr className="border-gray-200" />

        {/* Investment Suitability Disclaimer Section */}
        <div className="flex flex-col gap-6 md:flex-row md:gap-x-6 pb-8">
          <div className="w-full md:w-2/6">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Investment Suitability Disclaimer
            </label>
          </div>
          <div
            className="flex-1 w-full md:w-4/6"
            ref={(node) => {
              if (node) {
                fieldRefs.current.set('acknowledgedSuitabilityCriteria', node);
              } else {
                fieldRefs.current.delete('acknowledgedSuitabilityCriteria');
              }
            }}
          >
            <label className="flex items-start gap-3 cursor-pointer">
              <input
                type="checkbox"
                className="accent-primary-700 w-4 h-4 mt-0.5"
                checked={formData.acknowledgedSuitabilityCriteria}
                onChange={() =>
                  updateField(
                    'acknowledgedSuitabilityCriteria',
                    !formData.acknowledgedSuitabilityCriteria,
                  )
                }
              />
              <span className="text-sm text-gray-700">
                I acknowledge that I have reviewed the suitability criteria for
                this investment
              </span>
            </label>
            {errors.acknowledgedSuitabilityCriteria ? (
              <p className="mt-1 text-xs text-red-500">
                {errors.acknowledgedSuitabilityCriteria}
              </p>
            ) : null}
          </div>
        </div>

        <hr className="border-gray-200" />

        {/* Final Declaration and Attestation Section */}
        <div className="flex flex-col gap-6 md:flex-row md:gap-x-6 pb-8">
          <div className="w-full md:w-2/6">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Final Declaration and Attestation
            </label>
          </div>
          <div className="flex-1 w-full md:w-4/6 space-y-4">
            {/* Authorized Representative Name */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('authorizedRepresentativeName', node);
                } else {
                  fieldRefs.current.delete('authorizedRepresentativeName');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Authorized Representative Name{' '}
                <span className="text-red-500">*</span>
              </label>
              <TextInput
                label=""
                inputType="text"
                defaultValue={formData.authorizedRepresentativeName}
                placeholder="Enter authorized representative name"
                onInputChange={(value) =>
                  updateField('authorizedRepresentativeName', value)
                }
              />
              {errors.authorizedRepresentativeName ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.authorizedRepresentativeName}
                </p>
              ) : null}
            </div>

            {/* Position/Title */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set(
                    'authorizedRepresentativeTitleOrPosition',
                    node,
                  );
                } else {
                  fieldRefs.current.delete(
                    'authorizedRepresentativeTitleOrPosition',
                  );
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Position/Title of Authorized Representative{' '}
                <span className="text-red-500">*</span>
              </label>
              <TextInput
                label=""
                inputType="text"
                defaultValue={formData.authorizedRepresentativeTitleOrPosition}
                placeholder="Enter position or title"
                onInputChange={(value) =>
                  updateField('authorizedRepresentativeTitleOrPosition', value)
                }
              />
              {errors.authorizedRepresentativeTitleOrPosition ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.authorizedRepresentativeTitleOrPosition}
                </p>
              ) : null}
            </div>

            {/* Contact Email */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('authorizedRepresentativeEmail', node);
                } else {
                  fieldRefs.current.delete('authorizedRepresentativeEmail');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Contact Email of Authorized Representative{' '}
                <span className="text-red-500">*</span>
              </label>
              <TextInput
                label=""
                inputType="text"
                defaultValue={formData.authorizedRepresentativeEmail}
                placeholder="Enter contact email"
                onInputChange={(value) =>
                  updateField('authorizedRepresentativeEmail', value)
                }
              />
              {errors.authorizedRepresentativeEmail ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.authorizedRepresentativeEmail}
                </p>
              ) : null}
            </div>

            {/* Attestation Checkbox */}
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set(
                    'attestInformationAccurateAndVerifiable',
                    node,
                  );
                } else {
                  fieldRefs.current.delete(
                    'attestInformationAccurateAndVerifiable',
                  );
                }
              }}
            >
              <label className="flex items-start gap-3 cursor-pointer pt-2">
                <input
                  type="checkbox"
                  className="accent-primary-700 w-4 h-4 mt-0.5"
                  checked={formData.attestInformationAccurateAndVerifiable}
                  onChange={() =>
                    updateField(
                      'attestInformationAccurateAndVerifiable',
                      !formData.attestInformationAccurateAndVerifiable,
                    )
                  }
                />
                <span className="text-sm text-gray-700">
                  I attest that the information provided is accurate and
                  verifiable
                </span>
              </label>
              {errors.attestInformationAccurateAndVerifiable ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.attestInformationAccurateAndVerifiable}
                </p>
              ) : null}
            </div>
          </div>
        </div>
      </div>

      {/* Hidden file input for logo upload */}
      <input
        ref={fileInputRef}
        type="file"
        className="hidden"
        onChange={handleLogoUpload}
        accept=".jpg,.jpeg,.gif,.png,.pdf,.svg"
      />

      {/* Navigation Buttons */}
      <div className="w-full flex justify-end mt-10">
        <div className="flex self-end w-full max-w-2xl space-x-6">
          <div className="w-2/4">
            <ButtonSecondary
              label="Back"
              onclick={() => {
                navigate(-1);
              }}
            />
          </div>
          <div className="w-3/4">
            <Button
              label={isSubmitting ? 'Saving...' : 'Save and Continue'}
              onclick={() => {
                void handleSubmit();
              }}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
