"use client";

import PrimaryButton from "@/components/PrimaryButton";
import { useRouter } from "next/navigation";
import React from "react";
import { FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";
import AssetDetails from "../components/AssetDetails";
import IncomeSummary from "../components/IncomeSummary";
import CostExpenses from "../components/CostExpenses";
import Summary from "../components/Summary";

const IncomeDetailsPage = () => {
  const router = useRouter();
  const handleBack = () => {
    router.back();
  };
  return (
    <>
      <Header>
        <BackButtonLink onClick={handleBack}>
          <FaArrowLeft size={18} />
        </BackButtonLink>

        <PrimaryButton
          buttonStyle={{
            width: "100px",

            marginRight: "20px",
          }}
          onClick={() => router.push("/assetincomereport/editreport")}
        >
          Edit
        </PrimaryButton>
      </Header>
      <Container>
        {" "}
        <AssetDetails />
        <IncomeSummary />
        <CostExpenses />
        <Summary />
      </Container>
    </>
  );
};

export default IncomeDetailsPage;

const BackButtonLink = styled.button`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
  background: none;
  cursor: pointer;
  border: none;
`;
const Container = styled.section`
  padding: 32px;
  border-radius: 24px;
  background-color: #ffffff;
  min-width: 800px;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 30px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  gap: 30px;
  justify-content: space-between;
  margin-bottom: 20px;
`;
