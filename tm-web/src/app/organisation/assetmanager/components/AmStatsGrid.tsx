import React from "react";
import styled from "styled-components";
// import StatCard from "./StatsCard";
import { AmSummary } from "@/redux/api/assetManager";
import StatCard from "../../components/StatsCard";

interface AmStatsGridProps {
  summary?: AmSummary;
}

const formatAmount = (value?: number | string, currency = "NGN") => {
  const amount = Number(value ?? 0);
  return Number.isFinite(amount)
    ? new Intl.NumberFormat("en-NG", { style: "currency", currency }).format(
        amount,
      )
    : String(value ?? 0);
};

const AmStatsGrid = ({ summary }: AmStatsGridProps) => (
  <Grid>
    <StatCard label="Assets Under Management" value={summary?.assets_managed ?? 0} />
    <StatCard label="Pending Trustee Approvals" value={summary?.pending_trustee_approvals ?? 0} />

  
    <StatCard label="YTD Revenue" value={formatAmount(summary?.ytd_revenue)} />
  </Grid>
);

export default AmStatsGrid;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
`;
