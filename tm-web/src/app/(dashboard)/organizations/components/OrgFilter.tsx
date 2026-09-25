"use client";

import { DropdownSelect } from "@/components";
import CustomFilter from "@/components/CustomFilter";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import React, { useEffect, useState } from "react";
import styled from "styled-components";

interface OrgFilterProps {
  selectedFilters: Record<string, string>;
  setSelectedFilters: React.Dispatch<
    React.SetStateAction<Record<string, string>>
  >;
}

const OrgFilter: React.FC<OrgFilterProps> = ({
  selectedFilters,
  setSelectedFilters,
}) => {
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [selectValues, setSelectValues] = useState<any>({});

  useEffect(() => {
    if (isFilterOpen) setSelectValues(selectedFilters);
  }, [isFilterOpen, selectedFilters]);

  const labelToValue: Record<string, string> = {
    "Asset Manager": "ASSET_MANAGER",
    "Asset Custodian": "ASSET_CUSTODIAN",
    "Rating Agency": "RATING_AGENCY",
    Regulator: "REGULATOR",
    "Legal Agency": "LEGAL_AGENCY",
    "Professional Agency": "PROFESSIONAL_AGENCY",
  };
  const valueToLabel: Record<string, string> = Object.fromEntries(
    Object.entries(labelToValue).map(([k, v]) => [v, k])
  );
  const orgTypeOptions = Object.keys(labelToValue);

  const handleFilterApplication = () => {
    setIsFilterOpen(false);
    setSelectedFilters(selectValues);
  };

  const handleFilterReset = () => {
    setSelectValues({});
    setSelectedFilters({});

    console.log("Filters Reset");
  };
  return (
    <>
      <CustomFilter
        position={{ top: "270px", right: "30px" }}
        open={isFilterOpen}
        onClose={setIsFilterOpen}
      >
        <FilterContent>
          <DropdownSelect
            options={orgTypeOptions}
            placeholder="All"
            labelText="Organization Type"
            labelColor="#828282"
            placeholderColor="#00225A"
            // Convert the stored API value back to label for display
            value={valueToLabel[selectValues["type"] || ""] || ""}
            onSelect={(selectedLabel) =>
              setSelectValues((prev: any) => ({
                ...prev,
                type: labelToValue[selectedLabel] || "",
              }))
            }
          />

          {/* Apply and Reset Buttons */}
          <ButtonContainer>
            <SecondaryButton
              onClick={handleFilterReset}
              buttonStyle={{
                padding: " 8px 16px",
              }}
            >
              Reset
            </SecondaryButton>
            <PrimaryButton
              onClick={handleFilterApplication}
              buttonStyle={{
                padding: " 8px 16px",
              }}
            >
              Apply
            </PrimaryButton>
          </ButtonContainer>
        </FilterContent>
      </CustomFilter>
    </>
  );
};

export default OrgFilter;

const FilterContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 15px;
  marign-top: 15px;
`;

const SelectText = styled.p`
  font-size: 14px;
  color: #828282;
`;

const CalenderRow = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  cursor: pointer;
`;

const ButtonContainer = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 20px;
`;
