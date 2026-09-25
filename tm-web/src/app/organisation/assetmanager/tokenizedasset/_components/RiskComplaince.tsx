"use client";

import styled from "styled-components";
import { AiOutlineClockCircle } from "react-icons/ai";

const RiskAndCompliance = () => {
  return (
    <Wrapper>
      <Title>Risk and Compliance</Title>

      <Card>
        <Item>
          <Label>Risk Assessment Score</Label>

          <RiskRow>
            <RiskBadge>LOW</RiskBadge>
            <Score>18/100</Score>
          </RiskRow>
        </Item>

        <Item>
          <Label>Compliance Status</Label>

          <Status>
            <AiOutlineClockCircle size={14} />
            Pending
          </Status>
        </Item>

        <Item>
          <Label>Last Audit</Label>
          <Value>20 Feb, 2026</Value>
        </Item>
      </Card>
    </Wrapper>
  );
};

export default RiskAndCompliance;

const Wrapper = styled.div`
  margin-bottom: 32px;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin-bottom: 16px;
`;

const Card = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  background: #ffffff;
  border: 1px solid #e6e6e6;
  border-radius: 12px;
  padding: 20px;
`;

const Item = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.p`
  font-size: 13px;
  color: #828282;
`;

const Value = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
`;

const RiskRow = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;

const RiskBadge = styled.div`
  background: #ffe5d9;
  color: #ff6b3d;
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 6px;
`;

const Score = styled.div`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
`;

const Status = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
  background: #e9f2ff;
  color: #007cdf;
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 6px;
  width: fit-content;
`;
