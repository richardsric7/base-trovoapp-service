"use client";
import React, { useState } from "react";
import CustomTable from "@/components/CustomTable";
import Pagination from "@/components/CustomPagination";
import { IPartnerWallet, useGetPMExchangesQuery, useGetPMWalletsQuery } from "@/redux/api/publicMarkets";
import { Card, date, Input, Mono, Muted, num, Pill, Row, SectionNav, Select, StatCard, Stats, Title } from "../components/ui";

const WalletsPage = () => {
  const [f, setF] = useState({ exchange: "", status: "", search: "" });
  const [page, setPage] = useState(1);
  const { data, isFetching } = useGetPMWalletsQuery({ ...f, page, limit: 20 });
  const { data: ex } = useGetPMExchangesQuery({});
  const st = data?.data?.stats;
  const set = (k: keyof typeof f) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setF({ ...f, [k]: e.target.value });
    setPage(1);
  };
  return (
    <div>
      <SectionNav />
      <Stats>
        <StatCard title="Wallets Provisioned" value={num(st?.provisioned)} sub={st ? `${num(st.deployedWallets)} deployed on-chain` : undefined} />
        <StatCard title="NDPA Consent Recorded" value={st ? `${st.consentPercent}%` : "—"} color="#00A859" bg="#00A8591A" />
        <StatCard title="Tax Data Complete" value={st ? `${st.taxDataPercent}%` : "—"} color="#00A859" bg="#00A8591A" />
        <StatCard title="Rejected Requests (30d)" value={num(st?.rejected30d)} color="#BE3800" bg="#BE38001A" />
      </Stats>
      <Card>
        <Title>Wallet Provisioning</Title>
        <Muted>
          Wallets exchanges opened for their customers: one Safe per customer, owned by the Public Markets signers. Personal data is
          {data?.data?.personalDataVisible ? " shown because you hold VIEW_PUBLIC_MARKETS_PII." : " masked (VIEW_PUBLIC_MARKETS_PII shows it)."}
        </Muted>
        <Row>
          <Input $wide placeholder="External user ref, wallet id or address" value={f.search} onChange={set("search")} />
          <Select value={f.exchange} onChange={set("exchange")}>
            <option value="">All exchanges</option>
            {(ex?.data?.exchanges ?? []).map((e) => (
              <option key={e.serviceLinkId} value={e.serviceLinkId}>
                {e.name}
              </option>
            ))}
          </Select>
          <Select value={f.status} onChange={set("status")}>
            <option value="">All</option>
            <option value="active">Active</option>
            <option value="rejected">Rejected</option>
          </Select>
        </Row>
        <CustomTable
          columns={[
            { title: "External user ref", dataIndex: "externalUserRef", render: (v: string) => <Mono>{v}</Mono> },
            { title: "Exchange", dataIndex: "exchangeName" },
            { title: "Legal name", dataIndex: "legalName" },
            { title: "Tax ID", dataIndex: "taxIdentifier", render: (v: string) => (v ? <Mono>{v}</Mono> : "—") },
            { title: "Residency", dataIndex: "residencyCountry" },
            { title: "Nationality", dataIndex: "nationality" },
            { title: "NDPA consent", dataIndex: "ndpaConsent", render: (v: boolean) => <Pill status={v ? "COMPLETE" : "FAILED"} label={v ? "Recorded" : "Missing"} /> },
            { title: "Wallet", dataIndex: "walletAddress", render: (v: string, w: IPartnerWallet) => (v ? <Mono>{`${w.walletId} · ${v}`}</Mono> : "—") },
            {
              title: "Status",
              dataIndex: "status",
              render: (v: string, w: IPartnerWallet) => <span title={w.rejectionReason}><Pill status={v} label={v === "rejected" ? `Rejected · ${w.rejectionReason}` : undefined} /></span>,
            },
            { title: "Created", dataIndex: "createdAt", render: (v: string) => date(v) },
          ]}
          dataSource={data?.data?.wallets ?? []}
          isLoading={isFetching && !data}
        />
        <Pagination currentPage={page} totalCount={data?.data?.total ?? 0} pageSize={20} onPageChange={setPage} onPageSizeChange={() => undefined} />
      </Card>
    </div>
  );
};

export default WalletsPage;
