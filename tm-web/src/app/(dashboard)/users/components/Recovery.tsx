import React from "react";
import Stepper from "../../admin-users/components/Stepper";
import recoveryIcon from "@/assets/images/recoveryGroup.svg";
import { FaCircleCheck } from "react-icons/fa6";
import styled from "styled-components";
const Recovery = () => {
  const activities = [
    {
      user: recoveryIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Enabled account recovery",
      work: "",
    },
    {
      user: recoveryIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Recovered account",
      work: "",
    },
    {
      user: recoveryIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Recovered account",
      work: "",
    },
    {
      user: recoveryIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Recovered account",
      work: "",
    },
  ];

  return (
    <Container>
      <Badge>
        {" "}
        Enabled {/* <StyledIcon> */} <FaCircleCheck />
        {/* </StyledIcon> */}
      </Badge>
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
    </Container>
  );
};

const Container = styled.div`
  background-color: #fff;
  padding: 20px;
  border-radius: 24px 24px 0 0;
`;

export default Recovery;

const Badge = styled.div`
  background-color: #007cdf1a;
  color: #007cdf;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 8px;
  width: 10%;
  border-radius: 5px;

  font-size: 12px;
  font-weight: 400;
  text-align: center;
`;
const StyledIcon = styled.span`
  // width: 20px;
  // height: 20px;
  // margin: 0 auto;
`;
