"use client";
import React from "react";
import dynamic from "next/dynamic";
import styled from "styled-components";
import Card from "../../../components/Card";
import revenueIcon from "@/assets/images/revenueIcon.svg";
import moneyIcon from "@/assets/images/money.svg";
import moneyBagIcon from "@/assets/images/solar_money-bag-linear.svg";
import { ReportRange, useP2pRevenueReportQuery } from "@/redux/api/p2p";
import ReportSectionCard from "./ReportSectionCard";
import { formatAmount, formatDateLabel } from "./reportUtils";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

const RevenueReportSection = ({ range }: { range: ReportRange }) => {
  const { data, isLoading } = useP2pRevenueReportQuery({ range });
  const report = data?.data;

  const series = [
    { name: "Platform fee", data: report?.series?.map((p) => Number(p.platformFee)) ?? [] },
    { name: "Regulatory fee", data: report?.series?.map((p) => Number(p.regulatoryFee)) ?? [] },
    { name: "VAT", data: report?.series?.map((p) => Number(p.vat)) ?? [] },
  ];
  const categories = report?.series?.map((p) => formatDateLabel(p.date)) ?? [];

  const options: ApexCharts.ApexOptions = {
    chart: { id: "p2p-revenue-trend", type: "bar", stacked: true, fontFamily: "inherit", height: 320, toolbar: { show: false } },
    plotOptions: { bar: { columnWidth: "45%", borderRadius: 2 } },
    dataLabels: { enabled: false },
    xaxis: { categories, labels: { style: { colors: "#828282", fontSize: "12px" } } },
    colors: ["#007CDF", "#F5A623", "#9B51E0"],
    legend: { position: "top", horizontalAlign: "left" },
    grid: { borderColor: "#f1f1f1" },
  };

  return (
    <div>
      <MetricsCardContainer>
        <Card
          text="Total Revenue Collected"
          heading={isLoading ? "-" : formatAmount(report?.totalRevenue)}
          img={revenueIcon}
        />
        <Card text="Platform Fees" heading={isLoading ? "-" : formatAmount(report?.totalPlatformFee)} img={moneyIcon} />
        <Card
          text="Regulatory Fees"
          heading={isLoading ? "-" : formatAmount(report?.totalRegulatoryFee)}
          img={moneyBagIcon}
        />
        <Card text="VAT Collected" heading={isLoading ? "-" : formatAmount(report?.totalVat)} img={moneyIcon} />
      </MetricsCardContainer>

      <ReportSectionCard title="Fee Revenue Trend" subtitle="Collected on completed orders only, per day">
        <Chart options={options} series={series} type="bar" height={320} />
      </ReportSectionCard>
    </div>
  );
};

export default RevenueReportSection;

const MetricsCardContainer = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
`;
