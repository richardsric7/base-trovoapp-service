"use client";

import Link from "next/link";
import React from "react";
import { FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";
import SalesTable from "./components/SalesTable";
import Image from "next/image";
import salesTagIcon from "@/assets/images/hugeicons_sale-tag-01.svg";
import { Progress } from "antd";
import { useRouter } from "next/navigation";

const cardDetails = [
  {
    title: "Token for sale",
    value: 4000,
  },
  {
    title: "Price per Token",
    value: "10 CNGN",
  },
  {
    title: "Total Tokens Bought",
    value: "1,724 AE1",
  },
  {
    title: "Total Buyers",
    value: 9000,
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
      <Title>Primary Sales Details</Title>

      <CardContainer>
        {cardDetails.map((item) => (
          <InfoCard key={item.title}>
            <CardLabel>{item.title}</CardLabel>
            <CardValue>{item.value}</CardValue>
          </InfoCard>
        ))}
      </CardContainer>

      <SalesSection>
        <SalesTable />
        <Sidebar>
          <StatusCard>
            <StatusIcon>
              <Image
                src={salesTagIcon}
                alt="sales-icon-tag"
                height={20}
                width={20}
              />
            </StatusIcon>
            <div>
              <StatusTitle>Primary sale ongoing</StatusTitle>
              <StatusSubtitle>Ending in 3 days</StatusSubtitle>
            </div>
          </StatusCard>

          <SummaryCard>
            <SummaryTitle>Amount Raised</SummaryTitle>
            <SummaryRow>
              <SquareFilled />
              <SummaryText>Amount raised</SummaryText>
            </SummaryRow>
            <SummaryRow>
              <SquareEmpty />
              <SummaryText>Amount to be raised</SummaryText>
            </SummaryRow>
            <div>
              <ProgressLabel>
                <BoldSpan>35,000 CNGN</BoldSpan> raised out of{" "}
                <BoldSpan>80,000 CNGN</BoldSpan>
              </ProgressLabel>
              <Progress percent={50} showInfo={false} strokeColor="#00a85a" />
              <SummaryText>45,000 CNGN left</SummaryText>
            </div>
          </SummaryCard>

          <SummaryCard>
            <SummaryTitle>Token Sold</SummaryTitle>
            <SummaryRow>
              <SquareFilledBlue />
              <SummaryText>Token remaining for sale</SummaryText>
            </SummaryRow>
            <SummaryRow>
              <SquareEmpty />
              <SummaryText>Token sold out</SummaryText>
            </SummaryRow>
            <div>
              <ProgressLabel>
                <BoldSpan>20 AE1</BoldSpan> sold out of{" "}
                <BoldSpan>120 AE1</BoldSpan>
              </ProgressLabel>
              <Progress percent={50} showInfo={false} />
              <SummaryText>100 AE1 remaining for sale</SummaryText>
            </div>
          </SummaryCard>
        </Sidebar>
      </SalesSection>
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
