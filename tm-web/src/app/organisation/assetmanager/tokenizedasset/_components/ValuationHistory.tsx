"use client";

import styled from "styled-components";
import { AiOutlineCheckCircle, AiOutlineClockCircle } from "react-icons/ai";
import { showErrorToast, showSuccessToast } from "@/components";
import {
  useGetAssetManagerValuationHistoryQuery,
  useRequestIndependentAssetValuationMutation,
} from "@/redux/api/assetManager";

interface ValuationHistoryProps {
  assetId: string;
}

const formatLabel = (value: string) =>
  value.replaceAll("_", " ").replace(/\b\w/g, (letter) => letter.toUpperCase());

const formatDate = (value: string) => {
  const date = new Date(value);
  return Number.isNaN(date.getTime())
    ? "—"
    : date.toLocaleDateString("en-GB", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      });
};

const formatMoney = (value: string | number, currency: string) => {
  const amount = Number(value);
  if (!Number.isFinite(amount)) return `${value} ${currency}`;
  try {
    return new Intl.NumberFormat("en-NG", {
      style: "currency",
      currency,
      maximumFractionDigits: 2,
    }).format(amount);
  } catch {
    return `${amount.toLocaleString()} ${currency}`;
  }
};

const ValuationHistory = ({ assetId }: ValuationHistoryProps) => {
  const { data, isLoading, isFetching, isError, refetch } =
    useGetAssetManagerValuationHistoryQuery(assetId, { skip: !assetId });
  const [requestIndependentValuation, { isLoading: isRequesting }] =
    useRequestIndependentAssetValuationMutation();
  const records = data?.data.records ?? [];

  const handleIndependentRequest = async (valuationId: string) => {
    try {
      await requestIndependentValuation(valuationId).unwrap();
      showSuccessToast("Independent valuation requested successfully");
      await refetch();
    } catch (error: any) {
      showErrorToast(
        error?.data?.message ??
          error?.message ??
          error?.data?.error ??
          error?.error ??
          "Unable to request an independent valuation",
      );
    }
  };

  return (
    <Wrapper>
      <Title>Valuation History</Title>
      <Table aria-busy={isLoading || isFetching}>
        <Header>
          <span>Date</span>
          <span>Value</span>
          <span>Methodology</span>
          <span>Status</span>
          <span>Action</span>
        </Header>

        {isLoading || isFetching ? (
          <State>Loading valuation history...</State>
        ) : isError ? (
          <State>
            Unable to load valuation history.
            <RetryButton type="button" onClick={() => refetch()}>
              Try again
            </RetryButton>
          </State>
        ) : records.length === 0 ? (
          <State>No valuation history found.</State>
        ) : (
          records.map((item) => (
            <Row key={item.id}>
              <DateCell>{formatDate(item.valuation_date)}</DateCell>
              <Value>{formatMoney(item.valuation, item.currency)}</Value>
              <Source>{formatLabel(item.methodology || "Not specified")}</Source>
              <Status $status={item.status}>
                {item.status.toLowerCase() === "pending" ? (
                  <AiOutlineClockCircle size={14} />
                ) : (
                  <AiOutlineCheckCircle size={14} />
                )}
                {formatLabel(item.status)}
              </Status>
              <ActionButton
                type="button"
                disabled={
                  isRequesting || item.status.toLowerCase() !== "submitted"
                }
                onClick={() => handleIndependentRequest(item.id)}
              >
                {item.status.toLowerCase() === "independent_requested"
                  ? "Requested"
                  : "Request Independent"}
              </ActionButton>
            </Row>
          ))
        )}
      </Table>
    </Wrapper>
  );
};

export default ValuationHistory;

const Wrapper = styled.div`margin-top: 32px;`;
const Title = styled.h2`font-size: 20px; font-weight: 600; color: #00225a; margin-bottom: 16px;`;
const Table = styled.div`overflow: hidden; background: #fff; border: 1px solid #e6e6e6; border-radius: 12px;`;
const Header = styled.div`display: grid; grid-template-columns: 1.5fr 1.7fr 1.8fr 1.4fr 1.5fr; gap: 12px; padding: 16px 20px; font-size: 13px; color: #828282; border-bottom: 1px solid #eee;`;
const Row = styled.div`display: grid; grid-template-columns: 1.5fr 1.7fr 1.8fr 1.4fr 1.5fr; gap: 12px; align-items: center; padding: 18px 20px; border-bottom: 1px solid #f1f1f1; &:last-child { border-bottom: 0; }`;
const DateCell = styled.div`font-size: 14px; font-weight: 500; color: #00225a;`;
const Value = styled(DateCell)``;
const Source = styled(DateCell)``;
const Status = styled.div<{ $status: string }>`display: flex; align-items: center; gap: 6px; width: fit-content; padding: 4px 10px; border-radius: 6px; font-size: 12px; color: ${({ $status }) => $status.toLowerCase() === "rejected" ? "#b42318" : $status.toLowerCase() === "pending" ? "#9a6700" : "#1c9e67"}; background: ${({ $status }) => $status.toLowerCase() === "rejected" ? "#fee4e2" : $status.toLowerCase() === "pending" ? "#fff6dc" : "#dff5ea"};`;
const State = styled.div`display: flex; min-height: 120px; align-items: center; justify-content: center; gap: 12px; padding: 24px; color: #828282; font-size: 14px;`;
const RetryButton = styled.button`border: 1px solid #007cdf; border-radius: 7px; padding: 7px 12px; background: #fff; color: #007cdf; cursor: pointer;`;
const ActionButton = styled.button`width: fit-content; border: 1px solid #007cdf; border-radius: 7px; padding: 7px 10px; background: #fff; color: #007cdf; font-size: 12px; font-weight: 500; cursor: pointer; &:hover:not(:disabled) { background: #f0f8ff; } &:disabled { border-color: #d0d5dd; color: #98a2b3; cursor: not-allowed; }`;
