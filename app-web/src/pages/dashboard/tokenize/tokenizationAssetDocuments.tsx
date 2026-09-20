import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../../../components/button';
import ButtonSecondary from '../../../components/buttonSecondary';
import { hideLoader, showLoader } from '../../../utils/showToaster';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../../store/reduxStore';
import { Encryptor } from '../../../utils/encryptor';
import {
  setActiveTokenizedAsset,
  setTokenizationData,
} from '../../../store/appStateSlice';
import { TokenizedAsset } from '../../../types/tokenizedAsset';
import {
  useLazyFetchFormFieldsQuery,
  useUploadDocumentMutation,
  useDeleteDocumentMutation,
  useLazyGetTokenizationDetailQuery,
} from '../../../store/api/tokenizationApis';
import { useActiveTokenizedAsset } from './useActiveTokenizedAsset';
import { useTokenizationData } from './useTokenizationData';
type DocumentFile = {
  name?: string;
  documentType?: string;
  description?: string;
  required?: boolean;
  requiredForDebtAndHybrid?: boolean;
  public?: boolean;
  fileType?: string;
};

type DocumentSection = {
  sectionName?: string;
  files: DocumentFile[];
};

type UploadedDocument = {
  id: string;
  documentTitle: string;
  documentType: string;
  documentUrl: string;
};

function getStatusCode(responseData: unknown): number | undefined {
  const record = responseData as
    | {
        statusCode?: number;
        data?: {
          statusCode?: number;
          data?: { statusCode?: number };
        };
      }
    | undefined;

  return (
    record?.statusCode ??
    record?.data?.statusCode ??
    record?.data?.data?.statusCode
  );
}

function extractTokenizedAssetFromResponse(
  responseData: unknown,
): TokenizedAsset | undefined {
  const source = responseData as
    | {
        data?: {
          data?: TokenizedAsset;
        };
      }
    | { data?: TokenizedAsset }
    | TokenizedAsset
    | undefined;

  if (!source) {
    return undefined;
  }

  if ('data' in source && source.data) {
    const nested = source.data as { data?: TokenizedAsset } | TokenizedAsset;
    if (
      nested &&
      typeof nested === 'object' &&
      'data' in nested &&
      (nested as { data?: TokenizedAsset }).data
    ) {
      return (nested as { data?: TokenizedAsset }).data;
    }

    return nested as TokenizedAsset;
  }

  return source as TokenizedAsset;
}

function normalizeDocumentTypeValue(value?: string): string {
  return (value ?? '').trim().toLowerCase();
}

function normalizeUploadedDocuments(payload: unknown): UploadedDocument[] {
  if (!Array.isArray(payload)) {
    return [];
  }

  return payload.reduce<UploadedDocument[]>((accumulator, entry) => {
    if (!entry || typeof entry !== 'object' || Array.isArray(entry)) {
      return accumulator;
    }

    const source = entry as Record<string, unknown>;
    accumulator.push({
      id: String(source.id ?? ''),
      documentTitle: String(source.documentTitle ?? ''),
      documentType: String(source.documentType ?? ''),
      documentUrl: String(source.documentUrl ?? ''),
    });

    return accumulator;
  }, []);
}

function parseApiFormPayload(payload: unknown): unknown {
  if (typeof payload === 'string') {
    const trimmedPayload = payload.trim();

    if (!trimmedPayload) {
      return null;
    }

    try {
      return parseApiFormPayload(JSON.parse(trimmedPayload));
    } catch {
      return payload;
    }
  }

  if (payload && typeof payload === 'object' && !Array.isArray(payload)) {
    const source = payload as Record<string, unknown>;

    if (source.formString) {
      return parseApiFormPayload(source.formString);
    }

    if (
      source.data &&
      typeof source.data === 'object' &&
      !Array.isArray(source.data)
    ) {
      return parseApiFormPayload(
        (source.data as Record<string, unknown>).formString,
      );
    }
  }

  return payload;
}

