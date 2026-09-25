"use client";

import FundReleaseActor from "@/app/organisation/components/FundReleaseActor";
import { FormEvent, useMemo, useState } from "react";
import Link from "next/link";
import styled from "styled-components";
import { FiPlus } from "react-icons/fi";
import CustomFilter from "@/components/CustomFilter";
import Loader from "@/components/Loader";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import Pagination from "@/components/CustomPagination";
import { DropdownSelect, showErrorToast, showSuccessToast } from "@/components";
import FundCard from "@/app/organisation/trustee/fundmanagement/_components/FundCard";
import {
  AssetManagerCurrencyTotal,
  useCreateAssetManagerFundReleaseMutation,
  useGetAssetManagerAssetsQuery,
  useGetAssetManagerFundReleasesQuery,
  useGetAssetManagerFundSummaryQuery,
} from "@/redux/api/assetManager";
import {
  useGetStakeholderDocumentsQuery,
  useUploadStakeholderDocumentMutation,
} from "@/redux/api/sharedstakeholders";
import AmFundRealseRequestForm from "../components/AmFundRealseRequestForm";
import type {
  AssetOption,
  ReleaseForm,
} from "../components/AmFundRealseRequestForm";
const emptyForm: ReleaseForm = {
  assetId: "",
  amount: "",
  currency: "",
  purpose: "",
  receivingAccountName: "",
  receivingAccountNumber: "",
  receivingBank: "",
};
const statuses = [
  "submitted",
  "trustee_approved",
  "execution_pending",
  "processing",
  "completed",
  "trustee_rejected",
  "failed",
];
const formatLabel = (value: string) =>
  value.replaceAll("_", " ").replace(/\b\w/g, (letter) => letter.toUpperCase());
const formatTotals = (totals?: AssetManagerCurrencyTotal[]) =>
  !totals?.length
    ? "0"
    : totals
        .map(({ amount, currency }) =>
          `${Number(amount).toLocaleString(undefined, { maximumFractionDigits: 20 })} ${currency}`.trim(),
        )
        .join(" • ");
const formatDate = (value: string) =>
  new Date(value).toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
const errorMessage = (error: any) =>
  error?.data?.message ??
  error?.message ??
  error?.data?.error ??
  error?.error ??
  "Unable to create fund release request";

const statusLabel = (status: string) =>
  ({
    submitted: "Pending Review",
    trustee_approved: "Approved",
    execution_pending: "Awaiting Release",
    completed: "Released",
    trustee_rejected: "Rejected",
  })[status] ?? formatLabel(status);
const statusColor = (status: string) =>
  ["trustee_approved", "completed"].includes(status)
    ? "#22c55e"
    : ["trustee_rejected", "failed"].includes(status)
      ? "#ef4444"
      : "#facc15";

