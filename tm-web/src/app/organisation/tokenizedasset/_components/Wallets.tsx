"use client";

import React from "react";
import styled from "styled-components";
import logo from "@/assets/images/g647.svg";
import Image from "next/image";

const dummyWallets = [
  {
    id: 1,
    walletType: "Issuing Wallet",
    walletOwner: "atprofile_iky001issuer",
    balanceLabel: "Total Balance",
    balanceNGN: "0 NGN",
    balanceUSD: "0 USD",
  },
  {
    id: 2,
    walletType: "Distribution Wallet",
    walletOwner: "atprofile_iky001issuer-distribution",
    balanceLabel: "Total Balance",
    balanceNGN: "0 NGN",
    balanceUSD: "0 USD",
  },
  {
    id: 3,
    walletType: "Holding Wallet",
    walletOwner: "atprofile_iky001issuer-holding",
    balanceLabel: "Total Balance",
    balanceNGN: "0 NGN",
    balanceUSD: "0 USD",
  },
];

const Wallets: React.FC = () => {
  return (
    <Container>
      <Heading>Wallets</Heading>

      <WalletsContainer>
        {dummyWallets.map((wallet) => (
          <WalletCard key={wallet.id}>
            <WalletInfo>
              <WalletType>{wallet.walletType}</WalletType>

              <WalletOwner title={wallet.walletOwner}>
                {wallet.walletOwner}
              </WalletOwner>

              <BalanceLabel>{wallet.balanceLabel}</BalanceLabel>

              <Balance>{wallet.balanceNGN}</Balance>

              <BalanceUsd>{wallet.balanceUSD}</BalanceUsd>
            </WalletInfo>

            <ImageWrapper>
              <Image src={logo} alt="trovo-logo" />
            </ImageWrapper>
          </WalletCard>
        ))}
      </WalletsContainer>
    </Container>
  );
};

export default Wallets;

const Container = styled.div`
  margin: 50px 0;
  width: 100%;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #00225a;
  padding-bottom: 16px;
`;

const WalletsContainer = styled.div`
  display: flex;
  gap: 20px;
  width: 100%;
`;

const ImageWrapper = styled.div`
  width: 48px;
  height: 42px;
  margin-bottom: 10px;

  img {
    width: 100%;
    height: 100%;
    object-fit: contain;
    mix-blend-mode: overlay;
  }
`;

const WalletInfo = styled.div`
  display: flex;
  flex-direction: column;
  min-width: 0;
  flex: 1;
`;

const WalletCard = styled.div`
  flex: 1;
  min-width: 0;
  padding: 16px 10px;
  border-radius: 20px;
  background-color: #f2f6f9;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border: 1px solid #007cdf;
`;

const WalletType = styled.h2`
  font-size: 14px;
  font-weight: 400;
  color: #00225ab2;
  margin: 0;
  padding-bottom: 4px;
`;

const WalletOwner = styled.p`
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 4px;
  color: #00225ab2;

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const BalanceLabel = styled.p`
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 1px;
`;

const Balance = styled.p`
  font-size: 20px;
  font-weight: bold;
  color: #00225ab2;
  margin: 0;
`;

const BalanceUsd = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225ab2;
  margin: 0;
`;
