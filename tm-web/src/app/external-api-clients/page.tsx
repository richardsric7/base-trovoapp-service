"use client";

import styled from "styled-components";
import Image from "next/image";
import TransactionHistory from "./components/TransactionHistory";
import Revenue from "./components/Revenue";
import RecentlyAddedCustomers from "./components/RecentlyAddedCustomers";
import CardStats from "./components/StatsCard";
import usersGroupRounded from "@/assets/images/solar_users-group-rounded-linear.svg";
import clarityWalletLine from "@/assets/images/solar_money-bag-linear.svg";

const ExternalApiClientsDashboard = () => {
  const dummyStats = [
    {
      label: "Total Customers",
      value: "5000",
      img: usersGroupRounded,
      hasShadow: true,
    },
    {
      label: "Transactions",
      value: "1,000",
      img: clarityWalletLine,
      hasShadow: true,
      trendValue: "+21.01%",
      isPositive: true,
      period: "This week",
    },
    {
      label: "Revenue",
      value: "1,000",
      suffix: "CNGN",
      img: clarityWalletLine,
      hasShadow: true,
      trendValue: "+21.01%",
      isPositive: true,
      period: "This week",
    },
  ];

  return (
    <Container>
      <StatsSection>
        {dummyStats.map((stat, index) => (
          <CardStats
            key={index}
            text={stat.label}
            heading={stat.value}
            suffix={stat.suffix}
            img={<Image src={stat.img} alt="" width={24} height={24} />}
            hasShadow={stat.hasShadow}
            trendValue={stat.trendValue}
            isPositive={stat.isPositive}
            period={stat.period}
          />
        ))}
      </StatsSection>
      <SecondSection>
        <Revenue />
        <RecentlyAddedCustomers />
      </SecondSection>

      <TransactionHistory />
    </Container>
  );
};

export default ExternalApiClientsDashboard;

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 32px;
`;

const StatsSection = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
`;

const SecondSection = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
`;
