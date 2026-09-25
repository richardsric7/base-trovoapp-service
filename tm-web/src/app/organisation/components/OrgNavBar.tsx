import React, { useEffect, useRef, useState } from "react";
import { FaAngleDown, FaAngleUp } from "react-icons/fa6";
import profileIcon from "@/assets/images/profile.svg";
import settingIcon from "@/assets/images/setting-2.svg";
import logoutIcon from "@/assets/images/logout-icon.svg";
import { usePathname, useRouter } from "next/navigation";
import styled from "styled-components";
import Image from "next/image";
import { useDispatch } from "react-redux";
import Link from "next/link";
import { useGetStakeholderProfileQuery } from "@/redux/api/sharedstakeholders";
import { FaWallet } from "react-icons/fa6";
import LinkWalletModal from "./LinkWalletModal";
import { deleteCookie } from "cookies-next";
import { TOKEN_ORG } from "@/redux";
import { clearOrgAuth } from "@/redux/slices/orgSlice";
import { useLogoutOrganizationMemberMutation } from "@/redux/api/org/api";
import { orgApi } from "@/redux/baseApi/orgApi";
import { showSuccessToast } from "@/components";
const OrgNavBar = () => {
  const [openUserProfile, setOpenUserProfile] = useState<boolean>(false);
  const [openLinkWallet, setOpenLinkWallet] = useState(false);
  const dispatch = useDispatch();
  const router = useRouter();
  const pathname = usePathname();
  const dropdownRef = useRef<HTMLDivElement>(null);
  const [logoutOrganizationMember, { isLoading: isLoggingOut }] =
    useLogoutOrganizationMemberMutation();

  const { data: profileData, refetch: refetchProfile } =
    useGetStakeholderProfileQuery();
  const member = profileData?.data?.member;
  const wallet = profileData?.data?.wallet;
  const userName = member ? `${member.first_name} ${member.last_name}` : "";
  const userEmail = member?.email ?? "";

  useEffect(() => {
    // Load `userInfo` from `localStorage` on mount
    // const savedUser = localStorage.getItem(USER_DETAILS);
    // if (savedUser) {
    //   dispatch(setUser(JSON.parse(savedUser)));
    // }

    const handleScroll = () => setOpenUserProfile(false);
    window.addEventListener("scroll", handleScroll);

    setOpenUserProfile(false);

    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setOpenUserProfile(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);

    // Cleanup event listener
    return () => {
      window.removeEventListener("scroll", handleScroll);
      // Was addEventListener: every route change leaked two mousedown handlers,
      // so the tab slowed to a halt after minutes of navigation.
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [dispatch, pathname]);
  const handleLogout = async () => {
    try {
      await logoutOrganizationMember().unwrap();
    } catch {
      // Organization JWTs are stateless, so local cleanup still completes logout.
    }

    deleteCookie(TOKEN_ORG);
    dispatch(clearOrgAuth());
    dispatch(orgApi.util.resetApiState());
    setOpenUserProfile(false);
    router.replace("/organizations/login");
    showSuccessToast("Logged out");
  };
  return (
    <Container>
      <HeaderWrapper>
        <HeadingText>DashBoard</HeadingText>
        <UserInfoSection onClick={() => setOpenUserProfile(!openUserProfile)}>
          <UserInfoDetails>
            <Avatar></Avatar>
            <div>
              <UserDetails>
                <UserName>{userName}</UserName>
              </UserDetails>
              <UserEmail>{userEmail}</UserEmail>
            </div>
          </UserInfoDetails>
          {openUserProfile ? (
            <FaAngleUp color="#828282" cursor="pointer" />
          ) : (
            <FaAngleDown color="#828282" cursor="pointer" />
          )}
        </UserInfoSection>
        {openUserProfile && (
          <UserProfileDropdown ref={dropdownRef}>
            <UserInfoSections>
              <Avatar></Avatar>{" "}
              <div>
                <UserDetails>
                  <UserName>{userName}</UserName>
                </UserDetails>
                <UserEmail>{userEmail}</UserEmail>
              </div>
            </UserInfoSections>
            <Divider />
            <ProfileActions>
              <div>
                <ProfileAction>
                  <Image src={profileIcon} alt="profile-icon" />
                  <ActionText href="/organisation/profile">
                    Manage Profile
                  </ActionText>
                </ProfileAction>
                <ProfileAction>
                  <Image src={settingIcon} alt="profile-icon" />
                  <ActionText href="/organisation/profile#notifications">
                    {" "}
                    Settings
                  </ActionText>
                </ProfileAction>
                {!wallet?.is_wallet_linked && (
                  <ProfileAction
                    onClick={() => {
                      setOpenUserProfile(false);
                      setOpenLinkWallet(true);
                    }}
                  >
                    <WalletActionIcon />
                    <ActionButton type="button">Link Wallet</ActionButton>
                  </ProfileAction>
                )}
                {wallet?.is_wallet_linked && (
                  <ProfileAction>
                    <LinkedWalletIcon />
                    <LinkedWalletText>
                      {wallet.trovo_wallet_username || "Wallet linked"}
                    </LinkedWalletText>
                  </ProfileAction>
                )}
              </div>
              <Divider />
              <ProfileAction
                onClick={isLoggingOut ? undefined : handleLogout}
                aria-disabled={isLoggingOut}
              >
                <Image src={logoutIcon} alt="profile-icon" />
                <LogoutText>
                  {isLoggingOut ? "Logging out..." : "Logout"}
                </LogoutText>
              </ProfileAction>
            </ProfileActions>
          </UserProfileDropdown>
        )}
      </HeaderWrapper>
      <LinkWalletModal
        isOpen={openLinkWallet}
        onClose={() => setOpenLinkWallet(false)}
        onLinked={refetchProfile}
      />
    </Container>
  );
};

export default OrgNavBar;
const Container = styled.nav`
  position: fixed;
  height: 90px;
  top: 0;
  right: 0;
  width: calc(100% - 210px);
  display: block;
  background-color: #ffffff;
  z-index: 100;
`;
const HeaderWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
`;

const HeadingText = styled.h1`
  font-weight: 600;
  font-size: 32px;
  margin-left: 32px;
  letter-spacing: 0.1px;
  color: #00225a;
`;
const UserInfoSection = styled.div`
  display: flex;
  align-items: center;
  column-gap: 10px;
`;
const UserInfoSections = styled.div`
  display: flex;
  align-items: center;
  column-gap: 10px;
  margin-bottom: 20px;
`;
const UserInfoDetails = styled.div`
  display: flex;
  column-gap: 10px;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
`;
const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
  overflow: hidden;
`;

const UserDetails = styled.div`
  display: flex;
  align-items: center;
  column-gap: 6px;
`;
const UserName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;
const UserEmail = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const UserProfileDropdown = styled.div`
  position: absolute;
  top: 102px;
  right: 38px;
  background: #fff;
  box-shadow: 0px 4px 100px 0px #00000026;
  width: 300px;
  padding: 24px;
  gap: 16px;
  border-radius: 24px;
  z-index: 10000;
`;

const Divider = styled.div`
  background-color: #e5e5ef;
  height: 1px;
  width: calc(100% + 48px);
  margin: 10px -24px 10px -24px;
`;

const ProfileActions = styled.div`
  display: flex;
  column-gap: 10px;
  flex-direction: column;
  gap: 10px;
  margin-top: 18px;
`;

const ProfileAction = styled.div`
  display: flex;
  column-gap: 10px;
  align-items: center;
  cursor: pointer;
  margin-bottom: 20px;

  &:last-child {
    margin-bottom: 0;
  }
`;

const ActionText = styled(Link)`
  color: #828282;
  font-size: 16px;
  font-weight: 500;
  text-decoration: none;
`;
const LogoutText = styled.p`
  color: #828282;
  font-size: 16px;
  font-weight: 500;
  text-decoration: none;
`;
const WalletActionIcon = styled(FaWallet)`
  width: 20px;
  height: 20px;
  color: #828282;
`;
const LinkedWalletIcon = styled(WalletActionIcon)`
  color: #00875a;
`;
const ActionButton = styled.button`
  padding: 0;
  border: 0;
  background: transparent;
  color: #828282;
  font: inherit;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
`;
const LinkedWalletText = styled.span`
  color: #00875a;
  font-size: 14px;
  font-weight: 500;
`;
