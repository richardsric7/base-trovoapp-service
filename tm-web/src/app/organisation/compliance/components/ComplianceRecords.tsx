"use client";
import React, { useMemo, useState } from "react";
import styled from "styled-components";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import { StatusBadge } from "@/components";
import { useGetOrgComplianceQuery } from "@/redux/api/orgCompliance";
import { ComplianceRequirementInstance } from "@/redux/api/compliance/interface";

const formatDate = (value?: string) =>
  value ? new Date(value).toLocaleDateString() : "N/A";

const ComplianceRecords = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);

  const { data, isLoading, isFetching, isError, refetch } =
    useGetOrgComplianceQuery({ page: currentPage, limit: pageSize });

  const records = useMemo(
    () => (Array.isArray(data?.data?.records) ? data.data.records : []),
    [data],
  );
  const meta = data?.data?.meta;

  const columns = useMemo(
    () => [
      {
        title: "Category",
        dataIndex: "category",
      },
      {
        title: "Requirement",
        dataIndex: "requirement",
      },
      {
        title: "Due Date",
        dataIndex: "due_date",
        render: (value: ComplianceRequirementInstance["due_date"]) =>
          formatDate(value),
      },
      {
        title: "Status",
        dataIndex: "status",
        render: (status: ComplianceRequirementInstance["status"]) => (
          <StatusBadge status={status} />
        ),
      },
      {
        title: "Completed At",
        dataIndex: "completed_at",
        render: (value: ComplianceRequirementInstance["completed_at"]) =>
          formatDate(value),
      },
    ],
    [],
  );

  if (isError) {
    return (
      <Container>
        <Title>Compliance</Title>
        <StateMessage>
          Unable to load compliance records.
          <RetryButton onClick={() => refetch()}>Try again</RetryButton>
        </StateMessage>
      </Container>
    );
  }

  return (
    <Container>
      <Header>
        <Title>Compliance</Title>
        <Count>{meta?.total ?? 0} records</Count>
      </Header>

      <CustomTable
        columns={columns}
        dataSource={records}
        isLoading={isLoading || isFetching}
        totalItems={meta?.total ?? 0}
        pageSize={pageSize}
      />

      {(meta?.total ?? 0) > 0 && (
        <Pagination
          totalCount={meta?.total ?? 0}
          currentPage={currentPage}
          pageSize={pageSize}
          onPageChange={setCurrentPage}
          onPageSizeChange={setPageSize}
        />
      )}
    </Container>
  );
};

export default ComplianceRecords;

const Container = styled.div`
  background: #ffffff;
  border-radius: 24px;
  padding: 32px;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin: 0;
`;

const Count = styled.span`
  background: #f2f6f9;
  color: #007cdf;
  font-size: 12px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 999px;
`;

const StateMessage = styled.div`
  color: #00225a;
`;

const RetryButton = styled.button`
  display: block;
  margin-top: 12px;
  padding: 8px 16px;
  border: 0;
  border-radius: 8px;
  background: #007cdf;
  color: #fff;
  cursor: pointer;
`;
