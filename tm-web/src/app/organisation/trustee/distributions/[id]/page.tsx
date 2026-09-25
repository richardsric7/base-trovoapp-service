"use client";

import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import SecondaryButton from "@/components/SecondaryButton";
import {
  useAuthorizeTrusteeDistributionMutation,
  useGetTrusteeDistributionQuery,
  useRejectTrusteeDistributionMutation,
  useDownloadTrusteeDistributionPayoutsMutation,
} from "@/redux/api/trustees";
import { Dropdown, MenuProps, Modal, message } from "antd";
import { useParams, useRouter } from "next/navigation";
import { useMemo, useState } from "react";
import { AiFillClockCircle } from "react-icons/ai";
import { BsThreeDots } from "react-icons/bs";
import { FaArrowLeft, FaCheck, FaMagnifyingGlass } from "react-icons/fa6";
import styled from "styled-components";
import StepUpAuthorizationModal from "@/app/organisation/trustee/_components/StepUpAuthorizationModal";

const money = (amount?: string | number, currency?: string) => {
  const value = Number(amount);
  return `${Number.isFinite(value) ? value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) : (amount ?? "—")} ${currency ?? ""}`;
};

const date = (value?: string) =>
  value
    ? new Date(value).toLocaleDateString("en-GB", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      })
    : "—";

const shortenWallet = (value: string) =>
  value.length > 16 ? `${value.slice(0, 8)}...${value.slice(-6)}` : value;

