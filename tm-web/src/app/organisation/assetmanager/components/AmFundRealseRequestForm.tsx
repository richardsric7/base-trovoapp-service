import React from "react";
import styled from "styled-components";
import { showErrorToast } from "@/components";
import type { IStakeholderDocument } from "@/redux/api/sharedstakeholders/interface";

export type AssetOption = {
  id: string;
  code: string;
  name: string;
  currency: string;
};
export type ReleaseForm = {
  assetId: string;
  amount: string;
  currency: string;
  purpose: string;
  receivingAccountName: string;
  receivingAccountNumber: string;
  receivingBank: string;
};

interface AmFundRealseRequestFormProps {
  form: ReleaseForm;
  setForm: React.Dispatch<React.SetStateAction<ReleaseForm>>;
  assets: AssetOption[];
  assetsLoading: boolean;
  documents: IStakeholderDocument[];
  documentsLoading: boolean;
  documentFiles: File[];
  setDocumentFiles: React.Dispatch<React.SetStateAction<File[]>>;
  selectedDocumentIds: string[];
  setSelectedDocumentIds: React.Dispatch<React.SetStateAction<string[]>>;
  isCreating: boolean;
  isUploading: boolean;
  canCreate: boolean;
  setCreateOpen: React.Dispatch<React.SetStateAction<boolean>>;
  handleCreate: (event: React.FormEvent<HTMLFormElement>) => Promise<void>;
}

