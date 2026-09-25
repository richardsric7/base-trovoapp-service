"use client";
import CustomTable from "@/components/CustomTable";
import React, { useMemo } from "react";
import { useRouter } from "next/navigation";
import styled from "styled-components";
import { AiOutlineCheckCircle, AiOutlineClockCircle } from "react-icons/ai";
import { LaRecentAsset } from "@/redux/api/legalAdviser";

type StructuringStatus = "Pending Structuring" | "Structuring Completed";
interface RecentAssetAssignedProps {
  assets?: LaRecentAsset[];
}

const LaRecentAssetAssigned = ({ assets = [] }: RecentAssetAssignedProps) => {
  const router = useRouter();

  const columns = [
    {
      title: "Asset",
      dataIndex: "name",
      key: "name",
      width: "30%",
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
      width: "25%",
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      width: "25%",
      render: (status: StructuringStatus) => {
        const isPending = status === "Pending Structuring";

        return (
          <StatusBadge $status={status}>
            {isPending ? <AiOutlineClockCircle /> : <AiOutlineCheckCircle />}
            {status}
          </StatusBadge>
        );
      },
    },
    {
      title: "Action",
      dataIndex: "action",
      key: "action",
      width: "20%",
      render: (_: any, record: any) =>
        record.status === "Pending Structuring" ? (
          <ActionButton
            onClick={() =>
              router.push(
                `/organisation/legaladviser/structuring/${record.key}/complete`,
              )
            }
          >
            Complete Structuring
          </ActionButton>
        ) : (
          <NoAction>—</NoAction>
        ),
    },
  ];

  const dataSource = useMemo(() => {
    if (!assets || assets.length === 0) return [];
    return assets.map((asset) => ({
      key: asset.id,
      name: asset.assetCode || asset.assetName || "N/A",
      subname: asset.assetName || "N/A",
      category: asset.assetSector || asset.assetSubSector || "N/A",
      status:
        asset.assignment?.status === "active"
          ? "Structuring Completed"
          : "Pending Structuring",
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

export default LaRecentAssetAssigned;

const Container = styled.div`
  background-color: #fff;
  padding: 32px;
  height: 100%;
  gap: 20px;
  border-radius: 24px;

  table {
    width: 100%;
    table-layout: fixed;
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

const StatusBadge = styled.span<{ $status: StructuringStatus }>`
  font-size: 12px;
  font-weight: 500;
  padding: 6px 10px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  gap: 6px;

  background-color: ${({ $status }) =>
    $status === "Pending Structuring" ? "#F2F6F9" : "#00A85926"};

  color: ${({ $status }) =>
    $status === "Pending Structuring" ? "#007CDF" : "#00A859"};
`;

const ActionButton = styled.button`
  border: 1px solid #007cdf;
  color: #007cdf;
  background: #fff;
  border-radius: 7px;
  padding: 7px 12px;
  cursor: pointer;
  font: inherit;
  font-size: 12px;
  font-weight: 600;
`;

const NoAction = styled.span`
  color: #98a2b3;
`;
