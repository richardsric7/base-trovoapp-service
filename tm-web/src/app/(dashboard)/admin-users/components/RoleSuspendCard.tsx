import React from "react";
import styled from "styled-components";
import userRole from "@/assets/images/oui_app-users-roless.svg";
import iconStatus from "@/assets/images/lets-icons_status.svg";
import Image from "next/image";

interface RolesProps {
  handleChangeRole: () => void;
  handleRevoke: () => void;
}

const RoleSuspendCard: React.FC<RolesProps> = ({
  handleChangeRole,
  handleRevoke,
}) => {
  return (
    <RoleContainer>
      <RoleButton onClick={handleChangeRole}>
        <Image src={userRole} alt="trovo-icon" />
        Change Role
      </RoleButton>

      <RoleButton onClick={handleRevoke}>
        <Image src={iconStatus} alt="trovo-icon" />
        Revoke Suspension
      </RoleButton>
    </RoleContainer>
  );
};

export default RoleSuspendCard;

const RoleContainer = styled.div``;

const RoleButton = styled.button`
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
