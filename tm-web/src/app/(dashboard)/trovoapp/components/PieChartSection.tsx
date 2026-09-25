"use client";
import React from "react";
import styled from "styled-components";
import Chart from "react-apexcharts";

import { useWalletDistQuery } from "@/redux/api/usersmetrics";
const PieChartSection = () => {
  const { data, isLoading } = useWalletDistQuery();
  // console.log("pie", data);

  // const dummyData = {
  //   series: [11, 39, 26, 24],
  //   options: {
  //     chart: {
  //       type: "donut" as const,
  //       fontFamily: "inherit",
  //     },
  //     labels: ["Wallet A", "Wallet B", "Wallet C", "Wallet D"],
  //     colors: ["#004988", "#ACD1EF", "#62ACE8", "#007CDF"],
  //     legend: {
  //       show: false,
  //     },
  //     states: {
  //       hover: {
  //         filter: {
  //           type: "none",
  //         },
  //       },
  //     },
  //     plotOptions: {
  //       pie: {
  //         donut: {
  //           size: "55%",
  //           labels: {
  //             show: true,
  //             total: {
  //               show: true,

  //               formatter: (w: { globals: { seriesTotals: number[] } }) => {
  //                 const total = w.globals.seriesTotals.reduce(
  //                   (a: number, b: number) => a + b,
  //                   0
  //                 );
  //                 return total.toString();
  //               },

  //               label: "Total Users",
  //               style: {
  //                 fontSize: "14px",
  //                 fontWeight: 400,
  //                 color: "#828282",
  //               },
  //               offsetY: 40,
  //             },
  //             value: {
  //               show: true,
  //               offsetY: -10, // Move the number up
  //               fontSize: "22px",
  //               fontWeight: 600,
  //               color: "#00225A",
  //             },
  //           },
  //         },
  //       },
  //     },
  //     responsive: [
  //       {
  //         breakpoint: 480,
  //         options: {
  //           chart: {
  //             width: 300,
  //           },
  //           legend: {
  //             position: "bottom",
  //           },
  //         },
  //       },
  //     ],
  //   },
  // };

  const dummyData = {
    series: [11, 39, 26, 24],
    options: {
      chart: {
        type: "donut" as const,
        fontFamily: "inherit",
      },
      labels: ["Wallet A", "Wallet B", "Wallet C", "Wallet D"],
      colors: ["#004988", "#ACD1EF", "#62ACE8", "#007CDF"],
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
      plotOptions: {
        pie: {
          donut: {
            size: "55%",
            labels: {
              show: false, // Disable default labels
            },
          },
        },
      },
      responsive: [
        {
          breakpoint: 480,
          options: {
            chart: {
              width: 300,
            },
            legend: {
              position: "bottom",
            },
          },
        },
      ],
    },
  };

  const totalUsers = dummyData.series.reduce((a, b) => a + b, 0);

  return (
    <Container>
      <Text>Statistics</Text>
      <SubTitle>Wallet Version Distribution</SubTitle>
      <Chart
        options={dummyData.options}
        series={dummyData.series}
        type="donut"
        width="100%"
      />
      <CustomLabel>
        <LabelValue>{totalUsers}</LabelValue>
        <LabelText>Total Users</LabelText>
      </CustomLabel>

      <LegendContainer>
        <LegendContent>
          <LegendLeft>
            <LegendBox color="#004988" />
            <LegendTitle>Version 1.43</LegendTitle>
          </LegendLeft>
          <LegendText>412</LegendText>
        </LegendContent>
        <LegendContent>
          <LegendLeft>
            <LegendBox color="#007CDF" />
            <LegendTitle>Version 1.43</LegendTitle>
          </LegendLeft>
          <LegendText>412</LegendText>
        </LegendContent>
        <LegendContent>
          <LegendLeft>
            <LegendBox color="#62ACE8" />
            <LegendTitle>Version 1.43</LegendTitle>
          </LegendLeft>
          <LegendText>412</LegendText>
        </LegendContent>
        <LegendContent>
          <LegendLeft>
            <LegendBox color="#ACD1EF" />
            <LegendTitle>Version 1.43</LegendTitle>
          </LegendLeft>
          <LegendText>412</LegendText>
        </LegendContent>
      </LegendContainer>
    </Container>
  );
};

export default PieChartSection;

const Container = styled.div`
  display: flex;
  flex-direction: column;
  // align-items: center;
  justify-content: center;
  width: 100%;
  max-width: 400px;
  margin: auto;
  padding: 20px;
  background-color: white;
  position: relative;
  border-radius: 24px;
`;
const Text = styled.p`
  font-size: 16px;
  font-weight: 400;
  line-height: 20px;
  margin: 0;
  color: #828282;
`;
const SubTitle = styled.h3`
  font-size: 20px;
  font-weight: 600;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;

const CustomLabel = styled.div`
  position: absolute;
  display: flex;
  right: 160px;
  top: 195px;

  flex-direction: column;
  align-items: center;
  justify-content: center;
`;

const LabelValue = styled.div`
  font-size: 22px;
  font-weight: 600;
  color: #00225a;
`;

const LabelText = styled.div`
  font-size: 14px;
  font-weight: 400;
  color: #828282;
`;

const LegendContainer = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  align-items: center;
  gap: 10px;

  margin: 20px 0;
`;
const LegendLeft = styled.div`
  display: flex;
  //   align-items: center;
  gap: 5px;
`;
const LegendContent = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
  justify-content: space-between;
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

const LegendTitle = styled.span`
  color: #00225a;
  font-size: 14px;
  font-weight: 600;
  text-align: left;
`;
