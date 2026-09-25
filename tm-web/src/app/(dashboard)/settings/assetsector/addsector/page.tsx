"use client";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import React, { useState } from "react";
import styled from "styled-components";
import AddSubSectorModal from "../../components/AddSubSectorModal";
import { useRouter } from "next/navigation";
import { FaArrowLeft } from "react-icons/fa6";
import Image from "next/image";
import folderImg from "@/assets/images/Folder.svg";
const AddSectorPage = () => {
  const [isOpenSubSector, setIsOpenSubSector] = useState(false);
  const router = useRouter();

  const handleBack = () => {
    router.back();
  };
  return (
    <>
      <BackButtonLink onClick={handleBack}>
        <FaArrowLeft size={18} />
      </BackButtonLink>
      <PageContainer>
        <Title>Add Asset Sector</Title>

        <FormGroup>
          <Label>Asset Sector</Label>
          <FeeInputWrapper>
            <InputField name="" placeholder="Enter Sector Name" value="" />
          </FeeInputWrapper>
        </FormGroup>

        <SubTitle> Sub-sectors</SubTitle>

        <SecondaryButton
          buttonStyle={{ width: "100%", margin: "0" }}
          onClick={() => setIsOpenSubSector(!isOpenSubSector)}
        >
          + Add Sub-sector
        </SecondaryButton>

        <SectorContainer>
          <Image src={folderImg} alt="" width={50} height={50} />
          <SectorText>No Sub-sector added yet</SectorText>
        </SectorContainer>
        <PrimaryButton buttonStyle={{ width: "100%" }}>
          Save Asset Sector
        </PrimaryButton>

        {isOpenSubSector && (
          <AddSubSectorModal
            add={isOpenSubSector}
            setAdd={setIsOpenSubSector}
          />
        )}
      </PageContainer>
    </>
  );
};

export default AddSectorPage;

const BackButtonLink = styled.button`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
  background: none;
  cursor: pointer;
  border: none;
  padding: 24px;
`;

const PageContainer = styled.section`
  width: 480px;
  display: flex;
  margin: auto;
  flex-direction: column;
  gap: 16px;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  color: #00225a;
`;

const Label = styled.label`
  color: #828282;
  font-weight: 500;
  font-size: 16px;
  margin-top: 14px;
`;
const SubTitle = styled.h1`
  font-size: 14px;
  font-weight: 600;

  color: #00225a;
`;

const InputField = styled.input`
  font-size: 14px;
  width: 100%;
  outline: none;
  font-family: inherit;
  border: none;

  &::placeholder {
    color: #bdbdbd;
    font-family: inherit;
    font-size: 14px;
  }
`;
const FeeInputWrapper = styled.div`
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  display: flex;
  align-items: center;
`;
const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const SectorContainer = styled.div`
  background-color: #f2f6f9;
  display: flex;
  align-items: center;
  flex-direction: column;
  padding: 8px 16px;
  border-radius: 8px;
  width: 100%;
`;

const SectorText = styled.p`
  font-weight: 500;
  font-style: Medium;
  font-size: 14px;
  line-height: 24px;
  color: #00225a;
`;
