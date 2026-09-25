"use client";
import CustomFilter from "@/components/CustomFilter";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import React, { useState } from "react";
import styled from "styled-components";
import Calender from "../../admin-users/components/Calender";
import calenderIcon from "@/assets/images/calender.svg";
import Image from "next/image";

interface TokenFilterProps {
  //   onApplyFilters: (filters: any) => void;
  //   currentFilters: any;
}

const TokensFilter: React.FC<TokenFilterProps> = (
  {
    //   onApplyFilters,
    //   currentFilters,
  }
) => {
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [openCalender, setOpenCalender] = useState(false);
  const [selectedDate, setSelectedDate] = useState();
  // currentFilters.date || "dd-mm-yy"
  const [selectValues, setSelectValues] = useState();

  const handleDateSelection = (date: Date) => {
    const formattedDate = date.toLocaleDateString("en-GB");
    // setSelectedDate(formattedDate);
    setOpenCalender(false);
  };

  const handleFilterApplication = () => {
    // onApplyFilters({
    //   //   ...selectValues,
    //   date: selectedDate !== "dd-mm-yy" ? selectedDate : undefined,
    // });
    setIsFilterOpen(false);
  };

  const handleFilterReset = () => {
    // setSelectValues({});
    // setSelectedDate("dd-mm-yy");
    // onApplyFilters({});
  };

  return (
    <>
      <CustomFilter
        position={{ right: "40px", top: "500px" }}
        open={isFilterOpen}
        onClose={() => setIsFilterOpen(false)}
      >
        <FilterContent>
          {/* Date Selection */}
          <div>
            <SelectText>Created on</SelectText>
            <CalenderRow onClick={() => setOpenCalender(!openCalender)}>
              {!openCalender && selectedDate}
              <Image src={calenderIcon} alt="calender-icon" />
            </CalenderRow>
            {openCalender && <Calender onSelectDate={handleDateSelection} />}
          </div>

          {/* Apply and Reset Buttons */}
          <ButtonContainer>
            <SecondaryButton
              onClick={handleFilterReset}
              buttonStyle={{ padding: "8px 16px" }}
            >
              Reset
            </SecondaryButton>
            <PrimaryButton
              onClick={handleFilterApplication}
              buttonStyle={{ padding: "8px 16px" }}
            >
              Apply
            </PrimaryButton>
          </ButtonContainer>
        </FilterContent>
      </CustomFilter>
    </>
  );
};

export default TokensFilter;

const FilterContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-top: 15px;
`;

const SelectText = styled.p`
  font-size: 14px;
  color: #828282;
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
  align-items: center;
  justify-content: space-between;
  gap: 20px;
`;
