import React from "react";
import styled from "styled-components";
import StatCard from "../../components/StatsCard";
import {
  CustodianCurrencyAmount,
  CustodianDashboardSummary,
  toCurrencyAmounts,
} from "@/redux/api/assetCustodian";

interface AcStatsGridProps {
  summary?: CustodianDashboardSummary;
}

export const formatCurrency = (amount: number, currency: string) => {
  if (!currency) return amount.toLocaleString();

  try {
    return new Intl.NumberFormat("en-NG", {
      style: "currency",
      currency,
      maximumFractionDigits: 2,
    }).format(amount);
  } catch {
    return `${currency} ${amount.toLocaleString()}`;
  }
};

const CurrencyStat = ({
  label,
  amounts,
  footer,
}: {
  label: string;
  amounts?: CustodianCurrencyAmount[];
  footer?: string;
}) => {
  const [primary, ...rest] = toCurrencyAmounts(amounts);

  return (
    <Card>
      <Label>{label}</Label>
      <Value>
        {primary ? formatCurrency(primary.amount, primary.currency) : "—"}
      </Value>

      {rest.length > 0 && (
        <OtherCurrencies>
          {rest.map((item) => (
            <OtherCurrency key={item.currency}>
              {formatCurrency(item.amount, item.currency)}
            </OtherCurrency>
          ))}
        </OtherCurrencies>
      )}

      {footer && <Footer>{footer}</Footer>}
    </Card>
  );
};

const AcStatsGrid = ({ summary }: AcStatsGridProps) => {
  const segregatedAccounts = summary?.segregated_accounts ?? 0;

  return (
    <Grid>
      <StatCard
        label="Total Asset Count"
        value={summary?.total_asset_count ?? 0}
      />
      <CurrencyStat
        label="Assets Under Custody"
        amounts={summary?.assets_under_custody}
      />
      <CurrencyStat label="Fees Generated" amounts={summary?.fees_generated} />
      <CurrencyStat
        label="Account Balances"
        amounts={summary?.account_balances}
        footer={`${segregatedAccounts} segregated account${
          segregatedAccounts === 1 ? "" : "s"
        }`}
      />
      <StatCard
        label="Fund Release Requests"
        value={summary?.pending_fund_releases ?? 0}
        pending={summary?.pending_fund_releases ?? 0}
        showFooter
      />
      <StatCard
        label="Compliance Items"
        value={summary?.pending_compliance_items ?? 0}
        pending={summary?.pending_compliance_items ?? 0}
        showFooter
      />
    </Grid>
  );
};

export default AcStatsGrid;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
`;

const Card = styled.div`
  background: #f2f6f9;
  padding: 16px;
  border-radius: 10px;
`;

const Label = styled.div`
  font-size: 12px;
  color: #9a9a9a;
`;

const Value = styled.div`
  font-size: 20px;
  font-weight: 600;
  margin-top: 4px;
  color: #00225a;
`;

const OtherCurrencies = styled.div`
  margin-top: 6px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
`;

const OtherCurrency = styled.span`
  font-size: 12px;
  color: #00225a;
  background: #ffffff;
  border-radius: 6px;
  padding: 4px 8px;
`;

const Footer = styled.div`
  margin-top: 10px;
  font-size: 12px;
  color: #8a94a6;
`;
