"use client";
import React, { useState } from "react";
import Published from "./_components/Published";
import Drafts from "./_components/Drafts";
import styled from "styled-components";
import PrimaryButton from "@/components/PrimaryButton";
import { FaPlus } from "react-icons/fa6";
import { useRouter } from "next/navigation";
import { useGetJobsQuery } from "@/redux/api/jobs";

const JobsPage = () => {
  const router = useRouter();
  const [currentTab, setCurrentTab] = useState("pub");

  const { data: publishedData } = useGetJobsQuery({
    status: "published",
  });

  const { data: draftData } = useGetJobsQuery({
    status: "draft",
  });

  const tabs = [
    {
      key: "pub",
      label: `Published (${publishedData?.data?.length || 0})`,
    },
    {
      key: "drafts",
      label: `Drafts (${draftData?.data?.length || 0})`,
    },
  ];

  return (
    <Container>
      <Header>
        <Title>Jobs</Title>

        <PrimaryButton
          onClick={() => router.push("/jobs/new")}
          buttonStyle={{
            width: "20%",
            marginRight: "0rem",
            display: "flex",
            alignItems: "center",
            gap: "4px",
          }}
        >
          <FaPlus /> New Job Post
        </PrimaryButton>
      </Header>
      <CustomTab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
      >
        {currentTab === "pub" && <Published />}
        {currentTab === "drafts" && <Drafts />}
      </CustomTab>
    </Container>
  );
};

const CustomTab: React.FC<{
  tabs: { key: string; label: string }[];
  currentTab: string;
  setCurrentTab: (tab: string) => void;
  children?: React.ReactNode;
}> = ({ tabs, currentTab, setCurrentTab, children }) => {
  return (
    <TabWrapper>
      <TabNav>
        {tabs.map((tab) => (
          <TabButton
            key={tab.key}
            isActive={currentTab === tab.key}
            onClick={() => setCurrentTab(tab.key)}
          >
            {tab.label}
          </TabButton>
        ))}
      </TabNav>
      <TabContent>{children}</TabContent>
    </TabWrapper>
  );
};
const Container = styled.section`
  padding: 32px;
  background-color: #ffffff;
  border-radius: 12px;
`;

const Header = styled.header`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  width: 100%;
`;
const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  color: #00225a;
`;
const TabWrapper = styled.div`
  width: 100%;
`;

const TabNav = styled.div`
  display: flex;
  gap: 8px;
  background-color: transparent;
`;

const TabButton = styled.button<{ isActive: boolean }>`
  padding: 8px 24px;
  background-color: ${({ isActive }) => (isActive ? "#007CDF" : "#F2F6F9")};
  border: none;
  cursor: pointer;
  font-size: 16px;
  font-weight: 600;
  color: ${({ isActive }) => (isActive ? "#ffffff" : "#00225A")};
  border-radius: ${({ isActive }) => (isActive ? "6px" : "10px")};
  transition: all 0.3s ease;
  box-shadow: ${({ isActive }) =>
    isActive ? "0 2px 8px rgba(0, 34, 90, 0.2)" : "none"};

  &:hover {
    transform: translateY(-2px);
  }
`;

const TabContent = styled.div`
  padding: 20px 0;
  animation: fadeIn 0.3s ease-in;

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
`;

export default JobsPage;
