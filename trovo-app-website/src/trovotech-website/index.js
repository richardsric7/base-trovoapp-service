import React from "react";
import { Footer } from "./footer";
import { Navbar } from "./navbar";
import { styled } from "styled-components";
import bannerLogo from "../utils/assets/images/banner_logo.png";
import service1 from "../utils/assets/images/service1.svg";
import service2 from "../utils/assets/images/service2.svg";
import service3 from "../utils/assets/images/service3.svg";
import service4 from "../utils/assets/images/service4.svg";
import service5 from "../utils/assets/images/service5.svg";
import service6 from "../utils/assets/images/service6.svg";
import buildIcon1 from "../utils/assets/images/buildIcon1.svg";
import buildIcon2 from "../utils/assets/images/buildIcon2.svg";
import buildIcon3 from "../utils/assets/images/buildIcon3.svg";
import trovoOffice from "../utils/assets/images/trovo office 1.png";
import trovoView from "../utils/assets/images/image 6.png";
import leftChevron from "../utils/assets/images/leftChevron.svg";
import rightChevron from "../utils/assets/images/rightChevron.svg";
import { motion } from "framer-motion";

export const TrovotechWebsite = () => {
  const servicesCardDetails = [
    {
      image: service1,
      cardTitle: "Digital Assets Services",
      cardDescription:
        "We design, develop and deploy the following components of the digital assets...",
    },
    {
      image: service2,
      cardTitle: "Blockchain System Integration",
      cardDescription:
        "Are you looking to integrate blockchain into aspects of your business processes to increase ...",
    },
    {
      image: service3,
      cardTitle: "Blockchain System Consulting",
      cardDescription:
        "We provide blockchain/DLT advisory services to help you analyze how your business can be impacted by shelf ..",
    },
    {
      image: service4,
      cardTitle: "Program/Project Management",
      cardDescription:
        "We have seasoned project management professionals who can help you deploy your business ...",
    },
    {
      image: service5,
      cardTitle: "Resource Planning  and Management",
      cardDescription:
        "We also excel at helping businesses allocate tasks to team members based. ",
    },
    {
      image: service6,
      cardTitle: "Blockchain System Support",
      cardDescription:
        "We design, develop and deploy the following components of the digital assets...",
    },
  ];
  const featuresContents = [
    {
      icon: buildIcon1,
      featuresTitle: "Experience",
      featuresDescription: "Over 200 years combined experience.",
      color: "rgba(0, 124, 223, 0.05)",
    },
    {
      icon: buildIcon2,
      featuresTitle: "Professionalism",
      featuresDescription:
        "Smooth user experience from engagement to delivery.",
      color: "rgba(0, 168, 89, 0.05)",
    },
    {
      icon: buildIcon3,
      featuresTitle: "Speed ",
      featuresDescription:
        "We cut through the noise and deliver in timely fashion.",
      color: "rgba(254, 246, 247, 1)",
    },
  ];

  const variants = {
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
  const productsVariants = {
    hide: {
      x: 40,
      y: 20,
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
    <div>
      <Navbar />
      <BannerSection
        as={motion.section}
        variants={variants}
        initial={"hide"}
        exit={"hide"}
        whileInView={"show"}
      >
        <BannerImageContainer>
          <img src={bannerLogo} />
        </BannerImageContainer>
        <BannerHeader> trovotech</BannerHeader>
        <BannerPhrase>Build Web3 solutions without limits</BannerPhrase>
        <InputSection>
          <BarnerInputContainer>
            <BarnerInput placeholder="Be the first to try our products" />
            <p>|</p>
            <BarnerButton>TELL US</BarnerButton>
          </BarnerInputContainer>
        </InputSection>
      </BannerSection>
      <FeaturesSection>
        <FeaturesContainer>
          <FeaturesHeaderSection
            as={motion.div}
            variants={featuresVariants}
            initial="hide"
            whileInView="show"
            exit="hide"
          >
            <FeaturesTitle>Why buiild with us?</FeaturesTitle>
            <FeaturesSubHeading>
              We have experience with building layer 1 & 2 distributed ledger
              technologies (DLTs)
            </FeaturesSubHeading>
            <IconContainer>
              <img src={leftChevron} alt="" />
              <img src={rightChevron} alt="" />
            </IconContainer>
          </FeaturesHeaderSection>
          {featuresContents.map((featuresContent) => {
            return (
              <FeaturesCardContainer $backgroundColor={featuresContent.color}>
                <FeaturesCardImage src={featuresContent.icon} />
                <FeaturesCardTitle>
                  {featuresContent.featuresTitle}
                </FeaturesCardTitle>
                <FeaturesCardPhrase>
                  {featuresContent.featuresDescription}
                </FeaturesCardPhrase>
              </FeaturesCardContainer>
            );
          })}
        </FeaturesContainer>
      </FeaturesSection>
      <ProductsSection>
        <ProductsContainer>
          <ProductImageContainer
            as={motion.div}
            variants={featuresVariants}
            initial="hide"
            whileInView="show"
            exit="hide"
          >
            <ProductImage src={trovoView} />
          </ProductImageContainer>
          <ProductTextContainer
            as={motion.div}
            variants={productsVariants}
            initial="hide"
            whileInView="show"
            exit="hide"
          >
            <ProductHeader>Our Products-Trovo App</ProductHeader>
            <ProductPhrase>
              A decentralized, non-custodial multipurpose application with a
              leading feature of Asset Tokenization, provides shared access
              feature with multisig, multi-sub-accounts, and a whole lot of
              other features for secure storage and exchange of digital assets.
            </ProductPhrase>
            <ServiceCardButton>Explore More</ServiceCardButton>
          </ProductTextContainer>
        </ProductsContainer>
        <ChevronContainer>
          <img src={leftChevron} alt="" />
          <img src={rightChevron} alt="" />
        </ChevronContainer>
      </ProductsSection>
      <ServicesSection>
        <ServicesContainer>
          <motion.div
            variants={serviceVariants}
            initial="hide"
            whileInView="show"
            exit="hide"
          >
            <ServicesHeader>Our Services</ServicesHeader>
            <ServicesParagraph>
              At Trovotech we tap into our combined decades of experience in the
              technology industry to create blockchain and other Web 3 solutions
              that are fit for purpose. We help clients to analyze and review
              their current systems to understand areas that require
              optimization and help them achieve their business optimization
              goals. The Trovotech team is highly competent in various aspects
              of business technology not limited to Blockchain/Distributed
              Ledger Technology (DLT) only solutions
            </ServicesParagraph>
          </motion.div>
          <ServiceCardContainer
            as={motion.div}
            variants={serviceVariants}
            initial="hide"
            whileInView="show"
            exit="hide"
          >
            {servicesCardDetails.map((serviceCardDetail) => {
              return (
                <ServiceCardWrapper
                  as={motion.div}
                  variants={serviceVariants}
                  initial="hide"
                  whileInView="show"
                  exit="hide"
                >
                  <ServiceCardImageContainer>
                    <ServiceCardImage src={serviceCardDetail.image} />
                  </ServiceCardImageContainer>
                  <ServiceCardTextSection>
                    <ServiceCardTitle>
                      {serviceCardDetail.cardTitle}
                    </ServiceCardTitle>
                    <ServiceCardDescription>
                      {serviceCardDetail.cardDescription}
                    </ServiceCardDescription>
                  </ServiceCardTextSection>
                  <ServiceCardButton>Explore More</ServiceCardButton>
                </ServiceCardWrapper>
              );
            })}
          </ServiceCardContainer>
        </ServicesContainer>
      </ServicesSection>
      <ContactUsSection>
        <ContactUsContainer>
          <ContactUsTextSection>
            <ContactUsHeading>Get started with us</ContactUsHeading>
            <ContactUsPrompt>
              Do you wish to know more about how we can add value to your
              business? Reach out to us today.
            </ContactUsPrompt>
            <ContactUsButtonContainer>
              <CallButtion>Schedule a call</CallButtion>
              <EmailButton>Email Us</EmailButton>
            </ContactUsButtonContainer>
          </ContactUsTextSection>
          <ContactUsImageSection>
            <ContactUsImageContainer>
              <OfficeImage src={trovoOffice} />
            </ContactUsImageContainer>
          </ContactUsImageSection>
        </ContactUsContainer>
      </ContactUsSection>
      <Footer />
    </div>
  );
};

const BannerSection = styled.section`
  height: calc(100vh - 70px);
  background: #fff;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  align-items: center;
  width: 90%;
  margin: auto;
  padding-bottom: 70px;
  @media (max-width: 1024px) {
    padding-bottom: 0;
  }
  @media (max-width: 768px) {
    height: calc(85vh - 70px);
  }
`;
const BannerImageContainer = styled.div``;
const BannerHeader = styled.h1`
  color: #3a4c66;
  text-align: center;
  font-family: "MatahariExtended";
  font-size: 128px;
  font-weight: 600;
  line-height: 72px;
  letter-spacing: 0.169px;
  margin: 5px 0;

  @media (max-width: 768px) {
    font-size: 76px;
    line-height: 82px;
  }
  @media (max-width: 400px) {
    font-size: 60px;
    line-height: 68px;
  }
`;
const BannerPhrase = styled.p`
  color: #1b1d21;
  text-align: center;
  font-family: "MatahariExtended";
  font-size: 40px;
  font-weight: 400;
  line-height: 72px;
  letter-spacing: 0.169px;
  margin: 0;

  @media (max-width: 768px) {
    font-size: 30px;
    line-height: 45px;
  }
  @media (max-width: 400px) {
    font-size: 24px;
    line-height: 30px;
  }
`;
const BarnerInputContainer = styled.div`
  border-radius: 8px;
  background: #fff;
  border-radius: 6px;
  display: flex;
  column-gap: 30px;

  @media (max-width: 600px) {
    column-gap: 12px;
  }
`;

const FeaturesSection = styled.section`
  background-color: rgba(255, 255, 255, 1);
  padding: 100px 0 58px;

  @media (max-width: 768px) {
    padding: 70px 0 58px;
  }
`;
const FeaturesContainer = styled.div`
  display: flex;
  column-gap: 24px;
  width: 80%;
  margin: auto;

  @media (max-width: 1024px) {
    width: 90%;
  }
  @media (max-width: 768px) {
    flex-wrap: wrap;
    row-gap: 24px;
    justify-content: center;
  }
`;
const FeaturesHeaderSection = styled.div`
  width: 506px;
  margin-right: 30px;

  @media (max-width: 1024px) {
    margin-right: 0;
  }
  @media (max-width: 1024px) {
    margin-right: 0;
  }
  @media (max-width: 768px) {
    text-align: center;
  }
`;
const FeaturesTitle = styled.h3`
  font-family: "poppins", sans-serif;
  font-weight: 600;
  font-size: 40px;
  line-height: 52px;
  margin: 0;
`;
const FeaturesSubHeading = styled.p`
  font-family: "Inter", sans-serif;
  font-weight: 400;
  font-size: 20px;
  line-height: 36px;

  @media (max-width: 1024px) {
    font-size: 17px;
    line-height: 30px;
  }
`;
const FeaturesCardContainer = styled.div((props) => ({
  backgroundColor: props.$backgroundColor,
  padding: "32px 20px",
  width: "280px",
}));
const IconContainer = styled.div`
  display: flex;
  align-items: center;
  column-gap: 16px;

  @media (max-width: 768px) {
    justify-content: center;
  }
`;
const ChevronContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  column-gap: 15px;
  margin-top: 50px;
`;
const FeaturesCardImage = styled.img``;
const FeaturesCardTitle = styled.h5`
  font-weight: 600;
  color: #1b1d21;
  font-size: 20px;
  margin: 14px 0;
  line-height: 28px;
`;
const FeaturesCardPhrase = styled.p`
  color: #000000;
  font-size: 16px;
  line-height: 28px;
`;
const ProductsSection = styled.section`
  width: 80%;
  margin: auto;
  padding: 33px 0;

  @media (max-width: 1024px) {
    width: 90%;
  }
`;
const ProductsContainer = styled.div`
  display: flex;
  column-gap: 78px;
  justify-content: center;
  align-items: center;

  @media (max-width: 1024px) {
    flex-direction: column-reverse;
    row-gap: 30px;
    text-align: center;
    justify-content: center;
  }
`;
const ProductImageContainer = styled.div`
  flex: 1;
`;
const ProductImage = styled.img`
  width: 100%;
  height: 100%;
`;
const ProductTextContainer = styled.div`
  flex: 1;
`;
const ProductHeader = styled.h3`
  font-size: 40px;
  line-height: 52px;
  font-family: "poppins";
  font-weight: 600;
  margin: 0px, 0px, 0px, 24px;

  @media (max-width: 1024px) {
    font-size: 32px;
    line-height: 38px;
    margin: 0px;
  }
`;
const ProductPhrase = styled.p`
  font-size: 20px;
  line-height: 36px;
  font-family: "poppins";
  font-weight: 400;
  margin: 0, 0, 0, 24;
`;
const ServicesSection = styled.section`
  background-color: #f6f7fa;
`;
const ServicesContainer = styled.div`
  width: 80%;
  margin: 3px auto;
  padding: 33px 0 60px;

  @media (max-width: 1024px) {
    width: 90%;
  }
`;
const ServicesHeader = styled.h3`
  color: #1b1d21;
  font-family: "Poppins", sans-serif;
  font-weight: 600;
  line-height: 52px;
  font-size: 40px;
  text-align: center;
`;
const ServicesParagraph = styled.p`
  color: #787c84;
  font-size: 20px;
  font-weight: 400;
  font-family: "Inter", sans-serif;
  line-height: 36px;
`;

const ServiceCardWrapper = styled.div`
  display: flex;
  flex-direction: column;
  width: calc(33.33% - 40px);
  justify-content: flex-start;
  align-items: center;

  @media (max-width: 768px) {
    width: calc(50.33% - 40px);
  }
  @media (max-width: 600px) {
    width: calc(100% - 40px);
  }
`;
const ServiceCardContainer = styled.div`
  display: flex;
  column-gap: 30px;
  width: 100%;
  flex-wrap: wrap;
  row-gap: 30px;
  margin-top: 30px;
  justify-content: center;
  /* flex-wrap: wrap; */
`;
const ServiceCardImageContainer = styled.div``;
const ServiceCardImage = styled.img`
  width: 100%;
  height: 100%;
`;
const ServiceCardTextSection = styled.div``;
const ServiceCardTitle = styled.h5`
  font-size: 20px;
  font-weight: 600;
  margin: 12px 0;
  text-align: center;
  font-family: "Inter", sans-serif;
  color: #000000;
`;
const ServiceCardDescription = styled.p`
  font-weight: 400;
  font-size: 16px;
  line-height: 28px;
  text-align: center;
  font-family: "Inter", sans-serif;
  color: #000000;
`;
const ServiceCardButton = styled.button`
  background-color: #007cdf;
  color: #ffffff;
  border: none;
  padding: 11px 22px;
  width: 162px;
  border-radius: 8px;
  font-size: 18px;
  font-weight: 500;
  font-family: "Inter", sans-serif;
`;
const InputSection = styled.div`
  height: 150px;
  display: flex;
  align-items: flex-end;

  @media (max-width: 1024px) {
    height: 50px;
  }
`;
const BarnerButton = styled.button`
  border-radius: 4px;
  background: #007cdf;
  border: none;
  padding: 13px 41px;
  color: #fff;
  cursor: pointer;

  @media (max-width: 600px) {
    padding: 10px 10px;
  }
`;
const BarnerInput = styled.input`
  border: none;
  width: 285px;
  outline: none;
  font-weight: 400;
  font-size: 20px;
  line-height: 36px;
  font-family: "Inter", sans-serif;
  color: #1b1d21;

  @media (max-width: 600px) {
    width: 222px;
    font-size: 15px;
    line-height: 25px;
  }
`;
// const

const ContactUsSection = styled.section`
  background-color: #fff;
  width: 80%;
  margin: auto;
  padding-top: 135px;
  padding-bottom: 100px;

  @media (max-width: 1024px) {
    width: 90%;
    padding-top: 30px;
    padding-bottom: 30px;
  }
  @media (max-width: 600px) {
    width: 100%;
    padding-top: 30px;
    padding-bottom: 30px;
  }
`;
const ContactUsContainer = styled.div`
  background-color: #1d212f;
  display: flex;
  padding: 50px 45px;
  border-radius: 10px;
  column-gap: 90px;
  justify-content: center;
  align-items: center;

  @media (max-width: 768px) {
    column-gap: 25px;
  }
  @media (max-width: 600px) {
    flex-direction: column-reverse;
    row-gap: 30px;
  }
`;
const ContactUsTextSection = styled.div``;
const ContactUsImageSection = styled.div``;
const ContactUsImageContainer = styled.div`
  /* width: 432px;
  height: 357px; */
`;
const OfficeImage = styled.img`
  width: 100%;
  height: 100%;
  object-fit: contain;
`;
const ContactUsHeading = styled.h4`
  color: #fff;
  font-family: "Poppins", sans-serif;
  font-size: 40px;
  font-weight: 600;
  line-height: 52px;
  letter-spacing: 0.20000000298023224px;
  text-align: left;
  margin: 0;

  @media (max-width: 768px) {
    font-size: 28px;
    line-height: 36px;
  }
`;
const ContactUsPrompt = styled.p`
  font-family: "Inter";
  font-size: 20px;
  font-weight: 400;
  line-height: 36px;
  letter-spacing: 0px;
  text-align: left;
  color: #fff;
  margin: 24px 0 41px 0;

  @media (max-width: 768px) {
    font-size: 16px;
    line-height: 22px;
  }
`;
const ContactUsButtonContainer = styled.div`
  display: flex;
  column-gap: 24px;
`;
const CallButtion = styled.button`
  background-color: #007cdf;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 12px 20px;
`;
const EmailButton = styled.button`
  background-color: inherit;
  color: #fff;
  border: 1px solid #fff;
  border-radius: 8px;
  padding: 12px 25px;
`;
