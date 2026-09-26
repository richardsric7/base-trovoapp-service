"use client";

import React, { useState, useMemo } from "react";
import { useParams, useRouter } from "next/navigation";
import { FaArrowLeft } from "react-icons/fa6";
import { LuCalendarDays, LuMail, LuPhone } from "react-icons/lu";
import { MdVerified } from "react-icons/md";
import { BsThreeDotsVertical } from "react-icons/bs";
import Image from "next/image";
import styled from "styled-components";
import goldTag from "@/assets/images/goldtag.svg";
import tradeIcon from "@/assets/images/coin.svg";
import ToggleSwitch from "@/components/ToggleSwitch";
import statusIcon from "@/assets/images/status.svg";
import Tab from "@/components/Tab";
import P2pOverview from "../components/P2pOverview";
import P2pTransactionsTab from "../components/P2pTransactionsTab";
import { useP2pUsersQuery } from "@/redux/api/p2p";
import RecordViolation from "../components/RecordViolation";
import SuccessMessage from "@/components/SuccessMessage";
interface isActiveProps {
  isactive: boolean;
}

const P2pUserDetailsPage = () => {
  const router = useRouter();
  const params = useParams();

  const rawUsername = useParams().username;
  const username = decodeURIComponent(
    Array.isArray(rawUsername) ? rawUsername[0] : rawUsername
  );

  const { data, isLoading } = useP2pUsersQuery({
    username,
  });

  const user = data?.data?.data?.[0];

  const [showTooltip, setShowTooltip] = useState(false);
  const [enabled, setEnabled] = useState(false);
  const [suspendSuccess, setSuspendSuccess] = useState(false);

  const [currentTab, setCurrentTab] = useState("overview");
  const tabs = useMemo(
    () => [
      { key: "overview", label: "Overview" },
      { key: "transactions", label: "Transactions" },
    ],
    []
  );

  const [isOpen, setIsOpen] = useState(false);
  const handleClick = () => {
    setIsOpen(!isOpen);

    setShowTooltip(false);
  };

  const handleToggleTooltip = (e: React.MouseEvent) => {
    e.stopPropagation();
    setShowTooltip(!showTooltip);
  };

  if (isLoading) return <p>Loading....</p>;

  return (
    <Container>
      <BackLink onClick={() => router.back()}>
        <FaArrowLeft />
      </BackLink>

      <Row>
        <UserInfoContainer>
          <UserInfoSection>
            <Avatar />
            <div>
              <UserName>{user?.username}</UserName>
              {/* <FullName>{`${dummyUser.first_name} ${dummyUser.last_name}`}</FullName> */}
            </div>
          </UserInfoSection>

          <TooltipIconContainer>
            <BsThreeDotsVertical onClick={handleToggleTooltip} />
            {showTooltip && (
              <TooltipWrapper onClick={(e) => e.stopPropagation()}>
                <ToolTipButton>
                  <Image src={tradeIcon} alt="status-icon" />
                  <TooltipItem>Set Trading Fee</TooltipItem>
                </ToolTipButton>
                <ToolTipButton onClick={handleClick}>
                  <Image src={statusIcon} alt="status-icon" />
                  <TooltipItem>Record Violation</TooltipItem>
                </ToolTipButton>
              </TooltipWrapper>
            )}
          </TooltipIconContainer>
        </UserInfoContainer>

        <Wrapper>
          <UserRow>
            <UserStatus isactive={!user?.suspended}>
              {user?.suspended ? "Suspended" : "Active"}
            </UserStatus>

            <ConfirmTag>
              <Image src={goldTag} alt="gold-tag" width={20} height={20} />
              Gold Patron
            </ConfirmTag>

            <VerifiedTag>
              <p>{user?.isMerchant ? "Merchant" : "Non-Merchant"}</p>
              {user?.isMerchant && <MdVerified color="#007CDF" />}
            </VerifiedTag>
          </UserRow>

          <UserDetails>
            <DetalisRow>
              <LuCalendarDays color="#828282" />

              {user?.registrationDate && (
                <UserDetailsText>
                  Joined{" "}
                  {new Date(user?.registrationDate).toLocaleDateString(
                    "en-US",
                    {
                      year: "numeric",
                      month: "long",
                      day: "numeric",
                    }
                  )}
                </UserDetailsText>
              )}
            </DetalisRow>
            <DetalisRow>
              <LuMail color="#828282" />
              <UserDetailsText>{user?.email}</UserDetailsText>
            </DetalisRow>
            <DetalisRow>
              <LuPhone color="#828282" />
              <UserDetailsText>{user?.phone}</UserDetailsText>
            </DetalisRow>
          </UserDetails>
          <EnableMerchant>
            <MerchantText>Merchant Status</MerchantText>
            <MerchantRow>
              <ToggleSwitch label="" checked={enabled} onChange={setEnabled} />

              <Stat> {enabled ? "Enabled" : "Disabled"}</Stat>
            </MerchantRow>
          </EnableMerchant>
        </Wrapper>
      </Row>

      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabButtonStyle={{
          borderRadius: "12px",
          fontSize: "16px",
          color: "#4F4F4F",
          width: "90%",
          padding: "0px 1px",
          marginTop: "10px",
        }}
        activeTabButtonStyle={{
          backgroundColor: "#F2F6F9",
          color: "#00225A",
          border: "0px",
          width: "60%",
        }}
        tabContainerStyle={{
          width: "60%",
          borderRadius: "10px",
          height: "36px",
          padding: "1px",
          margin: "20px 0 10px 10px",
        }}
        tabContentStyle={{
          padding: "10px 10px 0 10px",
          background: "transparent",
        }}
      >
        {currentTab === "overview" && <P2pOverview />}

        {currentTab === "transactions" && <P2pTransactionsTab />}
      </Tab>

      {isOpen && (
        <RecordViolation
          isOpen={isOpen}
          setIsOpen={setIsOpen}
          setSuspendSuccess={setSuspendSuccess}
        />
      )}

      {suspendSuccess && (
        <SuccessMessage
          isOpen={suspendSuccess}
          setIsOpen={setSuspendSuccess}
          heading="Success"
          message="User has been successfully recorded a violation"
          email={user?.email}
        />
      )}
    </Container>
  );
};

