"use client";

import Image from "next/image";
import { useParams } from "next/navigation";
import React, { useState } from "react";
import styled from "styled-components";
import editIcon from "@/assets/images/fluent_edit-16-regular.svg";
import deletIcon from "@/assets/images/deletIcon.svg";
import avatar from "@/assets/images/realestate.png";
import Link from "next/link";
import { FaArrowLeft, FaCopy } from "react-icons/fa6";
import DeleteAssetModal from "../components/DeleteAssetModal";
const AssetDetailPage = () => {
  const params = useParams();
  const userId = params.userId;
  const [openModal, setOpenModal] = useState(false);
  return (
    <Container>
      <Header>
        <StyledLink href="/assetcuration">
          <FaArrowLeft />
        </StyledLink>
        <ActionsWrapper>
          <StyledLink href="/assetcuration/edittoken">
            <Image src={editIcon} alt="edit" width={16} height={16} />
          </StyledLink>
          <Image
            src={deletIcon}
            alt="delet-icon"
            width={16}
            height={16}
            onClick={() => setOpenModal(!openModal)}
          />
        </ActionsWrapper>
      </Header>
      <AssetsContent>
        <Image src={avatar} alt="edit" width={40} height={40} />
        <AssetName>AFT (Animal Farm Token)</AssetName>
        <AssetText>www.animalfarm.com</AssetText>
        <AssetDescription>
          AFT - Animal Farm tokens are fractional tokens that represent part
          ownership (via investment) of our Agricultural project at Farmers
          Guild.
        </AssetDescription>
        <div>
          <AssetTitle>Category</AssetTitle>
          <AssetText>Agriculture</AssetText>
        </div>

        <div>
          <AssetTitle>Price</AssetTitle>
          <AssetText>1 ANIMAL FARM TOKEN = 0.00150 USD</AssetText>
        </div>
        <div>
          <AssetTitle>Issuer Public Key</AssetTitle>
          <AssetText>
            GAXMB...A2CFJ{" "}
            <StyledSpan>
              {" "}
              <FaCopy />
            </StyledSpan>
          </AssetText>
        </div>

        <div>
          <AssetTitle>Contact Email</AssetTitle>
          <AssetText>animalfarm@gmail.com</AssetText>
        </div>
      </AssetsContent>
      {openModal && (
        <DeleteAssetModal openModal={openModal} setOpenModal={setOpenModal} />
      )}
    </Container>
  );
};

export default AssetDetailPage;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;
  border-radius: 24px;
  gap: 20px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 20px 0;
`;

const StyledLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
`;

const ActionsWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

const AssetsContent = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  padding-bottom: 30px;
`;

const AssetName = styled.h1`
  font-weight: 600;
  font-size: 24px;
  color: #00225a;
`;
const AssetText = styled.p`
  color: #00225a;
  font-weight: 400;
  cursor: pointer;
  line-height: 20px;
  letter-spacing: -0.5%;
  text-align: center;
  font-size: 14px;
`;
const AssetTitle = styled.h2`
  font-weight: 600;
  font-size: 16px;
  color: #00225a;
  text-align: center;
  line-height: 20px;
  letter-spacing: -0.5%;
`;

const AssetDescription = styled.p`
  color: #00225a;
  font-weight: 400;
  width: 70%;
  text-align: center;
  font-size: 14px;
  line-height: 25px;
  letter-spacing: -0.5%;
  text-align: center;
`;

const StyledSpan = styled.span``;
