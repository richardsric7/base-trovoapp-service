import { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import Header from '../../../components/header';
import Modal from '../../../components/modal';
import ButtonSecondary from '../../../components/buttonSecondary';
import { RootState } from '../../../store/reduxStore';
import { TokenizationData } from '../../../types/tokenizationData';
import { countries } from '../../../utils/countries';
import { capitalizeEachWord } from '../../../utils/capitalizeFirst';
import {
  formatHistoryNumber,
  formatNumberShort,
  getFiatValue,
} from '../../../utils/fomatNumber';
import { formatDateLong } from '../../../utils/formatDate';
import { useActiveTokenizedAsset } from './useActiveTokenizedAsset';

type DetailValue = string | Record<string, string>;
type DetailEntries = Record<string, DetailValue>;

function isZeroLikeValue(value: string): boolean {
  const normalized = value.replace(/,/g, '').trim().toLowerCase();

  if (!normalized) {
    return true;
  }

  // Match leading numeric portion in values such as "0", "0.00 usd", "0%".
  const numericMatch = normalized.match(/^[-+]?\d*\.?\d+/);
  if (!numericMatch) {
    return false;
  }

  const numericValue = Number(numericMatch[0]);
  return Number.isFinite(numericValue) && numericValue === 0;
}

function getTokenizationStatus(status: number): string {
  switch (status) {
    case 0:
      return 'Draft';
    case 1:
      return 'Awaiting Fee';
    case 4:
      return 'Approved';
    case 5:
      return 'Primary Sales';
    case 6:
      return 'Secondary Market';
    case 7:
      return 'Liquidated';
    case 8:
      return 'Refunded';
    default:
      return 'Processing';
  }
}

function toNumber(value: unknown): number {
  const numeric = Number(value);
  return Number.isFinite(numeric) ? numeric : 0;
}

function truncateToDecimals(value: number, decimals: number = 2): string {
  return value.toFixed(decimals);
}

function safeDate(value?: string): string {
  if (!value) {
    return 'Nill';
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return 'Nill';
  }

  return formatDateLong(date);
}

function AssetDetailItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-2xl w-full bg-white px-4 py-3">
      <p className="text-sm text-primary-500">{label}</p>
      <p className="text-sm md:text-base font-montserratSemiBold text-primary-800 break-words mt-1">
        {value}
      </p>
    </div>
  );
}

function AsseCategoryItem({
  label,
  image,
  onClick,
}: {
  label: string;
  image: string;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className="flex py-4 w-full items-center px-2 xl:px-5 rounded-xl bg-white"
    >
      <div className="flex space-x-4 items-center w-full">
        <img src={image} alt={label} className="h-10 w-10" />
        <p className="text-base font-montserratSemiBold text-primary-800 text-left">
          {label}
        </p>
      </div>
      <img src="/images/go.png" alt="Open" />
    </button>
  );
}

function DetailsModal({
  title,
  entries,
  isOpen,
  onClose,
}: {
  title: string;
  entries: DetailEntries;
  isOpen: boolean;
  onClose: () => void;
}) {
  const hideZeroLikeValues = title.trim().toLowerCase() === 'financial details';

  return (
    <Modal showModal={isOpen} onClose={onClose}>
      <div className="px-5 pb-6">
        <h3 className="font-montserratSemiBold text-lg text-primary-800 mb-4">
          {title}
        </h3>
        <div className="space-y-3 max-h-[60vh] overflow-y-auto">
          {Object.entries(entries)
            .filter(([, value]) => {
              if (typeof value === 'string') {
                return (
                  value.trim().length > 0 &&
                  (!hideZeroLikeValues || !isZeroLikeValue(value))
                );
              }

              return Object.values(value).some((nestedValue) => {
                return (
                  nestedValue.trim().length > 0 &&
                  (!hideZeroLikeValues || !isZeroLikeValue(nestedValue))
                );
              });
            })
            .map(([key, value]) => (
              <div key={key} className="rounded-lg bg-primary-100 px-4 py-3">
                <p className="text-primary-800 font-montserratSemiBold">
                  {key}
                </p>

                {typeof value === 'string' ? (
                  <p className="text-primary-800 whitespace-pre-wrap break-words mt-1">
                    {value}
                  </p>
                ) : (
                  <div className="mt-2 space-y-1">
                    {Object.entries(value).map(([nestedKey, nestedValue]) => {
                      if (
                        nestedValue.trim().length <= 0 ||
                        (hideZeroLikeValues && isZeroLikeValue(nestedValue))
                      ) {
                        return null;
                      }

                      return (
                        <div
                          key={`${key}-${nestedKey}`}
                          className="flex items-start justify-between gap-3"
                        >
                          <p className="text-md text-primary-800">
                            {nestedKey}
                          </p>
                          <p className="text-md font-montserratSemiBold text-primary-700 text-right break-words">
                            {nestedValue}
                          </p>
                        </div>
                      );
                    })}
                  </div>
                )}
              </div>
            ))}
        </div>
        <div className="mt-5">
          <ButtonSecondary label="Close" onclick={onClose} />
        </div>
      </div>
    </Modal>
  );
}

