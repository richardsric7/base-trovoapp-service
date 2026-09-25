"use client";

import React, { FormEvent, Suspense, useMemo, useState } from "react";
import { useSearchParams } from "next/navigation";
import styled from "styled-components";
import { FiExternalLink, FiFileText, FiPlus } from "react-icons/fi";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import CustomFilter from "@/components/CustomFilter";
import Loader from "@/components/Loader";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import SearchableAssetSelect from "@/app/organisation/assetmanager/components/SearchableAssetSelect";
import {
  DropdownSelect,
  SearchBar,
  showErrorToast,
  showSuccessToast,
} from "@/components";
import {
  AssetManagerReport,
  useGenerateAssetManagerReportMutation,
  useGetAssetManagerAssetsQuery,
  useGetAssetManagerReportTypesQuery,
  useGetAssetManagerReportsQuery,
  useSubmitAssetManagerReportMutation,
} from "@/redux/api/assetManager";
import {
  useDownloadStakeholderDocumentMutation,
  useUploadStakeholderDocumentMutation,
} from "@/redux/api/sharedstakeholders";

type AssetOption = { id: string; code: string; name: string };
type ReportForm = {
  assetId: string;
  reportType: string;
  title: string;
  file: File | null;
};

const initialForm: ReportForm = {
  assetId: "",
  reportType: "",
  title: "",
  file: null,
};

const getErrorMessage = (error: any, fallback: string) =>
  error?.data?.message ?? error?.error ?? fallback;

const formatDate = (value?: string | null) => {
  if (!value) return "—";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "—";
  return date.toLocaleDateString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
};

const formatLabel = (value: string) =>
  value.replaceAll("_", " ").replace(/\b\w/g, (letter) => letter.toUpperCase());

