"use client";

import { SearchBar } from "@/components";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import { ArrowLeft } from "iconsax-react";

import React, { useEffect, useMemo, useState } from "react";
import styled from "styled-components";
import TradesFilter from "../components/TradesFilter";
import { useRouter } from "next/navigation";
import { useTradesQuery } from "@/redux/api/p2p";
import SelectedFiltersComponent from "@/components/SelectedFilter";
import dayjs from "dayjs";

const AllTrades = () => {
  const [pageSize, setPageSize] = useState(10);
  const [currentPage, setCurrentPage] = useState(1);
  const [searchInput, setSearchInput] = useState<string>("");

  const [searchQuery, setSearchQuery] = useState<string>("");

  const [createdAt, setCreatedAt] = useState<string>("");
  const [selectedFilters, setSelectedFilters] = useState<
    Record<string, string>
  >({});

  const { data, isLoading, isFetching, refetch } = useTradesQuery({
    page: currentPage,
    pageSize,
    username: searchQuery || undefined,
    createdAt: createdAt || undefined,
  });

  useEffect(() => {
    if (searchInput === "") {
      setSearchQuery("");
      setCurrentPage(1);
    }
  }, [searchInput]);

  const columns = [
    {
      title: "Price",
      dataIndex: "price",
      key: "price",

      render: (key: string) => <Price>{key || "N/A"}</Price>,
    },

    {
      title: "Amount ",
      dataIndex: "amount",
      key: "amount",

      render: (amount: number) => <Text>{amount}</Text>,
    },
    {
      title: "Total ",
      dataIndex: "total",
      key: "total",

      render: (count: number) => <Text>{count}</Text>,
    },
    {
      title: "Traded On ",
      dataIndex: "tradedOn",
      key: "tradedOn",

      render: (date: number) =>
        new Date(date).toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        }),
    },

    {
      title: "Maker",
      dataIndex: "maker",
      key: "maker",

      render: (_: any, record: any) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserContent>
              <UserName>{record.merchantUsername || "N/A"}</UserName>
            </UserContent>
          </div>
        </UserInfoSection>
      ),
    },

    {
      title: "Taker",
      dataIndex: "taker",
      key: "taker",

      render: (_: any, record: any) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserContent>
              <UserName>{record.customerUsername || "N/A"}</UserName>
            </UserContent>
          </div>
        </UserInfoSection>
      ),
    },
  ];

  const dataSource = useMemo(() => {
    const tradesArray = data?.data?.data;

    if (!Array.isArray(tradesArray)) return [];

    return tradesArray.map((item) => ({
      key: item.id,
      price: `${item.price} ${item.currency}`,
      amount: `${item.specifiedAssetAmount} ${item.asset}`,
      total: `${item.paymentAmount} ${item.currency}`,
      tradedOn: item.createdAt ?? "N/A",
      merchantUsername: item.merchantUsername ?? "N/A",
      customerUsername: item.customerUsername ?? "N/A",
    }));
  }, [data]);

  const router = useRouter();
  const goBack = () => {
    router.back();
  };

  const handleSearch = () => {
    setSearchQuery(searchInput);
    setCurrentPage(1);
  };

  const applyCreatedAtFilter = (value: string) => {
    // tm-api's createdAt filter compares DATE(created_at) against this
    // value, so it must be a plain YYYY-MM-DD date, not a full ISO
    // datetime (which would never match).
    const dateOnly = dayjs(value).format("YYYY-MM-DD");
    setCreatedAt(dateOnly);

    const formatted = new Date(dateOnly).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });

    setSelectedFilters((prev) => ({
      ...prev,
      TradedOn: ` ${formatted}`,
    }));

    setCurrentPage(1);
    refetch();
  };

  const handleRemoveFilter = (key: string) => {
    // Remove from visual filter tags
    setSelectedFilters((prev) => {
      const updated = { ...prev };
      delete updated[key];
      return updated;
    });

    // Reset filter value
    if (key === "TradedOn") {
      setCreatedAt("");
    }

    setCurrentPage(1);
    refetch();
  };

  return (
    <Container>
      <BackButton onClick={goBack}>
        <ArrowLeft color="#00225A" />
      </BackButton>
      <Heading>Trades</Heading>
      <HeaderContainer>
        <SearchBar
          customWidth="320px"
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          handleSearch={handleSearch}
        />
        <TradesFilter
          createdAt={createdAt}
          setCreatedAt={applyCreatedAtFilter}
        />
      </HeaderContainer>

      <SelectedFiltersComponent
        selectedFilters={selectedFilters}
        handleRemoveFilter={handleRemoveFilter}
      />

      {/* <TableContainer> */}
      {isLoading ? (
        <UserName>Loading trades...</UserName>
      ) : searchInput && isFetching ? (
        <>
          <UserName>Searching user...</UserName>
        </>
      ) : searchInput && data?.data?.total === 0 ? (
        <>
          <UserName>No user match your search.</UserName>
        </>
      ) : (
        <>
          <CustomTable columns={columns} dataSource={dataSource} />
          <Pagination
            currentPage={currentPage}
            totalCount={data?.data?.total ?? 0}
            pageSize={pageSize}
            onPageChange={setCurrentPage}
            onPageSizeChange={setPageSize}
            isFetching={isFetching}
          />{" "}
        </>
      )}
    </Container>
  );
};

export default AllTrades;

const Container = styled.section`
  padding: 32px;
  border-radius: 24px;
  background-color: #ffffff;
  opacity: 0px;

  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 12%;
  }

  th:nth-child(2),
  td:nth-child(2) {
    width: 15%;
  }

  th:nth-child(3),
  td:nth-child(3) {
    width: 12%;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 14%;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 17%;
  }

  th:nth-child(6),
  td:nth-child(6) {
    width: 12%;
  }

  th {
    color: #828282;
    font-size: 14px;
    font-weight: 500;
  }
  td {
    color: #00225a;
    font-weight: 400;
  }
`;

const BackButton = styled.div`
  cursor: pointer;
`;
const Heading = styled.h1`
  color: #00225a;
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  text-align: left;
  margin-top: 10px;
`;
const HeaderContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 0;
`;
const UserInfoSection = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

const UserContent = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

const UserName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const Price = styled.div`
  font-size: 14px;
  font-weight: 500;
  line-height: 16px;
  color: #00a859;
`;

const Text = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 16px;
  color: #00225a;
`;
