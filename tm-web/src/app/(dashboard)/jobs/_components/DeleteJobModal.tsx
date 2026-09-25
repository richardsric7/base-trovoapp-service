"use client";

import { Modal, showErrorToast, showSuccessToast } from "@/components";
import Image from "next/image";
import React from "react";
import styled from "styled-components";
import deleted from "@/assets/images/deleted.svg";
import { useRouter } from "next/navigation";
import { useDeleteJobMutation } from "@/redux/api/jobs";

interface DeleteJobProps {
  openModal: boolean;
  setOpenModal: React.Dispatch<React.SetStateAction<boolean>>;
  jobTitle: string;
  jobId: string | number;
}

const DeleteJobModal: React.FC<DeleteJobProps> = ({
  openModal,
  setOpenModal,
  jobTitle,
  jobId,
}) => {
  const [deleteJob, { isLoading }] = useDeleteJobMutation();
  const router = useRouter();
  const handleSubmit = async () => {
    try {
      await deleteJob(jobId).unwrap();
      setOpenModal(false);
      router.push("/jobs");
      showSuccessToast("Job post deleted successfully!");
    } catch (error) {
      showErrorToast("Failed to delete job post");
      console.error("Delete failed", error);
    }
  };

  return (
    <Modal title="" isOpen={openModal} onClose={() => setOpenModal(false)}>
      <ModalContent>
        <Image src={deleted} alt="delete-icon" width={80} height={80} />
        <Title>Delete Job</Title>
        <Text>
          Are you sure you want to delete this job from the job posts? This will
          be removed entirely from the website.
        </Text>
        <DeleteButton onClick={handleSubmit}>
          {isLoading ? "Deleting" : "Yes, Delete"}
        </DeleteButton>
      </ModalContent>
    </Modal>
  );
};

export default DeleteJobModal;

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
