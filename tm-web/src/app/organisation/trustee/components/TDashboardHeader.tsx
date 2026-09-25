"use client";
import React from "react";
import styled from "styled-components";

import { useGetOrganizationDetailsDataQuery } from "@/redux/api/org";
import { TrusteeSummary } from "@/redux/api/trustees";
import TStatsGrid from "./TStatsGrid";

interface TDashboardHeaderProps {
  summary?: TrusteeSummary;
}

const TDashboardHeader = ({ summary }: TDashboardHeaderProps) => {
  const { data, isLoading } = useGetOrganizationDetailsDataQuery();

  const orgName = data?.data?.organization?.name ?? "";
  if (isLoading) return <Subtitle>Loading...</Subtitle>;

  return (
    <Wrapper>
      <LeftHeader>
        <Avatar>{orgName.charAt(0).toUpperCase() || "?"}</Avatar>
        <div>
          <Title>{orgName}</Title>
          <Subtitle>
            {" "}
            {data?.data?.organization?.type
              ?.toLowerCase()
              .split("_")
              .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
              .join(" ")}{" "}
            | {data?.data?.organization?.id?.slice(0, 7)}
          </Subtitle>
        </div>
      </LeftHeader>
      <TStatsGrid summary={summary} />
    </Wrapper>
  );
};

export default TDashboardHeader;

const Wrapper = styled.div`
  background: #fff;
  padding: 32px 16px 24px 16px;
  border-radius: 10px;
  margin-bottom: 20px;
`;

const LeftHeader = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 18px;
`;

const Avatar = styled.div`
  width: 38px;
  height: 38px;
  border-radius: 50%;
  background: #007cdf;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
`;

const Title = styled.p`
  font-weight: 600;
  font-size: 24px;
  color: #00225a;
`;

const Subtitle = styled.p`
  font-size: 15px;
  color: #00225a;
  font-weight: 500;
`;
