"use client";

import React from "react";
import { BiCheck } from "react-icons/bi";
import styled from "styled-components";

export interface StepProps {
  action: string;
  date: string;
}

export interface StepItemProps {
  step: StepProps;
  index: number;
  currentStep: number;
  stepsLength: number;
}

function StepperItems({
  step,
  index,
  currentStep,
  stepsLength,
}: StepItemProps) {
  // ICON LOGIC (matches design)
  let icon = null;

  if (index <= currentStep) {
    icon = <BiCheck size={14} />;
  }

  return (
    <StepContainer>
      <StepCircleContainer>
        <StepCircle $active={index <= currentStep}>{icon}</StepCircle>

        {/* line */}
        {index !== stepsLength - 1 && <Line />}
      </StepCircleContainer>

      <ContentContainer>
        <TextContainer>
          <StepText $active={index <= currentStep}>{step.action}</StepText>

          {/* ONLY active step shows date */}
          {index === currentStep && <StepDate>{step.date}</StepDate>}
        </TextContainer>
      </ContentContainer>
    </StepContainer>
  );
}

export default React.memo(StepperItems);

const StepContainer = styled.div`
  display: flex;
  align-items: flex-start;
  gap: 10px;
  margin-bottom: 20px;
`;

const StepCircleContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const StepCircle = styled.div<{ $active: boolean }>`
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: ${(p) => (p.$active ? "#22C55E" : "#D1D5DB")};
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const Line = styled.div`
  width: 2px;
  height: 32px;
  background: #e5e7eb;
  margin-top: 6px;
`;

const ContentContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

const TextContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

const StepText = styled.p<{ $active: boolean }>`
  font-size: 14px;
  font-weight: 500;
  color: ${(p) => (p.$active ? "#00225a" : "#9ca3af")};
`;

const StepDate = styled.p`
  font-size: 12px;
  color: #9ca3af;
  margin-top: 2px;
`;
