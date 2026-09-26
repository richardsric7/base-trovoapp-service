"use client";
import React from "react";
import dynamic from "next/dynamic";
import styled from "styled-components";
import { ReportRange, useP2pDistributionReportQuery, IVolumeItem, ICountItem } from "@/redux/api/p2p";
import ReportSectionCard from "./ReportSectionCard";
import { CHART_COLORS, formatAmount, titleCase } from "./reportUtils";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

const barOptions = (categories: string[]): ApexCharts.ApexOptions => ({
  chart: { type: "bar", fontFamily: "inherit", height: 280, toolbar: { show: false } },
  plotOptions: { bar: { horizontal: true, borderRadius: 2, distributed: true } },
  dataLabels: { enabled: false },
  legend: { show: false },
  xaxis: { categories, labels: { style: { colors: "#828282", fontSize: "12px" } } },
  colors: CHART_COLORS,
  grid: { borderColor: "#f1f1f1" },
});

const VolumeBar = ({ title, items }: { title: string; items: IVolumeItem[] | undefined }) => {
  const top = (items ?? []).slice(0, 8);
  const series = [{ name: "Trades", data: top.map((i) => i.count) }];
  const categories = top.map((i) => i.key);
  return (
    <Column>
      <ColumnTitle>{title}</ColumnTitle>
      {top.length === 0 ? (
        <Empty>No trades in this period.</Empty>
      ) : (
        <>
          <Chart options={barOptions(categories)} series={series} type="bar" height={Math.max(200, top.length * 34)} />
          <VolumeList>
            {top.map((i) => (
              <VolumeRow key={i.key}>
                <span>{i.key}</span>
                <span>{formatAmount(i.volume)} traded</span>
              </VolumeRow>
            ))}
          </VolumeList>
        </>
      )}
    </Column>
  );
};

const StatusDonut = ({ items }: { items: ICountItem[] | undefined }) => {
  const data = items ?? [];
  const series = data.map((i) => i.count);
  const labels = data.map((i) => titleCase(i.key));
  const options: ApexCharts.ApexOptions = {
    chart: { type: "donut", fontFamily: "inherit" },
    labels,
    colors: CHART_COLORS,
    legend: { position: "bottom" },
    dataLabels: { enabled: true, formatter: (val: number) => `${val.toFixed(0)}%` },
  };
  return (
    <Column>
      <ColumnTitle>Order Status Mix</ColumnTitle>
      {data.length === 0 ? <Empty>No orders in this period.</Empty> : <Chart options={options} series={series} type="donut" height={300} />}
    </Column>
  );
};

const DistributionReportSection = ({ range }: { range: ReportRange }) => {
  const { data } = useP2pDistributionReportQuery({ range });
  const report = data?.data;

  return (
    <ReportSectionCard title="Market Distribution" subtitle="What's trading, in what currency, and where">
      <Grid>
        <StatusDonut items={report?.byStatus} />
        <VolumeBar title="Most Traded Assets" items={report?.byAsset} />
        <VolumeBar title="Most Traded Currencies" items={report?.byCurrency} />
        <VolumeBar title="Trades by Country" items={report?.byCountry} />
      </Grid>
    </ReportSectionCard>
  );
};

export default DistributionReportSection;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
`;

const Column = styled.div`
  background-color: #f9fafb;
  border-radius: 16px;
  padding: 20px;
`;

const ColumnTitle = styled.h4`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
  margin: 0 0 12px 0;
`;

const Empty = styled.p`
  font-size: 13px;
  color: #828282;
  padding: 24px 0;
  text-align: center;
`;

const VolumeList = styled.div`
  margin-top: 8px;
`;

const VolumeRow = styled.div`
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #4b5563;
  padding: 4px 0;
`;
