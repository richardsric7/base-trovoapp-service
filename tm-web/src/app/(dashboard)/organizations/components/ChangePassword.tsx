import { Modal, showErrorToast, showSuccessToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import Image from "next/image";
import React from "react";
import styled from "styled-components";
import passwordIcon from "@/assets/images/Password Manager.png";
import { useResetMemberPasswordMutation } from "@/redux/api/organizations";

interface ResetPasswordProps {
  isResetPassword: boolean;
  setIsRestPassword: React.Dispatch<React.SetStateAction<boolean>>;
  memberName: string;
  memberEmail: string;
  organizationId: string;
}

const ChangePassword: React.FC<ResetPasswordProps> = ({
  isResetPassword,
  setIsRestPassword,
  memberName,
  memberEmail,
  organizationId,
}) => {
  const [resetPassword, { isLoading }] = useResetMemberPasswordMutation();

  const resetMemberPassword = async () => {
    try {
      await resetPassword({
        email: memberEmail,
        organization_id: organizationId,
      }).unwrap();
      showSuccessToast("successfully sent!");
    } catch (error: any) {
      const msg =
        error?.data?.message ??
        error?.message ??
        error?.error ??
        "Something went wrong. Please try again.";
      showErrorToast(msg);
    }
  };
  return (
    <Modal
      title=""
      isOpen={isResetPassword}
      onClose={() => setIsRestPassword(false)}
      closeIconPosition="left"
    >
      <ModalContent>
        <Image src={passwordIcon} alt="password-icon" width={80} height={80} />
        <Title>Reset Password</Title>
        <Text>
          Resetting this password will log
          <StyledSpan> {memberName} </StyledSpan> out of all devices and require
          them to sign in again.
        </Text>
        <PrimaryButton
          buttonStyle={{ width: "100%" }}
          onClick={resetMemberPassword}
        >
          Continue
        </PrimaryButton>
      </ModalContent>
    </Modal>
  );
};

export default ChangePassword;

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
