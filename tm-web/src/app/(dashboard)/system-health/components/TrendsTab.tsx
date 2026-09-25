"use client";
import { useState } from "react";
import dynamic from "next/dynamic";
import styled from "styled-components";
import { useGetSystemTrendsQuery } from "@/redux/api/observability/api";
import { ISeries } from "@/redux/api/observability/interface";
import { Loading, Panel, PanelHint, PanelTitle, Unavailable } from "./shared";

// ApexCharts touches window on import, so it cannot render on the server.
const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

const WINDOWS = [
  { key: "1h", label: "Last hour" },
  { key: "6h", label: "Last 6 hours" },
  { key: "24h", label: "Last 24 hours" },
];

/** One colour per service, stable across all three charts so the eye can follow a service down the page. */
const PALETTE = ["#00225A", "#00A859", "#B26A00", "#7B2E86", "#BE3800"];

/**
 * Traffic and performance over time.
 *
 * The 95th percentile is charted rather than an average because an average
 * hides exactly the requests users complain about: if one request in twenty
 * takes four seconds, the mean stays comfortable and the p95 does not.
 */
const TrendsTab = () => {
  const [window, setWindow] = useState("6h");
  const { data, isLoading, isError } = useGetSystemTrendsQuery({ window });

  if (isLoading) return <Loading label="Loading metrics…" />;
  if (isError) return <Unavailable reason="Could not reach the admin API to load metrics." />;

  const payload = data?.data;
  if (!payload?.available) return <Unavailable reason={payload?.reason} />;

  const hasData =
    payload.request_rate.some((s) => s.points.length > 0) ||
    payload.error_rate.some((s) => s.points.length > 0) ||
    (payload.handled_failures ?? []).some((s) => s.points.length > 0);

  return (
    <Stack>
      <Header>
        <div>
          <PanelTitle>Trends</PanelTitle>
          <PanelHint>
            Traffic and responsiveness per service. Response time is the 95th percentile — the
            experience of the slowest one request in twenty, which an average would hide.
          </PanelHint>
        </div>
        <WindowGroup>
          {WINDOWS.map((w) => (
            <WindowButton key={w.key} $active={window === w.key} onClick={() => setWindow(w.key)}>
              {w.label}
            </WindowButton>
          ))}
        </WindowGroup>
      </Header>

      {!hasData ? (
        <Unavailable reason="No metrics have been recorded in this window yet." />
      ) : (
        <>
          <ChartPanel
            title="Requests per second"
            hint="How much traffic each service is handling."
            series={payload.request_rate}
            unit=""
            decimals={2}
          />
          <ChartPanel
            title="Error rate"
            hint="Percentage of requests failing with a server error. HTTP only — background work is in the panel below."
            series={payload.error_rate}
            unit="%"
            decimals={2}
          />
          <ChartPanel
            title="Handled failures per minute"
            hint="Failures the services caught rather than returning as a server error — most of them in background jobs that serve no request, so the error rate above cannot see them."
            series={payload.handled_failures ?? []}
            unit="/min"
            decimals={2}
          />
          <ChartPanel
            title="Response time (95th percentile)"
            hint="How slow the slowest one in twenty requests is."
            series={payload.latency_p95}
            unit="s"
            decimals={3}
          />
        </>
      )}
    </Stack>
  );
};

const ChartPanel = ({
  title,
  hint,
  series,
  unit,
  decimals,
}: {
  title: string;
  hint: string;
  series: ISeries[];
  unit: string;
  decimals: number;
}) => {
  const withData = series.filter((s) => s.points.length > 0);

  return (
    <Panel>
      <PanelTitle>{title}</PanelTitle>
      <PanelHint>{hint}</PanelHint>
      {withData.length === 0 ? (
        <NoSeries>No data for this metric in the selected window.</NoSeries>
      ) : (
        <Chart
          type="area"
          height={260}
          series={withData.map((s) => ({
            name: s.name,
            data: s.points.map((p) => ({ x: new Date(p.time).getTime(), y: Number(p.value.toFixed(decimals)) })),
          }))}
          options={{
            chart: { toolbar: { show: false }, zoom: { enabled: false }, fontFamily: "inherit" },
            colors: PALETTE,
            dataLabels: { enabled: false },
            stroke: { curve: "smooth", width: 2 },
            fill: { type: "gradient", gradient: { opacityFrom: 0.25, opacityTo: 0.02 } },
            xaxis: {
              type: "datetime",
              labels: { style: { colors: "#828282", fontSize: "11px" }, datetimeUTC: false },
              axisBorder: { show: false },
              axisTicks: { show: false },
            },
            yaxis: {
              labels: {
                style: { colors: "#828282", fontSize: "11px" },
                formatter: (v: number) => `${v?.toFixed(decimals)}${unit}`,
              },
              min: 0,
            },
            grid: { borderColor: "#e5e5ef", strokeDashArray: 4 },
            legend: { position: "bottom", horizontalAlign: "left", fontSize: "12px" },
            tooltip: {
              x: { format: "dd MMM HH:mm" },
              y: { formatter: (v: number) => `${v?.toFixed(decimals)}${unit}` },
            },
          }}
        />
      )}
    </Panel>
  );
};

export default TrendsTab;

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

const WindowGroup = styled.div`
  display: inline-flex;
  border: 1px solid #e5e5ef;
  border-radius: 8px;
  overflow: hidden;
`;

const WindowButton = styled.button<{ $active: boolean }>`
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  border: none;
  cursor: pointer;
  white-space: nowrap;
  color: ${({ $active }) => ($active ? "#ffffff" : "#828282")};
  background: ${({ $active }) => ($active ? "#00225a" : "#ffffff")};

  & + & {
    border-left: 1px solid #e5e5ef;
  }
`;

const NoSeries = styled.p`
  font-size: 14px;
  color: #828282;
  text-align: center;
  padding: 40px 0;
  margin: 0;
`;
