"use client";
import React from "react";
import dynamic from "next/dynamic";
import styled from "styled-components";
import Card from "../../../components/Card";
import newUsersIcon from "@/assets/images/profile-add.svg";
import merchantsIcon from "@/assets/images/merchant-icon.svg";
import { ReportRange, useP2pGrowthReportQuery } from "@/redux/api/p2p";
import ReportSectionCard from "./ReportSectionCard";
import { formatDateLabel } from "./reportUtils";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

const GrowthReportSection = ({ range }: { range: ReportRange }) => {
  const { data, isLoading } = useP2pGrowthReportQuery({ range });
  const report = data?.data;

  const series = [
    { name: "New offers", data: report?.series?.map((p) => p.newOffers) ?? [] },
    { name: "New merchants", data: report?.series?.map((p) => p.newMerchants) ?? [] },
  ];
  const categories = report?.series?.map((p) => formatDateLabel(p.date)) ?? [];

  const options: ApexCharts.ApexOptions = {
    chart: { id: "p2p-growth-trend", type: "bar", fontFamily: "inherit", height: 300, toolbar: { show: false } },
    plotOptions: { bar: { columnWidth: "40%", borderRadius: 2 } },
    dataLabels: { enabled: false },
    xaxis: { categories, labels: { style: { colors: "#828282", fontSize: "12px" } } },
    colors: ["#007CDF", "#ACD1EF"],
    legend: { position: "top", horizontalAlign: "left" },
    grid: { borderColor: "#f1f1f1" },
  };

  return (
    <div>
      <MetricsCardContainer>
        <Card text="New Offers" heading={isLoading ? "-" : report?.totalNewOffers ?? 0} img={merchantsIcon} />
        <Card text="First-Time Merchants" heading={isLoading ? "-" : report?.totalNewMerchants ?? 0} img={newUsersIcon} />
      </MetricsCardContainer>

      <ReportSectionCard title="Marketplace Growth" subtitle="New listings and first-time merchants, per day">
        <Chart options={options} series={series} type="bar" height={300} />
      </ReportSectionCard>
    </div>
  );
};

export default GrowthReportSection;

const MetricsCardContainer = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  gap: 10px;

  > div {
    max-width: 280px;
  }
`;
