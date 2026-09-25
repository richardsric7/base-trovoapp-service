"use client";
import Tab from "@/components/Tab";
import { useState } from "react";

import TransactionRecordsTable from "./components/TransactionRecordsTable";
import AssetRecordsTable from "./components/AccessRecordsTable";
import ComplianceTemplatesTab from "./components/ComplianceTemplatesTab";
import ComplianceRequirementsQueue from "./components/ComplianceRequirementsQueue";

const ComplaincePage = () => {
  const [currentTab, setCurrentTab] = useState("accessRecords");

  const tabs = [
    { key: "accessRecords", label: "Access Records" },
    { key: "transactionRecords", label: "Transaction Records" },
    { key: "complianceTemplates", label: "Compliance Templates" },
    { key: "complianceRequirements", label: "Compliance Requirements" },
  ];
  return (
    <>
      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabContainerStyle={{
          width: "100%",
          maxWidth: "700px",
        }}
      >
        {currentTab === "accessRecords" && <AssetRecordsTable />}
        {currentTab === "transactionRecords" && <TransactionRecordsTable />}
        {currentTab === "complianceTemplates" && <ComplianceTemplatesTab />}
        {currentTab === "complianceRequirements" && (
          <ComplianceRequirementsQueue />
        )}
      </Tab>
    </>
  );
};

export default ComplaincePage;
