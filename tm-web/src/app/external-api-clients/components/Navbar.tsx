"use client";
import React, { useState, useRef, useEffect, useMemo } from "react";
import styled from "styled-components";
import Image from "next/image";
import profileIcon from "@/assets/images/profile.svg";
import notificationIcon from "@/assets/images/notificationBellIcon.svg";
import { FaAngleDown, FaAngleUp } from "react-icons/fa6";
import { useAppSelector } from "@/lib/hooks";

export const Navbar = () => {
  const [showDropdown, setShowDropdown] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const userInfo = useAppSelector((state) => state.auth.user);

  const displayInfo = useMemo(() => {
    if (!userInfo) return { name: "Obi Enechi", email: "admin@trovotech.com" };

    // Handle different structures of userInfo
    const user = (userInfo as any).userInfo || userInfo;

    return {
      name: user?.first_name
        ? `${user.first_name} ${user.last_name || ""}`
        : "Obi Enechi",
      email: user?.email || "admin@trovotech.com",
    };
  }, [userInfo]);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setShowDropdown(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <Container>
      <LeftSection>
        <PageTitle>Dashboard</PageTitle>
      </LeftSection>
      <RightSection>
        <UserSection
          onClick={() => setShowDropdown(!showDropdown)}
          ref={dropdownRef}
        >
          <Avatar />
          <UserDetails>
            <UserName>{displayInfo.name}</UserName>
            <UserEmail>{displayInfo.email}</UserEmail>
          </UserDetails>
          {showDropdown ? <FaAngleUp /> : <FaAngleDown />}

          {showDropdown && (
            <Dropdown>
              <DropdownItem>Profile Settings</DropdownItem>
              <DropdownItem>Logout</DropdownItem>
            </Dropdown>
          )}
        </UserSection>
      </RightSection>
    </Container>
  );
};

const Container = styled.nav`
  height: 90px;
  background-color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 40px;
  position: fixed;
  top: 0;
  right: 0;
  width: calc(100% - 252px);
  z-index: 90;
  border-bottom: 1px solid #f0f0f0;
`;

const LeftSection = styled.div``;

const PageTitle = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin: 0;
`;

const RightSection = styled.div`
  display: flex;
  align-items: center;
  gap: 24px;
`;

const UserSection = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  position: relative;
  padding: 8px;
  border-radius: 12px;
  transition: background 0.2s;

  &:hover {
    background: #f8fbff;
  }
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: #007cdf;
`;

const UserDetails = styled.div`
  display: flex;
  flex-direction: column;
`;

const UserName = styled.span`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
`;

const UserEmail = styled.span`
  font-size: 12px;
  color: #00225a;
`;

const Dropdown = styled.div`
  position: absolute;
  top: 100%;
  right: 0;
  background: white;
  box-shadow: 0px 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  padding: 8px;
  width: 200px;
  margin-top: 8px;
`;

const DropdownItem = styled.div`
  padding: 12px 16px;
  font-size: 14px;
  color: #828282;
  border-radius: 8px;

  &:hover {
    background: #f8fbff;
    color: #007cdf;
  }
`;
