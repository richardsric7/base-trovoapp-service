"use client";

import React, { FormEvent, useMemo, useState } from "react";
import styled from "styled-components";
import { FiDollarSign, FiPlus, FiSend, FiX } from "react-icons/fi";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import Loader from "@/components/Loader";
import CustomFilter from "@/components/CustomFilter";
import PrimaryActionButton from "@/components/PrimaryButton";
import SecondaryActionButton from "@/components/SecondaryButton";
import SearchableAssetSelect from "@/app/organisation/assetmanager/components/SearchableAssetSelect";
import OrgTokenizedAssetsStats from "@/app/organisation/assetmanager/tokenizedasset/_components/OrgTokenizedAssetsStats";
import {
  DropdownSelect,
  SearchBar,
  showErrorToast,
  showSuccessToast,
} from "@/components";
import {
  AssetManagerRevenueRecord,
  useGetAssetManagerAssetsQuery,
  useGetAssetManagerRevenueQuery,
  useRecordAssetManagerRevenueMutation,
  useSubmitAssetManagerRevenueDistributionMutation,
} from "@/redux/api/assetManager";

type AssetOption = { id: string; code: string; name: string; currency: string };
type RevenueForm = {
  assetId: string;
  amount: string;
  currency: string;
  source: string;
  periodStart: string;
  periodEnd: string;
  collectedAt: string;
};

const today = new Date().toISOString().slice(0, 10);
const emptyForm: RevenueForm = {
  assetId: "",
  amount: "",
  currency: "",
  source: "",
  periodStart: "",
  periodEnd: "",
  collectedAt: today,
};
const errorMessage = (error: any, fallback: string) =>
  error?.data?.message ?? error?.error ?? fallback;
const label = (value: string) =>
  value.replaceAll("_", " ").replace(/\b\w/g, (c) => c.toUpperCase());
const dateLabel = (value?: string) =>
  value
    ? new Date(value).toLocaleDateString("en-GB", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      })
    : "—";
const money = (amount: string | number, currency: string) => {
  const value = Number(amount);
  if (!Number.isFinite(value)) return `${currency} ${amount}`;
  try {
    return new Intl.NumberFormat("en-NG", {
      style: "currency",
      currency,
      maximumFractionDigits: 2,
    }).format(value);
  } catch {
    return `${currency} ${value.toLocaleString()}`;
  }
};

