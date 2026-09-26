import React from "react";
import dynamic from "next/dynamic";
import styled from "styled-components";
import { useP2pHourlyActivityReportQuery } from "@/redux/api/p2p";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

const formatHour = (h: number) => `${h.toString().padStart(2, "0")}:00`;

const TradingHoursChart = () => {
  // 90 days gives a representative sample of "which hour is usually
  // busiest" without needing a date picker - the underlying report can
  // take a `range` param if a narrower/wider window is ever needed here.
  const { data, isLoading } = useP2pHourlyActivityReportQuery({ range: "90d" });
  const report = data?.data;

  const series = [
    {
      name: "Orders",
      data: (report?.series ?? []).map((p) => ({ x: formatHour(p.hour), y: p.count })),
    },
  ];

  const options: ApexCharts.ApexOptions = {
    chart: {
      id: "trading-hours",
      type: "area",
      fontFamily: "inherit",
      height: 350,
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
        formatter: (value) => `${value.toLocaleString()} order${value === 1 ? "" : "s"}`,
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
          {
            offset: 0,
            color: "#62ACE8",
            opacity: 1,
          },
          {
            offset: 50,
            color: "#ACD1EF",
            opacity: 0.5,
          },
          {
            offset: 100,
            color: "#FFFFFF",
            opacity: 0.1,
          },
        ],
      },
    },
    grid: {
      borderColor: "#f1f1f1",
      xaxis: {
        lines: {
          show: false,
        },
      },
      yaxis: {
        lines: {
          show: true,
        },
      },
    },
  };

  return (
    <Container>
      <FlexContent>
        <div>
          <Text>Statistics</Text>
          <TotalUsers>Trading Hours (UTC, last 90 days)</TotalUsers>
        </div>
      </FlexContent>
      <Status>
        Peak Trading Hour: <StyledSpan>{isLoading ? "-" : formatHour(report?.peakHour ?? 0)}</StyledSpan>
      </Status>
      <Chart options={options} series={series} type="area" height={380} />
    </Container>
  );
};

export default TradingHoursChart;

const Container = styled.div`
  background-color: #ffffff;
  padding: 40px 32px 40px 32px;
  border-radius: 24px;
  margin: 20px 0;
`;

const Text = styled.p`
  font-size: 16px;
  font-weight: 400;
  line-height: 20px;
  margin: 0;
  color: #828282;
`;

const Status = styled.p`
  font-size: 14px;
  font-weight: 400;
  color: #828282;
  margin: 0;
`;

const StyledSpan = styled.span`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  margin: 0;
`;

const TotalUsers = styled.h3`
  font-size: 20px;
  font-weight: 600;
  line-height: 28px;
  color: #00225a;
`;

const FlexContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 20px;
`;
