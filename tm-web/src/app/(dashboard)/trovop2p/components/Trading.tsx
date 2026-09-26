"use client";

import React from "react";
import Card from "../../components/Card";
import styled from "styled-components";
import totalTradesIcon from "@/assets/images/icon-park-outline_stock-market.svg";
import averageTradeIcon from "@/assets/images/carbon_chart-average.svg";
import tradeSuccessIcon from "@/assets/images/copy-success.svg";
import disputesIcon from "@/assets/images/message-question.svg";
import TradingHoursChart from "./TradingHoursChart";
import TradingVolumeChart from "./TradingVolumeChart";
import TopTraders from "./TopTraders";
import Trades from "./Trades";
import { useTradeStatisticsQuery } from "@/redux/api/p2p";
const Trading = () => {
  const { data } = useTradeStatisticsQuery();
  return (
    <Container>
      <MetricsCardContainer>
        <Card
          text="Total Completed Trades"
          heading={data?.data?.completedTrades ?? 0}
          img={totalTradesIcon}
        />

        <Card
          text="Total Trades"
          heading={data?.data?.totalTrades ?? 0}
          img={averageTradeIcon}
        />

        <Card
          text="Trade Success Rate"
          heading={Number((data?.data?.tradeSuccessRate ?? 0).toFixed(2))}
          img={tradeSuccessIcon}
        />

        <Card
          text="Open Disputes"
          heading={data?.data?.openDisputes ?? 0}
          img={disputesIcon}
        />
      </MetricsCardContainer>
      <TradingHoursChart />

      <Content>
        <TradingVolumeChart />
        <TopTraders />
      </Content>
      <Trades />
    </Container>
  );
};

export default Trading;
const Container = styled.div``;

const MetricsCardContainer = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
  background-color: transparent;
`;

const Content = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
`;
