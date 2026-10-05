"use client";
import React, { useState } from "react";
import { useParams } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import { Modal } from "@/components/CustomModal";
import {
  ExchangeInput,
  IWebhookDelivery,
  useGetPMExchangeQuery,
  useGetPMPermissionsQuery,
  useGetPMSettingsQuery,
  usePmExchangeActionMutation,
  useReplayPMWebhookMutation,
  useRequestPMWithdrawalMutation,
  useUpdatePMExchangeMutation,
} from "@/redux/api/publicMarkets";
import { act, Back, Button, Card, date, Field, FormGrid, Grid, LinkButton, Mono, Muted, ngn, Note, num, Pill, Row, Select, Stat, SubTitle, tierNames, Title } from "../../components/ui";

const ExchangePage = () => {
  const { id } = useParams<{ id: string }>();
  const [deliveries, setDeliveries] = useState("");
  const { data, isLoading } = useGetPMExchangeQuery({ id: decodeURIComponent(id), deliveries }, { skip: !id });
  const { data: me } = useGetPMPermissionsQuery();
  const [action, { isLoading: acting }] = usePmExchangeActionMutation();
  const [withdraw] = useRequestPMWithdrawalMutation();
  const [replay] = useReplayPMWebhookMutation();
  const [edit, setEdit] = useState(false);
  const [secret, setSecret] = useState("");
  const e = data?.data;
  if (isLoading || !e) return <Card>{isLoading ? "Loading…" : "Exchange not found."}</Card>;
  const manage = !!me?.data?.manage;

  const rotate = async () => {
    const res = await act(() => action({ id: e.serviceLinkId, action: "rotate-secret" }).unwrap(), {
      confirm: "Issue a new signing secret? The current one keeps working for 24 hours.",
    });
    if (res?.data?.signingSecret) setSecret(res.data.signingSecret);
  };

  return (
    <div>
      <Back href="/publicmarkets/exchanges" label="Exchange Partners" />
      <Card>
        <Row style={{ justifyContent: "space-between", marginTop: 0 }}>
          <div>
            <Title>
              {e.name || e.serviceLinkId} <Pill status={e.status} /> <Pill status={e.environment} />
            </Title>
            <Muted>Partner since {date(e.createdAt)}</Muted>
          </div>
          {manage && (
            <Row style={{ margin: 0 }}>
              <Button onClick={() => setEdit(true)}>Edit</Button>
              <Button onClick={rotate} disabled={acting}>
                Rotate credentials
              </Button>
              <Button
                onClick={() =>
                  act((amount) => withdraw({ id: e.serviceLinkId, amount: amount ?? "" }).unwrap(), {
                    reason: `Amount to pay back to ${e.fundingAddress} (balance ${ngn(e.balance)})?`,
                  })
                }
              >
                Pay balance back
              </Button>
              {e.status === "active" ? (
                <Button
                  $variant="danger"
                  onClick={() => act(() => action({ id: e.serviceLinkId, action: "suspend" }).unwrap(), { confirm: "Suspend? Its API calls are refused with 403." })}
                >
                  Suspend
                </Button>
              ) : (
                <Button $variant="primary" onClick={() => act(() => action({ id: e.serviceLinkId, action: "activate" }).unwrap(), { confirm: "Reactivate this exchange?" })}>
                  Activate
                </Button>
              )}
            </Row>
          )}
        </Row>
        {e.previousSecretValidUntil && new Date(e.previousSecretValidUntil) > new Date() && (
          <Note>The previous signing secret keeps working until {date(e.previousSecretValidUntil)}.</Note>
        )}
        <SubTitle style={{ marginTop: 20 }}>Integration</SubTitle>
        <Grid>
          <Stat label="Service link" value={<Mono>{e.serviceLinkId}</Mono>} />
          <Stat label="Signing" value="HMAC-SHA256 · Authorization: HMAC {key}:{signature}" />
          <Stat label="Callback URL" value={<Mono>{e.callbackUrl}</Mono>} />
          <Stat label="Rate-limit tier" value={e.rateLimitTier ? `${e.rateLimitTier} · ${e.rateLimitPerMinute} req/min` : "Default"} />
          <Stat label="Revenue share tier" value={e.revenueShareTier || "—"} />
          <Stat label="Confirmation SLA" value={e.confirmationSlaHours ? `${e.confirmationSlaHours}h` : "Settings' SLA"} />
          <Stat label="Technical contact" value={e.techContact || "—"} />
          <Stat label="Wallets provisioned" value={num(e.wallets)} />
          <Stat label="Orders (30 days)" value={num(e.orders30d)} />
        </Grid>
        <SubTitle style={{ marginTop: 20 }}>Prefunded balance</SubTitle>
        <Grid>
          <Stat label="Balance" value={ngn(e.balance)} />
          <Stat label="Funding wallet" value={<Mono>{e.fundingAddress}</Mono>} />
        </Grid>
      </Card>

      <Card>
        <Row style={{ justifyContent: "space-between", marginTop: 0 }}>
          <SubTitle>Webhook deliveries</SubTitle>
          <Row style={{ margin: 0 }}>
            <Select value={deliveries} onChange={(ev) => setDeliveries(ev.target.value)}>
              <option value="">All</option>
              <option value="PENDING">Pending</option>
              <option value="DELIVERED">Delivered</option>
              <option value="DEAD_LETTER">Dead-letter</option>
            </Select>
            {manage && e.deadLetters > 0 && (
              <Button
                onClick={() => act(() => action({ id: e.serviceLinkId, action: "replay-dead-letters" }).unwrap(), { confirm: `Replay ${e.deadLetters} dead-lettered deliveries?` })}
              >
                Replay all dead letters
              </Button>
            )}
          </Row>
        </Row>
        <Muted>At least once, with exponential backoff, then dead-lettered.</Muted>
        <CustomTable
          columns={[
            { title: "Event", dataIndex: "event", render: (v: string) => <Mono>{v}</Mono> },
            { title: "Reference", dataIndex: "reference", render: (v: string, w: IWebhookDelivery) => <Mono>{[v, w.walletId, w.assetCode].filter(Boolean).join(" · ")}</Mono> },
            { title: "Attempts", dataIndex: "attempts" },
            { title: "Last response", dataIndex: "lastResponseCode", render: (v: number, w: IWebhookDelivery) => (v ? String(v) : w.lastError || "—") },
            { title: "Status", dataIndex: "status", render: (v: string) => <Pill status={v} /> },
            {
              title: "Confirmed",
              dataIndex: "confirmedAt",
              render: (v: string, w: IWebhookDelivery) => (w.needsConfirmation ? (v ? date(v) : w.escalatedAt ? <Pill status="ESCALATED" label="past SLA" /> : "awaiting") : ""),
            },
            {
              title: "",
              dataIndex: "id",
              render: (v: string, w: IWebhookDelivery) =>
                manage && w.status === "DEAD_LETTER" ? <LinkButton onClick={() => act(() => replay(v).unwrap(), { success: "Replayed" })}>Replay</LinkButton> : "",
            },
          ]}
          dataSource={e.deliveries}
        />
      </Card>

      <Card>
        <SubTitle>Balance movements</SubTitle>
        <CustomTable
          columns={[
            { title: "When", dataIndex: "createdAt", render: (v: string) => date(v) },
            { title: "Kind", dataIndex: "kind" },
            { title: "Amount", dataIndex: "amount", render: (v: string) => <span style={{ color: Number(v) < 0 ? "#BE3800" : "#00A859" }}>{ngn(v)}</span> },
            { title: "Balance after", dataIndex: "balanceAfter", render: (v: string) => ngn(v) },
            { title: "Reference", dataIndex: "reference", render: (v: string) => <Mono>{v}</Mono> },
            { title: "Note", dataIndex: "note" },
          ]}
          dataSource={e.ledger}
        />
        {e.withdrawals.length > 0 && (
          <>
            <SubTitle style={{ marginTop: 18 }}>Withdrawal requests</SubTitle>
            <CustomTable
              columns={[
                { title: "Requested", dataIndex: "createdAt", render: (v: string) => date(v) },
                { title: "Amount", dataIndex: "target", render: (v: string) => ngn(v.split(":")[1]) },
                { title: "By", dataIndex: "requestedBy" },
                { title: "Result", dataIndex: "result", render: (v: string, j: { doneAt?: string | null }) => (j.doneAt ? v : "Waiting for the engine…") },
              ]}
              dataSource={e.withdrawals}
            />
          </>
        )}
      </Card>

      {edit && <EditExchange id={e.serviceLinkId} current={e} onClose={() => setEdit(false)} />}
      {secret && (
        <Modal title="New signing secret" isOpen onClose={() => setSecret("")} actions={<Button $variant="primary" onClick={() => setSecret("")}>I have stored it</Button>}>
          <Note $tone="warn">Shown once. The previous secret works for 24 more hours.</Note>
          <p style={{ marginTop: 16 }}>
            <Mono>{secret}</Mono>
          </p>
        </Modal>
      )}
    </div>
  );
};

