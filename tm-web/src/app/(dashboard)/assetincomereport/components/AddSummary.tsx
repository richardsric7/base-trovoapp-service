import React from "react";
import styled from "styled-components";

const AddSummary = () => {
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
        <FormGroup>
          <Label>
            Please request for approval from the following users to save this
            report.
          </Label>
          <FeeInputWrapper2>
            <InputText>Onoja</InputText>
            <InputText>Nancy</InputText>
            <InputText>Obiaruku</InputText>
          </FeeInputWrapper2>
        </FormGroup>
      </div>
    </>
  );
};

export default AddSummary;

const Content = styled.div`
  background-color: #f2f6f9;
  padding: 6px 16px 16px 16px;
  border-radius: 12px;
`;

const Heading = styled.h1`
  font-size: 16px;
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
  color: #191919;
  text-align: center;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const FeeInputWrapper2 = styled.div`
  padding: 10px;
  border: 1px solid #007cdf;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
`;

const InputText = styled.p`
  color: #ffffff;
  background-color: #007cdf;
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: 400;
  font-size: 14px;
`;

const Label = styled.p`
  color: #191919;
  font-weight: 500;
  font-size: 14px;
  margin-top: 14px;
`;
