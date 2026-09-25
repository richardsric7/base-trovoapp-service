"use client";
import { Modal } from "@/components";
import React, { useState } from "react";
import PrimaryButton from "@/components/PrimaryButton";
import Image from "next/image";
import asset from "@/assets/images/realestate.png";

import styled from "styled-components";

interface AssetCurationProps {
  openModal: boolean;
  setOpenModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const CompleteCurationModal: React.FC<AssetCurationProps> = ({
  openModal,
  setOpenModal,
}) => {
  const [showAllAssets, setShowAllAssets] = useState(false);

  const handleSubmit = () => {
    setOpenModal(false);
  };

  const assets = [asset, asset, asset, asset, asset, asset, asset];

  const displayedAssets = showAllAssets ? assets : assets.slice(0, 5);
  const remainingAssets = assets.length - 5;

  return (
    <Modal
      title="Asset Curation"
      isOpen={openModal}
      onClose={() => setOpenModal(false)}
    >
      <Text>
        You are about to complete your curation of the following assets
      </Text>
      <ImageWrapper>
        {displayedAssets.map((_, index) => (
          <ImageContainer key={index}>
            <Image src={asset} alt={`asset-${index}`} width={40} height={40} />
          </ImageContainer>
        ))}
        {!showAllAssets && remainingAssets > 0 && (
          <AssetCircle onClick={() => setShowAllAssets(true)}>
            +{remainingAssets}
          </AssetCircle>
        )}
      </ImageWrapper>
      {showAllAssets && (
        <ViewLessButton onClick={() => setShowAllAssets(false)}>
          View Less
        </ViewLessButton>
      )}
      <div>
        <Label>Curation Title</Label>
        <StyledInput />
      </div>
      <PrimaryButton onClick={handleSubmit}>Complete Curation</PrimaryButton>
    </Modal>
  );
};

export default CompleteCurationModal;

const Label = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  letter-spacing: 0.10000000149011612px;
  color: #00225a;
`;

const Text = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  margin: 0;
  color: #00225a;
`;

const ImageWrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  margin-top: 10px;
`;

const StyledInput = styled.input`
  width: 95%;
  padding: 12px;
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 400;
  color: #000;
  background-color: transparent;
  display: block;
  outline: none;

  &:focus {
    outline: none;
    border-color: #007cdf;
    background-color: #ffffff;
  }
`;
const ImageContainer = styled.div`
  margin-left: -8px;
  &:first-child {
    margin-left: 0;
  }
`;

const AssetCircle = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
  color: #00225a;
  margin-left: -8px;
  cursor: pointer;
`;
const ViewLessButton = styled.button`
  background: none;
  border: none;
  color: #007cdf;
  font-size: 14px;
  cursor: pointer;
  margin-top: 10px;
  padding: 5px 10px;
  border-radius: 5px;
  transition: background-color 0.3s;

  &:hover {
    background-color: #f0f0f0;
  }
`;
