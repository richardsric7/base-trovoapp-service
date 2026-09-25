"use client";
import React, { useState } from "react";
import styled from "styled-components";
import CustomPDFViewer from "./CustomPDFViewer";
import { TokenizationRecord } from "@/redux/api/assettokenization";
import Image from "next/image";
import pdfIcon from "@/assets/images/pdfdoc.svg";
interface DocumentType {
  id: number;
  title: string;
  file: string;
  size: string;
}

interface PaymentProps {
  asset?: TokenizationRecord;
}
const Payment: React.FC<PaymentProps> = ({ asset }) => {
  const [selectedDocument, setSelectedDocument] = useState<DocumentType | null>(
    null
  );

  // Static document with PDF from public folder
  // const [documents] = useState<DocumentType[]>([
  //   {
  //     id: 1,

  //     file: "/file-sample_150kB.pdf", // Assuming your PDF is in public folder named sample.pdf
  //     size: "2.5 MB",
  //   },
  // ]);
  const paymentReceiptDocuments = asset?.ProofOfPaymentDocuments?.map(
    (doc, index) => ({
      id: doc?.id,
      title: `Proof of Payment ${index + 1}`,
      file: doc?.documentUrl,
      size: "", // size is not provided in your data
    })
  );

  const closePdfViewer = () => {
    setSelectedDocument(null);
  };

  const paymentData = [
    {
      id: 1,
      title: "Preferred Payment Method",
      subTitle: asset?.assetQuoteCurrency || "N/A",
    },
    {
      id: 2,
      title: "Wallet Address",
      subTitle: asset?.issuingWalletAddress || "N/A",
    },

    {
      id: 3,
      title: "Amount",
      subTitle: asset?.issuingWalletAddress || "N/A",
    },
    {
      id: 4,
      title: "Payment Receipt",

      subTitle: (
        <>
          {paymentReceiptDocuments?.map((document: any) => (
            <DocumentWrapper
              key={document.id}
              onClick={() => setSelectedDocument(document)}
            >
              <DocumentContent>
                <Image src={pdfIcon} alt="PDF" width={18} height={18} />
                <DocumentTitle>{document.title}</DocumentTitle>
              </DocumentContent>
              {document.size && <DocumentSize>{document.size}</DocumentSize>}
            </DocumentWrapper>
          ))}
        </>
      ),

      // subTitle: (
      //   // <DocumentWrapper
      //   //   onClick={() => setSelectedDocument(paymentReceiptDocument)}
      //   // >
      //   //   <DocumentContent>
      // <Image src={pdfIcon} alt="PDF" width={18} height={18} />
      //   //     <DocumentTitle>{paymentReceiptDocument?.title}</DocumentTitle>
      //   //   </DocumentContent>
      //   //   <DocumentSize>{paymentReceiptDocument?.size}</DocumentSize>
      //   // </DocumentWrapper>

      // ),
    },
  ];

  return (
    <>
      <Heading>Tokenization Application Payment</Heading>
      <Wrapper>
        {paymentData.map((data) => (
          <Card key={data.id}>
            <Label>{data.title}</Label>
            <Value>{data.subTitle}</Value>
          </Card>
        ))}
      </Wrapper>

      {selectedDocument && (
        <CustomPDFViewer
          pdfDocument={selectedDocument}
          onClose={closePdfViewer}
        />
      )}
    </>
  );
};

export default Payment;

// Styled Components
const Wrapper = styled.section`
  padding: 20px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
  padding: 16px 0;
`;

const Card = styled.div`
  display: flex;
  flex-direction: column;
`;

const Label = styled.p`
  font-size: 14px;
  font-weight: 400;
  color: #828282;
  margin: 0;
`;

const Value = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  margin: 8px 0;
  overflow-wrap: break-word;
  word-break: break-word;
  white-space: normal;
`;

const DocumentWrapper = styled.div`
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-bottom: 6px;

  &:hover {
    opacity: 0.8;
  }
`;

const DocumentTitle = styled.span`
  color: #00225a;
  font-size: 14px;
  font-weight: 500;
`;

const DocumentSize = styled.span`
  color: #828282;
  font-size: 12px;
`;
const DocumentContent = styled.div`
  display: flex;
  align-items: center;
  row-gap: 10px;
  cursor: pointer;
`;
