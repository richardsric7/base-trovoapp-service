import React, { useState } from "react";
import styled from "styled-components";
import { Modal, showErrorToast, showSuccessToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import { useLiftUserSuspensionMutation } from "@/redux/api/users";

interface RevokeUserProps {
  revokeSuspension: boolean;
  setRevokeSuspension: React.Dispatch<React.SetStateAction<boolean>>;
  userEmail: string;
  username: string;
}

const RevokeUserSuspension: React.FC<RevokeUserProps> = ({
  revokeSuspension,
  setRevokeSuspension,
  userEmail,
  username,
}) => {
  const [reason, setReason] = useState<string>("");
  const [liftUserSuspension, { isLoading }] = useLiftUserSuspensionMutation();

  const handleLiftSuspension = async () => {
    if (!reason.trim()) {
      showErrorToast("A reason for lifting the suspension is required.");
      return;
    }
    try {
      await liftUserSuspension({ email: userEmail, reason: reason.trim() }).unwrap();
      showSuccessToast("Suspension has been lifted.");
      setReason("");
      setRevokeSuspension(false);
    } catch (error: any) {
      showErrorToast(error?.data?.error || "Failed to lift suspension.");
    }
  };

  return (
    <Modal
      title="Lift suspension"
      isOpen={revokeSuspension}
      onClose={() => setRevokeSuspension(false)}
    >
      <Text>
        Lift the suspension on <StyledSpan>{username}</StyledSpan>? They will
        immediately be able to use their wallet(s) again.
      </Text>
      <FormGroup>
        <FormLabel>Reason for lifting the suspension *</FormLabel>
        <TextArea
          placeholder="Describe why this suspension is being lifted..."
          value={reason}
          onChange={(e) => setReason(e.target.value)}
        />
      </FormGroup>
      <PrimaryButton onClick={handleLiftSuspension} disabled={isLoading}>
        {isLoading ? "Lifting suspension..." : "Yes, lift suspension"}
      </PrimaryButton>
    </Modal>
  );
};

export default RevokeUserSuspension;

const Text = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  letter-spacing: 0.10000000149011612px;
  text-align: center;
  color: #00225a;
`;

const StyledSpan = styled.span`
  font-size: 14px;
  font-weight: 600;
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
