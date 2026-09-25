import React from "react";
import { Modal } from "@/components";
import { ArrowDown } from "iconsax-react";
import styled from "styled-components";

interface TokenHoldersModalProps {
  isOpen: boolean;
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

interface TokenTransaction {
  description: string;
  date: string;
  amount: string;
}

const transactions: TokenTransaction[] = [
  {
    description: "Primary sale purchase",
    date: "9:30 AM, 12 Oct, 2023",
    amount: "+15 AEI",
  },
  {
    description: "Sold to chidinma_investments on TrovoP2P",
    date: "9:30 AM, 12 Oct, 2023",
    amount: "-15 AEI",
  },
  {
    description: "Received from Tom_secondary on Trovo App",
    date: "9:30 AM, 12 Oct, 2023",
    amount: "+15 AEI",
  },
  {
    description: "Sent to chidinma_investments on Trovo App",
    date: "9:30 AM, 12 Oct, 2023",
    amount: "-15 AEI",
  },
  {
    description: "Primary sale purchase",
    date: "9:30 AM, 12 Oct, 2023",
    amount: "+15 AEI",
  },
];

const TokenHoldersModal: React.FC<TokenHoldersModalProps> = ({
  isOpen,
  setIsOpen,
}) => {
  return (
    <Modal
      title="Token Activity"
      isOpen={isOpen}
      onClose={() => setIsOpen(false)}
      closeIconPosition="right"
    >
      <ModalBody>
        <Header>
          <UserProfile>
            <UserAvatar />
            <UserDetails>
              <UserId>Folrence</UserId>
              <UserName>Florence Zach</UserName>
            </UserDetails>
          </UserProfile>
          <UserBalance>
            <BalanceLabel>Total</BalanceLabel>
            <BalanceValue>100 AEI</BalanceValue>
          </UserBalance>
        </Header>
        <Divider />

        <TransactionList>
          {transactions.map((txn, idx) => {
            const isNegative = txn.amount.startsWith("-");
            return (
              <TransactionRow key={idx}>
                <TransactionInfo>
                  <ArrowCircle $negative={isNegative}>
                    <ArrowDown color={isNegative ? "#BE3800" : "#00A859"} />
                  </ArrowCircle>
                  <div>
                    <TransactionDescription>
                      {txn.description}
                    </TransactionDescription>
                    <TransactionDate>{txn.date}</TransactionDate>
                  </div>
                </TransactionInfo>
                <TransactionAmount $negative={isNegative}>
                  {txn.amount}
                </TransactionAmount>
              </TransactionRow>
            );
          })}
        </TransactionList>
      </ModalBody>
    </Modal>
  );
};

export default TokenHoldersModal;

const ModalBody = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const Divider = styled.div`
  height: 1px;
  width: 100%;
  background-color: #e5e5ef;
  margin: 10px 0;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const UserProfile = styled.div`
  display: flex;
  align-items: center;
  gap: 0.5rem;
`;

const UserAvatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const UserDetails = styled.div`
  display: flex;
  flex-direction: column;
`;

const UserId = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  margin: 0;
`;

const UserName = styled.p`
  font-size: 10px;
  font-weight: 400;
  letter-spacing: 0.1px;
  color: #00225a;
  margin: 0;
`;

const UserBalance = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-end;
`;

const BalanceLabel = styled.span`
  font-weight: 400;
  font-size: 11px;
  color: #4f4f4f;
`;

const BalanceValue = styled.span`
  font-weight: 600;
  font-size: 14px;
  color: #00225a;
`;

const TransactionList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 14px;
`;

const TransactionRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const TransactionInfo = styled.div`
  display: flex;
  align-items: center;
  gap: 0.5rem;
`;

const ArrowCircle = styled.div<{ $negative: boolean }>`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: ${({ $negative }) => ($negative ? "#BE38001A" : "#00a8591a")};
  display: flex;
  align-items: center;
  justify-content: center;
`;

const TransactionDescription = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  margin: 0;
`;

const TransactionDate = styled.p`
  font-size: 10px;
  font-weight: 400;
  letter-spacing: 0.1px;
  color: #4f4f4f;
  margin: 0;
`;

const TransactionAmount = styled.p<{ $negative: boolean }>`
  font-weight: 600;
  font-size: 14px;
  color: ${({ $negative }) => ($negative ? "#BE3800" : "#00A859")};
  margin: 0;
`;
