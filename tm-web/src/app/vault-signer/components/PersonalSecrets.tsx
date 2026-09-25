"use client";

import React, { useState } from "react";
import styled from "styled-components";
import { showSuccessToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import { FaWallet } from "react-icons/fa6";
import PersonalSecretModal from "./PersonalSecretModal";
import { getErrorMessage } from "./MySignerSecrets";
import { showErrorToast } from "@/components";
import type {
  ICreatePersonalEnvRequest,
  ICreatePersonalEnvResponse,
  IPersonalEnvListItem,
  IPersonalEnvListResponse,
} from "@/redux/api/vaultSigner/interface";

// See MySignerSecrets.tsx for why these are loose structural types rather than `typeof
// useXQuery` — the dashboard's baseApi hooks and the org portal's orgApi hooks are bound to
// different BaseQueryFn generics and aren't directly assignable to one another under `typeof`.
type UseListPersonalEnvs = (
  arg?: void,
  options?: any,
) => { data?: IPersonalEnvListResponse; isLoading: boolean; refetch: () => void };

type UseCreatePersonalEnv = () => readonly [
  (arg: { prefix: string; payload: ICreatePersonalEnvRequest }) => {
    unwrap: () => Promise<ICreatePersonalEnvResponse>;
  },
  { isLoading: boolean },
];

type UseDeletePersonalEnv = () => readonly [
  (prefix: string) => { unwrap: () => Promise<void> },
  { isLoading: boolean },
];

interface WalletGate {
  // Only the organisation portal supplies this — an org member's TrovoWalletUsername is
  // opt-in, so personal secrets are gated behind the same "wallet linked" state OrgNavBar
  // already surfaces. Trovo Admins never hit this gate (AdminUser.Username is required), so
  // the dashboard page renders PersonalSecrets without this prop at all.
  isLinked: boolean;
  walletUsername?: string;
  onLinkWallet: () => void;
}

interface PersonalSecretsProps {
  useListPersonalEnvs: UseListPersonalEnvs;
  useCreatePersonalEnv: UseCreatePersonalEnv;
  useDeletePersonalEnv: UseDeletePersonalEnv;
  walletGate?: WalletGate;
}

const PersonalSecrets: React.FC<PersonalSecretsProps> = ({
  useListPersonalEnvs,
  useCreatePersonalEnv,
  useDeletePersonalEnv,
  walletGate,
}) => {
  const gated = !!walletGate && !walletGate.isLinked;

  const { data, isLoading, refetch } = useListPersonalEnvs(undefined, { skip: gated });
  const [createPersonalEnv, { isLoading: isCreating }] = useCreatePersonalEnv();
  const [deletePersonalEnv, { isLoading: isDeleting }] = useDeletePersonalEnv();

  const [createTarget, setCreateTarget] = useState<IPersonalEnvListItem | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<IPersonalEnvListItem | null>(null);

  const items = data?.data ?? [];

  const handleCreate = async (newValue: string) => {
    if (!createTarget) return;
    try {
      await createPersonalEnv({ prefix: createTarget.prefix, payload: { newValue } }).unwrap();
      showSuccessToast(`${createTarget.friendlyLabel} created`);
      setCreateTarget(null);
      refetch();
    } catch (error) {
      showErrorToast(getErrorMessage(error, "Unable to create personal secret"));
    }
  };

  const handleDelete = async () => {
    if (!deleteTarget) return;
    try {
      await deletePersonalEnv(deleteTarget.prefix).unwrap();
      showSuccessToast(`${deleteTarget.friendlyLabel} deleted`);
      setDeleteTarget(null);
      refetch();
    } catch (error) {
      showErrorToast(getErrorMessage(error, "Unable to delete personal secret"));
    }
  };

  return (
    <SectionContainer>
      <SectionHeader>
        <SectionTitle>Personal Secrets</SectionTitle>
        <SectionSubtitle>
          Create-only values scoped to you. There is no way to read a value back once saved — to
          change one, delete it and create it again.
        </SectionSubtitle>
      </SectionHeader>

      {gated ? (
        <GateCard>
          <GateIcon>
            <FaWallet />
          </GateIcon>
          <GateText>Connect your profile to the Trovo app to enable personal secrets.</GateText>
          <PrimaryButton buttonStyle={{ width: "auto", padding: "10px 24px" }} onClick={walletGate!.onLinkWallet}>
            Link Wallet
          </PrimaryButton>
        </GateCard>
      ) : (
        <List>
          {!isLoading &&
            items.map((item) => (
              <Row key={item.prefix}>
                <RowInfo>
                  <RowLabel>{item.friendlyLabel}</RowLabel>
                  <RowKey>{item.vaultKey}</RowKey>
                </RowInfo>
                <RowStatus>
                  <StatusPill exists={item.exists}>{item.exists ? "Exists" : "Not set"}</StatusPill>
                  {item.exists ? (
                    <ActionButton danger onClick={() => setDeleteTarget(item)}>
                      Delete
                    </ActionButton>
                  ) : (
                    <ActionButton onClick={() => setCreateTarget(item)}>Create</ActionButton>
                  )}
                </RowStatus>
              </Row>
            ))}
          {!isLoading && items.length === 0 && <EmptyText>No personal secret prefixes configured.</EmptyText>}
        </List>
      )}

      {createTarget && (
        <PersonalSecretModal
          mode="create"
          friendlyLabel={createTarget.friendlyLabel}
          isOpen={!!createTarget}
          isSubmitting={isCreating}
          onClose={() => setCreateTarget(null)}
          onConfirm={handleCreate}
        />
      )}
      {deleteTarget && (
        <PersonalSecretModal
          mode="delete"
          friendlyLabel={deleteTarget.friendlyLabel}
          isOpen={!!deleteTarget}
          isSubmitting={isDeleting}
          onClose={() => setDeleteTarget(null)}
          onConfirm={handleDelete}
        />
      )}
    </SectionContainer>
  );
};

export default PersonalSecrets;

const SectionContainer = styled.section`
  background-color: #ffffff;
  border-radius: 24px;
  padding: 24px;
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

const List = styled.div`
  display: flex;
  flex-direction: column;
`;

const Row = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 0;
  border-bottom: 1px solid #e5e5ef;

  &:last-child {
    border-bottom: none;
  }
`;

const RowInfo = styled.div`
  display: flex;
  flex-direction: column;
  gap: 2px;
`;

const RowLabel = styled.span`
  color: #00225a;
  font-size: 14px;
  font-weight: 600;
`;

const RowKey = styled.span`
  color: #828282;
  font-size: 12px;
`;

const RowStatus = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

const StatusPill = styled.span<{ exists: boolean }>`
  font-size: 12px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 999px;
  background: ${(props) => (props.exists ? "#e9f8f1" : "#f2f6f9")};
  color: ${(props) => (props.exists ? "#00875a" : "#828282")};
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

const EmptyText = styled.p`
  color: #828282;
  font-size: 13px;
  text-align: center;
  padding: 24px 0;
`;

const GateCard = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 48px 20px;
  text-align: center;
`;

const GateIcon = styled.span`
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: #eaf5fd;
  color: #007cdf;
  font-size: 18px;
`;

const GateText = styled.p`
  color: #00225a;
  font-size: 14px;
  max-width: 360px;
  margin: 0;
`;
