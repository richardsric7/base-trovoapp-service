import React, { useState } from "react";
import styled from "styled-components";
import { Modal, showErrorToast, showSuccessToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import { useSuspendUserMutation } from "@/redux/api/users";

interface SuspendUserModalProps {
  isOpen: boolean;
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  userEmail: string;
}

export const SuspendUserModal: React.FC<SuspendUserModalProps> = ({
  isOpen,
  setIsOpen,
  userEmail,
}) => {
  const [reason, setReason] = useState<string>("");
  const [suspendUser, { isLoading }] = useSuspendUserMutation();

  const handleSuspendUser = async () => {
    if (!reason.trim()) {
      showErrorToast("A reason for the suspension is required.");
      return;
    }
    try {
      await suspendUser({ email: userEmail, reason: reason.trim() }).unwrap();
      showSuccessToast("User has been suspended.");
      setReason("");
      setIsOpen(false);
    } catch (error: any) {
      showErrorToast(error?.data?.error || "Failed to suspend user.");
    }
  };

  return (
    <Modal title="Suspend user" isOpen={isOpen} onClose={() => setIsOpen(false)}>
      <SuspendCard>
        <FormGroup>
          <FormLabel>Reason for suspension *</FormLabel>
          <TextArea
            placeholder="Describe why this user is being suspended..."
            value={reason}
            onChange={(e) => setReason(e.target.value)}
          />
          <Hint>
            While suspended, this user&apos;s wallet(s) cannot be used in any
            transaction (payments, swaps, P2P, shared access, etc).
          </Hint>
        </FormGroup>
      </SuspendCard>
      <PrimaryButton onClick={handleSuspendUser} disabled={isLoading}>
        {isLoading ? "Suspending..." : "Suspend user"}
      </PrimaryButton>
    </Modal>
  );
};
export default SuspendUserModal;

const SuspendCard = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
`;

const FormGroup = styled.div`
  margin-bottom: 16px;
`;

const FormLabel = styled.label`
  margin-bottom: 8px;
  font-style: normal;
  font-weight: 500;
  font-size: 14px;
  line-height: 24px;
  color: #00225a;
  display: block;
`;

const TextArea = styled.textarea`
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  padding: 8px 16px;
  width: 96%;
  height: 100px;
  resize: vertical;
  font-family: inherit;
  font-size: 14px;
`;

const Hint = styled.p`
  font-size: 12px;
  color: #828282;
  margin: 8px 0 0 0;
`;
