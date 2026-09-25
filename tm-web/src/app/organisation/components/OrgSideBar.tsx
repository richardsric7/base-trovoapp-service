"use client";
import React from "react";

import Image from "next/image";
import Logo from "@/assets/images/trovologosvg.svg";
import styled from "styled-components";
import { usePathname, useRouter } from "next/navigation";
import { useGetOrganizationDetailsDataQuery } from "@/redux/api/org/api";
import { getOrganizationNavigation } from "./organizationNavigation";

interface StyledProps {
  $isActive?: boolean;
}
interface SubMenuTextProps {
  $isActive?: boolean;
}
const OrgSideBar = () => {
  const router = useRouter();

  const { data } = useGetOrganizationDetailsDataQuery();
  const orgName = data?.data.organization.name ?? "";
  const stakeholderType =
    data?.data.organization.stakeholder_type || data?.data.organization.type;
  const navigationItems = getOrganizationNavigation(stakeholderType);
  const orgType = data?.data?.organization?.type
    ?.toLowerCase()
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");

  const pathname = usePathname();

  // Simplified isActive check
  const isActive = (path: string): boolean => {
    if (path === "/organisation") {
      return pathname === "/organisation";
    }
    return pathname.startsWith(path);
  };

  // logout function
  const handleLogout = () => {
    // router.push("/organization/login");
    // showSuccessToast("Logged out");
  };
  return (
    <Container>
      <LogoContainer>
        <Image src={Logo} alt="logo" width={180} height={40} priority />
      </LogoContainer>
      <ScrollableMenuContainer>
        <MenuContainer>
          <MenuItems>
            {navigationItems.map((item) => {
              const isItemActive = item.exact
                ? pathname === item.link
                : isActive(item.link);
              return (
                <MenuItem
                  key={item.link}
                  onClick={() => router.push(item.link)}
                  $isActive={isItemActive}
                >
                  <IconWrapper $isActive={isItemActive}>
                    <Image src={item.icon} alt="menu-icon" />
                  </IconWrapper>
                  <MenuText $isActive={isItemActive}>{item.title}</MenuText>
                </MenuItem>
              );
            })}
          </MenuItems>
        </MenuContainer>
        <OrgInfoSection>
          <SidebarLogoCircle>
            {orgName.charAt(0).toUpperCase() || "R"}
          </SidebarLogoCircle>
          <SidebarOrgDetails>
            <SidebarOrgName>
              {orgName || "Rand Merchant Bank..."}
            </SidebarOrgName>
            <SidebarOrgRole>{orgType || "Asset custodian"}</SidebarOrgRole>
          </SidebarOrgDetails>
        </OrgInfoSection>
      </ScrollableMenuContainer>
    </Container>
  );
};

export default OrgSideBar;

const Container = styled.aside`
  height: 100vh;
  width: 252px;
  padding: 40px 10px 0px 10px;
  position: fixed;
  background-color: #ffffff;
  left: 0;
  top: 0;
  bottom: 0;
  box-shadow: 4px 0px 12px 0px #acd1ef1a;
  display: flex;
  flex-direction: column;
`;

const LogoContainer = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
  column-gap: 8px;
  margin-bottom: 8px;
`;

// New ScrollableMenuContainer to handle overflow
const ScrollableMenuContainer = styled.div`
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  overflow-y: auto;
  min-height: 0;

  /* For Firefox */
  scrollbar-width: 0;
  scrollbar-color: #007cdf;

  /* For Chrome, Safari, and Opera */
  &::-webkit-scrollbar {
    width: 0;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background-color: #007cdf;
    border-radius: 1px;
  }
`;

const MenuItems = styled.div``;
const MenuContainer = styled.div``;
const MenuText = styled.div<StyledProps>`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  letter-spacing: 0.10000000149011612px;
  text-align: center;
  color: ${(props) => (props.$isActive ? "#007cdf" : "#828282")};
  transition: color 0.2s ease;
`;

const SubMenuText = styled(MenuText)<SubMenuTextProps>`
  font-size: 12px;
  color: ${(props) => (props.$isActive ? "#007cdf" : "#828282")};
`;
const ContentWrapper = styled.div`
  display: flex;
  align-items: center;
  column-gap: 8px;
`;
const MenuContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  column-gap: 4px;
  padding: 14px 2px;
  background-color: transparent;
  border: none;
  cursor: pointer;
  width: 100%;
  flex: 1;
`;
const IconWrapper = styled.div<StyledProps>`
  display: flex;
  align-items: center;
  ${(props) =>
    props.$isActive &&
    `
    filter: brightness(0) saturate(100%) invert(32%) sepia(71%)
      saturate(1863%) hue-rotate(176deg) brightness(91%) contrast(93%);
  `}
`;
const MenuItem = styled.div<StyledProps>`
  display: flex;
  align-items: center;
  column-gap: 6px;
  padding: 14px 4px;
  background-color: transparent;
  border: none;
  cursor: pointer;
  width: 100%;
  color: ${(props) => (props.$isActive ? "#007cdf" : "#828282")};
  flex: 1;
  text-align: left;

  &:hover {
    ${MenuText}, ${IconWrapper} {
      color: #007cdf;
      filter: brightness(0) saturate(100%) invert(32%) sepia(71%)
        saturate(1863%) hue-rotate(176deg) brightness(91%) contrast(93%);
    }
  }
`;
const SubMenuItem = styled(MenuItem)`
  padding: 10px 30px;
`;

const SubMenu = styled.div`
  display: flex;
  flex-direction: column;
  position: relative;
  padding-left: 1px;

  &::before {
    content: "";
    position: absolute;
    left: 20px;
    top: 10px;
    height: calc(100% - 20px);
    width: 1px;
    background-color: #e5e5ef;
  }
`;

const OrgInfoSection = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px;
  background: #f8fbff;
  border-radius: 12px;
  margin-top: auto;
  margin-bottom: 20px;

  @media (min-width: 1400px) {
    position: absolute;
    bottom: 20px;
    left: 10px;
    right: 10px;
  }
`;

const SidebarLogoCircle = styled.div`
  background: #007cdf;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #ffffff;
  font-weight: 600;
  font-size: 16px;
  flex-shrink: 0;
`;

const SidebarOrgDetails = styled.div`
  display: flex;
  flex-direction: column;
  overflow: hidden;
`;

const SidebarOrgName = styled.p`
  font-weight: 600;
  font-size: 13px;
  line-height: 1.2;
  color: #00225a;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const SidebarOrgRole = styled.p`
  font-weight: 400;
  font-size: 11px;
  color: #828282;
  margin: 0;
`;
