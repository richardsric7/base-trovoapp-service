import { Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import React from "react";
import styled from "styled-components";
interface AddSubSectorProps {
  add: boolean;
  setAdd: React.Dispatch<React.SetStateAction<boolean>>;
}

const AddSubSectorModal: React.FC<AddSubSectorProps> = ({ add, setAdd }) => {
  return (
    <Modal
      title="Add Subsector"
      isOpen={add}
      onClose={() => setAdd(false)}
      closeIconPosition="left"
    >
      <FormGroup>
        <Label>Asset Sector</Label>
        <Text>Real Estate</Text>
      </FormGroup>

      <FormGroup>
        <Label>Sub-sector</Label>
        <Input type="text" placeholder="Enter sub-sector ..." value="" />
      </FormGroup>
      <PrimaryButton>Save Sub-sector</PrimaryButton>
    </Modal>
  );
};

export default AddSubSectorModal;

const Input = styled.input`
  width: 100%;
  height: 48px;
  justify-content: space-between;
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  font-family: inherit;
  margin: 10px 0;
  padding: 0 6px;
`;

const Label = styled.p`
  color: #828282;
  font-weight: 500;
  font-size: 14px;
  margin-top: 14px;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Text = styled.p`
  font-weight: 600;
  color: #00225a;
  font-size: 14px;
  line-height: 24px;
  letter-spacing: 0.1px;
`;
