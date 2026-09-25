"use client";

import React, { useState } from "react";
import styled from "styled-components";
import Tab from "@/components/Tab";
import dynamic from "next/dynamic";

const TrovoAppPage = () => {
  // Dynamically import the components with SSR disabled
  const TrovoAppUsers = dynamic(
    () => import("./trovoapp/components/TrovoAppUsers"),
    { ssr: false }
  );
  const TrovoAppTransactions = dynamic(
    () => import("./trovoapp/components/TrovoAppTransactions"),
    { ssr: false }
  );
  const TrovoAppCustomerSupport = dynamic(
    () => import("./trovoapp/components/TrovoAppCustomerSupport"),
    { ssr: false }
  );
  const [currentTab, setCurrentTab] = useState("users");
  const tabs = [
    { key: "users", label: "Users" },
    // { key: "transactions", label: "Transactions" },
    // { key: "customer-support", label: "Customer Support" },
  ];
  return (
    <Container>
      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabContainerStyle={{
          width: "100%",
          maxWidth: "100px",
          // maxWidth: "405px",
        }}
        tabContentStyle={{
          backgroundColor: "transparent",
          padding: 0,
        }}
      >
        {currentTab === "users" && <TrovoAppUsers />}
        {/* {currentTab === "transactions" && <TrovoAppTransactions />}
        {currentTab === "customer-support" && <TrovoAppCustomerSupport />} */}
      </Tab>
    </Container>
  );
};

export default TrovoAppPage;

const Container = styled.section`
  padding: 20px 10px;
  background-color: #f5f5f5;
`;
