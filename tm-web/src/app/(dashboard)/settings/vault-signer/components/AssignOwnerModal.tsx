"use client";

import React, { useState } from "react";
import styled from "styled-components";
import { Modal } from "@/components/CustomModal";
import PrimaryButton from "@/components/PrimaryButton";
import { showErrorToast, showSuccessToast } from "@/components";
import {
  useCreateAssignmentMutation,
  useEditAssignmentMutation,
} from "@/redux/api/vaultSigner/api";
import type { IAssignment, OwnerType } from "@/redux/api/vaultSigner/interface";
import { getErrorMessage } from "@/app/vault-signer/components/MySignerSecrets";

interface AssignOwnerModalProps {
  managedSecretId: string;
  isOpen: boolean;
  onClose: () => void;
  // Present only when reassigning an existing assignment (PATCH); absent for a brand new
  // assignment (POST). Either way, no position field is ever collected here — the backend
  // never accepts one at creation or reassignment time, it's claimed by the owner's own
  // first PUT (see the backend plan, Section 5e).
  //
  // (ManagedSecretID, OwnerRefID) is unique on the backend (Section 4/12 of the backend plan)
  // — an owner can hold at most one assignment on a given managed secret. There's no reliable
  // way to pre-check this client-side: the admin types a username/email (`ownerIdentifier`),
  // and only the backend resolves that to the stable `OwnerRefID` the constraint actually keys
  // on (`services.ResolveOwnerRefID`). So this form always submits and lets the server be the
  // judge — both create and reassign surface the backend's `409 owner_already_assigned` via
  // the shared `getErrorMessage` helper below, same as every other error in this feature.
  editingAssignment?: IAssignment;
}

const AssignOwnerModal: React.FC<AssignOwnerModalProps> = ({
  managedSecretId,
  isOpen,
  onClose,
  editingAssignment,
}) => {
  const isEdit = !!editingAssignment;
  const [ownerIdentifier, setOwnerIdentifier] = useState(editingAssignment?.ownerLabel ?? "");
  const [ownerType, setOwnerType] = useState<OwnerType>(editingAssignment?.ownerType ?? "trovo_admin");

  const [createAssignment, { isLoading: isCreating }] = useCreateAssignmentMutation();
  const [editAssignment, { isLoading: isEditing }] = useEditAssignmentMutation();
  const isLoading = isCreating || isEditing;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!ownerIdentifier.trim()) return;

    try {
      if (isEdit) {
        await editAssignment({
          managedSecretId,
          assignmentId: editingAssignment!.id,
          ownerIdentifier: ownerIdentifier.trim(),
          ownerType,
        }).unwrap();
        showSuccessToast(
          "Owner reassigned — the new owner picks up the slot the next time they set a value.",
        );
      } else {
        await createAssignment({
          managedSecretId,
          ownerIdentifier: ownerIdentifier.trim(),
          ownerType,
        }).unwrap();
        showSuccessToast("Assignment created");
      }
      onClose();
    } catch (error) {
      showErrorToast(getErrorMessage(error, "Unable to save assignment"));
    }
  };

  return (
    <Modal
      title={isEdit ? "Reassign Owner" : "Assign Owner"}
      isOpen={isOpen}
      onClose={onClose}
      style={{ width: "40%" }}
    >
      <Form onSubmit={handleSubmit}>
        <Field>
          <Label htmlFor="ownerType">Owner Type</Label>
          <Select
            id="ownerType"
            value={ownerType}
            onChange={(e) => setOwnerType(e.target.value as OwnerType)}
          >
            <option value="trovo_admin">Trovo Admin</option>
            <option value="org_member">Organization Member</option>
          </Select>
        </Field>
        <Field>
          <Label htmlFor="ownerIdentifier">
            {ownerType === "trovo_admin" ? "Username" : "Email"}
          </Label>
          <Input
            id="ownerIdentifier"
            value={ownerIdentifier}
            onChange={(e) => setOwnerIdentifier(e.target.value)}
            placeholder={ownerType === "trovo_admin" ? "e.g. obi" : "e.g. member@acme-trust.com"}
            autoComplete="off"
          />
        </Field>
        {isEdit && (
          <Hint>
            This doesn&apos;t touch Vault or the chain — the current signer key (if any) stays
            exactly as it is until the new owner submits their own value.
          </Hint>
        )}
        <PrimaryButton buttonStyle={{ width: "100%" }} disabled={isLoading}>
          {isLoading ? "Saving..." : isEdit ? "Reassign" : "Assign"}
        </PrimaryButton>
      </Form>
    </Modal>
  );
};

export default AssignOwnerModal;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 18px;
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
const inputStyles = `
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
const Input = styled.input`
  ${inputStyles}
`;
const Select = styled.select`
  ${inputStyles}
  background: #fff;
`;
const Hint = styled.span`
  color: #828282;
  font-size: 12px;
`;
