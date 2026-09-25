import { Modal } from "@/components";
import React from "react";
import { FaExchangeAlt } from "react-icons/fa";
import { FaArrowDown, FaArrowUp } from "react-icons/fa6";
import { MdOutlineContentCopy } from "react-icons/md";
import { toast } from "react-toastify";
import styled from "styled-components";

interface TransactionModalProps {
  openModal: boolean;
  setOpenModal: React.Dispatch<React.SetStateAction<boolean>>;
  transaction: any;
}

const TransactionModal: React.FC<TransactionModalProps> = ({
  openModal,
  setOpenModal,
  transaction,
}) => {
  const handleCopy = (text: string) => {
    if (!text) {
      toast.error("Nothing to copy.");
      return;
    }
    navigator.clipboard
      .writeText(text)
      .then(() => toast.success("copied to clipboard!"))
      .catch(() => toast.error("Failed to copy."));
  };

  const getTransactionIcon = () => {
    const category = transaction?.category?.toUpperCase() || "";
    if (category.startsWith("SWAP")) {
      return <FaExchangeAlt color="#007CDF" />;
    }

    switch (category) {
      case "SENT":
        return <FaArrowUp color="#BE3800" />;
      case "PAYMENT":
        return <FaArrowDown color="#00A859" />;
      default:
        return <FaArrowDown color="gray" />;
    }
  };

  const getTransactionColor = () => {
    const category = transaction?.category?.toUpperCase() || "";
    if (category.startsWith("SWAP")) return "#007CDF";

    switch (category) {
      case "SENT":
        return "#BE3800";
      case "PAYMENT":
        return "#00A859";
      default:
        return "gray";
    }
  };

  const extractNameFromBrackets = (text?: string) => {
    const match = text?.match(/\[(.*?)\]/);
    return match ? match[1] : text;
  };

  const formatAmount = (amount: number | string) => {
    const num = typeof amount === "string" ? parseFloat(amount) : amount;
    if (isNaN(num)) return "--";
    return new Intl.NumberFormat("en-US", {
      minimumFractionDigits: 2, // Always show at least 2
      maximumFractionDigits: 2,
    }).format(num);
  };

  return (
    <Modal title="" isOpen={openModal} onClose={() => setOpenModal(false)}>
      <ModalTitle>Transaction Details</ModalTitle>
      <ModalHeader>
        <HeaderDetails>
          <IconWrapper
            style={{ backgroundColor: `${getTransactionColor()}33` }}
          >
            {getTransactionIcon()}
          </IconWrapper>
          <div>
            <TransactionStatus style={{ color: getTransactionColor() }}>
              {transaction.category}
            </TransactionStatus>
            <TransactionTimestamp>{transaction.date}</TransactionTimestamp>
          </div>
        </HeaderDetails>
        <TransactionValue style={{ color: getTransactionColor() }}>
          {formatAmount(transaction.amount)}{" "}
          <span>{transaction.assetCode}</span>
        </TransactionValue>
      </ModalHeader>

      <ModalBody>
        <DetailsSection>
          <Label>Blockchain ID</Label>
          <CopyableText onClick={() => handleCopy(transaction.transactionId)}>
            {transaction.transactionId.substring(0, 15)}...
            <MdOutlineContentCopy color="#00225A" size={20} />
          </CopyableText>
        </DetailsSection>

        <SectionDivider />
        <Section>
          <SectionTitle>Received On</SectionTitle>
          <DetailsSection>
            <Label>Wallet</Label>
            <CopyableTextBlue
              onClick={() => handleCopy(transaction.to || transaction.to)}
            >
              {extractNameFromBrackets(transaction.to)?.substring(0, 12) ||
                "--"}

              <MdOutlineContentCopy size={20} />
            </CopyableTextBlue>
          </DetailsSection>
          <DetailsSection>
            <Label>Public Key</Label>
            <CopyableTextBlue
              onClick={() => handleCopy(transaction.toPublicKey)}
            >
              {transaction.toPublicKey.substring(0, 12)}...{" "}
              <MdOutlineContentCopy size={20} />
            </CopyableTextBlue>
          </DetailsSection>
        </Section>

        <SectionDivider />
        <Section>
          <SectionTitle>From</SectionTitle>
          <DetailsSection>
            <Label>Wallet</Label>
            <CopyableTextBlue
              onClick={() => handleCopy(transaction.from || transaction.from)}
            >
              {extractNameFromBrackets(transaction.from) || " --"}

              <MdOutlineContentCopy size={20} />
            </CopyableTextBlue>
          </DetailsSection>
          <DetailsSection>
            <Label>Public Key</Label>
            <CopyableTextBlue
              onClick={() => handleCopy(transaction.fromPublicKey)}
            >
              {transaction.fromPublicKey?.substring(0, 12)}...
              <MdOutlineContentCopy size={20} />
            </CopyableTextBlue>
          </DetailsSection>
        </Section>
      </ModalBody>
    </Modal>
  );
};

export default TransactionModal;

const ModalTitle = styled.h3`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
  margin-bottom: 10px;
`;

const ModalHeader = styled.div`
  background-color: #f2f6f9;
  padding: 16px;
  gap: 10px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 30px;
`;

const IconWrapper = styled.div`
  background-color: #00a85933;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
`;

const TransactionStatus = styled.h4`
  font-size: 14px;
  font-weight: 500;
  line-height: 24px;
  color: #00a859;
`;

const TransactionTimestamp = styled.p`
  font-size: 12px;
  font-weight: 400;
  color: #828282;
`;

const TransactionValue = styled(TransactionStatus)``;

const HeaderDetails = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

const ModalBody = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

const DetailsSection = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  column-gap: 50px;
`;

const Label = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #828282;
`;

const CopyableText = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #007cdf;
  display: flex;
  aligin-items: center;
  gap: 4px;
  cursor: pointer;
`;

const CopyableTextBlue = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  display: flex;
  aligin-items: center;
  gap: 4px;
  cursor: pointer;
`;

const SectionDivider = styled.div`
  width: 100%;
  height: 1px;
  background-color: #e5e5ef;
  margin: 12px 0;
`;

const SectionTitle = styled(ModalTitle)`
  margin-bottom: 0;
`;

const Section = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
`;
