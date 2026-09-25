"use client";
import React, { SetStateAction, useEffect, useMemo, useState } from "react";
import styled from "styled-components";
import { useRouter } from "next/navigation";
import { deleteCookie, getCookie, hasCookie, setCookie } from "cookies-next";
import { CookieType } from "@/app/CookieType";
import { TOKEN, useLoginMutation, useVerifyLoginQuery } from "@/redux";
import { showErrorToast, showSuccessToast } from "@/components";
import { useDispatch } from "react-redux";
import { setToken, setUser } from "@/redux/slices/authSlice";
import { GrRefresh } from "react-icons/gr";
import logo from "@/assets/images/logo-qr-code.svg";
import Image from "next/image";
import { BASE_URL } from "@/config";

type ILoginRes = {
  dynamicLink: string;
  loginId: string;
  qrCode: string;
  username: string;
};

interface QrCodeContainerProps {
  showRetry: boolean;
}

const useDeviceType = () => {
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const handleResize = () => {
      setIsMobile(window.innerWidth < 768);
    };
    handleResize();
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  return isMobile;
};

export default function ScanQr() {
  const isMobile = useDeviceType();
  const router = useRouter();
  const [loginProps, setLoginProps] = useState<ILoginRes>(() => {
    if (hasCookie(CookieType.LoginResults)) {
      const loginResults = getCookie(CookieType.LoginResults)!;
      return JSON.parse(loginResults) as ILoginRes;
    }
    return {} as ILoginRes;
  });

  const [refreshQrCode, setRefreshQrCode] = useState(false);
  const [userLogin, { isLoading }] = useLoginMutation();

  const handleLogin = () => {
    const payload = {
      username: loginProps?.username!,
      loginType: "user",
    };
    userLogin(payload)
      .unwrap()
      .then((result) => {
        setCookie(
          CookieType.LoginResults,
          JSON.stringify({
            ...result,
            username: loginProps?.username,
          }),
        );
        setLoginProps({ ...result, username: loginProps.username! });
        setRefreshQrCode(false);
      })
      .catch((err) => {
        // Show the backend's human-readable message rather than its machine-readable error
        // slug (or the raw error object, which showErrorToast can't render as text).
        showErrorToast(
          (err as IErrorResponse)?.message ??
            (err as IErrorResponse)?.error ??
            "Login failed",
        );
      });
  };
  const handleRefreshQrCode = () => {
    setRefreshQrCode(true);
  };

  const dispatch = useDispatch();
  const { data: loginDetails, error, refetch: refetchLoginStatus } = useVerifyLoginQuery(
    {
      targetUser: `${loginProps?.username}@trovo`,
      loginID: loginProps?.loginId ?? "",
    },
    {
      pollingInterval: 5000,
      // Requires setupListeners(store.dispatch) (redux/store.ts) to have any effect. Doesn't
      // stop the browser throttling this poll's timer while the tab is backgrounded, but
      // removes the compounding delay on return: the check reruns the instant this tab regains
      // focus instead of waiting on the (possibly still-throttled) next interval tick.
      refetchOnFocus: true,
    },
  );

  // Fast path alongside the poll above, not a replacement for it: the backend already knows
  // the moment the wallet approves (LoginCallback calls BroadcastToLoginID), but until now
  // nothing on the frontend could receive that push. An open SSE connection isn't subject to
  // the same background-tab timer throttling that makes the 5s poll above pause when this tab
  // isn't focused, so this is what actually fixes "login doesn't complete until I switch back
  // to the tab" rather than just softening it. The pushed event only signals "check now" (it
  // carries the username, not the tokens) — refetchLoginStatus() re-runs the same verify call
  // the poll already does.
  const loginIdForStream = loginProps?.loginId;
  useEffect(() => {
    if (!loginIdForStream || loginDetails?.data) return;

    const source = new EventSource(`${BASE_URL}/login/stream/${loginIdForStream}`);

    source.onmessage = () => {
      refetchLoginStatus();
    };
    source.onerror = () => {
      // Connection dropped or was never supported (e.g. a proxy blocking long-lived
      // responses) — the 5s poll above is still running and remains the source of truth,
      // so there's nothing to recover here beyond closing this one out.
      source.close();
    };

    return () => {
      source.close();
    };
  }, [loginIdForStream, loginDetails?.data, refetchLoginStatus]);

  useEffect(() => {
    if (loginDetails && loginDetails.data) {
      showSuccessToast("wallet connected successfully");
      dispatch(
        setToken({
          accessToken: loginDetails.data.accessToken,
          refreshToken: loginDetails.data.refreshToken,
          // Enables the proactive refresh in getPreloadedState/scheduleTokenRefresh, which
          // otherwise never runs and leaves the app depending entirely on the reactive 401
          // path. No JWT decoding available client-side, so this assumes the same 15-minute
          // validity scheduleTokenRefresh's own refresh branch already assumes.
          expiresAt: Date.now() + 15 * 60 * 1000,
        }),
      );
      setCookie(TOKEN, loginDetails.data.accessToken);
      dispatch(setUser(loginDetails.data.userInfo));

      router.push("/");
      deleteCookie(CookieType.LoginResults);
    }
  }, [loginDetails]);
  useEffect(() => {
    if (
      (error as IErrorResponse)?.message?.includes(
        "does not have pending login session",
      )
    ) {
      handleRefreshQrCode();
    }
  }, [error]);
  const handleVerification = () => {
    window.location.href = loginProps?.dynamicLink as string;
  };
  return (
    <Container>
      <ContentContainer>
        <Title>Login with Trovo App</Title>
        <ImageContainer>
          {/* <img src={loginProps?.qrCode} alt="logo" width={220} height={220} /> */}
          <QrCodeContainer showRetry={refreshQrCode}>
            <StyledImage src={loginProps?.qrCode} alt="QR Code" />
            <LogoImage src={logo} alt="Trovotech Logo" />
          </QrCodeContainer>
          {refreshQrCode && (
            <ButtonContainer>
              <RetryButton onClick={handleLogin}>
                <GrRefreshIconContainer>
                  <GrRefresh color="#191919" size={35} fontWeight={700} />
                </GrRefreshIconContainer>
                <Text>Click to reload QR code</Text>
              </RetryButton>
            </ButtonContainer>
          )}
        </ImageContainer>
        <InfoText>
          Point your phone to this screen to scan this qrcode from your Trovo
          App .
        </InfoText>
        {isMobile && (
          <LoginButton onClick={() => handleVerification()}>Login</LoginButton>
        )}
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
const ContentContainer = styled.div``;
const Title = styled.h2`
  color: #191919;
  font-weight: 600;
  font-size: 32px;
  line-height: 39px;
  text-align: center;
  padding-bottom: 10px;
`;
const ImageContainer = styled.div`
  text-align: center;
  margin-bottom: 24px;
  position: relative;
`;
const Text = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  text-align: center;
  color: #ffffff;
  padding-top: 15px;
  font-family: inherit;
`;
const ButtonContainer = styled.div`
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  width: 220px;
  height: 220px;
  background-color: tansparent;
  z-index: 2;
  margin: auto;
  display: flex;
  justify-content: center;
  align-items: center;
`;
const RetryButton = styled.button`
  background-color: #191919;
  padding: 10px;
  cursor: pointer;
  width: 140px;
  height: 140px;
  border: none;
  color: white;
  font-size: 20px;
  border-radius: 10px;
  z-index: 3;
  font-family: inherit;
`;
const InfoText = styled.p`
  font-size: 14px;
  font-weight: 600;
  line-height: 24px;
  color: #4f4f4f;
  text-align: center;
  width: 70%;
  margin: auto;
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
`;
const GrRefreshIconContainer = styled.div`
  background-color: white;
  border-radius: 50%;
  padding: 5px;
  width: 58px;
  height: 52px;
  border-radius: 10px;
  display: inline-flex;
  justify-content: center;
  align-items: center;
`;

const StyledImage = styled.img`
  width: 220px;
  height: 220px;
`;

const LogoImage = styled(Image)`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 50px;
  height: 50px;
  border-radius: 50px;
  background-color: white;
  padding: 5px;
`;
const QrCodeContainer = styled.div<QrCodeContainerProps>`
  position: relative;
  display: inline-block;
  ${({ showRetry }) =>
    showRetry &&
    `
    & img {
      opacity: 0.1;
    }
  `}
`;
