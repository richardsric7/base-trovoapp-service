"use client";

import { useEffect, useMemo, useState } from "react";
import styled from "styled-components";
import * as yup from "yup";
import { useFormik } from "formik";

import { DropdownSelect, Modal, StatusBadge, showErrorToast, showSuccessToast } from "@/components";
import CustomTable from "@/components/CustomTable";
import CustomFilter from "@/components/CustomFilter";
import Pagination from "@/components/CustomPagination";
import PrimaryButton from "@/components/PrimaryButton";
import { useGetOrganizationsListQuery } from "@/redux/api/organizations";
import {
  useGetComplianceRequirementsQuery,
  useGetComplianceRequirementDetailQuery,
  useReviewComplianceRequirementMutation,
  useStartComplianceRequirementReviewMutation,
  useDownloadComplianceRequirementDocumentMutation,
} from "@/redux/api/admin";
import {
  ComplianceRequirementInstance,
  ComplianceRequirementStatus,
} from "@/redux/api/compliance/interface";
import AdHocComplianceForm from "./AdHocComplianceForm";

const STATUS_OPTIONS = Object.values(ComplianceRequirementStatus);

const humanize = (value: string) =>
  value
    .toLowerCase()
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");

const formatDate = (value?: string) =>
  value ? new Date(value).toLocaleDateString() : "N/A";

interface ReviewFormValues {
  decision: "approved" | "rejected" | "";
  rejection_reason: string;
}

const reviewValidationSchema = yup.object({
  decision: yup.string().oneOf(["approved", "rejected"]).required(),
  rejection_reason: yup.string().when("decision", {
    is: "rejected",
    then: (schema) => schema.trim().required("Rejection reason is required"),
    otherwise: (schema) => schema.notRequired(),
  }),
});

