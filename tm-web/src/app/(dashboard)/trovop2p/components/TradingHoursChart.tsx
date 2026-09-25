import React, { useState } from "react";
import dynamic from "next/dynamic";
import styled from "styled-components";
import { DropdownSelect } from "@/components";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

const TradingHoursChart = () => {
  const [selectedOption, setSelectedOption] =
    useState<string>("Thu 21 Sept, 2023");

  const series = [
    {
      name: "Trading Hours",
      data: [
        { x: "0:00", y: 2000 },
        { x: "1:00", y: 2500 },
        { x: "2:00", y: 3000 },
        { x: "3:00", y: 4500 },
        { x: "4:00", y: 4000 },
        { x: "5:00", y: 3500 },
        { x: "6:00", y: 3800 },
        { x: "7:00", y: 4200 },
        { x: "8:00", y: 3900 },
        { x: "9:00", y: 4300 },
        { x: "10:00", y: 4800 },
        { x: "11:00", y: 4100 },
        { x: "12:00", y: 3700 },
      ],
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
      categories: [
        "0:00",
        "1:00",
        "2:00",
        "3:00",
        "4:00",
        "5:00",
        "6:00",
        "7:00",
        "8:00",
        "9:00",
        "10:00",
        "11:00",
        "12:00",
      ],
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
      min: 1000,
      max: 5000,
      tickAmount: 4,
      labels: {
        formatter: (value) => `${value / 1000}k`,
        style: {
          colors: "#828282",
          fontSize: "12px",
        },
      },
    },
    tooltip: {
      x: {
        format: "HH:mm",
      },
      y: {
        formatter: (value) => `${value.toLocaleString()} hours`,
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
          <div>
            <Text>Statistics</Text>
            <TotalUsers>Trading Hours</TotalUsers>
          </div>
        </div>
        <DropdownSelect
          options={["Thu 21 Sept, 2023", "Thu 21 Sept, 2024"]}
          placeholder="Select Date"
          labelText=""
          value={selectedOption}
          onSelect={(item) => setSelectedOption(item)}
          backgroundColor="#F2F6F9"
          borderless={true}
          iconColor="#00225A"
        />
      </FlexContent>
      <Status>
        Peak Trading Hours: <StyledSpan>10:00</StyledSpan>
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
