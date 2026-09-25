"use client";

import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import Tab from "@/components/Tab";
import {
  useAuthorizeTrusteeDistributionMutation,
  useGetTrusteeDistributionHistoryQuery,
  useGetTrusteeDistributionsQuery,
  useRejectTrusteeDistributionMutation,
} from "@/redux/api/trustees";
import { ITrusteeDistribution } from "@/redux/api/trustees/interface";
import { Dropdown, Modal, message } from "antd";
import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { BsThreeDots } from "react-icons/bs";
import StepUpAuthorizationModal from "@/app/organisation/trustee/_components/StepUpAuthorizationModal";
import styled from "styled-components";

type DistributionView = "pending" | "history";

const formatAmount = (amount: string | number, currency: string) => {
  const numericAmount = Number(amount);
  return `${Number.isFinite(numericAmount) ? numericAmount.toLocaleString() : amount} ${currency}`;
};

const formatDate = (date: string) =>
  new Date(date).toLocaleDateString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });

const getStatusTone = (status: string) => {
  const normalized = status.toLowerCase();
  if (["authorized", "approved", "completed", "paid"].includes(normalized)) {
    return { color: "#00875a", background: "#e9f8f1" };
  }
  if (["rejected", "failed"].includes(normalized)) {
    return { color: "#be3800", background: "#fff1eb" };
  }
  return { color: "#9a6700", background: "#fff8db" };
};

