import React from "react";
import { Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import styled from "styled-components";

interface RevokeUserProps {
  revokeSuspension: boolean;
  setRevokeSuspension: React.Dispatch<React.SetStateAction<boolean>>;
}
const RevokeUserSuspension: React.FC<RevokeUserProps> = ({
  revokeSuspension,
  setRevokeSuspension,
}) => {
  return (
    <>
      <Modal
        title="Revoke Suspension"
        isOpen={revokeSuspension}
        onClose={() => setRevokeSuspension(false)}
      >
        <Text>
          Are you sure you want to revoke the suspension on user
          <StyledSpan> Kennis Maduka? </StyledSpan>
        </Text>
        <PrimaryButton onClick={() => setRevokeSuspension(!revokeSuspension)}>
          Yes, revoke suspension
        </PrimaryButton>
      </Modal>
    </>
  );
};

export default RevokeUserSuspension;

const Title = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  text-align: center;
  color: #00225a;
  margin: 0;
`;

const Text = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  letter-spacing: 0.10000000149011612px;
  text-align: center;
  color: #00225a;
`;

const StyledSpan = styled.span`
  font-size: 14px;
  font-weight: 600;
`;
