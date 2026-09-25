"use client";
import React, { useState } from "react";
import styled from "styled-components";
import { FaEllipsisVertical } from "react-icons/fa6";
import settingsIcon from "@/assets/images/settings-2.svg";
import shieldIcon from "@/assets/images/shield-slash.svg";

import lockIcon from "@/assets/images/lock.svg";

import Image from "next/image";

import { useGetOrganizationMembersQuery } from "@/redux/api/org";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import { SearchBar } from "@/components";
import MembersFilter from "@/app/(dashboard)/organizations/components/MembersFilter";
import InviteMemberModal from "@/app/(dashboard)/organizations/components/InviteMemberModal";


type Member = {
  id: string;
  role: string;
  created_at?: string;
  status?: string;
  first_name?: string;
  last_name?: string;
  email?: string;
};

const TeamMembersPage: React.FC = () => {
  const [currentPage, setCurrentPage] = React.useState(1);
  const [pageSize, setPageSize] = React.useState(10);
  const [isRestPassword, setIsRestPassword] = useState(false);
  const [isDisableUser, setIsDisableUser] = useState(false);
  const [isChangeRole, setIsChangeRole] = useState(false);

  const [activeTulipId, setActiveTulipId] = useState<number | null>(null);
  const { data, isLoading } = useGetOrganizationMembersQuery({
    page: currentPage,
    pageSize,
  });

  const members: Member[] | undefined = data?.data?.members;
  const pagination = data?.data?.pagination;

  const formatRole = (role?: string) =>
    (role ?? "")
      .toLowerCase()
      .split("_")
      .map((w) => (w ? w[0].toUpperCase() + w.slice(1) : w))
      .join(" ");

  const columns = [
    {
      title: "Name",
      dataIndex: "admin",
      key: "admin",
      render: (_: any, record: any) => (
        <AssetIdSection>
          <BoldText2>
            {record?.admin ? record?.admin.charAt(0).toUpperCase() : "?"}
          </BoldText2>
          <div>
            <AssetName>{record.admin}</AssetName>
            <AssetSubName>{record.adminEmail}</AssetSubName>
          </div>
        </AssetIdSection>
      ),
    },
    {
      title: "Role",
      dataIndex: "role",
      key: "role",
      render: (_: any, record: any) => (
        <TypeWrapper>{formatRole(record.role)}</TypeWrapper>
      ),
    },
    {
      title: "Added On",
      dataIndex: "addedOn",
      key: "addedOn",
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (_: any, record: any) => (
        <TypeWrapper>{record.status}</TypeWrapper>
      ),
    },

    {
      title: "Actions",
      dataIndex: "actions",
      render: (_: any, record: any) => {
        const isOpen = activeTulipId === record.id;

        return (
          <ActionsWrapper
            onClick={(e) => {
              e.stopPropagation();
              setActiveTulipId((prev) =>
                prev === record.id ? null : record.id,
              );
            }}
          >
            <FaEllipsisVertical />
            {isOpen && (
              <ActionContent>
                <EditLink
                  onClick={() => {
                    setIsChangeRole(true);
                    setActiveTulipId(null);
                  }}
                >
                  <Image
                    src={settingsIcon}
                    alt="setting-icon"
                    width={16}
                    height={16}
                  />
                  Change Role
                </EditLink>

                <EditLink
                  onClick={() => {
                    setIsRestPassword(true);
                    setActiveTulipId(null);
                  }}
                >
                  <Image
                    src={lockIcon}
                    alt="lock-icon"
                    width={16}
                    height={16}
                  />
                  Reset Password
                </EditLink>
                <DeleteBtn
                  onClick={() => {
                    setIsDisableUser(true);
                    setActiveTulipId(null);
                  }}
                >
                  <Image
                    src={shieldIcon}
                    alt="shield-icon"
                    width={16}
                    height={16}
                  />
                  Disable User
                </DeleteBtn>
              </ActionContent>
            )}
          </ActionsWrapper>
        );
      },
    },
  ];

  const dataSource = React.useMemo(() => {
    if (!members) return [];
    return members.map((m) => ({
      id: m.id,
      role: m.role,
      addedOn: m.created_at
        ? new Date(m.created_at).toLocaleDateString("en-GB", {
            day: "2-digit",
            month: "short",
            year: "numeric",
          })
        : "-",
      status: m.status ?? "Pending",
      admin: `${m.first_name ?? ""} ${m.last_name ?? ""}`.trim() || "-",
      adminEmail: m.email ?? "-",
    }));
  }, [members]);

  return (
    <PageContainer>
      <Header>
        <Title>Members</Title>

        <InviteMemberModal />
      </Header>

      <FiltersSection>
        <SearchBar />
        <MembersFilter />
      </FiltersSection>

      <MemberTableContainer>
        {isLoading ? (
          <p>Loading members...</p>
        ) : dataSource.length === 0 ? (
          <p>No members found.</p>
        ) : (
          <>
            <CustomTable columns={columns} dataSource={dataSource} />
            <Pagination
              currentPage={pagination?.page ?? currentPage}
              totalCount={pagination?.total ?? 0}
              pageSize={pagination?.pageSize ?? pageSize}
              onPageChange={setCurrentPage}
              onPageSizeChange={setPageSize}
            />
          </>
        )}
      </MemberTableContainer>
    </PageContainer>
  );
};

export default TeamMembersPage;

const PageContainer = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
  margin-top: 30px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  color: #00225a;
`;

const MemberTableContainer = styled.div`
  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 35%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 20%;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 20%;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 15%;
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

const TypeWrapper = styled.div`
  font-weight: 400;
  font-size: 12px;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: center;
  color: #007cdf;
  background-color: #007cdf1a;
  gap: 10px;
  display: inline-flex;
  border-radius: 8px;
  padding: 8px;
`;

const AssetIdSection = styled.div`
  display: flex;
  column-gap: 8px;
  align-items: center;
`;

const BoldText2 = styled.div`
  background-color: #acd1ef;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #191919;
  font-weight: 600;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;
`;

const AssetName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const AssetSubName = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #5a6b85;
`;

const FiltersSection = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 25px;
  margin: 10px 0 30px 0;
`;

const ActionsWrapper = styled.div`
  position: relative;
  display: inline-block;
  cursor: pointer;
  padding: 12px;
  margin: -8px;
`;

const ActionContent = styled.div`
  position: absolute;
  right: 0;
  top: 100%;
  margin-top: 8px;
  width: 200px;
  background-color: #ffffff;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  padding: 8px 12px;
  z-index: 1000;
`;

const DeleteBtn = styled.div`
  font-weight: 400;
  font-size: 14px;
  line-height: 28px;
  letter-spacing: 0%;
  display: flex;
  align-items: center;
  gap: 2px;
  color: #be3800;
`;

const EditLink = styled.div`
  font-weight: 400;
  font-size: 14px;
  line-height: 28px;
  letter-spacing: 0%;
  display: flex;
  align-items: center;
  gap: 2px;
  color: #00225a;
  text-decoration: none;
`;
