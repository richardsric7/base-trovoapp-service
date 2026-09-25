"use client";

import React from "react";
import styled from "styled-components";
import AmDashboardHeader from "./components/AmDashboardHeader";

import AmPortfolioSummary from "./components/AmPortfolioSummary";
import AmRecentActivity from "./components/AmRecentActivities";
import AmRecentAssetAssigned from "./components/AmRecentAssetAssigned";
import { useGetAssetManagerDashboardQuery } from "@/redux/api/assetManager";

const AssetManagerDashBoard = () => {
  const { data, isLoading, isError, refetch } =
    useGetAssetManagerDashboardQuery();

  if (isLoading) return <StateMessage>Loading dashboard...</StateMessage>;

  if (isError) {
    return (
      <StateMessage>
        Unable to load dashboard.
        <RetryButton onClick={() => refetch()}>Try again</RetryButton>
      </StateMessage>
    );
  }

  const dashboard = data?.data;

  return (
    <Container>
      <TopSection>
        <AmDashboardHeader
          summary={dashboard?.summary}
        />
      </TopSection>

      <MiddleSection>
        <AmRecentAssetAssigned assets={dashboard?.recent_assets} />
        <AmPortfolioSummary summary={dashboard?.portfolio_summary} />
      </MiddleSection>

      <AmRecentActivity activities={dashboard?.recent_activity} />
    </Container>
  );
};

export default AssetManagerDashBoard;

const Container = styled.div`
  padding: 24px;

  min-height: 100vh;
`;

const TopSection = styled.div`
  margin-bottom: 20px;
`;

const MiddleSection = styled.div`
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
`;

const StateMessage = styled.div`
  padding: 40px 24px;
  color: #00225a;
`;

const RetryButton = styled.button`
  display: block;
  margin-top: 12px;
  padding: 8px 16px;
  border: 0;
  border-radius: 8px;
  background: #007cdf;
  color: #fff;
  cursor: pointer;
`;
