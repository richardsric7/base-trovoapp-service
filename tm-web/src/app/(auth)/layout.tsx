"use client";
import React from "react";
import styled from "styled-components";
import Logo from "@/assets/images/trovologosvg.svg";
import authImage from "@/assets/images/authdashboard.png";
import backgroundImage from "@/assets/images/authbg.png";
import Image from "next/image";
export default function AuthLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <AuthContainer>
      <AuthContentRow>
        <AuthWrapper>
          <LogoContainer>
            <Image src={Logo} alt="Trovotech- logo" />
          </LogoContainer>
          {children}
        </AuthWrapper>
        <InfoContainer>
          <Title>All Trovotech&apos;s Products Management in one place</Title>

          <StyledImageContainer>
            <StyledImage
              src={authImage}
              alt="Trovotech-Manager-Dashboard- illustration"
              width={510}
              height={300}
            />
          </StyledImageContainer>
        </InfoContainer>
      </AuthContentRow>
    </AuthContainer>
  );
}

const AuthContainer = styled.section`
  background: url(${backgroundImage.src}) center/cover no-repeat;
  background-color: #007cdf;
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  overflow: hidden;
  box-sizing: border-box;
`;

const AuthContentRow = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  align-items: center;
  width: 90%;
  max-width: 1100px;
  margin: auto;
  overflow: hidden;
  gap: 10px;
  border: 2px solid #fff;
  border-radius: 16px;
  height: 550px;

  box-sizing: border-box;
`;
const AuthWrapper = styled.div`
  background-color: #fff;
  height: 122%;
  width: 100%;
  padding: 30px;
  z-index: 1;
  display: flex;
  flex-direction: column;
  border-radius: 16px 0 0 16px;
`;

const InfoContainer = styled.div`
  height: 100%;
  padding: 30px;
  border-radius: 16px 0 0 16px;
  position: relative;
`;

const LogoContainer = styled.div`
  padding: 55px 10px 10px 10px;
  margin-top: 30px;
`;

const Title = styled.p`
  font-size: 40px;
  font-weight: 600;
  line-height: 48.76px;
  color: #ffffff;
`;

const StyledImageContainer = styled.div`
  border-top: 6px solid #191919;
  border-left: 6px solid #191919;
  border-bottom: 6px solid #191919;
  border-radius: 16px 0 0 16px;
  position: absolute;
  right: 0px;
  margin-top: 20px;
  max-width: 510px;
  width: 100%;
  overflow: hidden; /* Ensure content stays within the border */
  display: flex;
  justify-content: center;
  align-items: center;
`;

const StyledImage = styled(Image)`
  width: 100%;
  // height: auto;
`;
