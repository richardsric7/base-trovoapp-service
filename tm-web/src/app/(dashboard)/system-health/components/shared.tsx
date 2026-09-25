"use client";
import styled from "styled-components";

/**
 * Shared presentation for the System Health screens.
 *
 * Colours, radii and chip shape are taken from the dashboard's existing
 * components rather than invented - the status chip matches the one the Audit
 * Trail table uses (#00A859 / #BE3800 over a 1A-alpha ground, 8px radius), so
 * "healthy" looks the same here as everywhere else in Trovo Manager.
 *
 * The one addition is amber for "degraded", a state the rest of the dashboard
 * has no equivalent of: a service that is impaired but still serving is
 * neither the green nor the red the existing pages need.
 */
export const statusColor = (status?: string) => {
  switch (status) {
    case "up":
      return { fg: "#00A859", bg: "#00A8591A" };
    case "degraded":
      return { fg: "#B26A00", bg: "#B26A001A" };
    case "down":
      return { fg: "#BE3800", bg: "#BE38001A" };
    case "skipped":
      return { fg: "#6B7280", bg: "#E5E7EB" };
    default:
      return { fg: "#6B7280", bg: "#E5E7EB" };
  }
};

export const StatusPill = styled.span<{ status?: string }>`
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  text-transform: capitalize;
  white-space: nowrap;
  color: ${({ status }) => statusColor(status).fg};
  background-color: ${({ status }) => statusColor(status).bg};
`;

export const Panel = styled.div`
  background: #ffffff;
  border: 1px solid #e5e5ef;
  border-radius: 10px;
  padding: 20px;
`;

export const PanelTitle = styled.h2`
  font-size: 16px;
  font-weight: 600;
  color: #00225a;
  margin: 0 0 4px;
`;

export const PanelHint = styled.p`
  font-size: 14px;
  color: #828282;
  margin: 0 0 20px;
  line-height: 1.5;
`;

export const Grid = styled.div<{ min?: string }>`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(${({ min }) => min || "320px"}, 1fr));
  gap: 16px;
`;

/**
 * Shown when a screen's backing service is not configured or unreachable.
 * Deliberately explicit: an empty chart reads as "nothing is wrong", which is
 * the opposite of the truth when the data source is down.
 */
export const Unavailable = ({ reason }: { reason?: string }) => (
  <UnavailableBox>
    <strong>This view is unavailable</strong>
    <span>{reason || "The data source for this screen is not configured in this environment."}</span>
  </UnavailableBox>
);

const UnavailableBox = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 32px;
  border: 1px dashed #e5e5ef;
  border-radius: 10px;
  text-align: center;
  color: #828282;

  strong {
    font-size: 14px;
    color: #00225a;
  }
  span {
    font-size: 14px;
    line-height: 1.5;
  }
`;

export const Loading = ({ label }: { label?: string }) => (
  <LoadingBox>{label || "Loading…"}</LoadingBox>
);

const LoadingBox = styled.div`
  padding: 32px;
  text-align: center;
  color: #828282;
  font-size: 14px;
`;

/** Human-readable uptime, e.g. "2d 4h" — seconds are noise at this scale. */
export const formatUptime = (seconds: number) => {
  if (!seconds || seconds < 0) return "—";
  const d = Math.floor(seconds / 86400);
  const h = Math.floor((seconds % 86400) / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  if (d > 0) return `${d}d ${h}h`;
  if (h > 0) return `${h}h ${m}m`;
  return `${m}m`;
};

/** Latency with the unit that keeps it readable. */
export const formatLatency = (ms: number) => {
  if (ms === undefined || ms === null) return "—";
  if (ms < 1) return "<1ms";
  if (ms < 1000) return `${ms}ms`;
  return `${(ms / 1000).toFixed(1)}s`;
};

export const formatTime = (iso?: string) => {
  if (!iso) return "—";
  const d = new Date(iso);
  return Number.isNaN(d.getTime()) ? "—" : d.toLocaleString();
};

/** "3 minutes ago" — what an operator actually wants when scanning errors. */
export const relativeTime = (iso?: string) => {
  if (!iso) return "—";
  const then = new Date(iso).getTime();
  if (Number.isNaN(then)) return "—";
  const secs = Math.floor((Date.now() - then) / 1000);
  if (secs < 60) return "just now";
  if (secs < 3600) return `${Math.floor(secs / 60)}m ago`;
  if (secs < 86400) return `${Math.floor(secs / 3600)}h ago`;
  return `${Math.floor(secs / 86400)}d ago`;
};