export default function RevenuePage() {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState("all");
  const [draftStatus, setDraftStatus] = useState("all");
  const [filterOpen, setFilterOpen] = useState(false);
  const [recordOpen, setRecordOpen] = useState(false);
  const [distributionOpen, setDistributionOpen] = useState(false);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [form, setForm] = useState<RevenueForm>(emptyForm);
  const [scheduledDate, setScheduledDate] = useState("");

  const { data, isLoading, isFetching, isError, refetch } =
    useGetAssetManagerRevenueQuery({ page, limit: pageSize });
  const { data: assetsResponse, isLoading: assetsLoading } =
    useGetAssetManagerAssetsQuery({ page: 1, limit: 100 });
  const [recordRevenue, { isLoading: isRecording }] =
    useRecordAssetManagerRevenueMutation();
  const [submitDistribution, { isLoading: isSubmitting }] =
    useSubmitAssetManagerRevenueDistributionMutation();

  const records = data?.data?.records ?? [];
  const meta = data?.data?.meta;
  const assets = useMemo<AssetOption[]>(() => {
    const raw = assetsResponse?.data as any;
    const list = Array.isArray(raw) ? raw : (raw?.records ?? raw?.assets ?? []);
    return list
      .map((asset: any) => ({
        id: String(asset.asset_id ?? asset.assetId ?? asset.id ?? ""),
        code: String(asset.asset_code ?? asset.assetCode ?? ""),
        name: String(asset.asset_name ?? asset.assetName ?? ""),
        currency: String(
          asset.assetQuoteCurrency ??
            asset.asset_quote_currency ??
            asset.currency ??
            "",
        ),
      }))
      .filter((asset: AssetOption) => asset.id);
  }, [assetsResponse]);

  const visibleRecords = useMemo(() => {
    const term = search.trim().toLowerCase();
    return records.filter(
      (record) =>
        (status === "all" || record.status === status) &&
        (!term ||
          record.asset_code?.toLowerCase().includes(term) ||
          record.source.toLowerCase().includes(term) ||
          record.currency.toLowerCase().includes(term)),
    );
  }, [records, search, status]);
  const selected = records.filter((record) => selectedIds.includes(record.id));
  // Rounded to currency minor-unit precision immediately after summing — plain float addition
  // on money values can produce noise like 3050.3049999999998, and this total is submitted
  // verbatim as the real distribution amount a trustee later authorizes.
  const selectedTotal = Number(
    selected
      .reduce((sum, record) => sum + Number(record.amount || 0), 0)
      .toFixed(2),
  );
  const selectionKey = selected[0]
    ? `${selected[0].asset_id}:${selected[0].currency}`
    : "";
  const eligibleVisible = visibleRecords.filter(
    (record) => record.status === "recorded",
  );
  const allEligibleSelected =
    eligibleVisible.length > 0 &&
    eligibleVisible.every((record) => selectedIds.includes(record.id));

  const toggleRecord = (record: AssetManagerRevenueRecord) => {
    if (record.status !== "recorded") return;
    setSelectedIds((current) => {
      if (current.includes(record.id))
        return current.filter((id) => id !== record.id);
      if (
        selectionKey &&
        selectionKey !== `${record.asset_id}:${record.currency}`
      ) {
        showErrorToast("Select records from the same asset and currency");
        return current;
      }
      return [...current, record.id];
    });
  };

  const submitRecord = async (event: FormEvent) => {
    event.preventDefault();
    if (!form.assetId) {
      showErrorToast("Select an asset");
      return;
    }
    if (form.periodEnd < form.periodStart) {
      showErrorToast("Period end must be on or after period start");
      return;
    }
    try {
      await recordRevenue({
        asset_id: form.assetId,
        amount: form.amount,
        currency: form.currency.trim().toUpperCase(),
        source: form.source.trim(),
        period_start: form.periodStart,
        period_end: form.periodEnd,
        collected_at: form.collectedAt,
      }).unwrap();
      showSuccessToast("Revenue recorded successfully");
      setForm(emptyForm);
      setRecordOpen(false);
      setPage(1);
    } catch (error) {
      showErrorToast(errorMessage(error, "Unable to record revenue"));
    }
  };

  const submitProposal = async (event: FormEvent) => {
    event.preventDefault();
    if (!selected[0]) return;
    try {
      await submitDistribution({
        asset_id: selected[0].asset_id,
        revenue_record_ids: selectedIds,
        amount: selectedTotal.toString(),
        currency: selected[0].currency,
        source: Array.from(new Set(selected.map((item) => item.source))).join(
          ", ",
        ),
        scheduled_date: scheduledDate,
      }).unwrap();
      showSuccessToast("Revenue submitted for distribution");
      setSelectedIds([]);
      setScheduledDate("");
      setDistributionOpen(false);
    } catch (error) {
      showErrorToast(errorMessage(error, "Unable to submit distribution"));
    }
  };

  const columns = [
    {
      title: "",
      dataIndex: "select",
      render: (_: unknown, record: AssetManagerRevenueRecord) => (
        <Checkbox
          type="checkbox"
          aria-label={`Select ${record.asset_code || "revenue"}`}
          checked={selectedIds.includes(record.id)}
          disabled={record.status !== "recorded"}
          onClick={(event) => event.stopPropagation()}
          onChange={() => toggleRecord(record)}
        />
      ),
    },
    {
      title: "Asset",
      dataIndex: "asset_code",
      render: (value: string) => value || "—",
    },
    {
      title: "Source",
      dataIndex: "source",
      render: (value: string) => label(value),
    },
    {
      title: "Revenue period",
      dataIndex: "period_start",
      render: (_: unknown, record: AssetManagerRevenueRecord) =>
        `${dateLabel(record.period_start)} – ${dateLabel(record.period_end)}`,
    },
    { title: "Collected", dataIndex: "collected_at", render: dateLabel },
    {
      title: "Amount",
      dataIndex: "amount",
      render: (value: string, record: AssetManagerRevenueRecord) => (
        <Amount>{money(value, record.currency)}</Amount>
      ),
    },
    {
      title: "Status",
      dataIndex: "status",
      render: (value: string) => (
        <Status $status={value}>{label(value)}</Status>
      ),
    },
  ];

  if (isLoading) return <Loader />;
  return (
    <Page>
      <StatsContainer>
        <Header>
          <div>
            <Title>Revenue</Title>
          </div>
          <HeaderActions>
            {selectedIds.length > 0 && (
              <OutlineButton onClick={() => setDistributionOpen(true)}>
                <FiSend /> Submit distribution ({selectedIds.length})
              </OutlineButton>
            )}
            <PrimaryButton onClick={() => setRecordOpen(true)}>
              <FiPlus /> Record revenue
            </PrimaryButton>
          </HeaderActions>
        </Header>
        <StatsContent>
          <OrgTokenizedAssetsStats
            text="Total Revenue Records"
            count={String(meta?.total ?? records.length)}
            dotColor="#007CDF"
          />
          <OrgTokenizedAssetsStats
            text="Awaiting Distribution"
            count={String(
              records.filter((r) => r.status === "recorded").length,
            )}
            dotColor="#FFCC00"
          />
          <OrgTokenizedAssetsStats
            text="Submitted"
            count={String(
              records.filter((r) => r.status === "submitted_for_distribution")
                .length,
            )}
            dotColor="#00A859"
          />
        </StatsContent>
      </StatsContainer>

      <Card>
        <TableHeader>
          <div>
            <Title>Revenue Records</Title>
          </div>
        </TableHeader>
        <Toolbar>
          <SearchBar
            customWidth="320px"
            placeholder="Search asset, source or currency"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          <CustomFilter
            position={{ top: "112px", right: "32px" }}
            open={filterOpen}
            onClose={setFilterOpen}
          >
            <DropdownSelect
              labelText="Status"
              placeholder="All statuses"
              value={
                draftStatus === "all"
                  ? "All statuses"
                  : draftStatus === "recorded"
                    ? "Awaiting distribution"
                    : label(draftStatus)
              }
              options={[
                "All statuses",
                "Awaiting distribution",
                "Submitted",
                "Draft",
              ]}
              onSelect={(_, index) =>
                setDraftStatus(
                  ["all", "recorded", "submitted_for_distribution", "draft"][
                    index
                  ],
                )
              }
            />
            <FilterActions>
              <SecondaryActionButton
                buttonStyle={{ width: "100%", margin: 0 }}
                onClick={() => {
                  setDraftStatus("all");
                  setStatus("all");
                }}
              >
                Reset
              </SecondaryActionButton>
              <PrimaryActionButton
                buttonStyle={{ width: "100%", margin: 0 }}
                onClick={() => {
                  setStatus(draftStatus);
                  setFilterOpen(false);
                }}
              >
                Apply
              </PrimaryActionButton>
            </FilterActions>
          </CustomFilter>
        </Toolbar>
        {eligibleVisible.length > 0 && (
          <SelectionBar>
            <label>
              <Checkbox
                type="checkbox"
                checked={allEligibleSelected}
                onChange={() =>
                  setSelectedIds(
                    allEligibleSelected ? [] : eligibleVisible.map((r) => r.id),
                  )
                }
              />{" "}
              Select all eligible records
            </label>
            {selectedIds.length > 0 && (
              <button onClick={() => setSelectedIds([])}>
                Clear selection
              </button>
            )}
          </SelectionBar>
        )}
        {isError ? (
          <State>
            <strong>Unable to load revenue</strong>
            <span>Please check your connection and try again.</span>
            <OutlineButton onClick={() => refetch()}>Try again</OutlineButton>
          </State>
        ) : records.length === 0 ? (
          <State>
            <SummaryIcon>
              <FiDollarSign />
            </SummaryIcon>
            <strong>No revenue recorded yet</strong>
            <span>Record your first asset revenue collection.</span>
            <PrimaryButton onClick={() => setRecordOpen(true)}>
              <FiPlus /> Record revenue
            </PrimaryButton>
          </State>
        ) : (
          <>
            <TableWrap $fetching={isFetching}>
              <CustomTable
                columns={columns}
                dataSource={visibleRecords}
                onRowClick={toggleRecord}
              />
            </TableWrap>
            <Pagination
              currentPage={meta?.page ?? page}
              totalCount={meta?.total ?? records.length}
              pageSize={meta?.limit ?? pageSize}
              onPageChange={setPage}
              onPageSizeChange={(size) => {
                setPageSize(size);
                setPage(1);
              }}
              isFetching={isFetching}
            />
          </>
        )}
      </Card>

      {recordOpen && (
        <Overlay onMouseDown={() => !isRecording && setRecordOpen(false)}>
          <Modal onMouseDown={(e) => e.stopPropagation()}>
            <ModalHeader>
              <div>
                <ModalTitle>Record revenue</ModalTitle>
                <Subtitle>Add revenue received from a managed asset.</Subtitle>
              </div>
              <IconButton
                onClick={() => setRecordOpen(false)}
                aria-label="Close"
              >
                <FiX />
              </IconButton>
            </ModalHeader>
            <Form onSubmit={submitRecord}>
              <Field>
                <Label htmlFor="revenue-asset">Asset <RequiredMark>*</RequiredMark></Label>
                <SearchableAssetSelect
                  id="revenue-asset"
                  assets={assets}
                  value={form.assetId}
                  disabled={assetsLoading}
                  loading={assetsLoading}
                  onChange={(assetId) => {
                    const asset = assets.find(
                      (item) => item.id === assetId,
                    );
                    setForm({
                      ...form,
                      assetId,
                      currency: asset?.currency || form.currency,
                    });
                  }}
                />
              </Field>
              <Grid>
                <Field>
                  <Label>Amount <RequiredMark>*</RequiredMark></Label>
                  <Input
                    required
                    min="0.01"
                    step="0.01"
                    type="number"
                    placeholder="0.00"
                    value={form.amount}
                    onChange={(e) =>
                      setForm({ ...form, amount: e.target.value })
                    }
                  />
                </Field>
                <Field>
                  <Label>Currency <RequiredMark>*</RequiredMark></Label>
                  <Input
                    required
                    maxLength={12}
                    placeholder="NGN"
                    value={form.currency}
                    onChange={(e) =>
                      setForm({
                        ...form,
                        currency: e.target.value.toUpperCase(),
                      })
                    }
                  />
                </Field>
              </Grid>
              <Field>
                <Label>Revenue source <RequiredMark>*</RequiredMark></Label>
                <Input
                  required
                  placeholder="e.g. Rent"
                  value={form.source}
                  onChange={(e) => setForm({ ...form, source: e.target.value })}
                />
              </Field>
              <Grid>
                <Field>
                  <Label>Period start <RequiredMark>*</RequiredMark></Label>
                  <Input
                    required
                    type="date"
                    value={form.periodStart}
                    onChange={(e) =>
                      setForm({ ...form, periodStart: e.target.value })
                    }
                  />
                </Field>
                <Field>
                  <Label>Period end <RequiredMark>*</RequiredMark></Label>
                  <Input
                    required
                    type="date"
                    min={form.periodStart}
                    value={form.periodEnd}
                    onChange={(e) =>
                      setForm({ ...form, periodEnd: e.target.value })
                    }
                  />
                </Field>
              </Grid>
              <Field>
                <Label>Collection date <RequiredMark>*</RequiredMark></Label>
                <Input
                  required
                  type="date"
                  value={form.collectedAt}
                  onChange={(e) =>
                    setForm({ ...form, collectedAt: e.target.value })
                  }
                />
              </Field>
              <ModalActions>
                <CancelButton
                  type="button"
                  onClick={() => setRecordOpen(false)}
                >
                  Cancel
                </CancelButton>
                <PrimaryButton type="submit" disabled={isRecording}>
                  {isRecording ? "Recording..." : "Record revenue"}
                </PrimaryButton>
              </ModalActions>
            </Form>
          </Modal>
        </Overlay>
      )}

      {distributionOpen && selected[0] && (
        <Overlay
          onMouseDown={() => !isSubmitting && setDistributionOpen(false)}
        >
          <Modal onMouseDown={(e) => e.stopPropagation()}>
            <ModalHeader>
              <div>
                <ModalTitle>Submit distribution</ModalTitle>
                <Subtitle>
                  Review the selected revenue before submitting.
                </Subtitle>
              </div>
              <IconButton
                onClick={() => setDistributionOpen(false)}
                aria-label="Close"
              >
                <FiX />
              </IconButton>
            </ModalHeader>
            <Review>
              <ReviewRow>
                <span>Asset</span>
                <strong>
                  {selected[0].asset_code || selected[0].asset_id}
                </strong>
              </ReviewRow>
              <ReviewRow>
                <span>Revenue records</span>
                <strong>{selected.length}</strong>
              </ReviewRow>
              <ReviewRow>
                <span>Source</span>
                <strong>
                  {Array.from(
                    new Set(selected.map((r) => label(r.source))),
                  ).join(", ")}
                </strong>
              </ReviewRow>
              <ReviewTotal>
                <span>Total distribution</span>
                <strong>{money(selectedTotal, selected[0].currency)}</strong>
              </ReviewTotal>
            </Review>
            <Form onSubmit={submitProposal}>
              <Field>
                <Label>Scheduled distribution date <RequiredMark>*</RequiredMark></Label>
                <Input
                  required
                  type="date"
                  min={today}
                  value={scheduledDate}
                  onChange={(e) => setScheduledDate(e.target.value)}
                />
              </Field>
              <Notice>
                This proposal will be sent to the assigned trustee for
                authorization.
              </Notice>
              <ModalActions>
                <CancelButton
                  type="button"
                  onClick={() => setDistributionOpen(false)}
                >
                  Cancel
                </CancelButton>
                <PrimaryButton type="submit" disabled={isSubmitting}>
                  {isSubmitting ? "Submitting..." : "Submit distribution"}
                </PrimaryButton>
              </ModalActions>
            </Form>
          </Modal>
        </Overlay>
      )}
    </Page>
  );
}

