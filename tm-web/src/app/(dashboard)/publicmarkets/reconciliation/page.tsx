"use client";
import React from "react";
import { useRouter } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import {
  IReconRow,
  useGetPMCustodiansQuery,
  useGetPMJobsQuery,
  useGetPMPermissionsQuery,
  useGetPMReconciliationQuery,
  useRequestPMPositionFeedMutation,
  useRunPMReconciliationMutation,
} from "@/redux/api/publicMarkets";
import { act, Button, Card, date, Grid, Muted, Note, num, Pill, Row, SectionNav, Stat, StatCard, Stats, SubTitle, Title } from "../components/ui";

const ReconciliationPage = () => {
  const router = useRouter();
  const { data, isFetching } = useGetPMReconciliationQuery(undefined, { pollingInterval: 30000 });
  const { data: jobs } = useGetPMJobsQuery(undefined, { pollingInterval: 15000 });
  const { data: me } = useGetPMPermissionsQuery();
  const { data: custodians } = useGetPMCustodiansQuery();
  const [run, { isLoading }] = useRunPMReconciliationMutation();
  const [feed] = useRequestPMPositionFeedMutation();
  const r = data?.data;
  const drifts = (r?.rows ?? []).filter((x) => x.result !== "MATCHED");
  const mocks = (custodians?.data ?? []).filter((c) => c.integration?.mode === "MOCK");

  return (
    <div>
      <SectionNav />
      <Stats>
        <StatCard title="Assets Checked" value={r?.checked ?? "—"} />
        <StatCard title="Matched" value={r?.matched ?? "—"} color="#00A859" bg="#00A8591A" />
        <StatCard title="Drift Detected" value={r?.drift ?? "—"} color="#BE3800" bg="#BE38001A" />
        <StatCard title="Stale Prices" value={r?.stalePrices ?? "—"} color="#B26A00" bg="#B26A001A" />
      </Stats>
      <Card>
        <Row style={{ justifyContent: "space-between", marginTop: 0 }}>
          <div>
            <Title>Reconciliation</Title>
            <Muted>
              Daily, fail-closed: token supply must not exceed the Custodian&apos;s position, the ownership ledger must equal supply, and a position must exist. Any
              drift halts the asset; it resumes only when a fresh run matches. Last run {date(r?.lastRunAt ?? undefined)} · next run {date(r?.nextRunAt)}.
            </Muted>
          </div>
          {me?.data?.manage && (
            <Row style={{ margin: 0 }}>
              {mocks.map((c) => (
                <Button key={c.id} onClick={() => act(() => feed({ target: c.integration!.code }).unwrap(), { success: `Position feed requested from ${c.name}` })}>
                  Mock feed: {c.integration!.code}
                </Button>
              ))}
              <Button $variant="primary" disabled={isLoading} onClick={() => act(() => run({ target: "ALL" }).unwrap(), { confirm: "Reconcile every asset now?" })}>
                Run now
              </Button>
            </Row>
          )}
        </Row>
        <div style={{ height: 16 }} />
        <CustomTable
          columns={[
            { title: "Asset", dataIndex: "assetCode", render: (v: string) => <strong>{v}</strong> },
            { title: "Custodian", dataIndex: "custodianName" },
            { title: "Token supply", dataIndex: "tokenSupply", render: (v: string) => num(v) },
            { title: "Custodian position", dataIndex: "custodianPosition", render: (v: string, x: IReconRow) => `${num(v)} (${x.positionSource || "—"})` },
            {
              title: "Δ",
              dataIndex: "delta",
              render: (v: string) => <span style={{ fontWeight: 700, color: Number(v) > 0 ? "#BE3800" : "#00A859" }}>{Number(v) > 0 ? `+${num(v)}` : "0"}</span>,
            },
            { title: "Inventory", dataIndex: "inventory", render: (v: string) => num(v) },
            { title: "Σ ledger", dataIndex: "ledgerTotal", render: (v: string, x: IReconRow) => `${num(v)}${Number(x.ledgerDelta) ? ` (Δ ${x.ledgerDelta})` : ""}` },
            { title: "Price", dataIndex: "priceStale", render: (v: boolean) => <Pill status={v ? "Stale" : "Fresh"} /> },
            { title: "Result", dataIndex: "result", render: (v: string, x: IReconRow) => <span title={x.detail}><Pill status={v} /></span> },
            { title: "Checked", dataIndex: "createdAt", render: (v: string) => date(v) },
          ]}
          dataSource={r?.rows ?? []}
          isLoading={isFetching && !data}
          onRowClick={(x: IReconRow) => router.push(`/publicmarkets/assets/${x.assetId}`)}
        />
      </Card>

      {drifts.map((x) => (
        <Card key={x.id}>
          <SubTitle>
            {x.assetCode}: <Pill status={x.result} /> — diagnosis
          </SubTitle>
          <Grid>
            <Stat label="Detail" value={x.detail || "—"} />
            <Stat label="Position as of" value={`${date(x.positionAsOf ?? undefined)} · ${x.positionSource || "—"}`} />
            <Stat label="Asset" value={<Pill status={x.assetStatus} />} />
          </Grid>
          {x.ordersSince.length > 0 && (
            <>
              <Note>Orders settled after the Custodian&apos;s position was taken; the next statement should include them.</Note>
              <CustomTable
                columns={[
                  { title: "Order", dataIndex: "id" },
                  { title: "Type", dataIndex: "type" },
                  { title: "Quantity", dataIndex: "quantity", render: (v: string) => num(v) },
                  { title: "Settled", dataIndex: "settledAt", render: (v: string) => date(v) },
                ]}
                dataSource={x.ordersSince}
              />
            </>
          )}
          {me?.data?.manage && (
            <Row>
              <Button onClick={() => router.push(`/publicmarkets/assets/${x.assetId}`)}>Record the Custodian&apos;s position</Button>
              <Button $variant="primary" onClick={() => act(() => run({ target: x.assetCode }).unwrap(), { confirm: `Re-run reconciliation for ${x.assetCode}?` })}>
                Re-run for {x.assetCode}
              </Button>
            </Row>
          )}
        </Card>
      ))}

      <Card>
        <SubTitle>Requests to the engine</SubTitle>
        <CustomTable
          columns={[
            { title: "Requested", dataIndex: "createdAt", render: (v: string) => date(v) },
            { title: "Job", dataIndex: "job" },
            { title: "Target", dataIndex: "target" },
            { title: "By", dataIndex: "requestedBy" },
            { title: "Result", dataIndex: "result", render: (v: string, j: { doneAt?: string | null }) => (j.doneAt ? v || "done" : "Waiting for the engine…") },
          ]}
          dataSource={(jobs?.data ?? []).slice(0, 20)}
        />
      </Card>
    </div>
  );
};

export default ReconciliationPage;
