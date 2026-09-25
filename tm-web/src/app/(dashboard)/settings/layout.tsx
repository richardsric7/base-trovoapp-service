"use client";
import React from "react";
import styled from "styled-components";
import Link from "next/link";
import { usePathname } from "next/navigation";

export default function SettingLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const pathname = usePathname();

  return (
    <SettingContainer>
      <SettingsMenu>
        <MenuContainer>
          <MenuTitle>Tokenization Settings</MenuTitle>
          <MenuContent>
            <MenuLink href="/settings" pathname={pathname}>
              Minting Initiators
            </MenuLink>
            <MenuLink href="/settings/mintingapprovers" pathname={pathname}>
              Minting Approvers
            </MenuLink>
            {/* <MenuLink
              href="/settings/tokenizationstakeholders"
              pathname={pathname}
            >
              Stakeholders
            </MenuLink> */}
            <MenuLink href="/settings/assetsector" pathname={pathname}>
              Asset Sectors
            </MenuLink>
          </MenuContent>
        </MenuContainer>

        <MenuContainer>
          <MenuTitle>Countries Settings</MenuTitle>
          <MenuContent>
            <MenuLink href="/settings/countries" pathname={pathname}>
              Countries
            </MenuLink>
          </MenuContent>
        </MenuContainer>

        <MenuContainer>
          <MenuTitle>Fees Settings</MenuTitle>
          <MenuContent>
            <MenuLink href="/settings/trovoappfees" pathname={pathname}>
              Trovo App Fees
            </MenuLink>
            <MenuLink href="/settings/trovop2pfees" pathname={pathname}>
              TrovoP2P Fees
            </MenuLink>
          </MenuContent>
        </MenuContainer>

        <MenuContainer>
          <MenuTitle>Vault Signer Settings</MenuTitle>
          <MenuContent>
            <MenuLink href="/settings/vault-signer" pathname={pathname}>
              Managed Secrets & Audit Log
            </MenuLink>
          </MenuContent>
        </MenuContainer>
      </SettingsMenu>
      <SettingContent>{children}</SettingContent>
    </SettingContainer>
  );
}

const SettingContainer = styled.section`
  padding: 10px;
  display: flex;
  align-items: flex-start;
  gap: 32px;
`;

const SettingsMenu = styled.div`
  background-color: #fff;
  padding: 20px;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  gap: 36px;
  width: 244px;
`;

const SettingContent = styled.div`
  background-color: #fff;
  padding: 32px;
  border-radius: 24px;
  flex: 1;
`;

const MenuTitle = styled.p`
  font-weight: 600;
  font-size: 16px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #00225a;
  width: 100%;
`;

const MenuContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

const MenuContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 24px;
  width: 100%;
`;

// Active-aware Link Component
interface MenuLinkProps {
  href: string;
  pathname: string;
  children: React.ReactNode;
}

const MenuLink = ({ href, pathname, children }: MenuLinkProps) => {
  const isActive = pathname === href;
  return (
    <StyledMenuLink href={href} $active={isActive}>
      {children}
    </StyledMenuLink>
  );
};

const StyledMenuLink = styled(Link)<{ $active?: boolean }>`
  font-weight: 500;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: ${({ $active }) => ($active ? "#007cdf" : "#00225a")};
  text-decoration: none;
  padding: ${({ $active }) => ($active ? "6px 6px" : "6px 6px")};
  border-radius: 6px;
  background-color: ${({ $active }) => ($active ? "#f2f6f9" : "transparent")};
  transition:
    background-color 0.2s,
    color 0.2s,
    padding 0.2s;
  display: inline-block;
  width: fit-content;
  min-width: 0;

  &:hover {
    background-color: #f2f6f9;
    color: #007cdf;
    padding: 6px;
  }
`;
