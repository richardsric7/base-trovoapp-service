"use client";

import { DropdownSelect } from "@/components";
import { useCountiresStatisticsQuery } from "@/redux/api/usersmetrics";
import React, { useState } from "react";
import Chart from "react-apexcharts";
import styled from "styled-components";

const CountriesStatistics = () => {
  const [selectedOption, setSelectedOption] = useState<string>("September");
  const { data, isLoading } = useCountiresStatisticsQuery();
  // console.log("countries data", data);

  const countriesData = data?.data
    ? [...data.data].sort((a, b) => b.Count - a.Count)
    : [];
  const series = [
    {
      name: "Users",
      data: countriesData.map(
        (item: { Country: string; Count: number }) => item.Count,
      ),
    },
  ];
  const options: ApexCharts.ApexOptions = {
    chart: {
      type: "bar",
      height: 380,
      fontFamily: "inherit",
      toolbar: {
        show: false, // Disable the toolbar
      },
    },
    plotOptions: {
      bar: {
        horizontal: true,
        barHeight: "85%",
        borderRadius: 0,
        columnWidth: "90%", // Reduce this to make bars closer
        distributed: true,

        borderRadiusApplication: "around",
        dataLabels: {
          position: "top",
        },

        colors: {
          ranges: [
            {
              from: 0,
              to: 50,
              color: "#ACD1EF", // Light blue for "Others"
            },
          ],
          backgroundBarRadius: 8,
        },
      },
    },
    dataLabels: {
      enabled: true,
      style: {
        fontSize: "12px",
        colors: ["#333"],
      },
      formatter: (value: number) => `${value}`,
    },
    colors: ["#007CDF"], // Blue for "Top Countries"

    xaxis: {
      categories: countriesData.map((item) => item.Country || "Unknown"),
    },
    tooltip: {
      y: {
        formatter: (value: number) => `${value} users`,
      },
    },
    legend: {
      show: false,
      position: "top",
      horizontalAlign: "left",
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
      <Statistics>
        <>
          <HeaderContent>
            <div>
              <Text>Statistics</Text>
              <SubTitle>Users by Country</SubTitle>
            </div>
            <DropdownSelect
              options={["October", "November"]}
              placeholder="This Week"
              labelText=""
              selectColor="#828282"
              value={selectedOption}
              onSelect={(item) => setSelectedOption(item)}
              backgroundColor="#F2F6F9"
              borderless={true}
              iconColor="#828282"
            />
          </HeaderContent>
          <LegendContainer>
            <LegendContent>
              <LegendBox color="#007CDF" />
              <LegendText>Top Countries</LegendText>
            </LegendContent>
            <LegendContent>
              <LegendBox color="#ACD1EF" />
              <LegendText>Other Countries</LegendText>
            </LegendContent>
          </LegendContainer>
        </>

        <ChartContainer>
          <Chart options={options} series={series} type="bar" height={320} />
        </ChartContainer>
      </Statistics>

      <FeaturesContainer>
        <HeaderContent>
          <SubTitle>Popular Features</SubTitle>
          <DropdownSelect
            options={["October", "November"]}
            placeholder="This Week"
            labelText=""
            selectColor="#828282"
            value={selectedOption}
            onSelect={(item) => setSelectedOption(item)}
            backgroundColor="#F2F6F9"
            borderless={true}
            iconColor="#828282"
          />
        </HeaderContent>
        <Line></Line>
        <Features>
          {/* First Feature */}
          <FeatureContent>
            <Content>
              <Eclipse></Eclipse>
              <FeatureText>Asset Tokenization</FeatureText>
            </Content>
            <Sessions>5000 sessions</Sessions>
          </FeatureContent>

          {/* Second Feature */}
          <FeatureContentSecondary>
            <Content>
              <Eclipse></Eclipse>
              <FeatureTextSecondary>Swap</FeatureTextSecondary>
            </Content>
            <SessionsSecondary>4200 sessions</SessionsSecondary>
          </FeatureContentSecondary>

          {/* Third Feature */}
          <FeatureContentTertiary>
            <Content>
              <Eclipse></Eclipse>
              <FeatureTextTertiary>Trovo Voting</FeatureTextTertiary>
            </Content>
            <SessionsTertiary>3800 sessions</SessionsTertiary>
          </FeatureContentTertiary>
        </Features>
      </FeaturesContainer>
    </Container>
  );
};

export default CountriesStatistics;

const Container = styled.div`
  width: 60%;
`;

const Statistics = styled.div`
  background-color: #fff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
  margin-bottom: 10px;
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
const SubTitle = styled.h3`
  font-size: 20px;
  font-weight: 600;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;
const Features = styled.div`
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
`;
const FeaturesContainer = styled.div`
  background-color: #fff;
  padding: 40px 32px;
  gap: 20px;
  border-radius: 24px;
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

const Line = styled.div`
  height: 1px;
  background-color: #e5e5ef;
  width: 100%;
  margin-top: 5px;
`;

const FeatureContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #acd1ef;
  padding: 8px 12px 8px 12px;
  border-radius: 12px;
`;

const Content = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;
const FeatureContentSecondary = styled(FeatureContent)`
  background-color: #f8e5c1;
  border-radius: 8px;
  padding: 10px;
`;

const FeatureContentTertiary = styled(FeatureContent)`
  background-color: #ffeeeeb2;
  border-radius: 8px;
  padding: 10px;
`;

const FeatureText = styled.div`
  color: #004988;
`;

const FeatureTextSecondary = styled(FeatureText)`
  color: #4f4f4f;
`;

const FeatureTextTertiary = styled(FeatureText)`
  color: #77213b;
`;

const Eclipse = styled.div`
  width: 8px;
  height: 8px;
  border: 1px solid #828282;
  border-radius: 50%;
`;

const Sessions = styled.div`
  background-color: #ffffff73;
  padding: 4px 8px 4px 8px;
  border-radius: 8px;
`;

const SessionsSecondary = styled(Sessions)`
  color: #4f4f4f;
`;

const SessionsTertiary = styled(Sessions)`
  color: #77213b;
`;
const ChartContainer = styled.div`
  margin-top: 10px;
`;
