import React, { useState } from "react";
import styled from "styled-components";
import Chart from "react-apexcharts";
import { DropdownSelect } from "@/components";
const SteppedLineChart = () => {
  const [selectedOption, setSelectedOption] =
    useState<string>("Thu 21 Sept, 2023");
  const series = [
    {
      name: "Uptime",
      data: [34, 44, 54, 21, 12, 43, 33, 23, 66, 66, 58],
    },
    {
      name: "Downtime",
      data: [34, 44, 54, 21, 12, 43, 33, 23, 66, 66, 58],
    },
  ];

  const options: ApexCharts.ApexOptions = {
    chart: {
      type: "line",
      height: 350,
      fontFamily: "inherit",
      toolbar: {
        show: false,
      },
    },
    stroke: {
      curve: "stepline",
    },
    colors: ["#00A859", "#BE3800"],
    dataLabels: {
      enabled: false,
    },

    markers: {
      size: 5,
      colors: undefined,
      strokeColors: "#fff",
      strokeWidth: 2,
      hover: {
        sizeOffset: 4,
      },
    },
  };

  return (
    <Container>
      <FlexContent>
        <div>
          <div>
            <Text>Activity</Text>
            <TotalUsers>Uptime and Downtime History </TotalUsers>
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
        Current Status: <StyledSpan>Running and Healthy</StyledSpan>
      </Status>
      <Chart options={options} series={series} type="line" height={380} />
    </Container>
  );
};

export default SteppedLineChart;

const Container = styled.div``;

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
  line-height: 20px;
  color: #828282;
`;

const StyledSpan = styled.span`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  color: #00a859;
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
