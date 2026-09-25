"use client";
import React, { useMemo, useState } from "react";
import styled from "styled-components";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import { FaEllipsisVertical } from "react-icons/fa6";

import RoleManagement from "./RoleManagement";
import SuspendAdmin from "./SuspendAdmin";
import RoleSuspendCard from "./RoleSuspendCard";
import ChangeRole from "./ChangeRole";
import RevokeSuspension from "./RevokeSuspension";
import { useRouter } from "next/navigation";
import { useListAdminsQuery } from "@/redux/api/admin/admin";
import TruncatedText from "@/hooks/useTruncate";
import SuccessMessage from "@/components/SuccessMessage";
import { User } from "@/redux/api/admin";

interface AdminUserTableProps {
  searchTerm: string;
}

export default function AdminUserTable({ searchTerm }: AdminUserTableProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [showRoleManagement, setShowRoleManagement] = useState(false);
  const [showRoleSuspend, setShowRoleSuspend] = useState(false);

  const [revoke, setRevoke] = useState<boolean>(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [suspendAdmin, setSuspendAdmin] = useState<boolean>(false);
  const [changeAdminRole, setChangeAdminRole] = useState<boolean>(false);
  const [activeRowId, setActiveRowId] = useState<number | null>(null);
  const [showSuspendSuccess, setShowSuspendSuccess] = useState(false);
  const [suspendedEmail, setSuspendedEmail] = useState<string | null>(null);

  const router = useRouter();

  const { data, isLoading, isFetching, error } = useListAdminsQuery(
    {
      page: currentPage,
      pageSize,
    },
    {
      skip: false, // Only fetch when needed
    }
  );
  const adminData = data?.admins || [];
  console.log("adminspage", data);

  const columns = [
    {
      title: "Users",
      dataIndex: "Username",
      key: "Username",
      width: "20%",
      render: (_: any, record: User) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserName>
              {" "}
              <TruncatedText text={record.Username || "N/A"} maxLength={15} />
            </UserName>
          </div>
        </UserInfoSection>
      ),
    },
    {
      title: "Email Address",
      dataIndex: "Email",
      key: "Email",
      render: (email: string) => <UserEmail>{email || "N/A"}</UserEmail>,
    },
    {
      title: "Role",
      dataIndex: "Role",
      key: "Role",
      render: (role: string) => (
        <RoleDisplay>{role ? role.replace(/_/g, " ") : "N/A"}</RoleDisplay>
      ),
    },
    {
      title: "Created On",
      dataIndex: "CreatedAt",
      key: "CreatedAt",
      render: (CreatedAt: number) =>
        new Date(CreatedAt).toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        }),
    },
    {
      title: "Account Status",
      dataIndex: "Status",
      key: "Status",
      render: (_: any, record: User) => {
        const normalizedStatus =
          record.Status === "SUSPENDED" ? "SUSPENDED" : "ACTIVE";

        return (
          <UserStatus status={normalizedStatus}>{normalizedStatus}</UserStatus>
        );
      },
    },
    {
      title: "Last Activity",
      dataIndex: "UpdatedAt",
      key: "UpdatedAt",
      render: (UpdatedAt: number) =>
        new Date(UpdatedAt).toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        }),
    },
    {
      title: "Actions",
      dataIndex: "actions",
      key: "actions",
      render: (_: any, record: User) => (
        <ActionWrapper
          onClick={(e) => {
            e.stopPropagation();
            handleUserAction(record.Status, record);
          }}
        >
          <FaEllipsisVertical style={{ cursor: "pointer" }} />

          {activeRowId === record.ID && showRoleManagement && (
            <InlineMenu>
              <RoleManagement
                handleSuspend={handleSuspend}
                handleChangeRole={handleChangeRole}
              />
            </InlineMenu>
          )}

          {activeRowId === record.ID && showRoleSuspend && (
            <InlineMenu>
              <RoleSuspendCard
                handleChangeRole={handleChangeRole}
                handleRevoke={handleRevoke}
              />
            </InlineMenu>
          )}
        </ActionWrapper>
      ),
    },
  ];

  const handleSuspend = () => {
    setSuspendAdmin(!suspendAdmin);
  };

  const handleChangeRole = () => {
    setChangeAdminRole(!changeAdminRole);
  };

  const handleRevoke = () => {
    setRevoke(!revoke);
  };

  const handleUserAction = (status: string, user: User) => {
    setSelectedUser(user);
    setActiveRowId(user.ID);

    const normalizedStatus = status === "SUSPENDED" ? "SUSPENDED" : "ACTIVE";

    if (normalizedStatus === "ACTIVE") {
      setShowRoleManagement(activeRowId !== user.ID);
      setShowRoleSuspend(false);
    } else {
      setShowRoleSuspend(activeRowId !== user.ID);
      setShowRoleManagement(false);
    }
  };
  const filteredAdmin = useMemo(() => {
    if (!data?.admins) return [];
    return data.admins.filter(
      (admin) =>
        admin.FirstName.toLowerCase().includes(searchTerm.toLowerCase()) ||
        admin.Email.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [data?.admins, searchTerm]);

  const paginatedData = useMemo(() => {
    return filteredAdmin;
  }, [filteredAdmin]);

  const handlePageChange = (pageNumber: number) => {
    setCurrentPage(pageNumber); // This will trigger data fetching for the specific page
  };
  if (isLoading) {
    return (
      <Container>
        <UserName>Loading Admin users...</UserName>
      </Container>
    );
  }

  if (error) {
    return (
      <Container>
        <UserName>Error loading admin data</UserName>
      </Container>
    );
  }

  return (
    <>
      <Container>
        {suspendAdmin && selectedUser && (
          <SuspendAdmin
            suspendAdmin={suspendAdmin}
            setSuspendAdmin={setSuspendAdmin}
            adminEmail={selectedUser.Email}
            onSuccess={(email: string) => {
              // close the modal
              setSuspendAdmin(false);
              // show success message
              setSuspendedEmail(email);
              setShowSuspendSuccess(true);
            }}
          />
        )}
        {changeAdminRole && selectedUser && (
          <ChangeRole
            changeAdminRole={changeAdminRole}
            setChangeAdminRole={setChangeAdminRole}
            adminEmail={selectedUser.Email}
          />
        )}

        {revoke && selectedUser && (
          <RevokeSuspension
            revoke={revoke}
            setRevoke={setRevoke}
            adminName={`${selectedUser.FirstName} ${selectedUser.LastName}`}
            adminEmail={selectedUser.Email}
          />
        )}
        {searchTerm && filteredAdmin.length === 0 ? (
          <>
            <UserName>No users match your search. </UserName>
          </>
        ) : (
          <>
            <CustomTable
              columns={columns}
              dataSource={paginatedData}
              onRowClick={(record) => {
                if (record?.Username) {
                  router.push(
                    `/admin-users/${record.Username}?page=${currentPage}&pageSize=${pageSize}`
                  );
                }
              }}
            />

            <Pagination
              currentPage={currentPage}
              totalCount={data?.pagination.total || 0}
              pageSize={pageSize}
              onPageSizeChange={setPageSize}
              onPageChange={handlePageChange}
              isFetching={isFetching}
            />
          </>
        )}
      </Container>

      {showSuspendSuccess && suspendedEmail && (
        <SuccessMessage
          isOpen={showSuspendSuccess}
          setIsOpen={setShowSuspendSuccess}
          heading="Success"
          message="You have successfully suspended the admin user"
          email={suspendedEmail}
        />
      )}
    </>
  );
}

const Container = styled.section`
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
    width: 30%;
  }

  th:nth-child(3),
  td:nth-child(3) {
    width: 20%;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 12%;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 20%;
    text-align: center;
  }

  th:nth-child(6),
  td:nth-child(6) {
    width: 20%;
    text-align: center;
  }

  th:nth-child(7),
  td:nth-child(7) {
    width: 4%;
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
`;

const UserEmail = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  margin: 0;
`;

const UserStatus = styled.p<{ status: string }>`
  color: ${({ status }) => (status === "ACTIVE" ? "#00A859" : "#FF4D4D")};
  font-weight: 500;
  background-color: ${({ status }) =>
    status === "ACTIVE" ? "#00A8591A" : "#BE38001A"};
  padding: 8px;
  width: 120px;
  border-radius: 8px;
  text-align: center;
  width: 100px;
  margin: 0 auto;
`;
const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const RoleDisplay = styled.span``;

const ActionWrapper = styled.div`
  position: relative;
  display: inline-block;
  cursor: pointer;
  padding: 12px;
  margin: -8px;
`;

const InlineMenu = styled.div`
  position: absolute;
  top: 100%;
  right: 0;
  width: 190px;
  margin-top: 8px;
  z-index: 1;
  background-color: white;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  padding: 8px 12px;
`;
