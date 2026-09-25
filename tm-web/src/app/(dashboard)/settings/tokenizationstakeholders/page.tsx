"use client";
import React, { useState } from "react";
import { SearchBar } from "@/components";
import styled from "styled-components";
import PrimaryButton from "@/components/PrimaryButton";
import { FaPlus } from "react-icons/fa6";
import Link from "next/link";
import StackholdersTable from "./components/StackholdersTable";
const TokenizationStakeHolders = () => {
  return (
    <PageContainer>
      <Header>
        <Title> Tokenization Stakeholders</Title>

        <StyledLink href="/settings/tokenizationstakeholders/addstakeholder">
          <PrimaryButton
            buttonStyle={{
              width: "auto",
              padding: "10px 20px",
              display: "flex",
              alignItems: "center",
              gap: "4px",
            }}
          >
            <FaPlus /> Stakeholder
          </PrimaryButton>
        </StyledLink>
      </Header>

      <FiltersSection>
        <SearchBar />
        {/* <AssetCurationFilter /> */}
      </FiltersSection>

      <StackholdersTable />
    </PageContainer>
  );
};

export default TokenizationStakeHolders;

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