const ReviewPanel = ({
  requirementId,
  onDone,
}: {
  requirementId: string;
  onDone: () => void;
}) => {
  const { data, isLoading } = useGetComplianceRequirementDetailQuery(
    requirementId,
  );
  const [reviewRequirement, { isLoading: isReviewing }] =
    useReviewComplianceRequirementMutation();
  const [startReview] = useStartComplianceRequirementReviewMutation();
  const [downloadDocument, { isLoading: isDownloading }] =
    useDownloadComplianceRequirementDocumentMutation();

  const requirement = data?.data;

  // GET is read-only server-side; a submitted item only enters review once
  // the admin actually opens it, via this explicit transition.
  useEffect(() => {
    if (requirement?.status === ComplianceRequirementStatus.Submitted) {
      startReview(requirementId);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [requirement?.status, requirementId]);

  const handleDownload = async (documentId: string, filename?: string) => {
    try {
      const blob = await downloadDocument({
        requirementId,
        documentId,
      }).unwrap();
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = filename || "document";
      document.body.appendChild(link);
      try {
        link.click();
      } finally {
        link.remove();
        window.setTimeout(() => URL.revokeObjectURL(url), 1000);
      }
    } catch (err: any) {
      showErrorToast(
        err?.data?.message || err?.error || "Unable to download document.",
      );
    }
  };

  const { values, errors, touched, handleChange, setFieldValue, handleSubmit } =
    useFormik<ReviewFormValues>({
      initialValues: { decision: "", rejection_reason: "" },
      validationSchema: reviewValidationSchema,
      onSubmit: async (formValues) => {
        try {
          await reviewRequirement({
            id: requirementId,
            payload: {
              decision: formValues.decision as "approved" | "rejected",
              ...(formValues.decision === "rejected"
                ? { rejection_reason: formValues.rejection_reason.trim() }
                : {}),
            },
          }).unwrap();
          showSuccessToast(
            formValues.decision === "approved"
              ? "Requirement approved."
              : "Requirement rejected.",
          );
          onDone();
        } catch (err: any) {
          showErrorToast(
            err?.data?.message || err?.error || "Failed to submit review.",
          );
        }
      },
    });

  if (isLoading) return <StateMessage>Loading requirement...</StateMessage>;
  if (!requirement) return <StateMessage>Requirement not found.</StateMessage>;

  const submission = requirement.submission_data;

  return (
    <DetailContainer>
      <DetailRow>
        <DetailLabel>Category</DetailLabel>
        <DetailValue>{requirement.category}</DetailValue>
      </DetailRow>
      <DetailRow>
        <DetailLabel>Requirement</DetailLabel>
        <DetailValue>{requirement.requirement}</DetailValue>
      </DetailRow>
      <DetailRow>
        <DetailLabel>Status</DetailLabel>
        <DetailValue>
          <StatusBadge status={requirement.status} />
        </DetailValue>
      </DetailRow>
      <DetailRow>
        <DetailLabel>Organisation</DetailLabel>
        <DetailValue>{requirement.org_id}</DetailValue>
      </DetailRow>

      <SubmissionBox>
        <DetailLabel>Submission</DetailLabel>
        {submission ? (
          requirement.input_type === "document_upload" &&
          submission.document_id ? (
            <DownloadLink
              type="button"
              onClick={() =>
                handleDownload(
                  String(submission.document_id),
                  submission.original_filename
                    ? String(submission.original_filename)
                    : undefined,
                )
              }
              disabled={isDownloading}
            >
              {isDownloading ? "Downloading..." : "Download submitted document"}
            </DownloadLink>
          ) : requirement.input_type === "text" && submission.text ? (
            <p>{String(submission.text)}</p>
          ) : (
            <pre>{JSON.stringify(submission, null, 2)}</pre>
          )
        ) : (
          <StateMessage>No submission yet.</StateMessage>
        )}
      </SubmissionBox>

      {requirement.rejection_reason && (
        <DetailRow>
          <DetailLabel>Previous rejection reason</DetailLabel>
          <DetailValue>{requirement.rejection_reason}</DetailValue>
        </DetailRow>
      )}

      <ReviewForm onSubmit={handleSubmit}>
        <DecisionRow>
          <DecisionButton
            type="button"
            $active={values.decision === "approved"}
            $variant="approve"
            onClick={() => setFieldValue("decision", "approved")}
          >
            Approve
          </DecisionButton>
          <DecisionButton
            type="button"
            $active={values.decision === "rejected"}
            $variant="reject"
            onClick={() => setFieldValue("decision", "rejected")}
          >
            Reject
          </DecisionButton>
        </DecisionRow>

        {values.decision === "rejected" && (
          <FormGroup>
            <Label>Rejection Reason *</Label>
            <TextArea
              name="rejection_reason"
              value={values.rejection_reason}
              onChange={handleChange}
              rows={3}
              placeholder="Explain what needs to change so the organisation can resubmit"
            />
            {touched.rejection_reason && errors.rejection_reason && (
              <ErrorMessage>{errors.rejection_reason}</ErrorMessage>
            )}
          </FormGroup>
        )}

        <PrimaryButton
          disabled={!values.decision || isReviewing}
          buttonStyle={{ margin: "16px 0 0 0", width: "200px", height: "44px" }}
        >
          {isReviewing ? "Submitting..." : "Submit Review"}
        </PrimaryButton>
      </ReviewForm>
    </DetailContainer>
  );
};

const ComplianceRequirementsQueue = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [filterOpen, setFilterOpen] = useState(false);
  const [statusFilter, setStatusFilter] = useState("");
  const [categoryFilter, setCategoryFilter] = useState("");
  const [orgFilter, setOrgFilter] = useState("");
  const [selectedRequirementId, setSelectedRequirementId] = useState<
    string | null
  >(null);
  const [isAdHocModalOpen, setIsAdHocModalOpen] = useState(false);

  const { data: orgsData } = useGetOrganizationsListQuery({
    page: 1,
    pageSize: 100,
  });
  const organisations = orgsData?.organizations ?? [];

  const { data, isLoading, isFetching, isError, refetch } =
    useGetComplianceRequirementsQuery({
      page: currentPage,
      limit: pageSize,
      ...(statusFilter ? { status: statusFilter } : {}),
      ...(categoryFilter ? { category: categoryFilter } : {}),
      ...(orgFilter ? { org_id: orgFilter } : {}),
    });

  const records = useMemo(
    () => (Array.isArray(data?.data?.records) ? data.data.records : []),
    [data],
  );
  const meta = data?.data?.meta;

  const orgName = (orgId: string) =>
    organisations.find((org) => org.id === orgId)?.name || orgId;

  const columns = [
    {
      title: "Organisation",
      dataIndex: "org_id",
      render: (orgId: string) => orgName(orgId),
    },
    { title: "Category", dataIndex: "category" },
    { title: "Requirement", dataIndex: "requirement" },
    {
      title: "Status",
      dataIndex: "status",
      render: (status: string) => <StatusBadge status={status} />,
    },
    {
      title: "Submitted At",
      dataIndex: "updated_at",
      render: (value: string) => formatDate(value),
    },
  ];

  if (isError) {
    return (
      <Container>
        <Title>Compliance Requirements</Title>
        <StateMessage>
          Unable to load compliance requirements.
          <RetryButton onClick={() => refetch()}>Try again</RetryButton>
        </StateMessage>
      </Container>
    );
  }

  return (
    <Container>
      <Header>
        <div>
          <Title>Compliance Requirements</Title>
          <SubTitle>
            Review submitted documents and information, and approve or reject
            each item.
          </SubTitle>
        </div>
        <HeaderActions>
          <CustomFilter
            open={filterOpen}
            onClose={setFilterOpen}
            position={{ top: "60px", right: "0px" }}
          >
            <FilterGroup>
              <DropdownSelect
                labelText="Status"
                placeholder="Any status"
                value={statusFilter ? humanize(statusFilter) : ""}
                options={STATUS_OPTIONS.map(humanize)}
                onSelect={(label) => {
                  const found = STATUS_OPTIONS.find(
                    (opt) => humanize(opt) === label,
                  );
                  setStatusFilter(found || "");
                }}
              />
              <DropdownSelect
                labelText="Organisation"
                placeholder="Any organisation"
                value={orgName(orgFilter) === orgFilter ? "" : orgName(orgFilter)}
                options={organisations.map((org) => org.name)}
                onSelect={(name) => {
                  const found = organisations.find((org) => org.name === name);
                  setOrgFilter(found?.id || "");
                }}
              />
              <FilterInputGroup>
                <Label>Category</Label>
                <Input
                  value={categoryFilter}
                  onChange={(e) => setCategoryFilter(e.target.value)}
                  placeholder="e.g. KYC"
                />
              </FilterInputGroup>
            </FilterGroup>
          </CustomFilter>

          <PrimaryButton
            onClick={() => setIsAdHocModalOpen(true)}
            buttonStyle={{ margin: 0, width: "220px", height: "44px" }}
          >
            Assign Ad-hoc Requirement
          </PrimaryButton>
        </HeaderActions>
      </Header>

      <CustomTable
        columns={columns}
        dataSource={records}
        isLoading={isLoading || isFetching}
        totalItems={meta?.total ?? 0}
        pageSize={pageSize}
        onRowClick={(record: ComplianceRequirementInstance) =>
          setSelectedRequirementId(record.id)
        }
      />

      {(meta?.total ?? 0) > 0 && (
        <Pagination
          totalCount={meta?.total ?? 0}
          currentPage={currentPage}
          pageSize={pageSize}
          onPageChange={setCurrentPage}
          onPageSizeChange={setPageSize}
        />
      )}

      <Modal
        title="Requirement Details"
        isOpen={!!selectedRequirementId}
        onClose={() => setSelectedRequirementId(null)}
      >
        {selectedRequirementId && (
          <ReviewPanel
            requirementId={selectedRequirementId}
            onDone={() => setSelectedRequirementId(null)}
          />
        )}
      </Modal>

      <Modal
        title="Assign Ad-hoc Requirement"
        isOpen={isAdHocModalOpen}
        onClose={() => setIsAdHocModalOpen(false)}
      >
        <AdHocComplianceForm onCreated={() => setIsAdHocModalOpen(false)} />
      </Modal>
    </Container>
  );
};

export default ComplianceRequirementsQueue;

const Container = styled.div`
  background: #ffffff;
  border-radius: 24px;
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 24px;
`;

const Header = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
  position: relative;
`;

const HeaderActions = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  position: relative;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin: 0 0 4px 0;
`;

const SubTitle = styled.p`
  font-size: 14px;
  color: #828282;
  margin: 0;
`;

const StateMessage = styled.div`
  color: #00225a;
`;

const RetryButton = styled.button`
  display: block;
  margin-top: 12px;
  padding: 8px 16px;
  border: 0;
  border-radius: 8px;
  background: #007cdf;
  color: #fff;
  cursor: pointer;
`;

const FilterGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 220px;
`;

const FilterInputGroup = styled.div`
  display: flex;
  flex-direction: column;
`;

const Label = styled.label`
  font-weight: 500;
  font-size: 14px;
  margin-bottom: 8px;
  color: #00225a;
`;

const Input = styled.input`
  height: 40px;
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  padding: 0 12px;
  font-size: 14px;
  color: #00225a;
  font-family: inherit;
  outline: none;
`;

const DetailContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

const DetailRow = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;

const DetailLabel = styled.span`
  font-size: 12px;
  font-weight: 600;
  color: #828282;
  text-transform: uppercase;
`;

const DetailValue = styled.span`
  font-size: 14px;
  color: #00225a;
`;

const DownloadLink = styled.button`
  align-self: flex-start;
  border: none;
  background: transparent;
  color: #007cdf;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  padding: 0;
  text-decoration: underline;

  &:disabled {
    opacity: 0.6;
    cursor: default;
  }
`;

const SubmissionBox = styled.div`
  background: #f2f6f9;
  border-radius: 12px;
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;

  pre {
    white-space: pre-wrap;
    word-break: break-word;
    margin: 0;
    font-size: 13px;
  }
`;

const ReviewForm = styled.form`
  border-top: 1px solid #f2f2f2;
  padding-top: 16px;
`;

const DecisionRow = styled.div`
  display: flex;
  gap: 12px;
`;

const DecisionButton = styled.button<{
  $active: boolean;
  $variant: "approve" | "reject";
}>`
  flex: 1;
  height: 44px;
  border-radius: 10px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  border: 1px solid
    ${({ $variant }) => ($variant === "approve" ? "#1D8A46" : "#EB5757")};
  background: ${({ $active, $variant }) =>
    $active
      ? $variant === "approve"
        ? "#1D8A46"
        : "#EB5757"
      : "transparent"};
  color: ${({ $active, $variant }) =>
    $active ? "#ffffff" : $variant === "approve" ? "#1D8A46" : "#EB5757"};
`;

const FormGroup = styled.div`
  margin-top: 16px;
  display: flex;
  flex-direction: column;
`;

const TextArea = styled.textarea`
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  padding: 12px;
  font-size: 14px;
  color: #00225a;
  font-family: inherit;
  resize: vertical;
  outline: none;
`;

const ErrorMessage = styled.p`
  font-size: 12px;
  color: #be3800;
  margin: 6px 0 0 0;
`;
