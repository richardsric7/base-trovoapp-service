"use client";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import TruncatedText from "@/hooks/useTruncate";
import { useRecentRegQuery } from "@/redux/api/usersmetrics";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Router } from "next/router";
import { useMemo, useState } from "react";
import styled from "styled-components";

const RecentRegistration = () => {
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const { data, isLoading, isFetching } = useRecentRegQuery(
    {
      page: currentPage,
      pageSize,
    },
    {
      skip: false, // Only fetch when needed
      refetchOnMountOrArgChange: false, // Prevent automatic refetching
    }
  );
  console.log("recent", data);
  const columns = [
    {
      title: "Name",
      dataIndex: "full_name",
      render: (_: any, record: any) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserName>
              <TruncatedText text={`${record.full_name}`} maxLength={25} />
            </UserName>
          </div>
        </UserInfoSection>
      ),
    },
    {
      title: "UserName",
      dataIndex: "username",
      key: "username",
      render: (_: any, record: any) => (
        <UserEmail>
          <TruncatedText text={record.username} maxLength={15} />
        </UserEmail>
      ),
    },

    {
      title: "Public key",
      dataIndex: "public_key",
      key: "public_key",
      render: (_: any, record: any) => (
        <PublicKey>
          <TruncatedText text={record.public_key} maxLength={15} />
        </PublicKey>
      ),
    },

    {
      title: "Email Address",
      dataIndex: "email",
      key: "email",
      render: (_: any, record: any) => (
        <TruncatedText text={record.email} maxLength={25} />
      ),
    },
    {
      title: "Phone number",
      dataIndex: "phone_number",
      key: "number",
    },
    {
      title: "Date",
      dataIndex: "date",
      render: (date: string) =>
        new Date(date).toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        }),
      key: "created_at",
    },
  ];

  const dataSource = useMemo(() => {
    // console.log("Raw API Data:", data);

    if (!data || !Array.isArray(data.data.data)) {
      console.error("Invalid data format:", data);
      return [];
    }

    return data.data.data.map((user) => ({
      ...user,
      key: user.id,
      full_name: user.full_name,
      username: user.username,
      public_key: user.public_key,
      email: user.email,
      phone_number: user.phone_number,
      date: user.date,
    }));
  }, [data]);

  const handlePageChange = (pageNumber: number) => {
    setCurrentPage(pageNumber); // This will trigger data fetching for the specific page
  };
  if (isLoading) {
    return (
      <Container>
        <UserName>Loading users...</UserName>
      </Container>
    );
  }

  return (
    <Container>
      <Title>Recent Registrations</Title>
      <>
        <CustomTable
          columns={columns}
          dataSource={dataSource}
          onRowClick={(record) => {
            router.push(`/users/${record.username}`);
          }}
          // isLoading={isLoading}
        />
        <Pagination
          totalCount={data?.data?.total || 0}
          onPageChange={handlePageChange}
          currentPage={currentPage}
          pageSize={pageSize}
          onPageSizeChange={setPageSize}
          isFetching={isFetching}
        />
      </>
    </Container>
  );
};

export default RecentRegistration;

const Container = styled.section`
  padding: 32px;
  background-color: #ffffff;
  border-radius: 12px;
margin-top:20px;
  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 11%;
    //  text-align: center;
   }

  th:nth-child(2),
  td:nth-child(2) {
    width: 6%;
    
    //  text-align: center;
  }

  th:nth-child(3),
  td:nth-child(3) {
    width: 8%;
    //  text-align: center;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 8%;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 8%;
    text-align: center;
  }

  th:nth-child(6),
  td:nth-child(6) {
    width: 4%;
    text-align: center;
  }

 

  th {
    color: #828282;
    font-size: 14px;
    font-weight: 500;
    text-align center;
  }:
`;

const Title = styled.h2`
  color: #00225a;
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  text-align: left;
  margin-bottom: 8px;
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
const UserInfoSection = styled.div`
  display: flex;
  column-gap: 10px;
  align-items: center;
  cursor: pointer;
`;
const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  max-width: 40px;
  background-color: rebeccapurple;
`;
const UserName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
`;
const UserEmail = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  margin: 0;
`;

const PublicKey = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #007cdf;
`;
