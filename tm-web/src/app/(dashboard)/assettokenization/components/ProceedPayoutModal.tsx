import { Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import React from "react";
import styled from "styled-components";

interface ProceedPayoutModalProps {
  isOpen: boolean;
  setOpenProceeds: React.Dispatch<React.SetStateAction<boolean>>;
  setOnSuccess: React.Dispatch<React.SetStateAction<boolean>>;
}

const ProceedPayoutModal = ({
  isOpen,
  setOpenProceeds,
  setOnSuccess,
}: ProceedPayoutModalProps) => {
  const handleSaveProceeds = () => {
    setOnSuccess(true);
    setOpenProceeds(false);
  };

  return (
    <Modal
      title="Payout Proceeds"
      isOpen={isOpen}
      onClose={() => setOpenProceeds(false)}
      closeIconPosition="right"
    >
      <ModalContent>
        <FormGroup>
          <Label>Enter Amount</Label>
          <InputBox>
            <AmountInput name="amount" value="" placeholder="00" />

            <AmountInfo>
              <Avatar />
              <CurrencyText>Trov</CurrencyText>
            </AmountInfo>
          </InputBox>
        </FormGroup>
        <FormGroup>
          <Label>Payout Description</Label>
          <DescriptionInput
            name="description"
            value=""
            placeholder="Enter Description"
          />
        </FormGroup>
        <FormGroup>
          <Label>
            Please request approval from the following users to proceed
          </Label>
          <ApproversBox>
            <Approver>obiaruku</Approver>
            <Approver>obiaruku</Approver>
            <Approver>obiaruku</Approver>
          </ApproversBox>
        </FormGroup>
        <PrimaryButton
          onClick={handleSaveProceeds}
          buttonStyle={{ width: "100%" }}
        >
          Save Changes
        </PrimaryButton>
      </ModalContent>
    </Modal>
  );
};

export default ProceedPayoutModal;

const ModalContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  color: #00225a;
  font-weight: 500;
  font-size: 16px;
  line-height: 20px;
`;

const InputBox = styled.div`
  padding: 10px;
  border: 1px solid #bdbdbd;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 400;
  width: 100%;
  color: #00225a;
  font-family: inherit;
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const AmountInput = styled.input`
  font-size: 14px;
  font-weight: 400;
  width: 100%;
  border: none;
  outline: none;
  color: #00225a;
  font-family: inherit;
`;

const AmountInfo = styled.div`
  display: flex;
  align-items: center;
  gap: 2px;
`;

const Avatar = styled.div`
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const CurrencyText = styled.p`
  font-weight: 600;
  font-size: 14px;
  line-height: 24px;
  color: #00225a;
`;

const DescriptionInput = styled.input`
  padding: 6px 4px 60px 4px;
  border: 1px solid #bdbdbd;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 400;
  width: 100%;
  outline: none;
  color: #00225a;
  font-family: inherit;
`;

const ApproversBox = styled.div`
  padding: 10px;
  border: 1px solid #007cdf;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 400;
  width: 100%;
  display: flex;
  align-items: center;
  gap: 4px;
`;

const Approver = styled.div`
  background-color: #007cdf;
  color: #fff;
  border-radius: 4px;
  padding: 4px 8px;
  text-align: center;
`;
