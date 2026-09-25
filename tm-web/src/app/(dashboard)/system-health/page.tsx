"use client";
import { useState } from "react";
import dynamic from "next/dynamic";
import styled from "styled-components";
import Tab from "@/components/Tab";

// Declared at module scope: calling dynamic() during render creates a new
// component type each time, which remounts the tab and restarts its fetch.
const HealthTab = dynamic(() => import("./components/HealthTab"));
const IssuesTab = dynamic(() => import("./components/IssuesTab"));
const TrendsTab = dynamic(() => import("./components/TrendsTab"), { ssr: false });
const TraceTab = dynamic(() => import("./components/TraceTab"));

/**
 * System Health — what is happening across the Trovo platform, without needing
 * access to the engineering tools.
 *
 * Each tab is loaded on demand: the charts pull in a charting library, and the
 * trace screen is rarely the first thing opened, so neither should slow down
 * the health view that people actually land on.
 */
const SystemHealthPage = () => {
  const [currentTab, setCurrentTab] = useState("health");
  const tabs = [
    { key: "health", label: "Health" },
    { key: "issues", label: "Issues" },
    { key: "trends", label: "Trends" },
    { key: "trace", label: "Trace" },
  ];

  return (
    <Container>
      <Title>System Health</Title>
      <Subtitle>
        Live status of every Trovo service, the errors users are hitting, and how the platform is
        performing.
      </Subtitle>

      {/* The active tab is rendered outside <Tab> rather than as its children.
          Passing it as children mounted the component more than once, so the
          instance on screen was not the one whose data had arrived - the panel
          sat on its loading state while the store already held the result. */}
      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabContainerStyle={{ width: "100%", maxWidth: "420px" }}
        tabContentStyle={{ backgroundColor: "transparent", padding: 0 }}
      />
      <TabBody>
        {currentTab === "health" && <HealthTab />}
        {currentTab === "issues" && <IssuesTab />}
        {currentTab === "trends" && <TrendsTab />}
        {currentTab === "trace" && <TraceTab />}
      </TabBody>
    </Container>
  );
};

export default SystemHealthPage;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;
  border-radius: 24px;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  color: #00225a;
`;

const Subtitle = styled.p`
  font-size: 14px;
  color: #828282;
  margin: 8px 0 24px;
  max-width: 640px;
  line-height: 1.5;
`;

const TabBody = styled.div`
  padding-top: 24px;
`;
