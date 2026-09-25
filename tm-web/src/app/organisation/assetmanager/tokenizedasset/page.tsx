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

import { useGetAssetManagerAssetsQuery } from "@/redux/api/assetManager";
import Loader from "@/components/Loader";

import EmptyTableData from "@/components/EmptyTableData";
import PrimaryButton from "@/components/PrimaryButton";

const TokenizedAssetPage = () => {
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchTerm, setSearchTerm] = useState("");

  const { data, isLoading, isFetching } = useGetAssetManagerAssetsQuery({
    page: currentPage,
    limit: pageSize,
    search: searchTerm || undefined,
  });

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

  const assetManagerAssets = useMemo(() => {
    const assets = data?.data?.records ?? [];

    return assets.map((asset) => {
      const compliance = asset.complianceStatus ?? "";
      const custody = asset.custodyStatus ?? "";
      const valuation = asset.assetCurrentValue;
      const currency = asset.assetQuoteCurrency ?? "NGN";
      const updatedAt = asset.updatedAt ?? asset.createdAt;

      return {
        id: asset.id,
        name: asset.assetName ?? asset.assetCode ?? "N/A",
        subname: asset.assetCode ?? "N/A",
        sector: asset.assetSector ?? "N/A",
        valuation:
          valuation !== undefined && valuation !== null
            ? `${formatAmount(Number(valuation))} ${currency}`
            : `0.00 ${currency}`,
        holders: String(asset.tokenHolderCount ?? 0),
        compliance: compliance
          ? String(compliance).replaceAll("_", " ")
          : "N/A",
        complianceType: String(compliance).toLowerCase(),
        date: updatedAt ? formatDate(updatedAt) : "N/A",
        custody: custody ? String(custody).replaceAll("_", " ") : "N/A",
        custodyType: String(custody).toLowerCase(),
      };
    });
  }, [data]);
  const totalAssetCount =
    data?.data?.summary?.total ??
    data?.data?.meta?.total ??
    assetManagerAssets.length;
  const activeAssetCount = data?.data?.summary?.active ?? 0;
  const liquidatedAssetCount = data?.data?.summary?.liquidated ?? 0;

  // const assetManagerAssets = data?.data;

  const dummyAssets = useMemo(() => {
    return [
      {
        id: 1,
        name: "Atlantis 1",
        subname: "Atlantis Developers",
        sector: "Real Estate",
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
        sector: "Agriculture",
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
        sector: "Real Estate",
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
      title: "Asset Name",
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
      title: "Sector",
      dataIndex: "sector",
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
  if (isLoading || isFetching) {
    return <Loader />;
  }

  return (
    <>
      <CardContainer>
        <HeadingContainer>
          <Title>Tokenized Asset</Title>

          <ButtonContainer>
            {/* <PrimaryButton>Export</PrimaryButton> */}
            <PrimaryButton
              onClick={() =>
                router.push("/organisation/assetmanager/reports?generate=true")
              }
            >
              Generate Report
            </PrimaryButton>
          </ButtonContainer>
        </HeadingContainer>
        <CardContent>
          <OrgTokenizedAssetsStats
            text="Total Asset Count"
            count={String(totalAssetCount)}
            dotColor="#007CDF"
          />

          <OrgTokenizedAssetsStats
            text="Active Assets"
            count={String(activeAssetCount)}
            dotColor="#00A859"
          />

          <OrgTokenizedAssetsStats
            text="Liquidated Assets"
            count={String(liquidatedAssetCount)}
            dotColor="#BE3800"
          />
        </CardContent>
      </CardContainer>

      {assetManagerAssets.length === 0 ? (
        <EmptyTableData />
      ) : (
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
              dataSource={assetManagerAssets}
              // dataSource={dummyAssets}
              isLoading={false}
              onRowClick={(record) => {
                router.push(
                  `/organisation/assetmanager/tokenizedasset/${record.id}`,
                );
              }}
            />
            <Pagination
              currentPage={currentPage}
              totalCount={totalAssetCount}
              pageSize={pageSize}
              onPageChange={setCurrentPage}
              onPageSizeChange={setPageSize}
            />
          </TableContainer>
        </Container>
      )}
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
