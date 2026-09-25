"use client";

import { useState } from "react";
import { Modal } from "antd";
import { FiX, FiAlertCircle } from "react-icons/fi";
import styled from "styled-components";
import { useRejectTrusteeDueDiligenceMutation } from "@/redux/api/trustees";

interface Props {
  open: boolean;
  onClose: () => void;
  assetId: string;
  assetCode?: string;
  onSuccess?: () => void;
}

export const RejectDueDiligenceModal = ({
  open,
  onClose,
  assetId,
  assetCode,
  onSuccess,
}: Props) => {
  const [reason, setReason] = useState("");
  const [rejectDueDiligence, { isLoading }] =
    useRejectTrusteeDueDiligenceMutation();

  const handleReject = async () => {
    if (!reason.trim()) return;
    try {
      await rejectDueDiligence({
        assetId,
        payload: { reason },
      }).unwrap();
      setReason("");
      onSuccess?.();
      onClose();
    } catch (err) {
      console.error("Reject due diligence failed:", err);
    }
  };

  const handleClose = () => {
    setReason("");
    onClose();
  };

  return (
    <Modal
      open={open}
      onCancel={handleClose}
      footer={null}
      centered
      width={480}
      closeIcon={<FiX size={20} />}
    >
      <Wrapper>
        {/* Icon */}
        <IconWrapper>
          <FiAlertCircle size={28} color="#BE3800" />
        </IconWrapper>

        <Title>Reject Due Diligence</Title>
        <Description>
          You are about to reject the due diligence for asset{" "}
          {assetCode ? <strong>{assetCode}</strong> : "this asset"}. Please
          provide a clear reason before proceeding.
        </Description>

        <Field>
          <Label>Reason for Rejection</Label>
          <TextArea
            placeholder="Describe why you are rejecting this due diligence…"
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            rows={5}
          />
          {!reason.trim() && reason.length > 0 && (
            <ErrorText>Reason cannot be empty.</ErrorText>
          )}
        </Field>

        <ActionRow>
          <CancelBtn onClick={handleClose} disabled={isLoading}>
            Cancel
          </CancelBtn>
          <RejectBtn
            onClick={handleReject}
            disabled={!reason.trim() || isLoading}
          >
            {isLoading ? "Rejecting…" : "Reject Asset"}
          </RejectBtn>
        </ActionRow>
      </Wrapper>
    </Modal>
  );
};



const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 4px 0 8px;
`;

const IconWrapper = styled.div`
  width: 52px;
  height: 52px;
  background: #fff1ee;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const Title = styled.h2`
  font-size: 18px;
  font-weight: 600;
  color: #00225a;
  margin: 0;
`;

const Description = styled.p`
  font-size: 14px;
  color: #667085;
  line-height: 1.6;
`;

const Field = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
`;

const Label = styled.label`
  font-size: 13px;
  font-weight: 500;
  color: #00225a;
`;

const TextArea = styled.textarea`
  width: 100%;
  border-radius: 10px;
  border: 1px solid #e5e7eb;
  padding: 12px;
  resize: none;
  font-size: 14px;
  font-family: inherit;
  color: #344054;
  outline: none;
  transition: border-color 0.2s;

  &:focus {
    border-color: #007cdf;
  }

  &::placeholder {
    color: #9ca3af;
  }
`;

const ErrorText = styled.span`
  font-size: 12px;
  color: #be3800;
`;

const ActionRow = styled.div`
  display: flex;
  gap: 12px;
  margin-top: 4px;
`;

const CancelBtn = styled.button`
  flex: 1;
  padding: 12px;
  border-radius: 10px;
  border: 1px solid #e5e7eb;
  background: #fff;
  color: #344054;
  font-size: 14px;
  font-family: inherit;
  cursor: pointer;
  transition: background 0.15s;

  &:hover:not(:disabled) {
    background: #f9fafb;
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`;

const RejectBtn = styled.button`
  flex: 1;
  padding: 12px;
  border-radius: 10px;
  border: none;
  background: #be3800;
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  transition: background 0.15s, opacity 0.15s;

  &:hover:not(:disabled) {
    background: #a83000;
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`;
