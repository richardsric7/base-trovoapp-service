"use client";
import React, { useState } from "react";
import { useRouter } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import { IFeeCollection, IPayout, useGetPayoutFeesReportQuery, useGetPayoutsReportQuery } from "@/redux/api/proceedPayouts";
import { amount, Card, date, Grid, Input, Mono, Muted, Row, Select, short, Stat, STATUS_LABELS, StatusPill, SubTitle } from "./ui";

const toCSV = (rows: (string | number)[][]) =>
  rows.map((r) => r.map((v) => `"${String(v ?? "").replace(/"/g, '""')}"`).join(",")).join("\n");

const download = (name: string, rows: (string | number)[][]) => {
  const url = URL.createObjectURL(new Blob([toCSV(rows)], { type: "text/csv" }));
  const a = document.createElement("a");
  a.href = url;
  a.download = name;
  a.click();
  URL.revokeObjectURL(url);
};

// Payouts report (per payout and per currency) and the payout fees and VAT
// collected.
const Reports = () => {
  const router = useRouter();
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [asset, setAsset] = useState("");
  const [status, setStatus] = useState("");
  const filters = { from, to, asset: asset.trim() };
  const { data: payoutsData, isFetching: loadingPayouts } = useGetPayoutsReportQuery({ ...filters, status });
  const { data: feesData, isFetching: loadingFees } = useGetPayoutFeesReportQuery(filters);
  const report = payoutsData?.data;
  const fees = feesData?.data;

  const payoutColumns = [
    { title: "Registered", dataIndex: "createdAt", render: (v: string) => date(v) },
    { title: "Asset", dataIndex: "assetCode" },
    { title: "Batch", dataIndex: "batch" },
    { title: "Amount", dataIndex: "totalAmount", render: (v: string, r: IPayout) => amount(v, r.payoutAssetCode) },
    { title: "To holders", dataIndex: "holderPayable", render: (v: string) => amount(v) },
    { title: "Paid", dataIndex: "paid", render: (v: string) => amount(v) },
    { title: "Fee", dataIndex: "fee", render: (v: string) => amount(v) },
    { title: "VAT", dataIndex: "vat", render: (v: string) => amount(v) },
    { title: "Holders paid", dataIndex: "paidCount", render: (_: any, r: IPayout) => `${r.paidCount} (${r.failedCount} failed)` },
    { title: "Status", dataIndex: "status", render: (v: string) => <StatusPill status={v} /> },
  ];

  const feeColumns = [
    { title: "Date", dataIndex: "createdAt", render: (v: string) => date(v) },
    { title: "Asset", dataIndex: "assetCode" },
    { title: "Payout", dataIndex: "payoutBatch", render: (v: string) => (v || "").replace(/^payout /, "") },
    { title: "Type", dataIndex: "feeType", render: (v: string) => (v === "PROCEED_PAYOUT_FEE_VAT" ? "VAT" : "Fee") },
    { title: "Amount", dataIndex: "amount", render: (v: number, r: IFeeCollection) => amount(v, r.payoutAssetCode) },
    { title: "Paid to", dataIndex: "destinationWallet", render: (v: string) => <Mono>{short(v)}</Mono> },
    { title: "Transaction", dataIndex: "transactionHash", render: (v?: string) => <Mono>{short(v)}</Mono> },
  ];

  const exportPayouts = () =>
    download("payouts.csv", [
      ["registered", "asset", "batch", "currency", "amount", "to_holders", "paid", "fee", "vat", "holders_paid", "holders_failed", "status"],
      ...(report?.payouts ?? []).map((p) => [
        p.createdAt, p.assetCode, p.batch, p.payoutAssetCode, p.totalAmount, p.holderPayable, p.paid, p.fee, p.vat, p.paidCount, p.failedCount, p.status,
      ]),
    ]);
  const exportFees = () =>
    download("payout-fees.csv", [
      ["date", "asset", "payout", "type", "amount", "currency", "paid_to", "transaction"],
      ...(fees?.collections ?? []).map((c) => [
        c.createdAt, c.assetCode, c.payoutBatch, c.feeType, c.amount, c.payoutAssetCode, c.destinationWallet, c.transactionHash ?? "",
      ]),
    ]);

  return (
    <>
      <Card>
        <Row>
          <span>From</span>
          <Input type="date" aria-label="From" value={from} onChange={(e) => setFrom(e.target.value)} />
          <span>To</span>
          <Input type="date" aria-label="To" value={to} onChange={(e) => setTo(e.target.value)} />
          <Input placeholder="Asset code or ID" value={asset} onChange={(e) => setAsset(e.target.value)} />
        </Row>
      </Card>

      <Card>
        <SubTitle>Payouts</SubTitle>
        <Row>
          <Select value={status} onChange={(e) => setStatus(e.target.value)}>
            <option value="">All statuses</option>
            {Object.entries(STATUS_LABELS).map(([k, v]) => (
              <option key={k} value={k}>
                {v}
              </option>
            ))}
          </Select>
          <button onClick={exportPayouts} disabled={!report?.payouts?.length} style={{ marginLeft: "auto" }}>
            Export CSV
          </button>
        </Row>
        {(report?.totals ?? []).map((t) => (
          <Grid key={t.currency} style={{ marginBottom: 16 }}>
            <Stat label={`${t.currency} payouts`} value={t.payouts} />
            <Stat label="Authorized" value={amount(t.total, t.currency)} />
            <Stat label="To holders" value={amount(t.holderPayable, t.currency)} />
            <Stat label="Paid to holders" value={amount(t.paidToHolders, t.currency)} />
            <Stat label="Fees" value={amount(t.fees, t.currency)} />
            <Stat label="VAT" value={amount(t.vat, t.currency)} />
            <Stat label="Holders paid / failed" value={`${t.holdersPaid} / ${t.holdersFailed}`} />
          </Grid>
        ))}
        {report && Object.keys(report.byStatus ?? {}).length > 0 && (
          <Muted>
            {Object.entries(report.byStatus)
              .map(([s, n]) => `${STATUS_LABELS[s as keyof typeof STATUS_LABELS] ?? s}: ${n}`)
              .join(" · ")}
          </Muted>
        )}
        <CustomTable
          columns={payoutColumns}
          dataSource={report?.payouts ?? []}
          isLoading={loadingPayouts && !report}
          onRowClick={(r: IPayout) => router.push(`/dividendandyield/${r.id}`)}
        />
      </Card>

      <Card>
        <SubTitle>Payout fees and VAT</SubTitle>
        <Muted>Payout processing fees and the VAT on them, as paid on-chain to the fee and VAT wallets.</Muted>
        <Row>
          {Object.entries(fees?.totalFees ?? {}).map(([c, v]) => (
            <Stat key={`f${c}`} label={`Fees (${c})`} value={amount(v, c)} />
          ))}
          {Object.entries(fees?.totalVat ?? {}).map(([c, v]) => (
            <Stat key={`v${c}`} label={`VAT (${c})`} value={amount(v, c)} />
          ))}
          <button onClick={exportFees} disabled={!fees?.collections?.length} style={{ marginLeft: "auto" }}>
            Export CSV
          </button>
        </Row>
        {(fees?.totals ?? []).length > 0 && (
          <CustomTable
            columns={[
              { title: "Asset", dataIndex: "assetCode" },
              { title: "Currency", dataIndex: "payoutAssetCode" },
              { title: "Fees", dataIndex: "fees", render: (v: number) => amount(v) },
              { title: "VAT", dataIndex: "vat", render: (v: number) => amount(v) },
              { title: "Lines", dataIndex: "count" },
            ]}
            dataSource={fees?.totals ?? []}
          />
        )}
        <CustomTable columns={feeColumns} dataSource={fees?.collections ?? []} isLoading={loadingFees && !fees} />
      </Card>
    </>
  );
};

export default Reports;
