"use client";
import React, { useEffect, useState } from "react";
import styled from "styled-components";
import { useParams, useRouter } from "next/navigation";
import { FaArrowLeft } from "react-icons/fa6";
import CustomTable from "@/components/CustomTable";
import { showErrorToast, showSuccessToast } from "@/components";
import {
  IPayoutBatch,
  PayoutStatus,
  useGetPayoutQuery,
  usePayoutActionMutation,
  useSetPayoutFeeMutation,
} from "@/redux/api/proceedPayouts";
import PayoutListTable from "../components/PayoutListTable";
import {
  amount,
  Button,
  Card,
  date,
  errorMessage,
  Grid,
  Input,
  Mono,
  Muted,
  Note,
  Row,
  Select,
  short,
  Stat,
  StatusPill,
  SubTitle,
  Title,
} from "../components/ui";

// statuses the engine is working through: refresh while in them
const WORKING: PayoutStatus[] = ["PREPARE_REQUESTED", "PREPARING", "FUNDING_CHECK_REQUESTED", "PAYING"];

type Action = "prepare" | "approve" | "reject" | "confirm-funding" | "pause" | "resume" | "cancel" | "retry-failed";

const ACTIONS: { action: Action; label: string; when: PayoutStatus[]; reason?: string; confirm?: string; variant?: "primary" | "danger" }[] = [
  { action: "prepare", label: "Prepare schedule", when: ["REGISTERED"], variant: "primary", confirm: "Snapshot the token holders and lock the schedule for approval?" },
  { action: "approve", label: "Approve schedule", when: ["LOCKED"], variant: "primary", confirm: "Approve this schedule, its amounts and its fee?" },
  { action: "confirm-funding", label: "Confirm funding", when: ["APPROVED"], variant: "primary", confirm: "Is the payout Safe funded? The engine checks it and starts paying." },
  { action: "prepare", label: "Prepare again", when: ["LOCKED", "APPROVED"], confirm: "Prepare the schedule again? Its approvals are dropped." },
  { action: "reject", label: "Reject schedule", when: ["LOCKED", "APPROVED"], reason: "Why is the schedule rejected?" },
  { action: "pause", label: "Pause", when: ["PAYING", "FUNDING_CHECK_REQUESTED"], reason: "Why pause this payout?" },
  { action: "resume", label: "Resume", when: ["PAUSED"], variant: "primary", confirm: "Resume? The payout Safe's balance is checked again first." },
  { action: "retry-failed", label: "Retry failed transfers", when: ["COMPLETED_WITH_FAILURES"], variant: "primary", confirm: "Pay the failed transfers again (after a funding check)?" },
  {
    action: "cancel",
    label: "Cancel payout",
    when: ["REGISTERED", "PREPARE_REQUESTED", "PREPARING", "LOCKED", "APPROVED", "FUNDING_CHECK_REQUESTED", "PAUSED"],
    reason: "Why cancel this payout for good? Holders already paid stay paid.",
    variant: "danger",
  },
];

