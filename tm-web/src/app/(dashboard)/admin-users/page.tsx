"use client";

import React, { useState } from "react";
import styled from "styled-components";

import Permissions from "./components/PermissionsTab";
import Tab from "@/components/Tab";
import AdminUsersTab from "./components/AdminUsersTab";
const AdminManagementpage = () => {
  const [currentTab, setCurrentTab] = useState("adminusers");
  const tabs = [
    { key: "adminusers", label: "AdminUsers" },
    // { key: "permissions", label: "Permissions" },
  ];

  return (
    <>
      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabContainerStyle={{
          width: "100%",
          maxWidth: "150px",
        }}
      >
        {currentTab === "adminusers" && <AdminUsersTab />}
        {/* {currentTab === "permissions" && <Permissions />} */}
      </Tab>
    </>
  );
};

export default AdminManagementpage;

const Container = styled.section``;

const TabContainer = styled.div`
  display: flex;
  margin-bottom: 24px;
  background-color: #ffffff;
  width: 30%;
  height: 40px;
  padding: 4px;
  border-radius: 12px;
`;
