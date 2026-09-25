"use client";
import React from "react";
import styled from "styled-components";
import { Sidebar } from "./components/Sidebar";
import { Navbar } from "./components/Navbar";

export default function ExternalApiClientsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <LayoutContainer>
      <Sidebar />
      <ContentWrapper>
        <Navbar />
        <ContentContainer>{children}</ContentContainer>
      </ContentWrapper>
    </LayoutContainer>
  );
}

const LayoutContainer = styled.main`
  display: flex;
  min-height: 100vh;
`;

const ContentWrapper = styled.div`
  flex: 1;
  margin-left: 252px;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
`;

const ContentContainer = styled.div`
  margin-top: 90px;
  padding: 32px;
  flex: 1;
`;