const AmFundRealseRequestForm = ({
  form,
  setForm,
  assets,
  assetsLoading,
  documents,
  documentsLoading,
  documentFiles,
  setDocumentFiles,
  selectedDocumentIds,
  setSelectedDocumentIds,
  isCreating,
  isUploading,
  canCreate,
  setCreateOpen,
  handleCreate,
}: AmFundRealseRequestFormProps) => {
  return (
    <Overlay onMouseDown={() => !isCreating && setCreateOpen(false)}>
      <ModalCard onMouseDown={(event) => event.stopPropagation()}>
        <ModalHeader>
          <div>
            <ModalTitle>Create Fund Release</ModalTitle>
            <ModalSubtitle>
              Request a release for one of your managed assets.
            </ModalSubtitle>
          </div>
          <CloseButton onClick={() => setCreateOpen(false)}>×</CloseButton>
        </ModalHeader>
        <Form onSubmit={handleCreate}>
          <Field>
            <FormLabel htmlFor="release-asset">
              Asset <RequiredMark>*</RequiredMark>
            </FormLabel>
            <InputSelect
              id="release-asset"
              value={form.assetId}
              disabled={assetsLoading}
              onChange={(event) => {
                const asset = assets.find(
                  (item) => item.id === event.target.value,
                );
                setForm((current) => ({
                  ...current,
                  assetId: event.target.value,
                  currency: asset?.currency || current.currency,
                }));
                setSelectedDocumentIds([]);
                setDocumentFiles([]);
              }}
            >
              <option value="">
                {assetsLoading ? "Loading assets..." : "Select an asset"}
              </option>
              {assets.map((asset) => (
                <option key={asset.id} value={asset.id}>
                  {asset.code}
                  {asset.name ? ` — ${asset.name}` : ""}
                </option>
              ))}
            </InputSelect>
          </Field>
          <Field>
            <FormLabel htmlFor="release-amount">
              Amount <RequiredMark>*</RequiredMark>
            </FormLabel>
            <TextInput
              id="release-amount"
              type="number"
              min="0.01"
              step="0.01"
              value={form.amount}
              onChange={(event) =>
                setForm((current) => ({
                  ...current,
                  amount: event.target.value,
                }))
              }
            />
          </Field>
          <Field>
            <FormLabel htmlFor="release-currency">
              Currency <RequiredMark>*</RequiredMark>
            </FormLabel>
            <TextInput
              id="release-currency"
              placeholder="e.g. NGN"
              value={form.currency}
              onChange={(event) =>
                setForm((current) => ({
                  ...current,
                  currency: event.target.value,
                }))
              }
            />
          </Field>
          <Field>
            <FormLabel htmlFor="release-purpose">
              Purpose <RequiredMark>*</RequiredMark>
            </FormLabel>
            <TextArea
              id="release-purpose"
              rows={3}
              value={form.purpose}
              onChange={(event) =>
                setForm((current) => ({
                  ...current,
                  purpose: event.target.value,
                }))
              }
            />
          </Field>
          <Field>
            <FormLabel htmlFor="receiving-bank">Receiving bank</FormLabel>
            <TextInput
              id="receiving-bank"
              value={form.receivingBank}
              onChange={(event) =>
                setForm((current) => ({
                  ...current,
                  receivingBank: event.target.value,
                }))
              }
            />
          </Field>
          <Field>
            <FormLabel htmlFor="receiving-account-name">
              Receiving account name
            </FormLabel>
            <TextInput
              id="receiving-account-name"
              value={form.receivingAccountName}
              onChange={(event) =>
                setForm((current) => ({
                  ...current,
                  receivingAccountName: event.target.value,
                }))
              }
            />
          </Field>
          <Field>
            <FormLabel htmlFor="receiving-account-number">
              Receiving account number
            </FormLabel>
            <TextInput
              id="receiving-account-number"
              inputMode="numeric"
              value={form.receivingAccountNumber}
              onChange={(event) =>
                setForm((current) => ({
                  ...current,
                  receivingAccountNumber: event.target.value,
                }))
              }
            />
          </Field>
          <Field>
            <FormLabel htmlFor="release-document">
              Supporting documents
            </FormLabel>
            <FileInput
              id="release-document"
              type="file"
              multiple
              accept="application/pdf,image/jpeg,image/png"
              disabled={isUploading}
              onChange={(event) => {
                const input = event.currentTarget;
                const files = Array.from(input.files ?? []);
                const unsupported = files.find(
                  (file) =>
                    !["application/pdf", "image/jpeg", "image/png"].includes(
                      file.type,
                    ),
                );
                if (unsupported) {
                  showErrorToast(
                    `${unsupported.name}: only PDF, JPEG and PNG documents are supported`,
                  );
                  input.value = "";
                  return;
                }
                const oversized = files.find(
                  (file) => file.size > 10 * 1024 * 1024,
                );
                if (oversized) {
                  showErrorToast(
                    `${oversized.name}: document must be 10 MB or smaller`,
                  );
                  input.value = "";
                  return;
                }
                setDocumentFiles((current) => {
                  const known = new Set(
                    current.map(
                      (file) =>
                        `${file.name}-${file.size}-${file.lastModified}`,
                    ),
                  );
                  return [
                    ...current,
                    ...files.filter(
                      (file) =>
                        !known.has(
                          `${file.name}-${file.size}-${file.lastModified}`,
                        ),
                    ),
                  ];
                });
                input.value = "";
              }}
            />
            <FileHint>
              Add one or more PDF, JPEG or PNG files. You can click the picker
              again to add more. Maximum size is 10 MB per file.
            </FileHint>
            {documentFiles.length > 0 && (
              <SelectedFiles>
                {documentFiles.map((file) => (
                  <SelectedFile
                    key={`${file.name}-${file.size}-${file.lastModified}`}
                  >
                    <span>{file.name}</span>
                    <RemoveFile
                      type="button"
                      onClick={() =>
                        setDocumentFiles((current) =>
                          current.filter((item) => item !== file),
                        )
                      }
                    >
                      Remove
                    </RemoveFile>
                  </SelectedFile>
                ))}
              </SelectedFiles>
            )}
          </Field>
          {form.assetId && (
            <Field>
              <FormLabel>Previously uploaded supporting documents</FormLabel>
              {documentsLoading ? (
                <DocumentNote>Loading documents…</DocumentNote>
              ) : documents.length === 0 ? (
                <DocumentNote>
                  No supporting documents found for this asset.
                </DocumentNote>
              ) : (
                <DocumentList>
                  {documents.map((document) => (
                    <DocumentOption key={document.id}>
                      <input
                        type="checkbox"
                        checked={selectedDocumentIds.includes(document.id)}
                        onChange={() =>
                          setSelectedDocumentIds((current) =>
                            current.includes(document.id)
                              ? current.filter((id) => id !== document.id)
                              : [...current, document.id],
                          )
                        }
                      />
                      <span>
                        {document.title || document.original_filename}
                      </span>
                    </DocumentOption>
                  ))}
                </DocumentList>
              )}
            </Field>
          )}
          <ModalActions>
            <CancelButton type="button" onClick={() => setCreateOpen(false)}>
              Cancel
            </CancelButton>
            <CreateButton
              type="submit"
              disabled={!canCreate || isCreating || isUploading}
            >
              {isUploading
                ? "Uploading..."
                : isCreating
                  ? "Creating..."
                  : "Create Request"}
            </CreateButton>
          </ModalActions>
        </Form>
      </ModalCard>
    </Overlay>
  );
};

