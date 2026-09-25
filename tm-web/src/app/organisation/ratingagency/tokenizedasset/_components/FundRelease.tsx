import styled from "styled-components";

import { useState } from "react";
import { useRouter } from "next/navigation";
import FundFilter from "@/app/organisation/assetcustodian/fundmanagement/_components/LegacyFundFilter";

const activities = [
  {
    id: 1,
    title: "ATL",
    subtitle: "Atlantis",
    requestId: "FR-4491",
    description: "1,000,000.00 NGN",
    purpose: "Milestone Verification",
    requestedBy: "Florencez",
    status: "Pending Review",
    date: "6th Jun, 2025",
  },
  {
    id: 2,

    title: "ATL",
    subtitle: "Atlantis",
    requestId: "FR-4491",
    description: " ₦120,000,000 ",
    purpose: "Milestone Payment",
    requestedBy: "Florencez",
    status: "Pending Review",
    date: "6th Jun, 2025",
  },
  {
    id: 3,
    requestId: "FR-4491",
    title: "ATL",
    subtitle: "Atlantis",
    description: "₦120,000,000",
    purpose: "Income Distribution",
    requestedBy: "Florencez",
    status: "Pending Review",
    date: "6th Jun, 2025",
  },
];

const FundRelease = () => {
  const [filters, setFilters] = useState({});
  const router = useRouter();
  const getRoute = (item: (typeof activities)[number]) => {
    switch (item.purpose) {
      case "Income Distribution":
        return `/organisation/fundmanagement/incomedistribution/${item.id}`;

      default:
        return `/organisation/fundmanagement/${item.id}`;
    }
  };

  return (
    <Container>
      <Header>
        <Title>Fund Release Requests</Title>
        <FundFilter filters={filters} onApply={setFilters} />
      </Header>
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

            <ReviewButton onClick={() => router.push(getRoute(item))}>
              Review
            </ReviewButton>
          </TopRow>

          <Description>{item.description}</Description>

          <DetailsRow>
            <Detail>
              <Label>Request ID</Label>
              <Value>{item.requestId}</Value>
            </Detail>

            <Divider />

            <Detail>
              <Label>Purpose</Label>
              <Value>{item.purpose}</Value>
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

export default FundRelease;

const Container = styled.div`
  background: #fff;
  padding: 24px;
  border-radius: 16px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
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

const Description = styled.p`
  margin-top: 10px;

  color: #00225a;

  font-weight: 600;
  font-size: 16px;

  line-height: 140%;
  letter-spacing: 0.1px;
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
  min-width: 144px;
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
