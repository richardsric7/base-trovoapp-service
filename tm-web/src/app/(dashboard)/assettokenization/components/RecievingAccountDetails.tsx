"use client";
import { TokenizationRecord } from "@/redux/api/assettokenization";
import React from "react";
import styled from "styled-components";
import AssestCard from "./AssetCard";
interface AccountDetailsProps {
  asset?: TokenizationRecord;
}
const RecievingAccountDetails: React.FC<AccountDetailsProps> = ({ asset }) => {
  const assetData = [
    {
      title: "Bank",
      subTitle: asset?.bankInfo?.bankName || "N/A",
    },
    {
      title: "Account Name",
      subTitle: asset?.beneficiaryName?.trim() || "N/A",
    },
    {
      title: "Account Number",
      subTitle: asset?.accountNumber || "N/A",
    },
  ];

  return (
    <>
      <Heading>Receiving Account Details</Heading>
      <Wrapper>
        {assetData.map((data, index) => (
          <AssestCard key={index} label={data.title} value={data.subTitle} />
        ))}
      </Wrapper>
    </>
  );
};

export default RecievingAccountDetails;
const Wrapper = styled.section`
  padding: 20px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
  padding: 16px 0;
`;
