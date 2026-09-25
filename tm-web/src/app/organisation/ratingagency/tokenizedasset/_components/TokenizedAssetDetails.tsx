"use client";
import styled from "styled-components";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { FaArrowLeft } from "react-icons/fa6";
import Image from "next/image";
import AssetCard from "@/app/(dashboard)/assettokenization/components/AssetCard";
import profile from "../../../../assets/images/profileblue.svg";
import SecondaryButton from "@/components/SecondaryButton";
import { FaRegEdit } from "react-icons/fa";
import DetailsSection from "../../../components/DetailsSection";
import PrimaryButton from "@/components/PrimaryButton";
import { AiOutlineClockCircle } from "react-icons/ai";
import AssetTokenizer from "./AssetTokenizer";
import AssignedStakeHolder from "./AssignedStakeholders";
import RiskAndCompliance from "../../../components/AssetRiskAndCompliance";
import ValuationHistory from "./ValuationHistory";
import Wallets from "./Wallets";
import ExemptedCountries from "./ExemptedCountries";
import AssetDocuments from "./AssetDocuments";
import AssetProtection from "./AssetProtections";
import StepItem from "./StepItem";
import { BsThreeDots } from "react-icons/bs";
import { Dropdown, MenuProps, Badge } from "antd";
import { useState } from "react";
import { IStakeholderAssetDetailResponse } from "@/redux/api/sharedstakeholders";

const TOKENIZATION_STATUS_LABELS: Record<number, string> = {
  0: "Pending Vetting",
  1: "Awaiting Payment",
  2: "Payment Made",
  3: "Processing",
  4: "Tokenization Approved",
  5: "Primary Sale Started",
  6: "Secondary Sales",
  7: "Liquidated",
  8: "Refunded",
};

