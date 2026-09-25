"use client";
import styled from "styled-components";
import { useGetSystemHealthQuery } from "@/redux/api/observability/api";
import { IServiceHealth, IStream } from "@/redux/api/observability/interface";
import {
  Grid,
  Loading,
  Panel,
  PanelHint,
  PanelTitle,
  StatusPill,
  Unavailable,
  formatLatency,
  formatTime,
  formatUptime,
} from "./shared";

/**
 * Live status of every backend and the host they share.
 *
 * Polls every 30 seconds: fast enough that the screen is trustworthy during an
 * incident, slow enough not to add meaningful load. The backend caches for 15s,
 * so several people watching this page do not multiply into upstream queries.
 */
const HealthTab = () => {
  const { data, isLoading, isError, refetch } = useGetSystemHealthQuery(undefined, {
    pollingInterval: 30000,
  });

  if (isLoading) return <Loading label="Checking every service…" />;
  if (isError || !data?.data) {
    return <Unavailable reason="Could not reach the admin API to check service health." />;
  }

  const { services, host, checked_at, grafana_url } = data.data;

  return (
    <Stack>
      <Header>
        <div>
          <PanelTitle>Services</PanelTitle>
          <PanelHint>
            Each service reports its own dependencies. A dependency marked optional can fail
            without taking the service down — it degrades instead.
          </PanelHint>
        </div>
        <Meta>
          <span>Checked {formatTime(checked_at)}</span>
          <RefreshButton onClick={() => refetch()}>Refresh</RefreshButton>
        </Meta>
      </Header>

      {services.length === 0 ? (
        <Unavailable reason="No services are configured for health checks in this environment." />
      ) : (
        <Grid min="340px">
          {services.map((service) => (
            <ServiceCard key={service.name} service={service} />
          ))}
        </Grid>
      )}

      <Panel>
        <PanelTitle>Host</PanelTitle>
        <PanelHint>
          The machine every service runs on. Disk is worth watching — a full disk stops
          deployments before it stops the applications.
        </PanelHint>
        {!host.available ? (
          <Unavailable reason="Host metrics are not configured in this environment." />
        ) : (
          <Grid min="200px">
            <Metric label="CPU used" value={host.cpu_percent} warn={80} critical={90} />
            <Metric label="Memory used" value={host.memory_percent} warn={80} critical={90} />
            <Metric label="Disk used" value={host.disk_percent} warn={75} critical={85} />
          </Grid>
        )}
      </Panel>

      {grafana_url && (
        <FooterLink href={grafana_url} target="_blank" rel="noopener noreferrer">
          Open Grafana for full metrics and history →
        </FooterLink>
      )}
    </Stack>
  );
};

const ServiceCard = ({ service }: { service: IServiceHealth }) => (
  <Panel>
    <CardTop>
      <div>
        <ServiceName>{service.name}</ServiceName>
        <ServiceMeta>
          {service.version ? `build ${service.version.slice(0, 7)}` : "version unknown"}
          {service.uptime_seconds > 0 && ` · up ${formatUptime(service.uptime_seconds)}`}
        </ServiceMeta>
      </div>
      <StatusPill status={service.status}>{service.status}</StatusPill>
    </CardTop>

    {service.error && <ErrorNote>{service.error}</ErrorNote>}

    {service.stream && <StreamPanel stream={service.stream} />}

    {service.dependencies.length > 0 && (
      <DependencyList>
        {service.dependencies.map((dep) => (
          <DependencyRow key={dep.name}>
            <DependencyName>
              {dep.name.replace(/_/g, " ")}
              {dep.optional && <Optional>optional</Optional>}
            </DependencyName>
            <DependencyRight>
              {dep.status === "up" && <Latency>{formatLatency(dep.latency_ms)}</Latency>}
              <StatusPill status={dep.status}>{dep.status}</StatusPill>
            </DependencyRight>
          </DependencyRow>
        ))}
      </DependencyList>
    )}
  </Panel>
);

/**
 * Progress for a worker service.
 *
 * A background worker fails silently - a dead stream returns no errors to
 * anyone - so the count of processed operations and how long since the last one
 * are the only signals that say whether it is really working. Shown above the
 * dependency list because it matters more than any of them.
 */
const StreamPanel = ({ stream }: { stream: IStream }) => {
  const label =
    stream.status === "stalled"
      ? "Stalled"
      : stream.status === "skipped"
      ? "Starting up"
      : "Processing";

  return (
    <Stream status={stream.status}>
      <StreamTop>
        <StreamLabel>{label}</StreamLabel>
        <StatusPill status={stream.status === "stalled" ? "down" : stream.status === "skipped" ? "skipped" : "up"}>
          {stream.status === "stalled" ? "stalled" : stream.status === "skipped" ? "idle" : "live"}
        </StatusPill>
      </StreamTop>
      <StreamFigures>
        <span>
          <StreamNumber>{stream.operations_processed.toLocaleString()}</StreamNumber> operations
        </span>
        {stream.last_operation_at && (
          <StreamAgo>last {formatSince(stream.seconds_since_last_operation)}</StreamAgo>
        )}
      </StreamFigures>
      {stream.note && <StreamNote>{stream.note}</StreamNote>}
    </Stream>
  );
};

