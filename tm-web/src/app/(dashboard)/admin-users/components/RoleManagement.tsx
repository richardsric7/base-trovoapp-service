import React from "react";
import userRole from "@/assets/images/oui_app-users-roless.svg";
import changeIcon from "@/assets/images/lets-icons_statuss.svg";
import Image from "next/image";
import styled from "styled-components";

interface RolesProps {
  handleSuspend: () => void;
  handleChangeRole: () => void;
}

const RoleManagement: React.FC<RolesProps> = ({
  handleSuspend,
  handleChangeRole,
}) => {
  return (
    <RoleContainer>
      <RoleChangeButton onClick={handleChangeRole}>
        <Image src={userRole} alt="trovo-icon" />
        Change Role
      </RoleChangeButton>
      <RoleSuspendButton onClick={handleSuspend}>
        {" "}
        <Image src={changeIcon} alt="trovo-icon" />
        Suspend Admin
      </RoleSuspendButton>
    </RoleContainer>
  );
};

export default RoleManagement;

const RoleContainer = styled.div``;

const RoleChangeButton = styled.button`
  color: #00225a;
  font-size: 14px;
  font-weight: 400;
  line-height: 28px;
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 10px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-family: inherit;
`;

const RoleSuspendButton = styled.button`
  color: #be3800;
  font-size: 14px;
  font-weight: 400;
  line-height: 28px;
  border: none;
  background: transparent;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-family: inherit;
`;
