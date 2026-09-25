"use client"
import SecondaryButton from "@/components/SecondaryButton";
import Image from "next/image";
import React from "react";
import styled from "styled-components";
import profile from "../../../../assets/images/profileblue.svg";
import { TokenizationRecord } from "@/redux/api/assettokenization";
import { useRouter } from "next/navigation";
import { useFetchUserByCriteriaQuery } from "@/redux/api/users";
import { MdVerified } from "react-icons/md";

interface AssetTokenizerProps {
  asset?: TokenizationRecord;
}
const AssetTokenizer: React.FC<AssetTokenizerProps> = ({ asset }) => {
  const { data: userData, isLoading } = useFetchUserByCriteriaQuery(
    { username: asset?.initiatorUsername || "" },
    { skip: !asset?.initiatorUsername }
  );
  const router = useRouter();
  return (
    <>
      <HeadingContent>
        <Heading>Asset Tokenizer</Heading>
        <SecondaryButton
          buttonStyle={{
            width: "150px",
            marginRight: "20px",
            display: "flex",
            alignItems: "center",
            gap: "2px",
            marginBottom: "16px",
          }}
          onClick={() => {
            if (asset?.initiatorUsername) {
              router.push(`/users/${asset.initiatorUsername}`);
            }
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
              {asset?.initiatorUsername}

              <VerifiedTag>
                {userData?.data?.user_info?.email ? (
                  ""
                ) : (
                  <MdVerified color="#007CDF" />
                )}
              </VerifiedTag>
            </UserName>

            {(userData?.data?.user_info?.first_name ||
              userData?.data?.user_info?.last_name) && (
              <FullName>
                {userData.data.user_info.first_name}{" "}
                {userData.data.user_info.last_name}
              </FullName>
            )}
          </div>
        </AssetIdSection>

        <UserDetails>
          <div>
            <Title>Email Address</Title>

            {userData?.data?.user_info?.email && (
              <Text>{userData.data.user_info.email}</Text>
            )}
          </div>

          <div>
            <Title>Wallet Address</Title>
            {userData?.data?.user_info?.email && (
              <Text>{userData.data.user_info.address}</Text>
            )}
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
  flex-direction: column;
  gap: 15px;
  margin-top: 15px;
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

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 600px;
  display: block;
`;

const VerifiedTag = styled.div``;
