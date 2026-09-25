"use client";

import { DropdownSelect } from "@/components";
import CustomFilter from "@/components/CustomFilter";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import { Input } from "antd";
import { useState } from "react";
import styled from "styled-components";

export interface JobFilterValues {
  location?: string;
  type?: string;
  experience?: string;
}

const JobFilter = ({
  onFilter,
}: {
  onFilter?: (values: JobFilterValues) => void;
}) => {
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [values, setValues] = useState<JobFilterValues>({});

  const applyFilter = () => {
    onFilter?.(values);
    setIsFilterOpen(false);
  };

  const resetFilter = () => {
    setValues({});
    onFilter?.({});
    setIsFilterOpen(false);
  };

  return (
    <>
      <CustomFilter
        open={isFilterOpen}
        onClose={setIsFilterOpen}
        position={{ top: "360px", right: "30px" }}
      >
        <FilterContent>
          {/* LOCATION */}
          <DropdownSelect
            labelText="Location"
            placeholder="All"
            options={["remote", "onsite", "hybrid"]}
            value={values.location || ""}
            onSelect={(item) => setValues({ ...values, location: item })}
          />

          {/* EXPERIENCE */}
          <div>
            <SelectText>Years of Experience</SelectText>
            <Input
              placeholder="e.g. 3+ years"
              value={values.experience || ""}
              onChange={(e) =>
                setValues({
                  ...values,
                  experience: e.target.value,
                })
              }
            />
          </div>

          {/* WORK TYPE */}
          <DropdownSelect
            labelText="Work Type"
            placeholder="All"
            options={["fulltime", "parttime", "contract"]}
            value={values.type || ""}
            onSelect={(item) => setValues({ ...values, type: item })}
          />

          {/* ACTIONS */}
          <ButtonContainer>
            <SecondaryButton onClick={resetFilter}>Reset</SecondaryButton>
            <PrimaryButton onClick={applyFilter}>Apply</PrimaryButton>
          </ButtonContainer>
        </FilterContent>
      </CustomFilter>
    </>
  );
};

export default JobFilter;

const FilterContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

const SelectText = styled.p`
  font-size: 14px;
  color: #828282;
  margin-bottom: 6px;
`;

const ButtonContainer = styled.div`
  display: flex;
  gap: 12px;
`;
