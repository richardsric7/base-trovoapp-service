"use client";
import React from "react";
import dynamic from "next/dynamic";
import styled from "styled-components";
import Card from "../../../components/Card";
import disputesIcon from "@/assets/images/message-question.svg";
import successIcon from "@/assets/images/copy-success.svg";
import averageIcon from "@/assets/images/carbon_chart-average.svg";
import { ReportRange, useP2pDisputeReportQuery } from "@/redux/api/p2p";
import ReportSectionCard from "./ReportSectionCard";
import { CHART_COLORS, formatDateLabel, formatDuration, titleCase } from "./reportUtils";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

const DisputeReportSection = ({ range }: { range: ReportRange }) => {
  const { data, isLoading } = useP2pDisputeReportQuery({ range });
  const report = data?.data;

  const series = [
    { name: "Opened", data: report?.series?.map((p) => p.opened) ?? [] },
    { name: "Resolved", data: report?.series?.map((p) => p.resolved) ?? [] },
  ];
  const categories = report?.series?.map((p) => formatDateLabel(p.date)) ?? [];

  const options: ApexCharts.ApexOptions = {
    chart: { id: "p2p-dispute-trend", type: "line", fontFamily: "inherit", height: 300, toolbar: { show: false } },
    dataLabels: { enabled: false },
    stroke: { curve: "smooth", width: 2 },
    xaxis: { categories, labels: { style: { colors: "#828282", fontSize: "12px" } } },
    colors: ["#EB5757", "#00A859"],
    legend: { position: "top", horizontalAlign: "left" },
    grid: { borderColor: "#f1f1f1" },
  };

  return (
    <div>
      <MetricsCardContainer>
        <Card text="Disputes Opened" heading={isLoading ? "-" : report?.totalOpened ?? 0} img={disputesIcon} />
        <Card text="Still Open" heading={isLoading ? "-" : report?.stillOpen ?? 0} img={disputesIcon} />
        <Card
          text="Resolution Rate"
          heading={isLoading ? "-" : `${(report?.resolutionRate ?? 0).toFixed(1)}%`}
          img={successIcon}
        />
        <Card
          text="Avg. Resolution Time"
          heading={isLoading ? "-" : formatDuration(report?.averageResolutionTimeSeconds ?? 0)}
          img={averageIcon}
        />
      </MetricsCardContainer>

      <ReportSectionCard title="Dispute Health" subtitle="Opened vs. resolved disputes, per day">
        <Chart options={options} series={series} type="line" height={300} />
        <BreakdownGrid>
          <div>
            <ColumnTitle>By Subject</ColumnTitle>
            {(report?.bySubject ?? []).length === 0 ? (
              <Empty>No disputes in this period.</Empty>
            ) : (
              (report?.bySubject ?? []).map((item, i) => (
                <BreakdownRow key={item.key}>
                  <Dot color={CHART_COLORS[i % CHART_COLORS.length]} />
                  <span>{titleCase(item.key)}</span>
                  <Count>{item.count}</Count>
                </BreakdownRow>
              ))
            )}
          </div>
          <div>
            <ColumnTitle>By Resolution</ColumnTitle>
            {(report?.byResolution ?? []).length === 0 ? (
              <Empty>No resolved disputes in this period.</Empty>
            ) : (
              (report?.byResolution ?? []).map((item, i) => (
                <BreakdownRow key={item.key}>
                  <Dot color={CHART_COLORS[i % CHART_COLORS.length]} />
                  <span>{titleCase(item.key)}</span>
                  <Count>{item.count}</Count>
                </BreakdownRow>
              ))
            )}
          </div>
        </BreakdownGrid>
      </ReportSectionCard>
    </div>
  );
};

export default DisputeReportSection;

const MetricsCardContainer = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
`;

const BreakdownGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  margin-top: 20px;
`;

const ColumnTitle = styled.h4`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
  margin: 0 0 12px 0;
`;

const BreakdownRow = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0;
  font-size: 13px;
  color: #4b5563;
`;

const Dot = styled.span<{ color: string }>`
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: ${(p) => p.color};
  flex-shrink: 0;
`;

const Count = styled.span`
  margin-left: auto;
  font-weight: 600;
  color: #00225a;
`;

const Empty = styled.p`
  font-size: 13px;
  color: #828282;
`;
