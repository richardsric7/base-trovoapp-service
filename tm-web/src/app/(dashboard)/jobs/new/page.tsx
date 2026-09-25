"use client";

import React, { useMemo, useState } from "react";
import styled from "styled-components";
import { useRouter } from "next/navigation";
import { FaArrowLeftLong, FaPlus } from "react-icons/fa6";
import { Button, Input, Radio, Space, Typography, message } from "antd";

import dynamic from "next/dynamic";

import {
  usePostJobMutation,
  usePublishJobMutation,
  useUpdateJobMutation,
} from "@/redux/api/jobs/api";
import { FiEye, FiFileText } from "react-icons/fi";
import { PiAirplane } from "react-icons/pi";
import { showErrorToast, showSuccessToast } from "@/components";
import {
  JobPostPayload,
  LocationType,
  SectionKey,
  WorkType,
} from "@/redux/api/jobs";
import SectionEditor from "../_components/SectionEditor";

const { TextArea } = Input;
const { Title } = Typography;

const DEFAULT_SECTIONS: Record<SectionKey, string[]> = {
  "About The Role": [],
  "What You'll Own": [],
  "What We're Looking For": [],
  "What Success Looks Like": [],
};

const APPLY_LOCATION = "Nigeria, Remote";

function normalizeWorkTypeUIToApi(value: string): WorkType {
  const v = value.toLowerCase().replace(/[\s-]+/g, "");
  if (v === "fulltime") return "fulltime";
  if (v === "parttime") return "parttime";
  return "contract";
}

function normalizeLocationUIToApi(value: string): LocationType {
  const v = value.toLowerCase().replace(/[\s-]+/g, "");
  if (v === "remote") return "remote";
  if (v === "onsite") return "onsite";
  return "hybrid";
}

const JobPreview = dynamic(() => import("../_components/JobPreview"), {
  ssr: false,
  loading: () => null,
});

