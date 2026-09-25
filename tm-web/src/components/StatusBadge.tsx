import styled from "styled-components";

interface StatusStyle {
  background: string;
  color: string;
}

const STATUS_STYLES: Record<string, StatusStyle> = {
  draft: { background: "#F2F2F2", color: "#4F4F4F" },
  submitted: { background: "#E6F3FF", color: "#007CDF" },
  under_review: { background: "#FFF4E0", color: "#B25E00" },
  approved: { background: "#E6F7EC", color: "#1D8A46" },
  rejected: { background: "#FDEBEA", color: "#EB5757" },
  pending: { background: "#F2F6F9", color: "#00225A" },
  overdue: { background: "#FDEBEA", color: "#EB5757" },
  waived: { background: "#F1E9FB", color: "#7C4DBE" },
};

const DEFAULT_STYLE: StatusStyle = { background: "#F2F6F9", color: "#00225A" };

const humanize = (value: string) =>
  value
    .toLowerCase()
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");

export interface StatusBadgeProps {
  status: string;
}

export const StatusBadge = ({ status }: StatusBadgeProps) => {
  const key = (status || "").toLowerCase();
  const style = STATUS_STYLES[key] || DEFAULT_STYLE;

  return (
    <Badge $background={style.background} $color={style.color}>
      {humanize(key || "unknown")}
    </Badge>
  );
};

export default StatusBadge;

const Badge = styled.span<{ $background: string; $color: string }>`
  display: inline-block;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
  background: ${({ $background }) => $background};
  color: ${({ $color }) => $color};
`;
