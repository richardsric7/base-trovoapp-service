import React from "react";
import styled from "styled-components";

const ActivityCard = () => {
  return (
    <Card>
      <Header>
        <div>
          <strong>Phase 1 foundation works completed</strong>
          <div>Milestone Verification</div>
        </div>

        <Button>Review</Button>
      </Header>
    </Card>
  );
};

export default ActivityCard;

const Card = styled.div`
  background: #f9fafc;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 10px;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
`;

const Button = styled.button`
  background: #2d7ef7;
  color: white;
  border: none;
  padding: 6px 14px;
  border-radius: 6px;
  cursor: pointer;
`;