export default AmFundRealseRequestForm;

const Overlay = styled.div`
  position: fixed;
  inset: 0;
  z-index: 10000;
  display: grid;
  place-items: center;
  padding: 24px;
  background: rgba(0, 34, 90, 0.42);
`;
const ModalCard = styled.div`
  width: min(560px, 100%);
  max-height: 90vh;
  overflow-y: auto;
  border-radius: 20px;
  padding: 28px;
  background: #fff;
`;
const ModalHeader = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 24px;
`;
const ModalTitle = styled.h2`
  margin: 0;
  color: #00225a;
  font-size: 22px;
`;
const ModalSubtitle = styled.p`
  margin: 6px 0 0;
  color: #828282;
  font-size: 14px;
`;
const CloseButton = styled.button`
  border: 0;
  background: transparent;
  color: #667085;
  font-size: 26px;
  cursor: pointer;
`;
const Form = styled.form`
  display: grid;
  gap: 18px;
`;
const Field = styled.div`
  display: grid;
  gap: 8px;
`;
const FormLabel = styled.label`
  color: #00225a;
  font-size: 14px;
  font-weight: 500;
`;
const RequiredMark = styled.span`
  color: #d92d20;
`;
const inputStyles = `width: 100%; border: 1px solid #d0d5dd; border-radius: 8px; padding: 11px 12px; background: #fff; color: #00225a; font: inherit; outline: none; &:focus { border-color: #007cdf; }`;
const TextInput = styled.input`
  ${inputStyles}
`;
const TextArea = styled.textarea`
  ${inputStyles};
  resize: vertical;
`;
const InputSelect = styled.select`
  ${inputStyles}
`;
const FileInput = styled.input`
  ${inputStyles};
  padding: 8px 10px;
`;
const FileHint = styled.span`
  font-size: 12px;
  color: #828282;
`;
const SelectedFiles = styled.div`
  display: grid;
  gap: 4px;
  font-size: 12px;
  color: #00225a;
`;
const SelectedFile = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
`;
const RemoveFile = styled.button`
  border: 0;
  background: transparent;
  color: #b42318;
  font: inherit;
  cursor: pointer;
`;
const DocumentNote = styled.span`
  font-size: 13px;
  color: #828282;
`;
const DocumentList = styled.div`
  display: grid;
  gap: 8px;
  max-height: 128px;
  overflow-y: auto;
  padding: 10px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
`;
const DocumentOption = styled.label`
  display: flex;
  align-items: center;
  gap: 8px;
  color: #00225a;
  font-size: 13px;
  cursor: pointer;
`;
const ModalActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 8px;
`;
const CancelButton = styled.button`
  border: 1px solid #d0d5dd;
  border-radius: 8px;
  min-height: 40px;
  padding: 0 18px;
  background: #fff;
  color: #344054;
  font: inherit;
  cursor: pointer;
`;

const CreateButton = styled.button`
  background: #007cdf;
  border: none;
  padding: 8px 18px;
  border-radius: 8px;
  color: #fff;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  &:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }
`;
