"use client";
import React from "react";
import styled from "styled-components";
import { usePathname, useRouter } from "next/navigation";
import { showErrorToast, showSuccessToast } from "@/components";
import { errorMessage } from "../../dividendandyield/components/ui";

// The generic page pieces are shared with the proceeds payout pages.
export {
  amount,
  Button,
  Card,
  date,
  errorMessage,
  Grid,
  Input,
  LinkButton,
  Mono,
  Muted,
  Note,
  Row,
  Select,
  short,
  Stat,
  SubTitle,
  Title,
} from "../../dividendandyield/components/ui";

const TONES: Record<string, [string, string]> = {
  good: ["#00A859", "#00A8591A"],
  bad: ["#BE3800", "#BE38001A"],
  warn: ["#B26A00", "#B26A001A"],
  info: ["#007CDF", "#007CDF1A"],
  muted: ["#475467", "#F2F4F7"],
};

const GOOD = ["LIVE", "MATCHED", "complete", "COMPLETE", "SETTLED", "PROCESSED", "DELIVERED", "DISTRIBUTED", "PAID", "active", "up", "EXECUTED", "INTERNAL", "Fresh"];
const BAD = ["HALTED", "DRIFT", "LEDGER_DRIFT", "NO_POSITION", "ERROR", "rejected", "failed", "REJECTED", "FAILED", "DEAD_LETTER", "ESCALATED", "suspended", "down", "CANCELLED", "cancelled", "Stale"];
const WARN = ["SETUP", "AWAITING_APPROVAL", "SNAPSHOTTED", "NEEDS_MANUAL", "degraded", "pending-execution", "awaiting-payment", "HANDLED", "sandbox"];

export const tone = (s: string) =>
  GOOD.includes(s) ? TONES.good : BAD.includes(s) ? TONES.bad : WARN.includes(s) ? TONES.warn : s ? TONES.info : TONES.muted;

export const Pill = ({ status, label }: { status: string; label?: string }) => {
  const [fg, bg] = tone(status);
  return (
    <PillBox $fg={fg} $bg={bg} $raw={label === undefined}>
      {label ?? status.replace(/_/g, " ")}
    </PillBox>
  );
};

const PillBox = styled.span<{ $fg: string; $bg: string; $raw: boolean }>`
  display: inline-block;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  color: ${(p) => p.$fg};
  background: ${(p) => p.$bg};
  white-space: nowrap;
  text-transform: ${(p) => (p.$raw ? "capitalize" : "none")};
`;

export const ngn = (v?: string | number) => {
  if (v === undefined || v === null || v === "") return "—";
  const n = Number(v);
  return Number.isFinite(n) ? `₦${n.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}` : String(v);
};

export const num = (v?: string | number, max = 4) => {
  if (v === undefined || v === null || v === "") return "—";
  const n = Number(v);
  return Number.isFinite(n) ? n.toLocaleString(undefined, { maximumFractionDigits: max }) : String(v);
};

// StatCard is one of the headline numbers at the top of a page.
export const StatCard = ({ title, value, sub, color = "#00225A", bg = "#F2F6F9" }: { title: string; value: React.ReactNode; sub?: React.ReactNode; color?: string; bg?: string }) => (
  <StatBox>
    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start" }}>
      <div>
        <StatTitle>{title}</StatTitle>
        <StatNumber style={{ color }}>{value}</StatNumber>
      </div>
      <Dot style={{ background: bg, color }}>●</Dot>
    </div>
    {sub && <StatSub style={{ color }}>{sub}</StatSub>}
  </StatBox>
);

export const Stats = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
`;

const StatBox = styled.div`
  background: #ffffff;
  border-radius: 20px;
  padding: 20px;
`;
const StatTitle = styled.div`
  font-size: 14px;
  color: #828282;
`;
const StatNumber = styled.div`
  font-size: 26px;
  font-weight: 700;
  margin-top: 6px;
`;
const StatSub = styled.div`
  font-size: 13px;
  font-weight: 600;
  margin-top: 10px;
`;
const Dot = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  font-size: 12px;
`;

// Tabs within a Public Markets page, kept in the URL (?tab=).
export const Tabs = ({ tabs, current, onChange }: { tabs: { key: string; label: string }[]; current: string; onChange: (k: string) => void }) => (
  <TabRow>
    {tabs.map((t) => (
      <TabButton key={t.key} $on={t.key === current} onClick={() => onChange(t.key)}>
        {t.label}
      </TabButton>
    ))}
  </TabRow>
);

