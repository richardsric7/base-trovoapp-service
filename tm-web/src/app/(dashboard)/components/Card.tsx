"use client";
import Image from "next/image";
import React, { useState } from "react";
import styled from "styled-components";
interface CardDetailsProps {
  text: string;
  heading: string | number;
  img: string;
  hasShadow?: boolean;
}

const Card: React.FC<CardDetailsProps> = ({
  text,
  heading,
  img,
  hasShadow,
}) => {
  return (
    <CardContainer hasShadow={hasShadow}>
      <CardHeader>
        <div>
          <CardText>{text}</CardText>
          <CardTitle>{heading}</CardTitle>
        </div>
        <StyledImage>
          <Image src={img} alt="" />
        </StyledImage>
      </CardHeader>
    </CardContainer>
  );
};

export default Card;

const CardContainer = styled.div<{ hasShadow?: boolean }>`
  width: 100%;
  padding: 8px;
  border-radius: 10px;
  background-color: #ffffff;

  box-shadow: ${({ hasShadow }) =>
    hasShadow ? "0 4px 12px rgba(0, 0, 0, 0.08)" : "none"};
`;

const CardHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 16px;
`;
const CardText = styled.p`
  font-size: 14px;
  font-weight: 400;
  margin: 0;
  letter-spacing: 0.25px;
  text-align: left;
  color: #828282;
`;
const CardTitle = styled.h2`
  font-size: 24px;
  font-weight: 600;
  margin: 0;
  padding-top: 4px;
  letter-spacing: 0.1px;
  text-align: left;
  color: #00225a;
`;
const StyledImage = styled.div`
  width: 40px;
  height: 40px;
  background-color: #f2f6f9;
  display: flex;
  align-items: center;
  justify-content: center;
`;
