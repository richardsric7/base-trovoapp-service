"use client";

import { SearchBar } from "@/components";
import CustomerSupportTable from "./component/CustomerSupportTable";
import styled from "styled-components";
import SupportFilter from "./component/SupportFilter";

const CustomerSupportPage = () => {
  return (
    <Wrapper>
      <Header>
        <TitleSection>
          <Title>Support Tickets</Title>
          <TotalCount>
            Total: <span>500 Ticket</span>
          </TotalCount>
        </TitleSection>

        <StatsHeader>
          <StatItem color="#00225a">
            <h3>256</h3>
            <p>Resolved</p>
          </StatItem>
          <Line />
          <StatItem color="#00225a">
            <h3>21</h3>
            <p>Pending</p>
          </StatItem>
          <Line />
          <StatItem color="#BE3800">
            <h3>223</h3>
            <p>Open</p>
          </StatItem>
        </StatsHeader>
      </Header>
      <SubHeader>
        <SearchBar />
        <SupportFilter />
      </SubHeader>
      <CustomerSupportTable />
    </Wrapper>
  );
};

export default CustomerSupportPage;

const Wrapper = styled.section`
  background-color: #ffffff;
  padding: 32px;
  border-radius: 24px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const SubHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 25px;
  margin: 10px 0 30px 0;
`;

const TitleSection = styled.div``;

const Title = styled.h1`
  color: #00225a;
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  line-height: 30px;
`;

const TotalCount = styled.p`
  color: #828282;
  font-size: 16px;
  font-weight: 400;
  margin: 0;
  line-height: 28px;

  span {
    color: #00225a;
    font-size: 16px;
    margin: 0;
    font-weight: 500;
  }
`;
const StatsHeader = styled.div`
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
`;

const StatItem = styled.div<{ color: string }>`
  text-align: center;
  flex: 1;

  h3 {
    color: ${({ color }) => color};
    font-size: 24px;
    margin: 0;
  }

  p {
    font-size: 14px;
    color: #828282;
  }
`;

const Line = styled.div`
  width: 1px;
  height: 40px;
  background-color: #e5e5ef;
  align-self: center;
`;
