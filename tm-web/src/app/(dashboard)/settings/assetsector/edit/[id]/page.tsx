"use client";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import React, { useState } from "react";
import styled from "styled-components";

import { useParams, useRouter } from "next/navigation";
import { FaArrowLeft } from "react-icons/fa6";
import { BiEdit } from "react-icons/bi";
import { RiDeleteBin6Line } from "react-icons/ri";
import AddSubSectorModal from "../../../components/AddSubSectorModal";
import SuccessMessage from "@/components/SuccessMessage";

const EditSectorPage = () => {
  const [isOpenSubSector, setIsOpenSubSector] = useState(false);
  const [isSuccessOpen, setIsSuccessOpen] = useState(false);
  const [sectorName, setSectorName] = useState("Real Estate");

  const { id } = useParams();
  const router = useRouter();
  const sectors = ["Industrial buildings", "Rsidentail buildings"];

  const handleSave = () => {
    // TODO: call API here
    setIsSuccessOpen(true);
  };

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
            <InputField name="" placeholder="Real Estate" value="" />
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
          {sectors.map((sec, idx) => (
            <SectorContent>
              <SectorTitle>{sec}</SectorTitle>
              <SectorIconContainer>
                <EditLink>
                  <BiEdit color=" #007cdf" size={20} />
                </EditLink>

                <DeleteBtn>
                  {" "}
                  <RiDeleteBin6Line color="#be3800" size={20} />
                </DeleteBtn>
              </SectorIconContainer>
            </SectorContent>
          ))}
        </SectorContainer>
        <PrimaryButton onClick={handleSave}>Save Asset Sector</PrimaryButton>

        {isOpenSubSector && (
          <AddSubSectorModal
            add={isOpenSubSector}
            setAdd={setIsOpenSubSector}
          />
        )}

        {isSuccessOpen && (
          <SuccessMessage
            isOpen={isSuccessOpen}
            setIsOpen={(val) => {
              setIsSuccessOpen(val);
              if (!val) {
                router.push("/settings/assetsector");
              }
            }}
            heading="Success!"
            message={`You have successfully updated ${sectorName}.`}
          />
        )}
      </PageContainer>
    </>
  );
};

export default EditSectorPage;

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
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: 8px;
  width: 100%;
`;
const SectorContent = styled.div`
  background-color: #f2f6f9;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  border-radius: 8px;
  width: 100%;
`;

const SectorTitle = styled.p`
  font-weight: 400;
  font-size: 14px;
  color: #00225a;
  line-height: 24px;
  letter-spacing: 0.1px;
`;

const SectorIconContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

const EditLink = styled.div`
  border: 1px solid #007cdf;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
`;

const DeleteBtn = styled.div`
  border: 1px solid #be3800;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
`;
