import React from "react";
import styled from "styled-components";
import { Modal } from "@/components";
import Image from "next/image";
import successIcon from "@/assets/images/success.png";
import { useRouter } from "next/navigation";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";

interface EditCurationProps {
  openModal: boolean;
  setOpenModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const CompleteEdit: React.FC<EditCurationProps> = ({
  openModal,
  setOpenModal,
}) => {
  const router = useRouter();

  return (
    <Modal
      title="Curation Complete"
      isOpen={openModal}
      onClose={() => setOpenModal(false)}
    >
      <Container>
        <Text>Curated assets have been successfully updated.</Text>
        <Image src={successIcon} alt="success-icon" />

        <ButtonContainer>
          <PrimaryButton
            onClick={() => router.push("/assetcuration")}
          >
            Go Back
          </PrimaryButton>
          <SecondaryButton onClick={() => setOpenModal(false)}>
            Continue Editing
          </SecondaryButton>
        </ButtonContainer>
      </Container>
    </Modal>
  );
};

export default CompleteEdit;

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin: 0 auto;
  align-items: center;
  justify-content: center;
  text-align: center;
`;

const Text = styled.h1`
  font-size: 18px;
  font-weight: 400;
  line-height: 20px;
  color: #00225a;
  letter-spacing: 0.1px;
`;

const ButtonContainer = styled.div`
  display: flex;
  width: 70%;
  margin: 0 auto;
  gap: 10px;
`;
