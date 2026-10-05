"use client";
import React, { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import { Modal } from "@/components/CustomModal";
import {
  IHolder,
  IReconRun,
  useGetPMAssetQuery,
  useGetPMPermissionsQuery,
  usePmAssetActionMutation,
  useRecordPMPositionMutation,
  useRegisterPMContractMutation,
  useRunPMReconciliationMutation,
  useSetPMPriceMutation,
} from "@/redux/api/publicMarkets";
import { act, Back, Button, Card, date, Field, FormGrid, Grid, Mono, Muted, ngn, Note, num, Pill, Row, Stat, SubTitle, Title } from "../../components/ui";

type Dialog = "" | "contract" | "price" | "position";

const AssetPage = () => {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const { data, isLoading } = useGetPMAssetQuery(id, { skip: !id });
  const { data: me } = useGetPMPermissionsQuery();
  const [assetAction, { isLoading: acting }] = usePmAssetActionMutation();
  const [reconcile] = useRunPMReconciliationMutation();
  const [dialog, setDialog] = useState<Dialog>("");
  const a = data?.data;
  const manage = !!me?.data?.manage;

  if (isLoading || !a) return <Card>{isLoading ? "Loading…" : "Asset not found."}</Card>;
  const drift = a.position !== "" && Number(a.supply) > Number(a.position);

  return (
    <div>
      <Back href="/publicmarkets/assets" label="Public Market Assets" />
      <Card>
        <Row style={{ justifyContent: "space-between", marginTop: 0 }}>
          <div>
            <Title>
              {a.assetCode} <Pill status={a.status} />
            </Title>
            <Muted>{a.instrumentName}</Muted>
          </div>
          {manage && (
            <Row style={{ margin: 0 }}>
              {a.status === "SETUP" && (
                <>
                  <Button onClick={() => setDialog("contract")}>Register token</Button>
                  <Button
                    $variant="primary"
                    disabled={acting}
                    onClick={() => act(() => assetAction({ id: a.id, action: "go-live" }).unwrap(), { confirm: `Open ${a.assetCode} for creation and redemption?` })}
                  >
                    Take live
                  </Button>
                </>
              )}
              {a.status === "LIVE" && (
                <Button
                  $variant="danger"
                  disabled={acting}
                  onClick={() =>
                    act((reason) => assetAction({ id: a.id, action: "halt", reason }).unwrap(), {
                      reason: `Why halt ${a.assetCode}? New orders are rejected; holders' balances are not frozen.`,
                    })
                  }
                >
                  Halt creation &amp; redemption
                </Button>
              )}
              {a.status === "HALTED" && (
                <Button
                  $variant="primary"
                  disabled={acting}
                  onClick={() =>
                    act(() => assetAction({ id: a.id, action: "resume" }).unwrap(), {
                      confirm: "Reconcile now and reopen only if supply, ledger and the Custodian position match?",
                    })
                  }
                >
                  Resume after resolution
                </Button>
              )}
              <Button onClick={() => setDialog("price")}>Set price</Button>
              <Button onClick={() => setDialog("position")}>Record position</Button>
              <Button onClick={() => act(() => reconcile({ target: a.assetCode }).unwrap(), { success: "Reconciliation requested" })}>Reconcile now</Button>
            </Row>
          )}
        </Row>
        {a.status === "HALTED" && (
          <Note $tone="warn">
            <b>Halted {a.haltedAt ? `since ${date(a.haltedAt)}` : ""}.</b> {a.haltReason}
          </Note>
        )}
        {a.status === "SETUP" && (
          <div style={{ marginTop: 18 }}>
            <SubTitle>Setup</SubTitle>
            {a.setupSteps.map((s) => (
              <div key={s.label} style={{ display: "flex", gap: 10, padding: "6px 0", alignItems: "center" }}>
                <span style={{ color: s.done ? "#00A859" : "#B26A00", fontWeight: 700 }}>{s.done ? "✓" : "○"}</span>
                <span style={{ fontWeight: 600 }}>{s.label}</span>
                <span style={{ color: "#667085", fontSize: 13 }}>{s.detail}</span>
              </div>
            ))}
          </div>
        )}
      </Card>

      <Card>
        <SubTitle>Asset information</SubTitle>
        <Grid>
          <Stat label="Market" value={`${a.market} · ${a.assetType}`} />
          <Stat label="ISIN" value={a.isin} />
          <Stat label="Unit" value={a.unitDescription || "—"} />
          <Stat label="Minimum buy" value={ngn(a.minimumBuy)} />
          <Stat label="Trade fee" value={a.feePercent ? `${a.feePercent}%` : "Settings' fee"} />
          <Stat label="Inventory target" value={`${num(a.inventoryTargetUnits)} units`} />
          <Stat label="Token contract" value={a.contractAddress ? <Mono>{a.contractAddress}</Mono> : "Not registered"} />
          <Stat label="Issuing Safe" value={a.issuingSafeAddress ? <Mono>{a.issuingSafeAddress}</Mono> : "—"} />
          <Stat label="Token decimals" value={a.contractAddress ? a.tokenDecimals : "—"} />
        </Grid>
      </Card>

      <Card>
        <SubTitle>Custody &amp; execution</SubTitle>
        <Grid>
          <Stat label="Custodian" value={a.custodianName || "—"} />
          <Stat label="Nominee" value={a.custodian?.nomineeName || "—"} />
          <Stat label="Omnibus account" value={a.omnibusReference || "—"} />
          <Stat label="Custodian integration" value={a.custodian ? `${a.custodian.mode}${a.custodian.transport ? ` · ${a.custodian.transport}` : ""}` : "—"} />
          <Stat label="Dealing Member" value={a.dealingMemberName || "—"} />
          <Stat label="CSCS member code" value={a.dealingMember?.cscsMemberCode || "—"} />
          <Stat label="Dealing Member integration" value={a.dealingMember?.mode || "—"} />
          <Stat label="Settlement" value={a.market === "NGX" ? "T+2 · confirmed by the Custodian" : "Confirmed by the Custodian"} />
        </Grid>
      </Card>

      <Card>
        <SubTitle>Supply &amp; reconciliation</SubTitle>
        <Grid>
          <Stat label="Token supply (Σ ledger)" value={num(a.supply)} />
          <Stat
            label="Custodian position"
            value={
              <span style={{ color: drift ? "#BE3800" : undefined }}>
                {a.position ? `${num(a.position)} · ${a.positionSource} · ${date(a.positionAsOf ?? undefined)}` : "None recorded"}
              </span>
            }
          />
          <Stat label="Beneficial owners" value={`${num(a.owners)} wallets`} />
          <Stat label="Last price" value={Number(a.lastPrice) > 0 ? `${ngn(a.lastPrice)} · ${a.priceSource}` : "—"} />
          <Stat label="Market value" value={ngn(a.marketValue)} />
          <Stat label="Open orders" value={a.openOrders} />
          <Stat label="Last reconciliation" value={a.lastReconResult ? <Pill status={a.lastReconResult} /> : "—"} />
        </Grid>
        <div style={{ marginTop: 18 }}>
          <CustomTable
            columns={[
              { title: "Run", dataIndex: "createdAt", render: (v: string) => date(v) },
              { title: "Supply", dataIndex: "tokenSupply", render: (v: string) => num(v) },
              { title: "Position", dataIndex: "custodianPosition", render: (v: string) => num(v) },
              { title: "Inventory", dataIndex: "inventory", render: (v: string) => num(v) },
              { title: "Σ ledger", dataIndex: "ledgerTotal", render: (v: string) => num(v) },
              { title: "Result", dataIndex: "result", render: (v: string, r: IReconRun) => <span title={r.detail}><Pill status={v} /></span> },
            ]}
            dataSource={a.reconciliationRuns}
          />
        </div>
      </Card>

      <Card>
        <SubTitle>Beneficial owners</SubTitle>
        <CustomTable
          columns={[
            { title: "Wallet", dataIndex: "walletAddress", render: (v: string) => <Mono>{v}</Mono> },
            { title: "Channel", dataIndex: "channel", render: (v: string, h: IHolder) => (v === "EXCHANGE" ? h.exchangeName || "Exchange" : v === "TROVO_APP" ? "Trovo App" : "Platform") },
            { title: "Balance", dataIndex: "balance", render: (v: string) => num(v) },
            { title: "% of supply", dataIndex: "percentOfSupply", render: (v: string) => `${v}%` },
            { title: "", dataIndex: "substantial", render: (v: boolean) => (v ? <Pill status="NEEDS_MANUAL" label="Above threshold: CAC disclosure due" /> : "") },
          ]}
          dataSource={a.holders}
        />
        <Muted style={{ marginTop: 10 }}>
          Exchange-provisioned wallets belong to individual customers, so they appear here like Trovo App wallets. Platform wallets hold working inventory.
        </Muted>
      </Card>

      <Card>
        <Row style={{ justifyContent: "space-between", marginTop: 0 }}>
          <SubTitle>Corporate actions</SubTitle>
          <Button onClick={() => router.push(`/publicmarkets/corporate-actions?asset=${a.assetCode}`)}>All corporate actions</Button>
        </Row>
        <CustomTable
          columns={[
            { title: "Event", dataIndex: "eventType" },
            { title: "Record date", dataIndex: "recordDate" },
            { title: "Pay date", dataIndex: "payDate" },
            { title: "Per unit", dataIndex: "amountPerUnit", render: (v: string) => (Number(v) > 0 ? ngn(v) : "—") },
            { title: "Status", dataIndex: "status", render: (v: string) => <Pill status={v} /> },
          ]}
          dataSource={a.corporateActions}
          onRowClick={(c: { id: string }) => router.push(`/publicmarkets/corporate-actions/${c.id}`)}
        />
      </Card>

      {dialog && <AssetDialog kind={dialog} id={a.id} code={a.assetCode} onClose={() => setDialog("")} />}
    </div>
  );
};

const AssetDialog = ({ kind, id, code, onClose }: { kind: Exclude<Dialog, "">; id: string; code: string; onClose: () => void }) => {
  const [register, r1] = useRegisterPMContractMutation();
  const [setPrice, r2] = useSetPMPriceMutation();
  const [record, r3] = useRecordPMPositionMutation();
  const [v, setV] = useState<Record<string, string>>({});
  const bind = (k: string) => ({ value: v[k] ?? "", onChange: (e: React.ChangeEvent<HTMLInputElement>) => setV({ ...v, [k]: e.target.value }) });
  const busy = r1.isLoading || r2.isLoading || r3.isLoading;

  const submit = async () => {
    let res;
    if (kind === "contract") {
      res = await act(() => register({ id, contractAddress: v.contract ?? "", issuingSafeAddress: v.safe ?? "" }).unwrap(), { success: "Token verified and registered" });
    } else if (kind === "price") {
      res = await act(() => setPrice({ id, price: v.price ?? "" }).unwrap(), { success: "Price recorded" });
    } else {
      res = await act(() => record({ id, unitsHeld: v.units ?? "", asOf: v.asOf ?? "", reference: v.reference ?? "" }).unwrap(), { success: "Position recorded" });
    }
    if (res) onClose();
  };

  const titles = { contract: `Register ${code}'s token`, price: `Manual price for ${code}`, position: `Custodian position for ${code}` };
  return (
    <Modal
      title={titles[kind]}
      isOpen
      onClose={onClose}
      style={{ maxWidth: 560, width: "95%" }}
      actions={
        <Button $variant="primary" disabled={busy} onClick={submit}>
          Save
        </Button>
      }
    >
      {kind === "contract" && (
        <FormGrid>
          <Field label="Token contract" hint="Deployed TokenizedAsset ERC-20, nothing minted yet">
            <input {...bind("contract")} placeholder="0x…" />
          </Field>
          <Field label="Issuing Safe" hint="Must be the token's owner; signed by the Public Markets signers">
            <input {...bind("safe")} placeholder="0x…" />
          </Field>
        </FormGrid>
      )}
      {kind === "price" && (
        <Field label="Price (NGN)" hint="Recorded with source MANUAL. Use it when no feed is connected or to correct one.">
          <input {...bind("price")} placeholder="221.50" />
        </Field>
      )}
      {kind === "position" && (
        <FormGrid>
          <Field label="Units held" hint="Whole units, from the Custodian's statement">
            <input {...bind("units")} />
          </Field>
          <Field label="As of" hint="YYYY-MM-DD">
            <input {...bind("asOf")} placeholder="2026-09-29" />
          </Field>
          <Field label="Statement reference">
            <input {...bind("reference")} placeholder="POS_20260929.csv" />
          </Field>
        </FormGrid>
      )}
    </Modal>
  );
};

export default AssetPage;
