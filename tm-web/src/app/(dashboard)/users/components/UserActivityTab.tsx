import React, { useState } from "react";
import Stepper from "../../admin-users/components/Stepper";
import groupIcon from "@/assets/images/Group.png";
import logoutIcon from "@/assets/images/circum_logout.svg";
import logoutIcon2 from "@/assets/images/circum_logout2.svg";
import accessIcon from "@/assets/images/access.svg";
import walletIcon from "@/assets/images/clarity_wallet-line.svg";
import UserWallet from "./UserWallet";
import styled from "styled-components";
import {
  IUserInfo,
  IWallet,
  useFetchUserByCriteriaQuery,
} from "@/redux/api/users";
import { useParams } from "next/navigation";
const UserActivityTab = () => {
  const params = useParams();
  const username = Array.isArray(params.username)
    ? params.username[0]
    : params.username;

  const { data, isLoading, error } = useFetchUserByCriteriaQuery({
    username: username || "",
  });
  const [visibleActivities, setVisibleActivities] = useState(3);

  // Ensure userInfo exists before rendering the component
  if (isLoading || !data?.data?.user_info) {
    return <div className="p-4 text-gray-500">Loading wallets...</div>;
  }

  const userInfo: IUserInfo = data.data.user_info;
  const wallets: IWallet[] = data.data.wallets;

  const activities = [
    {
      user: logoutIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Logged Out from Trovo wallet",
      work: "",
    },
    {
      user: groupIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Tokenized asset",
      work: "“Atlantis Real Estate”",
    },
    {
      user: accessIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Granted viewer access to",
      work: "“Ephizee, Tom, Tunde, Kennis“",
    },
    {
      user: walletIcon,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Added new minting subwallet",
      work: "“Atlantis _minting”",
    },
    {
      user: logoutIcon2,
      time: "9:45 AM, 2 Apr, 2024",
      action: "Signed In to Trovo wallet",
      work: "",
    },
  ];

  const handleSeeMore = () => {
    setVisibleActivities((prevVisible) =>
      prevVisible + 3 > activities.length ? activities.length : prevVisible + 3
    );
  };

  return (
    <Container>
      <HeaderContainer>
        <UserWallet wallets={wallets ?? []} userInfo={userInfo} />
      </HeaderContainer>
      <Content>
        <ContentHeading>
          <Heading>User Activity</Heading>
          {visibleActivities < activities.length && (
            <MoreBtn onClick={handleSeeMore}>See All</MoreBtn>
          )}
        </ContentHeading>
        {activities.slice(0, visibleActivities).map((activity, index) => (
          <Stepper
            key={index}
            avatar={activity.user}
            time={activity.time}
            action={activity.action}
            spanText={activity.work}
            isLast={
              index === Math.min(visibleActivities, activities.length) - 1
            }
          />
        ))}
      </Content>
    </Container>
  );
};

export default UserActivityTab;

const Container = styled.div``;

const HeaderContainer = styled.div`
  margin-bottom: 20px;
  background-color: #fff;
  padding: 20px;
  border-radius: 24px;
`;

const Content = styled.div`
  background-color: #fff;
  padding: 20px;
  border-radius: 24px 24px 0 0;
`;
const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  color: #00225a;
  margin-bottom: 10px;
`;
const ContentHeading = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 40px 0 20px 0;
`;
const MoreBtn = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  letter-spacing: 0.10000000149011612px;
  color: #007cdf;
  cursor: pointer;
`;
