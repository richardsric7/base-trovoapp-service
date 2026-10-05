"use client";
import CustomTable from "@/components/CustomTable";
import { useGetPMHealthQuery } from "@/redux/api/publicMarkets";
import { Card, date, Muted, Pill, SubTitle } from "../../publicmarkets/components/ui";

// The Public Markets engine's background jobs (they run inside app-backend).
const PublicMarketsJobs = () => {
  const { data, isLoading, isError } = useGetPMHealthQuery(undefined, { pollingInterval: 30000 });
  const h = data?.data;
  if (isError) return null; // not a Trovo admin, or Public Markets unavailable
  return (
    <Card style={{ marginTop: 20 }}>
      <SubTitle>Public Markets background jobs</SubTitle>
      <Muted>
        Queues: {h?.pendingEvents ?? 0} partner event(s) to process, {h?.pendingWebhooks ?? 0} webhook(s) to deliver, {h?.pendingMockEvents ?? 0} mock partner
        callback(s) due.
      </Muted>
      <div style={{ height: 12 }} />
      <CustomTable
        columns={[
          { title: "Job", dataIndex: "job", render: (v: string) => <strong>{v}</strong> },
          { title: "Last run", dataIndex: "lastRunAt", render: (v: string) => date(v) },
          { title: "Detail", dataIndex: "detail" },
          { title: "Status", dataIndex: "status", render: (v: string) => <Pill status={v} /> },
        ]}
        dataSource={h?.jobs ?? []}
        isLoading={isLoading}
      />
    </Card>
  );
};

export default PublicMarketsJobs;
