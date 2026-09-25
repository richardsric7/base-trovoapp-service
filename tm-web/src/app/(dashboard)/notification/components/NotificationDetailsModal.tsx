import styled from "styled-components";

import React from "react";
import { Modal } from "@/components";
import previewImg from "@/assets/images/woman-unlocking-phone1.png";
import Image from "next/image";
import logo from "@/assets/images/TransperantLogo12.svg";

import { AiOutlineCheckCircle } from "react-icons/ai";
interface NotificationDetailsProps {
  notificationDetails: boolean;
  setNotificationDetails: React.Dispatch<React.SetStateAction<boolean>>;
}

const NotificationDetailsModal: React.FC<NotificationDetailsProps> = ({
  notificationDetails,
  setNotificationDetails,
}) => {
  return (
    <StyledModal
      title=""
      isOpen={notificationDetails}
      onClose={() => setNotificationDetails(false)}
    >
      <Card>
        <Header>
          <IconWrapper>
            <Image src={logo} alt="logo" />
          </IconWrapper>
          <HeaderContent>
            <Title>Upgrade to Trovo Diamond</Title>
            <Subtitle>
              Unlock the full potential of the Tropisoft system with Trovo
              Diamond
            </Subtitle>
            <UpgradeButton>Upgrade</UpgradeButton>
          </HeaderContent>
          <StyledImage src={previewImg} alt="tovotech" />
        </Header>
        <MetricsTitle>Details</MetricsTitle>
        <Grid>
          <GridItem>
            <Label>Notification type</Label>
            <Value>Feature</Value>
          </GridItem>
          <GridItem>
            <Label>Recipients</Label>
            <Value>Trovo Platinum Users (100)</Value>
          </GridItem>
          <GridItem>
            <Label>Scheduled date</Label>
            <Value>26 May, 2024</Value>
          </GridItem>
          <GridItem>
            <Label>Delivery status</Label>
            <StatusIndicator>
              <AiOutlineCheckCircle color="#00A859" size={12} />
              <StatusValue>Sent</StatusValue>
            </StatusIndicator>
          </GridItem>
        </Grid>

        <MetricsSection>
          <MetricsTitle>Metrics</MetricsTitle>
          <Grid>
            <GridItem>
              <Label>Open rate</Label>
              <Value>80%</Value>
            </GridItem>
            <GridItem>
              <Label>Click through rate</Label>
              <Value>32.8%</Value>
            </GridItem>
          </Grid>
        </MetricsSection>
      </Card>
    </StyledModal>
  );
};

export default NotificationDetailsModal;

const StyledModal = styled(Modal)`
  .modal-content {
    margin: 0 auto;
    max-width: 800px;
    width: 100%;
    height: 100vh;
    border-radius: 0;
    overflow-y: auto;
  }
`;
const Card = styled.div``;

const Header = styled.div`
  display: flex;
  align-items: flex-start;
  margin-bottom: 20px;
  //   margin-top: 20px;
  gap: 16px;
  border-radius: 16px;

  background-color: #dbecf9;
  padding: 10px;
`;

const IconWrapper = styled.div`
  background-color: #dbeafe;
  padding: 8px;
  border-radius: 8px;

  svg {
    width: 24px;
    height: 24px;
    color: #3b82f6;
  }
`;

const HeaderContent = styled.div`
  margin-left: 16px;
`;

const Title = styled.h2`
  font-size: 18px;
  font-weight: 600;
  margin: 0;
`;

const Subtitle = styled.p`
  font-size: 14px;
  color: #6b7280;
  margin: 4px 0 0 0;
`;

const UpgradeButton = styled.button`
  width: 100%;
  background-color: #3b82f6;
  color: white;
  padding: 8px 16px;
  border-radius: 6px;
  border: none;
  font-weight: 500;
  cursor: pointer;
  margin-bottom: 24px;
  margin-top: 24px;
  transition: background-color 0.2s;

  &:hover {
    background-color: #2563eb;
  }
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
`;

const GridItem = styled.div`
  margin-bottom: 16px;
`;

const Label = styled.p`
  color: #6b7280;
  font-size: 14px;
  margin: 0 0 4px 0;
`;

const Value = styled.p`
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  line-height: 24px;
  text-align: left;
  color: #007cdf;
`;

const StatusValue = styled.p`
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  line-height: 24px;
  text-align: left;
  color: #00a859;
`;

const StatusIndicator = styled.div`
  display: flex;
  align-items: center;
`;

const MetricsSection = styled.div`
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
`;

const MetricsTitle = styled.h3`
  font-size: 16px;
  font-weight: 600;
  line-height: 16px;
  margin: 0 0 16px 0;
  color: #00225a;
`;

const StyledImage = styled(Image)``;
