"use client";
import PrimaryButton from "@/components/PrimaryButton";
import React from "react";
import styled from "styled-components";

const TrovoP2pFeesPage = () => {
  return (
    <PageContainer>
      <Title>TrovoP2p Fees</Title>

      <div>
        <SubTitle>Trovo P2P fee collection and distribution</SubTitle>
        <FormGroup>
          <Label>Maker fee</Label>
          <FeeInputWrapper>
            <FeeInputField name="" placeholder="0" value="" />
            <DividerGroup>
              <VerticalDivider />
              <FeeLabel>%</FeeLabel>
            </DividerGroup>
          </FeeInputWrapper>
        </FormGroup>
        <FormGroup>
          <Label>Taker fee</Label>
          <FeeInputWrapper>
            <FeeInputField name="" placeholder="0" value="" />
            <DividerGroup>
              <VerticalDivider />
              <FeeLabel>%</FeeLabel>
            </DividerGroup>
          </FeeInputWrapper>
        </FormGroup>

        <PrimaryButton buttonStyle={{ width: "100%" }}>
          Save Changes
        </PrimaryButton>
      </div>
    </PageContainer>
  );
};

export default TrovoP2pFeesPage;
const PageContainer = styled.section``;
const Title = styled.h1`
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  color: #00225a;
`;

const SubTitle = styled.h1`
  font-size: 14px;
  font-weight: 600;
  margin: 10px 0;
  color: #00225a;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  color: #828282;
  font-weight: 500;
  font-size: 14px;
  margin-top: 4px;
`;

const VerticalDivider = styled.div`
  width: 1px;
  height: 24px;
  background-color: #e0e0e0;
`;

const FeeLabel = styled.p`
  font-weight: 500;
  font-size: 14px;
  color: #191919;
`;

const DividerGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const FeeInputField = styled.input`
  font-size: 14px;
  width: 100%;
  outline: none;
  font-family: inherit;
  border: none;
`;

const FeeInputWrapper = styled.div`
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  display: flex;
  align-items: center;
`;
