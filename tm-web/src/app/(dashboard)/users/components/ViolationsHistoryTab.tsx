import { SearchBar } from "@/components";
import React, { useState } from "react";
import { FaDotCircle } from "react-icons/fa";
import { FaAngleRight, FaPlus } from "react-icons/fa6";
import styled from "styled-components";
import TransactionsFilter from "./TransactionsFilter";
import { LuDot } from "react-icons/lu";
import SuspendUserModal from "./SuspendUserModal";
import ViolationsFilter from "./ViolationsFilter";

interface ViolationItemProps {
  isActive?: boolean;
}

const ViolationsHistory = () => {
  const violationsData = [
    {
      id: 1,
      title: "Fraudulent Token Sale",
      status: "Suspension Active",
      date: "23 Sep 2023, 8:17 AM",
      description:
        "Created fake asset tokens and attempted to sell them to other users.",
      admin: {
        name: "Tom Anaja",
        role: "Admin",
      },
    },
    {
      id: 2,
      title: "Payment Fraud",
      status: "Suspension Resolved",
      date: "15 Aug 2023, 10:45 AM",
      description: "Attempted unauthorized payment transactions.",
      admin: {
        name: "Sarah Johnson",
        role: "Senior Admin",
      },
    },
    {
      id: 3,
      title: "Account Manipulation",
      status: " Not Suspended",
      date: "05 Oct 2023, 3:22 PM",
      description:
        "Unauthorized modifications to account settings and permissions.",
      admin: {
        name: "Mike Rodriguez",
        role: "Security Manager",
      },
    },
  ];

  const [selectedViolation, setSelectedViolation] = useState(violationsData[0]);
  const [isOpen, setIsOpen] = useState<boolean>(false);

  return (
    <Wrapper>
      <LeftPanel>
        <Header>
          <Heading>Violation History</Heading>
          <AddButton onClick={() => setIsOpen(true)}>
            <FaPlus color="#fff" />
          </AddButton>
          <SuspendUserModal
            isOpen={isOpen}
            setIsOpen={setIsOpen}
            userEmail={""}
          />
        </Header>

        <StatsContainer>
          <StatCard>
            <StatHeader>
              <LuDot color="#BE3800" size={30} />
              <Label>Suspensions</Label>
            </StatHeader>
            <StatValue>8</StatValue>
          </StatCard>
          <StatCard>
            <StatHeader>
              <LuDot color="#00A859" size={30} />
              <Label>Suspensions Resolved</Label>
            </StatHeader>
            <StatValue>5</StatValue>
          </StatCard>
          <StatCard>
            <StatHeader>
              <LuDot color="#007CDF" size={30} />
              <Label>No Suspensions</Label>
            </StatHeader>
            <StatValue>2</StatValue>
          </StatCard>
        </StatsContainer>

        <SearchAndFilter>
          <SearchBar customWidth="320px" />
          <ViolationsFilter />
        </SearchAndFilter>

        <ViolationsList>
          {violationsData.map((violation) => (
            <ViolationItem
              key={violation.id}
              onClick={() => setSelectedViolation(violation)}
              isActive={selectedViolation.id === violation.id}
            >
              <div>
                <ViolationTitle>{violation.title}</ViolationTitle>
                <Details>
                  <Timestamp>{violation.date}</Timestamp>
                  <Line></Line>
                  <StatusTag status={violation.status}>
                    {violation.status}
                  </StatusTag>
                </Details>
              </div>
              {selectedViolation.id === violation.id ? "" : <FaAngleRight />}
            </ViolationItem>
          ))}
        </ViolationsList>
      </LeftPanel>

      <RightPanel>
        <Heading>{selectedViolation.title}</Heading>
        <DetailsContainer>
          <DetailItem>
            <DetailText>Date:</DetailText>
            <DetailValue>{selectedViolation.date}</DetailValue>
          </DetailItem>
          <DetailItem>
            <DetailText>Description:</DetailText>
            <DetailValue>{selectedViolation.description}</DetailValue>
          </DetailItem>
          <DetailItem>
            <DetailText>Status:</DetailText>
            <StatusTag status={selectedViolation.status}>
              {selectedViolation.status}
            </StatusTag>
          </DetailItem>
          <DetailItem>
            <DetailText>By:</DetailText>
            <DetailValueContainer>
              <Avatar></Avatar>
              <DetailContent>
                <DetailValue2>{selectedViolation.admin.name}</DetailValue2>
                <DetailText>{selectedViolation.admin.role}</DetailText>
              </DetailContent>
            </DetailValueContainer>
          </DetailItem>
        </DetailsContainer>
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
  // width: 80%;
`;

const RightPanel = styled.div`
  background-color: #fff;
  border-radius: 24px 24px 0 0;
  padding: 20px;
  // width: 100%;
`;

const DetailsContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding-top: 20px;
`;

const AddButton = styled.div`
  background-color: #007cdf;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  cursor: pointer;
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
  justify-content: space-between;
  padding: 8px 0;
  border-radius: 12px;
  gap: 8px;
`;

const StatHeader = styled.div`
  display: flex;
  align-items: center;
`;

const StatCard = styled.div`
  padding: 8px 4px;
  border-radius: 12px;
  background-color: #f2f6f9;
  width: 193px;
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
  display: inline-block;
  margin-left: 30px;
`;

const SearchAndFilter = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
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
  justify-content: space-between;
  padding: 5px 0;
`;

const Timestamp = styled.p`
  font-size: 12px;
  color: #00225a;
`;

const getStatusColor = (status: string) => {
  switch (status) {
    case "Not Suspended":
      return "#00225A";
    case "Suspension Resolved":
      return "#00A859";
    case "Suspension Active":
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

const DetailValue2 = styled.p`
  color: #00225a;
  font-size: 14px;
  font-weight: 400;
`;

const DetailItem = styled.div`
  display: flex;
  align-items: flex-start; 
  justify-content: 
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
const DetailValueContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;

const DetailContent = styled.div``;

const Line = styled.div`
  width: 13px;
  height: 1px;
  background-color: #bdbdbd;
  transform: rotate(-90deg);
`;
const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;
