import React from "react";
import styled from "styled-components";
import { ReportRange } from "@/redux/api/p2p";
import { REPORT_RANGE_OPTIONS } from "./reportUtils";

const ReportRangeSelector = ({
  value,
  onChange,
}: {
  value: ReportRange;
  onChange: (range: ReportRange) => void;
}) => {
  return (
    <TabsContainer>
      {REPORT_RANGE_OPTIONS.map((opt) => (
        <TabButton key={opt.value} $active={value === opt.value} onClick={() => onChange(opt.value)}>
          {opt.label}
        </TabButton>
      ))}
    </TabsContainer>
  );
};

export default ReportRangeSelector;

const TabsContainer = styled.div`
  display: flex;
  gap: 4px;
  background-color: #f5f5f5;
  padding: 4px;
  border-radius: 12px;
  width: fit-content;
`;

const TabButton = styled.button<{ $active?: boolean }>`
  padding: 8px 16px;
  border: none;
  background-color: ${(props) => (props.$active ? "white" : "transparent")};
  color: ${(props) => (props.$active ? "#00225a" : "#828282")};
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;

  &:hover {
    background-color: ${(props) => (props.$active ? "white" : "#e0e0e0")};
  }
`;
