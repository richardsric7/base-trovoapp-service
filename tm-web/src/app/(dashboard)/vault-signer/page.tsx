"use client";

import React from "react";
import styled from "styled-components";
import MySignerSecrets from "@/app/vault-signer/components/MySignerSecrets";
import PersonalSecrets from "@/app/vault-signer/components/PersonalSecrets";
import {
  useListMySecretsQuery,
  usePutMySecretValueMutation,
  useListPersonalEnvsQuery,
  useCreatePersonalEnvMutation,
  useDeletePersonalEnvMutation,
} from "@/redux/api/vaultSigner/api";

// Dashboard side of the self-service vault-signer pages (/me/vault-signer/*), authenticated via
// baseApi's Bearer-JWT axios interceptor. A Trovo Admin's AdminUser.Username is always present,
// so there is no wallet-link gate here — that only applies to org members on the organisation
// portal counterpart (src/app/organisation/vault-signer/page.tsx).
const VaultSignerPage = () => {
  return (
    <PageContainer>
      <Title>Vault Signer</Title>
      <MySignerSecrets
        useListMySecrets={useListMySecretsQuery}
        usePutMySecretValue={usePutMySecretValueMutation}
      />
      <PersonalSecrets
        useListPersonalEnvs={useListPersonalEnvsQuery}
        useCreatePersonalEnv={useCreatePersonalEnvMutation}
        useDeletePersonalEnv={useDeletePersonalEnvMutation}
      />
    </PageContainer>
  );
};

export default VaultSignerPage;

const PageContainer = styled.section``;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 20px 0;
  color: #00225a;
`;
