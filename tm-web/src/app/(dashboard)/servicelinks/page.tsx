"use client";
import React, { useState } from "react";
import styled from "styled-components";
import Link from "next/link";
import { FaPlus } from "react-icons/fa6";
import { SearchBar, showErrorToast, showSuccessToast } from "@/components";
import { DropdownSelect } from "@/components/Dropdown";
import PrimaryButton from "@/components/PrimaryButton";
import ToggleSwitch from "@/components/ToggleSwitch";
import CustomTable from "@/components/CustomTable";
import {
  IServiceLink,
  useGetServiceLinksQuery,
  useSetServiceLinkInactiveMutation,
} from "@/redux/api/serviceLinks";

const STATUS_FILTER_OPTIONS = ["All statuses", "Active", "Inactive"];
const PAGE_SIZE = 20;

const maskApiKey = (key: string) => (key.length <= 8 ? key : `${key.slice(0, 4)}...${key.slice(-4)}`);

const ServiceLinksPage = () => {
  const [search, setSearch] = useState("");
  const [searchTrigger, setSearchTrigger] = useState("");
  const [statusFilter, setStatusFilter] = useState(STATUS_FILTER_OPTIONS[0]);
  const [page, setPage] = useState(1);

  const inactive = statusFilter === "Active" ? "false" : statusFilter === "Inactive" ? "true" : undefined;

  const { data, isLoading } = useGetServiceLinksQuery({
    page,
    pageSize: PAGE_SIZE,
    ownerUsername: searchTrigger || undefined,
    inactive,
  });
  const [setInactive] = useSetServiceLinkInactiveMutation();

  const handleSearchInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearch(e.target.value);
    if (e.target.value.trim() === "") setSearchTrigger("");
  };
  const handleSearch = () => {
    setPage(1);
    setSearchTrigger(search);
  };

  const handleCopyApiKey = (apiKey: string) => {
    navigator.clipboard
      .writeText(apiKey)
      .then(() => showSuccessToast("API key copied to clipboard"))
      .catch(() => showErrorToast("Failed to copy API key"));
  };

  const columns = [
    {
      title: "Service",
      dataIndex: "shortName",
      render: (_: any, record: IServiceLink) => (
        <div>
          <ServiceName>{record.shortName}</ServiceName>
          <ServiceLongName>{record.longName || "—"}</ServiceLongName>
        </div>
      ),
    },
    {
      title: "Owner",
      dataIndex: "ownerUsername",
      render: (_: any, record: IServiceLink) => (
        <div>
          <OwnerUsername>{record.ownerUsername}</OwnerUsername>
          <OwnerEmail>{record.owner?.email || "user not found"}</OwnerEmail>
        </div>
      ),
    },
    {
      title: "API Key",
      dataIndex: "apiKey",
      render: (_: any, record: IServiceLink) => (
        <ApiKeyPill onClick={() => handleCopyApiKey(record.apiKey)} title="Click to copy">
          {maskApiKey(record.apiKey)}
        </ApiKeyPill>
      ),
    },
    {
      title: "Verified",
      dataIndex: "verified",
      render: (v: number) => (v ? "Yes" : "No"),
    },
    {
      title: "Active",
      dataIndex: "inactive",
      render: (_: any, record: IServiceLink) => (
        <ToggleSwitch
          checked={!record.inactive}
          onChange={(checked) => setInactive({ id: record.id, inactive: !checked })}
        />
      ),
    },
    {
      title: "Actions",
      dataIndex: "actions",
      render: (_: any, record: IServiceLink) => (
        <StyledLink href={`/servicelinks/${record.id}`}>Edit</StyledLink>
      ),
    },
  ];

  return (
    <PageContainer>
      <Header>
        <TitleSection>
          <Title>Service Links</Title>
          <LinkCount>
            Total: <HighlightedText>{data?.data?.total ?? 0} integrations</HighlightedText>
          </LinkCount>
          <Text>
            White-label partner integrations - each has its own API key and permission
            set gating what it can do on behalf of its owner account.
          </Text>
        </TitleSection>
        <StyledButtonLink href="/servicelinks/addservicelink">
          <PrimaryButton
            buttonStyle={{
              width: "auto",
              padding: "10px 20px",
              display: "flex",
              alignItems: "center",
              gap: "4px",
            }}
          >
            Add Service Link <FaPlus />
          </PrimaryButton>
        </StyledButtonLink>
      </Header>

      <FiltersSection>
        <SearchBar
          value={search}
          onChange={handleSearchInputChange}
          handleSearch={handleSearch}
          placeholder="Search by owner username"
        />
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

export default ServiceLinksPage;

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

const LinkCount = styled.p`
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

const ServiceName = styled.p`
  font-weight: 600;
  font-size: 14px;
  color: #00225a;
  margin: 0;
`;

const ServiceLongName = styled.p`
  font-size: 12px;
  color: #828282;
  margin: 0;
`;

const OwnerUsername = styled.p`
  font-weight: 500;
  font-size: 14px;
  color: #00225a;
  margin: 0;
`;

const OwnerEmail = styled.p`
  font-size: 12px;
  color: #828282;
  margin: 0;
`;

const ApiKeyPill = styled.button`
  font-family: monospace;
  font-size: 12px;
  color: #00225a;
  background-color: #f2f6f9;
  border: none;
  border-radius: 6px;
  padding: 4px 8px;
  cursor: pointer;

  &:hover {
    background-color: #e5eef5;
  }
`;
