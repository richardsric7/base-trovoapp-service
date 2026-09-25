"use client";

import React, { useState } from "react";
import styled from "styled-components";
import CustomTable from "@/components/CustomTable";
import PrimaryButton from "@/components/PrimaryButton";
import { FaPlus } from "react-icons/fa6";
import AssignOwnerModal from "./AssignOwnerModal";
import DeleteAssignmentConfirm from "./DeleteAssignmentConfirm";
import AccessDenied from "./AccessDenied";
import { useListAssignmentsQuery } from "@/redux/api/vaultSigner/api";
import type { IAssignment } from "@/redux/api/vaultSigner/interface";

interface AssignmentsTableProps {
  managedSecretId: string;
}

const AssignmentsTable: React.FC<AssignmentsTableProps> = ({ managedSecretId }) => {
  const { data, isLoading, error } = useListAssignmentsQuery(managedSecretId);
  const [assignOpen, setAssignOpen] = useState(false);
  const [editingAssignment, setEditingAssignment] = useState<IAssignment | null>(null);
  const [deletingAssignment, setDeletingAssignment] = useState<IAssignment | null>(null);

  if ((error as any)?.status === 403) {
    return <AccessDenied />;
  }

  const rows = data?.data ?? [];

  const columns = [
    {
      title: "Signer Position",
      dataIndex: "position",
      render: (value: number | null) => (value === null ? "Unclaimed" : `Position ${value}`),
    },
    { title: "Owner", dataIndex: "ownerLabel" },
    { title: "Owner Type", dataIndex: "ownerType" },
    {
      title: "Assigned At",
      dataIndex: "assignedAt",
      render: (value: string) => (value ? new Date(value).toLocaleString() : "—"),
    },
    {
      title: "Actions",
      dataIndex: "actions",
      render: (_: unknown, record: IAssignment) => (
        <ActionsWrapper>
          <ActionButton onClick={() => setEditingAssignment(record)}>Reassign</ActionButton>
          <ActionButton danger onClick={() => setDeletingAssignment(record)}>
            Remove
          </ActionButton>
        </ActionsWrapper>
      ),
    },
  ];

  return (
    <div>
      <Header>
        <div />
        <PrimaryButton
          buttonStyle={{
            width: "auto",
            padding: "10px 20px",
            display: "flex",
            alignItems: "center",
            gap: "4px",
          }}
          onClick={() => setAssignOpen(true)}
        >
          <FaPlus /> Assign Owner
        </PrimaryButton>
      </Header>

      <CustomTable columns={columns} dataSource={rows} isLoading={isLoading} />

      {assignOpen && (
        <AssignOwnerModal
          managedSecretId={managedSecretId}
          isOpen={assignOpen}
          onClose={() => setAssignOpen(false)}
        />
      )}
      {editingAssignment && (
        <AssignOwnerModal
          managedSecretId={managedSecretId}
          isOpen={!!editingAssignment}
          onClose={() => setEditingAssignment(null)}
          editingAssignment={editingAssignment}
        />
      )}
      {deletingAssignment && (
        <DeleteAssignmentConfirm
          managedSecretId={managedSecretId}
          assignment={deletingAssignment}
          isOpen={!!deletingAssignment}
          onClose={() => setDeletingAssignment(null)}
        />
      )}
    </div>
  );
};

export default AssignmentsTable;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: flex-end;
  margin-bottom: 16px;
`;

const ActionsWrapper = styled.div`
  display: flex;
  gap: 10px;
`;

const ActionButton = styled.button<{ danger?: boolean }>`
  background: transparent;
  border: 1px solid ${(props) => (props.danger ? "#d92d20" : "#007cdf")};
  color: ${(props) => (props.danger ? "#d92d20" : "#007cdf")};
  border-radius: 8px;
  padding: 6px 14px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
`;
