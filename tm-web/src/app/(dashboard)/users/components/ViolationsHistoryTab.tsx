import React, { useMemo, useState } from "react";
import styled from "styled-components";
import { useGetUserSuspensionHistoryQuery } from "@/redux/api/users";
import { IUserSuspensionHistoryEntry } from "@/redux/api/users/interface";

interface ViolationItemProps {
  isActive?: boolean;
}

interface ViolationsHistoryProps {
  userEmail: string;
}

const formatDate = (iso: string) =>
  iso
    ? new Date(iso).toLocaleString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "numeric",
        minute: "2-digit",
      })
    : "";

const ViolationsHistory: React.FC<ViolationsHistoryProps> = ({ userEmail }) => {
  const { data, isLoading } = useGetUserSuspensionHistoryQuery(userEmail, {
    skip: !userEmail,
  });
  const history = useMemo(
    () => [...(data ?? [])].sort((a, b) => b.suspensionDateTime.localeCompare(a.suspensionDateTime)),
    [data]
  );
  const [selected, setSelected] = useState<IUserSuspensionHistoryEntry | null>(null);
  const active = selected ?? history[0] ?? null;

  const suspendCount = history.filter((h) => h.actionType === "SUSPENDED").length;
  const liftCount = history.filter((h) => h.actionType === "ACTIVE").length;

  return (
    <Wrapper>
      <LeftPanel>
        <Header>
          <Heading>Suspension History</Heading>
        </Header>

        <StatsContainer>
          <StatCard>
            <Label>Times suspended</Label>
            <StatValue>{suspendCount}</StatValue>
          </StatCard>
          <StatCard>
            <Label>Times lifted</Label>
            <StatValue>{liftCount}</StatValue>
          </StatCard>
        </StatsContainer>

        {isLoading ? (
          <EmptyText>Loading history...</EmptyText>
        ) : history.length === 0 ? (
          <EmptyText>No suspension history for this user.</EmptyText>
        ) : (
          <ViolationsList>
            {history.map((entry) => (
              <ViolationItem
                key={entry.id}
                onClick={() => setSelected(entry)}
                isActive={active?.id === entry.id}
              >
                <div>
                  <ViolationTitle>
                    {entry.actionType === "SUSPENDED" ? "Suspended" : "Suspension lifted"}
                  </ViolationTitle>
                  <Details>
                    <Timestamp>{formatDate(entry.suspensionDateTime)}</Timestamp>
                    <StatusTag status={entry.actionType}>
                      {entry.actionType === "SUSPENDED" ? "Suspended" : "Active"}
                    </StatusTag>
                  </Details>
                </div>
              </ViolationItem>
            ))}
          </ViolationsList>
        )}
      </LeftPanel>

      <RightPanel>
        {active ? (
          <>
            <Heading>{active.actionType === "SUSPENDED" ? "Suspended" : "Suspension lifted"}</Heading>
            <DetailsContainer>
              <DetailItem>
                <DetailText>Date:</DetailText>
                <DetailValue>{formatDate(active.suspensionDateTime)}</DetailValue>
              </DetailItem>
              <DetailItem>
                <DetailText>Reason:</DetailText>
                <DetailValue>{active.reason}</DetailValue>
              </DetailItem>
              <DetailItem>
                <DetailText>Status:</DetailText>
                <StatusTag status={active.actionType}>
                  {active.actionType === "SUSPENDED" ? "Suspended" : "Active"}
                </StatusTag>
              </DetailItem>
              <DetailItem>
                <DetailText>By:</DetailText>
                <DetailValue>{active.actionPerformedBy}</DetailValue>
              </DetailItem>
            </DetailsContainer>
          </>
        ) : (
          <EmptyText>Select an entry to see its details.</EmptyText>
        )}
      </RightPanel>
    </Wrapper>
  );
};

export default ViolationsHistory;

const Wrapper = styled.section`
  padding: 20px 0;
  display: flex;
  justify-content: space-between;
  gap: 20px;
`;

const LeftPanel = styled.div`
  background-color: #fff;
  border-radius: 24px 24px 0 0;
  padding: 20px;
  min-width: 340px;
`;

const RightPanel = styled.div`
  background-color: #fff;
  border-radius: 24px 24px 0 0;
  padding: 20px;
  flex-grow: 1;
`;

const DetailsContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding-top: 20px;
`;

const Heading = styled.h3`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin-bottom: 10px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const StatsContainer = styled.div`
  display: flex;
  padding: 8px 0;
  border-radius: 12px;
  gap: 8px;
`;

const StatCard = styled.div`
  padding: 8px 12px;
  border-radius: 12px;
  background-color: #f2f6f9;
`;

const Label = styled.p`
  font-size: 14px;
  font-weight: 400;
  color: #00225a;
`;

const StatValue = styled.p`
  font-size: 18px;
  font-weight: 600;
  color: #00225a;
`;

const EmptyText = styled.p`
  color: #828282;
  font-size: 14px;
  padding: 20px 0;
`;

const ViolationsList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

const ViolationItem = styled.div<ViolationItemProps>`
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  padding: 10px;
  border-radius: 8px;
  background-color: ${(props) => (props.isActive ? "#f2f6f9" : "transparent")};

  &:hover {
    background-color: #f2f6f9;
  }
`;

const ViolationTitle = styled.p`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
`;

const Details = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 5px 0;
`;

const Timestamp = styled.p`
  font-size: 12px;
  color: #00225a;
`;

const getStatusColor = (status: string) => {
  switch (status) {
    case "ACTIVE":
      return "#00A859";
    case "SUSPENDED":
      return "#BE3800";
    default:
      return "#00225A";
  }
};
const StatusTag = styled.p<{ status: string }>`
  font-size: 12px;
  color: ${(props) => getStatusColor(props.status)};
  font-weight: 400;
`;
const DetailValue = styled.p`
  color: #00225a;
  font-size: 14px;
  font-weight: 400;
  flex-grow: 1;
  word-break: break-word;
  overflow-wrap: break-word;
  min-height: 20px;
`;

const DetailItem = styled.div`
  display: flex;
  align-items: flex-start;
  column-gap: 20px;
  margin-bottom: 10px;
`;

const DetailText = styled.p`
  font-size: 14px;
  font-weight: 400;
  color: #828282;
  min-width: 100px;
  flex-shrink: 0;
`;
