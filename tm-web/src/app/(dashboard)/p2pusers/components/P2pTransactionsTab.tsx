"use client";
import { SearchBar } from "@/components";
import CustomTable from "@/components/CustomTable";
import Link from "next/link";

import React, { useMemo, useState } from "react";
import { MdVerified } from "react-icons/md";
import styled from "styled-components";

import Pagination from "@/components/CustomPagination";
import TradeDetails from "./TradeDetails";
import P2pFilter from "./P2pFilter";

type TableData = {
  username: string;
  fullName: string;
  price: string;
  amount: string;
  total: number;
  tradedOn: string;
};
const P2pTransactionsTab = () => {
  const [pageSize, setPageSize] = useState(10);
  const [currentPage, setCurrentPage] = useState(1);
  const [openModal, setOpenModal] = useState<boolean>(false);
  const columns = [
    {
      title: "Maker",
      dataIndex: "maker",
      key: "maker",

      render: (_: any, record: TableData) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserContent>
              <UserName>{record.username || "N/A"}</UserName>
              <MdVerified color="#007CDF" />
            </UserContent>
            <FullName>{record.fullName || "N/A"}</FullName>
          </div>
        </UserInfoSection>
      ),
    },

    {
      title: "Taker",
      dataIndex: "taker",
      key: "taker",

      render: (_: any, record: TableData) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserContent>
              <UserName>{record.username || "N/A"}</UserName>
              <MdVerified color="#007CDF" />
            </UserContent>
            <FullName>{record.fullName || "N/A"}</FullName>
          </div>
        </UserInfoSection>
      ),
    },

    {
      title: "Amount ",
      dataIndex: "amount",
      key: "amount",

      render: (amount: number) => <Price>{amount}</Price>,
    },

    {
      title: "Value",
      dataIndex: "value",
      key: "value",

      render: (key: string) => <Text>{key || "N/A"}</Text>,
    },
    {
      title: "Fee ",
      dataIndex: "fee",
      key: "fee",

      render: (count: number) => <Text>{count}</Text>,
    },
    {
      title: "Date ",
      dataIndex: "date",
      key: "date",

      render: (count: number) => <Text>{count}</Text>,
    },

    {
      title: "Status",
      dataIndex: "status",
      key: "status",

      render: (count: number) => <Text>{count}</Text>,
    },
  ];

  // Dummy data for testing
  const dummyData = [
    {
      username: "JohnDoe",
      fullName: "John Doe",
      amount: "+2455 TROV",
      value: "250 NGN",
      fee: "15 NGN",
      date: "23 Sep 2023",
      status: "Completed",
    },
    {
      username: "JaneSmith",
      fullName: "Jane Smith",
      amount: "+2455 TROV",
      value: "250 NGN",
      fee: "15 NGN",
      date: "23 Sep 2023",
      status: "Completed",
    },
    {
      username: "Marlone",
      fullName: "Mark Ovey",
      amount: "+2455 TROV",
      value: "250 NGN",
      fee: "15 NGN",
      date: "23 Sep 2023",
      status: "Completed",
    },

    {
      username: "JaneS",
      fullName: "Jane Smith",
      amount: "+2455 TROV",
      value: "250 NGN",
      fee: "15 NGN",
      date: "23 Sep 2023",
      status: "Completed",
    },
  ];

  const dataSource = useMemo(() => {
    return dummyData.map((item, index) => ({
      key: index + 1,
      username: item.username,
      fullName: item.fullName,
      amount: item.amount,
      value: item.value,
      fee: item.fee,
      date: item.date,
      status: item.status,
    }));
  }, []);

  return (
    <Container>
      <Heading>Transactions</Heading>
      <Header>
        <SearchBar />
        <P2pFilter />
      </Header>

      <CustomTable
        columns={columns}
        dataSource={dataSource}
        onRowClick={() => setOpenModal(!openModal)}
      />
      <Pagination
        currentPage={currentPage}
        totalCount={dataSource.length}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />

      {openModal && (
        <TradeDetails openModal={openModal} setOpenModal={setOpenModal} />
      )}
    </Container>
  );
};

export default P2pTransactionsTab;

const Container = styled.div`
  background-color: #fff;
  padding: 40px 32px;
  gap: 20px;
  border-radius: 24px;

  margin-top: 20px;
  table {
    width: 100%;
    table-layout: fixed;

    th:nth-child(1),
    td:nth-child(1) {
      width: 10%;
    }

    th:nth-child(2),
    td:nth-child(2) {
      width: 10%;
    }
    th:nth-child(3),
    td:nth-child(3) {
      width: 8%;
    }

    th:nth-child(4),
    td:nth-child(4) {
      width: 8%;
    }
    th:nth-child(5),
    td:nth-child(5) {
      width: 6%;
    }

    th:nth-child(6),
    td:nth-child(6) {
      width: 8%;
    }
    th:nth-child(7),
    td:nth-child(7) {
      width: 2%;
    }
  }
`;

const Heading = styled.h1`
  color: #00225a;
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  text-align: left;
  margin: 10px 0;
`;
const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
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

const FullName = styled.p`
  font-size: 11px;
  font-weight: 400;
  line-height: 16px;
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
