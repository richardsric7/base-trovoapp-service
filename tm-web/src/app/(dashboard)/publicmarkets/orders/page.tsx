"use client";
import React, { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import { Modal } from "@/components/CustomModal";
import {
  IBatchView,
  IInstruction,
  IPMOrder,
  useGetPMBatchesQuery,
  useGetPMInstructionsQuery,
  useGetPMOrderQuery,
  useGetPMOrdersQuery,
  useGetPMPermissionsQuery,
  usePmBatchActionMutation,
  usePmInstructionActionMutation,
  useRecordPMOutcomeMutation,
} from "@/redux/api/publicMarkets";
import {
  act,
  Button,
  Card,
  date,
  Field,
  FormGrid,
  Grid,
  Input,
  LinkButton,
  Mono,
  Muted,
  ngn,
  Note,
  num,
  ORDER_STATES,
  PATH_LABELS,
  Pill,
  Row,
  SectionNav,
  Select,
  Stat,
  SubTitle,
  Tabs,
  Title,
} from "../components/ui";

const Orders = () => {
  const params = useSearchParams();
  const router = useRouter();
  const tab = params.get("tab") ?? "all";
  const { data: approvals } = useGetPMBatchesQuery({ status: "AWAITING_APPROVAL", limit: 100 });
  const { data: escalations } = useGetPMInstructionsQuery({ status: "ESCALATED", limit: 100 });
  const tabs = [
    { key: "all", label: "All Orders" },
    { key: "approval", label: `Awaiting Approval (${approvals?.data?.total ?? 0})` },
    { key: "escalations", label: `Escalations (${escalations?.data?.total ?? 0})` },
    { key: "batches", label: "Net Batches" },
  ];
  return (
    <div>
      <SectionNav />
      <Card>
        <Title>Orders</Title>
        <Muted>Creation and redemption orders from the Trovo App and exchange partners run through the same engine.</Muted>
        <Tabs tabs={tabs} current={tab} onChange={(k) => router.push(`/publicmarkets/orders?tab=${k}`)} />
        {tab === "all" && <AllOrders batch={params.get("batch") ?? ""} />}
        {tab === "approval" && <Approvals rows={approvals?.data?.batches ?? []} />}
        {tab === "escalations" && <Escalations rows={escalations?.data?.instructions ?? []} />}
        {tab === "batches" && <Batches />}
      </Card>
    </div>
  );
};

const AllOrders = ({ batch }: { batch: string }) => {
  const [f, setF] = useState({ search: "", type: "", state: "", channel: "", path: "", asset: "", batch });
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [open, setOpen] = useState("");
  const { data, isFetching } = useGetPMOrdersQuery({ ...f, page, limit });
  const set = (k: keyof typeof f) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setF({ ...f, [k]: e.target.value });
    setPage(1);
  };
  return (
    <>
      <Row>
        <Input $wide placeholder="Order ID, wallet, username or external ref" value={f.search} onChange={set("search")} />
        <Input placeholder="Asset code" value={f.asset} onChange={set("asset")} />
        <Select value={f.type} onChange={set("type")}>
          <option value="">Buy &amp; sell</option>
          <option value="CREATION">Creation (buy)</option>
          <option value="REDEMPTION">Redemption (sell)</option>
        </Select>
        <Select value={f.channel} onChange={set("channel")}>
          <option value="">All channels</option>
          <option value="TROVO_APP">Trovo App</option>
          <option value="EXCHANGE">Exchanges</option>
        </Select>
        <Select value={f.path} onChange={set("path")}>
          <option value="">All paths</option>
          <option value="FAST">Fast</option>
          <option value="SLOW">Slow</option>
          <option value="NETTED">Netted</option>
        </Select>
        <Select value={f.state} onChange={set("state")}>
          <option value="">All states</option>
          <option value="open">Open</option>
          {ORDER_STATES.map((s) => (
            <option key={s} value={s}>
              {s}
            </option>
          ))}
        </Select>
        {f.batch && <LinkButton onClick={() => setF({ ...f, batch: "" })}>Batch {f.batch} ✕</LinkButton>}
      </Row>
      <CustomTable
        columns={[
          { title: "Order ID", dataIndex: "id", render: (v: string) => <strong>{v}</strong> },
          { title: "Type", dataIndex: "type" },
          { title: "Asset", dataIndex: "assetCode" },
          { title: "Channel", dataIndex: "channel", render: (v: string, o: IPMOrder) => (v === "EXCHANGE" ? `Exchange · ${o.externalOrderRef || o.partnerWalletId}` : o.username || "Trovo App") },
          { title: "Quantity", dataIndex: "quantity", render: (v: string) => num(v) },
          { title: "Amount", dataIndex: "amount", render: (v: string) => ngn(v) },
          { title: "Path", dataIndex: "path", render: (v: string) => PATH_LABELS[v] ?? v },
          { title: "State", dataIndex: "state", render: (v: string) => <Pill status={v} /> },
          { title: "Created", dataIndex: "createdAt", render: (v: string) => date(v) },
        ]}
        dataSource={data?.data?.orders ?? []}
        isLoading={isFetching && !data}
        onRowClick={(o: IPMOrder) => setOpen(o.id)}
      />
      <Pagination
        currentPage={page}
        totalCount={data?.data?.total ?? 0}
        pageSize={limit}
        onPageChange={setPage}
        onPageSizeChange={(n) => {
          setLimit(n);
          setPage(1);
        }}
      />
      {open && <OrderModal id={open} onClose={() => setOpen("")} />}
    </>
  );
};

