"use client";
import React, { useState } from "react";
import { useRouter } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import { Modal } from "@/components/CustomModal";
import { ExchangeInput, IExchangeRow, useGetPMExchangesQuery, useGetPMPermissionsQuery, useGetPMSettingsQuery, useOnboardPMExchangeMutation } from "@/redux/api/publicMarkets";
import { act, Button, Card, date, Field, FormGrid, Input, Mono, Muted, ngn, Note, num, Pill, Row, SectionNav, Select, tierNames, Title } from "../components/ui";

const ExchangesPage = () => {
  const router = useRouter();
  const [status, setStatus] = useState("");
  const [search, setSearch] = useState("");
  const { data, isFetching } = useGetPMExchangesQuery({ status, search });
  const { data: me } = useGetPMPermissionsQuery();
  const [open, setOpen] = useState(false);
  return (
    <div>
      <SectionNav />
      <Card>
        <Row style={{ justifyContent: "space-between", marginTop: 0 }}>
          <div>
            <Title>Exchange Partners</Title>
            <Muted>
              Exchanges call <code>/v1/trovo-api/public-markets</code> with their service link&apos;s API key and sign every request. They open a wallet per
              customer, pay creations from a prefunded CNGN balance and receive signed webhooks.
            </Muted>
          </div>
          {me?.data?.manage && (
            <Button $variant="primary" onClick={() => setOpen(true)}>
              Onboard Exchange
            </Button>
          )}
        </Row>
        <Row>
          <Input $wide placeholder="Search" value={search} onChange={(e) => setSearch(e.target.value)} />
          <Select value={status} onChange={(e) => setStatus(e.target.value)}>
            <option value="">All</option>
            <option value="active">Active</option>
            <option value="suspended">Suspended</option>
          </Select>
        </Row>
        <CustomTable
          columns={[
            { title: "Exchange", dataIndex: "name", render: (v: string, e: IExchangeRow) => <strong>{v || e.serviceLinkId}</strong> },
            { title: "Status", dataIndex: "status", render: (v: string, e: IExchangeRow) => <><Pill status={v} /> <Pill status={e.environment} /></> },
            { title: "Rate-limit tier", dataIndex: "rateLimitTier", render: (v: string, e: IExchangeRow) => (v ? `${v} · ${e.rateLimitPerMinute}/min` : "Default") },
            { title: "Callback URL", dataIndex: "callbackUrl", render: (v: string) => <Mono>{v}</Mono> },
            { title: "Balance", dataIndex: "balance", render: (v: string) => ngn(v) },
            { title: "Wallets", dataIndex: "wallets", render: (v: number) => num(v) },
            { title: "Orders (30d)", dataIndex: "orders30d", render: (v: number) => num(v) },
            { title: "Dead-letter", dataIndex: "deadLetters", render: (v: number) => (v ? <Pill status="DEAD_LETTER" label={String(v)} /> : "0") },
            { title: "Partner since", dataIndex: "createdAt", render: (v: string) => date(v) },
          ]}
          dataSource={data?.data?.exchanges ?? []}
          isLoading={isFetching && !data}
          onRowClick={(e: IExchangeRow) => router.push(`/publicmarkets/exchanges/${encodeURIComponent(e.serviceLinkId)}`)}
        />
      </Card>
      {open && <Onboard candidates={data?.data?.candidates ?? []} onClose={() => setOpen(false)} />}
    </div>
  );
};

const Onboard = ({ candidates, onClose }: { candidates: { id: string; shortName: string; longName: string }[]; onClose: () => void }) => {
  const { data: settings } = useGetPMSettingsQuery();
  const [onboard, { isLoading }] = useOnboardPMExchangeMutation();
  const [f, setF] = useState<ExchangeInput>({ environment: "sandbox" });
  const [secret, setSecret] = useState("");
  const set = (k: keyof ExchangeInput) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => setF({ ...f, [k]: e.target.value });
  const tiers = tierNames(settings?.data?.rateLimitTiers);
  const revenue = tierNames(settings?.data?.revenueShareTiers);
  const submit = async () => {
    const res = await act(() => onboard(f).unwrap(), { success: "Exchange onboarded" });
    if (res) setSecret(res.data.signingSecret);
  };
  if (secret) {
    return (
      <Modal title="Signing secret" isOpen onClose={onClose} style={{ maxWidth: 620, width: "95%" }} actions={<Button $variant="primary" onClick={onClose}>I have stored it</Button>}>
        <Note $tone="warn">This is the only time the secret is shown. Send it to the exchange&apos;s technical contact through a secure channel.</Note>
        <p style={{ marginTop: 16 }}>
          <Mono>{secret}</Mono>
        </p>
        <Muted>The exchange signs each request: hex(HMAC-SHA256(secret, timestamp + &quot;.&quot; + body)), and verifies our webhooks the same way.</Muted>
      </Modal>
    );
  }
  return (
    <Modal
      title="Onboard Exchange"
      isOpen
      onClose={onClose}
      style={{ maxWidth: 720, width: "95%" }}
      actions={
        <Button $variant="primary" disabled={isLoading} onClick={submit}>
          Onboard
        </Button>
      }
    >
      <Muted>The exchange must already have a verified service link (Service Links). Its KYC, AML and sanctions screening stay with the exchange.</Muted>
      <div style={{ height: 14 }} />
      <FormGrid>
        <Field label="Service link">
          <select value={f.serviceLinkId ?? ""} onChange={set("serviceLinkId")}>
            <option value="">Choose…</option>
            {candidates.map((c) => (
              <option key={c.id} value={c.id}>
                {c.longName || c.shortName}
              </option>
            ))}
          </select>
        </Field>
        <Field label="Environment">
          <select value={f.environment} onChange={set("environment")}>
            <option value="sandbox">Sandbox</option>
            <option value="production">Production</option>
          </select>
        </Field>
        <Field label="Callback URL" hint="Webhooks are POSTed here">
          <input value={f.callbackUrl ?? ""} onChange={set("callbackUrl")} placeholder="https://" />
        </Field>
        <Field label="Funding wallet" hint="Deposits are credited only from this address">
          <input value={f.fundingAddress ?? ""} onChange={set("fundingAddress")} placeholder="0x…" />
        </Field>
        <Field label="Rate-limit tier" hint="Tiers are set in Settings">
          <select value={f.rateLimitTier ?? ""} onChange={set("rateLimitTier")}>
            <option value="">Default</option>
            {tiers.map((t) => (
              <option key={t}>{t}</option>
            ))}
          </select>
        </Field>
        <Field label="Revenue share tier">
          <select value={f.revenueShareTier ?? ""} onChange={set("revenueShareTier")}>
            <option value="">None</option>
            {revenue.map((t) => (
              <option key={t}>{t}</option>
            ))}
          </select>
        </Field>
        <Field label="Technical contact">
          <input value={f.techContact ?? ""} onChange={set("techContact")} />
        </Field>
      </FormGrid>
    </Modal>
  );
};

export default ExchangesPage;
