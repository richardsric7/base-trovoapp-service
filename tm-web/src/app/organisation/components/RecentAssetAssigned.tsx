"use client";
import CustomTable from "@/components/CustomTable";
import React, { useMemo } from "react";
import styled from "styled-components";
import { AiOutlineCheckCircle, AiOutlineClockCircle } from "react-icons/ai";
import { AssetManagerDashboardAsset } from "@/redux/api/assetManager";

type VettingStatus = "Pending Vetting" | "Vetting Completed";
interface RecentAssetAssignedProps {
  assets?: AssetManagerDashboardAsset[];
}

const RecentAssetAssigned = ({ assets = [] }: RecentAssetAssignedProps) => {
  const columns = [
    {
      title: "User",
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
      render: (status: VettingStatus) => {
        const isPending = status === "Pending Vetting";

        return (
          <StatusBadge $status={status}>
            {isPending ? <AiOutlineClockCircle /> : <AiOutlineCheckCircle />}
            {isPending ? "Pending Vetting" : "Vetting Completed"}
          </StatusBadge>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 6 }, (_, index) => ({
      tokenId: index + 1,
      name: `ATL ${index + 1}`,
      subname: `Atalntis`,
      category: "Real Estate",
      status: index % 2 === 0 ? "Pending Vetting" : "Vetting Completed",
    }));
  }, []);
  // const dataSource = assets.map((asset) => ({
  //   tokenId: asset.asset_id,
  //   name: asset.asset_code || asset.asset_name || "N/A",
  //   subname: asset.asset_name || "N/A",
  //   category: asset.asset_sector || asset.asset_type || "N/A",
  //   status:
  //     asset.status?.toLowerCase() === "pending"
  //       ? "Pending Vetting"
  //       : "Vetting Completed",
  // }));
  return (
    <Container>
      <SubTitle>Recent Assets Assigned</SubTitle>
      <CustomTable columns={columns} dataSource={dataSource} />
    </Container>
  );
};

export default RecentAssetAssigned;

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
