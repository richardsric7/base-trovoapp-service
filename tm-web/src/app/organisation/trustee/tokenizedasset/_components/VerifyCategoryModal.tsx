"use client";

import React, { useState, useEffect } from "react";
import { Modal, Select, message } from "antd";
import { FiX, FiCheckCircle } from "react-icons/fi";
import styled from "styled-components";
import { useVerifyTrusteeDueDiligenceCategoryMutation } from "@/redux/api/trustees";

interface Props {
  open: boolean;
  onClose: () => void;
  assetId: string;
  categoryKey: string;
  categoryTitle: string;
  onSuccess?: () => void;
}

const STATUS_OPTIONS = [
  { value: "complete", label: "Complete (Verified)" },
  { value: "pending", label: "Pending" },
  { value: "failed", label: "Failed" },
  { value: "not_applicable", label: "Not Applicable" },
];

export const VerifyCategoryModal = ({
  open,
  onClose,
  assetId,
  categoryKey,
  categoryTitle,
  onSuccess,
}: Props) => {
  const [status, setStatus] = useState<string>("complete");
  const [notes, setNotes] = useState<string>("");

  const [verifyCategory, { isLoading }] =
    useVerifyTrusteeDueDiligenceCategoryMutation();

  useEffect(() => {
    if (open) {
      setStatus("complete");
      setNotes("");
    }
  }, [open]);

  const handleSubmit = async () => {
    if (!assetId || !categoryKey) return;
    try {
      await verifyCategory({
        assetId,
        payload: {
          category: categoryKey,
          status,
          notes: notes.trim(),
        },
      }).unwrap();

      message.success(`Updated ${categoryTitle} status to ${status.replace("_", " ")}.`);
      onSuccess?.();
      onClose();
    } catch (err) {
      console.error("Failed to verify category:", err);
      message.error("Failed to update category status. Please try again.");
    }
  };

  return (
    <Modal
      open={open}
      onCancel={onClose}
      footer={null}
      centered
      width={480}
      closeIcon={<FiX size={20} />}
    >
      <Wrapper>
        <IconWrapper>
          <FiCheckCircle size={28} color="#007cdf" />
        </IconWrapper>

        <Title>Verify Category</Title>
        <Subtitle>{categoryTitle}</Subtitle>
        <Description>
          Atomically update every due diligence item in this category to the selected status.
        </Description>

        <Field>
          <Label>Category Status</Label>
          <Select
            value={status}
            onChange={(val) => setStatus(val)}
            options={STATUS_OPTIONS}
            style={{ width: "100%" }}
            size="large"
          />
        </Field>

        <Field>
          <Label>Notes (Optional)</Label>
          <TextArea
            placeholder="Add any verification notes or observations…"
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
            rows={4}
          />
        </Field>

        <ActionRow>
          <CancelBtn onClick={onClose} disabled={isLoading}>
            Cancel
          </CancelBtn>
          <SubmitBtn onClick={handleSubmit} disabled={isLoading}>
            {isLoading ? "Updating…" : "Confirm Update"}
          </SubmitBtn>
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
  background: #e6f2ff;
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

const Subtitle = styled.div`
  font-size: 14px;
  font-weight: 600;
  color: #007cdf;
  margin-top: -10px;
`;

const Description = styled.p`
  font-size: 13px;
  color: #667085;
  line-height: 1.5;
  margin: 0;
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
  outline: none;
  transition: border-color 0.2s;

  &:focus {
    border-color: #007cdf;
  }
`;

const ActionRow = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 8px;
`;

const CancelBtn = styled.button`
  padding: 10px 18px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  background: #fff;
  color: #344054;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
`;

const SubmitBtn = styled.button`
  padding: 10px 18px;
  border-radius: 8px;
  border: none;
  background: #007cdf;
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
`;
