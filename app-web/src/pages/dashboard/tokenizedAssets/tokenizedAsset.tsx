import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import Button from '../../../components/button';
import Modal from '../../../components/modal';
import Header from '../../../components/header';
import { useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import Countdown from '../../../components/countdown';
import { formatDateLong } from '../../../utils/formatDate';
import { formatNumber, getFiatValue } from '../../../utils/fomatNumber';
import { capitalizeEachWord } from '../../../utils/capitalizeFirst';
import { Encryptor } from '../../../utils/encryptor';
import {
  useFetchExpressedInterestsQuery,
  useFetchTokenizationDataQuery,
  useFetchTokenizedAssetsQuery,
} from '../../../store/api/walletApis';
import { TokenizationData } from '../../../types/tokenizationData';
import { ExpressedInterest } from '../../../types/expressedInterest';
import ButtonSecondary from '../../../components/buttonSecondary';
import AssetProgressCard from '../../../components/progressbar';
import {
  BuyTokenModal,
  ExpressInterestModal,
  KycUnverifiedModal,
  P2pComingSoonModal,
  SoldOutModal,
} from '../../../components/tokenizedAssetActionModals';

type GridItemProps = {
  title: string;
  icon?: React.ReactNode;
  onClick: () => void;
};

type AssetPropItems = {
  title: string;
  props: { name: string; value: string }[];
};

function AsseCategoryItem({ icon, title, onClick }: GridItemProps) {
  return (
    <button
      onClick={onClick}
      className="flex py-4 w-full items-center px-2 xl:px-5 rounded-xl bg-white"
    >
      <div className="flex b space-x-4 items-center w-full">
        {icon}
        <p className="text-lg font-montserratSemiBold">{title}</p>
      </div>
      <img src="/images/go.png" alt="go" />
    </button>
  );
}

type AssetDetailItemProps = {
  title: string;
  value: string;
  value2?: string;
  value3?: string;
};

// eslint-disable-next-line
function AssetDetailItem({
  title,
  value,
  value2,
  value3,
}: AssetDetailItemProps) {
  return (
    <div className="rounded-2xl w-full bg-white px-4 py-3">
      <div className="flex  flex-col space-y-2 w-full">
        <p className="font-montserratSemiBold md:text-lg">{title}</p>
        <p className="text-sm">{value}</p>
        <p className="text-sm">{value2}</p>
        <p className="text-sm">{value3}</p>
      </div>
    </div>
  );
}

export default function TokenizedAssetDetailsView() {
  const [searchParams] = useSearchParams();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const primaryWallet = appUser.userWallets.find((w) => w.primaryWallet)!;
  const [asset, setAsset] = useState(
    useSelector((state: RootState) => state.appState.activeTokenizedAsset),
  );
  var tokenizationData = useSelector(
    (state: RootState) => state.appState.tokenizationData,
  );
  const [secretKey, setSecretKey] = useState('');

  useEffect(() => {
    (async () => {
      const encryptor = new Encryptor();
      setSecretKey(await encryptor.getSecretKey(appUser));
    })();
  }, []);

  const { data } = useFetchTokenizedAssetsQuery(
    {
      signer: primaryWallet?.signer ?? '',
      address: primaryWallet.address,
      secretKey,
      body: {
        status: searchParams.get('salesList'),
        filters: `assetCode=${searchParams.get('assetCode')}`,
      },
    },
    { skip: !secretKey || !primaryWallet || asset != undefined },
  );

  const { data: expressedInterests } = useFetchExpressedInterestsQuery(
    {
      signer: primaryWallet?.signer ?? '',
      address: primaryWallet.address,
      secretKey,
      body: { status: 1 },
    },
    { skip: !secretKey || !primaryWallet || asset != undefined },
  );

  useEffect(() => {
    if (!asset && data?.records && expressedInterests?.records) {
      const record = expressedInterests.records.find(
        (a: ExpressedInterest) => a.assetCode === data.records[0].assetCode,
      );

      const y = {
        ...data.records[0],
        expressedInterestAmount: record ? record.amount : 0,
      };

      setAsset(y);
    }
  }, [asset, data, expressedInterests]);

  if (!tokenizationData) {
    const { data } = useFetchTokenizationDataQuery(
      {
        signer: primaryWallet?.signer ?? '',
        address: primaryWallet.address,
        secretKey,
        body: { limit: 5 },
      },
      { skip: !secretKey || !primaryWallet },
    );

    tokenizationData = data as TokenizationData;
  }

  const navigate = useNavigate();
  const [showSubscribeModal, setShowSubscribeModal] = useState(false);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [showSoldOutModal, setShowSoldOutModal] = useState(false);
  const [showKycUnverifiedModal, setShowKycUnverifiedModal] = useState(false);
  const [showBuyTokenModal, setShowBuyTokenModal] = useState(false);
  const [showAssetDetailsModal, setShowAssetDetailsModal] = useState(false);
  const [showP2pComingSoonModal, setShowP2pComingSoonModal] = useState(false);
  const [showVerificationDocumentsModal, setShowVerificationDocumentsModal] =
    useState(false);
  const [assetPropItems, setAssetPropItems] = useState<AssetPropItems>();
  const [isFinancialAssetType] = useState(
    asset?.assetSector!.toLowerCase() == 'finance and investment markets',
  );
  var quoteCurrencyCode = '';
  var regulatorName = '';
  var fiatCurrency = '';
  var assetType = '';
  var assetBalance = 0;

  if (tokenizationData) {
    for (var i = 0; i < tokenizationData!.countryConfigs.length; i++) {
      if (
        tokenizationData!.countryConfigs[i]['countryCode']
          .toString()
          .toLowerCase() == asset?.assetCountryLocation.toString().toLowerCase()
      ) {
        quoteCurrencyCode =
          tokenizationData!.countryConfigs[i]['quoteCurrencyCode'];
        regulatorName = tokenizationData!.countryConfigs[i]['regulatorName'];
      }
    }

    for (var i = 0; i < tokenizationData!.tokenizationCurrencies.length; i++) {
      if (
        tokenizationData!.tokenizationCurrencies[i]['assetCode']
          .toString()
          .toLowerCase() == quoteCurrencyCode.toString().toLowerCase()
      ) {
        fiatCurrency = tokenizationData!.tokenizationCurrencies[i]['label']
          .toString()
          .toUpperCase();
      }
    }

    for (var i = 0; i < tokenizationData!.assetTypes.length; i++) {
      if (
        tokenizationData!.assetTypes[i]['id'].toString() == asset?.assetType
      ) {
        assetType = tokenizationData!.assetTypes[i]['assetType'];
      }
    }

    appUser.userWallets.forEach((wallet) => {
      wallet.claimedAssets?.forEach((a) => {
        if (
          asset?.assetCode != null &&
          a.assetCode?.toLowerCase() == asset.assetCode?.toLowerCase()
        ) {
          assetBalance += a.amount;
        }
      });
    });
  }
  const daysProgress = Math.floor(
    (new Date().getTime() -
      new Date(asset?.salesStart ?? new Date()).getTime()) /
      (1000 * 60 * 60 * 24),
  );

  const totalDays = Math.floor(
    (new Date(asset?.salesEnd ?? new Date()).getTime() -
      new Date(asset?.salesStart ?? new Date()).getTime()) /
      (1000 * 60 * 60 * 24),
  );

  const tokensRemaining =
    asset?.numberOfTokenToBeSold! - asset?.quantityOfTokensSold!;
  const tokensBought = asset?.quantityOfTokensSold ?? 0;
  const normalizedProgress = 1 - tokensBought / asset?.numberOfTokenToBeSold!; // Convert to 0-1 range
  const normalizedDaysProgress =
    daysProgress == totalDays ? 1 : 1 - daysProgress / totalDays; // Convert to 0-1 range

  function renderActionButton() {
    switch (asset?.assetTokenizationStatus) {
      case 5:
        return (
          <div className="flex w-full md:w-1/6 space-x-5">
            {tokensRemaining > 0 ? (
              <Button
                label="Buy"
                leftIcon={
                  <img
                    src="/images/buy_with_fiat.png"
                    height={24}
                    width={24}
                    alt="Buy token"
                  />
                }
                additionalClasses="bg-primary-600 text-white font-montserratSemiBold"
                onclick={() => {
                  if (!appUser.kycVerified) {
                    setShowKycUnverifiedModal(true);
                    return;
                  }

                  setShowBuyTokenModal(true);
                }}
              />
            ) : (
              <ButtonSecondary
                label="Sold Out"
                additionalClasses="bg-primary-600 text-white font-montserratSemiBold"
                onclick={() => {
                  setShowSoldOutModal(true);
                }}
              />
            )}
          </div>
        );
      case 4:
        return (
          <div className="flex w-full md:w-1/6 space-x-5">
            <Button
              label={
                asset.expressedInterestAmount > 0
                  ? 'Update Expressed Interest'
                  : 'Express Interest'
              }
              leftIcon={
                <img
                  src="/images/buy_with_fiat.png"
                  height={24}
                  width={24}
                  alt="Express Interest"
                />
              }
              additionalClasses="bg-primary-600 text-white font-montserratSemiBold"
              onclick={() => {
                setShowSubscribeModal(true);
              }}
            />
          </div>
        );
      default:
        return (
          <div className="flex w-full md:w-1/6 space-x-5">
            <Button
              label={'Buy on TrovoP2P'}
              leftIcon={
                <img
                  src="/images/buy_with_fiat.png"
                  height={24}
                  width={24}
                  alt="Buy ETH on TrovoP2P"
                />
              }
              additionalClasses="bg-primary-600 text-white font-montserratSemiBold"
              onclick={() => {
                setShowP2pComingSoonModal(true);
              }}
            />
          </div>
        );
    }
  }

  return !asset ? (
    <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
      <Header isHomeView />
      <div className="w-full px-3 md:px-5 md:pt-5 flex flex-col md:flex-row space-y-5 md:space-y-0 items-center justify-between">
        <button
          className="flex space-x-5 w-full md:w-auto items-center"
          type="button"
          onClick={() => navigate(-1)}
        >
          <img src="/images/arrowBack.png" alt="arrow back" />
          <p className="font-montserratSemiBold text-lg xl:text-xl">Go Back</p>
        </button>
      </div>
      <div className="w-full">
        <div className="md:px-5 w-full">
          <div className="w-full flex flex-col justify-center space-y-5 md:mb-10">
            <img className="self-center" src="/images/empty.png" alt="" />
            <p className="font-montserratSemiBold text-md text-center">
              Something went wrong.
            </p>
          </div>
        </div>
      </div>
    </div>
  ) : (
    <div className="flex text-primary-800 text-sm md:text-md flex-col space-y-5 p-3">
      <Header isHomeView />
      <div className="w-full px-3 md:px-5 md:pt-5 flex flex-col md:flex-row space-y-5 md:space-y-0 items-center justify-between">
        <div className="flex space-x-5 w-full md:w-auto items-center">
          <button type="button" onClick={() => navigate(-1)}>
            <img src="/images/arrowBack.png" alt="arrow back" />
          </button>
          <p className="font-montserratSemiBold text-lg xl:text-xl">
            {asset?.assetCode}
          </p>
        </div>
        {renderActionButton()}
      </div>
      <div className="w-full px-3 md:px-5 md:pt-5 flex flex-col items-center md:space-x-3 space-y-5">
        <img
          src={asset?.assetLogo}
          className="rounded-full"
          height="120"
          width="120"
          alt="asset logo"
        />
        <div className="w-full flex-col flex justify-center items-center space-y-4">
          <p className="flex items-center space-x-2 font-montserratSemiBold text-xl">
            {asset?.assetName}
          </p>
          {asset?.assetTokenizationStatus === 4 ? (
            <Countdown isColumn={true} startDate={new Date(asset.salesStart)} />
          ) : asset?.assetTokenizationStatus === 5 ? (
            <p className="bg-primary-100 p-3 rounded-md font-montserratSemiBold">
              Available
            </p>
          ) : null}
          {asset.expressedInterestAmount > 0 && (
            <p>
              You have indicated interest to invest{' '}
              <b className="font-montserratSemiBold">
                {formatNumber(asset.expressedInterestAmount ?? 0)}{' '}
                {asset.assetQuoteCurrency}
              </b>{' '}
              on this asset
            </p>
          )}
          <p className="text-center w-4/5">{asset?.assetDescription}</p>
        </div>

        {asset.assetTokenizationStatus == 5 && (
          <div className="w-full">
            <AssetProgressCard
              normalizedProgress={normalizedProgress}
              normalizedDaysProgress={normalizedDaysProgress}
              tokensRemaining={tokensRemaining}
              tokenizedAsset={asset}
              totalDays={totalDays}
              daysProgress={daysProgress}
            />
          </div>
        )}
      </div>
      <div className="flex space-x-8 px-3 md:px-5">
        <div className="flex w-4/6 space-y-3 py-7 px-2 xl:px-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <div className="flex items-center w-full justify-between">
            <p className="font-montserratSemiBold text-lg md:text-xl">
              Statistics
            </p>
          </div>
          <div className="w-full grid grid-cols-2 gap-4">
            <AssetDetailItem
              title={
                asset?.assetAlreadyExists == 1
                  ? 'Asset Value'
                  : 'Total Project Budget'
              }
              value={`${getFiatValue(asset.numberOfTokenToBeIssued * asset.pricePerToken)} ${fiatCurrency}`}
              value2=""
            />
            <AssetDetailItem
              title="Total Token Supply"
              value={`${getFiatValue(asset.numberOfTokenToBeIssued!)} ${asset.assetCode.toUpperCase()}`}
            />

            <AssetDetailItem
              title="Tokens for Sale"
              value={`${getFiatValue(asset.numberOfTokenToBeSold ?? 0)} ${asset.assetCode.toUpperCase()}`}
            />
            <AssetDetailItem
              title="Amount to be Raised"
              value={`${getFiatValue(asset.numberOfTokenToBeSold * asset.pricePerToken)} ${fiatCurrency}`}
            />

            <AssetDetailItem title="Funding Currency" value="CNGN" value2="" />
            <AssetDetailItem
              title="Price per Token"
              value={`${formatNumber(asset.pricePerToken!)} ${fiatCurrency}`}
            />
            <AssetDetailItem
              title="Total Tokens Held"
              value={`${getFiatValue(assetBalance)} ${fiatCurrency}`}
              value2=""
            />
            <AssetDetailItem
              title="Value of Tokens Held"
              value={`${getFiatValue(assetBalance * asset.pricePerToken!)} ${fiatCurrency}`}
            />
          </div>
        </div>
        <div className="flex w-2/6 space-y-3 py-7 px-2 xl:px-5 xl:space-y-5 rounded-lg flex-col items-center bg-primary-100">
          <div className="flex items-center w-full justify-between">
            <p className="font-montserratSemiBold text-lg md:text-xl">
              Details
            </p>
          </div>
          <AsseCategoryItem
            onClick={() => {
              setShowAssetDetailsModal(true);
              setAssetPropItems({
                title: 'Asset Information',
                props: [
                  { name: 'Sector', value: asset?.assetSector ?? '' },
                  {
                    name: 'Sub-Sector',
                    value: asset?.assetSubSector ?? '',
                  },
                  { name: 'Type', value: assetType ?? '' },
                  {
                    name: 'Asset Country',
                    value: asset?.assetCountryLocation ?? '',
                  },
                  {
                    name: 'Address',
                    value: isFinancialAssetType
                      ? ''
                      : (asset?.assetPhysicalAddress ?? ''),
                  },
                  {
                    name: 'Map Coordinates',
                    value: isFinancialAssetType
                      ? ''
                      : !asset?.assetLatitude ||
                          asset?.assetLatitude.length === 0
                        ? ''
                        : `Lat. ${asset?.assetLatitude}, Lon. ${asset?.assetLongitude}`,
                  },
                  {
                    name: 'Project Strategic Objectives',
                    value: asset?.projectStrategicObjectives ?? '',
                  },
                  {
                    name: 'Project Development Timeline',
                    value: asset?.projectDevelopmentTimeline ?? '',
                  },
                  {
                    name: 'Key Milestones & Dates',
                    value: asset?.projectKeyMilestoneAndDates ?? '',
                  },
                  {
                    name: 'Project Scope',
                    value: asset?.projectScope ?? '',
                  },
                  {
                    name: 'Project Economic Benefits',
                    value: asset?.projectEconomicBenefits ?? '',
                  },
                  {
                    name: 'Expected No. of Job to be Created',
                    value:
                      asset?.assetAlreadyExists === 0
                        ? asset?.projectExpectedNoOfJobs.toString()
                        : '',
                  },
                  {
                    name: 'Project Intended Social Benefits',
                    value: asset?.projectIntendedSocialBenefits ?? '',
                  },
                  {
                    name: 'On-Site Security Personnel Available',
                    value:
                      asset?.securityMeasuresOnSiteSecurityPersonnel === 1
                        ? 'Yes'
                        : '',
                  },
                  {
                    name: 'Service Level Agreements Available',
                    value: asset?.contractualProtectionSLA === 1 ? 'Yes' : '',
                  },
                  {
                    name: 'Surveillance Systems Available',
                    value: asset?.contractualProtectionSLA === 1 ? 'Yes' : '',
                  },
                  {
                    name: 'Perimeter Security Available',
                    value: asset?.contractualProtectionSLA === 1 ? 'Yes' : '',
                  },
                  {
                    name: 'Revenue Guarantees Available',
                    value:
                      asset?.contractualProtectionRevGuarantees === 1
                        ? 'Yes'
                        : '',
                  },
                  {
                    name: 'Critical Infrastructure Protections Available',
                    value:
                      asset?.securityMeasuresCriticalInfraProtections === 1
                        ? 'Yes'
                        : '',
                  },
                  {
                    name: 'Insurance Company Name',
                    value: asset?.insuranceCompanyName ?? '',
                  },
                  {
                    name: 'Technical Partners',
                    value: asset?.projectTechnicalPartners ?? '',
                  },
                  {
                    name: 'Financial Partners',
                    value: asset?.projectFinancialPartners ?? '',
                  },
                ],
              });
            }}
            icon={
              <img
                src="/images/asset_information.png"
                height={40}
                width={40}
                alt="Asset Information"
              />
            }
            title="Asset Information"
          />
          <AsseCategoryItem
            onClick={() => {
              setShowAssetDetailsModal(true);
              setAssetPropItems({
                title: 'Asset Token & Sale Information',
                props: [
                  {
                    name: 'Sales Window',
                    value: `${formatDateLong(asset?.salesStart!)} - ${formatDateLong(asset?.salesEnd!)}`,
                  },
                  {
                    name: 'Cap Amount',
                    value: `${getFiatValue(Number(asset?.capAmountInFiat))} ${fiatCurrency}`,
                  },
                  {
                    name: 'Cap Quantity',
                    value: `${getFiatValue(asset?.capQuantity!)} ${asset?.assetCode?.toUpperCase()}`,
                  },
                  {
                    name: 'Cap Duration',
                    value: `${asset?.capDurationInDays} days`,
                  },
                  {
                    name: 'Proceed Payout Cycle',
                    value: capitalizeEachWord(
                      asset?.proceedCycle?.toLowerCase() ?? '',
                    ),
                  },
                  {
                    name: 'Exempted Countries',
                    value: asset?.exemptedCountries?.replace(/,/g, ', ') ?? '',
                  },
                ],
              });
            }}
            icon={
              <img
                src="/images/asset_token_information.png"
                height={40}
                width={40}
                alt="Asset Token & Sale Information"
              />
            }
            title="Asset Token & Sale Information"
          />
          {asset?.assetAlreadyExists == 0 && (
            <>
              <AsseCategoryItem
                onClick={() => {
                  setShowAssetDetailsModal(true);
                  setAssetPropItems({
                    title: 'Asset Financial Information',
                    props: [
                      {
                        name: 'Estimated Project IRR',
                        value: formatNumber(asset?.estimatedProjectIRR ?? 0),
                      },
                      {
                        name: 'Estimated Project ROI',
                        value: formatNumber(asset?.estimatedProjectROI ?? 0),
                      },
                      {
                        name: 'Estimated Project NPV at Launch (Day 1)',
                        value: formatNumber(asset?.estimatedProjectNPV ?? 0),
                      },
                      {
                        name: 'Estimated Project Payback Periods (in Months)',
                        value:
                          asset?.estimatedProjectPaybackPeriodsInMonths.toString() ??
                          '',
                      },
                      {
                        name: 'All Key Assumptions Including Values Assumed',
                        value: asset?.keyAssumptionsList ?? '',
                      },
                    ],
                  });
                }}
                icon={
                  <img
                    src="/images/verification-documents.png"
                    height={40}
                    width={40}
                    alt="Asset Financial Information"
                  />
                }
                title="Asset Financial Information"
              />
              <AsseCategoryItem
                onClick={() => {
                  setShowAssetDetailsModal(true);
                  setAssetPropItems({
                    title: 'Project Risk Assessment',
                    props: [
                      {
                        name: 'Legal Risks Identified',
                        value: asset?.projectIdentifiedLegalRisks ?? '',
                      },
                      {
                        name: 'Regulatory Risks Identified',
                        value: asset?.projectIdentifiedRegulatoryRisks ?? '',
                      },
                      {
                        name: 'Operational/Execution Risks Identified',
                        value:
                          asset?.projectIdentifiedOperationalOrExecutionRisks ??
                          '',
                      },
                      {
                        name: 'Market Risks Identified',
                        value: asset?.projectIdentifiedMarketRisks ?? '',
                      },
                      {
                        name: 'Other Relevant Risks Identified',
                        value: asset?.projectIdentifiedOtherRelevantRisks ?? '',
                      },
                    ],
                  });
                }}
                icon={
                  <img
                    src="/images/verification-documents.png"
                    height={40}
                    width={40}
                    alt="Project Risk Assessment"
                  />
                }
                title="Project Risk Assessment"
              />
            </>
          )}
          <AsseCategoryItem
            onClick={() => {
              setShowAssetDetailsModal(true);
              setAssetPropItems({
                title: 'Stakeholders Information',
                props: [
                  {
                    name: 'Regulator',
                    value: capitalizeEachWord(
                      regulatorName?.toLowerCase() ?? '',
                    ),
                  },
                  {
                    name: 'Asset Custodian',
                    value:
                      asset?.approvedAssetCustodianInfo?.assetCustodianName ??
                      '',
                  },
                  {
                    name: 'Asset Manager',
                    value: asset?.assetManagerInfo?.assetManagerName ?? '',
                  },
                  {
                    name: 'Issuing House',
                    value: capitalizeEachWord(
                      asset?.assetIssuingHouseInfo?.assetIssuingHouseName?.toLowerCase() ??
                        '',
                    ),
                  },
                  {
                    name: 'Legal/Professional Advisor',
                    value: capitalizeEachWord(
                      asset?.legalAdvisor?.toLowerCase() ?? '',
                    ),
                  },
                  {
                    name: 'Rating Agency',
                    value: capitalizeEachWord(
                      asset?.ratingAgencyInfo.agencyName.toLowerCase() ?? '',
                    ),
                  },
                  {
                    name: 'Trustees',
                    value: capitalizeEachWord(
                      asset?.trusteeInfo.trusteeName.toLowerCase() ?? '',
                    ),
                  },
                ],
              });
            }}
            icon={
              <img
                src="/images/stakeholders.png"
                height={40}
                width={40}
                alt="Stakeholders Information"
              />
            }
            title="Stakeholders Information"
          />
          <AsseCategoryItem
            onClick={() => {
              setShowAssetDetailsModal(true);
              setAssetPropItems({
                title: 'Legal & Compliance Information',
                props: [
                  {
                    name: 'Free of liens, mortgages, and outstanding loans.',
                    value:
                      asset.IsFreeFromLiensAndEncumbrances === 1 ? 'Yes' : 'No',
                  },
                  {
                    name: 'Not pledged as collateral for any debts and has no use restrictions.',
                    value: asset.undertakingNotCollateral === 1 ? 'Yes' : 'No',
                  },
                  {
                    name: 'No third party has any claims, rights, or interests in this asset.',
                    value: asset.undertakingNoClaims === 1 ? 'Yes' : 'No',
                  },
                  {
                    name: 'Free of any foreclosure, bankruptcy proceedings, legal disputes, judgments, or court-ordered payments.',
                    value: asset.undertakingNoForeclosure === 1 ? 'Yes' : 'No',
                  },
                  {
                    name: 'Complies with all environmental and land use regulations and is free of violations.',
                    value: asset.complianceNoViolation === 1 ? 'Yes' : 'No',
                  },
                  {
                    name: 'All necessary permits, licenses, and approvals for the use and ownership of this asset are in place.',
                    value: asset.complianceNoViolation === 1 ? 'Yes' : 'No',
                  },
                  {
                    name: 'No unpaid taxes, utility bills, fees, or other property-related expenses associated with this asset.',
                    value:
                      asset.outstandingFinancialRespNoDebts === 1
                        ? 'Yes'
                        : 'No',
                  },
                  {
                    name: 'Asset does not have any hidden liabilities or obligations that have not been disclosed.',
                    value:
                      asset.outstandingFinancialRespNoHiddenLiabilities === 1
                        ? 'Yes'
                        : 'No',
                  },
                  {
                    name: 'Asset is adequately insured against risks such as fire, theft, and natural disasters.',
                    value:
                      asset.riskManagementFullyInsured === 1 ? 'Yes' : 'No',
                  },
                  {
                    name: 'The declared value of this asset reflects its current market value and condition.',
                    value:
                      asset.riskManagementDeclaredValue === 1 ? 'Yes' : 'No',
                  },
                  {
                    name: 'Asset is not affected by undisclosed easements, rights of way, expropriation, or condemnation.',
                    value:
                      asset.physicalConditionNoUndisclosedEasements === 1
                        ? 'Yes'
                        : 'No',
                  },
                  {
                    name: 'Asset is structurally sound and has no unresolved maintenance or safety issues.',
                    value: asset.physicalConditionSound === 1 ? 'Yes' : 'No',
                  },
                  {
                    name: 'Asset is not subject to any agreements, such as leases or contracts, that could limit its use or transfer.',
                    value: asset.physicalConditionNolease === 1 ? 'Yes' : 'No',
                  },
                ],
              });
            }}
            icon={
              <img
                src="/images/stakeholders.png"
                height={40}
                width={40}
                alt="Legal & Compliance Information"
              />
            }
            title="Legal & Compliance Information"
          />
          <AsseCategoryItem
            onClick={() => {
              setShowVerificationDocumentsModal(true);
            }}
            icon={
              <img
                src="/images/verification-documents.png"
                height={40}
                width={40}
                alt="Asset Documents"
              />
            }
            title="Asset Documents"
          />
        </div>
      </div>
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
      {/* Public key copied modal */}
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
      <BuyTokenModal
        show={showBuyTokenModal}
        onClose={() => {
          setShowBuyTokenModal(false);
        }}
        asset={asset}
      />
      <Modal
        showModal={showSuccessModal}
        onClose={() => {
          setShowSuccessModal(false);
        }}
      >
        <div className="flex flex-col space-y-5 items-center w-full py-10 justify-center">
          <img src="/images/success.png" alt="success" />
          <div className="flex flex-col text-center space-y-5 items-center w-2/3 md:px-10 justify-center">
            <p className="text-primary-800 text-md xl:text-lg font-montserratSemiBold">
              You have successfully expressed interest on {asset.assetCode}
            </p>
          </div>
        </div>
      </Modal>
      <Modal
        showModal={showAssetDetailsModal}
        onClose={() => {
          setShowAssetDetailsModal(false);
        }}
      >
        <div className="flex flex-col items-center w-full px-5 space-y-5 py-5 justify-center">
          <p className="text-primary-800 w-full text-lg font-montserratSemiBold">
            {assetPropItems?.title}
          </p>
          {assetPropItems?.props.map(
            (item, index) =>
              item.value && (
                <div
                  key={index}
                  className="flex w-full py-3 px-2 xl:px-5 xl:space-y-5 rounded-xl flex-col items-center bg-primary-100"
                >
                  <div className="flex flex-col space-y-2 w-full justify-between">
                    <p>{item.name}</p>
                    <p className="font-montserratSemiBold">{item.value}</p>
                  </div>
                </div>
              ),
          )}
          <div className="flex w-3/4 space-x-5">
            <Button
              label="Close"
              additionalClasses="bg-primary-600 text-white font-montserratSemiBold"
              onclick={() => {
                setShowAssetDetailsModal(false);
              }}
            />
          </div>
          <div />
        </div>
      </Modal>
      <Modal
        showModal={showVerificationDocumentsModal}
        onClose={() => {
          setShowVerificationDocumentsModal(false);
        }}
      >
        <div className="flex flex-col items-center w-full px-5 space-y-5 py-5 justify-center">
          <p className="text-primary-800 w-full text-lg font-montserratSemiBold">
            Asset Documents
          </p>
          <div className="flex w-full py-3 px-2 xl:px-5 xl:space-y-5 rounded-xl flex-col items-center bg-primary-100">
            {asset?.AssetTokenizationDocuments?.length > 0 ? (
              asset?.AssetTokenizationDocuments.map((item, index) => (
                <div
                  key={index}
                  className="flex flex-col space-y-2 w-full justify-between"
                >
                  <a href={item.documentUrl} target="_blank">
                    <p className="font-montserratSemiBold underline">
                      {item.documentTitle}
                    </p>
                  </a>
                </div>
              ))
            ) : (
              <></>
            )}
          </div>
          <div className="flex w-3/4 space-x-5">
            <Button
              label="Close"
              additionalClasses="bg-primary-600 text-white font-montserratSemiBold"
              onclick={() => {
                setShowVerificationDocumentsModal(false);
              }}
            />
          </div>
        </div>
      </Modal>
    </div>
  );
}
