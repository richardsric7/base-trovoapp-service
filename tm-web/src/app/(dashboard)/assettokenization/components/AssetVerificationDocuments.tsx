"use client";
import React, { useState, useEffect } from "react";
import styled from "styled-components";
import Image from "next/image";
import * as pdfjs from "pdfjs-dist";
import { TbTrash } from "react-icons/tb";

import pdfIcon from "@/assets/images/pdfdoc.svg";
import CustomPDFViewer from "./CustomPDFViewer";
import { TokenizationRecord } from "@/redux/api/assettokenization";
import SecondaryButton from "@/components/SecondaryButton";
import { FaPlus } from "react-icons/fa6";
import AddDocumentModal from "./AddDocumentModal";
// Explicitly set the worker source
pdfjs.GlobalWorkerOptions.workerSrc = `//cdnjs.cloudflare.com/ajax/libs/pdf.js/${pdfjs.version}/pdf.worker.min.js`;

// Set the worker source
pdfjs.GlobalWorkerOptions.workerSrc = "/pdf.worker.min.mjs";

interface DocumentType {
  id: number;
  title: string;
  file?: string | File; // Accepts both URL string and File object.
  fileUrl?: string;
  size: string;
}

interface AssetDocumentsProps {
  asset?: TokenizationRecord;
}
const AssetVerificationDocuments: React.FC<AssetDocumentsProps> = ({
  asset,
}) => {
  const [documents, setDocuments] = useState<DocumentType[]>([]);

  const [selectedDocument, setSelectedDocument] = useState<DocumentType | null>(
    null
  );
  const [addDocument, setAddDocument] = useState<boolean>(false);

  useEffect(() => {
    const loadDocuments = async () => {
      let docs: DocumentType[] = [];

      // Load Asset Verification Documents if they exist.
      if (
        asset &&
        (asset as any).AssetTokenizationDocuments &&
        (asset as any).AssetTokenizationDocuments.length > 0
      ) {
        const assetDocs: DocumentType[] = (
          asset as any
        ).AssetTokenizationDocuments.map((doc: any) => ({
          id: doc.id,
          title: doc.documentTitle,
          fileUrl: doc.documentUrl,
          size: "", // Size can be computed/fetched if needed.
        }));
        docs = [...docs, ...assetDocs];
      }

      // Load Proof of Payment Documents if they exist.
      if (
        asset &&
        asset.ProofOfPaymentDocuments &&
        asset.ProofOfPaymentDocuments.length > 0
      ) {
        const paymentDocs: DocumentType[] = asset.ProofOfPaymentDocuments.map(
          (doc: any, index: number) => ({
            id: doc.id,
            title: `Proof of Payment ${index + 1}`,
            fileUrl: doc.documentUrl,
            size: "", // Size is not provided in your data.
          })
        );
        docs = [...docs, ...paymentDocs];
      }

      // If no documents are available, leave the array empty.
      setDocuments(docs);
    };

    loadDocuments();
  }, [asset]);
  

  const openPdfViewer = (document: DocumentType) => {
    setSelectedDocument(document);
  };

  const closePdfViewer = () => {
    setSelectedDocument(null);
  };

  return (
    <>
      <HeadingContent>
        <Heading>Asset Verification Documents</Heading>
        <SecondaryButton
          buttonStyle={{
            marginRight: "20px",
            display: "flex",
            alignItems: "center",
            gap: "2px",
            marginBottom: "16px",
          }}
          onClick={() => setAddDocument(!addDocument)}
        >
          {" "}
          <FaPlus /> Add Document
        </SecondaryButton>
      </HeadingContent>

      {documents.length > 0 ? (
        <DocumentGrid>
          {documents.map((doc) => (
            <DocumentCard key={doc.id} onClick={() => openPdfViewer(doc)}>
              <DocumentHeader>
                <DocumentTitle>{doc.title}</DocumentTitle>
                <TbTrash color="#BE3800" />
              </DocumentHeader>
              <DocumentContent>
                <Image src={pdfIcon} alt="PDF" width={18} height={18} />
                <DocumentInfo>
                  <Title>{doc.title}</Title>
                  <DocumentSize>{doc.size}</DocumentSize>
                </DocumentInfo>
              </DocumentContent>
              <DocumentCreator>
                <CreatorTitle>
                  {doc.fileUrl
                    ? "Uploaded by Creator"
                    : "Uploaded by Creator or Admin"}
                </CreatorTitle>
                <CreatorName>{doc.fileUrl ? "C" : "C or A"} </CreatorName>
              </DocumentCreator>
            </DocumentCard>
          ))}
        </DocumentGrid>
      ) : (
        <NoDocumentsMessage>No document uploaded</NoDocumentsMessage>
      )}

      {selectedDocument && (
        <CustomPDFViewer
          pdfDocument={selectedDocument}
          onClose={closePdfViewer}
        />
      )}

      {asset && addDocument && (
        <AddDocumentModal
          addDocument={addDocument}
          setAddDocument={setAddDocument}
          asset={asset}
        />
      )}
    </>
  );
};

export default AssetVerificationDocuments;

const DocumentGrid = styled.section`
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;

  color: #00225a;
`;
const DocumentCard = styled.div`
  background-color: #f2f6f9;
  border-radius: 8px;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
`;
const DocumentContent = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
`;

const DocumentHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 4px;
  cursor: pointer;
`;
const DocumentTitle = styled.h1`
  font-size: 14px;
  font-weight: 400;
  line-height: 17.07px;
  color: #828282;
  margin: 4px 0;

  overflow-wrap: break-word;
  word-break: break-word;
  white-space: normal;

  text-transform: capitalize;
`;
const Title = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #191919;
  margin: 4px 0;
  text-transform: capitalize;
`;
const DocumentInfo = styled.div`
  margin-top: 10px;
`;

const DocumentSize = styled.p`
  font-size: 14px;
  font-weight: 400;
  margin: 0;
  color: #828282;
`;
const HeadingContent = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  // margin: 10px 0;
  width: 100%;
`;

const DocumentCreator = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const CreatorTitle = styled.p`
  background-color: #ffffff;
  border-radius: 8px;
  padding: 4px 8px;
  color: #191919;

  font-weight: 400;
  font-size: 12px;
  line-height: 100%;
  letter-spacing: 0%;
`;

const CreatorName = styled.p`
  background-color: #007cdf;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  font-size: 13px;
  font-weight: 600;
  text-align: center;
  justify-content: center;
  color: #ffffff;
`;
const NoDocumentsMessage = styled.p`
  font-size: 16px;
  font-weight: 500;
  color: #828282;
  text-align: center;
  padding: 20px;
`;
