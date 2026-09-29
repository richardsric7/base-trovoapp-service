"use client";
import React, { useState } from "react";
import styled from "styled-components";
import logo from "@/assets/images/g647.svg";
import Image from "next/image";
import {
  TokenizationRecord,
  useRegisterTokenContractMutation,
} from "@/redux/api/assettokenization";
import { useGetWalletBalancesQuery } from "@/redux/api/users";
import { showErrorToast, showSuccessToast } from "@/components";

interface WalletProps {
  asset?: TokenizationRecord;
}

const isAddress = (value: string) => /^0x[0-9a-fA-F]{40}$/.test(value.trim());

// A tokenized asset on Base has two distinct addresses:
//  - the issuing Safe (issuingWalletAddress): a multisig that owns and mints
//    the token and holds its unsold supply (treasury);
//  - the token contract (contractAddress): the B20 token itself.
const Wallets: React.FC<WalletProps> = ({ asset }) => {
  const issuingSafe = asset?.issuingWalletAddress;
  const tokenContract = asset?.contractAddress;
  // the contract can be (re)registered until the asset is minted
  const canRegister =
    !!asset?.id && !!issuingSafe && (asset?.assetTokenizationStatus ?? 0) <= 3;

  const [contractInput, setContractInput] = useState("");
  const [registerTokenContract, { isLoading: isRegistering }] =
    useRegisterTokenContractMutation();

  const { data: issuingSafeData } = useGetWalletBalancesQuery(issuingSafe!, {
    skip: !issuingSafe,
  });

  // Helper function to extract the amount for a given asset code.
  // For NGN, assetCode is "" and for USD we look for "USDT"
  const getBalance = (walletData: any, assetCode: string): string =>
    walletData?.data?.claimed?.find((item: any) => item.assetCode === assetCode)
      ?.amount || "0";

  const onRegister = async () => {
    if (!asset?.id) return;
    if (!isAddress(contractInput)) {
      showErrorToast("Enter the token contract's 0x-prefixed address.");
      return;
    }
    try {
      await registerTokenContract({
        tokenizedAssetID: asset.id,
        contractAddress: contractInput.trim(),
      }).unwrap();
      setContractInput("");
      showSuccessToast("Token contract verified and registered.");
    } catch (error: any) {
      const message =
        error?.data?.message ||
        error?.data?.error ||
        error?.message ||
        "The token contract could not be registered.";
      const match = String(message).match(/message:(.*?)(\]$|$)/);
      showErrorToast(match?.[1]?.trim() || String(message));
    }
  };

  return (
    <Container>
      <Heading>Wallets</Heading>
      <WalletsContainer>
        <WalletCard>
          <WalletInfo>
            <WalletType>Issuing Safe (minter &amp; treasury)</WalletType>
            <WalletOwner title={issuingSafe}>
              {asset?.issuingWalletAlias || "Not assigned"}
            </WalletOwner>
            <Address>{issuingSafe || "—"}</Address>
            <BalanceLabel>Total Balance</BalanceLabel>
            <Balance>{getBalance(issuingSafeData, "")} NGN</Balance>
            <BalanceUsd>{getBalance(issuingSafeData, "USDT")} USD</BalanceUsd>
          </WalletInfo>
          <ImageWrapper>
            <Image src={logo} alt="trovo-logo" />
          </ImageWrapper>
        </WalletCard>

        <WalletCard>
          <WalletInfo>
            <WalletType>Token Contract</WalletType>
            <WalletOwner>
              {tokenContract ? asset?.assetCode : "Not registered"}
            </WalletOwner>
            <Address>{tokenContract || "—"}</Address>
            {canRegister && (
              <>
                <Hint>
                  Deploy the {asset?.assetCode || "asset"} B20 token with the
                  issuing Safe as its owner / MINTER_ROLE holder, symbol{" "}
                  {asset?.assetCode || "= asset code"} and zero supply, then
                  register it here. It is verified on-chain before minting can
                  start.
                </Hint>
                <RegisterRow>
                  <AddressInput
                    placeholder="0x… token contract address"
                    value={contractInput}
                    onChange={(e) => setContractInput(e.target.value)}
                  />
                  <RegisterButton
                    type="button"
                    disabled={isRegistering || !contractInput}
                    onClick={onRegister}
                  >
                    {isRegistering
                      ? "Verifying…"
                      : tokenContract
                        ? "Replace"
                        : "Register"}
                  </RegisterButton>
                </RegisterRow>
              </>
            )}
          </WalletInfo>
        </WalletCard>
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
const Address = styled.p`
  font-size: 12px;
  font-family: monospace;
  color: #00225ab2;
  margin: 0 0 8px;
  overflow-wrap: anywhere;
`;
const Hint = styled.p`
  font-size: 12px;
  line-height: 18px;
  color: #00225ab2;
  margin: 4px 0 8px;
`;
const RegisterRow = styled.div`
  display: flex;
  gap: 8px;
`;
const AddressInput = styled.input`
  flex: 1;
  min-width: 0;
  padding: 8px 10px;
  border: 1px solid #007cdf;
  border-radius: 8px;
  font-family: monospace;
  font-size: 12px;
`;
const RegisterButton = styled.button`
  padding: 8px 14px;
  border: none;
  border-radius: 8px;
  background: #007cdf;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
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
