import React from "react";
import styled from "styled-components";
import StatCard from "../../components/StatsCard";
import { TrusteeSummary } from "@/redux/api/trustees";

interface TStatsGridProps {
  summary?: TrusteeSummary;
}

const TStatsGrid: React.FC<TStatsGridProps> = ({ summary }) => (
  <Grid>
    <StatCard
      label="Total Asset Count"
      value={summary?.assets_under_trust ?? 0}
    />
    <StatCard
      label="Assets Under Trust"
      value={summary?.assets_under_trust ?? 0}
    />
     <StatCard
      label="Fees Generated"
      value={ 0}
    />
    <StatCard
      label="Pending Due Diligence"
      value={summary?.pending_due_diligence ?? 0}
    />
     <StatCard label="Fund Release Approvals" value={summary?.pending_fund_release_approvals ?? 0} pending={summary?.pending_fund_release_approvals ?? 0} showFooter />
    <StatCard label="Distribution Requests" value={summary?.pending_distribution_authorizations ?? 0 } pending={summary?.pending_distribution_authorizations ?? 0 } showFooter />
    
    {/* <StatCard
      label="Fund Release Approvals"
      value={summary?.pending_fund_release_approvals ?? 0}
      
    />
    <StatCard
      label="Distribution Authorizations"
      value={summary?.pending_distribution_authorizations ?? 0}
    /> */}
  </Grid>
);

export default TStatsGrid;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;

  @media (max-width: 1024px) {
    grid-template-columns: repeat(2, 1fr);
  }

  @media (max-width: 600px) {
    grid-template-columns: 1fr;
  }
`;
