"use client";

import { DropdownSelect } from "@/components";
import CustomFilter from "@/components/CustomFilter";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import { DatePicker } from "antd";
import { useState } from "react";

import styled from "styled-components";

const SupportFilter = () => {
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [selectValues, setSelectValues] = useState<any>({});

  const handleFilterApplication = () => {
    // setSelectedFilters({ ...selectValues, date: selectedDate });
    setIsFilterOpen(false);
  };

  const handleFilterReset = () => {
    setSelectValues({});
    console.log("Filters Reset");
  };

  const filterCriteria = {
    issueType: ["Billing", "Suspended"],
    urgency: ["High", "Medium", "Normal"],
    status: ["Successful", "Pending", "Failed"],
  };
  return (
    <>
      <CustomFilter
        position={{ top: "270px", right: "30px" }}
        open={isFilterOpen}
        onClose={setIsFilterOpen}
      >
        <FilterContent>
          {/* Date Selection */}
          <div>
            <SelectText> Reported On</SelectText>

            <StyledDatePicker
              format="YYYY-MM-DD"
              style={{ width: "100%" }}
              placeholder="DD-MM-YYYY"
              size="large"
            />
          </div>

          {/* Date Selection */}
          <div>
            <SelectText> Last Updated On</SelectText>

            <StyledDatePicker
              format="YYYY-MM-DD"
              style={{ width: "100%" }}
              placeholder="DD-MM-YYYY"
              size="large"
            />
          </div>

          <DropdownSelect
            options={filterCriteria.issueType}
            labelText="Issue Type"
            placeholder="All"
            value={selectValues["issueType"] || ""}
            onSelect={(item) =>
              setSelectValues({ ...selectValues, issueType: item })
            }
          />

          <DropdownSelect
            options={filterCriteria.urgency}
            labelText="urgency"
            placeholder="All"
            value={selectValues["urgency"] || ""}
            onSelect={(item) =>
              setSelectValues({ ...selectValues, urgency: item })
            }
          />

          {/* Status Filter */}
          <DropdownSelect
            options={filterCriteria.status}
            placeholder="All"
            labelText=" Status"
            value={selectValues["status"] || ""}
            onSelect={(item) =>
              setSelectValues({ ...selectValues, status: item })
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

export default SupportFilter;
const FilterContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 15px;
  marign-top: 15px;
`;

const SelectText = styled.p`
  font-size: 14px;
  color: #828282;
  padding: 10px 0;
`;

const ButtonContainer = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 20px;
`;

const StyledDatePicker = styled(DatePicker)`
  width: 100%;
  padding: 0 4px !important;
  .ant-picker-input > input {
    padding: 10px;
    font-size: 14px;
    color: #00225a;
    font-family: inherit;
    border-radius: 10px;
  }
  &.ant-picker,
  &.ant-picker-focused {
    border: 1px solid #e0e0e0;
    border-radius: 10px;
    box-shadow: none;
    outline: none;
    transition: border 0.2s;
  }
  &:hover,
  &.ant-picker-focused {
    border: 1.5px solid #e0e0e0;
  }
`;
