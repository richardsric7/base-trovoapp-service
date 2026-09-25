"use client";
import Tab from "@/components/Tab";
import dynamic from "next/dynamic";
import { useState } from "react";
import styled from "styled-components";

// Dynamically import P2pUsers
const P2pUsers = dynamic(() => import("./components/P2pUsers"), {
  ssr: false,
});

// Static import (if needed dynamically, convert like above)
import P2pViolationLog from "./components/P2pViolationLog";

const TrovoP2pUsers = () => {
  const [currentTab, setCurrentTab] = useState("users");

  const tabs = [
    { key: "users", label: "Users" },
    { key: "violation", label: "Violations Log" },
  ];

  return (
    <Container>
      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabContainerStyle={{
          width: "100%",
          maxWidth: "250px",
        }}
        tabContentStyle={{
          backgroundColor: "transparent",
          padding: 0,
        }}
      >
        {currentTab === "users" && <P2pUsers />}
        {currentTab === "violation" && <P2pViolationLog />}
      </Tab>
    </Container>
  );
};

export default TrovoP2pUsers;

const Container = styled.section`
  padding: 20px 10px;
  background-color: #f5f5f5;
`;
