"use client";

/**
 * The UI shown when a route crashes. Shared by every error.tsx boundary so a
 * failure looks the same wherever it happens.
 *
 * Two things it deliberately does: offers a recovery action that re-renders the
 * segment rather than reloading the whole app, and shows the support reference
 * so a user can quote it - that reference is what ties their report to the
 * error record and, through it, to the backend request that failed.
 */
import styled from "styled-components";
import { getLastRequestId } from "@/observability";

interface ErrorStateProps {
  error: Error & { digest?: string };
  reset: () => void;
  /** What failed, e.g. "the organisations page". Keeps the message specific. */
  what?: string;
}

export const ErrorState = ({ error, reset, what }: ErrorStateProps) => {
  const requestId = getLastRequestId();
  const reference = error.digest ?? requestId;

  return (
    <Wrapper>
      <Title>Something went wrong</Title>
      <Message>
        {what
          ? `We could not load ${what}. The problem has been reported to the team.`
          : "This page could not be displayed. The problem has been reported to the team."}
      </Message>
      {reference && (
        <Reference>
          Reference: <code>{reference}</code>
        </Reference>
      )}
      <Actions>
        <PrimaryButton onClick={() => reset()}>Try again</PrimaryButton>
        <SecondaryButton onClick={() => window.location.assign("/")}>
          Go to dashboard
        </SecondaryButton>
      </Actions>
    </Wrapper>
  );
};

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  padding: 32px;
  text-align: center;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin: 0 0 12px;
`;

const Message = styled.p`
  font-size: 14px;
  line-height: 1.6;
  color: #5b6b8c;
  max-width: 420px;
  margin: 0 0 16px;
`;

const Reference = styled.p`
  font-size: 12px;
  color: #8a97b1;
  margin: 0 0 24px;

  code {
    font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
    background: #eef1f6;
    padding: 2px 6px;
    border-radius: 4px;
  }
`;

const Actions = styled.div`
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: center;
`;

const PrimaryButton = styled.button`
  padding: 10px 24px;
  font-size: 14px;
  font-weight: 500;
  color: #fff;
  background: #00225a;
  border: none;
  border-radius: 8px;
  cursor: pointer;

  &:hover {
    opacity: 0.9;
  }
`;

const SecondaryButton = styled(PrimaryButton)`
  color: #00225a;
  background: transparent;
  border: 1px solid #d6dbe6;
`;
