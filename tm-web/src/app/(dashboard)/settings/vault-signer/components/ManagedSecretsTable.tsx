"use client";

import React, { useState } from "react";
import styled from "styled-components";
import { useRouter } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import PrimaryButton from "@/components/PrimaryButton";
import TruncatedText from "@/hooks/useTruncate";
import { FaPlus } from "react-icons/fa6";
import RegisterManagedSecretModal from "./RegisterManagedSecretModal";
import AccessDenied from "./AccessDenied";
import { useListManagedSecretsQuery } from "@/redux/api/vaultSigner/api";
import type { IManagedSecretWithAssignments } from "@/redux/api/vaultSigner/interface";

const ManagedSecretsTable = () => {
  const router = useRouter();
  const { data, isLoading, error } = useListManagedSecretsQuery();
  const [registerOpen, setRegisterOpen] = useState(false);

  if ((error as any)?.status === 403) {
    return <AccessDenied />;
  }

  const rows = data?.data ?? [];

  const columns = [
    {
      title: "Label",
      dataIndex: "label",
      render: (_: unknown, record: IManagedSecretWithAssignments) => record.managedSecret.label,
    },
    {
      title: "Wallet Public Key",
      dataIndex: "walletPublicKey",
      render: (_: unknown, record: IManagedSecretWithAssignments) => (
        <TruncatedText text={record.managedSecret.wallet_public_key} maxLength={16} />
      ),
    },
    {
      title: "Active Signing Count",
      dataIndex: "activeSigningCount",
      render: (_: unknown, record: IManagedSecretWithAssignments) =>
        record.managedSecret.active_signing_count,
    },
    {
      title: "Assignments",
      dataIndex: "assignmentsCount",
      render: (_: unknown, record: IManagedSecretWithAssignments) => record.assignments.length,
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
          onClick={() => setRegisterOpen(true)}
        >
          <FaPlus /> Register Secret
        </PrimaryButton>
      </Header>

      <CustomTable
        columns={columns}
        dataSource={rows}
        isLoading={isLoading}
        onRowClick={(record) => router.push(`/settings/vault-signer/${record.managedSecret.id}`)}
      />

      {registerOpen && (
        <RegisterManagedSecretModal isOpen={registerOpen} onClose={() => setRegisterOpen(false)} />
      )}
    </div>
  );
};

export default ManagedSecretsTable;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: flex-end;
  margin-bottom: 16px;
`;
