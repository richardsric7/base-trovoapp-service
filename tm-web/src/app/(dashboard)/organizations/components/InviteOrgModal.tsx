"use client";

import { DropdownSelect, Modal, showErrorToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import SuccessMessage from "@/components/SuccessMessage";
import { useInviteOragnizationMutation } from "@/redux/api/organizations";

import React, { useState } from "react";
import { FaPlus } from "react-icons/fa6";
import styled from "styled-components";

const InviteOrgModal = () => {
  const [isInviteModalOpen, setIsInviteModalOpen] = useState(false);
  const [sendInvite, setSendInvite] = useState(false);
  const [invitedEmail, setInvitedEmail] = useState<string>("");
  const [orgData, setOrgData] = useState({
    admin_email: "",
    admin_first_name: "",
    admin_last_name: "",
    organization_name: "",
    address: "",
    country: "Nigeria",
    fee_fixed: "",
    fee_percent: "",
    stakeholder_type: "",
  });

  const [inviteOrganization, { isLoading, error }] =
    useInviteOragnizationMutation();
  const stakeholderOptions = [
    { label: "Asset Manager", value: "asset_manager" },
    { label: "Asset Issuing House", value: "asset_issuing_house" },
    { label: "Asset Custodian", value: "approved_asset_custodian" },
    { label: "Legal & Professionals", value: "legal_and_professionals" },
    { label: "Rating Agency", value: "rating_agency" },
    { label: "Trustees", value: "trustees" },
    { label: "Legal Adviser", value: "legal_adviser" },
    { label: "Financial Adviser", value: "financial_adviser" },
  ];

  const countryOptions = [{ label: "Nigeria", value: "NG" }];

  const handleSendInvite = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validate required fields
    if (
      !orgData.admin_email ||
      !orgData.admin_first_name ||
      !orgData.admin_last_name ||
      !orgData.organization_name ||
      !orgData.address ||
      !orgData.country ||
      !orgData.fee_fixed ||
      !orgData.fee_percent
    ) {
      showErrorToast("Please fill all required details!");
      return;
    }

    // Option 1: Send body directly (try this first)
    const bodyData = {
      organization_name: orgData.organization_name,
      address: orgData.address,
      country: orgData.country,
      admin_first_name: orgData.admin_first_name,
      admin_last_name: orgData.admin_last_name,
      admin_email: orgData.admin_email,
      fee_fixed: orgData.fee_fixed,
      fee_percent: orgData.fee_percent,
    };

    const requestPayload = {
      ...bodyData, // Spread body directly
      stakeholder_type: orgData.stakeholder_type || undefined,
    };

    try {
      const result = await inviteOrganization(requestPayload).unwrap();

      console.log(" SUCCESS! Response:", result);

      setInvitedEmail(orgData.admin_email);
      setSendInvite(true);
      setIsInviteModalOpen(false);

      // Reset form
      setOrgData({
        admin_email: "",
        admin_first_name: "",
        admin_last_name: "",
        organization_name: "",
        address: "",
        country: "NG",
        fee_fixed: "",
        fee_percent: "",
        stakeholder_type: "",
      });
    } catch (error: any) {
      const msg =
        error?.data?.message ??
        error?.message ??
        error?.error ??
        "Something went wrong. Please try again.";

      console.error("Toast message shown:", msg);
      showErrorToast(msg);
    }
  };

  return (
    <>
      <InviteButton onClick={() => setIsInviteModalOpen(true)}>
        <FaPlus />
        Invite Organization
      </InviteButton>
      {isInviteModalOpen && (
        <Modal
          style={{ width: "40%" }}
          title="Invite Organization"
          isOpen={isInviteModalOpen}
          closeIconPosition="left"
          onClose={() => setIsInviteModalOpen(false)}
        >
          <OrgInviteContainer>
            <form onSubmit={handleSendInvite}>
              <InviteFormContainer>
                <FormField>
                  <FormLabel>Organization Name</FormLabel>
                  <TextInput
                    name="orgName"
                    placeholder="Input organization's name"
                    value={orgData.organization_name}
                    onChange={(e) =>
                      setOrgData({
                        ...orgData,
                        organization_name: e.target.value,
                      })
                    }
                  />
                </FormField>

                <FormField>
                  <FormLabel>Address</FormLabel>
                  <TextInput
                    name="address"
                    placeholder="Enter address"
                    value={orgData.address}
                    onChange={(e) =>
                      setOrgData({ ...orgData, address: e.target.value })
                    }
                  />
                </FormField>

                <DropdownWrapper>
                  <FormLabel>Country</FormLabel>
                  <DropdownSelect
                    options={countryOptions.map((c) => c.label)}
                    placeholder="Select Country"
                    value={
                      countryOptions.find((c) => c.value === orgData.country)
                        ?.label || ""
                    }
                    onSelect={(label) => {
                      const found = countryOptions.find(
                        (c) => c.label === label,
                      );
                      setOrgData({
                        ...orgData,
                        country: found?.value || "",
                      });
                    }}
                  />
                </DropdownWrapper>

                <FeeContainer>
                  <FormFieldGroup>
                    <FormLabel>Fixed Fee</FormLabel>
                    <FeeInputWrapper>
                      <FeeInputField
                        placeholder="0"
                        name="fee_fixed"
                        value={orgData.fee_fixed}
                        onChange={(e) =>
                          setOrgData({ ...orgData, fee_fixed: e.target.value })
                        }
                      />
                      <DividerGroup>
                        <VerticalDivider />
                        <FeeLabel>CNGN</FeeLabel>
                      </DividerGroup>
                    </FeeInputWrapper>
                  </FormFieldGroup>

                  <FormFieldGroup>
                    <FormLabel>Fee Percentage</FormLabel>
                    <FeeInputWrapper>
                      <FeeInputField
                        placeholder="0"
                        name="fee_percent"
                        value={orgData.fee_percent}
                        onChange={(e) =>
                          setOrgData({
                            ...orgData,
                            fee_percent: e.target.value,
                          })
                        }
                      />
                      <DividerGroup>
                        <VerticalDivider />
                        <FeeLabel>%</FeeLabel>
                      </DividerGroup>
                    </FeeInputWrapper>
                  </FormFieldGroup>
                </FeeContainer>
                {/* <FormField>
                  <FormLabel>Percentage Fee</FormLabel>
                  <TextInput
                    name="fee_percent"
                    placeholder="e.g. 0.5"
                    value={orgData.fee_percent}
                    onChange={(e) =>
                      setOrgData({ ...orgData, fee_percent: e.target.value })
                    }
                  />
                </FormField> */}

                <DropdownWrapper>
                  <FormLabel>Organization Type</FormLabel>
                  <DropdownSelect
                    options={stakeholderOptions.map((s) => s.label)}
                    placeholder="Select Stakeholder Type"
                    value={
                      stakeholderOptions.find(
                        (s) => s.value === orgData.stakeholder_type,
                      )?.label || ""
                    }
                    onSelect={(label) => {
                      const found = stakeholderOptions.find(
                        (s) => s.label === label,
                      );
                      setOrgData({
                        ...orgData,
                        stakeholder_type: found?.value || "",
                      });
                    }}
                  />
                </DropdownWrapper>

                <AdminContainer>
                  <FormField>
                    <FormLabel>Admin First Name</FormLabel>
                    <TextInput
                      name="adminFirstName"
                      placeholder="Enter First Name"
                      value={orgData.admin_first_name}
                      onChange={(e) =>
                        setOrgData({
                          ...orgData,
                          admin_first_name: e.target.value,
                        })
                      }
                    />
                  </FormField>

                  <FormField>
                    <FormLabel>Admin Last Name</FormLabel>
                    <TextInput
                      name="adminLastName"
                      placeholder="Enter Last Name"
                      value={orgData.admin_last_name}
                      onChange={(e) =>
                        setOrgData({
                          ...orgData,
                          admin_last_name: e.target.value,
                        })
                      }
                    />
                  </FormField>
                </AdminContainer>

                <FormField>
                  <FormLabel>Admin Email</FormLabel>
                  <TextInput
                    name="adminEmail"
                    type="email"
                    placeholder="Enter Email"
                    value={orgData.admin_email}
                    onChange={(e) =>
                      setOrgData({ ...orgData, admin_email: e.target.value })
                    }
                  />
                </FormField>

                <PrimaryButton>
                  {" "}
                  {isLoading ? "Sending..." : "Send Invite"}
                </PrimaryButton>
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

export default InviteOrgModal;

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

const FeeContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
`;
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

const FormFieldGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const FeeInputField = styled.input`
  font-size: 14px;
  width: 100%;
  outline: none;
  font-family: inherit;
  border: none;
`;

const FeeInputWrapper = styled.div`
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  display: flex;
  align-items: center;
`;

const DividerGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const VerticalDivider = styled.div`
  width: 1px;
  height: 24px;
  background-color: #e0e0e0;
`;
const FeeLabel = styled.p`
  font-weight: 500;
  font-size: 14px;
  color: #191919;
`;