export default function AssetManagerFundManagementPage() {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [status, setStatus] = useState("");
  const [draftStatus, setDraftStatus] = useState("");
  const [filterOpen, setFilterOpen] = useState(false);
  const [createOpen, setCreateOpen] = useState(false);
  const [form, setForm] = useState<ReleaseForm>(emptyForm);
  const [selectedDocumentIds, setSelectedDocumentIds] = useState<string[]>([]);
  const [documentFiles, setDocumentFiles] = useState<File[]>([]);
  const { data, isLoading, isFetching, isError, refetch } =
    useGetAssetManagerFundReleasesQuery({
      page: currentPage,
      limit: pageSize,
      ...(status ? { status } : {}),
    });
  const {
    data: summaryResponse,
    isLoading: summaryLoading,
    isError: summaryError,
    refetch: refetchSummary,
  } = useGetAssetManagerFundSummaryQuery();
  const { data: assetsResponse, isLoading: assetsLoading } =
    useGetAssetManagerAssetsQuery({ page: 1, limit: 100 });
  const [createRelease, { isLoading: isCreating }] =
    useCreateAssetManagerFundReleaseMutation();
  const { data: documentsResponse, isLoading: documentsLoading } =
    useGetStakeholderDocumentsQuery(
      { page: 1, limit: 100, asset_id: form.assetId },
      { skip: !form.assetId },
    );
  const [uploadDocument, { isLoading: isUploading }] =
    useUploadStakeholderDocumentMutation();
  const records = data?.data.records ?? [];
  const meta = data?.data.meta;
  const summary = summaryResponse?.data;
  const documents = (documentsResponse?.data.records ?? []).filter(
    (document) => document.category === "fund_release_supporting",
  );
  const assets = useMemo<AssetOption[]>(
    () =>
      (assetsResponse?.data?.records ?? [])
        .map((asset) => ({
          id: String(asset.id ?? ""),
          code: String(asset.assetCode ?? ""),
          name: String(asset.assetName ?? ""),
          currency: String(asset.assetQuoteCurrency ?? ""),
        }))
        .filter((asset) => asset.id),
    [assetsResponse],
  );
  const canCreate = Boolean(
    form.assetId &&
    Number(form.amount) > 0 &&
    form.currency.trim() &&
    form.purpose.trim(),
  );

  const handleCreate = async (event: FormEvent) => {
    event.preventDefault();
    if (!canCreate) return;
    try {
      const documentIds = [...selectedDocumentIds];
      for (const documentFile of documentFiles) {
        const uploadPayload = new FormData();
        uploadPayload.append("document_file", documentFile);
        uploadPayload.append("asset_id", form.assetId);
        uploadPayload.append("category", "fund_release_supporting");
        uploadPayload.append("title", documentFile.name);
        const uploaded = await uploadDocument(uploadPayload).unwrap();
        documentIds.push(uploaded.data.id);
      }
      await createRelease({
        asset_id: form.assetId,
        amount: form.amount,
        currency: form.currency.trim().toUpperCase(),
        purpose: form.purpose.trim(),
        ...(form.receivingAccountName.trim()
          ? { receiving_account_name: form.receivingAccountName.trim() }
          : {}),
        ...(form.receivingAccountNumber.trim()
          ? { receiving_account_number: form.receivingAccountNumber.trim() }
          : {}),
        ...(form.receivingBank.trim()
          ? { receiving_bank: form.receivingBank.trim() }
          : {}),
        supporting_document_ids: documentIds,
      }).unwrap();
      showSuccessToast("Fund release request created successfully");
      setForm(emptyForm);
      setSelectedDocumentIds([]);
      setDocumentFiles([]);
      setCreateOpen(false);
    } catch (error) {
      showErrorToast(errorMessage(error));
    }
  };

  const cards = [
    {
      title: "Remaining Balance",
      amount: formatTotals(summary?.remaining_balance),
    },
    {
      title: "Total Requested",
      amount: formatTotals(summary?.requested.amounts),
    },
    {
      title: "Total Approved",
      amount: formatTotals(summary?.approved.amounts),
    },
    {
      title: "Total Released",
      amount: formatTotals(summary?.released.amounts),
    },
    { title: "Total Pending", amount: formatTotals(summary?.pending.amounts) },
    {
      title: "Total Rejected",
      amount: formatTotals(summary?.rejected.amounts),
    },
  ];
  if (isLoading || summaryLoading) return <Loader />;

  return (
    <Container>
      {summaryError && (
        <ErrorState>
          <span>Unable to load the fund management summary.</span>
          <RetryButton onClick={() => refetchSummary()}>Try again</RetryButton>
        </ErrorState>
      )}
      <CardContainer>
        {cards.map((item) => (
          <FundCard key={item.title} {...item} />
        ))}
      </CardContainer>
      <ReleaseContainer>
        <Header>
          <Title>Fund Release Requests</Title>
          <HeaderActions>
            <CustomFilter
              position={{ top: "400px", right: "30px" }}
              open={filterOpen}
              onClose={setFilterOpen}
            >
              <FilterContent>
                <DropdownSelect
                  options={["All", ...statuses.map(formatLabel)]}
                  labelText="Status"
                  placeholder="All"
                  value={draftStatus ? formatLabel(draftStatus) : "All"}
                  onSelect={(_, index) =>
                    setDraftStatus(index === 0 ? "" : statuses[index - 1])
                  }
                />
                <FilterActions>
                  <SecondaryButton
                    onClick={() => {
                      setDraftStatus("");
                      setStatus("");
                      setCurrentPage(1);
                      setFilterOpen(false);
                    }}
                  >
                    Reset
                  </SecondaryButton>
                  <PrimaryButton
                    onClick={() => {
                      setStatus(draftStatus);
                      setCurrentPage(1);
                      setFilterOpen(false);
                    }}
                  >
                    Apply
                  </PrimaryButton>
                </FilterActions>
              </FilterContent>
            </CustomFilter>
            <CreateButton onClick={() => setCreateOpen(true)}>
              <FiPlus /> Create Request
            </CreateButton>
          </HeaderActions>
        </Header>
        {isError && (
          <EmptyState>
            Unable to load fund release requests.{" "}
            <InlineRetry onClick={() => refetch()}>Try again</InlineRetry>
          </EmptyState>
        )}
        {!isError && records.length === 0 && (
          <EmptyState>No fund release requests found.</EmptyState>
        )}
        {!isError &&
          records.map((item) => (
            <ActivityCard key={item.id}>
              <DetailsLink
                href={`/organisation/assetmanager/fundmanagement/${encodeURIComponent(item.id)}`}
              >
                View Details
              </DetailsLink>
              <TopRow>
                <LeftInfo>
                  <Avatar>
                    {(item.asset_code || "AS").slice(0, 2).toUpperCase()}
                  </Avatar>
                  <ProjectInfo>
                    <ProjectName>
                      {item.asset_code || item.asset_id}
                    </ProjectName>
                    <ProjectSub>{item.purpose}</ProjectSub>
                  </ProjectInfo>
                </LeftInfo>
              </TopRow>
              <Description>
                {Number(item.amount).toLocaleString()} {item.currency}
              </Description>
              <DetailsRow>
                <Detail>
                  <DetailLabel>Request ID</DetailLabel>
                  <Value title={item.id}>{item.id.slice(0, 8)}…</Value>
                </Detail>
                <Divider />
                <Detail>
                  <DetailLabel>Purpose</DetailLabel>
                  <Value>{item.purpose}</Value>
                </Detail>
                <Divider />
                <Detail>
                  <DetailLabel>Requested by</DetailLabel>
                  <FundReleaseActor actor={item.requester} />
                </Detail>
                <Divider />
                {/* <Detail>
                  <DetailLabel>Reviewed by</DetailLabel>
                  <FundReleaseActor
                    actor={item.reviewer}
                    fallback={
                      item.reviewed_at || item.reviewed_by_member_id
                        ? "N/A"
                        : "Not yet reviewed"
                    }
                  />
                </Detail> */}
                {/* <Divider /> */}
                <Detail>
                  <DetailLabel>Documents</DetailLabel>
                  <Value>{item.supporting_document_ids?.length ?? 0}</Value>
                </Detail>
                <Divider />
                <Detail>
                  <DetailLabel>Status</DetailLabel>
                  <Status>
                    <StatusDot color={statusColor(item.status)} />
                    {statusLabel(item.status)}
                  </Status>
                </Detail>
                <Divider />
                <Detail>
                  <DetailLabel>Date</DetailLabel>
                  <Value>{formatDate(item.created_at)}</Value>
                </Detail>
              </DetailsRow>
            </ActivityCard>
          ))}
        {!isError && (meta?.total ?? 0) > 0 && (
          <Pagination
            currentPage={meta?.page ?? currentPage}
            totalCount={meta?.total ?? records.length}
            pageSize={meta?.limit ?? pageSize}
            onPageChange={setCurrentPage}
            onPageSizeChange={(size) => {
              setPageSize(size);
              setCurrentPage(1);
            }}
            isFetching={isFetching}
          />
        )}
      </ReleaseContainer>
      {createOpen && (
        <AmFundRealseRequestForm
          form={form}
          setForm={setForm}
          assets={assets}
          assetsLoading={assetsLoading}
          documents={documents}
          documentsLoading={documentsLoading}
          documentFiles={documentFiles}
          setDocumentFiles={setDocumentFiles}
          selectedDocumentIds={selectedDocumentIds}
          setSelectedDocumentIds={setSelectedDocumentIds}
          isCreating={isCreating}
          isUploading={isUploading}
          canCreate={canCreate}
          setCreateOpen={setCreateOpen}
          handleCreate={handleCreate}
        />
      )}
    </Container>
  );
}

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 30px;
`;
const DetailsLink = styled(Link)`
  float: right;
  background: #007cdf;
  border: none;
  padding: 8px 24px;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  font-family: inherit;
  cursor: pointer;
  font-weight: 600;
  text-decoration: none;
