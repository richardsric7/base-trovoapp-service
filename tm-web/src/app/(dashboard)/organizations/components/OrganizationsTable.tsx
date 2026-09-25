"use client";

import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import { useRouter } from "next/navigation";
import React, { useMemo, useState } from "react";
import styled from "styled-components";
import { RiDeleteBin6Line } from "react-icons/ri";
import { BiEdit } from "react-icons/bi";
import { FaEllipsisVertical } from "react-icons/fa6";
import DeleteOrgModal from "./DeleteOrgModal";
import EditOrgModal from "./EditOrgModal";
import { useGetOrganizationsListQuery } from "@/redux/api/organizations";
import { useEffect } from "react";
import { usePrettyType } from "./TypeFormatter";

interface OrgTableProps {
  searchTerm: string;
  filters?: Record<string, string>;
}

const OrganizationsTable: React.FC<OrgTableProps> = ({
  searchTerm,
  filters,
}) => {
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [activeTulipId, setActiveTulipId] = useState<number | null>(null);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isEditModal, setIsEditModal] = useState(false);
  const [selectedOrg, setSelectedOrg] = useState<any | null>(null);

  const { data, isLoading } = useGetOrganizationsListQuery({
    page: currentPage,
    pageSize,
    search: searchTerm,
    ...filters,
  });

  const prettyType = usePrettyType();

  // Close actions dropdown when modals open
  useEffect(() => {
    if (isEditModal || isDeleteModalOpen) {
      setActiveTulipId(null);
    }
  }, [isEditModal, isDeleteModalOpen]);

  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      render: (_: any, record: any) => {
        return (
          <AssetIdSection>
            <BoldText>
              {record.name ? record.name.charAt(0).toUpperCase() : "?"}
            </BoldText>

            <div>
              <AssetName>{record.name}</AssetName>
              <AssetSubName>{record.subname}</AssetSubName>
            </div>
          </AssetIdSection>
        );
      },
    },
    {
      title: "Type",
      dataIndex: "type",
      render: (_: any, record: any) => {
        return <TypeWrapper>{record.type}</TypeWrapper>;
      },
      key: "type",
    },
    {
      title: "Added On",
      dataIndex: "addedOn",
      key: "addedOn",
    },

    {
      title: "Members",
      dataIndex: "members",
      key: "members",
    },
    {
      title: "Admin",
      dataIndex: "admin",
      render: (_: any, record: any) => {
        return (
          <AssetIdSection>
            <BoldText2>
              {record.name ? record.admin.charAt(0).toUpperCase() : "?"}
            </BoldText2>

            <div>
              <AssetName>{record.admin}</AssetName>
              <TruncatedEmail title={record.adminEmail}>
                {record.adminEmail}
              </TruncatedEmail>
            </div>
          </AssetIdSection>
        );
      },
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
                    setSelectedOrg(record);
                    setIsEditModal(true);
                    setActiveTulipId(null);
                  }}
                >
                  <BiEdit />
                  Edit Organization
                </EditLink>
                <DeleteBtn
                  onClick={() => {
                    setSelectedOrg({
                      id: record.id,
                      name: record.name,
                    });

                    setIsDeleteModalOpen(true);
                    setActiveTulipId(null);
                  }}
                >
                  {" "}
                  <RiDeleteBin6Line />
                  Deactivate Organization
                </DeleteBtn>
              </ActionContent>
            )}
          </ActionsWrapper>
        );
      },
    },
  ];

  // Helper for date formatting
  const prettyDate = (dateStr: string) => {
    if (!dateStr) return "-";
    const date = new Date(dateStr);
    return date.toLocaleDateString("en-GB", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    }); // e.g. "23 Oct, 2023"
  };

  const organizations = data?.organizations ?? [];
  const dataSource = useMemo(() => {
    return organizations.map((org) => ({
      id: org.id,
      name: org.name,
      type: prettyType(org.type),
      status: org.status,
      rawType: org.type,
      email: org.email || "",
      address: org.address || "",
      country: org.country || "",
      stakeholder_type: org.stakeholder_type || "",
      fee_fixed: org.fee_fixed || "",
      fee_percent: org.fee_percent || "",
      addedOn: prettyDate(org.created_at),
      members: org.team_member_count ?? 0,
      admin: org.created_by_admin
        ? `${org.created_by_admin.first_name} ${org.created_by_admin.last_name}`
        : "N/A",
      adminEmail: org.created_by_admin?.email || "",
    }));
  }, [organizations]);

  const totalCount = data?.pagination?.total ?? 0;

  if (isLoading) {
    return <AletText>Loading Organizations...</AletText>;
  }

  if (!isLoading && dataSource.length === 0) {
    if (filters?.type) {
      return <AletText>No organizations found for the selected type.</AletText>;
    }
    return <AletText>No organizations found.</AletText>;
  }

  return (
    <Container>
      <CustomTable
        columns={columns}
        dataSource={dataSource}
        onRowClick={(record) => {
          router.push(`/organizations/${record.id}`);
        }}
      />
      <Pagination
        currentPage={currentPage}
        totalCount={totalCount}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />

      {isEditModal && selectedOrg && (
        <EditOrgModal
          isEditModal={isEditModal}
          setIsEditModal={setIsEditModal}
          selectedOrg={selectedOrg}
        />
      )}

      {isDeleteModalOpen && selectedOrg && (
        <DeleteOrgModal
          setIsDeleteModalOpen={setIsDeleteModalOpen}
          isDeleteModalOpen={isDeleteModalOpen}
          orgId={selectedOrg.id}
          orgName={selectedOrg.name}
        />
      )}
    </Container>
  );
};

export default OrganizationsTable;
const Container = styled.section`
  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 16%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 12%;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 8%;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 4%;
  }
  th:nth-child(5),
  td:nth-child(5) {
    width: 12%;
    text-align: center;
  }

  th:nth-child(6),
  td:nth-child(6) {
    width: 8%;
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

const AssetIdSection = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

const BoldText = styled.div`
  background-color: #007cdf;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #ffffff;
  font-weight: 600;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;
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

const TruncatedEmail = styled(AssetSubName)`
  max-width: 140px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
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
  width: 220px;
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

const TypeWrapper = styled.div`
  font-weight: 400;
  font-size: 12px;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: center;
  color: #007cdf;
  background-color: #007cdf1a;
  // width: 150px;
  gap: 10px;
  display: inline-flex;
  border-radius: 8px;
  padding: 8px;
`;

const AletText = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin-top: 14px;
  text-transform: capitalize;
`;
