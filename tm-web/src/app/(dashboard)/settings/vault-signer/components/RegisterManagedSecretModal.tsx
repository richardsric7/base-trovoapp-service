"use client";

import React, { useState } from "react";
import styled from "styled-components";
import { Modal } from "@/components/CustomModal";
import PrimaryButton from "@/components/PrimaryButton";
import { showErrorToast, showSuccessToast } from "@/components";
import { useRegisterManagedSecretMutation } from "@/redux/api/vaultSigner/api";
import type { IRegisterManagedSecretRequest } from "@/redux/api/vaultSigner/interface";
import { getErrorMessage } from "@/app/vault-signer/components/MySignerSecrets";

interface RegisterManagedSecretModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const emptyForm: IRegisterManagedSecretRequest = {
  label: "",
  vaultMount: "",
  vaultPath: "",
  vaultField: "",
  walletAddress: "",
  activeSigningCount: 3,
};

// Register-only: the CSV that lives at vaultMount/vaultPath/vaultField is expected to already
// hold real, pre-existing signer keys (Section 0/4 of the backend plan) — this form just points
// this feature at an existing Vault secret, it never seeds one.
const RegisterManagedSecretModal: React.FC<RegisterManagedSecretModalProps> = ({
  isOpen,
  onClose,
}) => {
  const [form, setForm] = useState<IRegisterManagedSecretRequest>(emptyForm);
  const [registerManagedSecret, { isLoading }] = useRegisterManagedSecretMutation();

  const update = (field: keyof IRegisterManagedSecretRequest) => (
    e: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const value = field === "activeSigningCount" ? Number(e.target.value) : e.target.value;
    setForm((prev) => ({ ...prev, [field]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await registerManagedSecret(form).unwrap();
      showSuccessToast("Managed secret registered");
      setForm(emptyForm);
      onClose();
    } catch (error) {
      showErrorToast(getErrorMessage(error, "Unable to register managed secret"));
    }
  };

  return (
    <Modal title="Register Managed Secret" isOpen={isOpen} onClose={onClose} style={{ width: "45%" }}>
      <Form onSubmit={handleSubmit}>
        <Field>
          <Label htmlFor="label">Label</Label>
          <Input id="label" value={form.label} onChange={update("label")} placeholder="e.g. Mainnet Payment Signers" required />
        </Field>
        <Field>
          <Label htmlFor="walletAddress">Wallet Address</Label>
          <Input id="walletAddress" value={form.walletAddress} onChange={update("walletAddress")} placeholder="0x..." required />
        </Field>
        <Row>
          <Field>
            <Label htmlFor="vaultMount">Vault Mount</Label>
            <Input id="vaultMount" value={form.vaultMount} onChange={update("vaultMount")} placeholder="secret" required />
          </Field>
          <Field>
            <Label htmlFor="vaultPath">Vault Path</Label>
            <Input id="vaultPath" value={form.vaultPath} onChange={update("vaultPath")} placeholder="base/mainnet-signers" required />
          </Field>
        </Row>
        <Row>
          <Field>
            <Label htmlFor="vaultField">Vault Field</Label>
            <Input id="vaultField" value={form.vaultField} onChange={update("vaultField")} placeholder="signers_csv" required />
          </Field>
          <Field>
            <Label htmlFor="activeSigningCount">Active Signing Count</Label>
            <Input
              id="activeSigningCount"
              type="number"
              min={1}
              value={form.activeSigningCount}
              onChange={update("activeSigningCount")}
              required
            />
          </Field>
        </Row>
        <Hint>
          The CSV at this Vault path must already hold at least 4 real, activated signer keys —
          this only registers the pointer, it never writes to Vault.
        </Hint>
        <PrimaryButton buttonStyle={{ width: "100%" }} disabled={isLoading}>
          {isLoading ? "Registering..." : "Register"}
        </PrimaryButton>
      </Form>
    </Modal>
  );
};

export default RegisterManagedSecretModal;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 18px;
`;
const Row = styled.div`
  display: flex;
  gap: 16px;

  > * {
    flex: 1;
  }
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
