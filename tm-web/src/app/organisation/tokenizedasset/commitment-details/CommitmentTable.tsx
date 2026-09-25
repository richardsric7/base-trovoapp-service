import { SearchBar } from "@/components";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import React, { useMemo, useState } from "react";
import styled from "styled-components";

interface Commitment {
  userId: number;
  name: string;
  subname: string;
  commitedAmount: string;
  date: string;
}

const CommitmentTable = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const columns = [
    {
      title: "Trovo User ID",
      dataIndex: "userId",
      render: (_: any, record: any) => {
        return (
          <UserIdSection>
            <Avatar />
            <div>
              <UserName>{record.name}</UserName>
              <UserSubName>{record.subname}</UserSubName>
            </div>
          </UserIdSection>
        );
      },
    },
    {
      title: "Commitment Amount",
      dataIndex: "commitedAmount",
      key: "commitedAmount",
    },

    {
      title: "Date",
      dataIndex: "date",
      key: "date",
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 10 }, (_, index) => ({
      userId: index + 1,
      name: `Florence ${index + 1}`,
      subname: `florence Zach`,
      commitedAmount: "1,000 CNGN",
      date: "23 Oct, 2023",
    }));
  }, []);

  const paginatedData = useMemo(() => {
    const start = (currentPage - 1) * pageSize;
    return dataSource.slice(start, start + pageSize);
  }, [dataSource, currentPage, pageSize]);

  return (
    <Container>
      <Title>Commitments History</Title>
      <SearchBar />
      <CustomTable columns={columns} dataSource={paginatedData} />
      <Pagination
        currentPage={currentPage}
        totalCount={dataSource.length}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />
    </Container>
  );
};

export default CommitmentTable;

const Container = styled.section`
  background-color: #fff;
  gap: 20px;
  border-radius: 24px;
  padding: 32px;
  width: 100%;

  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 18%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 19%;
    text-align: center;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 20%;
    text-align: center;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 12%;
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
  margin-bottom: 10px;
  color: #00225a;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
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

const UserSubName = styled.p`
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;
