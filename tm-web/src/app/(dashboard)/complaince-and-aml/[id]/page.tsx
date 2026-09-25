"use client";
import styled from "styled-components";

import { FaRegSmile, FaUserCircle } from "react-icons/fa";
import { FaArrowDown, FaArrowLeft } from "react-icons/fa6";
import { Dropdown } from "antd";
import { DownOutlined } from "@ant-design/icons";
import type { MenuProps } from "antd";
import { ImAttachment } from "react-icons/im";

import PrimaryButton from "@/components/PrimaryButton";
import {
  AiOutlineCheckCircle,
  AiOutlineClockCircle,
  AiOutlineCloseCircle,
} from "react-icons/ai";
import { useState } from "react";
import type { ReactElement } from "react";
import { useRouter } from "next/navigation";
import { TbSend2 } from "react-icons/tb";

type Status = {
  label: string;
  color: string;
  icon: ReactElement;
};

const TransactionDetailsPage = () => {
  const router = useRouter();
  const [status, setStatus] = useState<Status>({
    label: "Cleared",
    color: "#16b364",

    icon: <AiOutlineCheckCircle size={16} color="#16b364" />,
  });

  const items: MenuProps["items"] = [
    {
      key: "cleared",
      label: (
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: 8,
            color: "#00A859",
          }}
        >
          <AiOutlineCheckCircle size={16} color="#00A859" />
          Cleared
        </div>
      ),
      onClick: () =>
        setStatus({
          label: "Cleared",
          color: "#00A859",
          icon: <AiOutlineCheckCircle size={16} color="#00A859" />,
        }),
    },
    {
      key: "flagged",
      label: (
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: 8,
            color: "#ff4d4f",
          }}
        >
          <AiOutlineCloseCircle size={16} color="#ff4d4f" />
          Flagged
        </div>
      ),
      onClick: () =>
        setStatus({
          label: "Flagged",
          color: "#ff4d4f",
          icon: <AiOutlineCloseCircle size={16} color="#ff4d4f" />,
        }),
    },
    {
      key: "review",
      label: (
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: 8,

            color: "#007CDF",
          }}
        >
          <AiOutlineClockCircle size={16} color="#007CDF" />
          Under Review
        </div>
      ),
      onClick: () =>
        setStatus({
          label: "Under Review",
          color: "#007CDF",
          icon: <AiOutlineClockCircle size={16} color="#007CDF" />,
        }),
    },
  ];

  return (
    <Container>
      <FaArrowLeft onClick={() => router.back()} cursor="pointer" />
      <Header>
        <TitleSection>
          <TransactionID>Transaction ID 9085924u0u</TransactionID>
          <DateText>
            Date: <span>23 Sep 2023, 8:17:20</span>
          </DateText>
        </TitleSection>
        <PrimaryButton buttonStyle={{ width: "15%" }}>Update</PrimaryButton>
      </Header>
      <Content>
        <LeftSection>
          <SectionTitle>Comments</SectionTitle>

          <Comment>
            <User>
              <FaUserCircle size={22} color="#3F84F8" />
              <UserInfo>
                <Name>
                  Tom Onoja <span>(Admin)</span>
                </Name>
                <Timestamp>29 Sep, 2023 · 13:05</Timestamp>
              </UserInfo>
            </User>
            <Text>
              Lorem ipsum dolor sit amet consectetur. Sit facilisi ac urna
              tincidunt. Gravida sed auctor id commodo purus pharetra massa
              faucibus.
            </Text>
          </Comment>

          <Comment>
            <User>
              <FaUserCircle size={22} color="#3F84F8" />
              <UserInfo>
                <Name>
                  Obi Enechi <span>(Admin)</span>
                </Name>
                <Timestamp>28 Sep, 2023 · 16:05</Timestamp>
              </UserInfo>
            </User>
            <Text>
              Lorem ipsum dolor sit amet consectetur. Sit facilisi ac urna
              tincidunt. Gravida sed auctor id commodo purus pharetra massa
              faucibus.
            </Text>
          </Comment>

          <CommentInput>
            <textarea placeholder="Write your comment..." />

            <Actions>
              <FaRegSmile color="#828282" size={18} />

              <IconButton as="label">
                <ImAttachment color="#828282" size={18} />

                <input type="file" hidden />
              </IconButton>
            </Actions>

            <SendWrapper>
              <TbSend2 color="#828282" size={22} />
            </SendWrapper>
          </CommentInput>
        </LeftSection>

        <RightSection>
          <DetailsCard>
            <SectionTitle>Details</SectionTitle>
            <DetailRow1>
              <Label1>Price</Label1>
              <Value className="green">3,000.000 XBN</Value>
            </DetailRow1>

            <DetailRow1>
              <Label1>Transaction Type</Label1>
              <StatusPill>
                {" "}
                <FaArrowDown color="00A859" className="icon" /> Received
              </StatusPill>
            </DetailRow1>

            <DetailRow>
              <Label>Description</Label>
              <Input disabled value="Received from Obi" />
            </DetailRow>

            <DetailRow>
              <Label>Status</Label>
              <Dropdown menu={{ items }} trigger={["click"]}>
                <DropdownButton
                  style={{
                    color: status.color,
                  }}
                >
                  <div
                    style={{ display: "flex", alignItems: "center", gap: 6 }}
                  >
                    {status.icon}
                    {status.label}
                  </div>
                  <DownOutlined
                    style={{ fontSize: "12px", color: "#828282" }}
                  />
                </DropdownButton>
              </Dropdown>
            </DetailRow>

            <DetailRow>
              <Label>Cleared By</Label>
              <ClearedByBox>
                <User>
                  <FaUserCircle size={20} color="#3F84F8" />
                  <UserInfo>
                    <Name>Tom Onoja</Name>
                    <Role>Super Admin 29 Sep, 2023, 13:05</Role>
                  </UserInfo>
                </User>
              </ClearedByBox>
            </DetailRow>
          </DetailsCard>
        </RightSection>
      </Content>
    </Container>
  );
};

