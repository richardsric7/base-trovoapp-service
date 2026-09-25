"use client";

import { Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import Image from "next/image";
import React, { useEffect, useState } from "react";
import styled from "styled-components";
import goldTag from "../../../../assets/images/goldtag.svg";
interface EditMembershipProps {
  isOpen: boolean;
  onClose: () => void;
  plan?: string; // e.g. "Gold Plan"
  initialValues?: {
    monthly?: string;
    annual?: string;
    lifetime?: string;
  };
  onSave?: (values: {
    monthly: string;
    annual: string;
    lifetime: string;
  }) => void;
}

const Editmembership: React.FC<EditMembershipProps> = ({
  isOpen,
  onClose,
  plan = "Gold Plan",
  initialValues,
  onSave,
}) => {
  const [monthly, setMonthly] = useState("");
  const [annual, setAnnual] = useState("");
  const [lifetime, setLifetime] = useState("");

  useEffect(() => {
    setMonthly(initialValues?.monthly ?? "");
    setAnnual(initialValues?.annual ?? "");
    setLifetime(initialValues?.lifetime ?? "");
  }, [initialValues, isOpen]);

  const handleSave = () => {
    onSave?.({ monthly, annual, lifetime });
  };

  return (
    <Modal title="Edit MemberShip Fee" isOpen={isOpen} onClose={onClose}>
      <PlanBox>
        {" "}
        <Image src={goldTag} alt="gold-tag-icon" height={20} width={20} />{" "}
        {plan}
      </PlanBox>

      <div>
        <FormGroup>
          <Label>Monthly Subscription</Label>
          <FeeInputWrapper>
            <FeeInputField
              name=""
              placeholder="0"
              value={monthly}
              onChange={(e) => setMonthly(e.target.value)}
            />
            <DividerGroup>
              <VerticalDivider />
              <FeeLabel>CNGN</FeeLabel>
            </DividerGroup>
          </FeeInputWrapper>
        </FormGroup>
        <FormGroup>
          <Label>Annual Subscription</Label>
          <FeeInputWrapper>
            <FeeInputField
              name=""
              placeholder="0"
              value={annual}
              onChange={(e) => setAnnual(e.target.value)}
            />
            <DividerGroup>
              <VerticalDivider />
              <FeeLabel>CNGN</FeeLabel>
            </DividerGroup>
          </FeeInputWrapper>
        </FormGroup>

        <FormGroup>
          <Label>Life-time Subscription</Label>
          <FeeInputWrapper>
            <FeeInputField
              name=""
              placeholder="0"
              value={lifetime}
              onChange={(e) => setLifetime(e.target.value)}
            />
            <DividerGroup>
              <VerticalDivider />
              <FeeLabel>CNGN</FeeLabel>
            </DividerGroup>
          </FeeInputWrapper>
        </FormGroup>

        <PrimaryButton buttonStyle={{ width: "100%" }} onClick={handleSave}>
          Save Changes
        </PrimaryButton>
      </div>
    </Modal>
  );
};

export default Editmembership;

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
