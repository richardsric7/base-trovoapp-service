"use client";
import React from "react";
import styled from "styled-components";
import { useRouter } from "next/navigation";

export default function Login() {
  const router = useRouter();
  return (
    <Container>
      <ContentContainer>
        <Title>Hi There, Welcome!</Title>
        <FormGroup>
          <InputLabel>Email</InputLabel>
          <InputField placeholder="Input email address" />
        </FormGroup>
        <FormGroup>
          <InputLabel>Password</InputLabel>
          <InputField placeholder="Input password" />
        </FormGroup>
        <InfoText>Forgot password?</InfoText>
        <LoginButton onClick={() => router.push("/dashboard")}>
          Login
        </LoginButton>
      </ContentContainer>
    </Container>
  );
}
const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  margin: auto;
`;
const ContentContainer = styled.div``;
const FormGroup = styled.div`
  margin: 12px auto;
`;
const Title = styled.h2`
  color: #191919;
  font-weight: 600;
  font-size: 32px;
  line-height: 39px;
  padding-bottom: 10px;
`;
const InputLabel = styled.label`
  display: block;
  color: #00225a;
  margin-bottom: 12px;
  font-weight: 500;
`;
const InputField = styled.input`
  border: 1px solid #bdbdbd;
  width: 440px;
  padding: 13px 16px;
  border-radius: 12px;
  font-size: 18px;
  line-height: 21px;
  font-family: inherit;
`;
const InfoText = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 17.07px;
  color: #007cdf;
  text-align: end;
`;
const LoginButton = styled.button`
  display: block;
  background-color: #007cdf;
  border: none;
  padding: 11px;
  color: #fff;
  margin-top: 40px;
  width: 100%;
  border-radius: 10px;
  font-family: inherit;
`;
