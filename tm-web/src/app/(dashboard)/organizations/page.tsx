"use client";
import React, { useState } from "react";
import styled from "styled-components";
import OrganizationsTable from "./components/OrganizationsTable";
import { SearchBar } from "@/components";

import InviteOrgModal from "./components/InviteOrgModal";
import OrgFilter from "./components/OrgFilter";
import SelectedFiltersComponent from "@/components/SelectedFilter";

const valueToLabel = {
  ASSET_MANAGER: "Asset Manager",
  ASSET_CUSTODIAN: "Asset Custodian",
  RATING_AGENCY: "Rating Agency",
  REGULATOR: "Regulator",
  LEGAL_AGENCY: "Legal Agency",
  PROFESSIONAL_AGENCY: "Professional Agency",
};

const Oragnizationspage = () => {
  const [searchTerm, setSearchTerm] = useState<string>("");

  const [inputValue, setInputValue] = useState(""); // local input value

  const [selectedFilters, setSelectedFilters] = useState<
    Record<string, string>
  >({});

  // Handler for the search button click
  const handleSearch = () => {
    setSearchTerm(inputValue);
  };

  const handleSearchInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setInputValue(value);

    if (value.trim() === "") {
      setSearchTerm("");
    }
  };

  const handleRemoveFilter = (key: string) => {
    setSelectedFilters((prev) => {
      const newFilters = { ...prev };
      delete newFilters[key];
      return newFilters;
    });
  };
  return (
    <PageContainer>
      <Header>
        <Title>Organizations</Title>

        <InviteOrgModal />
      </Header>

      <FiltersSection>
        <SearchBar
          value={inputValue}
          onChange={handleSearchInputChange}
          handleSearch={handleSearch}
        />
        <OrgFilter
          setSelectedFilters={setSelectedFilters}
          selectedFilters={selectedFilters}
        />
      </FiltersSection>

      <SelectedFiltersComponent
        selectedFilters={selectedFilters}
        handleRemoveFilter={handleRemoveFilter}
        valueToLabel={valueToLabel}
      />
      <OrganizationsTable searchTerm={searchTerm} filters={selectedFilters} />
    </PageContainer>
  );
};
export default Oragnizationspage;

const PageContainer = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
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
