import React from "react";
import Stepper from "./Stepper";
import groupIcon from "@/assets/images/Group.png";
import styled from "styled-components";

const ActivityTab = () => {
  const activities = [
    {
      user: groupIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Resolve Appeal Order 9085924u0u in ",
      work: "Customer Support",
    },
    {
      user: groupIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Resolve Appeal Order 9085924u0u in ",
      work: "Customer Support",
    },
    {
      user: groupIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Resolve Appeal Order 9085924u0u in ",
      work: "Customer Support",
    },
    {
      user: groupIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Resolve Appeal Order 9085924u0u in ",
      work: "Customer Support",
    },
    {
      user: groupIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Resolve Appeal Order 9085924u0u in ",
      work: "Customer Support",
    },
  ];

  return (
    <div>
      <Heading>User Activity</Heading>
      {activities.map((activity, index) => (
        <Stepper
          key={index}
          avatar={activity.user}
          time={activity.time}
          action={activity.action}
          spanText={activity.work}
          isLast={index === activities.length - 1}
        />
      ))}
    </div>
  );
};

export default ActivityTab;

const Heading = styled.h1`
  font-weight: 600;
  font-style: SemiBold;
  font-size: 20px;
  leading-trim: NONE;
  line-height: 28px;
  letter-spacing: 0px;
  color: #00225a;
`;
