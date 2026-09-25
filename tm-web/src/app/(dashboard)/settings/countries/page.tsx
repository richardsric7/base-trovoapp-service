"use client";
import PrimaryButton from "@/components/PrimaryButton";
import Link from "next/link";
import React from "react";
import styled from "styled-components";
import { FaPlus } from "react-icons/fa6";
import { SearchBar } from "@/components";
import CountriesTable from "./components/CountriesTable";
import CountryFilter from "./components/CountryFilter";
const CountiresPage = () => {
  return (
    <PageContainer>
      <Header>
        <Title> Countries</Title>

        <StyledLink href="/countries/addcountry">
          <PrimaryButton
            buttonStyle={{
              width: "auto",
              padding: "10px 20px",
              display: "flex",
              alignItems: "center",
              gap: "4px",
            }}
          >
            <FaPlus /> Country
          </PrimaryButton>
        </StyledLink>
      </Header>

      <FiltersSection>
        <SearchBar />
        <CountryFilter />
      </FiltersSection>

      <CountriesTable />
    </PageContainer>
  );
};

export default CountiresPage;

const PageContainer = styled.section``;

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

const StyledLink = styled(Link)`
  text-decoration: none;
  width: 20%;
`;
