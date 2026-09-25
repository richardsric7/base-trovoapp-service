import React, { useState } from "react";
import styled from "styled-components";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";
import { format, differenceInDays } from "date-fns";

import { FaAngleDown } from "react-icons/fa";

export default function DateRangePicker() {
  const [startDate, setStartDate] = useState<Date | undefined>(undefined);
  const [endDate, setEndDate] = useState<Date | undefined>(undefined);
  const [daysDifference, setDaysDifference] = useState<number | null>(null);

  const handleDateChange = (dates: [Date | null, Date | null]) => {
    const [start, end] = dates;
    setStartDate(start ?? undefined);
    setEndDate(end ?? undefined);

    if (start && end) {
      const days = differenceInDays(end, start);
      setDaysDifference(days);
    }
  };

  const formatDate = (date: Date | undefined) => {
    return date ? format(date, "d MMM, yyyy") : "";
  };

  return (
    <>
      <FormInputContent>
        <Input
          type="text"
          placeholder="0"
          value={daysDifference !== null ? `${daysDifference} days` : ""}
          readOnly
        />
        <div>
          <DatePickerWrapper>
            <DatePicker
              selected={startDate}
              onChange={handleDateChange}
              startDate={startDate}
              endDate={endDate}
              selectsRange
              placeholderText="Days"
              dateFormat="d MMM, yyyy"
              customInput={<Input />}
            />
            <FaAngleDown />
          </DatePickerWrapper>
        </div>
      </FormInputContent>
      {startDate && endDate && (
        <DateText>
          {formatDate(startDate)} - {formatDate(endDate)}
        </DateText>
      )}
    </>
  );
}

const FormInputContent = styled.div`
  display: flex;
  align-items: center;
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  justify-content: space-between;
  gap: 10px;
`;

const Input = styled.input`
  padding: 8px;
  font-size: 16px;
  border: none;
  width: 100%;
`;

const DatePickerWrapper = styled.div`
  display: flex;
  align-items: center;

  background-color: #f2f6f9;

  padding: 8px;
  cursor: pointer;

  .react-datepicker-wrapper {
    width: auto;
    flex-grow: 1;
    cursor: pointer;
  }

  .react-datepicker__input-container input {
    width: 100%;
    font-size: 12px;
    outline: none;
    border: none;
    background: transparent;
  }
`;

const DateInput = styled(Input)`
  background-color: transparent;
`;

const DateText = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: right;
  color: #00225a;
`;
