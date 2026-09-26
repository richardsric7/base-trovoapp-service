import { ReportRange } from "@/redux/api/p2p";

// formatAmount renders a decimal-string amount (tm-api never sends
// floats for money, to avoid precision loss) as a locale-formatted
// number for display - never re-parsed back into anything sent to the
// server.
export const formatAmount = (value: string | undefined): string => {
  const n = Number(value ?? "0");
  if (!Number.isFinite(n)) return "0";
  return n.toLocaleString(undefined, { maximumFractionDigits: 2 });
};

// formatDuration turns a second count into a compact "Xd Yh", "Xh Ym", or
// "Xm" label - dispute resolution and order completion times are commonly
// in the hours-to-days range, so this stays readable without a full
// duration-formatting library.
export const formatDuration = (totalSeconds: number): string => {
  if (!totalSeconds || totalSeconds <= 0) return "-";
  const days = Math.floor(totalSeconds / 86400);
  const hours = Math.floor((totalSeconds % 86400) / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  if (days > 0) return `${days}d ${hours}h`;
  if (hours > 0) return `${hours}h ${minutes}m`;
  if (minutes > 0) return `${minutes}m`;
  return `${totalSeconds}s`;
};

// formatDateLabel turns a "YYYY-MM-DD" series date into a short axis
// label ("Sep 26").
export const formatDateLabel = (isoDate: string): string => {
  const d = new Date(isoDate + "T00:00:00Z");
  if (Number.isNaN(d.getTime())) return isoDate;
  return d.toLocaleDateString(undefined, { month: "short", day: "numeric", timeZone: "UTC" });
};

// titleCase turns a backend enum key ("BANK_TRANSFER", "IN_FAVOR_OF_BUYER")
// into a display label ("Bank transfer", "In favor of buyer").
export const titleCase = (key: string): string => {
  if (!key) return "Unknown";
  const words = key.toLowerCase().split("_");
  return words.map((w, i) => (i === 0 ? w.charAt(0).toUpperCase() + w.slice(1) : w)).join(" ");
};

export const REPORT_RANGE_OPTIONS: { label: string; value: ReportRange }[] = [
  { label: "7D", value: "7d" },
  { label: "30D", value: "30d" },
  { label: "90D", value: "90d" },
  { label: "1Y", value: "1y" },
];

export const CHART_COLORS = [
  "#007CDF",
  "#62ACE8",
  "#00A859",
  "#F5A623",
  "#EB5757",
  "#9B51E0",
  "#ACD1EF",
  "#00225A",
];
