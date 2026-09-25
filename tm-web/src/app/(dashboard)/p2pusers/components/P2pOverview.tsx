"use client";

import React from "react";
import Card from "../../components/Card";
import styled from "styled-components";
import usersGroupIcon from "@/assets/images/solar_users-group-rounded-linear.svg";
import activeIcon from "@/assets/images/healthicons_happy-outline.svg";
import sessionsIcon from "@/assets/images/fluent_data-usage-20-regular.svg";
import timeIcon from "@/assets/images/carbon_time.svg";
import PeerTraders from "./PeerTraders";
import UserTradeVolume from "./UserTradeVolume";

const P2pOverview = () => {
  return (
    <>
      <Container>
        <Title>User Statistics</Title>
        <MetricsCardContainer>
          <Card
            text="Total Completed Trades"
            heading={23}
            img={usersGroupIcon}
            hasShadow
          />

          <Card
            text="Average Trade Size"
            heading="$115k"
            img={sessionsIcon}
            hasShadow
          />

          <Card
            text="Trade Success Rate"
            heading={`${20}%`}
            img={activeIcon}
            hasShadow
          />
          <Card text="Trades on Appeal" heading={20} img={timeIcon} hasShadow />
        </MetricsCardContainer>
      </Container>
      <Content>
        <UserTradeVolume />
        <PeerTraders />
      </Content>
    </>
  );
};

export default P2pOverview;

const Container = styled.div`
  border-radius: 16px;
  padding: 24px;
  background-color: #ffffff;
`;

const MetricsCardContainer = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
`;

const Title = styled.h1`
  font-weight: 600;
  font-size: 20px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #00225a;
  margin-bottom: 8px;
`;

const Content = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
`;
