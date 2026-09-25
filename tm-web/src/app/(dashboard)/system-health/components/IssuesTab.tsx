"use client";
import { useState } from "react";
import styled from "styled-components";
import { AiOutlineCheckCircle, AiOutlineCloseCircle } from "react-icons/ai";
import CustomTable from "@/components/CustomTable";
import { useGetSystemIssuesQuery } from "@/redux/api/observability/api";
import { IIssue } from "@/redux/api/observability/interface";
import { Loading, PanelHint, PanelTitle, Unavailable, relativeTime } from "./shared";

/**
 * Grouped crash reports from every Trovo application.
 *
 * Grouped is the important word: a hundred occurrences of one bug is a single
 * row with a count, not a hundred rows. Sorted by most recent activity, which
 * is what someone asking "what is broken now" wants.
 *
 * Uses the dashboard's CustomTable so the columns, pagination and empty state
 * behave the same as every other table in Trovo Manager.
 */
const IssuesTab = () => {
  const [app, setApp] = useState("");
  const { data, isLoading, isError } = useGetSystemIssuesQuery(
    { app: app || undefined, limit: 50 },
    { pollingInterval: 60000 }
  );

  if (isLoading) return <Loading label="Fetching recent errors…" />;
  if (isError) return <Unavailable reason="Could not reach the admin API to fetch errors." />;

  const payload = data?.data;
  if (!payload?.available) return <Unavailable reason={payload?.reason} />;

  const apps = Array.from(new Set(payload.issues.map((i) => i.app))).sort();

  const columns = [
    {
      title: "Error",
      dataIndex: "title",
      key: "title",
      render: (_: string, record: IIssue) => (
        <ErrorCell>
          <ErrorTitle
            href={record.permalink || payload.glitchtip_url}
            target="_blank"
            rel="noopener noreferrer"
          >
            {record.title}
          </ErrorTitle>
          {record.culprit && <Culprit>{record.culprit}</Culprit>}
        </ErrorCell>
      ),
    },
    {
      title: "Application",
      dataIndex: "app",
      key: "app",
      render: (value: string) => <AppTag>{value}</AppTag>,
    },
    {
      title: "Events",
      dataIndex: "count",
      key: "count",
      render: (value: number) => <Count>{value.toLocaleString()}</Count>,
    },
    {
      title: "Last seen",
      dataIndex: "last_seen",
      key: "last_seen",
      render: (_: string, record: IIssue) => (
        <div>
          <LastSeen title={record.last_seen}>{relativeTime(record.last_seen)}</LastSeen>
          <FirstSeen>first {relativeTime(record.first_seen)}</FirstSeen>
        </div>
      ),
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (status: string) => {
        const isResolved = status === "resolved";
        return (
          <IssueStatus isResolved={isResolved}>
            {isResolved ? <AiOutlineCheckCircle size={18} /> : <AiOutlineCloseCircle size={18} />}
            {status}
          </IssueStatus>
        );
      },
    },
  ];

  return (
    <Stack>
      <Header>
        <div>
          <PanelTitle>Application errors</PanelTitle>
          <PanelHint>
            One row per distinct error, with how many times it has happened. Errors that stop
            appearing here have stopped happening.
          </PanelHint>
        </div>
        {apps.length > 1 && (
          <Select value={app} onChange={(e) => setApp(e.target.value)}>
            <option value="">All applications</option>
            {apps.map((a) => (
              <option key={a} value={a}>
                {a}
              </option>
            ))}
          </Select>
        )}
      </Header>

      <CustomTable
        columns={columns}
        dataSource={payload.issues}
        totalItems={payload.issues.length}
        pageSize={10}
        isLoading={false}
      />

      {payload.glitchtip_url && (
        <FooterLink href={payload.glitchtip_url} target="_blank" rel="noopener noreferrer">
          Open the crash reporter for stack traces and history →
        </FooterLink>
      )}
    </Stack>
  );
};

export default IssuesTab;

const Stack = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
`;

const Header = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
  flex-wrap: wrap;
`;

const Select = styled.select`
  padding: 8px 14px;
  font-size: 14px;
  font-family: inherit;
  color: #00225a;
  background: #ffffff;
  border: 1px solid #e5e5ef;
  border-radius: 8px;
  cursor: pointer;
`;

const ErrorCell = styled.div`
  max-width: 380px;
`;

const ErrorTitle = styled.a`
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  text-decoration: none;
  word-break: break-word;

  &:hover {
    text-decoration: underline;
  }
`;

const Culprit = styled.span`
  display: block;
  font-size: 12px;
  color: #828282;
  margin-top: 3px;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
`;

const AppTag = styled.span`
  display: inline-block;
  padding: 4px 10px;
  font-size: 13px;
  border-radius: 8px;
  background: #f4f6f9;
  color: #00225a;
  white-space: nowrap;
`;

const Count = styled.span`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  font-variant-numeric: tabular-nums;
`;

const LastSeen = styled.span`
  display: block;
  font-size: 14px;
  color: #00225a;
  white-space: nowrap;
`;

const FirstSeen = styled.span`
  display: block;
  font-size: 12px;
  color: #828282;
  margin-top: 2px;
  white-space: nowrap;
`;

// Matches the Audit Trail table's status chip so a resolved/unresolved error
// reads the same way as a successful/failed login elsewhere in the dashboard.
const IssueStatus = styled.p<{ isResolved: boolean }>`
  color: ${({ isResolved }) => (isResolved ? "#00A859" : "#BE3800")};
  background-color: ${({ isResolved }) => (isResolved ? "#00A8591A" : "#BE38001A")};
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border-radius: 8px;
  text-transform: capitalize;
  margin: 0;
  white-space: nowrap;
`;

const FooterLink = styled.a`
  font-size: 14px;
  color: #00225a;
  text-decoration: underline;
  font-weight: 500;
`;
