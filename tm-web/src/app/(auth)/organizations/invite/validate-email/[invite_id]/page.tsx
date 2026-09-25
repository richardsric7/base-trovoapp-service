"use client";

import React, { useEffect, useState } from "react";
import styled from "styled-components";
import { useRouter, useParams } from "next/navigation";
import { EyeInvisibleOutlined, EyeTwoTone } from "@ant-design/icons";
import { Input } from "antd";
import { Formik, Form, Field, ErrorMessage } from "formik";
import * as Yup from "yup";
import type { FieldInputProps } from "formik";

import {
  useResendTeamMemberOtpMutation,
  useSetupRootUserPasswordMutation,
  useValidateRootUserInviteMutation,
} from "@/redux/api/organizations";
import { showErrorToast, showSuccessToast } from "@/components";
import InviteExpiredPage from "@/app/(auth)/invite-expired/page";

// Validation schema
const SetPasswordSchema = Yup.object().shape({
  password: Yup.string()
    .min(8, "Password must be at least 8 characters")
    .max(32, "Password cannot exceed 32 characters")
    .matches(/[A-Z]/, "At least one uppercase letter required")
    .matches(/[a-z]/, "At least one lowercase letter required")
    .matches(/[0-9]/, "At least one number required")
    .matches(
      /[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]+/,
      "At least one special character required"
    )
    .required("Password is required"),
});

