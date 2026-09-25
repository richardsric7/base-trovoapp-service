"use client";
import PrimaryButton from "@/components/PrimaryButton";
import React from "react";
import { FaPlus } from "react-icons/fa6";
import styled from "styled-components";
import TokensTable from "./components/TokensTable";
import { SearchBar } from "@/components";
import TokensFilter from "./components/TokensFilter";
import Link from "next/link";

const ManageOtherTokenPage = () => {
  return (
    <PageContainer>
      <Header>
        <Title>Other Tokens</Title>
        <StyledLink href="othertoken/addtoken">
          <PrimaryButton
            buttonStyle={{
              width: "auto",
              padding: "10px 20px",
              display: "flex",
              alignItems: "center",
              gap: "4px",
            }}
          >
            Add New Token
            <FaPlus />
          </PrimaryButton>
        </StyledLink>
      </Header>

      <FiltersSection>
        <SearchBar />
        <TokensFilter />
      </FiltersSection>

      <TokensTable />
    </PageContainer>
  );
};

export default ManageOtherTokenPage;

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

const StyledLink = styled(Link)`
  text-decoration: none;
  width: 20%;
`;
