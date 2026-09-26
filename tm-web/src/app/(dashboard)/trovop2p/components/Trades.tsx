"use client";
import CustomTable from "@/components/CustomTable";
import { useTradesQuery } from "@/redux/api/p2p";
import Link from "next/link";

import React, { useMemo, useState } from "react";
import styled from "styled-components";

const Trades = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const { data, isLoading } = useTradesQuery({
    page: currentPage,
    pageSize,
  });
  
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
    const tradesArray = data?.data?.data.slice(0, 6);

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

  return (
    <Container>
      <Header>
        {" "}
        <SubTitle>Trades</SubTitle>
        <MoreBtn href="/trovop2p/alltrades">See all</MoreBtn>
      </Header>

      <CustomTable columns={columns} dataSource={dataSource} />
    </Container>
  );
};

export default Trades;

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
      //   width: 25%;
    }

    th:nth-child(2),
    td:nth-child(2) {
      //   width: 20%;
    }
    th:nth-child(3),
    td:nth-child(3) {
      //   width: 15%;
    }
  }
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;
const SubTitle = styled.h3`
  font-size: 20px;
  font-weight: 600;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;
const MoreBtn = styled(Link)`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  color: #007cdf;
  text-decoration: none;
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
