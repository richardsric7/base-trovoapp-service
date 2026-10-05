"use client";
import React, { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import { Modal } from "@/components/CustomModal";
import {
  ICustodianRow,
  IDealingMember,
  IPMSettings,
  PartnerInput,
  useConfigurePMCustodianMutation,
  useGetPMCustodiansQuery,
  useGetPMDealingMembersQuery,
  useGetPMPermissionsQuery,
  useGetPMSettingsQuery,
  useSavePMDealingMemberMutation,
  useUpdatePMSettingsMutation,
} from "@/redux/api/publicMarkets";
import { act, Button, Card, date, Field, FormGrid, LinkButton, Mono, Muted, Note, Pill, Row, SectionNav, SubTitle, Tabs, Title } from "../components/ui";

const Settings = () => {
  const params = useSearchParams();
  const router = useRouter();
  const tab = params.get("tab") ?? "dm";
  const { data: me } = useGetPMPermissionsQuery();
  const canEdit = !!me?.data?.settings;
  return (
    <div>
      <SectionNav />
      <Card>
        <Title>Public Markets Settings</Title>
        <Muted>{canEdit ? "Changes apply to the engine within a minute." : "Read only: changing settings needs the MANAGE_SETTINGS permission."}</Muted>
        <Tabs
          tabs={[
            { key: "dm", label: "Dealing Members" },
            { key: "cust", label: "Custodians" },
            { key: "thr", label: "Thresholds & Limits" },
            { key: "appr", label: "Approvers" },
            { key: "fees", label: "Fees & Tax" },
            { key: "market", label: "Market Hours" },
          ]}
          current={tab}
          onChange={(k) => router.push(`/publicmarkets/settings?tab=${k}`)}
        />
        {tab === "dm" && <DealingMembers canEdit={canEdit} />}
        {tab === "cust" && <Custodians canEdit={canEdit} />}
        {["thr", "appr", "fees", "market"].includes(tab) && <SettingsForm section={tab} canEdit={canEdit} />}
      </Card>
    </div>
  );
};

const integration = (p: { mode: string; authScheme?: string; baseUrl?: string }) => (
  <span>
    <Pill status={p.mode === "MOCK" ? "sandbox" : p.mode === "MANUAL" ? "SETUP" : "LIVE"} label={p.mode} />
    {p.mode === "REST" && (
      <small>
        {" "}
        {p.authScheme} · {p.baseUrl}
      </small>
    )}
  </span>
);

const DealingMembers = ({ canEdit }: { canEdit: boolean }) => {
  const { data } = useGetPMDealingMembersQuery();
  const [edit, setEdit] = useState<IDealingMember | "new" | null>(null);
  return (
    <>
      <Row style={{ justifyContent: "space-between" }}>
        <Muted>Licensed brokers that execute the net buys and sells on NGX and FMDQ. Trovotech works with several in parallel; each asset uses one.</Muted>
        {canEdit && (
          <Button $variant="primary" onClick={() => setEdit("new")}>
            Add Dealing Member
          </Button>
        )}
      </Row>
      <CustomTable
        columns={[
          { title: "Name", dataIndex: "dealingMemberName", render: (v: string, d: IDealingMember) => <strong>{v}{!d.active && " (inactive)"}</strong> },
          { title: "Code", dataIndex: "code", render: (v: string) => <Mono>{v}</Mono> },
          { title: "Country", dataIndex: "dealingMemberCountry" },
          { title: "CSCS member code", dataIndex: "cscsMemberCode", render: (v: string) => <Mono>{v || "—"}</Mono> },
          { title: "Fee", dataIndex: "feePercent", render: (v: number, d: IDealingMember) => `${v}%${d.feeFixed ? ` + ₦${d.feeFixed}` : ""}` },
          { title: "Integration", dataIndex: "mode", render: (_: string, d: IDealingMember) => integration(d) },
          { title: "Assets", dataIndex: "assets" },
          { title: "", dataIndex: "id", render: (_: number, d: IDealingMember) => canEdit && <LinkButton onClick={() => setEdit(d)}>Edit</LinkButton> },
        ]}
        dataSource={data?.data ?? []}
      />
      {edit && <PartnerModal kind="dm" dm={edit === "new" ? undefined : edit} onClose={() => setEdit(null)} />}
    </>
  );
};

const Custodians = ({ canEdit }: { canEdit: boolean }) => {
  const { data } = useGetPMCustodiansQuery();
  const [edit, setEdit] = useState<ICustodianRow | null>(null);
  return (
    <>
      <Muted>
        The Approved Asset Custodians, with their Public Markets integration. Each holds the shares through its own nominee company; CSCS is reached only through
        the Custodian. Until a Custodian&apos;s API is agreed, run it as MOCK (simulated) or MANUAL (Operations record confirmations from its portal).
      </Muted>
      <div style={{ height: 14 }} />
      <CustomTable
        columns={[
          { title: "Custodian", dataIndex: "name", render: (v: string) => <strong>{v}</strong> },
          { title: "Code", dataIndex: "integration", render: (v: ICustodianRow["integration"]) => (v ? <Mono>{v.code}</Mono> : "Not configured") },
          { title: "Nominee", dataIndex: "integration", render: (v: ICustodianRow["integration"]) => v?.nomineeName || "—" },
          { title: "Transport", dataIndex: "integration", render: (v: ICustodianRow["integration"]) => v?.transport || "—" },
          { title: "Integration", dataIndex: "integration", render: (v: ICustodianRow["integration"]) => (v ? integration(v) : "—") },
          { title: "Credentials", dataIndex: "integration", render: (v: ICustodianRow["integration"]) => <Mono>{v?.credentialsRef || "—"}</Mono> },
          { title: "Assets", dataIndex: "assets" },
          { title: "", dataIndex: "id", render: (_: number, c: ICustodianRow) => canEdit && <LinkButton onClick={() => setEdit(c)}>{c.integration ? "Edit" : "Configure"}</LinkButton> },
        ]}
        dataSource={data?.data ?? []}
      />
      {edit && <PartnerModal kind="cust" custodian={edit} onClose={() => setEdit(null)} />}
    </>
  );
};

const PartnerModal = ({ kind, dm, custodian, onClose }: { kind: "dm" | "cust"; dm?: IDealingMember; custodian?: ICustodianRow; onClose: () => void }) => {
  const [saveDM, r1] = useSavePMDealingMemberMutation();
  const [saveCust, r2] = useConfigurePMCustodianMutation();
  const c = custodian?.integration;
  const [f, setF] = useState<PartnerInput>(
    kind === "dm"
      ? {
          name: dm?.dealingMemberName,
          address: dm?.dealingMemberAddress,
          country: dm?.dealingMemberCountry ?? "NGA",
          cscsMemberCode: dm?.cscsMemberCode,
          requirementDocument: dm?.requirementDocument,
          feePercent: dm?.feePercent ?? 0,
          feeFixed: dm?.feeFixed ?? 0,
          code: dm?.code,
          mode: dm?.mode ?? "MOCK",
          authScheme: dm?.authScheme,
          baseUrl: dm?.baseUrl,
          credentialsRef: dm?.credentialsRef,
          active: dm?.active ?? true,
        }
      : {
          code: c?.code,
          nomineeName: c?.nomineeName,
          mode: c?.mode ?? "MOCK",
          transport: c?.transport,
          authScheme: c?.authScheme,
          baseUrl: c?.baseUrl,
          credentialsRef: c?.credentialsRef,
          active: c?.active ?? true,
        },
  );
  const set = (k: keyof PartnerInput) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) =>
    setF({ ...f, [k]: k === "feePercent" || k === "feeFixed" ? Number(e.target.value) : k === "active" ? e.target.value === "true" : e.target.value });
  const submit = async () => {
    const res =
      kind === "dm"
        ? await act(() => saveDM({ id: dm?.id, ...f }).unwrap(), { success: "Dealing Member saved" })
        : await act(() => saveCust({ id: custodian!.id, ...f }).unwrap(), { success: "Custodian saved" });
    if (res) onClose();
  };
  return (
    <Modal
      title={kind === "dm" ? (dm ? `Edit ${dm.dealingMemberName}` : "Add Dealing Member") : `Public Markets integration · ${custodian?.name}`}
      isOpen
      onClose={onClose}
      style={{ maxWidth: 760, width: "95%" }}
      actions={
        <Button $variant="primary" disabled={r1.isLoading || r2.isLoading} onClick={submit}>
          Save
        </Button>
      }
    >
      <FormGrid>
        {kind === "dm" && (
          <>
            <Field label="Dealing Member name">
              <input value={f.name ?? ""} onChange={set("name")} placeholder="Broker Partner 4 Ltd" />
            </Field>
            <Field label="Address">
              <input value={f.address ?? ""} onChange={set("address")} />
            </Field>
            <Field label="Country">
              <input value={f.country ?? ""} onChange={set("country")} />
            </Field>
            <Field label="CSCS member code">
              <input value={f.cscsMemberCode ?? ""} onChange={set("cscsMemberCode")} placeholder="CSCS-DM-0000" />
            </Field>
            <Field label="Fee (%)">
              <input type="number" step="0.01" value={f.feePercent ?? 0} onChange={set("feePercent")} />
            </Field>
            <Field label="Fixed fee (NGN)">
              <input type="number" value={f.feeFixed ?? 0} onChange={set("feeFixed")} />
            </Field>
            <Field label="Requirement document (URL)">
              <input value={f.requirementDocument ?? ""} onChange={set("requirementDocument")} />
            </Field>
          </>
        )}
        {kind === "cust" && (
          <>
            <Field label="Nominee company">
              <input value={f.nomineeName ?? ""} onChange={set("nomineeName")} placeholder="Custodian A Nominees Ltd" />
            </Field>
            <Field label="Transport (label)">
              <input value={f.transport ?? ""} onChange={set("transport")} placeholder="REST + webhooks / SFTP end-of-day file" />
            </Field>
          </>
        )}
        <Field label="Partner code" hint="Sent as X-Partner-Code on its callbacks">
          <input value={f.code ?? ""} onChange={set("code")} placeholder="CUSTA" />
        </Field>
        <Field label="Integration mode" hint="MOCK simulates the partner; MANUAL means Operations record outcomes">
          <select value={f.mode} onChange={set("mode")}>
            <option value="MOCK">MOCK</option>
            <option value="REST">REST</option>
            <option value="MANUAL">MANUAL</option>
          </select>
        </Field>
        {f.mode === "REST" && (
          <>
            <Field label="Base URL">
              <input value={f.baseUrl ?? ""} onChange={set("baseUrl")} placeholder="https://api.partner.ng" />
            </Field>
            <Field label="Auth scheme">
              <select value={f.authScheme ?? "HMAC"} onChange={set("authScheme")}>
                <option value="HMAC">HMAC</option>
                <option value="MTLS">mTLS</option>
                <option value="NONE">None</option>
              </select>
            </Field>
          </>
        )}
        <Field label="Credentials reference" hint="env:NAME or vault://path#FIELD — never the secret itself">
          <input value={f.credentialsRef ?? ""} onChange={set("credentialsRef")} placeholder="vault://pm/custodians/a#CUSTODIAN_A" />
        </Field>
        <Field label="Active">
          <select value={String(f.active)} onChange={set("active")}>
            <option value="true">Active</option>
            <option value="false">Inactive</option>
          </select>
        </Field>
      </FormGrid>
    </Modal>
  );
};

