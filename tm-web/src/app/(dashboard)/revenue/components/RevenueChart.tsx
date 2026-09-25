"use client";

import React, { useState } from "react";

import styled from "styled-components";
import Tab from "@/components/Tab";
import dynamic from "next/dynamic";

const RevenueBarChart = dynamic(() => import("./TodayRevChart"), {
  ssr: false,
});

const RevenueChart = () => {
  const [currentTab, setCurrentTab] = useState("today");

  const tabs = [
    { key: "today", label: "Today" },
    { key: "last7days", label: "Last 7 Days" },
    { key: "last30days", label: "Last 30 Days" },
  ];

  return (
    <Content>
      <ChartHeading>
        <div>
          <Text>Statistics</Text>
          <Title>Revenue Trends Over Time</Title>
        </div>

        <Tab
          tabs={tabs}
          currentTab={currentTab}
          setCurrentTab={setCurrentTab}
          tabButtonStyle={{
            borderRadius: "12px",
            fontSize: "16px",
            color: "#4F4F4F",
            padding: "6px 16px",
            width: "auto",
          }}
          activeTabButtonStyle={{
            backgroundColor: "#fff",
            color: "#00225A",
            border: "0px",
            width: "auto",
          }}
          tabContainerStyle={{
            width: "100%",
            display: "flex",
            gap: "8px",
            backgroundColor: "#F2F6F9",
            padding: 0,
          }}
          tabContentStyle={{
            backgroundColor: "transparent",
            padding: 0,
          }}
        ></Tab>
      </ChartHeading>
      {currentTab === "today" && <RevenueBarChart chartType="today" />}
      {currentTab === "last7days" && <RevenueBarChart chartType="last7days" />}
      {currentTab === "last30days" && (
        <RevenueBarChart chartType="last30days" />
      )}
    </Content>
  );
};

export default RevenueChart;

const Content = styled.div`
  padding: 24px;
  background-color: #fff;
  border-radius: 24px;
  margin-top: 20px;
`;

const ChartHeading = styled.div`
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const Text = styled.p`
  font-weight: 400;
  font-size: 16px;
  line-height: 20px;
  letter-spacing: 0%;
  color: #828282;
`;

const Title = styled.p`
  font-weight: 600;
  font-size: 20px;

  line-height: 28px;
  letter-spacing: 0px;
  color: #00225a;
`;
