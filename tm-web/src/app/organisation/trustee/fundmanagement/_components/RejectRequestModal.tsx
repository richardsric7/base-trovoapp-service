"use client";

import PrimaryButton from "@/components/PrimaryButton";
import { useRejectTrusteeFundReleaseMutation } from "@/redux/api/trustees";
import { IFundReleaseRecord } from "@/redux/api/trustees/interface";
import { message, Modal } from "antd";
import { useState } from "react";
import styled from "styled-components";

interface RejectRequestModalProps {
  open: boolean;
  onClose: () => void;
  record?: IFundReleaseRecord;
}

export const RejectRequestModal = ({
  open,
  onClose,
  record,
}: RejectRequestModalProps) => {
  const [reason, setReason] = useState("");
  const [reject, { isLoading }] = useRejectTrusteeFundReleaseMutation();

  const handleReject = async () => {
    if (!record?.id) return;
    if (!reason.trim()) {
      message.warning("Please provide a reason for rejecting this request.");
      return;
    }
    try {
      await reject({ requestId: record.id, payload: { reason } }).unwrap();
      message.success("Fund release request rejected.");
      setReason("");
      onClose();
    } catch {
      message.error("Failed to reject fund release. Please try again.");
    }
  };

  return (
    <Modal open={open} onCancel={onClose} footer={null} centered width={480}>
      <Wrapper>
        <Title>Reject Request</Title>

        <Description>
          You are about to reject the release of funds from the custody account.
          Please enter a reason for rejecting this request before proceeding.
        </Description>

        <Label>Reason</Label>
        <TextArea
          placeholder="Add a reason for rejecting this request"
          value={reason}
          onChange={(e) => setReason(e.target.value)}
        />

        <PrimaryButton
          buttonStyle={{ backgroundColor: "#BE3800" }}
          onClick={handleReject}
          disabled={isLoading}
        >
          {isLoading ? "Rejecting…" : "Reject Fund Release"}
        </PrimaryButton>
      </Wrapper>
    </Modal>
  );
};

const TextArea = styled.textarea`
  width: 100%;
  height: 120px;
  border-radius: 10px;
  border: 1px solid #e5e7eb;
  padding: 10px;
  resize: none;
  font-size: 14px;
  color: #00225a;
  outline: none;

  &:focus {
    border-color: #007cdf;
  }
`;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

const Title = styled.h2`
  font-size: 18px;
  font-weight: 600;
  color: #00225a;
`;

const Description = styled.p`
  font-size: 14px;
  color: #00225a;
`;

const Label = styled.p`
  font-size: 14px;
  color: #00225a;
  font-weight: 500;
  margin-bottom: 2px;
`;