type Key = keyof IPMSettings;
const SECTIONS: Record<string, { key: Key; label: string; hint?: string; type?: "number" | "text" | "list" }[]> = {
  thr: [
    { key: "netCreationThreshold", label: "Net daily creation approval threshold (NGN)", hint: "Per asset per day. Above it, net buys wait for approvals." },
    { key: "approvalsRequired", label: "Approvals required", type: "number" },
    { key: "batchIntervalMinutes", label: "Net batch interval (minutes)", type: "number" },
    { key: "priceStaleMinutes", label: "Price staleness (minutes, market hours)", type: "number" },
    { key: "instructionMaxAttempts", label: "Instruction retry cap", type: "number", hint: "Then moved to Orders › Escalations" },
    { key: "settlementSlaHours", label: "Settlement SLA (hours)", type: "number" },
    { key: "webhookMaxAttempts", label: "Webhook attempts before dead-letter", type: "number" },
    { key: "confirmationSlaHours", label: "Exchange confirmation SLA (hours)", type: "number", hint: "0: not set" },
    { key: "reconciliationHour", label: "Daily reconciliation hour (WAT)", type: "number" },
    { key: "substantialHoldingPercent", label: "Substantial holding disclosure (%)", hint: "CAMA s.120" },
    { key: "rateLimitTiers", label: "Rate-limit tiers", hint: 'JSON, requests per minute: {"Tier 1": 1200, "Tier 2": 3000}' },
  ],
  appr: [
    { key: "netCreationApprovers", label: "Net Creation Approvers", type: "list", hint: "Admin emails, comma-separated" },
    { key: "dividendApprovers", label: "Dividend Approvers", type: "list", hint: "Admin emails, comma-separated" },
  ],
  fees: [
    { key: "tradeFeePercent", label: "Creation / redemption fee (%)", hint: "R1: each way, unless the asset sets its own" },
    { key: "feeWallet", label: "Fee wallet", hint: "Fees are swept here from the treasury" },
    { key: "aumFeePercent", label: "Wrapper fee on AUM (% per year)", hint: "R2" },
    { key: "fxSpreadPercent", label: "FX spread (%)", hint: "R3" },
    { key: "revenueShareTiers", label: "Exchange revenue share tiers", hint: 'R4, JSON percent: {"Tier A": 20, "Tier B": 30}' },
    { key: "whtResidentPercent", label: "Withholding tax, residents (%)" },
    { key: "whtNonResidentPercent", label: "Withholding tax, non-residents (%)" },
    { key: "whtMissingTaxIdPercent", label: "Withholding tax, no tax ID (%)" },
    { key: "whtWallet", label: "Withholding tax wallet" },
  ],
  market: [
    { key: "ngxOpen", label: "NGX opens (WAT)" },
    { key: "ngxClose", label: "NGX closes (WAT)" },
    { key: "fmdqOpen", label: "FMDQ opens (WAT)" },
    { key: "fmdqClose", label: "FMDQ closes (WAT)" },
    { key: "marketHolidays", label: "Market holidays", type: "list", hint: "YYYY-MM-DD, comma-separated" },
  ],
};

