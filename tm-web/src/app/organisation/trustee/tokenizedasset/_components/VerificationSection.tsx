import styled from "styled-components";
import DocumentCard from "./DocumentCard";
import {
  AiOutlineCheckCircle,
  AiOutlineClockCircle,
  AiOutlineCloseCircle,
  AiOutlineMinusCircle,
} from "react-icons/ai";
import { FiCalendar } from "react-icons/fi";
import { IDueDiligenceItem } from "@/redux/api/trustees/interface";
import { IStakeholderAssetTokenizationDocument } from "@/redux/api/sharedstakeholders";

interface Props {
  title: string;
  button?: string;
  status?: string;
  updatedAt?: string;
  verifiedDate?: string;
  verifiedBy?: string;
  checklistItems?: IDueDiligenceItem[];
  documents?: IStakeholderAssetTokenizationDocument[];
  onVerify?: () => void;
}

const STATUS_COLORS: Record<string, { bg: string; color: string; border: string }> = {
  complete: { bg: "#ecfdf3", color: "#12b76a", border: "#abefc6" },
  approved: { bg: "#ecfdf3", color: "#12b76a", border: "#abefc6" },
  pending: { bg: "#f0f8ff", color: "#007cdf", border: "#b2ddff" },
  in_review: { bg: "#fffaeb", color: "#f79009", border: "#fedf89" },
  failed: { bg: "#fef3f2", color: "#be3800", border: "#fda29b" },
  rejected: { bg: "#fef3f2", color: "#be3800", border: "#fda29b" },
  not_applicable: { bg: "#f8f9fa", color: "#667085", border: "#d0d5dd" },
};

const formatStatusLabel = (statusStr?: string): string => {
  if (!statusStr) return "Pending";
  const s = statusStr.toLowerCase();
  if (s === "complete" || s === "approved") return "Complete";
  if (s === "pending") return "Pending";
  if (s === "in_review") return "In Review";
  if (s === "failed" || s === "rejected") return "Failed";
  if (s === "not_applicable") return "Not Applicable";
  return statusStr.replace(/_/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
};

const renderStatusIcon = (statusStr?: string) => {
  const s = statusStr?.toLowerCase();
  if (s === "complete" || s === "approved") {
    return <AiOutlineCheckCircle size={14} />;
  }
  if (s === "failed" || s === "rejected") {
    return <AiOutlineCloseCircle size={14} />;
  }
  if (s === "not_applicable") {
    return <AiOutlineMinusCircle size={14} />;
  }
  return <AiOutlineClockCircle size={14} />;
};

const formatDate = (dateStr?: string) => {
  if (!dateStr) return "";
  try {
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return dateStr;
    return d.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  } catch {
    return dateStr;
  }
};

const VerificationSection = ({
  title,
  button,
  status = "pending",
  updatedAt,
  verifiedBy,
  verifiedDate,
  checklistItems,
  documents = [],
  onVerify,
}: Props) => {
  const displayDate = formatDate(updatedAt || verifiedDate);
  const statusLabel = formatStatusLabel(status);
  const isPending =
    !status ||
    status.toLowerCase() === "pending" ||
    status.toLowerCase() === "in_review";
  const buttonText = button || (isPending ? "Mark as Completed" : statusLabel);
  const isComplete = ["complete", "approved"].includes(status.toLowerCase());

  const openDocument = (document: IStakeholderAssetTokenizationDocument) => {
    if (document.url) {
      window.open(document.url, "_blank", "noopener,noreferrer");
    }
  };

  return (
    <Section>
      <SectionHeader>
        <RightSection>
          <SectionTitle>{title}</SectionTitle>

          <MetaRow>
            <StatusBadge $status={status}>
              {renderStatusIcon(status)}
              <span>{statusLabel}</span>
            </StatusBadge>

            {displayDate && (
              <DateWrap>
                <FiCalendar size={12} />
                <span>Updated on {displayDate}</span>
              </DateWrap>
            )}

            {verifiedBy && (
              <VerifierWrap>
                <Avatar>{verifiedBy.charAt(0).toUpperCase()}</Avatar>
                <VerifierText>By {verifiedBy}</VerifierText>
              </VerifierWrap>
            )}
          </MetaRow>
        </RightSection>

        <SeeMore
          onClick={isComplete ? undefined : onVerify}
          $status={status}
          disabled={isComplete}
          aria-label={isComplete ? `${title} is completed` : buttonText}
        >
          {renderStatusIcon(status)} <span>{buttonText}</span>
        </SeeMore>
      </SectionHeader>

      <CardGrid>
        {documents.map((document) => (
          <DocumentCard
            key={document.id}
            title={document.title || document.document_type}
            fileUrl={document.url}
            onClick={document.url ? () => openDocument(document) : undefined}
          />
        ))}
        {documents.length === 0 && (
          <EmptyDocuments>No tokenized asset documents available.</EmptyDocuments>
        )}
      </CardGrid>
    </Section>
  );
};

export default VerificationSection;

const Section = styled.div`
  background: white;
  border-radius: 10px;
  padding: 18px;
  margin-bottom: 16px;
`;

const SectionHeader = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  align-items: center;
  gap: 4px;
`;

const SectionTitle = styled.h4`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
  margin: 0 0 4px 0;
`;

const MetaRow = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
`;

const StatusBadge = styled.div<{ $status?: string }>`
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 500;
  background: ${({ $status }) =>
    STATUS_COLORS[$status?.toLowerCase() || "pending"]?.bg || "#f0f8ff"};
  color: ${({ $status }) =>
    STATUS_COLORS[$status?.toLowerCase() || "pending"]?.color || "#007cdf"};
  border: 1px solid
    ${({ $status }) =>
      STATUS_COLORS[$status?.toLowerCase() || "pending"]?.border || "#b2ddff"};
`;

const DateWrap = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  color: #667085;
  font-size: 11px;
`;

const VerifierWrap = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

const Avatar = styled.div`
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #f2f4f7;
  color: #344054;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: 600;
`;

const VerifierText = styled.span`
  color: #667085;
  font-size: 11px;
`;

const SeeMore = styled.button<{ $status?: string }>`
  border: 1px solid
    ${({ $status }) =>
      STATUS_COLORS[$status?.toLowerCase() || "pending"]?.border || "#007cdf"};
  background: #fff;
  padding: 6px 12px;
  color: ${({ $status }) =>
    STATUS_COLORS[$status?.toLowerCase() || "pending"]?.color || "#007cdf"};
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;
  display: flex;
  gap: 6px;
  align-items: center;
  transition: all 0.2s ease;

  &:hover {
    background: ${({ $status }) =>
      STATUS_COLORS[$status?.toLowerCase() || "pending"]?.bg || "#f0f8ff"};
  }

  &:disabled {
    cursor: not-allowed;
    opacity: 0.7;
  }
`;

const CardGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 14px;
  border: 1px solid #e0e0e0;
  border-radius: 12px;
  padding: 24px;
`;

const EmptyDocuments = styled.p`
  grid-column: 1 / -1;
  margin: 0;
  color: #667085;
  font-size: 12px;
  text-align: center;
`;

const RightSection = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;
