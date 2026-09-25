"use client";
import React from "react";
import { AiOutlineCheckCircle, AiOutlineCloseCircle } from "react-icons/ai";
import styled from "styled-components";

import { useParams, useRouter } from "next/navigation";
import PayoutListTable from "../components/PayoutListTable";
import { FaArrowLeft } from "react-icons/fa6";

const DividendAndYielsDetailsPage = () => {
  const router = useRouter();
  const params = useParams();
  const id = params.id;

  const isProcessed = true;

  const dividendData = [
    { title: "Date", dataIndex: "23 Sep 2023", key: "payoutDate" },
    { title: "Token Holders", dataIndex: "1200", key: "tokenHolders" },
    { title: "Payout Per Token", dataIndex: "₦20.00", key: "payoutPerToken" },
    { title: "Gross Amount", dataIndex: "₦100,000,000.00", key: "grossAmount" },
    { title: "Tax (WHT)", dataIndex: "₦5,000,000.00", key: "tax" },
    { title: "Net Amount", dataIndex: "₦95,000,000.00", key: "netAmount" },
  ];

  const handleBack = () => {
    router.back();
  };
  return (
    <>
      <BackButtonLink onClick={handleBack}>
        <FaArrowLeft size={18} />
      </BackButtonLink>
      <HeadingContent>
        <Wrapper>
          <HeadingTitle>
            <Heading>TR123456</Heading>
            <AssetStatus
              textColor={isProcessed ? "#00A859" : "#FF4D4D"}
              backgroundColor={isProcessed ? "#00A8591A" : "#BE38001A"}
            >
              {isProcessed ? (
                <AiOutlineCheckCircle size={14} />
              ) : (
                <AiOutlineCloseCircle size={14} />
              )}
              {isProcessed ? "Paid" : "Failed"}
            </AssetStatus>
          </HeadingTitle>

          <Description>
            Lorem ipsum dolor sit amet consectetur. Bibendum nisl vivamus
            phasellus velit eleifend lacinia. Nibh elit vestibulum sed nunc
            vitae nibh suspendisse. Commodo erat sit amet duis at.
          </Description>

          <DataContainer>
            {dividendData.map((item) => (
              <DataCard key={item.key}>
                <DataTitle>{item.title}</DataTitle>
                <DataValue>{item.dataIndex}</DataValue>
              </DataCard>
            ))}
          </DataContainer>
        </Wrapper>
      </HeadingContent>

      <PayoutListTable />
    </>
  );
};

export default DividendAndYielsDetailsPage;

const BackButtonLink = styled.button`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
  background: none;
  cursor: pointer;
  border: none;
  padding: 24px;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
`;

const HeadingContent = styled.div`
  background-color: #ffffff;
  padding: 24px;
  border-radius: 20px;
`;

const HeadingTitle = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
`;

const Wrapper = styled.section`
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
`;

const Description = styled.p`
  font-weight: 400;
  font-size: 14px;
  line-height: 25.27px;
  letter-spacing: 1.9%;
  color: #828282;
  padding-bottom: 20px;
`;

const AssetStatus = styled.p<{ textColor: string; backgroundColor: string }>`
  color: ${({ textColor }) => textColor};
  background-color: ${({ backgroundColor }) => backgroundColor};
  font-weight: 500;
  padding: 8px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
`;

const DataContainer = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
`;

const DataCard = styled.div``;

const DataTitle = styled.p`
  font-weight: 400;
  font-size: 14px;
  color: #828282;
  padding-bottom: 4px;
`;

const DataValue = styled.p`
  font-weight: 500;
  font-size: 14px;
  color: #00225a;
`;
