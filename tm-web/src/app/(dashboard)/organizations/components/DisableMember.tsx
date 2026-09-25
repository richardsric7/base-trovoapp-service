import { Modal } from "@/components";
import React from "react";
import blockIcon from "@/assets/images/block.png";
import Image from "next/image";
import styled from "styled-components";
import PrimaryButton from "@/components/PrimaryButton";

interface DisableUserProps {
  isDisableUser: boolean;
  setIsDisableUser: React.Dispatch<React.SetStateAction<boolean>>;
  memberName: string;
}

const DisableMember: React.FC<DisableUserProps> = ({
  setIsDisableUser,
  isDisableUser,
  memberName,
}) => {
  return (
    <Modal
      title=""
      isOpen={isDisableUser}
      onClose={() => setIsDisableUser(false)}
      closeIconPosition="left"
    >
      <ModalContent>
        <Image src={blockIcon} alt="block-icon" width={80} height={80} />
        <Title>Disable User</Title>
        <Text>
          Resetting this password will log
          <StyledSpan> {memberName} </StyledSpan> out of all devices and require
          them to sign in again.
        </Text>
        <PrimaryButton buttonStyle={{ width: "100%" }}>Continue</PrimaryButton>
      </ModalContent>
    </Modal>
  );
};

export default DisableMember;
const ModalContent = styled.div`
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: 10px;
`;
const Title = styled.h1`
  font-weight: 600;
  font-size: 20px;
  line-height: 20px;
  letter-spacing: 0.1px;
  color: #00225a;
`;
const Text = styled.p`
  font-weight: 400;
  font-size: 14px;
  line-height: 24px;
  letter-spacing: -0.5%;
  text-align: center;
  color: #00225a;
`;

const StyledSpan = styled.span`
  font-weight: 600;
  font-size: 14px;
  line-height: 24px;
  letter-spacing: -0.5%;
  text-align: center;
`;
