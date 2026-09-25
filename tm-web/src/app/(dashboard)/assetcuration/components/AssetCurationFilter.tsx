"use client";
import React, { useState } from "react";
import FilterIcon from "@/assets/images/filter.svg";
import styled from "styled-components";
import Calender from "@/app/(dashboard)/admin-users/components/Calender";
import calenderIcon from "@/assets/images/calender.svg";
import SecondaryButton from "@/components/SecondaryButton";
import PrimaryButton from "@/components/PrimaryButton";
import Image from "next/image";
import { FaX } from "react-icons/fa6";

const AssetCurationFilter = () => {
  const [openFilter, setOpenFilter] = useState<boolean>(false);
  const [selectValues, setSelectValues] = useState<any>({});
  const [openCalender, setOpenCalender] = useState<boolean>(false);
  const [selectedDate, setSelectedDate] = useState<string>("dd-mm-yy");
  const [quantity, setQuantity] = useState(0);

  const handleFilterApplication = () => {
    console.log("Selected Filters: ", { ...selectValues, date: selectedDate });
    setOpenFilter(false);
  };

  const handleFilterReset = () => {
    setSelectValues({});
    setSelectedDate("dd-mm-yy");
    console.log("Filters Reset");
  };

  const handleDateSelection = (date: Date) => {
    const formattedDate = date.toLocaleDateString("en-GB");
    setSelectedDate(formattedDate);
    setOpenCalender(false);
    console.log("Selected Date: ", formattedDate);
  };

  return (
    <>
      <FilterButton
        onClick={() => setOpenFilter(!openFilter)}
        active={openFilter} // passing the active state
      >
        <StyledImage src={FilterIcon} alt="filter-icon" active={openFilter} />
        <Text active={openFilter}>Filter</Text>
      </FilterButton>

      {openFilter && (
        <FilterContainer>
          <FilterRow>
            <FilterText>Filter</FilterText>
            <FaX
              onClick={() => setOpenFilter(false)}
              cursor="pointer"
              size={12}
            />
          </FilterRow>
          <Divider />

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

            {/* Date Selection */}
            <div>
              <SelectText>Updated on</SelectText>
              <CalenderRow onClick={() => setOpenCalender(!openCalender)}>
                {!openCalender && selectedDate}
                <Image src={calenderIcon} alt="calender-icon" />
              </CalenderRow>
              {openCalender && <Calender onSelectDate={handleDateSelection} />}
            </div>
            <Divider />

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
        </FilterContainer>
      )}
    </>
  );
};

export default AssetCurationFilter;

const FilterButton = styled.button<{ active: boolean }>`
  border: 1px solid ${(props) => (props.active ? "#007CDF" : "#E0E0E0")};
  background-color: #ffffff;
  width: 110px;
  height: 40px;
  padding: 8px 16px;
  gap: 15px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
`;

const StyledImage = styled(Image)<{ active: boolean }>`
  filter: ${(props) =>
    props.active
      ? "brightness(0) saturate(100%) invert(35%) sepia(97%) saturate(1064%) hue-rotate(183deg) brightness(99%) contrast(101%)"
      : "none"};
`;

const Text = styled.h1<{ active: boolean }>`
  color: ${(props) => (props.active ? "#007CDF" : "#00225A")};
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  letter-spacing: 0.1px;
  margin: 0;
`;

const FilterContainer = styled.div`
  position: absolute;
  top: 280px;
  right: 40px;
  background-color: #ffffff;
  box-shadow: 0px 4px 100px 0px #00000026;
  width: 300px;
  padding: 24px;
  border-radius: 24px;
  z-index: 3000;
  margin-bottom: 20px;
  z-index: 10;
`;

const FilterRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const FilterText = styled.h3`
  font-size: 20px;
  font-weight: 700;
  color: #00225a;
`;
const Divider = styled.div`
  background-color: #e5e5ef;
  height: 1px;
  width: calc(100% + 48px);
  margin: 20px -24px 10px -24px;
`;
const FilterContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 15px;
  marign-top: 15px;
`;

const SelectText = styled.p`
  font-size: 14px;
  color: #828282;
`;

const LabelText = styled.p`
  font-size: 14px;
  color: #828282;
`;
const InputField = styled.input`
  padding: 8px 4px;
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  opacity: 0px;
  width: 100%;
  outline: none;
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
