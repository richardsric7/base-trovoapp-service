"use client";

import Tab from "@/components/Tab";
import React, { Suspense, useState } from "react";
import { useSearchParams } from "next/navigation";
import PayoutsList from "./components/PayoutsList";
import FeeAndEngine from "./components/FeeAndEngine";
import Reports from "./components/Reports";

// Proceeds payouts (dividends and yields of tokenized assets), paid by
// payout-engine and driven from here.
const ProceedPayouts = () => {
  const params = useSearchParams();
  const [currentTab, setCurrentTab] = useState("payouts");

  const tabs = [
    { key: "payouts", label: "Payouts" },
    { key: "fees", label: "Fee & engine" },
    { key: "reports", label: "Reports" },
  ];

  return (
    <Tab tabs={tabs} currentTab={currentTab} setCurrentTab={setCurrentTab} tabContainerStyle={{ width: "100%", maxWidth: "420px" }}>
      {currentTab === "payouts" && <PayoutsList asset={params.get("asset") ?? undefined} />}
      {currentTab === "fees" && <FeeAndEngine />}
      {currentTab === "reports" && <Reports />}
    </Tab>
  );
};

const ProceedPayoutsPage = () => (
  <Suspense fallback={null}>
    <ProceedPayouts />
  </Suspense>
);

export default ProceedPayoutsPage;
