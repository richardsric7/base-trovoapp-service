import { Modal } from "@/components";
import React from "react";
import { BiEdit } from "react-icons/bi";
import { RiDeleteBin6Line } from "react-icons/ri";
import styled from "styled-components";

interface ViewModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const ViewSubSectorModal: React.FC<ViewModalProps> = ({ isOpen, onClose }) => {
  const sectors = ["Industrial buildings", "Rsidentail buildings"];
  return (
    <Modal
      title="Sub-sectors"
      isOpen={isOpen}
      onClose={onClose}
      closeIconPosition="left"
    >
      <FormGroup>
        <Label>Asset Sector</Label>
        <Text>Real Estate</Text>
      </FormGroup>
      <SectorContainer>
        {sectors.map((sec, idx) => (
          <SectorContent>
            <SectorTitle>{sec}</SectorTitle>
            <SectorIconContainer>
              <EditLink>
                <BiEdit color=" #007cdf" size={20} />
              </EditLink>

              <DeleteBtn>
                {" "}
                <RiDeleteBin6Line color="#be3800" size={20} />
              </DeleteBtn>
            </SectorIconContainer>
          </SectorContent>
        ))}
      </SectorContainer>
    </Modal>
  );
};

export default ViewSubSectorModal;

const Label = styled.p`
  color: #828282;
  font-weight: 500;
  font-size: 14px;
  margin-top: 14px;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 10px;
`;

const Text = styled.p`
  font-weight: 600;
  color: #00225a;
  font-size: 14px;
  line-height: 24px;
  letter-spacing: 0.1px;
`;

const SectorContainer = styled.div`
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: 8px;
  width: 100%;
`;
const SectorContent = styled.div`
  background-color: #f2f6f9;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  border-radius: 8px;
  width: 100%;
`;

const SectorTitle = styled.p`
  font-weight: 400;
  font-size: 14px;
  color: #00225a;
  line-height: 24px;
  letter-spacing: 0.1px;
`;

const SectorIconContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const EditLink = styled.div`
  border: 1px solid #007cdf;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
`;

const DeleteBtn = styled.div`
  border: 1px solid #be3800;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
`;
