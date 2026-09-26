"use client";
import React, { useState } from "react";
import styled from "styled-components";
import PrimaryButton from "@/components/PrimaryButton";
import ToggleSwitch from "@/components/ToggleSwitch";
import { showErrorToast } from "@/components";
import { IServiceLink } from "@/redux/api/serviceLinks";

export type ServiceLinkFormValues = {
  ownerUsername: string;
  shortName: string;
  longName: string;
  loginPermission: boolean;
  paymentPermission: boolean;
  tokenInfoPermission: boolean;
  authorizationPermission: boolean;
  eventPermission: boolean;
  allowUserInfo: boolean;
  pushNotificationPermission: boolean;
  includePhoneNumbers: boolean;
  includeUserBalances: boolean;
  tokenizedAssetAuthorizationPermission: boolean;
  createUsersPermission: boolean;
  allowReferralForRegisteredUsers: boolean;
  verified: boolean;
  inactive: boolean;
};

const PERMISSION_FIELDS: {
  key: keyof ServiceLinkFormValues;
  label: string;
  hint: string;
}[] = [
  { key: "loginPermission", label: "Login", hint: "Initiate/approve QR-login flows" },
  { key: "paymentPermission", label: "Payment", hint: "Request payments from a user" },
  { key: "tokenInfoPermission", label: "Token info", hint: "Read tokenized-asset info" },
  { key: "authorizationPermission", label: "Authorization", hint: "2FA/authorization request+approval flows" },
  { key: "eventPermission", label: "Events", hint: "Generate and approve event links" },
  { key: "allowUserInfo", label: "User info", hint: "Fetch a user's info" },
  { key: "pushNotificationPermission", label: "Push notifications", hint: "Push notifications to a user" },
  { key: "includePhoneNumbers", label: "Include phone numbers", hint: "Include mobile numbers in user-info responses" },
  { key: "includeUserBalances", label: "Include balances", hint: "Include wallet balances in user-info responses" },
  {
    key: "tokenizedAssetAuthorizationPermission",
    label: "Tokenized-asset authorization",
    hint: "Authorize tokenized-asset transactions",
  },
  {
    key: "createUsersPermission",
    label: "Create users",
    hint: "Onboard/manage users via the trovo-api surface (KYC, minting, tokenization, payments)",
  },
  {
    key: "allowReferralForRegisteredUsers",
    label: "Auto-referral",
    hint: "Users onboarded through this link get it set as their referrer",
  },
];

const fromServiceLink = (link?: IServiceLink): ServiceLinkFormValues => ({
  ownerUsername: link?.ownerUsername ?? "",
  shortName: link?.shortName ?? "",
  longName: link?.longName ?? "",
  loginPermission: !!link?.loginPermission,
  paymentPermission: !!link?.paymentPermission,
  tokenInfoPermission: !!link?.tokenInfoPermission,
  authorizationPermission: !!link?.authorizationPermission,
  eventPermission: !!link?.eventPermission,
  allowUserInfo: !!link?.allowUserInfo,
  pushNotificationPermission: !!link?.pushNotificationPermission,
  includePhoneNumbers: !!link?.includePhoneNumbers,
  includeUserBalances: !!link?.includeUserBalances,
  tokenizedAssetAuthorizationPermission: !!link?.tokenizedAssetAuthorizationPermission,
  createUsersPermission: !!link?.createUsersPermission,
  allowReferralForRegisteredUsers: !!link?.allowReferralForRegisteredUsers,
  verified: !!link?.verified,
  inactive: !!link?.inactive,
});

