"use client";
import React, { useState } from "react";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import { showErrorToast, showSuccessToast } from "@/components";
import { IPayout, IPayoutItem, useGetPayoutItemsQuery, usePayoutItemActionMutation } from "@/redux/api/proceedPayouts";
import { amount, Card, date, errorMessage, Input, LinkButton, Mono, Row, Select, StatusPill, SubTitle } from "./ui";

// payout statuses in which admins may change schedule lines (tm-api agrees)
const EDITABLE = ["LOCKED", "APPROVED", "FUNDING_CHECK_REQUESTED", "PAYING", "PAUSED", "COMPLETED_WITH_FAILURES"];

// A payout's schedule: the fee and VAT lines, then the holders by amount,
// with the controls for single holders.
const PayoutListTable = ({ payout }: { payout: IPayout }) => {
  const [status, setStatus] = useState("");
  const [kind, setKind] = useState("");
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(25);
  const { data, isFetching } = useGetPayoutItemsQuery({ id: payout.id, status, kind, search: search.trim(), page, limit });
  const [itemAction] = usePayoutItemActionMutation();
  const editable = EDITABLE.includes(payout.status);
  const canMarkPaid = payout.status === "PAUSED" || payout.status === "COMPLETED_WITH_FAILURES";

  const act = async (item: IPayoutItem, action: "exclude" | "include" | "mark-paid") => {
    let reason: string | undefined;
    let reference: string | undefined;
    if (action === "exclude") {
      reason = window.prompt(`Why leave ${item.username || item.beneficiaryAddress} out of this payout? Its share stays in the payout Safe.`) ?? undefined;
      if (!reason?.trim()) return;
    } else if (action === "mark-paid") {
      reference = window.prompt("How was this holder paid? (e.g. the transaction or transfer reference)") ?? undefined;
      if (!reference?.trim()) return;
    } else if (!window.confirm("Include this holder in the payout again?")) {
      return;
    }
    try {
      await itemAction({ id: payout.id, itemId: item.id, action, reason, reference }).unwrap();
      showSuccessToast(action === "exclude" ? "Holder excluded" : action === "include" ? "Holder included" : "Holder marked paid");
    } catch (err) {
      showErrorToast(errorMessage(err, "The change failed"));
    }
  };

  const columns = [
    {
      title: "Beneficiary",
      dataIndex: "beneficiaryAddress",
      render: (v: string, r: IPayoutItem) => (
        <span>
          {r.kind !== "HOLDER" ? <strong>{r.kind === "FEE" ? "Payout fee" : "VAT on the fee"}</strong> : r.username ? <strong>{r.username}</strong> : null}
          {(r.kind !== "HOLDER" || r.username) && <br />}
          <Mono>{v}</Mono>
        </span>
      ),
    },
    {
      title: "Tokens held",
      dataIndex: "confirmedTokenizedAssetBalance",
      render: (v: number, r: IPayoutItem) => (r.kind === "HOLDER" ? amount(v) : "—"),
    },
    { title: "Amount", dataIndex: "amountToReceive", render: (v: number) => amount(v, payout.payoutAssetCode) },
    { title: "Status", dataIndex: "status", render: (v: string) => <StatusPill status={v} /> },
    {
      title: "Details",
      dataIndex: "reason",
      render: (v: string, r: IPayoutItem) => (
        <small>
          {v}
          {r.txHash && (
            <>
              {v && <br />}
              <Mono>{r.txHash}</Mono>
            </>
          )}
          {r.paidAt && <> · {date(r.paidAt)}</>}
        </small>
      ),
    },
    {
      title: "",
      dataIndex: "actions",
      render: (_: any, r: IPayoutItem) =>
        r.kind !== "HOLDER" ? null : (
          <span>
            {editable && (r.status === "PENDING" || r.status === "FAILED") && (
              <LinkButton $danger onClick={() => act(r, "exclude")}>
                Exclude
              </LinkButton>
            )}
            {editable && r.status === "EXCLUDED" && r.actionBy && <LinkButton onClick={() => act(r, "include")}>Include</LinkButton>}
            {canMarkPaid && (r.status === "FAILED" || r.status === "PENDING") && (
              <LinkButton onClick={() => act(r, "mark-paid")}>Mark paid</LinkButton>
            )}
          </span>
        ),
    },
  ];

  return (
    <Card>
      <SubTitle>Schedule</SubTitle>
      <Row>
        <Select
          value={status}
          onChange={(e) => {
            setStatus(e.target.value);
            setPage(1);
          }}
        >
          <option value="">All statuses</option>
          {["PENDING", "QUEUED", "PAID", "FAILED", "EXCLUDED", "SKIPPED"].map((s) => (
            <option key={s} value={s}>
              {s.charAt(0) + s.slice(1).toLowerCase()}
            </option>
          ))}
        </Select>
        <Select
          value={kind}
          onChange={(e) => {
            setKind(e.target.value);
            setPage(1);
          }}
        >
          <option value="">Holders, fee and VAT</option>
          <option value="HOLDER">Holders</option>
          <option value="FEE">Fee</option>
          <option value="VAT">VAT</option>
        </Select>
        <Input
          $wide
          placeholder="Search address or username"
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setPage(1);
          }}
        />
      </Row>
      <CustomTable columns={columns} dataSource={data?.data?.items ?? []} isLoading={isFetching && !data} />
      <Pagination
        currentPage={page}
        totalCount={data?.data?.total ?? 0}
        pageSize={limit}
        onPageChange={setPage}
        onPageSizeChange={(n) => {
          setLimit(n);
          setPage(1);
        }}
      />
    </Card>
  );
};

export default PayoutListTable;
