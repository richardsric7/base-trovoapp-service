"use client";
import React from "react";
import styled from "styled-components";
import { useRouter } from "next/navigation";

import { useFormik } from "formik";
import * as yup from "yup";
import { Input } from "antd";
import { showErrorToast, showSuccessToast } from "@/components";
import { EyeInvisibleOutlined, EyeTwoTone } from "@ant-design/icons";
import { setCookie } from "cookies-next";
import {
  IMemberLoginPayload,
  useLoginMemberMutation,
} from "@/redux/api/organizations";
import { CookieType } from "@/app/CookieType";
import { useDispatch } from "react-redux";
import { setOrgToken } from "@/redux/slices/orgSlice";
import { TOKEN_ORG } from "@/redux";
import { useLazyGetOrganizationDetailsDataQuery } from "@/redux/api/org/api";
import { getOrganizationDashboardRoute } from "@/utils/organizationRouting";

export default function MemberLogin() {
  const loginValidation = yup.object({
    email: yup.string().email("Invalid email").required("*required"),
    password: yup.string().required("*required"),
  });
  const dispatch = useDispatch();

  const router = useRouter();
  const [memberLogin, { isLoading }] = useLoginMemberMutation();
  const [getOrganizationDetails, { isFetching: isLoadingDetails }] =
    useLazyGetOrganizationDetailsDataQuery();

  const handleLogin = (values: IMemberLoginPayload) => {
    memberLogin(values)
      .unwrap()
      .then((result) => {
        showSuccessToast("Login successfull!");

        const token = result.data.token;
        // Store token in cookie (optional)
        setCookie(TOKEN_ORG, token, {
          maxAge: 60 * 60 * 24 * 7, // 7 days cookie expiry
        });

        // Dispatch token and user info to Redux store
        dispatch(setOrgToken({ accessToken: token, refreshToken: null }));

        // dispatch(setOrgToken(token));
        // dispatch(setOrgUser(userInfo));

        return getOrganizationDetails().unwrap();
      })
      .then((details) => {
        const stakeholderType =
          details.data.organization.stakeholder_type ||
          details.data.organization.type;

        router.replace(getOrganizationDashboardRoute(stakeholderType));
      })
      .catch((err) => {
        showErrorToast(
          err?.data?.message ||
            err?.message ||
            err?.data?.error ||
            err?.error ||
            "Login failed, please check your credentials.",
        );
      });
  };

  const { values, errors, touched, handleChange, handleSubmit, handleBlur } =
    useFormik({
      initialValues: { email: "", password: "" },
      validationSchema: loginValidation,
      onSubmit: handleLogin,
    });

  return (
    <Container>
      <ContentContainer onSubmit={handleSubmit}>
        <Title>Hi There, Welcome!</Title>
        <InputLabel>Email</InputLabel>
        <Input
          type="email"
          name="email"
          onChange={handleChange}
          onBlur={handleBlur}
          value={values.email}
          placeholder="Enter your email"
          size="large"
        />
        {touched.email && errors.email && <ErrorText>{errors.email}</ErrorText>}

        <InputLabel>Password</InputLabel>
        <Input.Password
          name="password"
          onChange={handleChange}
          onBlur={handleBlur}
          value={values.password}
          placeholder="Enter your password"
          iconRender={(visible) =>
            visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />
          }
          size="large"
        />
        {touched.password && errors.password && (
          <ErrorText>{errors.password}</ErrorText>
        )}

        <ForgetPassword>Forget Password</ForgetPassword>
        <LoginButton
          type="submit"
          disabled={isLoading || isLoadingDetails}
        >
          {isLoading || isLoadingDetails ? "Logging in..." : "Login"}
        </LoginButton>
      </ContentContainer>
    </Container>
  );
}

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  margin: auto;
`;

const ContentContainer = styled.form`
  display: flex;
  flex-direction: column;
  width: 440px;
`;

const Title = styled.h2`
  color: #191919;
  font-weight: 600;
  font-size: 32px;
  line-height: 39px;
  margin-bottom: 30px;
`;

const InputLabel = styled.label`
  display: block;
  color: #00225a;
  margin-bottom: 12px;
  margin-top: 20px;
`;

const InputField = styled.input`
  border: 1px solid #bdbdbd;
  width: 100%;
  padding: 13px 16px;
  border-radius: 12px;
  font-size: 18px;
  line-height: 21px;
  font-family: inherit;
`;

const ErrorText = styled.div`
  color: red;
  font-size: 13px;
  margin-top: 3px;
`;

const LoginButton = styled.button`
  display: block;
  background-color: #007cdf;
  border: none;
  padding: 11px;
  color: #fff;
  margin: 10px 0 14px 0;
  width: 100%;
  border-radius: 10px;
  cursor: pointer;
  font-family: inherit;
`;

const ForgetPassword = styled.button`
  display: block;
  background-color: transparent;
  border: none;
  // display: flex;
  align-self: flex-end;
  padding: 11px;
  color: #007cdf;
  margin-top: 0px;
  // width: 100%;
  cursor: pointer;
  font-family: inherit;
`;
