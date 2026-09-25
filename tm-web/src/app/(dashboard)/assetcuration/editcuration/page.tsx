"use client";
import React, { useState } from "react";
import styled from "styled-components";
import PrimaryButton from "@/components/PrimaryButton";
import Link from "next/link";

import { FaArrowLeft } from "react-icons/fa6";
import AddCurationTable from "../components/AddCurationTable";

import CompleteEdit from "../components/CompleteEdit";
const EditCurationpage = () => {
  const [openModal, setOpenModal] = useState<boolean>(false);
  return (
    <PageContainer>
      <BackLink href="/assetcuration">
        <FaArrowLeft />
      </BackLink>

      <Header>
        <TitleSection>
          <Title>Edit Curated Assets</Title>
          <Text>Please tick all the assets you would like to edit.</Text>
        </TitleSection>

        <PrimaryButton
          buttonStyle={{
            width: "auto",
            padding: "10px 20px",
            display: "flex",
            alignItems: "center",
            gap: "4px",
          }}
          onClick={() => setOpenModal(!openModal)}
        >
          Complete
        </PrimaryButton>
      </Header>
      <AddCurationTable />
      {openModal && (
        <CompleteEdit openModal={openModal} setOpenModal={setOpenModal} />
      )}
    </PageContainer>
  );
};

export default EditCurationpage;

const PageContainer = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  margin-bottom: 25px;
`;

const TitleSection = styled.div`
  flex-grow: 1;
`;
const Text = styled.p`
  font-size: 16px;
  font-weight: 400;
  margin: 0;
  color: #828282;
`;
const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  color: #00225a;
`;

const BackLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
  padding-bottom: 20px;
`;
