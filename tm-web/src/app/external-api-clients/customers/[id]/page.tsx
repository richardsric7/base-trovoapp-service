"use client";

import styled from "styled-components";
import { useMemo, useState } from "react";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import { SearchBar } from "@/components";
import { useParams, useRouter } from "next/navigation";
import { FaArrowLeft, FaRegEnvelope } from "react-icons/fa6";
import { AiOutlineCheckCircle } from "react-icons/ai";
import { BsTelephone } from "react-icons/bs";
import { LuCalendarDays } from "react-icons/lu";
import { GoArrowDown, GoArrowSwitch, GoArrowUp } from "react-icons/go";

const CustomerDetailsPage = () => {
  const { id } = useParams();
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const columns = [
    {
      title: "Transaction",
      dataIndex: "type",
      render: (_: any, record: any) => {
        return (
          <TransactionType type={record.type}>
            <IconCircle type={record.type}>
              {record.type === "Received" && <GoArrowDown />}
              {record.type === "Sent" && <GoArrowUp />}
              {record.type === "Swap" && <GoArrowSwitch />}
              {record.type === "Purchase" && "🛒"}
              {record.type === "Fee" && "⚙️"}
            </IconCircle>
            <div>
              <TypeText>{record.type}</TypeText>
              <SubText>{record.date}</SubText>
            </div>
          </TransactionType>
        );
      },
    },
    {
      title: "Asset",
      dataIndex: "asset",
    },
    {
      title: "From",
      dataIndex: "from",
      render: (_: any, record: any) => (
        <UserRow>
          <Avatar />
          <span>{record.from}</span>
        </UserRow>
      ),
    },
    {
      title: "To",
      dataIndex: "to",
      render: (_: any, record: any) => (
        <UserRow>
          <Avatar />
          <span>{record.to}</span>
        </UserRow>
      ),
    },
    {
      title: "Amount",
      dataIndex: "amount",
      render: (_: any, record: any) => (
        <AmountText type={record.type}>{record.amount}</AmountText>
      ),
    },
    {
      title: "",
      dataIndex: "action", // fake field
      render: () => <ViewLink>View</ViewLink>,
    },
  ];

  const dataSource = useMemo(() => {
    return [
      {
        type: "Received",
        asset: "TROV",
        from: "Marlon",
        to: "Florencez",
        amount: "2200 TROV",
        date: "4 Oct, 2023",
      },
      {
        type: "Sent",
        asset: "TROV",
        from: "Florencez",
        to: "Marlon",
        amount: "2200 TROV",
        date: "4 Oct, 2023",
      },
      {
        type: "Swap",
        asset: "XBN/TROV",
        from: "XBN",
        to: "TROV",
        amount: "200 TROV",
        date: "4 Oct, 2023",
      },
      {
        type: "Purchase",
        asset: "AET",
        from: "--",
        to: "Florencez",
        amount: "2200 TROV",
        date: "4 Oct, 2023",
      },
    ];
  }, []);

  return (
    <>
      <FaArrowLeft
        color="#00225A"
        onClick={() => router.back()}
        style={{ cursor: "pointer", marginBottom: "20px" }}
      />
      <Container>
        {/* PROFILE CARD */}
        <ProfileCard>
          <ProfileLeft>
            <AvatarLarge />
            <div>
              <Name>Florence</Name>
              <SubText>Florence Zach</SubText>
              <Verified>
                <AiOutlineCheckCircle /> Verified
              </Verified>
              <InfoRow>
                <InfoItem>
                  <LuCalendarDays color="#828282" />
                  Added 20 Jul, 2023
                </InfoItem>
                <InfoItem>
                  <FaRegEnvelope color="#828282" />
                  florencez@gmail.com
                </InfoItem>
                <InfoItem>
                  <BsTelephone color="#828282" /> +234 80123456789
                </InfoItem>
              </InfoRow>
            </div>
          </ProfileLeft>
        </ProfileCard>

        {/* TRANSACTION */}
        <Card>
          <HeaderRow>
            <Title>Transaction History</Title>
            <Actions>
              <SearchBar />
              <FilterButton>Filters</FilterButton>
            </Actions>
          </HeaderRow>

          <CustomTable columns={columns} dataSource={dataSource} />

          <Pagination
            totalCount={100}
            currentPage={currentPage}
            pageSize={pageSize}
            onPageChange={setCurrentPage}
            onPageSizeChange={setPageSize}
          />
        </Card>
      </Container>
    </>
  );
};

export default CustomerDetailsPage;

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
`;

const Card = styled.div`
  background: #fff;
  border-radius: 16px;
  padding: 20px;
`;

const ProfileCard = styled(Card)`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const ProfileLeft = styled.div`
  display: flex;
  gap: 16px;
`;

const Avatar = styled.div`
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #ccc;
`;

const AvatarLarge = styled(Avatar)`
  width: 48px;
  height: 48px;
`;

const Name = styled.h2`
  margin: 0;
  font-size: 20px;
  color: #00225a;
`;

const SubText = styled.p`
  margin: 0;
  font-size: 12px;
  color: #6b7280;
`;

const Verified = styled.span`
  font-size: 12px;
  color: #00a859;
  background: #00a8591a;
  display: inline-flex;
  border-radius: 8px;
  align-items: center;
  padding: 4px 8px;
  gap: 4px;
  margin-top: 4px;
`;

const InfoRow = styled.div`
  display: flex;
  gap: 16px;
  font-size: 12px;
  margin-top: 6px;
`;

const InfoItem = styled.p`
  display: flex;
  align-items: center;
  gap: 4px;
  color: #00225a;
`;

const HeaderRow = styled.div`
  display: flex;
  justify-content: space-between;
  flex-direction: column;
  margin-bottom: 16px;

  gap: 12px;
`;

const Title = styled.h3`
  font-size: 18px;
  color: #00225a;
`;

const Actions = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
`;

const FilterButton = styled.button`
  border: 1px solid #ddd;
  padding: 8px 12px;
  border-radius: 8px;
  background: #fff;
  font-family: inherit;
`;

const TransactionType = styled.div<{ type: string }>`
  display: flex;
  gap: 10px;
  align-items: center;
`;

const IconCircle = styled.div<{ type: string }>`
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: grid;
  place-items: center;

  background: ${(p) =>
    p.type === "Received"
      ? "#E6FFF0"
      : p.type === "Sent"
        ? "#FFEAE6"
        : "#EEF5FF"};

  color: ${(p) =>
    p.type === "Received"
      ? "#00A859"
      : p.type === "Sent"
        ? "#BE3800"
        : "#007CDF"};
`;

const TypeText = styled.p`
  margin: 0;
  font-weight: 500;
  color: #00225a;
`;

const UserRow = styled.div`
  display: flex;
  gap: 6px;
  align-items: center;
`;

const AmountText = styled.span<{ type: string }>`
  color: ${(p) =>
    p.type === "Received"
      ? "#00A859"
      : p.type === "Sent"
        ? "#BE3800"
        : "#007CDF"};
  font-weight: 500;
`;

const ViewLink = styled.span`
  color: #007cdf;
  cursor: pointer;
`;
