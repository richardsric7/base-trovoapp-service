"use client";
import CustomTable from "@/components/CustomTable";
import { useReferrerCountQuery } from "@/redux/api/usersmetrics";
import Link from "next/link";
import { useRouter } from "next/navigation";
import React, { useMemo } from "react";
import styled from "styled-components";

interface ReferrerData {
  referrer: string;
  count: number;
}

interface TableData {
  key: number;
  username: string; // map referrer to this
  count: number;
}

const TopReferrals = () => {
  const { data, isLoading } = useReferrerCountQuery();

  // console.log("ref", data);
  const router = useRouter();
  const columns = [
    {
      title: "Username",
      dataIndex: "username",
      key: "username",
      width: "100%",
      render: (_: any, record: TableData) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserName>
              {/* {record.username !== "noreferrer" ? record.username : "N/A"} */}

              {record.username || "N/A"}
            </UserName>
          </div>
        </UserInfoSection>
      ),
    },
    {
      title: "Total Referrals",
      dataIndex: "count",
      key: "count",
      width: "100%",
      render: (count: number) => (
        <Count>{count ? `${count} Referrals` : "N/A"}</Count>
      ),
    },
  ];

  const dataSource = useMemo(() => {
    if (!data || isLoading) {
      return [];
    }
    return data?.data?.map((item: ReferrerData, index: number) => ({
      key: index + 1,
      username: item.referrer,
      count: item.count,
    }));
  }, [data, isLoading]);

  return (
    <Container>
      <SubTitle>Top Referrals</SubTitle>
      <CustomTable
        columns={columns}
        dataSource={dataSource}
        onRowClick={(record) => {
          router.push(`/users/${record.username}`);
        }}
      />
    </Container>
  );
};

export default TopReferrals;

const Container = styled.div`
  background-color: #fff;
  padding: 40px 32px;
  gap: 20px;
  border-radius: 24px;
  width: 60%;
  table {
    width: 100%;
    table-layout: fixed;

    th:nth-child(1),
    td:nth-child(1) {
      width: 25%;
    }

    th:nth-child(2),
    td:nth-child(2) {
      width: 10%;
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

const Count = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const StyledButton = styled(Link)`
  background: none;
  border: none;
  padding: 0;
  margin: 0;
  text-align: left;
  cursor: pointer;
  text-decoration: none;
  display: block;
  color: #00225a;
`;
