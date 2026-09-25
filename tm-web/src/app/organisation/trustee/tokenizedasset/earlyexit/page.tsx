"use client";

import { SearchBar } from "@/components";
import React, { useMemo, useState } from "react";
import styled from "styled-components";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import { MdVerified } from "react-icons/md";
import TruncatedText from "@/hooks/useTruncate";
import { AiOutlineCheckCircle, AiOutlineClockCircle } from "react-icons/ai";
import ProcessExitModal from "./components/ProcessExitModal";
import ViewModal from "./components/ViewModal";
import { FaArrowLeft } from "react-icons/fa6";
import { useRouter } from "next/navigation";
import AssetCurationFilter from "@/app/(dashboard)/assetcuration/components/AssetCurationFilter";

const EarlyExitPage = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const router = useRouter();

  const [modal, setModal] = useState<{
    type: "view" | "process" | null;
    record: any | null;
  }>({ type: null, record: null });

  const openView = (record: any) => setModal({ type: "view", record });
  const openProcess = (record: any) => setModal({ type: "process", record });
  const closeModal = () => setModal({ type: null, record: null });

  const handleBack = () => {
    router.back();
  };

  const columns = [
    {
      title: "User",
      dataIndex: "user",
      render: (_: any, record: any) => {
        return (
          <InvestorInfoSection>
            <Avatar />
            <div>
              <UserNameWrapper>
                <UserName>
                  {" "}
                  <TruncatedText text={record.userName} maxLength={15} />
                  {record.kyc_level === 0 ? (
                    ""
                  ) : (
                    <MdVerified size={12} color="#007CDF" />
                  )}
                </UserName>
              </UserNameWrapper>

              <FullName>{record.fullName}</FullName>
            </div>
          </InvestorInfoSection>
        );
      },
    },

    {
      title: "Requested Time",
      dataIndex: "requestedTime",
      key: "requestedTime",
    },
    {
      title: "Token Quantity",
      dataIndex: "tokenQuantity",
      key: "tokenQuantity",
    },
    {
      title: " Estimated Payout",
      dataIndex: "estimatedPayout",
      key: "estimatedPayout",
    },

    {
      title: " Status",
      dataIndex: "status",
      key: "status",
      render: (_: any, record: any) => (
        <StatusWithIcon status={record.status} />
      ),
    },

    {
      title: "Action",
      dataIndex: "action",
      key: "action",
      render: (_: any, record: any) => {
        const isProcessed = record.status === "Processed";
        const label = isProcessed ? "View" : "Process";
        const onClick = isProcessed
          ? () => openView(record)
          : () => openProcess(record);

        return (
          <ActionButton onClick={onClick} aria-label={`${label} request`}>
            {label}
          </ActionButton>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 100 }, (_, index) => ({
      id: index + 1,
      userName: "Odogwu",
      fullName: "Obi Enechi",
      requestedTime: `23 Sep 2023, 8:17`,
      tokenQuantity: `10,000 TROV`,
      estimatedPayout: `10,000 TROV`,
      status: index % 2 ? "Processed" : "Pending",
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
        <Title>Early Exit Requests</Title>

        <SubHeader>
          <SearchBar />
          <AssetCurationFilter />
        </SubHeader>
        <CustomTable columns={columns} dataSource={paginatedData} />
        <Pagination
          currentPage={currentPage}
          totalCount={dataSource.length}
          pageSize={pageSize}
          onPageChange={setCurrentPage}
          onPageSizeChange={setPageSize}
        />
        <ProcessExitModal
          isOpen={modal.type === "process"}
          onClose={closeModal}
        />
        <ViewModal isOpen={modal.type === "view"} onClose={closeModal} />
      </Container>
    </>
  );
};

export default EarlyExitPage;

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

const Container = styled.section`
  background: #ffffff;
  padding: 32px;

  border-radius: 24px;

  table {
    width: 100%;
    table-layout: fixed;
    padding-top: 8px;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 14%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 12%;
    text-align: center;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 10%;
    text-align: center;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 8%;
    text-align: center;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 10%;
    text-align: center;
  }
  th:nth-child(6),
  td:nth-child(6) {
    width: 8%;
    text-align: center;
  }

  th:nth-child(7),
  td:nth-child(7) {
    width: 10%;
    text-align: center;
  }
  th:nth-child(8),
  td:nth-child(8) {
    width: 10%;
    text-align: center;
  }
  th:nth-child(9),
  td:nth-child(9) {
    width: 5%;
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
  margin: 0;
  color: #00225a;
`;

const InvestorInfoSection = styled.div`
  display: flex;
  column-gap: 10px;
  align-items: center;
  cursor: pointer;
`;
const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;
const UserNameWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;
const UserName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const FullName = styled.p`
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;
const SubHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 18px;
`;

const ActionButton = styled.button`
  background-color: transparent;
  border-radius: 8px;
  border: 1px solid #007cdf;
  padding: 10px 24px;
  color: #007cdf;
  font-family: inherit;
  cursor: pointer;

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`;

const AssetStatus = styled.p<{ textColor: string; backgroundColor: string }>`
  color: ${({ textColor }) => textColor};
  background-color: ${({ backgroundColor }) => backgroundColor};
  font-weight: 500;
  padding: 8px;
  border-radius: 8px;
  text-align: center;
  margin: 0 auto;
  display: inline-flex;
  align-items: center;
  gap: 4px;
`;

const StatusWithIcon = ({ status }: { status: "Processed" | "Pending" }) => {
  const isProcessed = status === "Processed";

  return (
    <AssetStatus
      textColor={isProcessed ? "#00A859" : "#007CDF"}
      backgroundColor={isProcessed ? "#00A8591A" : "#F2F6F9"}
    >
      {isProcessed ? (
        <AiOutlineCheckCircle size={14} />
      ) : (
        <AiOutlineClockCircle size={14} />
      )}
      {status}
    </AssetStatus>
  );
};
