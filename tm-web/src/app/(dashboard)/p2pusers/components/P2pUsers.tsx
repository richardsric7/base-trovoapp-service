"use client";

import React, { useState } from "react";
import styled from "styled-components";
import P2pUsersTable from "./P2pUsersTable";
import { SearchBar } from "@/components";
import P2pFilter from "./P2pFilter";

const P2pUsers = () => {
  const [searchTerm, setSearchTerm] = useState<string>("");
  const [searchInput, setSearchInput] = useState<string>("");

  const handleSearchInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearchInput(value);
    if (value.trim() === "") {
      setSearchTerm("");
      //   setCurrentPage(1);
    }
  };

  const handleSearch = () => {
    setSearchTerm(searchInput);
    // setCurrentPage(1);
    // refetch();
  };

  return (
    <PageContainer>
      <Title>Users</Title>

      <FiltersSection>
        <SearchBar
          value={searchInput}
          onChange={handleSearchInputChange}
          handleSearch={handleSearch}
        />
        <P2pFilter />
      </FiltersSection>

      <P2pUsersTable searchTerm={searchTerm} />
    </PageContainer>
  );
};

export default P2pUsers;

const PageContainer = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;

const FiltersSection = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 25px;
  margin: 10px 0 30px 0;
`;
