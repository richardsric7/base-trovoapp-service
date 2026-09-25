"use client";

import React from "react";
import styled from "styled-components";
import Image from "next/image";
import type { StaticImageData } from "next/image";

interface StepperProps {
  avatar: StaticImageData;
  time: string;
  action: string;
  spanText: string;
  isLast?: boolean; // Optional prop to handle the last step
}

const Stepper: React.FC<StepperProps> = ({
  avatar,
  time,
  action,
  spanText,
  isLast,
}) => {
  return (
    <StepContainer>
      <StepLeft>
        <Circle>
          <StyledImage src={avatar} alt="user-avatar" />
        </Circle>
        {!isLast && <Line />}
      </StepLeft>
      <StepContent>
        <TimeStamp>{time}</TimeStamp>
        <Description>
          {action} <span>{spanText}</span>
        </Description>
      </StepContent>
    </StepContainer>
  );
};

export default Stepper;

const StepContainer = styled.div`
  display: flex;
  align-items: flex-start;
  margin: 20px 0;
  gap: 4px;
`;

const StepLeft = styled.div`
  display: flex;
  flex-direction: column; /* Stack circle and line vertically */
  align-items: center; /* Center items horizontally */
`;

const Circle = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: #f2f6f9;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 5px; /* Space between circle and line */
`;

const StyledImage = styled(Image)`
  border-radius: 50%;
  object-fit: cover;
`;

const Line = styled.div`
  width: 1px;
  height: 40px; /* Adjust height as needed */
  background-color: #e5e5ef;
`;

const StepContent = styled.div`
  display: flex;
  flex-direction: column;
`;

const TimeStamp = styled.p`
  font-size: 12px;
  color: #4f4f4f;
  font-weight: 400;
  line-height: 14.63px;
  margin: 0;
`;

const Description = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 17.07px;
  color: #00225a;

  span {
    font-weight: 600;
  }
`;
