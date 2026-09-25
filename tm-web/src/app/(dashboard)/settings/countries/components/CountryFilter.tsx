import { DropdownSelect } from "@/components";
import CustomFilter from "@/components/CustomFilter";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import React, { useState } from "react";
import styled from "styled-components";

const CountryFilter = () => {
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [selectValues, setSelectValues] = useState<any>({});
  const filterCriteria = {
    Country: ["Nigeria", "Cameroon", "Gahna"],
  };

  const handleFilterApplication = () => {
    setIsFilterOpen(false);
  };

  const handleFilterReset = () => {
    setSelectValues({});

    console.log("Filters Reset");
  };
  return (
    <>
      <CustomFilter
        position={{ top: "270px", right: "30px" }}
        open={isFilterOpen}
        onClose={setIsFilterOpen}
      >
        <FilterContent>
          {/* Location Filter */}
          <DropdownSelect
            options={filterCriteria.Country}
            placeholder="All"
            labelText="Country"
            labelColor="#828282"
            placeholderColor="#00225A"
            value={selectValues["Country"] || ""}
            onSelect={(item) =>
              setSelectValues({ ...selectValues, Location: item })
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
  );
};

export default CountryFilter;

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
