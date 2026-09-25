"use client";

import Image from "next/image";
import React from "react";
import styled from "styled-components";
import copyIcon from "@/assets/images/document-copy.svg";

interface CardProps {
  title: string;
  content: string;
}

const WalletsComponent: React.FC<CardProps> = ({ title, content }) => {
  return (
    <CardContainer>
      <CardTitle>{title}</CardTitle>
      <Content>
        <CardContent>{content}</CardContent>
        <CopyButton>
          <Image src={copyIcon} alt="Copy" width={20} height={20} />
        </CopyButton>
      </Content>
    </CardContainer>
  );
};

export default WalletsComponent;

const CardContainer = styled.div`
  background-color: #f2f6f9;
  border-radius: 8px;
  padding: 16px;
`;

const CardTitle = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 25.27px;
  color: #828282;
`;
const CardContent = styled.h3`
  font-size: 14px;
  font-weight: 600;
  line-height: 25.27px;
  color: #00225a;
  word-break: break-word;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
  white-space: nowrap;
`;
const Content = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
`;
const CopyButton = styled.button`
  background-color: transparent;
  border: none;
  cursor: pointer;
`;
