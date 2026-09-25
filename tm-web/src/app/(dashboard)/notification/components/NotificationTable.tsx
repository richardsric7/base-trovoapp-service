"use client";

import CustomTable from "@/components/CustomTable";
import React, { useMemo, useState } from "react";
import {
  AiOutlineCheckCircle,
  AiOutlineClockCircle,
  AiOutlineCloseCircle,
} from "react-icons/ai";
import styled from "styled-components";
import NotificationDetailsModal from "./NotificationDetailsModal";

const NotificationTable = () => {
  const [notificationDetails, setNotificationDetails] =
    useState<boolean>(false);
  const columns = [
    {
      title: "Notification Title",
      dataIndex: "notificationTitle",
      width: "20%",
      render: (_: any, record: any) => (
        <NotificationTitle onClick={() => setNotificationDetails(record)}>
          {record.notificationTitle}
        </NotificationTitle>
      ),
    },

    {
      title: "Notification Type",
      dataIndex: "notificationType",
      key: "notificationType",
    },
    {
      title: "Recipients",
      dataIndex: "recipients",
      key: "recipients",
    },
    {
      title: "Delivery Status",
      dataIndex: "status",
      render: (_: any, record: any) => {
        return (
          <StatusWithIcon status={record.status} />
          // <StatusWithIcon status={record.status}>
          //   {record.status}
          // </StatusWithIcon>
        );
      },
      key: "status",
    },
    {
      title: "Scheduled Date",
      dataIndex: "scheduledDate",
      key: "scheduledDate",
    },
  ];

  const dataSource = useMemo(() => {
    return [
      {
        notificationTitle: "Version Update",
        notificationType: "Update",
        recipients: "All Users [200]",
        scheduledDate: "23 Sep 2023",
        status: "Failed",
      },
      {
        notificationTitle: "Upgrade to Trovo Diamond",
        notificationType: "Upgrade",
        recipients: "Trovo Platinum Users",
        scheduledDate: "23 Sep 2023",
        status: "Successful",
      },
      {
        notificationTitle: "Version Update",
        notificationType: "Update",
        recipients: "All Users [200]",
        scheduledDate: "23 Sep 2023",
        status: "Pending",
      },
    ];
  }, []);

  return (
    <Container>
      <CustomTable columns={columns} dataSource={dataSource} />
      {notificationDetails && (
        <NotificationDetailsModal
          notificationDetails={notificationDetails}
          setNotificationDetails={setNotificationDetails}
        />
      )}
    </Container>
  );
};

export default NotificationTable;

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
    width: 15%;
  }

  th:nth-child(3),
  td:nth-child(3) {
    width: 12%;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 14%;
    text-align: center;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 17%;
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

const DeliveryStatus = styled.p<{ status: string }>`
  color: ${({ status }) =>
    status === "Failed"
      ? "#FF4D4D"
      : status === "Successful"
      ? "#00A859"
      : "#007CDF"};
  font-weight: 500;
  background-color: ${({ status }) =>
    status === "Failed"
      ? "#BE38001A"
      : status === "Successful"
      ? "#00A8591A"
      : "#F2F6F9"};
  font-weight: 500;
  padding: 8px 0;
  width: 120px;
  border-radius: 8px;
  text-align: center;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
`;
const NotificationTitle = styled.p`
cursor:pointer`;
const StatusWithIcon = ({ status }: { status: string }) => {
  return (
    <DeliveryStatus status={status}>
      {status === "Successful" && <AiOutlineCheckCircle />}
      {status === "Pending" && <AiOutlineClockCircle />}
      {status === "Failed" && <AiOutlineCloseCircle />}
      {status}
    </DeliveryStatus>
  );
};
