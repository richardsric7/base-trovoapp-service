"use client";
import { SearchBar } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import React from "react";
import { FaPlus } from "react-icons/fa6";
import styled from "styled-components";
import CountryFilter from "../../countries/components/CountryFilter";
import AssetSectorTable from "../components/AssetSectorTable";
import { useRouter } from "next/navigation";

const AssetSectorPage = () => {
  const router = useRouter();
  return (
    <PageContainer>
      <Header>
        <Title>Asset Sectors</Title>

        <PrimaryButton
          buttonStyle={{
            width: "auto",
            padding: "10px 20px",
            display: "flex",
            alignItems: "center",
            gap: "4px",
            marginRight: "40px",
          }}
          onClick={() => router.push("/settings/assetsector/addsector")}
        >
          <FaPlus /> Sector
        </PrimaryButton>
      </Header>

      <FiltersSection>
        <SearchBar />
        <CountryFilter />
      </FiltersSection>
      <AssetSectorTable />
    </PageContainer>
  );
};

export default AssetSectorPage;

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
