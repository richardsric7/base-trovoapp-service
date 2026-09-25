"use client";

import { Modal, showErrorToast, showSuccessToast } from "@/components";
import {
  useRequestWalletLinkAuthorizationMutation,
  useVerifyWalletLinkMutation,
} from "@/redux/api/sharedstakeholders";
import { FormEvent, useEffect, useRef, useState } from "react";
import { FaCheck, FaWallet } from "react-icons/fa6";
import styled from "styled-components";

interface LinkWalletModalProps {
  isOpen: boolean;
  onClose: () => void;
  onLinked: () => void;
}

interface WalletChallenge {
  authId: string;
  dynamicLink?: string;
  qrCode?: string;
}

const getErrorMessage = (error: any, fallback: string) =>
  error?.data?.message ?? error?.data?.error ?? error?.error ?? fallback;

const normalizeWalletUsername = (username: string) => {
  const normalized = username.trim().toLowerCase();
  return normalized.endsWith("@trovo") ? normalized : `${normalized}@trovo`;
};

const getQrCodeSource = (qrCode?: string) => {
  if (!qrCode) return undefined;

  const source = qrCode.trim();
  const isImageSource = /^(data:|https?:\/\/|blob:|\/)/i.test(source);

  return isImageSource ? source : `data:image/png;base64,${source}`;
};

export default function LinkWalletModal({
  isOpen,
  onClose,
  onLinked,
}: LinkWalletModalProps) {
  const [walletUsername, setWalletUsername] = useState("");
  const [challenge, setChallenge] = useState<WalletChallenge>();
  const isOpenRef = useRef(isOpen);
  const [requestAuthorization, requestState] =
    useRequestWalletLinkAuthorizationMutation();
  const [verifyWallet, verifyState] = useVerifyWalletLinkMutation();

  useEffect(() => {
    isOpenRef.current = isOpen;
    if (!isOpen) {
      setWalletUsername("");
      setChallenge(undefined);
    }
  }, [isOpen]);

  const closeModal = () => {
    isOpenRef.current = false;
    onClose();
  };

  const startAuthorization = async (event: FormEvent) => {
    event.preventDefault();

    if (!walletUsername.trim()) {
      showErrorToast("Enter your Trovo Wallet username");
      return;
    }

    try {
      const normalizedUsername = normalizeWalletUsername(walletUsername);
      const response = await requestAuthorization({
        wallet_username: normalizedUsername,
      }).unwrap();

      if (!isOpenRef.current) return;

      setWalletUsername(normalizedUsername);
      setChallenge(response);
    } catch (error) {
      showErrorToast(
        getErrorMessage(error, "Unable to start wallet authorization"),
      );
    }
  };

  const verifyAuthorization = async () => {
    if (!challenge) return;

    try {
      const response = await verifyWallet({
        auth_id: challenge.authId,
        wallet_username: normalizeWalletUsername(walletUsername),
      }).unwrap();

      showSuccessToast(response.message || "Wallet linked successfully");
      onLinked();
      closeModal();
    } catch (error) {
      showErrorToast(
        getErrorMessage(error, "Wallet approval has not been verified yet"),
      );
    }
  };

  const qrSource = getQrCodeSource(challenge?.qrCode);

  return (
    <Modal
      title="Link Trovo Wallet"
      isOpen={isOpen}
      onClose={closeModal}
      style={{ maxWidth: "calc(100vw - 32px)" }}
    >
      {!challenge ? (
        <WalletForm onSubmit={startAuthorization}>
          <Intro>
            <WalletIcon>
              <FaWallet />
            </WalletIcon>
            <p>
              Connect your Trovo Wallet to approve secure organization actions.
            </p>
          </Intro>
          <Field>
            <Label htmlFor="wallet-username">Wallet username</Label>
            <Input
              id="wallet-username"
              value={walletUsername}
              onChange={(event) => setWalletUsername(event.target.value)}
              placeholder="e.g. chidinma@trovo"
              autoComplete="off"
              autoFocus
            />
            <Hint>You can enter the username with or without @trovo.</Hint>
          </Field>
          <PrimaryButton type="submit" disabled={requestState.isLoading}>
            {requestState.isLoading ? "Requesting..." : "Continue"}
          </PrimaryButton>
        </WalletForm>
      ) : (
        <ApprovalStep>
          <StepIcon>
            <FaCheck />
          </StepIcon>
          <StepTitle>Approve the request</StepTitle>
          <StepCopy>
            Open Trovo Wallet or scan the QR code, approve the request, then
            return here to verify it.
          </StepCopy>
          {qrSource && (
            <QRCode src={qrSource} alt="Trovo Wallet authorization QR code" />
          )}
          {challenge.dynamicLink && (
            <WalletLink
              href={challenge.dynamicLink}
              target="_blank"
              rel="noreferrer"
            >
              Open Trovo Wallet
            </WalletLink>
          )}
          <PrimaryButton
            type="button"
            onClick={verifyAuthorization}
            disabled={verifyState.isLoading}
          >
            {verifyState.isLoading
              ? "Verifying..."
              : "I have approved — Verify"}
          </PrimaryButton>
          <SecondaryButton
            type="button"
            onClick={() => setChallenge(undefined)}
          >
            Use another wallet
          </SecondaryButton>
        </ApprovalStep>
      )}
    </Modal>
  );
}

const WalletForm = styled.form`
  display: flex;
  flex-direction: column;
  gap: 20px;
`;
const Intro = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  color: #667085;
  p {
    margin: 0;
    font-size: 14px;
    line-height: 1.5;
  }
`;
const WalletIcon = styled.span`
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  flex-shrink: 0;
  background: #eaf5fd;
  color: #007cdf;
  font-size: 18px;
`;
const Field = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;
const Label = styled.label`
  color: #00225a;
  font-size: 14px;
  font-weight: 600;
`;
const Input = styled.input`
  height: 48px;
  border: 1px solid #d0d5dd;
  border-radius: 8px;
  padding: 0 14px;
  color: #00225a;
  font: inherit;
  outline: none;
  &:focus {
    border-color: #007cdf;
    box-shadow: 0 0 0 3px rgba(0, 124, 223, 0.12);
  }
`;
const Hint = styled.span`
  color: #828282;
  font-size: 12px;
`;
const PrimaryButton = styled.button`
  width: 100%;
  height: 48px;
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
const SecondaryButton = styled.button`
  border: 0;
  background: transparent;
  color: #667085;
  font: inherit;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
`;
const ApprovalStep = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  text-align: center;
`;
const StepIcon = styled.span`
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: #e9f8f1;
  color: #00875a;
`;
const StepTitle = styled.h3`
  margin: 0;
  color: #00225a;
  font-size: 18px;
`;
const StepCopy = styled.p`
  margin: 0;
  color: #667085;
  font-size: 13px;
  line-height: 1.55;
  max-width: 390px;
`;
const QRCode = styled.img`
  width: 190px;
  height: 190px;
  object-fit: contain;
  padding: 8px;
  border: 1px solid #e5e5ef;
  border-radius: 10px;
  background: #fff;
`;
const WalletLink = styled.a`
  color: #007cdf;
  font-size: 14px;
  font-weight: 600;
  text-decoration: none;
`;
