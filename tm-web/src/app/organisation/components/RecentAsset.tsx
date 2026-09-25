import React from "react";
import styled from "styled-components";

const Card = styled.div`
  background: white;
  padding: 20px;
  border-radius: 10px;
`;

const Title = styled.h3`
  margin-bottom: 16px;
`;

const Row = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  padding: 10px 0;
  border-bottom: 1px solid #f1f1f1;
`;

const Status = styled.span<{ type: string }>`
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  background: ${({ type }) => (type === "completed" ? "#e7f7ee" : "#eef3ff")};
  color: ${({ type }) => (type === "completed" ? "#1f9254" : "#3a7bfd")};
`;

const RecentAssets = () => {
  const assets = new Array(6).fill({
    user: "ATL",
    category: "Real Estate",
    status: "Pending vetting",
  });

  return (
    <Card>
      <Title>Recent Assets Assigned</Title>

      {assets.map((a, i) => (
        <Row key={i}>
          <div>{a.user}</div>
          <div>{a.category}</div>
          <Status type="pending">{a.status}</Status>
        </Row>
      ))}
    </Card>
  );
};

export default RecentAssets;
