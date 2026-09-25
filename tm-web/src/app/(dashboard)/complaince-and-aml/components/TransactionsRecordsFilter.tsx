"use client";

import styled from "styled-components";

import { DropdownSelect } from "@/components";
import CustomFilter from "@/components/CustomFilter";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import { DatePicker } from "antd";
import { useState } from "react";

const TransactionsRecordsFilter = () => {
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [selectValues, setSelectValues] = useState<any>({});

  const [price, setPrice] = useState("");

  const handleFilterApplication = () => {
    // setSelectedFilters({ ...selectValues, date: selectedDate });
    setIsFilterOpen(false);
  };

  const handleFilterReset = () => {
    setSelectValues({});
    console.log("Filters Reset");
  };

  const filterCriteria = {
    currency: ["N", "$"],

    transactionType: ["Recieved", "sent"],
    amlStatus: ["Successful", "Pending", "Failed"],
  };
  return (
    <>
      <CustomFilter
        position={{ top: "270px", right: "30px" }}
        open={isFilterOpen}
        onClose={setIsFilterOpen}
      >
        <FilterContent>
          <InputWrapper>
            <Label>Price</Label>
            <InputField
              name=""
              placeholder="0"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
            />
          </InputWrapper>

          <DropdownSelect
            options={filterCriteria.currency}
            labelText="Currency"
            placeholder="All"
            value={selectValues["currency"] || ""}
            onSelect={(item) =>
              setSelectValues({ ...selectValues, currency: item })
            }
          />

          <DropdownSelect
            options={filterCriteria.transactionType}
            labelText="Transaction Type"
            placeholder="All"
            value={selectValues["transactionType"] || ""}
            onSelect={(item) =>
              setSelectValues({ ...selectValues, transactionType: item })
            }
          />
          {/* Date Selection */}
          <div>
            <SelectText>Date</SelectText>

            <StyledDatePicker
              format="YYYY-MM-DD"
              style={{ width: "100%" }}
              placeholder="DD-MM-YYYY"
              size="large"
            />
          </div>

          {/* Status Filter */}
          <DropdownSelect
            options={filterCriteria.amlStatus}
            placeholder="All"
            labelText="AML Status"
            value={selectValues["status"] || ""}
            onSelect={(item) =>
              setSelectValues({ ...selectValues, amlStatus: item })
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

export default TransactionsRecordsFilter;

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

const InputField = styled.input`
  font-size: 14px;
  width: 100%;
  outline: none;
  font-family: inherit;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
`;

const InputWrapper = styled.div``;

const Label = styled.p`
  font-weight: 500;
  font-size: 14px;
  color: #828282;
  line-height: 24px;
  letter-spacing: 0.1px;
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
