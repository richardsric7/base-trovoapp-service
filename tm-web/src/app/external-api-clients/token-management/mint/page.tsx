"use client";

import { useRouter } from "next/navigation";
import PrimaryButton from "@/components/PrimaryButton";
import { useState } from "react";
import styled from "styled-components";
import { FaArrowLeft } from "react-icons/fa6";
import CreateTokenSuccessModal from "../../components/CreateTokenSuccessModal";
import MintingSuccessfulModal from "../../components/MintingSuccessfulModal";

const MintTokenPage = () => {
  const router = useRouter();
  const [mintingQuantity, setMintingQuantity] = useState("");
  const [open, setOpen] = useState(false);
  return (
    <PageWrapper>
      <BackButton onClick={() => router.back()}>
        <FaArrowLeft />
      </BackButton>

      <Card>
        <Title>Mint Token</Title>

        <FormGroup>
          <Label>Select Issuer Wallet</Label>
          <Select>
            <option>Select</option>
          </Select>
        </FormGroup>

        <FormGroup>
          <Label>Minting Quantity</Label>
          <Input
            placeholder="E.g 1000"
            value=""
            onChange={(e) => setMintingQuantity(e.target.value)}
          />
        </FormGroup>

        <PrimaryButton
          buttonStyle={{ width: "100%" }}
          onClick={() => {
            // open success modal
            setOpen(true);
          }}
        >
          Mint Token
        </PrimaryButton>
        <MintingSuccessfulModal open={open} onClose={() => setOpen(false)} />
      </Card>
    </PageWrapper>
  );
};

export default MintTokenPage;

const PageWrapper = styled.div`
  background: #fff;
  min-height: 100vh;

  border-radius: 12px;

  border-radius: 24px;
  padding: 32px;
`;

const Card = styled.div`
  width: 720px;
  margin: 0px auto;
  padding: 40px 48px;
`;

const BackButton = styled.div`
  cursor: pointer;
  font-size: 20px;
`;

const Title = styled.h2`
  text-align: left;
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 32px;
  color: #00225a;
`;

const FormGroup = styled.div`
  margin-bottom: 20px;
`;

const Label = styled.label`
  font-size: 14px;
  color: #6b7280;
  display: block;
  margin-bottom: 6px;
`;

const Input = styled.input`
  width: 100%;
  height: 44px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 0 12px;
  font-size: 14px;
  outline: none;

  &::placeholder {
    color: #9ca3af;
  }

  &:focus {
    border-color: #2563eb;
  }
`;

const Select = styled.select`
  width: 100%;
  height: 44px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 0 12px;
  font-size: 14px;
  color: #6b7280;
`;
