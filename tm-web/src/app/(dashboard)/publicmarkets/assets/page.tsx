"use client";
import React, { useState } from "react";
import { useRouter } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import { Modal } from "@/components/CustomModal";
import {
  AssetInput,
  IAssetRow,
  useCreatePMAssetMutation,
  useGetPMAssetsQuery,
  useGetPMCustodiansQuery,
  useGetPMDealingMembersQuery,
  useGetPMPermissionsQuery,
} from "@/redux/api/publicMarkets";
import { act, Button, Card, Field, FormGrid, Input, Muted, ngn, num, Pill, Row, SectionNav, Select, Title } from "../components/ui";

const AssetsPage = () => {
  const router = useRouter();
  const [filters, setFilters] = useState({ market: "", type: "", status: "", search: "" });
  const { data, isFetching } = useGetPMAssetsQuery(filters);
  const { data: me } = useGetPMPermissionsQuery();
  const [open, setOpen] = useState(false);
  const rows = data?.data ?? [];
  const set = (k: keyof typeof filters) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => setFilters({ ...filters, [k]: e.target.value });

  const columns = [
    {
      title: "Asset",
      dataIndex: "assetCode",
      render: (_: unknown, a: IAssetRow) => (
        <span>
          <strong>{a.assetCode}</strong>
          <br />
          <small>{a.instrumentName}</small>
        </span>
      ),
    },
    { title: "Market", dataIndex: "market", render: (_: unknown, a: IAssetRow) => `${a.market} · ${a.assetType}` },
    { title: "ISIN", dataIndex: "isin" },
    { title: "Custodian", dataIndex: "custodianName" },
    { title: "Dealing Member", dataIndex: "dealingMemberName" },
    { title: "Token supply", dataIndex: "supply", render: (v: string) => num(v) },
    {
      title: "Custodian position",
      dataIndex: "position",
      render: (v: string, a: IAssetRow) => {
        const drift = v && Number(a.supply) > Number(v);
        return <span style={{ color: drift ? "#BE3800" : undefined, fontWeight: drift ? 700 : 400 }}>{v ? num(v) : "—"}</span>;
      },
    },
    { title: "Price", dataIndex: "lastPrice", render: (v: string, a: IAssetRow) => (Number(v) > 0 ? `${ngn(v)}${a.priceFresh ? "" : " (stale)"}` : "—") },
    { title: "Status", dataIndex: "status", render: (v: string) => <Pill status={v} /> },
  ];

  return (
    <div>
      <SectionNav />
      <Card>
        <Row style={{ justifyContent: "space-between", marginTop: 0 }}>
          <Title>Public Market Assets</Title>
          {me?.data?.manage && (
            <Button $variant="primary" onClick={() => setOpen(true)}>
              Add Public Market Asset
            </Button>
          )}
        </Row>
        <Row>
          <Input $wide placeholder="Search by code, ISIN or name" value={filters.search} onChange={set("search")} />
          <Select value={filters.market} onChange={set("market")}>
            <option value="">All markets</option>
            <option value="NGX">NGX</option>
            <option value="FMDQ">FMDQ</option>
          </Select>
          <Select value={filters.type} onChange={set("type")}>
            <option value="">All classes</option>
            <option value="EQUITY">Equity</option>
            <option value="BOND">Bond</option>
          </Select>
          <Select value={filters.status} onChange={set("status")}>
            <option value="">All statuses</option>
            <option value="LIVE">Live</option>
            <option value="HALTED">Halted</option>
            <option value="SETUP">Setting up</option>
          </Select>
        </Row>
        <CustomTable columns={columns} dataSource={rows} isLoading={isFetching && !data} onRowClick={(a: IAssetRow) => router.push(`/publicmarkets/assets/${a.id}`)} />
        <Muted style={{ marginTop: 14 }}>
          Public Market Assets are tokenized NGX equities and FMDQ bonds, each held by one Custodian in an omnibus account at CSCS and bought or sold by one Dealing
          Member. They are separate from the Tokenized Assets statistics.
        </Muted>
      </Card>
      {open && <AddAsset onClose={() => setOpen(false)} onCreated={(id) => router.push(`/publicmarkets/assets/${id}`)} />}
    </div>
  );
};

