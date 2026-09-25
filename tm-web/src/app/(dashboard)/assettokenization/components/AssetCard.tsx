"use client";
import React from "react";
import styled from "styled-components";

interface AssetCardProps {
  label?: string;
  // value?: string | number;
  value?: React.ReactNode;
}

const AssetCard: React.FC<AssetCardProps> = ({ label, value }) => {
  return (
    <Card>
      <Label>{label}</Label>
      <Value>{value}</Value>
    </Card>
  );
};

export default AssetCard;

const Card = styled.div`
  display: flex;
  flex-direction: column;
`;

const Label = styled.p`
  font-size: 14px;
  font-weight: 400;
  color: #828282;
  margin: 0;
`;

const Value = styled.div`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  margin: 8px 0;
  overflow-wrap: break-word;
  word-break: break-word;
  white-space: normal;
`;
