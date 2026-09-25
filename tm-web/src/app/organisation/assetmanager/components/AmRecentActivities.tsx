import React from "react";
import styled from "styled-components";
import { AmRecentActivity } from "@/redux/api/assetManager";
import { getInitials } from "@/utils/getInitials";

interface RecentActivityProps {
  activities?: AmRecentActivity[];
}

const RecentActivity = ({ activities = [] }: RecentActivityProps) => {
  return (
    <Container>
      <Title>Recent Activity</Title>

      {activities.length === 0 ? (
        <EmptyState>No data available</EmptyState>
      ) : (
        activities.slice(0, 5).map((item, index) => (
          <ActivityCard key={item.id || index}>
            <TopRow>
              <LeftInfo>
                <Avatar>
                  {(item.actor_organization || item.title || "A")
                    .charAt(0)
                    .toUpperCase()}
                </Avatar>
                <ProjectInfo>
                  <ProjectName>{item.title || "Activity"}</ProjectName>
                  <ProjectSub>
                    {item.actor_organization || "Organization"}
                  </ProjectSub>
                </ProjectInfo>
              </LeftInfo>
            </TopRow>

            <DetailsRow>
              <Detail>
                <Label>Activity Type</Label>
                <Value>{item.type || item.entity_type || "N/A"}</Value>
              </Detail>

              <Divider />

              <Detail>
                <Label>Actor</Label>
                <User>
                  <UserAvatar aria-hidden="true">
                    {getInitials(
                      item.actor_name?.trim() || item.actor_organization,
                    )}
                  </UserAvatar>
                  {item.actor_name
                    ? `${item.actor_name} (${item.actor_organization || ""})`
                    : item.actor_organization || "N/A"}
                </User>
              </Detail>

              <Divider />

              <Detail>
                <Label>Status</Label>
                <Status>
                  <StatusDot />
                  {item.status || "N/A"}
                </Status>
              </Detail>

              <Divider />

              <Detail>
                <Label>Date</Label>
                <Value>
                  {item.created_at
                    ? new Date(item.created_at).toLocaleString(undefined, {
                        dateStyle: "medium",
                        timeStyle: "short",
                      })
                    : "N/A"}
                </Value>
              </Detail>
            </DetailsRow>
          </ActivityCard>
        ))
      )}
    </Container>
  );
};

export default RecentActivity;

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

const EmptyState = styled.div`
  padding: 40px 0;
  text-align: center;
  color: #8a94a6;
  font-size: 14px;
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
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  color: #fff;
  background: #007cdf;
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

const ReviewButton = styled.button`
  background: #007cdf;
  border: none;
  padding: 8px 24px;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  cursor: pointer;
  font-weight: 600;
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
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 10px;
  font-weight: 600;
  color: #191919;
  border-radius: 50%;
  background: #acd1ef;
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
