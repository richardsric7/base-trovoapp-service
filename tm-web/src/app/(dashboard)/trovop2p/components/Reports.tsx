"use client";
import React, { useState } from "react";
import styled from "styled-components";
import { ReportRange, useP2pVolumeReportQuery } from "@/redux/api/p2p";
import ReportRangeSelector from "./reports/ReportRangeSelector";
import VolumeReportSection from "./reports/VolumeReportSection";
import DistributionReportSection from "./reports/DistributionReportSection";
import DisputeReportSection from "./reports/DisputeReportSection";
import RevenueReportSection from "./reports/RevenueReportSection";
import GrowthReportSection from "./reports/GrowthReportSection";
import MerchantLeaderboardSection from "./reports/MerchantLeaderboardSection";

// Reports is the P2P Market Reports tab: trading volume, distribution,
// disputes, fee revenue, marketplace growth, and a full merchant
// leaderboard - every chart is backed by tm-api's /p2p/reports/* endpoints
// (real aggregations over the P2P module's own data), gated server-side
// on the ACCESS_REPORTS admin permission. A 403 from that gate (a real
// admin without report access, e.g. a VIEW_ONLY_ADMIN) renders a
// restricted-access message here instead of a page full of broken charts.
const Reports = () => {
  const [range, setRange] = useState<ReportRange>("30d");
  const { error, isLoading } = useP2pVolumeReportQuery({ range });
  const accessDenied = (error as any)?.status === 403;

  if (isLoading) {
    return <Loading>Loading reports...</Loading>;
  }

  if (accessDenied) {
    return (
      <Restricted>
        <RestrictedTitle>Access restricted</RestrictedTitle>
        <RestrictedBody>
          Your admin role doesn&apos;t have permission to view P2P market reports. Ask a super admin to grant
          you the &quot;Access Reports&quot; permission.
        </RestrictedBody>
      </Restricted>
    );
  }

  return (
    <Container>
      <Header>
        <Title>P2P Market Reports</Title>
        <ReportRangeSelector value={range} onChange={setRange} />
      </Header>

      <VolumeReportSection range={range} />
      <DistributionReportSection range={range} />
      <DisputeReportSection range={range} />
      <RevenueReportSection range={range} />
      <GrowthReportSection range={range} />
      <MerchantLeaderboardSection />
    </Container>
  );
};

export default Reports;

const Container = styled.div``;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 700;
  color: #00225a;
  margin: 0;
`;

const Loading = styled.div`
  padding: 60px;
  text-align: center;
  color: #828282;
`;

const Restricted = styled.div`
  background-color: #ffffff;
  border-radius: 24px;
  padding: 60px 32px;
  text-align: center;
  margin-top: 20px;
`;

const RestrictedTitle = styled.h3`
  color: #00225a;
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 8px 0;
`;

const RestrictedBody = styled.p`
  color: #828282;
  font-size: 14px;
  max-width: 480px;
  margin: 0 auto;
`;