function DocumentsModal({
  title,
  items,
  isOpen,
  onClose,
}: {
  title: string;
  items: Array<{ label: string; url: string }>;
  isOpen: boolean;
  onClose: () => void;
}) {
  return (
    <Modal showModal={isOpen} onClose={onClose}>
      <div className="px-5 pb-6">
        <h3 className="font-montserratSemiBold text-lg text-primary-800 mb-4">
          {title}
        </h3>
        {items.length === 0 ? (
          <p className="text-sm text-primary-500">No documents here...</p>
        ) : (
          <div className="space-y-2 max-h-[60vh] overflow-y-auto">
            {items.map((item) => (
              <button
                key={`${item.label}-${item.url}`}
                type="button"
                onClick={() =>
                  window.open(item.url, '_blank', 'noopener,noreferrer')
                }
                className="block text-left text-sm underline text-primary-800 break-all"
              >
                {item.label}
              </button>
            ))}
          </div>
        )}
        <div className="mt-5">
          <ButtonSecondary label="Close" onclick={onClose} />
        </div>
      </div>
    </Modal>
  );
}

export function TokenizationAssetDashboard() {
  const navigate = useNavigate();
  const appState = useSelector((state: RootState) => state.appState);
  const { activeTokenizedAsset: asset, isLoading } = useActiveTokenizedAsset();
  const tokenizationData = appState.tokenizationData as
    | TokenizationData
    | undefined;

  const [modalEntries, setModalEntries] = useState<DetailEntries>({});
  const [modalTitle, setModalTitle] = useState('');
  const [showDetailsModal, setShowDetailsModal] = useState(false);
  const [showDocumentsModal, setShowDocumentsModal] = useState(false);
  const [documentsModalTitle, setDocumentsModalTitle] = useState('');
  const [documentItems, setDocumentItems] = useState<
    Array<{
      label: string;
      url: string;
    }>
  >([]);

  const countryNameMap = useMemo(() => {
    return countries.reduce<Record<string, string>>((accumulator, country) => {
      accumulator[country.isoAlpha2.toLowerCase()] = country.name;
      return accumulator;
    }, {});
  }, []);

  const context = useMemo(() => {
    if (!asset || !tokenizationData) {
      return null;
    }

    const assetType =
      tokenizationData.assetTypes.find(
        (entry) => String(entry.id) === String(asset.assetType),
      )?.assetType ?? '';

    const countryConfig = tokenizationData.countryConfigs.find(
      (entry) =>
        entry.countryCode.toLowerCase() ===
        String(asset.assetCountryLocation ?? '').toLowerCase(),
    );

    const fiatCurrency =
      tokenizationData.tokenizationCurrencies
        .find(
          (entry) =>
            entry.assetCode.toLowerCase() ===
            String(countryConfig?.quoteCurrencyCode ?? '').toLowerCase(),
        )
        ?.label?.toUpperCase() ??
      String(asset.assetQuoteCurrency ?? '').toUpperCase();

    const regulatorName = countryConfig?.regulatorName ?? '';
    const tokenizationApplicationFee = toNumber(
      countryConfig?.tokenizationApplicationFee,
    );
    const tokenizationApplicationFeeAsset = String(
      countryConfig?.tokenizationApplicationFeeAsset ?? '',
    )
      .split(':')[0]
      .toUpperCase();

    const amountRaised =
      toNumber(asset.quantityOfTokensSold) * toNumber(asset.pricePerToken);
    const totalAmountToBeRaised =
      toNumber(asset.numberOfTokenToBeSold) * toNumber(asset.pricePerToken);

    const salesStartDate = new Date(asset.salesStart);
    const salesEndDate = new Date(asset.salesEnd);

    const daysProgress = Math.max(
      0,
      Math.floor(
        (Date.now() - salesStartDate.getTime()) / (1000 * 60 * 60 * 24),
      ),
    );
    const totalDays = Math.max(
      1,
      Math.floor(
        (salesEndDate.getTime() - salesStartDate.getTime()) /
          (1000 * 60 * 60 * 24),
      ),
    );

    const normalizedProgress =
      totalAmountToBeRaised > 0
        ? Math.min(1, amountRaised / totalAmountToBeRaised)
        : 0;

    const normalizedDaysProgress =
      daysProgress >= totalDays ? 1 : 1 - daysProgress / totalDays;

    const isFinancialAssetType =
      String(asset.assetSector ?? '').toLowerCase() ===
      'finance and investment markets';

    const exemptedCountries = String(asset.exemptedCountries ?? '')
      .replaceAll(' ', '')
      .split(',')
      .filter((entry) => entry.trim().length > 0)
      .map(
        (entry) => countryNameMap[entry.toLowerCase()] ?? entry.toUpperCase(),
      );

    const tokenizationFee =
      tokenizationData.tokenizationFees.find(
        (fee, index) =>
          index === toNumber(asset.tokenizationFeeId) ||
          fee.id === toNumber(asset.tokenizationFeeId),
      ) ?? tokenizationData.tokenizationFees[0];

    const fiatPercentage = toNumber(tokenizationFee?.feeFiatPercentage);
    const fiatFeeCap = toNumber(tokenizationFee?.feeFiatCap);
    const feeInfo = (toNumber(asset.assetCurrentValue) * fiatPercentage) / 100;
    const selectedTokenizationFee = Math.max(fiatFeeCap, feeInfo);

    const feeInAsset = toNumber(asset.feeInAsset);
    const vatInAsset = feeInAsset * 0.075;
    const totalFeeInAsset = feeInAsset + vatInAsset;

    const equivalentFeesInToken =
      totalFeeInAsset * toNumber(asset.pricePerToken);

    const totalFee =
      toNumber(asset.SECTokenizationFeeValue) +
      toNumber(asset.custodianFeeValue) +
      toNumber(asset.assetManagerFeeValue) +
      toNumber(asset.issuingHouseFeeValue) +
      toNumber(asset.legalAndProfessionalFeeValue) +
      toNumber(asset.ratingAgencyFeeValue) +
      toNumber(asset.trusteeFeeValue) +
      toNumber(asset.vatValue) +
      selectedTokenizationFee;

    const tokensNotForSale =
      toNumber(asset.numberOfTokenToBeIssued) -
      toNumber(asset.numberOfTokenToBeSold) -
      toNumber(asset.feeInAsset) -
      toNumber(asset.vatInAsset);

    return {
      assetType,
      fiatCurrency,
      regulatorName,
      tokenizationApplicationFee,
      tokenizationApplicationFeeAsset,
      amountRaised,
      totalAmountToBeRaised,
      normalizedProgress,
      normalizedDaysProgress,
      daysProgress,
      totalDays,
      isFinancialAssetType,
      exemptedCountries,
      totalFee,
      tokensNotForSale,
      selectedTokenizationFee,
      feeInAsset,
      vatInAsset,
      totalFeeInAsset,
      equivalentFeesInToken,
    };
  }, [asset, tokenizationData, countryNameMap]);

  if (isLoading) {
    return (
      <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
        <Header />
        <div className="rounded-2xl bg-primary-100 p-8 text-center">
          <p className="font-montserratSemiBold text-md text-primary-800">
            Loading tokenized asset details...
          </p>
        </div>
      </div>
    );
  }

  if (!asset) {
    return (
      <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
        <Header />
        <div className="rounded-2xl bg-primary-100 p-8 text-center">
          <img className="mx-auto mb-4" src="/images/empty.png" alt="Empty" />
          <p className="font-montserratSemiBold text-md text-primary-800">
            No tokenized asset selected.
          </p>
          <p className="text-sm mt-2 text-primary-500">
            Go back and select an asset to view details.
          </p>
          <div className="max-w-[240px] mx-auto mt-5">
            <ButtonSecondary label="Go Back" onclick={() => navigate(-1)} />
          </div>
        </div>
      </div>
    );
  }

  if (!context) {
    return (
      <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
        <Header />
        <div className="rounded-2xl bg-primary-100 p-8 text-center">
          <p className="font-montserratSemiBold text-md text-primary-800">
            Tokenization metadata is unavailable.
          </p>
          <p className="text-sm mt-2 text-primary-500">
            Refresh tokenization data from the tokenize page and try again.
          </p>
          <div className="max-w-[240px] mx-auto mt-5">
            <ButtonSecondary label="Go Back" onclick={() => navigate(-1)} />
          </div>
        </div>
      </div>
    );
  }

  const openDetailsModal = (title: string, entries: DetailEntries) => {
    setModalTitle(title);
    setModalEntries(entries);
    setShowDetailsModal(true);
  };

  const openDocumentsModal = (
    title: string,
    items: Array<{ label: string; url: string }>,
  ) => {
    setDocumentsModalTitle(title);
    setDocumentItems(items);
    setShowDocumentsModal(true);
  };

  const salesWindowLabel = `${safeDate(asset.salesStart)} to ${safeDate(asset.salesEnd)}`;
  const tokensRemaining =
    toNumber(asset.numberOfTokenToBeSold) -
    toNumber(asset.quantityOfTokensSold);

  return (
    <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
      <Header isHomeView />

      <div className="w-full px-3 md:px-5 md:pt-5 flex flex-col md:flex-row space-y-5 md:space-y-0 items-center justify-between">
        <div className="flex space-x-5 w-full md:w-auto items-center">
          <button type="button" onClick={() => navigate(-1)}>
            <img src="/images/arrowBack.png" alt="arrow back" />
          </button>
          <p className="font-montserratSemiBold text-lg xl:text-xl">
            Asset Details
          </p>
        </div>
      </div>

      <div className="w-full px-3 md:px-5 md:pt-5 flex flex-col items-center md:space-x-3 space-y-5">
        <img
          src={asset.assetLogo || '/images/trovoLogo.png'}
          className="rounded-full h-[120px] w-[120px] object-cover"
          alt="asset logo"
          onError={(event) => {
            (event.target as HTMLImageElement).src = '/images/trovoLogo.png';
          }}
        />
        <div className="w-full flex-col flex justify-center items-center space-y-3">
          <p className="flex items-center space-x-2 font-montserratSemiBold text-xl">
            {(asset.assetCode || '').toUpperCase()}
          </p>
          <div className="flex flex-wrap items-center justify-center gap-2">
            <p className="bg-primary-100 py-2 px-4 rounded-full font-montserratSemiBold text-sm">
              {asset.assetAlreadyExists === 1
                ? 'Existing Asset'
                : 'Upcoming Asset'}
            </p>
            <div className="w-full md:w-auto">
              <div className="rounded-full bg-primary-100 px-4 py-2 text-center">
                <span className="text-sm font-montserratSemiBold text-green-600">
                  {getTokenizationStatus(asset.assetTokenizationStatus)}
                </span>
              </div>
            </div>
          </div>
          <div></div>
          <div></div>
          <p className="font-montserratSemiBold text-lg xl:text-xl">
            Description
          </p>
          <p className="text-center w-full md:w-4/5">
            {asset.assetDescription}
          </p>
        </div>
        <div></div>
        <div></div>
        <div className="w-full rounded-lg bg-primary-100 p-4">
          <p className="font-montserratSemiBold text-primary-800">
            Amount Raised
          </p>
          <div className="w-full mt-2 h-2 rounded-full bg-white overflow-hidden">
            <div
              className="h-full bg-green-500"
              style={{
                width: `${Math.min(100, context.normalizedProgress * 100)}%`,
              }}
            />
          </div>
          <p className="text-xs mt-2 text-primary-600">
            {formatHistoryNumber(context.amountRaised, 99000000000)}{' '}
            {asset.assetQuoteCurrency} raised out of{' '}
            {formatHistoryNumber(context.totalAmountToBeRaised, 99000000000)}{' '}
            {asset.assetQuoteCurrency}
          </p>

          {asset.assetTokenizationStatus === 5 && context.daysProgress >= 0 && (
            <>
              <p className="font-montserratSemiBold text-primary-800 mt-4">
                Sales Window Progress
              </p>
              <div className="w-full mt-2 h-2 rounded-full bg-white overflow-hidden">
                <div
                  className="h-full bg-primary-800"
                  style={{
                    width: `${Math.min(100, context.normalizedDaysProgress * 100)}%`,
                  }}
                />
              </div>
              <p className="text-xs mt-2 text-primary-600">
                {Math.max(0, context.totalDays - context.daysProgress)} days
                left out of {context.totalDays} days
              </p>
            </>
          )}

          {asset.assetTokenizationStatus === 5 && (
            <p className="text-xs mt-3 text-primary-600">
              Tokens remaining:{' '}
              {formatNumberShort(Math.max(0, tokensRemaining))}{' '}
              {asset.assetCode}
            </p>
          )}
        </div>
      </div>

      <div className="flex flex-col md:flex-row space-y-5 md:space-y-0 md:space-x-8 px-3 md:px-5">
        <div className="flex w-full md:w-4/6 space-y-3 py-7 px-2 xl:px-5 rounded-lg flex-col items-center bg-primary-100">
          <div className="flex items-center w-full justify-between">
            <p className="font-montserratSemiBold text-lg md:text-xl">
              Statistics
            </p>
          </div>
          <div className="w-full grid grid-cols-1 md:grid-cols-2 gap-4">
            <AssetDetailItem
              label={
                asset.assetAlreadyExists === 1
                  ? 'Asset Value'
                  : 'Total Project Budget'
              }
              value={`${getFiatValue(toNumber(asset.numberOfTokenToBeIssued) * toNumber(asset.pricePerToken))} ${context.fiatCurrency}`}
            />
            <AssetDetailItem
              label="Total Token Supply"
              value={`${getFiatValue(toNumber(asset.numberOfTokenToBeIssued))} ${asset.assetCode}`}
            />
            <AssetDetailItem
              label="Tokens not for Sale"
              value={`${getFiatValue(context.tokensNotForSale)} ${asset.assetCode}`}
            />
            <AssetDetailItem
              label={
                asset.assetAlreadyExists === 1
                  ? 'Amount Retained'
                  : 'Equity Contributed'
              }
              value={`${getFiatValue(toNumber(asset.assetOwnerRetainedOrContributedValue))} ${context.fiatCurrency}`}
            />
            <AssetDetailItem
              label="Tokens for Sale"
              value={`${getFiatValue(toNumber(asset.numberOfTokenToBeSold))} ${asset.assetCode}`}
            />
            <AssetDetailItem
              label="Amount to be Raised"
              value={`${getFiatValue(context.totalAmountToBeRaised)} ${context.fiatCurrency}`}
            />
            <AssetDetailItem label="Funding Currency" value="CNGN" />
            <AssetDetailItem
              label="Price per Token"
              value={`${truncateToDecimals(toNumber(asset.pricePerToken), 6)} ${context.fiatCurrency}`}
            />
            <AssetDetailItem
              label="No. of Interest Expressed"
              value={`${toNumber(asset.numberOfExpressedInterests)} users`}
            />
            <AssetDetailItem
              label="Purchase Commitments"
              value={`${getFiatValue(toNumber(asset.purchaseCommitments))} ${asset.assetQuoteCurrency}`}
            />
            <AssetDetailItem
              label="Total Quantity Sold"
              value={`${truncateToDecimals(toNumber(asset.quantityOfTokensSold), 6)} ${asset.assetCode}`}
            />
            <AssetDetailItem
              label="Total Amount Raised"
              value={`${getFiatValue(context.amountRaised)} ${asset.assetQuoteCurrency}`}
            />
            <AssetDetailItem
              label="No. of Users Bought"
              value={`${toNumber(asset.numberOfSubscribers)} users`}
            />
            <AssetDetailItem
              label="% of Amount Raised"
              value={`${formatNumberShort(
                context.totalAmountToBeRaised > 0
                  ? (context.amountRaised / context.totalAmountToBeRaised) * 100
                  : 0,
              )} %`}
            />
          </div>
        </div>

        <div className="flex w-full md:w-2/6 space-y-3 py-7 px-2 xl:px-5 rounded-lg flex-col items-center bg-primary-100">
          <div className="flex items-center w-full justify-between">
            <p className="font-montserratSemiBold text-lg md:text-xl">
              Details
            </p>
          </div>
          <AsseCategoryItem
            label="Asset Information"
            image="/images/asset_information.png"
            onClick={() => {
              const entries: DetailEntries = {
                Status:
                  asset.assetAlreadyExists === 1 ? 'Existing' : 'Upcoming',
                Sector: asset.assetSector ?? '',
                'Sub-Sector': asset.assetSubSector ?? '',
                Type: context.assetType,
                'Asset Original Value': `${getFiatValue(toNumber(asset.assetCurrentValue))} ${context.fiatCurrency}`,
                'Incurred Costs Outside Valuation': `${getFiatValue(toNumber(asset.assetMscCostOutisdeOfValuation))} ${context.fiatCurrency}`,
                'Asset Country':
                  countryNameMap[
                    String(asset.assetCountryLocation ?? '').toLowerCase()
                  ] ?? String(asset.assetCountryLocation ?? ''),
                Address: asset.assetPhysicalAddress ?? '',
                'Project Strategic Objectives':
                  asset.projectStrategicObjectives ?? '',
                'Project Development Timeline':
                  asset.projectDevelopmentTimeline ?? '',
                'Key Milestones & Dates':
                  asset.projectKeyMilestoneAndDates ?? '',
                'Project Scope': asset.projectScope ?? '',
                'Project Economic Benefits':
                  asset.projectEconomicBenefits ?? '',
                'Expected No. of Job to be Created':
                  toNumber(asset.projectExpectedNoOfJobs) > 0
                    ? String(asset.projectExpectedNoOfJobs)
                    : '',
                'Project Intended Social Benefits':
                  asset.projectIntendedSocialBenefits ?? '',
                'Technical Partners': asset.projectTechnicalPartners ?? '',
                'Financial Partners': asset.projectFinancialPartners ?? '',
                'Original Asset Owner / Project Sponsor':
                  asset.assetOwnerName ?? '',
                'Date Submitted': safeDate(asset.dateSubmitted),
                'Approval Date': [4, 5, 6, 7, 8].includes(
                  asset.assetTokenizationStatus,
                )
                  ? safeDate(asset.dateOfApproval)
                  : 'Nill',
              };
              openDetailsModal('Asset Information', entries);
            }}
          />

          <AsseCategoryItem
            label="Asset Token & Sale Information"
            image="/images/asset_token_information.png"
            onClick={() => {
              const entries: DetailEntries = {
                'Minting Date': [5, 6].includes(asset.assetTokenizationStatus)
                  ? safeDate(asset.mintingDate)
                  : 'Nill',
                'Sales Window': salesWindowLabel,
                'Cap Amount': `${getFiatValue(toNumber(asset.capAmountInFiat))} ${context.fiatCurrency}`,
                'Cap Quantity': `${getFiatValue(toNumber(asset.capQuantity))} ${(asset.assetCode || '').toUpperCase()}`,
                'Cap Duration': `${toNumber(asset.capDurationInDays)} days`,
                'Proceed Payout Cycle': capitalizeEachWord(
                  asset.proceedCycle ?? '',
                ),
                'Exempted Countries': context.exemptedCountries.join(', '),
              };
              openDetailsModal('Asset Token & Sale Information', entries);
            }}
          />

          {asset.assetAlreadyExists === 0 && (
            <>
              <AsseCategoryItem
                label="Asset Financial Information"
                image="/images/token.png"
                onClick={() => {
                  const entries: DetailEntries = {
                    'Estimated Project IRR': `${formatNumberShort(toNumber(asset.estimatedProjectIRR))} %`,
                    'Estimated Project ROI': `${formatNumberShort(toNumber(asset.estimatedProjectROI))} %`,
                    'Estimated Project NPV at Launch (Day 1)': `${getFiatValue(toNumber(asset.estimatedProjectNPV))} ${context.fiatCurrency}`,
                    'Estimated Project Payback Periods (in Months)': String(
                      toNumber(asset.estimatedProjectPaybackPeriodsInMonths),
                    ),
                    'All Key Assumptions Including Values Assumed':
                      asset.keyAssumptionsList ?? '',
                  };
                  openDetailsModal('Asset Financial Information', entries);
                }}
              />

              <AsseCategoryItem
                label="Project Risk Assessment"
                image="/images/alert.png"
                onClick={() => {
                  const entries: DetailEntries = {
                    'Legal Risks Identified':
                      asset.projectIdentifiedLegalRisks ?? '',
                    'Regulatory Risks Identified':
                      asset.projectIdentifiedRegulatoryRisks ?? '',
                    'Operational/Execution Risks Identified':
                      asset.projectIdentifiedOperationalOrExecutionRisks ?? '',
                    'Market Risks Identified':
                      asset.projectIdentifiedMarketRisks ?? '',
                    'Other Relevant Risks Identified':
                      asset.projectIdentifiedOtherRelevantRisks ?? '',
                  };
                  openDetailsModal('Project Risk Assessment', entries);
                }}
              />
            </>
          )}

          <AsseCategoryItem
            label="Stakeholders Information"
            image="/images/stakeholders.png"
            onClick={() => {
              const entries: DetailEntries = {
                Regulator: capitalizeEachWord(
                  context.regulatorName.toLowerCase(),
                ),
                'Asset Custodian':
                  asset.approvedAssetCustodianInfo?.assetCustodianName ?? '',
                'Asset Manager': asset.assetManagerInfo?.assetManagerName ?? '',
                'Issuing House':
                  capitalizeEachWord(
                    asset.assetIssuingHouseInfo?.assetIssuingHouseName ?? '',
                  ) ?? '',
                'Legal/Professional Advisor':
                  asset.legalAndProfesionalPartnerInfo?.partnerName ?? '',
                'Rating Agency': asset.ratingAgencyInfo?.agencyName ?? '',
                Trustee: asset.trusteeInfo?.trusteeName ?? '',
              };
              openDetailsModal('Stakeholders Information', entries);
            }}
          />

          {!context.isFinancialAssetType && (
            <AsseCategoryItem
              label="Legal & Compliance Information"
              image="/images/verification-documents.png"
              onClick={() => {
                const entries: DetailEntries = {
                  'I confirm that this asset is free of liens, mortgages, and outstanding loans.':
                    asset.IsFreeFromLiensAndEncumbrances == 1 ? 'Yes' : 'No',
                  'I confirm that this asset is not pledged as collateral for any debts and has no use restrictions.':
                    asset.undertakingNotCollateral == 1 ? 'Yes' : 'No',
                  'I confirm that no third party has any claims, rights, or interests in this asset.':
                    asset.undertakingNoClaims == 1 ? 'Yes' : 'No',
                  'I confirm that this asset is free of any foreclosure, bankruptcy proceedings, legal disputes, judgments, or court-ordered payments.':
                    asset.undertakingNoForeclosure == 1 ? 'Yes' : 'No',
                  'I confirm that this asset complies with all environmental and land use regulations and is free of violations.':
                    asset.complianceNoViolation == 1 ? 'Yes' : 'No',
                  'I confirm that all necessary permits, licenses, and approvals for the use and ownership of this asset are in place.':
                    asset.complianceNoViolation == 1 ? 'Yes' : 'No',
                  'I confirm that there are no unpaid taxes, utility bills, fees, or other property-related expenses associated with this asset.':
                    asset.outstandingFinancialRespNoDebts == 1 ? 'Yes' : 'No',
                  'I confirm that this asset does not have any hidden liabilities or obligations that have not been disclosed. ':
                    asset.outstandingFinancialRespNoHiddenLiabilities == 1
                      ? 'Yes'
                      : 'No',
                  'I confirm that this asset is adequately insured against risks such as fire, theft, and natural disasters.':
                    asset.riskManagementFullyInsured == 1 ? 'Yes' : 'No',
                  'I confirm that the declared value of this asset reflects its current market value and condition.':
                    asset.riskManagementDeclaredValue == 1 ? 'Yes' : 'No',
                  'I confirm that this asset is not affected by undisclosed easements, rights of way, expropriation, or condemnation.':
                    asset.physicalConditionNoUndisclosedEasements == 1
                      ? 'Yes'
                      : 'No',
                  'I confirm that this asset is structurally sound and has no unresolved maintenance or safety issues.':
                    asset.physicalConditionSound == 1 ? 'Yes' : 'No',
                  'I confirm that this asset is not subject to any agreements, such as leases or contracts, that could limit its use or transfer.':
                    asset.physicalConditionNolease == 1 ? 'Yes' : 'No',
                };
                openDetailsModal('Legal & Compliance Information', entries);
              }}
            />
          )}

          <AsseCategoryItem
            label="Asset Documents"
            image="/images/verification-documents.png"
            onClick={() => {
              const items = (asset.AssetTokenizationDocuments ?? []).map(
                (document) => ({
                  label: document.documentTitle || 'Untitled document',
                  url: document.documentUrl || '',
                }),
              );
              openDocumentsModal(
                'Asset Documents',
                items.filter((entry) => entry.url.trim().length > 0),
              );
            }}
          />

          <AsseCategoryItem
            label="Proof of Payment Documents"
            image="/images/verification-documents.png"
            onClick={() => {
              const items = (asset.ProofOfPaymentDocuments ?? [])
                .map((document) => {
                  const entry = document as {
                    documentTitle?: string;
                    documentUrl?: string;
                  };
                  return {
                    label:
                      entry.documentTitle ||
                      entry.documentUrl ||
                      'Proof document',
                    url: entry.documentUrl || '',
                  };
                })
                .filter((entry) => entry.url.trim().length > 0);
              openDocumentsModal('Proof of Payment Documents', items);
            }}
          />

          {asset.vettingStatus === 1 && (
            <AsseCategoryItem
              label="Fees Details"
              image="/images/token.png"
              onClick={() => {
                const entries: DetailEntries = {
                  'Application Fee': `${formatNumberShort(
                    context.tokenizationApplicationFee,
                  )} ${context.tokenizationApplicationFeeAsset}`,
                  'Fees in Fiat': {
                    'Additional Cost': `${getFiatValue(toNumber(asset.assetMscCostOutisdeOfValuation))} ${context.fiatCurrency}`,
                    'Tokenization Fee': `${formatNumberShort(context.selectedTokenizationFee)} ${context.fiatCurrency}`,
                    'SEC Fee': `${formatNumberShort(toNumber(asset.SECTokenizationFeeValue))} ${context.fiatCurrency}`,
                    'Custody Fee': `${formatNumberShort(toNumber(asset.custodianFeeValue))} ${context.fiatCurrency}`,
                    'Management Fee': `${formatNumberShort(toNumber(asset.assetManagerFeeValue))} ${context.fiatCurrency}`,
                    'Issuing House Fee': `${formatNumberShort(toNumber(asset.issuingHouseFeeValue))} ${context.fiatCurrency}`,
                    'Legal/Professional Fee': `${formatNumberShort(toNumber(asset.legalAndProfessionalFeeValue))} ${context.fiatCurrency}`,
                    'Rating Agency Fee': `${formatNumberShort(toNumber(asset.ratingAgencyFeeValue))} ${context.fiatCurrency}`,
                    'Trustee Fee': `${formatNumberShort(toNumber(asset.trusteeFeeValue))} ${context.fiatCurrency}`,
                    VAT: `${formatNumberShort(toNumber(asset.vatValue))} ${context.fiatCurrency}`,
                    Total: `${formatNumberShort(context.totalFee)} ${context.fiatCurrency}`,
                  },
                  'Fees in Asset': {
                    'Tokenization Fee': `${formatNumberShort(context.feeInAsset)} ${(asset.assetCode || '').toUpperCase()}`,
                    VAT: `${formatNumberShort(context.vatInAsset)} ${(asset.assetCode || '').toUpperCase()}`,
                    Total: `${formatNumberShort(context.totalFeeInAsset)} ${(asset.assetCode || '').toUpperCase()}`,
                  },
                  'Current Asset Value': {
                    'Original valuation': `${formatNumberShort(toNumber(asset.assetCurrentValue))} ${context.fiatCurrency}`,
                    'Additional cost': `${formatNumberShort(toNumber(asset.assetMscCostOutisdeOfValuation))} ${context.fiatCurrency}`,
                    'Fees in fiat': `${formatNumberShort(context.totalFee)} ${context.fiatCurrency}`,
                    'Equiv. fees in token': `${formatNumberShort(context.equivalentFeesInToken)} ${context.fiatCurrency}`,
                    Total: `${getFiatValue(toNumber(asset.numberOfTokenToBeIssued) * toNumber(asset.pricePerToken))} ${context.fiatCurrency}`,
                  },
                };
                openDetailsModal('Financial Details', entries);
              }}
            />
          )}
        </div>
      </div>

      <DetailsModal
        title={modalTitle}
        entries={modalEntries}
        isOpen={showDetailsModal}
        onClose={() => setShowDetailsModal(false)}
      />

      <DocumentsModal
        title={documentsModalTitle}
        items={documentItems}
        isOpen={showDocumentsModal}
        onClose={() => setShowDocumentsModal(false)}
      />
    </div>
  );
}

export default TokenizationAssetDashboard;
