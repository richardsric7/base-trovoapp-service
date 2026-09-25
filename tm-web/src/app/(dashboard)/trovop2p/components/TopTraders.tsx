"use client";
import CustomTable from "@/components/CustomTable";
import { ITopTrader, useTradeStatisticsQuery } from "@/redux/api/p2p";
import React, { useMemo } from "react";
// import { MdVerified } from "react-icons/md";
import styled from "styled-components";

const TopTraders = () => {
  const { data, isLoading } = useTradeStatisticsQuery();

  const columns = [
    {
      title: "User",
      dataIndex: "offer_maker",
      key: "offer_maker",

      render: (_: any, record: ITopTrader) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserContent>
              <UserName>{record.offer_maker || "N/A"}</UserName>
              {/* <MdVerified color="#007CDF" /> */}
            </UserContent>
            {/* <FullName>{record.fullName || "N/A"}</FullName> */}
          </div>
        </UserInfoSection>
      ),
    },
    // {
    //   title: "Public Key",
    //   dataIndex: "publicKey",
    //   key: "publicKey",

    //   render: (key: string) => <PublicKey>{key || "N/A"}</PublicKey>,
    // },
    {
      title: "Total Trades",
      dataIndex: "trades",
      key: "trades",

      render: (count: number) => <Trades>{count}</Trades>,
    },
  ];

  // Dummy data for testing
  const dummyData = [
    {
      username: "JohnDoe",
      fullName: "John Doe",
      publicKey: "CD345678SF",
      trades: "$4,500.00",
    },
    {
      username: "JaneSmith",
      fullName: "Jane Smith",
      publicKey: "CD345678SF",
      trades: "$4,500.00",
    },
    {
      username: "Marlone",
      fullName: "Mark Ovey",
      publicKey: "CD345678SF",
      trades: "$4,500.00",
    },

    {
      username: "JaneS",
      fullName: "Jane Smith",
      publicKey: "CD345678SF",
      trades: "$4,500.00",
    },
    {
      username: "JohnDoe",
      fullName: "John Doe",
      publicKey: "CD345678SF",
      trades: "$4,500.00",
    },
    {
      username: "Marlone",
      fullName: "Mark Ovey",
      publicKey: "CD345678SF",
      trades: "$4,500.00",
    },
  ];

  const dataSource = useMemo(() => {
    return (
      data?.top_traders?.map((item, index) => ({
        key: index + 1,
        offer_maker: item.offer_maker,
        trades: item.order_count,
      })) ?? []
    );
  }, [data]);

  if (isLoading) return <p>Loading top traders...</p>;

  return (
    <Container>
      <SubTitle>Top Traders</SubTitle>
      <CustomTable columns={columns} dataSource={dataSource} />
    </Container>
  );
};

export default TopTraders;

const Container = styled.div`
  background-color: #fff;
  padding: 40px 32px;
  gap: 20px;
  border-radius: 24px;
  width: 40%;
  margin-top: 20px;
  table {
    width: 100%;
    table-layout: fixed;

    th:nth-child(1),
    td:nth-child(1) {
      width: 75%;
    }

    th:nth-child(2),
    td:nth-child(2) {
      width: 25%;
      text-align: center;
    }
  }
`;

const SubTitle = styled.h3`
  font-size: 20px;
  font-weight: 600;
  line-height: 28px;
  margin: 0;
  color: #00225a;
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

const PublicKey = styled.div`
  font-size: 14px;
  font-weight: 500;
  line-height: 16px;
  color: #007cdf;
`;

const Trades = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 16px;
  color: #00225a;
`;