export default function DistributionDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const id = String(params?.id ?? "");
  const [authorizeOpen, setAuthorizeOpen] = useState(false);
  const [rejectOpen, setRejectOpen] = useState(false);
  const [reason, setReason] = useState("");
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [search, setSearch] = useState("");

  const detail = useGetTrusteeDistributionQuery(id, { skip: !id });
  const [authorize, authorizeState] = useAuthorizeTrusteeDistributionMutation();
  const [reject, rejectState] = useRejectTrusteeDistributionMutation();
  const [downloadPayouts, { isLoading: isDownloading }] =
    useDownloadTrusteeDistributionPayoutsMutation();

  const record = detail.data?.data.distribution;
  const breakdown = detail.data?.data.breakdown;
  const payout = detail.data?.data.payout;
  const isPending = record?.status === "proposed";

  const payouts = useMemo(
    () =>
      (payout?.payouts ?? []).map((item) => ({
        key: item.id,
        name: "Token Holder",
        username: shortenWallet(item.beneficiary_address),
        wallet: shortenWallet(item.beneficiary_address),
        tokens: Number(item.confirmed_token_balance).toLocaleString(),
        amount: money(item.amount, item.currency),
        status: item.cannot_receive_asset
          ? "Unable to receive"
          : item.paid
            ? "Paid"
            : "Pending",
      })),
    [payout?.payouts],
  );

  const filteredPayouts = useMemo(() => {
    const term = search.trim().toLowerCase();
    return payouts.filter(
      (item) =>
        !term ||
        item.name.toLowerCase().includes(term) ||
        item.username.toLowerCase().includes(term) ||
        item.wallet.toLowerCase().includes(term),
    );
  }, [payouts, search]);
  const visiblePayouts = filteredPayouts.slice(
    (page - 1) * pageSize,
    page * pageSize,
  );

  const submitAuthorize = async (challengeId: string) => {
    if (!record) return;
    try {
      await authorize({ distributionId: record.id, challengeId }).unwrap();
      message.success("Distribution authorized successfully.");
      setAuthorizeOpen(false);
    } catch (error) {
      message.error("Unable to authorize this distribution.");
      throw error;
    }
  };

  const submitReject = async () => {
    if (!record || !reason.trim()) {
      message.warning("Please provide a rejection reason.");
      return;
    }
    try {
      await reject({
        distributionId: record.id,
        payload: { reason: reason.trim() },
      }).unwrap();
      message.success("Distribution rejected successfully.");
      setRejectOpen(false);
    } catch {
      message.error("Unable to reject this distribution.");
    }
  };

  const handleDownloadPayouts = async () => {
    if (!record || isDownloading) return;
    try {
      const blob = await downloadPayouts(record.id).unwrap();
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `distribution-payouts-${record.id.replace(/[^a-zA-Z0-9_-]/g, "_")}.csv`;
      document.body.appendChild(link);
      try {
        link.click();
      } finally {
        link.remove();
        window.setTimeout(() => URL.revokeObjectURL(url), 1000);
      }
    } catch {
      message.error("Unable to download the payout list. Please try again.");
    }
  };

  const menuItems: MenuProps["items"] = [
    {
      key: "download",
      label: isDownloading ? "Downloading..." : "Download Payout List",
      disabled: isDownloading,
      onClick: handleDownloadPayouts,
    },
    ...(isPending
      ? [
          {
            key: "reject",
            label: <span style={{ color: "#be3800" }}>Reject Request</span>,
            onClick: () => setRejectOpen(true),
          },
        ]
      : []),
  ];

  const columns = [
    {
      title: "Investor",
      dataIndex: "name",
      render: (_: string, item: (typeof payouts)[number]) => (
        <Investor>
          <InvestorAvatar>{item.name.charAt(0)}</InvestorAvatar>
          <div>
            <InvestorName>{item.name}</InvestorName>
            <InvestorHandle>{item.username}</InvestorHandle>
          </div>
        </Investor>
      ),
    },
    { title: "Wallet / Bank Details", dataIndex: "wallet" },
    { title: "Tokens Held", dataIndex: "tokens" },
    { title: "Amount", dataIndex: "amount" },
    {
      title: "Status",
      dataIndex: "status",
      render: (status: string) => (
        <PayoutStatus>
          <AiFillClockCircle />
          {status}
        </PayoutStatus>
      ),
    },
  ];

  if (detail.isLoading) {
    return <StateCard>Loading distribution details...</StateCard>;
  }
  if (detail.isError || !record || !breakdown) {
    return <StateCard>Distribution not found.</StateCard>;
  }

  const steps = [
    {
      label: "Submitted",
      complete: true,
      detail: `By asset manager · ${date(record.created_at)}`,
    },
    {
      label: "Request Authorized",
      complete: ["authorized", "processing", "completed"].includes(
        record.status,
      ),
      detail: "By asset trustee",
    },
    {
      label: "Funds Released",
      complete: ["processing", "completed"].includes(record.status),
    },
    {
      label: "Processing",
      complete: ["processing", "completed"].includes(record.status),
    },
    { label: "Completed", complete: record.status === "completed" },
    { label: "Confirmed", complete: false },
  ];

  return (
    <Page>
      <BackButton onClick={() => router.back()}>
        <FaArrowLeft /> Back
      </BackButton>
      <Layout>
        <TimelineCard>
          {steps.map((step, index) => (
            <TimelineItem key={step.label}>
              <TimelineRail>
                <StepDot $complete={step.complete}>
                  {step.complete && <FaCheck size={9} />}
                </StepDot>
                {index < steps.length - 1 && <StepLine />}
              </TimelineRail>
              <StepContent>
                <StepLabel $complete={step.complete}>{step.label}</StepLabel>
                {step.detail && <StepDetail>{step.detail}</StepDetail>}
              </StepContent>
            </TimelineItem>
          ))}
        </TimelineCard>

        <Main>
          <DetailsCard>
            <CardHeader>
              <SectionTitle>Fund Distribution Request</SectionTitle>
              <Actions>
                {isPending && (
                  <PrimaryAction onClick={() => setAuthorizeOpen(true)}>
                    Authorize Request
                  </PrimaryAction>
                )}

                <Dropdown
                  menu={{ items: menuItems }}
                  trigger={["click"]}
                  placement="bottomRight"
                >
                  <MoreButton aria-label="More distribution actions">
                    <BsThreeDots />
                  </MoreButton>
                </Dropdown>
              </Actions>
            </CardHeader>

            <RequestPanel>
              <AssetIdentity>
                <AssetAvatar>{record.asset_code.slice(0, 2)}</AssetAvatar>
                <div>
                  <AssetName>{record.asset_code}</AssetName>
                  <AssetSub>Tokenized asset</AssetSub>
                </div>
              </AssetIdentity>
              <HeroAmount>{money(record.amount, record.currency)}</HeroAmount>
              <InfoGrid>
                <Info>
                  <InfoLabel>Request ID</InfoLabel>
                  <InfoValue title={record.id}>
                    {record.id.slice(0, 8)}...
                  </InfoValue>
                </Info>
                <Info>
                  <InfoLabel>Purpose</InfoLabel>
                  <InfoValue>{record.source}</InfoValue>
                </Info>
                <Info>
                  <InfoLabel>Status</InfoLabel>
                  <StatusValue>
                    <StatusDot />
                    {record.status.replaceAll("_", " ")}
                  </StatusValue>
                </Info>
                <Info>
                  <InfoLabel>Requested by</InfoLabel>
                  <Requester>
                    <SmallAvatar>AM</SmallAvatar>
                    <InfoValue>Asset Manager</InfoValue>
                  </Requester>
                </Info>
                <Info>
                  <InfoLabel>Proposed Date</InfoLabel>
                  <InfoValue>{date(record.created_at)}</InfoValue>
                </Info>
              </InfoGrid>
            </RequestPanel>

            <BreakdownTitle>Distribution Breakdown</BreakdownTitle>
            <BreakdownCard>
              <BreakdownStats>
                <Info>
                  <InfoLabel>Total Tokens</InfoLabel>
                  <InfoValue>
                    {Number(breakdown.total_tokens).toLocaleString()}
                  </InfoValue>
                </Info>
                <Info>
                  <InfoLabel>Token Holders</InfoLabel>
                  <InfoValue>
                    {breakdown.total_token_holders.toLocaleString()}
                  </InfoValue>
                </Info>
                <Info>
                  <InfoLabel>Distribution per Token</InfoLabel>
                  <InfoValue>
                    {money(
                      breakdown.distribution_amount_per_token.amount,
                      breakdown.distribution_amount_per_token.currency,
                    )}
                  </InfoValue>
                </Info>
              </BreakdownStats>
              <Formula>
                <FormulaRow>
                  <span>Net Income:</span>
                  <strong>
                    {money(
                      breakdown.net_income.amount,
                      breakdown.net_income.currency,
                    )}
                  </strong>
                </FormulaRow>
                <FormulaRow>
                  <span>÷ Total Tokens:</span>
                  <strong>
                    {Number(breakdown.total_tokens).toLocaleString()}
                  </strong>
                </FormulaRow>
                <FormulaDivider />
                <FormulaRow>
                  <span>= Distribution per Token:</span>
                  <strong>
                    {money(
                      breakdown.distribution_amount_per_token.amount,
                      breakdown.distribution_amount_per_token.currency,
                    )}
                  </strong>
                </FormulaRow>
              </Formula>
            </BreakdownCard>
          </DetailsCard>

          <PayoutCard>
            <SectionTitle>Payout List</SectionTitle>
            <SearchBox>
              <input
                value={search}
                onChange={(event) => {
                  setSearch(event.target.value);
                  setPage(1);
                }}
                placeholder="Start search"
              />
              <FaMagnifyingGlass />
            </SearchBox>
            <TableWrap>
              <CustomTable columns={columns} dataSource={visiblePayouts} />
            </TableWrap>
            <PaginationWrap>
              <Pagination
                currentPage={page}
                totalCount={filteredPayouts.length}
                pageSize={pageSize}
                onPageChange={setPage}
                onPageSizeChange={(size) => {
                  setPageSize(size);
                  setPage(1);
                }}
              />
            </PaginationWrap>
          </PayoutCard>
        </Main>
      </Layout>

      <StepUpAuthorizationModal
        open={authorizeOpen}
        onClose={() => setAuthorizeOpen(false)}
        action="distribution.authorize"
        entityType="distribution"
        entityId={record.id}
        title="Authorize Distribution"
        description="Confirm the details, then approve this distribution with your linked Trovo Wallet."
        confirmLabel="Authorize Distribution"
        onAuthorize={submitAuthorize}
        isAuthorizing={authorizeState.isLoading}
      >
        <AssetBanner>
          <AssetAvatar>{record.asset_code.slice(0, 2)}</AssetAvatar>
          <div>
            <AssetName>{record.asset_code}</AssetName>
            <AssetSub>Tokenized asset</AssetSub>
          </div>
        </AssetBanner>
        <ModalDetails>
          <FormulaRow>
            <span>Amount</span>
            <strong>{money(record.amount, record.currency)}</strong>
          </FormulaRow>
          <FormulaRow>
            <span>Purpose</span>
            <strong>{record.source}</strong>
          </FormulaRow>
          <FormulaRow>
            <span>Proposed Date</span>
            <strong>{date(record.created_at)}</strong>
          </FormulaRow>
          <FormulaRow>
            <span>Requested by</span>
            <strong>Asset Manager</strong>
          </FormulaRow>
        </ModalDetails>
      </StepUpAuthorizationModal>

      <Modal
        open={rejectOpen}
        onCancel={() => setRejectOpen(false)}
        footer={null}
        centered
        width={500}
      >
        <ModalBody>
          <ModalTitle>Reject Request</ModalTitle>
          <ModalText>
            Provide a reason for rejecting this distribution request.
          </ModalText>
          <ReasonArea
            value={reason}
            onChange={(event) => setReason(event.target.value)}
            placeholder="Enter rejection reason"
          />
          <DangerAction disabled={rejectState.isLoading} onClick={submitReject}>
            {rejectState.isLoading ? "Rejecting..." : "Reject Request"}
          </DangerAction>
        </ModalBody>
      </Modal>
    </Page>
  );
}

