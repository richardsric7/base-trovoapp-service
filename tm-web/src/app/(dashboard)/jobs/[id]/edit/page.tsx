"use client";

import React, { useEffect, useMemo, useState } from "react";
import styled from "styled-components";
import { useParams, useRouter } from "next/navigation";
import { FaArrowLeftLong, FaPlus } from "react-icons/fa6";
import { Button, Input, Radio, Space } from "antd";

import dynamic from "next/dynamic";

import {
  useGetJobByIdQuery,
  useUpdateJobMutation,
  usePublishJobMutation,
} from "@/redux/api/jobs/api";
import { FiEye, FiFileText, FiTrash2 } from "react-icons/fi";
import { PiAirplane } from "react-icons/pi";
import { showErrorToast, showSuccessToast } from "@/components";
import {
  JobPostPayload,
  LocationType,
  SectionKey,
  WorkType,
} from "@/redux/api/jobs";
import { ArrowLeft } from "iconsax-react";
import SectionEditor from "../../_components/SectionEditor";

const { TextArea } = Input;

const DEFAULT_SECTIONS: Record<SectionKey, string[]> = {
  "What We're Looking For": [],
  "About The Role": [],
  "What You'll Own": [],
  "What Success Looks Like": [],
};

const APPLY_EMAIL = "hiring@trovotech.io";
// const APPLY_LOCATION = "Nigeria, Remote";

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

function denormalizeWorkTypeApiToUI(value: WorkType): string {
  if (value === "fulltime") return "Full-time";
  if (value === "parttime") return "Part-time";
  return "Contract";
}

function denormalizeLocationApiToUI(value: LocationType): string {
  if (value === "remote") return "Remote";
  if (value === "onsite") return "Onsite";
  return "Hybrid";
}

function sanitizeBullet(text: string) {
  return text.replace(/\s+/g, " ").trim();
}

const JobPreview = dynamic(() => import("../../_components/JobPreview"), {
  ssr: false,
  loading: () => null,
});

