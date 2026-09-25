"use client";
import React, { useEffect, useState } from "react";
import styled from "styled-components";
import Logo from "@/assets/images/trovologosvg.svg";
import logoutIcon from "@/assets/images/Icon.svg";
import Image from "next/image";
import Dashboard from "@/assets/images/dashboardIcon.svg";
import appmanagent from "@/assets/images/mage_dashboard.svg";
import p2pIcon from "@/assets/images/money-change.svg";
import organizationaIcon from "@/assets/images/buildings-2.svg";
import AuditTrail from "@/assets/images/AuditTrailIcon.svg";
// Reuses the dashboard glyph: System Health is a status overview, and adding a
// new icon asset for it is a design decision, not an engineering one.
import SystemHealth from "@/assets/images/mage_dashboard.svg";
import asseticon from "@/assets/images/chainlink-(link).svg";
import settingIcon from "@/assets/images/setting-2.svg";
import AdminManagement from "@/assets/images/adminManagementIcon.svg";
import Payments from "@/assets/images/money.svg";
import Compliance from "@/assets/images/complianceIcon.svg";
import Support from "@/assets/images/streamline_interface-help-customer-support-2-customer-headphones-headset-help-microphone-phone-person-support.svg";
import NotificationIcon from "@/assets/images/notificationBellIcon.svg";
import Revenue from "@/assets/images/revenueIcon.svg";
import VaultSignerIcon from "@/assets/images/key-square.svg";
import { FaAngleUp, FaAngleDown } from "react-icons/fa";
import breifcaseIcon from "@/assets/images/briefcase.svg";
import { usePathname, useRouter } from "next/navigation";
import { useLogoutMutation, TOKEN } from "@/redux";
import { setToken, setUser } from "@/redux/slices/authSlice";
import { clearStorage } from "@/utils";
import { deleteCookie } from "cookies-next";
import { showSuccessToast } from "@/components";
import { useAppDispatch, useAppSelector } from "@/lib/hooks";

interface StyledProps {
  $isActive?: boolean;
}
interface SubMenuTextProps {
  $isActive?: boolean;
}

