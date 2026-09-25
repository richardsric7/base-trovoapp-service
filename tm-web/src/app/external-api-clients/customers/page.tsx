"use client";

import { SearchBar } from "@/components";
import { useMemo } from "react";
import styled from "styled-components";
import CustomerFilter from "../components/CustomerFilter";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import {
  AiOutlineClockCircle,
  AiOutlineCheckCircle,
  AiOutlineCloseCircle,
} from "react-icons/ai";
import { useState } from "react";
import { useRouter } from "next/navigation";

interface CustomersStatusProps {
  isActive?: boolean;
}

const CutomersPage = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const router = useRouter();
  const totalUsers = 100;
  const [filters, setFilters] = useState({});

  const columns = [
    {
      title: "Customer",
      dataIndex: "customer",
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
      title: "Email Address",
      dataIndex: "email",
      key: "email",
    },

    {
      title: "Phone Number",
      dataIndex: "phoneNumber",
      key: "phoneNumber",
    },
    {
      title: "Registration Date",
      dataIndex: "registrationDate",
      key: "registrationDate",
    },
    {
      title: "Kyc Status",
      dataIndex: "kycStatus",
      key: "kycStatus",
      render: (_: any, record: any) => {
        const status: string = record.kycStatus || "Pending";
        const key = status.toLowerCase();
        let Icon = AiOutlineClockCircle;
        if (key === "verified") Icon = AiOutlineCheckCircle;
        if (key === "rejected" || key === "reject") Icon = AiOutlineCloseCircle;

        return (
          <StatusPill status={key}>
            <Icon />
            {status}
          </StatusPill>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 10 }, (_, index) => ({
      userId: index + 1,
      name: `Florence ${index + 1}`,
      subname: `florence Zach`,
      email: `florence${index + 1}@example.com`,
      phoneNumber: "+234 80123456789",
      registrationDate: "23 Sep 2025",
      kycStatus: [
        "Pending",
        "Rejected",
        "Verified",
        "Unverified",
        "Active",
        "Inactive",
      ][index % 3],
    }));
  }, []);
  return (
    <Container>
      <Title>Customers</Title>
      <HeaderRow>
        <SearchBar />
        <CustomerFilter filters={filters} onApply={setFilters} />
      </HeaderRow>
      <div>
        <CustomTable
          columns={columns}
          dataSource={dataSource}
          onRowClick={(record) =>
            router.push(`/external-api-clients/customers/${record.userId}`)
          }
        />
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

export default CutomersPage;

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

const StatusPill = styled.div<{ status?: string }>`
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  svg {
    display: inline-block;
  }
  color: ${(p) =>
    p.status === "verified"
      ? "#00A859"
      : p.status === "rejected"
        ? "#BE3800"
        : "#007CDF"};
  background: ${(p) =>
    p.status === "verified"
      ? "#E6FFF0"
      : p.status === "rejected"
        ? "#BE38001A"
        : "#F2F6F9"};
`;
