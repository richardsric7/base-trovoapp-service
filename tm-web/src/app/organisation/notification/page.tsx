"use client";

import {
  useGetStakeholderNotificationsQuery,
  useMarkStakeholderNotificationReadMutation,
} from "@/redux/api/sharedstakeholders";
import { format, formatDistanceToNow } from "date-fns";
import styled, { keyframes } from "styled-components";

const titleCase = (value: string) =>
  value.replace(/[._-]+/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase());

const getTime = (value: string) => {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return { relative: "Date unavailable", absolute: value };
  return {
    relative: formatDistanceToNow(date, { addSuffix: true }),
    absolute: format(date, "dd MMM, yyyy h:mm a"),
  };
};

export default function NotificationPage() {
  const { data, isLoading, isFetching, isError, refetch } =
    useGetStakeholderNotificationsQuery();
  const [markRead, { isLoading: isMarkingRead, originalArgs: markingId }] =
    useMarkStakeholderNotificationReadMutation();
  const records = data?.data.records ?? [];

  const handleMarkRead = async (id: string) => {
    try {
      await markRead(id).unwrap();
    } catch {
      // The query remains unchanged and the user can retry from the same control.
    }
  };

  return (
    <PageContainer>
      <Tabs aria-label="Notification views">
        <ActiveTab type="button">General</ActiveTab>
      </Tabs>

      <Feed aria-busy={isLoading || isFetching}>
        {isLoading ? (
          Array.from({ length: 6 }).map((_, index) => (
            <SkeletonRow key={index} aria-hidden="true">
              <SkeletonAvatar />
              <SkeletonText />
            </SkeletonRow>
          ))
        ) : isError ? (
          <MessageState role="alert">
            <MessageTitle>We couldn&apos;t load notifications</MessageTitle>
            <MessageText>Please check your connection and try again.</MessageText>
            <RetryButton type="button" onClick={() => refetch()}>Try again</RetryButton>
          </MessageState>
        ) : records.length === 0 ? (
          <MessageState>
            <MessageTitle>No notifications yet</MessageTitle>
            <MessageText>New stakeholder notifications will appear here.</MessageText>
          </MessageState>
        ) : (
          records.map((notification) => {
            const time = getTime(notification.created_at);
            const isRead = Boolean(notification.read_at);
            const isCurrentMutation = isMarkingRead && markingId === notification.id;

            return (
              <NotificationItem key={notification.id} $isRead={isRead}>
                <Avatar aria-hidden="true">
                  {(notification.type || notification.title || "N").charAt(0).toUpperCase()}
                </Avatar>
                <Content>
                  <TitleRow>
                    <NotificationTitle>{notification.title}</NotificationTitle>
                    {!isRead && <UnreadDot title="Unread" />}
                  </TitleRow>
                  <NotificationMessage>{notification.message}</NotificationMessage>
                  <Metadata>
                    <span>{titleCase(notification.type || "Notification")}</span>
                    <span aria-hidden="true">•</span>
                    <time dateTime={notification.created_at} title={time.absolute}>
                      {time.relative}
                    </time>
                  </Metadata>
                </Content>
                <ReadButton
                  type="button"
                  disabled={isRead || isCurrentMutation}
                  aria-label={isRead ? "Notification already read" : "Mark notification as read"}
                  title={isRead ? "Read" : "Mark as read"}
                  onClick={() => handleMarkRead(notification.id)}
                >
                  {isCurrentMutation ? "Marking..." : isRead ? "Read" : "Mark as read"}
                </ReadButton>
              </NotificationItem>
            );
          })
        )}
      </Feed>
    </PageContainer>
  );
}

const shimmer = keyframes`
  from { background-position: 100% 0; }
  to { background-position: -100% 0; }
`;

const PageContainer = styled.section`
  padding: 0 24px 40px;
  color: #00225a;
  @media (max-width: 768px) { padding: 0 12px 28px; }
`;
const Tabs = styled.div`
  width: min(130px, 100%);
  padding: 4px;
  margin-bottom: 22px;
  border-radius: 12px;
  background: #ffffff;
`;
const ActiveTab = styled.button`
  width: 100%;
  padding: 8px 18px;
  border: 0;
  border-radius: 9px;
  background: #f2f6f9;
  color: #00225a;
  font: inherit;
  font-size: 13px;
  cursor: default;
`;
const Feed = styled.div`
  padding: 18px;
  border-radius: 22px;
  background: #ffffff;
  @media (max-width: 768px) { padding: 10px; border-radius: 16px; }
`;
const NotificationItem = styled.article<{ $isRead: boolean }>`
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 68px;
  margin-bottom: 8px;
  padding: 11px 14px;
  border-radius: 13px;
  background: ${({ $isRead }) => ($isRead ? "#f8fafc" : "#f2f6f9")};
  opacity: ${({ $isRead }) => ($isRead ? 0.72 : 1)};
  &:last-child { margin-bottom: 0; }
`;
const Avatar = styled.div`
  display: grid;
  flex: 0 0 42px;
  width: 42px;
  height: 42px;
  place-items: center;
  border-radius: 50%;
  background: #f2d3ad;
  color: #9a673c;
  font-weight: 600;
`;
const Content = styled.div`flex: 1; min-width: 0;`;
const TitleRow = styled.div`display: flex; align-items: center; gap: 7px;`;
const NotificationTitle = styled.h2`margin: 0; color: #007cdf; font-size: 14px; font-weight: 500;`;
const UnreadDot = styled.span`width: 7px; height: 7px; flex: 0 0 7px; border-radius: 50%; background: #007cdf;`;
const NotificationMessage = styled.p`margin: 3px 0; color: #43536a; font-size: 12px; line-height: 1.4;`;
const Metadata = styled.div`display: flex; flex-wrap: wrap; gap: 5px; color: #8a94a6; font-size: 11px;`;
const ReadButton = styled.button`
  flex: 0 0 auto;
  min-width: 92px;
  min-height: 32px;
  padding: 7px 12px;
  border: 1px solid #007cdf;
  border-radius: 9px;
  background: #007cdf;
  color: #ffffff;
  font: inherit;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  &:hover:not(:disabled), &:focus-visible { background: #006bbf; }
  &:disabled { border-color: #b8c3cd; background: #b8c3cd; color: #ffffff; cursor: default; }
`;
const MessageState = styled.div`padding: 64px 20px; text-align: center;`;
const MessageTitle = styled.h2`margin: 0 0 8px; color: #00225a; font-size: 17px;`;
const MessageText = styled.p`margin: 0; color: #828282; font-size: 13px;`;
const RetryButton = styled.button`margin-top: 18px; padding: 9px 20px; border: 0; border-radius: 8px; background: #007cdf; color: #fff; font: inherit; cursor: pointer;`;
const SkeletonRow = styled.div`display: flex; align-items: center; gap: 12px; min-height: 68px; margin-bottom: 8px; padding: 11px 14px; border-radius: 13px; background: #f2f6f9;`;
const SkeletonAvatar = styled.div`width: 42px; height: 42px; flex: 0 0 42px; border-radius: 50%; background: linear-gradient(90deg, #e3eaf0 25%, #f2f6f9 50%, #e3eaf0 75%); background-size: 200% 100%; animation: ${shimmer} 1.4s infinite;`;
const SkeletonText = styled.div`width: min(520px, 70%); height: 30px; border-radius: 6px; background: linear-gradient(90deg, #e3eaf0 25%, #f8fafc 50%, #e3eaf0 75%); background-size: 200% 100%; animation: ${shimmer} 1.4s infinite;`;
