"use client";

import React, { useState } from "react";
import styled from "styled-components";
import { Modal } from "@/components/CustomModal";
import PrimaryButton from "@/components/PrimaryButton";

interface SetSignerValueModalProps {
  label: string;
  isOpen: boolean;
  isSubmitting: boolean;
  onClose: () => void;
  onSubmit: (newValue: string) => void;
}

const SetSignerValueModal: React.FC<SetSignerValueModalProps> = ({
  label,
  isOpen,
  isSubmitting,
  onClose,
  onSubmit,
}) => {
  const [value, setValue] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!value.trim()) return;
    onSubmit(value.trim());
  };

  return (
    <Modal title={label} isOpen={isOpen} onClose={onClose} style={{ width: "40%" }}>
      <Form onSubmit={handleSubmit}>
        <Field>
          <Label htmlFor="signer-value">Signer</Label>
          <Input
            id="signer-value"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            placeholder="e.g. SABC... or GABC..."
            autoComplete="off"
            autoFocus
          />
          <Hint>
            An active position triggers an on-chain signer rotation; a spare position is a plain
            validated write.
          </Hint>
        </Field>
        <PrimaryButton buttonStyle={{ width: "100%" }} disabled={isSubmitting}>
          {isSubmitting ? "Submitting..." : "Submit"}
        </PrimaryButton>
      </Form>
    </Modal>
  );
};

export default SetSignerValueModal;

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
