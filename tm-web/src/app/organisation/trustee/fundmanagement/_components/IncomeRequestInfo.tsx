"use client";
import { MdVerified } from "react-icons/md";
import styled from "styled-components";
import DistributionBreakdown from "./DistributionBreakdown";

const IncomeRequestInfo = () => {
  return (
    <Container>
      <Card>
        <TopSection>
          <LeftTop>
            <AssetIcon />
            <div>
              <AssetName>Atlantis 1</AssetName>
            </div>
          </LeftTop>
        </TopSection>

        <Amount>1,000,000.00 NGN</Amount>

        <Grid>
          <InfoBlock>
            <Label>Request ID</Label>
            <Value>AE1</Value>
          </InfoBlock>

          <InfoBlock>
            <Label>Purpose</Label>
            <Value>Milestone Payment</Value>
          </InfoBlock>

          <InfoBlock>
            <Label>Status</Label>
            <StatusRow>
              <StatusDot />
              <StatusText>Pending Review</StatusText>
            </StatusRow>
          </InfoBlock>

          <InfoBlock>
            <Label>Requested by</Label>
            <UserRow>
              <UserAvatar />
              <div>
                <UserName>Florencez</UserName>
                <MdVerified color="#007CDF" size={10} />
              </div>
            </UserRow>
          </InfoBlock>

          <InfoBlock>
            <Label>Date</Label>
            <Value>20th Feb, 2026</Value>
          </InfoBlock>
        </Grid>
      </Card>

      <DistributionBreakdown />
    </Container>
  );
};

export default IncomeRequestInfo;

const Container = styled.div`
  margin-top: 16px;
  //   border-radius: 16px;
  //   padding: 24px;
  //   background: #fff;
`;

const Card = styled.div`
  margin-top: 16px;
`;

const TopSection = styled.div`
  display: flex;
  justify-content: space-between;
`;

const LeftTop = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

const AssetIcon = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #f2f6f9;
`;

const AssetName = styled.h3`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
`;

const Amount = styled.p`
  font-size: 20px;
  font-weight: 600;
  margin-top: 16px;
  color: #00225a;
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-top: 24px;
`;

const InfoBlock = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
`;

const Label = styled.p`
  font-size: 14px;
  font-weight: 400;
  color: #828282;
  margin: 0;
`;

const Value = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  margin: 8px 0;
  overflow-wrap: break-word;
  word-break: break-word;
  white-space: normal;
`;

const StatusRow = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

const StatusDot = styled.div`
  width: 6px;
  height: 6px;
  background: #f59e0b;
  border-radius: 50%;
`;

const StatusText = styled.span`
  font-size: 13px;
  color: #00225a;
`;

const UserRow = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const UserAvatar = styled.div`
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #007cdf;
`;

const UserName = styled.span`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
`;
