import React, { useEffect, useRef, useState } from "react";
import { Navbar } from "./navbar";
import { Footer } from "./footer";
import heroImage from "../utils/assets/images/hero-image.png";
import multiImage from "../utils/assets/images/multi-image.png";
import maleAvatar from "../utils/assets/images/male-avatar.svg";
import accountRecovery from "../utils/assets/images/accountRecovery.svg";
import digitalAssets from "../utils/assets/images/digitalAssets.svg";
import multiWallet from "../utils/assets/images/multiWallet.svg";
import multiWallet1 from "../utils/assets/images/multiWallet1.png";
import multiWallet2 from "../utils/assets/images/multiWallet2.png";
import multiWallet3 from "../utils/assets/images/multiWallet3.png";
import multiWallet4 from "../utils/assets/images/multiWallet4.png";
import multiWallet5 from "../utils/assets/images/multiWallet5.png";
import multiWallet6 from "../utils/assets/images/multiWallet6.png";
import multiWallet7 from "../utils/assets/images/multiWallet7.png";
import downloadPromptImage from "../utils/assets/images/downloadImage.png";
import nfts from "../utils/assets/images/nfts.svg";
import playStoreIcon from "../utils/assets/images/googlePlayIcon.svg";
import iosAppIcon from "../utils/assets/images/appStoreIcon.svg";
import swapAssets from "../utils/assets/images/swapAssets.svg";
import downloadInstruction1 from "../utils/assets/images/downloadInstruction1.svg";
import downloadInstruction2 from "../utils/assets/images/downloadInstruction2.svg";
import downloadInstruction3 from "../utils/assets/images/downloadInstruction3.svg";
import downloadInstruction4 from "../utils/assets/images/downloadInstruction4.svg";
import tokenizedAssets from "../utils/assets/images/tokenizedAssets.svg";
import { styled } from "styled-components";
import { theme } from "../utils/theme";
import { ServiceSlider } from "../utils/slider";
import Paginate from "./pagination";
import { motion } from "framer-motion";
// import { useLocation } from "react-router-dom";

