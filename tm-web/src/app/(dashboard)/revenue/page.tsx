"use client";

import Tab from "@/components/Tab";
import React, { useState } from "react";
import styled from "styled-components";
import RevenueChart from "./components/RevenueChart";
import dynamic from "next/dynamic";

const TodayTab = dynamic(() => import("./components/TodayTab"), {
  ssr: false,
});

const RevenuePage = () => {
  const [currentTab, setCurrentTab] = useState("today");

  const tabs = [
    { key: "today", label: "Today" },
    { key: "last7days", label: "Last 7 Days" },
    { key: "last30days", label: "Last 30 Days" },
    { key: "daterange", label: "Date Range" },
  ];
  return (
    <>
      <Container>
        <ContentWrapper>
          <Heading>Revenue Inflow</Heading>
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
        </ContentWrapper>
        {currentTab === "today" && <TodayTab />}
        {currentTab === "last7days" && <TodayTab />}
        {currentTab === "last30days" && <TodayTab />}
        {currentTab === "daterange" && <TodayTab />}
      </Container>

      <RevenueChart />
    </>
  );
};

export default RevenuePage;

const Container = styled.section`
  padding: 20px 10px;
  background-color: #fff;
  border-radius: 24px;
`;

const ContentWrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const Heading = styled.h1`
  font-weight: 600;
  font-size: 20px;
  line-height: 28px;
  color: #00225a;
`;