const ReportsPageContent = () => {
  const searchParams = useSearchParams();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState("all");
  const [type, setType] = useState("all");
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [draftStatus, setDraftStatus] = useState("all");
  const [draftType, setDraftType] = useState("all");
  const [showGenerate, setShowGenerate] = useState(
    () => searchParams.get("generate") === "true",
  );
  const [reportToSubmit, setReportToSubmit] =
    useState<AssetManagerReport | null>(null);
  const [form, setForm] = useState<ReportForm>(initialForm);

  const { data, isLoading, isFetching, isError, refetch } =
    useGetAssetManagerReportsQuery({ page: currentPage, limit: pageSize });
  const { data: assetsResponse, isLoading: isLoadingAssets } =
    useGetAssetManagerAssetsQuery({ page: 1, limit: 100 });
  const { data: reportTypesResponse, isLoading: isLoadingReportTypes } =
    useGetAssetManagerReportTypesQuery();
  const [uploadDocument, { isLoading: isUploading }] =
    useUploadStakeholderDocumentMutation();
  const [downloadDocument] = useDownloadStakeholderDocumentMutation();
  const [openingDocumentId, setOpeningDocumentId] = useState<string | null>(null);
  const [generateReport, { isLoading: isGenerating }] =
    useGenerateAssetManagerReportMutation();
  const [submitReport, { isLoading: isSubmitting }] =
    useSubmitAssetManagerReportMutation();

  const assets = useMemo<AssetOption[]>(() => {
    const response = assetsResponse?.data as any;
    const records = Array.isArray(response)
      ? response
      : Array.isArray(response?.records)
        ? response.records
        : Array.isArray(response?.assets)
          ? response.assets
          : [];

    return records
      .map((asset: any) => ({
        id: String(asset.asset_id ?? asset.assetId ?? asset.id ?? ""),
        code: String(asset.asset_code ?? asset.assetCode ?? ""),
        name: String(asset.asset_name ?? asset.assetName ?? ""),
      }))
      .filter((asset: AssetOption) => asset.id && asset.code);
  }, [assetsResponse]);

  const reports = data?.data?.records ?? [];
  const meta = data?.data?.meta;
  const reportTypes = reportTypesResponse?.data ?? [];
  const visibleReports = useMemo(() => {
    const term = search.trim().toLowerCase();
    return reports.filter((report) => {
      const matchesSearch =
        !term ||
        report.title.toLowerCase().includes(term) ||
        report.asset_code.toLowerCase().includes(term);
      const matchesStatus = status === "all" || report.status === status;
      const matchesType = type === "all" || report.report_type === type;
      return matchesSearch && matchesStatus && matchesType;
    });
  }, [reports, search, status, type]);

  const selectedAsset = assets.find((asset) => asset.id === form.assetId);
  const canGenerate = Boolean(
    selectedAsset && form.reportType && form.title.trim() && form.file,
  );

  const handleGenerate = async (event: FormEvent) => {
    event.preventDefault();
    if (!canGenerate || !selectedAsset) return;
    try {
      const uploadPayload = new FormData();
      uploadPayload.append("document_file", form.file!);
      uploadPayload.append("asset_id", selectedAsset.id);
      uploadPayload.append("category", "asset_report");
      uploadPayload.append("title", form.title.trim());
      const uploadedDocument = await uploadDocument(uploadPayload).unwrap();
      await generateReport({
        asset_id: selectedAsset.id,
        asset_code: selectedAsset.code,
        report_type: form.reportType,
        title: form.title.trim(),
        document_id: uploadedDocument.data.id,
      }).unwrap();
      showSuccessToast("Report generated successfully");
      setForm(initialForm);
      setShowGenerate(false);
    } catch (error) {
      showErrorToast(getErrorMessage(error, "Unable to generate report"));
    }
  };

  const handleSubmit = async () => {
    if (!reportToSubmit) return;
    try {
      await submitReport(reportToSubmit.id).unwrap();
      showSuccessToast("Report submitted successfully");
      setReportToSubmit(null);
    } catch (error) {
      showErrorToast(getErrorMessage(error, "Unable to submit report"));
    }
  };

  const handleOpenDocument = async (documentId: string) => {
    const previewWindow = window.open("", "_blank");
    if (!previewWindow) {
      showErrorToast("Allow pop-ups to view this document");
      return;
    }
    try {
      setOpeningDocumentId(documentId);
      const blob = await downloadDocument(documentId).unwrap();
      const objectUrl = URL.createObjectURL(blob);
      previewWindow.opener = null;
      previewWindow.location.href = objectUrl;
      window.setTimeout(() => URL.revokeObjectURL(objectUrl), 60_000);
    } catch (error) {
      previewWindow.close();
      showErrorToast(getErrorMessage(error, "Unable to open report document"));
    } finally {
      setOpeningDocumentId(null);
    }
  };

  const columns = [
    { title: "Report title", dataIndex: "title" },
    { title: "Asset", dataIndex: "asset_code" },
    {
      title: "Report type",
      dataIndex: "report_type",
      render: (value: string) => formatLabel(value),
    },
    {
      title: "Created",
      dataIndex: "created_at",
      render: (value: string) => formatDate(value),
    },
    {
      title: "Submitted",
      dataIndex: "submitted_at",
      render: (value: string | null) => formatDate(value),
    },
    {
      title: "Status",
      dataIndex: "status",
      render: (value: string) => (
        <StatusBadge $status={value}>{formatLabel(value)}</StatusBadge>
      ),
    },
    {
      title: "File",
      dataIndex: "file_url",
      render: (_: unknown, report: AssetManagerReport) => (
        <Actions onClick={(event) => event.stopPropagation()}>
          {report.document_id ? (
            <ViewButton
              type="button"
              disabled={openingDocumentId === report.document_id}
              onClick={() => void handleOpenDocument(report.document_id!)}
            >
              <FiExternalLink />
              {openingDocumentId === report.document_id ? "Opening..." : "View file"}
            </ViewButton>
          ) : report.file_url ? (
            <ViewLink
              href={report.file_url}
              target="_blank"
              rel="noopener noreferrer"
            >
              <FiExternalLink /> View file
            </ViewLink>
          ) : (
            <NoAction>—</NoAction>
          )}
        </Actions>
      ),
    },
    {
      title: "Action",
      dataIndex: "action",
      render: (_: unknown, report: AssetManagerReport) =>
        report.status !== "submitted" ? (
          <SubmitButton
            onClick={(event) => {
              event.stopPropagation();
              setReportToSubmit(report);
            }}
          >
            Submit
          </SubmitButton>
        ) : (
          <NoAction>—</NoAction>
        ),
    },
  ];

  if (isLoading) return <Loader />;

  return (
    <Page>
      <Card>
        <Header>
          <div>
            <Title>Reports</Title>
          </div>
          <PrimaryAction onClick={() => setShowGenerate(true)}>
            <FiPlus /> Generate Report
          </PrimaryAction>
        </Header>

        <Filters>
          <SearchBar
            customWidth="320px"
            placeholder="Search by title or asset code"
            value={search}
            onChange={(event) => setSearch(event.target.value)}
          />
          <CustomFilter
            position={{ top: "205px", right: "56px" }}
            open={isFilterOpen}
            onClose={setIsFilterOpen}
          >
            <DropdownSelect
              value={draftType === "all" ? "All" : formatLabel(draftType)}
              options={["All", ...reportTypes.map(({ label }) => label)]}
              placeholder="All"
              labelText="Report Type"
              onSelect={(_, index) =>
                setDraftType(index === 0 ? "all" : reportTypes[index - 1]?.value ?? "all")
              }
            />
            <DropdownSelect
              value={draftStatus === "all" ? "All" : formatLabel(draftStatus)}
              options={["All", "Generated", "Submitted"]}
              placeholder="All"
              labelText="Status"
              onSelect={(_, index) =>
                setDraftStatus(["all", "generated", "submitted"][index])
              }
            />
            <FilterActions>
              <SecondaryButton
                onClick={() => {
                  setDraftType("all");
                  setDraftStatus("all");
                  setType("all");
                  setStatus("all");
                }}
                buttonStyle={{ width: "100%", margin: 0 }}
              >
                Reset
              </SecondaryButton>
              <PrimaryButton
                onClick={() => {
                  setType(draftType);
                  setStatus(draftStatus);
                  setIsFilterOpen(false);
                }}
                buttonStyle={{ width: "100%", margin: 0 }}
              >
                Apply
              </PrimaryButton>
            </FilterActions>
          </CustomFilter>
        </Filters>

        {isError ? (
          <State>
            <strong>Unable to load reports</strong>
            <span>Please check your connection and try again.</span>
            <RetryButton onClick={() => refetch()}>Try again</RetryButton>
          </State>
        ) : reports.length === 0 ? (
          <State>
            <EmptyIcon>
              <FiFileText />
            </EmptyIcon>
            <strong>No reports yet</strong>
            <span>Generate your first report for a managed asset.</span>
            <PrimaryAction onClick={() => setShowGenerate(true)}>
              <FiPlus /> Generate Report
            </PrimaryAction>
          </State>
        ) : (
          <>
            <TableWrap $fetching={isFetching}>
              <CustomTable columns={columns} dataSource={visibleReports} />
            </TableWrap>
            <Pagination
              currentPage={meta?.page ?? currentPage}
              totalCount={meta?.total ?? reports.length}
              pageSize={meta?.limit ?? pageSize}
              onPageChange={setCurrentPage}
              onPageSizeChange={(size) => {
                setPageSize(size);
                setCurrentPage(1);
              }}
              isFetching={isFetching}
            />
          </>
        )}
      </Card>

      {showGenerate && (
        <Overlay onMouseDown={() => !isGenerating && setShowGenerate(false)}>
          <ModalCard onMouseDown={(event) => event.stopPropagation()}>
            <ModalHeader>
              <div>
                <ModalTitle>Generate Report</ModalTitle>
                <ModalSubtitle>
                  Add a report for one of your assets.
                </ModalSubtitle>
              </div>
              <CloseButton
                aria-label="Close"
                onClick={() => setShowGenerate(false)}
              >
                ×
              </CloseButton>
            </ModalHeader>
            <Form onSubmit={handleGenerate}>
              <Field>
                <Label htmlFor="report-asset">Asset <RequiredMark>*</RequiredMark></Label>
                <SearchableAssetSelect
                  id="report-asset"
                  assets={assets}
                  value={form.assetId}
                  disabled={isLoadingAssets || isUploading}
                  loading={isLoadingAssets}
                  onChange={(assetId) =>
                    setForm((current) => ({
                      ...current,
                      assetId,
                    }))
                  }
                />
              </Field>
              <Field>
                <Label htmlFor="report-type">Report type <RequiredMark>*</RequiredMark></Label>
                <InputSelect
                  id="report-type"
                  value={form.reportType}
                  onChange={(event) =>
                    setForm((current) => ({
                      ...current,
                      reportType: event.target.value,
                    }))
                  }
                >
                  <option value="">
                    {isLoadingReportTypes ? "Loading report types..." : "Select report type"}
                  </option>
                  {reportTypes.map((reportType) => (
                    <option key={reportType.value} value={reportType.value}>
                      {reportType.label}
                    </option>
                  ))}
                </InputSelect>
              </Field>
              <Field>
                <Label htmlFor="report-title">Report title <RequiredMark>*</RequiredMark></Label>
                <TextInput
                  id="report-title"
                  placeholder="Enter report title"
                  value={form.title}
                  onChange={(event) =>
                    setForm((current) => ({
                      ...current,
                      title: event.target.value,
                    }))
                  }
                />
              </Field>
              <Field>
                <Label htmlFor="report-file">Report document <RequiredMark>*</RequiredMark></Label>
                <FileInput
                  id="report-file"
                  type="file"
                  accept="application/pdf,image/jpeg,image/png"
                  disabled={isUploading}
                  onChange={(event) => {
                    const file = event.currentTarget.files?.[0] ?? null;
                    if (file && !["application/pdf", "image/jpeg", "image/png"].includes(file.type)) {
                      showErrorToast("Only PDF, JPEG and PNG documents are supported");
                      event.currentTarget.value = "";
                      return;
                    }
                    if (file && file.size > 10 * 1024 * 1024) {
                      showErrorToast("Document must be 10 MB or smaller");
                      event.currentTarget.value = "";
                      return;
                    }
                    setForm((current) => ({ ...current, file }));
                  }}
                />
                <FileHint>PDF, JPEG or PNG, up to 10 MB.</FileHint>
              </Field>
              <ModalActions>
                <CancelButton
                  type="button"
                  onClick={() => setShowGenerate(false)}
                >
                  Cancel
                </CancelButton>
                <PrimaryAction
                  type="submit"
                  disabled={!canGenerate || isGenerating || isUploading}
                >
                  {isUploading ? "Uploading..." : isGenerating ? "Generating..." : "Generate Report"}
                </PrimaryAction>
              </ModalActions>
            </Form>
          </ModalCard>
        </Overlay>
      )}

      {reportToSubmit && (
        <Overlay onMouseDown={() => !isSubmitting && setReportToSubmit(null)}>
          <ConfirmCard onMouseDown={(event) => event.stopPropagation()}>
            <ConfirmIcon>
              <FiFileText />
            </ConfirmIcon>
            <ModalTitle>Submit report?</ModalTitle>
            <ConfirmText>
              You are about to submit <strong>“{reportToSubmit.title}”</strong>{" "}
              to the assigned trustee. A report cannot be submitted twice.
            </ConfirmText>
            <ModalActions>
              <CancelButton onClick={() => setReportToSubmit(null)}>
                Cancel
              </CancelButton>
              <PrimaryAction onClick={handleSubmit} disabled={isSubmitting}>
                {isSubmitting ? "Submitting..." : "Submit Report"}
              </PrimaryAction>
            </ModalActions>
          </ConfirmCard>
        </Overlay>
      )}
    </Page>
  );
};

