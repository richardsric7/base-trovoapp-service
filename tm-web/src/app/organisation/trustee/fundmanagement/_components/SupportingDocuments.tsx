"use client";

import DocumentSection from "../../tokenizedasset/_components/DocumentSection";
import { IFundReleaseSupportingDocument } from "@/redux/api/trustees";
import { useDownloadStakeholderDocumentMutation } from "@/redux/api/sharedstakeholders";
import { message } from "antd";
import { useState } from "react";

interface SupportingDocumentsProps {
  documents: IFundReleaseSupportingDocument[];
}

const formatSize = (size: number) => {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / (1024 * 1024)).toFixed(1)} MB`;
};

const SupportingDocuments = ({ documents }: SupportingDocumentsProps) => {
  const [downloadDocument] = useDownloadStakeholderDocumentMutation();
  const [openingDocumentId, setOpeningDocumentId] = useState<string | null>(
    null,
  );

  const openDocument = async (document: IFundReleaseSupportingDocument) => {
    if (document.external_url) {
      window.open(document.external_url, "_blank", "noopener,noreferrer");
      return;
    }

    const previewWindow = window.open("", "_blank");
    if (!previewWindow) {
      message.error("Allow pop-ups to view this document.");
      return;
    }

    try {
      setOpeningDocumentId(document.id);
      const response = await downloadDocument(document.id).unwrap();
      const blob =
        response.type || !document.mime_type
          ? response
          : new Blob([response], { type: document.mime_type });
      const objectUrl = URL.createObjectURL(blob);
      previewWindow.opener = null;
      previewWindow.location.href = objectUrl;
      window.setTimeout(() => URL.revokeObjectURL(objectUrl), 60_000);
    } catch {
      previewWindow.close();
      message.error("Unable to open this document. Please try again.");
    } finally {
      setOpeningDocumentId(null);
    }
  };

  const items = documents.map((document) => ({
    title: document.title || document.original_filename || "Document",
    size: formatSize(document.size_bytes ?? 0),
    fileUrl: document.external_url || document.download_path,
    isLoading: openingDocumentId === document.id,
    onClick: () => void openDocument(document),
  }));

  return (
    <DocumentSection
      title={documents.length ? "Supporting Documents" : "No supporting documents"}
      documents={items}
      showActions={false}
    />
  );
};

export default SupportingDocuments;
