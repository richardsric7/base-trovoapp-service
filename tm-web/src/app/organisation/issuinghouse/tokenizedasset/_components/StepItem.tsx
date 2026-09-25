import React from "react";
import { BiCheck } from "react-icons/bi";
import { FiClock } from "react-icons/fi";
import styled from "styled-components";

export interface StepProps {
  action: string;
  date: string;
  subset?: string[];
}
interface StepCircleStyleProps {
  $active: boolean;
  $disapproved?: boolean;
  $processing?: boolean;
}

export interface StepItemProps {
  step: StepProps;
  index: number;
  currentStep: number;
  processingStep: number | null;
  isDueDiligenceDisapproved: boolean;
}

const StepItem: React.FC<StepItemProps> = ({
  step,
  index,
  currentStep,
  processingStep,
  isDueDiligenceDisapproved,
}) => {
  let icon = null;

  if (
    index === 4 &&
    processingStep === 4 &&
    currentStep === 4 &&
    !isDueDiligenceDisapproved
  ) {
    icon = <FiClock />;
  } else if (
    (index === 4 && isDueDiligenceDisapproved) ||
    index <= currentStep
  ) {
    icon = <BiCheck />;
  }

  return (
    <StepContainer>
      <StepCircleContainer>
        <StepCircle
          $active={
            index <= currentStep || (index === 4 && isDueDiligenceDisapproved)
          }
          $disapproved={index === 5 && isDueDiligenceDisapproved}
          $processing={
            index === 4 &&
            processingStep === 4 &&
            currentStep === 4 &&
            !isDueDiligenceDisapproved
          }
        >
          {icon}
        </StepCircle>

        {index < 10 && <Line />}
      </StepCircleContainer>

      <ContentContainer>
        <TextContainer>
          <StepText $isDisapproved={index === 5 && isDueDiligenceDisapproved}>
            {index === 5 && isDueDiligenceDisapproved
              ? "Tokenization Failed"
              : step.action}
          </StepText>

          <StepDate>{step.date}</StepDate>
        </TextContainer>
      </ContentContainer>
    </StepContainer>
  );
};

export default StepItem;

const StepContainer = styled.div`
  display: flex;
  flex-direction: row;
  gap: 8px;
  margin-bottom: 16px;
  width: 100%;
  min-width: 0;
`;

const StepCircleContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const StepCircle = styled.div<StepCircleStyleProps>`
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: ${(props) =>
    props.$disapproved
      ? "#FF0000"
      : props.$processing
        ? "#007CDF"
        : props.$active
          ? "#4CAF50"
          : "#D9D9D9"};
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
`;

const ContentContainer = styled.div`
  display: flex;
  flex: 1;
  min-width: 0;
  flex-direction: column;
`;

const TextContainer = styled.div`
  display: flex;
  width: 100%;
  min-width: 0;
  flex-direction: column;
  align-items: flex-start;
`;

const StepText = styled.p<{ $isDisapproved?: boolean }>`
  width: 100%;
  margin: 0;
  font-size: 14px;
  white-space: normal;
  overflow-wrap: anywhere;
  word-break: break-word;
  color: ${(props) => (props.$isDisapproved ? "#FF0000" : "#00225a")};
  font-weight: 500;
  line-height: 18px;
  text-align: left;
`;

const StepDate = styled.p`
  width: 100%;
  margin: 4px 0 0;
  overflow-wrap: anywhere;
  word-break: break-word;
  font-size: 12px;
  line-height: 16px;
  font-weight: 400;
  color: #828282;
`;

const StepSubSet = styled.p`
  font-size: 12px;
  font-weight: 400;
  color: #828282;
  white-space: pre-line;
  margin-top: 4px;
`;

const Line = styled.div`
  flex-grow: 1;
  height: 40px;
  background-color: #7c7c7c99;
  margin: 8px 0;
  width: 1px;
`;
