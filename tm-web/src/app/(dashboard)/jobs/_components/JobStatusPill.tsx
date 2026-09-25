"use client";
import styled from "styled-components";
import { Spin } from "antd";

export type JobStatus = "published" | "disabled";
interface Props {
  status: JobStatus;
  loading?: boolean;
  onPublish: () => void;
  onDisable: () => void;
}

const JobStatusPill = ({ status, loading, onPublish, onDisable }: Props) => {
  return (
    <Wrapper>
      <SectionTitle>Job Status</SectionTitle>
      <Pill>
        <PillButton
          $active={status === "published"}
          onClick={onPublish}
          disabled={loading || status === "published"}
        >
          {loading && status !== "published" ? (
            <Spin size="small" />
          ) : (
            "Published"
          )}
        </PillButton>

        <PillButton
          $active={status === "disabled"}
          onClick={onDisable}
          disabled={loading || status === "disabled"}
        >
          Disabled
        </PillButton>
      </Pill>
    </Wrapper>
  );
};

export default JobStatusPill;
const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-top: 12px;
`;
const SectionTitle = styled.h3`
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #00225a;
`;

const Pill = styled.div`
  display: inline-flex;
  background: #e8ecf4;
  border-radius: 999px;
  padding: 4px;
  width: fit-content;
`;

const PillButton = styled.button<{ $active: boolean }>`
  padding: 8px 24px;
  border-radius: 40px;
  border: none;
  font-size: 14px;
  font-weight: 400;
  cursor: pointer;
  min-width: 96px;
  transition: all 0.2s ease;

  background: ${({ $active }) => ($active ? "#007CDF" : "transparent")};
  color: ${({ $active }) => ($active ? "#ffffff" : "#4F4F4F")};

  &:disabled {
    cursor: not-allowed;
    opacity: 0.7;
  }
`;
