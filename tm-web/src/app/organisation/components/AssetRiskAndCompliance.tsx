"use client";

import styled from "styled-components";
import type { IStakeholderAssetDetailData } from "@/redux/api/sharedstakeholders";

const humanize = (value: string) =>
  value.replace(/[_-]+/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase());

export default function AssetRiskAndCompliance({
  data,
}: {
  data?: IStakeholderAssetDetailData["risk_and_compliance"];
}) {
  const score = data?.risk_assessment_score;
  const maximum = data?.maximum_risk_score;
  const date = data?.last_audit_date ? new Date(data.last_audit_date) : null;
  const auditDate = date && Number.isFinite(date.getTime())
    ? date.toLocaleDateString("en-GB", { day: "numeric", month: "short", year: "numeric" })
    : "Not available";

  return (
    <Wrapper>
      <Title>Risk and Compliance</Title>
      <Card>
        <Item>
          <Label>Risk Assessment Score</Label>
          <RiskRow>
            {data?.risk_level && <Badge>{humanize(data.risk_level)}</Badge>}
            <Value>
              {typeof score === "number" && Number.isFinite(score)
                ? `${score}${typeof maximum === "number" && maximum > 0 ? `/${maximum}` : ""}`
                : "Not assessed"}
            </Value>
          </RiskRow>
        </Item>
        <Item>
          <Label>Compliance Status</Label>
          <Value>{data?.compliance_status ? humanize(data.compliance_status) : "Not available"}</Value>
        </Item>
        <Item>
          <Label>Last Audit</Label>
          <Value>{auditDate}</Value>
        </Item>
      </Card>
    </Wrapper>
  );
}

const Wrapper = styled.div`
  margin-bottom: 32px;
  font-family: inherit;
`;
const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin-bottom: 16px;
`;
const Card = styled.div`
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 20px;
  background: #fff;
  border: 1px solid #e6e6e6;
  border-radius: 12px;
  padding: 20px;
  @media (max-width: 640px) {
    grid-template-columns: 1fr;
  }
`;
const Item = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;
const Label = styled.p`
  margin: 0;
  font-size: 13px;
  color: #828282;
`;
const Value = styled.p`
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  overflow-wrap: anywhere;
`;
const RiskRow = styled.div`
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
`;
const Badge = styled.span`
  background: #f2f6f9;
  color: #00225a;
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 6px;
`;

