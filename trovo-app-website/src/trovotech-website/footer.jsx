import React from "react";
import { styled } from "styled-components";
import { theme } from "../utils/theme";
import trovoLogo from "../utils/assets/images/footerLogo.svg";
import emailIcon from "../utils/assets/images/footerEmail.svg";
import facebookLogo from "../utils/assets/images/footerFacebook.svg";
import githubIcon from "../utils/assets/images/footerGithub.svg";
import locationIcon from "../utils/assets/images/footerLocation.svg";
import twitterLogo from "../utils/assets/images/footerTwitter.svg";
import instagramLogo from "../utils/assets/images/footerInstagram.svg";
import linkedinLogo from '../utils/assets/images/footerLinkedIn.svg'
// import linkedinLogo from "../utils/assets/images/footerLinkedIn.svg";
import { motion } from "framer-motion";

export const Footer = () => {
  const serviceVariants = {
    hide: {
      opacity: 0,
      y: 50,
    },
    show: {
      opacity: 1,
      y: 0,
      transition: {
        duration: 1,
        stiffness: 200,
      },
    },
  };
  return (
    <FooterLayout id="contact">
      <FooterContainer
            as={motion.div}
            variants={serviceVariants}
            initial="hide"
            whileInView="show"
            exit="hide">
        <ImageGroup>
          <LogoImage src={trovoLogo} />
          <SocialGroup>
            <SocialsLinks
              href="https://www.facebook.com/trovotech/"
              target="blank"
            >
              <SocialImages src={facebookLogo} />
            </SocialsLinks>
            <SocialsLinks
              href="https://twitter.com/Trovotechio/"
              target="blank"
            >
              <SocialImages src={twitterLogo} />
            </SocialsLinks>
            <SocialsLinks
              href="https://twitter.com/Trovotechio/"
              target="blank"
            >
              <SocialImages src={instagramLogo} />
            </SocialsLinks>
            <SocialsLinks
              href="https://www.linkedin.com/company/trovotech/"
              target="blank"
            >
              <SocialImages src={linkedinLogo} />
            </SocialsLinks>
            <SocialImages src={githubIcon} />
          </SocialGroup>
        </ImageGroup>
        <InfoGroup>
          <FooterTitle>Company</FooterTitle>
          <FooterText>About Us</FooterText>
          {/* <FooterText>News</FooterText> */}
          {/* <FooterText>Careers</FooterText> */}
          <FooterText
            href="https://drive.google.com/file/d/16O590h9g8XeHwp9tAjU0hJJ3Z1eYlySk/view"
            target="blank"
          >
            Documentation
          </FooterText>
        </InfoGroup>
        <InfoGroup>
          <FooterTitle>Quick Links</FooterTitle>
          <FooterText>Privacy Policy</FooterText>
          <FooterText>Terms of Use</FooterText>
          <FooterText>Demo</FooterText>
        </InfoGroup>
        <InfoGroup>
          <FooterTitle>Contact Us</FooterTitle>
          <ContactRow>
            <SocialImages src={locationIcon} />
            <FooterText>Lagos</FooterText>
          </ContactRow>
          <ContactRow>
            <SocialImages src={emailIcon} />
            <FooterText>info@trovotech.io</FooterText>
          </ContactRow>
        </InfoGroup>
      </FooterContainer>
      <CopyrightSection>
        <CopyrightText>Copyright 2023. All Rights Reserved.</CopyrightText>
      </CopyrightSection>
    </FooterLayout>
  );
};

const FooterLayout = styled.section`
  background-color: #1D212F;
`;
const FooterContainer = styled.div`
  display: flex;
  align-items: flex-start;
  column-gap: 120px;
  padding: 80px;
  flex-wrap: wrap;
  row-gap: 32px;

  @media (min-width: 1450px) {
    max-width: 1440px;
    margin: auto;
  }
  @media (max-width: 700px) {
    justify-content: center;
  }
  @media (max-width: 600px) {
    text-align: center;
  }
`;
const ImageGroup = styled.div`
  display: flex;
  flex-direction: column;
  row-gap: 64px;

  @media (max-width: 600px) {
    row-gap: 24px;
  }
`;
const LogoImage = styled.img``;
const SocialGroup = styled.div`
  display: flex;
  column-gap: 20px;
  justify-content: flex-end;
  width: 100%;
`;
const SocialImages = styled.img``;
const InfoGroup = styled.div`
  display: flex;
  flex-direction: column;
  row-gap: 24px;
`;

const FooterTitle = styled.h5`
  margin: 0;
  color: ${theme.colors.white};
  font-size: 18px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 600;
  line-height: normal;
  opacity: 0.800000011920929;
`;

const FooterText = styled.a`
  margin: 0;
  color: ${theme.colors.white};
  font-size: 14px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  opacity: 0.800000011920929;
  text-decoration: none;
`;
const CopyrightSection = styled.div`
  padding: 40px 80px;
  border-top: 2px solid rgba(255, 255, 255, 0.2);
  @media (min-width: 1450px) {
    max-width: 1440px;
    margin: auto;
  }
`;
const CopyrightText = styled.p`
  margin: 0;
  color: #fff;
  font-size: 14px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 500;
  line-height: 20px;
`;
const ContactRow = styled.div`
  display: flex;
  column-gap: 8px;

  @media (max-width: 600px) {
    text-align: center;
    justify-content: center;
  }
`;
const SocialsLinks = styled.a`
  text-decoration: none;
`;