const OrderModal = ({ id, onClose }: { id: string; onClose: () => void }) => {
  const { data } = useGetPMOrderQuery(id);
  const d = data?.data;
  const o = d?.order;
  return (
    <Modal title={`Order ${id}`} isOpen onClose={onClose} style={{ maxWidth: 820, width: "95%" }}>
      {!o ? (
        "Loading…"
      ) : (
        <>
          <Grid>
            <Stat label="State" value={<Pill status={o.state} />} />
            <Stat label="Type" value={`${o.type} · ${PATH_LABELS[o.path] ?? o.path}`} />
            <Stat label="Asset" value={o.assetCode} />
            <Stat label="Channel" value={o.channel === "EXCHANGE" ? d.exchangeName || "Exchange" : "Trovo App"} />
            <Stat label="Customer" value={o.username || o.externalOrderRef || o.partnerWalletId || "—"} />
            <Stat label="Wallet" value={<Mono>{o.walletAddress}</Mono>} />
            <Stat label="Amount" value={ngn(o.amount)} />
            <Stat label="Trovo fee" value={ngn(o.fee)} />
            <Stat label={o.type === "CREATION" ? "Invested" : "Paid out"} value={ngn(o.netAmount)} />
            <Stat label="Quantity" value={num(o.quantity)} />
            <Stat label="Reference price" value={ngn(o.referencePrice)} />
            <Stat label="Executed price" value={o.executedPrice ? ngn(o.executedPrice) : "—"} />
            <Stat label="Net batch" value={o.batchId || "—"} />
            <Stat label="Payment tx" value={o.paymentTxHash ? <Mono>{o.paymentTxHash}</Mono> : "—"} />
            <Stat label="Mint / burn tx" value={o.tokenTxHash ? <Mono>{o.tokenTxHash}</Mono> : "—"} />
            <Stat label="Payout tx" value={o.payoutTxHash ? <Mono>{o.payoutTxHash}</Mono> : "—"} />
          </Grid>
          {d.lastError && (
            <Note $tone="warn">
              Last error ({d.attempts} attempt{d.attempts === 1 ? "" : "s"}): {d.lastError}
            </Note>
          )}
          <SubTitle style={{ marginTop: 20 }}>Timeline</SubTitle>
          {d.events.map((e, i) => (
            <div key={i} style={{ display: "flex", gap: 12, padding: "6px 0", borderTop: "1px solid #eef2f6", fontSize: 13 }}>
              <span style={{ width: 160, color: "#667085" }}>{date(e.at)}</span>
              <Pill status={e.state} />
              <span>{e.note}</span>
            </div>
          ))}
          {d.instructions.length > 0 && (
            <>
              <SubTitle style={{ marginTop: 20 }}>Batch instructions</SubTitle>
              <CustomTable
                columns={[
                  { title: "Instruction", dataIndex: "id", render: (v: string) => <Mono>{v}</Mono> },
                  { title: "Counterparty", dataIndex: "partnerName" },
                  { title: "Side", dataIndex: "side" },
                  { title: "Quantity", dataIndex: "quantity" },
                  { title: "Status", dataIndex: "status", render: (v: string) => <Pill status={v} /> },
                ]}
                dataSource={d.instructions}
              />
            </>
          )}
        </>
      )}
    </Modal>
  );
};

