"use client";
import React, { useState } from "react";
import FilterIcon from "@/assets/images/filter.svg";
import styled from "styled-components";
import { DropdownSelect } from "@/components";
import SecondaryButton from "@/components/SecondaryButton";
import PrimaryButton from "@/components/PrimaryButton";
import Image from "next/image";
import { FaX } from "react-icons/fa6";
// import { DatePicker } from "antd";
import { DatePicker } from "antd";
import { useGetTokenizationParamsQuery } from "@/redux/api/assettokenization";

interface AssetFilterProps {
  onApply: (filters: any) => void;
}

const statusOptions = [
  { label: "All", value: "" },
  { label: "Pending Vetting", value: "0" },
  { label: "Awaiting Payment", value: "1" },
  { label: "Payment Made", value: "2" },
  { label: "Processing", value: "3" },
  { label: "Tokenization Approved", value: "4" },
  { label: "Primary Sale Started", value: "5" },
  { label: "Secondary Sales", value: "6" },
  { label: "Liquidated", value: "7" },
  { label: "Refunded", value: "8" },
];

const defaultSectors = [
  "Real Estate",
  "Commodity",
  "Energy",
  "Infrastructure",
  "Agriculture",
  "Technology",
];

const AssetFilter = ({ onApply }: AssetFilterProps) => {
  const { data: paramsData } = useGetTokenizationParamsQuery();
  const [openFilter, setOpenFilter] = useState<boolean>(false);
  const [selectValues, setSelectValues] = useState<any>({
    assetType: "",
    assetTokenizationStatus: "",
    createdBetween: "",
    initiatorUsername: "",
    assetSector: "",
  });

  const assetTypes =
    paramsData?.assetTypes?.map((type) => type.assetType) || [];

  // Use API sectors, then derived sectors, then default sectors as last resort
  const derivedSectors = Array.from(
    new Set(
      paramsData?.assetTypes?.map((type) => type.assetType.split(" ")[0]) || [],
    ),
  );

  const assetSectors = paramsData?.assetSectors?.length
    ? paramsData.assetSectors.map((sector) => sector.sector)
    : derivedSectors.length
      ? derivedSectors
      : defaultSectors;

  const handleFilterApplication = () => {
    const formattedFilters: any = {
      ...selectValues,
      assetTokenizationStatus:
        selectValues.assetTokenizationStatus || undefined,
    };

    onApply(formattedFilters);
    setOpenFilter(false);
    console.log("Applying filters:", formattedFilters);
  };

  const handleFilterReset = () => {
    const resetValues = {
      assetType: "",
      assetTokenizationStatus: "",
      createdBetween: "",
      initiatorUsername: "",
      assetSector: "",
    };
    setSelectValues(resetValues);
    onApply(resetValues);
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
            {/* Asset Type Filter */}
            <DropdownSelect
              options={["All", ...assetTypes]}
              placeholder="All"
              labelText="Asset Type"
              value={selectValues["assetType"] || "All"}
              onSelect={(item) =>
                setSelectValues({
                  ...selectValues,
                  assetType: item === "All" ? "" : item,
                })
              }
            />

            {/* Asset Sector Filter */}
            <DropdownSelect
              options={["All", ...assetSectors]}
              placeholder="All"
              labelText="Asset Sector"
              value={selectValues["assetSector"] || "All"}
              onSelect={(item) =>
                setSelectValues({
                  ...selectValues,
                  assetSector: item === "All" ? "" : item,
                })
              }
            />

            {/* Status Filter */}
            <DropdownSelect
              options={statusOptions.map((opt) => opt.label)}
              placeholder="All"
              labelText="Status"
              value={
                statusOptions.find(
                  (opt) => opt.value === selectValues.assetTokenizationStatus,
                )?.label || "All"
              }
              onSelect={(item) =>
                setSelectValues({
                  ...selectValues,
                  assetTokenizationStatus:
                    statusOptions.find((opt) => opt.label === item)?.value ||
                    "",
                })
              }
            />

            {/* Initiator Name */}
            <div>
              <SelectText>Initiator Name</SelectText>
              <TextInput
                placeholder="Enter initiator username"
                value={selectValues.initiatorUsername}
                onChange={(e) =>
                  setSelectValues({
                    ...selectValues,
                    initiatorUsername: e.target.value,
                  })
                }
              />
            </div>

            {/* Date Selection */}
            <div>
              <SelectText>Created Between</SelectText>
              <DatePicker.RangePicker
                format="YYYY-MM-DD"
                style={{ width: "100%" }}
                size="large"
                onChange={(dates, dateStrings) => {
                  setSelectValues({
                    ...selectValues,
                    createdBetween: dateStrings.join("|"),
                  });
                }}
              />
            </div>

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

export default AssetFilter;

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
  top: 400px;
  right: 40px;
  background-color: #ffffff;
  box-shadow: 0px 4px 100px 0px #00000026;
  width: 280px;
  padding: 30px 24px;
  border-radius: 24px;
  z-index: 3000;
  margin-bottom: 20px;
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
  margin-top: 15px;
`;

const SelectText = styled.p`
  font-size: 14px;
  color: #828282;
  padding-bottom: 6px;
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

const TextInput = styled.input`
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  width: 100%;
  outline: none;
  font-family: inherit;
  color: #00225a;

  &::placeholder {
    color: #828282;
  }
`;
