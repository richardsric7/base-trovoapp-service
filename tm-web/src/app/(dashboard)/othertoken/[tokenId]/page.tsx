"use client";

import Image from "next/image";
import { useParams } from "next/navigation";
import React, { useState } from "react";
import styled from "styled-components";
import editIcon from "@/assets/images/fluent_edit-16-regular.svg";
import deleteIcon from "@/assets/images/deletIcon.svg";
import trov from "@/assets/images/TROVTokenicon.svg";
import Link from "next/link";
import { FaArrowLeft, FaCopy } from "react-icons/fa6";
import DeleteTokenModal from "../../assetcuration/components/DeleteAssetModal";

const TokenDetailsPage = () => {
  const params = useParams();
  const tokenId = params.tokenId;
  const [isModalOpen, setIsModalOpen] = useState(false);

  return (
    <TokenContainer>
      <TokenHeader>
        <BackLink href="/othertoken">
          <FaArrowLeft />
        </BackLink>
        <TokenActions>
          <EditLink href="/othertoken/edittoken">
            <Image src={editIcon} alt="Edit Token" width={16} height={16} />
          </EditLink>
          <DeleteIcon
            src={deleteIcon}
            alt="Delete Token"
            width={16}
            height={16}
            onClick={() => setIsModalOpen(!isModalOpen)}
          />
        </TokenActions>
      </TokenHeader>
      <TokenContent>
        <Image src={trov} alt="TROV Token" width={40} height={40} />
        <TokenName>TROV Token</TokenName>
        <TokenWebsite>www.trovotech.io</TokenWebsite>
        <TokenDescription>
          TROV token (TROV) is the utility token that powers the Trovotech
          ecosystem. TROV token is used to access discounts, voting rights,
          airdrops, NFTs, and other community incentives.
        </TokenDescription>

        <TokenInfo>
          <TokenLabel>Price</TokenLabel>
          <TokenValue>1 TROV = 0.00300 USD</TokenValue>
        </TokenInfo>

        <TokenInfo>
          <TokenLabel>Issuer Public Key</TokenLabel>
          <TokenValue>
            GAXMB...A2CFJ{" "}
            <CopyIcon>
              <FaCopy />
            </CopyIcon>
          </TokenValue>
        </TokenInfo>

        <TokenInfo>
          <TokenLabel>Contact Email</TokenLabel>
          <TokenValue>admin@trovotech.io</TokenValue>
        </TokenInfo>
      </TokenContent>

      {isModalOpen && (
        <DeleteTokenModal
          openModal={isModalOpen}
          setOpenModal={setIsModalOpen}
        />
      )}
    </TokenContainer>
  );
};

export default TokenDetailsPage;


const TokenContainer = styled.section`
  background: #ffffff;
  padding: 32px;
  border-radius: 24px;
  gap: 20px;
`;

const TokenHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 20px 0;
`;

const BackLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
`;

const TokenActions = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

const EditLink = styled(Link)`
  display: flex;
  align-items: center;
`;

const DeleteIcon = styled(Image)`
  cursor: pointer;
`;

const TokenContent = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  padding-bottom: 30px;
`;

const TokenName = styled.h1`
  font-weight: 600;
  font-size: 24px;
  color: #00225a;
`;

const TokenWebsite = styled.p`
  color: #00225a;
  font-weight: 400;
  cursor: pointer;
  line-height: 20px;
  letter-spacing: -0.5%;
  text-align: center;
  font-size: 14px;
`;

const TokenLabel = styled.h2`
  font-weight: 600;
  font-size: 16px;
  color: #00225a;
  text-align: center;
  line-height: 20px;
  letter-spacing: -0.5%;
`;

const TokenValue = styled.p`
  color: #00225a;
  font-weight: 400;
  font-size: 14px;
  text-align: center;
  line-height: 20px;
  letter-spacing: -0.5%;
`;

const TokenDescription = styled.p`
  color: #00225a;
  font-weight: 400;
  width: 70%;
  text-align: center;
  font-size: 14px;
  line-height: 25px;
  letter-spacing: -0.5%;
`;

const TokenInfo = styled.div`
  text-align: center;
  margin-bottom: 8px;
`;

const CopyIcon = styled.span`
  cursor: pointer;
  margin-left: 5px;
`;
