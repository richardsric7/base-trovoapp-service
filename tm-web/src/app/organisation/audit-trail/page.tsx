"use client";

import {
  IStakeholderAuditTrailRecord,
  useGetStakeholderAuditTrailQuery,
} from "@/redux/api/sharedstakeholders";
import { format, formatDistanceToNow } from "date-fns";
import { useState } from "react";
import { FiMoreHorizontal } from "react-icons/fi";
import styled, { keyframes } from "styled-components";

const titleCase = (value: string) =>
  value
    .replace(/[._-]+/g, " ")
    .replace(/\b\w/g, (character) => character.toUpperCase());

const getActionSummary = (record: IStakeholderAuditTrailRecord) => {
  const [resource, operation] = record.action.split(".");
  const subject = titleCase(record.entity_type || resource || "record");
  const pastTenseActions: Record<string, string> = {
    approve: "approved",
    authorize: "authorized",
    create: "created",
    delete: "deleted",
    execute: "executed",
    generate: "generated",
    record: "recorded",
    reject: "rejected",
    submit: "submitted",
    update: "updated",
  };
  const normalizedOperation = (
    operation ||
    record.action ||
    "update"
  ).toLowerCase();
  const action =
    pastTenseActions[normalizedOperation] ||
    titleCase(normalizedOperation).toLowerCase();

  return `${subject} was ${action}.`;
};

const formatTimestamp = (value: string) => {
  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return { relative: "Date unavailable", absolute: value };
  }

  return {
    relative: formatDistanceToNow(date, { addSuffix: true }),
    absolute: format(date, "dd MMM, yyyy h:mm a"),
  };
};

const formatDetailValue = (key: string, value: unknown) => {
  if (typeof value === "string" && key.toLowerCase().endsWith("_at")) {
    const date = new Date(value);

    if (!Number.isNaN(date.getTime())) {
      return format(date, "dd MMM, yyyy h:mm a");
    }
  }

  if (typeof value === "object" && value !== null) {
    return JSON.stringify(value);
  }

  return String(value ?? "—");
};

export default function StakeholderAuditTrailPage() {
  const { data, isLoading, isFetching, isError, refetch } =
    useGetStakeholderAuditTrailQuery();
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const records = data?.data.records ?? [];

  return (
    <PageContainer>
      <Tabs aria-label="Audit trail views">
        <ActiveTab type="button">General</ActiveTab>
      </Tabs>

      <Feed aria-busy={isLoading || isFetching}>
        {isLoading ? (
          Array.from({ length: 6 }).map((_, index) => (
            <SkeletonRow key={index} aria-hidden="true">
              <SkeletonAvatar />
              <SkeletonText />
            </SkeletonRow>
          ))
        ) : isError ? (
          <MessageState role="alert">
            <MessageTitle>We couldn&apos;t load the audit trail</MessageTitle>
            <MessageText>
              Please check your connection and try again.
            </MessageText>
            <RetryButton type="button" onClick={() => refetch()}>
              Try again
            </RetryButton>
          </MessageState>
        ) : records.length === 0 ? (
          <MessageState>
            <MessageTitle>No audit activity yet</MessageTitle>
            <MessageText>
              New stakeholder activity will appear here.
            </MessageText>
          </MessageState>
        ) : (
          records.map((record) => {
            const timestamp = formatTimestamp(record.created_at);
            const expanded = expandedId === record.id;

            return (
              <AuditItem key={record.id}>
                <AuditRow>
                  <Avatar aria-hidden="true">
                    {record.actor_role?.charAt(0).toUpperCase() || "A"}
                  </Avatar>
                  <AuditContent>
                    <Summary>{getActionSummary(record)}</Summary>
                    <Metadata>
                      <span>
                        {titleCase(record.actor_role || "Stakeholder")}
                      </span>
                      <span aria-hidden="true">•</span>
                      <time
                        dateTime={record.created_at}
                        title={timestamp.absolute}
                      >
                        {timestamp.relative}
                      </time>
                    </Metadata>
                  </AuditContent>
                  <MoreButton
                    type="button"
                    aria-label={`${expanded ? "Hide" : "Show"} audit details`}
                    aria-expanded={expanded}
                    aria-controls={`audit-details-${record.id}`}
                    onClick={() => setExpandedId(expanded ? null : record.id)}
                  >
                    <FiMoreHorizontal size={18} />
                  </MoreButton>
                </AuditRow>

                {expanded && (
                  <Details id={`audit-details-${record.id}`}>
                    <DetailRow>
                      <DetailLabel>Entity</DetailLabel>
                      <DetailValue>{titleCase(record.entity_type)}</DetailValue>
                    </DetailRow>
                    <DetailRow>
                      <DetailLabel>Reference ID</DetailLabel>
                      <DetailValue>{record.entity_id}</DetailValue>
                    </DetailRow>
                    <DetailRow>
                      <DetailLabel>Date</DetailLabel>
                      <DetailValue>{timestamp.absolute}</DetailValue>
                    </DetailRow>
                    {record.after_state &&
                      Object.entries(record.after_state).map(([key, value]) => (
                        <DetailRow key={key}>
                          <DetailLabel>{titleCase(key)}</DetailLabel>
                          <DetailValue>{formatDetailValue(key, value)}</DetailValue>
                        </DetailRow>
                      ))}
                  </Details>
                )}
              </AuditItem>
            );
          })
        )}
      </Feed>
    </PageContainer>
  );
}

