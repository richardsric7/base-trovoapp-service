"use client";
import React, { useState } from "react";
import styled from "styled-components";
import Link from "next/link";
import { FaPlus } from "react-icons/fa6";
import { SearchBar } from "@/components";
import { DropdownSelect } from "@/components/Dropdown";
import PrimaryButton from "@/components/PrimaryButton";
import ToggleSwitch from "@/components/ToggleSwitch";
import CustomTable from "@/components/CustomTable";
import {
  ICuratedAsset,
  useGetCuratedAssetsQuery,
  useSetCuratedAssetP2PEnabledMutation,
  useSetCuratedAssetInactiveMutation,
} from "@/redux/api/curatedAssets";

const P2P_FILTER_OPTIONS = ["All assets", "P2P enabled", "P2P disabled"];
const STATUS_FILTER_OPTIONS = ["All statuses", "Active", "Inactive"];
const PAGE_SIZE = 20;

const AssetCurationPage = () => {
  const [search, setSearch] = useState("");
  const [searchTrigger, setSearchTrigger] = useState("");
  const [p2pFilter, setP2pFilter] = useState(P2P_FILTER_OPTIONS[0]);
  const [statusFilter, setStatusFilter] = useState(STATUS_FILTER_OPTIONS[0]);
  const [page, setPage] = useState(1);

  const p2pEnabled =
    p2pFilter === "P2P enabled" ? "true" : p2pFilter === "P2P disabled" ? "false" : undefined;
  const inactive =
    statusFilter === "Active" ? "false" : statusFilter === "Inactive" ? "true" : undefined;

  const { data, isLoading } = useGetCuratedAssetsQuery({
    page,
    pageSize: PAGE_SIZE,
    assetCode: searchTrigger || undefined,
    p2pEnabled,
    inactive,
  });
  const [setP2PEnabled] = useSetCuratedAssetP2PEnabledMutation();
  const [setInactive] = useSetCuratedAssetInactiveMutation();

  const handleSearchInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearch(e.target.value);
    if (e.target.value.trim() === "") setSearchTrigger("");
  };
  const handleSearch = () => {
    setPage(1);
    setSearchTrigger(search);
  };

  const columns = [
    {
      title: "Asset",
      dataIndex: "assetCode",
      render: (_: any, record: ICuratedAsset) => (
        <div>
          <AssetCode>{record.assetCode}</AssetCode>
          <AssetName>{record.assetName || "—"}</AssetName>
        </div>
      ),
    },
    { title: "Class ID", dataIndex: "assetClassId" },
    { title: "Decimals", dataIndex: "decimalPlaces" },
    {
      title: "Withdrawable",
      dataIndex: "withdrawable",
      render: (v: number) => (v ? "Yes" : "No"),
    },
    {
      title: "Active",
      dataIndex: "inactive",
      render: (_: any, record: ICuratedAsset) => (
        <ToggleSwitch
          checked={!record.inactive}
          onChange={(checked) => setInactive({ id: record.id, inactive: !checked })}
        />
      ),
    },
    {
      title: "P2P Enabled",
      dataIndex: "p2pEnabled",
      render: (_: any, record: ICuratedAsset) => (
        <ToggleSwitch
          checked={record.p2pEnabled}
          onChange={(checked) => setP2PEnabled({ id: record.id, p2pEnabled: checked })}
        />
      ),
    },
    {
      title: "Actions",
      dataIndex: "actions",
      render: (_: any, record: ICuratedAsset) => (
        <StyledLink href={`/assetcuration/${record.id}`}>Edit</StyledLink>
      ),
    },
  ];

  return (
    <PageContainer>
      <Header>
        <TitleSection>
          <Title>Curated Assets</Title>
          <AssetCount>
            Total: <HighlightedText>{data?.data?.total ?? 0} assets</HighlightedText>
          </AssetCount>
          <Text>
            An asset must be P2P Enabled here before it can be used to create a P2P offer or be
            found in marketplace search.
          </Text>
        </TitleSection>
        <StyledButtonLink href="/assetcuration/addcuration">
          <PrimaryButton
            buttonStyle={{
              width: "auto",
              padding: "10px 20px",
              display: "flex",
              alignItems: "center",
              gap: "4px",
            }}
          >
            Add Asset <FaPlus />
          </PrimaryButton>
        </StyledButtonLink>
      </Header>

      <FiltersSection>
        <SearchBar value={search} onChange={handleSearchInputChange} handleSearch={handleSearch} />
        <FilterGroup>
          <DropdownSelect
            value={statusFilter}
            options={STATUS_FILTER_OPTIONS}
            onSelect={(item) => {
              setPage(1);
              setStatusFilter(item);
            }}
            placeholder="Status"
            labelText=""
            backgroundColor="#F2F6F9"
            borderless
            iconColor="#00225A"
          />
          <DropdownSelect
            value={p2pFilter}
            options={P2P_FILTER_OPTIONS}
            onSelect={(item) => {
              setPage(1);
              setP2pFilter(item);
            }}
            placeholder="Filter"
            labelText=""
            backgroundColor="#F2F6F9"
            borderless
            iconColor="#00225A"
          />
        </FilterGroup>
      </FiltersSection>

      <CustomTable
        columns={columns}
        dataSource={data?.data?.data ?? []}
        totalItems={data?.data?.total ?? 0}
        pageSize={PAGE_SIZE}
        isLoading={isLoading}
        onPageChange={setPage}
      />
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
  align-items: flex-start;
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
  font-size: 14px;
  font-weight: 400;
  color: #828282;
  margin: 0;
  line-height: 22px;
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

const FilterGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

const StyledButtonLink = styled(Link)`
  text-decoration: none;
`;

const StyledLink = styled(Link)`
  color: #007cdf;
  font-weight: 600;
  font-size: 13px;
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
`;

const AssetCode = styled.p`
  font-weight: 600;
  font-size: 14px;
  color: #00225a;
  margin: 0;
`;

const AssetName = styled.p`
  font-size: 12px;
  color: #828282;
  margin: 0;
`;
