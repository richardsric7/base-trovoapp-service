"use client";

import { DropdownSelect } from "@/components";
import React, { useState } from "react";
import styled from "styled-components";

const Settings = () => {
  const [selectValues, setSelectValues] = useState<any>({});
  const [dateValue, setDateValue] = useState("");
  const [timeValue, setTimeValue] = useState("");

  const RecipientsType = ["Update", "Upgrade", "Feature"];

  return (
    <>
      <Title>Recipients</Title>
      <DropdownSelect
        options={RecipientsType}
        placeholder="All Users"
        labelText="Select Recipient user group"
        labelColor="#828282"
        placeholderColor="#00225A"
        value={selectValues["RecipientsType"] || ""}
        onSelect={(item) =>
          setSelectValues({ ...selectValues, RecipientsType: item })
        }
      />

      <CreateGroup>Create new group</CreateGroup>

      <Timing>
        <Title>Timing</Title>
        <DateLabel>Schedule Date</DateLabel>
        <InputWrapper>
          <StyledInput
            type="date"
            value={dateValue}
            onChange={(e) => setDateValue(e.target.value)}
          />
          {dateValue === "" && <Placeholder>Select Date</Placeholder>}
        </InputWrapper>

        <Label>Schedule Time</Label>
        <InputWrapper>
          <StyledInput
            type="time"
            value={timeValue}
            onChange={(e) => setTimeValue(e.target.value)}
          />
          {timeValue === "" && (
            <Placeholder>Enter Notification Time</Placeholder>
          )}
        </InputWrapper>
      </Timing>
    </>
  );
};

export default Settings;

const Title = styled.h1`
  font-size: 20px;
  font-weight: 700;
  line-height: 28px;
  text-align: left;
  color: #00225a;
  margin: 0;
`;

const CreateGroup = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 24px;
  color: #007cdf;
  text-align: left;
  text-decoration: underline;
`;
const Timing = styled.div`
  margin: 20px 0;
`;
const DateLabel = styled.p`
  font-size: 14px;
  font-weight: 500;
  margin: 8px 0;
  letter-spacing: 0.1px;
  text-align: left;
  color: #00225a;
`;

const Label = styled.p`
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.1px;
  text-align: left;
  color: #828282;
`;

const InputWrapper = styled.div`
  position: relative;
  width: 100%;
`;

const StyledInput = styled.input`
  width: 95%;
  border: 1px solid #e0e0e0;
  padding: 8px 12px;
  border-radius: 10px;
  outline: none;
  cursor: pointer;
  /* Hide the default date/time text */
  color: transparent; /* Make the text transparent */

  &::-webkit-input-placeholder {
    color: transparent; /* Hide default placeholder */
  }

  &:focus {
    color: #000; /* Show text color on focus */
  }

  &:focus + div {
    display: none; /* Hide placeholder when input is focused or has value */
  }
`;

const Placeholder = styled.div`
  position: absolute;
  left: 12px;
  top: 10px;
  color: #000;
  pointer-events: none;
  transition: 0.2s ease all;
  font-size: 14px;
  font-weight: 500;
  opacity: 0.6;
`;
