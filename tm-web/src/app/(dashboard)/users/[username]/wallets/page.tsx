"use client";
import React from "react";
import logo from "@/assets/images/g647.svg";
import WhiteLogo from "@/assets/images/whitetrovoicon.svg";
import styled from "styled-components";
import Image from "next/image";
import Link from "next/link";
import { FaAngleRight, FaArrowLeft } from "react-icons/fa6";

import { useFetchUserByCriteriaQuery } from "@/redux/api/users";
import { useParams, useRouter } from "next/navigation";

interface WalletCardProps {
  blueBackground?: boolean;
}

export default function AllWallets() {
  const params = useParams();

  const username = params?.username
    ? Array.isArray(params.username)
      ? params.username[0]
      : params.username
    : undefined;

  const queryParams = username ? { username } : {};

  const { data, isLoading, isError } = useFetchUserByCriteriaQuery(queryParams);
  const router = useRouter();

  if (isLoading) return <div>Loading...</div>;
  if (isError) return <div>Error fetching wallets</div>;

  const wallets = data?.data?.wallets || [];
  if (!wallets || wallets.length === 0) {
    return <div>No wallets found</div>;
  }

  const sortedWallets = [...wallets].sort(
    (a, b) => b.primaryWallet - a.primaryWallet
  );
  const backLink = () => {
    router.back();
  };

  return (
    <Container>
      <HeadingContent>
        <BackLink onClick={backLink}>
          <FaArrowLeft />
        </BackLink>
      </HeadingContent>
      <Heading>Wallets</Heading>
      <WalletsContainer>
        {sortedWallets.map((wallet) => (
          <StyledLink
            key={wallet.address}
            href={`/users/${username}/wallets/${wallet.alias}`}
          >
            <WalletCard blueBackground={wallet.primaryWallet === 1}>
              <WalletOwnerContainer>
                <WalletOwner>{wallet.alias}</WalletOwner>
                <StyledLink href={`/users/${username}/wallets/${wallet.alias}`}>
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
    </Container>
  );
}

const Container = styled.section`
  background-color: #fff;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 24px;
`;

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
  margin-bottom: 10px;
`;

const BackLink = styled.button`
  text-decoration: none;
  display: block;
  color: #000000;
  cursor: pointer;
  padding-bottom: 10px;
  background-color: transparent;
  border: none;
  font-family: inherit;
`;
