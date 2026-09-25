"use client";
import React from "react";
import styled from "styled-components";
import { SideBar } from "./components/SideBar";
import { NavBar } from "./components/NavBar";

export default function DashboardLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <LayoutContainer>
      <LayoutWrapper>
        <SideBar />
        <ContentWrapper>
          <NavBar />
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
  margin-left: 252px;
`;
const ContentWrapper = styled.section`
  background-color: #f5f5f5;
  padding: 6px;
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 100vh;
  overflow-x: hidden;
  min-width: 0;
`;
const ContentContainer = styled.div`
  margin-top: 100px;
  margin-left: 10px;
  margin-right: 10px;
  height: 100vh;
  border-radius: 24px;
  min-height: fit-content;
  // width: 100%;
  // max-width: 100%;
  // overflow-x: hidden;
  // box-sizing: border-box;
  // min-width: 0;
`;
