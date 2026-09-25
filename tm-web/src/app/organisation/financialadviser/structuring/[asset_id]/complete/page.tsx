"use client";

import React, { FormEvent, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import styled from "styled-components";
import { FaArrowLeft } from "react-icons/fa6";
import { BsCheckCircle } from "react-icons/bs";
import { showErrorToast, showSuccessToast } from "@/components";
import { useCompleteAssetStructuringMutation } from "@/redux/api/financialAdviser";

type FormShape = {
  title: string;
  document_id: string;
  file_url: string;
  notes: string;
};

const initialForm: FormShape = {
  title: "",
  document_id: "",
  file_url: "",
  notes: "",
};

const getErrorMessage = (error: any, fallback: string) =>
  error?.data?.message ?? error?.error ?? fallback;

const CompleteStructuringPage = () => {
  const params = useParams();
  const router = useRouter();
  const assetId = params?.asset_id as string;

  const [form, setForm] = useState<FormShape>(initialForm);
  const [completeStructuring, { isLoading }] =
    useCompleteAssetStructuringMutation();

  const isValidUrl = (() => {
    try {
      const url = new URL(form.file_url);
      return url.protocol === "http:" || url.protocol === "https:";
    } catch {
      return false;
    }
  })();

  const canSubmit = Boolean(
    form.title.trim() &&
      form.document_id.trim() &&
      form.notes.trim() &&
      isValidUrl,
  );

  const setField = <K extends keyof FormShape>(key: K, value: FormShape[K]) =>
    setForm((prev) => ({ ...prev, [key]: value }));

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    if (!canSubmit || !assetId) return;

    try {
      await completeStructuring({
        assetId,
        title: form.title.trim(),
        document_id: form.document_id.trim(),
        file_url: form.file_url.trim(),
        notes: form.notes.trim(),
      }).unwrap();
      showSuccessToast("Structuring completed successfully");
      router.push("/organisation/financialadviser");
    } catch (error) {
      showErrorToast(getErrorMessage(error, "Unable to complete structuring"));
    }
  };

  return (
    <Page>
      <BackButton type="button" onClick={() => router.back()}>
        <FaArrowLeft />
      </BackButton>

      <Card>
        <Title>Complete Asset Structuring</Title>
        <Subtitle>
          Submit the financial structuring document to mark this asset as
          structured.
        </Subtitle>

        <Form onSubmit={handleSubmit}>
          <Field>
            <Label htmlFor="structuring-title">Title *</Label>
            <TextInput
              id="structuring-title"
              placeholder="Enter document title"
              value={form.title}
              onChange={(event) => setField("title", event.target.value)}
            />
          </Field>

          <Field>
            <Label htmlFor="structuring-document-id">Document ID *</Label>
            <TextInput
              id="structuring-document-id"
              placeholder="Enter document ID"
              value={form.document_id}
              onChange={(event) =>
                setField("document_id", event.target.value)
              }
            />
          </Field>

          <Field>
            <Label htmlFor="structuring-file-url">File URL *</Label>
            <TextInput
              id="structuring-file-url"
              type="url"
              placeholder="https://example.com/structuring-doc.pdf"
              value={form.file_url}
              onChange={(event) => setField("file_url", event.target.value)}
            />
            {form.file_url && !isValidUrl && (
              <FieldError>Enter a valid HTTP or HTTPS URL.</FieldError>
            )}
          </Field>

          <Field>
            <Label htmlFor="structuring-notes">Notes *</Label>
            <TextAreaInput
              id="structuring-notes"
              rows={4}
              placeholder="Add structuring notes"
              value={form.notes}
              onChange={(event) => setField("notes", event.target.value)}
            />
          </Field>

          <SubmitButton type="submit" disabled={!canSubmit || isLoading}>
            <BsCheckCircle size={16} />
            {isLoading ? "Submitting..." : "Complete Structuring"}
          </SubmitButton>
        </Form>
      </Card>
    </Page>
  );
};

export default CompleteStructuringPage;

const Page = styled.main`
  padding: 24px;
  min-height: 100vh;
`;

const BackButton = styled.button`
  border: 0;
  background: transparent;
  color: #00225a;
  cursor: pointer;
  margin-bottom: 16px;
  font-size: 18px;
`;

const Card = styled.section`
  background: #fff;
  border-radius: 24px;
  padding: 32px;
  max-width: 560px;
  margin: 0 auto;
`;

const Title = styled.h1`
  margin: 0;
  color: #00225a;
  font-size: 22px;
  text-align: center;
`;

const Subtitle = styled.p`
  margin: 8px 0 0;
  color: #828282;
  font-size: 14px;
  text-align: center;
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 18px;
  margin-top: 26px;
`;

const Field = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  font-size: 14px;
  font-weight: 500;
  color: #191919;
`;

const fieldCss = `height: 46px; width: 100%; border: 1px solid #e0e0e0; border-radius: 8px; padding: 0 12px; outline: none; font: inherit; color: #00225a; background: #fff;`;

const TextInput = styled.input`
  ${fieldCss}
  &:focus {
    border-color: #007cdf;
  }
`;

const TextAreaInput = styled.textarea`
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 12px;
  outline: none;
  font: inherit;
  color: #00225a;
  resize: vertical;

  &:focus {
    border-color: #007cdf;
  }
`;

const FieldError = styled.span`
  font-size: 12px;
  color: #b42318;
`;

const SubmitButton = styled.button`
  margin-top: 10px;
  background: #007cdf;
  color: white;
  border: none;
  border-radius: 10px;
  padding: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  cursor: pointer;
  font-family: inherit;
  font-weight: 600;
  font-size: 14px;

  &:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }
`;
