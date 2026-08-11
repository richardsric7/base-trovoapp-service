import React from "react";
import { styled } from "styled-components";
import userAvatar from "../utils/assets/images/user-avatar.svg";
import starRating from "../utils/assets/images/rating-star.svg";

export const ServiceSlide = ({ trovoService }) => {
  return (
    <ServiceCard>
      <ServiceImage src={trovoService.image} />
      <ServiceTitle>{trovoService.title}</ServiceTitle>
      <ServiceText>{trovoService.description}</ServiceText>
    </ServiceCard>
  );
};

export const FeedbackCard = ({ feedbackDetail }) => {
  return (
    <UserFeedback>
      <RatingContainer>
        <StarRating src={starRating} />
        <StarRating src={starRating} />
        <StarRating src={starRating} />
        <StarRating src={starRating} />
        <StarRating src={starRating} />
      </RatingContainer>
      <FeedbackText>{feedbackDetail.description}</FeedbackText>
      <UserDetails>
        <UserAvatar src={userAvatar} />
        <UserDescription>
          <UserName>{feedbackDetail.username}</UserName>
          <UserDesignation>{feedbackDetail.designation}</UserDesignation>
        </UserDescription>
      </UserDetails>
    </UserFeedback>
  );
};

const UserFeedback = styled.div`
  margin-top: 100px;
`;
const RatingContainer = styled.div``;
const StarRating = styled.img``;
const FeedbackText = styled.p`
  color: #1b1d21;
  font-size: 14px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 400;
  line-height: 28px;

  @media (max-width: 1000px) {
    width: 90%;
  }
`;
const UserDetails = styled.div`
  display: flex;
  align-items: center;
  column-gap: 16px;
`;
const UserDescription = styled.div``;
const UserAvatar = styled.img``;
const UserName = styled.h6`
  color: #1b1d21;
  font-size: 14px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 600;
  line-height: normal;
  margin: 0;
`;
const UserDesignation = styled.p`
  color: #8d8e90;
  font-size: 12px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  margin: 0;
`;
const ServiceCard = styled.div`
  padding: 24px 20px;
  border-radius: 10px;
  background: rgba(0, 73, 136, 0.03);
  max-width: 280px;
  width: calc(50.33% - 20px);
  margin: auto;

  @media (max-width: 450px) {
    width: calc(80.33% - 20px);
  }
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
