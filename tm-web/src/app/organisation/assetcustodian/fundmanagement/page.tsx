"use client";
import styled from "styled-components";
import FundCard from "./_components/FundCard";
import FundRelease from "./_components/FundRelease";

const FundManagementPage = () => {
  const fundCardContent = [
    {
      title: "Total Amount from Primary Sale ",
      amount: "—",
    },

    {
      title: "Total Income from Asset",
      amount: "—",
    },
    {
      title: "Total Amount Paid to Issuers",
      amount: "—",
    },
    {
      title: "Total Paid to Investors ",
      amount: "—",
    },
    {
      title: "Milestone Payment Balance ",
      amount: "—",
    },
    {
      title: "Total Fees Generated",
      amount: "—",
    },
  ];
  return (
    <Container>
      <CardContainer>
        <FundCard title="Total Amount Processed" amount="—" />
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
