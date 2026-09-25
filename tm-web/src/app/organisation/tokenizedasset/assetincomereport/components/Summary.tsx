import React from "react";
import styled from "styled-components";

const Summary = () => {
  const amountDetails = [
    {
      text: "Net Distributable Income",
      value: "₦400,000,000.00",
    },

    {
      text: "Total Tokens in Circulation ",
      value: "50,000,000",
    },
    {
      text: "Dividend Per Token",
      value: "₦8.00",
    },
  ];

  const approvers = ["onoja", "Nancy", "Obiaruku"];
  return (
    <>
      <div>
        <Heading>Summary</Heading>

        <Content>
          {amountDetails.map((item) => (
            <TextContent key={item.text}>
              <Text>{item.text}</Text>
              <Value>{item.value}</Value>
            </TextContent>
          ))}
        </Content>
      </div>

      <div>
        <Heading>Approvers</Heading>

        <Content2>
          <ApproversList>
            {approvers.map((item, idx) => (
              <ApproversNames key={idx}>{item}</ApproversNames>
            ))}
          </ApproversList>
          <ApproversNum>3/5 Approvers</ApproversNum>
        </Content2>
      </div>
    </>
  );
};

export default Summary;

const Content = styled.div`
  background-color: #f2f6f9;
  padding: 6px 16px 16px 16px;
  border-radius: 12px;
`;

const Content2 = styled.div`
  background-color: #f2f6f9;
  padding: 16px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
`;

const ApproversList = styled.div`
  display: flex;
  align-items: center;

  gap: 6px;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  margin-bottom: 6px;
  color: #00225a;
`;
const TextContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
`;

const Text = styled.p`
  font-weight: 500;
  font-size: 14px;
  letter-spacing: 0px;
  color: #828282;
  padding-bottom: 16px;
`;

const Value = styled.p`
  font-weight: 500;
  font-size: 14px;
  color: #00225a;
  text-align: center;
`;

const ApproversNames = styled.p`
  color: #ffffff;
  background-color: #007cdf;
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: 400;
  font-size: 14px;
`;

const ApproversNum = styled.p`
  font-size: 14px;
  color: #00225a;
`;
