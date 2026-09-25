"use client";

import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import PrimaryButton from "@/components/PrimaryButton";
import { useRouter } from "next/navigation";
import { useMemo, useState } from "react";
import styled from "styled-components";

const ApikeysPage = () => {
  const [currentPage, setCurrentPage] = useState(1);

  const [pageSize, setPageSize] = useState(10);

  const router = useRouter();
  const totalUsers = 100;

  const columns = [
    {
      title: "Key Name",
      dataIndex: "keyName",
      key: "keyName",
    },

    {
      title: "Key Type",
      dataIndex: "keyType",
      key: "keyType",
    },

    {
      title: "Date Created",
      dataIndex: "date",
      key: "date",
    },
    {
      title: "Created By",
      dataIndex: "createdBy",
      render: (_: any, record: any) => {
        return (
          <UserIdSection>
            <Avatar />
            <div>
              <UserName>{record.name}</UserName>
              <UserEmail>{record.email}</UserEmail>
            </div>
          </UserIdSection>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 10 }, (_, index) => ({
      userId: index + 1,
      name: `John Doe`,
      email: `john@gmail.com`,
      keyName: `Key ${index + 1}`,
      keyType: `API`,
      date: "23 Sep 2025",
    }));
  }, []);
  return (
    <Container>
      <HeaderRow>
        <Title>API Keys</Title>
        <PrimaryButton
          onClick={() => router.push("/external-api-clients/api-keys/create")}
          buttonStyle={{
            width: "fit-content",
            flexShrink: 0,
            padding: "0px 16px",
          }}
        >
          Create API Key
        </PrimaryButton>
      </HeaderRow>
      <div>
        <CustomTable columns={columns} dataSource={dataSource} />
        <Pagination
          totalCount={totalUsers}
          currentPage={currentPage}
          pageSize={pageSize}
          onPageChange={setCurrentPage}
          onPageSizeChange={setPageSize}
        />
      </div>
    </Container>
  );
};

export default ApikeysPage;

const Container = styled.section`
  background-color: #fff;
  gap: 20px;
  border-radius: 16px;
  padding: 24px;

  width: 100%;

  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 4%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 8%;
    text-align: center;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 10%;
    text-align: center;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 6%;
    text-align: center;
  }

  th {
    color: #828282;
    font-size: 14px;
    font-weight: 500;
  }
  td {
    color: #00225a;
    font-weight: 400;
  }
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  margin: 0;
  flex: 1;
  color: #00225a;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;
const HeaderRow = styled.div`
  display: flex;
  align-items: center;
  margin-bottom: 20px;

  justify-content: space-between;

  width: 100%;
`;

const UserIdSection = styled.div`
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

const UserEmail = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;
