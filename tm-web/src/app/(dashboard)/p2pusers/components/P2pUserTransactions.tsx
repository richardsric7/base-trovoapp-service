import CustomTable from "@/components/CustomTable";
import Link from "next/link";
import React, { useMemo } from "react";
import { MdVerified } from "react-icons/md";
import styled from "styled-components";

type TableData = {
  username: string;
  fullName: string;
  price: string;
  amount: string;
  total: number;
  tradedOn: string;
};

const P2pUserTransactions = () => {
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

      render: (count: number) => <Text>{count}</Text>,
    },

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
  ];

  // Dummy data for testing
  const dummyData = [
    {
      username: "JohnDoe",
      fullName: "John Doe",
      price: "0.510595 NGN",
      amount: "90,450.00 TROV",
      total: "45,270.41 NGN",
      tradedOn: "23 Sep 2023, 8:17",
    },
    {
      username: "JaneSmith",
      fullName: "Jane Smith",
      price: "0.510595 NGN",
      amount: "90,450.00 TROV",
      total: "45,270.41 NGN",
      tradedOn: "23 Sep 2023, 8:17",
    },
    {
      username: "Marlone",
      fullName: "Mark Ovey",
      price: "0.510595 NGN",
      amount: "90,450.00 TROV",
      total: "45,270.41 NGN",
      tradedOn: "23 Sep 2023, 8:17",
    },

    {
      username: "JaneS",
      fullName: "Jane Smith",
      price: "0.510595 NGN",
      amount: "90,450.00 TROV",
      total: "45,270.41 NGN",
      tradedOn: "23 Sep 2023, 8:17",
    },
  ];

  const dataSource = useMemo(() => {
    return dummyData.map((item, index) => ({
      key: index + 1,
      username: item.username,
      fullName: item.fullName,
      price: item.price,
      amount: item.amount,
      total: item.total,
      tradedOn: item.tradedOn,
    }));
  }, []);
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

export default P2pUserTransactions;

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
