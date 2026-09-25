import styled from "styled-components";

const DistributionBreakdown = () => {
  return (
    <Wrapper>
      <Title>Distribution Breakdown</Title>

      <Card>
        <TopStats>
          <Stat>
            <Label>Total Tokens</Label>
            <Value>10,000,000.00</Value>
          </Stat>

          <Stat>
            <Label>Token Holders</Label>
            <Value>1,546</Value>
          </Stat>

          <Stat>
            <Label>Distribution per Token</Label>
            <Value>₦4.20</Value>
          </Stat>
        </TopStats>

        <BreakdownBox>
          <Row>
            <RowLabel>Net Income:</RowLabel>
            <RowValue>₦42,000,000</RowValue>
          </Row>

          <Row>
            <RowLabel>÷ Total Tokens:</RowLabel>
            <RowValue>10,000,000</RowValue>
          </Row>

          <Divider />

          <Row>
            <RowLabel>= Distribution per Token:</RowLabel>
            <HighlightValue>₦4.20</HighlightValue>
          </Row>
        </BreakdownBox>
      </Card>
    </Wrapper>
  );
};

export default DistributionBreakdown;

const Wrapper = styled.div`
  margin-top: 24px;
`;

const Title = styled.h3`
  font-size: 18px;
  font-weight: 600;
  color: #1a2b49;
  margin-bottom: 16px;
`;

const Card = styled.div`
  background: #ffffff;
  border-radius: 12px;
  padding: 10px;
  border: 1px solid #eef2f6;
`;

const TopStats = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
`;

const Stat = styled.div`
  display: flex;
  flex-direction: column;
`;

const Label = styled.div`
  font-size: 12px;
  color: #8a94a6;
  margin-bottom: 4px;
`;

const Value = styled.div`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
`;

const BreakdownBox = styled.div`
  background: #eef2f6;
  border-radius: 10px;
  padding: 16px;
`;

const Row = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
`;

const RowLabel = styled.div`
  font-size: 14px;
  color: #6b7280;
`;

const RowValue = styled.div`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
`;

const HighlightValue = styled.div`
  font-size: 14px;
  font-weight: 700;
  color: #00225a;
`;

const Divider = styled.div`
  height: 1px;
  background: #d1d5db;
  margin: 10px 0;
`;
