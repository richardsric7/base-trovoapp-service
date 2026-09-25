"use client";
import React from "react";
import styled from "styled-components";
import { FaCheckCircle } from "react-icons/fa";
import { LuClock3 } from "react-icons/lu";

interface ComplianceSidebarProps {
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

const ComplianceSidebar = ({
  activeTab,
  setActiveTab,
}: ComplianceSidebarProps) => {
  const menuItems = [
    {
      id: "details",
      label: "Organization Details",
      icon: <FaCheckCircle color="#27AE60" />,
    },
    { id: "documents", label: "Documents", icon: <LuClock3 color="#007CDF" /> },
    {
      id: "records",
      label: "Compliance Records",
      icon: <LuClock3 color="#007CDF" />,
    },
  ];

  return (
    <SidebarContainer>
      {menuItems.map((item) => (
        <MenuItem
          key={item.id}
          $isActive={activeTab === item.id}
          onClick={() => setActiveTab(item.id)}
        >
          <IconWrapper>{item.icon}</IconWrapper>
          <Label>{item.label}</Label>
        </MenuItem>
      ))}
    </SidebarContainer>
  );
};

export default ComplianceSidebar;

const SidebarContainer = styled.div`
  width: 280px;
  background: #ffffff;
  border-radius: 24px;
  padding: 24px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  height: fit-content;
`;

const MenuItem = styled.div<{ $isActive: boolean }>`
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-radius: 12px;
  cursor: pointer;
  background: ${(props) => (props.$isActive ? "#F2F6F9" : "transparent")};
  transition: all 0.2s ease;

  &:hover {
    background: #f8fbff;
  }
`;

const IconWrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
`;

const Label = styled.span`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
`;
