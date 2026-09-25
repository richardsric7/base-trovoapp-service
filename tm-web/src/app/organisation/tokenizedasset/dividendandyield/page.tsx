"use client";

import Tab from "@/components/Tab";
import React, { useState } from "react";

import { useRouter } from "next/navigation";
import { FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";
import DividendTable from "./components/DividendTable";
import YieldTable from "./components/YieldTable";

const DividendAndYieldPage = () => {
  const [currentTab, setCurrentTab] = useState("dividend");
  const router = useRouter();

  const tabs = [
    { key: "dividend", label: "Dividend" },
    { key: "yield", label: "yield" },
  ];

  const handleBack = () => {
    router.back();
  };
  return (
    <>
      <BackButtonLink onClick={handleBack}>
        <FaArrowLeft size={18} />
      </BackButtonLink>
      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabContainerStyle={{
          width: "100%",
          maxWidth: "200px",
        }}
      >
        {currentTab === "dividend" && <DividendTable />}
        {currentTab === "yield" && <YieldTable />}
      </Tab>
    </>
  );
};

export default DividendAndYieldPage;

const BackButtonLink = styled.button`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
  background: none;
  cursor: pointer;
  border: none;
  padding: 24px;
`;