const PayoutDetailPage = () => {
  const router = useRouter();
  const params = useParams();
  const id = Number(params.id);
  const [poll, setPoll] = useState(0);
  const { data, isLoading } = useGetPayoutQuery(id, { skip: !id, pollingInterval: poll });
  const [runAction, { isLoading: acting }] = usePayoutActionMutation();
  const [setFee, { isLoading: savingFee }] = useSetPayoutFeeMutation();
  const p = data?.data;

  const [feeType, setFeeType] = useState<"FIXED" | "PERCENT">("FIXED");
  const [feeValue, setFeeValue] = useState("0");
  const [feeCap, setFeeCap] = useState("0");
  useEffect(() => {
    if (!p) return;
    setPoll(WORKING.includes(p.status) ? 10000 : 0);
    setFeeType(p.feeType || "FIXED");
    setFeeValue(p.feeValue || "0");
    setFeeCap(p.feeCap || "0");
  }, [p]);

  if (isLoading || !p) {
    return <Card>{isLoading ? "Loading…" : "Payout not found."}</Card>;
  }

  const run = async (a: (typeof ACTIONS)[number]) => {
    let reason: string | undefined;
    if (a.reason) {
      reason = window.prompt(a.reason) ?? undefined;
      if (!reason?.trim()) return;
    } else if (a.confirm && !window.confirm(a.confirm)) {
      return;
    }
    try {
      const res = await runAction({ id, action: a.action, reason }).unwrap();
      showSuccessToast(res.message || "Done");
    } catch (err) {
      showErrorToast(errorMessage(err, "The action failed"));
    }
  };

  const saveFee = async () => {
    const msg =
      p.status === "REGISTERED"
        ? "Set this payout's fee? It applies when the schedule is prepared."
        : "Change the fee? The schedule is prepared again and needs new approvals.";
    if (!window.confirm(msg)) return;
    try {
      await setFee({ id, feeType, feeValue: feeValue.trim() || "0", feeCap: feeType === "PERCENT" ? feeCap.trim() || "0" : "0" }).unwrap();
      showSuccessToast("Fee set");
    } catch (err) {
      showErrorToast(errorMessage(err, "Could not set the fee"));
    }
  };

  const code = p.payoutAssetCode;
  const feeEditable = ["REGISTERED", "LOCKED", "APPROVED"].includes(p.status);
  const scheduled = p.lockedAt && !p.lockedAt.startsWith("0001");
  const feeTerms =
    p.feeType === "PERCENT"
      ? `${p.feeValue}% of the payout${Number(p.feeCap) > 0 ? `, at most ${amount(p.feeCap, code)}` : ""}`
      : amount(p.feeValue, code);

  const batchColumns = [
    { title: "Batch", dataIndex: "id" },
    { title: "Sent", dataIndex: "createdAt", render: (v: string) => date(v) },
    { title: "Transfers", dataIndex: "itemCount" },
    { title: "Safe nonce", dataIndex: "safeNonce" },
    { title: "Transaction", dataIndex: "txHash", render: (v: string) => <Mono>{short(v)}</Mono> },
    { title: "Gas", dataIndex: "gasUsed", render: (v: number) => (v ? v.toLocaleString() : "—") },
    { title: "Status", dataIndex: "status", render: (v: string, r: IPayoutBatch) => (
      <span>
        <StatusPill status={v} />
        {r.error && <small> {r.error}</small>}
      </span>
    ) },
  ];

  return (
    <>
      <Back onClick={() => router.push("/dividendandyield")}>
        <FaArrowLeft size={18} />
      </Back>
      <Card>
        <Header>
          <div>
            <Title>
              {p.assetCode || p.tokenizedAssetId} · {amount(p.totalAmount, code)}
            </Title>
            <Muted>
              {p.assetName} · payout #{p.id} ({p.batch}) · distribution {p.distributionId}
            </Muted>
          </div>
          <StatusPill status={p.status} />
        </Header>
        {p.note && <Note $tone={p.status === "APPROVED" && p.note.startsWith("Funding not confirmed") ? "warn" : "info"}>{p.note}</Note>}
        {p.engine && (p.engine.halted || !p.engine.online) && (
          <Note $tone="warn">
            {p.engine.halted ? `The payout engine is stopped: ${p.engine.haltReason}` : "The payout engine is not responding; actions wait until it is back."}
          </Note>
        )}
        <Row>
          {ACTIONS.filter((a) => a.when.includes(p.status)).map((a) => (
            <Button key={`${a.action}-${a.label}`} $variant={a.variant} disabled={acting} onClick={() => run(a)}>
              {a.label}
            </Button>
          ))}
        </Row>
      </Card>

      <Card>
        <SubTitle>Amounts</SubTitle>
        <Grid>
          <Stat label="Authorized" value={amount(p.totalAmount, code)} />
          <Stat label="Processing fee" value={scheduled ? amount(p.fee, code) : feeTerms} />
          <Stat label={`VAT on the fee${p.vatPercent ? ` (${p.vatPercent}%)` : ""}`} value={scheduled ? amount(p.vat, code) : "at the asset country's rate"} />
          <Stat label="To holders" value={scheduled ? amount(p.holderPayable, code) : "—"} />
          <Stat label="Per token" value={p.amountPerToken ? amount(p.amountPerToken, code) : "—"} />
          <Stat label="Paid to holders" value={amount(p.paid, code)} />
          <Stat label="Kept in the Safe" value={scheduled ? amount(p.retained, code) : "—"} />
          <Stat label="Holders" value={p.holderCount ? `${p.holderCount} (${p.paidCount} paid, ${p.failedCount} failed, ${p.excludedCount} excluded)` : "—"} />
          <Stat label="Snapshot block" value={p.snapshotBlock || (p.scannedBlock ? `scanning… ${p.scannedBlock}` : "—")} />
          <Stat label="Payout Safe" value={<Mono>{p.payoutSafeAddress || "—"}</Mono>} />
          <Stat label="Fee wallet" value={<Mono>{p.feeWallet || "—"}</Mono>} />
          <Stat label="VAT wallet" value={<Mono>{p.vatWallet || "—"}</Mono>} />
        </Grid>
      </Card>

      <Card>
        <SubTitle>Processing fee</SubTitle>
        <Muted>
          {feeTerms}
          {p.feeSetBy ? ` · set by ${p.feeSetBy}` : " · the default"}. The fee is approved with the schedule; the admin who
          sets it cannot approve.
        </Muted>
        {feeEditable && (
          <Row>
            <Select value={feeType} onChange={(e) => setFeeType(e.target.value as "FIXED" | "PERCENT")}>
              <option value="FIXED">Fixed fee</option>
              <option value="PERCENT">Percent fee</option>
            </Select>
            <Input placeholder={feeType === "PERCENT" ? "Percent" : `Amount (${code})`} value={feeValue} onChange={(e) => setFeeValue(e.target.value)} />
            {feeType === "PERCENT" && <Input placeholder={`Cap (${code}, 0 = none)`} value={feeCap} onChange={(e) => setFeeCap(e.target.value)} />}
            <Button disabled={savingFee} onClick={saveFee}>
              {savingFee ? "Saving…" : "Set fee"}
            </Button>
          </Row>
        )}
      </Card>

      <Card>
        <SubTitle>
          Approvals ({p.status === "LOCKED" ? `${p.approvals} of ${p.approvalsRequired}` : p.approvalsRequired + " required"})
        </SubTitle>
        <Muted>
          Prepared by {p.preparedBy || "—"}. Distinct admins approve the locked schedule; the admin who prepared it cannot.
          {p.scheduleChecksum && (
            <>
              {" "}
              Schedule checksum <Mono>{short(p.scheduleChecksum)}</Mono>.
            </>
          )}
        </Muted>
        <List>
          {(p.approvalList ?? []).map((a) => (
            <li key={a.id}>
              {a.adminEmail} · {date(a.createdAt)}
              {a.scheduleChecksum !== p.scheduleChecksum && " (an earlier schedule)"}
            </li>
          ))}
          {!(p.approvalList ?? []).length && <li>No approvals yet.</li>}
        </List>
      </Card>

      <PayoutListTable payout={p} />

      {(p.batches ?? []).length > 0 && (
        <Card>
          <SubTitle>Batches</SubTitle>
          <CustomTable columns={batchColumns} dataSource={p.batches} />
        </Card>
      )}
    </>
  );
};

export default PayoutDetailPage;

const Back = styled.button`
  display: block;
  color: #000000;
  background: none;
  cursor: pointer;
  border: none;
  padding: 24px;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
`;

const List = styled.ul`
  margin: 12px 0 0;
  padding-left: 20px;
  color: #00225a;
  font-size: 14px;
  line-height: 24px;
`;
