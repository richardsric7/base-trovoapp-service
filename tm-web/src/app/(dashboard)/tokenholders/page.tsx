"use client";

import { SearchBar } from "@/components";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import { useRouter } from "next/navigation";
import React, { useMemo, useState } from "react";
import { FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";
import TokenHoldersModal from "../assettokenization/components/TokenHoldersModal";

const TokenHoldersPage = () => {
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const handleBack = () => {
    router.back();
  };

  const cardDetails = [
    {
      title: "Token for sale",
      value: 4000,
    },

    {
      title: "Total Buyers",
      value: 9000,
    },
  ];

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
      title: " Quantity Held",
      dataIndex: "quantityHeld",
      key: "quantityHeld",
    },
    {
      title: "Last Activity",
      dataIndex: "date",
      key: "date",
    },

    {
      title: "Action",
      dataIndex: "",
      render: (_: any, record: any) => {
        return (
          <Actionwrapper onClick={() => setIsOpen(!isOpen)}>
            View Token Activity
          </Actionwrapper>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 100 }, (_, index) => ({
      userId: index + 1,
      name: `Florence ${index + 1}`,
      subname: `florence Zach`,
      quantityHeld: "1000 CNGN",
      date: "23 Sep 2023, 8:17",
    }));
  }, []);

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);

  return (
    <>
      <BackButtonLink onClick={handleBack}>
        <FaArrowLeft size={18} />
      </BackButtonLink>
      <Container>
        <Title>Token Holders</Title>
        <CardContainer>
          {cardDetails.map((item) => (
            <InfoCard key={item.title}>
              <CardLabel>{item.title}</CardLabel>
              <CardValue>{item.value}</CardValue>
            </InfoCard>
          ))}
        </CardContainer>

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

      {isOpen && <TokenHoldersModal isOpen={isOpen} setIsOpen={setIsOpen} />}
    </>
  );
};

export default TokenHoldersPage;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;

  table {
    width: 100%;
    table-layout: fixed;
    padding-top: 8px;
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
    width: 18%;
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
  margin-bottom: 10px;
  color: #00225a;
`;

const BackButtonLink = styled.button`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
  background: none;
  cursor: pointer;
  border: none;
  padding: 24px;
`;

const CardContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  width: 100%;
  margin-bottom: 14px;
`;
const InfoCard = styled.div`
  background-color: #f2f6f9;
  border-radius: 10px;
  padding: 8px 16px;
  width: 80%;
`;

const CardLabel = styled.p`
  font-weight: 400;
  font-size: 14px;
  line-height: 20px;
  color: #828282;
  padding-bottom: 10px;
`;

const CardValue = styled.p`
  font-weight: 600;
  font-size: 24px;
  line-height: 24px;
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

const Actionwrapper = styled.button`
  background-color: transparent;
  border-radius: 8px;
  border: 1px solid #007cdf;
  padding: 10px 24px;
  color: #007cdf;
  font-family: inherit;
  cursor: pointer;
`;
