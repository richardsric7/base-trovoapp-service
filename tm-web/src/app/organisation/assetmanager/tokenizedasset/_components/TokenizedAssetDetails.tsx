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
import { FormEvent, useState } from "react";
import { IStakeholderAssetDetailResponse } from "@/redux/api/sharedstakeholders";
import Loader from "@/components/Loader";
import { useUploadStakeholderDocumentMutation } from "@/redux/api/sharedstakeholders";
import { useSubmitAssetManagerValuationMutation } from "@/redux/api/assetManager";
import { showErrorToast, showSuccessToast } from "@/components";

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

const emptyValuationForm = {
  valuation: "",
  currency: "",
  methodology: "",
  valuationDate: "",
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
  const data = assetDetail?.data;
  const router = useRouter();
  const [valuationOpen, setValuationOpen] = useState(false);
  const [valuationForm, setValuationForm] = useState(emptyValuationForm);
  const [valuationFile, setValuationFile] = useState<File | null>(null);
  const [uploadDocument, { isLoading: isUploadingValuation }] =
    useUploadStakeholderDocumentMutation();
  const [submitValuation, { isLoading: isSubmittingValuation }] =
    useSubmitAssetManagerValuationMutation();

  if (isLoading) return <Loader />;
  if (!data) return <EmptyState>Unable to load asset details.</EmptyState>;

  const asset = {
    assetCode: data.assetCode || "N/A",
    assetName: data.assetName || "N/A",
    assetSector: data.assetSector || "N/A",
    assetSubSector: data.assetSubSector || "N/A",
    assetType: data.assetType || "N/A",
    assetTokenizationStatus: data.assetTokenizationStatus,
    assetQuoteCurrency: data.assetQuoteCurrency || "N/A",
    assetValue: data.assetCurrentValue ?? 0,
    totalTokenValue: data.valueOfTokenizedAsset ?? 0,
    tokensToBeIssued: data.numberOfTokenToBeIssued ?? 0,
    tokenSymbol: data.assetCode || "Token",
    assetLogo: data.asset_profile?.logo_url || null,
    description: data.asset_profile?.description || "N/A",
    createdAt: data.createdAt
      ? new Date(data.createdAt).toLocaleDateString()
      : "N/A",
  };

  const handleMenuClick: MenuProps["onClick"] = ({ key }) => {
    switch (key) {
      case "1": // Primary Sale Details
        router.push(
          "/organisation/assetmanager/tokenizedasset/primarysalesdetails",
        );
        break;
      case "2": // Commitment Details
        router.push(
          "/organisation/assetmanager/tokenizedasset/commitment-details",
        );
        break;
      case "3": // View Token Holders
        router.push("/organisation/assetmanager/tokenizedasset/tokenholders");
        break;
      case "4": // Early Exit
        router.push("/organisation/assetmanager/tokenizedasset/earlyexit");
        break;
      case "5": // Asset Income Report
        router.push(
          "/organisation/assetmanager/tokenizedasset/assetincomereport",
        );
        break;
      case "6":
        setValuationForm((current) => ({
          ...current,
          currency:
            asset.assetQuoteCurrency === "N/A"
              ? current.currency
              : asset.assetQuoteCurrency,
        }));
        setValuationOpen(true);
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

    {
      key: "5",
      label: "Asset Income Report",
    },
    {
      key: "6",
      label: "Submit Valuation",
    },
  ];

  const canSubmitValuation = Boolean(
    Number(valuationForm.valuation) > 0 &&
    valuationForm.currency.trim() &&
    valuationForm.methodology.trim() &&
    valuationForm.valuationDate &&
    valuationFile,
  );

  const handleValuationSubmit = async (event: FormEvent) => {
    event.preventDefault();
    if (!canSubmitValuation || !valuationFile) return;
    try {
      const documentPayload = new FormData();
      documentPayload.append("document_file", valuationFile);
      documentPayload.append("asset_id", assetId);
      documentPayload.append("category", "valuation_report");
      documentPayload.append("title", valuationFile.name);
      const uploadedDocument = await uploadDocument(documentPayload).unwrap();

      await submitValuation({
        asset_id: assetId,
        valuation: valuationForm.valuation,
        currency: valuationForm.currency.trim().toUpperCase(),
        methodology: valuationForm.methodology.trim(),
        valuation_date: valuationForm.valuationDate,
        report_document_id: uploadedDocument.data.id,
      }).unwrap();
      showSuccessToast("Valuation submitted successfully");
      setValuationForm(emptyValuationForm);
      setValuationFile(null);
      setValuationOpen(false);
    } catch (error: any) {
      showErrorToast(
        error?.data?.message ??
          error?.message ??
          error?.data?.error ??
          error?.error ??
          "Unable to submit valuation",
      );
    }
  };

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
    {
      label: "Asset Status",
      value:
        TOKENIZATION_STATUS_LABELS[asset.assetTokenizationStatus] ?? "Unknown",
    },
    {
      label: "Offering Type",
      value: data.asset_profile?.offering_type || "N/A",
    },
    { label: "Country", value: data.asset_profile?.country || "N/A" },
    {
      label: "Physical Asset Address",
      value: data.asset_profile?.physical_address || "N/A",
    },
    { label: "Longitude", value: data.asset_profile?.longitude || "N/A" },
    { label: "Latitude", value: data.asset_profile?.latitude || "N/A" },
    { label: "Compliance Status", value: data.complianceStatus || "N/A" },
    { label: "Custody Status", value: data.custodyStatus || "N/A" },
    {
      label: "Token Holders",
      value: data.tokenHolderCount?.toLocaleString() ?? "0",
    },
  ];

  const assetOwner = [
    { label: "Ownership", value: data.ownership?.type || "N/A" },
    { label: "Owner Type", value: data.ownership?.kind || "N/A" },
    { label: "Website", value: data.asset_profile?.website || "N/A" },
    {
      label: "Physical Address",
      value: data.asset_profile?.physical_address || "N/A",
    },
  ];

  const assetValue = [
    {
      label: "Current Asset Value",
      value: `${formatAmount(asset.assetValue)} ${asset.assetQuoteCurrency}`,
    },
    {
      label: "Cost Outside Valuation",
      value: `${formatAmount(data.ownership?.cost_outside_valuation ?? 0)} ${asset.assetQuoteCurrency}`,
    },
    {
      label: "Value to be retained",
      value: `${formatAmount(data.ownership?.retained_or_contributed_value ?? 0)} ${asset.assetQuoteCurrency}`,
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
    { label: "Bank", value: data.receiving_account?.bank_name || "N/A" },
    {
      label: "Account Name",
      value: data.receiving_account?.account_name || "N/A",
    },
    {
      label: "Account Number",
      value: data.receiving_account?.account_number || "N/A",
    },
  ];

  const steps =
    data.tokenization_timeline?.stages.map((stage) => ({
      action: stage.description,
      date:
        stage.current && data.tokenization_timeline?.date_approved
          ? new Date(data.tokenization_timeline.date_approved).toLocaleString()
          : stage.status_id === 0 && data.tokenization_timeline?.date_submitted
            ? new Date(
                data.tokenization_timeline.date_submitted,
              ).toLocaleString()
            : "",
    })) ?? [];

  return (
    <TokenizationStep>
      <StepperContainer>
        {steps.map((step, index) => (
          <StepItem
            key={index}
            step={step}
            index={index}
            currentStep={data.tokenization_timeline?.current_status_id ?? 0}
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
                //   router.push("/organisation/assetmanager/tokenizedasset/dividendandyield");
                // }}
              >
                Update
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
              <Avatar />
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
          <AssetProtection protection={data.asset_protection} />
        </Section>

        {/* Asset Token & Sales */}
        <Section>
          <DetailsSection title="Asset Token & Sales" items={tokenSales} showEdit={false} />
        </Section>

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
          <AssetDocuments documents={data.tokenization_documents} />
        </Section>

        <Section>
        <RiskAndCompliance data={data?.risk_and_compliance} />
        </Section>

        <Section>
          <ValuationHistory assetId={assetId} />
        </Section>

        {valuationOpen && (
          <ModalOverlay
            onMouseDown={() =>
              !isSubmittingValuation &&
              !isUploadingValuation &&
              setValuationOpen(false)
            }
          >
            <ModalCard onMouseDown={(event) => event.stopPropagation()}>
              <ModalHeader>
                <div>
                  <ModalTitle>Submit Valuation</ModalTitle>
                  <ModalSubtitle>
                    Add the asset valuation and supporting report.
                  </ModalSubtitle>
                </div>
                <CloseButton
                  type="button"
                  onClick={() => setValuationOpen(false)}
                >
                  ×
                </CloseButton>
              </ModalHeader>
              <ValuationForm onSubmit={handleValuationSubmit}>
                <Field>
                  <FieldLabel htmlFor="valuation-amount">
                    Valuation amount <RequiredMark>*</RequiredMark>
                  </FieldLabel>
                  <Input
                    id="valuation-amount"
                    type="number"
                    min="0.01"
                    step="0.01"
                    value={valuationForm.valuation}
                    onChange={(event) =>
                      setValuationForm((current) => ({
                        ...current,
                        valuation: event.target.value,
                      }))
                    }
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="valuation-currency">
                    Currency <RequiredMark>*</RequiredMark>
                  </FieldLabel>
                  <Input
                    id="valuation-currency"
                    placeholder={
                      asset.assetQuoteCurrency === "N/A"
                        ? "e.g. NGN"
                        : asset.assetQuoteCurrency
                    }
                    value={valuationForm.currency}
                    onChange={(event) =>
                      setValuationForm((current) => ({
                        ...current,
                        currency: event.target.value,
                      }))
                    }
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="valuation-methodology">
                    Methodology <RequiredMark>*</RequiredMark>
                  </FieldLabel>
                  <Input
                    id="valuation-methodology"
                    placeholder="e.g. Income approach"
                    value={valuationForm.methodology}
                    onChange={(event) =>
                      setValuationForm((current) => ({
                        ...current,
                        methodology: event.target.value,
                      }))
                    }
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="valuation-date">
                    Valuation date <RequiredMark>*</RequiredMark>
                  </FieldLabel>
                  <Input
                    id="valuation-date"
                    type="date"
                    value={valuationForm.valuationDate}
                    onChange={(event) =>
                      setValuationForm((current) => ({
                        ...current,
                        valuationDate: event.target.value,
                      }))
                    }
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="valuation-report">
                    Valuation report <RequiredMark>*</RequiredMark>
                  </FieldLabel>
                  <Input
                    id="valuation-report"
                    type="file"
                    accept="application/pdf,image/jpeg,image/png"
                    onChange={(event) => {
                      const file = event.target.files?.[0] ?? null;
                      if (
                        file &&
                        ![
                          "application/pdf",
                          "image/jpeg",
                          "image/png",
                        ].includes(file.type)
                      ) {
                        showErrorToast(
                          "Only PDF, JPEG and PNG reports are supported",
                        );
                        event.target.value = "";
                        setValuationFile(null);
                        return;
                      }
                      if (file && file.size > 10 * 1024 * 1024) {
                        showErrorToast(
                          "Valuation report must be 10 MB or smaller",
                        );
                        event.target.value = "";
                        setValuationFile(null);
                        return;
                      }
                      setValuationFile(file);
                    }}
                  />
                  <FieldHint>PDF, JPEG or PNG. Maximum size 10 MB.</FieldHint>
                </Field>
                <ModalActions>
                  <CancelButton
                    type="button"
                    onClick={() => setValuationOpen(false)}
                  >
                    Cancel
                  </CancelButton>
                  <SubmitButton
                    type="submit"
                    disabled={
                      !canSubmitValuation ||
                      isUploadingValuation ||
                      isSubmittingValuation
                    }
                  >
                    {isUploadingValuation
                      ? "Uploading report..."
                      : isSubmittingValuation
                        ? "Submitting..."
                        : "Submit Valuation"}
                  </SubmitButton>
                </ModalActions>
              </ValuationForm>
            </ModalCard>
          </ModalOverlay>
        )}
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
  box-sizing: border-box;
  display: flex;
  gap: 20px;
  width: 100%;
  max-width: 100%;
  min-width: 0;
  align-items: flex-start;
`;

const Container = styled.section`
  box-sizing: border-box;
  flex: 1 1 0;
  padding: 20px;
  border-radius: 24px;
  background-color: #ffffff;
  width: auto;
  max-width: 100%;
  min-width: 0;
`;

const ButtonContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const Section = styled.div`
  margin-top: 40px;
  min-width: 0;
  max-width: 100%;
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
  box-sizing: border-box;
  width: 100%;
  max-width: 100%;
  min-width: 0;
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
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
  box-sizing: border-box;
  min-width: 0;
  max-width: 100%;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
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
const EmptyState = styled.div`
  padding: 48px;
  text-align: center;
  color: #828282;
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
const ModalOverlay = styled.div`
  position: fixed;
  inset: 0;
  z-index: 10000;
  display: grid;
  place-items: center;
  padding: 24px;
  background: rgba(0, 34, 90, 0.42);
`;
const ModalCard = styled.div`
  box-sizing: border-box;
  width: min(560px, 100%);
  max-height: 90vh;
  overflow-y: auto;
  padding: 28px;
  border-radius: 20px;
  background: #fff;
`;
const ModalHeader = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 24px;
`;
const ModalTitle = styled.h2`
  margin: 0;
  color: #00225a;
  font-size: 22px;
`;
const ModalSubtitle = styled.p`
  margin: 6px 0 0;
  color: #828282;
  font-size: 14px;
`;
const CloseButton = styled.button`
  border: 0;
  background: transparent;
  color: #667085;
  font-size: 26px;
  cursor: pointer;
`;
const ValuationForm = styled.form`
  display: grid;
  gap: 18px;
`;
const Field = styled.div`
  display: grid;
  gap: 8px;
`;
const FieldLabel = styled.label`
  color: #00225a;
  font-size: 14px;
  font-weight: 500;
`;
const RequiredMark = styled.span`
  color: #d92d20;
`;
const Input = styled.input`
  box-sizing: border-box;
  width: 100%;
  border: 1px solid #d0d5dd;
  border-radius: 8px;
  padding: 11px 12px;
  background: #fff;
  color: #00225a;
  font: inherit;
  outline: none;
  &:focus {
    border-color: #007cdf;
  }
`;
const FieldHint = styled.span`
  font-size: 12px;
  color: #828282;
`;
const ModalActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 8px;
`;
const CancelButton = styled.button`
  border: 1px solid #d0d5dd;
  border-radius: 8px;
  min-height: 42px;
  padding: 0 18px;
  background: #fff;
  color: #344054;
  font: inherit;
  cursor: pointer;
`;
const SubmitButton = styled.button`
  border: 0;
  border-radius: 8px;
  min-height: 42px;
  padding: 0 18px;
  background: #007cdf;
  color: #fff;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
  &:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }
`;
