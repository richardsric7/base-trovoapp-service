"use client";

import React, { useState } from "react";
import styled from "styled-components";
import SecondaryButton from "@/components/SecondaryButton";
import { BsThreeDots } from "react-icons/bs";
import PrimaryButton from "@/components/PrimaryButton";
import { Badge, Dropdown, MenuProps } from "antd";
import StepperItems from "../../_components/StepperItems";
import { AuthorizeRequestModal } from "../../_components/AuthorizeRequestModal";
import { RejectRequestModal } from "../../_components/RejectRequestModal";
import PayoutList from "../../_components/PayoutList";
import IncomeRequestInfo from "../../_components/IncomeRequestInfo";
import { useRouter } from "next/navigation";
import { FaArrowLeft } from "react-icons/fa6";

const IncomeDistribution = () => {
  const [isAuthorizeOpen, setIsAuthorizeOpen] = useState(false);
  const [isRejectOpen, setIsRejectOpen] = useState(false);

  const router = useRouter();

  const steps = [
    { action: "Submitted", date: "27 Jan 2026, 02:00 PM" },
    { action: "Request Authorized", date: "" },
    { action: "Funds Released", date: "" },
    { action: "Processing", date: "" },
    { action: "Completed", date: "" },
    { action: "Confirmed", date: "" },
  ];

  const menuItems: MenuProps["items"] = [
    {
      key: "1",
      label: "Download Payout List",
    },

    {
      key: "2",
      label: <span style={{ color: "#BE3800" }}>Reject Request</span>,
      onClick: () => setIsRejectOpen(true),
    },
  ];
  return (
    <>
      <BackButton onClick={() => router.back()}>
        {" "}
        <FaArrowLeft color="#002255a" size={22} />
      </BackButton>
      <Container>
        <Sidebar>
          {steps.map((step, index) => (
            <StepperItems
              key={index}
              step={step}
              index={index}
              currentStep={0}
              stepsLength={steps.length}
            />
          ))}
        </Sidebar>
        <Content>
          <Section>
            <HeaderRow>
              <div>
                <Title>Fund Release Request</Title>
              </div>

              <ButtonContainer>
                <PrimaryButton
                  buttonStyle={{ width: "170px" }}
                  onClick={() => setIsAuthorizeOpen(true)}
                >
                  Authorize Request
                </PrimaryButton>
                <DropdownStyles>
                  <Dropdown
                    menu={{ items: menuItems }}
                    trigger={["click"]}
                    placement="bottomRight"
                    overlayClassName="due-diligence-dropdown"
                  >
                    <div>
                      <DotBadge>
                        <SecondaryButton buttonStyle={{ width: "48px" }}>
                          <BsThreeDots />
                        </SecondaryButton>
                      </DotBadge>
                    </div>
                  </Dropdown>
                </DropdownStyles>
              </ButtonContainer>
            </HeaderRow>

            <IncomeRequestInfo />
          </Section>

          {/* <Section> */}
          <PayoutList />
          {/* </Section> */}
        </Content>

        <AuthorizeRequestModal
          open={isAuthorizeOpen}
          onClose={() => setIsAuthorizeOpen(false)}
        />

        <RejectRequestModal
          open={isRejectOpen}
          onClose={() => setIsRejectOpen(false)}
        />
      </Container>
    </>
  );
};

export default IncomeDistribution;

const BackButton = styled.button`
  border: 0px;
  padding: 4px 20px;
  background-color: transparent;
  cursor: pointer;
`;

const Container = styled.div`
  display: flex;
  gap: 20px;
  padding: 20px;
  width: 100%;
  min-width: 60%;
  min-height: 100vh;
`;
const Content = styled.div`
  flex: 1;
  background: #fff;
  padding: 32px;
  border-radius: 24px;
`;
const Title = styled.h3`
  font-weight: 600;
  font-style: SemiBold;
  font-size: 24px;
  color: #00225a;
`;

const Sidebar = styled.div`
  width: 220px;
  background: white;
  flex-shrink: 0;
  padding: 24px;
  height: fit-content;
  border-radius: 16px;
`;

const Step = styled.div<{ active?: boolean }>`
  padding: 10px;
  margin-bottom: 10px;
  border-left: 3px solid ${(p) => (p.active ? "#2563eb" : "#ddd")};
  color: ${(p) => (p.active ? "#2563eb" : "#999")};
`;
const Header = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
`;

const Section = styled.div`
  margin-top: 0px;
  background-color: #fff;
  padding: 24px;
  border-radius: 16px;
`;

const ButtonContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
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

const HeaderRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
`;
const DotBadge = styled(Badge)`
  .ant-badge-dot {
    top: 20px;
    right: 2px;
    width: 8px;
    height: 8px;
  }
`;