const EditExchange = ({ id, current, onClose }: { id: string; current: ExchangeInput & { confirmationSlaHours: number }; onClose: () => void }) => {
  const { data: settings } = useGetPMSettingsQuery();
  const [update, { isLoading }] = useUpdatePMExchangeMutation();
  const [f, setF] = useState<ExchangeInput>({
    callbackUrl: current.callbackUrl,
    environment: current.environment,
    rateLimitTier: current.rateLimitTier,
    revenueShareTier: current.revenueShareTier,
    fundingAddress: current.fundingAddress,
    techContact: current.techContact,
    confirmationSlaHours: current.confirmationSlaHours,
  });
  const set = (k: keyof ExchangeInput) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) =>
    setF({ ...f, [k]: k === "confirmationSlaHours" ? Number(e.target.value) : e.target.value });
  const submit = async () => {
    const res = await act(() => update({ id, ...f }).unwrap(), { success: "Exchange updated" });
    if (res) onClose();
  };
  return (
    <Modal title="Edit exchange" isOpen onClose={onClose} style={{ maxWidth: 720, width: "95%" }} actions={<Button $variant="primary" disabled={isLoading} onClick={submit}>Save</Button>}>
      <FormGrid>
        <Field label="Callback URL">
          <input value={f.callbackUrl ?? ""} onChange={set("callbackUrl")} />
        </Field>
        <Field label="Environment">
          <select value={f.environment} onChange={set("environment")}>
            <option value="sandbox">Sandbox</option>
            <option value="production">Production</option>
          </select>
        </Field>
        <Field label="Funding wallet">
          <input value={f.fundingAddress ?? ""} onChange={set("fundingAddress")} />
        </Field>
        <Field label="Rate-limit tier" hint="Sets the service link's requests per minute">
          <select value={f.rateLimitTier ?? ""} onChange={set("rateLimitTier")}>
            <option value="">Default</option>
            {tierNames(settings?.data?.rateLimitTiers).map((t) => (
              <option key={t}>{t}</option>
            ))}
          </select>
        </Field>
        <Field label="Revenue share tier">
          <select value={f.revenueShareTier ?? ""} onChange={set("revenueShareTier")}>
            <option value="">None</option>
            {tierNames(settings?.data?.revenueShareTiers).map((t) => (
              <option key={t}>{t}</option>
            ))}
          </select>
        </Field>
        <Field label="Confirmation SLA (hours)" hint="0: the settings' SLA">
          <input type="number" value={f.confirmationSlaHours ?? 0} onChange={set("confirmationSlaHours")} />
        </Field>
        <Field label="Technical contact">
          <input value={f.techContact ?? ""} onChange={set("techContact")} />
        </Field>
      </FormGrid>
    </Modal>
  );
};

export default ExchangePage;
