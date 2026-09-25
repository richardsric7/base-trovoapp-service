"use client";

import React, { useState } from "react";
import styled from "styled-components";
import { ArrowLeft, Briefcase, Location, Calendar, Clock } from "iconsax-react";
import { PiAirplane } from "react-icons/pi";
import { useParams, useRouter } from "next/navigation";
import SecondaryButton from "@/components/SecondaryButton";
import PrimaryButton from "@/components/PrimaryButton";
import { BiEdit } from "react-icons/bi";
import { TbTrash } from "react-icons/tb";
import DeleteJobModal from "../_components/DeleteJobModal";
import {
  useDeleteJobMutation,
  useGetJobByIdQuery,
  usePublishJobMutation,
  useJobStatusMutation,
} from "@/redux/api/jobs";
import { showErrorToast, showSuccessToast } from "@/components";
import JobStatusPill from "../_components/JobStatusPill";

const JobDetailsPage = () => {
  const router = useRouter();
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const { id } = useParams<{ id: string }>();
  const { data, isLoading, isError } = useGetJobByIdQuery(id);
  const [deleteJob, { isLoading: isDeleting }] = useDeleteJobMutation();
  const [publishJob, { isLoading: isPublishing }] = usePublishJobMutation();
  const [updateJobStatus, { isLoading: isStatusUpdating }] =
    useJobStatusMutation();

  if (isLoading) {
    return <Centered>Loading job details…</Centered>;
  }

  if (isError || !data?.data) {
    return <Centered>Failed to load job details.</Centered>;
  }

  const handlePublishJob = async (jobId: string) => {
    try {
      await publishJob(jobId).unwrap();
      showSuccessToast("Job published successfully");
      router.push("/jobs");
    } catch (error) {
      showErrorToast("Failed to publish job");
    }
  };

  const handlePublishStatus = async () => {
    try {
      await updateJobStatus({ id, status: "published" }).unwrap();
      showSuccessToast("Job published successfully");
    } catch (error) {
      showErrorToast("Failed to publish job");
    }
  };

  const handleDisableStatus = async () => {
    try {
      await updateJobStatus({ id, status: "disabled" }).unwrap();
      showSuccessToast("Job disabled successfully");
    } catch (error) {
      showErrorToast("Failed to disable job");
    }
  };

  const job = data.data;
  const capitalizeFirstLetter = (str: string): string =>
    str.charAt(0).toUpperCase() + str.slice(1);

  return (
    <Container>
      {/* HEADER */}
      <Header>
        <BackButton onClick={() => router.push("/jobs")}>
          <ArrowLeft size={18} />
        </BackButton>

        <Actions>
          <SecondaryButton
            onClick={() => router.push(`/jobs/${id}/edit`)}
            buttonStyle={{ display: "flex", alignItems: "center", gap: "4px" }}
          >
            <BiEdit size={16} />
            Edit
          </SecondaryButton>
          <SecondaryButton
            onClick={() => setIsDeleteModalOpen(true)}
            buttonStyle={{
              border: "1px solid #BE3800",
              color: "#BE3800",
              display: "flex",
              alignItems: "center",
              gap: "4px",
            }}
          >
            <TbTrash size={16} />
            Delete
          </SecondaryButton>

          {job.status === "draft" && (
            <PrimaryButton
              onClick={() => handlePublishJob(job.id)}
              buttonStyle={{ display: "flex", gap: "4px" }}
            >
              <PiAirplane size={16} />
              Publish Job
            </PrimaryButton>
          )}
        </Actions>
      </Header>

      {/* JOB CARD */}
      <Card>
        <Section>
          <Title>{job.heading}</Title>

          {/* META */}
          <MetaRow>
            <MetaItem>
              <Clock size={16} color="#BDBDBD" />
              {capitalizeFirstLetter(job.work_type)}
            </MetaItem>
            <MetaItem>
              <Location size={16} color="#BDBDBD" />
              {capitalizeFirstLetter(job.location)}
            </MetaItem>
            <MetaItem>
              <Briefcase size={16} color="#BDBDBD" />

              {job.years_of_experience}
            </MetaItem>
          </MetaRow>
        </Section>
        <Section>
          {(job.status === "published" || job.status === "disabled") && (
            <JobStatusPill
              status={job.status as "published" | "disabled"}
              loading={isStatusUpdating}
              onPublish={handlePublishStatus}
              onDisable={handleDisableStatus}
            />
          )}
        </Section>

        {/* SECTIONS */}
        <Section>
          <SectionTitle>Company Overview</SectionTitle>
          <CompanyText>{job.company_overview}</CompanyText>
        </Section>
        {job.sections?.["About The Role"] && (
          <Section>
            <SectionTitle>About The Role</SectionTitle>
            <List>
              {job.sections["About The Role"].map((item: any, i: any) => (
                <li key={i}>{item}</li>
              ))}
            </List>
          </Section>
        )}

        {job.sections?.["What You'll Own"] && (
          <Section>
            <SectionTitle>What You&apos;ll Own</SectionTitle>
            <List>
              {job.sections["What You'll Own"].map((item: any, i: any) => (
                <li key={i}>{item}</li>
              ))}
            </List>
          </Section>
        )}

        {job.sections?.["What We're Looking For"] && (
          <Section>
            <SectionTitle>What We&apos;re Looking For</SectionTitle>
            <List>
              {job.sections["What We're Looking For"].map(
                (item: any, i: any) => (
                  <li key={i}>{item}</li>
                ),
              )}
            </List>
          </Section>
        )}

        {job.sections?.["What Success Looks Like"] && (
          <Section>
            <SectionTitle>What Success Looks Like</SectionTitle>
            <List>
              {job.sections["What Success Looks Like"].map(
                (item: any, i: any) => (
                  <li key={i}>{item}</li>
                ),
              )}
            </List>
          </Section>
        )}

        <Section>
          <SectionTitle>Application</SectionTitle>
          <Text>
            {" "}
            Apply To:
            <span> {job.application?.apply} </span>
          </Text>
          <Text>
            Email Subject: <span>{job.application?.subject}</span>
          </Text>
        </Section>
      </Card>

      {isDeleteModalOpen && (
        <DeleteJobModal
          openModal={isDeleteModalOpen}
          setOpenModal={setIsDeleteModalOpen}
          jobTitle={job.heading}
          jobId={job.id}
        />
      )}
    </Container>
  );
};

