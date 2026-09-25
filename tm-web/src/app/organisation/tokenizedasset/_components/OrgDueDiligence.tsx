"use client";

import React, { useState } from "react";
import styled from "styled-components";
import VerificationSection from "./VerificationSection";
import DueDiligenceScore from "./DueDiligenceScoreCard";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import { BsThreeDots } from "react-icons/bs";
import { AiOutlineClockCircle } from "react-icons/ai";
import { Badge, Dropdown, MenuProps } from "antd";

const DueDiligenceTab = () => {
  const [openMenu, setOpenMenu] = useState(false);

  const menuItems: MenuProps["items"] = [
    {
      key: "1",
      label: "Add Notes",
    },
    {
      key: "2",
      label: "Upload Document",
    },
    {
      key: "3",
      label: "Download Documents",
    },

    {
      key: "4",
      label: "Request Document",
    },

    {
      key: "5",
      label: "Generate a Report",
    },

    {
      key: "6",
      label: "Schedule a Meeting",
    },

    {
      key: "7",
      label: <span style={{ color: "#BE3800" }}>Reject Asset</span>,
    },
  ];
  return (
    <Container>
      <HeaderRow>
        <div>
          <Title>Due Diligence</Title>
          <Status>
            <AiOutlineClockCircle />
            Pending
          </Status>
        </div>

        <ButtonContainer>
          <PrimaryButton buttonStyle={{ width: "154px" }}>
            Approve Asset
          </PrimaryButton>
          <DropdownStyles>
            <Dropdown
              menu={{ items: menuItems }}
              trigger={["click"]}
              placement="bottomRight"
              overlayClassName="due-diligence-dropdown"
            >
              <div>
                <DotBadge dot>
                  <SecondaryButton buttonStyle={{ width: "48px" }}>
                    <BsThreeDots />
                  </SecondaryButton>
                </DotBadge>
              </div>
            </Dropdown>
          </DropdownStyles>
        </ButtonContainer>
      </HeaderRow>

      <DueDiligenceScore score={7} />

      <VerificationSection
        title="Legal Verification"
        button="Verified"
        verifiedDate="20 Feb, 2026"
        verifiedBy=" Anna Quinn"
      />

      <VerificationSection
        title="Financial Verification"
        button="Mark as Verified"
      />

      <VerificationSection
        title="Operational Verification"
        button="Mark as Verified"
      />

      <VerificationSection
        title="Technical Due Diligence"
        button="Verified"
        verifiedDate="20 Feb, 2026"
        verifiedBy="Anna Quinn"
      />

      <VerificationSection
        title="Risk Assessment"
        button="Verified"
        verifiedDate="20 Feb, 2026"
        verifiedBy="Anna Quinn"
      />
    </Container>
  );
};

export default DueDiligenceTab;

const Container = styled.div`
  padding: 24px;
  background-color: #fff;
  border-radius: 24px;
  background-color: #fff;
`;

const HeaderRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
`;
const ButtonContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
`;

const Status = styled.div`
  background-color: #f2f6f9;
  color: #007cdf;
  font-size: 12px;
  font-weight: 500;
  padding: 6px 10px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
`;

const DotBadge = styled(Badge)`
  .ant-badge-dot {
    top: 20px;
    right: 2px;
    width: 8px;
    height: 8px;
  }
`;

const DropdownStyles = styled.div`
  .due-diligence-dropdown .ant-dropdown-menu-item {
    font-family: "Montserrat", sans-serif !important;
    font-weight: 400 !important;
    font-size: 14px !important;
    line-height: 28px !important;
    letter-spacing: 0% !important;
    color: #00225a !important;
  }

  .due-diligence-dropdown .ant-dropdown-menu-item:hover {
    background: #f2f6f9;
    color: #00225a;
  }
`;