const AddAsset = ({ onClose, onCreated }: { onClose: () => void; onCreated: (id: string) => void }) => {
  const { data: custodians } = useGetPMCustodiansQuery();
  const { data: dms } = useGetPMDealingMembersQuery();
  const [create, { isLoading }] = useCreatePMAssetMutation();
  const [f, setF] = useState<AssetInput>({ market: "NGX", assetType: "EQUITY" });
  const set = (k: keyof AssetInput) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) =>
    setF({ ...f, [k]: k === "custodianId" || k === "dealingMemberId" ? Number(e.target.value) : e.target.value });
  const configured = (custodians?.data ?? []).filter((c) => c.integration?.active);
  const activeDMs = (dms?.data ?? []).filter((d) => d.active);

  const submit = async () => {
    const res = await act(() => create(f).unwrap(), { success: "Asset created in setup" });
    if (res?.data) onCreated(res.data.id);
  };

  return (
    <Modal
      title="Add Public Market Asset"
      isOpen
      onClose={onClose}
      style={{ maxWidth: 760, width: "95%" }}
      actions={
        <Button $variant="primary" disabled={isLoading} onClick={submit}>
          Create asset
        </Button>
      }
    >
      <FormGrid>
        <Field label="Asset code" hint="The token's code, e.g. DANGCEM-T">
          <input value={f.assetCode ?? ""} onChange={set("assetCode")} placeholder="DANGCEM-T" />
        </Field>
        <Field label="Ticker" hint="Defaults to the code without -T">
          <input value={f.ticker ?? ""} onChange={set("ticker")} placeholder="DANGCEM" />
        </Field>
        <Field label="Market">
          <select value={f.market} onChange={set("market")}>
            <option value="NGX">NGX</option>
            <option value="FMDQ">FMDQ</option>
          </select>
        </Field>
        <Field label="Asset class">
          <select value={f.assetType} onChange={set("assetType")}>
            <option value="EQUITY">Equity</option>
            <option value="BOND">Bond / T-Bill</option>
          </select>
        </Field>
        <Field label="Instrument name">
          <input value={f.instrumentName ?? ""} onChange={set("instrumentName")} placeholder="Dangote Cement Plc" />
        </Field>
        <Field label="Short name">
          <input value={f.shortName ?? ""} onChange={set("shortName")} placeholder="Dangote Cement" />
        </Field>
        <Field label="ISIN" hint="12 characters">
          <input value={f.isin ?? ""} onChange={set("isin")} placeholder="NGDANGCEM008" />
        </Field>
        <Field label="Sector">
          <input value={f.sector ?? ""} onChange={set("sector")} />
        </Field>
        <Field label="Custodian" hint="Configured in Settings › Custodians">
          <select value={f.custodianId ?? ""} onChange={set("custodianId")}>
            <option value="">Choose…</option>
            {configured.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name} ({c.integration?.code})
              </option>
            ))}
          </select>
        </Field>
        <Field label="Dealing Member">
          <select value={f.dealingMemberId ?? ""} onChange={set("dealingMemberId")}>
            <option value="">Choose…</option>
            {activeDMs.map((d) => (
              <option key={d.id} value={d.id}>
                {d.dealingMemberName}
              </option>
            ))}
          </select>
        </Field>
        <Field label="Omnibus account reference" hint="Defaults to POOL-<TICKER>-01">
          <input value={f.omnibusReference ?? ""} onChange={set("omnibusReference")} />
        </Field>
        <Field label="Minimum buy (NGN)">
          <input value={f.minimumBuy ?? ""} onChange={set("minimumBuy")} placeholder="1000" />
        </Field>
        <Field label="Inventory target (units)" hint="Bought on top of a net shortfall so later buys fill instantly">
          <input value={f.inventoryTargetUnits ?? ""} onChange={set("inventoryTargetUnits")} placeholder="0" />
        </Field>
        <Field label="Trade fee (%)" hint="Empty: the settings' fee">
          <input value={f.feePercent ?? ""} onChange={set("feePercent")} />
        </Field>
      </FormGrid>
    </Modal>
  );
};

export default AssetsPage;