function normalizeDocumentDefinition(payload: unknown): DocumentSection[] {
  const parsedPayload = parseApiFormPayload(payload);

  if (
    !parsedPayload ||
    typeof parsedPayload !== 'object' ||
    Array.isArray(parsedPayload)
  ) {
    return [];
  }

  const source = parsedPayload as Record<string, unknown>;
  const documentsSource = source.documents;

  if (
    !documentsSource ||
    typeof documentsSource !== 'object' ||
    Array.isArray(documentsSource)
  ) {
    return [];
  }

  return Object.entries(documentsSource as Record<string, unknown>).reduce<
    DocumentSection[]
  >((accumulator, [, sectionValue]) => {
    if (
      !sectionValue ||
      typeof sectionValue !== 'object' ||
      Array.isArray(sectionValue)
    ) {
      return accumulator;
    }

    const section = sectionValue as Record<string, unknown>;
    const filesSource = Array.isArray(section.files) ? section.files : [];

    accumulator.push({
      sectionName: String(section.sectionName ?? 'Supporting Documents'),
      files: filesSource.reduce<DocumentFile[]>(
        (fileAccumulator, fileValue) => {
          if (
            !fileValue ||
            typeof fileValue !== 'object' ||
            Array.isArray(fileValue)
          ) {
            return fileAccumulator;
          }

          const fileEntry = fileValue as Record<string, unknown>;
          fileAccumulator.push({
            name: fileEntry.name ? String(fileEntry.name) : undefined,
            documentType: fileEntry.documentType
              ? String(fileEntry.documentType)
              : undefined,
            description: fileEntry.description
              ? String(fileEntry.description)
              : undefined,
            required: Boolean(fileEntry.required),
            requiredForDebtAndHybrid: Boolean(
              fileEntry.requiredForDebtAndHybrid,
            ),
            public: Boolean(fileEntry.public),
            fileType: fileEntry.fileType
              ? String(fileEntry.fileType)
              : undefined,
          });

          return fileAccumulator;
        },
        [],
      ),
    });

    return accumulator;
  }, []);
}

