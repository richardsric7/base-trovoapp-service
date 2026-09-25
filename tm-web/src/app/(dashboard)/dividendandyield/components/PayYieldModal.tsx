import { Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import { DatePicker } from "antd";

import styled from "styled-components";
import dayjs from "dayjs";
import customParseFormat from "dayjs/plugin/customParseFormat";

dayjs.extend(customParseFormat);

const dateFormat = "YYYY-MM-DD";

interface PayYieldModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const PayYieldModal: React.FC<PayYieldModalProps> = ({ isOpen, onClose }) => {
  return (
    <>
      <Modal title="Yield Payout" isOpen={isOpen} onClose={onClose}>
        <FormGroup>
          <Label> Upcoming Payout Amount</Label>
          <FeeInputWrapper>
            <DividerGroup>
              <FeeLabel>₦</FeeLabel>
              <VerticalDivider />
            </DividerGroup>
            <FeeInputField name="" placeholder=" 100,000,000.00" value="" />
          </FeeInputWrapper>
        </FormGroup>
        <FormGroup>
          <Label>Payout Due Date</Label>
          <StyledDatePicker defaultValue={dayjs("2025-08-03", dateFormat)} />
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

export default PayYieldModal;

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

const StyledDatePicker = styled(DatePicker)`
  width: 100%;
  height: 44px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 0 10px;
  display: flex;
  align-items: center;
    color: #00225a;
  font-family: "Montserrat", sans-serif;

  .ant-picker-input > input {
    font-size: 14px;
    font-family: inherit;
    color:  &::placeholder {
    color: #00225a;
   
    font-family: "Montserrat", sans-serif;
  }
`;
