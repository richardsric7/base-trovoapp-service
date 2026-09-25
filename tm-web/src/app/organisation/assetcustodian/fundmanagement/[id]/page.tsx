"use client";

import PrimaryButton from "@/components/PrimaryButton";
import { useGetCustodianFundReleaseQuery } from "@/redux/api/org";
import { Alert, Button, Spin } from "antd";
import styled from "styled-components";
import SupportingDocuments from "../_components/SupportingDocuments";
import RequestInfo from "../_components/RequestInfo";
import StepperItems from "../_components/StepperItems";
import { useState } from "react";
import ExecutionModal from "../_components/ExecutionModal";

import { FaArrowLeft } from "react-icons/fa6";
import { useRouter, useParams } from "next/navigation";
import AssetCard from "@/app/(dashboard)/assettokenization/components/AssetCard";

const FundManagementDetailsPage = () => {
  const router = useRouter();
  const params = useParams();
  const id = params?.id as string;

  const { data, isLoading, isError, refetch } = useGetCustodianFundReleaseQuery(id, { skip: !id });
  const record = data?.data;

  const [isAuthorizeOpen, setIsAuthorizeOpen] = useState(false);


  const status = record?.status?.trim().toLowerCase() ?? "";
  const canReview = status === "execution_pending" || status === "processing";
  const currentStep = (() => {
    switch (status) {
      case "execution_pending":
      case "trustee_approved":
        return 1;
      case "processing":
        return 3;
      case "completed":
        return 4;
      case "failed":
        return 3;
      default:
        return 0;
    }
  })();

  const steps = [
    {
      action: "Submitted",
      date: record?.created_at
        ? new Date(record.created_at).toLocaleString()
        : "",
    },
    {
      action: "Request Authorized",
      date: record?.reviewed_at
        ? new Date(record.reviewed_at).toLocaleString()
        : "",
    },
    { action: "Funds Released", date: "" },
    { action: "Processing", date: "" },
    { action: "Completed", date: "" },
    { action: "Confirmed", date: "" },
  ];

  const receivingAccount = [
    { label: "Bank", value: record?.receiving_account?.bank || "—" },
    {
      label: "Account Name",
      value: record?.receiving_account?.account_name || "—",
    },
    {
      label: "Account Number",
      value: record?.receiving_account?.account_number || "—",
    },
  ];

  if (isLoading) {
    return (
      <LoadingWrapper>
        <Spin size="large" />
      </LoadingWrapper>
    );
  }

  if (isError || !record) {
    return <Alert type="error" message="Unable to load fund release request." action={<Button onClick={() => refetch()}>Retry</Button>} />;
  }

  return (
    <>
      <BackButton onClick={() => router.back()}>
        {" "}
        <FaArrowLeft color="#00225a" size={22} />
      </BackButton>

      <Container>
        <Sidebar>
          {steps.map((step, index) => (
            <StepperItems
              key={index}
              step={step}
              index={index}
              currentStep={currentStep}
              stepsLength={steps.length}
            />
          ))}
        </Sidebar>

        <Content>
          <Section>
            <HeaderRow>
              <div>
                <Title>Fund Release Request</Title>
              </div>

              <ButtonContainer>
                {canReview && (
                  <>
                    <PrimaryButton
                      buttonStyle={{ width: "170px" }}
                      onClick={() => setIsAuthorizeOpen(true)}
                    >
                      {status === "processing" ? "Update Status" : "Release Funds"}
                    </PrimaryButton>

                  </>
                )}
              </ButtonContainer>
            </HeaderRow>

            <RequestInfo record={record} />
          </Section>

          {/* Receiving Account Details */}
          <SectionOne>
            <Heading2>Receiving Account Details</Heading2>
            <Wrapper>
              {receivingAccount.map((item, index) => (
                <AssetCard key={index} label={item.label} value={item.value} />
              ))}
            </Wrapper>
          </SectionOne>

          {/* Supporting Documents */}
          <SectionOne>
            <SupportingDocuments
              documents={record?.supporting_documents ?? []}
            />
          </SectionOne>
        </Content>

        {canReview && (
          <>
            <ExecutionModal
              open={isAuthorizeOpen}
              onClose={() => setIsAuthorizeOpen(false)}
              record={record}
            />
          </>
        )}
      </Container>
    </>
  );
};

export default FundManagementDetailsPage;

const LoadingWrapper = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
`;

const BackButton = styled.button`
  border: 0px;
  padding: 4px 20px;
  background-color: transparent;
  cursor: pointer;
`;

const Container = styled.div`
  display: flex;
  gap: 20px;
  padding: 20px;
  width: 100%;
  min-width: 60%;
  min-height: 100vh;
`;

const Content = styled.div`
  flex: 1;
  background: #fff;
  padding: 32px;
  border-radius: 24px;
`;

const Sidebar = styled.div`
  width: 220px;
  background: white;
  padding: 24px;
  flex-shrink: 0;
  height: fit-content;
  border-radius: 16px;
`;

const Wrapper = styled.section`
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
`;

const Heading2 = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin-bottom: 16px;
`;

const Title = styled.h3`
  font-weight: 600;
  font-size: 24px;
  color: #00225a;
`;

const Section = styled.div``;

const SectionOne = styled.div`
  margin-top: 40px;
`;

const ButtonContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
`;

const RejectButton = styled.button`
  width: 150px;
  height: 48px;
  margin-top: 20px;
  border: 1px solid #be3800;
  border-radius: 8px;
  background: #fff;
  color: #be3800;
  font: inherit;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const HeaderRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
`;
