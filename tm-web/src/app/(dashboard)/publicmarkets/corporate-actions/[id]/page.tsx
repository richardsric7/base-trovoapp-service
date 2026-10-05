"use client";
import React, { useState } from "react";
import { useParams } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import { IEntitlement, useGetPMActionQuery, useGetPMPermissionsQuery, usePmActionActionMutation } from "@/redux/api/publicMarkets";
import { act, Back, Button, Card, date, Grid, Mono, Muted, ngn, Note, num, Pill, Row, Select, Stat, SubTitle, Title } from "../../components/ui";

const CHANNEL: Record<string, string> = { TROVO_APP: "Trovo App", EXCHANGE: "Exchange", PLATFORM: "Platform (retained)" };

const ActionPage = () => {
  const { id } = useParams<{ id: string }>();
  const [status, setStatus] = useState("");
  const [page, setPage] = useState(1);
  const { data, isLoading } = useGetPMActionQuery({ id, status, page, limit: 50 }, { skip: !id, pollingInterval: 15000 });
  const { data: me } = useGetPMPermissionsQuery();
  const [run, { isLoading: acting }] = usePmActionActionMutation();
  const c = data?.data;
  if (isLoading || !c) return <Card>{isLoading ? "Loading…" : "Corporate action not found."}</Card>;
  const approved = c.approvals.filter((a) => a.checksum === c.snapshotChecksum).length;
  const iApproved = c.approvals.some((a) => a.approver === me?.data?.email);

  return (
    <div>
      <Back href="/publicmarkets/corporate-actions" label="Corporate Actions" />
      <Card>
        <Row style={{ justifyContent: "space-between", marginTop: 0 }}>
          <div>
            <Title>
              {c.assetCode} {c.eventType.toLowerCase()} <Pill status={c.status} />
            </Title>
            <Muted>
              {c.description} · {c.id}
            </Muted>
          </div>
          {me?.data?.manage && (
            <Row style={{ margin: 0 }}>
              {c.status === "SNAPSHOTTED" && me.data.dividendApprover && !iApproved && (
                <Button
                  $variant="primary"
                  disabled={acting}
                  onClick={() =>
                    act(() => run({ id: c.id, action: "approve" }).unwrap(), {
                      confirm: `Approve paying ${ngn(c.netAmount)} net (${ngn(c.whtAmount)} withholding tax) to ${c.holderCount} holders? Your approval is for this exact snapshot.`,
                    })
                  }
                >
                  Approve distribution
                </Button>
              )}
              {["ANNOUNCED", "SNAPSHOTTED", "APPROVED", "NEEDS_MANUAL"].includes(c.status) && (
                <Button $variant="danger" disabled={acting} onClick={() => act((reason) => run({ id: c.id, action: "cancel", reason }).unwrap(), { reason: "Why cancel?" })}>
                  Cancel
                </Button>
              )}
            </Row>
          )}
        </Row>
        {c.status === "NEEDS_MANUAL" && (
          <Note $tone="warn">Bonus issues, rights and splits change the supply and are not automated: handle them with engineering and the Custodian.</Note>
        )}
        {c.note && <Note>{c.note}</Note>}
        <div style={{ height: 16 }} />
        <Grid>
          <Stat label="Record date" value={c.recordDate} />
          <Stat label="Pay date" value={c.payDate || "—"} />
          <Stat label="Amount per unit" value={Number(c.amountPerUnit) > 0 ? ngn(c.amountPerUnit) : "—"} />
          <Stat label="Source" value={c.source === "MANUAL" ? `Trovo Manager (${c.declaredBy})` : `Custodian webhook ${c.sourceReference}`} />
          <Stat label="Snapshot" value={c.snapshotAt ? `${date(c.snapshotAt)} · block ${c.recordBlock}` : "After the record date"} />
          <Stat label="Holders" value={num(c.holderCount)} />
          <Stat label="Eligible units" value={num(c.eligibleUnits)} />
          <Stat label="Retained (platform)" value={num(c.retainedUnits)} />
          <Stat label="Gross" value={ngn(c.grossAmount)} />
          <Stat label="Withholding tax" value={ngn(c.whtAmount)} />
          <Stat label="Net" value={ngn(c.netAmount)} />
          <Stat label="Approvals" value={`${approved} of ${c.approvalsRequired || "—"}`} />
          <Stat label="Funded" value={date(c.fundedAt ?? undefined)} />
          <Stat label="WHT transfer" value={c.whtTxHash ? <Mono>{c.whtTxHash}</Mono> : "—"} />
          <Stat label="Checksum" value={c.snapshotChecksum ? <Mono>{c.snapshotChecksum.slice(0, 16)}…</Mono> : "—"} />
        </Grid>
        {c.approvals.length > 0 && (
          <Muted style={{ marginTop: 12 }}>
            Approved by{" "}
            {c.approvals.map((a) => `${a.approver}${a.checksum !== c.snapshotChecksum ? " (older snapshot, not counted)" : ""}`).join(", ")}. Approvers:{" "}
            {c.approvers.join(", ") || "none configured"}.
          </Muted>
        )}
      </Card>

      {c.byChannel.length > 0 && (
        <Card>
          <SubTitle>By channel</SubTitle>
          <CustomTable
            columns={[
              { title: "Channel", dataIndex: "channel", render: (v: string) => CHANNEL[v] ?? v },
              { title: "Recipients", dataIndex: "recipients", render: (v: number) => num(v) },
              { title: "Gross", dataIndex: "gross", render: (v: string) => ngn(v) },
              { title: "WHT", dataIndex: "wht", render: (v: string) => ngn(v) },
              { title: "Net", dataIndex: "net", render: (v: string) => ngn(v) },
              { title: "Paid", dataIndex: "paid", render: (v: number) => num(v) },
              { title: "Failed", dataIndex: "failed", render: (v: number) => num(v) },
            ]}
            dataSource={c.byChannel}
          />
        </Card>
      )}

      <Card>
        <Row style={{ justifyContent: "space-between", marginTop: 0 }}>
          <SubTitle>Entitlements</SubTitle>
          <Select value={status} onChange={(e) => setStatus(e.target.value)}>
            <option value="">All</option>
            {["PENDING", "QUEUED", "PAID", "FAILED", "RETAINED"].map((s) => (
              <option key={s} value={s}>
                {s.toLowerCase()}
              </option>
            ))}
          </Select>
        </Row>
        <CustomTable
          columns={[
            { title: "Wallet", dataIndex: "walletAddress", render: (v: string) => <Mono>{v}</Mono> },
            { title: "Holder", dataIndex: "username", render: (v: string, e: IEntitlement) => v || e.partnerWalletId || CHANNEL[e.channel] },
            { title: "Units", dataIndex: "units", render: (v: string) => num(v) },
            { title: "Gross", dataIndex: "grossAmount", render: (v: string) => ngn(v) },
            { title: "WHT", dataIndex: "whtPercent", render: (v: string, e: IEntitlement) => `${v}% · ${ngn(e.whtAmount)}` },
            { title: "Net", dataIndex: "netAmount", render: (v: string) => ngn(v) },
            { title: "Status", dataIndex: "status", render: (v: string, e: IEntitlement) => <span title={e.note}><Pill status={v} /></span> },
          ]}
          dataSource={c.entitlements}
        />
        <Pagination currentPage={page} totalCount={c.totalEntitlements} pageSize={50} onPageChange={setPage} onPageSizeChange={() => undefined} />
      </Card>
    </div>
  );
};

export default ActionPage;
