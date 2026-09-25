"use client";

import {
  StakeholderAuthorizationAction,
  StakeholderAuthorizationEntity,
  useCreateStakeholderAuthorizationMutation,
  useGetStakeholderProfileQuery,
  useVerifyStakeholderAuthorizationMutation,
} from "@/redux/api/sharedstakeholders";
import { Modal, message } from "antd";
import { useEffect, useMemo, useState } from "react";
import { FaCheck, FaWallet } from "react-icons/fa6";
import styled from "styled-components";
import Image from "next/image";
import logo from "@/assets/images/logo-qr-code.svg";

interface Props {
  open: boolean;
  onClose: () => void;
  action: StakeholderAuthorizationAction;
  entityType: StakeholderAuthorizationEntity;
  entityId?: string;
  title: string;
  description: string;
  confirmLabel: string;
  onAuthorize: (challengeId: string) => Promise<void>;
  isAuthorizing?: boolean;
  children?: React.ReactNode;
}

const errorText = (error: any, fallback: string) =>
  error?.data?.message ?? error?.data?.error ?? error?.error ?? fallback;

const getQrCodeSource = (qrCode?: string) => {
  if (!qrCode) return undefined;

  const source = qrCode.trim();
  const isImageSource = /^(data:|https?:\/\/|blob:|\/)/i.test(source);

  return isImageSource ? source : `data:image/png;base64,${source}`;
};

export default function StepUpAuthorizationModal({
  open,
  onClose,
  action,
  entityType,
  entityId,
  title,
  description,
  confirmLabel,
  onAuthorize,
  isAuthorizing,
  children,
}: Props) {
  const { data: profile, isLoading: profileLoading } =
    useGetStakeholderProfileQuery(undefined, { skip: !open });
  const [createChallenge, createState] =
    useCreateStakeholderAuthorizationMutation();
  const [verifyChallenge, verifyState] =
    useVerifyStakeholderAuthorizationMutation();
  const [challenge, setChallenge] = useState<{
    id: string;
    auth_id: string;
    dynamic_link?: string;
    qr_code?: string;
    expires_at: string;
  }>();
  const [verified, setVerified] = useState(false);

  useEffect(() => {
    if (!open) {
      setChallenge(undefined);
      setVerified(false);
    }
  }, [open]);

  const wallet = profile?.data.wallet;
  const linked = Boolean(wallet?.is_wallet_linked);
  const expired = useMemo(
    () =>
      Boolean(
        challenge && new Date(challenge.expires_at).getTime() <= Date.now(),
      ),
    [challenge],
  );

  const start = async () => {
    if (!entityId) return;
    try {
      const response = await createChallenge({
        action,
        entity_type: entityType,
        entity_id: entityId,
      }).unwrap();
      setChallenge(response.data);
      setVerified(false);
    } catch (error) {
      message.error(errorText(error, "Unable to start wallet authorization."));
    }
  };

  const verify = async () => {
    if (!challenge) return;
    try {
      const response = await verifyChallenge({
        challengeId: challenge.id,
        auth_id: challenge.auth_id,
      }).unwrap();
      if (response.data.status === "verified") {
        setVerified(true);
        message.success("Wallet authorization verified.");
      }
    } catch (error) {
      message.error(
        errorText(error, "Wallet approval has not been verified yet."),
      );
    }
  };

  const finish = async () => {
    if (!challenge || !verified) return;
    await onAuthorize(challenge.id);
  };

  const qrSource = getQrCodeSource(challenge?.qr_code);

  return (
    <Modal
      open={open}
      onCancel={onClose}
      footer={null}
      centered
      width={540}
      destroyOnClose
    >
      <Content>
        <Title>{title}</Title>
        <Description>{description}</Description>
        {children}

        {profileLoading && <Notice>Checking your linked wallet...</Notice>}
        {!profileLoading && !linked && (
          <WalletWarning>
            <FaWallet />
            <div>
              <strong>Link your Trovo Wallet first</strong>
              <p>
                This action requires a linked Trovo Wallet on your stakeholder
                profile.
              </p>
            </div>
          </WalletWarning>
        )}
        {linked && !challenge && (
          <PrimaryButton
            onClick={start}
            disabled={createState.isLoading || !entityId}
          >
            {createState.isLoading
              ? "Starting..."
              : "Continue with Trovo Wallet"}
          </PrimaryButton>
        )}
        {linked && challenge && !verified && (
          <WalletStep>
            <StepTitle>Approve in Trovo Wallet</StepTitle>
            <StepCopy>
              Scan the QR code or open the wallet request, approve it, then
              verify below. The request expires in 10 minutes.
            </StepCopy>
            {qrSource && (
              <QRCodeContainer>
                <QRCode
                  src={qrSource}
                  alt="Trovo Wallet authorization QR code"
                />
                <LogoImage src={logo} alt="Trovotech logo" />
              </QRCodeContainer>
            )}
            {challenge.dynamic_link && (
              <WalletLink
                href={challenge.dynamic_link}
                target="_blank"
                rel="noreferrer"
              >
                Open Trovo Wallet
              </WalletLink>
            )}
            <PrimaryButton
              onClick={expired ? start : verify}
              disabled={verifyState.isLoading || createState.isLoading}
            >
              {verifyState.isLoading
                ? "Verifying..."
                : expired
                  ? "Create a new request"
                  : "I have approved — Verify"}
            </PrimaryButton>
          </WalletStep>
        )}
        {verified && (
          <VerifiedBox>
            <FaCheck />
            <span>Wallet authorization verified</span>
          </VerifiedBox>
        )}
        {verified && (
          <PrimaryButton onClick={finish} disabled={isAuthorizing}>
            {isAuthorizing ? "Authorizing..." : confirmLabel}
          </PrimaryButton>
        )}
      </Content>
    </Modal>
  );
}

