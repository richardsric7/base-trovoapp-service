"use client";

import React from "react";
import styled from "styled-components";
import { useParams, useRouter } from "next/navigation";
import { FaAngleLeft } from "react-icons/fa6";
import AssignmentsTable from "../components/AssignmentsTable";
import { useListManagedSecretsQuery } from "@/redux/api/vaultSigner/api";

// There is no GET /admin/vault-signer/managed-secrets/:id — only the list endpoint (which
// already carries every managed secret's assignments) and the assignments-only endpoint for
// one secret. This page reads the secret's own label/detail out of the already-fetched list
// query (same VAULT_SIGNER_MANAGED_SECRETS-tagged cache the list page populates) rather than
// inventing a singular fetch the backend doesn't expose.
const ManagedSecretAssignmentsPage = () => {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const managedSecretId = params.id;
  const { data } = useListManagedSecretsQuery();

  const secret = data?.data.find((row) => row.managedSecret.id === managedSecretId)?.managedSecret;

  return (
    <PageContainer>
      <BackLink onClick={() => router.push("/settings/vault-signer")}>
        <FaAngleLeft /> Managed Secrets
      </BackLink>
      <Title>{secret ? secret.label : "Assignments"}</Title>
      {secret && <Subtitle>{secret.wallet_address}</Subtitle>}
      <AssignmentsTable managedSecretId={managedSecretId} />
    </PageContainer>
  );
};

export default ManagedSecretAssignmentsPage;

const PageContainer = styled.section``;

const BackLink = styled.button`
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  gap: 6px;
  color: #007cdf;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  padding: 0;
  margin-bottom: 12px;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 4px 0;
  color: #00225a;
`;

const Subtitle = styled.p`
  font-size: 13px;
  color: #828282;
  margin: 0 0 20px 0;
  font-family: monospace;
`;
