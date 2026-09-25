"use client";

import React from "react";
import styled from "styled-components";
import CustomTable from "@/components/CustomTable";
import TruncatedText from "@/hooks/useTruncate";
import AccessDenied from "./AccessDenied";
import { useListAuditLogQuery } from "@/redux/api/vaultSigner/api";
import type { IAuditLogRow } from "@/redux/api/vaultSigner/interface";

// Read-only — the backend caps this at the 200 most recent rows server-side, no pagination
// params on GET /admin/vault-signer/audit-log.
const AuditLogTable = () => {
  const { data, isLoading, error } = useListAuditLogQuery();

  if ((error as any)?.status === 403) {
    return <AccessDenied />;
  }

  const rows = data?.data ?? [];

  const columns = [
    { title: "Kind", dataIndex: "kind" },
    {
      title: "Actor",
      dataIndex: "actor",
      render: (_: unknown, record: IAuditLogRow) => `${record.actor_id} (${record.actor_type})`,
    },
    {
      title: "Position / Vault Key",
      dataIndex: "positionOrKey",
      render: (_: unknown, record: IAuditLogRow) =>
        record.position != null ? `Position ${record.position}` : record.vault_key ?? "—",
    },
    {
      title: "Stellar Tx",
      dataIndex: "stellarTx",
      render: (_: unknown, record: IAuditLogRow) =>
        record.stellar_tx_hash ? (
          <TxWrapper>
            <TruncatedText text={record.stellar_tx_hash} maxLength={12} />
            <TxStatus status={record.stellar_tx_status}>{record.stellar_tx_status}</TxStatus>
          </TxWrapper>
        ) : (
          "—"
        ),
    },
    {
      title: "Changed At",
      dataIndex: "changed_at",
      render: (value: string) => (value ? new Date(value).toLocaleString() : "—"),
    },
  ];

  return <CustomTable columns={columns} dataSource={rows} isLoading={isLoading} />;
};

export default AuditLogTable;

const TxWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const TxStatus = styled.span<{ status?: string }>`
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 999px;
  background: ${(props) => (props.status === "success" ? "#e9f8f1" : "#fdecea")};
  color: ${(props) => (props.status === "success" ? "#00875a" : "#d92d20")};
`;
