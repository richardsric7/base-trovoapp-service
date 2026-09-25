"use client";

import React from "react";
import ReactApexChart from "react-apexcharts";
import { ApexOptions } from "apexcharts";
import styled from "styled-components";

interface RevenueBarChartProps {
  chartType: "today" | "last7days" | "last30days";
}

const RevenueBarChart: React.FC<RevenueBarChartProps> = ({ chartType }) => {
  const getSeries = () => {
    const commonSeries = [
      {
        name: "Trovo Wallet Swap Fee",
        today: [1900000, 1500000, 1200000, 1300000, 1100000, 1700000, 2000000],
        last7days: [
          2800000, 3000000, 2500000, 4000000, 2600000, 3700000, 3200000,
        ],
        last30days: [
          1000000, 1100000, 1200000, 1300000, 1400000, 1500000, 1200000,
        ],
      },
      {
        name: "Trovo P2P Fee",
        today: [1200000, 1300000, 1000000, 1500000, 900000, 1600000, 1400000],
        last7days: [
          3300000, 2800000, 3000000, 3600000, 3400000, 2900000, 3100000,
        ],
        last30days: [
          1300000, 1400000, 1350000, 1450000, 1380000, 1480000, 1420000,
        ],
      },
      {
        name: "Trovo Patron Membership Fee",
        today: [800000, 850000, 900000, 950000, 1000000, 1100000, 1200000],
        last7days: [
          2100000, 2200000, 2000000, 2500000, 2400000, 2600000, 2300000,
        ],
        last30days: [
          1000000, 900000, 1100000, 1300000, 1200000, 1400000, 1250000,
        ],
      },
      {
        name: "Trovo Wallet Recovery Subscription Fee",
        today: [700000, 800000, 850000, 600000, 500000, 900000, 1000000],
        last7days: [
          2000000, 1800000, 1500000, 2500000, 2200000, 2300000, 2400000,
        ],
        last30days: [
          900000, 1100000, 1200000, 1000000, 1050000, 1300000, 1400000,
        ],
      },
      {
        name: "Asset Tokenization Fee",
        today: [500000, 550000, 700000, 650000, 800000, 950000, 1000000],
        last7days: [
          1700000, 1600000, 1800000, 1900000, 2000000, 1950000, 1850000,
        ],
        last30days: [900000, 850000, 880000, 890000, 870000, 890000, 860000],
      },
      {
        name: "Shared Access Payment Transaction Fee",
        today: [300000, 350000, 400000, 380000, 370000, 420000, 430000],
        last7days: [
          1200000, 1300000, 1400000, 1500000, 1350000, 1250000, 1450000,
        ],
        last30days: [700000, 750000, 770000, 800000, 780000, 760000, 740000],
      },
      {
        name: "Sub Wallet Creation Fee",
        today: [200000, 180000, 220000, 240000, 230000, 210000, 250000],
        last7days: [900000, 850000, 920000, 880000, 950000, 970000, 960000],
        last30days: [600000, 580000, 590000, 610000, 620000, 600000, 630000],
      },
      {
        name: "Wrapped Asset Withdrawal Fee",
        today: [100000, 120000, 140000, 160000, 180000, 150000, 170000],
        last7days: [500000, 600000, 650000, 700000, 620000, 680000, 660000],
        last30days: [400000, 420000, 430000, 440000, 410000, 450000, 460000],
      },
      {
        name: "Statement of Account Generation Fee",
        today: [90000, 85000, 95000, 100000, 110000, 105000, 95000],
        last7days: [300000, 350000, 330000, 360000, 320000, 340000, 310000],
        last30days: [200000, 210000, 220000, 230000, 240000, 250000, 260000],
      },
    ];

    return commonSeries.map((item) => ({
      name: item.name,
      data: item[chartType],
    }));
  };

  const options: ApexOptions = {
    chart: {
      type: "bar",
      stacked: true,
      toolbar: { show: false },
    },
    colors: [
      "#00225A",
      "#004988",
      "#007CDF",
      "#62ACE8",
      "#D5E1F8",
      "#4F9A94",
      "#95C2BF",
      "#C08B8C",
      "#D9B9BA",
    ],
    plotOptions: {
      bar: {
        horizontal: false,
        columnWidth: "60%",
        borderRadius: 6,
      },
    },
    dataLabels: { enabled: false },
    grid: {
      show: true,
      borderColor: "#E0E0E0",
      strokeDashArray: 3,
    },

    xaxis: {
      categories: [
        "08/01",
        "09/01",
        "10/01",
        "11/01",
        "12/01",
        "13/01",
        "14/01",
      ],
      labels: {
        show: true,
        style: {
          fontSize: "12px",
          fontFamily: "inherit",
          colors: "#828282",
        },
      },
      axisBorder: { show: false },
      axisTicks: { show: false },
    },

    yaxis: {
      labels: {
        formatter: (val: number) => {
          if (val >= 1_000_000) return `$${val / 1_000_000}M`;
          if (val >= 1_000) return `$${val / 1_000}K`;
          return `$${val}`;
        },
        style: {
          fontFamily: "inherit",
          fontSize: "14px",
          colors: "#828282",
        },
      },
    },
    tooltip: {
      y: {
        formatter: (val: number) => `$${val.toLocaleString()}`,
      },
    },
    legend: {
      position: "top",
      horizontalAlign: "left",
      fontFamily: "inherit",

      fontSize: "14px",
    },
    fill: {
      opacity: 1,
    },
  };

  return (
    <StyledChartWrapper>
      <ReactApexChart
        options={options}
        series={getSeries()}
        type="bar"
        height={500}
      />
    </StyledChartWrapper>
  );
};

export default RevenueBarChart;

const StyledChartWrapper = styled.div`
  .apexcharts-legend {
    display: grid !important;
    grid-template-columns: repeat(3, 1fr) !important;
    gap: 12px 24px !important;
    justify-items: start !important;
    align-items: center !important;
    padding-bottom: 8px;
  }

  .apexcharts-legend-text {
    font-family: inherit !important;
    font-size: 14px !important;
    color: #4f4f4f !important;
  }
`;