export default TransactionDetailsPage;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;
  border-radius: 24px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const BackArrow = styled.div`
  font-size: 24px;
  cursor: pointer;
`;

const TitleSection = styled.div`
  flex: 1;
  padding-left: 16px;
`;

const TransactionID = styled.h2`
  margin: 0;
  font-weight: 700;
  font-style: Bold;
  font-size: 24px;

  color: #00225a;
`;

const DateText = styled.p`
  font-size: 14px;
  color: #999;
  margin-top: 4px;

  span {
    font-weight: 500;
    color: #00225a;
  }
`;

const Content = styled.div`
  display: flex;
  margin-top: 32px;
  gap: 32px;
`;

const LeftSection = styled.div`
  flex: 2;
`;

const SectionTitle = styled.h3`
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #00225a;
`;

const Comment = styled.div`
  background: #fff;
  border-radius: 10px;
  padding: 12px 16px;
  margin-bottom: 16px;
`;

const User = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const UserInfo = styled.div``;

const Name = styled.p`
  font-weight: 600;
  font-size: 14px;
  margin: 0;
  color: #00225a;

  span {
    font-weight: 400;
    color: #00225a;
    margin-left: 4px;
  }
`;

const Role = styled.p`
  margin: 0;
  font-size: 13px;
  color: #00225a;
`;

const Timestamp = styled.p`
  font-size: 12px;
  color: #00225a;
  margin: 2px 0 0;
`;

const Text = styled.p`
  font-size: 14px;
  color: #00225a;
  margin-top: 8px;
`;

const CommentInput = styled.div`
  margin-top: 24px;
  position: relative;
  border-radius: 10px;
  border: 1px solid #ccc;
  min-height: 120px;

  textarea {
    width: 100%;
    padding: 10px;
    border: none;
    outline: none;
    resize: none;
    min-height: 80px;
    font-family: inherit;
  }
`;

const SendWrapper = styled.div`
  position: absolute;
  right: 16px;
  bottom: 10px;
  cursor: pointer;
`;

const RightSection = styled.div`
  flex: 1;
  background: #f2f6f9;
  border-radius: 10px;
`;

const DetailsCard = styled.div`
  padding: 20px;
`;

const DetailRow1 = styled.div`
  margin-bottom: 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const DetailRow = styled.div`
  margin-bottom: 18px;
`;

const Label = styled.p`
  font-size: 14px;
  color: #00225a;
  margin-bottom: 6px;

  font-weight: 600;
`;

const Label1 = styled.p`
  font-size: 14px;
  color: #828282;
  margin-bottom: 6px;

  font-weight: 400;
`;

const Value = styled.p`
  font-size: 16px;
  font-weight: 600;
  color: #2e8b57;

  &.green {
    color: #00a859;
  }
`;

const StatusPill = styled.p`
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 500;
  color: #00a859;
  .icon {
    background-color: #00a85933;
    color: #00a859;
    border-radius: 50%;
    padding: 4px;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
`;

const Input = styled.input`
  width: 100%;
  padding: 8px 12px;
  border-radius: 6px;
  border: none;
  font-family: inherit;
  background: #fff;
  color: #00225a;
`;

const ClearedByBox = styled.div`
  background: #fff;
  padding: 12px;
  border-radius: 8px;
`;

const DropdownButton = styled.button`
  background: #fff;
  color: #3f84f8;
  padding: 12px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
  border: none;
  cursor: pointer;
  font-family: inherit;
  width: 100%;
`;
const Actions = styled.div`
  display: flex;
  gap: 4px;
  position: absolute;
  left: 16px;
  bottom: 10px;
`;

const IconButton = styled.button`
  border: none;
  background: transparent;
  cursor: pointer;
  color: #666;

  &:hover {
    color: #007cdf;
  }
`;
