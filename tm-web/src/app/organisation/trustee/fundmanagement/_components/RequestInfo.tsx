"use client";
import styled from "styled-components";
import { IFundReleaseRecord } from "@/redux/api/trustees/interface";
import ActorAvatar from "@/app/organisation/components/ActorAvatar";

interface RequestInfoProps {
  record?: IFundReleaseRecord;
}

const RequestInfo = ({ record }: RequestInfoProps) => {
  const requesterName = record?.requester?.name?.trim() ||
    [record?.requester?.first_name, record?.requester?.last_name]
      .filter(Boolean).join(" ").trim();

  const formatAmount = (amount: string | number, currency: string) =>
    `${Number(amount).toLocaleString()} ${currency}`;

  const formatDate = (iso: string) =>
    new Date(iso).toLocaleDateString("en-GB", {
      day: "numeric",
      month: "short",
      year: "numeric",
    });

  const getStatusColor = (status: string) => {
    switch (status?.toLowerCase()) {
      case "approved":
        return "#22c55e";
      case "rejected":
        return "#ef4444";
      default:
        return "#f59e0b";
    }
  };

  const getStatusLabel = (status: string) => {
    switch (status?.toLowerCase()) {
      case "submitted":
        return "Pending Review";
      case "approved":
        return "Approved";
      case "rejected":
        return "Rejected";
      default:
        return status ?? "—";
    }
  };

  return (
    <>
      <Card>
        <TopSection>
          <LeftTop>
            <AssetIcon>
              {record?.asset_code?.slice(0, 2).toUpperCase() ?? "—"}
            </AssetIcon>
            <div>
              <AssetName>{record?.asset_code ?? "—"}</AssetName>
            </div>
          </LeftTop>
        </TopSection>

        <Amount>
          {record ? formatAmount(record.amount, record.currency) : "—"}
        </Amount>

        <Grid>
          <InfoBlock>
            <Label>Request ID</Label>
            <Value title={record?.id}>{record?.id?.slice(0, 8) ?? "—"}…</Value>
          </InfoBlock>

          <InfoBlock>
            <Label>Purpose</Label>
            <Value>{record?.purpose ?? "—"}</Value>
          </InfoBlock>

          <InfoBlock>
            <Label>Requested by</Label>
            <UserRow>
              <ActorAvatar name={requesterName || record?.requester?.email || record?.requester?.organization_name} />
              <UserDetails>
                <UserName>
                  {requesterName || record?.requester?.email || record?.requester?.organization_name || "—"}
                </UserName>
                <UserEmail>{record?.requester?.email ?? "—"}</UserEmail>
              </UserDetails>
              <LegacyUserName title={record?.requester_member_id}>
                {record?.requester_member_id?.slice(0, 4) ?? "—"}…
              </LegacyUserName>
            </UserRow>
          </InfoBlock>

          <InfoBlock>
            <Label>Status</Label>
            <StatusRow>
              <StatusDot color={getStatusColor(record?.status ?? "")} />
              <StatusText>{getStatusLabel(record?.status ?? "")}</StatusText>
            </StatusRow>
          </InfoBlock>

          <InfoBlock>
            <Label>Requester organization</Label>
            <Value>{record?.requester?.organization_name ?? "—"}</Value>
          </InfoBlock>

          {/* <InfoBlock>
            <Label>Requester email</Label>
            <Value>{record?.requester?.email ?? "—"}</Value>
          </InfoBlock> */}

          <InfoBlock>
            <Label>Asset</Label>
            <Value>{record?.asset_code ?? "—"}</Value>
          </InfoBlock>

          <InfoBlock>
            <Label>Date</Label>
            <Value>
              {record?.created_at ? formatDate(record.created_at) : "—"}
            </Value>
          </InfoBlock>
        </Grid>
      </Card>
    </>
  );
};

export default RequestInfo;

const Card = styled.div``;

const TopSection = styled.div`
  display: flex;
  justify-content: space-between;
`;

const LeftTop = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

const AssetIcon = styled.div`
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

const AssetName = styled.h3`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
`;

const Amount = styled.p`
  font-size: 20px;
  font-weight: 600;
  margin-top: 16px;
  color: #00225a;
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-top: 24px;
`;

const InfoBlock = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
`;

const Label = styled.p`
  font-size: 14px;
  font-weight: 400;
  color: #828282;
  margin: 0;
`;

const Value = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  margin: 8px 0;
  overflow-wrap: break-word;
  word-break: break-word;
  white-space: normal;
`;

const StatusRow = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

const StatusDot = styled.div<{ color: string }>`
  width: 6px;
  height: 6px;
  background: ${(p) => p.color};
  border-radius: 50%;
`;

const StatusText = styled.span`
  font-size: 13px;
  color: #00225a;
`;

const UserRow = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
`;

const UserName = styled.span`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
`;

const UserDetails = styled.div`
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 2px;
`;

const UserEmail = styled.span`
  color: #828282;
  font-size: 12px;
  overflow-wrap: anywhere;
`;

const LegacyUserName = styled(UserName)`
  display: none;
`;
