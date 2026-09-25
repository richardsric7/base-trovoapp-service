"use client";
import React, { useEffect, useMemo, useRef, useState } from "react";
import Link from "next/link";
import { FaArrowLeft, FaCopy } from "react-icons/fa6";

import { LuCalendarDays, LuMail, LuPhone, LuSend } from "react-icons/lu";
import { MdVerified } from "react-icons/md";
import styled from "styled-components";
import { BsThreeDotsVertical } from "react-icons/bs";
import { SuspendUserModal } from "../components";
import TransactionsTab from "../components/TransactionsTab";
import Recovery from "../components/Recovery";
import Tab from "@/components/Tab";

import { useFetchUserByCriteriaQuery } from "@/redux/api/users";
import { useParams, useRouter } from "next/navigation";
import goldTag from "@/assets/images/goldtag.svg";
import Image from "next/image";
import copyIcon from "@/assets/images/document-copy.svg";
import statusIcon from "@/assets/images/status.svg";
import { toast } from "react-toastify";
import UserActivityTab from "../components/UserActivityTab";

import ViolationsHistory from "../components/ViolationsHistoryTab";
interface isActiveProps {
  isactive: boolean;
}

export default function ViewSingleUser() {
  const node = useRef<HTMLDivElement>(null);
  const router = useRouter();
  const [viewSuspensionModal, setViewSuspensionModal] = useState(false);
  const [showTooltip, setShowTooltip] = useState(false);
  const [currentTab, setCurrentTab] = useState("overview");
  const tabs = useMemo(
    () => [
      { key: "overview", label: "Overview" },
      { key: "recovery", label: "Recovery" },
      { key: "transactions", label: "Transactions" },
      { key: "violations", label: "Violations" },
    ],
    []
  );

  const params = useParams();
  const username = Array.isArray(params.username)
    ? params.username[0]
    : params.username;

  const { data, isLoading, isError, error } = useFetchUserByCriteriaQuery({
    username,
  });

  console.log("data", data);
  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (isError) {
    return <div>Error fetching user data</div>;
  }

  const handleCopy = (text: string) => {
    if (!text) {
      toast.error("Nothing to copy.");
      return;
    }
    navigator.clipboard
      .writeText(text)
      .then(() => toast.success("Public key copied to clipboard!"))
      .catch(() => toast.error("Failed to copy public key."));
  };

  const handleClick = (
    e: React.MouseEvent<HTMLDivElement | HTMLButtonElement, MouseEvent>
  ) => {
    e.stopPropagation();
    setViewSuspensionModal(true);

    setShowTooltip(!showTooltip);
  };
  const handleToggleTooltip = (e: React.MouseEvent) => {
    e.stopPropagation();
    setShowTooltip(!showTooltip);
  };
  const backLink = () => {
    router.back();
  };

  return (
    <Container>
      <BackLink onClick={backLink}>
        <FaArrowLeft />
      </BackLink>
      <Row>
        <div>
          <UserInfoContainer>
            <UserInfoSection>
              <Avatar> </Avatar>
              <div>
                <UserName>{data?.data?.user_info?.username}</UserName>
                <UserNameRow>
                  <FullName>
                    {data?.data?.user_info?.first_name}{" "}
                    {data?.data?.user_info?.last_name}
                  </FullName>
                </UserNameRow>
              </div>
            </UserInfoSection>
            <TooltipIconContainer>
              <BsThreeDotsVertical onClick={handleToggleTooltip} />
              {showTooltip && (
                <Tooltip
                  toolTipText="Record Voilation"
                  handleClick={handleClick}
                />
              )}
            </TooltipIconContainer>
          </UserInfoContainer>
          <Wrapper>
            <UserRow>
              <VerifiedTag>
                {data?.data?.user_info?.verified ? (
                  <>
                    <p>Verified</p>
                    <MdVerified color="#007CDF" />
                  </>
                ) : (
                  <>
                    <p>Not Verified</p>
                    {/* <MdVerified color="#007CDF" /> */}
                  </>
                )}
              </VerifiedTag>

              <ConfirmTag>
                <Image
                  src={goldTag}
                  alt="trovoplantag"
                  width={20}
                  height={20}
                />
                Gold Patron
              </ConfirmTag>
              <UserStatus isactive={data?.data?.user_info?.suspended === 0}>
                {data?.data?.user_info?.suspended === 0
                  ? "Active"
                  : "Suspended"}
              </UserStatus>
            </UserRow>

            <UserDetails>
              <div>
                <LuCalendarDays color="#828282" />

                <UserDetailsText>
                  Joined{" "}
                  {data?.data?.user_info?.created_at
                    ? new Date(
                        data.data.user_info.created_at
                      ).toLocaleDateString("en-US", {
                        year: "numeric",
                        month: "long",
                        day: "numeric",
                      })
                    : "N/A"}
                </UserDetailsText>
              </div>
              <div>
                <LuMail color="#828282" />
                <UserDetailsText>
                  {data?.data?.user_info?.email}
                </UserDetailsText>
              </div>

              <div>
                <LuPhone color="#828282" />
                <UserDetailsText>
                  {data?.data?.user_info?.mobile}
                </UserDetailsText>
              </div>
            </UserDetails>

            <PublicKeyContainer>
              <Text>Public Key</Text>
              <TagContainer>
                <Tag>{data?.data?.user_info?.address}</Tag>
                <Image
                  src={copyIcon}
                  alt="copy-icon"
                  onClick={() =>
                    handleCopy(data?.data?.user_info?.address || "N/A")
                  }
                  style={{ cursor: "pointer" }}
                />
              </TagContainer>
            </PublicKeyContainer>
          </Wrapper>
        </div>
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
        {currentTab === "overview" && <UserActivityTab />}
        {currentTab === "recovery" && <Recovery />}
        {currentTab === "transactions" && <TransactionsTab />}
        {currentTab === "violations" && <ViolationsHistory />}
      </Tab>

      {viewSuspensionModal && (
        <SuspendUserModal
          isOpen={viewSuspensionModal}
          setIsOpen={setViewSuspensionModal}
          userEmail={data?.data?.user_info?.email || "N/A"}
        />
      )}
    </Container>
  );
}

