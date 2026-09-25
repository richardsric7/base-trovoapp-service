"use client";

import React from "react";
import styled from "styled-components";
import dynamic from "next/dynamic";
import type { ApexOptions } from "apexcharts";
import { AmPortfolioSummary } from "@/redux/api/assetManager";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

interface PortfolioSummaryProps {
  summary?: AmPortfolioSummary;
}

const AmPortfolioSummaryComponent = ({ summary }: PortfolioSummaryProps) => {
  const categories = summary?.categories || [];
  const series = categories.map((cat) => cat.asset_count ?? 0);
  const labels = categories.map((cat) => cat.category || "Uncategorized");
  const totalAssets =
    summary?.asset_count ??
    categories.reduce((total, cat) => total + (cat.asset_count ?? 0), 0);

  const colors = ["#1877F2", "#5DA3D9", "#9BBBD6", "#0C4A87"];

  const options: ApexOptions = {
    chart: {
      type: "donut",
    },

    colors,

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
              label: "Assets managed",
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

      {categories.length === 0 ? (
        <EmptyState>No data available</EmptyState>
      ) : (
        <>
          <ChartWrapper>
            <Chart
              options={options}
              series={series}
              type="donut"
              height={260}
            />
          </ChartWrapper>

          <BreakdownHeader>
            <span>Category breakdown</span>
          </BreakdownHeader>

          <List>
            {categories.map((cat, index) => {
              const mainVal = cat.values?.[0];
              const valText = mainVal
                ? `${mainVal.currency || ""} ${Number(
                    mainVal.amount || 0,
                  ).toLocaleString()}`
                : `${cat.asset_count ?? 0} asset(s)`;

              return (
                <Row key={`${cat.category}-${index}`}>
                  <Left>
                    <Dot color={colors[index % colors.length]} />
                    {cat.category || "Uncategorized"}
                  </Left>
                  <Value>{valText}</Value>
                </Row>
              );
            })}
          </List>
        </>
      )}
    </Container>
  );
};

export default AmPortfolioSummaryComponent;

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

const EmptyState = styled.div`
  padding: 40px 0;
  text-align: center;
  color: #8a94a6;
  font-size: 14px;
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
