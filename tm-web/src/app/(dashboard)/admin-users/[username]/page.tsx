"use client";
import Image from "next/image";
import Link from "next/link";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import { FaAngleRight, FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";
import mailIcon from "@/assets/images/mail.svg";
import calenderIcon from "@/assets/images/solar_calendar-linear.svg";
import userIcon from "@/assets/images/oui_app-users-roles.svg";

import { useState } from "react";
import Tab from "@/components/Tab";
import AdminActionsTab from "../components/AdminActionsTab";
import ActivityTab from "../components/ActivityTab";
import { FaEllipsisVertical } from "react-icons/fa6";
import SuspendAdmin from "../components/SuspendAdmin";
import RoleManagement from "../components/RoleManagement";
import ChangeRole from "../components/ChangeRole";
import { useListAdminsQuery, User } from "@/redux/api/admin";
import SuccessMessage from "@/components/SuccessMessage";
import { MdVerified } from "react-icons/md";
import goldTag from "@/assets/images/goldtag.svg";
import ViolationsHistory from "../components/VoilationsHistory";
import RevokeSuspension from "../components/RevokeSuspension";
import RoleSuspendCard from "../components/RoleSuspendCard";

export default function AdminDetail() {
  const params = useParams();
  const username = params.username;
  let searchParams = useSearchParams();
  // read ?page=2 from the URL
  const page = Number(searchParams.get("page") || "1");
  const pageSize = Number(searchParams.get("pageSize") ?? 10);

  const {
    data: adminsData,
    isLoading,
    error,
  } = useListAdminsQuery({ page, pageSize });
  const admin = adminsData?.admins?.find((a) => a.Username === params.username);

  const [open, setOpen] = useState<boolean>(false);

  const [suspendAdmin, setSuspendAdmin] = useState<boolean>(false);
  const [changeAdminRole, setChangeAdminRole] = useState<boolean>(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [revoke, setRevoke] = useState<boolean>(false);
  const router = useRouter();
  const [showSuspendSuccess, setShowSuspendSuccess] = useState(false);

  const [suspendedEmail, setSuspendedEmail] = useState<string | null>(null);

  const handleSuspend = () => {
    setSuspendAdmin(!suspendAdmin);
  };

  const handleChangeRole = () => {
    setChangeAdminRole(!changeAdminRole);
  };
  const handleRevoke = () => {
    setRevoke(!revoke);
  };

  const adminCreatedAt = admin?.CreatedAt
    ? new Date(admin.CreatedAt).toLocaleDateString("en-US", {
        year: "numeric",
        month: "long",
        day: "numeric",
      })
    : "N/A";

  const details = [
    { icon: mailIcon, text: admin?.Email || "N/A" },
    {
      icon: calenderIcon,
      text: `Created on ${adminCreatedAt}`,
    },
    {
      icon: userIcon,
      text:
        admin?.Role?.split("_")
          .map(
            (word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()
          )
          .join(" ") || "N/A",
    },
  ];
  const [currentTab, setCurrentTab] = useState("activity");
  const tabs = [
    { key: "activity", label: "Activity" },
    { key: "adminaction", label: "Administrative Action" },
    { key: "violations", label: "Violations History" },
  ];

  const handleSuspendSuccess = (email: string) => {
    // 1) close the modal
    setSuspendAdmin(false);
    // 2) then show the success dialog
    setSuspendedEmail(email);
    setShowSuspendSuccess(true);
  };

  const viewProfile = () => {
    if (admin?.Username) {
      router.push(`/users/${admin?.Username}`);
    }
  };
  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading admin data</div>;

  return (
    <Container>
      <StyledLink href="/admin-users">
        <FaArrowLeft />
      </StyledLink>

      <Content>
        <AdminContainer>
          <StyledImage></StyledImage>

          <div>
            <Header>
              <div>
                <UserName>
                  {admin?.Username
                    ? admin?.Username.charAt(0).toUpperCase() +
                      admin?.Username.slice(1).toLowerCase()
                    : "N/A"}
                </UserName>

                <FullName>
                  <UserNames>
                    {admin?.FirstName
                      ? admin?.FirstName.charAt(0).toUpperCase() +
                        admin?.FirstName.slice(1).toLowerCase()
                      : "N/A"}
                  </UserNames>
                  <UserNames>
                    {admin?.LastName
                      ? admin.LastName.charAt(0).toUpperCase() +
                        admin?.LastName.slice(1).toLowerCase()
                      : "N/A"}
                  </UserNames>
                </FullName>
              </div>
            </Header>

            {open && (
              <InlineMenu>
                {admin?.Status === "SUSPENDED" ? (
                  <RoleSuspendCard
                    handleChangeRole={handleChangeRole}
                    handleRevoke={handleRevoke}
                  />
                ) : (
                  <RoleManagement
                    handleSuspend={handleSuspend}
                    handleChangeRole={handleChangeRole}
                  />
                )}
              </InlineMenu>
            )}
            {suspendAdmin && selectedUser && (
              <SuspendAdmin
                suspendAdmin={suspendAdmin}
                setSuspendAdmin={setSuspendAdmin}
                adminEmail={selectedUser.Email}
                onSuccess={handleSuspendSuccess}
              />
            )}
            {changeAdminRole && selectedUser && (
              <ChangeRole
                changeAdminRole={changeAdminRole}
                setChangeAdminRole={setChangeAdminRole}
                adminEmail={selectedUser.Email}
              />
            )}
            {revoke && selectedUser && (
              <RevokeSuspension
                revoke={revoke}
                setRevoke={setRevoke}
                adminName={`${selectedUser.FirstName} ${selectedUser.LastName}`}
                adminEmail={selectedUser.Email}
              />
            )}
            <AdminRow>
              <VerifiedTag>
                <p>Not Verified</p>
                <MdVerified color="#007CDF" />
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
              <UserStatus status={admin?.Status}>
                {admin?.Status
                  ? admin.Status.charAt(0).toUpperCase() +
                    admin.Status.slice(1).toLowerCase()
                  : "N/A"}
              </UserStatus>
            </AdminRow>
            <AdminInfoDetails>
              {details.map((list, index) => (
                <InfoContent key={index}>
                  <Image src={list.icon} alt={list.text} />
                  <InfoText>{list.text}</InfoText>
                </InfoContent>
              ))}
            </AdminInfoDetails>
            <ProfileButton onClick={viewProfile}>
              {" "}
              View as trovo user <FaAngleRight />{" "}
            </ProfileButton>
          </div>
        </AdminContainer>
        <ActionWrapper>
          <FaEllipsisVertical
            cursor="pointer"
            onClick={() => {
              setOpen(!open);
              setSelectedUser(admin!);
            }}
          />
        </ActionWrapper>
      </Content>

      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabButtonStyle={{
          borderRadius: "12px",
          fontSize: "16px",

          color: "#4F4F4F",
          width: "90%",
          padding: "0px",
          marginTop: "10px",
        }}
        activeTabButtonStyle={{
          backgroundColor: "#F2F6F9",
          color: "#00225A",
          border: "0px",
          width: "70%",
        }}
        tabContainerStyle={{
          width: "70%",
          borderRadius: "10px",
          height: "36px",
          padding: "1px",
          margin: "10px 0",
        }}
        tabContentStyle={{
          padding: "10px 10px 0 10px",
        }}
      >
        {currentTab === "activity" && <ActivityTab />}
        {currentTab === "adminaction" && <AdminActionsTab />}
        {currentTab === "violations" && (
          <ViolationsHistory adminEmail={admin?.Email} />
        )}
      </Tab>

      {showSuspendSuccess && suspendedEmail && (
        <SuccessMessage
          isOpen={showSuspendSuccess}
          setIsOpen={(open) => {
            setShowSuspendSuccess(open);
            if (!open) {
              // now we can tear down the modal
              setSuspendAdmin(false);
              router.push("/admin-users");
            }
          }}
          heading="Success"
          message="You have successfully suspended the admin user"
          email={suspendedEmail}
        />
      )}
    </Container>
  );
}

const Container = styled.section`
  padding: 32px;
  gap: 40px;
  border-radius: 24px;
`;
const Content = styled.div`
  display: flex;
  column-gap: 20px;
  justify-content: space-between;
  background-color: #ffffff;
  padding: 24px;
  border-radius: 16px;
`;

const StyledLink = styled(Link)`
  text-decoration: none
  display: block;
  color: #000000;
  width: 20px;
  padding-bottom: 30px;
`;

const AdminContainer = styled.div`
  display: flex;
  column-gap: 16px;
`;

const Header = styled.div`
  display: flex;
  column-gap: 16px;
`;
const FullName = styled.div`
  display: flex;
  align-items: center;
  column-gap: 4px;
`;

const UserName = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 16px;
  letter-spacing: 0.1px;
  margin: 0;
  padding-bottom: 8px;
  color: #00225a;
`;

const UserNames = styled.h1`
  font-size: 16px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  margin: 0;
  padding-bottom: 8px;
  color: #00225a;
`;

const StyledImage = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const AdminInfoDetails = styled.div`
  display: flex;
  column-gap: 16px;
  align-items: center;
`;

const InfoContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  column-gap: 8px;
`;
const InfoText = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 17.07px;
  text-align: left;
  color: #00225a;
`;

const ActionWrapper = styled.div`
  position: relative;
  display: inline-block;
  cursor: pointer;
  padding: 12px;
  margin: -8px;
`;

const InlineMenu = styled.div`
  position: absolute;
  top: 35%;
  right: 4%;
  width: 190px;
  margin-top: 8px;
  z-index: 1;
  background-color: white;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  padding: 8px 12px;
`;

const ProfileButton = styled.button`
  background-color: #f2f6f9;
  height: 48px;
  border-radius: 12px;
  padding: 12px;
  border: 0;
  margin-top: 8px;
  font-family: inherit;
  color: #00225a;
  font-weight: 500;
  font-size: 14px;
  line-height: 24px;
  letter-spacing: 0.25px;
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
`;
const ConfirmTag = styled.div`
  color: #ca8503;
  padding: 4px 16px;
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
const AdminRow = styled.div`
  display: flex;
  align-items: center;
  justify-content: start;
  margin-bottom: 10px;
  gap: 10px;
  width: 58%;
`;

const UserStatus = styled.div<{ status?: string }>`
  color: ${({ status }) => (status === "SUSPENDED" ? "#BE3800" : "#00A859")};
  background-color: ${({ status }) =>
    status === "SUSPENDED" ? "#BE38001A" : "#00a8591a"};
  font-weight: 500;
  padding: 4px 10px;
  border-radius: 6px;
  display: inline-flex;
  text-align: center;
  font-size: 12px;
  margin: 0 auto;
`;
