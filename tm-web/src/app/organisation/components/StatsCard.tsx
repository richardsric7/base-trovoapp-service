import React from "react";
import { FaArrowRight } from "react-icons/fa6";
import { LuDot } from "react-icons/lu";
import styled from "styled-components";

interface Props {
  label: string;
  value: string | number;
  showFooter?: boolean;
  pending?: string | number;
}

const StatCard: React.FC<Props> = ({ label, value, showFooter, pending = 0 }) => {
  return (
    <Card>
      <Label>{label}</Label>
      <Value>{value}</Value>

      {showFooter && (
        <CardFooter>
          <CardStatus>
            <Dot></Dot>
            <CardText>{pending} Pending</CardText>
          </CardStatus>

          <CardLink>
            View all <FaArrowRight />
          </CardLink>
        </CardFooter>
      )}
    </Card>
  );
};

export default StatCard;

const Card = styled.div`
  background: #f2f6f9;
  padding: 16px;
  border-radius: 10px;
`;

const Label = styled.div`
  font-size: 12px;
  color: #9a9a9a;
`;

const Value = styled.div`
  font-size: 20px;
  font-weight: 600;
  margin-top: 4px;
`;

const CardFooter = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 18px;
`;

const Dot = styled.div`
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #ffcc00;
`;

const CardStatus = styled.div`
  background-color: #ffffff;
  border-radius: 8px;
  padding: 8px;
  gap: 8px;
  font-size: 14px;
  color: #00225a;
  display: flex;
  align-items: center;
`;

const CardText = styled.p`
  font-size: 14px;
  color: #00225a;
`;

const CardLink = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  color: #00225a;
  font-size: 14px;
`;
