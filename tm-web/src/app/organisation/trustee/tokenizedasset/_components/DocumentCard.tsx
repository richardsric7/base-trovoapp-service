"use client";

import React from "react";
import styled from "styled-components";
import Image from "next/image";
import pdfIcon from "@/assets/images/pdfdoc.svg";

interface DocumentCardProps {
  title: string;
  size?: string;
  fileUrl?: string;
  onClick?: () => void;
  isLoading?: boolean;
}

const DocumentCard: React.FC<DocumentCardProps> = ({
  title,
  size,
  fileUrl,
  onClick,
  isLoading = false,
}) => {
  return (
    <Card
      type="button"
      onClick={onClick}
      disabled={!onClick || isLoading}
      aria-busy={isLoading}
      aria-label={`${isLoading ? "Opening" : "View"} ${title}`}
    >
      <Header>
        <HeaderTitle>{title}</HeaderTitle>
      </Header>

      <Content>
        <Image src={pdfIcon} alt="PDF" width={18} height={18} />

        <Info>
          <FileTitle>{isLoading ? "Opening document…" : title}</FileTitle>
          {size && <FileSize>{size}</FileSize>}
        </Info>
      </Content>

      <Creator>
        <CreatorLabel>
          {fileUrl ? "Uploaded by Creator" : "Uploaded by Creator or Admin"}
        </CreatorLabel>

        <CreatorAvatar>C</CreatorAvatar>
      </Creator>
    </Card>
  );
};

export default DocumentCard;

const Card = styled.button`
  font-family: inherit;
  background: #f2f6f9;
  border: 0;
  border-radius: 8px;
  padding: 12px;
  cursor: pointer;
  width: 100%;
  text-align: left;

  display: flex;
  flex-direction: column;
  gap: 10px;

  &:disabled {
    cursor: not-allowed;
  }
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const HeaderTitle = styled.p`
  font-size: 12px;
  font-weight: 600;
  color: #00225a;
  margin: 0;
`;

const Content = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const Info = styled.div`
  display: flex;
  flex-direction: column;
`;

const FileTitle = styled.p`
  font-size: 12px;
  font-weight: 500;
  color: #4f4f4f;
  margin: 0;
`;

const FileSize = styled.p`
  font-size: 11px;
  color: #9aa6b2;
  margin: 0;
`;

const Creator = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const CreatorLabel = styled.p`
  font-size: 11px;
  color: #191919;
  margin: 0;
  background: #fff;
  border-radius: 8px;
  padding: 4px 8px;
`;

const CreatorAvatar = styled.div`
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #007cdf;
  color: #fff;

  display: flex;
  align-items: center;
  justify-content: center;

  font-size: 11px;
  font-weight: 600;
`;
