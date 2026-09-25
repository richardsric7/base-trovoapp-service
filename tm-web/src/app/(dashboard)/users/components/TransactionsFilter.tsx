
import React, { useState } from "react";
import styled from "styled-components";
import Calender from "@/app/(dashboard)/admin-users/components/Calender";
import calenderIcon from "@/assets/images/calender.svg";
import { DropdownSelect } from "@/components";
import SecondaryButton from "@/components/SecondaryButton";
import PrimaryButton from "@/components/PrimaryButton";
import Image from "next/image";
import CustomFilter from "@/components/CustomFilter";
import { useGetPaymentHistoryQuery } from "@/redux/api/users";
import { useParams } from "next/navigation";

interface TransactionsFilterProps {
  onApplyFilters: (filters: any) => void;
  currentFilters: any;
}

const TransactionsFilter: React.FC<TransactionsFilterProps> = ({
  onApplyFilters,
  currentFilters,
}) => {
  const params = useParams();
  const username = Array.isArray(params.username)
    ? params.username[0]
    : params.username;

  const { data } = useGetPaymentHistoryQuery({
    from: username,
    page: 1,
    pageSize: 100,
  });

  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [selectValues, setSelectValues] = useState(currentFilters);
  const [openCalender, setOpenCalender] = useState(false);
  const [selectedDate, setSelectedDate] = useState(
    currentFilters.date || "dd-mm-yy"
  );

  const fromOptions = data?.data?.data?.map((item: any) => item.from) || [];
  const categoryOptions =
    data?.data?.data?.map((item: any) => item.transactionType) || [];

  const handleFilterApplication = () => {
    onApplyFilters({
      ...selectValues,
      date: selectedDate !== "dd-mm-yy" ? selectedDate : undefined,
    });
    setIsFilterOpen(false);
  };

  const handleFilterReset = () => {
    setSelectValues({});
    setSelectedDate("dd-mm-yy");
    onApplyFilters({});
  };

  const handleDateSelection = (date: Date) => {
    const formattedDate = date.toLocaleDateString("en-GB");
    setSelectedDate(formattedDate);
    setOpenCalender(false);
  };

  return (
    <>
      <CustomFilter
        position={{ right: "40px", top: "500px" }}
        open={isFilterOpen}
        onClose={setIsFilterOpen}
      >
        <FilterContent>
          <DropdownSelect
            options={fromOptions}
            placeholder="All"
            labelText="From"
            labelColor="#828282"
            placeholderColor="#00225A"
            value={selectValues.from || ""}
            onSelect={(item) =>
              setSelectValues({ ...selectValues, from: item })
            }
          />

          <DropdownSelect
            options={[
              "PAYMENT",
              "SWAP XBN>TROV",
              "SWAP TROV>USDT",
              "SWAP TROV>XBN",
            ]}
            placeholder="All"
            labelText="Category"
            value={selectValues.category || ""}
            onSelect={(item) =>
              setSelectValues({ ...selectValues, category: item })
            }
          />

          <div>
            <SelectText>Created on</SelectText>
            <CalenderRow onClick={() => setOpenCalender(!openCalender)}>
              {selectedDate}
              <Image src={calenderIcon} alt="calender-icon" />
            </CalenderRow>
            {openCalender && <Calender onSelectDate={handleDateSelection} />}
          </div>

          <ButtonContainer>
            <SecondaryButton
              onClick={handleFilterReset}
              buttonStyle={{
                padding: "6px 20px",
              }}
            >
              Reset
            </SecondaryButton>
            <PrimaryButton
              onClick={handleFilterApplication}
              buttonStyle={{
                padding: "8px 16px",
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

export default TransactionsFilter;

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
