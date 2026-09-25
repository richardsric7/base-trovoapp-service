"use client";

import { Modal } from "antd";
import styled from "styled-components";
import Image, { StaticImageData } from "next/image";
import { IoMdCopy } from "react-icons/io";
import PrimaryButton from "@/components/PrimaryButton";
import successIcon from "@/assets/images/success-icon.svg";

interface InfoItem {
  label: string;
  value: string;
  copyable?: boolean;
}

interface ReusableSuccessModalProps {
  open: boolean;
  onClose: () => void;
  title?: string;
  description?: string;
  buttonText?: string;
  icon?: StaticImageData | string;
  infoItems?: InfoItem[];
  width?: number;
}

export default function ReusableSuccessModal({
  open,
  onClose,
  title = "Successful!",
  description = "Operation completed successfully.",
  buttonText = "Done",
  icon = successIcon,
  infoItems = [],
  width = 420,
}: ReusableSuccessModalProps) {
  const handleCopy = async (value: string) => {
    try {
      await navigator.clipboard.writeText(value);
    } catch (error) {
      console.error("Failed to copy text:", error);
    }
  };

  return (
    <Modal
      open={open}
      onCancel={onClose}
      footer={null}
      centered
      width={width}
      closeIcon={<span style={{ fontSize: "16px" }}>×</span>}
    >
      <ModalContent>
        <Title>{title}</Title>

        <Image src={icon} alt="modal-icon" width={80} height={80} />

        <Description>{description}</Description>

        {infoItems.length > 0 && (
          <InfoBox>
            {infoItems.map((item, index) => (
              <InfoItemWrapper key={`${item.label}-${index}`}>
                <Label>{item.label}</Label>

                {item.copyable ? (
                  <CopyRow>
                    <Value title={item.value}>{item.value}</Value>
                    <CopyIcon onClick={() => handleCopy(item.value)}>
                      <IoMdCopy color="#00225A" />
                    </CopyIcon>
                  </CopyRow>
                ) : (
                  <Value title={item.value}>{item.value}</Value>
                )}
              </InfoItemWrapper>
            ))}
          </InfoBox>
        )}

        <PrimaryButton onClick={onClose}>{buttonText}</PrimaryButton>
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

const InfoItemWrapper = styled.div`
  &:not(:last-child) {
    margin-bottom: 16px;
  }
`;

const Label = styled.p`
  font-size: 12px;
  color: #00225a;
  margin-bottom: 4px;
`;

const Value = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  display: flex;
  align-items: center;
  gap: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  margin: 0;
`;

const CopyRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
`;

const CopyIcon = styled.span`
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
`;
