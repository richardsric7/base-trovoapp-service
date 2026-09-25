"use client";
import Tab from "@/components/Tab";
import { useState } from "react";
import MinitngHistory from "../components/MinitngHistory";
import Tokens from "../components/Tokens";

const TokenManagementPage = () => {
  const [currentTab, setCurrentTab] = useState("tokens");

  const tabs = [
    { key: "tokens", label: "Tokens" },
    { key: "minitng-history", label: "Minitng History" },
  ];
  return (
    <>
      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabContainerStyle={{
          width: "100%",
          maxWidth: "260px",
        }}
      >
        {currentTab === "tokens" && <Tokens />}
        {currentTab === "minitng-history" && <MinitngHistory />}
      </Tab>
    </>
  );
};

export default TokenManagementPage;
