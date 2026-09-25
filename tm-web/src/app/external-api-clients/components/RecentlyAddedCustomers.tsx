import styled from "styled-components";
import { useMemo } from "react";
import CustomTable from "@/components/CustomTable";
import { AiOutlineCheckCircle, AiOutlineClockCircle } from "react-icons/ai";

type KycStatus = "Pending" | " Completed";
const RecentlyAddedCustomers = () => {
  const columns = [
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
      title: "Balance",
      dataIndex: "balance",
      key: "balance",
    },

    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      width: "100%",
      render: (status: KycStatus) => {
        const isPending = status === "Pending";

        return (
          <StatusBadge $status={status}>
            {isPending ? <AiOutlineClockCircle /> : <AiOutlineCheckCircle />}
            {isPending ? "Pending" : "Completed"}
          </StatusBadge>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 6 }, (_, index) => ({
      tokenId: index + 1,
      customer: `Marlon${index + 1}`,
      subname: `Mark Ovey`,
      balance: "3,500 CNGN",
      status: index % 2 === 0 ? "Pending " : " Completed",
    }));
  }, []);
  return (
    <Container>
      <SubTitle>Recent Assets Assigned</SubTitle>
      <CustomTable columns={columns} dataSource={dataSource} />
    </Container>
  );
};

export default RecentlyAddedCustomers;
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
      width: 18%;

      text-align: center;
    }

    th:nth-child(3),
    td:nth-child(3) {
      width: 20%;

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

const StatusBadge = styled.span<{ $status: KycStatus }>`
  font-size: 12px;
  font-weight: 500;
  padding: 6px 10px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  gap: 6px;

  background-color: ${({ $status }) =>
    $status === "Pending" ? "#F2F6F9" : "#00A85926"};

  color: ${({ $status }) => ($status === "Pending" ? "#007CDF" : "#00A859")};
`;