const NewJobPostPage = () => {
  const router = useRouter();
  const [isPreviewOpen, setIsPreviewOpen] = useState(false);

  const [postJob, { isLoading: isPostLoading }] = usePostJobMutation();
  const [updateJob, { isLoading: isUpdateLoading }] = useUpdateJobMutation();
  const [publishJob, { isLoading: isPublishLoading }] = usePublishJobMutation();

  const isPublishing = isPublishLoading;

  const [heading, setHeading] = useState("");
  const [companyOverview, setCompanyOverview] = useState(
    "Trovotech Ltd is a SEC - regulated digital asset and tokenization platform building secure, compliant, and accessible investment infrastructure for real-world assets. Our mission is to build compliant, intuitive, and scalable technology that enables the tokenization, management and trade of real-world assets, expanding access to capital and investment opportunities within emerging and global markets. We operate at the intersection of finance, regulation, and technology, and we are intentionally building a high-performance, values-driven team to support our next phase of growth. Trovotech is an equal opportunity employer and we value diversity. We do not discriminate on the basis of religion, colour, gender, race, nationality, age, marital status, veteran status, or disability status.",
  );
  const [companyLocation, setCompanyLocation] = useState("");
  const [yearsOfExperience, setYearsOfExperience] = useState("");
  const [workTypeUI, setWorkTypeUI] = useState("Full-time");
  const [locationUI, setLocationUI] = useState("Remote");
  const [applyTo, setApplyTo] = useState(
    "Send CV + Short Note to hiring@trovotech.io",
  );
  const [subject, setSubject] = useState("");

  const [createdJobId, setCreatedJobId] = useState<string | null>(null);

  const [sections, setSections] =
    useState<Record<SectionKey, string[]>>(DEFAULT_SECTIONS);

  const payload: JobPostPayload = useMemo(() => {
    return {
      heading: heading.trim(),
      company_overview: companyOverview.trim(),
      location: normalizeLocationUIToApi(locationUI),
      years_of_experience: yearsOfExperience.trim(),
      work_type: normalizeWorkTypeUIToApi(workTypeUI),
      sections: {
        "About The Role": sections["About The Role"],
        "What We're Looking For": sections["What We're Looking For"],

        "What You'll Own": sections["What You'll Own"],
        "What Success Looks Like": sections["What Success Looks Like"],
      },
      application: {
        apply: applyTo.trim() || "Send CV + Short Note to hiring@trovotech.io",
        location: companyLocation.trim() || APPLY_LOCATION,
        subject: subject.trim() || `${heading || "Position"} - Trovotech`,
      },
    };
  }, [
    heading,
    companyOverview,
    companyLocation,
    applyTo,
    subject,
    yearsOfExperience,
    workTypeUI,
    locationUI,
    sections,
  ]);

  const handleSaveDraft = async () => {
    const minimalOk =
      payload.heading ||
      payload.company_overview ||
      payload.work_type ||
      payload.location ||
      Object.values(payload.sections).some((arr) => arr.length > 0);
    if (!minimalOk) {
      message.warning(
        "Add at least a title or overview before saving a draft.",
      );
      return;
    }

    try {
      if (createdJobId) {
        await updateJob({ id: createdJobId, ...payload }).unwrap();
        showSuccessToast("Draft updated successfully.");
      } else {
        const response = await postJob(payload).unwrap();
        setCreatedJobId(response?.data?.id);
        showSuccessToast("Draft saved successfully.");
      }
    } catch (e: any) {
      showErrorToast(e?.data?.message || "Failed to save draft.");
    }
  };

  const handlePublishJob = async () => {
    if (!createdJobId) {
      showErrorToast("Please save the job as a draft first.");
      return;
    }

    try {
      await publishJob(createdJobId).unwrap();
      showSuccessToast("Job published successfully");
    } catch {
      showErrorToast("Failed to publish job");
    }
  };

  return (
    <Container>
      <Header>
        <BackButton onClick={() => router.back()}>
          <FaArrowLeftLong />
        </BackButton>

        <ActionButtons>
          <SecondaryBtn
            onClick={handleSaveDraft}
            loading={isPostLoading || isUpdateLoading}
          >
            <FiFileText />
            {isPostLoading || isUpdateLoading ? "Saving..." : "Save as Draft"}
          </SecondaryBtn>
          <SecondaryBtn onClick={() => setIsPreviewOpen(true)}>
            <FiEye /> Preview
          </SecondaryBtn>
          <PrimaryBtn onClick={handlePublishJob}>
            <PiAirplane />
            {isPublishing ? "Publishing..." : "Publish Job"}
          </PrimaryBtn>
        </ActionButtons>
      </Header>
      <Content>
        <FormTitle>New Job Post</FormTitle>

        <Block>
          <Label>Position Title</Label>
          <StyledInput
            placeholder="e.g Technical Lead"
            value={heading}
            onChange={(e) => setHeading(e.target.value)}
          />
        </Block>

        <Block>
          <Label>Company Overview</Label>
          <StyledTextArea
            rows={4}
            value={companyOverview}
            onChange={(e) => setCompanyOverview(e.target.value)}
            placeholder="Describe the company’s mission"
          />
        </Block>

        <Block>
          <Label>Years of Experience</Label>
          <StyledInput
            value={yearsOfExperience}
            onChange={(e) => setYearsOfExperience(e.target.value)}
            placeholder="e.g 5"
          />
        </Block>

        <Block>
          <Label>Company Location</Label>
          <StyledInput
            value={companyLocation}
            onChange={(e) => setCompanyLocation(e.target.value)}
            placeholder="e.g Lagos, Nigeria"
          />
        </Block>

        <Block>
          <Label>Employment Type</Label>
          <Radio.Group
            value={workTypeUI}
            onChange={(e) => setWorkTypeUI(e.target.value)}
          >
            <Space direction="vertical">
              <StyledRadio value="Full-time">Full-time</StyledRadio>
              <StyledRadio value="Part-time">Part-time</StyledRadio>
              <StyledRadio value="Contract">Contract</StyledRadio>
            </Space>
          </Radio.Group>
        </Block>

        <Block>
          <Label>Select Work Model</Label>
          <Radio.Group
            value={locationUI}
            onChange={(e) => setLocationUI(e.target.value)}
          >
            <Space direction="vertical">
              <StyledRadio value="Onsite">Onsite</StyledRadio>
              <StyledRadio value="Remote">Remote</StyledRadio>
              <StyledRadio value="Hybrid">Hybrid</StyledRadio>
            </Space>
          </Radio.Group>
        </Block>

        <SectionEditor
          title="About The Role"
          items={sections["About The Role"]}
          onChange={(items) =>
            setSections((p) => ({ ...p, "About The Role": items }))
          }
        />

        <SectionEditor
          title="What You'll Own"
          items={sections["What You'll Own"]}
          onChange={(items) =>
            setSections((p) => ({ ...p, "What You'll Own": items }))
          }
        />

        <SectionEditor
          title="What We're Looking For"
          items={sections["What We're Looking For"]}
          onChange={(items) =>
            setSections((p) => ({ ...p, "What We're Looking For": items }))
          }
        />

        <SectionEditor
          title="What Success Looks Like"
          items={sections["What Success Looks Like"]}
          onChange={(items) =>
            setSections((p) => ({ ...p, "What Success Looks Like": items }))
          }
        />

        <Block>
          <Label>Application</Label>

          <StyledApplicationInput
            addonBefore="Apply to"
            placeholder="Enter email address"
            value={applyTo}
            onChange={(e) => setApplyTo(e.target.value)}
            style={{ marginBottom: 14 }}
          />

          <StyledApplicationInput
            addonBefore="Subject"
            value={subject || `${heading || "Position"} - Trovotech`}
            onChange={(e) => setSubject(e.target.value)}
          />
        </Block>
      </Content>
      <JobPreview
        payload={payload}
        open={isPreviewOpen}
        onCancel={() => setIsPreviewOpen(false)}
        onPublish={handlePublishJob}
        isPublishing={isPublishing}
      />
    </Container>
  );
};

