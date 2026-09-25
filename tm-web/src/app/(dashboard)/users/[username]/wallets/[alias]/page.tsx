"use client";

import React from "react";
import styled from "styled-components";

import WhiteLogo from "@/assets/images/whitetrovoicon.svg";
import Image from "next/image";

import asset from "@/assets/images/assetsImg.png";
import { FaArrowLeft } from "react-icons/fa6";
import { useParams, useRouter } from "next/navigation";
import WalletsComponent from "../../../components/WalletsComponent";
import SharedAccess from "../../../components/SharedAccess";
import UserAssets from "../../../components/UserAssets";
import { useFetchUserByCriteriaQuery } from "@/redux/api/users";

const WalletDetailsPage = () => {
 
  const { username, alias } = useParams() as {
    username: string;
    alias: string;
  };

  // Fetch user details using the username parameter.
  const { data, isLoading, isError } = useFetchUserByCriteriaQuery({
    username,
  });
  const router = useRouter();

  if (isLoading) return <div>Loading...</div>;
  if (isError || !data) return <div>Error fetching data</div>;

  // Find the wallet that matches the alias from the URL.
  const wallet = data.data.wallets.find(
    (wallet: any) => wallet.alias === alias
  );

  if (!wallet) return <div>No wallet found with alias: {alias}</div>;

  // Get the wallet's public key from the fetched data.
  const publicKey = wallet.publicKey;

  const access = [
    {
      name: "Obi Enechi",
      username: "obi",
      image: "",
    },

    {
      name: "Obi Enechi",
      username: "obi",
      image: "",
    },

    {
      name: "Obi Enechi",
      username: "obi",
      image: "",
    },

    // Add more asset objects as needed
  ];

  const assetTokens = [
    {
      name: "AFT",
      value: "100 CGCN",
      image: asset,
      amount: "2,900",
      val: "209,300CNGN",
    },
    {
      name: "AE1",
      value: "100 CGCN",
      image: asset,
      amount: "2,900",
      val: "209,300CNGN",
    },
    {
      name: "AFT",
      value: "100 CGCN",
      image: asset,
      amount: "2,900",
      val: "209,300CNGN",
    },
    {
      name: "AE1",
      value: "100 CGCN",
      image: asset,
      amount: "2,900",
      val: "209,300CNGN",
    },

    {
      name: "AE1",
      value: "100 CGCN",
      image: asset,
      amount: "2,900",
      val: "209,300CNGN",
    },
    {
      name: "AFT",
      value: "100 CGCN",
      image: asset,
      amount: "2,900",
      val: "209,300CNGN",
    },
    {
      name: "AE1",
      value: "100 CGCN",
      image: asset,
      amount: "2,900",
      val: "209,300CNGN",
    },

    {
      name: "AFT",
      value: "100 CGCN",
      image: asset,
      amount: "2,900",
      val: "209,300CNGN",
    },
  ];

  const otherAssets = [
    {
      name: "CVT",
      value: "100 CGCN",
      image: "/path/to/cvt-image.svg",
      amount: "2,900",
      val: "209,300CNGN",
    },
    {
      name: "DFT",
      value: "100 CGCN",
      image: "/path/to/dft-image.svg",
      amount: "2,900",
      val: "209,300CNGN",
    },

    {
      name: "AFT",
      value: "100 CGCN",
      image: asset,
      amount: "2,900",
      val: "209,300CNGN",
    },
    {
      name: "AE1",
      value: "100 CGCN",
      image: asset,
      amount: "2,900",
      val: "209,300CNGN",
    },
    {
      name: "AFT",
      value: "100 CGCN",
      image: asset,
      amount: "2,900",
      val: "209,300CNGN",
    },
    {
      name: "AE1",
      value: "100 CGCN",
      image: asset,
      amount: "2,900",
      val: "209,300CNGN",
    },
    {
      name: "DFT",
      value: "100 CGCN",
      image: "/path/to/dft-image.svg",
      amount: "2,900",
      val: "209,300CNGN",
    },
  ];

  const handleGoBack = () => {
    // console.log("Navigating back...");
    // console.log("History length:", window.history.length);
    try {
      router.back();
    } catch (error) {
      //   console.error("Error while navigating back:", error);
      router.push("/users");
    }
  };

  return (
    <>
      <BackIcon onClick={handleGoBack} />

      <Container>
        <LeftPanel>
          <WalletDetails>
            <WalletsComponent title="Alias" content={alias} />
            <WalletsComponent title="Public Key" content={publicKey} />
            <WalletCard>
              <div>
                <BalanceLabel>Total Balance</BalanceLabel>

                <Balance>2,082,898 NGN</Balance>
                <BalanceUsd>4,014 USD</BalanceUsd>
              </div>
              <ImageWrapper>
                <Image
                  src={WhiteLogo}
                  alt="trovo-logo"
                  width={500}
                  height={500}
                />
              </ImageWrapper>
            </WalletCard>
          </WalletDetails>
          <SharedAccess access={access} />
        </LeftPanel>
        <RightPanel>
          <UserAssets assetTokens={assetTokens} otherAssets={otherAssets} />
        </RightPanel>
      </Container>
    </>
  );
};

export default WalletDetailsPage;

const Container = styled.div`
  padding: 30px;
  display: flex;

  gap: 30px;
  justify-content: space-between;
`;
const LeftPanel = styled.div`
  width: 50%;
`;
const RightPanel = styled.div`
  width: 70%;
`;
const WalletDetails = styled.div`
  background-color: #fff;
  border-radius: 24px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

const WalletCard = styled.div`
  padding: 16px;
  gap: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-radius: 20px;
  background-color: #007cdf;
  box-sizing: border-box;
  color: #ffffff;

  img {
    filter: brightness(1.2);
  }
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

const BackIcon = styled(FaArrowLeft)`
  color: #00225a;
  size: 25px;
  cursor: pointer;
  margin: 15px 30px 5px 30px;
`;
