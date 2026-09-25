"use client";

import { Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import TruncatedText from "@/hooks/useTruncate";
import React from "react";
import { MdVerified } from "react-icons/md";
import styled from "styled-components";

interface ViewProps {
  isOpen: boolean;
  onClose: () => void;
}

const ViewModal: React.FC<ViewProps> = ({ isOpen, onClose }) => {
  const amountDetails = [
    {
      text: "Token Quantity",
      value: "1,000 TROV",
    },

    {
      text: "Current NAV/Token",
      value: "₦1,000",
    },
    {
      text: "Early Exit Penalty + Fees",
      value: "₦30,000",
    },

    {
      text: "Payout Token Price",
      value: "₦970",
    },

    {
      text: "Net Payout",
      value: "₦970,000",
    },
  ];

  const payoutDetails = [
    {
      text: "Payout Method",
      value: "Wallet",
    },

    {
      text: "Wallet ID",
      value: "Johnnydoe",
    },
    {
      text: "Settlement ETA",
      value: "48 hours",
    },
  ];
  return (
    <>
      <Modal title="Early Exit Details" isOpen={isOpen} onClose={onClose}>
        <ModalContent>
          <Title>Requested By</Title>
          <UserInfoContainer>
            <Avatar />
            <div>
              <UserNameWrapper>
                <UserName>
                  {" "}
                  <TruncatedText text="Florence" maxLength={15} />
                  {/* {kyc_level === 0 ? (
                        ""
                      ) : ( */}
                  <MdVerified size={12} color="#007CDF" />
                  {/* )} */}
                </UserName>
              </UserNameWrapper>

              <FullName>Florence Zach</FullName>
            </div>
          </UserInfoContainer>

          <Divider />
          <Content>
            <SubTitle>Quantity/Amount Details</SubTitle>
            {amountDetails.map((item) => (
              <TextContent key={item.text}>
                <Text>{item.text}</Text>
                <Value>{item.value}</Value>
              </TextContent>
            ))}
          </Content>

          <Divider />
          <Content>
            <SubTitle>Payout Details</SubTitle>
            {payoutDetails.map((items) => (
              <TextContent key={items.text}>
                <Text>{items.text}</Text>
                <Value>{items.value}</Value>
              </TextContent>
            ))}
          </Content>
        </ModalContent>
        <UserNameWrapper>
          <Title>Processed By</Title>
          <UserInfoContainer>
            <Avatar />
            <div>
              <UserName>
                {" "}
                <TruncatedText text="Obi" maxLength={15} />
              </UserName>

              <FullName>Obi Enechi</FullName>
            </div>
          </UserInfoContainer>
        </UserNameWrapper>
      </Modal>
    </>
  );
};

export default ViewModal;

const ModalContent = styled.div``;

const Title = styled.h1`
  font-weight: 600;
  font-size: 16px;
  color: #00225a;
`;

const SubTitle = styled.p`
  font-weight: 600;
  font-size: 16px;
  color: #00225a;
  padding-bottom: 10px;
`;

const UserInfoContainer = styled.div`
  background-color: #f2f6f9;
  height: 68px;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  column-gap: 10px;
  align-items: center;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const UserNameWrapper = styled.div`
  margin-top: 10px;
`;
const UserName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const FullName = styled.p`
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const Divider = styled.div`
  background-color: #e5e5ef;
  width: 100%;
  height: 1px;
  margin: 20px 0;
`;

const Content = styled.div``;

const TextContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
`;

const Text = styled.p`
  font-weight: 500;
  font-size: 14px;
  letter-spacing: 0px;
  color: #828282;
  padding-bottom: 16px;
`;

const Value = styled.p`
  font-weight: 500;
  font-size: 14px;
  color: #00225a;
  text-align: center;
`;
