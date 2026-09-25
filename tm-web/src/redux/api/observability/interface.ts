export interface IDependency {
  name: string;
  status: "up" | "down" | "degraded" | "skipped";
  latency_ms: number;
  error?: string;
  optional?: boolean;
}

/**
 * Reported only by worker services. The payment history engine streams
 * operations from the blockchain, so "the process is up" says nothing about
 * whether it is still writing payment history - this does.
 */
export interface IStream {
  status: "up" | "stalled" | "skipped";
  operations_processed: number;
  seconds_since_last_operation: number;
  last_operation_at?: string;
  note?: string;
}

export interface IServiceHealth {
  name: string;
  status: "up" | "down" | "degraded";
  version: string;
  uptime_seconds: number;
  dependencies: IDependency[];
  /** Present only for worker services. */
  stream?: IStream;
  error?: string;
}

export interface IHostHealth {
  cpu_percent: number | null;
  memory_percent: number | null;
  disk_percent: number | null;
  available: boolean;
}

export interface IHealthResponse {
  services: IServiceHealth[];
  host: IHostHealth;
  checked_at: string;
  grafana_url?: string;
}

export interface IIssue {
  id: string;
  title: string;
  culprit?: string;
  level: string;
  status: string;
  count: number;
  app: string;
  first_seen: string;
  last_seen: string;
  permalink?: string;
}

export interface IIssuesResponse {
  issues: IIssue[];
  available: boolean;
  reason?: string;
  glitchtip_url?: string;
}

export interface ISeriesPoint {
  time: string;
  value: number;
}

export interface ISeries {
  name: string;
  points: ISeriesPoint[];
}

export interface ITrendsResponse {
  request_rate: ISeries[];
  error_rate: ISeries[];
  handled_failures: ISeries[];
  latency_p95: ISeries[];
  window: string;
  available: boolean;
  reason?: string;
}

export interface ILogLine {
  time: string;
  service: string;
  level?: string;
  method?: string;
  path?: string;
  status?: number;
  duration_ms?: number;
  message: string;
}

export interface ITraceResponse {
  request_id: string;
  lines: ILogLine[];
  available: boolean;
  reason?: string;
}

export interface IRecentErrorsResponse {
  lines: ILogLine[];
  available: boolean;
  reason?: string;
}

export interface ICapabilities {
  health: boolean;
  issues: boolean;
  trends: boolean;
  trace: boolean;
  grafana_url?: string;
  glitchtip_url?: string;
}

/** The admin API wraps every observability payload in a data envelope. */
export interface IEnvelope<T> {
  data: T;
}
