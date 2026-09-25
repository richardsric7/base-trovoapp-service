"use client";

import React, { useEffect, useState } from "react";
import styled from "styled-components";

const InviteExpiredPage = () => {
  const [orgName, setOrgName] = useState("the organization");

  useEffect(() => {
    const name = sessionStorage.getItem("inviteOrgName");
    if (name && name.trim()) setOrgName(name);
  }, []);

  return (
    <Container>
      <ContentContainer>
        <Title>Invite Expired</Title>

        <InfoText>
          This invitation is no longer valid. To join
          <StyledSpan> {orgName}</StyledSpan> on Trovo Manager, ask the
          Organization&apos;s Admin to resend the invite or contact our{" "}
          <ResendButton>Support</ResendButton>.
        </InfoText>
      </ContentContainer>
    </Container>
  );
};

export default InviteExpiredPage;
const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  margin: auto;
`;
const ContentContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

const Title = styled.h2`
  color: #191919;
  font-weight: 600;
  font-size: 32px;
  line-height: 39px;
`;

const InfoText = styled.p`
  font-size: 16px;
  font-weight: 400;
  line-height: 160%;
  color: #191919;
  width: 89%;
  flex-wrap: non-wrap;
`;

const StyledSpan = styled.span`
  font-size: 16px;
  font-weight: 700;
  line-height: 17.07px;
  color: #191919;
`;

const ResendButton = styled.button`
  line-height: 21px;
  font-family: inherit;
  font-size: 14px;
  background: none;
  border: 0;
  color: #007cdf;
  cursor: pointer;
`;
