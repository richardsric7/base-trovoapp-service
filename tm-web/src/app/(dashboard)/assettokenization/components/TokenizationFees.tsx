"use client";

import React from "react";
import styled from "styled-components";
import AssetCard from "./AssetCard";
import { TokenizationRecord } from "@/redux/api/assettokenization";

interface TokenizationFeesProps {
  asset?: TokenizationRecord;
}

const TokenizationFees: React.FC<TokenizationFeesProps> = ({ asset }) => {
  const assetData = [
    {
      id: 1,
      title: "Asset Tokenization Fee",
      subTitle: asset?.feeInFiat
        ? `${asset.feeInFiat.toLocaleString()} NGN`
        : "N/A",
    },
    {
      id: 2,
      title: "SEC Regulatory Fee",
      subTitle:
        asset?.SECTokenizationFeeValue !== undefined
          ? `${asset.SECTokenizationFeeValue.toLocaleString()} NGN`
          : "N/A",
    },
    {
      id: 3,
      title: "Asset Custody Fee",
      subTitle: asset?.custodianFeeValue
        ? `${asset.custodianFeeValue.toLocaleString()} NGN`
        : "N/A",
    },
    {
      id: 4,
      title: "Asset Management Fee",
      subTitle: asset?.assetManagerFeeValue
        ? `${asset.assetManagerFeeValue.toLocaleString()} NGN`
        : "N/A",
    },
    {
      id: 5,
      title: "Issuing House Fee",
      subTitle: asset?.issuingHouseFeeValue
        ? `${asset.issuingHouseFeeValue.toLocaleString()} NGN`
        : "N/A",
    },

    {
      id: 6,
      title: "Legal Fee",
      subTitle: asset?.legalAndProfessionalFeeValue
        ? `${asset.legalAndProfessionalFeeValue.toLocaleString()} NGN`
        : "N/A",
    },

    {
      id: 7,
      title: "Rating Agency Fee",
      subTitle: asset?.ratingAgencyFeeValue
        ? `${asset.ratingAgencyFeeValue.toLocaleString()} NGN`
        : "N/A",
    },

    {
      id: 8,
      title: "VAT",
      subTitle: asset?.vatValue
        ? `${asset.vatValue.toLocaleString()} NGN`
        : "N/A",
    },

    {
      id: 9,
      title: "Trustee Fee",
      subTitle: asset?.trusteeFeeValue
        ? `${asset.trusteeFeeValue.toLocaleString()} NGN`
        : "N/A",
    },
  ];

  return (
    <>
      <Heading>Tokenization Fees In Fiat</Heading>

      <Wrapper>
        {assetData.map((data) => (
          <AssetCard key={data.id} label={data.title} value={data.subTitle} />
        ))}
      </Wrapper>
    </>
  );
};

export default TokenizationFees;

const Wrapper = styled.section`
  padding: 24px;
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
  margin: 16px 0px;
`;
