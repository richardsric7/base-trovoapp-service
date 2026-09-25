"use client";

import React, { useState } from "react";
import styled from "styled-components";
import { Modal } from "@/components/CustomModal";
import SuccessMessage from "@/components/SuccessMessage";
import { useInviteAdminMutation } from "@/redux/api/admin/admin";
import { IInviteAdmin } from "@/redux/api/admin/interface";
import { DropdownSelect, showErrorToast } from "@/components";

const InviteAdminUser = () => {
  const [isInviteModalOpen, setIsInviteModalOpen] = useState<boolean>(false);
  const [sendInvite, setSendInvite] = useState<boolean>(false);
  const [invitedEmail, setInvitedEmail] = useState<string | null>(null);

  const [adminData, setAdminData] = useState<IInviteAdmin>({
    admin_role: "",
    email_or_username: "",
  });
  const [inviteAdmin, { isLoading, error, data }] = useInviteAdminMutation();

  const roleOptions = ["SUPER_ADMIN", "EDIT_LEVEL_ADMIN", "VIEW_ONLY_ADMIN"];

  const handleSendInvite = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const result = await inviteAdmin(adminData).unwrap();
      setInvitedEmail(adminData.email_or_username); // Store the invited email
      setSendInvite(true);
      setIsInviteModalOpen(false);
      setAdminData({
        admin_role: "",
        email_or_username: "",
      });
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
    <>
      <InviteButton onClick={() => setIsInviteModalOpen(true)}>
        Invite Users
      </InviteButton>
      {isInviteModalOpen && !sendInvite && (
        <Modal
          style={{ width: "40%" }}
          title="Invite Users"
          isOpen={isInviteModalOpen}
          onClose={() => setIsInviteModalOpen(false)}
        >
          <AdminInviteContainer>
            <form onSubmit={handleSendInvite}>
              <InviteFormContainer>
                <InputContainer>
                  <input
                    placeholder="Enter email or username"
                    type="text"
                    value={adminData.email_or_username}
                    onChange={(e) =>
                      setAdminData({
                        ...adminData,
                        email_or_username: e.target.value,
                      })
                    }
                  />
                  <DropdownWrapper>
                    <DropdownSelect
                      options={roleOptions}
                      placeholder="Select Role"
                      labelText=""
                      labelColor="#828282"
                      placeholderColor="#00225A"
                      value={adminData.admin_role}
                      onSelect={(item) =>
                        setAdminData({ ...adminData, admin_role: item })
                      }
                    />
                  </DropdownWrapper>
                </InputContainer>
                <InviteButton type="submit" disabled={isLoading}>
                  {isLoading ? "Inviting..." : " Send Invite"}
                </InviteButton>
              </InviteFormContainer>
            </form>
          </AdminInviteContainer>
        </Modal>
      )}
      {sendInvite && invitedEmail && (
        <SuccessMessage
          isOpen={sendInvite}
          setIsOpen={setSendInvite}
          heading="Success"
          message="You have successfully sent an invite to"
          email={invitedEmail}
        />
      )}
    </>
  );
};
export default InviteAdminUser;
const InviteButton = styled.button`
  background-color: #007cdf;
  height: 48px;
  width: 130px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  border: none;
  color: #ffffff;
  cursor: pointer;
`;

const AdminInviteContainer = styled.div`
  background-color: #ffffff;
  border-radius: 24px;
  padding: 6px;
  width: 100%;
`;

const InviteFormContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
  // gap: 4px;
`;

const InputContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  // padding-top: 10px;
  padding: 4px 2px 2px 12px;
  gap: 0;
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  width: 100%;
  position: relative;
  overflow: visible;

  input {
    border: none;
    outline: none;
    font-family: inherit;
    flex-grow: 1;
  }

  placeholder {
    font-size: 14px;
    font-weight: 500;
    line-height: 24px;
    padding-right: 4px;
    letter-spacing: 0.1px;
  }
`;

const DropdownWrapper = styled.div`
  width: 50%;
  position: relative;
  overflow: visible;
  z-index: 1000;
`;
