import React from "react";
import styled, { keyframes } from "styled-components";
import { HiOutlineSearchCircle } from "react-icons/hi";

interface EmptyStateProps {
  title?: string;
  message?: string;
  icon?: React.ReactNode;
}

const EmptyState = ({
  title = "No Assets Found",
  message = "Try adjusting your filters or search term to find what you're looking for.",
  icon,
}: EmptyStateProps) => {
  return (
    <Container>
      <IconWrapper>
        {icon || <HiOutlineSearchCircle color="#007CDF" size={64} />}
      </IconWrapper>
      <Title>{title}</Title>
      <Message>{message}</Message>
    </Container>
  );
};

export default EmptyState;

const fadeIn = keyframes`
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
`;

const Container = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
  width: 100%;
  animation: ${fadeIn} 0.5s ease-out forwards;
`;

const IconWrapper = styled.div`
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f2f6f9 0%, #e5eef5 100%);
  width: 60px;
  height: 60px;
  border-radius: 50%;
  box-shadow: 0 4px 12px rgba(0, 124, 223, 0.05);
`;

const Title = styled.h3`
  font-size: 22px;
  font-weight: 700;
  color: #00225a;
  margin-bottom: 12px;
  letter-spacing: -0.2px;
`;

const Message = styled.p`
  font-size: 15px;
  color: #6c757d;
  max-width: 360px;
  line-height: 1.6;
`;
