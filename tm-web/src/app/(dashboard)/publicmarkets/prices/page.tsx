"use client";
import React from "react";
import { useRouter } from "next/navigation";
import CustomTable from "@/components/CustomTable";
import { IExecution, IPriceRow, useGetPMPricesQuery } from "@/redux/api/publicMarkets";
import { Card, date, Mono, Muted, ngn, num, Pill, SectionNav, Title } from "../components/ui";

const age = (m: number) => (m < 0 ? "never" : m < 60 ? `${m}m` : m < 60 * 48 ? `${Math.round(m / 60)}h` : `${Math.round(m / 1440)}d`);

const PricesPage = () => {
  const router = useRouter();
  const { data, isFetching } = useGetPMPricesQuery(undefined, { pollingInterval: 30000 });
  return (
    <div>
      <SectionNav />
      <Card>
        <Title>Price Oracle</Title>
        <Muted>
          Reference prices per asset: the feed during market hours; after the close, the last close (labelled synthetic). Set a manual price on an asset&apos;s page
          when no feed is connected.
        </Muted>
        <div style={{ height: 16 }} />
        <CustomTable
          columns={[
            { title: "Asset", dataIndex: "assetCode", render: (v: string) => <strong>{v}</strong> },
            { title: "Last price", dataIndex: "lastPrice", render: (v: string) => (Number(v) > 0 ? ngn(v) : "—") },
            { title: "Previous close", dataIndex: "previousClose", render: (v: string) => (Number(v) > 0 ? ngn(v) : "—") },
            { title: "Source", dataIndex: "source", render: (v: string) => <Mono>{v || "—"}</Mono> },
            { title: "Market", dataIndex: "marketOpen", render: (v: boolean, p: IPriceRow) => `${p.market} · ${v ? "open" : "closed"}` },
            { title: "Captured", dataIndex: "capturedAt", render: (v: string) => date(v) },
            { title: "Freshness", dataIndex: "fresh", render: (v: boolean, p: IPriceRow) => <Pill status={v ? "Fresh" : "Stale"} label={`${v ? "Fresh" : "Stale"} · ${age(p.ageMinutes)}`} /> },
          ]}
          dataSource={data?.data?.prices ?? []}
          isLoading={isFetching && !data}
          onRowClick={(p: IPriceRow) => router.push(`/publicmarkets/assets/${p.assetId}`)}
        />
      </Card>
      <Card>
        <Title>Dealing Member executions vs reference price</Title>
        <Muted>Every fill against the reference price taken when its net batch was built.</Muted>
        <div style={{ height: 16 }} />
        <CustomTable
          columns={[
            { title: "DM order", dataIndex: "instructionId", render: (v: string) => <Mono>{v}</Mono> },
            { title: "Dealing Member", dataIndex: "dealingMember" },
            { title: "Asset", dataIndex: "assetCode" },
            { title: "Side", dataIndex: "side" },
            { title: "Quantity", dataIndex: "quantity", render: (v: string) => num(v) },
            { title: "Executed", dataIndex: "executedPrice", render: (v: string) => ngn(v) },
            { title: "Reference", dataIndex: "referencePrice", render: (v: string) => ngn(v) },
            {
              title: "Deviation",
              dataIndex: "deviationPercent",
              render: (v: string, x: IExecution) => <span title={date(x.executedAt ?? undefined)}>{Number(v) > 0 ? `+${v}` : v}%</span>,
            },
          ]}
          dataSource={data?.data?.executions ?? []}
        />
      </Card>
    </div>
  );
};

export default PricesPage;
