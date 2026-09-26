import dynamic from "next/dynamic";
import React, { useState } from "react";
import styled from "styled-components";
import { ReportRange, useP2pVolumeReportQuery } from "@/redux/api/p2p";
import { formatAmount, formatDateLabel } from "./reports/reportUtils";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

type TabName = "daily" | "weekly" | "monthly";

const TAB_RANGE: Record<TabName, ReportRange> = {
  daily: "7d",
  weekly: "30d",
  monthly: "1y",
};

const TradingVolumeChart = () => {
  const [activeTab, setActiveTab] = useState<TabName>("monthly");
  const { data, isLoading } = useP2pVolumeReportQuery({ range: TAB_RANGE[activeTab] });
  const report = data?.data;

  const categories = report?.series?.map((p) => formatDateLabel(p.date)) ?? [];
  const series = [
    {
      name: "Completed trading volume",
      data: report?.series?.map((p) => Number(p.completedVolume)) ?? [],
    },
  ];

  const options: ApexCharts.ApexOptions = {
    chart: {
      id: "trading-volume",
      type: "area",
      fontFamily: "inherit",
      height: 360,
      toolbar: {
        show: false,
      },
    },
    dataLabels: {
      enabled: false,
    },
    stroke: {
      curve: "smooth",
      width: 2,
    },
    xaxis: {
      type: "category",
      categories,
      labels: {
        style: {
          colors: "#828282",
          fontSize: "12px",
        },
        rotate: 0,
      },
      axisBorder: {
        show: false,
      },
      axisTicks: {
        show: false,
      },
    },
    yaxis: {
      labels: {
        style: {
          colors: "#828282",
          fontSize: "12px",
        },
      },
    },
    tooltip: {
      y: {
        formatter: (value) => formatAmount(String(value)),
      },
    },
    fill: {
      type: "gradient",
      gradient: {
        type: "vertical",
        shadeIntensity: 1,
        opacityFrom: 0.7,
        opacityTo: 0.2,
        stops: [0, 100],
        colorStops: [
          { offset: 0, color: "#62ACE8", opacity: 1 },
          { offset: 50, color: "#ACD1EF", opacity: 0.5 },
          { offset: 100, color: "#FFFFFF", opacity: 0.1 },
        ],
      },
    },
    grid: {
      borderColor: "#f1f1f1",
      xaxis: {
        lines: { show: false },
      },
      yaxis: {
        lines: { show: true },
      },
    },
  };

  return (
    <Container>
      <FlexContent>
        <div>
          <Text>Trading Volume</Text>
          <VolContent>
            <TotalUsers>{isLoading ? "-" : formatAmount(report?.completedVolume)}</TotalUsers>
          </VolContent>
        </div>
        <TabsContainer>
          <TabButton
            $active={activeTab === "daily"}
            onClick={() => setActiveTab("daily")}
          >
            Daily
          </TabButton>
          <TabButton
            $active={activeTab === "weekly"}
            onClick={() => setActiveTab("weekly")}
          >
            Weekly
          </TabButton>
          <TabButton
            $active={activeTab === "monthly"}
            onClick={() => setActiveTab("monthly")}
          >
            Monthly
          </TabButton>
        </TabsContainer>
      </FlexContent>

      <Chart options={options} series={series} type="area" height={360} />
    </Container>
  );
};

export default TradingVolumeChart;

const Container = styled.div`
  background-color: #ffffff;
  padding: 40px 32px;
  border-radius: 24px;
  margin: 20px 0;
  width: 60%;
`;

const Text = styled.p`
  font-size: 16px;
  font-weight: 400;
  margin: 0;
  color: #828282;
`;

const TotalUsers = styled.h3`
  font-size: 32px;
  font-weight: 700;
  line-height: 48px;
  color: #1e1b39;
`;

const FlexContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 20px;
  gap: 6px;
`;

const VolContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 4px;
  gap: 6px;
`;

const TabsContainer = styled.div`
  display: flex;
  gap: 10px;
  background-color: #f5f5f5;
  padding: 4px;
  border-radius: 12px;
`;

const TabButton = styled.div<{ $active?: boolean }>`
  padding: 8px 16px;
  background-color: ${(props) => (props.$active ? "white" : "transparent")};
  color: ${(props) => (props.$active ? "#00225a" : "#828282")};
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 400;
  transition: all 0.3s ease;

  &:hover {
    background-color: ${(props) => (props.$active ? "white" : "#e0e0e0")};
  }
`;
