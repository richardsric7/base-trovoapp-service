"use client";
import React, { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import styled from "styled-components";
import { SearchBar } from "@/components";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import {
  AiOutlineCheckCircle,
  AiOutlineClockCircle,
  AiOutlineCloseCircle,
} from "react-icons/ai";
import OrgTokenizedAssetsStats from "./_components/OrgTokenizedAssetsStats";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";

// Dummy statistics data
const dummyStats = {
  total: {
    count: 156,
    total_tokenized_value: 45000000.0,
  },
  approved: {
    count: 89,
    total_tokenized_value: 28500000.0,
  },
  pending: {
    count: 52,
    total_tokenized_value: 14200000.0,
  },
  rejected: {
    count: 15,
    total_tokenized_value: 2300000.0,
  },
};

// Dummy table data

const TokenizedAssetPage = () => {
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchTerm, setSearchTerm] = useState("");

  const formatAmount = (value?: number) =>
    typeof value === "number"
      ? value.toLocaleString(undefined, {
          minimumFractionDigits: 2,
          maximumFractionDigits: 2,
        })
      : "0.00";

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const dummyAssets = useMemo(() => {
    return [
      {
        id: 1,
        name: "Atlantis 1",
        subname: "Atlantis Developers",
        category: "Real Estate",
        valuation: "10,000,000 NGN",
        holders: "130",
        compliance: "Verified",
        complianceType: "verified",
        date: "22 Oct, 2023",
        custody: "Active",
        custodyType: "active",
      },
      {
        id: 2,
        name: "Animal Farm",
        subname: "Greg Donald",
        category: "Agriculture",
        valuation: "10,000,000 NGN",
        holders: "130",
        compliance: "Pending",
        complianceType: "pending",
        date: "20 Oct, 2023",
        custody: "Liquidated",
        custodyType: "liquidated",
      },
      {
        id: 3,
        name: "Atlantis 1",
        subname: "Atlantis Developers",
        category: "Real Estate",
        valuation: "10,000,000 NGN",
        holders: "130",
        compliance: "Rejected",
        complianceType: "rejected",
        date: "22 Oct, 2023",
        custody: "Inactive",
        custodyType: "inactive",
      },
    ];
  }, []);

  const columns = [
    {
      title: "Asset ID",
      dataIndex: "asset",
      render: (_: any, record: any) => (
        <AssetCell>
          <Avatar />
          <div>
            <AssetName>{record.name}</AssetName>
            <SubText>{record.subname}</SubText>
          </div>
        </AssetCell>
      ),
    },
    {
      title: "Category",
      dataIndex: "category",
    },
    {
      title: "Current Valuation",
      dataIndex: "valuation",
    },
    {
      title: "Token Holders",
      dataIndex: "holders",
    },
    {
      title: "Compliance Status",
      dataIndex: "compliance",
      render: (_: any, record: any) => (
        <ComplianceStatus type={record.complianceType}>
          {record.complianceType === "verified" && <AiOutlineCheckCircle />}
          {record.complianceType === "pending" && <AiOutlineClockCircle />}
          {record.complianceType === "rejected" && <AiOutlineCloseCircle />}
          <span>{record.compliance}</span>
        </ComplianceStatus>
      ),
    },
    {
      title: "Last Updated On",
      dataIndex: "date",
    },
    {
      title: "Custody Status",
      dataIndex: "custody",
      render: (_: any, record: any) => (
        <CustodyStatus type={record.custodyType}>
          {record.custody}
        </CustodyStatus>
      ),
    },
  ];

  return (
    <>
      <CardContainer>
        <HeadingContainer>
          <Title>Tokenized Asset</Title>

          <ButtonContainer>
            <PrimaryButton>Export</PrimaryButton>
            <SecondaryButton>Generate Report</SecondaryButton>
          </ButtonContainer>
        </HeadingContainer>
        <CardContent>
          <OrgTokenizedAssetsStats
            text="Total Asset Count"
            count="12"
            dotColor="#007CDF"
          />

          <OrgTokenizedAssetsStats
            text="Active Assets"
            count="10"
            dotColor="#00A859"
          />

          <OrgTokenizedAssetsStats
            text="Liquidated Assets"
            count="8"
            dotColor="#BE3800"
          />
        </CardContent>
      </CardContainer>
      <Container>
        <Header>
          <TitleSection>
            <Title>Tokenized Assets</Title>
          </TitleSection>
        </Header>

        <FiltersSection>
          <SearchBar
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </FiltersSection>

        <TableContainer>
          <CustomTable
            columns={columns}
            dataSource={dummyAssets}
            isLoading={false}
            onRowClick={(record) => {
              router.push(`/organisation/tokenizedasset/${record.id}`);
            }}
          />
          <Pagination
            currentPage={currentPage}
            totalCount={dummyAssets.length}
            pageSize={pageSize}
            onPageChange={setCurrentPage}
            onPageSizeChange={setPageSize}
          />
        </TableContainer>
      </Container>
    </>
  );
};

export default TokenizedAssetPage;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
`;

const CardContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 20px;
  background-color: #fff;
  border-radius: 24px;
  margin-bottom: 20px;
`;

const HeadingContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
`;

const ButtonContainer = styled.div`
  display: flex;
  align-item: center;
  justify-content: center;
  gap: 10px;
  width: 320px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const CardContent = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
  gap: 20px;
`;

const TitleSection = styled.div``;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  line-height: 1;
  margin: 0;
  color: #00225a;
`;

const FiltersSection = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 25px;
  margin: 10px 0 30px 0;
`;

const TableContainer = styled.section`
  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 20%;
  }

  th:nth-child(2),
  td:nth-child(2) {
    width: 12%;
  }

  th:nth-child(3),
  td:nth-child(3) {
    width: 14%;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 12%;
    text-align: center;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 15%;
    text-align: center;
  }

  th:nth-child(6),
  td:nth-child(6) {
    width: 20%;
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

const AssetCell = styled.div`
  display: flex;
  gap: 10px;
  align-items: center;
`;

const SubText = styled.p`
  font-size: 12px;
  color: #9ca3af;
  margin: 0;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;

const AssetName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const ComplianceStatus = styled.div<{ type: string }>`
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  border-radius: 8px;
  padding: 6px 10px;

  background: ${(p) =>
    p.type === "verified"
      ? "#dcfce7"
      : p.type === "pending"
        ? "#eaf3ff"
        : "#fee2e2"};

  color: ${(p) =>
    p.type === "verified"
      ? "#166534"
      : p.type === "pending"
        ? "#007cdf"
        : "#b91c1c"};
`;

const CustodyStatus = styled.div<{ type: string }>`
  display: inline-flex;
  justify-content: center;
  font-size: 12px;
  border-radius: 8px;
  padding: 6px 12px;
  width: fit-content;

  background: ${(p) =>
    p.type === "active"
      ? "#dcfce7"
      : p.type === "inactive"
        ? "#e5e7eb"
        : "#fee2e2"};

  color: ${(p) =>
    p.type === "active"
      ? "#166534"
      : p.type === "inactive"
        ? "#6b7280"
        : "#b91c1c"};
`;
