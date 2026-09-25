"use client";
import { Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import {
  TokenizationRecord,
  UpdateTokenizationPayload,
} from "@/redux/api/assettokenization";
import React, { useEffect, useState } from "react";
import styled from "styled-components";

interface EditAssetValueProps {
  isOpen: boolean;
  asset?: TokenizationRecord;
  onClose: () => void;
  onSave?: (updatedData: Partial<UpdateTokenizationPayload>) => void;
}

const EditAssetValue: React.FC<EditAssetValueProps> = ({
  isOpen,
  asset,
  onClose,
  onSave,
}) => {
  const [editedAsset, setEditedAsset] = useState<
    Partial<UpdateTokenizationPayload>
  >(asset || {});

  useEffect(() => {
    if (asset) {
      setEditedAsset(asset as Partial<UpdateTokenizationPayload>);
    }
  }, [asset]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    let newValue: string | number = value;
    if (
      name === "assetCurrentValue" ||
      name === "assetMscCostOutisdeOfValuation" ||
      name === "assetOwnerRetainedOrContributedValue"
    ) {
      newValue = Number(value);
    }
    setEditedAsset((prev) => ({ ...prev, [name]: newValue }));
  };

  return (
    <>
      <Modal title="Edit Asset Value" isOpen={isOpen} onClose={onClose}>
        <ModalContent>
          <FormGroup>
            <Label>Current Asset Value</Label>
            <Input
              name="assetCurrentValue"
              value={editedAsset.assetCurrentValue || ""}
              onChange={handleChange}
            />
          </FormGroup>
          <FormGroup>
            <Label>Cost Outside Valuation</Label>
            <Input
              name="assetMscCostOutisdeOfValuation"
              value={editedAsset.assetMscCostOutisdeOfValuation || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <FormGroup>
            <Label>Value to be Retained</Label>
            <Input
              name="assetOwnerRetainedOrContributedValue"
              value={editedAsset.assetOwnerRetainedOrContributedValue || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <PrimaryButton
            onClick={() => onSave && onSave(editedAsset)}
            buttonStyle={{ width: "100%" }}
          >
            Save Changes
          </PrimaryButton>
        </ModalContent>
      </Modal>
    </>
  );
};

export default EditAssetValue;

const ModalContent = styled.div``;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  color: #00225a;
  font-weight: 500;
  font-size: 16px;
  line-height: 24px;
`;

const Input = styled.input`
  padding: 10px;
  border: 1px solid #bdbdbd;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 400;
  width: 100%;
  outline: none;
  color: #00225a;
  font-family: inherit;
`;
