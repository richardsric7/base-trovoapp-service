import PrimaryButton from "@/components/PrimaryButton";
import React from "react";
import styled from "styled-components";

interface CardFundProps {
  title: string;
  amount: string;
}

const FundCard = ({ title, amount }: CardFundProps) => {
  return (
    <CardContainer>
      <div>
        <Title>{title}</Title>
        <Amount>{amount}</Amount>
      </div>
      {/* <ViewButton>View Statement</ViewButton> */}
    </CardContainer>
  );
};

export default FundCard;

const CardContainer = styled.div`
  background-color: #ffffff;
  border-radius: 10px;
  padding: 12px 16px;

  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 56px;
`;

const Title = styled.p`
  font-weight: 400;
  font-size: 14px;
  line-height: 20px;
  letter-spacing: 0.25px;
  color: #828282;
`;

const Amount = styled.p`
  font-weight: 600;
  font-style: SemiBold;
  font-size: 24px;
  line-height: 24px;
  color: #00225a;
`;

const ViewButton = styled.button`
  background-color: #007cdf;
  color: #fff;
  width: 134px;
  height: 32px;
  gap: 10px;
  border-radius: 8px;
  padding: 6px 12px;
  border: none;
  font-size: 12px;
  font-family: inherit;
  cursor: pointer;

  margin-top: auto;
`;
