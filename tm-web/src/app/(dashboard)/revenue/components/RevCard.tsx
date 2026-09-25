import React from "react";

import { LuDot } from "react-icons/lu";
import styled from "styled-components";

interface RevCardProps {
  text: string;
  value: string;
  badge: string;
  dotColor?: string;
}

const RevCard: React.FC<RevCardProps> = ({ text, value, badge, dotColor }) => {
  return (
    <Card>
      <CardHeading>
        <LuDot size={30} color={dotColor || "#00225A"} />

        <CardText> {text}</CardText>
      </CardHeading>

      <CardContent>
        <CardValue>{value}</CardValue>
        <CardBadge>{badge}</CardBadge>
      </CardContent>
    </Card>
  );
};

export default RevCard;

const Card = styled.div`
  background-color: #f2f6f9;
  border-radius: 12px;
  padding: 12px;
`;

const CardHeading = styled.div`
  display: flex;
  align-items: center;
  padding-bottom: 10px;
`;

const CardContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const CardText = styled.p`
  font-weight: 500;
  font-size: 14px;
  letter-spacing: 0.1px;
  color: #828282;
`;
const CardValue = styled.p`
  font-weight: 600;
  font-size: 24px;
  letter-spacing: 0.1px;
  color: #00225a;
`;

const CardBadge = styled.div`
  background-color: #ffffff;
  border-radius: 8px;
  padding: 4px 8px;
  color: #004988;
  font-weight: 400;
  font-size: 12px;
`;
