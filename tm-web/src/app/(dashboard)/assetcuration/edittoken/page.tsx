"use client";

import PrimaryButton from "@/components/PrimaryButton";
import SuccessMessage from "@/components/SuccessMessage";
import Link from "next/link";
import React, { useState } from "react";
import { FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";

const EditTokenPage = () => {
  const [openModal, setOpenModal] = useState(false);
  return (
    <Container>
      <StyledLink href="/assetcuration">
        <FaArrowLeft />
      </StyledLink>
      <Heading>Edit Token</Heading>
      <Text>Update the token details and save your changes.</Text>

      <Form>
        <div>
          <Label>Asset Code</Label>
          <Input placeholder="AE1" />
        </div>

        <div>
          <Label>Asset Issuer Public Key</Label>
          <Input placeholder="GAXMBA2CFJGAXMBA2CFJGAXMBA2CFJGAXMBA2CFJ" />
        </div>

        <div>
          <Label>Website URL</Label>
          <Input placeholder="www.atlantisestates.com" />
        </div>
        <PrimaryButton
          buttonStyle={{ width: "49%" }}
          onClick={() => setOpenModal(!openModal)}
        >
          {" "}
          Continue
        </PrimaryButton>
      </Form>
      {openModal && (
        <SuccessMessage
          isOpen={openModal}
          setIsOpen={setOpenModal}
          message="AE1 token has been successfully updated. The changes will be visible on the Trovo Wallet App"
          email="AE1"
        />
      )}
    </Container>
  );
};

export default EditTokenPage;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;
  border-radius: 24px;
  gap: 20px;
  //   display: flex;
  //   flex-direction: column;
  //   gap: 6px;
  //   align-items: center;
`;
const Heading = styled.h1`
  font-weight: 700;
  font-size: 24px;
  line-height: 28px;
  letter-spacing: 0px;
  color: #00225a;
  text-align: center;
`;

const Text = styled.p`
  font-weight: 400;
  font-size: 16px;
  line-height: 28px;
  letter-spacing: 0px;
  text-align: center;
  color: #828282;
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: center;
  padding: 10px 0;
`;

const Label = styled.h2`
  font-weight: 500;
  font-size: 16px;
  line-height: 24px;
  letter-spacing: 0%;
  color: #00225a;
`;

const Input = styled.input`
  width: 500px;
  height: 48px;
  gap: 10px;
  border-radius: 12px;
  border: 1px solid #bdbdbd;
  font-family: inherit;
  padding: 0 4px;

  placeholder {
    color: #8b90a0;
    font-family: Montserrat;
    font-weight: 400;
    font-size: 18px;
    line-height: 21.94px;
    letter-spacing: 0%;
  }
`;
const StyledLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
`;