const EditJobPage = () => {
  const router = useRouter();
  const { id } = useParams<{ id: string }>();
  const [isPreviewOpen, setIsPreviewOpen] = useState(false);

  const { data: jobData, isLoading: isJobLoading } = useGetJobByIdQuery(id);
  const [updateJob, { isLoading: isUpdateLoading }] = useUpdateJobMutation();
  const [publishJob, { isLoading: isPublishLoading }] = usePublishJobMutation();

  const isSaving = isUpdateLoading;
  const isPublishing = isPublishLoading || isUpdateLoading;

  // Form state
  const [heading, setHeading] = useState("");
  const [companyOverview, setCompanyOverview] = useState("");
  const [companyLocation, setCompanyLocation] = useState("");
  const [yearsOfExperience, setYearsOfExperience] = useState("1-3 years");
  const [workTypeUI, setWorkTypeUI] = useState("Full-time");
  const [locationUI, setLocationUI] = useState("Remote");
  const [applyTo, setApplyTo] = useState(APPLY_EMAIL);
  const [subject, setSubject] = useState("");
  const [sections, setSections] =
    useState<Record<SectionKey, string[]>>(DEFAULT_SECTIONS);

  useEffect(() => {
    if (jobData?.data) {
      const job = jobData.data;
      setHeading(job.heading || "");
      setCompanyOverview(job.company_overview || "");
      setCompanyLocation(job.application?.location || "");
      setYearsOfExperience(job.years_of_experience || "1-3 years");
      setWorkTypeUI(denormalizeWorkTypeApiToUI(job.work_type));
      setLocationUI(denormalizeLocationApiToUI(job.location));
      setApplyTo(job.application?.apply || APPLY_EMAIL);
      setSubject(job.application?.subject || "");
      setSections({
        ...DEFAULT_SECTIONS,
        ...job.sections,
      });
    }
  }, [jobData]);

  const payload: JobPostPayload = useMemo(() => {
    const work_type = normalizeWorkTypeUIToApi(workTypeUI);
    const location = normalizeLocationUIToApi(locationUI);

    return {
      heading: heading.trim(),
      company_overview: companyOverview.trim(),
      location,
      years_of_experience: yearsOfExperience.trim(),
      work_type,
      sections: {
        "About The Role": sections["About The Role"]
          .map(sanitizeBullet)
          .filter(Boolean),
        "What You'll Own": sections["What You'll Own"]
          .map(sanitizeBullet)
          .filter(Boolean),
        "What We're Looking For": sections["What We're Looking For"]
          .map(sanitizeBullet)
          .filter(Boolean),
        "What Success Looks Like": sections["What Success Looks Like"]
          .map(sanitizeBullet)
          .filter(Boolean),
      },
      application: {
        apply: applyTo.trim(),
        location: companyLocation.trim() || "",
        subject:
          subject.trim() || `Application for ${heading.trim() || "Role"}`,
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

  const handleUpdate = async () => {
    try {
      await updateJob({ id, ...payload }).unwrap();
      showSuccessToast("Job updated successfully.");
      router.back();
    } catch (e: any) {
      showErrorToast(e?.data?.message || "Failed to update job.");
    }
  };

  if (isJobLoading) {
    return <div>Loading Job Details...</div>;
  }

  return (
    <Container>
      <Header>
        <BackButton onClick={() => router.back()} aria-label="Back">
          <ArrowLeft size={18} />
        </BackButton>

        <ActionButtons>
          <SecondaryBtn onClick={handleUpdate} loading={isSaving}>
            <FiFileText size={16} />
            Save Changes
          </SecondaryBtn>
          <SecondaryBtn onClick={() => setIsPreviewOpen(true)}>
            <FiEye size={16} />
            Preview
          </SecondaryBtn>
        </ActionButtons>
      </Header>

      <ContentArea>
        <FormTitle>Edit Job Post</FormTitle>

        <FormSection>
          <FormGroup>
            <Label>Position Title</Label>
            <StyledInput
              placeholder="e.g Technical lead"
              value={heading}
              onChange={(e) => setHeading(e.target.value)}
            />
          </FormGroup>

          <FormGroup>
            <Label>Company Overview</Label>
            <StyledTextArea
              rows={4}
              placeholder="Describe the company's mission"
              value={companyOverview}
              onChange={(e) => setCompanyOverview(e.target.value)}
            />
          </FormGroup>

          <FormGroup>
            <Label>Years of Experience</Label>
            <StyledInput
              value={yearsOfExperience}
              onChange={(e) => setYearsOfExperience(e.target.value)}
              placeholder="e.g 5"
            />
          </FormGroup>

          <FormGroup>
            <Label>Company Location</Label>
            <StyledInput
              placeholder="e.g Lagos, Nigeria"
              value={companyLocation}
              onChange={(e) => setCompanyLocation(e.target.value)}
            />
          </FormGroup>

          <FormGroup>
            <Label>Employment Type</Label>
            <Radio.Group
              value={workTypeUI}
              onChange={(e) => setWorkTypeUI(e.target.value)}
            >
              <Space direction="vertical">
                {" "}
                <StyledRadio value="Full-time">Full-time</StyledRadio>
                <StyledRadio value="Part-time">Part-time</StyledRadio>
                <StyledRadio value="Contract">Contract</StyledRadio>
              </Space>
            </Radio.Group>
          </FormGroup>

          <FormGroup>
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
          </FormGroup>

          <SectionEditor
            title="What You'll Own"
            items={sections["What You'll Own"]}
            onChange={(items) =>
              setSections((p) => ({ ...p, "What You'll Own": items }))
            }
            placeholder="Enter a point"
          />

          <SectionEditor
            title="What We're Looking For"
            items={sections["What We're Looking For"]}
            onChange={(items) =>
              setSections((p) => ({ ...p, "What We're Looking For": items }))
            }
            placeholder="Enter a point"
          />

          <SectionEditor
            title="What Success Looks Like"
            items={sections["What Success Looks Like"]}
            onChange={(items) =>
              setSections((p) => ({ ...p, "What Success Looks Like": items }))
            }
            placeholder="Enter a point"
          />

          <ApplicationBlock>
            <Label>Application</Label>
            <AppInputGroup>
              <AppInputLabel>Apply to</AppInputLabel>
              <AppInputField
                placeholder="Enter email address"
                value={applyTo}
                onChange={(e) => setApplyTo(e.target.value)}
              />
            </AppInputGroup>
            <AppInputGroup>
              <AppInputLabel>Subject</AppInputLabel>
              <AppInputField
                placeholder="Enter application subject"
                value={subject}
                onChange={(e) => setSubject(e.target.value)}
              />
            </AppInputGroup>
          </ApplicationBlock>
        </FormSection>
      </ContentArea>

      {/* Preview Modal */}
      <JobPreview
        payload={payload}
        open={isPreviewOpen}
        onCancel={() => setIsPreviewOpen(false)}
        showPublishButton={false}
        onPublish={() => {}}
      />
    </Container>
  );
};

export default EditJobPage;

const Container = styled.div`
  padding: 40px;
  background: #ffffff;
  min-height: 100vh;
  border-radius: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const ContentArea = styled.div`
  width: 100%;
  max-width: 800px;

  padding: 40px;
`;

const Header = styled.div`
  width: 100%;
  max-width: 800px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
`;

const BackButton = styled.button`
  background: transparent;
  border: none;
  cursor: pointer;
  color: #00225a;
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s;
  &:hover {
    transform: translateX(-4px);
  }
`;

const ActionButtons = styled.div`
  display: flex;
  gap: 16px;
  align-items: center;
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

const FormSection = styled.div`
  display: flex;
  flex-direction: column;
  gap: 24px;
`;

const FormTitle = styled.h1`
  font-size: 24px;
  font-weight: 700;
  color: #00225a;
  margin: 0 0 32px 0;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  margin-bottom: 4px;
`;

const StyledInput = styled(Input)`
  height: 48px;
  border-radius: 8px;
  border: 1px solid #e4e7ec;
  font-size: 14px;
  color: #00225a;

  &::placeholder {
    color: #98a2b3;
  }

  &:focus,
  &:hover {
    border-color: #007cdf;
    box-shadow: 0 0 0 2px rgba(0, 124, 223, 0.1);
  }
`;

const StyledTextArea = styled(TextArea)`
  border-radius: 8px;
  border: 1px solid #e4e7ec;
  font-size: 14px;
  padding: 12px;

  &:focus,
  &:hover {
    border-color: #007cdf;
    box-shadow: 0 0 0 2px rgba(0, 124, 223, 0.1);
  }
`;

const StyledRadio = styled(Radio)`
  .ant-radio-inner {
    border-color: #d0d5dd;
  }
  .ant-radio-checked .ant-radio-inner {
    border-color: #007cdf;
    background-color: #007cdf;
  }
  span {
    font-size: 14px;
    font-weight: 500;
    color: #344054;
  }
`;

const SectionCard = styled.div`
  display: flex;
  flex-direction: column;
  gap: 12px;
`;

const AddRow = styled.div`
  display: flex;
  gap: 12px;
`;

const AddInput = styled(Input)`
  height: 44px;
  border-radius: 8px;
  border: 1px solid #e4e7ec;
`;

const AddButton = styled(Button)`
  height: 44px;
  border-radius: 8px;
  background-color: #007cdf;
  color: white;
  border: none;
  font-weight: 600;
  &:hover {
    background-color: #006abc !important;
    color: #fff !important;
  }
`;

const ItemsWrap = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 4px;
`;

const ItemRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f2f4f7;
  padding: 12px 16px;
  border-radius: 8px;
`;

const ItemText = styled.span`
  font-size: 14px;
  color: #344054;
`;

const RemoveBtn = styled(Button)`
  color: #be3800;
  border: 1px solid ##be3800;
  background: white;
  height: 48px;
  width: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  &:hover {
    color: #b42318 !important;
    border-color: #fda29b !important;
  }
`;

const ApplicationBlock = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 8px;
`;

const AppInputGroup = styled.div`
  display: flex;
  align-items: center;
  border: 1px solid #e4e7ec;
  border-radius: 8px;
  height: 48px;
  overflow: hidden;
`;

const AppInputLabel = styled.div`
  background: #f9fafb;
  padding: 0 16px;
  height: 100%;
  display: flex;
  align-items: center;
  font-size: 14px;
  font-weight: 500;
  color: #475467;
  border-right: 1px solid #e4e7ec;
  min-width: 100px;
`;

const AppInputField = styled(Input)`
  border: none !important;
  height: 100%;
  &:focus,
  &:hover {
    border: none !important;
    box-shadow: none !important;
  }
`;
