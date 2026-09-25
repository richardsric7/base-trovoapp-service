"use client";

import { useUsersMetricsQuery } from "@/redux/api/usersmetrics";
import React from "react";
import styled from "styled-components";
import Card from "../../components/Card";
import usersGroupIcon from "@/assets/images/solar_users-group-rounded-linear.svg";
import activeIcon from "@/assets/images/healthicons_happy-outline.svg";
import sessionsIcon from "@/assets/images/fluent_data-usage-20-regular.svg";
import timeIcon from "@/assets/images/carbon_time.svg";
import PieChartSection from "./PieChartSection";
import BarChart from "./BarChart";
import TopReferrals from "./TopReferrals";
import CountriesStatistics from "./CountriesStatistics";
import RecentRegisteration from "./RecentRegisteration";
const TrovoAppUsers = () => {
  const { data, isLoading, isError, error } = useUsersMetricsQuery();

  // const percentageChanges = data?.data?.weekly_percentage_change || {
  //   active_users_percentage_change: 0,
  //   new_users_percentage_change: 0,
  //   total_sessions_percentage_change: 0,
  //   average_time_percentage_change: 0,
  // };

  return (
    <Container>
      <MetricsCardContainer>
        <Card
          text="Total Registered Users"
          heading={data?.data?.total_users ?? 0}
          img={usersGroupIcon}
        />

        <Card
          text="Total Verified Users"
          heading={data?.data?.total_sessions ?? 0}
          img={sessionsIcon}
        />

        <Card
          text="Active Users"
          heading={data?.data?.total_active_users ?? 0}
          img={activeIcon}
        />
        <Card
          text="Average time per Session"
          heading={data?.data?.average_time_per_session ?? 0}
          img={timeIcon}
        />
      </MetricsCardContainer>
      <ChartsContainer>
        <BarChart
          data={{
            data: {
              daily_active_and_new_users: { active_users: 100, new_users: 20 },
              weekly_active_and_new_users: { active_users: 500, new_users: 80 },
              monthly_active_and_new_users: {
                active_users: 2000,
                new_users: 400,
              },
              total_users: data?.data?.total_users ?? 0,
            },
          }}
        />
        <PieChartSection />
      </ChartsContainer>

      <RefSection>
        <TopReferrals />
        <CountriesStatistics />
      </RefSection>
      <RecentRegisteration />
    </Container>
  );
};

export default TrovoAppUsers;

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

const RefSection = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
  background-color: transparent;
`;
