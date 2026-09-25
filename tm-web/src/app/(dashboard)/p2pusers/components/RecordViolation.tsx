import React, { useState } from "react";
import styled from "styled-components";
import { FaCheck } from "react-icons/fa";
import { Modal, DropdownSelect } from "@/components";

import PrimaryButton from "@/components/PrimaryButton";

interface RecordViolationsProps {
  isOpen: boolean;
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  setSuspendSuccess: React.Dispatch<React.SetStateAction<boolean>>;
}

interface CheckBoxProps {
  isChecked: boolean;
}

const violationLabels = [
  { label: "Inappropriate Language", value: "language" },
  { label: "Fraudulent Activity", value: "fraud" },
  { label: "Spam or Harassment", value: "spam" },
];

const RecordViolation: React.FC<RecordViolationsProps> = ({
  isOpen,
  setIsOpen,
  setSuspendSuccess,
}) => {
  const [violationValue, setViolationValue] = useState<string>("");
  const [note, setNote] = useState<string>("");
  const [isChecked, setIsChecked] = useState<boolean>(false);

  const toggleCheckbox = () => setIsChecked(!isChecked);

  const handleSuspendUser = () => {
    if (!violationValue || !note) {
      alert("Please select a violation and provide a description.");
      return;
    }
    // Simulate API call
    setTimeout(() => {
      setSuspendSuccess(true);
      setIsOpen(false);
    }, 1000);
  };

  const violationLabels = [
    "Inappropriate Language",
    "Fraudulent Activity",
    "Spam or Harassment",
  ];
  return (
    <>
      {isOpen && (
        <Modal
          title="Record Violation"
          isOpen={isOpen}
          onClose={() => setIsOpen(false)}
        >
          <SuspendCard>
            <FormGroup>
              <FormLabel>Violation</FormLabel>
              <DropdownSelect
                options={violationLabels}
                placeholder="Select"
                placeholderColor="#828282"
                labelText=""
                value={violationValue}
                onSelect={setViolationValue}
              />
            </FormGroup>
            <FormGroup>
              <FormLabel>Description</FormLabel>
              <TextArea
                placeholder="Write a description..."
                value={note}
                onChange={(e) => setNote(e.target.value)}
                type="text"
              />
            </FormGroup>
            <FormGroupCheck>
              <CheckBox isChecked={isChecked} onClick={toggleCheckbox}>
                {isChecked && <FaCheck color="#fff" size={10} />}
              </CheckBox>
              <FormLabelCheck>Suspend user</FormLabelCheck>
            </FormGroupCheck>
          </SuspendCard>
          <PrimaryButton onClick={handleSuspendUser}>Submit</PrimaryButton>
        </Modal>
      )}
    </>
  );
};

export default RecordViolation;

const SuspendCard = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
`;

const FormGroup = styled.div`
  margin-bottom: 16px;
`;

const FormGroupCheck = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;
const FormLabelCheck = styled.p`
  font-style: normal;
  font-weight: 500;
  font-size: 14px;
  line-height: 24px;
  color: #00225a;
  display: block;
`;
const FormLabel = styled.label`
  margin-bottom: 8px;
  font-style: normal;
  font-weight: 500;
  font-size: 14px;
  line-height: 24px;
  color: #00225a;
  display: block;
`;

const TextArea = styled.input`
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  padding: 8px 16px;
  width: 96%;
  height: 100px;
  resize: vertical;
  font-family: montserrat;
`;

const CheckBox = styled.div<CheckBoxProps>`
  width: 16px;
  height: 16px;
  border: 1.5px solid ${(props) => (props.isChecked ? "none" : "#bdbdbd")};
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  background-color: ${(props) => (props.isChecked ? "#007CDF" : "transparent")};
  color: white;
  transition: all 0.2s ease-in-out;
`;
