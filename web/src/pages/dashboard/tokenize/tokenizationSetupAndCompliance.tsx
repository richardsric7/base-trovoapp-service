import { useEffect, useMemo, useRef, useState } from 'react';
import Dropdown from '../../../components/dropdown';
import Button from '../../../components/button';
import { useNavigate } from 'react-router-dom';
import ButtonSecondary from '../../../components/buttonSecondary';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import {
  excludedSubmissionFields,
  TokenizedAsset,
} from '../../../types/tokenizedAsset';
import { Encryptor } from '../../../utils/encryptor';
import {
  hideLoader,
  showLoader,
  showNotification,
} from '../../../utils/showToaster';
import { useSubmitTokenizationMutation } from '../../../store/api/tokenizationApis';
import { setActiveTokenizedAsset } from '../../../store/appStateSlice';
import { useActiveTokenizedAsset } from './useActiveTokenizedAsset';
import { useTokenizationData } from './useTokenizationData';

type OptionItem = {
  text: string;
  value: string;
};
type TokenizedAssetFormPayload = Partial<TokenizedAsset>;

const TOKENIZATION_REQUIREMENTS_URL =
  'https://tokenization-requirements-app-xu8c6.ondigitalocean.app';

const AVAILABLE_FINANCIAL_ASSET_SUBSECTORS = [
  'Investment Funds / Collective Investment Schemes',
  'Asset-Backed and Securitized Products',
  'Capital Markets - Equity & Fixed Income',
  'Commodity Markets',
];

