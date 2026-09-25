"use client";

import styled from "styled-components";
import { IFundReleaseRecord } from "@/redux/api/trustees/interface";

export default function FundReleaseAuthorizationDetails({ record }: { record?: IFundReleaseRecord }) {
  const amount = record
    ? `${Number(record.amount).toLocaleString()} ${record.currency ?? ""}`
    : "—";

  return <>
    <AssetBox>
      <Avatar>{record?.asset_code?.slice(0, 2).toUpperCase() ?? "—"}</Avatar>
      <div><AssetName>{record?.asset_code ?? "—"}</AssetName><SubText>{record?.purpose ?? "—"}</SubText></div>
    </AssetBox>
    <Details>
      <Row><span>Amount</span><strong>{amount}</strong></Row>
      <Row><span>Purpose</span><strong>{record?.purpose ?? "—"}</strong></Row>
      <Row><span>Status</span><strong>{record?.status ?? "—"}</strong></Row>
    </Details>
  </>;
}

const AssetBox=styled.div`background:#f2f6f9;border-radius:12px;padding:12px;display:flex;align-items:center;gap:10px;`;
const Avatar=styled.div`width:36px;height:36px;border-radius:50%;background:#007cdf;color:#fff;font-size:12px;font-weight:700;display:grid;place-items:center;`;
const AssetName=styled.p`margin:0;font-size:14px;font-weight:600;color:#00225a;`;
const SubText=styled.p`margin:3px 0 0;font-size:12px;color:#6b7280;`;
const Details=styled.div`display:flex;flex-direction:column;gap:10px;padding:4px 0;`;
const Row=styled.div`display:flex;justify-content:space-between;gap:20px;font-size:13px;color:#828282;strong{color:#00225a;text-align:right;}`;
