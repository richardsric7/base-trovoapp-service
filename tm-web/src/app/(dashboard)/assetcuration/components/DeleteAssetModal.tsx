import { Modal } from "@/components";
import Image from "next/image";
import deleted from "@/assets/images/deleted.svg";

import styled from "styled-components";
interface DeleteAssetProps {
  openModal: boolean;
  setOpenModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const DeleteAssetModal: React.FC<DeleteAssetProps> = ({
  openModal,
  setOpenModal,
}) => {
  const handleSubmit = () => {
    setOpenModal(false);
  };
  return (
    <Modal title="" isOpen={openModal} onClose={() => setOpenModal(false)}>
      <ModalContent>
        <Image src={deleted} alt="delete-icon" width={70} height={70} />
        <Title>Delete Asset Token</Title>
        <Text>
          Are you sure you want to delete <StyledSpan>AE1</StyledSpan> token
          from the curated asset tokens list? This would also be removed from
          the Trovo App
        </Text>
        <DeleteButton onClick={handleSubmit}>Yes, Delete</DeleteButton>
      </ModalContent>
    </Modal>
  );
};

export default DeleteAssetModal;
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
`;
