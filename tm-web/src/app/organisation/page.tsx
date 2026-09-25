"use client";

import React, { useEffect } from "react";
import styled from "styled-components";
import { useRouter } from "next/navigation";
import { useGetOrganizationDetailsDataQuery } from "@/redux/api/org/api";
import { getOrganizationDashboardRoute } from "@/utils/organizationRouting";
import DashboardHeader from "./components/DashboardHeader";
import RecentAssetAssigned from "./components/RecentAssetAssigned";
import PortfolioSummary from "./components/PortfolioSummary";
import RecentActivity from "./components/RecentActivities";

const StakeHoldersDashBoard = () => {
  const router = useRouter();
  const { data } = useGetOrganizationDetailsDataQuery();

  useEffect(() => {
    const stakeholderType =
      data?.data.organization.stakeholder_type || data?.data.organization.type;
    const dashboardRoute = getOrganizationDashboardRoute(stakeholderType);

    if (dashboardRoute !== "/organisation") {
      router.replace(dashboardRoute);
    }
  }, [data, router]);

  return (
    <Container>
      <TopSection>
        <DashboardHeader />
      </TopSection>

      <MiddleSection>
        <RecentAssetAssigned />
        <PortfolioSummary />
      </MiddleSection>

      <RecentActivity />
    </Container>
  );
};

export default StakeHoldersDashBoard;

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
