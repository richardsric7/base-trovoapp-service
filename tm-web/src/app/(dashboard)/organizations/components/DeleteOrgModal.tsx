import { Modal, showErrorToast, showSuccessToast } from "@/components";
import Image from "next/image";
import { useRouter } from "next/navigation";
import React from "react";
import styled from "styled-components";
import deleted from "@/assets/images/deleted.svg";
import { useDeactivateOrganizationMutation } from "@/redux/api/organizations";

interface DeleteOrgProps {
  isDeleteModalOpen: boolean;
  setIsDeleteModalOpen: React.Dispatch<React.SetStateAction<boolean>>;
  orgId: string;
  orgName: string;
}
const DeleteOrgModal: React.FC<DeleteOrgProps> = ({
  isDeleteModalOpen,
  setIsDeleteModalOpen,
  orgId,
  orgName,
}) => {
  const [deactivateOrganization, { isLoading }] =
    useDeactivateOrganizationMutation();

  const handleDeactivate = async () => {
    if (!orgId) return;
    try {
      const response = await deactivateOrganization(orgId).unwrap();

      setIsDeleteModalOpen(false);
      showSuccessToast("Successfully Deactivated!");
    } catch (error: any) {
      console.error("Failed to deactivate organization!:", error);
      showErrorToast(error);
    }
  };

  return (
    <Modal
      title=""
      isOpen={isDeleteModalOpen}
      onClose={() => setIsDeleteModalOpen(false)}
    >
      <ModalContent>
        <Image src={deleted} alt="delete-icon" width={80} height={80} />
        <Title>Deactivate Organization</Title>
        <Text>
          Are you sure you want to delete <StyledSpan> {orgName} </StyledSpan>{" "}
          from organization list?
        </Text>
        <DeleteButton onClick={handleDeactivate} disabled={isLoading}>
          {isLoading ? "Deactivating" : "Yes, Deactivate"}
        </DeleteButton>
      </ModalContent>
    </Modal>
  );
};

export default DeleteOrgModal;

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

const DeleteButton = styled.button`
  width: 50%;
  height: 48px;
  gap: 10px;
  font-family: inherit;
  border-radius: 10px;
  background-color: #be3800;
  border: none;
  color: #fff;
  cursor: pointer;
`;