const InviteAcceptScreen = () => {
  const router = useRouter();
  const params = useParams();
  const inviteId = params.invite_id as string;
  const [validateInvite, { isLoading: isValidating }] =
    useValidateRootUserInviteMutation();

  const [resendOtp, { isLoading: isResending }] =
    useResendTeamMemberOtpMutation();

  // State for the flow
  const [inviteValidated, setInviteValidated] = useState(false);
  const [inviteExpired, setInviteExpired] = useState(false);
  const [inviteInfo, setInviteInfo] = useState({
    org: "",
    email: "",
    role: "",
  });
  const [otp, setOtp] = useState("");
  const [setupRootUserPassword, { isLoading }] =
    useSetupRootUserPasswordMutation();

  // Only validate ONCE on mount
  useEffect(() => {
    if (!inviteId) return;
    let cancelled = false;

    (async () => {
      try {
        const data = await validateInvite({ inviteId }).unwrap();
        if (cancelled) return;
        setInviteValidated(true);
        setInviteInfo({
          org: data.organization_name,
          email: data.admin_email,
          role: data.role,
        });
        sessionStorage.setItem("inviteOrgName", data.organization_name);
        showSuccessToast("Invite validated!");
      } catch (error: any) {
        if (cancelled) return;
        setInviteExpired(true);
        showErrorToast(
          error?.data?.message ??
            error?.message ??
            error?.error ??
            "Invite link is expired or invalid."
        );
      }
    })();

    // Cleanup in case component unmounts
    return () => {
      cancelled = true;
    };
  }, [inviteId]);

  // Don't re-validate! Only show the OTP UI after first successful validation.

  if (inviteExpired) {
    return <InviteExpiredPage />;
  }

  if (isValidating && !inviteValidated) {
    return (
      <Container>
        <Title>Validating your invite...</Title>
      </Container>
    );
  }

  function maskEmail(email: any) {
    if (!email) return "";
    const [name, domain] = email.split("@");
    if (name.length <= 2) return email;
    return (
      name[0] +
      "*".repeat(name.length - 2) +
      name[name.length - 1] +
      "@" +
      domain
    );
  }

  function formatRole(role: any) {
    if (!role) return "";
    // Split by underscores, lowercase all, capitalize first letter
    return role
      .toLowerCase()
      .split("_")
      .map((word: any) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(" ");
  }

  // Handle resend OTP
  const handleResend = async () => {
    try {
      await resendOtp({ email: inviteInfo.email, token: inviteId }).unwrap();
      showSuccessToast("OTP resent!");
    } catch (error: any) {
      showErrorToast(
        error?.data?.message ??
          error?.message ??
          error?.error ??
          "Failed to resend code."
      );
    }
  };

  if (inviteValidated) {
    return (
      <Container>
        <ContentContainer>
          <Title>Hi There, Welcome!</Title>
          <InfoText>
            You&apos;ve been invited to join{" "}
            <StyledSpan>{inviteInfo.org}</StyledSpan> on Trovo Manager as{" "}
            <StyledSpan>{formatRole(inviteInfo.role)}</StyledSpan>.<br />
            Check your email{" "}
            <StyledSpan>{maskEmail(inviteInfo.email)}</StyledSpan> for a 6-digit
            code to verify your invite.
          </InfoText>
          {/* OTP input */}
          <Input.OTP
            value={otp}
            onChange={(val) => setOtp(val)}
            formatter={(str) => str.toUpperCase()}
          />

          <Formik
            initialValues={{ password: "" }}
            validationSchema={SetPasswordSchema}
            onSubmit={async (values, { setSubmitting }) => {
              //  API logic
              try {
                await setupRootUserPassword({
                  inviteId,
                  otp,
                  password: values.password,
                }).unwrap();

                showSuccessToast(
                  "Password set successfully. You can now log in."
                );

                router.push("/organizations/login");
              } catch (error: any) {
                showErrorToast(
                  error?.data?.message ||
                    error?.data?.error ||
                    "Failed to set password. Please try again."
                );
              } finally {
                setSubmitting(false);
              }
            }}
          >
            {({ isSubmitting, isValid, touched }) => (
              <Form>
                <FormGroup>
                  <Label>Password</Label>
                  <Field name="password">
                    {({ field }: { field: FieldInputProps<string> }) => (
                      <Input.Password
                        {...field}
                        placeholder="Create password"
                        iconRender={(visible) =>
                          visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />
                        }
                        size="large"
                      />
                    )}
                  </Field>
                  <ErrorMessage name="password" component={ErrorText} />
                </FormGroup>
                <LoginButton
                  disabled={isSubmitting || !isValid || !touched.password}
                >
                  Finish & Login
                </LoginButton>
              </Form>
            )}
          </Formik>
          <InfoText2>
            Didn&apos;t receive a code?{" "}
            <ResendButton onClick={handleResend}>Resend</ResendButton>
          </InfoText2>
        </ContentContainer>
      </Container>
    );
  }

  // Default: waiting for validation
  return (
    <Container>
      <InfoText>Organization info will appear here after validation.</InfoText>
    </Container>
  );
};

export default InviteAcceptScreen;
const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  margin: auto;
`;
const ContentContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

const Title = styled.h2`
  color: #191919;
  font-weight: 600;
  font-size: 32px;
  line-height: 39px;
`;

const InfoText = styled.p`
  font-size: 16px;
  font-weight: 400;
  line-height: 160%;
  color: #191919;
  width: 89%;
  flex-wrap: non-wrap;
`;

const InfoText2 = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 160%;
  color: #191919;
`;

const StyledSpan = styled.span`
  font-size: 16px;
  font-weight: 700;
  line-height: 17.07px;
  color: #191919;
`;

const ResendButton = styled.button`
  line-height: 21px;
  font-family: inherit;
  font-size: 14px;
  background: none;
  border: 0;
  color: #007cdf;
  cursor: pointer;
`;
const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  font-weight: 500;
  font-size: 16px;

  color: #00225a;
`;

const LoginButton = styled.button`
  display: block;
  background-color: #007cdf;
  font-size: 16px;
  border: none;
  padding: 11px;
  color: #fff;
  margin-top: 20px;
  width: 100%;
  border-radius: 10px;
  cursor: pointer;
  font-family: inherit;
`;

const ErrorText = styled.div`
  color: red;
  font-size: 13px;
  margin-top: 2px;
`;