const ServiceLinkForm = ({
  initialLink,
  submitLabel,
  submitting,
  onSubmit,
}: {
  initialLink?: IServiceLink;
  submitLabel: string;
  submitting: boolean;
  onSubmit: (values: ServiceLinkFormValues) => void;
}) => {
  const [values, setValues] = useState<ServiceLinkFormValues>(fromServiceLink(initialLink));
  const isEditing = !!initialLink;

  const set = <K extends keyof ServiceLinkFormValues>(key: K, value: ServiceLinkFormValues[K]) =>
    setValues((prev) => ({ ...prev, [key]: value }));

  const handleTextChange = (key: "shortName" | "longName" | "ownerUsername") => (
    e: React.ChangeEvent<HTMLInputElement>
  ) => set(key, e.target.value as never);

  const handleSubmit = () => {
    if (!values.ownerUsername.trim()) {
      showErrorToast("Owner username is required");
      return;
    }
    if (!values.shortName.trim()) {
      showErrorToast("Short name is required");
      return;
    }
    onSubmit(values);
  };

  return (
    <FormGrid>
      <Field>
        <Label>Owner username *</Label>
        <Input
          value={values.ownerUsername}
          onChange={handleTextChange("ownerUsername")}
          disabled={isEditing}
          placeholder="e.g. partnerco"
        />
        {!isEditing && (
          <FieldHint>
            Must already exist in the users table - the account&apos;s wallet
            address is copied onto this service link automatically.
          </FieldHint>
        )}
      </Field>
      <Field>
        <Label>Short name *</Label>
        <Input value={values.shortName} onChange={handleTextChange("shortName")} placeholder="e.g. partnerco-app" />
      </Field>
      <FullWidthField>
        <Label>Long name</Label>
        <Input value={values.longName} onChange={handleTextChange("longName")} placeholder="e.g. PartnerCo Wallet App" />
      </FullWidthField>

      <FullWidthField>
        <SectionLabel>Permissions</SectionLabel>
        <PermissionsGrid>
          {PERMISSION_FIELDS.map((field) => (
            <PermissionCard key={field.key}>
              <ToggleSwitch
                id={field.key}
                label={field.label}
                checked={values[field.key] as boolean}
                onChange={(v) => set(field.key, v as never)}
              />
              <PermissionHint>{field.hint}</PermissionHint>
            </PermissionCard>
          ))}
        </PermissionsGrid>
      </FullWidthField>

      <StatusRow>
        <ToggleSwitch
          id="verified"
          label="Verified"
          checked={values.verified}
          onChange={(v) => set("verified", v as never)}
        />
        <ToggleSwitch
          id="inactive-as-active"
          label="Active"
          checked={!values.inactive}
          onChange={(checked) => set("inactive", !checked as never)}
        />
      </StatusRow>

      <SubmitRow>
        <PrimaryButton onClick={handleSubmit} disabled={submitting} buttonStyle={{ width: "auto", padding: "10px 32px" }}>
          {submitting ? "Saving..." : submitLabel}
        </PrimaryButton>
      </SubmitRow>
    </FormGrid>
  );
};

export default ServiceLinkForm;

const FormGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  max-width: 900px;
`;

const Field = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
`;

const FullWidthField = styled(Field)`
  grid-column: 1 / -1;
`;

const Label = styled.label`
  font-size: 14px;
  color: #828282;
`;

const SectionLabel = styled.h3`
  font-size: 16px;
  font-weight: 600;
  color: #00225a;
  margin: 0;
`;

const FieldHint = styled.p`
  font-size: 12px;
  color: #828282;
  margin: 0;
`;

const Input = styled.input`
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  outline: none;
  font-size: 14px;

  &:disabled {
    background-color: #f5f5f5;
    color: #828282;
  }
`;

const PermissionsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-top: 12px;
`;

const PermissionCard = styled.div`
  background-color: #f2f6f9;
  border-radius: 10px;
  padding: 12px 16px;
`;

const PermissionHint = styled.p`
  font-size: 12px;
  color: #828282;
  margin: 6px 0 0 0;
`;

const StatusRow = styled.div`
  grid-column: 1 / -1;
  display: flex;
  gap: 32px;
  padding: 12px 0;
  border-top: 1px solid #e5e5ef;
  border-bottom: 1px solid #e5e5ef;
`;

const SubmitRow = styled.div`
  grid-column: 1 / -1;
  padding-top: 8px;
`;
