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
import type { DatePickerProps } from "antd";
import dayjs from "dayjs";
import { GetFiatPaymentsParams } from "@/redux/api/payment";

interface CustomerFilterProps {
  filters: Partial<GetFiatPaymentsParams>;
  onApply: (filters: Partial<GetFiatPaymentsParams>) => void;
}

const CustomerFilter: React.FC<CustomerFilterProps> = ({
  filters,
  onApply,
}) => {
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [selectValues, setSelectValues] = useState<any>({
    KycStatus:
      filters.status === "COMPLETED"
        ? "Completed"
        : filters.status === "PENDING"
          ? "Pending"
          : filters.status === "FAILED"
            ? "Failed"
            : "",
    date: filters.created_at || "",
  });

  const handleDateChange = (date: any, dateString: string | string[]) => {
    if (typeof dateString === "string") {
      setSelectValues({ ...selectValues, date: dateString });
    }
  };

  const handleFilterApplication = () => {
    const newFilters: Partial<GetFiatPaymentsParams> = {};

    if (selectValues.KycStatus) {
      const statusMap: Record<string, any> = {
        Verified: "VERIFIED",
        Pending: "PENDING",
        Rejected: "REJECTED",
        Unverified: "UNVERIFIED",
        Active: "ACTIVE",
        Inactive: "INACTIVE",
      };
      newFilters.status = statusMap[selectValues.KycStatus];
    }

    if (selectValues.date) {
      newFilters.created_at = dayjs(selectValues.date, "YYYY-MM-DD").format(
        "YYYY-MM-DD",
      );
    }

    onApply(newFilters);
    setIsFilterOpen(false);
  };

  const handleFilterReset = () => {
    setSelectValues({});
    onApply({});
    setIsFilterOpen(false);
  };

  const filterCriteria = {
    KycStatus: [
      "Verified",
      "Pending",
      "Rejected",
      "Unverified",
      "Active",
      "Inactive",
    ],
  };

  return (
    <div>
      {" "}
      <>
        <CustomFilter
          position={{ top: "240px", right: "30px" }}
          open={isFilterOpen}
          onClose={setIsFilterOpen}
        >
          <FilterContent>
            {/* Date Selection */}
            <div>
              <SelectText> Registration Date</SelectText>

              <StyledDatePicker
                value={
                  selectValues.date
                    ? dayjs(selectValues.date, "YYYY-MM-DD")
                    : null
                }
                onChange={handleDateChange}
                format="YYYY-MM-DD"
                style={{ width: "100%" }}
                placeholder="YYYY-MM-DD"
                size="large"
              />
            </div>

            {/* Status Filter */}
            <DropdownSelect
              options={filterCriteria.KycStatus}
              placeholder="All"
              labelText="Kyc Status"
              value={selectValues["KycStatus"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, KycStatus: item })
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

export default CustomerFilter;

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
