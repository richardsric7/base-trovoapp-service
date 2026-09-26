"use client";
import React from "react";
import dynamic from "next/dynamic";
import styled from "styled-components";
import Card from "../../../components/Card";
import stockIcon from "@/assets/images/icon-park-outline_stock-market.svg";
import moneyBagIcon from "@/assets/images/solar_money-bag-linear.svg";
import successIcon from "@/assets/images/copy-success.svg";
import { ReportRange, useP2pVolumeReportQuery } from "@/redux/api/p2p";
import ReportSectionCard from "./ReportSectionCard";
import { formatAmount, formatDateLabel } from "./reportUtils";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

const VolumeReportSection = ({ range }: { range: ReportRange }) => {
  const { data, isLoading } = useP2pVolumeReportQuery({ range });
  const report = data?.data;

  const series = [
    { name: "Total orders", data: report?.series?.map((p) => p.totalOrders) ?? [] },
    { name: "Completed orders", data: report?.series?.map((p) => p.completedOrders) ?? [] },
  ];
  const categories = report?.series?.map((p) => formatDateLabel(p.date)) ?? [];

  const options: ApexCharts.ApexOptions = {
    chart: { id: "p2p-volume-trend", type: "area", fontFamily: "inherit", height: 320, toolbar: { show: false } },
    dataLabels: { enabled: false },
    stroke: { curve: "smooth", width: 2 },
    xaxis: { categories, labels: { style: { colors: "#828282", fontSize: "12px" } } },
    yaxis: { labels: { style: { colors: "#828282", fontSize: "12px" } } },
    colors: ["#62ACE8", "#00A859"],
    legend: { position: "top", horizontalAlign: "left" },
    fill: { type: "gradient", gradient: { shadeIntensity: 1, opacityFrom: 0.5, opacityTo: 0.1 } },
    grid: { borderColor: "#f1f1f1" },
  };

  const completionRate =
    report && report.totalOrders > 0 ? ((report.completedOrders / report.totalOrders) * 100).toFixed(1) : "0.0";

  return (
    <div>
      <MetricsCardContainer>
        <Card text="Total Orders" heading={isLoading ? "-" : report?.totalOrders ?? 0} img={stockIcon} />
        <Card text="Completed Orders" heading={isLoading ? "-" : report?.completedOrders ?? 0} img={successIcon} />
        <Card
          text="Completed Volume (payment amount)"
          heading={isLoading ? "-" : formatAmount(report?.completedVolume)}
          img={moneyBagIcon}
        />
        <Card text="Completion Rate" heading={isLoading ? "-" : `${completionRate}%`} img={successIcon} />
      </MetricsCardContainer>

      <ReportSectionCard title="Trading Volume Trend" subtitle="Orders created vs. completed, per day">
        <Chart options={options} series={series} type="area" height={320} />
      </ReportSectionCard>
    </div>
  );
};

export default VolumeReportSection;

const MetricsCardContainer = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
`;
