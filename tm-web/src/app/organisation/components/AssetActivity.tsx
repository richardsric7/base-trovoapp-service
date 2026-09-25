"use client";

import { useGetStakeholderAssetActivitiesQuery } from "@/redux/api/sharedstakeholders";
import styled from "styled-components";
import { useState } from "react";
import Pagination from "@/components/CustomPagination";
import ActorAvatar from "./ActorAvatar";

interface AssetActivityProps {
  assetId: string;
  assetCode?: string;
  assetName?: string;
}

const humanize = (value?: string) =>
  value
    ? value
        .replace(/[._-]+/g, " ")
        .replace(/\b\w/g, (letter) => letter.toUpperCase())
    : "N/A";

const stateText = (value: unknown) =>
  typeof value === "string" || typeof value === "number"
    ? String(value).trim() || undefined
    : undefined;

const AssetActivityContent = ({
  assetId,
  assetCode,
  assetName,
}: AssetActivityProps) => {
  const [{ page, pageSize }, setPagination] = useState({ page: 1, pageSize: 10 });
  const setPage = (nextPage: number) => {
    setPagination((current) =>
      // The custom component also emits the old page after a size change.
      current.pageSize !== pageSize
        ? current
        : { ...current, page: nextPage },
    );
  };
  const { currentData: data, isFetching, isError, refetch } = useGetStakeholderAssetActivitiesQuery(
    { assetId, page, limit: pageSize },
    { skip: !assetId },
  );
  const activities = data?.data.records ?? [];

  return (
    <Container>
      <Title>Recent Activity</Title>

      {isFetching && <StateMessage role="status">Loading activity...</StateMessage>}
      {!isFetching && isError && (
        <StateMessage>
          Unable to load activity. <RetryButton onClick={() => refetch()}>Try again</RetryButton>
        </StateMessage>
      )}
      {!isFetching && !isError && activities.length === 0 && (
        <StateMessage>
          {page === 1 ? "No activity has been recorded for this asset." : "No activity on this page."}
          {page > 1 && <RetryButton onClick={() => setPage(1)}>Return to first page</RetryButton>}
        </StateMessage>
      )}

      {!isFetching && !isError && activities.map((item) => {
        const description =
          stateText(item.metadata?.message) ??
          stateText(item.metadata?.description) ??
          stateText(item.after_state?.description) ??
          humanize(item.action);
        const status =
          stateText(item.after_state?.status) ??
          stateText(item.metadata?.status) ??
          "Not provided";

        return (
          <ActivityCard key={item.id}>
            <TopRow>
              <LeftInfo>
                <Avatar />
                <ProjectInfo>
                  <ProjectName>{assetCode || "Asset"}</ProjectName>
                  <ProjectSub>{assetName || "Tokenized asset"}</ProjectSub>
                </ProjectInfo>
              </LeftInfo>

            </TopRow>

            <Description>{description}</Description>

            <DetailsRow>
              <Detail>
                <Label>Activity</Label>
                <Value>{humanize(item.action)}</Value>
              </Detail>

              <Divider />

              <Detail>
                <Label>Performed by</Label>
                <User>
                  <ActorAvatar name={item.actor_name?.trim() || item.actor_organization} />
                  <ProjectInfo>
                    <span>{item.actor_name?.trim() || humanize(item.actor_role)}</span>
                    {item.actor_organization?.trim() && <ProjectSub>{item.actor_organization}</ProjectSub>}
                  </ProjectInfo>
                </User>
              </Detail>

              <Divider />

              <Detail>
                <Label>Status</Label>
                <Status>
                  {humanize(status)}
                </Status>
              </Detail>

              <Divider />

              <Detail>
                <Label>Date</Label>
                <Value>{new Date(item.created_at).toLocaleDateString()}</Value>
              </Detail>
            </DetailsRow>
          </ActivityCard>
        );
      })}
      {!isFetching && !isError && data && data.data.meta.total > 0 && (
        <Pagination
          currentPage={page}
          pageSize={pageSize}
          totalCount={data.data.meta.total}
          onPageChange={setPage}
          onPageSizeChange={(size) => setPagination({ page: 1, pageSize: size })}
          isFetching={isFetching}
        />
      )}
    </Container>
  );
};

export default function AssetActivity(props: AssetActivityProps) {
  return <AssetActivityContent key={props.assetId} {...props} />;
}

const Container = styled.div`
  background: #fff;
  padding: 24px;
  border-radius: 16px;
`;

const Title = styled.h3`
  font-size: 18px;
  font-weight: 600;
  color: #1a2b49;
  margin-bottom: 20px;
`;

const ActivityCard = styled.div`
  background: #eef2f6;
  padding: 18px;
  border-radius: 12px;
  margin-bottom: 16px;
`;

const StateMessage = styled.div`
  padding: 24px;
  border-radius: 12px;
  background: #eef2f6;
  color: #6b7280;
  text-align: center;
`;

const RetryButton = styled.button`
  border: 0;
  background: transparent;
  color: #007cdf;
  cursor: pointer;
  font: inherit;
  font-weight: 600;
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
  background: #d9b48f;
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

const Description = styled.div`
  margin-top: 10px;
  font-size: 14px;
  color: #00225a;
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
  min-width: 160px;
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

const User = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
  font-size: 14px;
  color: #1a2b49;
`;

const Status = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 4px;
  font-size: 14px;
  color: #1a2b49;
`;
