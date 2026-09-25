"use client";

import { Modal } from "antd";
import styled from "styled-components";
import { BsTrash } from "react-icons/bs";
import { FiX } from "react-icons/fi";
import { BsCheckCircle } from "react-icons/bs";
import { IoIosCloseCircleOutline } from "react-icons/io";

export const MilestoneCompletionModal = ({ open, onClose }: any) => {
  return (
    <Modal
      open={open}
      onCancel={onClose}
      footer={null}
      centered
      width={520}
      closeIcon={<FiX size={20} />}
      styles={{
        body: {
          height: "80vh", // 👈 viewport height
          overflowY: "auto", // 👈 enable scroll
        },
      }}
    >
      <Wrapper>
        <Title>Milestone Completion</Title>

        {/* Milestone */}
        <Section>
          <Label>Milestone</Label>

          <MilestoneCard>
            <MilestoneTitle>Permits</MilestoneTitle>
            <MilestoneDate>September, 2026</MilestoneDate>

            <Status>
              <Dot />
              Pending
            </Status>
          </MilestoneCard>
        </Section>

        {/* Amount */}
        <Section>
          <Label>Request Funds</Label>

          <InputWrapper>
            <Input placeholder="0.00" />
            <Currency>NGN</Currency>
          </InputWrapper>
        </Section>

        {/* Images */}
        <Section>
          <Label>Supporting Images</Label>

          <UploadBox>
            <span>
              <UploadText>Upload</UploadText> or drag and drop files here
            </span>
          </UploadBox>

          <FileItem>
            <FileLeft>
              <FileIcon />
              <span>Foundation.jpg</span>
            </FileLeft>
            <FileSize>102KB</FileSize>
            <Trash>
              <BsTrash />
            </Trash>
          </FileItem>
        </Section>

        {/* Documents */}
        <Section>
          <Label>Supporting Documents</Label>

          <UploadBox>
            <span>
              <UploadText>Upload</UploadText> or drag and drop files here
            </span>
          </UploadBox>

          <FileItem>
            <FileLeft>
              <PdfIcon />
              <span>Atlantis plan document</span>
            </FileLeft>

            <ProgressWrapper>
              <ProgressBar />
              <span>99KB of 110KB</span>
            </ProgressWrapper>

            <CloseBtn>
              <FiX />
            </CloseBtn>
          </FileItem>

          <FileItem>
            <FileLeft>
              <DocIcon />
              <span>Atlantis permit</span>
            </FileLeft>
            <FileSize>102KB</FileSize>
            <Trash>
              <BsTrash />
            </Trash>
          </FileItem>
        </Section>

        {/* BUTTON */}
        <SubmitBtn>
          <BsCheckCircle size={16} />
          Complete Milestone
        </SubmitBtn>
      </Wrapper>
    </Modal>
  );
};
const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 18px;
`;

const Title = styled.h2`
  font-size: 18px;
  font-weight: 600;
  color: #00225a;
`;

const Section = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.p`
  font-size: 13px;
  color: #00225a;
  font-weight: 500;
`;

const MilestoneCard = styled.div`
  background: #f2f6f9;
  border-radius: 12px;
  padding: 16px;
`;

const MilestoneTitle = styled.p`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
`;

const MilestoneDate = styled.p`
  font-size: 12px;
  color: #9ca3af;
  margin-top: 4px;
`;

const Status = styled.div`
  margin-top: 8px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #ffffff;
  color: #191919;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
`;

const Dot = styled.div`
  width: 6px;
  height: 6px;
  background: #f59e0b;
  border-radius: 50%;
`;

const InputWrapper = styled.div`
  display: flex;
  align-items: center;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 8px;
`;

const Input = styled.input`
  flex: 1;
  border: none;
  outline: none;
  font-family: inherit;
`;

const Currency = styled.div`
  font-size: 14px;
  color: #00225a;
  font-weight: 500;
`;

const UploadBox = styled.div`
  border: 1px dashed #007cdf;
  border-radius: 10px;
  padding: 16px;
  background-color: #f2f6f9;
  text-align: center;
  color: #6b7280;
  font-size: 13px;
`;

const UploadText = styled.span`
  color: #007cdf;
  font-weight: 500;
`;

const FileItem = styled.div`
  background: #f2f6f9;
  border-radius: 10px;
  padding: 10px 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const FileLeft = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const FileIcon = styled.div`
  width: 24px;
  height: 24px;
  background: #dbeafe;
  border-radius: 6px;
`;

const PdfIcon = styled.div`
  width: 24px;
  height: 24px;
  background: #fee2e2;
  border-radius: 6px;
`;

const DocIcon = styled.div`
  width: 24px;
  height: 24px;
  background: #e0e7ff;
  border-radius: 6px;
`;

const FileSize = styled.span`
  font-size: 12px;
  color: #9ca3af;
`;

const Trash = styled.div`
  cursor: pointer;
  color: #9ca3af;
`;

const CloseBtn = styled.div`
  cursor: pointer;
`;

const ProgressWrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;

const ProgressBar = styled.div`
  width: 120px;
  height: 4px;
  background: #e5e7eb;
  border-radius: 4px;
  position: relative;

  &::after {
    content: "";
    position: absolute;
    width: 70%;
    height: 100%;
    background: #007cdf;
    border-radius: 4px;
  }
`;

const SubmitBtn = styled.button`
  margin-top: 10px;
  background: #007cdf;
  color: white;
  border: none;
  border-radius: 10px;
  padding: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  cursor: pointer;
  font-family: inherit;
`;
