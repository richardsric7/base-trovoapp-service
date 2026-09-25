"use client";
import styled from "styled-components";
import { Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import Image from "next/image";
import { TokenizationRecord } from "@/redux/api/assettokenization";

interface DisapproveProps {
  disapproval: boolean;
  setDisapproval: React.Dispatch<React.SetStateAction<boolean>>;
  onSubmit: () => void;
  asset?: TokenizationRecord;
  reasonInput: string;
  setReasonInput: React.Dispatch<React.SetStateAction<string>>;
}
const DisapproveModal: React.FC<DisapproveProps> = ({
  disapproval,
  setDisapproval,
  onSubmit,
  asset,
  reasonInput,
  setReasonInput,
}) => {
  const handleSubmit = () => {
    onSubmit();
    setDisapproval(false);
  };

  return (
    <>
      <Modal
        title="Disapprove Tokenization"
        isOpen={disapproval}
        onClose={() => setDisapproval(false)}
      >
        <Text>
          You are about to disapprove the tokenization of the followin asset
        </Text>
        <Content>
          {asset?.assetLogo ? (
            <Image
              src={asset?.assetLogo}
              alt="Asset Logo"
              width={40}
              height={40}
              style={{ borderRadius: "50%" }}
            />
          ) : (
            <Avatar />
          )}

          <div>
            <AssetTitle>{asset?.assetCode}</AssetTitle>
            <AssetSubTitle>{asset?.assetName}</AssetSubTitle>
          </div>
        </Content>
        <div>
          <Label>Reason for disapproval</Label>
          <TextArea
            type="text"
            value={reasonInput}
            onChange={(e) => setReasonInput(e.target.value)}
          />
        </div>

        <PrimaryButton onClick={handleSubmit}>
          Complete Disapproval
        </PrimaryButton>
      </Modal>
    </>
  );
};

export default DisapproveModal;

const Label = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  letter-spacing: 0.10000000149011612px;

  color: #00225a;
`;

const AssetTitle = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;

  margin: 0;
  color: #00225a;
`;
const Text = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;

  margin: 0;
  color: #00225a;
`;
const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;
const Content = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 12px;
`;
const AssetSubTitle = styled.p`
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;
const TextArea = styled.input`
  width: 95%;
  padding: 12px;
  margin: 8px auto 0;
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 400;
  color: #000;
  background-color: transparent;
  resize: vertical;
  min-height: 120px;
  display: block;
  font-family:inherit;
    outline: none

  &:focus {
    outline: none;
    border-color: #007cdf;
    background-color: #ffffff;
  }
`;
const FileInputWrapper = styled.label`
  display: flex;
  align-items: center;
  gap: 2px;
  margin-top: 6px;
  padding: 8px 12px;
  cursor: pointer;
  color: #007cdf;
  font-size: 14px;
  font-weight: 600;
`;

const HiddenFileInput = styled.input`
  display: none;
`;
const FileList = styled.div`
  margin-top: 10px;
  padding: 0;
  list-style: none;
`;

const FileItem = styled.div`
  background: #f2f6f9;
  padding: 8px;
  border-radius: 10px;
  margin-top: 5px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
`;

const FileName = styled.p`
  font-size: 14px;
  gap: 8px;
  font-weight: 500;
  color: #00225a;
`;

const FileSize = styled.p`
  font-size: 12px;
  gap: 8px;
  font-weight: 400;
  color: #00000099;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
`;
