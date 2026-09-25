"use client";

import styled from "styled-components";
import FundCard from "@/app/organisation/assetcustodian/fundmanagement/_components/FundCard";
import FundRelease from "./FundRelease";

const OrgFundManagement = () => {
  const fundCardContent = [
    {
      title: "Total Amount from Primary Sale ",
      amount: "8,000 NGN",
    },

    {
      title: "Total Income from Asset",
      amount: "1,000 NGN",
    },
    {
      title: "Total Amount Paid to Issuers",
      amount: "8,000 NGN",
    },
    {
      title: "Total Paid to Investors ",
      amount: "10,000 NGN",
    },
    {
      title: "Milestone Payment Balance ",
      amount: "2,000 NGN",
    },
    {
      title: "Total Fees Generated",
      amount: "8,000 NGN",
    },
  ];

  return (
    <Container>
      <CardContainer>
        <FundCard title="Total Amount Processed" amount="10,000 NGN" />
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

export default OrgFundManagement;

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
