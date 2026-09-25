import { Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import Image from "next/image";
import React, { useState } from "react";
import styled from "styled-components";
import { FaCopy } from "react-icons/fa";
import trov from "@/assets/images/TROVTokenicon.svg";
import SuccessMessage from "@/components/SuccessMessage";

interface AssetModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const AddTokenModal: React.FC<AssetModalProps> = ({ isOpen, onClose }) => {
  const [success, setSuccess] = useState(false);

  return (
    <Modal title="" isOpen={isOpen} onClose={onClose}>
      <Container>
        <Text>
          Do you wish to add this asset <StyledSpan>[TOM]</StyledSpan> to the
          other assets curation? This will make this asset appear on the list of
          assets on Trovo App and will enable users to start transacting with
          it.
        </Text>
      </Container>

      <CenteredContainer>
        <TokenName>TROV TOKEN</TokenName>
        <Image src={trov} alt="TROV Token" width={70} height={70} />
      </CenteredContainer>

      <TokenInfo>
        <TokenLabel>Issuer Public Key</TokenLabel>
        <TokenValue>
          GAXMB...A2CFJ{" "}
          <CopyIcon>
            <FaCopy />
          </CopyIcon>
        </TokenValue>
      </TokenInfo>

      <PrimaryButton
        onClick={() => setSuccess(true)}
        buttonStyle={{ width: "100%" }}
      >
        Add Asset
      </PrimaryButton>

      {success && (
        <SuccessMessage
          isOpen={success}
          setIsOpen={setSuccess}
          heading="Success!"
          message="You have successfully added TROV token to the curation of other assets on the Trovo App"
          email=""
        />
      )}
    </Modal>
  );
};

export default AddTokenModal;

const Container = styled.div`
  background-color: #f2f6f9;
  padding: 20px 32px;
  border-radius: 20px;
`;

const Text = styled.p`
  color: #00225a;
  font-weight: 400;
  font-size: 14px;
  line-height: 17.07px;
  text-align: center;
`;

const StyledSpan = styled.span`
  color: #00225a;
  font-weight: 600;
  font-size: 14px;
  line-height: 17.07px;
  text-align: center;
`;

const TokenName = styled.h3`
  font-weight: 600;
  font-size: 24px;
  color: #00225a;
`;

const TokenInfo = styled.div`
  text-align: center;
  margin-bottom: 8px;
`;

const TokenLabel = styled.h4`
  font-weight: 600;
  font-size: 16px;
  color: #00225a;
  text-align: center;
  line-height: 20px;
  letter-spacing: -0.5%;
`;

const TokenValue = styled.p`
  color: #00225a;
  font-weight: 400;
  font-size: 14px;
  text-align: center;
  line-height: 20px;
  letter-spacing: -0.5%;
`;

const CopyIcon = styled.span`
  margin-left: 5px;
  cursor: pointer;
  color: #00225a;
`;
const CenteredContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  gap: 8px;
  margin: 10px 0;
`;
