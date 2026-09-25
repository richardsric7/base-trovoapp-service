"use client";
import React, { useState } from "react";
import styled from "styled-components";
import Calender from "@/app/(dashboard)/admin-users/components/Calender";
import calenderIcon from "@/assets/images/calender.svg";
import { DropdownSelect } from "@/components";
import SecondaryButton from "@/components/SecondaryButton";
import PrimaryButton from "@/components/PrimaryButton";
import Image from "next/image";
import CustomFilter from "@/components/CustomFilter";
import { DatePicker } from "antd";
const P2pFilter = () => {
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
    ActivationStatus: ["Active", "Suspended"],
    MerchantStatus: ["Merchant", "Non-Merchant"],
  };

  return (
    <div>
      {" "}
      <>
        <CustomFilter
          position={{ top: "270px", right: "30px" }}
          open={isFilterOpen}
          onClose={setIsFilterOpen}
        >
          <FilterContent>
            {/* Date Selection */}
            <div>
              <SelectText>Registration Date</SelectText>

              <StyledDatePicker
                format="YYYY-MM-DD"
                style={{ width: "100%" }}
                placeholder="DD-MM-YYYY"
                size="large"
              />
            </div>
            {/* Activation Filter */}
            <DropdownSelect
              options={filterCriteria.ActivationStatus}
              labelText="Activation Status"
              placeholder="All"
              value={selectValues["ActivationStatus"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, ActivationStatus: item })
              }
            />

            {/* Status Filter */}
            <DropdownSelect
              options={filterCriteria.MerchantStatus}
              placeholder="All"
              labelText="Merchant Status"
              value={selectValues["MerchantStatus"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, MerchantStatus: item })
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
    </div>
  );
};

export default P2pFilter;

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
