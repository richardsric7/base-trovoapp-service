import Carousel from "react-multi-carousel";
import "react-multi-carousel/lib/styles.css";
import arrowRight from "./assets/images/arrowRIght.svg";
import arrowLeft from "./assets/images/arrowLeft.svg";
import styled from "styled-components";
import { FeedbackCard, ServiceSlide } from "../trovotech-main/carouselCards";
import chevronLeft from "../utils/assets/images/chevron-left.svg";
import chevronRight from "../utils/assets/images/chevron-right.svg";
import "./slider.css"
import accountRecovery from "../utils/assets/images/accountRecovery.svg";
import digitalAssets from "../utils/assets/images/digitalAssets.svg";
import multiWallet from "../utils/assets/images/multiWallet.svg";
import nfts from "../utils/assets/images/nfts.svg";
import swapAssets from "../utils/assets/images/swapAssets.svg";
import tokenizedAssets from "../utils/assets/images/tokenizedAssets.svg";


export const ServiceSlider = () => {
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
          title: "NFTs",
          description: "NFTs minted on Bantu and other supported blockchains",
        },
      ];
  const responsive = {
    superLargeDesktop: {
      breakpoint: { max: 4000, min: 3000 },
      items: 0,
    },
    desktop: {
        breakpoint: { max: 3000, min: 1300 },
        items: 0,
    },
    tablet: {
      breakpoint: { max: 1300, min: 700 },
      items: 0,
    },
    largeMoblie: {
      breakpoint: { max: 700, min: 514 },
      items: 1,
    },
    mobile: {
      breakpoint: { max: 514, min: 0 },
      items: 1,
    },
  };

  return (
    <Carousel
    //   customLeftArrow={<Previous src={arrowLeft} alt="" />}
    //   customRightArrow={<NextButton src={arrowRight} alt="" />}
      responsive={responsive}
    >
      {trovoServices.map((trovoService, i) => (
        <ServiceSlide trovoService={trovoService} />
      ))}
    </Carousel>
  );
};
export const FeedbackSlider = () => {
  const feedbackDetails = [
    {
      description:
        "The Trovotech team have a good understanding of client needs as well as the technical know how to build the right business products att the least possible time to market.",
      username: "Nkem Ezeugo", designation: 'Director at NGX'
    },
    {
      description:
        "The Trovotech team have a good understanding of client needs as well as the technical know how to build the right business products att the least possible time to market.",
      username: "Nkem Ezeugo", designation: 'Director at NGX'
    },
    
  ];
  const responsive = {
    superLargeDesktop: {
      // the naming can be any, depends on you.
      breakpoint: { max: 4000, min: 3000 },
      items: 2,
    },
    desktop: {
      breakpoint: { max: 3000, min: 1300 },
      items: 2,
    },
    tablet: {
      breakpoint: { max: 1300, min: 840 },
      items: 2,
    },
    largeMoblie: {
      breakpoint: { max: 840, min: 514 },
      items: 1,
    },
    mobile: {
      breakpoint: { max: 514, min: 0 },
      items: 1,
    },
  };

  return (
    <Carousel
      customLeftArrow={<Previous src={chevronLeft} alt="" />}
      customRightArrow={<NextButton src={chevronRight} alt="" />}
      responsive={responsive}
      itemClass={'carouselClass'}
      partialVisbile={false}


    >
      {feedbackDetails.map((feedbackDetail, i) => (
      <div>
        <FeedbackCard feedbackDetail={feedbackDetail}  />
      </div>
       ))}
    </Carousel>
  );
};

const NextButton = styled.img`
  cursor: pointer;
  position: absolute;
  z-index: 100;
  right: 1%;
  top: 1%;
  /* margin-right: -20px; */
  @media (max-width: 1440px) {
    /* right:-1%; */
  }
  @media (max-width: 100px) {
    display: none;
  }
`;

const Previous = styled.img`
  cursor: pointer;
  position: absolute;
  right: 12%;
  top: 1%;
  /* transform: rotate(-180deg); */
  @media (max-width: 1000px) {
    display: none;
  }
`;
