"use client";
import PrimaryButton from "@/components/PrimaryButton";
import { useRouter } from "next/navigation";
import { FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";
import { useState } from "react";
import CreateTokenSuccessModal from "@/app/external-api-clients/components/CreateTokenSuccessModal";

const CreateTokenPage = () => {
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [tokenCode, setTokenCode] = useState("");
  const [publicKey, setPublicKey] = useState("");
  return (
    <PageWrapper>
      <BackButton onClick={() => router.back()}>
        <FaArrowLeft />
      </BackButton>

      <Card>
        <Title>Create Token</Title>

        <FormGroup>
          <Label>Select Issuer Wallet</Label>
          <Select>
            <option>Select</option>
          </Select>
        </FormGroup>

        <FormGroup>
          <Label>Token Code</Label>
          <Input
            placeholder="E.g ATL"
            value={tokenCode}
            onChange={(e) => setTokenCode(e.target.value)}
          />
        </FormGroup>

        <FormGroup>
          <Label>Token Logo</Label>
          <UploadBox>
            <UploadText>
              <span>Browse file</span> to select file
            </UploadText>
            <UploadSubText>10MB max file size</UploadSubText>
          </UploadBox>
        </FormGroup>

        <FormGroup>
          <Label>Issuer Public Key</Label>
          <Input placeholder="Enter public key" />
        </FormGroup>

        <FormGroup>
          <Label>Issuer Secret Key</Label>
          <Input placeholder="Enter secret key" type="password" />
        </FormGroup>

        <PrimaryButton
          buttonStyle={{ width: "100%" }}
          onClick={() => {
            // open success modal
            setOpen(true);
          }}
        >
          Create Token
        </PrimaryButton>
        <CreateTokenSuccessModal
          open={open}
          onClose={() => setOpen(false)}
          tokenCode={tokenCode}
          publicKey={publicKey}
        />
      </Card>
    </PageWrapper>
  );
};

export default CreateTokenPage;

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

const UploadBox = styled.div`
  border: 1px dashed #d1d5db;
  border-radius: 8px;
  height: 120px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  cursor: pointer;
`;

const UploadText = styled.span`
  font-size: 13px;
  color: #6b7280;

  span {
    color: #2563eb;
    cursor: pointer;
  }
`;

const UploadSubText = styled.span`
  font-size: 12px;
  color: #9ca3af;
`;
