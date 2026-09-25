"use client";
import React, { useMemo, useState } from "react";
import styled from "styled-components";
import InviteMemberModal from "./InviteMemberModal";
import { SearchBar } from "@/components";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import { FaEllipsisVertical } from "react-icons/fa6";
import ChangePassword from "./ChangePassword";
import DisableMember from "./DisableMember";
import ChangeMemberRole from "./ChangeMemberRole";
import settingsIcon from "@/assets/images/settings-2.svg";
import shieldIcon from "@/assets/images/shield-slash.svg";

import lockIcon from "@/assets/images/lock.svg";

import Image from "next/image";
import MembersFilter from "./MembersFilter";
import { useGetOrganizationMembersQuery } from "@/redux/api/organizations";

interface MembersProps {
  organizationId: string;
}

const MembersTable: React.FC<MembersProps> = ({ organizationId }) => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [activeTulipId, setActiveTulipId] = useState<number | null>(null);
  const [isRestPassword, setIsRestPassword] = useState(false);
  const [isDisableUser, setIsDisableUser] = useState(false);
  const [isChangeRole, setIsChangeRole] = useState(false);

  const { data, isLoading } = useGetOrganizationMembersQuery({
    organization_id: organizationId,
    page: 1,
    pageSize: 10,
  });

  // const { data, isLoading } = useGetOrganizationMembersQuery(
  //   { organization_id: organizationId, page: 1, pageSize: 10 },
  //   { skip: !organizationId },
  // );

  const members = data?.data.members;
  const pagination = data?.data.pagination;

  const formatRole = (role: string) => {
    return role
      .toLowerCase()
      .split("_")
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(" ");
  };

  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      render: (_: any, record: any) => {
        return (
          <AssetIdSection>
            <BoldText2>
              {record.admin ? record.admin.charAt(0).toUpperCase() : "?"}
            </BoldText2>

            <div>
              <AssetName>{record.admin}</AssetName>
              <AssetSubName>{record.adminEmail}</AssetSubName>
            </div>
          </AssetIdSection>
        );
      },
    },
    {
      title: "Role",
      dataIndex: "role",
      render: (_: any, record: any) => {
        return <TypeWrapper>{formatRole(record.role)}</TypeWrapper>;
      },
      key: "role",
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

  const dataSource = useMemo(() => {
    if (!members) return [];

    return members.map((m: any) => ({
      id: m.id,
      role: m.role,

      addedOn: new Date(m.created_at).toLocaleDateString("en-GB", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      }),
      status: m.status,
      admin: `${m.first_name} ${m.last_name}`,
      adminEmail: m.email,
    }));
  }, [members]);

  return (
    <>
      <PageContainer>
        <Header>
          <Title>Members</Title>

          <InviteMemberModal organizationId={organizationId} />
        </Header>

        <FiltersSection>
          <SearchBar />
          <MembersFilter />
        </FiltersSection>

        <MemberTableContainer>
          {" "}
          <CustomTable columns={columns} dataSource={dataSource} />
          <Pagination
            currentPage={pagination?.page || 1}
            totalCount={pagination?.total || 0}
            pageSize={pagination?.pageSize || 10}
            onPageChange={setCurrentPage}
            onPageSizeChange={setPageSize}
          />
        </MemberTableContainer>
      </PageContainer>

      {isRestPassword && (
        <ChangePassword
          isResetPassword={isRestPassword}
          setIsRestPassword={setIsRestPassword}
          memberName="John Doe"
          memberEmail="john@gmail.com"
          organizationId="2234"
        />
      )}
      {isDisableUser && (
        <DisableMember
          setIsDisableUser={setIsDisableUser}
          isDisableUser={isDisableUser}
          memberName="John Doe"
        />
      )}

      {isChangeRole && (
        <ChangeMemberRole
          setIsChangeRole={setIsChangeRole}
          isChangeRole={isChangeRole}
          role="admin"
        />
      )}
    </>
  );
};

export default MembersTable;

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

const FiltersSection = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 25px;
  margin: 10px 0 30px 0;
`;

const MemberTableContainer = styled.div`
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
    width: 15%;
    // text-align: center;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 10%;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 10%;
  }
  th:nth-child(5),
  td:nth-child(5) {
    width: 4%;
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

const AssetIdSection = styled.div`
  display: flex;
  column-gap: 4px;
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
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #00225a;
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
