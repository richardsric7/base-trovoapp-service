import CustomFilter from "@/components/CustomFilter";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import React, { useState } from "react";
import styled from "styled-components";
import dayjs from "dayjs";
import { DatePicker } from "antd";
import type { DatePickerProps } from "antd";
interface TradesFilterProps {
  createdAt: string | undefined;
  setCreatedAt: (val: string) => void;
}

const TradesFilter: React.FC<TradesFilterProps> = ({
  createdAt,
  setCreatedAt,
}) => {
  const [openFilter, setOpenFilter] = useState<boolean>(false);

  const [selectedDate, setSelectedDate] = useState<string>("");

  const handleDateChange: DatePickerProps["onChange"] = (date, dateString) => {
    if (typeof dateString === "string") {
      setSelectedDate(dateString);
    }
  };

  const handleApply = () => {
    setCreatedAt(
      selectedDate ? dayjs(selectedDate, "DD-MM-YYYY").format("YYYY-MM-DD") : ""
    );

    console.log("filter applied");
    setOpenFilter(false);
  };

  const handleReset = () => {
    setSelectedDate("");
    setCreatedAt("");

    setOpenFilter(false);
  };
  return (
    <CustomFilter
      position={{ top: "280px", right: "12px" }}
      open={openFilter}
      onClose={setOpenFilter}
    >
      {/* Date Selection */}
      <div>
        <Label>Traded on</Label>

        <DatePicker
          onChange={handleDateChange}
          format="DD-MM-YYYY"
          value={selectedDate ? dayjs(selectedDate, "DD-MM-YYYY") : null}
          placeholder="dd-mm-yy"
          style={{ width: "100%" }}
        />
      </div>
      {/* Apply and Reset Buttons */}
      <ButtonContainer>
        <SecondaryButton
          onClick={handleReset}
          buttonStyle={{
            padding: " 8px 16px",
          }}
        >
          Reset
        </SecondaryButton>
        <PrimaryButton
          onClick={handleApply}
          buttonStyle={{
            padding: " 8px 16px",
          }}
        >
          Apply
        </PrimaryButton>
      </ButtonContainer>
    </CustomFilter>
  );
};

export default TradesFilter;

const ButtonContainer = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 20px;
`;
const Label = styled.p`
  font-size: 14px;
  color: #828282;
  font-weight: 500;
  padding: 8px 0;
  margin-bottom: 4px;
`;
