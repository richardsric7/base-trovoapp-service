"use client";

import { DropdownSelect, Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import React, { useState } from "react";
import styled from "styled-components";

interface ChangeRoleProps {
  isChangeRole: boolean;
  setIsChangeRole: React.Dispatch<React.SetStateAction<boolean>>;
  role?: string;
  // name?: string;
  // email?: string;
}

const ChangeMemberRole: React.FC<ChangeRoleProps> = ({
  setIsChangeRole,
  isChangeRole,
  role,
  // name?: string;
  // email?: string;
}) => {
  const [memberData, setMemberData] = useState({
    memberRole: role ?? "",
  });
  const memberRoles = ["Admin", "viewer", "Auditor", "Payout"];

  const handleUpdateRole = () => {
    setIsChangeRole(false);
  };
  return (
    <Modal
      title="Change Role"
      isOpen={isChangeRole}
      onClose={() => setIsChangeRole(false)}
      closeIconPosition="left"
    >
      <ModalContent>
        <UserDetailsSection>
          <BoldText>
            H{/* {record.name ? record.admin.charAt(0).toUpperCase() : "?"} */}
          </BoldText>

          <div>
            <UserName>Hannah Doe</UserName>
            <UserEmail>john@gmail.com</UserEmail>
          </div>
        </UserDetailsSection>
        <DropdownWrapper>
          <FormLabel>Role</FormLabel>
          <DropdownSelect
            options={memberRoles}
            placeholder="Select role"
            labelText=""
            labelColor="#828282"
            placeholderColor="#00225A"
            value={memberData.memberRole}
            onSelect={(item) =>
              setMemberData({ ...memberData, memberRole: item })
            }
          />
        </DropdownWrapper>
        <PrimaryButton>Update Role</PrimaryButton>
      </ModalContent>
    </Modal>
  );
};

export default ChangeMemberRole;

const ModalContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding-top: 8px;
`;

const UserDetailsSection = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

const BoldText = styled.div`
  background-color: #acd1ef;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #191919;
  font-weight: 600;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;
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
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const DropdownWrapper = styled.div`
  width: 100%;
`;
const FormLabel = styled.label`
  font-weight: 500;
  font-size: 16px;
  color: #00225a;
  padding: 4px 0;
  line-height: 24px;
`;
