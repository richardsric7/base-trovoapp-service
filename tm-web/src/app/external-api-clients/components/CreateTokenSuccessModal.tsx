"use client";

import { Modal } from "antd";
import styled from "styled-components";
import PrimaryButton from "@/components/PrimaryButton";
import Image from "next/image";
import successIcon from "@/assets/images/success-icon.svg";
import { IoMdCopy } from "react-icons/io";
interface Props {
  open: boolean;
  onClose: () => void;
  tokenCode: string;
  publicKey: string;
}

export default function SuccessModal({
  open,
  onClose,
  tokenCode,
  publicKey,
}: Props) {
  const handleCopy = () => {
    navigator.clipboard.writeText(publicKey);
  };

  return (
    <Modal
      open={open}
      onCancel={onClose}
      footer={null}
      centered
      width={420}
      closeIcon={<span style={{ fontSize: "16px" }}>×</span>}
    >
      <ModalContent>
        <Title>Successful!</Title>

        <Image src={successIcon} alt="success-icon" width={80} height={80} />

        <Description>
          You&apos;ve successfully created a token with the following details:
        </Description>

        <InfoBox>
          <Label>Token Code</Label>
          <Value>ATL</Value>

          <Label>Public Key</Label>
          <CopyRow>
            <Value style={{ marginBottom: 0 }}>
              AS6HFTRHVU9876YNHGDLJFUEVSMKEF....EUIE75SDCVBNMU70453EHNKKKNLOJHDSSO98
            </Value>
            <CopyIcon onClick={handleCopy}>
              <IoMdCopy color="#00225A" />
            </CopyIcon>
          </CopyRow>
        </InfoBox>

        <PrimaryButton onClick={onClose}>Done</PrimaryButton>
      </ModalContent>
    </Modal>
  );
}

const ModalContent = styled.div`
  text-align: center;
  padding: 10px 10px 0;
`;

const Title = styled.h2`
  font-size: 22px;
  font-weight: 600;
  color: #1f3b64;
  margin-bottom: 20px;
`;

const Description = styled.p`
  font-size: 14px;
  color: #00225a;
  margin-bottom: 20px;
`;

const InfoBox = styled.div`
  background: #f2f6f9;
  border-radius: 10px;
  padding: 16px;
  text-align: left;
  margin-bottom: 24px;
`;

const Label = styled.p`
  font-size: 12px;
  color: #00225a;
`;

const Value = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
`;

const CopyRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 4px;
`;

const CopyIcon = styled.span`
  cursor: pointer;
  font-size: 16px;
`;