/** "12s ago" / "4m ago" - the unit that keeps the number readable. */
const formatSince = (seconds: number) => {
  if (seconds < 60) return `${seconds}s ago`;
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  return `${Math.floor(seconds / 3600)}h ago`;
};

const Metric = ({
  label,
  value,
  warn,
  critical,
}: {
  label: string;
  value: number | null;
  warn: number;
  critical: number;
}) => {
  // A null reading means the metric is genuinely absent. Showing 0% would be a
  // lie, and a reassuring one.
  const display = value === null ? "—" : `${value.toFixed(1)}%`;
  const tone = value === null ? "none" : value >= critical ? "critical" : value >= warn ? "warn" : "ok";

  return (
    <MetricBox>
      <MetricLabel>{label}</MetricLabel>
      <MetricValue tone={tone}>{display}</MetricValue>
      {value !== null && (
        <MetricBar>
          <MetricFill style={{ width: `${Math.min(value, 100)}%` }} tone={tone} />
        </MetricBar>
      )}
    </MetricBox>
  );
};

export default HealthTab;

const Stack = styled.div`
  display: flex;
  flex-direction: column;
  gap: 24px;
`;

const Header = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
  flex-wrap: wrap;
`;

const Meta = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  color: #828282;
  white-space: nowrap;
`;

const RefreshButton = styled.button`
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  color: #00225a;
  background: #ffffff;
  border: 1px solid #e5e5ef;
  border-radius: 8px;
  cursor: pointer;

  &:hover {
    background: #f4f6f9;
  }
`;

const CardTop = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
`;

const ServiceName = styled.h3`
  font-size: 16px;
  font-weight: 600;
  color: #00225a;
  margin: 0 0 4px;
`;

const ServiceMeta = styled.p`
  font-size: 12px;
  color: #828282;
  margin: 0;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
`;

const ErrorNote = styled.p`
  font-size: 12px;
  color: #be3800;
  background: #be38001a;
  padding: 8px 12px;
  border-radius: 8px;
  margin: 0 0 12px;
  word-break: break-word;
`;

const Stream = styled.div<{ status: string }>`
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 12px;
  background: ${({ status }) => (status === "stalled" ? "#BE38000D" : "#f4f6f9")};
  border: 1px solid ${({ status }) => (status === "stalled" ? "#BE380033" : "#e5e5ef")};
`;

const StreamTop = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
`;

const StreamLabel = styled.span`
  font-size: 13px;
  font-weight: 500;
  color: #00225a;
`;

const StreamFigures = styled.div`
  display: flex;
  align-items: baseline;
  gap: 12px;
  flex-wrap: wrap;
  font-size: 13px;
  color: #828282;
`;

const StreamNumber = styled.span`
  font-size: 15px;
  font-weight: 600;
  color: #00225a;
  font-variant-numeric: tabular-nums;
`;

const StreamAgo = styled.span`
  font-size: 13px;
  color: #828282;
  white-space: nowrap;
`;

const StreamNote = styled.p`
  font-size: 12px;
  color: #be3800;
  margin: 0;
  line-height: 1.45;
`;

const DependencyList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  border-top: 1px solid #e5e5ef;
  padding-top: 12px;
`;

const DependencyRow = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
`;

const DependencyName = styled.span`
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #00225a;
  text-transform: capitalize;
`;

const Optional = styled.span`
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.4px;
  color: #828282;
  border: 1px solid #e5e5ef;
  border-radius: 4px;
  padding: 1px 5px;
`;

const DependencyRight = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;

const Latency = styled.span`
  font-size: 13px;
  color: #828282;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
`;

const MetricBox = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const MetricLabel = styled.span`
  font-size: 14px;
  color: #828282;
`;

const MetricValue = styled.span<{ tone: string }>`
  font-size: 26px;
  font-weight: 700;
  color: ${({ tone }) =>
    tone === "critical" ? "#BE3800" : tone === "warn" ? "#B26A00" : tone === "ok" ? "#00225a" : "#828282"};
`;

const MetricBar = styled.div`
  height: 6px;
  border-radius: 999px;
  background: #f4f6f9;
  overflow: hidden;
`;

const MetricFill = styled.div<{ tone: string }>`
  height: 100%;
  border-radius: 999px;
  background: ${({ tone }) =>
    tone === "critical" ? "#BE3800" : tone === "warn" ? "#E0A33E" : "#00A859"};
  transition: width 0.4s ease;
`;

const FooterLink = styled.a`
  font-size: 14px;
  color: #00225a;
  text-decoration: underline;
  font-weight: 500;
`;