export default NewJobPostPage;

const Container = styled.div`
  padding: 32px;
  background: #fff;
  border-radius: 24px;
`;

const Content = styled.div`
  max-width: 520px;
  margin: 0 auto;
`;
const Header = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 24px;
`;

const BackButton = styled.button`
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
`;

const ActionButtons = styled.div`
  display: flex;
  gap: 10px;
`;

const PrimaryBtn = styled(Button)`
  background: #007cdf;
  color: #fff;
  height: 48px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 12px;
`;

const SecondaryBtn = styled(Button)`
  border: 1px solid #007cdf;
  color: #007cdf;
  height: 48px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 12px;
`;

const FormTitle = styled.h2`
  margin-bottom: 24px;
  font-size: 24px;
  font-weight: 700;
  color: #00225a;
`;

const Block = styled.div`
  margin-bottom: 20px;
`;

const Label = styled.label`
  display: block;
  font-weight: 500;
  margin-bottom: 8px;
  color: #00225a;
  font-size: 14px;
`;

const StyledInput = styled(Input)`
  height: 42px;
`;

const StyledTextArea = styled(TextArea)``;

const StyledRadio = styled(Radio)``;
const StyledApplicationInput = styled(Input)`
  .ant-input-group-addon {
    background: #FAFAFC
    border: 1px solid #ececec;
    border-right: none;
    padding: 0 16px;
    font-size: 14px;
    font-weight: 400;
    color: #00225A;
    min-width: 96px;
  }

  .ant-input {
    height: 48px;
    border: 1px solid #ececec;
    border-radius: 0 12px 12px 0;
    font-size: 14px;
    padding: 0 14px;

    &::placeholder {
      color: #b0b0b0;
    }
  }

  .ant-input-group-addon:first-child {
    border-radius: 12px 0 0 12px;
  }

  &:hover .ant-input,
  &:hover .ant-input-group-addon {
    border-color: #d1d5db;
  }

  &.ant-input-group-wrapper-focused .ant-input,
  &.ant-input-group-wrapper-focused .ant-input-group-addon {
    border-color: #2563eb;
    box-shadow: none;
  }
`;
const BlockTitle = styled.h3`
  font-size: 18px;
  font-weight: 600;
`;
