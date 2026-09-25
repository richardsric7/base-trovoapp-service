import React, { useState } from "react";
import styled from "styled-components";
import { DropdownSelect, Modal } from "@/components";
import SuccessMessage from "@/components/SuccessMessage";
import PrimaryButton from "@/components/PrimaryButton";
import { useSuspendUserMutation } from "@/redux/api/users";
import { FaCheck } from "react-icons/fa6";

interface SuspendUserModalProps {
  isOpen: boolean;
  setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  userEmail: string;
}
interface CheckBoxProps {
  isChecked: boolean;
}
export const SuspendUserModal: React.FC<SuspendUserModalProps> = ({
  isOpen,
  setIsOpen,
  userEmail,
}) => {
  const [showModal, setShowModal] = useState<boolean>(true);
  const [suspendSuccess, setSuspendSuccess] = useState<boolean>(false);
  const [violationValue, setViolationValue] = useState<string>(""); // Track selected label
  const [note, setNote] = useState<string>("");
  const [isChecked, setIsChecked] = useState<boolean>(false);
  const VIOLATION_OPTIONS = [
    { id: 1, label: "Violation of terms of service" },
    { id: 2, label: "Violation of community guidelines" },
    { id: 3, label: "Violation of KYC/AML policy" },
    { id: 4, label: "Violation of security policy" },
    { id: 5, label: "Violation of privacy policy" },
    { id: 6, label: "Violation of trading policy" },
    { id: 7, label: "Violation of payment policy" },
    { id: 8, label: "Violation of dispute resolution policy" },
    { id: 9, label: "Violation of support policy" },
  ];

  const toggleCheckbox = () => {
    setIsChecked((prev) => !prev);
  };
  // Map options to labels for DropdownSelect compatibility
  const violationLabels = VIOLATION_OPTIONS.map((option) => option.label);

  const [suspendUserMutation] = useSuspendUserMutation();

  const handleSuspendUser = async () => {
    // Find the suspension reason ID based on the selected label
    const selectedOption = VIOLATION_OPTIONS.find(
      (option) => option.label === violationValue
    );
    const suspension_reason_id = selectedOption ? selectedOption.id : null;

    if (!suspension_reason_id) {
      console.error("Please select a valid suspension reason.");
      return;
    }

    try {
      const payload = {
        email: userEmail,
        suspension_note: note,
        suspension_reason_id, // Add the suspension ID to the payload
      };
      console.log("Suspension payload:", payload);

      const response = await suspendUserMutation(payload).unwrap();
      console.log("Suspension response:", response);

      setSuspendSuccess(true);
      setIsOpen(false);
    } catch (error: unknown) {
      console.error("Failed to suspend user:", error);
      alert("Failed to suspend user.");
    }
  };

  return (
    <>
      {showModal && (
        <Modal
          title="Record Violation"
          isOpen={isOpen}
          onClose={() => setIsOpen(false)}
        >
          <SuspendCard>
            <FormGroup>
              <FormLabel>Violation </FormLabel>
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
                placeholder=""
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
      {suspendSuccess && (
        <SuccessMessage
          isOpen={suspendSuccess}
          setIsOpen={setSuspendSuccess}
          heading="Success"
          message="User has been successfully suspended"
          email={userEmail}
        />
      )}
    </>
  );
};
export default SuspendUserModal;

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