export const Content = styled.div`
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 16px;
`;
export const Title = styled.h2`
  margin: 0;
  color: #00225a;
  font-size: 20px;
`;
export const Description = styled.p`
  margin: 0;
  color: #667085;
  font-size: 13px;
  line-height: 1.5;
`;
const Notice = styled.div`
  padding: 16px;
  border-radius: 10px;
  background: #f2f6f9;
  color: #667085;
  text-align: center;
`;
const WalletWarning = styled.div`
  display: flex;
  gap: 12px;
  padding: 16px;
  border-radius: 10px;
  background: #fff4e8;
  color: #9a4d00;
  svg {
    margin-top: 3px;
    flex-shrink: 0;
  }
  p {
    margin: 4px 0 0;
    font-size: 12px;
    line-height: 1.5;
  }
`;
const WalletStep = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-radius: 12px;
  background: #f5f8fa;
`;
const StepTitle = styled.h3`
  margin: 0;
  color: #00225a;
  font-size: 15px;
`;
const StepCopy = styled.p`
  margin: 0;
  color: #667085;
  font-size: 12px;
  text-align: center;
  line-height: 1.5;
`;
const QRCodeContainer = styled.div`
  position: relative;
  display: inline-block;
`;
const QRCode = styled.img`
  width: 180px;
  height: 180px;
  object-fit: contain;
  background: #fff;
  border-radius: 8px;
  padding: 8px;
`;
const LogoImage = styled(Image)`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: #fff;
  padding: 5px;
`;
const WalletLink = styled.a`
  display: none;
  color: #007cdf;
  font-size: 13px;
  font-weight: 600;
  text-decoration: none;

  @media (max-width: 767px) {
    display: inline-block;
  }
`;
export const PrimaryButton = styled.button`
  width: 100%;
  height: 44px;
  border: 0;
  border-radius: 8px;
  background: #007cdf;
  color: #fff;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
  &:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }
`;
const VerifiedBox = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 13px;
  border-radius: 9px;
  background: #e9f8f1;
  color: #00875a;
  font-weight: 600;
`;
