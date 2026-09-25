"use client";
import { SearchBar } from "@/components";
import { useMemo, useState } from "react";
import styled from "styled-components";
import AssetRecordsFilter from "./AccessRecordsFilter";
import CustomTable from "@/components/CustomTable";
import { useRouter } from "next/navigation";
import { MdVerified } from "react-icons/md";
import { AiOutlineCheckCircle, AiOutlineCloseCircle } from "react-icons/ai";
import Pagination from "@/components/CustomPagination";
interface UserStatusProps {
  isActive?: boolean;
  isFailed?: boolean;
}

const AssetRecordsTable = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [activeId, setActiveId] = useState<number | null>(null);
  const router = useRouter();

  const dataSource = useMemo(() => {
    return Array.from({ length: 100 }, (_, index) => ({
      id: index + 1,
      username: "Marlon",
      fullname: "Mark Ovey",
      email: "kennismaduka@gmail.com",
      method: "Password",
      address: "192.168.1.1",
      date: `30 Jun, 2025 08:30 AM `,
      status: index % 3 === 0 ? "successful" : "failed",
    }));
  }, []);

  const selectedRecord = dataSource.find((item) => item.id === activeId);

  const columns = [
    {
      title: "Users",
      dataIndex: "username",
      render: (_: any, record: any) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserNameWrapper>
              <UserName>{record.username}</UserName>
              <MdVerified color="#007cdf" />
            </UserNameWrapper>
            <FullName>{record.fullname}</FullName>
          </div>
        </UserInfoSection>
      ),
    },

    {
      title: "Email Address",
      dataIndex: "email",
      key: "email",
    },

    {
      title: "Ip Address",
      dataIndex: "address",
      key: "adress",
    },

    {
      title: "Authentication Method",
      dataIndex: "method",
      key: "method",
    },

    {
      title: " Login Date/Time",
      dataIndex: "date",
      key: "date",
    },

    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (status: string) => {
        const isSuccessful = status === "successful";
        const isFailed = status === "failed";

        return (
          <PaymentStatus isActive={isSuccessful} isFailed={isFailed}>
            {isSuccessful && <AiOutlineCheckCircle size={18} />}

            {isFailed && <AiOutlineCloseCircle size={18} />}
            {status}
          </PaymentStatus>
        );
      },
    },
  ];

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);
  return (
    <Container>
      <Header>
        <TitleSection>
          <Title>Access History</Title>
          <SubTitle>Total Login Attempts: 20,134</SubTitle>
        </TitleSection>

        <RightSection>
          <TitleSection>
            <Title>10,000</Title>
            <SubTitle>Logins with Passwords</SubTitle>
          </TitleSection>

          <TitleSection>
            <Title>10,134</Title>
            <SubTitle>Logins with Biometrics</SubTitle>
          </TitleSection>
        </RightSection>
      </Header>
      <SubHeader>
        <SearchBar />
        <AssetRecordsFilter />
      </SubHeader>
      <TableContainer>
        <CustomTable columns={columns} dataSource={paginatedData} />
        <Pagination
          currentPage={currentPage}
          totalCount={dataSource.length}
          pageSize={pageSize}
          onPageChange={setCurrentPage}
          onPageSizeChange={setPageSize}
        />
      </TableContainer>
    </Container>
  );
};

export default AssetRecordsTable;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;

  border-radius: 24px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 2px;
  color: #00225a;
`;

const SubTitle = styled.p`
  font-weight: 400;
  font-size: 16px;
  color: #828282;
`;

const TitleSection = styled.div``;

const RightSection = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
`;

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
    width: 14%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 12%;
    text-align: center;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 8%;
    text-align: center;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 10%;
    text-align: center;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 12%;
    text-align: center;
  }

  th:nth-child(6),
  td:nth-child(6) {
    width: 10%;
    text-align: center;
  }

  th:nth-child(7),
  td:nth-child(7) {
    width: 4%;
    text-align: center;
  }
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

const UserName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
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

const PaymentStatus = styled.p<UserStatusProps>`
  color: ${({ isActive, isFailed }) =>
    isActive ? "#00A859" : isFailed ? "#BE3800" : "#6B7280"};
  background-color: ${({ isActive, isFailed }) =>
    isActive ? "#00A8591A" : isFailed ? "#BE38001A" : "#E5E7EB"};
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border-radius: 8px;
  text-align: center;
  margin: 0 auto;
  text-transform: capitalize;
`;
