import React from "react";
import Stepper from "./Stepper";
import groupIcon from "@/assets/images/8.png";
import styled from "styled-components";

const AdminActionsTab = () => {
  const activities = [
    {
      user: groupIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Obi Enechi changed role to ",
      work: "Super Admin",
    },
    {
      user: groupIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Obi Enechi changed role to  ",
      work: "Super Admin",
    },
  ];

  return (
    <div>
      <Heading>Administrative Actions</Heading>
      {activities.map((activity, index) => (
        <Stepper
          key={index}
          avatar={activity.user}
          time={activity.time}
          action={activity.action}
          spanText={activity.work}
          isLast={index === activities.length - 1} // Check if it's the last item
        />
      ))}
    </div>
  );
};

export default AdminActionsTab;
const Heading = styled.h1`
  font-weight: 600;
  font-style: SemiBold;
  font-size: 20px;
  leading-trim: NONE;
  line-height: 28px;
  letter-spacing: 0px;
  color: #00225a;
`;