// const ToolTipContainer = styled.div``
const Tooltip = ({
  toolTipText,
  handleClick,
}: {
  toolTipText: string;
  handleClick: (
    e: React.MouseEvent<HTMLDivElement | HTMLButtonElement, MouseEvent>
  ) => void;
}) => {
  return (
    <ToolTipContainer>
      <>
        <ToolTipButton onClick={handleClick}>
          <Image src={statusIcon} alt="status-icon" />
          <ToolTipText>{toolTipText}</ToolTipText>
        </ToolTipButton>
      </>
    </ToolTipContainer>
  );
};

const Container = styled.div`
  padding: 20px 20px 0 20px;
`;
const UserInfoContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
`;
const UserInfoSection = styled.div`
  display: flex;
  column-gap: 10px;
  align-items: center;
  margin-top: 20px;
`;
const Row = styled.div`
  padding: 20px;
  background: #ffffff;
  border-radius: 24px;
`;

const UserRow = styled.div`
  display: flex;
  align-items: center;
  justify-content: start;
  margin-bottom: 10px;
  gap: 10px;
  width: 38%;
`;

const Avatar = styled.div`
  width: 48px;
  height: 48px;
  border-radius: 24px;
  background-color: rebeccapurple;
`;
const TooltipIconContainer = styled.div`
  position: relative;
  cursor: pointer;
`;
const UserNameRow = styled.div`
  display: flex;
  align-items: center;
  margin-bottom: 1px;
`;
const UserName = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 28px;

  color: #00225a;
  text-transform: capitalize;
`;
const FullName = styled.p`
  font-size: 14px;
  font-weight: 400;

  text-align: left;
  margin: 0;
  color: #00225a;
`;
const UserStatus = styled.div.attrs<{ isactive: boolean }>({
  isactive: undefined,
})<isActiveProps>`
  color: ${(props) => (props.isactive ? "#BE3800" : "#00A859")};
  font-weight: 500;
  background-color: #00a8591a;
  padding: 4px 10px;
  border-radius: 6px;
  display: inline-flex;
  text-align: center;
  font-size: 12px;
  margin: 0 auto;
`;

const ConfirmTag = styled.div`
  color: #ca8503;
  padding: 4px 10px;
  background-color: #ffc54d33;
  text-align: center;
  border-radius: 6px;
  font-size: 12px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
`;
const VerifiedTag = styled.div`
  padding: 4px 10px;
  border-radius: 6px;
  column-gap: 4px;
  background-color: #007cdf1a;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  font-size: 12px;
  font-weight: 400;
  line-height: 16px;
  margin: 0 auto;
  color: #007cdf;
`;
const UserDetails = styled.div`
  display: flex;
  align-items: center;
  column-gap: 24px;
  margin-left: 15px;

  div {
    display: flex;
    align-items: center;
    column-gap: 4px;
  }
`;
const UserDetailsText = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 17.07px;
  text-align: left;
  color: #00225a;
`;

const BackLink = styled.div`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
  padding-bottom: 10px;
`;

const ToolTipContainer = styled.div`
  position: absolute;
  top: 75%;
  right: 3px;
  width: 200px;
  background-color: #fff;
  border-radius: 16px;
  box-shadow: 0px 4px 8px 0px #0000001a;
  padding: 16px;
`;
const ToolTipButton = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
`;

const ToolTipText = styled.p`
  margin: 0;
  color: #00225a;
  font-size: 14px;
  font-weight: 400;
  line-height: 28px;
  text-align: left;
`;

const PublicKeyContainer = styled.div`
  background: #f2f6f9;
  border-radius: 12px;
  padding: 12px;
  margin-top: 10px;
  margin-left: 12px;
`;

const Text = styled.p`
  font-size: 14px;
  font-weight: 400;
  margin-bottom: 6px;
  color: #828282;
`;

const Tag = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
`;
const TagContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 20px;
`;
const Wrapper = styled.div`
  margin-left: 40px;
`;
