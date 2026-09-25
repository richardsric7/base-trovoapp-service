"use client";

import React from "react";
import styled from "styled-components";
import dynamic from "next/dynamic";
import type { ApexOptions } from "apexcharts";
import { AssetManagerPortfolioItem } from "@/redux/api/assetManager";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

interface PortfolioSummaryProps {
  items?: AssetManagerPortfolioItem[];
}

const PortfolioSummary = ({ items = [] }: PortfolioSummaryProps) => {
  const series = items.map((item) => item.percentage ?? item.count ?? 0);
  const labels = items.map(
    (item) => item.asset_sector || item.label || "Uncategorized",
  );
  const totalAssets = items.reduce((total, item) => total + (item.count ?? 0), 0);

  const options: ApexOptions = {
    chart: {
      type: "donut",
    },

    colors: ["#1877F2", "#5DA3D9", "#9BBBD6", "#0C4A87"],

    labels,

    dataLabels: {
      enabled: true,
      formatter: (val) => `${Math.round(Number(val))}%`,
      style: {
        fontSize: "14px",
        colors: ["#fff"],
      },
    },

    plotOptions: {
      pie: {
        donut: {
          size: "65%",
          labels: {
            show: true,
            name: {
              show: false,
            },
            value: {
              show: true,
              fontSize: "32px",
              fontWeight: 700,
              offsetY: -5,
              formatter: () => totalAssets.toLocaleString(),
            },
            total: {
              show: true,
              label: "Assets in custody",
              fontSize: "12px",
              color: "#6b7280",
            },
          },
        },
      },
    },

    legend: {
      show: false,
    },

    stroke: {
      width: 0,
    },
  };

  return (
    <Container>
      <Title>Portfolio Summary</Title>

      <ChartWrapper>
        {series.length > 0 ? (
          <Chart options={options} series={series} type="donut" height={260} />
        ) : (
          <EmptyText>No portfolio data available</EmptyText>
        )}
      </ChartWrapper>

      <BreakdownHeader>
        <span>Category breakdown</span>

        <Arrows>
          <Arrow>{"‹"}</Arrow>
          <Arrow>{"›"}</Arrow>
        </Arrows>
      </BreakdownHeader>

      <List>
        {items.map((item, index) => (
          <Row key={`${item.asset_sector || item.label}-${index}`}>
            <Left>
              <Dot color={["#0C4A87", "#1877F2", "#5DA3D9", "#9BBBD6"][index % 4]} />
              {item.asset_sector || item.label || "Uncategorized"}
            </Left>
            <Value>{item.count ?? 0}</Value>
          </Row>
        ))}
      </List>
    </Container>
  );
};

export default PortfolioSummary;

const Container = styled.div`
  background: #fff;
  padding: 24px;
  border-radius: 16px;
  width: 380px;
`;

const Title = styled.h3`
  font-size: 18px;
  font-weight: 600;
  color: #1a2b49;
  margin-bottom: 20px;
`;

const ChartWrapper = styled.div`
  display: flex;
  justify-content: center;
`;

const EmptyText = styled.p`
  min-height: 260px;
  display: flex;
  align-items: center;
  color: #8a94a6;
`;

const BreakdownHeader = styled.div`
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
  font-weight: 600;
  color: #1a2b49;
`;

const Arrows = styled.div`
  display: flex;
  gap: 8px;
  color: #007cdf;
`;

const Arrow = styled.div`
  width: 28px;
  height: 28px;
  background: #f2f6f9;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
`;

const List = styled.div`
  margin-top: 12px;
`;

const Row = styled.div`
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #e5e5ef;
`;

const Left = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  color: #00225a;
  font-size: 14px;
`;

const Dot = styled.div<{ color: string }>`
  width: 10px;
  height: 10px;
  border-radius: 3px;
  background: ${(props) => props.color};
`;

const Value = styled.div`
  font-weight: 600;
  color: #00225a;
  font-size: 14px;
`;