export const LandingPage = ({ scrollTo }) => {
  const [activeStep, setActiveStep] = useState(1);
  const [currentPage, setCurrentPage] = useState(1);
  const [stepsPerPage] = useState(4);
  const downloadRef = useRef(null);

  useEffect(() => {
    if (scrollTo) {
      const el = document.getElementById(scrollTo);
      if (el) {
        el.scrollIntoView({ behavior: "smooth" });
      }
    }
  }, [scrollTo]);

  // const location = useLocation();

  // useEffect(() => {
  //   if (location.hash) {
  //     const id = location.hash.replace("#", "");
  //     setTimeout(() => {
  //       const el = document.getElementById(id);
  //       if (el) el.scrollIntoView({ behavior: "smooth" });
  //     }, 100);
  //   }
  // }, [location.hash]);

  const trovoServices = [
    {
      image: tokenizedAssets,
      title: "Tokenized Assets",
      description: "Tokenize real assets and expose them to the global market.",
    },
    {
      image: multiWallet,
      title: "Multi Wallet",
      description: "Run your business with multi access wallets.",
    },
    {
      image: digitalAssets,
      title: "Digital Assets",
      description: "Become a merchant and make market for digital assets",
    },
    {
      image: swapAssets,
      title: "Swap Assets",
      description: "Swap digital assets on the go.",
    },
    {
      image: accountRecovery,
      title: "Account Recovery",
      description: "Recover account even when you lose your secret keys .",
    },
    {
      image: nfts,
      title: "Earn Passive Income",
      description: "Refer users and share in our transaction fess for life.",
    },
  ];

  const tokenSteps = [
    {
      id: 1,
      image: multiImage,
      title: "STEP 1",
      description:
        "Asset owners/Project promoters apply to tokenize assets on the Trovo App.",
    },
    {
      id: 2,
      image: multiWallet1,
      title: "STEP 2",
      description:
        "Due diligence is carried out on the asset and the owner to meet all regulatory requirements.",
    },
    {
      id: 3,
      image: multiWallet2,
      title: "STEP 3",
      description:
        "If approval to tokenize is given by the SEC, digital tokens that represent the real asset are created on the blockchain.",
    },
    {
      id: 4,
      image: multiWallet3,
      title: "STEP 4",
      description:
        "The digital tokens are sold to vetted and verified investors in primary offering",
    },
    {
      id: 5,
      image: multiWallet4,
      title: "STEP 5",
      description:
        "The real asset underlying the tokens are held by a custodian or trust, which is managed by a third party. The custodian/trust ensures the assets are properly cared for and that the rights of token holders are protected.The real asset underlying the tokens are held by a custodian or trust, which is managed by a third party. The custodian/trust ensures the assets are properly cared for and that the rights of token holders are protected.",
    },
    {
      id: 6,
      image: multiWallet5,
      title: "STEP 6",
      description:
        "Investors can then trade the digital tokens on a global secondary market which allows for increased liquidity in the market for the asset./",
    },
    {
      id: 7,
      image: multiWallet6,
      title: "STEP 7",
      description:
        "Proceeds are paid out to investors as the underlying asset generates revenue.",
    },
    {
      id: 8,
      image: multiWallet7,
      title: "STEP 8",
      description:
        "Proceeds of sale are distributed to investors when underlying asset is sold.",
    },
  ];
  const handleClick = (id) => {
    setActiveStep(id);
  };
  const indexOfLastStep = currentPage * stepsPerPage;
  const indexOfFirstStep = indexOfLastStep - stepsPerPage;
  const currentSteps = tokenSteps.slice(indexOfFirstStep, indexOfLastStep);
  const paginate = (pageNumber) => {
    setCurrentPage(pageNumber);
  };
  const previousPage = () => {
    if (currentPage !== 1) {
      setCurrentPage(currentPage - 1);
    }
  };
  console.log({ indexOfFirstStep });
  const nextPage = () => {
    if (currentPage !== Math.ceil(tokenSteps.length / stepsPerPage)) {
      setCurrentPage(currentPage + 1);
    }
  };
  const homeVariants = {
    hide: {
      opacity: 0,
      y: -150,
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
  const featuresVariants = {
    hide: {
      x: -40,
      y: -20,
      opacity: 0,
    },
    show: {
      x: 0,
      y: 0,
      opacity: 1,
      transition: {
        duration: 1,
      },
    },
  };

  return (
    <div>
      <Navbar />
      <HeroSection
        as={motion.section}
        variants={homeVariants}
        initial={"hide"}
        exit={"hide"}
        whileInView={"show"}
        id="home"
      >
        <HeroTextSection>
          <GradientBox></GradientBox>
          <HeroHeading>
            Unlock liquidity in illiquid assets and transcend  borders with{" "}
            <TrovoText>
              {" "}
              Trovo App <MaleAvatar src={maleAvatar} />{" "}
            </TrovoText>
          </HeroHeading>

          <HeroDescription>
            A decentralized and non-custodial (self custody) wallet with a lot
            of innovative features for asset tokenization, secure storage,
            transfer, and exchange of digital assets.
          </HeroDescription>
          <ButtonContainer>
            <TrialButton
              as="button"
              onClick={() => {
                if (downloadRef.current) {
                  downloadRef.current.scrollIntoView({ behavior: "smooth" });
                }
              }}
            >
              Join Now
            </TrialButton>
          </ButtonContainer>
        </HeroTextSection>
        <HeroImageSection>
          <HeroImage src={heroImage} />
        </HeroImageSection>
      </HeroSection>
      <ServiceSection
        as={motion.div}
        variants={serviceVariants}
        initial="hide"
        whileInView="show"
        exit="hide"
        id="features"
      >
        <ServiceHeading>Do more with Trovo App</ServiceHeading>
        <ServiceCardContainer
          as={motion.div}
          variants={serviceVariants}
          initial="hide"
          whileInView="show"
          exit="hide"
        >
          {trovoServices.map((trovoService, index) => {
            return (
              <ServiceCard
                as={motion.div}
                variants={serviceVariants}
                initial="hide"
                whileInView="show"
                exit="hide"
              >
                <ServiceImage src={trovoService.image} />
                <ServiceTitle>{trovoService.title}</ServiceTitle>
                <ServiceText>{trovoService.description}</ServiceText>
              </ServiceCard>
            );
          })}
        </ServiceCardContainer>
        <MobileServiceContainer>
          <ServiceSlider />
        </MobileServiceContainer>
      </ServiceSection>
      <MultiSection ref={downloadRef} id="download">
        <MultiHeader>Asset Tokenization</MultiHeader>
        <MultiPhraseParagraph>
          Asset Tokenization is the process of representing physical or
          intangible assets as digital tokens on a blockchain, allowing the
          assets to be easily bought, sold, and traded on a global market. This
          is providing capital for asset owners and investment opportunity for
          everyone.
        </MultiPhraseParagraph>
        <TokenizationHeading>How It Works</TokenizationHeading>

        <MultiRow>
          <MultiImageGroup
            as={motion.div}
            variants={featuresVariants}
            initial="hide"
            whileInView="show"
            exit="hide"
          >
            <MultiImage src={tokenSteps[activeStep - 1].image} />
          </MultiImageGroup>
          <MultiTextRow>
            {indexOfFirstStep !== 0 && (
              <Paginate
                stepsPerPage={stepsPerPage}
                previousPage={previousPage}
                nextPage={nextPage}
                totalSteps={tokenSteps.length}
                paginate={paginate}
                showLastPage={true}
                showNextPage={false}
              />
            )}
            {currentSteps.map((tokenStep, index) => {
              console.log(activeStep, index);
              return (
                <MultiTextGroup
                  isActive={activeStep === tokenStep.id}
                  onClick={() => handleClick(tokenStep.id)}
                >
                  <MultiHeading isActive={activeStep === tokenStep.id}>
                    {tokenStep.title}
                  </MultiHeading>
                  <MultiPhrase isActive={activeStep === tokenStep.id}>
                    {tokenStep.description}
                  </MultiPhrase>
                </MultiTextGroup>
              );
            })}
            {indexOfFirstStep === 0 && (
              <Paginate
                stepsPerPage={stepsPerPage}
                previousPage={previousPage}
                nextPage={nextPage}
                totalSteps={tokenSteps.length}
                paginate={paginate}
                showNextPage={true}
                showLastPage={false}
              />
            )}
          </MultiTextRow>
        </MultiRow>
        {/* <div style={{width:'100%', display: 'flex', alignItems: 'center', justifyContent: 'center', marginTop: '60px'}}>
            <TrialButton target="_blank" href="https://trovo-app-web-uidfp.ondigitalocean.app/">Launch  Now</TrialButton>
          </div> */}
      </MultiSection>
      <DownloadSection>
        <DownloadRow>
          <DownloadTextSection>
            <DownloadHeading>
              Embrace the Future of Asset Tokenization
            </DownloadHeading>
            <DownloadPhrase>
              Step into the world of limitless possibilities with Trovo App, the
              ultimate blockchain-based non-custodial multi-access wallet. Take
              full control of your digital assets, access global markets, and
              experience the power of secure and transparent transactions like
              never before.
            </DownloadPhrase>
            <DownloadInstructions>
              <DownloadIcons src={downloadInstruction4} />
              <InstructionText>Download Trovo App.</InstructionText>
            </DownloadInstructions>
            <DownloadInstructions>
              <DownloadIcons src={downloadInstruction1} />
              <InstructionText>
                Register an account and backup your secret key.
              </InstructionText>
            </DownloadInstructions>
            <DownloadInstructions>
              <DownloadIcons src={downloadInstruction2} />
              <InstructionText>
                Join the revolution and unleash the true potential of blockchain
                technology.
              </InstructionText>
            </DownloadInstructions>
            <DownloadInstructions>
              <DownloadIcons src={downloadInstruction3} />
              <InstructionText>
                Refer your friends and earn commission from our fees.
              </InstructionText>
            </DownloadInstructions>
            <DownloadButtonContainer>
              <MobileAppButton
                href="https://testflight.apple.com/join/7ka7gHgQ"
                target="_blank"
                rel="noopener noreferrer"
              >
                <ButtonImage src={iosAppIcon} />
                <ButtonTextSection>
                  <CtaText>Download on the</CtaText>
                  <AppName> App Store</AppName>
                </ButtonTextSection>
              </MobileAppButton>
              <MobileAppButton
                href="https://play.google.com/store/apps/details?id=com.trovo.wallet"
                target="_blank"
                rel="noopener noreferrer"
              >
                <ButtonImage src={playStoreIcon} />
                <ButtonTextSection>
                  <CtaText>Get it on</CtaText>
                  <AppName>Google Play</AppName>
                </ButtonTextSection>
              </MobileAppButton>
            </DownloadButtonContainer>
          </DownloadTextSection>
          <DownloadImageSection>
            <DownloadBanner src={downloadPromptImage} />{" "}
          </DownloadImageSection>
        </DownloadRow>
      </DownloadSection>
      <Footer />
    </div>
  );
};

const HeroSection = styled.section`
  display: flex;
  column-gap: 25px;
  margin: 40px 80px 80px 120px;
  align-items: center;
  padding-top: 100px;

  @media (min-width: 1450px) {
    max-width: 1440px;
    margin: auto;
  }
  @media (max-width: 1000px) {
    padding-top: 30px;
  }

  @media (max-width: 992px) {
    flex-direction: column;
    row-gap: 40px;
  }
  @media (max-width: 600px) {
    margin: 48px 30px 40px 60px;
  }
`;
const HeroTextSection = styled.div`
  position: relative;
  flex: 1;
`;
const GradientBox = styled.div`
  border-radius: 80px 0px 0px 80px;
  background: linear-gradient(90deg, #acd1ef 0%, #e7f5fd 100%);
  box-shadow: 4px 0px 139px 0px rgba(0, 73, 136, 0.25);
  width: 167px;
  height: 90px;
  position: absolute;
  top: -3%;
  left: -15%;
  z-index: -1;
  @media (max-width: 1690px) {
    left: -3%;
  }
  @media (max-width: 1440px) {
    left: -15%;
  }

  @media (max-width: 768px) {
    width: 150px;
    height: 51px;
    top: -2%;
    left: -9%;
  }
`;
const HeroHeading = styled.h1`
  color: ${theme.colors.trovoNavyBlue};
  font-size: 44px;
  font-family: "MatahariExtended";
  font-style: normal;
  font-weight: 900;
  line-height: 72px;
  letter-spacing: 0.169px;
  margin-top: 0;

  @media (max-width: 768px) {
    font-size: 36px;
    line-height: 140%;
  }
`;
const TrovoText = styled.span`
  border-radius: 200px;
  background: linear-gradient(85deg, #00a859 0%, rgba(94, 199, 149, 0.56) 100%);
  width: 352px;
  height: 80px;
  display: inline-flex;
  justify-content: end;
  align-items: center;
  padding: 0;
  column-gap: 15px;
  color: #fff;

  @media (max-width: 768px) {
    width: 310px;
    height: 51px;
  }
`;
const MaleAvatar = styled.img`
  @media (max-width: 768px) {
    width: 50px;
    height: 50px;
  }
`;
const HeroDescription = styled.p`
  color: #000;
  font-size: 18px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  margin-bottom: 24px;

  @media (max-width: 600px) {
    font-size: 14px;
  }
`;
const HeroImageSection = styled.div`
  flex: 1;
`;
const HeroImage = styled.img`
  height: auto;
  max-width: 100%;
  object-fit: contain;
`;
const ServiceSection = styled.section`
  padding-top: 50px;
`;
const ServiceHeading = styled.h2`
  color: ${theme.colors.trovoBlue};
  font-size: 32px;
  font-family: "MatahariExtended";
  font-style: normal;
  font-weight: 800;
  line-height: 52px;
  letter-spacing: 0.2px;
  margin: 32px auto;
  text-align: center;

  @media (max-width: 600px) {
    font-size: 24px;
  }
`;
const MobileServiceContainer = styled.div`
  width: 100%;
  margin: auto;
`;
const ServiceCardContainer = styled.div`
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  row-gap: 24px;
  max-width: 1126px;
  margin: auto;
  column-gap: 80px;

  @media (max-width: 700px) {
    display: none;
  }
`;
const ServiceCard = styled.div`
  padding: 24px 20px;
  border-radius: 10px;
  background: rgba(0, 73, 136, 0.03);
  max-width: 280px;
  width: calc(30.33% - 20px);
`;
const ServiceImage = styled.img`
  height: auto;
  max-width: 100%;
  object-fit: contain;
`;
const ServiceTitle = styled.h5`
  color: #1b1d21;
  font-size: 18px;
  font-family: "Poppins";
  font-style: normal;
  font-weight: 600;
  line-height: 28px;
  margin: 12px 0;

  @media (max-width: 600px) {
    font-size: 16px;
  }
`;
const ServiceText = styled.p`
  color: #000;
  font-size: 14px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 400;
  line-height: 28px;
  opacity: 0.6000000238418579;
  margin: 0;
`;
const MultiSection = styled.section`
  padding: 80px;
  width: 75%;
  margin: auto;

  @media (min-width: 1450px) {
    max-width: 1440px;
    margin: auto;
  }
  @media (max-width: 600px) {
    padding: 24px 16px 50px;
  }
`;
const MultiHeader = styled.h3`
  color: ${theme.colors.trovoBlue};
  font-size: 32px;
  font-family: "MatahariExtended";
  font-style: normal;
  font-weight: 800;
  line-height: 52px;
  letter-spacing: 0.2px;
  margin: 0;
  text-align: center;

  @media (max-width: 600px) {
    font-size: 24px;
  }
`;

const MultiPhraseParagraph = styled.p`
  font-family: "Montserrat";
  font-size: 18px;
  font-weight: 400;
  line-height: 28px;
  text-align: center;
  color: #191919;
`;
const TokenizationHeading = styled.h5`
  font-family: "Montserrat";
  font-size: 24px;
  font-weight: 700;
  line-height: 52px;
  letter-spacing: 0.20000000298023224px;
  text-align: left;
  color: #00225a;
  text-align: center;
`;
const MultiRow = styled.div`
  display: flex;
  justify-content: space-between;
  column-gap: 50px;
  align-items: center;

  @media (max-width: 768px) {
    flex-direction: column;
    row-gap: 24px;
  }
`;
const MultiImageGroup = styled.div`
  /* flex: 1; */
  margin-right: 80px;
`;
const MultiImage = styled.img`
  width: 100%;
  height: 100%;
  object-fit: contain;
`;

const MultiTextRow = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;

  @media (max-width: 768px) {
    text-align: center;
  }
`;

const MultiTextGroup = styled.div`
  border-left: ${({ isActive }) =>
    isActive ? "3px solid #26649A" : "3px solid #f2f6f9;"};
  padding-left: 80px;
  padding: 25px 0px 25px 80px;
  cursor: pointer;
`;
const MultiHeading = styled.h5`
  color: ${({ isActive }) => (isActive ? "#004988" : "#7ba0bc")};
  font-size: 14px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 700;
  line-height: 30px;
  margin: 0;
  opacity: ${({ isActive }) => (isActive ? "unset" : "0.6000000238418579")};

  @media (max-width: 600px) {
    font-size: 16px;
  }
`;
const MultiPhrase = styled.p`
  color: ${({ isActive }) => (isActive ? "#191919" : "#000")};
  font-size: 14px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 400;
  line-height: 30px;
  margin: 0;
  opacity: ${({ isActive }) => (isActive ? "unset" : "0.6000000238418579")};

  @media (max-width: 600px) {
    font-size: 16px;
  }
`;
const DownloadSection = styled.section`
  padding: 78px 120px;

  @media (min-width: 1440px) {
    max-width: 1440px;
    margin: auto;
  }
  @media (max-width: 1250px) {
    padding: 78px 100px;
  }
  @media (max-width: 1150px) {
    padding: 78px 40px;
  }
  @media (max-width: 600px) {
    padding: 80px 20px;
  }
`;
const DownloadRow = styled.div`
  display: flex;
  column-gap: 69px;
  background-color: #d5e7f6;
  padding: 46px 0 0px 80px;
  border-radius: 24px;

  @media (max-width: 1250px) {
    padding-left: 40px;
  }
  @media (max-width: 1150px) {
    flex-direction: column-reverse;
    padding: 40px;
  }
  @media (max-width: 600px) {
    padding: 24px 16px 50px;
  }
`;
const DownloadTextSection = styled.div`
  flex: 1;
`;
const DownloadImageSection = styled.div`
  flex: 1;
  text-align: end;

  @media (max-width: 1150px) {
    margin: auto;
  }
`;
const DownloadHeading = styled.h4`
  color: ${theme.colors.trovoNavyBlue};
  font-feature-settings: "clig" off, "liga" off;
  font-family: "MatahariExtended";
  font-size: 32px;
  font-style: normal;
  font-weight: 800;
  line-height: 52px; /* 162.5% */
  letter-spacing: 0.2px;
  margin: 0;
`;
const DownloadPhrase = styled.p`
  color: ${theme.colors.trovoNavyBlue};
  font-size: 14px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 500;
  line-height: 24px;
`;
const DownloadInstructions = styled.div`
  display: flex;
  column-gap: 7px;
`;
const DownloadIcons = styled.img``;
const InstructionText = styled.p`
  color: ${theme.colors.trovoNavyBlue};
  font-family: "Montserrat";
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
`;
const DownloadButtonContainer = styled.div`
  display: flex;
  column-gap: 24px;
  margin: 20px 0;

  @media (max-width: 690px) {
    flex-direction: column;
    row-gap: 20px;
  }
`;
const ButtonContainer = styled.div`
  @media (max-width: 768px) {
    width: 100%;
    margin: auto;
    display: flex;
    justify-content: center;
  }
`;
const DownloadBanner = styled.img`
  margin-top: -125px;
  height: auto;
  max-width: 100%;
  object-fit: contain;
`;
const MobileAppButton = styled.a`
  display: flex;
  border: none;
  border-radius: 12px;
  column-gap: 10px;
  align-items: center;
  padding: 8px 28px;
  justify-content: center;
  cursor: pointer;
  text-decoration: none;

  background: ${theme.colors.trovoBlue};
`;
const CtaText = styled.p`
  margin: 0;
  color: ${theme.colors.white};
  font-family: "Montserrat";
  font-size: 12px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  letter-spacing: 0.1px;
  text-align: left;
`;
const AppName = styled.h6`
  margin: 0;
  color: ${theme.colors.white};
  font-family: "Montserrat";
  font-size: 14px;
  font-style: normal;
  font-weight: 600;
  line-height: 20px;
  text-align: left;
  letter-spacing: 0.1px;
`;
const ButtonTextSection = styled.div``;
const ButtonImage = styled.img`
  height: auto;
  max-width: 100%;
  object-fit: contain;
`;
const TrialButton = styled.a`
  padding: 12px 10px;
  background-color: ${theme.colors.trovoBlue};
  color: ${theme.colors.white};
  font-size: 14px;
  font-family: "Montserrat";
  border: none;
  border-radius: 10px;
  font-style: normal;
  font-weight: 600;
  line-height: 26px;
  width: 230px;
  display: block;
  text-align: center;
  text-decoration: none;
  cursor: pointer;
`;
