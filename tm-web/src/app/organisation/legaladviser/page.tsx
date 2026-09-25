"use client";

import React from "react";
import styled from "styled-components";
import LaDashboardHeader from "./components/LaDashboardHeader";

import LaPortfolioSummary from "./components/LaPortfolioSummary";
import LaRecentActivity from "./components/LaRecentActivities";
import LaRecentAssetAssigned from "./components/LaRecentAssetAssigned";
import { useGetLegalAdviserDashboardQuery } from "@/redux/api/legalAdviser";

const LegalAdviserDashBoard = () => {
  const { data, isLoading, isError, refetch } =
    useGetLegalAdviserDashboardQuery();

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
        <LaDashboardHeader
          assetCount={dashboard?.portfolio_summary?.asset_count}
          recentAssets={dashboard?.recent_assets}
        />
      </TopSection>

      <MiddleSection>
        <LaRecentAssetAssigned assets={dashboard?.recent_assets} />
        <LaPortfolioSummary summary={dashboard?.portfolio_summary} />
      </MiddleSection>

      <LaRecentActivity activities={dashboard?.recent_activity} />
    </Container>
  );
};

export default LegalAdviserDashBoard;

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
