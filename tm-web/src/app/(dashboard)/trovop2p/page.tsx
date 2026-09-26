"use client";
import Tab from "@/components/Tab";
import React, { Suspense, useState } from "react";
import styled from "styled-components";
import { useSearchParams } from "next/navigation";

import dynamic from "next/dynamic";

const VALID_TABS = ["statistics", "trading", "reports"];

function TrovoP2PInner() {
  // Supports deep-linking straight to a tab (e.g. "Most Traded Assets"'s
  // "See the full market reports" link on the Statistics tab).
  const searchParams = useSearchParams();
  const tabParam = searchParams.get("tab");
  const [currentTab, setCurrentTab] = useState(
    tabParam && VALID_TABS.includes(tabParam) ? tabParam : "statistics",
  );

  const TrovoP2PStatistics = dynamic(() => import("./components/Statistics"), {
    ssr: false,
  });

  const TrovoP2PTrading = dynamic(() => import("./components/Trading"), {
    ssr: false,
  });
  const TrovoP2PReports = dynamic(() => import("./components/Reports"), {
    ssr: false,
  });
  const tabs = [
    { key: "statistics", label: "Statistics" },
    { key: "trading", label: "Trading" },
    { key: "reports", label: "Reports" },
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
        {currentTab === "reports" && <TrovoP2PReports />}
      </Tab>
    </Container>
  );
}

export default function TrovoP2P() {
  return (
    <Suspense fallback={null}>
      <TrovoP2PInner />
    </Suspense>
  );
}

const Container = styled.section`
  padding: 20px 10px;
  background-color: #f5f5f5;
`;
