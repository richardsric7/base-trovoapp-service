import React from "react";
import styled from "styled-components";
import AssetCard from "../../assettokenization/components/AssetCard";

const IncomeSummary = () => {
  const assetDetails = [
    {
      id: 1,
      title: "Total Revenue Generated",

      text: "₦400,000,000.00",
    },
    {
      id: 2,
      title: "Currency",

      text: "Naira",
    },
  ];
  return (
    <div>
      <Heading>Income Summary</Heading>
      <Wrapper>
        <AssetDetailsContent>
          {assetDetails.map((data) => (
            <AssetCard key={data.id} label={data.title} value={data.text} />
          ))}
        </AssetDetailsContent>
      </Wrapper>

      <Wrapper2>
        <AssetCard
          label="Revenue Notes"
          value="Lorem ipsum dolor sit amet consectetur. Bibendum nisl vivamus phasellus velit eleifend lacinia. Nibh elit vestibulum sed nunc vitae nibh suspendisse. Commodo erat sit amet duis at."
        />
      </Wrapper2>
    </div>
  );
};

export default IncomeSummary;

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

const Wrapper2 = styled.section`
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  margin-top: 20px;
`;
