"use client";

import { Modal, showErrorToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import SuccessMessage from "@/components/SuccessMessage";
import {
  useGetOrganizationDetailsDataQuery,
  useInviteMemberMutation,
} from "@/redux/api/org";
import { IInviteMember } from "@/redux/api/organizations/interface";
import { useAppSelector } from "@/lib/hooks";
import React, { useState } from "react";
import { FaPlus } from "react-icons/fa6";
import styled from "styled-components";

interface InviteMemberProps {
  organizationId?: string;
}

const InviteMemberModal: React.FC<InviteMemberProps> = ({ organizationId }) => {
  const [isInviteModalOpen, setIsInviteModalOpen] = useState(false);
  const [sendInvite, setSendInvite] = useState(false);
  const [invitedEmail, setInvitedEmail] = useState<string>("");
  const [memberData, setMemberData] = useState<
    Omit<IInviteMember, "organization_id">
  >({
    first_name: "",
    last_name: "",
    email: "",
  });

  const { data: orgDetails } = useGetOrganizationDetailsDataQuery();
  const trovoAdminUser = useAppSelector((state) => state.auth.user);
  const [inviteMember, { isLoading }] = useInviteMemberMutation();

  const handleSendInvite = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!memberData.email || !memberData.first_name || !memberData.last_name) {
      showErrorToast("please fill all details!");
      return;
    }

    try {
      const payload: IInviteMember = {
        ...memberData,
      };

      // Determine organization_id
      // 1. Use passed prop if available (Trovo Admin flow usually)
      // 2. Fallback to fetched org details (Organization Admin or Trovo Admin viewing an org)
      const targetOrgId = organizationId || orgDetails?.data?.organization?.id;

      if (targetOrgId && trovoAdminUser) {
        payload.organization_id = targetOrgId;
      }

      await inviteMember(payload).unwrap();
      setInvitedEmail(memberData.email);
      setSendInvite(true);
      setIsInviteModalOpen(false);
      setMemberData({
        first_name: "",
        last_name: "",
        email: "",
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
        <FaPlus />
        Invite Member
      </InviteButton>
      {isInviteModalOpen && (
        <Modal
          style={{ width: "40%" }}
          title="Invite Member"
          isOpen={isInviteModalOpen}
          closeIconPosition="left"
          onClose={() => setIsInviteModalOpen(false)}
        >
          <OrgInviteContainer>
            <form onSubmit={handleSendInvite}>
              <InviteFormContainer>
                <AdminContainer>
                  <FormField>
                    <FormLabel> First Name</FormLabel>
                    <TextInput
                      name="memberFirstName"
                      placeholder="Enter member’s first name"
                      value={memberData.first_name}
                      onChange={(e) =>
                        setMemberData({
                          ...memberData,
                          first_name: e.target.value,
                        })
                      }
                    />
                  </FormField>

                  <FormField>
                    <FormLabel>Last Name</FormLabel>
                    <TextInput
                      name="memberLastName"
                      placeholder="Enter member’s last name"
                      value={memberData.last_name}
                      onChange={(e) =>
                        setMemberData({
                          ...memberData,
                          last_name: e.target.value,
                        })
                      }
                    />
                  </FormField>
                </AdminContainer>

                <FormField>
                  <FormLabel> Email Address</FormLabel>
                  <TextInput
                    name="memberEmail"
                    type="email"
                    placeholder="Enter Email"
                    value={memberData.email}
                    onChange={(e) =>
                      setMemberData({
                        ...memberData,
                        email: e.target.value,
                      })
                    }
                  />
                </FormField>

                {/* <DropdownWrapper>
                  <FormLabel>Role</FormLabel>
                  <DropdownSelect
                    options={memberRoles}
                    placeholder="Select role"
                    labelText=""
                    labelColor="#828282"
                    placeholderColor="#00225A"
                    value={memberData.role}
                    onSelect={(item) =>
                      setMemberData({ ...memberData, role: item })
                    }
                  />
                </DropdownWrapper> */}

                <PrimaryButton>Send Invite</PrimaryButton>
              </InviteFormContainer>
            </form>
          </OrgInviteContainer>
        </Modal>
      )}
      {sendInvite && (
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

export default InviteMemberModal;

const InviteButton = styled.button`
  background-color: #007cdf;
  height: 48px;
  width: 200px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  border: none;
  color: #ffffff;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 2px;
  justify-content: center;
  font-family: inherit;
`;

const OrgInviteContainer = styled.div``;

const InviteFormContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: 10px 0;
`;

const DropdownWrapper = styled.div`
  width: 100%;
`;

const TextInput = styled.input`
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  width: 100%;
  outline: none;
  font-family: inherit;
`;

const FormLabel = styled.label`
  font-weight: 500;
  font-size: 16px;
  color: #00225a;
  padding: 4px 0;
  line-height: 24px;
`;

const FormField = styled.div``;

const AdminContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
`;