const SettingsForm = ({ section, canEdit }: { section: string; canEdit: boolean }) => {
  const { data } = useGetPMSettingsQuery();
  const s = data?.data;
  if (!s) return <Muted>Loading…</Muted>;
  // keyed on the saved version, so the form starts again from what was saved
  return <SettingsFields key={`${section}-${s.updatedAt}`} section={section} canEdit={canEdit} saved={s} />;
};

const SettingsFields = ({ section, canEdit, saved: s }: { section: string; canEdit: boolean; saved: IPMSettings }) => {
  const [save, { isLoading }] = useUpdatePMSettingsMutation();
  const [f, setF] = useState<Partial<IPMSettings>>(s);
  const fields = SECTIONS[section];
  const submit = () => {
    const body: Partial<IPMSettings> = {};
    for (const fd of fields) {
      const v = f[fd.key];
      (body as Record<string, unknown>)[fd.key] = fd.type === "number" ? Number(v) : v ?? "";
    }
    act(() => save(body).unwrap(), { success: "Settings saved" });
  };
  return (
    <>
      {section === "appr" && (
        <Note>
          Approvers must be listed by their Trovo Manager email. A net batch above the threshold, and every dividend distribution, needs {s.approvalsRequired} of
          them.
        </Note>
      )}
      <div style={{ height: 14 }} />
      <FormGrid>
        {fields.map((fd) => (
          <Field key={fd.key} label={fd.label} hint={fd.hint}>
            {fd.type === "list" || fd.key === "rateLimitTiers" || fd.key === "revenueShareTiers" ? (
              <textarea
                rows={3}
                disabled={!canEdit}
                value={String(f[fd.key] ?? "")}
                onChange={(e) => setF({ ...f, [fd.key]: e.target.value })}
              />
            ) : (
              <input
                type={fd.type === "number" ? "number" : "text"}
                disabled={!canEdit}
                value={String(f[fd.key] ?? "")}
                onChange={(e) => setF({ ...f, [fd.key]: e.target.value })}
              />
            )}
          </Field>
        ))}
      </FormGrid>
      <Row style={{ justifyContent: "space-between" }}>
        <Muted>
          Last changed {date(s.updatedAt)} {s.updatedBy && `by ${s.updatedBy}`}
        </Muted>
        {canEdit && (
          <Button $variant="primary" disabled={isLoading} onClick={submit}>
            Save Changes
          </Button>
        )}
      </Row>
      {section === "fees" && (
        <>
          <SubTitle style={{ marginTop: 10 }}>Funding currency</SubTitle>
          <Muted>
            Orders are paid in {s.fundingAssetCode}. Exchange partners prefund a {s.fundingAssetCode} balance with the treasury.
          </Muted>
        </>
      )}
    </>
  );
};

const SettingsPage = () => (
  <Suspense fallback={null}>
    <Settings />
  </Suspense>
);

export default SettingsPage;