export default JobDetailsPage;

const Container = styled.div`
  background: #ffffff;
  border-radius: 12px;
  padding: 32px;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
`;

const BackButton = styled.button`
  display: flex;
  align-items: center;
  gap: 8px;
  background: transparent;
  border: none;
  color: #00225a;
  font-weight: 500;
  cursor: pointer;
`;

const Actions = styled.div`
  display: flex;
  gap: 12px;
  width: 48%;
  justify-content: flex-end;
`;

const Card = styled.div`
  width: 80%;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 600;
  color: #00225a;
  margin-bottom: 8px;
`;

const MetaRow = styled.div`
  display: flex;
  flex-direction: column;

  flex-wrap: wrap;
  gap: 16px;

  color: #6b7280;
  font-size: 14px;
`;

const MetaItem = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 400;
  color: #00225a;
`;

const Section = styled.section`
  margin-bottom: 24px;
`;

const SectionTitle = styled.h3`
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #00225a;
`;

const Text = styled.p`
  font-size: 14px;
  line-height: 1.7;
  color: #00225a;
  font-weight: 600;

  span {
    font-weight: 400;
    color: #00225a;
  }
`;

const CompanyText = styled.p`
  font-size: 14px;
  line-height: 1.7;
  color: #00225a;
`;

const List = styled.ul`
  padding-left: 20px;
  color: #00225a;
  font-size: 14px;
  line-height: 1.7;

  li {
    margin-bottom: 6px;
  }
`;

const Centered = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: center;
  margin-top: 20px;
`;
