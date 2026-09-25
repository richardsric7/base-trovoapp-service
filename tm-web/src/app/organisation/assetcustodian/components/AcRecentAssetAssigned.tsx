"use client";
import CustomTable from "@/components/CustomTable";
import React, { useMemo } from "react";
import styled from "styled-components";
import { AiOutlineCheckCircle, AiOutlineClockCircle } from "react-icons/ai";
import { CustodianDashboardAsset } from "@/redux/api/assetCustodian";

interface AcRecentAssetAssignedProps {
  assets?: CustodianDashboardAsset[];
  isLoading?: boolean;
}

// Mirrors the tokenization status codes used across the asset tokenization
// views (see AssetsTable's statusMapping) since the dashboard payload only
// exposes the raw assetTokenizationStatus code, not a label.
const TOKENIZATION_STATUS_LABELS: Record<number, string> = {
  0: "Pending Vetting",
  1: "Awaiting Payment",
  2: "Payment Made",
  3: "Processing",
  4: "Tokenization Approved",
  5: "Primary Sale Started",
  6: "Secondary Sales",
  7: "Liquidated",
  8: "Refunded",
};

const humanize = (value?: string) =>
  value
    ? value
        .toLowerCase()
        .split("_")
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(" ")
    : undefined;

const resolveAssetStatus = (asset: CustodianDashboardAsset) =>
  humanize(asset.assignment?.status) ??
  TOKENIZATION_STATUS_LABELS[asset.assetTokenizationStatus ?? -1];

const AcRecentAssetAssigned = ({
  assets = [],
  isLoading = false,
}: AcRecentAssetAssignedProps) => {
  const columns = [
    {
      title: "Asset",
      dataIndex: "name",
      key: "name",
      width: "100%",
      render: (_: any, record: any) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserName>{record.name || "N/A"}</UserName>
            <SubName>{record.subname || "N/A"}</SubName>
          </div>
        </UserInfoSection>
      ),
    },
    {
      title: "Category",
      dataIndex: "category",
      key: "category",
      width: "100%",
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      width: "100%",
      render: (status: string) => {
        const isPending = status?.toLowerCase().includes("pending");

        return (
          <StatusBadge $pending={isPending}>
            {isPending ? <AiOutlineClockCircle /> : <AiOutlineCheckCircle />}
            {status || "N/A"}
          </StatusBadge>
        );
      },
    },
  ];

  const dataSource = useMemo(
    () =>
      assets.map((asset, index) => ({
        key: asset.id ?? index,
        name: asset.assetCode || asset.assetName || "N/A",
        subname: asset.assetName || "N/A",
        category: asset.assetSector || asset.assetType || "N/A",
        status: resolveAssetStatus(asset) || "N/A",
      })),
    [assets],
  );

  return (
    <Container>
      <SubTitle>Recent Assets Assigned</SubTitle>
      <CustomTable
        columns={columns}
        dataSource={dataSource}
        isLoading={isLoading}
        totalItems={dataSource.length}
      />
    </Container>
  );
};

export default AcRecentAssetAssigned;

const Container = styled.div`
  background-color: #fff;
  padding: 32px;
  height: 100%;
  gap: 20px;
  border-radius: 24px;

  table {
    width: 100%;
    table-layout: fixed;

    th:nth-child(1),
    td:nth-child(1) {
      width: 14%;
    }

    th:nth-child(2),
    td:nth-child(2) {
      width: 18%;

      text-align: center;
    }

    th:nth-child(3),
    td:nth-child(3) {
      width: 20%;

      text-align: center;
    }
  }
`;

const SubTitle = styled.h3`
  font-size: 20px;
  font-weight: 600;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;

const UserInfoSection = styled.div`
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

const SubName = styled.p`
  font-size: 10px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const Avatar = styled.div`
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const StatusBadge = styled.span<{ $pending: boolean }>`
  font-size: 12px;
  font-weight: 500;
  padding: 6px 10px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  gap: 6px;

  background-color: ${({ $pending }) => ($pending ? "#F2F6F9" : "#00A85926")};

  color: ${({ $pending }) => ($pending ? "#007CDF" : "#00A859")};
`;
