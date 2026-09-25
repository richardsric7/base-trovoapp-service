"use client";

import React, { useState } from "react";
import styled from "styled-components";
import AssetCard from "./AssetCard";
import { TokenizationRecord } from "@/redux/api/assettokenization";

interface ApplicationPaymentProps {
  asset?: TokenizationRecord;
}

const ApplicationPayment: React.FC<ApplicationPaymentProps> = ({ asset }) => {
  const assetData = [
    {
      id: 1,
      title: "Payment Address",
      subTitle: "0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
    },
    {
      id: 2,
      title: "Reciveing Address",
      subTitle: "0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
    },
    {
      id: 3,
      title: "Amount Paid",
      subTitle: asset?.tokenizationApplicationFee
        ? ` ${asset.tokenizationApplicationFee.toFixed(2)} TROV`
        : "Amount Paid: —",
    },
    {
      id: 4,
      title: "Blockchain Proof",
      subTitle: "https//gjdjnjfujksbjfckvnvf",
    },
  ];

  return (
    <>
      <Heading>Application Fee Payment</Heading>

      <Wrapper>
        {assetData.map((data) => (
          <AssetCard
            key={data.id}
            label={data.title}
            value={
              data.title === "Blockchain Proof" ? (
                <BlockProof>{data.subTitle}</BlockProof>
              ) : (
                data.subTitle
              )
            }
          />
        ))}
      </Wrapper>
    </>
  );
};

export default ApplicationPayment;

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

const BlockProof = styled.p`
  color: #007cdf;
`;
