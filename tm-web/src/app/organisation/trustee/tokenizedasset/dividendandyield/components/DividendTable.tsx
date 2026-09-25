"use client";

import { SearchBar } from "@/components";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import React, { useMemo, useState } from "react";
import styled from "styled-components";
import AssetCurationFilter from "@/app/(dashboard)/assetcuration/components/AssetCurationFilter";
import { useRouter } from "next/navigation";
import PrimaryButton from "@/components/PrimaryButton";
import { AiOutlineCheckCircle, AiOutlineCloseCircle } from "react-icons/ai";
import PayDividendModal from "./PayDividendModal";

const DividendTable = () => {
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [isPay, setIsPay] = useState(false);
  const columns = [
    {
      title: "Payout Date",
      dataIndex: "payoutDate",
      key: "payoutDate",
    },

    {
      title: "Payout Ref",
      dataIndex: "payoutRef",
      key: "payoutRef",
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

    {
      title: "Action",
      dataIndex: "",
      render: (_: any, record: any) => {
        return (
          <Actionwrapper
            onClick={() =>
              router.push(
                `/organisation/trustee/tokenizedasset/dividendandyield/${record.id}`
              )
            }
          >
            View
          </Actionwrapper>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 100 }, (_, index) => ({
      id: index + 1,
      payoutDate: "23 Sep 2023",
      payoutRef: `TR123456`,
      tokenHolders: `1200`,
      payoutPerToken: `₦20.00`,
      grossAmount: `₦100,000,000.00`,
      tax: "₦5,000,000.00",
      netAmount: `₦95,000,000.00`,
      status: "Processed",
    }));
  }, []);

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);

  return (
    <Container>
      <Header>
        <Title>Dividend Payout History</Title>
        <PrimaryButton
          buttonStyle={{ width: "14%", marginRight: "-6px" }}
          onClick={() => setIsPay(true)}
        >
          Pay Dividend
        </PrimaryButton>
      </Header>

      <SubHeader>
        <SearchBar />
        <AssetCurationFilter />
      </SubHeader>
      <CustomTable columns={columns} dataSource={paginatedData} />
      <Pagination
        currentPage={currentPage}
        totalCount={dataSource.length}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />
      <PayDividendModal isOpen={isPay} onClose={() => setIsPay(false)} />
    </Container>
  );
};

export default DividendTable;

const Container = styled.section`
  background: #ffffff;
  padding: 4px 20px 10px 20px;

  border-radius: 24px;

  table {
    width: 100%;
    table-layout: fixed;
    padding-top: 8px;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 7%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 8%;
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
    width: 10%;
    text-align: center;
  }
  th:nth-child(8),
  td:nth-child(8) {
    width: 10%;
    text-align: center;
  }
  th:nth-child(9),
  td:nth-child(9) {
    width: 5%;
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
const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
  gap: 86px;
`;

const SubHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 18px;
`;
const Actionwrapper = styled.button`
  background-color: transparent;
  border-radius: 8px;
  border: 1px solid #007cdf;
  padding: 10px 24px;
  color: #007cdf;
  font-family: inherit;
  cursor: pointer;
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

const StatusWithIcon = ({ status }: { status: "Processed" | "Failed" }) => {
  const isProcessed = status === "Processed";

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
