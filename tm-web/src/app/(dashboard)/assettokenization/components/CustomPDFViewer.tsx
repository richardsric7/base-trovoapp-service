"use client";
import React, { useEffect, useRef, useState } from "react";
import * as pdfjs from "pdfjs-dist";
import styled from "styled-components";
import { FaAngleLeft, FaAngleRight } from "react-icons/fa6";
import { GrDocumentDownload } from "react-icons/gr";
import PrimaryButton from "@/components/PrimaryButton";

interface DocumentType {
  id: number;
  title: string;
  file?: string | File; // Accepts both URL string and File object.
  fileUrl?: string;
  size: string;
}

interface CustomPDFViewerProps {
  pdfDocument: DocumentType;
  onClose: () => void;
}

const CustomPDFViewer: React.FC<CustomPDFViewerProps> = ({
  pdfDocument,
  onClose,
}) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [pdf, setPdf] = useState<pdfjs.PDFDocumentProxy | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [numPages, setNumPages] = useState(0);
  const [scale, setScale] = useState(1.5);
  const [isImage, setIsImage] = useState(false);
  const [imageSrc, setImageSrc] = useState<string | null>(null);
  const [loadingError, setLoadingError] = useState<string | null>(null);
  const renderTaskRef = useRef<pdfjs.RenderTask | null>(null);

  const getPdfSourceUrl = () => {
    const pdfSource = pdfDocument.fileUrl || pdfDocument.file;
    if (!pdfSource) return null;
    if (typeof pdfSource === "string") return pdfSource;
    if (pdfSource instanceof File) return URL.createObjectURL(pdfSource);
    return null;
  };

  const checkIsImage = (urlOrFile: string | File | undefined) => {
    if (!urlOrFile) return false;
    if (urlOrFile instanceof File) {
      return urlOrFile.type.startsWith("image/");
    }
    const imageUrlPatterns = /\.(jpg|jpeg|png|gif|webp|bmp|svg)(\?.*)?$/i;
    return imageUrlPatterns.test(urlOrFile);
  };

  useEffect(() => {
    const loadDocument = async () => {
      setLoadingError(null);
      setIsImage(false);
      setImageSrc(null);
      setPdf(null);

      try {
        const pdfSource = pdfDocument.fileUrl || pdfDocument.file;
        if (!pdfSource) {
          throw new Error("No document source provided");
        }

        if (checkIsImage(pdfSource)) {
          setIsImage(true);
          const src =
            typeof pdfSource === "string"
              ? pdfSource
              : URL.createObjectURL(pdfSource);
          setImageSrc(src);
          return;
        }

        // It's a PDF - try loading it
        let urlToFetch: string;
        if (typeof pdfSource === "string") {
          urlToFetch = pdfSource;
        } else if (pdfSource instanceof File) {
          urlToFetch = URL.createObjectURL(pdfSource);
        } else {
          throw new Error("Invalid document source type");
        }

        console.log("Fetching PDF from:", urlToFetch);
        try {
          const response = await fetch(urlToFetch);
          if (!response.ok) {
            throw new Error(
              `Network response was not ok: ${response.statusText}`,
            );
          }
          const arrayBuffer = await response.arrayBuffer();
          const loadedPdf = await pdfjs.getDocument({ data: arrayBuffer })
            .promise;
          setPdf(loadedPdf);
          setNumPages(loadedPdf.numPages);
        } catch (fetchErr: any) {
          console.error("Fetch failed, likely CORS:", fetchErr);
          if (
            fetchErr.name === "TypeError" &&
            fetchErr.message === "Failed to fetch"
          ) {
            setLoadingError("CORS_BLOCK");
          } else {
            setLoadingError(fetchErr.message || "Error loading PDF");
          }
        }
      } catch (error: any) {
        console.error("Error in loadDocument:", error);
        setLoadingError(error.message || "An unknown error occurred");
      }
    };

    loadDocument();
  }, [pdfDocument]);

  const renderPage = async () => {
    if (!pdf || !canvasRef.current) return;

    if (renderTaskRef.current) {
      renderTaskRef.current.cancel();
    }

    try {
      const page = await pdf.getPage(currentPage);
      const canvas = canvasRef.current;
      const context = canvas.getContext("2d");
      const viewport = page.getViewport({ scale });

      canvas.height = viewport.height;
      canvas.width = viewport.width;

      const renderContext = {
        canvasContext: context!,
        viewport: viewport,
      };

      renderTaskRef.current = page.render(renderContext);
      await renderTaskRef.current.promise;
    } catch (error) {
      if ((error as any).name === "RenderingCancelledException") {
        console.log("Render task was cancelled");
      } else {
        console.error("Error rendering page:", error);
      }
    }
  };

  useEffect(() => {
    if (!isImage) {
      renderPage();
    }
  }, [pdf, currentPage, scale, isImage]);

  const changePage = (offset: number) => {
    setCurrentPage((prevPage) => {
      const newPage = prevPage + offset;
      return newPage > 0 && newPage <= numPages ? newPage : prevPage;
    });
  };

  const downloadPDF = async () => {
    const pdfSource = pdfDocument.fileUrl || pdfDocument.file;
    if (!pdfSource) return;

    let blob: Blob;
    try {
      if (typeof pdfSource === "string") {
        const response = await fetch(pdfSource);
        if (!response.ok) {
          throw new Error(
            `Network response was not ok: ${response.statusText}`,
          );
        }
        blob = await response.blob();
      } else if (pdfSource instanceof File) {
        blob = pdfSource;
      } else {
        throw new Error("Invalid document source type");
      }
    } catch (error) {
      console.error("Error downloading document:", error);
      // Fallback: open in new tab if fetch fails
      if (typeof pdfSource === "string") {
        window.open(pdfSource, "_blank");
      }
      return;
    }

    const blobUrl = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = blobUrl;
    link.download = pdfDocument.title || "document";
    document.body.appendChild(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(blobUrl);
  };

  return (
    <PDFViewerOverlay>
      <PDFViewerContent>
        <PDFHeader>
          <CloseButton onClick={onClose}>&times;</CloseButton>
          <PDFTitle>{pdfDocument.title}</PDFTitle>
          <DownloadText onClick={downloadPDF}>
            <GrDocumentDownload />
            Download
          </DownloadText>
        </PDFHeader>
        <PDFCanvas>
          {loadingError === "CORS_BLOCK" ? (
            <ErrorContainer>
              <ErrorTitle>CORS Restriction</ErrorTitle>
              <ErrorMessage>
                This document cannot be loaded directly in the viewer due to
                security restrictions (CORS).
              </ErrorMessage>
              <PrimaryButton
                onClick={() => window.open(getPdfSourceUrl() || "", "_blank")}
              >
                Open in New Tab
              </PrimaryButton>
            </ErrorContainer>
          ) : loadingError ? (
            <ErrorContainer>
              <ErrorTitle>Error Loading Document</ErrorTitle>
              <ErrorMessage>{loadingError}</ErrorMessage>
            </ErrorContainer>
          ) : isImage ? (
            <img
              src={imageSrc || ""}
              alt={pdfDocument.title}
              style={{
                maxWidth: "100%",
                height: "auto",
                maxHeight: "80vh",
                objectFit: "contain",
                padding: "20px 0",
              }}
            />
          ) : (
            <canvas ref={canvasRef} />
          )}
        </PDFCanvas>
        {!isImage && !loadingError && (
          <PDFControls>
            <ControlButton
              onClick={() => changePage(-1)}
              disabled={currentPage === 1}
            >
              <FaAngleLeft />
            </ControlButton>
            <PageInfo>
              Page {currentPage} of {numPages}
            </PageInfo>
            <ControlButton
              onClick={() => changePage(1)}
              disabled={currentPage === numPages}
            >
              <FaAngleRight />
            </ControlButton>
          </PDFControls>
        )}
      </PDFViewerContent>
    </PDFViewerOverlay>
  );
};

