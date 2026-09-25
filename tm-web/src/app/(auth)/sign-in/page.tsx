"use client";
import React, { useEffect } from "react";
import styled from "styled-components";
import { useRouter } from "next/navigation";
import { ILogin, useLoginMutation } from "@/redux/api/auth";
import { useFormik } from "formik";
import * as yup from "yup";
import { showErrorToast, showSuccessToast } from "@/components";
import { setCookie } from "cookies-next";
import { getStorage, clearAdminSession } from "@/utils";
import { TOKEN } from "@/redux/constants";
import { CookieType } from "../../CookieType";

export default function Auth() {
  // Belt-and-braces: if a user lands here with an expired admin token still in
  // storage, wipe it so nothing on this page (or the next load) trips over it.
  useEffect(() => {
    const token = getStorage<{ expiresAt?: number }>(TOKEN);
    if (token?.expiresAt && token.expiresAt <= Date.now()) {
      clearAdminSession();
    }
  }, []);

  const loginValidation = yup.object({
    username: yup.string().required("*required"),
  });

  const router = useRouter();
  const [userLogin, { isLoading }] = useLoginMutation();
  const handleLogin = (values: ILogin) => {
    const loginType = values.username === "user" ? "admin" : "user"; // Determine login type
    if (loginType !== "user") {
      console.error("Access denied. Only admins are allowed to log in.");
      showErrorToast("Only users are allowed to log in.");
      return;
    }

    // const payload = { username: values.username + "@trovo", loginType: "user" };
    const payload = { username: values.username + "@trovo", loginType };
    userLogin(payload)
      .unwrap()
      .then((result) => {
        showSuccessToast("username entered successfuly");
        router.push("/qr-scan");
        setCookie(
          CookieType.LoginResults,
          JSON.stringify({
            ...result,
            username: values.username,
          })
        );
      })
      .catch((err) => {
        showErrorToast(err?.data?.message ?? err?.message ?? err?.error ?? "Login failed");
      });
  };

  const { values, handleChange, handleSubmit } = useFormik({
    initialValues: { username: "", loginType: "" },
    validationSchema: loginValidation,
    onSubmit: handleLogin,
  });
  return (
    <Container>
      <ContentContainer onSubmit={handleSubmit}>
        <Title>Hi, Welcome Back!</Title>
        <InputLabel>Enter your Trovo App username</InputLabel>
        <InputField
          onChange={handleChange("username")}
          value={values.username}
          placeholder="trovo username"
        />
        <InfoText>Input the username without typing the “@” symbols</InfoText>
        <LoginButton type="submit">
          {isLoading ? "loading..." : "Login"}
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
const ContentContainer = styled.form``;
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
`;
const InputField = styled.input`
  border: 1px solid #bdbdbd;
  width: 440px;
  padding: 13px 16px;
  border-radius: 12px;
  font-size: 18px;
  line-height: 21px;
  font-family: inherit;
`;
const InfoText = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 17.07px;
  color: #00225a;
  margin-top: 5px;
`;
const LoginButton = styled.button`
  display: block;
  background-color: #007cdf;
  border: none;
  padding: 11px;
  color: #fff;
  margin-top: 40px;
  width: 100%;
  border-radius: 10px;
  cursor: pointer;
  font-family: inherit;
`;
