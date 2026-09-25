"use client";
import Image from "next/image";
import React from "react";
import styled from "styled-components";
import { FiArrowUpRight, FiChevronDown } from "react-icons/fi";
interface CardDetailsProps {
  text: string;
  heading: string | number;
  img: string | React.ReactNode;
  hasShadow?: boolean;
  trendValue?: string;
  isPositive?: boolean;
  suffix?: string;
  period?: string;
}

const CardStats = ({
  text,
  heading,
  img,
  hasShadow,
  trendValue,
  isPositive = true,
  suffix,
  period,
}: CardDetailsProps) => {
  return (
    <CardContainer hasShadow={hasShadow}>
      <CardHeader>
        <LeftContainer>
          <CardText>{text}</CardText>
          <CardTitle>
            {heading} {suffix && <Suffix>{suffix}</Suffix>}
          </CardTitle>
        </LeftContainer>
        <StyledImage>
          {typeof img === "string" ? (
            <Image src={img} alt="" width={24} height={24} />
          ) : (
            img
          )}
        </StyledImage>
      </CardHeader>
      {(trendValue || period) && (
        <CardFooter>
          {trendValue && (
            <TrendContainer isPositive={isPositive}>
              {trendValue} <FiArrowUpRight />
            </TrendContainer>
          )}
          {period && (
            <PeriodDropdown>
              {period} <FiChevronDown />
            </PeriodDropdown>
          )}
        </CardFooter>
      )}
    </CardContainer>
  );
};

export default CardStats;

const CardContainer = styled.div<{ hasShadow?: boolean }>`
  width: 100%;
  padding: 24px;
  border-radius: 12px;
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
  gap: 20px;
  box-shadow: ${({ hasShadow }) =>
    hasShadow ? "0 4px 20px rgba(0, 0, 0, 0.05)" : "none"};
  border: 1px solid #f2f2f2;
`;

const CardHeader = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
`;

const LeftContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const CardText = styled.p`
  font-size: 14px;
  font-weight: 400;
  margin: 0;
  color: #828282;
`;

const CardTitle = styled.h2`
  font-size: 28px;
  font-weight: 700;
  margin: 0;
  color: #00225a;
  display: flex;
  align-items: baseline;
  gap: 8px;
`;

const Suffix = styled.span`
  font-size: 18px;
  font-weight: 600;
  color: #00225a;
`;

const StyledImage = styled.div`
  width: 48px;
  height: 48px;
  background-color: #f2f6f9;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  img {
    width: 24px;
    height: 24px;
  }
`;

const CardFooter = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: auto;
`;

const TrendContainer = styled.div<{ isPositive?: boolean }>`
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  font-weight: 600;
  color: ${({ isPositive }) => (isPositive ? "#219653" : "#eb5757")};

  svg {
    font-size: 16px;
  }
`;

const PeriodDropdown = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background-color: #f2f6f9;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  color: #4f4f4f;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #e8eff5;
  }

  svg {
    font-size: 14px;
    color: #4f4f4f;
  }
`;