const ReportsPage = () => (
  <Suspense fallback={<Loader />}>
    <ReportsPageContent />
  </Suspense>
);

export default ReportsPage;

const Page = styled.main`
  padding: 24px;
  min-height: 100vh;
`;
const Card = styled.section`
  background: #fff;
  border-radius: 24px;
  padding: 32px;
`;
const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
`;
const Title = styled.h1`
  margin: 0;
  color: #00225a;
  font-size: 24px;
  line-height: 32px;
`;
const Subtitle = styled.p`
  margin: 6px 0 0;
  color: #828282;
  font-size: 14px;
`;
const PrimaryAction = styled.button`
  border: 0;
  border-radius: 8px;
  min-height: 44px;
  padding: 0 18px;
  background: #007cdf;
  color: #fff;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  &:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }
`;
const Filters = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  margin: 28px 0 20px;
`;
const FilterActions = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
`;
const TableWrap = styled.div<{ $fetching: boolean }>`
  opacity: ${({ $fetching }) => ($fetching ? 0.55 : 1)};
  transition: opacity 0.15s;
  overflow-x: auto;
`;
const StatusBadge = styled.span<{ $status: string }>`
  display: inline-flex;
  padding: 6px 10px;
  border-radius: 999px;
  font-weight: 500;
  color: ${({ $status }) => ($status === "submitted" ? "#16794b" : "#9a6700")};
  background: ${({ $status }) =>
    $status === "submitted" ? "#e8f8f0" : "#fff6dc"};
