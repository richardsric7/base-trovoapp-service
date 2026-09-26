"use client";
import React, { useMemo, useState } from "react";
import styled from "styled-components";
import CustomTable from "@/components/CustomTable";
import { IMerchantLeaderboardEntry, useP2pMerchantLeaderboardQuery } from "@/redux/api/p2p";
import ReportSectionCard from "./ReportSectionCard";
import { formatAmount } from "./reportUtils";

type SortBy = "completedTrades" | "completedVolume" | "completionRate" | "disputesOpened";

const SORT_OPTIONS: { label: string; value: SortBy }[] = [
  { label: "Completed Trades", value: "completedTrades" },
  { label: "Trade Volume", value: "completedVolume" },
  { label: "Completion Rate", value: "completionRate" },
  { label: "Disputes Opened", value: "disputesOpened" },
];

const PAGE_SIZE = 10;

const MerchantLeaderboardSection = () => {
  const [sortBy, setSortBy] = useState<SortBy>("completedTrades");
  const [page, setPage] = useState(1);
  const { data, isLoading } = useP2pMerchantLeaderboardQuery({ page, pageSize: PAGE_SIZE, sortBy });

  const columns = useMemo(
    () => [
      {
        title: "Merchant",
        dataIndex: "merchantUsername",
        render: (_: any, record: IMerchantLeaderboardEntry) => <Bold>{record.merchantUsername || "N/A"}</Bold>,
      },
      { title: "Completed Trades", dataIndex: "completedTrades" },
      {
        title: "Trade Volume",
        dataIndex: "completedTradeVolume",
        render: (v: string) => formatAmount(v),
      },
      {
        title: "Completion Rate",
        dataIndex: "completionRate",
        render: (v: string) => `${Number(v).toFixed(1)}%`,
      },
      { title: "Cancelled", dataIndex: "cancelledOrders" },
      { title: "Expired", dataIndex: "expiredOrders" },
      { title: "Disputes Opened", dataIndex: "disputesOpened" },
      {
        title: "Disputes Lost",
        dataIndex: "disputesResolvedAgainstMerchant",
      },
    ],
    [],
  );

  return (
    <ReportSectionCard
      title="Merchant Leaderboard"
      subtitle="Full merchant performance report, sortable"
      action={
        <SortSelect value={sortBy} onChange={(e) => { setSortBy(e.target.value as SortBy); setPage(1); }}>
          {SORT_OPTIONS.map((o) => (
            <option key={o.value} value={o.value}>
              Sort by: {o.label}
            </option>
          ))}
        </SortSelect>
      }
    >
      <CustomTable
        columns={columns}
        dataSource={data?.data?.data ?? []}
        totalItems={data?.data?.total ?? 0}
        pageSize={PAGE_SIZE}
        isLoading={isLoading}
        onPageChange={setPage}
      />
    </ReportSectionCard>
  );
};

export default MerchantLeaderboardSection;

const Bold = styled.span`
  font-weight: 600;
  color: #00225a;
`;

const SortSelect = styled.select`
  padding: 8px 12px;
  border-radius: 8px;
  border: 1px solid #e5e5ef;
  background-color: #f2f6f9;
  color: #00225a;
  font-size: 13px;
`;
