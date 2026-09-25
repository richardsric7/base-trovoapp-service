"use client";

import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import React, { useMemo, useState } from "react";
import nigFlag from "@/assets/images/twemoji_flag-nigeria.svg";
import Image from "next/image";
import styled from "styled-components";
import { FaCheck } from "react-icons/fa";

interface AssetData {
  userId: number;
  name: string;
  subname: string;
  category: string;
  location: string;
  email: string;
  createdOn: string;
}

interface ColumnType {
  title: string;
  dataIndex: string;
  key?: string;
  width?: string;
  render?: (text: string, record: AssetData) => React.ReactNode;
}

const AddCurationTable: React.FC = () => {
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [selectedRows, setSelectedRows] = useState<number[]>([]);

  // Toggle check state
  const toggleRowSelection = (id: number) => {
    setSelectedRows((prevSelected) =>
      prevSelected.includes(id)
        ? prevSelected.filter((rowId) => rowId !== id)
        : [...prevSelected, id]
    );
  };

  const columns: ColumnType[] = [
    {
      title: "Asset ID",
      dataIndex: "assetsId",
      width: "20%",
      render: (_, record: AssetData) => {
        return (
          <AssetWrapper>
            <CheckBox
              isChecked={selectedRows.includes(record.userId)}
              onClick={() => toggleRowSelection(record.userId)}
            >
              {selectedRows.includes(record.userId) && (
                <FaCheck color="#004988" />
              )}
            </CheckBox>
            <AssetIdSection>
              <Avatar />
              <div>
                <AssetName>{record.name}</AssetName>
                <AssetSubName>{record.subname}</AssetSubName>
              </div>
            </AssetIdSection>
          </AssetWrapper>
        );
      },
    },
    {
      title: "Category",
      dataIndex: "category",
      key: "category",
    },
    {
      title: "Location",
      dataIndex: "location",
      render: (_, record: AssetData) => {
        return (
          <LocationWrapper>
            <Image src={nigFlag} alt="Nigeria-flag" />
            {record.location}
          </LocationWrapper>
        );
      },
      key: "location",
    },
    {
      title: "Email Address",
      dataIndex: "email",
      key: "email",
    },
    {
      title: "Created On",
      dataIndex: "createdOn",
      key: "createdOn",
    },
  ];

  const dataSource: AssetData[] = useMemo(() => {
    return Array.from({ length: 100 }, (_, index) => ({
      userId: index + 1,
      name: `Atlantis ${index + 1}`,
      subname: `Atlantis Developers ${index + 1}`,
      category: "Real Estate",
      location: "Nigeria",
      email: `kennismaduka${index + 1}@gmail.com`,
      createdOn: "23 Oct, 2023",
    }));
  }, []);

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);

  return (
    <Container>
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

export default AddCurationTable;

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
    width: 14%;
  }

  th:nth-child(3),
  td:nth-child(3) {
    width: 14%;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 18%;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 10%;
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

const AssetIdSection = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

const AssetWrapper = styled.div`
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

const AssetName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const AssetSubName = styled.p`
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const LocationWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

interface CheckBoxProps {
  isChecked: boolean;
}

const CheckBox = styled.div<CheckBoxProps>`
  width: 16px;
  height: 16px;
  border: 1.5px solid ${(props) => (props.isChecked ? "#004988" : "#bdbdbd")};

  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  background-color: ${(props) =>
    props.isChecked ? "transparent" : "transparent"};
  color: white;
  transition: all 0.2s ease-in-out;
`;
