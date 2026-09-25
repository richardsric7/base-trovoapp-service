import React from "react";
import styled from "styled-components";
import { FaX } from "react-icons/fa6";
import Image from "next/image";
import { Modal } from "@/components";
import sucessIcon from "@/assets/images/success.png";
import type { ReactNode } from "react";

interface SuccessMessageProps {
  isOpen: boolean;
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  heading?: string;
  message: string | ReactNode;
  email?: string;
  iconSrc?: string;
}

const SuccessMessage: React.FC<SuccessMessageProps> = ({
  isOpen,
  setIsOpen,
  heading = "Success",
  message = "You have successfully sent an invite to",
  email,
  iconSrc,
}) => {
  if (!isOpen) return null;

  const renderBoldText = (text: string, boldText: string[]) => {
    if (!boldText.length) return text;
  };

  return (
    <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} title="">
      <Container>
        <CloseIcon onClick={() => setIsOpen(false)} />
        <SuccessHeading>{heading}</SuccessHeading>

        <SuccessCard>
          <Image src={sucessIcon} alt="success-icon" />
          <SuccessText>
            {message}
            {email && (
              <>
                <br />
                <StyledSpan>{email}</StyledSpan>
              </>
            )}
          </SuccessText>
        </SuccessCard>
      </Container>
    </Modal>
  );
};

export default SuccessMessage;

const Container = styled.div`
  // background-color: #ffffff;
  // border-radius: 24px;
  // padding: 20px;
  // width: 480px;
  // position: relative;
`;

const SuccessCard = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin: 0 auto;
  align-items: center;
  justify-content: center;
  text-align: center;
`;
const ModalHeader = styled.div`
  display: flex;
  align-items: center;
  gap: 12rem;
`;

const SuccessHeading = styled.h1`
  font-size: 18px;
  font-weight: 600;
  line-height: 20px;
  color: #00225a;
  letter-spacing: 0.1px;
  text-align: center;
  padding-bottom: 6px;
`;

const SuccessText = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  color: #00225a;
  letter-spacing: 0.1px;
`;

const StyledSpan = styled.span`
  font-weight: 600;
`;

const CloseIcon = styled(FaX)`
  position: absolute;
  top: 20px;
  right: 20px;
  cursor: pointer;
  font-size: 16px;
  color: #00225a;
`;
