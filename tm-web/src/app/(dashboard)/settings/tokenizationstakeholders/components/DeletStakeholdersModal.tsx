"use client";

import { Modal, showErrorToast, showSuccessToast } from "@/components";
import Image from "next/image";
import deleted from "@/assets/images/deleted.svg";

import styled from "styled-components";
import {
  PartnerType,
  useDeletePartnerMutation,
} from "@/redux/api/tokenizationstakeholders";
import { useRouter } from "next/navigation";
interface DeleteAssetProps {
  openModal: boolean;
  setOpenModal: React.Dispatch<React.SetStateAction<boolean>>;
  stakeholderName?: string;
  stakeholderType: PartnerType;
  stakeholderId: number;
}

const DeleteStakeholder: React.FC<DeleteAssetProps> = ({
  openModal,
  setOpenModal,
  stakeholderName,
  stakeholderType,
  stakeholderId,
}) => {
  const router = useRouter();
  const [deletePartner, { isLoading }] = useDeletePartnerMutation();
  const handleSubmit = async () => {
    try {
      await deletePartner({
        type: stakeholderType,
        id: stakeholderId,
      }).unwrap();
      setOpenModal(false);
      showSuccessToast("Stakeholder deleted");
      router.push("/settings/tokenizationstakeholders");
    } catch (error: any) {
      const errorString =
        error?.response?.data?.message ||
        error?.data?.message ||
        error?.message ||
        error?.response?.data?.error ||
        error?.data?.error ||
        error?.error ||
        "An error occurred while updating.";
      const messageMatch = errorString.match(/message:(.*?)(\]$|$)/);
      let finalErrorMessage = errorString;
      if (messageMatch && messageMatch[1]) {
        finalErrorMessage = messageMatch[1].trim();
      }
      // Show the error toast.
      showErrorToast(finalErrorMessage);
    }
  };
  return (
    <Modal title="" isOpen={openModal} onClose={() => setOpenModal(false)}>
      <ModalContent>
        <Image src={deleted} alt="delete-icon" width={80} height={80} />
        <Title>Delete Stakeholder</Title>
        <Text>
          Are you sure you want to delete{" "}
          <StyledSpan> {stakeholderName} </StyledSpan> from the Tokenization
          Stakeholders list?
        </Text>
        <DeleteButton onClick={handleSubmit}>
          {isLoading ? "Deleting" : "Yes, Delete"}
        </DeleteButton>
      </ModalContent>
    </Modal>
  );
};

export default DeleteStakeholder;
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
