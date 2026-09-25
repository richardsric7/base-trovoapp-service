"use client";
import ApexCharts from "apexcharts";
import { Select } from "antd";
import dynamic from "next/dynamic";
import React, { useState } from "react";
import styled from "styled-components";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

type TabName = "weekly" | "monthly" | "yearly";

const FeesGenerated = () => {
  const [activeTab, setActiveTab] = useState<TabName>("monthly");

  const getCategories = (tab: TabName): string[] => {
    switch (tab) {
      case "weekly":
        return ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

      case "monthly":
        return [
          "JAN",
          "FEB",
          "MAR",
          "APR",
          "MAY",
          "JUN",
          "JUL",
          "AUG",
          "SEP",
          "OCT",
          "NOV",
          "DEC",
        ];

      case "yearly":
        return ["2020", "2021", "2022", "2023", "2024", "2025"];
    }
  };

  const getData = (tab: TabName): number[] => {
    switch (tab) {
      case "weekly":
        return [12000, 18000, 15000, 22000, 25000, 30000, 28000];

      case "monthly":
        return [
          50000, 45000, 60000, 70000, 80000, 85000, 70000, 60000, 75000, 90000,
          95000, 100000,
        ];

      case "yearly":
        return [50000, 40000, 300000, 350000, 200000, 250000];
    }
  };

  const categories = getCategories(activeTab);
  const data = getData(activeTab);

  const series = [
    {
      name: "Fees Generated",
      data,
    },
  ];

  const options: ApexCharts.ApexOptions = {
    chart: {
      id: "fees-generated-chart",
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
        formatter: (value) => `${value.toLocaleString()} NGN`,
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
  const handleChange = (value: TabName) => {
    setActiveTab(value);
  };

  return (
    <Container>
      <FlexContent>
        <div>
          <Text>Fees Generated</Text>
          <VolContent>
            <TotalUsers>N100,000,000</TotalUsers>

            <Text>Last year</Text>
          </VolContent>
        </div>
        <SelectContainer>
          <Select
            defaultValue={activeTab}
            style={{ width: 120 }}
            onChange={handleChange}
            options={[
              { value: "yearly", label: "Yearly" },
              { value: "weekly", label: "Weekly" },
              { value: "monthly", label: "Monthly" },
            ]}
          />
        </SelectContainer>
      </FlexContent>

      <Chart options={options} series={series} type="area" height={360} />
    </Container>
  );
};

export default FeesGenerated;

const Container = styled.div`
  background-color: #ffffff;
  padding: 32px;
  border-radius: 24px;
  width: 100%;
  height: 100%;
`;

const Text = styled.p`
  font-size: 16px;
  font-weight: 400;
  margin: 0;
  color: #828282;
`;

const TotalUsers = styled.h3`
  font-size: 20px;
  font-weight: 700;
  line-height: 48px;
  color: #1e1b39;
`;

const FlexContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
`;

const VolContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 4px;
  gap: 6px;
`;

const SelectContainer = styled.div`
  gap: 10px;
  background-color: #f2f6f9;
  padding: 4px;
  border-radius: 12px;
`;
