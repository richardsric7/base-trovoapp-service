"use client";

import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import { useRouter } from "next/navigation";
import React, { useMemo, useState } from "react";
import { MdVerified } from "react-icons/md";
import styled from "styled-components";
import {
  AiOutlineCheckCircle,
  AiOutlineClockCircle,
  AiOutlineCloseCircle,
} from "react-icons/ai";
import PaymentModal from "./PaymentModal";
import { useGetFiatPaymentsQuery } from "@/redux/api/payment";
import { GetFiatPaymentsParams } from "@/redux/api/payment";

interface PaymentTableProps {
  search: string;
  filters: Partial<GetFiatPaymentsParams>;
}

interface UserStatusProps {
  isActive?: boolean;
  isProcessing?: boolean;
  isFailed?: boolean;
}

const PaymentTable: React.FC<PaymentTableProps> = ({ search, filters }) => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const { data, isLoading, isFetching, error } = useGetFiatPaymentsQuery({
    page: currentPage,
    page_size: pageSize,
    username: search || undefined,
    ...filters,
  });
  const [activeId, setActiveId] = useState<string | null>(null);
  const router = useRouter();

  type UiStatus = "successful" | "pending";

  function deriveStatus(record: any): UiStatus {
    if (record.record_type === "invoice") {
      return record.status === "COMPLETED" ? "successful" : "pending";
    }

    // payment → money already move
    return "successful";
  }

  const dataSource = useMemo(() => {
    return (
      data?.data?.data?.map((record) => ({
        id: record.id,
        username: record.username,
        amount: record.amount,
        paymentType:
          record.payment_type.charAt(0).toUpperCase() +
          record.payment_type.slice(1).toLowerCase(),
        date: record.created_at,
        status: record.record_type === "invoice" ? record.status : "COMPLETED",
        serviceProvider: record.service_provider,
        raw: record,
      })) ?? []
    );
  }, [data]);

  const selectedRecord = dataSource.find((item) => item.id === activeId);

  const columns = [
    {
      title: "Amount",
      dataIndex: "amount",
      render: (amount: string) =>
        `${Number(amount).toLocaleString("en-US", {
          minimumFractionDigits: 2,
        })} NGN`,
      key: "amount",
    },
    {
      title: "Users",
      dataIndex: "username",
      render: (_: any, record: any) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserNameWrapper>
              <UserName>{record.username}</UserName>
              {/* <MdVerified color="#007cdf" /> */}
            </UserNameWrapper>
            {/* <FullName>{record.fullname}</FullName> */}
          </div>
        </UserInfoSection>
      ),
    },
    {
      title: "Payment Type",
      dataIndex: "paymentType",
      key: "paymentType",
    },
    {
      title: "Date",
      dataIndex: "date",

      render: (created_at: string) =>
        new Date(created_at).toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        }),

      key: "date",
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (status: string) => {
        const isSuccessful = status === "COMPLETED";
        const isProcessing = status === "PENDING";
        const isFailed = status === "FAILED";

        return (
          <PaymentStatus
            isActive={isSuccessful}
            isProcessing={isProcessing}
            isFailed={isFailed}
          >
            {isSuccessful && <AiOutlineCheckCircle size={18} />}
            {isProcessing && <AiOutlineClockCircle size={18} />}
            {isFailed && <AiOutlineCloseCircle size={18} />}
            {status?.toLowerCase()}
          </PaymentStatus>
        );
      },
    },
    {
      title: "Action",
      dataIndex: "",
      render: (_: any, record: any) => (
        <Actionwrapper onClick={() => setActiveId(record.id)}>
          View
        </Actionwrapper>
      ),
    },
  ];

  if (isLoading || isFetching) {
    return (
      <TableContainer>
        <UserName>Loading...</UserName>
      </TableContainer>
    );
  }

  if (dataSource.length === 0) {
    return (
      <TableContainer>
        <UserName style={{ textAlign: "center", padding: "40px" }}>
          No search or filter results found
        </UserName>
      </TableContainer>
    );
  }
  return (
    <TableContainer>
      <CustomTable columns={columns} dataSource={dataSource} />
      <Pagination
        currentPage={currentPage}
        totalCount={data?.data.total ?? 0}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />

      {selectedRecord && (
        <PaymentModal
          isOpen={!!selectedRecord}
          onClose={() => setActiveId(null)}
          record={selectedRecord}
        />
      )}
    </TableContainer>
  );
};

export default PaymentTable;

const TableContainer = styled.div`
  margin-top: 12px;
`;

const UserNameWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
  max-width: 40px;
`;

const UserName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
`;

const UserInfoSection = styled.div`
  display: flex;
  column-gap: 10px;
  align-items: center;
  cursor: pointer;
`;

const FullName = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  margin: 0;
  text-transform: capitalize;
`;

const Actionwrapper = styled.button`
  background-color: transparent;
  border: none;
  color: #007cdf;
  font-family: inherit;
  cursor: pointer;
  text-decoration: underline;
`;

const PaymentStatus = styled.p<UserStatusProps>`
  color: ${({ isActive, isProcessing, isFailed }) =>
    isActive
      ? "#00A859"
      : isProcessing
        ? "#007CDF"
        : isFailed
          ? "#BE3800"
          : "#6B7280"};
  background-color: ${({ isActive, isProcessing, isFailed }) =>
    isActive
      ? "#00A8591A"
      : isProcessing
        ? "#F2F6F9"
        : isFailed
          ? "#BE38001A"
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
