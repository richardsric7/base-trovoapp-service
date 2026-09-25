"use client";

import React from "react";
import styled from "styled-components";
import AcDashboardHeader from "./components/AcDashboardHeader";
import AcPortfolioSummary from "./components/AcPortfolioSummary";
import AcRecentActivity from "./components/AcRecentActivities";
import AcRecentAssetAssigned from "./components/AcRecentAssetAssigned";
import { useGetCustodianDashboardQuery } from "@/redux/api/assetCustodian";

const AssetCustodianDashBoard = () => {
  const { data, isLoading, isFetching, isError, refetch } =
    useGetCustodianDashboardQuery();
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
        <AcDashboardHeader summary={dashboard?.summary} />
      </TopSection>

      <MiddleSection>
        <AcRecentAssetAssigned
          assets={dashboard?.recent_assets}
          isLoading={isFetching}
        />
        <AcPortfolioSummary portfolio={dashboard?.portfolio_summary} />
      </MiddleSection>

      <AcRecentActivity activities={dashboard?.recent_activity} />
    </Container>
  );
};

export default AssetCustodianDashBoard;

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
