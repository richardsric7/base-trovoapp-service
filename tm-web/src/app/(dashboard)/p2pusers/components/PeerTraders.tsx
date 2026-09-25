import CustomTable from "@/components/CustomTable";
import React, { useMemo } from "react";
import { MdVerified } from "react-icons/md";
import styled from "styled-components";

const PeerTraders = () => {
  const columns = [
    {
      title: "",
      dataIndex: "offer_maker",
      key: "offer_maker",

      render: (_: any, record: any) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserContent>
              <UserName>{record.offer_maker || "N/A"}</UserName>
              <MdVerified color="#007CDF" />
            </UserContent>
            <FullName>{record.fullName || "N/A"}</FullName>
          </div>
        </UserInfoSection>
      ),
    },

    {
      title: "",
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
      trades: "$4,500.00",
    },
    {
      username: "JaneSmith",
      fullName: "Jane Smith",
      trades: "$4,500.00",
    },
    {
      username: "Marlone",
      fullName: "Mark Ovey",
      trades: "$4,500.00",
    },

    {
      username: "JaneS",
      fullName: "Jane Smith",
      trades: "$4,500.00",
    },
    {
      username: "JohnDoe",
      fullName: "John Doe",
      trades: "$4,500.00",
    },
    {
      username: "Marlone",
      fullName: "Mark Ovey",
      trades: "$4,500.00",
    },
  ];

  const dataSource = useMemo(() => {
    return (
      dummyData?.map((item, index) => ({
        key: index + 1,
        offer_maker: item.username,
        fullName: item.fullName,
        trades: item.trades,
      })) ?? []
    );
  }, [dummyData]);

  return (
    <Container>
      <SubTitle>Top Peer Traders</SubTitle>
      <CustomTable columns={columns} dataSource={dataSource} />
    </Container>
  );
};

export default PeerTraders;

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

const Trades = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 16px;
  color: #00225a;
`;
