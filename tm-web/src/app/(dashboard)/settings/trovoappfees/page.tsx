"use client";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import React, { useState } from "react";
import styled from "styled-components";
import OpenFeeModal from "../components/OpenFeeModal";
import goldTag from "../../../../assets/images/goldtag.svg";
import Image from "next/image";

import { FaRegEdit } from "react-icons/fa";
import { RiDeleteBin6Line } from "react-icons/ri";
import Editmembership from "../components/Editmembership";
const TrovoAppFeesPage = () => {
  const [isopenFeeModal, setIsOpenFeeModal] = useState(false);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [goldPlanFees, setGoldPlanFees] = useState({
    monthly: "10",
    annual: "10",
    lifetime: "10",
  });

  const inputData = [
    {
      key: "walletSwapFee",
      label: " Trovo Wallet Swap fee collection and distribution",
      unit: "%",
    },
    {
      key: "sharedAccessFee",
      label:
        "Shared access payment transaction fee collection and distribution ",
      unit: "%",
    },
    {
      key: "walletRecoveryFee",
      label:
        "Trovo Wallet Recovery subscription fee collection and distribution ",
      unit: "TROV",
    },
    {
      key: "subWalletFee",
      label: "Sub wallet creation fee collection and distribution ",
      unit: "TROV",
    },
    {
      key: "assetWithdrawalFee",
      label: "Wrapped Asset Withdrawal fee collection and distribution",
      unit: "%",
    },
    {
      key: "statementGenerationFee",
      label: "Statement of account generation fee collection and distribution",
      unit: "TROV",
    },
  ];

  const [formValues, setFormValues] = useState(
    Object.fromEntries(inputData.map((item) => [item.key, ""]))
  );

  const handleInputChange =
    (key: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
      setFormValues((prev) => ({ ...prev, [key]: e.target.value }));
    };

  return (
    <PageContainer>
      <Title>TrovoApp Fees</Title>

      <div>
        <SubTitle>Trovo App</SubTitle>

        {inputData.map((item) => (
          <FormGroup key={item.key}>
            <Label>{item.label}</Label>
            <FeeInputWrapper>
              <FeeInputField
                name={item.key}
                placeholder="0"
                value={formValues[item.key]}
                onChange={handleInputChange(item.key)}
              />
              <DividerGroup>
                <VerticalDivider />
                <FeeLabel>{item.unit}</FeeLabel>
              </DividerGroup>
            </FeeInputWrapper>
          </FormGroup>
        ))}
        <div>
          <SubTitle> Membership fee collection and distribution TROV</SubTitle>

          <SecondaryButton
            buttonStyle={{ width: "100%" }}
            onClick={() => setIsOpenFeeModal(!isopenFeeModal)}
          >
            + Add Fee
          </SecondaryButton>
        </div>
        <FeesDisplay>
          <FeeHeader>
            <PlanBox>
              {" "}
              <Image
                src={goldTag}
                alt="gold-tag-icon"
                height={20}
                width={20}
              />{" "}
              Gold Plan
            </PlanBox>

            <ButtonBox>
              <IconBox
                borderColor="#007CDF"
                onClick={() => setIsEditOpen(true)}
                aria-label="Edit plan"
              >
                <FaRegEdit color="#007CDF" />
              </IconBox>
              <IconBox borderColor="#BE3800">
                <RiDeleteBin6Line color="#BE3800" />
              </IconBox>
            </ButtonBox>
          </FeeHeader>

          <FeeBody>
            <div>
              <Text>Monthly Subscription</Text>
              <Value>30CNGN</Value>
            </div>
            <div>
              <Text>Annual Subscription</Text>
              <Value>70CNGN</Value>
            </div>
            <div>
              <Text>Life-time Subscription</Text>
              <Value>10CNGN</Value>
            </div>
          </FeeBody>
        </FeesDisplay>
        <PrimaryButton buttonStyle={{ width: "100%" }}>
          Save Fee Changes
        </PrimaryButton>
      </div>

      {isopenFeeModal && (
        <OpenFeeModal
          isOpen={isopenFeeModal}
          onClose={() => setIsOpenFeeModal(false)}
        />
      )}

      <Editmembership
        isOpen={isEditOpen}
        onClose={() => setIsEditOpen(false)}
        plan="Gold Plan"
        initialValues={goldPlanFees}
        onSave={(vals:any) => {
          setGoldPlanFees(vals); // update local state (or call API)
          setIsEditOpen(false);
        }}
      />
    </PageContainer>
  );
};

export default TrovoAppFeesPage;
const PageContainer = styled.section`
  width: 70%;
  margin: 0 auto;
`;
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

const ButtonBox = styled.div`
  display: flex;
  gap: 8px;
`;

const IconBox = styled.div<{ borderColor?: string }>`
  width: 40px;
  height: 40px;
  border: 1px solid ${({ borderColor }) => borderColor || "#007cdf"};
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
`;

const FeesDisplay = styled.div`
  background-color: #f2f6f9;
  margin-top: 10px;
  opacity: 1;
  border-radius: 8px;
  padding: 16px;
`;

const FeeHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
`;
const PlanBox = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 600;
  font-style: SemiBold;
  font-size: 14px;
  color: #00225a;
  line-height: 24px;
  letter-spacing: 0.1px;
`;

const FeeBody = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
`;

const Text = styled.p`
  color: #828282;
  font-weight: 400;
  font-size: 14px;
  padding-bottom: 6px;
  line-height: 100%;
  letter-spacing: 0%;
`;

const Value = styled.p`
  color: #00225a;
  font-weight: 500;
  font-style: Medium;
  font-size: 14px;
  leading-trim: NONE;
  line-height: 100%;
  letter-spacing: 0%;
`;
