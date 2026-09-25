"use client";

import React, { useState } from "react";
import styled from "styled-components";
import CustomTable from "@/components/CustomTable";
import PrimaryButton from "@/components/PrimaryButton";
import EmptyState from "@/components/EmptyState";
import { showErrorToast, showSuccessToast } from "@/components";
import SetSignerValueModal from "./SetSignerValueModal";
import type {
  IMySecretListItem,
  IMySecretListResponse,
  IPutSecretRequest,
  IPutSecretResponse,
} from "@/redux/api/vaultSigner/interface";

// Reads any error the axios interceptors relay through RTK Query's `error` shape.
// The vault-signer backend replies with { error, message } (gin.H) for its structured cases —
// e.g. 409 { "error": "owner_already_assigned", "message": "This owner already has an
// assignment on this managed secret." } — where `error` is a machine-readable slug and
// `message` is the sentence meant for a person. `message` is checked first so the reader sees
// that sentence, not the slug; `error` is still the fallback for the plain `{"error": err.Error()}`
// shape this codebase's non-vault-signer handlers use, where `error` *is* the message (the same
// shape LinkWalletModal.tsx already reads from).
export const getErrorMessage = (error: any, fallback: string) =>
  error?.data?.message ?? error?.data?.error ?? error?.error ?? fallback;

// Loosely typed against the *shape* this component actually calls, not RTK Query's exact
// generic instantiation — the dashboard's baseApi hooks and the organisation portal's orgApi
// hooks are bound to two different BaseQueryFn types (see redux/api/vaultSigner/api.ts vs
// orgApi.ts) and are not directly assignable to one another under `typeof`, even though both
// satisfy this calling shape at runtime.
type UseListMySecrets = (
  arg?: void,
  options?: any,
) => { data?: IMySecretListResponse; isLoading: boolean; error?: unknown; refetch: () => void };

type UsePutMySecretValue = () => readonly [
  (arg: { secretId: string; payload: IPutSecretRequest }) => {
    unwrap: () => Promise<IPutSecretResponse>;
  },
  { isLoading: boolean },
];

interface MySignerSecretsProps {
  useListMySecrets: UseListMySecrets;
  usePutMySecretValue: UsePutMySecretValue;
}

const MySignerSecrets: React.FC<MySignerSecretsProps> = ({
  useListMySecrets,
  usePutMySecretValue,
}) => {
  const { data, isLoading, refetch } = useListMySecrets();
  const [putMySecretValue, { isLoading: isSubmitting }] = usePutMySecretValue();
  const [activeSecret, setActiveSecret] = useState<IMySecretListItem | null>(null);

  const secrets = data?.data ?? [];

  const handleSubmit = async (newValue: string) => {
    if (!activeSecret) return;
    try {
      const result = await putMySecretValue({
        secretId: activeSecret.managedSecretId,
        payload: { newValue },
      }).unwrap();
      const info = result.data;
      showSuccessToast(
        info.stellarTxHash
          ? `Signer value updated — on-chain swap submitted (tx ${info.stellarTxHash.slice(0, 10)}...)`
          : "Signer value updated",
      );
      setActiveSecret(null);
      refetch();
    } catch (error) {
      // 409 no_position_available and validation failures are real, expected outcomes here —
      // surface the backend's own message rather than a generic failure.
      showErrorToast(getErrorMessage(error, "Unable to update signer value"));
    }
  };

  const columns = [
    { title: "Label", dataIndex: "label" },
    {
      title: "Signer Position",
      dataIndex: "position",
      render: (value: number | null) => (value === null ? "Unclaimed" : `Position ${value}`),
    },
    {
      title: "Action",
      dataIndex: "managedSecretId",
      render: (_: string, record: IMySecretListItem) => (
        <ActionButton onClick={() => setActiveSecret(record)}>Set Value</ActionButton>
      ),
    },
  ];

  return (
    <SectionContainer>
      <SectionHeader>
        <SectionTitle>My Signer Keys</SectionTitle>
        <SectionSubtitle>
          Slots you have been assigned on shared, admin-managed signer secrets. Submitting a
          value here rotates the on-chain signer for active positions.
        </SectionSubtitle>
      </SectionHeader>

      {!isLoading && secrets.length === 0 ? (
        <EmptyState
          title="No signer slots yet"
          message="An admin hasn't assigned you to a managed signer secret yet."
        />
      ) : (
        <CustomTable columns={columns} dataSource={secrets} isLoading={isLoading} />
      )}

      {activeSecret && (
        <SetSignerValueModal
          label={activeSecret.label}
          isOpen={!!activeSecret}
          isSubmitting={isSubmitting}
          onClose={() => setActiveSecret(null)}
          onSubmit={handleSubmit}
        />
      )}
    </SectionContainer>
  );
};

export default MySignerSecrets;

const SectionContainer = styled.section`
  background-color: #ffffff;
  border-radius: 24px;
  padding: 24px;
  margin-bottom: 24px;
`;

const SectionHeader = styled.div`
  margin-bottom: 20px;
`;

const SectionTitle = styled.h2`
  font-size: 20px;
  font-weight: 700;
  margin: 0 0 4px 0;
  color: #00225a;
`;

const SectionSubtitle = styled.p`
  font-size: 13px;
  color: #828282;
  margin: 0;
  max-width: 640px;
`;

const ActionButton = styled.button`
  background: transparent;
  border: 1px solid #007cdf;
  color: #007cdf;
  border-radius: 8px;
  padding: 6px 14px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
`;
