"use client";
import React, { useState } from "react";
import styled from "styled-components";
import { SearchBar } from "@/components";
import AssetFilter from "./components/AssetFilter";
import AssetsTable from "./components/AssetsTable";
import Image from "next/image";
import ttassetIcon from "@/assets/images/ttasset.svg";
import { CiCircleCheck } from "react-icons/ci";
import { FiClock } from "react-icons/fi";
import { IoIosCloseCircleOutline } from "react-icons/io";
import TokenizedAssetStat from "./components/TokeniziedAssetStat";
import { useGetTokenizationStatisticsQuery } from "@/redux/api/assettokenization";

const TokenizedAssetsPage = () => {
  const { data } = useGetTokenizationStatisticsQuery();
  const [filters, setFilters] = useState<any>({
    assetType: "",
    assetTokenizationStatus: undefined,
    createdBetween: "",
    initiatorUsername: "",
    assetSector: "",
  });
  const [searchTerm, setSearchTerm] = useState("");

  const formatAmount = (value?: number) =>
    typeof value === "number"
      ? value.toLocaleString(undefined, {
          minimumFractionDigits: 2,
          maximumFractionDigits: 2,
        })
      : "0.00";

  const handleApplyFilters = (newFilters: any) => {
    setFilters(newFilters);
  };

  const handleSearch = (value: string) => {
    setSearchTerm(value);
  };

  return (
    <>
      <CardContainer>
        <TokenizedAssetStat
          title="Total Tokenized Assets"
          value={data?.data?.total?.count ?? 0}
          icon={
            <Image src={ttassetIcon} alt="asset-icon" width={24} height={24} />
          }
          color="#00225A"
          bgColor="#F2F6F9"
          amount={`${formatAmount(data?.data?.total?.totalTokenizedValue)} NGN`}
        />
        <TokenizedAssetStat
          title="Approved Tokenized Assets"
          value={data?.data?.approved?.count ?? 0}
          icon={<CiCircleCheck color="#00A859" size={24} />}
          color="#00A859"
          bgColor="#00A8591A"
          amount={`${formatAmount(data?.data?.approved?.totalTokenizedValue)} NGN`}
        />
        <TokenizedAssetStat
          title="Pending Tokenized Assets"
          value={data?.data?.pending?.count ?? 0}
          icon={<FiClock color="#007CDF" size={24} />}
          color="#007CDF"
          bgColor="#F2F6F9"
          amount={`${formatAmount(data?.data?.pending?.totalTokenizedValue)} NGN`}
        />
        <TokenizedAssetStat
          title="Rejected Tokenized Assets"
          value={data?.data?.rejected?.count ?? 0}
          icon={<IoIosCloseCircleOutline color="#BE3800" size={24} />}
          color="#BE3800"
          bgColor="#BE38001A"
          amount={`${formatAmount(data?.data?.rejected?.totalTokenizedValue)} NGN`}
        />
      </CardContainer>
      <Container>
        <Header>
          <TitleSection>
            <Title>Tokenized Assets</Title>
          </TitleSection>
        </Header>

        <FiltersSection>
          <SearchBar
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
          <AssetFilter onApply={handleApplyFilters} />
        </FiltersSection>

        <AssetsTable
          filters={{
            ...filters,
            ...(searchTerm ? { assetName: searchTerm } : {}),
          }}
        />
      </Container>
    </>
  );
};

export default TokenizedAssetsPage;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
`;

const CardContainer = styled.div`
  display: flex;
  gap: 10px;
  padding: 20px 0;

  width: 100%;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const TitleSection = styled.div``;

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
