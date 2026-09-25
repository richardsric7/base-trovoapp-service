"use client";
import { Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import {
  TokenizationRecord,
  UpdateTokenizationPayload,
} from "@/redux/api/assettokenization";
import React, { useEffect, useState } from "react";
import styled from "styled-components";

interface EditOwnerProps {
  isOpen: boolean;
  asset?: TokenizationRecord;
  onClose: () => void;
  onSave?: (updatedData: Partial<UpdateTokenizationPayload>) => void;
}

const EditAssetOwner: React.FC<EditOwnerProps> = ({
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

    setEditedAsset((prev) => ({ ...prev, [name]: newValue }));
  };

  return (
    <>
      <Modal
        title="Edit Asset Owner"
        isOpen={isOpen}
        onClose={onClose}
        closeIconPosition="right"
      >
        <ModalContent>
          <FormGroup>
            <Label>OwnerShip</Label>
            <Input
              name="ownershipKind"
              value={editedAsset.ownershipKind || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <FormGroup>
            <Label>Third Party Ownership Type</Label>
            <Input
              name="ownershipType"
              value={editedAsset.ownershipType || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <FormGroup>
            <Label>Organizations Name</Label>
            <Input
              name="assetOwnerName"
              value={editedAsset.assetOwnerName || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <FormGroup>
            <Label>Organizations Address</Label>
            <Input
              name="assetOwnerAddress"
              value={editedAsset.assetOwnerAddress || ""}
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

export default EditAssetOwner;

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
