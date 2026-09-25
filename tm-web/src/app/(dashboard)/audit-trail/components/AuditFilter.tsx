"use client";
import { DropdownSelect } from "@/components";
import CustomFilter from "@/components/CustomFilter";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import { useState } from "react";
import { DatePicker } from "antd";
import styled from "styled-components";

interface AuditFilterProps {
  onApply: (filters: { status?: string; date?: string }) => void;
}

const AuditFilter = ({ onApply }: AuditFilterProps) => {
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [selectValues, setSelectValues] = useState<any>({});
  const [selectedDate, setSelectedDate] = useState<string>("");

  const handleFilterApplication = () => {
    onApply({
      // Reconcile the UI's title-case with the API's canonical lowercase.
      status: selectValues["status"]
        ? String(selectValues["status"]).toLowerCase()
        : undefined,
      date: selectedDate || undefined,
    });
    setIsFilterOpen(false);
  };

  const handleFilterReset = () => {
    setSelectValues({});
    setSelectedDate("");
    onApply({ status: undefined, date: undefined });
  };

  const filterCriteria = {
    user: ["Active", "Suspended"],
    status: ["Successful", "Pending", "Failed"],
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
            {/* PaymentTypeFilter */}
            <DropdownSelect
              options={filterCriteria.user}
              labelText="User"
              placeholder="All"
              value={selectValues["user"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, user: item })
              }
            />
            {/* Date Selection */}
            <div>
              <SelectText> Date</SelectText>

              <StyledDatePicker
                format="YYYY-MM-DD"
                style={{ width: "100%" }}
                placeholder="YYYY-MM-DD"
                size="large"
                onChange={(_date: unknown, dateString: string | string[]) =>
                  setSelectedDate(
                    Array.isArray(dateString) ? dateString[0] ?? "" : dateString,
                  )
                }
              />
            </div>

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
    </div>
  );
};

export default AuditFilter;

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
