"use client";

import SecondaryButton from "@/components/SecondaryButton";
import { useGetAssetManagerFundReleaseQuery } from "@/redux/api/assetManager";
import { Spin } from "antd";
import styled from "styled-components";
import SupportingDocuments from "@/app/organisation/trustee/fundmanagement/_components/SupportingDocuments";
import RequestInfo from "@/app/organisation/trustee/fundmanagement/_components/RequestInfo";
import StepperItems from "@/app/organisation/trustee/fundmanagement/_components/StepperItems";
import { FaArrowLeft } from "react-icons/fa6";
import { useRouter, useParams } from "next/navigation";
import AssetCard from "@/app/(dashboard)/assettokenization/components/AssetCard";

const FundManagementDetailsPage = () => {
  const router = useRouter();
  const params = useParams();
  const id = params?.id as string;

  const { data, isLoading, isError, refetch } =
    useGetAssetManagerFundReleaseQuery(id, { skip: !id });
  const record = data?.data;

  const status = record?.status?.trim().toLowerCase() ?? "";
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
    { label: "Bank", value: record?.receiving_account?.bank || "_" },
    {
      label: "Account Name",
      value: record?.receiving_account?.account_name || "_",
    },
    {
      label: "Account Number",
      value: record?.receiving_account?.account_number || "_",
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
    return (
      <LoadingWrapper>
        <div>
          <p>Unable to load this fund release request.</p>
          <SecondaryButton onClick={() => refetch()}>Try again</SecondaryButton>
          <SecondaryButton
            onClick={() =>
              router.push("/organisation/assetmanager/fundmanagement")
            }
          >
            Back to requests
          </SecondaryButton>
        </div>
      </LoadingWrapper>
    );
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

const HeaderRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
`;
