"use client";
import React, { useMemo, useState } from "react";
import styled from "styled-components";

import { StatusBadge, showErrorToast, showSuccessToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import {
  useGetOrgComplianceQuery,
  useSubmitOrgComplianceItemMutation,
  useUploadComplianceDocumentMutation,
  useDownloadComplianceDocumentMutation,
} from "@/redux/api/orgCompliance";
import {
  ComplianceInputType,
  ComplianceRequirementInstance,
  ComplianceRequirementStatus,
} from "@/redux/api/compliance/interface";

const READ_ONLY_STATUSES: string[] = [
  ComplianceRequirementStatus.Approved,
  ComplianceRequirementStatus.Waived,
];

const RequirementCard = ({ item }: { item: ComplianceRequirementInstance }) => {
  const [submitItem, { isLoading: isSubmitting }] =
    useSubmitOrgComplianceItemMutation();
  const [uploadDocument, { isLoading: isUploading }] =
    useUploadComplianceDocumentMutation();
  const [downloadDocument, { isLoading: isDownloading }] =
    useDownloadComplianceDocumentMutation();

  const [textValue, setTextValue] = useState(
    (item.submission_data?.text as string) || "",
  );
  const [structuredValue, setStructuredValue] = useState(
    item.submission_data ? JSON.stringify(item.submission_data, null, 2) : "",
  );

  const isReadOnly = READ_ONLY_STATUSES.includes(item.status);
  const isRejected = item.status === ComplianceRequirementStatus.Rejected;

  const handleFileUpload = async (
    e: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      const formData = new FormData();
      formData.append("documentFile", file);
      formData.append("title", item.requirement);

      const uploadResult = await uploadDocument(formData).unwrap();
      const document = uploadResult?.data;

      await submitItem({
        itemId: item.id,
        payload: {
          submission_data: {
            document_id: document?.id,
            original_filename: document?.original_filename || file.name,
            mime_type: document?.mime_type || file.type,
          },
        },
      }).unwrap();

      showSuccessToast("Document submitted successfully.");
    } catch (err: any) {
      showErrorToast(
        err?.data?.message || err?.error || "Failed to submit document.",
      );
    }
  };

  const handleDownload = async () => {
    const documentId = item.submission_data?.document_id as string | undefined;
    if (!documentId) return;
    try {
      const blob = await downloadDocument(documentId).unwrap();
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download =
        (item.submission_data?.original_filename as string) || "document";
      document.body.appendChild(link);
      try {
        link.click();
      } finally {
        link.remove();
        window.setTimeout(() => URL.revokeObjectURL(url), 1000);
      }
    } catch {
      showErrorToast("Unable to download the document. Please try again.");
    }
  };

  const handleTextSubmit = async () => {
    try {
      await submitItem({
        itemId: item.id,
        payload: { submission_data: { text: textValue.trim() } },
      }).unwrap();
      showSuccessToast("Submitted successfully.");
    } catch (err: any) {
      showErrorToast(
        err?.data?.message || err?.error || "Failed to submit.",
      );
    }
  };

  const handleStructuredSubmit = async () => {
    try {
      let parsed: Record<string, unknown>;
      try {
        parsed = structuredValue.trim()
          ? JSON.parse(structuredValue)
          : {};
      } catch {
        showErrorToast("Structured data must be valid JSON.");
        return;
      }
      await submitItem({
        itemId: item.id,
        payload: { submission_data: parsed },
      }).unwrap();
      showSuccessToast("Submitted successfully.");
    } catch (err: any) {
      showErrorToast(
        err?.data?.message || err?.error || "Failed to submit.",
      );
    }
  };

  return (
    <Card>
      <CardHeader>
        <div>
          <CardTitle>{item.requirement}</CardTitle>
          <CardCategory>{item.category}</CardCategory>
        </div>
        <StatusBadge status={item.status} />
      </CardHeader>

      {isRejected && item.rejection_reason && (
        <RejectionNotice>
          Rejected: {item.rejection_reason}
        </RejectionNotice>
      )}

      {!!item.submission_data?.document_id && (
        <DownloadLink
          type="button"
          onClick={handleDownload}
          disabled={isDownloading}
        >
          {isDownloading
            ? "Downloading..."
            : `Download ${(item.submission_data?.original_filename as string) || "document"}`}
        </DownloadLink>
      )}

      {!isReadOnly && (
        <SubmissionControl>
          {item.input_type === ComplianceInputType.DocumentUpload && (
            <>
              <FileInput
                type="file"
                onChange={handleFileUpload}
                disabled={isUploading || isSubmitting}
              />
              {(isUploading || isSubmitting) && (
                <HelperText>Uploading...</HelperText>
              )}
            </>
          )}

          {item.input_type === ComplianceInputType.Text && (
            <>
              <TextArea
                rows={3}
                value={textValue}
                onChange={(e) => setTextValue(e.target.value)}
                placeholder="Enter the requested information"
              />
              <PrimaryButton
                onClick={handleTextSubmit}
                disabled={isSubmitting || !textValue.trim()}
                buttonStyle={{ margin: "8px 0 0 0", width: "160px", height: "40px" }}
              >
                {isRejected ? "Resubmit" : "Submit"}
              </PrimaryButton>
            </>
          )}

          {item.input_type === ComplianceInputType.StructuredForm && (
            <>
              {/* Phase 1 placeholder: no per-item structured schema is defined
                  on ComplianceTemplateItem yet, so structured_form items are
                  captured as free-form JSON pending a real schema field. */}
              <TextArea
                rows={4}
                value={structuredValue}
                onChange={(e) => setStructuredValue(e.target.value)}
                placeholder='{"field": "value"}'
              />
              <PrimaryButton
                onClick={handleStructuredSubmit}
                disabled={isSubmitting}
                buttonStyle={{ margin: "8px 0 0 0", width: "160px", height: "40px" }}
              >
                {isRejected ? "Resubmit" : "Submit"}
              </PrimaryButton>
            </>
          )}
        </SubmissionControl>
      )}
    </Card>
  );
};

