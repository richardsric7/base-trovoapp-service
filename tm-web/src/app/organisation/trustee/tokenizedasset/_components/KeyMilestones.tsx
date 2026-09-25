"use client";

import styled from "styled-components";
import { BsCheckCircle } from "react-icons/bs";
import { useState } from "react";
import { MilestoneCompletionModal } from "./MilestoneCompletionModal";

const milestones = [
  {
    title: "Land acquisition",
    date: "April, 2026",
    status: "Completed",
    type: "completed",
    action: "view",
  },
  {
    title: "Permits",
    date: "September, 2026",
    status: "Pending Authorization",
    type: "pending-auth",
    action: "view",
  },
  {
    title: "Permits",
    date: "December, 2026",
    status: "Pending",
    type: "pending",
    action: "complete",
  },
  {
    title: "Permits",
    date: "December, 2026",
    status: "Pending",
    type: "pending",
    action: "complete",
  },
  {
    title: "Permits",
    date: "December, 2026",
    status: "Pending",
    type: "pending",
    action: "complete",
  },
];

const KeyMilestones = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedMilestone, setSelectedMilestone] = useState<any>(null);

  return (
    <Container>
      <Title>Key Milestones & Dates</Title>

      <Wrapper>
        {milestones.map((item, index) => (
          <Card key={index}>
            <Left>
              <MilestoneTitle>{item.title}</MilestoneTitle>
              <Date>{item.date}</Date>

              <Status $type={item.type}>
                <Dot $type={item.type} />
                {item.status}
              </Status>
            </Left>

            <Right>
              {item.action === "view" ? (
                <OutlineBtn>View details</OutlineBtn>
              ) : (
                <PrimaryBtn
                  onClick={() => {
                    setSelectedMilestone(item);
                    setIsModalOpen(true);
                  }}
                >
                  <BsCheckCircle size={14} />
                  Mark as completed
                </PrimaryBtn>
              )}
            </Right>
          </Card>
        ))}
      </Wrapper>

      <MilestoneCompletionModal
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        milestone={selectedMilestone}
      />
    </Container>
  );
};

export default KeyMilestones;

const Container = styled.div`
  background: #ffffff;
  border-radius: 20px;
`;

const Title = styled.h2`
  font-size: 18px;
  font-weight: 600;
  color: #00225a;
  margin-bottom: 16px;
`;

const Wrapper = styled.div`
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 14px;
`;

const Card = styled.div`
  background: #f2f6f9;
  border-radius: 12px;
  padding: 16px 18px;
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const Left = styled.div`
  display: flex;
  flex-direction: column;
`;

const Right = styled.div`
  display: flex;
  align-items: center;
`;

const MilestoneTitle = styled.p`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
`;

const Date = styled.p`
  font-size: 12px;
  color: #9ca3af;
  margin-top: 4px;
`;

const Status = styled.div<{ $type: string }>`
  margin-top: 8px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  border-radius: 999px;
  padding: 4px 10px;
  background: ${(p) => (p.$type === "completed" ? "#fff" : "#ffff")};
  color: ${(p) => (p.$type === "completed" ? "#191919" : "#191919")};
`;

const Dot = styled.div<{ $type: string }>`
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: ${(p) => (p.$type === "completed" ? "#22c55e" : "#f59e0b")};
`;

const OutlineBtn = styled.button`
  border: 1px solid #007cdf;
  color: #007cdf;
  background: transparent;
  border-radius: 8px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
  font-family: inherit;
`;

const PrimaryBtn = styled.button`
  background: #007cdf;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 6px 12px;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-family: inherit;
`;
