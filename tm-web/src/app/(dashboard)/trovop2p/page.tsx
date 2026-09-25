"use client";
import Tab from "@/components/Tab";
import React, { useState } from "react";
import styled from "styled-components";

import dynamic from "next/dynamic";

export default function TrovoP2P() {
  const [currentTab, setCurrentTab] = useState("statistics");

  const TrovoP2PStatistics = dynamic(() => import("./components/Statistics"), {
    ssr: false,
  });

  const TrovoP2PTrading = dynamic(() => import("./components/Trading"), {
    ssr: false,
  });
  const tabs = [
    { key: "statistics", label: "Statistics" },
    { key: "trading", label: "Trading" },
  ];

  return (
    <Container>
      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabContainerStyle={{
          width: "100%",
          maxWidth: "205px",
        }}
        tabContentStyle={{
          backgroundColor: "transparent",
          padding: 0,
        }}
      >
        {currentTab === "statistics" && <TrovoP2PStatistics />}
        {currentTab === "trading" && <TrovoP2PTrading />}
      </Tab>
    </Container>
  );
}

const Container = styled.section`
  padding: 20px 10px;
  background-color: #f5f5f5;
`;
