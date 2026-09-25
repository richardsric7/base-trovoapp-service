"use client";
import { Modal } from "@/components";
import React from "react";
import styled from "styled-components";

interface AddFieldProps {
  isOpen: boolean;
  setIsOpenAdd: React.Dispatch<React.SetStateAction<boolean>>;
}

const AddFieldModal: React.FC<AddFieldProps> = ({ isOpen, setIsOpenAdd }) => {
  return (
    <>
      <Modal
        title="Add Field"
        isOpen={isOpen}
        onClose={() => setIsOpenAdd(false)}
        closeIconPosition="left"
      >
        <FormGroup>
          <Label>Title</Label>
          <FeeInputWrapper>
            <FeeInputField name="" placeholder="Enter title" value="" />
          </FeeInputWrapper>
        </FormGroup>
        <FormGroup>
          <Label> Fee</Label>
          <FeeInputWrapper>
            <DividerGroup>
              <FeeLabel>₦</FeeLabel>
              <VerticalDivider />
            </DividerGroup>
            <FeeInputField name="" placeholder="0" value="" />
          </FeeInputWrapper>
        </FormGroup>
      </Modal>
    </>
  );
};

export default AddFieldModal;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  color: #828282;
  font-weight: 500;
  font-size: 14px;
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
    color: #bdbdbd;
    font-family: inherit;
    font-size: 14px;
  }
`;

const FeeInputWrapper = styled.div`
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  display: flex;
  align-items: center;
`;