const Page = styled.main`
  padding: 24px;
  min-height: 100vh;
`;
const StatsContainer = styled.section`
  background: #fff;
  border-radius: 24px;
  padding: 20px;
  margin-bottom: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
`;
const Card = styled.section`
  position: relative;
  background: #fff;
  border-radius: 24px;
  padding: 32px;
`;
const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  @media (max-width: 760px) {
    flex-direction: column;
  }
`;
const TableHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
`;
const HeaderActions = styled.div`
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
`;
const Title = styled.h1`
  margin: 0;
  color: #00225a;
  font-size: 24px;
`;
const Subtitle = styled.p`
  margin: 6px 0 0;
  color: #828282;
  font-size: 14px;
`;
const button = `min-height:44px;padding:0 18px;border-radius:8px;font:inherit;font-weight:600;cursor:pointer;display:inline-flex;align-items:center;justify-content:center;gap:8px;`;
const PrimaryButton = styled.button`
  ${button}border:0;
  background: #007cdf;
  color: #fff;
  &:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }
`;
const OutlineButton = styled.button`
  ${button}border:1px solid #007cdf;
  background: #fff;
  color: #007cdf;
`;
const StatsContent = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
  gap: 20px;
  @media (max-width: 760px) {
    flex-direction: column;
  }