export function TokenizationAssetDocuments() {
  const [getForm, {}] = useLazyFetchFormFieldsQuery();
  const [uploadDocument] = useUploadDocumentMutation();
  const [deleteDocument] = useDeleteDocumentMutation();
  const [getTokenizationDetail] = useLazyGetTokenizationDetailQuery();
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const appUser = useSelector((state: RootState) => state.auth.user!);
  const { activeTokenizedAsset, assetId } = useActiveTokenizedAsset();
  const { tokenizationData: tokenizationDataFromHook } = useTokenizationData();
  const tokenizationData = tokenizationDataFromHook as
    | Record<string, unknown>
    | undefined;
  const fileInputRef = useRef<HTMLInputElement>(null);
  const tokenizedAssetId = assetId ?? activeTokenizedAsset?.id ?? '';

  const [documents, setDocuments] = useState<DocumentSection[]>([]);
  const [uploadedDocuments, setUploadedDocuments] = useState<
    UploadedDocument[]
  >([]);
  const [isLoading, setIsLoading] = useState(true);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [fundingStructure, setFundingStructure] = useState<number>(0);

  // Upload modal state
  const [showUploadModal, setShowUploadModal] = useState(false);
  const [uploadModalTitle, setUploadModalTitle] = useState('');
  const [uploadModalDocumentType, setUploadModalDocumentType] = useState('');
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);

  // Delete confirmation modal state
  const [documentToDelete, setDocumentToDelete] = useState<string | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  const syncStoreWithLatestDocuments = (detailData?: TokenizedAsset) => {
    if (!detailData) {
      return;
    }

    dispatch(setActiveTokenizedAsset(detailData));

    const baseTokenizationData =
      tokenizationData &&
      typeof tokenizationData === 'object' &&
      !Array.isArray(tokenizationData)
        ? tokenizationData
        : {};

    const mergedTokenizationData = {
      ...baseTokenizationData,
      AssetTokenizationDocuments: detailData.AssetTokenizationDocuments ?? [],
      fundingStructure: detailData.fundingStructure ?? 0,
    };

    dispatch(setTokenizationData(mergedTokenizationData as never));
  };

  const loadDocuments = async () => {
    setIsLoading(true);
    setErrorMessage(null);

    try {
      showLoader();
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);

      let resolvedPayload: unknown = null;
      try {
        const payload = {
          signer: appUser.primarySigner,
          address: appUser.address,
          secretKey: currentSecretKey,
          body: { formId: '1' },
        };

        const res = await getForm(payload);
        const responsePayload =
          res && typeof res === 'object' && 'data' in res && res.data
            ? (res.data as Record<string, unknown>).formString
            : null;

        resolvedPayload = parseApiFormPayload(responsePayload);
      } catch {
        // continue to the next endpoint if the request fails
      }

      const normalizedDocuments = normalizeDocumentDefinition(resolvedPayload);

      setDocuments(normalizedDocuments);
    } catch {
      setErrorMessage('Unable to load the document requirements.');
      setDocuments([]);
    } finally {
      hideLoader();
      setIsLoading(false);
    }
  };

  const refreshTokenizationDetail = async () => {
    if (!tokenizedAssetId) return;

    try {
      showLoader();
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);

      const res = await getTokenizationDetail({
        signer: appUser.primarySigner,
        address: appUser.address,
        secretKey: currentSecretKey,
        tokenizedAssetId,
      });

      if ('data' in res) {
        const responseData = res.data;
        const statusCode = getStatusCode(responseData);
        const detailData = extractTokenizedAssetFromResponse(responseData);

        if (statusCode === 200 || detailData) {
          const assetTokenizationDocuments =
            detailData?.AssetTokenizationDocuments ?? [];
          setUploadedDocuments(
            normalizeUploadedDocuments(assetTokenizationDocuments),
          );
          setFundingStructure(Number(detailData?.fundingStructure ?? 0));
          syncStoreWithLatestDocuments(detailData);
        }
      }
    } catch {
      // silently fail
    } finally {
      hideLoader();
    }
  };

  useEffect(() => {
    void loadDocuments();
  }, [appUser, getForm]);

  useEffect(() => {
    const documentsFromActiveAsset = normalizeUploadedDocuments(
      activeTokenizedAsset?.AssetTokenizationDocuments,
    );
    const documentsFromTokenizationData = normalizeUploadedDocuments(
      tokenizationData?.AssetTokenizationDocuments,
    );

    setUploadedDocuments(
      activeTokenizedAsset
        ? documentsFromActiveAsset
        : documentsFromTokenizationData,
    );

    setFundingStructure(
      Number(
        activeTokenizedAsset?.fundingStructure ??
          tokenizationData?.fundingStructure ??
          0,
      ),
    );
  }, [activeTokenizedAsset, tokenizationData]);

  const handleFileSelect = async (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const file = event.target.files?.[0];
    if (!file) return;

    // Reset file input
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }

    // Validate file size (max 900KB as in Dart code)
    if (file.size > 900000) {
      setUploadError('File size must be less than 900KB');
      return;
    }

    if (!tokenizedAssetId) {
      setUploadError(
        'Unable to find selected tokenized asset. Please reopen the application and try again.',
      );
      return;
    }

    setIsUploading(true);
    setUploadError(null);

    try {
      showLoader();
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);

      const shortId = crypto.randomUUID().split('-')[0];
      const documentTitle =
        `${uploadModalTitle.split('(')[0].trim()}-${shortId}.${file.name.split('.').pop()}`.toLowerCase();

      const res = await uploadDocument({
        signer: appUser.primarySigner,
        address: appUser.address,
        secretKey: currentSecretKey,
        tokenizedAssetId,
        documentTitle,
        documentType: uploadModalDocumentType,
        file,
      });

      if ('data' in res) {
        await refreshTokenizationDetail();
        setShowUploadModal(false);
      } else {
        setUploadError('Upload failed. Please try again.');
      }
    } catch {
      setUploadError('Sorry, something went wrong. Please try again.');
    } finally {
      setIsUploading(false);
      hideLoader();
    }
  };

  const handleDeleteDocument = (documentId: string) => {
    setDocumentToDelete(documentId);
  };

  const confirmDeleteDocument = async () => {
    if (!documentToDelete) return;

    setIsDeleting(true);

    try {
      showLoader();
      const encryptor = new Encryptor();
      const currentSecretKey = await encryptor.getSecretKey(appUser);

      const res = await deleteDocument({
        signer: appUser.primarySigner,
        address: appUser.address,
        secretKey: currentSecretKey,
        documentId: documentToDelete,
      });

      if ('data' in res) {
        await refreshTokenizationDetail();
        setDocumentToDelete(null);
      } else {
        setErrorMessage('Failed to delete document.');
        setDocumentToDelete(null);
      }
    } catch {
      setErrorMessage('Sorry, something went wrong. Please try again.');
      setDocumentToDelete(null);
    } finally {
      hideLoader();
      setIsDeleting(false);
    }
  };

  const handleViewDocument = (doc: UploadedDocument) => {
    const fileUrl = doc.documentUrl;
    if (!fileUrl) return;

    if (fileUrl.endsWith('.pdf')) {
      window.open(fileUrl, '_blank');
    } else {
      window.open(fileUrl, '_blank');
    }
  };

  const openUploadModal = (title: string, documentType: string) => {
    setUploadModalTitle(title);
    setUploadModalDocumentType(documentType);
    setUploadError(null);
    setShowUploadModal(true);
  };

  const getUploadedDocsForType = (documentType: string): UploadedDocument[] => {
    const normalizedType = normalizeDocumentTypeValue(documentType);
    return uploadedDocuments.filter((doc) => {
      return normalizeDocumentTypeValue(doc.documentType) === normalizedType;
    });
  };

  const renderDocumentCard = (section: DocumentSection) => {
    const visibleFiles = section.files.filter((file) => {
      if (fundingStructure === 0 && file.requiredForDebtAndHybrid) {
        return false;
      }

      return true;
    });

    const availableDocumentTypes = new Set(
      visibleFiles
        .map((file) => normalizeDocumentTypeValue(file.documentType))
        .filter((value) => value.length > 0),
    );
    const requiredDocumentTypes = new Set(
      visibleFiles
        .filter((file) => file.required)
        .map((file) => normalizeDocumentTypeValue(file.documentType))
        .filter((value) => value.length > 0),
    );
    const uploadedDocumentTypes = new Set(
      uploadedDocuments
        .map((doc) => normalizeDocumentTypeValue(doc.documentType))
        .filter((value) => value.length > 0),
    );

    const uploadedFileTypeCount = Array.from(availableDocumentTypes).filter(
      (documentType) => uploadedDocumentTypes.has(documentType),
    ).length;
    const uploadedRequiredFileTypeCount = Array.from(
      requiredDocumentTypes,
    ).filter((documentType) => uploadedDocumentTypes.has(documentType)).length;
    const totalDocuments = availableDocumentTypes.size;
    const totalRequiredFiles = requiredDocumentTypes.size;

    const progress =
      totalDocuments > 0 ? uploadedFileTypeCount / totalDocuments : 0;

    return (
      <div
        key={section.sectionName}
        className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm"
      >
        <div className="flex items-start justify-between gap-3">
          <div>
            <h3 className="font-semibold text-[#0E1F51]">
              {section.sectionName ?? 'Supporting Documents'}
            </h3>
            <p className="mt-1 text-sm text-gray-500">
              {`${uploadedFileTypeCount}/${totalDocuments} files uploaded (${uploadedRequiredFileTypeCount}/${totalRequiredFiles} required)`}
            </p>
          </div>
          <span className="rounded-full bg-primary-100 px-3 py-1 text-xs font-medium text-primary-700">
            {Math.round(progress * 100)}%
          </span>
        </div>

        <div className="mt-3 h-2 overflow-hidden rounded-full bg-gray-200">
          <div
            className="h-full rounded-full bg-primary-600 transition-all"
            style={{ width: `${Math.round(progress * 100)}%` }}
          />
        </div>

        <div className="mt-4 space-y-2">
          {section.files.map((file, index) => {
            // Skip debt and hybrid files if funding structure is equity
            if (fundingStructure === 0 && file.requiredForDebtAndHybrid) {
              return null;
            }

            const docsForType = getUploadedDocsForType(file.documentType ?? '');

            return (
              <div
                key={`${section.sectionName}-${file.documentType ?? `file-${index}`}-${index}`}
                className="rounded-lg border border-gray-100 bg-gray-50 p-3"
              >
                <div className="flex items-start justify-between gap-2">
                  <div className="flex-1">
                    <p className="text-sm font-medium text-gray-700">
                      {file.name ?? file.documentType ?? 'Document'}
                      {file.required ? (
                        <span className="ml-1 text-red-500">*</span>
                      ) : null}
                    </p>
                    {file.description ? (
                      <p className="mt-1 text-xs text-gray-500">
                        {file.description}
                      </p>
                    ) : null}
                  </div>
                  <span className="text-xs font-medium text-gray-500 whitespace-nowrap">
                    {file.required ? 'Required' : 'Optional'}
                  </span>
                </div>

                {/* Uploaded documents list */}
                {docsForType.length > 0 ? (
                  <div className="mt-2 space-y-1">
                    {docsForType.map((doc) => (
                      <div
                        key={doc.id}
                        className="flex items-center justify-between"
                      >
                        <button
                          type="button"
                          className="flex items-center gap-1 text-xs text-primary-600 hover:text-primary-800 underline"
                          onClick={() => handleViewDocument(doc)}
                        >
                          <svg
                            className="w-4 h-4"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                            />
                          </svg>
                          {doc.documentTitle}
                        </button>
                        <button
                          type="button"
                          className="text-red-500 hover:text-red-700 p-1"
                          onClick={() => handleDeleteDocument(doc.id)}
                          title="Delete document"
                        >
                          <svg
                            className="w-4 h-4"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                            />
                          </svg>
                        </button>
                      </div>
                    ))}
                  </div>
                ) : null}

                {/* Add File button */}
                <button
                  type="button"
                  className="mt-2 flex items-center justify-center gap-1 px-4 py-1.5 text-xs font-semibold text-gray-700 border border-primary-300 rounded-lg hover:bg-primary-50 transition-colors"
                  onClick={() =>
                    openUploadModal(
                      file.name ?? file.documentType ?? 'Document',
                      file.documentType ?? '',
                    )
                  }
                >
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                    />
                  </svg>
                  Add File
                </button>
              </div>
            );
          })}
        </div>
      </div>
    );
  };

  return (
    <div className="bg-white p-6 w-full mx-auto">
      <div className="space-y-6">
        {isLoading ? (
          <div className="rounded-lg border border-gray-200 bg-primary-100 p-6 text-sm text-gray-600">
            Loading document requirements...
          </div>
        ) : null}

        {errorMessage ? (
          <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-600">
            {errorMessage}
          </div>
        ) : null}

        {documents.length > 0 ? (
          <div className="rounded-lg border border-gray-200 bg-primary-100 p-6">
            <div className="mb-4">
              <h2 className="text-lg font-semibold text-[#0E1F51]">
                Upload the files required in each folder
              </h2>
              <p className="mt-1 text-sm text-gray-500">
                Review each folder and upload the supporting documents required
                for the asset.
              </p>
            </div>

            <div className="space-y-3">
              {documents.map((section) => renderDocumentCard(section))}
            </div>
          </div>
        ) : null}
      </div>

      <div className="w-full flex justify-end mt-10">
        <div className="flex self-end w-full max-w-2xl space-x-6">
          <div className="w-2/4">
            <ButtonSecondary
              label="Back"
              onclick={() => {
                navigate(-1);
              }}
            />
          </div>
          <div className="w-3/4">
            <Button
              label="Continue"
              onclick={() => {
                navigate(
                  `/dashboard/tokenize/apply/${tokenizedAssetId}/asset-token-information`,
                );
              }}
            />
          </div>
        </div>
      </div>

      {/* Hidden file input */}
      <input
        ref={fileInputRef}
        type="file"
        className="hidden"
        onChange={handleFileSelect}
        accept=".pdf,.jpg,.jpeg,.png,.doc,.docx"
      />

      {/* Upload Modal */}
      {showUploadModal ? (
        <div className="fixed inset-0 z-50 flex items-end justify-center sm:items-center">
          <div
            className="fixed inset-0 bg-black/50 transition-opacity"
            onClick={() => setShowUploadModal(false)}
          />
          <div className="relative w-full max-w-md bg-white rounded-t-2xl sm:rounded-2xl p-6 shadow-xl max-h-[80vh] overflow-y-auto">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-[#0E1F51] truncate pr-4">
                {uploadModalTitle}
              </h3>
              <button
                type="button"
                className="text-gray-400 hover:text-gray-600"
                onClick={() => setShowUploadModal(false)}
              >
                <svg
                  className="w-6 h-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>

            <p className="text-sm text-gray-500 mb-4">
              Select a file to upload. Files marked with * are required to
              proceed with your application.
            </p>

            {uploadError ? (
              <div className="mb-4 text-sm text-red-600">{uploadError}</div>
            ) : null}

            <div
              className="cursor-pointer"
              onClick={() => fileInputRef.current?.click()}
            >
              <div className="border-2 border-dashed border-primary-300 rounded-lg p-8 flex flex-col items-center justify-center gap-2 hover:bg-primary-50 transition-colors">
                <svg
                  className="w-8 h-8 text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                  />
                </svg>
                <span className="text-sm font-semibold text-gray-600">
                  {isUploading ? 'Uploading...' : 'Upload file here'}
                </span>
              </div>
            </div>
          </div>
        </div>
      ) : null}

      {/* Delete Confirmation Modal */}
      {documentToDelete ? (
        <div className="fixed inset-0 z-50 flex items-end justify-center sm:items-center">
          <div
            className="fixed inset-0 bg-black/50 transition-opacity"
            onClick={() => setDocumentToDelete(null)}
          />
          <div className="relative w-full max-w-md bg-white rounded-t-2xl sm:rounded-2xl p-6 shadow-xl">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-[#0E1F51]">
                Delete document
              </h3>
              <button
                type="button"
                className="text-gray-400 hover:text-gray-600"
                onClick={() => setDocumentToDelete(null)}
              >
                <svg
                  className="w-6 h-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>

            <p className="text-sm text-gray-500 mb-6">
              Are you sure you want to delete this document? This action cannot
              be undone.
            </p>

            <div className="flex justify-end gap-3">
              <button
                type="button"
                className="px-4 py-2 text-sm font-semibold text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                onClick={() => setDocumentToDelete(null)}
              >
                Cancel
              </button>
              <button
                type="button"
                className="px-4 py-2 text-sm font-semibold text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                onClick={() => void confirmDeleteDocument()}
                disabled={isDeleting}
              >
                {isDeleting ? 'Deleting...' : 'Delete'}
              </button>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
