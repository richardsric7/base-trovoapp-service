import CustomTable from "@/components/CustomTable";
import { useMemo } from "react";
import styled from "styled-components";

const TransactionHistory = () => {
  const columns = [
    {
      title: "Fee",
      dataIndex: "fee",
      key: "fee",
      render: (_: any, record: any) => <FeeText>{record.fee || "N/A"}</FeeText>,
    },
    {
      title: "Customer",
      dataIndex: "customer",
      key: "customer",

      render: (_: any, record: any) => (
        <UserInfoSection>
          <Avatar />
          <div>
            <UserName>{record.customer || "N/A"}</UserName>
            <SubName>{record.subname || "N/A"}</SubName>
          </div>
        </UserInfoSection>
      ),
    },

    {
      title: "Date",
      dataIndex: "date",
      key: "date",
    },
    {
      title: "Transaction ID",
      dataIndex: "transactionId",
      key: "transactionId",
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 6 }, (_, index) => ({
      tokenId: index + 1,
      customer: `Marlon${index + 1}`,
      subname: `Mark Ovey`,
      fee: "+10,000.00 NGN",
      date: "2026-03-31",
      transactionId: "TXN-1234567890",
    }));
  }, []);
  return (
    <Container>
      <SubTitle>Transaction History</SubTitle>
      <CustomTable columns={columns} dataSource={dataSource} />
    </Container>
  );
};

export default TransactionHistory;
const Container = styled.div`
  background-color: #fff;
  padding: 32px;
  height: 100%;
  gap: 20px;
  border-radius: 24px;

  table {
    width: 100%;
    table-layout: fixed;

    th:nth-child(1),
    td:nth-child(1) {
      width: 14%;
      //  text-align: center;
    }

    th:nth-child(2),
    td:nth-child(2) {
      width: 20%;

      text-align: center;
    }

    th:nth-child(3),
    td:nth-child(3) {
      width: 18%;

      text-align: center;
    }
    th:nth-child(4),
    td:nth-child(4) {
      width: 10%;

      text-align: center;
    }
  }
`;

const SubTitle = styled.h3`
  font-size: 20px;
  font-weight: 600;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;

const UserInfoSection = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
  justify-content: center;
`;

const UserName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const SubName = styled.p`
  font-size: 10px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const Avatar = styled.div`
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const FeeText = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00a859;
`;
