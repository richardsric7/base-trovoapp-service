"use client";
import React, { useState } from "react";
import { SearchBar } from "@/components";
import styled from "styled-components";
import AssetCurationFilter from "./components/AssetCurationFilter";
import AssetCurationTable from "./components/AssetCurationTable";
import PrimaryButton from "@/components/PrimaryButton";
import { FaPlus } from "react-icons/fa6";
import Link from "next/link";

const AssetCurationPage = () => {
  return (
    <PageContainer>
      <Header>
        <TitleSection>
          <Title> Curated Assets</Title>

          <AssetCount>
            Total: <HighlightedText>50 Curations</HighlightedText>
          </AssetCount>

          <Text>
            Drag the assets to rearrange them in your desired order on the Trovo
            app, or use the move icon in the Actions column
          </Text>
        </TitleSection>
        <StyledLink href="assetcuration/addcuration">
          <PrimaryButton
            buttonStyle={{
              width: "auto",
              padding: "10px 20px",
              display: "flex",
              alignItems: "center",
              gap: "4px",
            }}
          >
            Curate Assets <FaPlus />
          </PrimaryButton>
        </StyledLink>
      </Header>

      <FiltersSection>
        <SearchBar />
        <AssetCurationFilter />
      </FiltersSection>

      <AssetCurationTable />
    </PageContainer>
  );
};

export default AssetCurationPage;

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

const TitleSection = styled.div`
  flex-grow: 1;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;

const AssetCount = styled.p`
  font-size: 16px;
  font-weight: 400;
  color: #828282;
  margin: 0;
  line-height: 28px;
`;

const Text = styled.p`
  font-size: 16px;
  font-weight: 400;
  color: #828282;
  margin: 0;
  line-height: 28px;
  width: 90%;
`;

const HighlightedText = styled.span`
  font-weight: 500;
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
