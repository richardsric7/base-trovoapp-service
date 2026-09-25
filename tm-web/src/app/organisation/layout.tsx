"use client";
import React from "react";
import styled from "styled-components";
import OrgSideBar from "./components/OrgSideBar";
import OrgNavBar from "./components/OrgNavBar";

export default function OrganizationLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <LayoutContainer>
      <LayoutWrapper>
        <OrgSideBar />
        <ContentWrapper>
          <OrgNavBar />
          <ContentContainer>{children}</ContentContainer>
        </ContentWrapper>
      </LayoutWrapper>
    </LayoutContainer>
  );
}
const LayoutContainer = styled.main`
  min-height: 100vh;
  display: flex;
  flex-direction: column;
`;
const LayoutWrapper = styled.div`
  display: flex;
  flex: 1;
`;
const ContentWrapper = styled.section`
  background-color: #f5f5f5;
  padding: 6px;
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 100vh;
`;
const ContentContainer = styled.div`
  margin-top: 100px;
  height: 100vh;
  align-self: end;
  width: calc(100% - 260px);
  border-radius: 24px;
  min-height: fit-content;
`;
