"use client";
import Tab from "@/components/Tab";
import { useState } from "react";
import styled from "styled-components";
import UsersTable from "./components/UsersTable";
import ViolationsLog from "./components/ViolationsLog";
export default function Users() {
  const [currentTab, setCurrentTab] = useState("users");

  const tabs = [
    { key: "users", label: "Users" },
    { key: "voilations", label: "Violations Log" },
  ];

  return (
    <>
      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabContainerStyle={{
          width: "100%",
          maxWidth: "250px",
        }}
      >
        {currentTab === "users" && <UsersTable />}
        {currentTab === "voilations" && <ViolationsLog />}
      </Tab>
    </>
  );
}

const Container = styled.section`
  padding: 32px;
  background-color: #ffffff;
  border-radius: 12px;
`;
