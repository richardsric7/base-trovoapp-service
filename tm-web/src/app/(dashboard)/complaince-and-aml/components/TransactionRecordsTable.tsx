"use client";

import { SearchBar } from "@/components";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import { useRouter } from "next/navigation";
import { useMemo, useState } from "react";
import {
  AiOutlineCheckCircle,
  AiOutlineClockCircle,
  AiOutlineCloseCircle,
} from "react-icons/ai";
import styled from "styled-components";
import AssetRecordsFilter from "./AccessRecordsFilter";

import { BsArrowDownCircle, BsArrowUpCircle } from "react-icons/bs";
import { RiExchangeFill } from "react-icons/ri";
import TransactionsRecordsFilter from "./TransactionsRecordsFilter";

interface UserStatusProps {
  isSuccessful?: boolean;
  isFailed?: boolean;
  isReviewd?: boolean;
}
interface TransactionTypeProps {
  isRecieved?: boolean;
  isSent?: boolean;
  isSwapped?: boolean;
}
const TransactionRecordsTable = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [activeId, setActiveId] = useState<number | null>(null);
  const router = useRouter();

  const dataSource = useMemo(() => {
    return Array.from({ length: 100 }, (_, index) => ({
      id: index + 1,
      transactionId: "TRO23567",

      price: "3,000.000 XBN",

      description: "Received from Obi",
      date: `30 Jun, 2025 08:30 AM `,
      transactionType:
        index % 3 === 0 ? "Received" : index % 3 === 1 ? "Swap" : "Sent",
      status:
        index % 3 === 0
          ? "Cleared"
          : index % 3 === 1
          ? "Flagged"
          : "Under Review",
    }));
  }, []);

  const selectedRecord = dataSource.find((item) => item.id === activeId);

  const columns = [
    {
      title: "Transaction ID",
      dataIndex: "transactionId",
      key: "transactionId",
    },

    {
      title: "Price",
      dataIndex: "price",
      key: "price",
      render: (_: string, record: any) => {
        const { transactionType, price } = record;

        const isRecieved = transactionType === "Received";
        const isSent = transactionType === "Sent";
        const isSwapped = transactionType === "Swap";

        return (
          <Price isRecieved={isRecieved} isSent={isSent} isSwapped={isSwapped}>
            {price}
          </Price>
        );
      },
    },

    {
      title: "Transaction Type",

      dataIndex: "transactionType",
      key: "transactionType",
      render: (transactionType: string) => {
        const isRecieved = transactionType === "Received";

        const isSent = transactionType === "Sent";

        const isSwapped = transactionType === "Swap";

        return (
          <TransactionType
            isRecieved={isRecieved}
            isSent={isSent}
            isSwapped={isSwapped}
          >
            {isRecieved && <BsArrowUpCircle size={18} />}

            {isSent && <BsArrowDownCircle size={18} />}

            {isSwapped && <RiExchangeFill size={18} />}
            {transactionType}
          </TransactionType>
        );
      },
    },

    {
      title: "Description",
      dataIndex: "description",
      key: "description",
    },

    {
      title: "Date",
      dataIndex: "date",
      key: "date",
    },

    {
      title: "AML Status",
      dataIndex: "status",
      key: "status",
      render: (status: string) => {
        const isSuccessful = status === "Cleared";
        const isFailed = status === "Flagged";
        const isReviewd = status === "Under Review";

        return (
          <PaymentStatus
            isSuccessful={isSuccessful}
            isFailed={isFailed}
            isReviewd={isReviewd}
          >
            {isSuccessful && <AiOutlineCheckCircle size={18} />}
            {isReviewd && <AiOutlineClockCircle size={18} />}

            {isFailed && <AiOutlineCloseCircle size={18} />}
            {status}
          </PaymentStatus>
        );
      },
    },
  ];

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);
  return (
    <Container>
      <Header>
        <TitleSection>
          <Title>Transaction Records</Title>
          <SubTitle>Total: 30,230</SubTitle>
        </TitleSection>

        <RightSection>
          <SecondaryButton>Generate Report</SecondaryButton>

          <PrimaryButton>Export Report</PrimaryButton>
        </RightSection>
      </Header>
      <SubHeader>
        <SearchBar />
        <TransactionsRecordsFilter />
      </SubHeader>
      <TableContainer>
        <CustomTable
          columns={columns}
          dataSource={paginatedData}
          onRowClick={(record) =>
            router.push(`/complaince-and-aml/${record.id}`)
          }
        />
        <Pagination
          currentPage={currentPage}
          totalCount={dataSource.length}
          pageSize={pageSize}
          onPageChange={setCurrentPage}
          onPageSizeChange={setPageSize}
        />
      </TableContainer>
    </Container>
  );
};

export default TransactionRecordsTable;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;

  border-radius: 24px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 2px;
  color: #00225a;
`;

const SubTitle = styled.p`
  font-weight: 400;
  font-size: 16px;
  color: #828282;
`;

const TitleSection = styled.div``;

const RightSection = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  width: 35%;
`;

const SubHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 24px;
`;

const TableContainer = styled.div`
  margin-top: 12px;

  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 8%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 12%;
    text-align: center;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 8%;
    text-align: center;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 10%;
    text-align: center;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 12%;
    text-align: center;
  }

  th:nth-child(6),
  td:nth-child(6) {
    width: 10%;
    text-align: center;
  }

  th:nth-child(7),
  td:nth-child(7) {
    width: 4%;
    text-align: center;
  }
`;

const Price = styled.p<TransactionTypeProps>`
  color: ${({ isRecieved, isSent, isSwapped }) =>
    isRecieved
      ? "#00A859"
      : isSent
      ? "#BE3800"
      : isSwapped
      ? "#00225A"
      : "#6B7280"};
  font-weight: 600;
  text-align: center;
  margin: 0;
`;

const TransactionType = styled.p<TransactionTypeProps>`
  color: ${({ isRecieved, isSent, isSwapped }) =>
    isRecieved
      ? "#00A859"
      : isSent
      ? "#BE3800"
      : isSwapped
      ? "#00225A"
      : "#6B7280"};

  background-color: transparent; /* ✅ removes background entirely */
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 0;
  border-radius: 0;
  text-align: center;
  margin: 0 auto;
  text-transform: capitalize;

  svg {
    color: inherit;
  }
`;

const PaymentStatus = styled.p<UserStatusProps>`
  color: ${({ isSuccessful, isFailed, isReviewd }) =>
    isSuccessful
      ? "#00A859"
      : isFailed
      ? "#BE3800"
      : isReviewd
      ? "#007CDF"
      : "#6B7280"};
  background-color: ${({ isSuccessful, isFailed, isReviewd }) =>
    isSuccessful
      ? "#00A8591A"
      : isFailed
      ? "#BE38001A"
      : isReviewd
      ? "#F2F6F9"
      : "#E5E7EB"};
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border-radius: 8px;
  text-align: center;
  margin: 0 auto;
  text-transform: capitalize;
`;
