"use client";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import React, { useMemo, useState } from "react";
import styled from "styled-components";
import TruncatedText from "@/hooks/useTruncate";
import { useRouter } from "next/navigation";
import { useP2pUsersQuery } from "@/redux/api/p2p";
import { FaEllipsisVertical } from "react-icons/fa6";
import tradeIcon from "@/assets/images/coin.svg";
import statusIcon from "@/assets/images/status.svg";
import Image from "next/image";
import RecordViolation from "./RecordViolation";
import SuccessMessage from "@/components/SuccessMessage";
interface SearchTermProps {
  searchTerm: string;
}

const P2pUsersTable = ({ searchTerm }: SearchTermProps) => {
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [activeActionId, setActiveActionId] = useState<string | null>(null);

  const [isOpen, setIsOpen] = useState(false);
  const [suspendSuccess, setSuspendSuccess] = useState(false);
  // To pass user info for the modal:
  const [selectedUser, setSelectedUser] = useState<any>(null);

  const { data, isLoading, isFetching } = useP2pUsersQuery({
    page: currentPage,
    pageSize,
    search: searchTerm,
  });

  const totalUsers = data?.data?.total || 0;

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const columns = [
    {
      title: "User ID",
      dataIndex: "username",
      render: (_: any, record: any) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserNameWrapper>
              <UserName>
                <TruncatedText text={record.username} maxLength={15} />
                {/* <VerifiedIcon>
                  <MdVerified />
                </VerifiedIcon> */}
              </UserName>
            </UserNameWrapper>
          </div>
        </UserInfoSection>
      ),
    },
    {
      title: "Email Address",
      dataIndex: "email",
      key: "email",
      render: (email: string) => <TruncatedText text={email} maxLength={25} />,
    },
    {
      title: "Phone Number",
      dataIndex: "phone",
      key: "phone",
    },
    {
      title: "Registration Date",
      dataIndex: "registration_date",
      key: "registration_date",
      render: (date: number) =>
        new Date(date).toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        }),
    },
    {
      title: "Account Status",
      dataIndex: "account_status",
      key: "account_status",
      render: (status: number) => {
        const isActive = status === 0;
        return (
          <UserStatus isActive={isActive}>
            {isActive ? "Active" : "Suspended"}
          </UserStatus>
        );
      },
    },
    {
      title: "Market Makers",
      dataIndex: "merchant_status",
      key: "merchant_status",
      render: (status: number) => {
        const isMerchant = status === 1;

        return (
          <MerchantStatus isNormal={isMerchant}>
            {isMerchant ? "Merchant" : "Non-merchant"}
          </MerchantStatus>
        );
      },
    },

    {
      title: "Violations",
      dataIndex: "violations",
      render: (_: any, record: any) => {
        return (
          <ActionsWrapper
            onClick={(e) => {
              e.stopPropagation();
              setActiveActionId(
                activeActionId === record.id ? null : record.id
              );
            }}
          >
            <FaEllipsisVertical />
            {activeActionId === record.id && (
              <TooltipWrapper onClick={(e) => e.stopPropagation()}>
                <ToolTipButton>
                  <Image src={tradeIcon} alt="status-icon" />
                  <TooltipItem>Set Trading Fee</TooltipItem>
                </ToolTipButton>
                <ToolTipButton
                  onClick={() => {
                    setIsOpen(true);
                    setActiveActionId(null);
                  }}
                >
                  <Image src={statusIcon} alt="status-icon" />
                  <TooltipItem>Record Violation</TooltipItem>
                </ToolTipButton>
              </TooltipWrapper>
            )}
          </ActionsWrapper>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    return (
      data?.data?.data?.map((user) => ({
        id: user.id,
        username: user.username,
        email: user.email,
        phone: user.phone,
        registration_date: user.registration_date,
        account_status: user.account_status,
        merchant_status: user.merchant_status,
      })) || []
    );
  }, [data?.data?.data]);

  if (isLoading) {
    return <UserName>Loading users...</UserName>;
  }

  return (
    <Container>
      {searchTerm && isFetching ? (
        <UserName>Searching users...</UserName>
      ) : searchTerm && totalUsers === 0 ? (
        <UserName>No users match your search.</UserName>
      ) : (
        <>
          <CustomTable
            columns={columns}
            dataSource={dataSource}
            onRowClick={(record) => {
              const destination =
                record.account_status !== 0
                  ? `/p2pusers/suspendUser/${record.username}`
                  : `/p2pusers/${record.username}`;
              router.push(destination);
            }}
          />
          <Pagination
            totalCount={totalUsers}
            onPageChange={handlePageChange}
            currentPage={currentPage}
            pageSize={pageSize}
            onPageSizeChange={setPageSize}
            isFetching={isFetching}
          />
        </>
      )}

      {isOpen && (
        <RecordViolation
          isOpen={isOpen}
          setIsOpen={setIsOpen}
          setSuspendSuccess={setSuspendSuccess}
        />
      )}

      {suspendSuccess && (
        <SuccessMessage
          isOpen={suspendSuccess}
          setIsOpen={setSuspendSuccess}
          heading="Success"
          message="User has been successfully recorded a violation"
          // email={user?.email}
        />
      )}
    </Container>
  );
};

export default P2pUsersTable;

const Container = styled.section`
  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 14%;
  }

  th:nth-child(2),
  td:nth-child(2) {
    width: 10%;
  }

  th:nth-child(3),
  td:nth-child(3) {
    width: 10%;
    text-align: center;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 6%;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 8%;
    text-align: center;
  }

  th:nth-child(6),
  td:nth-child(6) {
    width: 8%;
    text-align: center;
  }

  th:nth-child(7),
  td:nth-child(7) {
    width: 4%;
    text-align: center;
  }

  th {
    color: #828282;
    font-size: 14px;
    font-weight: 500;
    // text-align: center;
  }
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
  background-color: rebeccapurple;
`;

const UserNameWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
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

const UserName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  margin: 0;
  text-transform: capitalize;
  display: flex;
  align-items: center;
  gap: 4px;
`;

const UserStatus = styled.p<{ isActive: boolean }>`
  color: ${(props) => (props.isActive ? "#00A859" : "#BE3800")};
  background-color: ${(props) => (props.isActive ? "#00A8591A" : "#BE38001A")};
  font-weight: 500;
  padding: 8px;
  width: 100px;
  border-radius: 8px;
  text-align: center;
  margin: 0 auto;
`;

const MerchantStatus = styled.p<{ isNormal: boolean }>`
  color: ${(props) => (props.isNormal ? "#007CDF" : "#4F4F4F")};
  background-color: ${(props) => (props.isNormal ? "#007CDF1A" : "#EFEFEF")};
  font-weight: 500;
  padding: 8px;
  width: 120px;
  border-radius: 8px;
  text-align: center;
  margin: 0 auto;
`;

const ActionsWrapper = styled.div`
  position: relative;
  display: inline-block;
  cursor: pointer;
  padding: 12px;
  margin: -8px;
`;

const TooltipWrapper = styled.div`
  position: absolute;
  top: 24px;
  right: 0;
  background: #fff;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  padding: 8px 0;
  z-index: 10;
  display: flex;
  flex-direction: column;
  width: 160px;
`;

const ToolTipButton = styled.div`
  padding: 8px 6px;
  display: flex;
  align-items: center;
  gap: 1px;
  cursor: pointer;
  &:hover {
    background-color: #f2f6f9;
  }
`;

const TooltipItem = styled.div`
  font-size: 14px;
  color: #00225a;
  cursor: pointer;
`;
