"use client";

import { SearchBar } from "@/components";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import React, { useMemo, useState } from "react";
import { AiOutlineCheckCircle, AiOutlineCloseCircle } from "react-icons/ai";
import styled from "styled-components";
import AdminFilter from "../../admin-users/components/AdminFilter";
import { MdVerified } from "react-icons/md";
import TruncatedText from "@/hooks/useTruncate";

const PayoutListTable = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const columns = [
    {
      title: "Investor",
      dataIndex: "investor",
      render: (_: any, record: any) => {
        return (
          <InvestorInfoSection>
            <Avatar />
            <div>
              <UserNameWrapper>
                <UserName>
                  {" "}
                  <TruncatedText text={record.userName} maxLength={15} />
                  {record.kyc_level === 0 ? (
                    ""
                  ) : (
                    <MdVerified size={12} color="#007CDF" />
                  )}
                </UserName>
              </UserNameWrapper>

              <FullName>{record.fullName}</FullName>
            </div>
          </InvestorInfoSection>
        );
      },
    },

    {
      title: "Wallet/Bank Details",
      dataIndex: "walletDetails",
      key: "walletDetails",
    },

    {
      title: "Token Holders",
      dataIndex: "tokenHolders",
      key: "tokenHolders",
    },

    {
      title: "Payout Per Token",
      dataIndex: "payoutPerToken",
      key: "payoutPerToken",
    },

    {
      title: " Gross Amount",
      dataIndex: "grossAmount",
      key: "grossAmount",
    },
    {
      title: "Tax(WHT)",
      dataIndex: "tax",
      key: "tax",
    },

    {
      title: " Net Amount",
      dataIndex: "netAmount",
      key: "netAmount",
    },

    {
      title: " Status",
      dataIndex: "status",
      key: "status",
      render: (_: any, record: any) => (
        <StatusWithIcon status={record.status} />
      ),
    },
  ];
  const dataSource = useMemo(() => {
    return Array.from({ length: 100 }, (_, index) => ({
      key: index + 1,
      payoutDate: "23 Sep 2023",
      userName: "Odogwu",
      fullName: "Obi Enechi",
      walletDetails: `12345678-First Bank`,
      tokenHolders: `1200`,
      payoutPerToken: `₦20.00`,
      grossAmount: `₦100,000,000.00`,
      tax: "₦5,000,000.00",
      netAmount: `₦95,000,000.00`,
      status: "Paid",
    }));
  }, []);

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);

  return (
    <Container>
      <Title> Payout List</Title>

      <SubHeader>
        <SearchBar />
        <AdminFilter />
      </SubHeader>
      <CustomTable columns={columns} dataSource={paginatedData} />
      <Pagination
        currentPage={currentPage}
        totalCount={dataSource.length}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />
    </Container>
  );
};

export default PayoutListTable;
const Container = styled.section`
  background: #ffffff;
  padding: 32px;
  margin-top: 30px;
  border-radius: 24px;

  table {
    width: 100%;
    table-layout: fixed;
    padding-top: 8px;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 14%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 12%;
    text-align: center;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 7%;
    text-align: center;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 8%;
    text-align: center;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 9%;
    text-align: center;
  }
  th:nth-child(6),
  td:nth-child(6) {
    width: 12%;
    text-align: center;
  }

  th:nth-child(7),
  td:nth-child(7) {
    width: 8%;
    text-align: center;
  }

  th:nth-child(8),
  td:nth-child(8) {
    width: 8%;
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
  margin: 0;
  color: #00225a;
`;

const SubHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 18px;
`;

const InvestorInfoSection = styled.div`
  display: flex;
  column-gap: 10px;
  align-items: center;
  cursor: pointer;
`;
const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const UserNameWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
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
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;
const AssetStatus = styled.p<{ textColor: string; backgroundColor: string }>`
  color: ${({ textColor }) => textColor};
  background-color: ${({ backgroundColor }) => backgroundColor};
  font-weight: 500;
  padding: 8px;
  border-radius: 8px;
  text-align: center;
  margin: 0 auto;
  display: inline-flex;
  align-items: center;
  gap: 4px;
`;

const StatusWithIcon = ({ status }: { status: "Paid" | "Failed" }) => {
  const isProcessed = status === "Paid";

  return (
    <AssetStatus
      textColor={isProcessed ? "#00A859" : "#FF4D4D"}
      backgroundColor={isProcessed ? "#00A8591A" : "#BE38001A"}
    >
      {isProcessed ? (
        <AiOutlineCheckCircle size={14} />
      ) : (
        <AiOutlineCloseCircle size={14} />
      )}
      {status}
    </AssetStatus>
  );
};
