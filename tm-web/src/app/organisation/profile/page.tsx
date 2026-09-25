"use client";

import { FormEvent, useEffect, useState } from "react";
import styled from "styled-components";
import { Switch } from "antd";
import {
  useGetStakeholderProfileQuery,
  useUpdateStakeholderNotificationPreferencesMutation,
  useUpdateStakeholderProfileMutation,
} from "@/redux/api/sharedstakeholders";
import { showErrorToast, showSuccessToast } from "@/components";

const errorMessage = (error: unknown) => {
  const apiError = error as {
    data?: { error?: string; message?: string };
    message?: string;
  };
  return apiError?.data?.message || apiError?.message || apiError?.data?.error || "Something went wrong";
};

const formatRole = (role: string) =>
  role
    .replace(/_/g, " ")
    .toLowerCase()
    .replace(/\b\w/g, (character) => character.toUpperCase());

export default function StakeholderProfilePage() {
  const { data, isLoading, isError, refetch } = useGetStakeholderProfileQuery();
  const [updateProfile, { isLoading: isSavingProfile }] =
    useUpdateStakeholderProfileMutation();
  const [updatePreferences, { isLoading: isSavingPreferences }] =
    useUpdateStakeholderNotificationPreferencesMutation();
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [emailEnabled, setEmailEnabled] = useState(true);
  const [inAppEnabled, setInAppEnabled] = useState(true);

  const profile = data?.data;

  useEffect(() => {
    if (profile?.member) {
      setFirstName(profile.member.first_name || "");
      setLastName(profile.member.last_name || "");
    }
  }, [profile?.member]);

  const saveProfile = async (event: FormEvent) => {
    event.preventDefault();
    if (!firstName.trim() || !lastName.trim()) {
      showErrorToast("First name and last name are required");
      return;
    }
    try {
      await updateProfile({
        first_name: firstName.trim(),
        last_name: lastName.trim(),
      }).unwrap();
      showSuccessToast("Profile updated");
    } catch (error) {
      showErrorToast(errorMessage(error));
    }
  };

  const updateNotificationPreference = async (
    preference: "email" | "inApp",
    checked: boolean,
  ) => {
    const previousEmailEnabled = emailEnabled;
    const previousInAppEnabled = inAppEnabled;
    const nextEmailEnabled = preference === "email" ? checked : emailEnabled;
    const nextInAppEnabled = preference === "inApp" ? checked : inAppEnabled;

    setEmailEnabled(nextEmailEnabled);
    setInAppEnabled(nextInAppEnabled);

    try {
      const response = await updatePreferences({
        email_enabled: nextEmailEnabled,
        in_app_enabled: nextInAppEnabled,
      }).unwrap();
      setEmailEnabled(response.data.email_enabled);
      setInAppEnabled(response.data.in_app_enabled);
      showSuccessToast("Notification preferences updated");
    } catch (error) {
      setEmailEnabled(previousEmailEnabled);
      setInAppEnabled(previousInAppEnabled);
      showErrorToast(errorMessage(error));
    }
  };

  if (isLoading) return <StateMessage>Loading profile...</StateMessage>;
  if (isError || !profile) {
    return (
      <StateMessage>
        Unable to load your profile.{" "}
        <button onClick={() => refetch()}>Retry</button>
      </StateMessage>
    );
  }

  return (
    <Page>
      <Header>
        <h1>Profile settings</h1>
        <p>Manage your stakeholder account and notification preferences.</p>
      </Header>

      <Card>
        <SectionTitle>Personal information</SectionTitle>
        <Form onSubmit={saveProfile}>
          <Field>
            <label htmlFor="first-name">First name</label>
            <input
              id="first-name"
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
            />
          </Field>
          <Field>
            <label htmlFor="last-name">Last name</label>
            <input
              id="last-name"
              value={lastName}
              onChange={(e) => setLastName(e.target.value)}
            />
          </Field>
          <Field>
            <label htmlFor="email">Email</label>
            <input id="email" value={profile.member.email} disabled />
          </Field>
          <Field>
            <label htmlFor="role">Role</label>
            <input id="role" value={formatRole(profile.member.role)} disabled />
          </Field>
          <Button type="submit" disabled={isSavingProfile}>
            {isSavingProfile ? "Saving..." : "Save changes"}
          </Button>
        </Form>
      </Card>

      <Card id="notifications">
        <SectionTitle>Notification preferences</SectionTitle>
        <Preference>
          <div>
            <strong>Email notifications</strong>
            <span>Receive stakeholder updates by email.</span>
          </div>
          <Switch
            id="email-notifications"
            checked={emailEnabled}
            onChange={(checked) =>
              updateNotificationPreference("email", checked)
            }
            disabled={isSavingPreferences}
          />
        </Preference>
        <Preference>
          <div>
            <strong>In-app notifications</strong>
            <span>Receive updates in the Trovo dashboard.</span>
          </div>
          <Switch
            id="in-app-notifications"
            checked={inAppEnabled}
            onChange={(checked) =>
              updateNotificationPreference("inApp", checked)
            }
            disabled={isSavingPreferences}
          />
        </Preference>
      </Card>
    </Page>
  );
}

const Page = styled.div`
  display: grid;
  gap: 24px;
  padding: 24px 24px 48px;
  color: #00225a;
`;
const Header = styled.header`
  h1 {
    margin: 0 0 8px;
    font-size: 28px;
  }
  p {
    margin: 0;
    color: #667085;
  }
`;
const Card = styled.section`
  background: #fff;
  border-radius: 20px;
  padding: 28px;
  scroll-margin-top: 120px;
`;
const SectionTitle = styled.h2`
  font-size: 20px;
  margin: 0 0 24px;
`;
const Form = styled.form`
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
  @media (max-width: 800px) {
    grid-template-columns: 1fr;
  }
`;
const Field = styled.div`
  display: grid;
  gap: 8px;
  label {
    font-size: 14px;
    font-weight: 500;
  }
  input {
    border: 1px solid #d0d5dd;
    border-radius: 8px;
    padding: 12px;
    color: #00225a;
    background: #fff;
  }
  input:disabled {
    background: #f2f4f7;
    color: #667085;
  }
`;
const Button = styled.button`
  justify-self: start;
  border: 0;
  border-radius: 8px;
  padding: 12px 20px;
  background: #0066ff;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  font-family: inherit;
`;
const Preference = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding: 18px 0;
  border-top: 1px solid #eaecf0;
  strong,
  span {
    display: block;
  }
  span {
    margin-top: 4px;
    color: #667085;
    font-size: 14px;
  }
`;
const StateMessage = styled.div`
  margin: 40px 24px;
  padding: 24px;
  border-radius: 16px;
  background: #fff;
  color: #475467;
  button {
    margin-left: 8px;
  }
`;