export const SideBar = () => {
  const sidebarItems = [
    {
      title: "Organizations",
      icon: organizationaIcon,
      link: "/organizations",
    },
    {
      title: "Admin Management",
      icon: AdminManagement,
      link: "/admin-users",
    },

    { title: "Revenue", icon: Revenue, link: "/revenue" },

    { title: "Vault Signer", icon: VaultSignerIcon, link: "/vault-signer" },

    { title: "Payments", icon: Payments, link: "/payment" },
    { title: "Customer Support", icon: Support, link: "/customer-support" },
    {
      title: "Compliance & AML",
      icon: Compliance,
      link: "/complaince-and-aml",
    },

    // { title: "Moderation", icon: Support, link: "/dashboard/users" },
    {
      title: "Notification",
      icon: NotificationIcon,
      link: "/notification",
    },

    {
      title: "Audit Trail",
      icon: AuditTrail,
      link: "/audit-trail",
    },

    {
      title: "System Health",
      icon: SystemHealth,
      link: "/system-health",
    },
    {
      title: "Jobs",
      icon: breifcaseIcon,
      link: "/jobs",
    },
    {
      title: "Settings",
      icon: settingIcon,
      link: "/settings",
    },
  ];
  const router = useRouter();

  const [isMenuOpen, setIsMenuOpen] = useState<string[]>([]);

  const pathname = usePathname();

  const getMenuSection = (path: string) => {
    if (path === "/" || path.startsWith("/trovop2p")) return "dashboard";
    if (path.startsWith("/users")) return "trovoAppManagement";
    if (path.startsWith("/p2pappeal") || path.startsWith("/p2pusers"))
      return "trovop2p";
    if (
      path.startsWith("/assettokenization") ||
      path.startsWith("/othertoken") ||
      path.startsWith("/assetcuration")
    )
      return "asset&tokens";
    return null;
  };

  useEffect(() => {
    const currentSection = getMenuSection(pathname);
    if (currentSection) {
      setIsMenuOpen([currentSection]); // Only keep current section open
    }
  }, [pathname]);

  const handleToggleMenu = (menuName: string) => {
    setIsMenuOpen((prev) => {
      if (prev.includes(menuName)) {
        // If menu is open, just close this one
        return prev.filter((name) => name !== menuName);
      } else {
        // If menu is closed, add it to open menus
        return [...prev, menuName];
      }
    });
  };

  // Simplified isActive check
  const isActive = (path: string): boolean => {
    if (path === "/") return pathname === "/";
    return pathname.startsWith(path);
  };

  // logout function
  const [onLogout] = useLogoutMutation();
  const dispatch = useAppDispatch();
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
    clearStorage();
    deleteCookie(TOKEN);
    router.push("/sign-in");
    showSuccessToast("Logged out");
  };

  return (
    <Container>
      <LogoContainer>
        <Image src={Logo} alt="logo" width={180} height={40} priority />
      </LogoContainer>
      <ScrollableMenuContainer>
        <MenuContainer>
          <MenuContent onClick={() => handleToggleMenu("dashboard")}>
            <ContentWrapper>
              <IconWrapper>
                <Image src={Dashboard} alt="admin-icon" />
              </IconWrapper>
              <MenuText>Dashboard</MenuText>
            </ContentWrapper>
            {isMenuOpen?.includes("dashboard") ? (
              <FaAngleUp color="#828282" cursor="pointer" />
            ) : (
              <FaAngleDown color="#828282" />
            )}
          </MenuContent>
          {isMenuOpen?.includes("dashboard") && (
            <SubMenu>
              <SubMenuItem onClick={() => router.push("/")}>
                <SubMenuText $isActive={isActive("/")}>Trovo App</SubMenuText>
              </SubMenuItem>
              <SubMenuItem onClick={() => router.push("/trovop2p")}>
                <SubMenuText $isActive={isActive("/trovop2p")}>
                  Trovo P2P
                </SubMenuText>
              </SubMenuItem>
            </SubMenu>
          )}

          <MenuContent onClick={() => handleToggleMenu("trovoAppManagement")}>
            <ContentWrapper>
              <IconWrapper>
                <Image src={appmanagent} alt="admin-icon" />
              </IconWrapper>
              <MenuText>TrovoApp Management</MenuText>
            </ContentWrapper>
            {isMenuOpen.includes("trovoAppManagement") ? (
              <FaAngleUp color="#828282" />
            ) : (
              <FaAngleDown color="#828282" />
            )}
          </MenuContent>
          {isMenuOpen?.includes("trovoAppManagement") && (
            <SubMenu>
              <SubMenuItem onClick={() => router.push("/users")}>
                <SubMenuText $isActive={isActive("/users")}>
                  User Management
                </SubMenuText>
              </SubMenuItem>
            </SubMenu>
          )}

          <MenuContent onClick={() => handleToggleMenu("trovop2p")}>
            <ContentWrapper>
              <IconWrapper>
                <Image src={p2pIcon} alt="admin-icon" />
              </IconWrapper>
              <MenuText>TrovoP2P Management</MenuText>
            </ContentWrapper>
            {isMenuOpen.includes("trovop2p") ? (
              <FaAngleUp color="#828282" />
            ) : (
              <FaAngleDown color="#828282" />
            )}
          </MenuContent>
          {isMenuOpen.includes("trovop2p") && (
            <SubMenu>
              <SubMenuItem onClick={() => router.push("/p2pusers")}>
                <SubMenuText $isActive={isActive("/p2pusers")}>
                  User Management
                </SubMenuText>
              </SubMenuItem>
              <SubMenuItem onClick={() => router.push("/p2pappeal")}>
                <SubMenuText $isActive={isActive("/p2pappeal")}>
                  P2P Appeal Resolution
                </SubMenuText>
              </SubMenuItem>

              <SubMenuItem onClick={() => router.push("/trades")}>
                <SubMenuText $isActive={isActive("/trades")}>
                  Trades
                </SubMenuText>
              </SubMenuItem>
            </SubMenu>
          )}

          <MenuContent onClick={() => handleToggleMenu("asset&tokens")}>
            <ContentWrapper>
              <IconWrapper>
                <Image src={asseticon} alt="asset&token-icon" />
              </IconWrapper>
              <MenuText>Assets & Tokens</MenuText>
            </ContentWrapper>
            {isMenuOpen.includes("asset&tokens") ? (
              <FaAngleUp color="#828282" />
            ) : (
              <FaAngleDown color="#828282" />
            )}
          </MenuContent>
          {isMenuOpen.includes("asset&tokens") && (
            <SubMenu>
              <SubMenuItem onClick={() => router.push("/assettokenization")}>
                <SubMenuText $isActive={isActive("/assettokenization")}>
                  Tokenized Assets
                </SubMenuText>
              </SubMenuItem>

              <SubMenuItem onClick={() => router.push("/othertoken")}>
                <SubMenuText $isActive={isActive("/othertoken")}>
                  Other Tokens
                </SubMenuText>
              </SubMenuItem>
              <SubMenuItem onClick={() => router.push("/assetcuration")}>
                <SubMenuText $isActive={isActive("/assetcuration")}>
                  Asset Curation
                </SubMenuText>
              </SubMenuItem>
            </SubMenu>
          )}

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
                    <Image src={item.icon} alt="menu-icon" />
                  </IconWrapper>
                  <MenuText $isActive={isItemActive}>{item.title}</MenuText>
                </MenuItem>
              );
            })}
          </MenuItems>
        </MenuContainer>
        <Logout onClick={handleLogout}>
          <IconWrapper>
            <Image src={logoutIcon} alt="logout-icon" />
          </IconWrapper>
          <MenuText>Logout</MenuText>
        </Logout>
      </ScrollableMenuContainer>
    </Container>
  );
};

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
const Logout = styled(MenuItem)`
  margin-top: 40px;
  bottom: 0;
  // @media (min-width: 1400px) {
  //   position: absolute;
  //   bottom: 0;
  // }
`;