const shimmer = keyframes`
  from { background-position: 100% 0; }
  to { background-position: -100% 0; }
`;

const PageContainer = styled.section`
  padding: 0 24px 40px;
  color: #00225a;

  @media (max-width: 768px) {
    padding: 0 12px 28px;
  }
`;

const Tabs = styled.div`
  width: min(130px, 100%);
  padding: 4px;
  margin-bottom: 22px;
  border-radius: 12px;
  background: #ffffff;
`;

const ActiveTab = styled.button`
  width: 100%;
  padding: 8px 18px;
  border: 0;
  border-radius: 9px;
  background: #f2f6f9;
  color: #00225a;
  font: inherit;
  font-size: 13px;
  cursor: default;
`;

const Feed = styled.div`
  padding: 18px;
  border-radius: 22px;
  background: #ffffff;

  @media (max-width: 768px) {
    padding: 10px;
    border-radius: 16px;
  }
`;

const AuditItem = styled.article`
  overflow: hidden;
  margin-bottom: 8px;
  border-radius: 13px;
  background: #f2f6f9;

  &:last-child {
    margin-bottom: 0;
  }
`;

const AuditRow = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 62px;
  padding: 10px 14px;
`;

const Avatar = styled.div`
  display: grid;
  flex: 0 0 42px;
  width: 42px;
  height: 42px;
  place-items: center;
  border-radius: 50%;
  background: #f2d3ad;
  color: #9a673c;
  font-weight: 600;
`;

const AuditContent = styled.div`
  flex: 1;
  min-width: 0;
`;

const Summary = styled.p`
  margin: 0 0 3px;
  color: #007cdf;
  font-size: 14px;
  line-height: 1.35;
`;

const Metadata = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  color: #8a94a6;
  font-size: 11px;
`;

const MoreButton = styled.button`
  display: grid;
  flex: 0 0 32px;
  width: 32px;
  height: 32px;
  place-items: center;
  border: 1px solid #007cdf;
  border-radius: 9px;
  background: #ffffff;
  color: #007cdf;
  cursor: pointer;

  &:hover,
  &:focus-visible {
    background: #eaf5ff;
  }
`;

const Details = styled.div`
  padding: 12px 68px 16px;
  border-top: 1px solid #dce8f1;
  background: #ffffff;

  @media (max-width: 768px) {
    padding: 14px;
  }
`;

const DetailRow = styled.div`
  display: grid;
  grid-template-columns: minmax(90px, 130px) 1fr;
  gap: 16px;
  padding: 7px 0;

  & + & {
    border-top: 1px solid #edf1f4;
  }
`;
const DetailLabel = styled.span`
  color: #828282;
  font-size: 12px;
`;
const DetailValue = styled.span`
  overflow-wrap: anywhere;
  color: #00225a;
  font-size: 12px;
  font-weight: 500;
`;

const MessageState = styled.div`
  padding: 64px 20px;
  text-align: center;
`;
const MessageTitle = styled.h2`
  margin: 0 0 8px;
  color: #00225a;
  font-size: 17px;
`;
const MessageText = styled.p`
  margin: 0;
  color: #828282;
  font-size: 13px;
`;
const RetryButton = styled.button`
  margin-top: 18px;
  padding: 9px 20px;
  border: 0;
  border-radius: 8px;
  background: #007cdf;
  color: #fff;
  font: inherit;
  cursor: pointer;
`;

const SkeletonRow = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 62px;
  margin-bottom: 8px;
  padding: 10px 14px;
  border-radius: 13px;
  background: #f2f6f9;
`;
const SkeletonAvatar = styled.div`
  width: 42px;
  height: 42px;
  flex: 0 0 42px;
  border-radius: 50%;
  background: linear-gradient(90deg, #e3eaf0 25%, #f2f6f9 50%, #e3eaf0 75%);
  background-size: 200% 100%;
  animation: ${shimmer} 1.4s infinite;
`;
const SkeletonText = styled.div`
  width: min(520px, 70%);
  height: 24px;
  border-radius: 6px;
  background: linear-gradient(90deg, #e3eaf0 25%, #f8fafc 50%, #e3eaf0 75%);
  background-size: 200% 100%;
  animation: ${shimmer} 1.4s infinite;
`;
