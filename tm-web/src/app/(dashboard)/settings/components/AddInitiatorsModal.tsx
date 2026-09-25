"use client";
import { Modal, showErrorToast, showSuccessToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import { useAddMintingUserMutation } from "@/redux/api/minting";

import React, { useState } from "react";
import styled from "styled-components";
interface AddInitiatorProps {
  add: boolean;
  setAdd: React.Dispatch<React.SetStateAction<boolean>>;
  setAdded: React.Dispatch<React.SetStateAction<boolean>>;
  username: string;
  setUsername: React.Dispatch<React.SetStateAction<string>>;
}

const AddInitiatorsModal: React.FC<AddInitiatorProps> = ({
  add,
  setAdd,
  setAdded,
  username,
  setUsername,
}) => {
  const [addinitiator, { isLoading: isAdding }] = useAddMintingUserMutation();

  const handleSubmit = async () => {
    try {
      await addinitiator({
        role: "initiator",
        username: username.trim(),
      }).unwrap();
      setAdded(true);
      setAdd(false);
      showSuccessToast(`Succefully added ${username} to initiator!`);
      setUsername("");
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
    <Modal title="Add Initiator" isOpen={add} onClose={() => setAdd(false)}>
      <Input
        type="text"
        placeholder="Enter user name ..."
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <PrimaryButton onClick={handleSubmit}>
        {isAdding ? "Adding Initiator..." : " Add Initiator"}
      </PrimaryButton>
    </Modal>
  );
};

export default AddInitiatorsModal;

const Input = styled.input`
  width: 100%;
  height: 48px;
  justify-content: space-between;
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  font-family: inherit;
  margin: 10px 0;
  padding: 0 6px;
`;
