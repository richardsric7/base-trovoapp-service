"use client";

import React, { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";
import PrimaryButton from "@/components/PrimaryButton";
import Wallets from "../components/Wallets";

import { AiOutlineCloseCircle } from "react-icons/ai";
import ExemptedCountries from "../components/ExemptedCountries";
import { showErrorToast, showSuccessToast } from "@/components";

import AssetValue from "../components/AssetValue";
import AssetToken from "../components/AssetToken";
import AssetVerificationDocuments from "../components/AssetVerificationDocuments";

import DisapproveModal from "../components/DisapproveModal";
import Payment from "../components/Payment";
import SuccessMessage from "@/components/SuccessMessage";
import InitiateMintingModal from "../components/InitiateMintingModal";
import {
  useApproveTokenizationAndMintMutation,
  useCompleteVettingMutation,
  useConfirmPaymentMutation,
  useDeleteTokenizationMutation,
  useFailDueDiligenceMutation,
  useGetTokenizationDetailQuery,
  useUpdateTokenizationInfoMutation,
} from "@/redux/api/assettokenization";
import AssetInformation from "../components/AssetInformation";
import { FiClock } from "react-icons/fi";
import RecievingAccountDetails from "../components/RecievingAccountDetails";
import StepItem from "../components/StepItem";
import AssetProtection from "../components/AssetProtection";
import TokenizationFees from "../components/TokenizationFees";
import AsignStakeHolders from "../components/AsignStakeHolders";
import AssignedStakeHolder from "../components/AssignedStakeHolder";
import ApplicationPayment from "../components/ApplicationPayment";
import AssetTokenizer from "../components/AssetTokenizer";
import AssetOwnership from "../components/AssetOwnership";
import ProceedPayoutModal from "../components/ProceedPayoutModal";
import AssetImages from "../components/AssetImages";
import FeeInAsset from "../components/FeeInAsset";
import { BsThreeDots } from "react-icons/bs";

export interface StepProps {
  action: string;
  date: string;
  subset?: string[];
}
export default function AssetTokenizationDetailsPage() {
  const [success, setSuccess] = useState<boolean>(false);
  const [onSuccess, setOnSuccess] = useState<boolean>(false);
  const [currentStep, setCurrentStep] = useState(0);
  const [openModal, setOpenModal] = useState<boolean>(false);
  const [openProceeds, setOpenProceeds] = useState(false);
  const [isDueDiligenceDisapproved, setIsDueDiligenceDisapproved] =
    useState(false);
  const [reasonInput, setReasonInput] = useState<string>("");
  const [disapproval, setDisapproval] = useState<boolean>(false);
  const [processingStep, setProcessingStep] = useState<number | null>(null);
  const [isInitiateMinting, setIsInitiateMinting] = useState<boolean>(false);
  const [showAssignStakeholdersModal, setShowAssignStakeholdersModal] =
    useState<boolean>(false);
  const [stakeholdersAssigned, setStakeholdersAssigned] =
    useState<boolean>(false);
  const [isOpenbox, setIsOpenBox] = useState<boolean>(false);

  const router = useRouter();
  const [assignedManagerId, setAssignedManagerId] = useState<
    number | undefined
  >();
  const [assignedCustodianId, setAssignedCustodianId] = useState<
    number | undefined
  >();
  const [assignedIssuingHouseId, setAssignedIssuingHouseId] = useState<
    number | undefined
  >();

  const [ratinAgencyId, setRatingAgencyId] = useState<number | undefined>();

  const [legalAndProfesionalPartnerId, setLegalAndProfesionalPartnerId] =
    useState<number | undefined>();
  const [trusteeId, setTrusteeId] = useState<number | undefined>();

  const [steps, setSteps] = useState<StepProps[]>([
    { action: "Submitted", date: "" },
    { action: "Vetting Completed", date: "" },
    { action: "Payment Made", date: "" },
    { action: "Payment Confirmed", date: "" },
    {
      action: "Processing",
      date: "",
      subset: ["Sec Approval", "Due Diligence"],
    },
    { action: "Tokenization Approved", date: "" },
    { action: "Minted", date: "" },
    { action: "Primary Sale Started", date: "" },
    { action: "Primary Sale Ended", date: "" },
    { action: "Secondary Market", date: "" },
    { action: "Liquidated", date: "" },
  ]);
  const [shouldPoll, setShouldPoll] = useState(false);

  const params = useParams();
  const assetName = params.id as string;

  const {
    data: tokenizationData,
    error,
    isLoading,
    refetch,
  } = useGetTokenizationDetailQuery(assetName, {
    pollingInterval: shouldPoll ? 5000 : 0,
    // pollingInterval: pollingInterval,
  });
  const asset = tokenizationData;

  const [completeVetting] = useCompleteVettingMutation();
  const [confirmPayment] = useConfirmPaymentMutation();
  const [approveTokenization] = useApproveTokenizationAndMintMutation();
  const [failDueDiligence] = useFailDueDiligenceMutation();
  const [updateTokenizationInfo] = useUpdateTokenizationInfoMutation();
  const [deleteTokenization, { isLoading: isDeleting }] =
    useDeleteTokenizationMutation();

  // Poll only when awaiting payment
  useEffect(() => {
    if (asset) {
      const needsPolling =
        asset.vettingStatus === 1 && asset.assetTokenizationStatus === 1;
      setShouldPoll(needsPolling);
    }
  }, [asset]);

  // Update steps based on asset status
  useEffect(() => {
    if (asset && !isDueDiligenceDisapproved) {
      let newStep = 0;
      //       // If vetting is done (and DB status is >= 1), we’re at least “Awaiting Payment”
      if (asset.vettingStatus === 1) {
        if (asset.assetTokenizationStatus === 1) newStep = 1;
        // If DB says 2 => "Awaiting Fee Payment Confirmation"
        if (asset.assetTokenizationStatus === 2) newStep = 2;
        // "Payment Confirmed / Awaiting DD"
        if (asset.assetTokenizationStatus === 3) {
          newStep = 4;
          setProcessingStep(4);
        }
        // If DB says 4 => "Approval / Minting"
        if (asset.assetTokenizationStatus === 4) newStep = 5;
        // If DB says 5 => "Primary Sales"
        if (asset.assetTokenizationStatus === 5) newStep = 7;
        // If DB says 6 => "Secondary Sales"
        if (asset.assetTokenizationStatus === 6) newStep = 9;
        // If DB says 7 => "Liquidated"
        if (asset.assetTokenizationStatus === 7) newStep = 10;
        // If DB says 8 => "Refunded", handle as needed
        // Possibly reset to 0 or show some "Refunded" step
      }

      setCurrentStep(newStep);
      if (newStep === 3) setProcessingStep(4);
    }
  }, [asset, isDueDiligenceDisapproved]);

  // Update step dates when currentStep changes
  useEffect(() => {
    if (currentStep >= 0) {
      const newDate = new Date().toLocaleDateString("en-US", {
        day: "numeric",
        month: "short",
        year: "numeric",
      });
      setSteps((prevSteps) =>
        prevSteps.map((step, index) =>
          index <= currentStep && !step.date ? { ...step, date: newDate } : step
        )
      );
    }
  }, [currentStep]);

  useEffect(() => {
    if (asset) {
      const hasAssignedStakeholders =
        !!asset.approvedAssetCustodianId &&
        !!asset.assetManagerId &&
        !!asset.assetIssuingHouseId;
      // Add more checks here if needed

      setStakeholdersAssigned(hasAssignedStakeholders);
    }
  }, [asset]);

  // handle update paid request
  const handleSave = async () => {
    if (!asset) return;
    const { id, createdAt, updatedAt, ...data } = asset;
    try {
      await updateTokenizationInfo({
        tokenizedAssetID: asset.id,
        data,
      }).unwrap();
    } catch (error: any) {
      const errorString =
        error?.response?.data?.message ||
        error?.data?.message ||
        error?.message ||
        error?.response?.data?.error ||
        error?.data?.error ||
        error?.error ||
        "An error occurred while updating.";

      const messageMatch = errorString.match(/message:(.*?)(\]$|$)/);
      let finalErrorMessage = errorString;
      if (messageMatch && messageMatch[1]) {
        finalErrorMessage = messageMatch[1].trim();
      }

      showErrorToast(finalErrorMessage);
    }
  };

  // Memoized event handlers
  const handleNextStepAfterSubmit = useCallback((newStep: number) => {
    if (newStep === undefined) {
      console.error("No new step provided.");
      return;
    }
    setCurrentStep(newStep);
    const newDate = new Date().toLocaleDateString("en-US", {
      day: "numeric",
      month: "short",
      year: "numeric",
    });
    setSteps((prevSteps) =>
      prevSteps.map((step, index) =>
        index === newStep && !step.date ? { ...step, date: newDate } : step
      )
    );
  }, []);

  const handleVettingClick = useCallback(
    async ({
      assetManagerId,
      approvedAssetCustodianId,
      assetIssuingHouseId,
      ratingAgencyId,
      legalAndProfesionalPartnerId,
      trusteeId,
      legalAdviserId,
      financialAdviserId,
    }: {
      assetManagerId?: number;
      approvedAssetCustodianId?: number;
      assetIssuingHouseId?: number;
      ratingAgencyId?: number;
      legalAndProfesionalPartnerId?: number;
      trusteeId?: number;
      legalAdviserId?: number;
      financialAdviserId?: number;
    }) => {
      if (!asset) return;
      // ✅ Runtime safety check
      if (
        assetManagerId === undefined ||
        approvedAssetCustodianId === undefined ||
        assetIssuingHouseId === undefined ||
        ratingAgencyId === undefined ||
        legalAndProfesionalPartnerId === undefined ||
        trusteeId === undefined
      ) {
        showErrorToast(
          "All stakeholders must be assigned before completing vetting."
        );
        return;
      }
      try {
        await completeVetting({
          tokenizedAssetID: asset.id,
          vettingPayload: {
            CountryCode: asset?.assetCountryLocation,
            approvedAssetCustodianId,
            assetManagerId,
            assetIssuingHouseId,
            ratingAgencyId,
            legalAndProfesionalPartnerId,
            trusteeId,
            ...(legalAdviserId !== undefined && { legalAdviserId }),
            ...(financialAdviserId !== undefined && { financialAdviserId }),
            assetQuoteCurrency: asset.assetQuoteCurrency,
            proceedPayoutCurrency: asset.proceedPayoutCurrency,
          },
        }).unwrap();
        handleNextStepAfterSubmit(1);
        await refetch();
      } catch (error: any) {
        const errorString =
          error?.response?.data?.message ||
          error?.data?.message ||
          error?.message ||
          error?.response?.data?.error ||
          error?.data?.error ||
          error?.error ||
          "An error occurred while updating.";

        const messageMatch = errorString.match(/message:(.*?)(\]$|$)/);
        let finalErrorMessage = errorString;
        if (messageMatch && messageMatch[1]) {
          finalErrorMessage = messageMatch[1].trim();
        }

        showErrorToast(finalErrorMessage);
      }
    },
    [asset, handleNextStepAfterSubmit, refetch]
  );

  const handlePaymentConfirm = useCallback(async () => {
    if (!asset) return;
    try {
      await confirmPayment(asset.id).unwrap();
      await refetch();
      setProcessingStep(4);
      handleNextStepAfterSubmit(4);

      // After confirming payment, update the tokenization info with current asset values.
      await handleSave();
    } catch (error: any) {
      const errorString =
        error?.response?.data?.message ||
        error?.data?.message ||
        error?.message ||
        error?.response?.data?.error ||
        error?.data?.error ||
        error?.error ||
        "An error occurred while updating.";

      const messageMatch = errorString.match(/message:(.*?)(\]$|$)/);
      let finalErrorMessage = errorString;
      if (messageMatch && messageMatch[1]) {
        finalErrorMessage = messageMatch[1].trim();
      }

      showErrorToast(finalErrorMessage);
    }
  }, [asset, handleNextStepAfterSubmit, refetch, handleSave]);

  const handleApproveDueDiligence = useCallback(async () => {
    if (!asset) return;
    const tokenizedAssetID = asset?.id;
    try {
      await approveTokenization(tokenizedAssetID).unwrap();
      await refetch();
      handleNextStepAfterSubmit(5);
    } catch (error: any) {
      const errorString =
        error?.response?.data?.message ||
        error?.data?.message ||
        error?.message ||
        error?.response?.data?.error ||
        error?.data?.error ||
        error?.error ||
        "An error occurred while updating.";

      const messageMatch = errorString.match(/message:(.*?)(\]$|$)/);
      let finalErrorMessage = errorString;
      if (messageMatch && messageMatch[1]) {
        finalErrorMessage = messageMatch[1].trim();
      }

      showErrorToast(finalErrorMessage);
    }
  }, [asset, handleNextStepAfterSubmit, refetch]);

  const handleDueDiligenceDisapprove = useCallback(async () => {
    if (!asset) return;
    if (!reasonInput.trim()) {
      showErrorToast("Please provide a reason for failing due diligence.");
      return;
    }
    try {
      await failDueDiligence({
        tokenizedAssetID: asset.id,
        reason: reasonInput,
      }).unwrap();
      const newDate = new Date().toLocaleDateString("en-US", {
        day: "numeric",
        month: "short",
        year: "numeric",
      });
      setSteps((prevSteps) => {
        const updatedSteps = [...prevSteps];
        if (updatedSteps[5]) {
          updatedSteps[5] = {
            ...updatedSteps[5],
            action: "Tokenization Failed",
            date: newDate,
          };
        }
        return updatedSteps;
      });
      setReasonInput("");
      setIsDueDiligenceDisapproved(true);
      setDisapproval(false);
      setProcessingStep(4);
    } catch (error: any) {
      const errorString =
        error?.response?.data?.message ||
        error?.data?.message ||
        error?.message ||
        error?.response?.data?.error ||
        error?.data?.error ||
        error?.error ||
        "An error occurred while updating.";

      const messageMatch = errorString.match(/message:(.*?)(\]$|$)/);
      let finalErrorMessage = errorString;
      if (messageMatch && messageMatch[1]) {
        finalErrorMessage = messageMatch[1].trim();
      }

      showErrorToast(finalErrorMessage);
    }
  }, [asset, reasonInput]);

  const handleRetryDueDiligence = useCallback(() => {
    setReasonInput("");
    setIsDueDiligenceDisapproved(false);
    setDisapproval(false);
    setProcessingStep(null);
    setCurrentStep(0);
    setSteps([
      { action: "Submitted", date: "" },
      { action: "Vetting Completed", date: "" },
      { action: "Payment Made", date: "" },
      { action: "Payment Confirmed", date: "" },
      {
        action: "Processing",
        date: "",
        subset: ["Sec Approval", "Due Diligence"],
      },
      { action: "Tokenization Approved", date: "" },
      { action: "Minted", date: "" },
      { action: "Primary Sale Started", date: "" },
      { action: "Primary Sale Ended", date: "" },
      { action: "Secondary Sales", date: "" },
      { action: "Liquidated", date: "" },
    ]);
    refetch();
  }, [refetch]);

  const toggleModal = useCallback((type: string) => {
    if (type === "completeVetting") {
      setShowAssignStakeholdersModal((prev) => !prev);
    }
    if (type === "dueDiligence") {
      setOpenModal((prev) => !prev);
    } else if (type === "disapproval") {
      setDisapproval((prev) => !prev);
    } else if (type === "initiateMinting") {
      setIsInitiateMinting((prev) => !prev);
    }
  }, []);

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading details</div>;

  return (
    <TokenizationStep>
      <StepWrapper>
        {steps.map((step, index) => (
          <StepItem
            key={index}
            step={step}
            index={index}
            currentStep={currentStep}
            processingStep={processingStep}
            isDueDiligenceDisapproved={isDueDiligenceDisapproved}
          />
        ))}
      </StepWrapper>
      <Container>
        <Header>
          <StyledLink href="/assettokenization">
            <FaArrowLeft />
          </StyledLink>

          <ButtonSection>
            {!isDueDiligenceDisapproved && asset?.vettingStatus !== 1 && (
              <PrimaryButton
                onClick={() => {
                  setShowAssignStakeholdersModal(true);
                }}
              >
                Complete Vetting
              </PrimaryButton>
            )}
            {asset?.assetTokenizationStatus === 2 && (
              <PrimaryButton onClick={handlePaymentConfirm}>
                Confirm Payment
              </PrimaryButton>
            )}
            {currentStep === 4 &&
              (isDueDiligenceDisapproved ? (
                <PrimaryButton onClick={handleRetryDueDiligence}>
                  Retry Due Diligence
                </PrimaryButton>
              ) : (
                <ButtonContainer>
                  <ApproveButton onClick={() => toggleModal("initiateMinting")}>
                    Approve Tokenization
                  </ApproveButton>

                  <DisapproveButton onClick={() => toggleModal("disapproval")}>
                    <AiOutlineCloseCircle color="#BE3800" />
                    Reject Tokenization
                  </DisapproveButton>
                </ButtonContainer>
              ))}
            {currentStep === 5 && (
              <ButtonContainer>
                <ApproveButton>Manage Schedule</ApproveButton>
              </ButtonContainer>
            )}

            {currentStep === 7 && (
              <ButtonContainer>
                <ApproveButton>Manage Schedule</ApproveButton>
                <SecondaryButton
                  onClick={() => router.push("/primarysalesdetails")}
                >
                  Primary sales Details
                </SecondaryButton>
              </ButtonContainer>
            )}

            {currentStep === 9 && (
              <ButtonContainer>
                <ApproveButton onClick={() => router.push("/dividendandyield")}>
                  Dividened & yield
                </ApproveButton>

                <IconBox onClick={() => setIsOpenBox(!isOpenbox)}>
                  <BsThreeDots color="#007CDF" />
                </IconBox>
              </ButtonContainer>
            )}
          </ButtonSection>
        </Header>

        {isOpenbox && (
          <OtherPages>
            <StyledLinks href="/earlyexit">
              {" "}
              Early Exit <StyledNote>2</StyledNote>
            </StyledLinks>

            <StyledLinks href="/assetincomereport">
              {" "}
              Asset Income Report
            </StyledLinks>
            <StyledLinks href="/tokenholders"> View Token Holders</StyledLinks>
          </OtherPages>
        )}
        {disapproval && (
          <DisapproveModal
            asset={asset}
            reasonInput={reasonInput}
            setReasonInput={setReasonInput}
            disapproval={disapproval}
            setDisapproval={setDisapproval}
            onSubmit={() => {
              handleDueDiligenceDisapprove();
              setSuccess(true);
            }}
          />
        )}

        {isInitiateMinting && (
          <InitiateMintingModal
            asset={asset}
            isInitiateMinting={isInitiateMinting}
            setIsInitiateMinting={setIsInitiateMinting}
            onSubmit={handleApproveDueDiligence}
            refetchAsset={refetch}
          />
        )}

        {success && (
          <SuccessMessage
            isOpen={success}
            setIsOpen={(val) => {
              setSuccess(val);
              if (!val) {
                router.push("/assettokenization");
              }
            }}
            heading="Success!"
            message="You have disapproved the due diligence for the asset"
            email={asset?.assetName}
          />
        )}

        {openProceeds && (
          <ProceedPayoutModal
            isOpen={openProceeds}
            setOpenProceeds={setOpenProceeds}
            setOnSuccess={setOnSuccess}
          />
        )}

        {onSuccess && (
          <SuccessMessage
            isOpen={onSuccess}
            setIsOpen={(val) => {
              setOnSuccess(val);
              if (!val) {
              }
            }}
            heading="Success!"
            message="Your Payout proceed Requests for Atlantis Asset have been submitted and are awaiting approval. You’ll be notified when all requests have been approved."
            email=""
          />
        )}

        {showAssignStakeholdersModal && (
          <AsignStakeHolders
            asset={asset}
            isOpen={showAssignStakeholdersModal}
            onClose={() => setShowAssignStakeholdersModal(false)}
            onAssignSuccess={(ids) => {
              const {
                assetManagerId,
                approvedAssetCustodianId,
                assetIssuingHouseId,
                ratingAgencyId,
                legalAndProfesionalPartnerId,
                trusteeId,
                legalAdviserId,
                financialAdviserId,
              } = ids;

              // Using the IDs from the callback directly, not from state
              console.log("Received stakeholder IDs:", ids);

              // Close modal first
              setShowAssignStakeholdersModal(false);

              // Update local state
              setAssignedManagerId(assetManagerId);
              setAssignedCustodianId(approvedAssetCustodianId);
              setAssignedIssuingHouseId(assetIssuingHouseId);
              setRatingAgencyId(ratingAgencyId);
              setLegalAndProfesionalPartnerId(legalAndProfesionalPartnerId);
              setTrusteeId(trusteeId);
              setStakeholdersAssigned(true);

              showSuccessToast(
                "Tokenization stakeholders assigned successfully."
              );

              //  Call handleVettingClick with the IDs directly from the callback
              handleVettingClick({
                assetManagerId,
                approvedAssetCustodianId,
                assetIssuingHouseId,
                ratingAgencyId,
                legalAndProfesionalPartnerId,
                trusteeId,
                legalAdviserId,
                financialAdviserId,
              }).then(() => {
                refetch();
              });
            }}
          />
        )}

        <AssetInformation asset={asset} />
        <AssetTokenizer asset={asset} />
        <AssetOwnership asset={asset} />
        <Wallets asset={asset} />
        <AssignedStakeHolder asset={asset} />
        <AssetValue asset={asset} />

        <AssetImages />
        <AssetProtection asset={asset} />
        <AssetToken asset={asset} />
        <RecievingAccountDetails asset={asset} />
        <ExemptedCountries asset={asset} />
        <AssetVerificationDocuments asset={asset} />
        <ApplicationPayment asset={asset} />
        <TokenizationFees asset={asset} />
        <FeeInAsset asset={asset} />
      </Container>
    </TokenizationStep>
  );
}

