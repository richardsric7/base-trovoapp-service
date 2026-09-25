"use client";

import ViewOnlyStakeholderDashboard from "../components/ViewOnlyStakeholderDashboard";
import { useGetIssuingHouseDashboardQuery } from "@/redux/api/sharedstakeholders";

export default function IssuingHouseDashboard() {
  return (
    <ViewOnlyStakeholderDashboard
      useDashboardQuery={useGetIssuingHouseDashboardQuery}
    />
  );
}
