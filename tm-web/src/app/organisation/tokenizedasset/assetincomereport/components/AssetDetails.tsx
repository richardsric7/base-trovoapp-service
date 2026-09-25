"use client";
import AssetCard from "@/app/(dashboard)/assettokenization/components/AssetCard";
import React from "react";
import styled from "styled-components";

const AssetDetails = () => {
  const assetDetails = [
    {
      id: 1,
      title: "Asset Name",

      text: "Atlantis Estate",
    },
    {
      id: 2,
      title: "Reporting Period",

      text: "01 Jun, 2025 - 30 Jun, 2025",
    },
    {
      id: 3,
      title: "Report Title",

      text: "June, 2025 Report",
    },
  ];
  return (
    <div>
      <Heading>Period & Asset Details</Heading>
      <Wrapper>
        <AssestInfoSection>
          <Avatar />

          <div>
            <AssetName>AEI</AssetName>
          </div>
        </AssestInfoSection>

        <AssetDetailsContent>
          {assetDetails.map((data) => (
            <AssetCard key={data.id} label={data.title} value={data.text} />
          ))}
        </AssetDetailsContent>
      </Wrapper>
    </div>
  );
};

export default AssetDetails;

const AssestInfoSection = styled.div`
  display: flex;
  column-gap: 8px;
  align-items: center;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;

const AssetName = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  margin: 4px 0;
  color: #00225a;
  text-transform: capitalize;
`;

const AssetDetailsContent = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  padding-top: 10px;
  width: 100%;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  margin-bottom: 6px;
  color: #00225a;
`;

const Wrapper = styled.section`
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
`;
