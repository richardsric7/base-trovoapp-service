import { Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import React from "react";
import styled from "styled-components";

interface PayDividendModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const PayDividendModal: React.FC<PayDividendModalProps> = ({
  isOpen,
  onClose,
}) => {
  return (
    <>
      <Modal title="Dividend Payout" isOpen={isOpen} onClose={onClose}>
        <FormGroup>
          <Label> Amount</Label>
          <FeeInputWrapper>
            <DividerGroup>
              <FeeLabel>₦</FeeLabel>
              <VerticalDivider />
            </DividerGroup>
            <FeeInputField name="" placeholder=" 400,000,000.00" value="" />
          </FeeInputWrapper>
        </FormGroup>

        <FormGroup>
          <Label>Payout Description</Label>
          <FeeInputWrapper>
            <FeeInputField
              name=""
              placeholder="Prefilled description"
              value=""
            />
          </FeeInputWrapper>
        </FormGroup>

        <FormGroup>
          <Label>
            Please request for approval from the following users to proceed
          </Label>
          <FeeInputWrapper2>
            <InputText>Onoja</InputText>
            <InputText>Nancy</InputText>
            <InputText>Obiaruku</InputText>
          </FeeInputWrapper2>
        </FormGroup>

        <PrimaryButton buttonStyle={{ width: "100%" }}>
          Request Approval
        </PrimaryButton>
      </Modal>
    </>
  );
};

export default PayDividendModal;
const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  color: #00225a;
  font-weight: 500;
  font-size: 16px;
  margin-top: 14px;
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

  &::placeholder {
    color: #00225a;
    font-family: inherit;
    font-size: 16px;
  }
`;

const FeeInputWrapper = styled.div`
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  display: flex;
  align-items: center;
`;

const FeeInputWrapper2 = styled.div`
  padding: 10px;
  border: 1px solid #007cdf;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
`;

const InputText = styled.p`
  color: #ffffff;
  background-color: #007cdf;
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: 400;
  font-size: 14px;
`;
