import React from "react";
import Link from "next/link";
import Card from "../../components/Card";
import styled from "styled-components";
import usersGroupIcon from "@/assets/images/solar_users-group-rounded-linear.svg";
import activeIcon from "@/assets/images/healthicons_happy-outline.svg";
import merchantsIcon from "@/assets/images/merchant-icon.svg";
import newUsersIcon from "@/assets/images/profile-add.svg";

import SteppedLineChart from "./SteppedLineChart";
import { useP2pMetricsQuery, useP2pDistributionReportQuery } from "@/redux/api/p2p";
import { formatAmount } from "./reports/reportUtils";

const Statistics = () => {
  const { data, isLoading } = useP2pMetricsQuery();
  // A 30-day snapshot of what's actually trading right now - the full,
  // date-range-filterable breakdown lives in the Reports tab; this is just
  // a quick-glance widget, not a substitute for it.
  const { data: distribution } = useP2pDistributionReportQuery({ range: "30d" });
  const topAssets = (distribution?.data?.byAsset ?? []).slice(0, 3);

  return (
    <Container>
      <MetricsCardContainer>
        <Card
          text="Total Offers"
          heading={data?.data?.totalOffers ?? 0}
          img={usersGroupIcon}
        />

        <Card
          text="Active Offers"
          heading={data?.data?.activeOffers ?? 0}
          img={merchantsIcon}
        />

        <Card
          text="Total Orders"
          heading={data?.data?.totalOrders ?? 0}
          img={activeIcon}
        />

        <Card
          text="Completed Orders"
          heading={data?.data?.completedOrders ?? 0}
          img={newUsersIcon}
        />
      </MetricsCardContainer>

      <ChartsContainer>
        <TradedAssetContainer>
          <Heading>Most Traded Assets (last 30 days)</Heading>
          <Line></Line>
          {topAssets.length === 0 ? (
            <EmptyText>No trades in the last 30 days.</EmptyText>
          ) : (
            topAssets.map((asset) => (
              <AssetDetails key={asset.key}>
                <AssetInfo>
                  <Name>{asset.key}</Name>
                </AssetInfo>
                <TradesCount>
                  {asset.count} trade{asset.count === 1 ? "" : "s"} · {formatAmount(asset.volume)} traded
                </TradesCount>
              </AssetDetails>
            ))
          )}
          <SeeMore href="/trovop2p?tab=reports">See the full market reports →</SeeMore>

        </TradedAssetContainer>
      </ChartsContainer>
      <SteppedLineContainer>
        <SteppedLineChart />
      </SteppedLineContainer>
    </Container>
  );
};

export default Statistics;
const Container = styled.div``;

const MetricsCardContainer = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
  background-color: transparent;
`;
const ChartsContainer = styled.div`
  display: flex;
  // align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
  background-color: transparent;
  // & > div {
  //   max-height: 400px; /* Set a consistent maximum height */
  // }
`;

const TradedAssetContainer = styled.div`
  padding: 40px 32px 40px 32px;
  gap: 16px;
  border-radius: 24px;
  background-color: #ffffff;
  margin-top: 10px;
`;
const AssetDetails = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #f2f6f9;
  padding: 8px 12px 8px 12px;
  border-radius: 12px;
  margin-top: 8px;
`;
const Heading = styled.h2`
  font-size: 20px;
  font-weight: 700;
  line-height: 28px;
  color: #00225a;
`;

const TradesCount = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 16px;
  color: #00225a;
  background-color: #ffffff;
  padding: 4px 8px 4px 8px;
  gap: 10px;
  border-radius: 8px;
`;

const AssetInfo = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;
const Name = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 16px;
  color: #004988;
`;

const EmptyText = styled.p`
  font-size: 13px;
  color: #828282;
  padding: 12px 0;
`;

const SeeMore = styled(Link)`
  display: inline-block;
  margin-top: 16px;
  font-size: 13px;
  font-weight: 600;
  color: #007cdf;
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
`;

const Line = styled.div`
  width: 100%;
  height: 1px;
  background-color: #e5e5ef;
  margin: 10px 0 15px 0;
`;

const SteppedLineContainer = styled.div`
  padding: 40px 32px 40px 32px;
  border-radius: 24px;
  background-color: #ffffff;
  margin-top: 10px;
`;
