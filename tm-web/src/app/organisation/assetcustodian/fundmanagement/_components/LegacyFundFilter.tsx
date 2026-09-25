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

interface PaymentFilterProps {
  filters: Partial<GetFiatPaymentsParams>;
  onApply: (filters: Partial<GetFiatPaymentsParams>) => void;
}

const FundFilter: React.FC<PaymentFilterProps> = ({ filters, onApply }) => {
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [selectValues, setSelectValues] = useState<any>({
    paymentType:
      filters.record_type === "invoice"
        ? "Invoice"
        : filters.record_type === "payment"
          ? "Payment"
          : "",

    paymentStatus:
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
    if (selectValues.paymentType === "Invoice")
      newFilters.record_type = "invoice";
    if (selectValues.paymentType === "Payment")
      newFilters.record_type = "payment";

    if (selectValues.paymentStatus) {
      const statusMap: Record<string, any> = {
        Completed: "COMPLETED",
        Pending: "PENDING",
        Failed: "FAILED",
      };
      newFilters.status = statusMap[selectValues.paymentStatus];
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
    paymentType: ["Invoice", "Payment"],
    paymentStatus: ["Completed", "Pending", "Failed"],
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
              options={filterCriteria.paymentType}
              labelText="Record Type"
              placeholder="All"
              value={selectValues["paymentType"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, paymentType: item })
              }
            />
            {/* Date Selection */}
            <div>
              <SelectText> Date</SelectText>

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
              options={filterCriteria.paymentStatus}
              placeholder="All"
              labelText="Payment Status"
              value={selectValues["paymentStatus"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, paymentStatus: item })
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

export default FundFilter;

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
