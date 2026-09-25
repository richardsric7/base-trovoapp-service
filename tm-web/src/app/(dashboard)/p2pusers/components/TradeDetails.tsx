import { Modal } from "@/components";
import React from "react";
import { MdOutlineContentCopy } from "react-icons/md";
import { toast } from "react-toastify";
import styled from "styled-components";

interface TransactionModalProps {
  openModal: boolean;
  setOpenModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const TradeDetails: React.FC<TransactionModalProps> = ({
  openModal,
  setOpenModal,
}) => {
  const transaction = {
    category: "Token Transfer",
    date: "May 16, 2025 10:00 AM",
    amount: "₦250,000",
    transactionId: "0x1234567890abcdef1234567890abcdef",
    to: "wallet_to_123456",
    toAddress: "pubkey_to_abcdef1234",
    from: "wallet_from_987654",
    fromAddress: "pubkey_from_5678abcd",
  };

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

  const getTransactionColor = () => {
    return "#00A859";
  };

  const getTransactionIcon = () => {
    return <span style={{ fontSize: 18 }}>🔁</span>;
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
            <TransactionTimestamp>
              {transaction.date}
              {/* November 12, 2023 9:14 am */}
            </TransactionTimestamp>
          </div>
        </HeaderDetails>
        <TransactionValue style={{ color: getTransactionColor() }}>
          {transaction.amount}
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

        <DetailsSection>
          <Label>Type</Label>
          <LabelText>Taker</LabelText>
        </DetailsSection>

        <DetailsSection>
          <Label>Status</Label>
          <LabelText>Completed</LabelText>
        </DetailsSection>

        <DetailsSection>
          <Label>Payment Method</Label>
          <LabelText>Bank Transfer</LabelText>
        </DetailsSection>

        <SectionDivider />
        <Section>
          <SectionTitle>User</SectionTitle>

          <DetailsSection>
            <Label>Name</Label>
            <LabelText>Florence Zach</LabelText>
          </DetailsSection>
          <DetailsSection>
            <Label>Wallet</Label>
            <CopyableTextBlue
              onClick={() =>
                handleCopy(transaction.to || transaction.toAddress)
              }
            >
              {transaction.to.substring(0, 12) ||
                transaction.toAddress.substring(0, 12)}
              ...
              <MdOutlineContentCopy size={20} />
            </CopyableTextBlue>
          </DetailsSection>
          <DetailsSection>
            <Label>Public Key</Label>
            <CopyableTextBlue
              onClick={() => handleCopy(transaction.toAddress)}
            >
              {transaction.toAddress.substring(0, 12)}...{" "}
              <MdOutlineContentCopy size={20} />
            </CopyableTextBlue>
          </DetailsSection>
        </Section>

        <SectionDivider />
        <Section>
          <SectionTitle>Sent To</SectionTitle>

          <DetailsSection>
            <Label>Peer</Label>
            <LabelText>Kennis Maduka</LabelText>
          </DetailsSection>
          <DetailsSection>
            <Label>Wallet</Label>
            <CopyableTextBlue
              onClick={() =>
                handleCopy(transaction.from || transaction.fromAddress)
              }
            >
              {transaction.from.substring(0, 12) ||
                transaction.fromAddress.substring(0, 12)}
              ...
              <MdOutlineContentCopy size={20} />
            </CopyableTextBlue>
          </DetailsSection>
          <DetailsSection>
            <Label>Public Key</Label>
            <CopyableTextBlue
              onClick={() => handleCopy(transaction.fromAddress)}
            >
              {transaction.fromAddress?.substring(0, 12)}...
              <MdOutlineContentCopy size={20} />
            </CopyableTextBlue>
          </DetailsSection>
        </Section>
      </ModalBody>
    </Modal>
  );
};

export default TradeDetails;

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

const LabelText = styled.p`
  font-size: 14px;
  font-weight: 500;

  color: #00225a;
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
