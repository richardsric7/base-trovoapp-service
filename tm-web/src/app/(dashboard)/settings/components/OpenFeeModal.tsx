import { Modal } from "@/components";
import Image from "next/image";
import React, { useState } from "react";
import styled from "styled-components";
import goldTag from "../../../../assets/images/goldtag.svg";
import platinumTag from "../../../../assets/images/platinumtag.svg";

import diamondTag from "../../../../assets/images/diamondtag.svg";

interface OpenFeesModalprops {
  isOpen: boolean;
  onClose: () => void;
}

interface PlanBoxProps {
  selected?: boolean;
}
const OpenFeeModal: React.FC<OpenFeesModalprops> = ({ isOpen, onClose }) => {
  const [selectedPlan, setSelectedPlan] = useState<string>("gold");

  const inputData = [
    {
      key: "month",
      label: "Monthly Subscription",
      unit: "CNGN",
    },
    {
      key: "annual",
      label: "Annual Subscription",
      unit: "CNGN",
    },
    {
      key: "lifetime",
      label: "Life-time Subscription",
      unit: "CNGN",
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
    <>
      <Modal
        title="Add Membership Fee"
        isOpen={isOpen}
        onClose={onClose}
        closeIconPosition="left"
      >
        {" "}
        <ModalContent>
          <ModalText>Select membership category</ModalText>

          <ButtonContainer>
            <PlanBox
              type="button"
              selected={selectedPlan === "gold"}
              onClick={() => setSelectedPlan("gold")}
            >
              {" "}
              <Image
                src={goldTag}
                alt="gold-tag-icon"
                height={20}
                width={20}
              />{" "}
              Gold Plan
            </PlanBox>

            <PlanBox
              selected={selectedPlan === "platinum"}
              onClick={() => setSelectedPlan("platinum")}
            >
              {" "}
              <Image
                src={platinumTag}
                alt="platinum-tag-icon"
                height={20}
                width={20}
              />{" "}
              Platinum Plan
            </PlanBox>
            <PlanBox
              type="button"
              selected={selectedPlan === "diamond"}
              onClick={() => setSelectedPlan("diamond")}
            >
              {" "}
              <Image
                src={diamondTag}
                alt="diamond-tag-icon"
                height={20}
                width={20}
              />{" "}
              Diamond Plan
            </PlanBox>
          </ButtonContainer>

          {inputData.map((item) => (
            <FormGroup key={item.key}>
              <Label>{item.label}</Label>
              <FeeInputWrapper>
                <FeeInputField
                  name={item.key}
                  placeholder="0.00"
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
        </ModalContent>{" "}
      </Modal>
    </>
  );
};

export default OpenFeeModal;

const ModalContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

const ModalText = styled.p`
  color: #828282;

  font-weight: 500;

  font-size: 14px;
  line-height: 24px;
  letter-spacing: 0.1px;
`;

const ButtonContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: 8px;
`;

const PlanBox = styled.button<PlanBoxProps>`
  // border: 1px solid #007cdf;
  border-radius: 32px;
  // color: #007cdf;

  border: 1px solid ${({ selected }) => (selected ? "#007cdf" : "#828282")};
  background: ${({ selected }) => (selected ? "" : "#fff")};
  color: ${({ selected }) => (selected ? "#007cdf" : "#828282")};
  font-weight: 400;
  font-size: 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 8px;
  width: 80%;
  height: 40px;
  cursor: pointer;
  font-family: inherit;
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
const FeeInputWrapper = styled.div`
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  display: flex;
  align-items: center;
`;
