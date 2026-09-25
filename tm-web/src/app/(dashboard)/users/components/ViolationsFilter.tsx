import { DropdownSelect } from "@/components";
import CustomFilter from "@/components/CustomFilter";
import React, { useState } from "react";
import styled from "styled-components";
import Calender from "@/app/(dashboard)/admin-users/components/Calender";
import calenderIcon from "@/assets/images/calender.svg";
import Image from "next/image";
import SecondaryButton from "@/components/SecondaryButton";
import PrimaryButton from "@/components/PrimaryButton";
const ViolationsFilter = () => {
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [selectValues, setSelectValues] = useState<any>({});
  const [openCalender, setOpenCalender] = useState<boolean>(false);
  const [selectedDate, setSelectedDate] = useState<string>("dd-mm-yy");

  const handleFilterApplication = () => {
    // setSelectedFilters({ ...selectValues, date: selectedDate });
    setIsFilterOpen(false);
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
    ViolationType: ["suspended", "suspended", "suspended"],
    Status: ["Active", "Inactive", "Suspended"],
  };

  return (
    <div>
      <>
        <CustomFilter
          position={{ top: "740px", right: "440px" }}
          open={isFilterOpen}
          onClose={setIsFilterOpen}
        >
          <FilterContent>
            {/* Location Filter */}
            <DropdownSelect
              options={filterCriteria.ViolationType}
              placeholder="All"
              labelText="Violation Type"
              labelColor="#828282"
              placeholderColor="#00225A"
              value={selectValues["Location"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, Location: item })
              }
            />

            {/* Date Selection */}
            <div>
              <SelectText>Registration Date</SelectText>

              <CalenderRow onClick={() => setOpenCalender(!openCalender)}>
                {!openCalender && selectedDate}
                <Image src={calenderIcon} alt="calender-icon" />
              </CalenderRow>
              {openCalender && <Calender onSelectDate={handleDateSelection} />}
            </div>
            {/* Activation Filter */}
            <DropdownSelect
              options={filterCriteria.Status}
              labelText="Status"
              placeholder="All"
              value={selectValues["status"] || ""}
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
        </CustomFilter>
      </>
    </div>
  );
};

export default ViolationsFilter;

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
