"use client";
import React, { useState } from "react";
import styled from "styled-components";
import Logo from "@/assets/images/trovologosvg.svg";
import logoutIcon from "@/assets/images/Icon.svg";
import Image from "next/image";
import DashboardIcon from "@/assets/images/dashboardIcon.svg";
import accessIcon from "@/assets/images/access.svg";
import securityUser from "@/assets/images/security-user.svg";
import emptyWallet from "@/assets/images/empty-wallet.svg";
import settingIcon from "@/assets/images/setting-2.svg";

import noteIcon from "@/assets/images/complianceIcon.svg";
import { usePathname, useRouter } from "next/navigation";
import keyIcon from "@/assets/images/key-square.svg";

interface StyledProps {
  $isActive?: boolean;
}

interface SubMenuTextProps {
  $isActive?: boolean;
}

export const Sidebar = () => {
  const router = useRouter();
  const pathname = usePathname();

  const sidebarItems = [
    {
      title: "Dashboard",
      icon: DashboardIcon,
      link: "/external-api-clients",
    },
    {
      title: "Customers",
      icon: securityUser,
      link: "/external-api-clients/customers",
    },
    {
      title: "Wallets",
      icon: emptyWallet,
      link: "/external-api-clients/wallets",
    },

    {
      title: "Token Management",
      icon: noteIcon,
      link: "/external-api-clients/token-management",
    },
    {
      title: "Api Keys",
      icon: keyIcon,
      link: "/external-api-clients/api-keys",
    },
    {
      title: "Settings",
      icon: settingIcon,
      link: "/external-api-clients/settings",
    },
  ];

  const isActive = (path: string): boolean => {
    // Dashboard (exact match only)
    if (path === "/external-api-clients") {
      return pathname === "/external-api-clients";
    }

    // Other pages (allow nested routes)
    return pathname.startsWith(path);
  };

  const handleLogout = () => {
    // router.push("/sign-in");
  };

  return (
    <Container>
      <LogoContainer onClick={() => router.push("/")}>
        <Image src={Logo} alt="logo" width={180} height={40} priority />
      </LogoContainer>
      <ScrollableMenuContainer>
        <MenuContainer>
          <MenuItems>
            {sidebarItems.map((item, index) => {
              const isItemActive = isActive(item.link);
              return (
                <MenuItem
                  key={index}
                  onClick={() => router.push(item.link)}
                  $isActive={isItemActive}
                >
                  <IconWrapper $isActive={isItemActive}>
                    <Image
                      src={item.icon}
                      alt={item.title}
                      width={20}
                      height={20}
                    />
                  </IconWrapper>
                  <MenuText $isActive={isItemActive}>{item.title}</MenuText>
                </MenuItem>
              );
            })}
          </MenuItems>
        </MenuContainer>
        <LogoutSection onClick={handleLogout}>
          <IconWrapper>
            <Image src={logoutIcon} alt="logout" width={20} height={20} />
          </IconWrapper>
          <MenuText>Logout</MenuText>
        </LogoutSection>

        <OrgInfoSection>
          <SidebarLogoCircle>R</SidebarLogoCircle>
          <SidebarOrgDetails>
            <SidebarOrgName>Rand Merchant Bank</SidebarOrgName>
            <SidebarOrgRole>Asset custodian</SidebarOrgRole>
          </SidebarOrgDetails>
        </OrgInfoSection>
      </ScrollableMenuContainer>
    </Container>
  );
};

const Container = styled.aside`
  height: 100vh;
  width: 252px;
  padding: 40px 10px 20px 10px;
  position: fixed;
  background-color: #ffffff;
  left: 0;
  top: 0;
  bottom: 0;
  box-shadow: 4px 0px 12px 0px #acd1ef1a;
  display: flex;
  flex-direction: column;
  z-index: 100;
`;

const LogoContainer = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
  column-gap: 8px;
  margin-bottom: 32px;
  cursor: pointer;
  padding-left: 10px;
`;

const ScrollableMenuContainer = styled.div`
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  overflow-y: auto;
  min-height: 0;

  &::-webkit-scrollbar {
    width: 0;
  }
`;

const MenuContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const MenuItems = styled.div``;

const MenuText = styled.span<StyledProps>`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  color: ${(props) => (props.$isActive ? "#007cdf" : "#828282")};
  transition: color 0.2s ease;
`;

const IconWrapper = styled.div<StyledProps>`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
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
  column-gap: 12px;
  padding: 12px 16px;
  background-color: ${(props) => (props.$isActive ? "#f8fbff" : "transparent")};
  border-radius: 12px;
  cursor: pointer;
  width: 100%;
  margin-bottom: 4px;
  transition: all 0.2s ease;

  &:hover {
    background-color: #f8fbff;
    ${MenuText} {
      color: #007cdf;
    }
    ${IconWrapper} {
      filter: brightness(0) saturate(100%) invert(32%) sepia(71%)
        saturate(1863%) hue-rotate(176deg) brightness(91%) contrast(93%);
    }
  }
`;

const LogoutSection = styled(MenuItem)`
  margin-top: 0;
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
