"use client";

import CustomTable from "@/components/CustomTable";
import React, { useMemo } from "react";
import styled from "styled-components";

const RecentTokenization = () => {
  const columns = [
    {
      title: "",
      dataIndex: "name",
      key: "name",
      width: "100%",
      render: (_: any, record: any) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserName>{record.name || "N/A"}</UserName>
            <SubName>{record.subname || "N/A"}</SubName>
          </div>
        </UserInfoSection>
      ),
    },
    {
      title: "",
      dataIndex: "assetType",
      key: "assetType",
      width: "100%",
      render: (_: any, record: any) => (
        <AssetType>{record.assetType || "N/A"}</AssetType>
      ),
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 5 }, (_, index) => ({
      tokenId: index + 1,
      name: `Atalntis ${index + 1}`,
      subname: `Tokenized on 20 Jul, 2025`,
      updatedOn: "20 jul, 2025",
      assetType: "Real Estate",
    }));
  }, []);

  const activityData = [
    {
      user: "John Doe",
      action: "approved",
      asset: "Atlantis Estate",
      date: "20th July 2025",
    },
    {
      user: "Jane Smith",
      action: "rejected",
      asset: "Oakwood Plaza",
      date: "21st July 2025",
    },
    {
      user: "John Doe",
      action: "approved",
      asset: "Atlantis Estate",
      date: "20th July 2025",
    },
    {
      user: "John Doe",
      action: "approved",
      asset: "Atlantis Estate",
      date: "20th July 2025",
    },
  ];
  return (
    <PageContainer>
      <Container>
        <SubTitle>Recent Tokenization Requests</SubTitle>
        <CustomTable columns={columns} dataSource={dataSource} />
      </Container>
      <ActivityCard>
        {activityData.map((item, idx) => (
          <React.Fragment key={idx}>
            <ActivityContent>
              <Avatar2 />
              <div>
                <ActivityText>
                  <StyledSpan>{item.user} </StyledSpan>
                  {item.action} tokenized asset{" "}
                  <StyledSpan> {item.asset}</StyledSpan>
                </ActivityText>
                <ActivityDate>{item.date}</ActivityDate>
              </div>
            </ActivityContent>
            {idx !== activityData.length - 1 && <Divider />}
          </React.Fragment>
        ))}
      </ActivityCard>
    </PageContainer>
  );
};

export default RecentTokenization;

const PageContainer = styled.div`
  width: 60%;
`;

const Container = styled.div`
  background-color: #fff;
  padding: 24px;
  gap: 20px;
  border-radius: 24px;

  table {
    width: 100%;
    table-layout: fixed;

    th:nth-child(1),
    td:nth-child(1) {
      width: 25%;
      //  text-align: center;
    }

    th:nth-child(2),
    td:nth-child(2) {
      width: 10%;

      //  text-align: center;
    }
  }
`;

const SubTitle = styled.h3`
  font-size: 20px;
  font-weight: 600;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;

const UserInfoSection = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

const UserName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const SubName = styled.p`
  font-size: 10px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const AssetType = styled.p`
  font-size: 12px;
  font-weight: 400;
  background-color: #f2f6f9;
  color: #007cdf;
  padding: 8px;
  display: inline-flex;
`;

const Avatar = styled.div`
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const ActivityContent = styled.div`
  display: flex;
  gap: 4px;
  align-items: center;
  padding-bottom: 8px;
`;

const ActivityCard = styled.div`
  background-color: #ffffff;
  border-radius: 24px;
  padding: 24px;
  margin-top: 10px;
`;

const ActivityText = styled.p`
  font-weight: 400;
  font-size: 14px;
  color: #00225a;
  letter-spacing: 0.1px;
`;

const StyledSpan = styled.span`
  font-weight: 600;
  font-size: 14px;
  letter-spacing: 0.1px;
`;

const ActivityDate = styled.p`
  color: #828282;
  font-weight: 400;
  font-size: 10px;
  letter-spacing: 0.1px;
`;

const Divider = styled.div`
  height: 1px;
  background: #f0f0f0;
  margin: 8px 0;
  width: 100%;
`;

const Avatar2 = styled.div`
  width: 32px;
  height: 32px;
  min-width: 32px;
  min-height: 32px;
  max-width: 32px;
  max-height: 32px;
  border-radius: 50%;
  background-color: rebeccapurple;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
`;
