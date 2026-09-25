"use client";

import React, { useState } from "react";
import styled from "styled-components";
import Tab from "@/components/Tab";
import ManagedSecretsTable from "./components/ManagedSecretsTable";
import AuditLogTable from "./components/AuditLogTable";

const VaultSignerSettingsPage = () => {
  const [currentTab, setCurrentTab] = useState("managedsecrets");
  const tabs = [
    { key: "managedsecrets", label: "Managed Secrets" },
    { key: "auditlog", label: "Audit Log" },
  ];

  return (
    <PageContainer>
      <Title>Vault Signer</Title>
      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabContainerStyle={{ width: "100%", maxWidth: "320px" }}
      >
        {currentTab === "managedsecrets" && <ManagedSecretsTable />}
        {currentTab === "auditlog" && <AuditLogTable />}
      </Tab>
    </PageContainer>
  );
};

export default VaultSignerSettingsPage;

const PageContainer = styled.section``;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 20px 0;
  color: #00225a;
`;
