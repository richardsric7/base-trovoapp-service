"use client";

import SecondaryButton from "@/components/SecondaryButton";
import React from "react";
import { FaRegEdit } from "react-icons/fa";
import styled from "styled-components";

const dummyStakeholders = {
  assetCustodian: "SunTrust Capital Partners",
  assetManager: "Brikkle Asset Management",
  regulator: "Securities and Exchange Commission SEC Nigeria",
  issuingHouse: "Test Issuing House Ltd",
  legalAdviser: "LexTrust Legal Advisory Ltd",
  ratingAgency: "Global Rating Agency",
  trustee: "Sun Trustees Ltd",
};

const AssignedStakeHolder: React.FC = () => {
  const stakeHoldersData = [
    {
      id: 1,
      title: "Asset Custodian",
      name: dummyStakeholders.assetCustodian,
    },
    {
      id: 2,
      title: "Asset Manager",
      name: dummyStakeholders.assetManager,
    },
    {
      id: 3,
      title: "Regulator",
      name: dummyStakeholders.regulator,
    },
    {
      id: 4,
      title: "Issuing House",
      name: dummyStakeholders.issuingHouse,
    },
    {
      id: 5,
      title: "Legal Adviser",
      name: dummyStakeholders.legalAdviser,
    },
    {
      id: 6,
      title: "Rating Agency",
      name: dummyStakeholders.ratingAgency,
    },
    {
      id: 7,
      title: "Trustee",
      name: dummyStakeholders.trustee,
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
            gap: "4px",
            marginBottom: "16px",
          }}
        >
          <FaRegEdit />
          Edit
        </SecondaryButton>
      </HeadingContent>

      <Wrapper>
        {stakeHoldersData.map((item) => (
          <Card key={item.id}>
            <Text>{item.title}</Text>

            <CardContent>
              <BoldText>
                {item.name ? item.name.charAt(0).toUpperCase() : "?"}
              </BoldText>

              <Title title={item.name}>{item.name}</Title>
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
  border-radius: 8px;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 6px;
`;

const Text = styled.p`
  font-weight: 400;
  font-size: 14px;
  color: #828282;
`;

const CardContent = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
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
`;

const Title = styled.h2`
  font-weight: 600;
  font-size: 14px;
  color: #00225a;

  max-width: 140px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
`;
