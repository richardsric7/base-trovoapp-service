"use client";

import React, { useState } from "react";
import { FaArrowLeft } from "react-icons/fa6";
import trov from "@/assets/images/TROVTokenicon.svg";
import styled from "styled-components";
import Image from "next/image";
import PrimaryButton from "@/components/PrimaryButton";
import SuccessMessage from "@/components/SuccessMessage";
import { MdOutlineModeEditOutline } from "react-icons/md";
import { GoTrash } from "react-icons/go";
import Link from "next/link";

const EditTokenPage = () => {
  const [saveEdit, setSaveEdit] = useState(false);
  return (
    <Wrapper>
      <BackLink href="/othertoken">
        <FaArrowLeft />
      </BackLink>
      <Card>
        <Title>Edit Token</Title>
        <Subtitle>Update the token details and save your changes.</Subtitle>

        <FormGroup>
          <Label>Asset Logo</Label>
          <LogoUpload>
            <LogoPlaceholder>
              <Image src={trov} alt="TROV Token" width={70} height={70} />
            </LogoPlaceholder>
            <div>
              <Text>Trovo-logo.png</Text>
              <ButtonGroup>
                <Button>
                  <MdOutlineModeEditOutline />
                  Replace
                </Button>
                <RemoveButton>
                  {" "}
                  <GoTrash />
                  Remove
                </RemoveButton>
              </ButtonGroup>
            </div>
          </LogoUpload>
        </FormGroup>

        <FormGroup>
          <Label>Asset Code</Label>
          <Input defaultValue="TROV" disabled />
        </FormGroup>

        <FormGroup>
          <Label>Asset Issuer Public Key</Label>
          <Input defaultValue="GAXMB2AFCJGAXMB2AFCJGAXMB2AFCJ" disabled />
        </FormGroup>

        <FormGroup>
          <Label>Asset Name</Label>
          <Input defaultValue="TROV" />
        </FormGroup>

        <FormGroup>
          <Label>Description</Label>
          <TextArea defaultValue="TROV token (TROV) is the utility token that powers the Trovotech ecosystem. TROV token is used to access discounts, voting rights, airdrops, NFTs, and other community incentives." />
        </FormGroup>

        <FormGroup>
          <Label>Contact Email</Label>
          <Input defaultValue="www.trovotech.io" />
        </FormGroup>

        <FormGroup>
          <Label>Transactions</Label>
          <CheckboxWrapper>
            <CheckboxLabel>
              <Checkbox type="checkbox" defaultChecked /> Send
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
          onClick={() => setSaveEdit(!saveEdit)}
          buttonStyle={{ width: "100%" }}
        >
          Save Changes
        </PrimaryButton>
        {saveEdit && (
          <SuccessMessage
            isOpen={saveEdit}
            setIsOpen={setSaveEdit}
            heading="Success"
            message="TROV token has been successfully updated. The changes will be visible on the Trovo Wallet App"
            email=""
          />
        )}
      </Card>
    </Wrapper>
  );
};

export default EditTokenPage;

const Wrapper = styled.div`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
`;

const Card = styled.section`
  width: 500px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin: auto;
`;

const Title = styled.h2`
  color: #00225a;
  text-align: center;
  font-weight: 700;
  font-size: 24px;
  margin: 0;
  letter-spacing: 0px;
`;

const Text = styled.p`
  font-weight: 400;
  font-size: 14px;
  color: #00225a;
  padding-bottom: 6px;
`;
const Subtitle = styled.p`
  font-size: 14px;
  font-weight: 400;
  color: #828282;
  text-align: center;
  margin-bottom: 15px;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  color: #00225a;
  font-weight: 500;
  font-size: 16px;
  line-height: 24px;
  letter-spacing: 0%;
`;

const Input = styled.input`
  padding: 10px;
  border: 1px solid #bdbdbd;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 400;
  width: 100%;
  outline: none;
  color: #00225a;
  font-family: inherit;
`;

const TextArea = styled.textarea`
  padding: 10px;
  border: 1px solid #bdbdbd;
  border-radius: 12px;
  font-size: 14px;
  width: 100%;
  outline: none;
  font-weight: 400;
  color: #00225a;
  height: 80px;
  resize: none;
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

const LogoUpload = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  border-radius: 12px;
  background-color: #f2f6f9;
  border: 1px solid #007cdf;
  padding: 14px 16px;
`;

const LogoPlaceholder = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  overflow: hidden;
`;

const ButtonGroup = styled.div`
  display: flex;
  gap: 8px;
`;

const Button = styled.button`
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid #0066ff;
  background: none;
  cursor: pointer;
  font-family: inherit;
  color: #007cdf;
  display:flex;
  items-center;
  gap:2px;

  &:hover {
    background: #f0f6ff;
  }
`;

const RemoveButton = styled(Button)`
  border-color: #ff4d4f;
  color: #ff4d4f;
  &:hover {
    background: #fff0f0;
  }
`;

const BackLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
`;
