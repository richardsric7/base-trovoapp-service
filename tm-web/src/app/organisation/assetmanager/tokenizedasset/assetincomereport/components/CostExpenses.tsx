import AssetCard from "@/app/(dashboard)/assettokenization/components/AssetCard";
import React from "react";
import styled from "styled-components";

const CostExpenses = () => {
  const assetDetails = [
    {
      id: 1,
      title: "Platform Fee",

      text: "₦400,000,000.00",
    },
    {
      id: 2,
      title: "Trustee Fee",

      text: "₦400,000,000.00",
    },

    {
      id: 3,
      title: "Asset Manager Fee",

      text: "₦400,000,000.00",
    },

    {
      id: 4,
      title: "Custodian Fee",

      text: "₦400,000,000.00",
    },

    {
      id: 5,
      title: "Regulatory Fee",

      text: "₦400,000,000.00",
    },
    {
      id: 6,
      title: "Maintenance and Operations",

      text: "₦400,000,000.00",
    },

    {
      id: 7,
      title: "Insurance Premium",

      text: "₦400,000,000.00",
    },
    {
      id: 8,
      title: "Other Costs",

      text: "₦400,000,000.00",
    },
  ];
  return (
    <div>
      <Heading>Cost & Breakdown Expenses</Heading>
      <Wrapper>
        <AssetDetailsContent>
          {assetDetails.map((data) => (
            <AssetCard key={data.id} label={data.title} value={data.text} />
          ))}
        </AssetDetailsContent>
      </Wrapper>

      <Wrapper2>
        <AssetCard
          label="Cost Notes"
          value="Lorem ipsum dolor sit amet consectetur. Bibendum nisl vivamus phasellus velit eleifend lacinia. Nibh elit vestibulum sed nunc vitae nibh suspendisse. Commodo erat sit amet duis at."
        />
      </Wrapper2>
    </div>
  );
};

export default CostExpenses;

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
