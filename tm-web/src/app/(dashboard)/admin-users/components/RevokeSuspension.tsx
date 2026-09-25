import { Modal, showErrorToast, showSuccessToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import { useUnsuspendAdminMutation } from "@/redux/api/admin/admin";
import { useRouter } from "next/navigation";
import { Router } from "next/router";
import React from "react";

import styled from "styled-components";

interface RevokeAdminProps {
  revoke: boolean;
  setRevoke: React.Dispatch<React.SetStateAction<boolean>>;
  adminName: string;
  adminEmail: string;
}

const RevokeSuspension: React.FC<RevokeAdminProps> = ({
  revoke,
  setRevoke,
  adminName,
  adminEmail,
}) => {
  const [unsuspendAdmin, { isLoading, isError }] = useUnsuspendAdminMutation();
  const router = useRouter();
  const handleRevoke = async () => {
    try {
      await unsuspendAdmin({ email: adminEmail }).unwrap();
      setRevoke(false);
      showSuccessToast("Suspension successfully revoked!");
      router.push("/admin-users");
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
    <>
      <Modal
        title="Revoke Suspension"
        isOpen={revoke}
        onClose={() => setRevoke(false)}
      >
        <Text>
          Are you sure you want to revoke the suspension on admin user
          <StyledSpan> {adminName}</StyledSpan>
        </Text>
        <PrimaryButton onClick={handleRevoke}>
          {isLoading ? "Revoking..." : "Yes, revoke suspension"}
        </PrimaryButton>
        {isError && (
          <ErrorText>Failed to revoke suspension. Please try again.</ErrorText>
        )}
      </Modal>
    </>
  );
};

export default RevokeSuspension;

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

const ErrorText = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  letter-spacing: 0.10000000149011612px;
  text-align: center;
  color: #be3800;
`;
