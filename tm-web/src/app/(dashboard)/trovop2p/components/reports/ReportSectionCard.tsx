import React from "react";
import styled from "styled-components";

const ReportSectionCard = ({
  title,
  subtitle,
  action,
  children,
}: {
  title: string;
  subtitle?: string;
  action?: React.ReactNode;
  children: React.ReactNode;
}) => {
  return (
    <Container>
      <Header>
        <div>
          <Title>{title}</Title>
          {subtitle && <Subtitle>{subtitle}</Subtitle>}
        </div>
        {action}
      </Header>
      {children}
    </Container>
  );
};

export default ReportSectionCard;

const Container = styled.div`
  background-color: #ffffff;
  padding: 32px;
  border-radius: 24px;
  margin-top: 20px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
`;

const Title = styled.h3`
  font-size: 18px;
  font-weight: 600;
  line-height: 24px;
  margin: 0;
  color: #00225a;
`;

const Subtitle = styled.p`
  font-size: 13px;
  font-weight: 400;
  margin: 4px 0 0 0;
  color: #828282;
`;
