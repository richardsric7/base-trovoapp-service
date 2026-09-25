"use client";
import styled from "styled-components";

interface TokenizedAssetStatProps {
  title: string;
  value: number;
  icon: React.ReactNode;
  color: string;
  bgColor: string;
  amount?: string;
}

const TokenizedAssetStat = ({
  title,
  value,
  icon,
  color,
  bgColor,
  amount,
}: TokenizedAssetStatProps) => {
  return (
    <StatCard>
      <StatHeader>
        <div>
          <StatText>{title}</StatText>
          <StatValue $color={color}>{value}</StatValue>{" "}
        </div>
        <IconWrapper $bgColor={bgColor}>{icon}</IconWrapper>
      </StatHeader>

      <StatAmount $color={color}>{amount}</StatAmount>
    </StatCard>
  );
};

export default TokenizedAssetStat;

const StatCard = styled.div`
  border-radius: 10px;
  background-color: #ffffff;
  padding: 8px 16px;

  flex: 1;
  max-width: 100%;
`;

const StatHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
`;

const StatText = styled.p`
  font-size: 14px;
  font-weight: 400;
  margin: 0;
  letter-spacing: 0.25px;
  text-align: left;
  color: #828282;

  @media (max-width: 1500px) {
    font-size: 12px;
  }
`;

const StatValue = styled.h2<{ $color?: string }>`
  font-size: 24px;
  font-weight: 600;
  margin: 0;
  color: ${({ $color }) => $color || "#00225A"};

  @media (max-width: 1500px) {
    font-size: 20px;
  }
`;

const StatAmount = styled.h2<{ $color?: string }>`
  font-size: 18px;
  font-weight: 500;
  margin-top: 24px;
  color: ${({ $color }) => $color || "#00225A"};
  @media (max-width: 1500px) {
    font-size: 14px;
  }
`;

const IconWrapper = styled.div<{ $bgColor?: string }>`
  width: 40px;
  height: 40px;
  background-color: ${({ $bgColor }) => $bgColor || "#F2F6F9"};
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
`;
