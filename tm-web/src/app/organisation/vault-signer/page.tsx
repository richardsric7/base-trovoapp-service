"use client";

import React, { useState } from "react";
import styled from "styled-components";
import MySignerSecrets from "@/app/vault-signer/components/MySignerSecrets";
import PersonalSecrets from "@/app/vault-signer/components/PersonalSecrets";
import LinkWalletModal from "@/app/organisation/components/LinkWalletModal";
import { useGetStakeholderProfileQuery } from "@/redux/api/sharedstakeholders";
import {
  useListMySecretsOrgQuery,
  usePutMySecretValueOrgMutation,
  useListPersonalEnvsOrgQuery,
  useCreatePersonalEnvOrgMutation,
  useDeletePersonalEnvOrgMutation,
} from "@/redux/api/vaultSigner/orgApi";

// Organisation-portal side of the self-service vault-signer pages (/me/vault-signer/*),
// authenticated via orgApi's raw-org-token axios interceptor. Personal secrets are keyed by
// OrganizationMember.TrovoWalletUsername on the backend, so unlinked members get the backend's
// 403 trovo_wallet_not_linked — this page checks the same wallet.is_wallet_linked flag
// OrgNavBar.tsx already surfaces and reuses its existing LinkWalletModal instead of a dead-end
// error panel, so the client-side check and the backend gate agree on the same source of truth.
const OrganisationVaultSignerPage = () => {
  const [openLinkWallet, setOpenLinkWallet] = useState(false);
  const { data: profileData, refetch: refetchProfile } = useGetStakeholderProfileQuery();
  const wallet = profileData?.data?.wallet;

  return (
    <PageContainer>
      <Title>Vault Signer</Title>
      <MySignerSecrets
        useListMySecrets={useListMySecretsOrgQuery}
        usePutMySecretValue={usePutMySecretValueOrgMutation}
      />
      <PersonalSecrets
        useListPersonalEnvs={useListPersonalEnvsOrgQuery}
        useCreatePersonalEnv={useCreatePersonalEnvOrgMutation}
        useDeletePersonalEnv={useDeletePersonalEnvOrgMutation}
        walletGate={{
          isLinked: !!wallet?.is_wallet_linked,
          walletUsername: wallet?.trovo_wallet_username,
          onLinkWallet: () => setOpenLinkWallet(true),
        }}
      />
      <LinkWalletModal
        isOpen={openLinkWallet}
        onClose={() => setOpenLinkWallet(false)}
        onLinked={refetchProfile}
      />
    </PageContainer>
  );
};

export default OrganisationVaultSignerPage;

const PageContainer = styled.section``;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 20px 0;
  color: #00225a;
`;
