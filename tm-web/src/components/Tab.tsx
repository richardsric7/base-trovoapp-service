"use client";
import React from "react";
import styled from "styled-components";

type TabProps = {
  tabs: { key: string; label: string }[];
  currentTab: string;
  setCurrentTab: (tab: string) => void;
  children?: React.ReactNode;
  tabContainerStyle?: React.CSSProperties;
  tabButtonStyle?: React.CSSProperties;
  activeTabButtonStyle?: React.CSSProperties;
  tabContentStyle?: React.CSSProperties;
  tabContentBackground?: string;
};

const Tab: React.FC<TabProps> = ({
  tabs,
  currentTab,
  setCurrentTab,
  children,
  tabContainerStyle,
  tabButtonStyle,
  activeTabButtonStyle,
  tabContentStyle,
  tabContentBackground,
}) => {
  return (
    <div>
      <TabContainer style={tabContainerStyle}>
        {tabs.map((tab) => (
          <TabButton
            key={tab.key}
            isTabActive={currentTab === tab.key}
            onClick={() => setCurrentTab(tab.key)}
            style={
              currentTab === tab.key ? activeTabButtonStyle : tabButtonStyle
            }
          >
            {tab.label}
          </TabButton>
        ))}
      </TabContainer>
      <TabContent
        style={tabContentStyle}
        tabContentBackground={tabContentBackground}
      >
        {children}
      </TabContent>
    </div>
  );
};

export default Tab;

const TabContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
  background-color: #ffffff;
  width: 40%;
  // height: 40px;
  padding: 4px;
  border-radius: 12px;

  @media (max-width: 1800px) {
    width: 27%; // Adjust width for screens up to 1800px
  }

  @media (max-width: 1400px) {
    width: 38%; // Further adjust width for smaller screens
  }
`;

const TabButton = styled.div<{ isTabActive: boolean }>`
  padding: 4px 16px;
  margin: 4px 8px 4px 4px;

  background-color: ${({ isTabActive }) =>
    isTabActive ? "#F2F6F9" : "transparent"};
  color: ${({ isTabActive }) => (isTabActive ? "#00225A" : "#333333")};
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  letter-spacing: 0.1px;
  text-align: center;
`;

const TabContent = styled.div<{ tabContentBackground?: string }>`
  padding: 16px;
  background-color: ${({ tabContentBackground }) =>
    tabContentBackground || "#ffffff"};
  border-radius: 12px;
`;
