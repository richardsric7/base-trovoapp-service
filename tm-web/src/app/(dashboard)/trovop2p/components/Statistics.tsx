import React from "react";
import Card from "../../components/Card";
import styled from "styled-components";
import usersGroupIcon from "@/assets/images/solar_users-group-rounded-linear.svg";
import activeIcon from "@/assets/images/healthicons_happy-outline.svg";
import merchantsIcon from "@/assets/images/merchant-icon.svg";
import newUsersIcon from "@/assets/images/profile-add.svg";

import P2pBarChart from "./P2pBarChart";
import Image from "next/image";
import trovIcon from "@/assets/images/TROVTokenicon.svg";
import xbnIcon from "@/assets/images/bantuIcon.svg";
import CircularGauge from "./CirclularGuage";
import SteppedLineChart from "./SteppedLineChart";
import { useP2pMetricsQuery } from "@/redux/api/p2p";

const Statistics = () => {
  const { data, isLoading } = useP2pMetricsQuery();

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
        <P2pBarChart
          data={{
            data: {
              totalOrders: data?.data?.totalOrders ?? 0,
            },
          }}
        />

        <div>
          <RoundContainer>
            <Stats>Statistics</Stats>
            <Heading>Service Resource Utilization</Heading>
            <Line></Line>
            <CircularGauge
              value={72}
              size={300}
              thickness={16}
              primaryColor="#0088ff"
              backgroundColor="#e0e0e0"
            />
          </RoundContainer>
          <TradedAssetContainer>
            <Heading>Most Traded Assets</Heading>

            <Line></Line>
            <AssetDetails>
              <AssetInfo>
                <Image src={trovIcon} alt="trov-token-icon" />
                <Name>Trov</Name>
              </AssetInfo>
              <TradesCount>500 Trades</TradesCount>
            </AssetDetails>
            <AssetDetails>
              <AssetInfo>
                <Image src={xbnIcon} alt="trov-token-icon" />
                <Name>XBN</Name>
              </AssetInfo>
              <TradesCount>492 Trades</TradesCount>
            </AssetDetails>
          </TradedAssetContainer>
        </div>
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

const RoundContainer = styled.div`
  padding: 40px 32px 40px 32px;
  gap: 48px;
  border-radius: 24px;
  background-color: #ffffff;
  margin-top: 20px;
`;
const Stats = styled.p`
  font-size: 16px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  text-underline-position: from-font;
  text-decoration-skip-ink: none;
  color: #828282;
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