const DistributionsPage = () => {
  const router = useRouter();
  const [view, setView] = useState<DistributionView>("pending");
  const [selected, setSelected] = useState<ITrusteeDistribution>();
  const [action, setAction] = useState<"authorize" | "reject">();
  const [reason, setReason] = useState("");
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);

  const pendingQuery = useGetTrusteeDistributionsQuery({
    page: view === "pending" ? page : 1,
    limit: pageSize,
  });
  const historyQuery = useGetTrusteeDistributionHistoryQuery({
    page: view === "history" ? page : 1,
    limit: pageSize,
  });
  const [authorize, authorizeState] = useAuthorizeTrusteeDistributionMutation();
  const [reject, rejectState] = useRejectTrusteeDistributionMutation();

  const activeQuery = view === "pending" ? pendingQuery : historyQuery;
  const records = activeQuery.data?.data?.records ?? [];
  const activeMeta = activeQuery.data?.data?.meta;
  const pendingCount = pendingQuery.data?.data?.meta?.total ?? 0;
  const historyCount = historyQuery.data?.data?.meta?.total ?? 0;
  const isSubmitting = authorizeState.isLoading || rejectState.isLoading;

  const closeModal = () => {
    if (isSubmitting) return;
    setSelected(undefined);
    setAction(undefined);
    setReason("");
  };

  const openAction = (
    record: ITrusteeDistribution,
    nextAction: "authorize" | "reject",
  ) => {
    setSelected(record);
    setAction(nextAction);
  };

  const handleSubmit = async (challengeId?: string) => {
    if (!selected || !action) return;
    if (action === "reject" && !reason.trim()) {
      message.warning(
        "Please provide a reason for rejecting this distribution.",
      );
      return;
    }

    try {
      if (action === "authorize") {
        if (!challengeId) return;
        await authorize({ distributionId: selected.id, challengeId }).unwrap();
        message.success("Distribution authorized successfully.");
      } else {
        await reject({
          distributionId: selected.id,
          payload: { reason: reason.trim() },
        }).unwrap();
        message.success("Distribution rejected successfully.");
      }
      setSelected(undefined);
      setAction(undefined);
      setReason("");
    } catch (error) {
      message.error(`Failed to ${action} distribution. Please try again.`);
      throw error;
    }
  };

  const columns = useMemo(
    () => [
      {
        title: "Asset",
        dataIndex: "asset_code",
        render: (_: string, record: ITrusteeDistribution) => (
          <AssetCell>
            <AssetAvatar>
              {record.asset_code.slice(0, 2).toUpperCase()}
            </AssetAvatar>
            <div>
              <AssetCode>{record.asset_code}</AssetCode>
              <MutedText title={record.asset_id}>
                {record.asset_id.slice(0, 8)}...
              </MutedText>
            </div>
          </AssetCell>
        ),
      },
      {
        title: "Source",
        dataIndex: "source",
      },
      {
        title: "Amount",
        dataIndex: "amount",
        render: (_: string, record: ITrusteeDistribution) => (
          <Amount>{formatAmount(record.amount, record.currency)}</Amount>
        ),
      },
      {
        title: "Scheduled date",
        dataIndex: "scheduled_date",
        render: (value: string) => formatDate(value),
      },
      {
        title: "Status",
        dataIndex: "status",
        render: (value: string) => {
          const tone = getStatusTone(value);
          return (
            <Status $color={tone.color} $background={tone.background}>
              {value.replaceAll("_", " ")}
            </Status>
          );
        },
      },
      ...(view === "pending"
        ? [
            {
              title: "Action",
              dataIndex: "actions",
              render: (_: unknown, record: ITrusteeDistribution) => (
                <ActionMenu onClick={(event) => event.stopPropagation()}>
                  <Dropdown
                    trigger={["click"]}
                    placement="bottomRight"
                    menu={{
                      items: [
                        {
                          key: "authorize",
                          label: "Authorize Distribution",
                          onClick: () => openAction(record, "authorize"),
                        },
                        {
                          key: "reject",
                          label: (
                            <span style={{ color: "#be3800" }}>
                              Reject Distribution
                            </span>
                          ),
                          onClick: () => openAction(record, "reject"),
                        },
                      ],
                    }}
                  >
                    <EllipsisButton aria-label="Distribution actions">
                      <BsThreeDots size={18} />
                    </EllipsisButton>
                  </Dropdown>
                </ActionMenu>
              ),
            },
          ]
        : []),
    ],
    [view],
  );

  return (
    <Page>
      {/* <Header>
        <div>
          <Title>Distributions</Title>
          <Subtitle>Review and authorize distributions to token holders.</Subtitle>
        </div>
        <SummaryCard>
          <SummaryLabel>Awaiting authorization</SummaryLabel>
          <SummaryValue>{pendingCount}</SummaryValue>
        </SummaryCard>
      </Header> */}

      <Tab
        tabs={[
          {
            key: "pending",
            label: `Pending Distributions (${pendingCount})`,
          },
          {
            key: "history",
            label: `Distribution History (${historyCount})`,
          },
        ]}
        currentTab={view}
        setCurrentTab={(tab) => {
          setView(tab as DistributionView);
          setPage(1);
        }}
        tabContainerStyle={{ width: "100%", maxWidth: "400px" }}
        tabContentStyle={{ padding: "24px", minHeight: "380px" }}
      >
        <TableBody>
          {activeQuery.isError && (
            <ErrorState>
              <p>We could not load distributions.</p>
              <RetryButton onClick={() => activeQuery.refetch()}>
                Try again
              </RetryButton>
            </ErrorState>
          )}
          {!activeQuery.isError && (
            <TableWrapper>
              <CustomTable
                columns={columns}
                dataSource={records}
                isLoading={activeQuery.isLoading || activeQuery.isFetching}
                onRowClick={(record) => {
                  sessionStorage.setItem(
                    `trustee-distribution:${record.id}`,
                    JSON.stringify(record),
                  );
                  router.push(
                    `/organisation/trustee/distributions/${record.id}`,
                  );
                }}
              />
            </TableWrapper>
          )}
          {(activeMeta?.total ?? 0) > 0 && (
            <PaginationArea>
              <Pagination
                currentPage={activeMeta?.page ?? page}
                totalCount={activeMeta?.total ?? 0}
                pageSize={activeMeta?.limit ?? pageSize}
                onPageChange={setPage}
                onPageSizeChange={(size) => {
                  setPageSize(size);
                  setPage(1);
                }}
                isFetching={activeQuery.isFetching}
              />
            </PaginationArea>
          )}
        </TableBody>
      </Tab>

      {action === "authorize" && (
        <StepUpAuthorizationModal
          open={Boolean(selected)}
          onClose={closeModal}
          action="distribution.authorize"
          entityType="distribution"
          entityId={selected?.id}
          title="Authorize Distribution"
          description="Confirm the details, then approve this distribution with your linked Trovo Wallet."
          confirmLabel="Authorize Distribution"
          onAuthorize={handleSubmit}
          isAuthorizing={isSubmitting}
        >
          <Details>
            <DetailRow>
              <span>Asset</span>
              <strong>{selected?.asset_code}</strong>
            </DetailRow>
            <DetailRow>
              <span>Amount</span>
              <strong>
                {selected && formatAmount(selected.amount, selected.currency)}
              </strong>
            </DetailRow>
            <DetailRow>
              <span>Scheduled date</span>
              <strong>{selected && formatDate(selected.scheduled_date)}</strong>
            </DetailRow>
          </Details>
        </StepUpAuthorizationModal>
      )}
      <Modal
        open={action === "reject" && Boolean(selected)}
        onCancel={closeModal}
        footer={null}
        centered
        width={500}
      >
        <ModalContent>
          <ModalTitle>
            {action === "authorize"
              ? "Authorize distribution"
              : "Reject distribution"}
          </ModalTitle>
          <ModalCopy>
            {action === "authorize"
              ? "Confirm that you want to authorize this distribution. This will move it to the next stage of processing."
              : "Tell the proposer why this distribution cannot be authorized."}
          </ModalCopy>
          <Details>
            <DetailRow>
              <span>Asset</span>
              <strong>{selected?.asset_code}</strong>
            </DetailRow>
            <DetailRow>
              <span>Source</span>
              <strong>{selected?.source}</strong>
            </DetailRow>
            <DetailRow>
              <span>Amount</span>
              <strong>
                {selected && formatAmount(selected.amount, selected.currency)}
              </strong>
            </DetailRow>
            <DetailRow>
              <span>Scheduled date</span>
              <strong>{selected && formatDate(selected.scheduled_date)}</strong>
            </DetailRow>
          </Details>
          {action === "reject" && (
            <Field>
              <label htmlFor="distribution-reason">Reason</label>
              <TextArea
                id="distribution-reason"
                value={reason}
                maxLength={500}
                onChange={(event) => setReason(event.target.value)}
                placeholder="Enter a reason for rejecting this distribution"
              />
              <CharacterCount>{reason.length}/500</CharacterCount>
            </Field>
          )}
          <ModalActions>
            <CancelButton onClick={closeModal} disabled={isSubmitting}>
              Cancel
            </CancelButton>
            <ConfirmButton
              $danger={action === "reject"}
              onClick={() => handleSubmit()}
              disabled={isSubmitting}
            >
              {isSubmitting
                ? "Please wait..."
                : action === "authorize"
                  ? "Authorize distribution"
                  : "Reject distribution"}
            </ConfirmButton>
          </ModalActions>
        </ModalContent>
      </Modal>
    </Page>
  );
};

