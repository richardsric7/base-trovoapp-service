"use client";
import SecondaryButton from "@/components/SecondaryButton";
import { TokenizationRecord } from "@/redux/api/assettokenization";
import React from "react";
import { FaRegEdit } from "react-icons/fa";
import styled from "styled-components";

interface AssignedStakeHolderProps {
  asset?: TokenizationRecord;
}

const AssignedStakeHolder: React.FC<AssignedStakeHolderProps> = ({ asset }) => {
  const truncateWords = (text: string, limit: number) => {
    return text.length > limit ? text.slice(0, limit) + "..." : text;
  };

  const stakeHoldersData = [
    {
      id: 1,
      title: "Asset Custodian",
      Name: asset?.approvedAssetCustodianInfo?.assetCustodianName,
    },
    {
      id: 2,
      title: "Asset Manager",
      Name: asset?.assetManagerInfo?.assetManagerName,
    },
    {
      id: 3,
      title: "Regulator",
      Name: "Securities and Exchange Commission SEC Nigeria",
    },

    {
      id: 4,
      title: "Issuing House",
      Name: asset?.assetIssuingHouseInfo?.assetIssuingHouseName,
    },
    {
      id: 5,
      title: "Legal Adviser",
      Name: asset?.legalAndProfesionalPartnerInfo?.partnerName,
    },
    {
      id: 6,
      title: "Rating Agency",
      Name: asset?.ratingAgencyInfo?.agencyName,
    },
    {
      id: 7,
      title: "Trustee",
      Name: asset?.trusteeInfo?.trusteeName,
    },
  ];

  return (
    <>
      <HeadingContent>
        <Heading>Assigned Stakeholders</Heading>

        <SecondaryButton
          buttonStyle={{
            width: "100px",
            marginRight: "20px",
            display: "flex",
            alignItems: "center",
            gap: "2px",
            marginBottom: "16px",
          }}
        >
          <FaRegEdit />
          Edit
        </SecondaryButton>
      </HeadingContent>

      <Wrapper>
        {stakeHoldersData.map((items) => (
          <Card key={items.id}>
            <Text>{items.title}</Text>
            <CardContent>
              <BoldText>
                {items.Name ? items.Name.charAt(0).toUpperCase() : "?"}
              </BoldText>
              <div>
                <Title title={items.Name}>{items.Name}</Title>
              </div>
            </CardContent>
          </Card>
        ))}
      </Wrapper>
    </>
  );
};

export default AssignedStakeHolder;

const Wrapper = styled.section`
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  width: 100%;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
`;

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const Card = styled.div`
  background-color: #f2f6f9;
  border: 1px solid #f3f2f2;
  height: 77px;
  gap: 8px;
  border-radius: 8px;
  padding: 8px;
  display: flex;

  flex-direction: column;
`;

const Text = styled.p`
  font-weight: 400;
  font-size: 14px;
  line-height: 100%;
  letter-spacing: -0.5%;
  color: #828282;
`;

const CardContent = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
  width: 100%;
`;

const BoldText = styled.div`
  background-color: #007cdf;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #ffffff;
  font-weight: 600;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;
`;

const Title = styled.h2`
  font-weight: 600;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #00225a;

  display: block;
  max-width: 140px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
`;
const Email = styled.p`
  font-weight: 400;
  font-size: 11px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #828282;
`;
