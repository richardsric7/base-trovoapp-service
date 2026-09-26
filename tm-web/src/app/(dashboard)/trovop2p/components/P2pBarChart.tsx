"use client";
import React, { useState } from "react";
import styled from "styled-components";
import Chart from "react-apexcharts";
import P2pChartFilter from "./P2pChartFilter";

type DateRangeOption =
  | "Today"
  | "Yesterday"
  | "Last week"
  | "Last month"
  | "Last 3 months"
  | "Last year"
  | "All time"
  | "Custom";

type DisplayUnit = "hours" | "days" | "weeks" | "months" | "years";
interface BarChartProps {
  data?: {
    data?: {
      totalOrders?: number;
    };
  };
}

const unitLabels: Record<DisplayUnit, string> = {
  hours: "Hours",
  days: "Days",
  weeks: "Weeks",
  months: "Months",
  years: "Years",
};
const P2pBarChart: React.FC<BarChartProps> = ({ data }) => {
  const [selectedRange, setSelectedRange] = useState<DateRangeOption>("Today");
  const [displayUnit, setDisplayUnit] = useState<DisplayUnit>("hours");
  // Dummy dataset
  const dummyData = {
    daily_active_and_new_users: { active_users: 120, new_users: 45 },
    weekly_active_and_new_users: { active_users: 450, new_users: 180 },
    monthly_active_and_new_users: { active_users: 2000, new_users: 650 },
    total_users: 5600,
  };

  const monthLabels = [
    "Jan",
    "Feb",
    "Mar",
    "Apr",
    "May",
    "Jun",
    "Jul",
    "Aug",
    "Sep",
    "Oct",
    "Nov",
    "Dec",
  ];

  // categories per range
  const getCategories = (
    range: DateRangeOption,
    unit: DisplayUnit
  ): string[] => {
    if (range === "Today" || range === "Yesterday") {
      return unit === "hours"
        ? Array.from({ length: 24 }, (_, i) => `${i}:00`)
        : ["Day"];
    }
    if (range === "Last week") {
      return unit === "days"
        ? ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]
        : ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];
    }
    if (range === "Last month") {
      return ["Week 1", "Week 2", "Week 3", "Week 4"];
    }
    if (range === "Last 3 months" || range === "Last year") {
      return monthLabels;
    }
    if (range === "All time") {
      return unit === "years"
        ? ["2021", "2022", "2023", "2024", "2025"]
        : ["All Time"];
    }
    return ["Custom"];
  };

  // Prepare series data
  const getSeriesData = (range: DateRangeOption, unit: DisplayUnit) => {
    let activeUsers = 0;
    let newUsers = 0;

    if (range === "Today" || range === "Yesterday") {
      activeUsers = dummyData.daily_active_and_new_users.active_users;
      newUsers = dummyData.daily_active_and_new_users.new_users;
    } else if (range === "Last week") {
      activeUsers = dummyData.weekly_active_and_new_users.active_users;
      newUsers = dummyData.weekly_active_and_new_users.new_users;
    } else {
      activeUsers = dummyData.monthly_active_and_new_users.active_users;
      newUsers = dummyData.monthly_active_and_new_users.new_users;
    }

    const categories = getCategories(range, unit);
    return {
      activeUsers: new Array(categories.length).fill(activeUsers),
      newUsers: new Array(categories.length).fill(newUsers),
      categories,
    };
  };

  const { activeUsers, newUsers, categories } = getSeriesData(
    selectedRange,
    displayUnit
  );

  const series = [
    { name: "Active Users", data: activeUsers },
    { name: "New Users", data: newUsers },
  ];

  const options: ApexCharts.ApexOptions = {
    chart: {
      type: "bar",
      height: 350,
      fontFamily: "inherit",
      toolbar: {
        show: false,
      },
    },
    plotOptions: {
      bar: {
        horizontal: false,
        columnWidth: "25%",
        borderRadius: 2,
        borderRadiusApplication: "end",
      },
    },
    dataLabels: {
      enabled: false,
    },
    stroke: {
      show: true,
      width: 2,
      colors: ["transparent"],
    },
    xaxis: {
      categories,
    },
    fill: {
      opacity: 1,
    },
    tooltip: {
      y: {
        formatter: (val: number) => `${val} users`,
      },
    },
    colors: ["#007CDF", "#ACD1EF"],
    legend: {
      show: false,
    },
    states: {
      hover: {
        filter: {
          type: "none",
        },
      },
    },
  };

  return (
    <Container>
      <div>
        <HeaderContent>
          <div>
            <Text>Overview</Text>
            <TotalUsers>
              Total Orders: {data?.data?.totalOrders ?? 0}{" "}
            </TotalUsers>
          </div>

          <P2pChartFilter
            onApply={(range, unit) => {
              setSelectedRange(range);
              setDisplayUnit(unit);
            }}
          />
        </HeaderContent>

        <LegendContainer>
          <LegendContent>
            <LegendBox color="#007CDF" />
            <LegendText>Active Users</LegendText>
          </LegendContent>
          <LegendContent>
            <LegendBox color="#ACD1EF" />
            <LegendText>New Users</LegendText>
          </LegendContent>
        </LegendContainer>
      </div>

      <Chart options={options} series={series} type="bar" height={380} />
    </Container>
  );
};

export default P2pBarChart;

const Container = styled.div`
  width: 100%;
  //   max-width: 800px;
  margin: 20px auto;
  padding: 20px;
  background: #fff;
  border-radius: 24px;
`;

const HeaderContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;
const Text = styled.p`
  font-size: 16px;
  font-weight: 400;
  line-height: 20px;
  margin: 0;
  color: #828282;
`;

const TotalUsers = styled.h3`
  font-size: 20px;
  font-weight: 600;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;
const TabsContainer = styled.div`
  margin-top: 10px;
  display: flex;
  justify-content: center;
  gap: 10px;
  background-color: #f5f5f5;
  border-raduis: 12px;
`;

const TabButton = styled.div<{ $active?: boolean }>`
  padding: 8px 16px;
  border: none;
  background-color: ${(props) => (props.$active ? "white" : "transparent")};
  color: ${(props) => (props.$active ? "#00225A" : "#333333")};
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 14px;
  font-weight: 400;
`;
const LegendContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 20px 0;
`;

const LegendContent = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

const LegendBox = styled.div<{ color: string }>`
  width: 16px;
  height: 16px;
  background-color: ${(props) => props.color};
  border-radius: 4px;
`;

const LegendText = styled.span`
  font-size: 14px;
  color: #333;
`;
