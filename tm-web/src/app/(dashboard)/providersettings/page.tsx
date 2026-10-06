"use client";
import React, { useState } from "react";
import styled from "styled-components";
import { showErrorToast, showSuccessToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import {
  KycProviderSettings,
  SecretStatus,
  StablerailSettings,
  useGetProviderSettingsQuery,
  useSaveKycProviderMutation,
  useSaveStablerailMutation,
} from "@/redux/api/providerSettings";

const errorMessage = (err: any, fallback: string) =>
  err?.data?.error || err?.data?.message || err?.error || fallback;

const providerNames: Record<string, string> = { sumsub: "Sumsub", doja: "Doja" };

const SecretState = ({ status }: { status?: SecretStatus }) =>
  status?.set ? <Set>Set ({status.hint})</Set> : <NotSet>Not set</NotSet>;

const KycProviderForm = ({ provider }: { provider: KycProviderSettings }) => {
  const [save, { isLoading }] = useSaveKycProviderMutation();
  const [token, setToken] = useState("");
  const [secretKey, setSecretKey] = useState("");
  const name = providerNames[provider.service_provider] ?? provider.service_provider;

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token.trim() && !secretKey.trim()) {
      showErrorToast("Enter a new token, a new secret key, or both.");
      return;
    }
    try {
      await save({ provider: provider.service_provider, token: token.trim(), secret_key: secretKey.trim() }).unwrap();
      showSuccessToast(`${name} credentials saved`);
      setToken("");
      setSecretKey("");
    } catch (err) {
      showErrorToast(errorMessage(err, `Could not save ${name} credentials`));
    }
  };

  return (
    <Card onSubmit={handleSave}>
      <CardTitle>{name}</CardTitle>
      <Row>
        <Label>App token</Label>
        <SecretState status={provider.token} />
      </Row>
      <Row>
        <Label>Secret key</Label>
        <SecretState status={provider.secret_key} />
      </Row>
      <Fields>
        <Input
          type="password"
          autoComplete="off"
          placeholder="New app token (leave empty to keep)"
          value={token}
          onChange={(e) => setToken(e.target.value)}
        />
        <Input
          type="password"
          autoComplete="off"
          placeholder="New secret key (leave empty to keep)"
          value={secretKey}
          onChange={(e) => setSecretKey(e.target.value)}
        />
        <PrimaryButton disabled={isLoading} buttonStyle={{ width: "auto", padding: "10px 20px" }}>
          {isLoading ? "Saving..." : "Save"}
        </PrimaryButton>
      </Fields>
    </Card>
  );
};

const StablerailForm = ({ current }: { current: StablerailSettings }) => {
  const [save, { isLoading }] = useSaveStablerailMutation();
  const [enabled, setEnabled] = useState(current.enabled);
  const [baseUrl, setBaseUrl] = useState(current.base_url);
  const [fintechId, setFintechId] = useState(current.fintech_id);
  const [apiKey, setApiKey] = useState("");

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await save({ enabled, base_url: baseUrl.trim(), fintech_id: fintechId.trim(), api_key: apiKey.trim() }).unwrap();
      showSuccessToast("Stablerail settings saved");
      setApiKey("");
    } catch (err) {
      showErrorToast(errorMessage(err, "Could not save Stablerail settings"));
    }
  };

  return (
    <Card onSubmit={handleSave}>
      <CardTitle>Stablerail (Naira bank deposits and withdrawals)</CardTitle>
      <Row>
        <Label>API key</Label>
        <SecretState status={current.api_key} />
      </Row>
      <Checkbox>
        <input type="checkbox" checked={enabled} onChange={(e) => setEnabled(e.target.checked)} />
        Bank deposits and withdrawals are on
      </Checkbox>
      <Fields>
        <Input
          placeholder="API address (e.g. https://api.stablerail.example.com)"
          value={baseUrl}
          onChange={(e) => setBaseUrl(e.target.value)}
          $wide
        />
        <Input placeholder="Fintech ID" value={fintechId} onChange={(e) => setFintechId(e.target.value)} />
        <Input
          type="password"
          autoComplete="off"
          placeholder="New API key (leave empty to keep)"
          value={apiKey}
          onChange={(e) => setApiKey(e.target.value)}
        />
        <PrimaryButton disabled={isLoading} buttonStyle={{ width: "auto", padding: "10px 20px" }}>
          {isLoading ? "Saving..." : "Save"}
        </PrimaryButton>
      </Fields>
    </Card>
  );
};

const ProviderSettingsPage = () => {
  const { data, isLoading, error } = useGetProviderSettingsQuery();

  return (
    <PageContainer>
      <Title>Provider Settings</Title>
      <Text>
        Credentials app-backend uses for outside providers: identity checks (KYC) and Naira bank
        deposits and withdrawals. Keys are never shown again after they are saved; type a new one to
        replace it. Changes apply to the next request and are recorded in the audit trail. Needs the
        MANAGE_SETTINGS permission.
      </Text>
      {isLoading && <Text>Loading...</Text>}
      {error ? <ErrorText>{errorMessage(error, "Could not load provider settings")}</ErrorText> : null}
      {data && (
        <>
          <SectionTitle>Identity checks (KYC)</SectionTitle>
          {data.kyc.map((p) => (
            <KycProviderForm key={p.service_provider} provider={p} />
          ))}
          <SectionTitle>Bank</SectionTitle>
          <StablerailForm current={data.stablerail} />
        </>
      )}
    </PageContainer>
  );
};

export default ProviderSettingsPage;

const PageContainer = styled.section`
  background: #ffffff;
  padding: 32px;
  border-radius: 24px;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;

const SectionTitle = styled.h2`
  font-size: 18px;
  font-weight: 700;
  margin: 28px 0 12px;
  color: #00225a;
`;

const Text = styled.p`
  font-size: 14px;
  color: #828282;
  margin: 8px 0 0;
  line-height: 22px;
  max-width: 820px;
`;

const ErrorText = styled.p`
  color: #d92d20;
  font-size: 14px;
  margin-top: 16px;
`;

const Card = styled.form`
  border: 1px solid #d9e1ec;
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 16px;
`;

const CardTitle = styled.h3`
  font-size: 16px;
  font-weight: 700;
  margin: 0 0 12px;
  color: #00225a;
`;

const Row = styled.div`
  display: flex;
  gap: 12px;
  font-size: 14px;
  margin-bottom: 6px;
`;

const Label = styled.span`
  width: 120px;
  color: #828282;
`;

const Set = styled.span`
  color: #027a48;
  font-weight: 600;
`;

const NotSet = styled.span`
  color: #b54708;
  font-weight: 600;
`;

const Checkbox = styled.label`
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #00225a;
  margin: 12px 0 0;
`;

const Fields = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  margin-top: 16px;
`;

const Input = styled.input<{ $wide?: boolean }>`
  flex: ${(p) => (p.$wide ? "1 1 320px" : "0 1 240px")};
  min-width: 0;
  padding: 10px 14px;
  border: 1px solid #d9e1ec;
  border-radius: 8px;
  font-size: 14px;
  color: #00225a;
  background: #f2f6f9;
`;
