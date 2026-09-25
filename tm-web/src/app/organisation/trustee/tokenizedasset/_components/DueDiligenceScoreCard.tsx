"use client";
import Image from "next/image";
import styled from "styled-components";
import noteIcon from "@/assets/images/note.svg";

interface Props {
  documentCount: number;
}

const DueDiligenceScore = ({ documentCount }: Props) => {
  return (
    <ScoreCard>
      <Image src={noteIcon} alt="note-icon" />
      <div>
        <ScoreLabel>Documents for Review</ScoreLabel>
        <Score>{documentCount}</Score>
      </div>
    </ScoreCard>
  );
};

export default DueDiligenceScore;

const ScoreCard = styled.div`
  background: #f2f6f9;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
`;

const ScoreLabel = styled.div`
  font-size: 12px;
  color: #7c8aa0;
  margin-bottom: 6px;
`;

const Score = styled.p`
  font-size: 24px;
  font-weight: 700;
  color: #191919;
`;