// Dummy asset data
const getDummyAsset = (id: string) => ({
  id: parseInt(id),
  assetCode: "IKY001",
  assetName: "Atlantis 1",
  assetSector: "Real Estate Sector",
  assetSubSector: "Residential Properties",
  assetType: "Single-family homes",
  assetStatus: "Existing",
  offeringType: "public",
  assetCountryLocation: "NG",
  assetAddress: "12, Clover Road Ikoyi",
  coordinates: {
    longitude: "6.8454548",
    latitude: "5.128468",
  },
  assetTokenizationStatus: 4,
  vettingStatus: 1,
  assetLogo: null,
  createdAt: "Dec 31, 2025",

  tokenizer: {
    name: "brikkle",
    profile: "BRIKKLE",
    email: "ayo@brikkle.co",
    walletAddress: "GAYBKWD76M7L7OFXQCTE6NGKQQQ5Y54DWSVDK4DNC4DZ7FW6KBZAAJEF",
  },

  owner: {
    wallets: {
      issuing: {
        address: "atprofile_iky001issuer",
        balanceNgn: "0",
        balanceUsd: "0",
      },
      distribution: {
        address: "atprofile_iky001issuer-distribution",
        balanceNgn: "0",
        balanceUsd: "0",
      },
      holding: {
        address: "atprofile_iky001issuer-holding",
        balanceNgn: "0",
        balanceUsd: "0",
      },
    },
  },

  assetValue: 200000000,
  costOutsideValuation: 10000000,
  valueRetained: 20000000,
  assetQuoteCurrency: "CNGN",

  tokenName: "IKY001",
  tokenSymbol: "IKY001",
  tokensToBeIssued: 5000,
  tokensForSale: 4642.56,
  tokensNotForSale: 357.44,
  tokenPrice: 55953.519,
  totalTokenValue: 279767593.525,
  preferredTokenizationFee: 25000000,
  salesStartDate: "Jan 6, 2026",
  salesEndDate: "Jan 30, 2026",
  capQuantity: 17.872,
  capAmount: 1000000,
  capDuration: "2 days",
  proceedPayoutCycle: "Monthly",
  payoutCurrency: "CNGN",
  amountToBeRaised: 259767593.527,

  stakeholders: {
    assetCustodian: "SunTrusts Capital Partners",
    assetManager: "testing-manager",
    regulator: "Securities and Exchange Commission SEC Nigeria",
    issuingHouse: "Test Issuing House Ltds",
    legalAdviser: "LexTrust Legal Advisory Ltd",
    ratingAgency: "Test Rating",
    trustee: "Sun Trustess ltd",
  },

  assetProtection: [
    "Insurance",
    "Contractual Protection",
    "Risk Sharing Mechanisms",
    "Governance and Oversight",
    "Environmental, Social and Governance (ESG) Safeguards",
    "Security Measures",
    "Legal/Financial Counsel",
  ],

  documents: [
    {
      name: "statutory licenses and permits/compliance certificates-968817d6.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "asset valuation certificate from licensed appraisers-2cec1d8b.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "financial advisors contact/proof-13ef6f7f.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "proof of asset existence-7240eebf.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "proof of additional cost incurred outside valuation-188996d7.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "upload signed declaration form-c7df5ad3.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "proof of asset address-566ed353.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "comprehensive insurance-da26fa71.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "insurance policies-59dd3bd8.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "proof of ownership-461a6cd8.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "original ownership agreements-c7368bea.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "asset status verification-961d85e1.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "proof of asset's condition-a59ab41c.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "contractual protections-0e3158c5.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "revenue guarantees-ea8433fb.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "asset financial performance report-31b63835.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "government issued id of original asset owner-ad1b9219.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "service level agreements-651e7aad.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "proof of compliance with local laws and regulation-f0a05944.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "proof of compliance with environmental standards and regulations-53c57a4a.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    {
      name: "legal advisors contact/proof-0bbb268b.png",
      type: "PDF",
      uploadedBy: "Creator",
    },
    { name: "Proof of Payment 1", type: "PDF", uploadedBy: "Creator" },
  ],

  receivingAccount: {
    bank: "STERLING",
    accountName: "BRIKKLE INC",
    accountNumber: "13456667778",
  },

  exemptedCountries: ["US"],
});

interface TokenizedAssetDetailsProps {
  assetDetail?: IStakeholderAssetDetailResponse;
  isLoading?: boolean;
}

const TokenizedAssetDetailsPage = ({
  assetDetail,
  isLoading,
}: TokenizedAssetDetailsProps) => {
  const params = useParams();
  const assetId = params.id as string;
  const dummyAsset = getDummyAsset(assetId);
  const data = assetDetail?.data;

  const asset = {
    ...dummyAsset,
    assetCode: data?.assetCode ?? dummyAsset.assetCode,
    assetName: data?.assetName ?? dummyAsset.assetName,
    assetSector: data?.assetSector ?? dummyAsset.assetSector,
    assetSubSector: data?.assetSubSector ?? dummyAsset.assetSubSector,
    assetType: data?.assetType ?? dummyAsset.assetType,
    assetTokenizationStatus:
      data?.assetTokenizationStatus ?? dummyAsset.assetTokenizationStatus,
    assetQuoteCurrency:
      data?.assetQuoteCurrency ?? dummyAsset.assetQuoteCurrency,
    assetValue: data?.assetCurrentValue ?? dummyAsset.assetValue,
    totalTokenValue: data?.valueOfTokenizedAsset ?? dummyAsset.totalTokenValue,
    tokensToBeIssued:
      data?.numberOfTokenToBeIssued ?? dummyAsset.tokensToBeIssued,
    tokenName: data?.assetCode ?? dummyAsset.tokenName,
    tokenSymbol: data?.assetCode ?? dummyAsset.tokenSymbol,
    createdAt: data?.createdAt
      ? new Date(data.createdAt).toLocaleDateString()
      : dummyAsset.createdAt,
    assetLogo: data?.asset_profile?.logo_url || null,
    description: data?.asset_profile?.description || "N/A",
    offeringType: data?.asset_profile?.offering_type || "N/A",
    assetCountryLocation: data?.asset_profile?.country || "N/A",
    assetAddress: data?.asset_profile?.physical_address || "N/A",
    coordinates: {
      longitude: data?.asset_profile?.longitude || "N/A",
      latitude: data?.asset_profile?.latitude || "N/A",
    },
    assetStatus: TOKENIZATION_STATUS_LABELS[data?.assetTokenizationStatus ?? 0] ?? "Unknown",
    costOutsideValuation: data?.ownership?.cost_outside_valuation ?? 0,
    valueRetained: data?.ownership?.retained_or_contributed_value ?? 0,
    tokensForSale: data?.offering?.tokens_for_sale ?? 0,
    tokensNotForSale: data?.offering?.tokens_not_for_sale ?? 0,
    tokenPrice: data?.offering?.price_per_token ?? 0,
    preferredTokenizationFee: data?.tokenization_fees?.fiat_amount ?? 0,
    salesStartDate: data?.offering?.sales_start ? new Date(data.offering.sales_start).toLocaleDateString() : "N/A",
    salesEndDate: data?.offering?.sales_end ? new Date(data.offering.sales_end).toLocaleDateString() : "N/A",
    capQuantity: data?.offering?.cap_quantity ?? 0,
    capAmount: data?.offering?.cap_amount ?? 0,
    capDuration: data?.offering ? `${data.offering.cap_duration_days} days` : "N/A",
    proceedPayoutCycle: data?.offering?.payout_cycle || "N/A",
    payoutCurrency: data?.offering?.payout_currency || "N/A",
    amountToBeRaised: data?.offering?.amount_to_be_raised ?? 0,
    receivingAccount: {
      bank: data?.receiving_account?.bank_name || "N/A",
      accountName: data?.receiving_account?.account_name || "N/A",
      accountNumber: data?.receiving_account?.account_number || "N/A",
    },
  };

  const router = useRouter();
  const [openMenu, setOpenMenu] = useState(false);

  const handleMenuClick: MenuProps["onClick"] = ({ key }) => {
    switch (key) {
      case "1": // Primary Sale Details
        router.push("/organisation/ratingagency/tokenizedasset/primarysalesdetails");
        break;
      case "2": // Commitment Details
        router.push("/organisation/ratingagency/tokenizedasset/commitment-details");
        break;
      case "3": // View Token Holders
        router.push("/organisation/ratingagency/tokenizedasset/tokenholders");
        break;
      case "4": // Early Exit
        router.push("/organisation/ratingagency/tokenizedasset/earlyexit");
        break;

      default:
        break;
    }
  };

  const menuItems: MenuProps["items"] = [
    {
      key: "1",
      label: "Primary Sale Details",
    },
    {
      key: "2",
      label: "Commitment Details",
    },

    {
      key: "3",
      label: "View Token Holders",
    },

    {
      key: "4",
      label: "Early Exit",
    },

  ];

  const formatAmount = (value: number) =>
    value.toLocaleString(undefined, {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });

  const assetInformation = [
    { label: "Asset Code", value: asset.assetCode },
    { label: "Date Created", value: asset.createdAt },
    { label: "Sector", value: asset.assetSector },
    { label: "SubSector", value: asset.assetSubSector },
    { label: "Type", value: asset.assetType },
    { label: "Asset Status", value: asset.assetStatus },
    { label: "Offering Type", value: asset.offeringType },
    { label: "Country", value: asset.assetCountryLocation },
    { label: "Physical Asset Address", value: asset.assetAddress },
    { label: "Longitude", value: asset.coordinates.longitude },
    { label: "Latitude", value: asset.coordinates.latitude },
  ];

  const assetOwner = [
    { label: "Ownership", value: data?.ownership?.type || "N/A" },
    { label: "Owner Type", value: data?.ownership?.kind || "N/A" },
    { label: "Website", value: data?.asset_profile?.website || "N/A" },
    { label: "Physical Address", value: data?.asset_profile?.physical_address || "N/A" },
  ];

  const assetValue = [
    {
      label: "Current Asset Value",
      value: `${formatAmount(asset.assetValue)} ${asset.assetQuoteCurrency}`,
    },
    {
      label: "Cost Outside Valuation",
      value: `${formatAmount(asset.costOutsideValuation)} NGN`,
    },
    {
      label: "Value to be retained",
      value: `${formatAmount(asset.valueRetained)} ${asset.assetQuoteCurrency}`,
    },
  ];

  const tokenSales = [
    {
      label: "Tokens to be Issued",
      value: `${formatAmount(data?.offering?.tokens_issued ?? asset.tokensToBeIssued)} ${asset.tokenSymbol}`,
    },
    {
      label: "Token for Sale",
      value: `${formatAmount(data?.offering?.tokens_for_sale ?? 0)} ${asset.tokenSymbol}`,
    },
    {
      label: "Token not for Sale",
      value: `${formatAmount(data?.offering?.tokens_not_for_sale ?? 0)} ${asset.tokenSymbol}`,
    },
    { label: "Asset Quote Currency", value: asset.assetQuoteCurrency },
    {
      label: "Price Per Token",
      value: `${formatAmount(data?.offering?.price_per_token ?? 0)} ${asset.assetQuoteCurrency}`,
    },
    {
      label: "Value of Total Tokens",
      value: `${formatAmount(asset.totalTokenValue)} ${asset.assetQuoteCurrency}`,
    },
    {
      label: "Preferred Tokenization Fee",
      value: `${formatAmount(data?.tokenization_fees?.fiat_amount ?? 0)} ${asset.assetQuoteCurrency}`,
    },
    {
      label: "Sales Start Date",
      value: data?.offering?.sales_start
        ? new Date(data.offering.sales_start).toLocaleDateString()
        : "N/A",
    },
    {
      label: "Sales End Date",
      value: data?.offering?.sales_end
        ? new Date(data.offering.sales_end).toLocaleDateString()
        : "N/A",
    },
    {
      label: "Cap Quantity",
      value: `${formatAmount(data?.offering?.cap_quantity ?? 0)} ${asset.tokenSymbol}`,
    },
    {
      label: "Cap Amount",
      value: `${formatAmount(data?.offering?.cap_amount ?? 0)} ${asset.assetQuoteCurrency}`,
    },
    {
      label: "Cap Duration",
      value: data?.offering
        ? `${data.offering.cap_duration_days} days`
        : "N/A",
    },
    {
      label: "Proceed Payout Cycle",
      value: data?.offering?.payout_cycle || "N/A",
    },
    {
      label: "Payout Currency",
      value: data?.offering?.payout_currency || "N/A",
    },
    {
      label: "Amount to be raised",
      value: `${formatAmount(data?.offering?.amount_to_be_raised ?? 0)} ${asset.assetQuoteCurrency}`,
    },
  ];
  const receivingAccount = [
    { label: "Bank", value: asset.receivingAccount.bank },
    { label: "Account Name", value: asset.receivingAccount.accountName },
    { label: "Account Number", value: asset.receivingAccount.accountNumber },
  ];

  const steps =
    data?.tokenization_timeline?.stages.map((stage) => ({
      action: stage.description,
      date:
        stage.status_id === 0 && data.tokenization_timeline?.date_submitted
          ? new Date(data.tokenization_timeline.date_submitted).toLocaleString()
          : stage.current && data.tokenization_timeline?.date_approved
            ? new Date(data.tokenization_timeline.date_approved).toLocaleString()
            : "",
    })) ?? [];

  return (
    <TokenizationStep>
      <StepperContainer>
        {steps.map((step, index) => (
          <StepItem
            key={step.action}
            step={step}
            index={index}
            currentStep={data?.tokenization_timeline?.current_status_id ?? 0}
            processingStep={null}
            isDueDiligenceDisapproved={false}
          />
        ))}
      </StepperContainer>
      <Container>
      {/* Asset Information */}
      <Section>
        <HeadingContent>
          <Heading>Asset Information</Heading>

          <ButtonContainer>
            {/* <PrimaryButton
              buttonStyle={{ width: "154px" }}
              // onClick={() => {
              //   router.push("/organisation/ratingagency/tokenizedasset/dividendandyield");
              // }}
            >
              Authorize
            </PrimaryButton> */}
            <DropdownStyles>
              <Dropdown
                menu={{ items: menuItems, onClick: handleMenuClick }}
                trigger={["click"]}
                placement="bottomRight"
                overlayClassName="due-diligence-dropdown"
              >
                <div>
                  <DotBadge dot>
                    <SecondaryButton buttonStyle={{ width: "48px" }}>
                      <BsThreeDots />
                    </SecondaryButton>
                  </DotBadge>
                </div>
              </Dropdown>
            </DropdownStyles>
          </ButtonContainer>
        </HeadingContent>
        <PropertyInfoSection>
          <AssetInfoHeader>
            {asset?.assetLogo ? (
              <Image
                src={asset.assetLogo}
                alt="Asset Logo"
                width={40}
                height={40}
                style={{ borderRadius: "50%" }}
                unoptimized={true}
              />
            ) : (
              <Avatar />
            )}
            <div>
              <AssetTitle style={{ whiteSpace: "pre-wrap" }}>
                {asset.assetName}
              </AssetTitle>
              <Status>
                <AiOutlineClockCircle />
                {TOKENIZATION_STATUS_LABELS[asset.assetTokenizationStatus] ??
                  "Unknown"}
              </Status>
            </div>
          </AssetInfoHeader>
          <Description>{asset.description}</Description>
          <AssetDetails>
            {assetInformation.map((item, index) => (
              <AssetCard key={index} label={item.label} value={item.value} />
            ))}
          </AssetDetails>
        </PropertyInfoSection>
      </Section>

      <Section>
        <AssetTokenizer username={data?.tokenizer_username} />
      </Section>

      <Section>
        <HeadingContent>
          <Heading>Asset Owner</Heading>
          {/* <SecondaryButton
            buttonStyle={{
              width: "100px",
              marginRight: "20px",
              display: "flex",
              alignItems: "center",
              gap: "2px",
              marginBottom: "16px",
            }}
          >
            <FaRegEdit />
            Edit
          </SecondaryButton> */}
        </HeadingContent>
        <Wrapper>
          {assetOwner.map((item, index) => (
            <AssetCard key={index} label={item.label} value={item.value} />
          ))}
        </Wrapper>
      </Section>

      <Section>
        <Wallets wallets={data?.wallets} />
      </Section>

      <Section>
        <AssignedStakeHolder stakeholders={data?.assigned_stakeholders} />
      </Section>

      {/* Asset Value */}
      <Section>
        <DetailsSection title="Asset Value" items={assetValue} showEdit={false} />
      </Section>

      {/* Asset Images */}
      {/* <Section>
        <HeadingContent>
          <Heading>Asset Images</Heading>
          <SecondaryButton
            buttonStyle={{
              width: "120px",
              marginRight: "20px",
              marginBottom: "16px",
            }}
          >
            Add Image
          </SecondaryButton>
        </HeadingContent>
        <Wrapper
          style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            height: "150px",
          }}
        >
          <p style={{ color: "#828282" }}>0 image uploaded</p>
        </Wrapper>
      </Section> */}
      {/* 
        <Section>
          <KeyMilestones />
        </Section> */}
      <Section>
        <AssetProtection protection={data?.asset_protection} />
      </Section>

      {/* Asset Token & Sales */}
      <Section>
        <DetailsSection title="Asset Token & Sales" items={tokenSales} showEdit={false} />
      </Section>

      {/* Receiving Account Details */}
      <Section>
        <Heading2>Receiving Account Details</Heading2>
        <Wrapper>
          {receivingAccount.map((item, index) => (
            <AssetCard key={index} label={item.label} value={item.value} />
          ))}
        </Wrapper>
      </Section>

      <Section>
        <ExemptedCountries countries={data?.exempted_countries} />
      </Section>

      <Section>
        <AssetDocuments documents={data?.tokenization_documents} />
      </Section>

      <Section>
        <RiskAndCompliance data={data?.risk_and_compliance} />
      </Section>

      {/* <Section>
          <ValuationHistory />
        </Section> */}
      </Container>
    </TokenizationStep>
  );
};

export default TokenizedAssetDetailsPage;
const StepperContainer = styled.div`
  width: 230px;
  flex-shrink: 0;
  padding: 20px;
  border-radius: 24px;
  background-color: #fff;
`;
const TokenizationStep = styled.div`
  display: flex;
  gap: 20px;
  width: 100%;
  max-width: 100%;
  min-width: 0;
  align-items: flex-start;

  & * {
    min-width: 0;
  }
`;

const Container = styled.section`
  padding: 24px;
  border-radius: 24px;
  background-color: #ffffff;
`;

const ButtonContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const Section = styled.div`
  margin-top: 40px;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
`;

const Heading2 = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin-bottom: 16px;
`;

const PropertyInfoSection = styled.div`
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  padding: 24px;
`;

const Wrapper = styled.section`
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
`;

const AssetInfoHeader = styled.div`
  display: flex;
  column-gap: 8px;
  align-items: center;
  margin-bottom: 8px;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;

const AssetTitle = styled.p`
  font-size: 14px;
  font-weight: 600;
  line-height: 16px;

  color: #00225a;
  text-transform: capitalize;
`;

const AssetDetails = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 28px;
  padding-top: 10px;
  width: 100%;
`;

const DotBadge = styled(Badge)`
  .ant-badge-dot {
    top: 20px;
    right: 2px;
    width: 8px;
    height: 8px;
  }
`;

const DropdownStyles = styled.div`
  .due-diligence-dropdown .ant-dropdown-menu-item {
    font-family: "Montserrat", sans-serif !important;
    font-weight: 400 !important;
    font-size: 14px !important;
    line-height: 28px !important;
    letter-spacing: 0% !important;
    color: #00225a !important;
  }

  .due-diligence-dropdown .ant-dropdown-menu-item:hover {
    background: #f2f6f9;
    color: #00225a;
  }
`;

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
`;

const Description = styled.p`
  font-size: 14px;
  color: #828282;
  margin: 8px 0;
  width: 100%;
  line-height: 25px;
`;

const Status = styled.div`
  background-color: #f2f6f9;
  color: #007cdf;
  font-size: 12px;
  font-weight: 500;
  padding: 6px 10px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-top: 2px;
`;
