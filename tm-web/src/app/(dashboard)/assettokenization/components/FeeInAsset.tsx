"use client";

import React from "react";
import styled from "styled-components";
import AssetCard from "./AssetCard";
import { TokenizationRecord } from "@/redux/api/assettokenization";

interface FeeInAssetProps {
  asset?: TokenizationRecord;
}

interface AssetData {
  id: number;
  title: string;
  subTitle: any;
}

const FeeInAsset: React.FC<FeeInAssetProps> = ({ asset }) => {
  const assetData: AssetData[] = [
    {
      id: 1,
      title: "Asset Tokenization Fee",
      subTitle: asset?.feeInAsset
        ? `${asset.feeInAsset.toLocaleString()} ${asset?.assetCode}`
        : "N/A",
    },

    {
      id: 2,
      title: "VAT",
      subTitle: asset?.vatInAsset
        ? `${asset.vatInAsset.toLocaleString()} ${asset?.assetCode}`
        : "N/A",
    },
  ];

  return (
    <>
      <HeadingContent>
        <Heading>Tokenization Fees In Asset</Heading>

        {/* <SecondaryButton
          onClick={() => {
            if (asset) {
              setShowEditModal(true);
            }
          }}
          buttonStyle={{
            width: "100px",
            display: "flex",
            alignItems: "center",
            gap: "2px",
            marginRight: "20px",
            marginBottom: "16px",
          }}
        >
          <FaRegEdit />
          Edit
        </SecondaryButton> */}
      </HeadingContent>
      <Wrapper>
        {assetData.map((data) => (
          <AssetCard key={data.id} label={data.title} value={data.subTitle} />
        ))}
      </Wrapper>
    </>
  );
};

export default FeeInAsset;

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

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;