const Approvals = ({ rows }: { rows: IBatchView[] }) => {
  const { data: me } = useGetPMPermissionsQuery();
  const [batchAction, { isLoading }] = usePmBatchActionMutation();
  const approver = !!me?.data?.netCreationApprover && !!me?.data?.manage;
  return (
    <>
      <Note>
        A session&apos;s net buy above the daily threshold needs approvals from the Net Creation Approvers (Settings) before the Dealing Member order is released.
        Rejecting a batch refunds its orders.
      </Note>
      <div style={{ height: 14 }} />
      <CustomTable
        columns={[
          { title: "Batch", dataIndex: "id", render: (v: string) => <strong>{v}</strong> },
          { title: "Asset", dataIndex: "assetCode" },
          { title: "Net quantity", dataIndex: "quantity", render: (v: string, b: IBatchView) => `${b.side} ${num(v)}` },
          { title: "Value", dataIndex: "value", render: (v: string) => ngn(v) },
          { title: "Threshold", dataIndex: "threshold", render: (v: string) => ngn(v) },
          {
            title: "Approvals",
            dataIndex: "approvals",
            render: (_: unknown, b: IBatchView) => (
              <span title={b.approvals.map((a) => a.approver).join(", ")}>
                {b.approvals.length} of {b.approvalsRequired}
              </span>
            ),
          },
          {
            title: "",
            dataIndex: "status",
            render: (_: unknown, b: IBatchView) =>
              approver ? (
                <Row style={{ margin: 0 }}>
                  <LinkButton
                    disabled={isLoading}
                    onClick={() =>
                      act(() => batchAction({ id: b.id, action: "approve" }).unwrap(), {
                        confirm: `Approve ${b.side} ${b.quantity} ${b.assetCode} (${ngn(b.value)}) at reference ${ngn(b.referencePrice)}?`,
                      })
                    }
                  >
                    Approve
                  </LinkButton>
                  <LinkButton
                    $danger
                    disabled={isLoading}
                    onClick={() => act((reason) => batchAction({ id: b.id, action: "reject", reason }).unwrap(), { reason: "Why reject? Its orders are refunded." })}
                  >
                    Reject
                  </LinkButton>
                </Row>
              ) : (
                <small>Approvers only</small>
              ),
          },
        ]}
        dataSource={rows}
      />
    </>
  );
};

const Escalations = ({ rows }: { rows: IInstruction[] }) => {
  const { data: me } = useGetPMPermissionsQuery();
  const [action, { isLoading }] = usePmInstructionActionMutation();
  const [batchAction] = usePmBatchActionMutation();
  const [outcome, setOutcome] = useState<IInstruction | null>(null);
  const manage = !!me?.data?.manage;
  return (
    <>
      <Note $tone="warn">
        Instructions to Custodians and Dealing Members retry with backoff; after the attempt cap they land here. Retry, record the partner&apos;s outcome (for a
        MANUAL partner, or a webhook that never came), roll an unexecuted batch to the next session, or mark it handled.
      </Note>
      <div style={{ height: 14 }} />
      <CustomTable
        columns={[
          { title: "Instruction", dataIndex: "id", render: (v: string) => <Mono>{v}</Mono> },
          { title: "Counterparty", dataIndex: "partnerName" },
          { title: "Side", dataIndex: "side" },
          { title: "Asset", dataIndex: "assetCode" },
          { title: "Quantity", dataIndex: "quantity", render: (v: string) => num(v) },
          { title: "Attempts", dataIndex: "attempts", render: (v: number, i: IInstruction) => `${v} / ${i.maxAttempts}` },
          { title: "Last error", dataIndex: "lastError" },
          {
            title: "",
            dataIndex: "status",
            render: (_: unknown, i: IInstruction) =>
              manage && (
                <Row style={{ margin: 0 }}>
                  <LinkButton disabled={isLoading} onClick={() => act(() => action({ id: i.id, action: "retry" }).unwrap(), { confirm: `Re-send ${i.id}?` })}>
                    Retry
                  </LinkButton>
                  <LinkButton onClick={() => setOutcome(i)}>Record outcome</LinkButton>
                  {i.kind === "DM_ORDER" && (
                    <LinkButton
                      onClick={() =>
                        act((reason) => batchAction({ id: i.batchId, action: "roll", reason }).unwrap(), {
                          reason: "Why roll the batch to the next session (e.g. market closed)?",
                        })
                      }
                    >
                      Roll to next session
                    </LinkButton>
                  )}
                  <LinkButton
                    $danger
                    onClick={() => act((reason) => action({ id: i.id, action: "handled", reason }).unwrap(), { reason: "How was it handled?" })}
                  >
                    Mark handled
                  </LinkButton>
                </Row>
              ),
          },
        ]}
        dataSource={rows}
      />
      {outcome && <OutcomeModal instruction={outcome} onClose={() => setOutcome(null)} />}
    </>
  );
};

