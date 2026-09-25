"use client";

import styled from "styled-components";
import { AiOutlineCheckCircle } from "react-icons/ai";

const dummyData = [
  {
    date: "20th Feb, 2026",
    value: "10,000,000.00 NGN",
    source: "Asset Manager",
    status: "Verified",
  },
  {
    date: "20th Feb, 2026",
    value: "10,000,000.00 NGN",
    source: "Asset Manager",
    status: "Verified",
  },
  {
    date: "20th Feb, 2026",
    value: "10,000,000.00 NGN",
    source: "Asset Manager",
    status: "Verified",
  },
  {
    date: "20th Feb, 2026",
    value: "10,000,000.00 NGN",
    source: "Asset Manager",
    status: "Verified",
  },
  {
    date: "20th Feb, 2026",
    value: "10,000,000.00 NGN",
    source: "Asset Manager",
    status: "Verified",
  },
];

const ValuationHistory = () => {
  return (
    <Wrapper>
      <Title>Valuation History</Title>

      <Table>
        <Header>
          <span>Date</span>
          <span>Value</span>
          <span>Source</span>
          <span>Status</span>
        </Header>

        {dummyData.map((item, index) => (
          <Row key={index}>
            <Date>{item.date}</Date>
            <Value>{item.value}</Value>
            <Source>{item.source}</Source>

            <Status>
              <AiOutlineCheckCircle size={14} />
              {item.status}
            </Status>
          </Row>
        ))}
      </Table>
    </Wrapper>
  );
};

export default ValuationHistory;

const Wrapper = styled.div`
  margin-top: 32px;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin-bottom: 16px;
`;

const Table = styled.div`
  background: #ffffff;
  border: 1px solid #e6e6e6;
  border-radius: 12px;
`;

const Header = styled.div`
  display: grid;
  grid-template-columns: 2fr 2fr 2fr 1fr;
  padding: 16px 20px;
  font-size: 13px;
  color: #828282;
  border-bottom: 1px solid #eeeeee;
`;

const Row = styled.div`
  display: grid;
  grid-template-columns: 2fr 2fr 2fr 1fr;
  padding: 18px 20px;
  border-bottom: 1px solid #f1f1f1;

  &:last-child {
    border-bottom: none;
  }
`;

const Date = styled.div`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
`;

const Value = styled.div`
  font-size: 14px;

  font-weight: 500;
  color: #00225a;
`;

const Source = styled.div`
  font-size: 14px;

  font-weight: 500;
  color: #00225a;
`;

const Status = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
  background: #dff5ea;
  color: #1c9e67;
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 6px;
  width: fit-content;
`;
