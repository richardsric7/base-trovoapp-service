"use client";

import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import { useGetAuditTrailQuery } from "@/redux/api/auditTrail/api";
import { IAuditRecord, IAuditTrailParams } from "@/redux/api/auditTrail/interface";
import { useState } from "react";
import { AiOutlineCheckCircle, AiOutlineCloseCircle } from "react-icons/ai";
import { MdVerified } from "react-icons/md";
import styled from "styled-components";
import ActivityModal from "./AcitvityModal";

interface UserStatusProps {
  isActive?: boolean;
  isFailed?: boolean;
}

interface AuditTableProps {
  filters: IAuditTrailParams;
}

const AuditTable = ({ filters }: AuditTableProps) => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [activeId, setActiveId] = useState<string | null>(null);

  const { data, isFetching, isError } = useGetAuditTrailQuery({
    ...filters,
    page: currentPage,
    pageSize,
  });

  const rows: IAuditRecord[] = data?.data?.data ?? [];
  const total = data?.data?.total ?? 0;

  const columns = [
    {
      title: "Users",
      dataIndex: "username",
      render: (_: any, record: IAuditRecord) => (
        <UserInfoSection onClick={() => setActiveId(record.id)}>
          <Avatar />
          <div>
            <UserNameWrapper>
              <UserName>{record.username || "—"}</UserName>
              <MdVerified color="#007cdf" />
            </UserNameWrapper>
            <FullName>{record.fullname}</FullName>
          </div>
        </UserInfoSection>
      ),
    },
    {
      title: "Action",
      dataIndex: "action",
      render: (_: any, record: IAuditRecord) => (
        <Actionwrapper onClick={() => setActiveId(record.id)}>
          {record.action || "—"}
        </Actionwrapper>
      ),
    },
    {
      title: "Date/Time",
      dataIndex: "date",
      key: "date",
    },
    {
      title: "Ip Address",
      dataIndex: "ip_address",
      key: "ip_address",
      render: (ip: string) => ip || "—",
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (status: string) => {
        const isSuccessful = status === "successful";
        const isFailed = status === "failed";

        return (
          <PaymentStatus isActive={isSuccessful} isFailed={isFailed}>
            {isSuccessful && <AiOutlineCheckCircle size={18} />}
            {isFailed && <AiOutlineCloseCircle size={18} />}
            {status}
          </PaymentStatus>
        );
      },
    },
  ];

  return (
    <TableContainer>
      {isError ? (
        // On error, show ONLY the error — never alongside the table's empty state.
        <EmptyState>Could not load the audit trail. Try again.</EmptyState>
      ) : (
        // CustomTable owns the loading and genuine empty-success states.
        <CustomTable columns={columns} dataSource={rows} isLoading={isFetching} />
      )}

      <Pagination
        currentPage={currentPage}
        totalCount={total}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />

      {activeId && (
        <ActivityModal
          id={activeId}
          isOpen={!!activeId}
          onClose={() => setActiveId(null)}
        />
      )}
    </TableContainer>
  );
};

export default AuditTable;

const TableContainer = styled.div`
  margin-top: 12px;

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
    width: 14%;
    text-align: center;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 16%;
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
`;

const EmptyState = styled.div`
  padding: 32px;
  text-align: center;
  color: #828282;
  font-size: 14px;
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
  color: #00225a;
  font-family: inherit;
  cursor: pointer;
  text-decoration: underline;
`;

const PaymentStatus = styled.p<UserStatusProps>`
  color: ${({ isActive, isFailed }) =>
    isActive ? "#00A859" : isFailed ? "#BE3800" : "#6B7280"};
  background-color: ${({ isActive, isFailed }) =>
    isActive ? "#00A8591A" : isFailed ? "#BE38001A" : "#E5E7EB"};
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
