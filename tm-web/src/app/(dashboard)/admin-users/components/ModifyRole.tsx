"use client";

import React, { useState } from "react";
import styled from "styled-components";
import { DropdownSelect, Modal } from "@/components";

import DateRangePicker from "./DateRangePicker";
import SuccessMessage from "@/components/SuccessMessage";
import PrimaryButton from "@/components/PrimaryButton";
interface ModifyAdminProps {
  modifyAdmin: boolean;
  setModifyAdmin: React.Dispatch<React.SetStateAction<boolean>>;
}
const ModifyRole: React.FC<ModifyAdminProps> = ({
  modifyAdmin,
  setModifyAdmin,
}) => {
  const [showModal, setShowModal] = useState<boolean>(true);
  const [suspendUser, setSuspendUser] = useState<boolean>(false);
  const [violationValue, setViolationValue] = useState<string>("");
  const handleSuspendAdmin = () => {
    setSuspendUser(true);
    setShowModal(false);
  };

  return (
    <>
      {showModal && (
        <Modal
          title="Modify Suspension"
          isOpen={modifyAdmin}
          onClose={() => setModifyAdmin(false)}
        >
       
          <SuspendCard>
            <FormGroup>
              <DropdownSelect
                options={[
                  "Abusing moderator privilege",
                  "Abusing moderator privilege",
                  "Abusing moderator privilege",
                ]}
                placeholder={"Select"}
                labelText="Violation"
                value={violationValue}
                onSelect={(item) => {
                  setViolationValue(item);
                }}
              />
            </FormGroup>
            <FormGroup>
              <FormLabel>Notes</FormLabel>
              <TextArea placeholder="" />
            </FormGroup>
            <FormGroup>
              <FormLabel>Duration</FormLabel>
              <DateRangePicker />
            </FormGroup>
          </SuspendCard>
          <PrimaryButton onClick={handleSuspendAdmin}>Update</PrimaryButton>
          {/* </Container> */}
        </Modal>
      )}
      {suspendUser && (
        <SuccessMessage
          isOpen={suspendUser}
          setIsOpen={setSuspendUser}
          heading="Success"
          message="You have successfully modified the admin user"
          email="Kennis Maduka"
        />
      )}
    </>
  );
};

export default ModifyRole;




const SuspendCard = styled.div``;

const DatePicker = styled.input`
  border: 1px solid #b1afaf;
  margin-bottom: 28px;
  width: 97%;
  border: 1px solid #e0e0e0;
  padding: 8px 12px;
  color: #828282;
  font-weight: 500;
  font-size: 14px;
  font-weight: 500;
  line-height: 24px;
  border-radius: 10px;
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

const TextArea = styled.textarea`
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  padding: 2px;
  width: 100%;
  outline: none;
  height: 100px;
`;

const FormGroup = styled.div`
  margin: 12px 0;
`;

const FormInputContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  border: 1px solid #e0e0e0;
  border-radius: 10px;

  input {
    outline: none;
    width: 100%;
    flex: 1;
    border: none;
  }
`;
