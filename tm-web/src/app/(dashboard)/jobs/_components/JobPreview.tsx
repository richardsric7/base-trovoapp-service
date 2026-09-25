"use client";

import React from "react";
import { Modal } from "antd";
import styled from "styled-components";
import { Timer1, Location, Briefcase } from "iconsax-react";

import PrimaryButton from "@/components/PrimaryButton";
import { PiAirplane } from "react-icons/pi";
import { JobPostPayload, SectionKey } from "@/redux/api/jobs";

interface JobPreviewProps {
  payload: JobPostPayload;
  open: boolean;
  onCancel: () => void;
  isPublishing?: boolean;
  onPublish: () => void;
  showPublishButton?: boolean;
}

const JobPreview: React.FC<JobPreviewProps> = ({
  payload,
  open,
  onCancel,
  isPublishing,
  onPublish,
  showPublishButton = true,
}) => {
  return (
    <StyledModal
      open={open}
      onCancel={onCancel}
      footer={null}
      width={920}
      closable={false}
    >
      <ModalHeader>
        {/* Left: Close */}
        <CloseBtn onClick={onCancel}>×</CloseBtn>

        {/* Right: Publish */}
        {showPublishButton && (
          <PrimaryButton
            onClick={onPublish}
            buttonStyle={{
              display: "flex",
              alignItems: "center",
              gap: 6,
              width: "20%",
              marginRight: "0rem",
            }}
          >
            <PiAirplane size={16} />
            {isPublishing ? "Publishing..." : "     Publish Job"}
          </PrimaryButton>
        )}
      </ModalHeader>

      {/* Content */}
      <PreviewWrap>
        <PreviewHeader>
          <PreviewTitle>{payload.heading || "Job Title"}</PreviewTitle>

          <PreviewMetaList>
            <PreviewMetaItem>
              <Timer1 size={18} color="#667085" />
              <span>
                {payload.work_type
                  ? payload.work_type.charAt(0).toUpperCase() +
                    payload.work_type.slice(1)
                  : "—"}
              </span>
            </PreviewMetaItem>
            <PreviewMetaItem>
              <Location size={18} color="#667085" />
              <span>
                {payload.location
                  ? payload.location.charAt(0).toUpperCase() +
                    payload.location.slice(1)
                  : "—"}
              </span>
            </PreviewMetaItem>
            <PreviewMetaItem>
              <Briefcase size={18} color="#667085" />
              <span>{payload.years_of_experience || "—"}</span>
            </PreviewMetaItem>
          </PreviewMetaList>
        </PreviewHeader>

        <PreviewSection>
          <SectionHeading>Company Overview</SectionHeading>
          <BodyText>{payload.company_overview || "—"}</BodyText>
        </PreviewSection>

        {(
          Object.entries(payload.sections) as Array<[SectionKey, string[]]>
        ).map(([title, items]) => (
          <PreviewSection key={title}>
            <SectionHeading>{title}</SectionHeading>
            {items.length === 0 ? (
              <BodyText>—</BodyText>
            ) : (
              <ul>
                {items.map((it, i) => (
                  <li key={`${title}-${i}`}>{it}</li>
                ))}
              </ul>
            )}
          </PreviewSection>
        ))}

        <PreviewSection>
          <SectionHeading>Application</SectionHeading>
          <ApplyBox>
            <div>
              <b>Apply:</b> {payload.application.apply}
            </div>

            <div>
              <b>Email Subject:</b> {payload.application.subject}
            </div>
          </ApplyBox>
        </PreviewSection>
      </PreviewWrap>
    </StyledModal>
  );
};

export default JobPreview;

const StyledModal = styled(Modal)`
  .ant-modal-content {
    padding: 32px;
    border-radius: 12px;
  }
`;

const ModalHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 28px;
`;

const CloseBtn = styled.button`
  background: none;
  border: none;
  font-size: 22px;
  cursor: pointer;
  color: #667085;
  padding: 0;
  line-height: 1;

  &:hover {
    color: #101828;
  }
`;

const PreviewWrap = styled.div`
  padding-top: 4px;

  ul {
    padding-left: 18px;
    margin: 10px 0 0 0;
    color: #00225a;
  }

  li {
    margin: 6px 0;
    line-height: 1.55;
  }
`;

const PreviewHeader = styled.div`
  margin-bottom: 20px;
`;

const PreviewTitle = styled.h2`
  font-size: 22px;
  font-weight: 800;
  color: #00225a;
  margin: 0 0 12px 0;
`;

const PreviewMetaList = styled.div`
  display: flex;
  flex-direction: column;

  flex-wrap: wrap;
  gap: 16px;

  color: #6b7280;
  font-size: 14px;
`;

const PreviewMetaItem = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 400;
  color: #00225a;
`;

const PreviewSection = styled.div`
  margin-bottom: 18px;
`;

const SectionHeading = styled.h3`
  font-size: 16px;
  font-weight: 600;
  color: #00225a;
  margin: 0 0 8px 0;
`;

const BodyText = styled.div`
  color: #00225a;
  font-size: 14px;
  line-height: 1.65;

  ul {
  }
`;

const ApplyBox = styled.div`
  color: #00225a;
  font-size: 14px;

  div {
    margin: 6px 0;
  }

  b {
    color: #00225a;
  }
`;
