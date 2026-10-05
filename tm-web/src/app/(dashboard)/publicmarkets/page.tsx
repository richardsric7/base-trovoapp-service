"use client";
import React from "react";
import { useRouter } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import { IAttentionItem, INetBatch, useGetPMOverviewQuery } from "@/redux/api/publicMarkets";
import { Card, date, LinkButton, Muted, ngn, Note, num, Panel, Pill, SectionNav, StatCard, Stats, SubTitle, Title, TwoCols } from "./components/ui";

// Where each attention item is handled.
const LINKS: Record<string, (t: string) => string> = {
  approval: () => "/publicmarkets/orders?tab=approval",
  escalation: () => "/publicmarkets/orders?tab=escalations",
  "dead-letter": (t) => `/publicmarkets/exchanges/${encodeURIComponent(t)}`,
  confirmation: () => "/publicmarkets/corporate-actions?tab=confirmations",
  disclosure: (t) => `/publicmarkets/assets/${encodeURIComponent(t)}`,
  halt: (t) => `/publicmarkets/assets/${encodeURIComponent(t)}`,
  drift: () => "/publicmarkets/reconciliation",
};

const PublicMarketsOverview = () => {
  const router = useRouter();
  const { data, isLoading } = useGetPMOverviewQuery(undefined, { pollingInterval: 30000 });
  const o = data?.data;

  return (
    <div>
      <SectionNav />
      <Stats>
        <StatCard title="Public Market Assets" value={o?.assets ?? "—"} sub={o ? `${o.liveAssets} live · ${ngn(o.marketValue)}` : undefined} />
        <StatCard title="Settled Orders (Today)" value={o?.settledToday ?? "—"} sub={o ? ngn(o.settledTodayValue) : undefined} color="#00A859" bg="#00A8591A" />
        <StatCard title="Orders In Progress" value={o?.inProgress ?? "—"} sub={o ? ngn(o.inProgressValue) : undefined} color="#007CDF" />
        <StatCard title="Halted Assets" value={o?.haltedAssets ?? "—"} color="#BE3800" bg="#BE38001A" />
      </Stats>

      {o?.halted?.map((a) => (
        <Note key={a.id} $tone="warn" style={{ marginTop: 0, marginBottom: 16 }}>
          <b>{a.assetCode} is halted.</b> {a.haltReason} Creation and redemption are stopped; holders&apos; balances are not frozen.{" "}
          <LinkButton onClick={() => router.push("/publicmarkets/reconciliation")}>Open reconciliation →</LinkButton>
        </Note>
      ))}

      <TwoCols>
        <Panel>
          <SubTitle>Needs attention</SubTitle>
          <CustomTable
            columns={[
              { title: "Item", dataIndex: "text" },
              { title: "Where", dataIndex: "where" },
            ]}
            dataSource={o?.attention ?? []}
            isLoading={isLoading}
            onRowClick={(r: IAttentionItem) => router.push((LINKS[r.kind] ?? (() => "/publicmarkets"))(r.target))}
          />
        </Panel>
        <Panel>
          <SubTitle>Background jobs</SubTitle>
          {(o?.jobs ?? []).map((j) => (
            <div key={j.job} style={{ display: "flex", justifyContent: "space-between", padding: "10px 0", borderTop: "1px solid #eef2f6" }}>
              <div>
                <div style={{ fontSize: 13, fontWeight: 600 }}>{j.job}</div>
                <div style={{ fontSize: 12, color: "#667085" }}>
                  {date(j.lastRunAt)} {j.detail && `· ${j.detail}`}
                </div>
              </div>
              <Pill status={j.status} />
            </div>
          ))}
          {!o?.jobs?.length && <Muted>The engine has not reported yet.</Muted>}
        </Panel>
      </TwoCols>

      <Card>
        <Title>Creation vs redemption (today)</Title>
        <Muted>The fast path fills from Custodian inventory; the slow path waits for the session&apos;s net batch to the Dealing Member.</Muted>
        <Stats style={{ marginTop: 16 }}>
          <StatCard title="Creation orders" value={num(o?.creationsToday)} sub={o ? `Fast path ${o.fastPercent}%` : undefined} />
          <StatCard title="Redemption orders" value={num(o?.redemptionsToday)} sub={o ? `Netted ${o.nettedPercent}%` : undefined} />
          <StatCard title="Net batches today" value={num(o?.batchesToday?.length ?? 0)} />
        </Stats>
        <CustomTable
          columns={[
            { title: "Batch", dataIndex: "id" },
            { title: "Side", dataIndex: "side" },
            { title: "Quantity", dataIndex: "quantity", render: (v: string) => num(v) },
            { title: "Value", dataIndex: "value", render: (v: string) => ngn(v) },
            { title: "Orders", dataIndex: "creationOrders", render: (_: unknown, b: INetBatch) => `${b.creationOrders} buy · ${b.redemptionOrders} sell` },
            { title: "Status", dataIndex: "status", render: (v: string) => <Pill status={v} /> },
          ]}
          dataSource={o?.batchesToday ?? []}
          onRowClick={(b: INetBatch) => router.push(`/publicmarkets/orders?batch=${encodeURIComponent(b.id)}`)}
        />
      </Card>
    </div>
  );
};

export default PublicMarketsOverview;