const Page = styled.main`
  padding: 0 20px 40px;
  color: #00225a;
`;
const BackButton = styled.button`
  border: 0;
  background: transparent;
  color: #00225a;
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  margin-bottom: 18px;
  font: inherit;
`;
const Layout = styled.div`
  display: grid;
  grid-template-columns: 210px minmax(0, 1fr);
  gap: 20px;
  align-items: start;
  @media (max-width: 1000px) {
    grid-template-columns: 1fr;
  }
`;
const TimelineCard = styled.aside`
  background: #fff;
  border-radius: 16px;
  padding: 24px 20px;
`;
const TimelineItem = styled.div`
  display: flex;
  gap: 10px;
  min-height: 70px;
`;
const TimelineRail = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;
const StepDot = styled.div<{ $complete: boolean }>`
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: ${(p) => (p.$complete ? "#21b573" : "#d8dde3")};
  color: #fff;
`;
const StepLine = styled.div`
  width: 2px;
  flex: 1;
  background: #e2e6ea;
  margin: 5px 0;
`;
const StepContent = styled.div`
  padding-top: 1px;
`;
const StepLabel = styled.p<{ $complete: boolean }>`
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  color: ${(p) => (p.$complete ? "#00225a" : "#98a2b3")};
`;
const StepDetail = styled.p`
  margin: 4px 0;
  color: #98a2b3;
  font-size: 10px;
  line-height: 1.5;
`;
const Main = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
  min-width: 0;
