import React from "react";
import styled from "styled-components";
import { LaRecentAsset } from "@/redux/api/legalAdviser";
import StatCard from "../../components/StatsCard";

interface LaStatsGridProps {
  assetCount?: number;
  recentAssets?: LaRecentAsset[];
}

const LaStatsGrid = ({ assetCount, recentAssets = [] }: LaStatsGridProps) => {
  const pendingStructuring = recentAssets.filter(
    (asset) => asset.assignment?.status !== "active",
  ).length;
  const completedStructuring = recentAssets.length - pendingStructuring;

  return (
    <Grid>
      <StatCard label="Assets Assigned" value={assetCount ?? recentAssets.length} />
      <StatCard label="Pending Structuring" value={pendingStructuring} />
      <StatCard label="Structuring Completed" value={completedStructuring} />
    </Grid>
  );
};

export default LaStatsGrid;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
`;