const Container = styled.section`
  padding: 32px;
  border-radius: 24px;
  background-color: #ffffff;
  // min-width: 750px;
  width: 100%;
  min-width: 70%;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  gap: 30px;
  justify-content: space-between;
`;

const TokenizationStep = styled.div`
  display: flex;
  gap: 20px;
`;

const StyledLink = styled(Link)`
  text-decoration: none;
  color: #000000;
  width: 20px;
`;

const StepWrapper = styled.div`
  display: flex;
  flex-direction: column;
  padding: 20px;
  border-radius: 24px;
  background-color: #ffffff;
  width: 230px;
  height: 100%;
`;

const ButtonSection = styled.div`
  display: flex;
  align-items: flex-start;
`;

const ButtonContainer = styled.div`
  display: flex;
  gap: 16px;
`;

const ApproveButton = styled.button`
  background-color: #007cdf;
  color: #fff;
  border: none;
  width: 208px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 14px;
  cursor: pointer;
  font-family: inherit;
`;

const DisapproveButton = styled.button`
  background-color: transparent;
  border: 1px solid #be3800;
  color: #be3800;
  width: 208px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 14px;
  cursor: pointer;
  font-family: inherit;
`;

const SecondaryButton = styled.button`
  background-color: transparent;
  border: 1px solid #007cdf;
  color: #007cdf;
  width: 208px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 14px;
  cursor: pointer;
  font-family: inherit;
`;

const IconBox = styled.div`
  width: 48px;
  height: 48px;
  border: 1px solid #007cdf;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  // position: realtive;
`;

const OtherPages = styled.div`
  position: absolute;
  background-color: #ffffff;
  box-shadow: 0px 4px 8px 0px #0000001a;
  width: 200px;
  top: 30%;
  right: 2%;
  border-radius: 16px;
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
`;
const StyledLinks = styled(Link)`
  font-weight: 400;
  color: #00225a;
  font-size: 14px;
  line-height: 28px;
  letter-spacing: 0%;
  text-decoration: none;
  padding: 2px;
  &:hover {
    padding: 6px;
    background-color: #f2f6f9;
    padding: 2px;
    border-radius: 6px;
  }
`;

const StyledNote = styled.span`
  width: 20px;
  height: 20px;
  background-color: #ed3137;
  color: #fff;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
`;
