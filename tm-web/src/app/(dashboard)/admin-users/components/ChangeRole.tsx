"use client";

import React, { useState } from "react";
import styled from "styled-components";
import {
  DropdownSelect,
  Modal,
  showErrorToast,
  showSuccessToast,
} from "@/components";
import { useUpdateAdminRoleMutation } from "@/redux/api/admin/admin";
import PrimaryButton from "@/components/PrimaryButton";

interface ChangeAdminProps {
  changeAdminRole: boolean;
  setChangeAdminRole: React.Dispatch<React.SetStateAction<boolean>>;
  adminEmail: string;
}

const ChangeRole: React.FC<ChangeAdminProps> = ({
  changeAdminRole,
  setChangeAdminRole,
  adminEmail,
}) => {
  const [updateAdminRole, { isLoading }] = useUpdateAdminRoleMutation();

  const [roleValue, setRoleValue] = useState<string>("");
  const [adminData, setAdminData] = useState({
    email: adminEmail,
    admin_role: "",
  });

  const handleChangeAdminRole = async () => {
    try {
      await updateAdminRole({
        email: adminData.email,
        admin_role: roleValue,
      }).unwrap();
      setChangeAdminRole(false); 

      showSuccessToast(`
        "Role updated successfully for:",
        ${(adminData.email, roleValue)}`);
    } catch (error: any) {
      const msg =
        error?.data?.message ??
        error?.message ??
        error?.error ??
        "Something went wrong. Please try again.";
      showErrorToast(msg);
    }
  };

  return (
    <Modal
      title="Change Role"
      isOpen={changeAdminRole}
      onClose={() => setChangeAdminRole(false)}
    >
      <ChangeCard>
        <Label>Email</Label>
        <StyledInput
          placeholder="Enter email"
          type="email"
          value={adminData.email}
          onChange={(e) =>
            setAdminData({ ...adminData, email: e.target.value })
          }
        />
        <FormGroup>
          <Label>Role</Label>
          <DropdownSelect
            options={["SUPER_ADMIN", "EDIT_LEVEL_ADMIN", "VIEW_ONLY_ADMIN"]}
            placeholder="Select"
            labelText=""
            value={roleValue}
            onSelect={(item) => {
              setRoleValue(item);
              setAdminData({ ...adminData, admin_role: item });
            }}
          />
        </FormGroup>
      </ChangeCard>

      <PrimaryButton onClick={handleChangeAdminRole}>
        {" "}
        {isLoading ? "Updating Role..." : "Update Role"}
      </PrimaryButton>
    </Modal>
  );
};

export default ChangeRole;

const FormGroup = styled.div`
  margin: 12px 0;
`;

const ChangeCard = styled.div`
  margin-bottom: 20px;
`;

const StyledInput = styled.input`
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  padding: 2px 6px;
  width: 100%;
  outline: none;
  height: 48px;
  font-family: inherit;

  &:disabled {
    background-color: #f5f5f5;
  }
`;
const Label = styled.div`
  margin-bottom: 8px;
  font-style: normal;
  font-weight: 500;
  font-size: 14px;
  line-height: 24px;
  color: #00225a;
  display: block;
`;
