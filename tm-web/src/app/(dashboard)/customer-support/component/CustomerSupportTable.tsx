"use client";

import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import { useRouter } from "next/navigation";
import { useMemo, useState } from "react";
import { MdVerified } from "react-icons/md";
import styled from "styled-components";

const CustomerSupportTable = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const router = useRouter();
  const dataSource = useMemo(() => {
    return Array.from({ length: 100 }, (_, index) => ({
      id: `T-00${index + 1}`,
      reportedBy: index % 2 === 0 ? "Kennis Moukla" : "Florence Zach",
      title: "Failed Asset Tokenization",
      reportedOn: "23 Sep 2023, 8:17",
      lastUpdated: "23 Sep 2023, 8:17",
      issueType: "Billing",
      urgency: index % 3 === 0 ? "High" : index % 3 === 1 ? "Medium" : "Normal",
      status:
        index % 3 === 0 ? "Pending" : index % 3 === 1 ? "Open" : "Resolved",
    }));
  }, []);

  const columns = [
    {
      title: "Ticket ID",
      dataIndex: "id",
    },
    {
      title: "Reported By",
      dataIndex: "reportedBy",
      render: (name: string) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserNameWrapper>
              <UserName>Kennis</UserName>
              <MdVerified color="#007cdf" />
            </UserNameWrapper>
            <FullName>Florence Zach</FullName>
          </div>
        </UserInfoSection>
      ),
    },
    {
      title: "Title",
      dataIndex: "title",
    },
    {
      title: "Reported On",
      dataIndex: "reportedOn",
    },
    {
      title: "Last Updated On",
      dataIndex: "lastUpdated",
    },
    {
      title: "Issue Type",
      dataIndex: "issueType",
    },
    {
      title: "Urgency",
      dataIndex: "urgency",
      render: (urgency: string) => (
        <UrgencyTag level={urgency}>
          <span className="dot" />
          {urgency}
        </UrgencyTag>
      ),
    },
    {
      title: "Status",
      dataIndex: "status",
      render: (status: string) => (
        <StatusPill status={status}>{status}</StatusPill>
      ),
    },
  ];

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);
  return (
    <TableContainer>
      <CustomTable
        columns={columns}
        dataSource={paginatedData}
        onRowClick={(record) => {
          if (record?.id) {
            router.push(`/customer-support/${record.id}`);
          }
        }}
      />
      <Pagination
        currentPage={currentPage}
        totalCount={dataSource.length}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />
    </TableContainer>
  );
};

export default CustomerSupportTable;

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
    width: 14%;
    text-align: center;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 12%;
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
    width: 12%;
    text-align: center;
  }

  th:nth-child(7),
  td:nth-child(7) {
    width: 12%;
    text-align: center;
  }
  th:nth-child(8),
  td:nth-child(8) {
    width: 12%;
    text-align: center;
  }
  th:nth-child(9),
  td:nth-child(9) {
    width: 12%;
    text-align: center;
  }
`;

const UserInfoSection = styled.div`
  display: flex;
  align-items: center;
  gap: 2px;
`;

const Avatar = styled.div`
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #007cdf;
`;

const UserNameWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 5px;
`;

const UserName = styled.p`
  font-weight: 600;
  font-size: 14px;
  margin: 0;
`;

const FullName = styled.p`
  font-size: 12px;
  color: #828282;
  margin: 0;
`;

const UrgencyTag = styled.div<{ level: string }>`
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 500;
  color: ${({ level }) =>
    level === "High" ? "#BE3800" : level === "Medium" ? "#FFCC00" : "#62ACE8"};

  .dot {
    width: 16px;
    height: 16px;
    border-radius: 4px;
    background-color: currentColor;
  }
`;

const StatusPill = styled.div<{ status: string }>`
  background-color: ${({ status }) =>
    status === "Pending"
      ? "#F2F6F9"
      : status === "Resolved"
      ? "#00A8591A"
      : "#BE38001A"};
  color: ${({ status }) =>
    status === "Pending"
      ? "#007CDF"
      : status === "Resolved"
      ? "#00A859"
      : "#BE3800"};
  padding: 5px 10px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  display: inline-block;
  text-align: center;
`;
