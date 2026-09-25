"use client";

import { useRouter } from "next/navigation";
import PrimaryButton from "@/components/PrimaryButton";
import { useState } from "react";
import styled from "styled-components";
import { FaArrowLeft } from "react-icons/fa6";

import ReusableSuccessModal from "../../components/ReusableSuccessModal";

const CreateApiKeysPage = () => {
  const router = useRouter();
  const [apiKeyName, setApiKeyName] = useState("API Key 1");
  const [apiKey, setApiKey] = useState(
    "AS6HFTRHVU9876YNHGDLJFUEVSMKEFIEUIE75SDCVBNMU70453EHNKKKNLOJHDSSO98....",
  );

  const [open, setOpen] = useState(false);
  return (
    <PageWrapper>
      <BackButton onClick={() => router.back()}>
        <FaArrowLeft />
      </BackButton>

      <Card>
        <Title>Create API Key</Title>

        <FormGroup>
          <Label>Select Api Key Type</Label>
          <Select>
            <option>Select</option>
          </Select>
        </FormGroup>

        <FormGroup>
          <Label>API Key Name</Label>
          <Input
            placeholder="E.g 1000"
            value=""
            onChange={(e) => setApiKeyName(e.target.value)}
          />
        </FormGroup>

        <PrimaryButton
          buttonStyle={{ width: "100%" }}
          onClick={() => {
            // open success modal
            setOpen(true);
          }}
        >
          Create API Key
        </PrimaryButton>
        <ReusableSuccessModal
          open={open}
          onClose={() => setOpen(false)}
          title="Successful!"
          description="You've successfully created your API key!This key is sensitive and will only be shown once. Make sure to copy it or save it somewhere safe."
          buttonText="Done"
          infoItems={[
            {
              label: "Token Code",
              value: apiKeyName,
            },
            {
              label: "Public Key",
              value: apiKey,
              copyable: true,
            },
          ]}
        />
      </Card>
    </PageWrapper>
  );
};

export default CreateApiKeysPage;

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
