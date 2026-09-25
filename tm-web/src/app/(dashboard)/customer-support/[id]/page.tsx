"use client";
import styled from "styled-components";
import { MdVerified } from "react-icons/md";
import { FaPaperPlane } from "react-icons/fa";
import { ArrowLeft } from "iconsax-react";
import { DropdownSelect } from "@/components";
import { useState } from "react";
import { useRouter } from "next/navigation";

const SupportDetailsPage = () => {
  const [selectValues, setSelectValues] = useState<any>({});
  const router = useRouter();
  const filterCriteria = {
    issueType: ["Billing", "Suspended"],
    urgency: ["High", "Medium", "Normal"],
    status: ["Successful", "Pending", "Failed"],
  };
  return (
    <PageWrapper>
      <BackArrow onClick={() => router.back()}>
        <ArrowLeft color="#00225A" />
      </BackArrow>
      <TopSection>
        <LeftHeaderContent>
          <TicketTitle>T-001</TicketTitle>
          <Subtitle>Failed Asset Tokenization</Subtitle>

          <MetaInfo>
            <p>
              Reported On: <span>23 Sep 2023, 8:17</span>
            </p>
            <p>
              Last Updated On: <span>29 Sep 2023, 13:05</span>
            </p>
          </MetaInfo>
        </LeftHeaderContent>

        <ResolveButton>Resolve Ticket</ResolveButton>
      </TopSection>

      <BodySection>
        {/* LEFT COLUMN */}
        <LeftContent>
          {/* MESSAGES */}
          {["23 Sep, 2023", "29 Sep, 2023"].map((date) => (
            <MessageGroup key={date}>
              <DateHeader>{date}</DateHeader>

              <MessageBubble>
                <UserAvatar />
                <MessageContent>
                  <UserRows>
                    <UserRow>
                      <Username>Kennis</Username>
                      <MdVerified size={12} color="#007cdf" />
                    </UserRow>
                    <Fullname>Kennis Maduka</Fullname>
                    <UserMeta>
                      <span>Customer</span>
                      <DividerDot />
                      <span>23 Sep 2023, 13:05</span>
                    </UserMeta>
                  </UserRows>
                  <MessageText>
                    Lorem ipsum dolor sit amet consectetur. Sit facilisi ac urna
                    tincidunt. Gravida sed auctor id commodo purus pharetra
                    massa faucibus.
                  </MessageText>
                </MessageContent>
              </MessageBubble>

              <MessageBubble>
                <UserAvatar />
                <MessageContent>
                  <UserRows>
                    <UserRow>
                      <Username>Tom</Username>
                      <MdVerified size={12} color="#007cdf" />
                    </UserRow>
                    <UserMeta>
                      <span>Super Admin</span>
                      <DividerDot />
                      <span>23 Sep 2023, 13:05</span>
                    </UserMeta>
                  </UserRows>
                  <MessageText>
                    Gravida sed auctor id commodo purus pharetra massa faucibus.
                  </MessageText>
                </MessageContent>
              </MessageBubble>
            </MessageGroup>
          ))}

          {/* RESPONSE BOX */}
          <ResponseSection>
            <InputBox placeholder="Write your response..." />
            <ResponseActions>
              <ActionIcon>
                <FaPaperPlane />
              </ActionIcon>
            </ResponseActions>
          </ResponseSection>
        </LeftContent>

        {/* RIGHT SIDEBAR */}
        <Sidebar>
          <DetailsTitle>Details</DetailsTitle>
          <DetailsNote>
            Select the options under each category that suit the current state
            of the ticket
          </DetailsNote>

          <SelectGroup>
            <DropdownSelect
              options={filterCriteria.issueType}
              labelText="Issue Type"
              placeholder="All"
              value={selectValues["issueType"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, issueType: item })
              }
            />
          </SelectGroup>

          <SelectGroup>
            <DropdownSelect
              options={filterCriteria.urgency}
              labelText="urgency"
              placeholder="All"
              value={selectValues["urgency"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, urgency: item })
              }
            />
          </SelectGroup>

          {/* Status Filter */}
          <SelectGroup>
            <DropdownSelect
              options={filterCriteria.status}
              placeholder="All"
              labelText=" Status"
              value={selectValues["status"] || ""}
              onSelect={(item) =>
                setSelectValues({ ...selectValues, status: item })
              }
            />
          </SelectGroup>
        </Sidebar>
      </BodySection>
    </PageWrapper>
  );
};

