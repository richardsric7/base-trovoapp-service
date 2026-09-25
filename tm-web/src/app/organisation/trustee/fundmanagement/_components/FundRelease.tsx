"use client";
import FundReleaseActor from "@/app/organisation/components/FundReleaseActor";
import styled from "styled-components";
import FundFilter from "./FundFilter";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useGetTrusteeFundReleasesQuery } from "@/redux/api/trustees";
import Pagination from "@/components/CustomPagination";
import {
  IFundReleaseRecord,
  ITrusteeFundReleaseListParams,
} from "@/redux/api/trustees/interface";

const FundRelease = () => {
  const [filters, setFilters] = useState<ITrusteeFundReleaseListParams>({});
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const router = useRouter();
  const { data, isLoading, isFetching, isError } =
    useGetTrusteeFundReleasesQuery({
      ...filters,
      page: currentPage,
      limit: pageSize,
    });

  const records: IFundReleaseRecord[] = data?.data?.records ?? [];
  const meta = data?.data?.meta;

  const formatAmount = (amount: string | number, currency: string) =>
    `${Number(amount).toLocaleString()} ${currency}`;

  const formatDate = (iso: string) =>
    new Date(iso).toLocaleDateString("en-GB", {
      day: "numeric",
      month: "short",
      year: "numeric",
    });

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case "approved":
        return "#22c55e";
      case "rejected":
        return "#ef4444";
      default:
        return "#facc15";
    }
  };

  const getStatusLabel = (status: string) => {
    switch (status.toLowerCase()) {
      case "submitted":
        return "Pending Review";
      case "approved":
        return "Approved";
      case "rejected":
        return "Rejected";
      default:
        return status;
    }
  };

  return (
    <Container>
      <Header>
        <Title>Fund Release Requests</Title>
        <FundFilter
          filters={filters}
          onApply={(nextFilters) => {
            setFilters(nextFilters);
            setCurrentPage(1);
          }}
        />
      </Header>

      {isLoading && <EmptyState>Loading fund release requests…</EmptyState>}
      {isError && (
        <EmptyState style={{ color: "#ef4444" }}>
          Failed to load requests. Please try again.
        </EmptyState>
      )}
      {!isLoading && !isError && records.length === 0 && (
        <EmptyState>No fund release requests found.</EmptyState>
      )}

      {records.map((item) => (
        <ActivityCard key={item.id}>
          <TopRow>
            <LeftInfo>
              <Avatar>{item.asset_code.slice(0, 2).toUpperCase()}</Avatar>
              <ProjectInfo>
                <ProjectName>{item.asset_code}</ProjectName>
                <ProjectSub>{item.purpose}</ProjectSub>
              </ProjectInfo>
            </LeftInfo>

            <ReviewButton
              onClick={() =>
                router.push(`/organisation/trustee/fundmanagement/${item.id}`)
              }
            >
              {item.reviewed_at ? "View" : "Review"}
            </ReviewButton>
          </TopRow>

          <Description>{formatAmount(item.amount, item.currency)}</Description>

          <DetailsRow>
            <Detail>
              <Label>Request ID</Label>
              <Value title={item.id}>{item.id.slice(0, 8)}…</Value>
            </Detail>

            <Divider />

            <Detail>
              <Label>Purpose</Label>
              <Value>{item.purpose}</Value>
            </Detail>

            <Divider />

            <Detail>
              <Label>Requested by</Label>
              <FundReleaseActor actor={item.requester} />
            </Detail>

            <Divider />

            {/* <Detail>
              <Label>Asset</Label>
              <Value>{item.asset_code}</Value>
            </Detail> */}

            {/* <Divider /> */}

            <Detail>
                <Label>Reviewed by</Label>
                <FundReleaseActor actor={item.reviewer} fallback={item.reviewed_at || item.reviewed_by_member_id ? "N/A" : "Not yet reviewed"} />
              </Detail>
              <Divider />
              <Detail>
                <Label>Status</Label>
              <Status>
                <StatusDot color={getStatusColor(item.status)} />
                {getStatusLabel(item.status)}
              </Status>
            </Detail>

            <Divider />

            <Detail>
              <Label>Date</Label>
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
    </Container>
  );
};

export default FundRelease;

const Container = styled.div`
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

const Title = styled.h3`
  font-size: 18px;
  font-weight: 600;
  color: #1a2b49;
  margin-bottom: 20px;
`;

const EmptyState = styled.p`
  text-align: center;
  color: #8a94a6;
  padding: 40px 0;
  font-size: 14px;
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

const ReviewButton = styled.button`
  background: #007cdf;
  border: none;
  padding: 8px 24px;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  cursor: pointer;
  font-weight: 600;
  font-family: inherit;
`;

const Description = styled.p`
  margin-top: 10px;
  color: #00225a;
  font-weight: 600;
  font-size: 16px;
  line-height: 140%;
  letter-spacing: 0.1px;
`;

const DetailsRow = styled.div`
  flex-wrap: wrap;
  row-gap: 16px;
  margin-top: 14px;
  background: #ffffff;
  padding: 14px;
  border-radius: 10px;
  display: flex;
  align-items: center;
`;

const Detail = styled.div`
  display: flex;
  flex-direction: column;
  min-width: 144px;
`;

const Label = styled.div`
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
  background: ${(p) => p.color};
`;
