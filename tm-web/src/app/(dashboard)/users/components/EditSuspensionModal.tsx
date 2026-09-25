import PrimaryButton from "@/components/PrimaryButton";
import { Modal } from "@/components";
import React, { useEffect, useState } from "react";
import styled from "styled-components";
import SuccessMessage from "@/components/SuccessMessage";

interface EditModalProps {
  isOpen: boolean;
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

export const EditSuspensionModal: React.FC<EditModalProps> = ({
  isOpen,
  setIsOpen,
}) => {
  const [editSuccess, setEditSuccess] = useState<boolean>(false);
  const [showModal, setShowModal] = useState<boolean>(true);
  const [startDate, setStartDate] = useState<string>("");
  const [endDate, setEndDate] = useState<string>("");

  const handleSubmit = () => {
    setShowModal(false);
    setEditSuccess(true);
  };

  return (
    <>
      {showModal && (
        <Modal
          title="Edit Suspension Duration"
          isOpen={isOpen}
          onClose={() => setIsOpen(false)}
        >
          <FormGroup>
            <Label>Duration</Label>
            <DateInputGroup>
              <Input
                type="date"
                required
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
              />
              <Dash>-</Dash>
              <Input
                type="date"
                required
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
              />
            </DateInputGroup>
          </FormGroup>

          <PrimaryButton onClick={handleSubmit}>Submit</PrimaryButton>
        </Modal>
      )}
      {editSuccess && (
        <SuccessMessage
          isOpen={editSuccess}
          setIsOpen={setEditSuccess}
          heading="Success"
          message="You have successfully changed the suspension duration for"
          email="odogwu" // You might want to pass the actual user email here
        />
      )}
    </>
  );
};

const FormGroup = styled.div``;

const DateInputGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  padding: 8px 12px;
`;

const Input = styled.input`
  flex: 1;
  height: 40px;
  border: none;
  outline: none;
  padding: 0 8px;
  cursor: pointer;
`;

const Dash = styled.span`
  font-size: 20px;
  align-self: center;
`;

const Label = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 24px;
  letter-spacing: 0.1px;
  text-align: left;
`;
