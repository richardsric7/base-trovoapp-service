"use client";
import styled from "styled-components";
import logo from "@/assets/images/g647.svg";
import Image from "next/image";
import { FaAngleLeft, FaAngleRight, FaArrowLeft } from "react-icons/fa6";
import { useRouter } from "next/navigation";
const WalletsPage = () => {
  const router = useRouter();
  const walletData = [
    {
      id: 1,
      walletType: "Main",
      balanceLabel: "Total Balance",
      balance: "2,082,898 NGN",
      usd: "4,014 USD",
      isPrimary: true,
    },
    {
      id: 2,
      walletType: "ATL",
      balanceLabel: "Total Balance",
      balance: "2,082,898 ATL",
      usd: "4,014 USD",
    },
    {
      id: 3,
      walletType: "ATL",
      balanceLabel: "Total Balance",
      balance: "2,082,898 ATL",
      usd: "4,014 USD",
    },
    {
      id: 4,
      walletType: "ATL",
      balanceLabel: "Total Balance",
      balance: "2,082,898 ATL",
      usd: "4,014 USD",
    },
    {
      id: 5,
      walletType: "ATL",
      balanceLabel: "Total Balance",
      balance: "2,082,898 ATL",
      usd: "4,014 USD",
    },
  ];
  return (
    <>
      <FaArrowLeft
        color="#00225A"
        onClick={() => router.back()}
        style={{ cursor: "pointer", marginBottom: "20px" }}
      />
      <Container>
        <Heading>Wallets</Heading>
        <WalletsContainer>
          {walletData.map((wallet) => (
            <WalletCard $primary={wallet.isPrimary}>
              <WalletInfo>
                <HeaderRow>
                  <WalletType $primary={wallet.isPrimary}>
                    {wallet.walletType}
                  </WalletType>

                  <FaAngleRight
                    color={wallet.isPrimary ? "#fff" : "#00225A"}
                    size={14}
                  />
                </HeaderRow>

                <BalanceLabel $primary={wallet.isPrimary}>
                  {wallet.balanceLabel}
                </BalanceLabel>

                <Balance $primary={wallet.isPrimary}>{wallet.balance}</Balance>

                <BalanceUsd $primary={wallet.isPrimary}>
                  {wallet.usd}
                </BalanceUsd>
              </WalletInfo>

              <Watermark>
                <Image src={logo} alt="logo" />
              </Watermark>
            </WalletCard>
          ))}
        </WalletsContainer>
      </Container>
    </>
  );
};

export default WalletsPage;

const Container = styled.div`
  background: #fff;
  border-radius: 16px;
  padding: 20px;
`;
const Heading = styled.h1`
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  letter-spacing: 0.1px;
  text-align: left;
  color: #00225a;
  padding-bottom: 16px;
`;
const WalletsContainer = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
`;

const WalletInfo = styled.div`
  display: flex;
  flex-direction: column;
  min-width: 0;
  flex: 1;
`;
const WalletCard = styled.div<{ $primary?: boolean }>`
  position: relative;
  padding: 20px;
  border-radius: 16px;
  background: ${(p) => (p.$primary ? "#1677c8" : "#eef2f6")};
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  min-height: 120px;
  overflow: hidden;
  transition: 0.2s ease;
  cursor: pointer;
  &:hover {
    transform: translateY(-2px);
  }
`;

const HeaderRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const Watermark = styled.div`
  position: absolute;
  right: 10px;
  bottom: 10px;
  opacity: 0.1;

  img {
    width: 40px;
    height: 40px;
  }
`;

const WalletType = styled.h2<{ $primary?: boolean }>`
  font-size: 14px;
  font-weight: 500;
  color: ${(p) => (p.$primary ? "#fff" : "#00225A")};
  margin: 0;
`;

const BalanceLabel = styled.p<{ $primary?: boolean }>`
  font-size: 12px;
  color: ${(p) => (p.$primary ? "#e6f0ff" : "#7b8a9a")};
  margin-top: 10px;
`;

const Balance = styled.p<{ $primary?: boolean }>`
  font-size: 18px;
  font-weight: 600;
  color: ${(p) => (p.$primary ? "#fff" : "#00225A")};
  margin: 4px 0;
`;

const BalanceUsd = styled.p<{ $primary?: boolean }>`
  font-size: 12px;
  color: ${(p) => (p.$primary ? "#e6f0ff" : "#7b8a9a")};
`;
