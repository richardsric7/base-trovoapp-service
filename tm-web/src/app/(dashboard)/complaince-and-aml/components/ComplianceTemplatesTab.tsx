"use client";

import { useMemo, useState } from "react";
import styled from "styled-components";

import { Modal } from "@/components";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import PrimaryButton from "@/components/PrimaryButton";
import { useGetComplianceTemplatesQuery } from "@/redux/api/admin";
import { ComplianceTemplate } from "@/redux/api/compliance/interface";
import { ORG_TYPES } from "@/constants/organizationTypes";
import ComplianceTemplateForm from "./ComplianceTemplateForm";

const orgTypeLabel = (value: string) =>
  ORG_TYPES.find((org) => org.value === value)?.label || value;

const formatDate = (value?: string) =>
  value ? new Date(value).toLocaleDateString() : "N/A";

const ComplianceTemplatesTab = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [editingTemplate, setEditingTemplate] =
    useState<ComplianceTemplate | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const { data, isLoading, isFetching, isError, refetch } =
    useGetComplianceTemplatesQuery({ page: currentPage, limit: pageSize });

  const records = useMemo(
    () => (Array.isArray(data?.data?.records) ? data.data.records : []),
    [data],
  );
  const meta = data?.data?.meta;

  const openCreateModal = () => {
    setEditingTemplate(null);
    setIsModalOpen(true);
  };

  const openEditModal = (template: ComplianceTemplate) => {
    setEditingTemplate(template);
    setIsModalOpen(true);
  };

  const columns = [
    {
      title: "Org Type",
      dataIndex: "org_type",
      render: (value: string) => orgTypeLabel(value),
    },
    { title: "Level", dataIndex: "level" },
    { title: "Version", dataIndex: "version" },
    {
      title: "Items",
      dataIndex: "items",
      render: (items: ComplianceTemplate["items"]) => items?.length ?? 0,
    },
    {
      title: "Updated At",
      dataIndex: "updated_at",
      render: (value: string) => formatDate(value),
    },
    {
      title: "",
      dataIndex: "id",
      render: (_: string, record: ComplianceTemplate) => (
        <EditLink onClick={() => openEditModal(record)}>Edit</EditLink>
      ),
    },
  ];

  if (isError) {
    return (
      <Container>
        <Title>Compliance Templates</Title>
        <StateMessage>
          Unable to load compliance templates.
          <RetryButton onClick={() => refetch()}>Try again</RetryButton>
        </StateMessage>
      </Container>
    );
  }

  return (
    <Container>
      <Header>
        <div>
          <Title>Compliance Templates</Title>
          <SubTitle>
            Define reusable KYC/KYB/compliance requirement sets by
            organisation type and level.
          </SubTitle>
        </div>
        <PrimaryButton
          onClick={openCreateModal}
          buttonStyle={{ margin: 0, width: "200px", height: "44px" }}
        >
          New Template
        </PrimaryButton>
      </Header>

      <CustomTable
        columns={columns}
        dataSource={records}
        isLoading={isLoading || isFetching}
        totalItems={meta?.total ?? 0}
        pageSize={pageSize}
        onRowClick={openEditModal}
      />

      {(meta?.total ?? 0) > 0 && (
        <Pagination
          totalCount={meta?.total ?? 0}
          currentPage={currentPage}
          pageSize={pageSize}
          onPageChange={setCurrentPage}
          onPageSizeChange={setPageSize}
        />
      )}

      <Modal
        title={editingTemplate ? "Edit Compliance Template" : "New Compliance Template"}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
      >
        <ComplianceTemplateForm
          editingTemplate={editingTemplate}
          onSaved={() => setIsModalOpen(false)}
        />
      </Modal>
    </Container>
  );
};

export default ComplianceTemplatesTab;

const Container = styled.div`
  background: #ffffff;
  border-radius: 24px;
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 24px;
`;

const Header = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin: 0 0 4px 0;
`;

const SubTitle = styled.p`
  font-size: 14px;
  color: #828282;
  margin: 0;
`;

const StateMessage = styled.div`
  color: #00225a;
`;

const RetryButton = styled.button`
  display: block;
  margin-top: 12px;
  padding: 8px 16px;
  border: 0;
  border-radius: 8px;
  background: #007cdf;
  color: #fff;
  cursor: pointer;
`;

const EditLink = styled.button`
  border: none;
  background: transparent;
  color: #007cdf;
  font-weight: 600;
  font-size: 13px;
  cursor: pointer;
  padding: 0;
`;
