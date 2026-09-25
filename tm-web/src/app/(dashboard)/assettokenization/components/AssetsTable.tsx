"use client";
import CustomTable from "@/components/CustomTable";
import Link from "next/link";
import React, { useMemo, useState } from "react";
import styled from "styled-components";
import Pagination from "@/components/CustomPagination";

import nigFlag from "@/assets/images/twemoji_flag-nigeria.svg";
import Image from "next/image";
import {
  AiOutlineCheckCircle,
  AiOutlineClockCircle,
  AiOutlineCloseCircle,
} from "react-icons/ai";
import { useRouter } from "next/navigation";
import { useGetTokenizationListQuery } from "@/redux/api/assettokenization";
import useFormatDate from "@/hooks/useFormatDate";
import { useFetchUserByCriteriaQuery } from "@/redux/api/users";
import EmailCell from "./EmailCell";
import EmptyState from "@/components/EmptyState";

interface AssetsTableProps {
  filters: any;
}

const AssetsTable = ({ filters }: AssetsTableProps) => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const formattedDate = useFormatDate();
  const {
    data: tokenizationData,
    error,
    isLoading,
    isFetching,
  } = useGetTokenizationListQuery({
    page: currentPage,
    limit: pageSize,
    ...filters,
  });

  const router = useRouter();

  const dataSource = tokenizationData?.records || [];
  const totalRecords = tokenizationData?.totalRecords || 0;

  // Updated status mapping that aligns with your process steps
  const statusMapping: {
    [key: number]: {
      label: string;
      icon: JSX.Element;
      textColor: string;
      backgroundColor: string;
    };
  } = {
    0: {
      label: "pending vetting",
      icon: <AiOutlineCheckCircle />,
      textColor: "#007CDF",
      backgroundColor: "#F2F6F9",
    },
    1: {
      label: "Awaiting Payment",
      icon: <AiOutlineCheckCircle />,
      textColor: "#007CDF",
      backgroundColor: "#F2F6F9",
    },
    2: {
      label: "Payment Made",
      icon: <AiOutlineCheckCircle />,
      textColor: "#007CDF",
      backgroundColor: "#F2F6F9",
    },
    3: {
      label: "Processing",
      icon: <AiOutlineClockCircle />,
      textColor: "#007CDF",
      backgroundColor: "#F2F6F9",
    },
    4: {
      label: "Tokenization Approved",
      icon: <AiOutlineCheckCircle />,
      textColor: "#00A859",
      backgroundColor: "#00A8591A",
    },
    5: {
      label: "Primary Sale Started",
      icon: <AiOutlineCheckCircle />,
      textColor: "#00A859",
      backgroundColor: "#00A8591A",
    },
    6: {
      label: "Secondary Sales",
      icon: <AiOutlineCheckCircle />,
      textColor: "#00A859",
      backgroundColor: "#00A8591A",
    },
    7: {
      label: "Liquidated",
      icon: <AiOutlineCloseCircle />,
      textColor: "#FF4D4D",
      backgroundColor: "#BE38001A",
    },
    8: {
      label: "Refunded",
      icon: <AiOutlineCloseCircle />,
      textColor: "#FF4D4D",
      backgroundColor: "#BE38001A",
    },
  };

  const StatusWithIcon = ({
    vettingStatus,
    assetTokenizationStatus,
  }: {
    vettingStatus: number | string;
    assetTokenizationStatus: number | string;
  }) => {
    // Convert values in case they come as strings
    const vetting = Number(vettingStatus);
    const tokenization = Number(assetTokenizationStatus);

    // Hard coded check for "pending vetting"
    if (vetting === 0 && tokenization === 1) {
      return (
        <AssetStatus textColor="#007CDF" backgroundColor="#F2F6F9">
          <AiOutlineCheckCircle />
          pending vetting
        </AssetStatus>
      );
    }

    // Fallback to using the mapping for other statuses
    const mapping = statusMapping[tokenization] || {
      label: "Unknown",
      icon: <AiOutlineClockCircle />,
      textColor: "#828282",
      backgroundColor: "#F2F6F9",
    };

    return (
      <AssetStatus
        textColor={mapping.textColor}
        backgroundColor={mapping.backgroundColor}
      >
        {mapping.icon}
        {mapping.label}
      </AssetStatus>
    );
  };

  //   const mapping = statusMapping[status] || {
  //     label: "Unknown",
  //     icon: <AiOutlineClockCircle />,
  //     textColor: "#828282",
  //     backgroundColor: "#F2F6F9",
  //   };
  //   return (
  //     <AssetStatus
  //       textColor={mapping.textColor}
  //       backgroundColor={mapping.backgroundColor}
  //     >
  //       {mapping.icon}
  //       {mapping.label}
  //     </AssetStatus>
  //   );
  // };

  const columns = useMemo(
    () => [
      {
        title: "Asset ID",
        dataIndex: "id",

        render: (_: any, record: any) => {
          return (
            <AssetIdSection>
              {record?.assetLogo ? (
                <Image
                  src={record?.assetLogo}
                  alt="Asset Logo"
                  width={40}
                  height={40}
                  style={{ borderRadius: "50%" }}
                  unoptimized={true} // Skip Next.js image optimization
                  onError={(e) => {
                    e.currentTarget.onerror = null; // Prevents infinite loop
                    e.currentTarget.src =
                      "https://media.istockphoto.com/id/1300845620/vector/user-icon-flat-isolated-on-white-background-user-symbol-vector-illustration.jpg?s=612x612&w=0&k=20&c=yBeyba0hUkh14_jgv1OKqIH0CCSWU_4ckRkAoy2p73o="; // Path to your fallback image
                  }}
                />
              ) : (
                <Avatar />
              )}

              <div>
                <AssetName>{record.assetCode}</AssetName>
                <AssetSubName>{record.assetName}</AssetSubName>
              </div>
            </AssetIdSection>
          );
        },
      },

      {
        title: "Category",
        dataIndex: "assetSector",
        render: (_: any, record: any) => <span>{record.assetSector}</span>,

        key: "assetSector",
      },

      {
        title: " Location",
        dataIndex: "assetCountryLocation",
        render: (_: any, record: any) => {
          return (
            <LocationWrapper>
              <Image src={nigFlag} alt="Nigeria-flag" />
              {record.assetCountryLocation}
            </LocationWrapper>
          );
        },
        key: "assetCountryLocation",
      },
      {
        title: "Email Address",
        dataIndex: "email",
        render: (_: any, record: any) => (
          <EmailCell initiatorUsername={record?.initiatorUsername} />
        ),
        key: "email",
      },

      {
        title: "Created On",
        dataIndex: "createdAt",
        render: (CreatedAt: any) => formattedDate(CreatedAt),

        key: "createdAt",
      },
      {
        title: " Status",
        dataIndex: "assetTokenizationStatus",
        render: (_: any, record: any) => (
          <StatusWithIcon
            vettingStatus={record.vettingStatus}
            assetTokenizationStatus={record.assetTokenizationStatus}
          />
        ),
        key: "assetTokenizationStatus",
      },
    ],
    [formattedDate],
  );

  if (error) return <div>Error loading data.</div>;

  return (
    <Container>
      {dataSource.length > 0 || isLoading || isFetching ? (
        <>
          <CustomTable
            columns={columns}
            dataSource={dataSource}
            isLoading={isLoading || isFetching}
            onRowClick={(record) => {
              router.push(`/assettokenization/${record.id}`);
            }}
          />
          {!isLoading && !isFetching && dataSource.length > 0 && (
            <Pagination
              currentPage={currentPage}
              totalCount={totalRecords}
              pageSize={pageSize}
              onPageChange={setCurrentPage}
              onPageSizeChange={setPageSize}
            />
          )}
        </>
      ) : (
        <EmptyState
          title="Asset Not Found"
          message="We couldn't find any assets matching your search or filters. Please try again with different criteria."
        />
      )}
    </Container>
  );
};

export default AssetsTable;
const Container = styled.section`
  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 12%;
  }

  th:nth-child(2),
  td:nth-child(2) {
    width: 15%;
  }

  th:nth-child(3),
  td:nth-child(3) {
    width: 12%;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 14%;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 15%;
    text-align: center;
  }

  th:nth-child(6),
  td:nth-child(6) {
    width: 16%;
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

const StyledLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #00225a;
`;
const AssetIdSection = styled.div`
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
const EmailAddress = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;
const LocationWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
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
