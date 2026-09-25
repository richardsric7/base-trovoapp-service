"use client";

import React from "react";
import styled from "styled-components";
import PrimaryButton from "@/components/PrimaryButton";
import { FaPlus } from "react-icons/fa6";
import Link from "next/link";
import NotificationFilter from "./components/NotificationFilter";
import { SearchBar } from "@/components";
import NotificationTable from "./components/NotificationTable";

const NotificationPage = () => {
  return (
    <PageContainer>
      <HeaderContent>
        <PageTitle>Notification</PageTitle>

        <AddNotificationLink href="notification/addnotification">
          <PrimaryButton
            buttonStyle={{
              width: "auto",
              padding: "10px 20px",
              display: "flex",
              alignItems: "center",
              gap: "4px",
            }}
          >
            <FaPlus /> Add New
          </PrimaryButton>
        </AddNotificationLink>
      </HeaderContent>

      <NotificationStats>
        <SentNotifications>Sent (13)</SentNotifications>
        <DraftNotifications>Draft (2)</DraftNotifications>
      </NotificationStats>

      <FilterContainer>
        <SearchBar 
        />
        <NotificationFilter />
      </FilterContainer>

      <NotificationTable />
    </PageContainer>
  );
};

export default NotificationPage;

const PageContainer = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;

  th {
    color: #828282;
    font-size: 14px;
    font-weight: 500;
  }
  td {
    color: #00225a;
    font-weight: 400;
  }
`;

const HeaderContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const PageTitle = styled.h1`
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;

const AddNotificationLink = styled(Link)`
  text-decoration: none;
`;

const NotificationStats = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const SentNotifications = styled.div`
  background-color: #007cdf;
  color: #ffffff;
  padding: 8px 24px 8px 24px;
  gap: 10px;
  font-size: 14px;
  border-radius: 10px;
`;

const DraftNotifications = styled.div`
  background-color: #f2f6f9;
  color: #00225a;
  padding: 8px 24px 8px 24px;
  gap: 10px;
  font-size: 14px;
  border-radius: 10px;
`;

const FilterContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 25px;
  margin: 10px 0 30px 0;
`;