export default SupportDetailsPage;

const PageWrapper = styled.section`
  background-color: #ffffff;
  padding: 24px;
  border-radius: 24px;
`;

const TopSection = styled.div`
  display: flex;
  //   align-items: center;
  justify-content: space-between;
`;

const BackArrow = styled.div`
  font-size: 20px;
  margin-right: 20px;
  margin-bottom: 10px;
  cursor: pointer;
`;

const LeftHeaderContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;

const TicketTitle = styled.h1`
  margin: 0;
  font-weight: 700;
  font-style: Bold;
  font-size: 24px;
  line-height: 28px;
  color: #00225a;
`;

const Subtitle = styled.p`
  color: #00225a;
  font-weight: 600;
  font-style: SemiBold;
  font-size: 16px;
  line-height: 28px;
`;

const ResolveButton = styled.button`
  background-color: #00a859;
  color: white;
  border: none;
  height: 48px;
  width: 148px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
`;

const BodySection = styled.div`
  display: flex;
  margin-top: 10px;
  gap: 40px;
`;

const LeftContent = styled.div`
  flex: 3;
`;

const MetaInfo = styled.div`
  font-size: 14px;
  color: #828282;
  margin-bottom: 20px;
  p {
    padding: 4px 0;
  }

  span {
    color: #00225a;
    font-weight: 500;
  }
`;

const MessageGroup = styled.div`
  margin-bottom: 30px;
`;

const DateHeader = styled.h4`
  color: #00225a;
  text-align: center;
  position: relative;
  margin: 32px 0;
  font-weight: 500;
  font-size: 16px;

  &::before,
  &::after {
    content: "";
    position: absolute;
    top: 50%;
    width: 40%;
    height: 1px;
    background-color: #e5e5ef;
  }

  &::before {
    left: 0;
    transform: translateY(-50%);
  }

  &::after {
    right: 0;
    transform: translateY(-50%);
  }
`;

const MessageBubble = styled.div`
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
`;

const UserAvatar = styled.div`
  width: 40px;
  height: 40px;
  background-color: #007cdf;
  border-radius: 50%;
`;

const MessageContent = styled.div`
  flex: 1;
`;

const UserRow = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const UserRows = styled.div``;
const Username = styled.p`
  font-size: 14px;
  margin: 0;
  font-weight: 500;
  font-style: Medium;
  font-size: 14px;
  line-height: 16px;
  color: #00225a;
`;

const Fullname = styled.p`
  font-weight: 400;
  font-style: Regular;
  font-size: 11px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #00225a;
`;

const MessageText = styled.p`
  font-size: 14px;
  margin-top: 4px;
  color: #00225a;
`;

const ResponseSection = styled.div`
  margin-top: 30px;
  padding-top: 20px;
`;

const InputBox = styled.textarea`
  width: 100%;
  border: 1px solid #e5e5e5;
  border-radius: 12px;
  padding: 16px;
  font-size: 14px;
  resize: none;
  height: 100px;
  font-family: inherit;
`;

const ResponseActions = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const ActionIcon = styled.div`
  background-color: #007cdf;
  color: white;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
`;

const Sidebar = styled.aside`
  flex: 1;
  background: #f2f6f9;
  padding: 24px;
  height: 430px;
  border-radius: 16px;
  font-size: 14px;
`;

const DetailsTitle = styled.h3`
  color: #00225a;
  margin-bottom: 4px;
  font-weight: 600;
  font-size: 20px;
  line-height: 28px;
`;

const DetailsNote = styled.p`
  font-size: 12px;
  color: #00225a;
  margin-bottom: 20px;
  font-weight: 400;
  font-size: 14px;
  line-height: 24px;
`;

const SelectGroup = styled.div`
  display: flex;
  flex-direction: column;
  margin-bottom: 20px;

  label {
    font-weight: 600;
    margin-bottom: 6px;
  }
`;

const UserMeta = styled.p`
  font-size: 12px;
  color: #828282;
  display: flex;
  align-items: center;
  gap: 6px;
`;

const DividerDot = styled.span`
  width: 1px;
  height: 12px;
  background-color: #e5e5ef;
  display: inline-block;
`;
