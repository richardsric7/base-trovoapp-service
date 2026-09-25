import React from "react";
import styled from "styled-components";

interface StatsProps {
  text: string;
  count: number | string;
  dotColor: string;
}

const OrgTokenizedAssetsStats = ({ text, count, dotColor }: StatsProps) => {
  return (
    <Card>
      <FlexText>
        <Dot color={dotColor} />
        <Text>{text}</Text>
      </FlexText>
      <Stats>{count}</Stats>
    </Card>
  );
};

export default OrgTokenizedAssetsStats;

const Card = styled.div`
  background-color: #f2f6f9;
  border-radius: 10px;
  padding: 16px;
  width: 100%;
`;

const Text = styled.p`
  color: #828282;
`;

const Stats = styled.p`
  color: #00225a;
  font-weight: 600;
  font-size: 24px;
  line-height: 24px;
  letter-spacing: 0.1px;
`;

const FlexText = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const Dot = styled.div<{ color: string }>`
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: ${({ color }) => color};
`;
