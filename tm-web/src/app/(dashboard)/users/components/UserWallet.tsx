import React, { useState } from "react";
import logo from "@/assets/images/g647.svg";
import WhiteLogo from "@/assets/images/whitetrovoicon.svg";
import styled from "styled-components";
import Image from "next/image";
import Link from "next/link";
import { FaAngleRight } from "react-icons/fa6";
import { useRouter } from "next/navigation";
import { IUserInfo, IWallet } from "@/redux/api/users";

interface WalletCardProps {
  blueBackground?: boolean;
}

interface UserWalletProps {
  wallets: IWallet[];

  userInfo: IUserInfo;
}

const UserWallet: React.FC<UserWalletProps> = ({ wallets, userInfo }) => {
  const router = useRouter();

  if (!wallets || wallets.length === 0) {
    return <div>No wallets found</div>;
  }
  const sortedWallets = [...wallets]
    .sort((a, b) => b.primaryWallet - a.primaryWallet)
    .slice(0, 4);
  return (
    <>
      <HeadingContent>
        {" "}
        <Heading>Wallets</Heading>
        <StyledButton
          onClick={() => router.push(`/users/${userInfo.username}/wallets`)}
        >
          See all
        </StyledButton>
      </HeadingContent>
      <WalletsContainer>
        {sortedWallets.map((wallet) => (
          <StyledLink
            href={`/users/${userInfo.username}/wallets/${wallet.alias}`}
          >
            <WalletCard
              blueBackground={wallet.primaryWallet === 1}
              key={wallet.address}
            >
              <WalletOwnerContainer>
                <WalletOwner>{wallet.alias}</WalletOwner>
                <StyledLink
                  href={`/users/${userInfo.username}/wallets/${wallet.alias}`}
                >
                  <FaAngleRight
                    color={wallet.primaryWallet === 1 ? "#ffffff" : "#00225A"}
                  />
                </StyledLink>
              </WalletOwnerContainer>

              <BalanceLabel>Total Balance</BalanceLabel>
              <WalletInfo>
                <div>
                  <Balance>2,082,898 NGN</Balance>
                  <BalanceUsd>4,014 USD</BalanceUsd>
                </div>
                <ImageWrapper>
                  <Image
                    src={wallet.primaryWallet === 1 ? WhiteLogo : logo}
                    alt="trovo-logo"
                  />
                </ImageWrapper>
              </WalletInfo>
            </WalletCard>
          </StyledLink>
        ))}
      </WalletsContainer>
    </>
  );
};

export default UserWallet;

const WalletsContainer = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(285px, 1fr));
  gap: 20px;
  background: "#fff";
  padding-top: 10px;
  margin-bottom: 20px;
  overflow-x: hidden;
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

const WalletCard = styled.div<WalletCardProps>`
  width: 285px;
  padding: 16px;
  gap: 16px;
  border-radius: 20px;
  background-color: ${(props) =>
    props.blueBackground ? "#007CDF" : "#E8ECF4"};
  box-sizing: border-box;
  color: ${(props) => (props.blueBackground ? "#ffffff" : "#00225AB2")};

  img {
    filter: brightness(1.2);
  }
`;

const WalletOwner = styled.p`
  font-size: 16px;
  font-weight: 600;
  line-height: 25.27px;
  letter-spacing: 0.019em;
  margin-bottom: 4px;
`;

const WalletOwnerContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;
const BalanceLabel = styled.p`
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.1px;
  margin-top: 10px;
`;

const Balance = styled.p`
  font-size: 20px;
  font-weight: 600;
  line-height: 24px;
  letter-spacing: 0.1px;
  text-align: left;
  font-weight: bold;
  margin: 0px;
`;

const BalanceUsd = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0px;
`;

const WalletInfo = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const StyledLink = styled(Link)`
  text-decoration: none;
`;

const StyledButton = styled.button`
  font-weight: 500;
  font-size: 14px;
  line-height: 20px;
  letter-spacing: 0.1px;
  color: #007cdf;
  background-color: transparent;
  border: none;
  font-family: inherit;
  cursor: pointer;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  color: #00225a;
  margin-bottom: 10px;
`;

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 0;
`;
