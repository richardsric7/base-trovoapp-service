"use client";
import React, { useEffect } from "react";
import styled from "styled-components";
import logo from "@/assets/images/g647.svg";
import Image from "next/image";
import { TokenizationRecord } from "@/redux/api/assettokenization";
import { useGetWalletBalancesQuery } from "@/redux/api/users";

interface WalletProps {
  asset?: TokenizationRecord;
}

const Wallets: React.FC<WalletProps> = ({ asset }) => {
  // console.log("Asset prop:", asset);

  const issuingWalletAddress = asset?.issuingWalletAddress;
  const marketMakingWallet = asset?.marketMakingWallet;
  const holdingWallet = asset?.walletToHoldAssetsNotForSale;
  // console.log("Using wallet key:", issuingWalletAddress);

  // Fetch wallet data
  const { data: issuingWalletData } = useGetWalletBalancesQuery(
    issuingWalletAddress!,
    {
      skip: !issuingWalletAddress,
    }
  );
  const { data: marketMakingWalletData } = useGetWalletBalancesQuery(
    marketMakingWallet!,
    {
      skip: !marketMakingWallet,
    }
  );
  const { data: holdingWalletData } = useGetWalletBalancesQuery(
    holdingWallet!,
    {
      skip: !holdingWallet,
    }
  );

  useEffect(() => {
    console.log("Issuing wallet data:", issuingWalletData?.claimed);
    console.log("Market making wallet data:", marketMakingWalletData);
    console.log("Holding wallet data:", holdingWalletData);
  }, [issuingWalletData, marketMakingWalletData, holdingWalletData]);

  // Helper function to extract the amount for a given asset code.
  // For NGN, assetCode is "" and for USD we look for "USDT"
  const getBalance = (walletData: any, assetCode: string): string =>
    walletData?.data?.claimed?.find((item: any) => item.assetCode === assetCode)
      ?.amount || "0";

  // Computing balances for each wallet
  const issuingNGN = getBalance(issuingWalletData, "");
  const issuingUSD = getBalance(issuingWalletData, "USDT");

  const marketMakingNGN = getBalance(marketMakingWalletData, "");
  const marketMakingUSD = getBalance(marketMakingWalletData, "USDT");

  const holdingNGN = getBalance(holdingWalletData, "");
  const holdingUSD = getBalance(holdingWalletData, "USDT");

  const walletData = [
    {
      id: 1,
      // walletType: "Minting Wallet",
      walletType: "Issuing Wallet",
      walletOwner: asset?.issuingWalletAlias,
      balanceLabel: "Total Balance",
      balanceNGN: issuingNGN + "NGN",
      balanceUSD: issuingUSD + " USD",
    },
    {
      id: 2,
      // walletType: "Distribution Wallet",
      walletType: "Distribution Wallet",
      walletOwner: asset?.issuingWalletAlias
        ? `${asset.issuingWalletAlias}-distribution`
        : "",

      balanceLabel: "Total Balance",
      balanceNGN: marketMakingNGN + "NGN",
      balanceUSD: marketMakingUSD + " USD",
    },
    {
      id: 3,
      walletType: "Holding Wallet",
      walletOwner: asset?.issuingWalletAlias
        ? `${asset.issuingWalletAlias}-holding`
        : "",

      balanceLabel: "Total Balance",
      balanceNGN: holdingNGN + "NGN",
      balanceUSD: holdingUSD + " USD",
    },
  ];

  return (
    <Container>
      <Heading>Wallets</Heading>
      <WalletsContainer>
        {walletData.map((wallet) => (
          <WalletCard key={wallet.id}>
            <WalletInfo>
              <WalletType>{wallet.walletType}</WalletType>
              <WalletOwner>{wallet.walletOwner}</WalletOwner>
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
  text-align: left;
  color: #00225a;
  padding-bottom: 16px;
`;
const WalletsContainer = styled.div`
  display: flex;
  gap: 20px;
  width: 100%;
  justify-content: space-between;
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
    opacity: 1px;
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
  max-width: 100%;
  padding: 16px 10px;
  gap: 16px;
  border-radius: 20px;
  opacity: 0px;
  background-color: #f2f6f9;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border: 1px solid #007cdf;
  overflow: hidden;
`;
const WalletType = styled.h2`
  font-size: 14px;
  font-weight: 400;
  line-height: 25.27px;
  letter-spacing: 0.019em;
  margin-bottom: 10px;
  color: #00225ab2;
  margin: 0px;
`;
const WalletOwner = styled.p`
  font-size: 16px;
  font-weight: 600;
  line-height: 25.27px;
  letter-spacing: 0.019em;
  margin-bottom: 4px;
  color: #00225ab2;
`;
const BalanceLabel = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 24px;
  letter-spacing: 0.1px;
  margin-bottom: 1px;
`;
const Balance = styled.p`
  font-size: 20px;
  font-weight: bold;
  line-height: 24px;
  letter-spacing: 0.1px;
  text-align: left;
  color: #00225ab2;
  margin: 0px;
`;
const BalanceUsd = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  letter-spacing: 0.1px;
  text-align: left;
  color: #00225ab2;
  margin: 0px;
`;