export default DistributionsPage;

const Page = styled.main`
  display: flex;
  flex-direction: column;
  gap: 24px;
`;
const Header = styled.header`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
`;
const Title = styled.h1`
  margin: 0;
  color: #00225a;
  font-size: 28px;
  font-weight: 600;
`;
const Subtitle = styled.p`
  margin: 6px 0 0;
  color: #667085;
  font-size: 14px;
`;
const SummaryCard = styled.div`
  min-width: 210px;
  padding: 18px 20px;
  border-radius: 16px;
  background: #fff;
  border: 1px solid #edf0f3;
`;
const SummaryLabel = styled.p`
  margin: 0;
  color: #667085;
  font-size: 12px;
`;
const SummaryValue = styled.p`
  margin: 5px 0 0;
  color: #00225a;
  font-size: 26px;
  font-weight: 600;
`;
const TableBody = styled.section`
  width: 100%;
  min-height: 332px;
  display: flex;
  flex-direction: column;
`;
const TableWrapper = styled.div`
  overflow-x: auto;
  table {
    min-width: 900px;
  }
  th,
  td {
    padding: 14px 12px;
  }
  th:last-child,
  td:last-child {
    text-align: right;
  }
`;
const PaginationArea = styled.div`
  margin-top: auto;
  padding-top: 48px;
`;
const AssetCell = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;
const AssetAvatar = styled.div`
  width: 38px;
  height: 38px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: #007cdf;
  color: #fff;
  font-weight: 700;
`;
const AssetCode = styled.p`
  margin: 0;
  font-weight: 600;
  font-size: 14px;
`;
const MutedText = styled.p`
  margin: 3px 0 0;
  color: #98a2b3;
  font-size: 11px;
`;
const Amount = styled.span`
  font-weight: 600;
`;
const Status = styled.span<{ $color: string; $background: string }>`
  display: inline-block;
  padding: 6px 10px;
  border-radius: 16px;
  color: ${(p) => p.$color};
  background: ${(p) => p.$background};
  text-transform: capitalize;
  font-weight: 600;
`;
const ActionMenu = styled.div`
  display: flex;
  justify-content: flex-end;
`;
const EllipsisButton = styled.button`
  width: 36px;
  height: 32px;
  border: 1px solid #d9e1e8;
  border-radius: 7px;
  background: #fff;
  color: #007cdf;
  display: grid;
  place-items: center;
  cursor: pointer;
`;
const ErrorState = styled.div`
  min-height: 260px;
  display: grid;
  place-content: center;
  justify-items: center;
  color: #667085;
`;
const RetryButton = styled.button`
  border: 0;
  background: #007cdf;
  color: #fff;
  border-radius: 8px;
  padding: 9px 16px;
  cursor: pointer;
`;
const ModalContent = styled.div`
  padding: 8px;
`;
const ModalTitle = styled.h2`
  margin: 0;
  color: #00225a;
  font-size: 20px;
`;
const ModalCopy = styled.p`
  margin: 10px 0 18px;
  color: #667085;
  line-height: 1.5;
`;
const Details = styled.div`
  padding: 14px;
  border-radius: 12px;
  background: #f7f9fb;
`;
const DetailRow = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 20px;
  padding: 6px 0;
  color: #667085;
  strong {
    color: #00225a;
    text-align: right;
  }
`;
const Field = styled.div`
  margin-top: 18px;
  label {
    display: block;
    color: #00225a;
    font-weight: 600;
    margin-bottom: 7px;
  }
`;
const TextArea = styled.textarea`
  width: 100%;
  min-height: 110px;
  resize: vertical;
  padding: 12px;
  border: 1px solid #d0d5dd;
  border-radius: 9px;
  font: inherit;
  color: #00225a;
  outline: none;
  &:focus {
    border-color: #007cdf;
  }
`;
const CharacterCount = styled.p`
  margin: 4px 0 0;
  text-align: right;
  color: #98a2b3;
  font-size: 11px;
`;
const ModalActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 22px;
`;
const ModalButton = styled.button`
  height: 42px;
  padding: 0 16px;
  border-radius: 8px;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
  &:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }
`;
const CancelButton = styled(ModalButton)`
  border: 1px solid #d0d5dd;
  color: #344054;
  background: #fff;
`;
const ConfirmButton = styled(ModalButton)<{ $danger: boolean }>`
  border: 0;
  color: #fff;
  background: ${(p) => (p.$danger ? "#be3800" : "#007cdf")};
`;
