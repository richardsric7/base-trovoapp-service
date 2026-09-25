import React from "react";
import styled from "styled-components";
import StatCard from "./StatsCard";
import {
  AssetManagerDashboardMetric,
  AssetManagerDashboardSummary,
} from "@/redux/api/assetManager";

interface StatsGridProps {
  summary?: AssetManagerDashboardSummary;
  currency?: string;
}

const getMetric = (metric?: AssetManagerDashboardMetric | number) =>
  typeof metric === "number" ? metric : metric?.total ?? metric?.value ?? 0;

const getPending = (metric?: AssetManagerDashboardMetric | number) =>
  typeof metric === "object" ? metric.pending ?? 0 : 0;

const formatAmount = (value?: number | string, currency = "NGN") => {
  const amount = Number(value ?? 0);
  return Number.isFinite(amount)
    ? new Intl.NumberFormat("en-NG", { style: "currency", currency }).format(
        amount,
      )
    : String(value ?? 0);
};

const StatsGrid = ({ summary, currency }: StatsGridProps) => (
  <Grid>
    <StatCard label="Total Asset Count" value={summary?.total_assets ?? 0} />
    <StatCard label="Asset Under Management" value={formatAmount(summary?.assets_under_management, currency)} />
    <StatCard label="Fees Generated" value={formatAmount(summary?.fees_generated, currency)} />
    <StatCard label="Milestone Verifications" value={getMetric(summary?.milestone_verifications)} pending={getPending(summary?.milestone_verifications)} showFooter />
    <StatCard label="Fund Release Requests" value={getMetric(summary?.fund_release_requests)} pending={getPending(summary?.fund_release_requests)} showFooter />
    <StatCard label="Distribution Requests" value={getMetric(summary?.distribution_requests)} pending={getPending(summary?.distribution_requests)} showFooter />
  </Grid>
);

export default StatsGrid;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
`;
