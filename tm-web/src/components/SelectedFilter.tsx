import React from "react";
import styled from "styled-components";
import { FaX } from "react-icons/fa6";

interface SelectedFiltersProps {
  selectedFilters: Record<string, string>;
  handleRemoveFilter: (key: string) => void;
  valueToLabel?: Record<string, string>;
}

const SelectedFiltersComponent: React.FC<SelectedFiltersProps> = ({
  selectedFilters,
  handleRemoveFilter,
  valueToLabel,
}) => {
  return (
    <SelectedFiltersContainer>
      {Object.entries(selectedFilters).length > 0 && (
        <StyledSearch>
          {Object.entries(selectedFilters).map(([key, value]) => {
            // For "type", show label; for others, show as-is
            const displayValue =
              key === "type" && valueToLabel?.[value]
                ? valueToLabel[value]
                : value;

            return (
              <SingleSearch key={key}>
                <KeyText>{key}:</KeyText>
                <ValueText>{displayValue}</ValueText>
                <RemoveButton
                  onClick={(e: React.MouseEvent) => {
                    e.preventDefault();
                    handleRemoveFilter(key);
                  }}
                  aria-label={`Remove ${key} filter`}
                >
                  <FaX />
                </RemoveButton>
              </SingleSearch>
            );
          })}
        </StyledSearch>
      )}
    </SelectedFiltersContainer>
  );
};

const SelectedFiltersContainer = styled.div`
  margin-top: 1rem;
`;

const StyledSearch = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;
const SingleSearch = styled.div`
  border: 1px solid #007cdf;
  background-color: #acd1ef;
  font-size: 12px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: center;
  border-radius: 16px;
  padding: 8px 16px;
  display: inline-flex;
  align-items: center;
  gap: 6px;

  transition: all 0.2s ease-in-out;
`;

const KeyText = styled.span`
  color: #4f4f4f;
  font-weight: 500;
`;

const ValueText = styled.span`
  color: #007cdf;
`;
const RemoveButton = styled.button`
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  padding: 0;
  margin-left: 4px;
  cursor: pointer;
  color: #007cdf;
  width: 10px;
  height: 10px;
`;

export default SelectedFiltersComponent;
