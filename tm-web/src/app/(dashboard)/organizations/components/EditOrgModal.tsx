"use client";

import {
  DropdownSelect,
  Modal,
  showErrorToast,
  showSuccessToast,
} from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import { useUpdateOrganizationInfoMutation } from "@/redux/api/organizations";
import React, { useEffect, useState } from "react";
import styled from "styled-components";

interface EditOrgProps {
  isEditModal: boolean;
  setIsEditModal: React.Dispatch<React.SetStateAction<boolean>>;
  selectedOrg: {
    id: string;
    name: string;
    email: string;
    address: string;
    type: string;
    rawType: string;
    country?: string;
    stakeholder_type: string;
    stakeholder_id?: string;
    admin: string;
    adminEmail: string;
    fee_fixed: string;
    fee_percent: string;
  };
}

const EditOrgModal: React.FC<EditOrgProps> = ({
  isEditModal,
  setIsEditModal,
  selectedOrg,
}) => {
  const [orgData, setOrgData] = useState({
    organization_name: "",
    address: "",
    country: "",
    admin_first_name: "",
    admin_last_name: "",
    fee_fixed: "",
    fee_percent: "",
    email: "",
    stakeholder_type: "",
    stakeholder_id: "",
  });

  const [updateOrganizationInfo, { isLoading }] =
    useUpdateOrganizationInfoMutation();

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

  useEffect(() => {
    if (selectedOrg && isEditModal) {
      // Use stakeholder_type directly from API (already in correct format)
      const stakeholderType = selectedOrg.stakeholder_type || "";

      setOrgData({
        organization_name: selectedOrg.name || "",
        address: selectedOrg.address || "",
        country: selectedOrg.country || "",
        admin_first_name: selectedOrg.admin?.split(" ")[0] || "",
        admin_last_name: selectedOrg.admin?.split(" ")[1] || "",
        email: selectedOrg.adminEmail || "",
        fee_fixed: selectedOrg.fee_fixed || "",
        fee_percent: selectedOrg.fee_percent || "",
        stakeholder_type: stakeholderType,
        stakeholder_id: selectedOrg.stakeholder_id || "",
      });
    }
  }, [selectedOrg, isEditModal]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!selectedOrg.stakeholder_type) {
      showErrorToast(
        "This organization cannot be edited. It doesn't have a linked stakeholder record.",
      );
      return;
    }

    const stakeholderType =
      orgData.stakeholder_type || selectedOrg.stakeholder_type;

    const requestBody = {
      organization_name: orgData.organization_name,
      address: orgData.address,
      country: orgData.country,
      admin_first_name: orgData.admin_first_name,
      admin_last_name: orgData.admin_last_name,
      fee_fixed: orgData.fee_fixed || "0",
      fee_percent: orgData.fee_percent || "0",
    };

    try {
      await updateOrganizationInfo({
        organizationId: selectedOrg.id, // path param
        stakeholderType, // query param
        body: requestBody, // body
      }).unwrap();

      showSuccessToast("Successfully updated!");
      setIsEditModal(false);
    } catch (error: any) {
      const msg =
        error?.data?.details ??
        error?.data?.message ??
        error?.message ??
        "Update failed";
      showErrorToast(msg);
    }
  };

  return (
    <>
      <Modal
        title="Edit Organization"
        isOpen={isEditModal}
        onClose={() => setIsEditModal(false)}
        closeIconPosition="right"
      >
        <ModalContent>
          <form onSubmit={handleSubmit}>
            <FormField>
              <FormLabel>Organization Name</FormLabel>
              <TextInput
                name="orgName"
                placeholder="Rand Merchant Nigeria Limited"
                value={orgData.organization_name}
                onChange={(e) =>
                  setOrgData({ ...orgData, organization_name: e.target.value })
                }
              />
            </FormField>

            <FormField>
              <FormLabel>Address</FormLabel>
              <TextInput
                name="address"
                placeholder=""
                value={orgData.address}
                onChange={(e) =>
                  setOrgData({ ...orgData, address: e.target.value })
                }
              />
            </FormField>

            <FormField>
              <FormLabel>Country</FormLabel>
              <TextInput
                name="country"
                value={orgData.country}
                onChange={(e) =>
                  setOrgData({ ...orgData, country: e.target.value })
                }
              />
            </FormField>

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
                    setOrgData({ ...orgData, fee_percent: e.target.value })
                  }
                />
                <DividerGroup>
                  <VerticalDivider />
                  <FeeLabel>%</FeeLabel>
                </DividerGroup>
              </FeeInputWrapper>
            </FormFieldGroup>

            {/* <DropdownWrapper>
              <FormLabel>Organization Type</FormLabel>
              <DropdownSelect
                options={stakeholderOptions.map((s) => s.label)}
                placeholder="Select Stakeholder Type"
                value={
                  stakeholderOptions.find(
                    (s) => s.value === orgData.stakeholder_type
                  )?.label || ""
                }
                onSelect={(label) => {
                  const found = stakeholderOptions.find(
                    (s) => s.label === label
                  );
                  setOrgData({
                    ...orgData,
                    stakeholder_type: found?.value || "",
                  });
                }}
              />
            </DropdownWrapper> */}

            <AdminContainer>
              <FormField>
                <FormLabel>Admin First Name</FormLabel>
                <TextInput
                  name="admin_first_name"
                  placeholder="Enter First Name"
                  value={orgData.admin_first_name}
                  onChange={(e) =>
                    setOrgData({ ...orgData, admin_first_name: e.target.value })
                  }
                />
              </FormField>

              <FormField>
                <FormLabel>Admin Last Name</FormLabel>
                <TextInput
                  name="admin_last_name"
                  placeholder="Enter Last Name"
                  value={orgData.admin_last_name}
                  onChange={(e) =>
                    setOrgData({ ...orgData, admin_last_name: e.target.value })
                  }
                />
              </FormField>
            </AdminContainer>

            <FormField>
              <FormLabel>Admin Email</FormLabel>
              <TextInput
                name="email"
                disabled
                type="email"
                placeholder="Enter Email"
                value={orgData.email}
                onChange={(e) =>
                  setOrgData({ ...orgData, email: e.target.value })
                }
              />
            </FormField>

            <PrimaryButton>
              {isLoading ? "Saving..." : "Save Changes"}
            </PrimaryButton>
          </form>
        </ModalContent>
      </Modal>
    </>
  );
};

export default EditOrgModal;

const ModalContent = styled.div``;

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

const FormField = styled.div`
  margin-bottom: 10px;
`;

const DropdownWrapper = styled.div`
  width: 100%;
  margin-bottom: 10px;
`;

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
  margin-bottom: 10px;
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
