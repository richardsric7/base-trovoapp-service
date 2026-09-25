"use client";
import { useState } from "react";
import styled from "styled-components";
import {
  useGetRecentErrorsQuery,
  useLazyGetRequestTraceQuery,
} from "@/redux/api/observability/api";
import { ILogLine } from "@/redux/api/observability/interface";
import { Loading, PanelHint, PanelTitle, Unavailable, formatLatency, relativeTime } from "./shared";

/** A request id looks like a UUID; anything else the backend will reject anyway. */
const looksLikeRequestId = (value: string) => /^[A-Za-z0-9._-]{1,64}$/.test(value.trim());

/**
 * Follow one request across every service that handled it.
 *
 * This is what the request-id plumbing was for: a user reports a problem and
 * quotes the reference from the error screen, and this turns it into the exact
 * sequence of log lines behind it — including the services the request passed
 * through on its way.
 */
const TraceTab = () => {
  const [input, setInput] = useState("");
  const [submitted, setSubmitted] = useState("");
  const [runTrace, traceResult] = useLazyGetRequestTraceQuery();
  const recent = useGetRecentErrorsQuery({ limit: 25 });

  const search = (value: string) => {
    const id = value.trim();
    if (!looksLikeRequestId(id)) return;
    setInput(id);
    setSubmitted(id);
    runTrace({ requestId: id });
  };

  const trace = traceResult.data?.data;
  const invalid = input.trim() !== "" && !looksLikeRequestId(input);

  return (
    <Stack>
      <div>
        <PanelTitle>Trace a request</PanelTitle>
        <PanelHint>
          Every API response carries an <code>X-Request-ID</code> header, and the web apps show
          it on their error screens. Paste one here to see everything that happened, across all
          services, in order.
        </PanelHint>
      </div>

      <SearchRow
        onSubmit={(e) => {
          e.preventDefault();
          search(input);
        }}
      >
        <Input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="e.g. f3bb087a-d3a0-41f0-871d-c9e0ce2f5142"
          aria-label="Request identifier"
          $invalid={invalid}
        />
        <SearchButton type="submit" disabled={!looksLikeRequestId(input)}>
          Trace
        </SearchButton>
      </SearchRow>
      {invalid && <InvalidNote>That does not look like a request identifier.</InvalidNote>}

      {traceResult.isFetching && <Loading label="Searching the logs…" />}

      {!traceResult.isFetching && submitted && trace && !trace.available && (
        <Unavailable reason={trace.reason} />
      )}

      {!traceResult.isFetching && trace?.available && (
        <>
          <ResultHeader>
            {trace.lines.length > 0
              ? `${trace.lines.length} log ${trace.lines.length === 1 ? "line" : "lines"} for ${submitted}`
              : `No log lines found for ${submitted}`}
          </ResultHeader>
          {trace.lines.length === 0 ? (
            <EmptyBox>
              <strong>Nothing found</strong>
              <span>
                Either the identifier is from outside the log retention window, or it was never
                seen by a service that ships logs.
              </span>
            </EmptyBox>
          ) : (
            <LogList lines={trace.lines} />
          )}
        </>
      )}

      <Divider />

      <div>
        <PanelTitle>Recent failed requests</PanelTitle>
        <PanelHint>
          The latest requests that returned an error, across every service. Click one to trace it.
        </PanelHint>
      </div>

      {recent.isLoading ? (
        <Loading />
      ) : !recent.data?.data?.available ? (
        <Unavailable reason={recent.data?.data?.reason} />
      ) : recent.data.data.lines.length === 0 ? (
        <EmptyBox>
          <strong>No failed requests in the last hour</strong>
          <span>Nothing has returned an error recently.</span>
        </EmptyBox>
      ) : (
        <RecentList>
          {recent.data.data.lines.map((line, i) => (
            <RecentRow key={`${line.time}-${i}`}>
              <RecentLeft>
                <StatusCode status={line.status}>{line.status || "—"}</StatusCode>
                <RecentPath>
                  {line.method} {line.path || line.message}
                </RecentPath>
              </RecentLeft>
              <RecentRight>
                <ServiceTag>{line.service}</ServiceTag>
                <Muted>{relativeTime(line.time)}</Muted>
              </RecentRight>
            </RecentRow>
          ))}
        </RecentList>
      )}
    </Stack>
  );
};

