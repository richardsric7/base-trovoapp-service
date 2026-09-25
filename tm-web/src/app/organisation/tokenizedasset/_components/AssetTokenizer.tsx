"use client";

import SecondaryButton from "@/components/SecondaryButton";
import Image from "next/image";
import React from "react";
import styled from "styled-components";
import profile from "../../../../assets/images/profileblue.svg";
import { useRouter } from "next/navigation";
import { MdVerified } from "react-icons/md";

const dummyUser = {
  username: "brikkle",
  first_name: "Ayo",
  last_name: "Ogunleye",
  email: "ayo@brikkle.co",
  address: "GAYBKWD76M7L7OFXQCTE6NGKQQQ5Y54DWSVDK4DNC4DZ7FW6KBZAAJEF",
  verified: true,
};

const AssetTokenizer: React.FC = () => {
  const router = useRouter();

  return (
    <>
      <HeadingContent>
        <Heading>Asset Tokenizer</Heading>

        <SecondaryButton
          buttonStyle={{
            width: "150px",
            display: "flex",
            alignItems: "center",
            gap: "6px",
            marginRight: "20px",
            marginBottom: "16px",
          }}
          onClick={() => {
            router.push(`/users/${dummyUser.username}`);
          }}
        >
          <Image src={profile} alt="profile" />
          View Profile
        </SecondaryButton>
      </HeadingContent>

      <Wrapper>
        <AssetIdSection>
          <Avatar />

          <div>
            <UserName>
              {dummyUser.username}

              {dummyUser.verified && (
                <VerifiedTag>
                  <MdVerified color="#007CDF" />
                </VerifiedTag>
              )}
            </UserName>

            <FullName>
              {dummyUser.first_name} {dummyUser.last_name}
            </FullName>
          </div>
        </AssetIdSection>

        <UserDetails>
          <div>
            <Title>Email Address</Title>
            <Text>{dummyUser.email}</Text>
          </div>

          <div>
            <Title>Wallet Address</Title>
            <Text>{dummyUser.address}</Text>
          </div>
        </UserDetails>
      </Wrapper>
    </>
  );
};

export default AssetTokenizer;

const Wrapper = styled.section`
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
  margin: 16px 0px;
`;

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const AssetIdSection = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;
const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;

const UserName = styled.div`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

const FullName = styled.p`
  font-size: 11px;
  font-weight: 400;
  line-height: 16px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const UserDetails = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 15px;
  gap: 100px;
`;

const Title = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #828282;
`;

const Text = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;

  word-break: break-all;
  overflow-wrap: break-word;

  display: block;
`;

const VerifiedTag = styled.div``;
