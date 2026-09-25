import { ArrowCircleUp, TrendUp } from "iconsax-react";
import dynamic from "next/dynamic";
import React, { useState } from "react";
import { BsArrowUpRightCircleFill } from "react-icons/bs";
import styled from "styled-components";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

type TabName = "daily" | "weekly" | "monthly";

const TradingVolumeChart = () => {
  const [activeTab, setActiveTab] = useState<TabName>("monthly");

  const getCategories = (tab: TabName): string[] => {
    switch (tab) {
      case "daily":
        return [
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
        ];
      case "weekly":
        return ["Week 1", "Week 2", "Week 3", "Week 4"];
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
    }
  };

  const getData = (tab: TabName): number[] => {
    switch (tab) {
      case "daily":
        return [
          2000, 2500, 3000, 4500, 4000, 3500, 3800, 4200, 3900, 4300, 4800,
          4100, 3700,
        ];
      case "weekly":
        return [12000, 15000, 18000, 20000];
      case "monthly":
        return [
          50000, 45000, 60000, 70000, 80000, 85000, 70000, 60000, 75000, 90000,
          95000, 100000,
        ];
    }
  };

  const categories = getCategories(activeTab);
  const data = getData(activeTab);

  const series = [
    {
      name: "Trading volume",
      data,
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

  const peakTradingHour = () => {
    const maxIndex = data.indexOf(Math.max(...data));
    return categories[maxIndex];
  };

  return (
    <Container>
      <FlexContent>
        <div>
          <Text>Trading Volume</Text>
          <VolContent>
            <TotalUsers>$12.7k</TotalUsers>
            <Status>
              {" "}
              <BsArrowUpRightCircleFill color="00A859" />
              <StyledSpan>1.3%</StyledSpan>
            </Status>
            <Text>VS LAST YEAR</Text>
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

const Status = styled.div`
  font-size: 14px;
  color: #00a859;
  margin: 0;
  display: flex;
  align-items: center;
`;

const StyledSpan = styled.span`
  font-size: 14px;
  font-weight: 500;
  color: #00a859;
`;