const OutcomeModal = ({ instruction: i, onClose }: { instruction: IInstruction; onClose: () => void }) => {
  const dm = i.kind === "DM_ORDER";
  const [record, { isLoading }] = useRecordPMOutcomeMutation();
  const [v, setV] = useState<Record<string, string>>({ status: dm ? "FILLED" : "settlement_final", executedQuantity: i.quantity, settledQuantity: i.quantity });
  const bind = (k: string) => ({ value: v[k] ?? "", onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => setV({ ...v, [k]: e.target.value }) });
  const submit = async () => {
    const res = await act(() => record({ id: i.id, status: v.status ?? "", ...v }).unwrap(), { success: "Recorded; the engine applies it within seconds" });
    if (res) onClose();
  };
  return (
    <Modal
      title={`${dm ? "Dealing Member fill" : "Custodian settlement"} · ${i.id}`}
      isOpen
      onClose={onClose}
      style={{ maxWidth: 620, width: "95%" }}
      actions={
        <Button $variant="primary" disabled={isLoading} onClick={submit}>
          Record
        </Button>
      }
    >
      <Muted>
        {i.partnerName} · {i.side} {i.quantity} {i.assetCode}. Enter exactly what the partner&apos;s portal or contract note shows.
      </Muted>
      <div style={{ height: 14 }} />
      <FormGrid>
        <Field label="Outcome">
          <select {...bind("status")}>
            {dm ? (
              <>
                <option value="FILLED">Filled</option>
                <option value="PARTIALLY_FILLED">Partially filled (escalates)</option>
                <option value="REJECTED">Rejected</option>
              </>
            ) : (
              <>
                <option value="settlement_final">Settled (settlement_final)</option>
                <option value="failed">Failed</option>
              </>
            )}
          </select>
        </Field>
        {dm ? (
          <>
            <Field label="Executed quantity">
              <input {...bind("executedQuantity")} />
            </Field>
            <Field label="Executed price (NGN)">
              <input {...bind("executedPrice")} />
            </Field>
          </>
        ) : (
          <>
            <Field label="Settled quantity">
              <input {...bind("settledQuantity")} />
            </Field>
            <Field label="Settlement date" hint="YYYY-MM-DD">
              <input {...bind("settlementDate")} />
            </Field>
            <Field label="Custodian reference">
              <input {...bind("custodianReference")} />
            </Field>
          </>
        )}
      </FormGrid>
    </Modal>
  );
};

const Batches = () => {
  const router = useRouter();
  const [status, setStatus] = useState("");
  const [page, setPage] = useState(1);
  const { data, isFetching } = useGetPMBatchesQuery({ status, page, limit: 20 });
  return (
    <>
      <Row>
        <Select value={status} onChange={(e) => setStatus(e.target.value)}>
          <option value="">All statuses</option>
          {["AWAITING_APPROVAL", "RELEASED", "EXECUTED", "SETTLED", "PROCESSED", "INTERNAL", "REJECTED", "FAILED"].map((s) => (
            <option key={s} value={s}>
              {s.replace(/_/g, " ").toLowerCase()}
            </option>
          ))}
        </Select>
      </Row>
      <CustomTable
        columns={[
          { title: "Batch", dataIndex: "id", render: (v: string) => <strong>{v}</strong> },
          { title: "Session", dataIndex: "sessionDate" },
          { title: "Net", dataIndex: "quantity", render: (v: string, b: IBatchView) => (b.side === "NONE" ? "Fully netted" : `${b.side} ${num(v)}`) },
          { title: "Reference", dataIndex: "referencePrice", render: (v: string) => ngn(v) },
          { title: "Executed", dataIndex: "executedPrice", render: (v: string) => (v ? ngn(v) : "—") },
          { title: "Orders", dataIndex: "creationOrders", render: (_: unknown, b: IBatchView) => `${b.creationOrders} buy · ${b.redemptionOrders} sell` },
          { title: "Status", dataIndex: "status", render: (v: string, b: IBatchView) => <span title={b.note}><Pill status={v} /></span> },
        ]}
        dataSource={data?.data?.batches ?? []}
        isLoading={isFetching && !data}
        onRowClick={(b: IBatchView) => router.push(`/publicmarkets/orders?tab=all&batch=${encodeURIComponent(b.id)}`)}
      />
      <Pagination currentPage={page} totalCount={data?.data?.total ?? 0} pageSize={20} onPageChange={setPage} onPageSizeChange={() => undefined} />
    </>
  );
};

const OrdersPage = () => (
  <Suspense fallback={null}>
    <Orders />
  </Suspense>
);

export default OrdersPage;
