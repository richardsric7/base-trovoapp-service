"use client";
import React, { useState } from "react";
import FilterIcon from "@/assets/images/filter.svg";
import styled from "styled-components";
import Calender from "@/app/(dashboard)/admin-users/components/Calender";
import calenderIcon from "@/assets/images/calender.svg";
import { DropdownSelect } from "@/components";
import SecondaryButton from "@/components/SecondaryButton";
import PrimaryButton from "@/components/PrimaryButton";
import Image from "next/image";
import { FaX } from "react-icons/fa6";

const AdminFilter = () => {
  const [openFilter, setOpenFilter] = useState<boolean>(false);
  const [selectValues, setSelectValues] = useState<any>({});
  const [openCalender, setOpenCalender] = useState<boolean>(false);
  const [selectedDate, setSelectedDate] = useState<string>("dd-mm-yy");

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

  const filterCriteria = {
    Role: ["Support", "Super Admin", "Compliance"],
    AccountStatus: ["Active", "Inactive", "Suspended"],
    Status: ["Pending", "Approved", "Rejected"],
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
            {/* Role Filter */}
            <DropdownSelect
              options={filterCriteria.Role}
              placeholder="All"
              labelText="Role"
              labelColor="#828282"
              placeholderColor="#00225A"
              value={selectValues["Category"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, Category: item })
              }
            />

            {/* Date Selection */}
            <div>
              <SelectText>Created on</SelectText>
              <CalenderRow onClick={() => setOpenCalender(!openCalender)}>
                {!openCalender && selectedDate}
                <Image src={calenderIcon} alt="calender-icon" />
              </CalenderRow>
              {openCalender && <Calender onSelectDate={handleDateSelection} />}
            </div>
            {/* Account Filter */}
            <DropdownSelect
              options={filterCriteria.AccountStatus}
              labelText="Account Status"
              placeholder="All"
              value={selectValues["Location"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, Location: item })
              }
            />

            {/* Status Filter */}
            <DropdownSelect
              options={filterCriteria.Status}
              placeholder="All"
              labelText="Status"
              value={selectValues["Status"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, Status: item })
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
        </FilterContainer>
      )}
    </>
  );
};

export default AdminFilter;

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
  top: 350px;
  right: 40px;
  background-color: #ffffff;
  box-shadow: 0px 4px 100px 0px #00000026;
  width: 280px;
  padding: 30px 24px;
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
  width: 100%;
  margin: 6px 0;
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