`;
const Card = styled.section`
  background: #fff;
  border-radius: 18px;
  padding: 24px;
`;
const DetailsCard = styled(Card)``;
const PayoutCard = styled(Card)``;
const CardHeader = styled.div`
  display: flex;
  // flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
`;
const SectionTitle = styled.h2`
  margin: 0;
  font-size: 20px;
  font-weight: 600;
`;
const Actions = styled.div`
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
`;
const PrimaryAction = styled.button`
  border: 0;
  background: #007cdf;
  color: #fff;
  border-radius: 8px;
  padding: 11px 18px;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
  &:disabled {
    opacity: 0.55;
  }
`;
const DangerAction = styled(PrimaryAction)`
  background: #be3800;
  width: 100%;
`;
const MoreButton = styled.button`
  width: 42px;
  height: 42px;
  border: 1px solid #007cdf;
  color: #007cdf;
  background: #fff;
  border-radius: 8px;
  display: grid;
  place-items: center;
  cursor: pointer;
`;
const RequestPanel = styled.div`
  border: 1px solid #e4e8ed;
  border-radius: 10px;
  padding: 16px 18px;
`;
const AssetIdentity = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;
const AssetBanner = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  background: #f1f5f8;
  border-radius: 10px;
  padding: 12px;
`;
const AssetAvatar = styled.div`
  width: 38px;
  height: 38px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: #ecd5ad;
  color: #855f27;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
`;
const AssetName = styled.p`
  margin: 0;
  font-weight: 600;
  font-size: 14px;
`;
const AssetSub = styled.p`
  margin: 3px 0 0;
  color: #8b95a5;
  font-size: 10px;
`;
const HeroAmount = styled.h3`
  font-size: 18px;
  margin: 18px 0 16px;
`;
const InfoGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 22px 32px;
  @media (max-width: 850px) {
    grid-template-columns: repeat(2, 1fr);
  }
