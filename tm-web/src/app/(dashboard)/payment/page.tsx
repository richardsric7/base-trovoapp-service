"use client";
import { SearchBar } from "@/components";
import React, { useState } from "react";
import styled from "styled-components";
import PaymentFilter from "./components/PaymentFilter";
import PaymentTable from "./components/PaymentTable";
import { GetFiatPaymentsParams } from "@/redux/api/payment";

const PaymentPage = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [filters, setFilters] = useState<Partial<GetFiatPaymentsParams>>({});
  const [searchInput, setSearchInput] = useState("");

  const handleSearch = () => {
    setSearchTerm(searchInput);
    // setCurrentPage(1);
    // refetch();
  };

  const handleSearchInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearchInput(value);

    if (value.trim() === "") {
      setSearchTerm("");
      // setCurrentPage(1);
    }
  };
  return (
    <Container>
      <Title>Payments</Title>
      <SubHeader>
        <SearchBar
          value={searchInput}
          onChange={handleSearchInputChange}
          placeholder="Search by username..."
          handleSearch={handleSearch}
        />
        <PaymentFilter filters={filters} onApply={setFilters} />
      </SubHeader>{" "}
      <PaymentTable search={searchTerm} filters={filters} />
    </Container>
  );
};

export default PaymentPage;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;

  border-radius: 24px;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  color: #00225a;
`;

const SubHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 18px;
`;
