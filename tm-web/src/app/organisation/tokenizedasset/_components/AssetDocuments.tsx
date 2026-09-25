"use client";

import DocumentSection from "./DocumentSection";

const dummyDocs = [
  {
    title: "Proof of Existence",
    fileName: "Atlantis Invoice.pdf",
    size: "203 kb",
    uploadedBy: "creator",
  },
  {
    title: "Proof of Asset Address",
    fileName: "Atlantis Invoice.pdf",
    size: "203 kb",
    uploadedBy: "admin",
  },
  {
    title: "Proof of Ownership",
    fileName: "Atlantis Invoice.pdf",
    size: "203 kb",
    uploadedBy: "admin",
  },
];

const AssetDocuments = () => {
  return (
    <>
      <DocumentSection
        title="Asset Information Documents"
        documents={[...dummyDocs, ...dummyDocs]}
        showActions={false}
      />

      <DocumentSection title="Asset Value Documents" documents={dummyDocs} />

      <DocumentSection
        title="Asset Protection Documents"
        documents={[...dummyDocs, ...dummyDocs]}
      />

      <DocumentSection
        title="Asset Verification Undertaking Documents"
        documents={[...dummyDocs, ...dummyDocs]}
      />
    </>
  );
};

export default AssetDocuments;