`;
const Actions = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;
const ViewLink = styled.a`
  color: #007cdf;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-weight: 500;
`;
const ViewButton = styled.button`
  border: 0;
  padding: 0;
  background: transparent;
  color: #007cdf;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font: inherit;
  font-weight: 500;
  cursor: pointer;
  &:disabled {
    cursor: wait;
    opacity: 0.65;
  }
`;
const SubmitButton = styled.button`
  border: 1px solid #007cdf;
  color: #007cdf;
  background: #fff;
  border-radius: 7px;
  padding: 7px 12px;
  cursor: pointer;
  font: inherit;
`;
const NoAction = styled.span`
  color: #98a2b3;
`;
const State = styled.div`
  min-height: 330px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: #00225a;
  gap: 8px;
  span {
    color: #828282;
    margin-bottom: 12px;
  }
`;
const EmptyIcon = styled.div`
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: #eef7ff;
  color: #007cdf;
  display: grid;
  place-items: center;
  font-size: 26px;
  margin-bottom: 6px;
`;
const RetryButton = styled.button`
  border: 1px solid #007cdf;
  color: #007cdf;
  background: #fff;
  border-radius: 8px;
  padding: 10px 18px;
  cursor: pointer;
`;
const Overlay = styled.div`
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(0, 20, 52, 0.48);
  display: grid;
  place-items: center;
  padding: 24px;
