"use client";

import React from "react";
import styled from "styled-components";
import { Modal } from "@/components/CustomModal";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import { showErrorToast, showSuccessToast } from "@/components";
import { useDeleteAssignmentMutation } from "@/redux/api/vaultSigner/api";
import type { IAssignment } from "@/redux/api/vaultSigner/interface";
import { getErrorMessage } from "@/app/vault-signer/components/MySignerSecrets";

interface DeleteAssignmentConfirmProps {
  managedSecretId: string;
  assignment: IAssignment;
  isOpen: boolean;
  onClose: () => void;
}

// The one destructive admin action in this feature with real on-chain consequences (backend
// plan Section 5e) — for a claimed position, this can remove a live signer from the wallet.
// 409 (would drop the CSV below the 4-entry floor) and 502 (on-chain/Vault failure, reverted)
// are both real, expected outcomes here, not bugs — surfaced via the same error-message
// precedence as everywhere else in this feature.
const DeleteAssignmentConfirm: React.FC<DeleteAssignmentConfirmProps> = ({
  managedSecretId,
  assignment,
  isOpen,
  onClose,
}) => {
  const [deleteAssignment, { isLoading }] = useDeleteAssignmentMutation();

  const handleDelete = async () => {
    try {
      const result = await deleteAssignment({
        managedSecretId,
        assignmentId: assignment.id,
      }).unwrap();
      const info = result.data;
      if (!info.claimed) {
        showSuccessToast("Assignment removed");
      } else {
        const onChain = info.onChainRemoval;
        showSuccessToast(
          onChain?.attempted
            ? `Assignment removed — signer removed on-chain (${onChain.stellarTxStatus}, tx ${onChain.stellarTxHash?.slice(0, 10)}...)`
            : `Assignment removed — position ${info.deletedPosition} renumbered, no on-chain signer to remove`,
        );
      }
      onClose();
    } catch (error) {
      showErrorToast(getErrorMessage(error, "Unable to remove assignment"));
    }
  };

  return (
    <Modal title="Remove Assignment" isOpen={isOpen} onClose={onClose} style={{ width: "35%" }}>
      <Body>
        <Text>
          {assignment.position === null ? (
            <>This slot has never claimed a position — this removes an empty assignment row only.</>
          ) : (
            <>
              This slot currently holds <strong>Signer Position {assignment.position}</strong>.
              If it holds an active, valid signer key, this may remove it from the wallet
              on-chain before the assignment is deleted.
            </>
          )}
        </Text>
        <ButtonRow>
          <SecondaryButton buttonStyle={{ width: "100%" }} onClick={onClose}>
            Cancel
          </SecondaryButton>
          <PrimaryButton
            buttonStyle={{ width: "100%", backgroundColor: "#d92d20" }}
            disabled={isLoading}
            onClick={handleDelete}
          >
            {isLoading ? "Removing..." : "Remove"}
          </PrimaryButton>
        </ButtonRow>
      </Body>
    </Modal>
  );
};

export default DeleteAssignmentConfirm;

const Body = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
`;
const Text = styled.p`
  color: #667085;
  font-size: 14px;
  line-height: 1.55;
  margin: 0;
`;
const ButtonRow = styled.div`
  display: flex;
  gap: 12px;
`;
