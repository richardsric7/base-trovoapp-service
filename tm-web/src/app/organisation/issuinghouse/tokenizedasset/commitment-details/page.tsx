"use client";

import React from "react";
import { FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";

import { useRouter } from "next/navigation";
import CommitmentTable from "./CommitmentTable";

const cardDetails = [
  {
    title: "Total Purchase Commitments",
    value: "1,000,000 CNGN",
  },
  {
    title: "Total Interests Expressed",
    value: "2,500",
  },
];

const PrimarySalesDetailsPage = () => {
  const router = useRouter();
  const goBack = () => {
    router.back();
  };
  return (
    <Container>
      <BackButton onClick={goBack}>
        <FaArrowLeft />
      </BackButton>
      <Title>Commitment Details</Title>

      <CardContainer>
        {cardDetails.map((item) => (
          <InfoCard key={item.title}>
            <CardLabel>{item.title}</CardLabel>
            <CardValue>{item.value}</CardValue>
          </InfoCard>
        ))}
      </CardContainer>
      <CommitmentTable />
    </Container>
  );
};

export default PrimarySalesDetailsPage;

const Container = styled.section`
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;

const BackButton = styled.button`
  text-decoration: none;
  color: #000000;
  width: 20px;
  display: block;
  font-family: inherit;
  border: none;
  background: none;
  cursor: pointer;
`;

const CardContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  width: 100%;
`;

const InfoCard = styled.div`
  background-color: #ffffff;
  border-radius: 10px;
  padding: 8px 16px;
  width: 80%;
`;

const CardLabel = styled.p`
  font-weight: 400;
  font-size: 14px;
  line-height: 20px;
  color: #828282;
  padding-bottom: 10px;
`;

const CardValue = styled.p`
  font-weight: 600;
  font-size: 24px;
  line-height: 24px;
  color: #00225a;
`;

const SalesSection = styled.div`
  display: flex;
  margin-top: 20px;
  gap: 20px;
`;

const Sidebar = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 40%;
`;

const StatusCard = styled.div`
  background-color: #1a253b;
  display: flex;
  align-items: center;
  gap: 8px;
  border-radius: 12px;
  padding: 10px;
`;

const StatusIcon = styled.div`
  background-color: #ffffff26;
  padding: 10px;
  border-radius: 50%;
`;

const StatusTitle = styled.p`
  color: #ffffff;
  font-weight: 500;
  font-size: 16px;
  line-height: 20px;
`;

const StatusSubtitle = styled.p`
  color: #ffffff;
  font-weight: 400;
  font-size: 12px;
  line-height: 20px;
`;

const SummaryCard = styled.div`
  background-color: #ffffff;
  border-radius: 24px;
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

const SummaryTitle = styled.p`
  color: #00225a;
  font-weight: 700;
  font-size: 20px;
  line-height: 28px;
`;

const SummaryRow = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const SquareFilled = styled.div`
  background-color: #00a85a;
  width: 16px;
  height: 16px;
  border-radius: 4px;
`;

const SquareFilledBlue = styled.div`
  background-color: #007cdf;
  width: 16px;
  height: 16px;
  border-radius: 4px;
`;

const SquareEmpty = styled.div`
  background-color: #f5f5f5;
  border: 1px solid #0000001a;
  width: 16px;
  height: 16px;
  border-radius: 4px;
`;

const SummaryText = styled.p`
  font-weight: 400;
  font-size: 14px;
  line-height: 20px;
  color: #828282;
`;

const ProgressLabel = styled.p`
  font-weight: 400;
  font-size: 14px;
  line-height: 16px;
  color: #00225a;
`;

const BoldSpan = styled.span`
  font-weight: 600;
  font-size: 14px;
  line-height: 16px;
  color: #00225a;
`;
