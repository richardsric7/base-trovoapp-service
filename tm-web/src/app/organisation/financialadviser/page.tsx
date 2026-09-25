"use client";

import React from "react";
import styled from "styled-components";
import FaDashboardHeader from "./components/FaDashboardHeader";

import FaPortfolioSummary from "./components/FaPortfolioSummary";
import FaRecentActivity from "./components/FaRecentActivities";
import FaRecentAssetAssigned from "./components/FaRecentAssetAssigned";
import { useGetFinancialAdviserDashboardQuery } from "@/redux/api/financialAdviser";

const FinancialAdviserDashBoard = () => {
  const { data, isLoading, isError, refetch } =
    useGetFinancialAdviserDashboardQuery();

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
        <FaDashboardHeader
          assetCount={dashboard?.portfolio_summary?.asset_count}
          recentAssets={dashboard?.recent_assets}
        />
      </TopSection>

      <MiddleSection>
        <FaRecentAssetAssigned assets={dashboard?.recent_assets} />
        <FaPortfolioSummary summary={dashboard?.portfolio_summary} />
      </MiddleSection>

      <FaRecentActivity activities={dashboard?.recent_activity} />
    </Container>
  );
};

export default FinancialAdviserDashBoard;

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
