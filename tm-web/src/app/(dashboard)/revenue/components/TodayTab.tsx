import React from "react";
import RevCard from "./RevCard";
import styled from "styled-components";

const TodayTab = () => {
  const todayData = [
    {
      text: "Trovo Wallet Swap Fee",
      value: "$456 000",
      badge: "Today",
    },
    {
      text: "Trovo Patron Membership",
      value: "$2 000",
      badge: "Today",
      dotColor: "#007CDF",
    },
    {
      text: "Trovo Wallet Recovery Subscription Fee",
      value: "$2 000",
      badge: "Today",
      dotColor: "#62ACE8",
    },

    {
      text: "Trovo P2P Fee",
      value: "$72 000",
      badge: "Today",
    },
    {
      text: "Asset Tokenization Fee",
      value: "$500 000",
      badge: "Today",
      dotColor: "#D5E1F8",
    },
    {
      text: "Shared Access Payment Transaction Fee",
      value: "$2 000",
      badge: "Today",
      dotColor: "#4F9A94",
    },
    {
      text: "Sub Wallet Creation Fee",
      value: "$119 000",
      badge: "Today",
      dotColor: "#95C2BF",
    },
    {
      text: "Wrapped Asset Withdrawal Fee",
      value: "$35 000",
      badge: "Today",
      dotColor: "#C08B8C",
    },
    {
      text: "Statement of account Generation Fee",
      value: "$2 000",
      badge: "Today",
      dotColor: "#D9B9BA",
    },
  ];

  return (
    <Container>
      {todayData.map((item, idx) => (
        <RevCard
          text={item.text}
          value={item.value}
          key={idx}
          badge={item.badge}
          dotColor={item.dotColor}
        />
      ))}
    </Container>
  );
};

export default TodayTab;

const Container = styled.div`
  display: grid;
  align-items: center;
  justify-content: center;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
`;