`;
const ModalCard = styled.div`
  width: min(100%, 540px);
  max-height: calc(100vh - 48px);
  overflow-y: auto;
  background: #fff;
  border-radius: 20px;
  padding: 28px;
  box-shadow: 0 24px 60px rgba(0, 34, 90, 0.2);
`;
const ConfirmCard = styled(ModalCard)`
  max-width: 440px;
  text-align: center;
`;
const ModalHeader = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 20px;
`;
const ModalTitle = styled.h2`
  margin: 0;
  color: #00225a;
  font-size: 22px;
`;
const ModalSubtitle = styled.p`
  margin: 5px 0 0;
  color: #828282;
  font-size: 14px;
`;
const CloseButton = styled.button`
  border: 0;
  background: transparent;
  color: #667085;
  font-size: 28px;
  line-height: 1;
  cursor: pointer;
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
const RequiredMark = styled.span`color: #d92d20;`;
const fieldCss = `height: 46px; width: 100%; border: 1px solid #e0e0e0; border-radius: 8px; padding: 0 12px; outline: none; font: inherit; color: #00225a; background: #fff;`;
const TextInput = styled.input`
  ${fieldCss} &:focus {
    border-color: #007cdf;
  }
`;
const InputSelect = styled.select`
  ${fieldCss} &:focus {
    border-color: #007cdf;
  }
`;
const FileInput = styled.input`
  ${fieldCss}
  height: auto;
  padding: 10px 12px;
`;
const FileHint = styled.span`
  color: #667085;
  font-size: 12px;
`;
const ModalActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 10px;
`;
const CancelButton = styled.button`
  min-height: 44px;
  padding: 0 18px;
  border-radius: 8px;
  border: 1px solid #d0d5dd;
  background: #fff;
  color: #00225a;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
`;
const ConfirmIcon = styled(EmptyIcon)`
  margin: 0 auto 18px;
`;
const ConfirmText = styled.p`
  color: #667085;
  line-height: 1.6;
  margin: 12px 0 22px;
`;
