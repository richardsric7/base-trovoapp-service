"use client";
import React from "react";
import styled from "styled-components";
import { PayoutItemStatus, PayoutStatus } from "@/redux/api/proceedPayouts";

// Shared pieces of the proceeds payout pages.

export const errorMessage = (err: any, fallback: string) =>
  err?.data?.error || err?.data?.message || err?.error || fallback;

export const STATUS_LABELS: Record<PayoutStatus, string> = {
  REGISTERED: "Registered",
  PREPARE_REQUESTED: "Preparation requested",
  PREPARING: "Preparing schedule",
  LOCKED: "Awaiting approval",
  APPROVED: "Approved",
  FUNDING_CHECK_REQUESTED: "Checking funding",
  PAYING: "Paying",
  PAUSED: "Paused",
  COMPLETED: "Completed",
  COMPLETED_WITH_FAILURES: "Completed with failures",
  CANCELLED: "Cancelled",
};

const TONES: Record<string, [string, string]> = {
  good: ["#00A859", "#00A8591A"],
  bad: ["#D92D20", "#D92D201A"],
  warn: ["#B54708", "#FEF0C7"],
  info: ["#007CDF", "#007CDF1A"],
  muted: ["#475467", "#F2F4F7"],
};

const tone = (s: string) => {
  if (["COMPLETED", "PAID", "MINED"].includes(s)) return TONES.good;
  if (["CANCELLED", "FAILED", "REVERTED", "COMPLETED_WITH_FAILURES"].includes(s)) return TONES.bad;
  if (["PAUSED", "LOCKED", "EXCLUDED", "DROPPED"].includes(s)) return TONES.warn;
  if (["PAYING", "PREPARING", "FUNDING_CHECK_REQUESTED", "PREPARE_REQUESTED", "APPROVED", "QUEUED", "SUBMITTED"].includes(s))
    return TONES.info;
  return TONES.muted;
};

export const StatusPill = ({ status }: { status: PayoutStatus | PayoutItemStatus | string }) => {
  const [fg, bg] = tone(status);
  return (
    <Pill $fg={fg} $bg={bg}>
      {STATUS_LABELS[status as PayoutStatus] ?? status.charAt(0) + status.slice(1).toLowerCase()}
    </Pill>
  );
};

export const amount = (v: string | number | undefined, code?: string) => {
  if (v === undefined || v === null || v === "") return "—";
  const n = Number(v);
  const s = Number.isFinite(n) ? n.toLocaleString(undefined, { maximumFractionDigits: 6 }) : String(v);
  return code ? `${s} ${code}` : s;
};

export const date = (v?: string) => (v && !v.startsWith("0001") ? new Date(v).toLocaleString() : "—");

export const short = (addr?: string) => (addr && addr.length > 14 ? `${addr.slice(0, 8)}…${addr.slice(-6)}` : addr || "—");

export const Card = styled.section`
  background: #ffffff;
  padding: 24px;
  border-radius: 24px;
  margin-bottom: 20px;
`;

export const Title = styled.h1`
  font-size: 22px;
  font-weight: 700;
  margin: 0;
  color: #00225a;
`;

export const SubTitle = styled.h2`
  font-size: 17px;
  font-weight: 700;
  margin: 0 0 12px;
  color: #00225a;
`;

export const Muted = styled.p`
  font-size: 14px;
  color: #667085;
  margin: 6px 0 0;
  line-height: 22px;
  max-width: 860px;
`;

export const Row = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  margin: 16px 0;
`;

export const Input = styled.input<{ $wide?: boolean }>`
  flex: ${(p) => (p.$wide ? "1 1 360px" : "0 1 200px")};
  min-width: 0;
  padding: 10px 14px;
  border: 1px solid #d9e1ec;
  border-radius: 8px;
  font-size: 14px;
  color: #00225a;
  background: #f2f6f9;
`;

export const Select = styled.select`
  padding: 10px 14px;
  border: 1px solid #d9e1ec;
  border-radius: 8px;
  font-size: 14px;
  color: #00225a;
  background: #f2f6f9;
`;

export const Button = styled.button<{ $variant?: "primary" | "danger" | "ghost" }>`
  padding: 9px 18px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid ${(p) => (p.$variant === "danger" ? "#D92D20" : "#007cdf")};
  color: ${(p) => (p.$variant === "primary" ? "#ffffff" : p.$variant === "danger" ? "#D92D20" : "#007cdf")};
  background: ${(p) => (p.$variant === "primary" ? "#007cdf" : "transparent")};
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`;

export const LinkButton = styled.button<{ $danger?: boolean }>`
  background: none;
  border: none;
  padding: 0 8px 0 0;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  color: ${(p) => (p.$danger ? "#D92D20" : "#007cdf")};
`;

export const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(190px, 1fr));
  gap: 16px 24px;
`;

export const Stat = ({ label, value }: { label: string; value: React.ReactNode }) => (
  <div>
    <StatLabel>{label}</StatLabel>
    <StatValue>{value}</StatValue>
  </div>
);

export const Note = styled.div<{ $tone?: "warn" | "info" }>`
  margin-top: 16px;
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 14px;
  line-height: 21px;
  color: ${(p) => (p.$tone === "warn" ? "#B54708" : "#00225a")};
  background: ${(p) => (p.$tone === "warn" ? "#FEF0C7" : "#EFF8FF")};
`;

export const Mono = styled.span`
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 12px;
  word-break: break-all;
`;

const StatLabel = styled.p`
  font-size: 13px;
  color: #828282;
  margin: 0 0 4px;
`;

const StatValue = styled.div`
  font-size: 15px;
  font-weight: 600;
  color: #00225a;
  word-break: break-all;
`;

const Pill = styled.span<{ $fg: string; $bg: string }>`
  display: inline-block;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  color: ${(p) => p.$fg};
  background: ${(p) => p.$bg};
  white-space: nowrap;
`;
