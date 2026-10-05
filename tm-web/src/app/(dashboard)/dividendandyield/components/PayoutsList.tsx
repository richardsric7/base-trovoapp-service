"use client";
import React, { useState } from "react";
import { useRouter } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import { IPayout, useGetPayoutsQuery } from "@/redux/api/proceedPayouts";
import { amount, Card, date, Input, Muted, Row, Select, STATUS_LABELS, StatusPill, Title } from "./ui";

// Payouts of authorized stakeholder distributions. A trustee authorizing a
// distribution registers its payout; admins then drive it from its page.
const PayoutsList = ({ asset: initialAsset }: { asset?: string }) => {
  const router = useRouter();
  const [status, setStatus] = useState("");
  const [asset, setAsset] = useState(initialAsset ?? "");
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const { data, isFetching } = useGetPayoutsQuery({ status, asset: asset.trim(), page, limit });
  const rows = data?.data?.payouts ?? [];

  const columns = [
    { title: "Registered", dataIndex: "createdAt", render: (v: string) => date(v) },
    {
      title: "Asset",
      dataIndex: "assetCode",
      render: (_: any, r: IPayout) => (
        <span>
          <strong>{r.assetCode || r.tokenizedAssetId}</strong>
          <br />
          <small>{r.assetName}</small>
        </span>
      ),
    },
    { title: "Batch", dataIndex: "batch" },
    { title: "Amount", dataIndex: "totalAmount", render: (v: string, r: IPayout) => amount(v, r.payoutAssetCode) },
    { title: "Per token", dataIndex: "amountPerToken", render: (v: string) => (v ? amount(v) : "—") },
    {
      title: "Holders",
      dataIndex: "holderCount",
      render: (_: any, r: IPayout) => (r.holderCount ? `${r.paidCount} / ${r.holderCount - r.excludedCount} paid` : "—"),
    },
    {
      title: "Approvals",
      dataIndex: "approvals",
      render: (_: any, r: IPayout) => (r.status === "LOCKED" ? `${r.approvals} / ${r.approvalsRequired}` : "—"),
    },
    { title: "Status", dataIndex: "status", render: (v: string) => <StatusPill status={v} /> },
  ];

  return (
    <Card>
      <Title>Proceeds payouts</Title>
      <Muted>
        A trustee authorizing a stakeholder distribution registers its payout here. Prepare the schedule (the
        engine snapshots the token holders), have it approved by other admins, fund the payout Safe and confirm
        funding: payout-engine then pays the holders in batches and notifies them.
      </Muted>
      <Row>
        <Select
          value={status}
          onChange={(e) => {
            setStatus(e.target.value);
            setPage(1);
          }}
        >
          <option value="">All statuses</option>
          {Object.entries(STATUS_LABELS).map(([k, v]) => (
            <option key={k} value={k}>
              {v}
            </option>
          ))}
        </Select>
        <Input
          placeholder="Asset code or ID"
          value={asset}
          onChange={(e) => {
            setAsset(e.target.value);
            setPage(1);
          }}
        />
      </Row>
      <CustomTable
        columns={columns}
        dataSource={rows}
        isLoading={isFetching && !data}
        onRowClick={(r: IPayout) => router.push(`/dividendandyield/${r.id}`)}
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
    </Card>
  );
};

export default PayoutsList;