const ComplianceSubmissions = () => {
  const { data, isLoading, isError, refetch } = useGetOrgComplianceQuery({
    page: 1,
    limit: 100,
  });

  const records = useMemo(
    () => (Array.isArray(data?.data?.records) ? data.data.records : []),
    [data],
  );

  if (isError) {
    return (
      <Container>
        <Title>Compliance Submissions</Title>
        <StateMessage>
          Unable to load your compliance requirements.
          <RetryButton onClick={() => refetch()}>Try again</RetryButton>
        </StateMessage>
      </Container>
    );
  }

  if (isLoading) {
    return (
      <Container>
        <Title>Compliance Submissions</Title>
        <StateMessage>Loading...</StateMessage>
      </Container>
    );
  }

  return (
    <Container>
      <Title>Compliance Submissions</Title>
      {records.length === 0 ? (
        <StateMessage>No compliance requirements assigned yet.</StateMessage>
      ) : (
        <CardList>
          {records.map((item) => (
            <RequirementCard key={item.id} item={item} />
          ))}
        </CardList>
      )}
    </Container>
  );
};

export default ComplianceSubmissions;

const Container = styled.div`
  background: #ffffff;
  border-radius: 24px;
  padding: 32px;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
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

const CardList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

const Card = styled.div`
  background: #fcfcfc;
  border: 1px solid #f2f6f9;
  border-radius: 12px;
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
`;

const CardHeader = styled.div`
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
`;

const CardTitle = styled.div`
  font-size: 16px;
  font-weight: 600;
  color: #00225a;
`;

const CardCategory = styled.div`
  font-size: 13px;
  color: #828282;
  margin-top: 2px;
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

const RejectionNotice = styled.div`
  background: #fdebea;
  color: #eb5757;
  border-radius: 8px;
  padding: 10px 14px;
  font-size: 13px;
`;

const SubmissionControl = styled.div`
  border-top: 1px solid #f2f2f2;
  padding-top: 12px;
  display: flex;
  flex-direction: column;
`;

const FileInput = styled.input`
  font-size: 13px;
`;

const HelperText = styled.span`
  font-size: 12px;
  color: #828282;
  margin-top: 6px;
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