`;
const SummaryIcon = styled.div`
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: #eef7ff;
  color: #007cdf;
  display: grid;
  place-items: center;
  font-size: 21px;
`;
const Toolbar = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin: 28px 0 16px;
  @media (max-width: 650px) {
    flex-direction: column;
  }
`;
const FilterActions = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
`;
const field = `height:46px;width:100%;border:1px solid #dfe4ea;border-radius:8px;padding:0 12px;outline:none;font:inherit;color:#00225a;background:#fff;&:focus{border-color:#007cdf;}`;
const Input = styled.input`
  ${field}
`;
const SelectionBar = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f4f9ff;
  color: #00225a;
  padding: 10px 14px;
  border-radius: 8px;
  margin-bottom: 8px;
  font-size: 13px;
  label {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  button {
    border: 0;
    background: none;
    color: #007cdf;
    cursor: pointer;
  }
`;
const Checkbox = styled.input`
  width: 16px;
  height: 16px;
  accent-color: #007cdf;
  cursor: pointer;
  &:disabled {
    cursor: not-allowed;
  }
`;
const TableWrap = styled.div<{ $fetching: boolean }>`
  overflow-x: auto;
  opacity: ${({ $fetching }) => ($fetching ? 0.55 : 1)};
  transition: opacity 0.15s;
  table {
    min-width: 850px;
  }
`;
const Amount = styled.strong`
  white-space: nowrap;
`;
const Status = styled.span<{ $status: string }>`
  display: inline-flex;
  padding: 6px 10px;
  border-radius: 999px;
  white-space: nowrap;
  color: ${({ $status }) =>
    $status === "submitted_for_distribution" ? "#16794b" : "#9a6700"};
  background: ${({ $status }) =>
    $status === "submitted_for_distribution" ? "#e8f8f0" : "#fff6dc"};
`;
const State = styled.div`
  min-height: 300px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: #00225a;
  gap: 10px;
  span {
    color: #828282;
    margin-bottom: 8px;
  }
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
const Modal = styled.div`
  width: min(100%, 560px);
  max-height: calc(100vh - 48px);
  overflow-y: auto;
  background: #fff;
  border-radius: 20px;
  padding: 28px;
  box-shadow: 0 24px 60px rgba(0, 34, 90, 0.2);
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
const IconButton = styled.button`
  border: 0;
  background: none;
  color: #667085;
  font-size: 23px;
  cursor: pointer;
`;
const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 18px;
  margin-top: 24px;
`;
const Grid = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
  @media (max-width: 560px) {
    grid-template-columns: 1fr;
  }
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
const ModalActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 8px;
`;
const CancelButton = styled.button`
  ${button}border:1px solid #d0d5dd;
  background: #fff;
  color: #00225a;
`;
const Review = styled.div`
  margin-top: 24px;
  border: 1px solid #e7edf5;
  border-radius: 12px;
  padding: 16px;
`;
const ReviewRow = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 20px;
  padding: 8px 0;
  color: #667085;
  strong {
    color: #00225a;
    text-align: right;
  }
`;
const ReviewTotal = styled(ReviewRow)`
  border-top: 1px solid #e7edf5;
  margin-top: 8px;
  padding-top: 16px;
  strong {
    font-size: 20px;
  }
`;
const Notice = styled.p`
  margin: 0;
  background: #f4f9ff;
  color: #52677f;
  padding: 12px;
  border-radius: 8px;
  font-size: 13px;
  line-height: 1.5;
`;
