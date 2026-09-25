"use client";

import React, { useState } from "react";
import styled from "styled-components";
import { Modal } from "@/components/CustomModal";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";

interface PersonalSecretCreateModalProps {
  mode: "create";
  friendlyLabel: string;
  isOpen: boolean;
  isSubmitting: boolean;
  onClose: () => void;
  onConfirm: (newValue: string) => void;
}

interface PersonalSecretDeleteModalProps {
  mode: "delete";
  friendlyLabel: string;
  isOpen: boolean;
  isSubmitting: boolean;
  onClose: () => void;
  onConfirm: () => void;
}

type PersonalSecretModalProps = PersonalSecretCreateModalProps | PersonalSecretDeleteModalProps;

// Create-only, matching the backend's own lifecycle: there is no update-in-place endpoint, so
// this modal never accepts an existing value to edit — a fresh value is always a brand new
// secret. Delete mode is a plain confirmation since the backend fully purges the secret's
// version history (not a soft delete).
const PersonalSecretModal: React.FC<PersonalSecretModalProps> = (props) => {
  const [value, setValue] = useState("");

  if (props.mode === "delete") {
    const { friendlyLabel, isOpen, isSubmitting, onClose, onConfirm } = props;
    return (
      <Modal title={`Delete ${friendlyLabel}`} isOpen={isOpen} onClose={onClose} style={{ width: "35%" }}>
        <ConfirmBody>
          <ConfirmText>
            This permanently purges the secret and its full version history from Vault — this is
            not a soft delete, and there is no way to recover the value afterward.
          </ConfirmText>
          <ButtonRow>
            <SecondaryButton buttonStyle={{ width: "100%" }} onClick={onClose}>
              Cancel
            </SecondaryButton>
            <PrimaryButton
              buttonStyle={{ width: "100%", backgroundColor: "#d92d20" }}
              disabled={isSubmitting}
              onClick={onConfirm}
            >
              {isSubmitting ? "Deleting..." : "Delete"}
            </PrimaryButton>
          </ButtonRow>
        </ConfirmBody>
      </Modal>
    );
  }

  const { friendlyLabel, isOpen, isSubmitting, onClose, onConfirm } = props;
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!value.trim()) return;
    onConfirm(value.trim());
  };

  return (
    <Modal title={friendlyLabel} isOpen={isOpen} onClose={onClose} style={{ width: "40%" }}>
      <Form onSubmit={handleSubmit}>
        <Field>
          <Label htmlFor="personal-secret-value">Signer</Label>
          <Input
            id="personal-secret-value"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            autoComplete="off"
            autoFocus
          />
          <Hint>
            This value can never be read back once saved — to change it later, delete and
            recreate it.
          </Hint>
        </Field>
        <PrimaryButton buttonStyle={{ width: "100%" }} disabled={isSubmitting}>
          {isSubmitting ? "Creating..." : "Create"}
        </PrimaryButton>
      </Form>
    </Modal>
  );
};

export default PersonalSecretModal;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 20px;
`;
const Field = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;
const Label = styled.label`
  color: #00225a;
  font-size: 14px;
  font-weight: 600;
`;
const Input = styled.input`
  height: 48px;
  border: 1px solid #d0d5dd;
  border-radius: 8px;
  padding: 0 14px;
  color: #00225a;
  font: inherit;
  outline: none;
  &:focus {
    border-color: #007cdf;
    box-shadow: 0 0 0 3px rgba(0, 124, 223, 0.12);
  }
`;
const Hint = styled.span`
  color: #828282;
  font-size: 12px;
`;
const ConfirmBody = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
`;
const ConfirmText = styled.p`
  color: #667085;
  font-size: 14px;
  line-height: 1.55;
  margin: 0;
`;
const ButtonRow = styled.div`
  display: flex;
  gap: 12px;
`;
