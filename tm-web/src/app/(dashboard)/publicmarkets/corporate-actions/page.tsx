"use client";
import React, { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import { Modal } from "@/components/CustomModal";
import {
  ICorporateAction,
  IConfirmationRow,
  IPartnerEvent,
  useDeclarePMActionMutation,
  useGetPMActionsQuery,
  useGetPMAssetsQuery,
  useGetPMConfirmationsQuery,
  useGetPMPartnerEventsQuery,
  useGetPMPermissionsQuery,
} from "@/redux/api/publicMarkets";
import { act, Button, Card, date, Field, FormGrid, Input, Muted, ngn, Note, num, Pill, Row, SectionNav, Select, Tabs, Title } from "../components/ui";

const CorporateActions = () => {
  const params = useSearchParams();
  const router = useRouter();
  const tab = params.get("tab") ?? "events";
  const { data: me } = useGetPMPermissionsQuery();
  const [declare, setDeclare] = useState(false);
  return (
    <div>
      <SectionNav />
      <Card>
        <Row style={{ justifyContent: "space-between", marginTop: 0 }}>
          <div>
            <Title>Corporate Actions</Title>
            <Muted>
              Dividends, coupons and other events on Public Market Assets, from Custodian webhooks or declared here. Distributions run on their own engine: a
              record-date snapshot of every holder, withholding tax per owner, approval of the snapshot, then payment from the treasury (app holders) or credit to
              the exchange&apos;s balance (exchange customers).
            </Muted>
          </div>
          {me?.data?.manage && (
            <Button $variant="primary" onClick={() => setDeclare(true)}>
              Declare Corporate Action
            </Button>
          )}
        </Row>
        <Tabs
          tabs={[
            { key: "events", label: "Events & distributions" },
            { key: "confirmations", label: "Exchange confirmations" },
            { key: "manual", label: "Recorded in Trovo Manager" },
          ]}
          current={tab}
          onChange={(k) => router.push(`/publicmarkets/corporate-actions?tab=${k}`)}
        />
        {tab === "events" && <Events asset={params.get("asset") ?? ""} />}
        {tab === "confirmations" && <Confirmations />}
        {tab === "manual" && <ManualEntries />}
      </Card>
      {declare && <DeclareModal onClose={() => setDeclare(false)} />}
    </div>
  );
};

const Events = ({ asset: initial }: { asset: string }) => {
  const router = useRouter();
  const [status, setStatus] = useState("");
  const [asset, setAsset] = useState(initial);
  const [page, setPage] = useState(1);
  const { data, isFetching } = useGetPMActionsQuery({ status, asset, page, limit: 20 });
  return (
    <>
      <Row>
        <Input placeholder="Asset code" value={asset} onChange={(e) => setAsset(e.target.value)} />
        <Select value={status} onChange={(e) => setStatus(e.target.value)}>
          <option value="">All statuses</option>
          {["ANNOUNCED", "SNAPSHOTTED", "APPROVED", "PAYING", "DISTRIBUTED", "NEEDS_MANUAL", "CANCELLED"].map((s) => (
            <option key={s} value={s}>
              {s.replace("_", " ").toLowerCase()}
            </option>
          ))}
        </Select>
      </Row>
      <CustomTable
        columns={[
          { title: "Asset", dataIndex: "assetCode", render: (v: string) => <strong>{v}</strong> },
          { title: "Event", dataIndex: "eventType", render: (v: string, c: ICorporateAction) => `${v}${c.description ? ` · ${c.description}` : ""}` },
          { title: "Source", dataIndex: "source", render: (v: string) => (v === "MANUAL" ? "Trovo Manager" : "Custodian webhook") },
          { title: "Record date", dataIndex: "recordDate" },
          { title: "Pay date", dataIndex: "payDate" },
          { title: "Per unit", dataIndex: "amountPerUnit", render: (v: string) => (Number(v) > 0 ? ngn(v) : "—") },
          { title: "Snapshot", dataIndex: "holderCount", render: (v: number, c: ICorporateAction) => (c.snapshotAt ? `${num(v)} holders · locked` : "—") },
          { title: "Net paid", dataIndex: "netAmount", render: (v: string) => (v ? ngn(v) : "—") },
          { title: "Status", dataIndex: "status", render: (v: string) => <Pill status={v} /> },
        ]}
        dataSource={data?.data?.actions ?? []}
        isLoading={isFetching && !data}
        onRowClick={(c: ICorporateAction) => router.push(`/publicmarkets/corporate-actions/${c.id}`)}
      />
      <Pagination currentPage={page} totalCount={data?.data?.total ?? 0} pageSize={20} onPageChange={setPage} onPageSizeChange={() => undefined} />
    </>
  );
};

const Confirmations = () => {
  const router = useRouter();
  const { data, isFetching } = useGetPMConfirmationsQuery();
  return (
    <>
      <Note $tone="warn">
        Exchanges confirm each wallet&apos;s <code>dividend.paid</code> through <code>POST /v1/trovo-api/public-markets/confirmations</code>. Confirmations missing
        after the SLA (Settings) are escalated by the hourly sweep.
      </Note>
      <div style={{ height: 14 }} />
      <CustomTable
        columns={[
          { title: "Exchange", dataIndex: "exchangeName" },
          { title: "Asset", dataIndex: "assetCode" },
          { title: "Wallets paid", dataIndex: "walletsPaid", render: (v: number) => num(v) },
          { title: "Confirmed", dataIndex: "confirmed", render: (v: number) => num(v) },
          { title: "Outstanding", dataIndex: "outstanding", render: (v: number) => num(v) },
          { title: "Due", dataIndex: "oldestDueAt", render: (v: string) => date(v) },
          {
            title: "Status",
            dataIndex: "escalated",
            render: (v: number, r: IConfirmationRow) => <Pill status={v > 0 ? "ESCALATED" : r.outstanding > 0 ? "pending" : "COMPLETE"} label={v > 0 ? `${v} past SLA` : undefined} />,
          },
        ]}
        dataSource={data?.data ?? []}
        isLoading={isFetching && !data}
        onRowClick={(r: IConfirmationRow) => router.push(`/publicmarkets/exchanges/${encodeURIComponent(r.serviceLinkId)}`)}
      />
    </>
  );
};

const ManualEntries = () => {
  const [page, setPage] = useState(1);
  const { data } = useGetPMPartnerEventsQuery({ source: "MANUAL", page, limit: 20 }, { pollingInterval: 10000 });
  return (
    <>
      <Muted>What Operations recorded here (declared corporate actions, partner outcomes) and how the engine applied it.</Muted>
      <div style={{ height: 14 }} />
      <CustomTable
        columns={[
          { title: "Recorded", dataIndex: "createdAt", render: (v: string) => date(v) },
          { title: "By", dataIndex: "source", render: (v: string) => v.replace(/^MANUAL:?/, "") },
          { title: "Kind", dataIndex: "kind" },
          { title: "Reference", dataIndex: "reference" },
          { title: "Result", dataIndex: "result", render: (v: string, e: IPartnerEvent) => (e.processedAt ? v : "Waiting for the engine…") },
        ]}
        dataSource={data?.data?.events ?? []}
      />
      <Pagination currentPage={page} totalCount={data?.data?.total ?? 0} pageSize={20} onPageChange={setPage} onPageSizeChange={() => undefined} />
    </>
  );
};

const DeclareModal = ({ onClose }: { onClose: () => void }) => {
  const { data: assets } = useGetPMAssetsQuery({});
  const [declare, { isLoading }] = useDeclarePMActionMutation();
  const [v, setV] = useState<Record<string, string>>({ eventType: "DIVIDEND" });
  const bind = (k: string) => ({ value: v[k] ?? "", onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => setV({ ...v, [k]: e.target.value }) });
  const cash = v.eventType === "DIVIDEND" || v.eventType === "COUPON";
  const submit = async () => {
    const res = await act(
      () =>
        declare({
          assetCode: v.assetCode ?? "",
          eventType: v.eventType,
          recordDate: v.recordDate ?? "",
          payDate: v.payDate,
          amountPerUnit: v.amountPerUnit,
          description: v.description,
          reference: v.reference,
        }).unwrap(),
      { success: "Recorded: it appears in the list within seconds" },
    );
    if (res) onClose();
  };
  return (
    <Modal
      title="Declare Corporate Action"
      isOpen
      onClose={onClose}
      style={{ maxWidth: 680, width: "95%" }}
      actions={
        <Button $variant="primary" disabled={isLoading} onClick={submit}>
          Declare
        </Button>
      }
    >
      <FormGrid>
        <Field label="Asset">
          <select {...bind("assetCode")}>
            <option value="">Choose…</option>
            {(assets?.data ?? []).map((a) => (
              <option key={a.id} value={a.assetCode}>
                {a.assetCode} · {a.instrumentName}
              </option>
            ))}
          </select>
        </Field>
        <Field label="Event type" hint="Bonus, rights and splits are recorded for manual handling">
          <select {...bind("eventType")}>
            {["DIVIDEND", "COUPON", "BONUS", "RIGHTS", "SPLIT"].map((t) => (
              <option key={t} value={t}>
                {t}
              </option>
            ))}
          </select>
        </Field>
        <Field label="Record date" hint="YYYY-MM-DD (WAT); holders at its end are entitled">
          <input {...bind("recordDate")} placeholder="2026-10-10" />
        </Field>
        {cash && (
          <>
            <Field label="Pay date" hint="YYYY-MM-DD">
              <input {...bind("payDate")} placeholder="2026-10-24" />
            </Field>
            <Field label="Amount per unit (NGN)" hint="Gross, before withholding tax">
              <input {...bind("amountPerUnit")} placeholder="1.00" />
            </Field>
          </>
        )}
        <Field label="Description">
          <input {...bind("description")} placeholder="Interim dividend 2026" />
        </Field>
        <Field label="Custodian / registrar reference">
          <input {...bind("reference")} />
        </Field>
      </FormGrid>
    </Modal>
  );
};

const CorporateActionsPage = () => (
  <Suspense fallback={null}>
    <CorporateActions />
  </Suspense>
);

export default CorporateActionsPage;