`;
const CardContainer = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  width: 100%;
`;
const ReleaseContainer = styled.div`
  background: #fff;
  padding: 24px;
  border-radius: 16px;
`;
const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
`;
const HeaderActions = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;
const Title = styled.h3`
  font-size: 18px;
  font-weight: 600;
  color: #1a2b49;
  margin-bottom: 20px;
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
const ActivityCard = styled.div`
  background: #eef2f6;
  padding: 18px;
  border-radius: 12px;
  margin-bottom: 16px;
`;
const TopRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;
const LeftInfo = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;
const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #007cdf;
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
`;
const ProjectInfo = styled.div`
  display: flex;
  flex-direction: column;
`;
const ProjectName = styled.div`
  font-weight: 600;
  color: #1a2b49;
`;
const ProjectSub = styled.div`
  font-size: 12px;
  color: #6b7280;
`;
const Description = styled.p`
  margin-top: 10px;
  color: #00225a;
  font-weight: 600;
  font-size: 16px;
  line-height: 140%;
`;
const DetailsRow = styled.div`
  margin-top: 14px;
  background: #fff;
  padding: 14px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  overflow-x: auto;
`;
const Detail = styled.div`
  display: flex;
  flex-direction: column;
  min-width: 144px;
`;
const DetailLabel = styled.div`
  font-size: 12px;
  color: #8a94a6;
`;
const Value = styled.div`
  font-size: 14px;
  color: #00225a;
  margin-top: 4px;
  font-weight: 500;
`;
const Divider = styled.div`
  width: 1px;
  min-width: 1px;
  height: 32px;
  background: #e5e7eb;
  margin: 0 20px;
`;
const Status = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 4px;
  font-size: 14px;
  color: #1a2b49;
`;
const StatusDot = styled.div<{ color: string }>`
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: ${(props) => props.color};
`;
const EmptyState = styled.p`
  text-align: center;
  color: #8a94a6;
  padding: 40px 0;
  font-size: 14px;
`;
const InlineRetry = styled.button`
  border: 0;
  background: transparent;
  color: #007cdf;
  cursor: pointer;
  font: inherit;
`;
const ErrorState = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 16px;
  border: 1px solid #fecdca;
  border-radius: 10px;
  background: #fef3f2;
  color: #b42318;
`;
const RetryButton = styled.button`
  border: 1px solid #b42318;
  border-radius: 8px;
  padding: 8px 14px;
  background: #fff;
  color: #b42318;
  font: inherit;
  cursor: pointer;
`;
const FilterContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 15px;
`;
const FilterActions = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 20px;
`;
