"use client";

import styled from "styled-components";
import { NoteRemove } from "iconsax-react";

interface EmptyTableDataProps {
  title?: string;
  description?: string;
}

const EmptyTableData = ({
  title = "No data yet",
  description = "Records will appear here once available.",
}: EmptyTableDataProps) => (
  <Container role="status">
    <Title>{title}</Title>
    <Description>{description}</Description>
  </Container>
);

export default EmptyTableData;

const Container = styled.div`
  display: flex;
  min-height: 190px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  text-align: center;
`;

const IconWrapper = styled.div`
  display: grid;
  width: 44px;
  height: 44px;
  margin-bottom: 12px;
  place-items: center;
  border: 1px solid #e5e5ef;
  border-radius: 12px;
  background: #f8f9fc;
  color: #667085;
`;

const Title = styled.p`
  margin: 0;
  color: #00225a;
  font-size: 14px;
  font-weight: 600;
  line-height: 20px;
`;

const Description = styled.p`
  margin: 4px 0 0;
  color: #828282;
  font-size: 12px;
  font-weight: 400;
  line-height: 18px;
`;
