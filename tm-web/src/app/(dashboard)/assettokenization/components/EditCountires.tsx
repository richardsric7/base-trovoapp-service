"use client";
import { Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import {
  TokenizationRecord,
  UpdateTokenizationPayload,
} from "@/redux/api/assettokenization";
import React, { useEffect, useState } from "react";
import styled from "styled-components";

interface EditCountriesProps {
  isOpen: boolean;
  asset?: TokenizationRecord;
  onClose: () => void;
  onSave?: (updatedData: Partial<UpdateTokenizationPayload>) => void;
}

const EditCountires: React.FC<EditCountriesProps> = ({
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

    setEditedAsset((prev) => ({ ...prev, [name]: value }));
  };
  return (
    <>
      <Modal title="Edit Countries" isOpen={isOpen} onClose={onClose}>
        <ModalContent>
          <FormGroup>
            <Label>Exempted Countries</Label>
            <Input
              name="exemptedCountries"
              value={editedAsset.exemptedCountries || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <PrimaryButton
            onClick={() => {
              if (onSave) {
                const cleaned = (editedAsset.exemptedCountries || "")
                  .toUpperCase()
                  .replace(/\s+/g, ",") // Replace all whitespace with commas
                  .replace(/,+/g, ",") // Replace multiple commas with a single one
                  .replace(/^,|,$/g, ""); // Trim leading/trailing commas

                onSave({ ...editedAsset, exemptedCountries: cleaned });
              }
            }}
            buttonStyle={{ width: "100%" }}
          >
            Save Changes
          </PrimaryButton>
        </ModalContent>
      </Modal>
    </>
  );
};

export default EditCountires;

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
