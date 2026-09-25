"use client";

import { SearchBar } from "@/components";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";

import React, { useMemo, useState } from "react";
import { AiFillClockCircle } from "react-icons/ai";
import styled from "styled-components";

interface PayoutStatusProps {
  isActive?: boolean;
}

const PayoutList = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const totalUsers = 100;

  const columns = [
    {
      title: "Investor",
      dataIndex: "investor",
      render: (_: any, record: any) => {
        return (
          <UserIdSection>
            <Avatar />
            <div>
              <UserName>{record.name}</UserName>
              <UserSubName>{record.subname}</UserSubName>
            </div>
          </UserIdSection>
        );
      },
    },
    {
      title: "Wallet /Bank Details",
      dataIndex: "wallet",
      key: "wallet",
    },

    {
      title: "Tokens Held",
      dataIndex: "tokensHeld",
      key: "tokensHeld",
    },
    {
      title: "Amount",
      dataIndex: "amount",
      key: "amount",
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (_: any, record: any) => (
        <StatusPill>
          <AiFillClockCircle />
          {record.status}
        </StatusPill>
      ),
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 10 }, (_, index) => ({
      userId: index + 1,
      name: `Florence ${index + 1}`,
      subname: `florence Zach`,
      amount: "₦20,000.00",
      tokensHeld: "1000",
      wallet: "12345678-First Bank",
      status: "Pending",
    }));
  }, []);
  return (
    <Container>
      <Title>Payout List</Title>
      <HeaderRow>
        <SearchBar />
      </HeaderRow>
      <div>
        <CustomTable columns={columns} dataSource={dataSource} />
        <Pagination
          totalCount={totalUsers}
          currentPage={currentPage}
          pageSize={pageSize}
          onPageChange={setCurrentPage}
          onPageSizeChange={setPageSize}
        />
      </div>
    </Container>
  );
};

export default PayoutList;

const Container = styled.section`
  background-color: #fff;
  gap: 20px;
  border-radius: 16px;
  padding: 24px;
  // margin-top: 10px;
  width: 100%;

  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 18%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 19%;
    text-align: center;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 10%;
    text-align: center;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 14%;
    text-align: center;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 16%;
    text-align: center;
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

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  margin-bottom: 10px;
  color: #00225a;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;
const HeaderRow = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
`;

const UserIdSection = styled.div`
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

const UserSubName = styled.p`
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const PayoutStatus = styled.p<PayoutStatusProps>`
  color: ${(props) => (props.isActive ? "#00A859" : "#BE3800")};
  background-color: ${(props) => (props.isActive ? "#00A8591A" : "#BE38001A")};
  font-weight: 500;
  display: inline-flex;
  padding: 8px;
  border-radius: 8px;
  text-align: center;
  margin: 0 auto;
`;
const StatusPill = styled.div`
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #eaf3ff;
  color: #007cdf;
  padding: 6px 10px;
  border-radius: 999px;
  font-size: 12px;
`;
