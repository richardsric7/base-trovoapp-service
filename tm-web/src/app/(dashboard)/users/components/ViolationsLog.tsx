import { SearchBar } from "@/components";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import React, { useMemo, useState } from "react";
import { MdVerified } from "react-icons/md";
import styled from "styled-components";
import UsersFilter from "./UsersFilter";
import { FaX } from "react-icons/fa6";

interface ViolationData {
  userId: number;
  name: string;
  subname: string;
  voilationtype: string;
  description: string;
  status: string;
  date: string;
}
interface ColumnType {
  title: string;
  dataIndex: string;
  key?: string;
  width?: string;
  render?: (text: string, record: ViolationData) => React.ReactNode;
}

const ViolationsLog = () => {
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [selectedFilters, setSelectedFilters] = useState<
    Record<string, string>
  >({});

  const handleRemoveFilter = (key: string) => {
    setSelectedFilters((prev) => {
      const newFilters = { ...prev };
      delete newFilters[key];
      return newFilters;
    });
  };
  const columns: ColumnType[] = [
    {
      title: "Trovo User ID",
      dataIndex: "trovouserId",
      width: "20%",
      render: (_, record: ViolationData) => {
        return (
          <ViolationWrapper>
            <ViolationIdSection>
              <Avatar />
              <div>
                <ViolationName>
                  {record.name}
                  <VerifiedIcon>
                    <MdVerified />
                  </VerifiedIcon>
                </ViolationName>
                <ViolationSubName>{record.subname}</ViolationSubName>
              </div>
            </ViolationIdSection>
          </ViolationWrapper>
        );
      },
    },
    {
      title: "Voilation Type",
      dataIndex: "voilationtype",
      key: "voilationtype",
    },

    {
      title: "Description",
      dataIndex: "description",
      key: "description",
    },

    {
      title: "Date",
      dataIndex: "date",
      key: "date",
    },
    {
      title: "Status",
      dataIndex: "status",
      render: (_: any, record: ViolationData) => {
        return (
          <StatusWrapper status={record.status}>{record.status}</StatusWrapper>
        );
      },
      key: "status",
    },
  ];

  const dataSource: ViolationData[] = useMemo(() => {
    return Array.from({ length: 100 }, (_, index) => ({
      userId: index + 1,
      name: `Florence ${index + 1}`,
      subname: `Florence Zach ${index + 1}`,
      voilationtype: "Fraudulent Token Sale",
      description: "Created fake asset tokens and attempted to sell the...",
      date: "23 Oct, 2023",
      status:
        index % 3 === 0
          ? "Suspended"
          : index % 3 === 1
          ? "No Suspension"
          : "Suspension Resolved",
    }));
  }, []);

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);
  return (
    <Container>
      <Content>
        <div>
          <PageTitle>Users</PageTitle>
          <SearchBar
            customWidth="320px"
            // placeholder="Start search..."
            // value={searchInput}
            // onChange={handleSearchInputChange}
            // handleSearch={handleSearch}
          />
        </div>
        <UsersFilter setSelectedFilters={setSelectedFilters} />
      </Content>
      <SelectedFilters>
        {Object.entries(selectedFilters).length > 0 && (
          <StyledSearch>
            {Object.entries(selectedFilters).map(([key, value]) => (
              <SingleSearch key={key}>
                <KeyText>{key}:</KeyText>
                <ValueText>{value}</ValueText>
                <RemoveButton
                  onClick={(e: React.MouseEvent) => {
                    e.preventDefault();
                    handleRemoveFilter(key);
                  }}
                  aria-label={`Remove ${key} filter`}
                >
                  <FaX />
                </RemoveButton>
              </SingleSearch>
            ))}
          </StyledSearch>
        )}
      </SelectedFilters>
      <CustomTable columns={columns} dataSource={paginatedData} />
      <Pagination
        currentPage={currentPage}
        totalCount={dataSource.length}
        pageSize={pageSize}
        onPageSizeChange={setPageSize}
        onPageChange={setCurrentPage}
      />
    </Container>
  );
};

export default ViolationsLog;

const Container = styled.section`
  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 20%;
  }

  th:nth-child(2),
  td:nth-child(2) {
    width: 18%;
  }

  th:nth-child(3),
  td:nth-child(3) {
    width: 16%;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 12%;
    text-align: center;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 10%;
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

const ViolationIdSection = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

const ViolationWrapper = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;

const ViolationName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
  display: flex;
  align-items: center;
  gap: 4px;
`;

const ViolationSubName = styled.p`
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const StatusWrapper = styled.p<{ status: string }>`
  color: ${({ status }) =>
    status === "Suspended"
      ? "#BE3800"
      : status === "No Suspension"
      ? "#00225A"
      : "#00A859"};
  font-weight: 500;
  text-align: center;
  margin: 0 auto;
`;

const VerifiedIcon = styled.div`
  color: #007cdf;
  display: flex;
  align-items: center;
  svg {
    width: 12px;
    height: 12px;
  }
`;

const Content = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 10px;
`;

const PageTitle = styled.h2`
  color: #00225a;
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  text-align: left;
  margin-bottom: 10px;
`;

const SelectedFilters = styled.div`
  margin: 20px 0;
  padding: 10px;
  // background-color: #f8f9fa;
  // border-radius: 8px;

  h4 {
    margin: 0 0 8px 0;
    color: #00225a;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    color: #333;
  }
`;

const StyledSearch = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;
const SingleSearch = styled.div`
  border: 1px solid #007cdf;
  background-color: #acd1ef;
  font-size: 12px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: center;
  border-radius: 16px;
  padding: 8px 16px;
  display: inline-flex;
  align-items: center;
  gap: 6px;

  transition: all 0.2s ease-in-out;
`;

const KeyText = styled.span`
  color: #4f4f4f;
  font-weight: 500;
`;

const ValueText = styled.span`
  color: #007cdf;
`;
const RemoveButton = styled.button`
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  padding: 0;
  margin-left: 4px;
  cursor: pointer;
  color: #007cdf;
  width: 10px;
  height: 10px;
`;