export default CustomPDFViewer;

const PDFViewerOverlay = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
`;

const PDFViewerContent = styled.div`
  background-color: white;
  border-radius: 8px;
  width: 90%;
  height: 90%;
  display: flex;
  flex-direction: column;
`;

const PDFHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  color: #00225a;
  background-color: #f5f5f5;
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
`;

const PDFTitle = styled.h2`
  margin: 0;
  font-size: 18px;
  color: #00225a;
`;

const PageInfo = styled.span`
  font-size: 14px;
  color: #00225a;
  font-weight: 600;
  line-height: 24px;
  text-align: center;
`;

const CloseButton = styled.button`
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #00225a;
`;

const PDFCanvas = styled.div`
  flex-grow: 1;
  overflow: auto;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #e0e0e0;
  padding: 40px 0;
`;

const PDFControls = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 15px;
  background-color: #f5f5f5;
  border-bottom-left-radius: 8px;
  border-bottom-right-radius: 8px;
`;

const ControlButton = styled.button`
  margin: 0 5px;
  padding: 8px 15px;
  background-color: #f2f6f9;
  border: none;
  border-radius: 8px;
  color: #00225a;
  cursor: pointer;
  transition: background-color 0.3s ease;

  &:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
  }
`;

const DownloadText = styled.p`
  color: #007cdf;
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
`;

const ErrorContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 40px;
  text-align: center;
  color: #00225a;
`;

const ErrorTitle = styled.h3`
  font-size: 20px;
  font-weight: 600;
  margin: 0;
`;

const ErrorMessage = styled.p`
  font-size: 16px;
  max-width: 400px;
  margin: 0;
  color: #828282;
`;
