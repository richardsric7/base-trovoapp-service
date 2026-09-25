"use client";

import React, { useEffect, useMemo, useRef, useState } from "react";
import styled from "styled-components";
import Image from "next/image";
import profileIcon from "@/assets/images/profile.svg";
import settingIcon from "@/assets/images/setting-2.svg";
import logoutIcon from "@/assets/images/logout-icon.svg";
import Link from "next/link";
import { FaAngleDown, FaAngleUp } from "react-icons/fa6";
import { useLogoutMutation, USER_DETAILS, TOKEN } from "@/redux";
import { setToken, setUser } from "@/redux/slices/authSlice";
import { clearStorage } from "@/utils";
import { deleteCookie } from "cookies-next";
import { usePathname, useRouter } from "next/navigation";
import { showSuccessToast } from "@/components";
import { useAppDispatch, useAppSelector } from "@/lib/hooks";

export const NavBar = () => {
  const [openUserProfile, setOpenUserProfile] = useState<boolean>(false);
  const [onLogout] = useLogoutMutation();
  const router = useRouter();
  const dispatch = useAppDispatch();
  const pathname = usePathname();
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Load `userInfo` from `localStorage` on mount
    const savedUser = localStorage.getItem(USER_DETAILS);
    if (savedUser) {
      dispatch(setUser(JSON.parse(savedUser)));
    }

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

  const userInfo = useAppSelector((state) => state.auth.user);

  // Helper function to safely access nested user properties
  const getUserDisplayInfo = () => {
    if (!userInfo) return { firstName: "", lastName: "", email: "", image: "" };

    // this describes or handles the  possible structures of userInfo
    const user = userInfo.userInfo || userInfo;

    return {
      firstName: user?.first_name || "",
      lastName: user?.last_name || "",
      email: user?.email || "",
      // image: user?.imageThumbnail || "",
    };
  };
  const { firstName, lastName, email } = useMemo(
    () => getUserDisplayInfo(),
    [userInfo]
  );

  // console.log("User info after refresh:", userInfo);

  const handleLogout = async () => {
    try {
      await onLogout().unwrap();
    } catch {
      // Local cleanup below still completes the logout even if the API call failed
      // (offline, brief 5xx, already-expired token) — the user should never end up stuck
      // signed in just because this one request didn't succeed.
    }

    dispatch(setToken(null));
    dispatch(setUser(null));
    setOpenUserProfile(false);
    clearStorage();
    deleteCookie(TOKEN);
    router.push("/sign-in");
    showSuccessToast("Logged out");
  };

  return (
    <Container>
      {/* <HeaderWrapper> */}
      {/* <SearchContainer>
          <SearchBar placeholder="Search" />
        </SearchContainer> */}
      <UserInfoSection onClick={() => setOpenUserProfile(!openUserProfile)}>
        <UserInfoDetails>
          <Avatar></Avatar>
          <div>
            <UserDetails>
              <UserName>{firstName}</UserName>
              <UserName>{lastName}</UserName>
            </UserDetails>
            <UserEmail>{email}</UserEmail>
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
            <Avatar>
              {" "}
              {/* {image && <Image src={image} alt={firstName} objectFit="cover" />} */}
            </Avatar>{" "}
            <div>
              <UserDetails>
                <UserName>{firstName}</UserName>
                <UserName>{lastName}</UserName>
              </UserDetails>
              <UserEmail>{email}</UserEmail>
            </div>
          </UserInfoSections>
          <Divider />
          <ProfileActions>
            <div>
              <ProfileAction>
                <Image src={profileIcon} alt="profile-icon" />
                <ActionText href="">Manage Profile</ActionText>
              </ProfileAction>
              <ProfileAction>
                <Image src={settingIcon} alt="profile-icon" />
                <ActionText href=""> Settings</ActionText>
              </ProfileAction>
            </div>
            <Divider />
            <ProfileAction onClick={handleLogout}>
              <Image src={logoutIcon} alt="profile-icon" />
              <LogoutText>Logout</LogoutText>
            </ProfileAction>
          </ProfileActions>
        </UserProfileDropdown>
      )}
      {/* </HeaderWrapper> */}
    </Container>
  );
};

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
  // padding: 25px 38px;
`;

const SearchContainer = styled.div``;
const UserInfoSection = styled.div`
  display: flex;
  align-items: center;
  column-gap: 10px;
  position: absolute;
  top: 25px;
  right: 25px;
  // margin-bottom: 20px;
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