`;
const Info = styled.div``;
const InfoLabel = styled.p`
  margin: 0 0 5px;
  color: #98a2b3;
  font-size: 11px;
`;
const InfoValue = styled.p`
  margin: 0;
  font-size: 13px;
  font-weight: 600;
`;
const StatusValue = styled(InfoValue)`
  display: flex;
  align-items: center;
  gap: 6px;
  text-transform: capitalize;
`;
const StatusDot = styled.span`
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #f5b800;
`;
const Requester = styled.div`
  display: flex;
  align-items: center;
  gap: 7px;
`;
const SmallAvatar = styled.div`
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: #007cdf;
  color: #fff;
  font-size: 8px;
`;
const BreakdownTitle = styled.h3`
  font-size: 16px;
  margin: 26px 0 12px;
`;
const BreakdownCard = styled.div`
  border: 1px solid #e4e8ed;
  border-radius: 10px;
  padding: 16px 18px;
`;
const BreakdownStats = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-bottom: 14px;
`;
const Formula = styled.div`
  background: #f0f5f9;
  border-radius: 7px;
  padding: 12px 14px;
`;
const FormulaRow = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 24px;
  padding: 5px 0;
  color: #667085;
  font-size: 12px;
  strong {
    color: #00225a;
    text-align: right;
  }
`;
const FormulaDivider = styled.div`
  height: 1px;
  background: #cad8e4;
  margin: 5px 0;
`;
const SearchBox = styled.div`
  margin: 18px 0;
  width: 280px;
  height: 38px;
  border: 1px solid #dfe4ea;
  border-radius: 8px;
  display: flex;
  align-items: center;
  padding: 0 10px;
  color: #007cdf;
  input {
    border: 0;
    outline: 0;
    flex: 1;
    font: inherit;
    font-size: 12px;
  }
`;
const TableWrap = styled.div`
  overflow-x: auto;
  table {
    min-width: 760px;
  }
  th,
  td {
    padding: 13px 8px;
  }
`;
const Investor = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;
const InvestorAvatar = styled.div`
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: #007cdf;
  color: #fff;
`;
const InvestorName = styled.p`
  margin: 0;
  font-size: 12px;
  font-weight: 600;
`;
const InvestorHandle = styled.p`
  margin: 2px 0 0;
  color: #98a2b3;
  font-size: 9px;
`;
const PayoutStatus = styled.span`
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: #007cdf;
  background: #edf6fd;
  padding: 5px 8px;
  border-radius: 12px;
`;
const PaginationWrap = styled.div`
  padding-top: 36px;
`;
const ModalBody = styled.div`
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 16px;
`;
const ModalTitle = styled.h2`
  margin: 0;
  font-size: 20px;
`;
const ModalText = styled.p`
  margin: 0;
  color: #667085;
  line-height: 1.5;
  font-size: 13px;
`;
const ModalDetails = styled.div`
  padding: 8px 0;
`;
const ReasonArea = styled.textarea`
  min-height: 120px;
  border: 1px solid #d0d5dd;
  border-radius: 8px;
  padding: 12px;
  font: inherit;
  resize: vertical;
`;
const StateCard = styled.div`
  background: #fff;
  border-radius: 16px;
  padding: 60px;
  text-align: center;
  color: #667085;
`;