const TabRow = styled.div`
  display: flex;
  gap: 4px;
  border-bottom: 1px solid #e4e7ec;
  margin: 18px 0;
  overflow-x: auto;
`;
const TabButton = styled.button<{ $on: boolean }>`
  border: none;
  background: none;
  padding: 10px 14px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  color: ${(p) => (p.$on ? "#007cdf" : "#667085")};
  border-bottom: 2px solid ${(p) => (p.$on ? "#007cdf" : "transparent")};
`;

// Field is a labelled form control.
export const Field = ({ label, hint, children }: { label: string; hint?: string; children: React.ReactNode }) => (
  <FieldBox>
    <label>{label}</label>
    {children}
    {hint && <small>{hint}</small>}
  </FieldBox>
);

const FieldBox = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
  label {
    font-size: 13px;
    font-weight: 600;
    color: #344054;
  }
  small {
    font-size: 12px;
    color: #667085;
  }
  input,
  select,
  textarea {
    padding: 10px 12px;
    border: 1px solid #d9e1ec;
    border-radius: 8px;
    font-size: 14px;
    color: #00225a;
    background: #f2f6f9;
    width: 100%;
    box-sizing: border-box;
  }
`;

export const FormGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
`;

export const Panel = styled.div`
  background: #ffffff;
  border-radius: 24px;
  padding: 20px 24px;
`;

export const TwoCols = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(360px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
`;

export const Back = ({ href, label }: { href: string; label: string }) => {
  const router = useRouter();
  return (
    <BackLink onClick={() => router.push(href)}>
      ‹ {label}
    </BackLink>
  );
};

const BackLink = styled.button`
  border: none;
  background: none;
  color: #007cdf;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  padding: 0 0 14px;
`;

// The section's own navigation (the sidebar lists the same pages).
const PAGES = [
  { href: "/publicmarkets", label: "Overview" },
  { href: "/publicmarkets/assets", label: "Assets" },
  { href: "/publicmarkets/orders", label: "Orders" },
  { href: "/publicmarkets/reconciliation", label: "Reconciliation" },
  { href: "/publicmarkets/corporate-actions", label: "Corporate Actions" },
  { href: "/publicmarkets/exchanges", label: "Exchange Partners" },
  { href: "/publicmarkets/wallets", label: "Wallet Provisioning" },
  { href: "/publicmarkets/prices", label: "Price Oracle" },
  { href: "/publicmarkets/settings", label: "Settings" },
];
export const PUBLIC_MARKETS_PAGES = PAGES;

export const SectionNav = () => {
  const path = usePathname();
  const router = useRouter();
  const active = (href: string) => (href === "/publicmarkets" ? path === href : path.startsWith(href));
  return (
    <TabRow style={{ marginTop: 0 }}>
      {PAGES.map((p) => (
        <TabButton key={p.href} $on={active(p.href)} onClick={() => router.push(p.href)}>
          {p.label}
        </TabButton>
      ))}
    </TabRow>
  );
};

// act runs a mutation with an optional confirmation or reason prompt and
// toasts the outcome. It returns the response, or undefined.
export async function act<R extends { message?: string }>(
  run: (reason?: string) => Promise<R>,
  opts: { confirm?: string; reason?: string; success?: string; failure?: string } = {},
): Promise<R | undefined> {
  let reason: string | undefined;
  if (opts.reason) {
    reason = window.prompt(opts.reason) ?? undefined;
    if (!reason?.trim()) return undefined;
  } else if (opts.confirm && !window.confirm(opts.confirm)) {
    return undefined;
  }
  try {
    const res = await run(reason);
    showSuccessToast(opts.success ?? res?.message ?? "Done");
    return res;
  } catch (err) {
    showErrorToast(errorMessage(err, opts.failure ?? "The action failed"));
    return undefined;
  }
}

export const ORDER_STATES = [
  "awaiting-payment",
  "queued",
  "filled-from-inventory",
  "pending-execution",
  "executed",
  "submitted",
  "chain_final",
  "settlement_final",
  "complete",
  "rejected",
  "failed",
  "cancelled",
];

export const PATH_LABELS: Record<string, string> = { FAST: "Fast (inventory)", SLOW: "Slow (next session)", NETTED: "Netted", "": "—" };

// tierNames lists the tiers of a settings JSON map ({"Tier 1": 1200}).
export const tierNames = (json?: string) => {
  try {
    return Object.keys(JSON.parse(json || "{}"));
  } catch {
    return [];
  }
};
