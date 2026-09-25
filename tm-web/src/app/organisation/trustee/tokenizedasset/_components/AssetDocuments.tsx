"use client";

import DocumentSection from "./DocumentSection";
import { IStakeholderAssetTokenizationDocument } from "@/redux/api/sharedstakeholders";

interface AssetDocumentsProps {
  documents?: IStakeholderAssetTokenizationDocument[];
}

const AssetDocuments = ({ documents = [] }: AssetDocumentsProps) => {
  const tokenizationDocuments = documents.map((document) => ({
    title: document.title || `Tokenization Document ${document.id}`,
    fileName: document.title || `Tokenization Document ${document.id}`,
    uploadedBy: "creator",
    fileUrl: document.url,
    onClick: document.url
      ? () => window.open(document.url, "_blank", "noopener,noreferrer")
      : undefined,
  }));

  return (
    <DocumentSection
      title="Asset Verification Documents"
      documents={tokenizationDocuments}
      showActions={false}
    />
  );
};

export default AssetDocuments;