export default P2pUserDetailsPage;
const Container = styled.div`
  padding: 20px;
`;

const BackLink = styled.div`
  color: #000;
  cursor: pointer;
  margin-bottom: 10px;
`;

const Row = styled.div`
  padding: 20px;
  background: #fff;
  border-radius: 24px;
`;

const Avatar = styled.div`
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const UserInfoContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const UserInfoSection = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

const UserName = styled.h1`
  font-size: 20px;
  color: #00225a;
  text-transform: capitalize;
  margin: 0;
`;

const DetalisRow = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const TooltipIconContainer = styled.div`
  position: relative;
  cursor: pointer;
`;

const TooltipWrapper = styled.div`
  position: absolute;
  top: 24px;
  right: 0;
  background: #fff;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  padding: 8px 0;
  z-index: 10;
  display: flex;

  flex-direction: column;
  width: 160px;
`;

const TooltipItem = styled.div`
  font-size: 14px;
  color: #00225a;
  cursor: pointer;
`;
const ToolTipButton = styled.div`
  padding: 8px 6px;
  display: flex;
  align-items: center;
  gap: 1px;
  cursor: pointer;
  &:hover {
    background-color: #f2f6f9;
  }
`;

const UserRow = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 16px 0;
`;

const VerifiedTag = styled.div`
  background-color: #007cdf1a;
  color: #007cdf;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
`;

const ConfirmTag = styled.div`
  background-color: #ffc54d33;
  color: #ca8503;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
`;

const UserStatus = styled.div<isActiveProps>`
  background-color: ${({ isactive }) => (isactive ? "#00a8591a" : "#be38001a")};
  color: ${({ isactive }) => (isactive ? "#00A859" : "#BE3800")};
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
`;

const UserDetails = styled.div`
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
`;

const UserDetailsText = styled.p`
  font-size: 14px;
  color: #00225a;
`;

const Wrapper = styled.div`
  margin-top: 20px;
`;

const EnableMerchant = styled.div`
  background-color: #f2f6f9;
  gap: 8px;
  border-radius: 16px;
  padding-top: 16px;
  padding-right: 24px;
  padding-bottom: 16px;
  padding-left: 24px;

  width: 206px;
  margin: 6px 0;
`;

const MerchantRow = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 4px;
`;

const MerchantText = styled.p`
  font-weight: 400;
  font-size: 12px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #828282;
`;

const Stat = styled.p`
  font-weight: 400;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #00225a;
`;
