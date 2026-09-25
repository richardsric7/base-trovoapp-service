"use client";

import React from "react";
import styled from "styled-components";

const activities = [
  {
    title: "ATL",
    subtitle: "Atlantis",
    description: "Custody confirmation completed for tokenized asset",
    activity: "Custody Confirmation",
    requestedBy: "Florencez",
    status: "Completed",
    date: "6th Jun, 2025",
  },
];

const OrgActivity = () => {
  return (
    <Container>
      <Title>Recent Activity</Title>

      {activities.map((item, index) => (
        <ActivityCard key={index}>
          <TopRow>
            <LeftInfo>
              <Avatar />
              <ProjectInfo>
                <ProjectName>{item.title}</ProjectName>
                <ProjectSub>{item.subtitle}</ProjectSub>
              </ProjectInfo>
            </LeftInfo>
          </TopRow>

          <Description>{item.description}</Description>

          <DetailsRow>
            <Detail>
              <Label>Activity</Label>
              <Value>{item.activity}</Value>
            </Detail>

            <Divider />

            <Detail>
              <Label>Requested by</Label>
              <User>
                <UserAvatar />
                {item.requestedBy}
              </User>
            </Detail>

            <Divider />

            <Detail>
              <Label>Status</Label>
              <Status>
                <StatusDot />
                {item.status}
              </Status>
            </Detail>

            <Divider />

            <Detail>
              <Label>Date</Label>
              <Value>{item.date}</Value>
            </Detail>
          </DetailsRow>
        </ActivityCard>
      ))}
    </Container>
  );
};

export default OrgActivity;

const Container = styled.div`
  background: #fff;
  padding: 24px;
  border-radius: 16px;
`;

const Title = styled.h3`
  font-size: 18px;
  font-weight: 600;
  color: #1a2b49;
  margin-bottom: 20px;
`;

const ActivityCard = styled.div`
  background: #eef2f6;
  padding: 18px;
  border-radius: 12px;
  margin-bottom: 16px;
`;

const TopRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const LeftInfo = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #d9b48f;
`;

const ProjectInfo = styled.div`
  display: flex;
  flex-direction: column;
`;

const ProjectName = styled.div`
  font-weight: 600;
  color: #1a2b49;
`;

const ProjectSub = styled.div`
  font-size: 12px;
  color: #6b7280;
`;

const Description = styled.div`
  margin-top: 10px;
  font-size: 14px;
  color: #00225a;
`;

const DetailsRow = styled.div`
  margin-top: 14px;
  background: #ffffff;
  padding: 14px;
  border-radius: 10px;
  display: flex;
  align-items: center;
`;

const Detail = styled.div`
  display: flex;
  flex-direction: column;
  min-width: 160px;
`;

const Label = styled.div`
  font-size: 12px;
  color: #8a94a6;
`;

const Value = styled.div`
  font-size: 14px;
  color: #00225a;
  margin-top: 4px;
  font-weight: 500;
`;

const Divider = styled.div`
  width: 1px;
  height: 32px;
  background: #e5e7eb;
  margin: 0 20px;
`;

const User = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
  font-size: 14px;
  color: #1a2b49;
`;

const UserAvatar = styled.div`
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #60a5fa;
`;

const Status = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 4px;
  font-size: 14px;
  color: #1a2b49;
`;

const StatusDot = styled.div`
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #facc15;
`;
