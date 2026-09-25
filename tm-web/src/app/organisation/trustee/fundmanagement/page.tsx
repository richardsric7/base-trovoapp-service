"use client";
import styled from "styled-components";
import FundCard from "./_components/FundCard";
import FundRelease from "./_components/FundRelease";
import Loader from "@/components/Loader";
import {
  ICurrencyTotal,
  useGetTrusteeFundManagementSummaryQuery,
} from "@/redux/api/trustees";

const formatTotals = (
  totals: ICurrencyTotal[] | undefined,
  maximumSignificantDigits?: number
) => {
  if (!totals?.length) return "0";

  return totals
    .map(({ amount, currency }) => {
      const numericAmount = Number(amount);
      const formattedAmount = Number.isFinite(numericAmount)
        ? numericAmount.toLocaleString(undefined, {
            ...(maximumSignificantDigits
              ? { maximumSignificantDigits }
              : { maximumFractionDigits: 20 }),
          })
        : amount;

      return `${formattedAmount} ${currency}`.trim();
    })
    .join(" • ");
};

const FundManagementPage = () => {
  const { data, isLoading, isError, refetch } =
    useGetTrusteeFundManagementSummaryQuery();
  const summary = data?.data;
  const fundCardContent = [
    {
      title: "Total Amount from Primary Sale ",
      amount: formatTotals(summary?.total_amount_from_primary_sales),
    },

    {
      title: "Total Income from Asset",
      amount: formatTotals(summary?.total_income_from_assets),
    },
    {
      title: "Total Amount Paid to Issuers",
      amount: formatTotals(summary?.total_amount_paid_to_issuers),
    },
    {
      title: "Total Paid to Investors ",
      amount: formatTotals(summary?.total_paid_to_investors),
    },
    {
      title: "Milestone Payment Balance ",
      amount: formatTotals(summary?.milestone_payment_balance),
    },
    {
      title: "Total Fees Generated",
      amount: formatTotals(summary?.total_fees_generated, 15),
    },
  ];

  if (isLoading) return <Loader />;

  return (
    <Container>
      {isError && (
        <ErrorState>
          <span>Unable to load the fund management summary.</span>
          <RetryButton type="button" onClick={() => refetch()}>
            Try again
          </RetryButton>
        </ErrorState>
      )}
      <CardContainer>
        <FundCard
          title="Total Amount Processed"
          amount={formatTotals(summary?.total_amount_processed)}
        />
        <CardContent>
          {fundCardContent.map((item, idx) => (
            <FundCard title={item.title} amount={item.amount} key={idx} />
          ))}
        </CardContent>
      </CardContainer>

      <FundRelease />
    </Container>
  );
};

export default FundManagementPage;

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 30px;
`;
const CardContainer = styled.div`
  display: flex;
  gap: 16px;
  width: 100%;
`;

const CardContent = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  flex: 1;
`;

const ErrorState = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 16px;
  border: 1px solid #fecdca;
  border-radius: 10px;
  background: #fef3f2;
  color: #b42318;
`;

const RetryButton = styled.button`
  border: 1px solid #b42318;
  border-radius: 8px;
  padding: 8px 14px;
  background: #fff;
  color: #b42318;
  font: inherit;
  cursor: pointer;
`;
