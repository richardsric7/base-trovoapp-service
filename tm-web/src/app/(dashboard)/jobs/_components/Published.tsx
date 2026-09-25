"use client";

import { SearchBar } from "@/components";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";

import { useMemo, useState } from "react";
import styled from "styled-components";
import JobFilter from "./JobFilter";
import { useRouter } from "next/navigation";
import { useGetJobsQuery } from "@/redux/api/jobs";

const Published = () => {
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchTerm, setSearchTerm] = useState("");
  const [searchInput, setSearchInput] = useState<string>("");
  const [filter, setFilter] = useState({
    type: "",
    location: "",
    experience: "",
  });
  const { data, isLoading, isFetching, error, refetch } = useGetJobsQuery({
    status: "published",
    search: searchTerm,
    location: filter.location || undefined,
    work_type: filter.type || undefined,
    years_of_experience: filter.experience || undefined,
    page: currentPage,
    pageSize,
  });

  const handleSearchInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearchInput(value);

    if (value.trim() === "") {
      setSearchTerm("");
      setCurrentPage(1);
    }
  };

  const handleSearch = () => {
    setSearchTerm(searchInput);
    setCurrentPage(1);
    refetch();
  };
  const handleFilterChange = (value: any) => {
    setFilter(value);
    setCurrentPage(1);
  };

  const columns = [
    {
      title: "Role",
      dataIndex: "role",
      key: "role",
    },

    {
      title: "Type",
      dataIndex: "type",
      key: "type",
    },

    {
      title: "Location",
      dataIndex: "location",
      key: "location",
    },
    {
      title: " Created On",
      dataIndex: "date",
      key: "date",
    },
  ];

  const dataSource = useMemo(() => {
    if (!data?.data) return [];

    return data?.data?.map((job: any) => ({
      id: job.id,
      // username: "Admin",
      // fullname: "System Admin",
      role: job.heading,
      type: job.work_type?.charAt(0).toUpperCase() + job.work_type?.slice(1),
      location: job.location?.charAt(0).toUpperCase() + job.location?.slice(1),
      date: new Date(job.created_at).toLocaleDateString("en-GB", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      }),
    }));
  }, [data]);

  const paginatedDataSource = useMemo(() => {
    const start = (currentPage - 1) * pageSize;
    const end = start + pageSize;

    return dataSource.slice(start, end);
  }, [dataSource, currentPage, pageSize]);

  const hasActiveQuery =
    Boolean(searchTerm) ||
    Boolean(filter.location) ||
    Boolean(filter.type) ||
    Boolean(filter.experience);

  return (
    <div>
      <Container>
        <SubHeader>
          <SearchBar
            value={searchInput}
            onChange={handleSearchInputChange}
            handleSearch={handleSearch}
          />
          <JobFilter onFilter={handleFilterChange} />
        </SubHeader>

        {isLoading ? (
          <LoadingText>Loading jobs...</LoadingText>
        ) : isFetching ? (
          <LoadingText>Searching jobs...</LoadingText>
        ) : hasActiveQuery && dataSource.length === 0 ? (
          <LoadingText>No jobs match your search or filters.</LoadingText>
        ) : (
          <TableContainer>
            <CustomTable
              columns={columns}
              dataSource={paginatedDataSource}
              isLoading={isLoading}
              onRowClick={(record) => router.push(`/jobs/${record.id}`)}
            />
            <Pagination
              currentPage={currentPage}
              totalCount={data?.data?.length || 0}
              pageSize={pageSize}
              onPageChange={setCurrentPage}
              onPageSizeChange={setPageSize}
            />
          </TableContainer>
        )}
      </Container>
    </div>
  );
};

export default Published;

const Container = styled.section``;

const SubHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 24px;
`;

const TableContainer = styled.div`
  margin-top: 12px;

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
    width: 10%;
    text-align: center;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 16%;
    text-align: center;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 14%;
    text-align: center;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 12%;
    text-align: center;
  }

  // th:nth-child(6),
  // td:nth-child(6) {
  //   width: 10%;
  //   text-align: center;
  // }

  // th:nth-child(7),
  // td:nth-child(7) {
  //   width: 4%;
  //   text-align: center;
  // }
`;

const UserNameWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
  max-width: 40px;
`;

const LoadingText = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: center;
  margin-top: 20px;
`;

const UserInfoSection = styled.div`
  display: flex;
  column-gap: 10px;
  align-items: center;
  cursor: pointer;
`;

const FullName = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  margin: 0;
  text-transform: capitalize;
`;
