"use client";
import React, { useState } from "react";
import styled from "styled-components";
import ComplianceSidebar from "./components/ComplianceSidebar";
import ComplianceSubmissions from "./components/ComplianceSubmissions";
import ComplianceRecords from "./components/ComplianceRecords";

const CompliancePage = () => {
  const [activeTab, setActiveTab] = useState("documents");

  return (
    <PageContainer>
      <ContentWrapper>
        <ComplianceSidebar activeTab={activeTab} setActiveTab={setActiveTab} />
        <MainContent>
          {activeTab === "documents" ? (
            <ComplianceSubmissions />
          ) : activeTab === "records" ? (
            <ComplianceRecords />
          ) : (
            <Placeholder>
              <Title>Organization Details</Title>
              <p>Content for {activeTab} will go here.</p>
            </Placeholder>
          )}
        </MainContent>
      </ContentWrapper>
    </PageContainer>
  );
};

export default CompliancePage;

const PageContainer = styled.div`
  padding: 24px;
  width: 100%;
`;

const ContentWrapper = styled.div`
  display: flex;
  gap: 32px;
  align-items: flex-start;
`;

const MainContent = styled.div`
  flex: 1;
  display: flex;
`;

const Placeholder = styled.div`
  background: #ffffff;
  border-radius: 24px;
  padding: 32px;
  flex: 1;
  min-height: 400px;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin: 0 0 16px 0;
`;
