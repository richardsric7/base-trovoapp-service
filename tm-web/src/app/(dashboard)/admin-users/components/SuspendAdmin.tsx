"use client";
import React, { useState } from "react";
import styled from "styled-components";
import { DropdownSelect, Modal, showErrorToast } from "@/components";
import SuccessMessage from "@/components/SuccessMessage";
import { useSuspendAdminMutation } from "@/redux/api/admin/admin";
import PrimaryButton from "@/components/PrimaryButton";

interface SuspendAdminProps {
  suspendAdmin: boolean;
  setSuspendAdmin: React.Dispatch<React.SetStateAction<boolean>>;
  adminEmail: string;
  onSuccess: (email: string) => void;
}

const SuspendAdmin: React.FC<SuspendAdminProps> = ({
  suspendAdmin,
  setSuspendAdmin,
  adminEmail,
  onSuccess,
}) => {
  const [violationValue, setViolationValue] = useState<string>("");

  const [suspensionReasonId, setSuspensionReasonId] = useState<number>(0);

  const [note, setNote] = useState<string>("");
  const [adminData, setAdminData] = useState<{ email: string }>({
    email: adminEmail,
  });

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
  const [suspendAdminMutation] = useSuspendAdminMutation();
  // const { data: violationValue, isLoading } = useGetVoilationsQuery();

  const handleSuspendAdmin = async () => {
    if (!adminData.email) {
      showErrorToast("Email is required to suspend an admin.");
      return;
    }
    if (!violationValue) {
      showErrorToast("Please select a suspension reason.");
      return;
    }
    if (!note.trim()) {
      showErrorToast("Please provide a note explaining the suspension.");
      return;
    }

    //Find the suspension reason ID based on the selected label
    const selectedOption = VIOLATION_OPTIONS.find(
      (option) => option.label === violationValue
    );
    const suspension_reason_id = selectedOption?.id;
    if (!suspension_reason_id) {
      console.error("Invalid suspension reason selected.");

      return;
    }
    try {
      await suspendAdminMutation({
        email: adminData.email,
        suspension_note: note,
        suspension_reason_id,
      });
      // setSuspensionSuccess(true);
      onSuccess(adminEmail);
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
        title="Suspend Admin User"
        isOpen={suspendAdmin}
        onClose={() => setSuspendAdmin(false)}
      >
        <SuspendCard>
          <FormGroup>
            <FormLabel>Email</FormLabel>
            <StyledInput
              placeholder="Enter email"
              type="email"
              value={adminData.email}
              onChange={(e) =>
                setAdminData({ ...adminData, email: e.target.value })
              }
            />
          </FormGroup>
          <FormGroup>
            <FormLabel>Select Suspension</FormLabel>

            <DropdownSelect
              options={VIOLATION_OPTIONS.map((opt) => opt.label)}
              placeholder={"Select"}
              labelText=""
              value={violationValue}
              onSelect={(item) => {
                setViolationValue(item);
              }}
            />
          </FormGroup>

          <FormGroup>
            <FormLabel>Notes</FormLabel>
            <TextArea
              type="text"
              placeholder="Type the reason here"
              value={note}
              onChange={(e) => setNote(e.target.value)}
            />
          </FormGroup>
        </SuspendCard>
        <PrimaryButton onClick={handleSuspendAdmin}>Submit</PrimaryButton>
      </Modal>
    </>
  );
};

export default SuspendAdmin;

const SuspendCard = styled.div``;

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
  padding: 2px;
  width: 100%;
  outline: none;
  height: 100px;
  font-family: inherit;
`;

const FormGroup = styled.div`
  margin: 12px 0;
`;

const StyledInput = styled.input`
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  padding: 2px 8px;
  width: 100%;
  outline: none;
  height: 48px;
  font-family: inherit;

  &:disabled {
    background-color: #f5f5f5;
  }
`;
