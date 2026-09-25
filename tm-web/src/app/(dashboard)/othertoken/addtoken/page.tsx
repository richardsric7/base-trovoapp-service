"use client";

import PrimaryButton from "@/components/PrimaryButton";
import Link from "next/link";
import React, { useState } from "react";
import { FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";
import AddTokenModal from "../components/AddTokenModal";

const AddTokenPage = () => {
  const [addToken, setAddToken] = useState(false);
  return (
    <Container>
      <BackLink href="/othertoken">
        <FaArrowLeft />
      </BackLink>
      <Card>
        <Title>Create Token</Title>
        <Subtitle>
          Please provide the following information to create and curate the
          token on the Trovo app
        </Subtitle>

        <FormGroup>
          <Label>Asset Code</Label>
          <Input placeholder="Enter asset code" />
        </FormGroup>

        <FormGroup>
          <Label>Asset Issuer Public Key</Label>
          <Input placeholder="Enter asset issuer public key" />
        </FormGroup>

        <FormGroup>
          <Label>Transactions</Label>
          <CheckboxWrapper>
            <CheckboxLabel>
              <Checkbox type="checkbox" /> Send
            </CheckboxLabel>
            <CheckboxLabel>
              <Checkbox type="checkbox" defaultChecked /> Receive
            </CheckboxLabel>
            <CheckboxLabel>
              <Checkbox type="checkbox" defaultChecked /> Deposit
            </CheckboxLabel>
            <CheckboxLabel>
              <Checkbox type="checkbox" defaultChecked /> Withdraw
            </CheckboxLabel>
          </CheckboxWrapper>
        </FormGroup>

        <FormGroup>
          <Label>Token Type</Label>
          <RadioWrapper>
            <RadioLabel>
              <Radio type="radio" name="tokenType" /> Stable Coin
            </RadioLabel>
            <RadioLabel>
              <Radio type="radio" name="tokenType" defaultChecked /> Other Token
            </RadioLabel>
          </RadioWrapper>
        </FormGroup>

        <PrimaryButton
          onClick={() => setAddToken(!addToken)}
          buttonStyle={{ width: "100%" }}
        >
          {" "}
          Continue
        </PrimaryButton>
      </Card>
      {addToken && (
        <AddTokenModal isOpen={addToken} onClose={() => setAddToken(false)} />
      )}
    </Container>
  );
};

export default AddTokenPage;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
`;

const Card = styled.section`
  width: 480px;
  display: flex;
  margin: auto;
  flex-direction: column;
  gap: 16px;
`;

const Title = styled.h2`
  color: #00225a;
  text-align: center;
  font-weight: 700;
  font-size: 24px;
  line-height: 28px;
  letter-spacing: 0px;
`;

const Subtitle = styled.p`
  font-size: 14px;
  font-weight: 400;
  color: #828282;
  text-align: center;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
`;

const Input = styled.input`
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 8px;
  font-size: 14px;
  width: 100%;
  outline: none;
  font-family: inherit;
`;

const CheckboxWrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const CheckboxLabel = styled.label`
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
`;

const Checkbox = styled.input`
  width: 16px;
  height: 16px;
`;

const RadioWrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const RadioLabel = styled.label`
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
`;

const Radio = styled.input`
  width: 16px;
  height: 16px;
`;

const BackLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
`;
