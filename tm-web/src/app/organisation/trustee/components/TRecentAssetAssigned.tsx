"use client";
import CustomTable from "@/components/CustomTable";
import React, { useMemo } from "react";
import styled from "styled-components";
import { AiOutlineCheckCircle, AiOutlineClockCircle } from "react-icons/ai";
import { TrusteeRecentAsset } from "@/redux/api/trustees";

type VettingStatus = "Pending Vetting" | "Vetting Completed";

interface RecentAssetAssignedProps {
  assets?: TrusteeRecentAsset[];
}

const TRecentAssetAssigned = ({ assets = [] }: RecentAssetAssignedProps) => {
  const columns = [
    {
      title: "Asset",
      dataIndex: "name",
      key: "name",
      width: "35%",
      render: (_: any, record: any) => (
        <UserInfoSection>
          <Avatar>{(record.name || "A").charAt(0).toUpperCase()}</Avatar>
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
      width: "35%",
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      width: "30%",
      render: (status: VettingStatus) => {
        const isPending = status === "Pending Vetting";

        return (
          <StatusBadge $status={status}>
            {isPending ? <AiOutlineClockCircle /> : <AiOutlineCheckCircle />}
            {isPending ? "Pending Vetting" : "Active"}
          </StatusBadge>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    if (!assets || assets.length === 0) return [];
    return assets.map((asset) => ({
      key: asset.id,
      name: asset.assetCode || asset.assetName || "N/A",
      subname: asset.assetName || "N/A",
      category: asset.assetSector || "N/A",
      status: asset.assignment?.status === "active" ? "Vetting Completed" : "Pending Vetting",
    }));
  }, [assets]);

  return (
    <Container>
      <SubTitle>Recent Assets Assigned</SubTitle>
      {dataSource.length === 0 ? (
        <EmptyState>No data available</EmptyState>
      ) : (
        <CustomTable columns={columns} dataSource={dataSource} />
      )}
    </Container>
  );
};

export default TRecentAssetAssigned;

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
      //  text-align: center;
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

const EmptyState = styled.div`
  padding: 40px 0;
  text-align: center;
  color: #8a94a6;
  font-size: 14px;
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
  background-color: #007cdf;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
`;

const StatusBadge = styled.span<{ $status: VettingStatus }>`
  font-size: 12px;
  font-weight: 500;
  padding: 6px 10px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  gap: 6px;

  background-color: ${({ $status }) =>
    $status === "Pending Vetting" ? "#F2F6F9" : "#00A85926"};

  color: ${({ $status }) =>
    $status === "Pending Vetting" ? "#007CDF" : "#00A859"};
`;
