"use client";
import { SearchBar } from "@/components";
import { useRouter } from "next/navigation";
import React, { useMemo, useState } from "react";
import { FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";
import AssetCurationFilter from "../assetcuration/components/AssetCurationFilter";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import PrimaryButton from "@/components/PrimaryButton";

const AssetIncomeReportPage = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const router = useRouter();

  const handleBack = () => {
    router.back();
  };

  const columns = [
    {
      title: "Title",
      dataIndex: "title",
      key: "title",
    },
    {
      title: "Start Date",
      dataIndex: "startDate",
      key: "startDate",
    },
    {
      title: "endDate",
      dataIndex: "endDate",
      key: "endDate",
    },

    {
      title: "Circulated Tokens",
      dataIndex: "circulatedTokens",
      key: "circulatedTokens",
    },

    {
      title: "Revenue Generated",
      dataIndex: "revenueGenerated",
      key: "revenueGenerated",
    },
    {
      title: "Action",
      dataIndex: "",
      render: (_: any, record: any) => {
        return (
          <Actionwrapper
            onClick={() => router.push(`/assetincomereport/${record.id}`)}
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
      title: "June 2025",
      startDate: `01 Jun 2025`,
      endDate: `30 Jun 2025`,
      circulatedTokens: "50,000",
      revenueGenerated: "₦400,000,000.00",
    }));
  }, []);

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);

  return (
    <>
      <BackButtonLink onClick={handleBack}>
        <FaArrowLeft size={18} />
      </BackButtonLink>
      <Container>
        <Header>
          <Title>Asset Income Reports</Title>
          <PrimaryButton
            buttonStyle={{ width: "18%", marginRight: "-6px" }}
            onClick={() => router.push("/assetincomereport/generatereport")}
          >
            Generate Report
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
      </Container>
    </>
  );
};

export default AssetIncomeReportPage;

const BackButtonLink = styled.button`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
  background: none;
  cursor: pointer;
  border: none;
  padding: 24px;
`;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;

  border-radius: 24px;

  table {
    width: 100%;
    table-layout: fixed;
    padding-top: 8px;
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
    width: 10%;
    text-align: center;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 8%;
    text-align: center;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 10%;
    text-align: center;
  }
  th:nth-child(6),
  td:nth-child(6) {
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
