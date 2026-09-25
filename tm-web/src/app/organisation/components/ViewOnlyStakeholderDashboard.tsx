"use client";

import React from "react";
import styled from "styled-components";
import { useGetOrganizationDetailsDataQuery } from "@/redux/api/org";
import type { IStakeholderDashboardResponse } from "@/redux/api/sharedstakeholders";
import StatCard from "./StatsCard";
import TRecentAssetAssigned from "../trustee/components/TRecentAssetAssigned";
import TPortfolioSummary from "../trustee/components/TPortfolioSummary";
import TRecentActivity from "../trustee/components/TRecentActivities";

interface DashboardQueryResult {
  data?: IStakeholderDashboardResponse;
  isLoading: boolean;
  isError: boolean;
  refetch: () => unknown;
}

interface Props {
  useDashboardQuery: () => DashboardQueryResult;
}

const labelFromKey = (key: string) =>
  key
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");

const ViewOnlyStakeholderDashboard = ({ useDashboardQuery }: Props) => {
  const { data, isLoading, isError, refetch } = useDashboardQuery();
  const { data: organizationData } = useGetOrganizationDetailsDataQuery();

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
  const organization = organizationData?.data?.organization;
  const summaryItems = Object.entries(dashboard?.summary ?? {}).filter(
    ([, value]) => ["string", "number"].includes(typeof value),
  );

  if (summaryItems.length === 0) {
    summaryItems.push([
      "assigned_assets",
      dashboard?.portfolio_summary?.asset_count ?? dashboard?.recent_assets?.length ?? 0,
    ]);
  }

  return (
    <Container>
      <Header>
        <Identity>
          <Avatar>{organization?.name?.charAt(0).toUpperCase() || "?"}</Avatar>
          <div>
            <Title>{organization?.name || "Organization"}</Title>
            <Subtitle>
              {labelFromKey(dashboard?.role || organization?.type || "stakeholder")}
              {organization?.id ? ` | ${organization.id.slice(0, 7)}` : ""}
            </Subtitle>
          </div>
        </Identity>
        <StatsGrid>
          {summaryItems.map(([key, value]) => (
            <StatCard key={key} label={labelFromKey(key)} value={String(value)} />
          ))}
        </StatsGrid>
      </Header>

      <MiddleSection>
        <TRecentAssetAssigned assets={dashboard?.recent_assets} />
        <TPortfolioSummary summary={dashboard?.portfolio_summary} />
      </MiddleSection>

      <TRecentActivity activities={dashboard?.recent_activity} />
    </Container>
  );
};

export default ViewOnlyStakeholderDashboard;

const Container = styled.div`
  padding: 24px;
  min-height: 100vh;
`;

const Header = styled.div`
  background: #fff;
  padding: 32px 16px 24px;
  border-radius: 10px;
  margin-bottom: 20px;
`;

const Identity = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 18px;
`;

const Avatar = styled.div`
  width: 38px;
  height: 38px;
  border-radius: 50%;
  background: #007cdf;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
`;

const Title = styled.p`
  margin: 0;
  font-weight: 600;
  font-size: 24px;
  color: #00225a;
`;

const Subtitle = styled.p`
  margin: 0;
  font-size: 15px;
  color: #00225a;
  font-weight: 500;
`;

const StatsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;

  @media (max-width: 900px) {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
`;

const MiddleSection = styled.div`
  display: grid;
  grid-template-columns: minmax(0, 2fr) 380px;
  gap: 20px;
  margin-bottom: 20px;

  @media (max-width: 1100px) {
    grid-template-columns: 1fr;
  }
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
