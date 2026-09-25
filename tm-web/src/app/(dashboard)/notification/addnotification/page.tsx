"use client";

import { useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { Setting2 } from "iconsax-react";
import { FaArrowLeft } from "react-icons/fa6";
import { FiFileText } from "react-icons/fi";
import styled from "styled-components";

import PrimaryButton from "@/components/PrimaryButton";
import SuccessMessage from "@/components/SuccessMessage";
import previewImg from "@/assets/images/woman-unlocking-phone1.png";
import dynamic from "next/dynamic";

// Define TypeScript interfaces for props
interface IconWrapperProps {
  isActive: boolean;
}

interface DescriptionTextProps {
  isActive: boolean;
}

const Content = dynamic(() => import("../components/Content"), {
  ssr: false, // Important for Quill editor
});

const Settings = dynamic(() => import("../components/Settings"), {
  ssr: false,
});

const AddNotificationPage = () => {
  const [activeSection, setActiveSection] = useState<string>("Content");
  const [publish, setPublish] = useState<boolean>(false);

  const publishNotification = () => {
    console.log("ok published");
    setPublish(true);
  };
  return (
    <Container>
      <BackButton href="/notification">
        <FaArrowLeft />
      </BackButton>
      <HeaderTitle>Add New Notification</HeaderTitle>

      <MainContent>
        <SidebarMenu>
          <IconWrapper
            isActive={activeSection === "Content"}
            onClick={() => setActiveSection("Content")}
          >
            <FiFileText
              size="20"
              color={activeSection === "Content" ? "#007CDF" : "#828282"}
            />
            <DescriptionText isActive={activeSection === "Content"}>
              Content
            </DescriptionText>
          </IconWrapper>

          <IconWrapper
            isActive={activeSection === "Setting"}
            onClick={() => setActiveSection("Setting")}
          >
            <Setting2
              size="20"
              color={activeSection === "Setting" ? "#007CDF" : "#828282"}
            />
            <DescriptionText isActive={activeSection === "Setting"}>
              Setting
            </DescriptionText>
          </IconWrapper>
        </SidebarMenu>

        <ContentSection>
          {activeSection === "Content" && <Content />}
          {activeSection === "Setting" && <Settings />}
        </ContentSection>

        <PreviewSection>
          <PreviewHeaderWrapper>
            <PreviewTitle>Preview</PreviewTitle>
            <PrimaryButton
              onClick={publishNotification}
              buttonStyle={{
                padding: "8px 10px",
                width: "200px",
              }}
            >
              Publish Notification
            </PrimaryButton>
          </PreviewHeaderWrapper>
          <PreviewContent>
            <StyledImage src={previewImg} alt="Preview" />
            <PreviewTitle>Upgrade to Trovo Diamond</PreviewTitle>
            <PreviewText>
              Lorem ipsum dolor sit amet consectetur. Duis fringilla elementum
              viverra sed condimentum id risus lorem dictumst
            </PreviewText>
            <PrimaryButton>Upgrade</PrimaryButton>
          </PreviewContent>
        </PreviewSection>

        {publish && (
          <SuccessMessage
            isOpen={publish}
            setIsOpen={setPublish}
            heading="Success"
            message="You have successfully published the Notification titled "
            email="Upgrade Version"
          />
        )}
      </MainContent>
    </Container>
  );
};

export default AddNotificationPage;
const Container = styled.section``;

const HeaderTitle = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 20px 0;
  color: #00225a;
`;

const BackButton = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
  padding-bottom: 20px;
  margin-top: 20px;
`;

const MainContent = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 20px;
`;

const SidebarMenu = styled.div`
  background: #ffffff;
  border-radius: 24px;
  width: 150px;
  height: 300px;
  padding: 24px 16px;
`;

const ContentSection = styled.div`
  background: #ffffff;
  padding: 30px;
  width: 80%;
  height: 800px;
  margin: 0 auto;
  gap: 20px;
  border-radius: 24px;
`;

const PreviewSection = styled.div`
  background: #ffffff;
  padding: 30px;
  gap: 20px;
  width: 28%;
  border-radius: 24px;
`;

const PreviewHeaderWrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: 30px;
`;

const PreviewTitle = styled.p`
  font-size: 20px;
  font-weight: 700;
  color: #00225a;
  margin: 20px 0 0 0;
  text-align: center;
`;

const PreviewContent = styled.div`
  background-color: #f2f6f9;
  border-radius: 24px;
  margin-top: 50px;
  padding: 24px 8px;
`;

const DescriptionText = styled.p<DescriptionTextProps>`
  font-size: 16px;
  font-weight: 500;
  color: ${(props) => (props.isActive ? "#007CDF" : "#828282")};
  margin: 2px 0;
`;

const IconWrapper = styled.div<IconWrapperProps>`
  cursor: pointer;
  color: ${(props) => (props.isActive ? "#007CDF" : "#828282")};
  background-color: ${(props) => (props.isActive ? "#DBECF9" : "transparent")};
  width: 120px;
  height: 120px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  border-radius: 16px;
`;

const PreviewText = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 24px;
  letter-spacing: 0.1px;
  text-align: center;
  color: #00225a;
`;

const StyledImage = styled(Image)`
  width: 134px;
  height: 100px;
  margin: 0 auto 10px auto;
  display: flex;
  align-items: center;
  justify-content: center;
`;