const LogList = ({ lines }: { lines: ILogLine[] }) => (
  <LogWrap>
    {lines.map((line, i) => (
      <LogRow key={`${line.time}-${i}`} level={line.level}>
        <LogTime>{new Date(line.time).toLocaleTimeString()}</LogTime>
        <ServiceTag>{line.service}</ServiceTag>
        <LogBody>
          {line.method && line.path ? (
            <>
              <strong>
                {line.method} {line.path}
              </strong>
              {line.status ? <StatusCode status={line.status}>{line.status}</StatusCode> : null}
              {line.duration_ms ? <Muted>{formatLatency(line.duration_ms)}</Muted> : null}
            </>
          ) : (
            <RawMessage>{line.message}</RawMessage>
          )}
        </LogBody>
      </LogRow>
    ))}
  </LogWrap>
);

export default TraceTab;

const Stack = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

const SearchRow = styled.form`
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
`;

const Input = styled.input<{ $invalid?: boolean }>`
  flex: 1;
  min-width: 280px;
  padding: 10px 16px;
  font-size: 13px;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  color: #00225a;
  border: 1px solid ${({ $invalid }) => ($invalid ? "#BE3800" : "#e5e5ef")};
  border-radius: 8px;
  outline: none;

  &:focus {
    border-color: #00225a;
  }
`;

const SearchButton = styled.button`
  padding: 10px 24px;
  font-size: 14px;
  font-family: inherit;
  font-weight: 500;
  color: #ffffff;
  background: #00225a;
  border: none;
  border-radius: 8px;
  cursor: pointer;

  &:disabled {
    background: #e5e5ef;
    cursor: not-allowed;
  }
`;

const InvalidNote = styled.p`
  font-size: 13px;
  color: #be3800;
  margin: -8px 0 0;
`;

const ResultHeader = styled.p`
  font-size: 13px;
  color: #5b6b8c;
  margin: 0;
`;

const LogWrap = styled.div`
  border: 1px solid #e5e5ef;
  border-radius: 12px;
  overflow: hidden;
`;

const LogRow = styled.div<{ level?: string }>`
  display: flex;
  align-items: baseline;
  gap: 12px;
  padding: 10px 16px;
  border-bottom: 1px solid #e5e5ef;
  background: ${({ level }) => (level === "error" ? "#be38000a" : "#ffffff")};

  &:last-child {
    border-bottom: none;
  }
`;

const LogTime = styled.span`
  font-size: 12px;
  color: #a0a8bb;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  white-space: nowrap;
`;

const ServiceTag = styled.span`
  display: inline-block;
  padding: 2px 8px;
  font-size: 11px;
  border-radius: 5px;
  background: #f4f6f9;
  color: #00225a;
  white-space: nowrap;
`;

const LogBody = styled.div`
  display: flex;
  align-items: baseline;
  gap: 10px;
  flex-wrap: wrap;
  font-size: 13px;
  color: #00225a;
  min-width: 0;

  strong {
    font-weight: 500;
    font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  }
`;

const RawMessage = styled.span`
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 12px;
  color: #5b6b8c;
  word-break: break-all;
`;

const StatusCode = styled.span<{ status?: number }>`
  display: inline-block;
  padding: 1px 8px;
  border-radius: 5px;
  font-size: 12px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: ${({ status }) => (!status ? "#6B7280" : status >= 500 ? "#BE3800" : status >= 400 ? "#B26A00" : "#00A859")};
  background: ${({ status }) => (!status ? "#E5E7EB" : status >= 500 ? "#BE38001A" : status >= 400 ? "#B26A001A" : "#00A8591A")};
`;

const Muted = styled.span`
  font-size: 12px;
  color: #828282;
  white-space: nowrap;
`;

const Divider = styled.hr`
  border: none;
  border-top: 1px solid #e5e5ef;
  margin: 8px 0;
`;

const RecentList = styled.div`
  border: 1px solid #e5e5ef;
  border-radius: 12px;
  overflow: hidden;
`;

const RecentRow = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 12px 16px;
  border-bottom: 1px solid #e5e5ef;

  &:last-child {
    border-bottom: none;
  }
`;

const RecentLeft = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
`;

const RecentPath = styled.span`
  font-size: 13px;
  color: #00225a;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
`;

const RecentRight = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  white-space: nowrap;
`;

const EmptyBox = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 32px;
  text-align: center;
  border: 1px dashed #e5e5ef;
  border-radius: 12px;

  strong {
    font-size: 14px;
    color: #00225a;
  }
  span {
    font-size: 13px;
    color: #5b6b8c;
    line-height: 1.5;
  }
`;
