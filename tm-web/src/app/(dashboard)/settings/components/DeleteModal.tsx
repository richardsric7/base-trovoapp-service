"use client";

import { Modal, showErrorToast, showSuccessToast } from "@/components";
import Image from "next/image";
import React from "react";
import styled from "styled-components";
import deleted from "@/assets/images/deleted.svg";
import { useDeleteMintingUserMutation } from "@/redux/api/minting";

interface DeletemintingUsersProps {
  openModal: boolean;
  setOpenModal: React.Dispatch<React.SetStateAction<boolean>>;
  username: string;
  id: number;
  role: string;
}

const DeleteModal: React.FC<DeletemintingUsersProps> = ({
  openModal,
  setOpenModal,
  username,
  id,
  role,
}) => {
  const [deleteInitiators, { isLoading }] = useDeleteMintingUserMutation();
  console.log("ids", id, role);

  const handleSubmit = async () => {
    try {
      await deleteInitiators({
        id: id,
        role: role,
      }).unwrap();
      setOpenModal(false);
      showSuccessToast("Successfully Deleted!");
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
    <Modal title="" isOpen={openModal} onClose={() => setOpenModal(false)}>
      <ModalContent>
        <Image src={deleted} alt="delete-icon" width={80} height={80} />
        <Title>Delete {role === "initiator" ? "Initiator" : "Approver"}</Title>
        <Text>
          Are you sure you want to delete <StyledSpan> {username} </StyledSpan>{" "}
          from the list of {role}?
        </Text>
        <DeleteButton onClick={handleSubmit}>
          {isLoading ? "Deleting" : "Yes, Delete"}
        </DeleteButton>
      </ModalContent>
    </Modal>
  );
};

export default DeleteModal;

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