export function TokenizationSetupAndCompliance() {
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const dispatch = useDispatch();
  const { activeTokenizedAsset, assetId } = useActiveTokenizedAsset();
  const { tokenizationData } = useTokenizationData();
  const availableFinancialAssetTypes = useSelector(
    (state: RootState) => state.appState.availableFinancialAssetTypes,
  );

  const [selectedAssetSectorId, setSelectedAssetSectorId] = useState('');
  const [selectedAssetSubSectorId, setSelectedAssetSubSectorId] = useState('');
  const [selectedAssetTypeId, setSelectedAssetTypeId] = useState('');
  const [selectedCountry, setSelectedCountry] = useState('');

  const [fundingStructure, setFundingStructure] = useState('equity');
  const [assetStatus, setAssetStatus] = useState('existing');
  const [offeringType, setOfferingType] = useState('public');
  const [hasDocuments, setHasDocuments] = useState('yes');
  const [confirm, setConfirm] = useState(true);
  const [
    acceptTokenizationTermsAndAgreement,
    setAcceptTokenizationTermsAndAgreement,
  ] = useState(false);
  const [equityPercentage, setEquityPercentage] = useState('');
  const [debtPercentage, setDebtPercentage] = useState('');
  const [proceedPayoutCurrency, setProceedPayoutCurrency] = useState('');
  const [assetQuoteCurrency, setAssetQuoteCurrency] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const fieldRefs = useRef<Map<string, HTMLElement>>(new Map());

  const [submitTokenization] = useSubmitTokenizationMutation();
  const navigate = useNavigate();

  const isFinanceAndInvestmentMarketsSector =
    selectedAssetSectorId.trim().toLowerCase() ===
    'finance and investment markets';

  const assetSectorOptions = useMemo<OptionItem[]>(() => {
    if (!tokenizationData?.assetSectors) {
      return [];
    }

    return tokenizationData.assetSectors.map((sector) => ({
      text: sector.sector,
      value: sector.sector,
    }));
  }, [tokenizationData]);

  const assetSubSectorOptions = useMemo<OptionItem[]>(() => {
    if (!tokenizationData?.assetSubSectors || !selectedAssetSectorId) {
      return [];
    }

    return tokenizationData.assetSubSectors
      .filter((item) => item.assetSectorId === selectedAssetSectorId)
      .filter((item) => {
        if (!isFinanceAndInvestmentMarketsSector) {
          return true;
        }

        return AVAILABLE_FINANCIAL_ASSET_SUBSECTORS.includes(item.subSector);
      })
      .map((item) => ({
        text: item.subSector,
        value: item.subSector,
      }));
  }, [
    tokenizationData,
    selectedAssetSectorId,
    isFinanceAndInvestmentMarketsSector,
  ]);

  const assetTypeOptions = useMemo<OptionItem[]>(() => {
    if (!tokenizationData?.assetTypes || !selectedAssetSubSectorId) {
      return [];
    }

    return tokenizationData.assetTypes
      .filter((item) => item.assetSubSectorId === selectedAssetSubSectorId)
      .filter((item) => {
        if (!isFinanceAndInvestmentMarketsSector) {
          return true;
        }

        return availableFinancialAssetTypes.includes(String(item.id));
      })
      .map((item) => ({
        text: item.assetType,
        value: String(item.id),
      }))
      .sort((a, b) => a.text.localeCompare(b.text));
  }, [
    tokenizationData,
    selectedAssetSubSectorId,
    isFinanceAndInvestmentMarketsSector,
    availableFinancialAssetTypes,
  ]);

  const countryOptions = useMemo<OptionItem[]>(() => {
    if (!tokenizationData?.countries) {
      return [];
    }

    const countries = tokenizationData.countries.filter(
      (item) => item.countryCode.toUpperCase() === 'NG',
    );

    const source =
      countries.length > 0 ? countries : tokenizationData.countries;

    return source.map((country) => ({
      text: country.countryName,
      value: country.countryCode,
    }));
  }, [tokenizationData]);

  const selectedSectorOption = useMemo(
    () =>
      assetSectorOptions.find((item) => item.value === selectedAssetSectorId),
    [assetSectorOptions, selectedAssetSectorId],
  );

  const selectedSubSectorOption = useMemo(
    () =>
      assetSubSectorOptions.find(
        (item) => item.value === selectedAssetSubSectorId,
      ),
    [assetSubSectorOptions, selectedAssetSubSectorId],
  );

  const selectedTypeOption = useMemo(
    () => assetTypeOptions.find((item) => item.value === selectedAssetTypeId),
    [assetTypeOptions, selectedAssetTypeId],
  );

  const selectedCountryOption = useMemo(
    () => countryOptions.find((item) => item.value === selectedCountry),
    [countryOptions, selectedCountry],
  );

  useEffect(() => {
    if (!activeTokenizedAsset) {
      return;
    }

    setSelectedAssetSectorId(activeTokenizedAsset.assetSector ?? '');
    setSelectedAssetSubSectorId(activeTokenizedAsset.assetSubSector ?? '');
    setSelectedAssetTypeId(String(activeTokenizedAsset.assetType ?? ''));
    setSelectedCountry(activeTokenizedAsset.assetCountryLocation ?? '');

    setAssetStatus(
      (activeTokenizedAsset.assetAlreadyExists ?? 0) === 1
        ? 'existing'
        : 'toBuild',
    );

    const incomingFundingStructure = Number(
      activeTokenizedAsset.fundingStructure ?? 0,
    );
    if (incomingFundingStructure === 2) {
      setFundingStructure('hybrid');
    } else if (incomingFundingStructure === 1) {
      setFundingStructure('debt');
    } else {
      setFundingStructure('equity');
    }

    setEquityPercentage(
      activeTokenizedAsset.equityPercentage !== undefined
        ? String(activeTokenizedAsset.equityPercentage)
        : '',
    );
    setDebtPercentage(
      activeTokenizedAsset.debtPercentage !== undefined
        ? String(activeTokenizedAsset.debtPercentage)
        : '',
    );

    setOfferingType(
      activeTokenizedAsset.offeringType === 'private' ? 'private' : 'public',
    );
    setHasDocuments(
      (activeTokenizedAsset.assetManagerId ?? 1) === 0 ? 'no' : 'yes',
    );
    setConfirm((activeTokenizedAsset.agreeTransferTitleToCustodian ?? 0) === 1);
    setAcceptTokenizationTermsAndAgreement(
      (activeTokenizedAsset.acceptTokenizationTermsAndAgreement ?? 0) === 1,
    );

    setProceedPayoutCurrency(activeTokenizedAsset.proceedPayoutCurrency ?? '');
    setAssetQuoteCurrency(activeTokenizedAsset.assetQuoteCurrency ?? '');
  }, [activeTokenizedAsset]);

  useEffect(() => {
    if (!selectedCountry || !tokenizationData?.countryConfigs) {
      return;
    }

    const config = tokenizationData.countryConfigs.find(
      (item) =>
        item.countryCode.toLowerCase() === selectedCountry.trim().toLowerCase(),
    );

    if (config) {
      setProceedPayoutCurrency(config.quoteCurrencyCode ?? '');
      setAssetQuoteCurrency(config.quoteCurrencyCode ?? '');
    }
  }, [selectedCountry, tokenizationData]);

  const getRequirementsUrl = () => {
    if (!selectedAssetSectorId) {
      return '';
    }

    const normalizedSector = selectedAssetSectorId
      .trim()
      .replace(/\s+/g, '-')
      .toLowerCase();

    if (isFinanceAndInvestmentMarketsSector) {
      const formName = String(
        (activeTokenizedAsset as any)?.assetFormName ?? '',
      )
        .trim()
        .replace(/\s+/g, '-')
        .toLowerCase();

      if (formName) {
        return `${TOKENIZATION_REQUIREMENTS_URL}/#/${normalizedSector}/${formName}`;
      }

      return `${TOKENIZATION_REQUIREMENTS_URL}/#/${normalizedSector}`;
    }

    const statusPath =
      assetStatus === 'existing' ? 'existing-assets' : 'non-existing-assets';
    return `${TOKENIZATION_REQUIREMENTS_URL}/#/${normalizedSector}/${statusPath}`;
  };

  const validate = () => {
    const nextErrors: Record<string, string> = {};

    if (!selectedAssetSectorId) {
      nextErrors.assetSector = 'Please select an asset sector';
    }

    if (!selectedAssetSubSectorId) {
      nextErrors.assetSubSector = 'Please select an asset sub-sector';
    }

    if (!selectedAssetTypeId) {
      nextErrors.assetType = 'Please select an asset type';
    }

    if (!selectedCountry) {
      nextErrors.country = 'Please select country location';
    }

    if (hasDocuments !== 'yes') {
      nextErrors.documents =
        'You need to acquire all the documents listed in the tokenization requirements document before you can proceed.';
    }

    if (!confirm) {
      nextErrors.confirm =
        'Tokenizing your asset requires ownership transfer of the asset to a licensed trustee. You must agree to continue.';
    }

    if (!acceptTokenizationTermsAndAgreement) {
      nextErrors.acceptTerms =
        'You must accept the tokenization terms and conditions.';
    }

    if (!isFinanceAndInvestmentMarketsSector && fundingStructure === 'hybrid') {
      const equityValue = Number(equityPercentage);
      const debtValue = Number(debtPercentage);

      if (!equityPercentage || Number.isNaN(equityValue)) {
        nextErrors.equityPercentage = 'Equity percentage is required';
      } else if (equityValue < 0 || equityValue > 100) {
        nextErrors.equityPercentage =
          'Equity percentage must be between 0 and 100';
      }

      if (!debtPercentage || Number.isNaN(debtValue)) {
        nextErrors.debtPercentage = 'Debt percentage is required';
      } else if (debtValue < 0 || debtValue > 100) {
        nextErrors.debtPercentage = 'Debt percentage must be between 0 and 100';
      }
    }

    setErrors(nextErrors);
    const hasErrors = Object.keys(nextErrors).length > 0;
    if (hasErrors) {
      scrollToFirstError(nextErrors);
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

  const clearError = (field: string) => {
    if (!errors[field]) {
      return;
    }

    setErrors((current) => {
      const next = { ...current };
      delete next[field];
      return next;
    });
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
    if (!tokenizationData) {
      showNotification(
        'error',
        'Tokenization metadata is unavailable. Refresh and try again.',
      );
      return;
    }

    if (!validate()) {
      return;
    }

    setIsSubmitting(true);
    showLoader();

    try {
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);
      const fundingStructureCode =
        fundingStructure === 'hybrid' ? 2 : fundingStructure === 'debt' ? 1 : 0;

      const existingPayload: TokenizedAssetFormPayload = {
        ...(activeTokenizedAsset ?? {}),
      };

      const payloadBody = {
        ...existingPayload,
        id: activeTokenizedAsset?.id,
        assetSector: selectedAssetSectorId,
        assetSubSector: selectedAssetSubSectorId,
        assetType: selectedAssetTypeId,
        offeringType: offeringType === 'private' ? 'private' : 'public',
        approvedAssetCustodianId: 1,
        assetManagerId: 1,
        agreeTransferTitleToCustodian: confirm ? 1 : 0,
        assetAlreadyExists:
          isFinanceAndInvestmentMarketsSector || assetStatus === 'existing'
            ? 1
            : 0,
        secApproval: 0,
        secApprovalIdNumber: '',
        assetCountryLocation: selectedCountry,
        proceedPayoutCurrency,
        assetQuoteCurrency,
        fundingStructure: fundingStructureCode,
        equityPercentage:
          fundingStructure === 'hybrid' ? Number(equityPercentage || 0) : 0,
        debtPercentage:
          fundingStructure === 'hybrid' ? Number(debtPercentage || 0) : 0,
        acceptTokenizationTermsAndAgreement: acceptTokenizationTermsAndAgreement
          ? 1
          : 0,
      };

      const payload = excludeSubmissionFields(payloadBody);

      console.log('this is payload ====>', payload);

      const res = await submitTokenization({
        signer: appUser.primarySigner,
        publicKey: appUser.publicKey,
        secretKey: currentSecretKey,
        body: payload,
      });

      hideLoader();
      console.log('res here ====> ', res);

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

        showNotification('success', 'Setup and compliance saved successfully.');
        navigate(`/dashboard/tokenize/apply/${savedAssetId}/asset-information`);
        return;
      }

      if ('error' in res) {
        const errorData = (res as any).error;
        showNotification(
          'error',
          errorData?.data?.message ??
            'Unable to save setup and compliance. Please try again.',
        );
        return;
      }

      showNotification(
        'error',
        'Unable to save setup and compliance. Please try again.',
      );
    } catch (error: any) {
      hideLoader();
      showNotification(
        'error',
        error?.message ??
          'Unable to save setup and compliance. Please try again.',
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="bg-white p-6 w-full mx-auto">
      {/* Asset Classification */}
      <div className="space-y-6 margin">
        <div className="flex gap-x-6">
          <div className="w-2/6">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Asset Classification
            </label>
            <p className="text-sm text-gray-500 mb-2 max-w-sm ">
              Please select the classification from the list that suits your
              asset the most
            </p>
          </div>
          {/* <div className="grid md:grid-cols-3 gap-4"> */}
          <div className="flex-1 gap-y-6 flex flex-col w-2/6">
            <div
              className=""
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('assetSector', node);
                } else {
                  fieldRefs.current.delete('assetSector');
                }
              }}
            >
              <label className="block font-medium text-[#0E1F51] mb-1">
                Asset Sector
              </label>
              <Dropdown
                label="Asset Sector"
                options={assetSectorOptions}
                defaultValue={selectedSectorOption}
                onSelect={(item: OptionItem) => {
                  setSelectedAssetSectorId(item.value);
                  setSelectedAssetSubSectorId('');
                  setSelectedAssetTypeId('');
                  clearError('assetSector');
                }}
              />
              {errors.assetSector ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.assetSector}
                </p>
              ) : null}
            </div>
            <div className="flex align-center flex-1 gap-x-6">
              <div
                className="flex-1"
                ref={(node) => {
                  if (node) {
                    fieldRefs.current.set('assetSubSector', node);
                  } else {
                    fieldRefs.current.delete('assetSubSector');
                  }
                }}
              >
                <label className="block font-medium text-[#0E1F51] mb-1">
                  Asset Sub-sector
                </label>
                <Dropdown
                  label="Asset Sub-sector"
                  options={assetSubSectorOptions}
                  defaultValue={selectedSubSectorOption}
                  onSelect={(item: OptionItem) => {
                    setSelectedAssetSubSectorId(item.value);
                    setSelectedAssetTypeId('');
                    clearError('assetSubSector');
                  }}
                />
                {errors.assetSubSector ? (
                  <p className="mt-1 text-xs text-red-500">
                    {errors.assetSubSector}
                  </p>
                ) : null}
              </div>
              <div
                className="flex-1"
                ref={(node) => {
                  if (node) {
                    fieldRefs.current.set('assetType', node);
                  } else {
                    fieldRefs.current.delete('assetType');
                  }
                }}
              >
                <label className="block font-medium text-[#0E1F51] mb-1">
                  Asset Type
                </label>
                <Dropdown
                  label="Asset Type"
                  options={assetTypeOptions}
                  defaultValue={selectedTypeOption}
                  onSelect={(item: OptionItem) => {
                    setSelectedAssetTypeId(item.value);
                    clearError('assetType');
                  }}
                />
                {errors.assetType ? (
                  <p className="mt-1 text-xs text-red-500">
                    {errors.assetType}
                  </p>
                ) : null}
              </div>
            </div>
          </div>
        </div>
        <hr />
        <div className="flex gap-x-6">
          <div className="w-2/6">
            <h3 className="block font-medium text-[#0E1F51] mb-1">
              Asset Location
            </h3>
            <p className="text-sm text-gray-500 mb-2">
              Select Country where the asset is located
            </p>
          </div>
          <div
            className="flex-1 w-4/6"
            ref={(node) => {
              if (node) {
                fieldRefs.current.set('country', node);
              } else {
                fieldRefs.current.delete('country');
              }
            }}
          >
            <label className="block font-medium text-[#0E1F51] mb-1">
              Country
            </label>
            <Dropdown
              label="Asset Location"
              options={countryOptions}
              defaultValue={selectedCountryOption}
              onSelect={(item: OptionItem) => {
                setSelectedCountry(item.value);
                clearError('country');
              }}
            />
            {errors.country ? (
              <p className="mt-1 text-xs text-red-500">{errors.country}</p>
            ) : null}
          </div>
        </div>
        <hr />
        {!isFinanceAndInvestmentMarketsSector ? (
          <>
            <div className="flex gap-x-6">
              <div className="w-2/6">
                <label className="block font-medium text-[#0E1F51] mb-1">
                  Asset Status
                </label>
                <p className="text-sm text-gray-500 mb-2">
                  Select what applies to the current status of your asset
                </p>
              </div>
              <div className="flex-1 w-4/6">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="radio"
                    className="accent-primary-700"
                    checked={assetStatus === 'existing'}
                    onChange={() => setAssetStatus('existing')}
                  />
                  Asset is already existing
                </label>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="radio"
                    className="accent-primary-700"
                    checked={assetStatus === 'toBuild'}
                    onChange={() => setAssetStatus('toBuild')}
                  />
                  Asset is yet to be built/acquired
                </label>
              </div>
            </div>
            <hr />
            <div className="flex gap-x-6">
              <div className="w-2/6">
                <label className="block font-medium text-[#0E1F51] mb-1">
                  Funding Structure
                </label>
                <p className="text-sm text-gray-500 mb-2">
                  Select the type of funding structure that applies to your
                  asset
                </p>
              </div>
              <div className="flex-1 w-4/6">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="radio"
                    className="accent-primary-700"
                    checked={fundingStructure === 'equity'}
                    onChange={() => setFundingStructure('equity')}
                  />
                  Equity
                </label>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="radio"
                    className="accent-primary-700"
                    checked={fundingStructure === 'debt'}
                    onChange={() => setFundingStructure('debt')}
                  />
                  Debt
                </label>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="radio"
                    className="accent-primary-700"
                    checked={fundingStructure === 'hybrid'}
                    onChange={() => setFundingStructure('hybrid')}
                  />
                  Hybrid
                </label>

                {fundingStructure === 'hybrid' ? (
                  <div className="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div
                      ref={(node) => {
                        if (node) {
                          fieldRefs.current.set('equityPercentage', node);
                        } else {
                          fieldRefs.current.delete('equityPercentage');
                        }
                      }}
                    >
                      <label className="block font-medium text-[#0E1F51] mb-1">
                        Equity percentage (%)
                      </label>
                      <input
                        type="number"
                        className="mt-1 h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 outline-none focus:border-primary-600"
                        placeholder="How much (%)"
                        value={equityPercentage}
                        onChange={(event) => {
                          setEquityPercentage(event.target.value);
                          clearError('equityPercentage');
                        }}
                      />
                      {errors.equityPercentage ? (
                        <p className="mt-1 text-xs text-red-500">
                          {errors.equityPercentage}
                        </p>
                      ) : null}
                    </div>
                    <div
                      ref={(node) => {
                        if (node) {
                          fieldRefs.current.set('debtPercentage', node);
                        } else {
                          fieldRefs.current.delete('debtPercentage');
                        }
                      }}
                    >
                      <label className="block font-medium text-[#0E1F51] mb-1">
                        Debt percentage (%)
                      </label>
                      <input
                        type="number"
                        className="mt-1 h-12 w-full rounded-md border border-gray-200 bg-primary-100 px-3 text-gray-700 outline-none focus:border-primary-600"
                        placeholder="How much (%)"
                        value={debtPercentage}
                        onChange={(event) => {
                          setDebtPercentage(event.target.value);
                          clearError('debtPercentage');
                        }}
                      />
                      {errors.debtPercentage ? (
                        <p className="mt-1 text-xs text-red-500">
                          {errors.debtPercentage}
                        </p>
                      ) : null}
                    </div>
                  </div>
                ) : null}
              </div>
            </div>
            <hr />
          </>
        ) : null}
        {/* Offering Type */}
        <div className="flex gap-x-6">
          <div className="w-2/6">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Offering Type
            </label>
            <p className="text-sm text-gray-500 max-w-sm mb-2">
              Select the type of offering for your tokenization{' '}
              <span className="italic">
                (private offering is restricted to closed groups)
              </span>
            </p>
          </div>
          <div className="flex-1 w-4/6">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                className="accent-primary-700"
                checked={offeringType === 'public'}
                onChange={() => setOfferingType('public')}
              />
              Public
            </label>
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                className="accent-primary-700"
                checked={offeringType === 'private'}
                onChange={() => setOfferingType('private')}
              />
              Private
            </label>
          </div>
        </div>
        <hr />
        {/* Required Documents */}
        <div className="flex gap-x-6">
          <div className="w-2/6">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Required Documents
            </label>
            <p className="text-sm text-gray-500 mb-2">
              Please select the option after checking the requirements
            </p>
          </div>
          <div className="flex-1 w-4/6">
            <label className="flex items-center cursor-pointer">
              Do you have all the documents listed{' '}
              <span className="mx-1 text-primary-800 underline">
                <a
                  href={getRequirementsUrl() || '#'}
                  target="_blank"
                  rel="noreferrer"
                  onClick={(event) => {
                    if (!selectedAssetSectorId) {
                      event.preventDefault();
                      showNotification(
                        'info',
                        'You must select an asset sector before you can view the requirements.',
                      );
                    }
                  }}
                >
                  here
                </a>
              </span>{' '}
              on the requirements for tokenization
            </label>
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                className="accent-primary-700"
                checked={hasDocuments === 'yes'}
                onChange={() => {
                  setHasDocuments('yes');
                  clearError('documents');
                }}
              />
              Yes
            </label>
            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('documents', node);
                } else {
                  fieldRefs.current.delete('documents');
                }
              }}
            >
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  className="accent-primary-700"
                  checked={hasDocuments === 'no'}
                  onChange={() => {
                    setHasDocuments('no');
                    clearError('documents');
                  }}
                />
                No
              </label>
              {errors.documents ? (
                <p className="mt-1 text-xs text-red-500">{errors.documents}</p>
              ) : null}
            </div>

            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('confirm', node);
                } else {
                  fieldRefs.current.delete('confirm');
                }
              }}
            >
              <label className="flex items-center mt-10 gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  className="accent-primary-700"
                  checked={confirm}
                  onChange={() => {
                    setConfirm(!confirm);
                    clearError('confirm');
                  }}
                />
                I confirm that tokenizing this asset requires transferring it to
                a licensed Asset Custodian
              </label>
              {errors.confirm ? (
                <p className="mt-1 text-xs text-red-500">{errors.confirm}</p>
              ) : null}
            </div>

            <div
              ref={(node) => {
                if (node) {
                  fieldRefs.current.set('acceptTerms', node);
                } else {
                  fieldRefs.current.delete('acceptTerms');
                }
              }}
            >
              <label className="flex items-center mt-10 gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  className="accent-primary-700"
                  checked={acceptTokenizationTermsAndAgreement}
                  onChange={() => {
                    setAcceptTokenizationTermsAndAgreement(
                      !acceptTokenizationTermsAndAgreement,
                    );
                    clearError('acceptTerms');
                  }}
                />
                I agree to the Trovotech Tokenization Terms and Conditions
              </label>
              {errors.acceptTerms ? (
                <p className="mt-1 text-xs text-red-500">
                  {errors.acceptTerms}
                </p>
              ) : null}
            </div>
          </div>
        </div>
      </div>

      {/* Buttons */}
      <div className="w-full flex justify-end mt-10">
        <div className="flex self-end w-3/6 space-x-6">
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
              disabled={isSubmitting}
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
