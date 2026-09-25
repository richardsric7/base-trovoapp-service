"use client";
import { DropdownSelect, Modal, showErrorToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import {
  TokenizationRecord,
  useGetTokenizationParamsQuery,
  useUploadTokenizationDocumentMutation,
} from "@/redux/api/assettokenization";
import React, { useState } from "react";
import styled from "styled-components";

interface AddDocumentProps {
  addDocument: boolean;
  setAddDocument: (value: boolean) => void;
  asset: TokenizationRecord;
}
const AddDocumentModal: React.FC<AddDocumentProps> = ({
  addDocument,
  setAddDocument,
  asset,
}) => {
  const [selectValues, setSelectValues] = useState<any>({});
  const [title, setTitle] = useState("");
  const [file, setFile] = useState<File | null>(null);

  const [uploadDocument, { isLoading }] =
    useUploadTokenizationDocumentMutation();

  const { data: tokenizationParams } = useGetTokenizationParamsQuery();

  const documentTypeOptions =
    tokenizationParams?.tokenizationDocumentTypes.map((type) => ({
      label: type.documentType.replace(/([A-Z])/g, " $1").trim(),
      value: type.documentType,
    })) ?? [];

  const filterCriteria = {
    documentType: documentTypeOptions,
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selected = e.target.files?.[0];
    if (selected) {
      setFile(selected);
      console.log("Selected File", selected); // Should show: name, type, size
    }
  };
  const handleSave = async () => {
    if (!file || !title || !selectValues.documentType) return;

    try {
      await uploadDocument({
        tokenizedAssetID: asset.id,
        documentType: selectValues.documentType,
        documentTitle: title,
        documentFile: file,
      });

      setAddDocument(false);
      setFile(null);
      setTitle("");
      setSelectValues({});
    } catch (error: any) {
      const errorString =
        error?.response?.data?.message ||
        error?.data?.message ||
        error?.message ||
        error?.response?.data?.error ||
        error?.data?.error ||
        error?.error ||
        "An error occurred while updating.";

      const messageMatch = errorString.match(/message:(.*?)(\]$|$)/);
      let finalErrorMessage = errorString;
      if (messageMatch && messageMatch[1]) {
        finalErrorMessage = messageMatch[1].trim();
      }

      showErrorToast(finalErrorMessage);
    }
    console.log("Uploading with file:", file);
    console.log("Type of file:", typeof file); // should be object
    console.log("Is file instanceof File:", file instanceof File); //must be true
  };

  return (
    <>
      <Modal
        title="Add Asset Information Document"
        isOpen={addDocument}
        onClose={() => setAddDocument(false)}
        closeIconPosition="right"
      >
        <ModalContent>
          <HiddenFileInput
            id="fileUpload"
            type="file"
            onChange={handleFileChange}
          />
          <FileInputLabel htmlFor="fileUpload">
            {file?.name || "Upload or drag and drop files here"}
          </FileInputLabel>

          <Content>
            <div>
              <Text>Document Title</Text>
              <TitleInput
                placeholder="Tilte"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
              />
            </div>

            <DropdownSelect
              options={documentTypeOptions.map((opt) => opt.label)}
              placeholder="Select Type"
              labelText=""
              value={
                documentTypeOptions.find(
                  (opt) => opt.value === selectValues["documentType"],
                )?.label || ""
              }
              onSelect={(item) => {
                const selected = documentTypeOptions.find(
                  (opt) => opt.label === item,
                );
                setSelectValues({
                  ...selectValues,
                  documentType: selected?.value,
                });
              }}
            />
          </Content>

          <PrimaryButton buttonStyle={{ width: "100%" }} onClick={handleSave}>
            {isLoading ? "Saving..." : "Save Document"}
          </PrimaryButton>
        </ModalContent>
      </Modal>
    </>
  );
};

export default AddDocumentModal;

const ModalContent = styled.div`
  padding: 10px;
`;

const HiddenFileInput = styled.input`
  display: none;
`;

const FileInputLabel = styled.label`
  display: block;
  background-color: #f2f6f9;
  border: 1px solid #007cdf80;
  width: 100%;
  height: 48px;
  border-radius: 10px;
  color: #007cdf;
  font-size: 14px;
  text-align: center;
  line-height: 48px;
  cursor: pointer;
  font-family: inherit;
  transition: background-color 0.2s ease;

  &:hover {
    background-color: #e6eff5;
  }
`;

const Input = styled.input`
  background-color: #f2f6f9;
  border: 1px solid #007cdf80;
  width: 100%;
  height: 48px;
  border-radius: 10px;

  &::placeholder {
    color: #007cdf;
    font-size: 14px;
    opacity: 1;
    text-align: center;
    font-family: inherit;
  }
`;

const TitleInput = styled.input`
  background-color: #f2f6f9;
  border: 1px solid #007cdf80;
  width: 100%;
  height: 48px;
  border-radius: 10px;
  padding: 0 4px;

  &::placeholder {
    color: #007cdf;
    font-size: 14px;
    padding: 0 4px;
    font-family: inherit;
  }
`;

const Content = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding-top: 25px;
`;
const Text = styled.p`
  font-size: 14px;
  padding-bottom: 4px;
`;
