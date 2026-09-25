"use client";

import ViewOnlyStakeholderDashboard from "../components/ViewOnlyStakeholderDashboard";
import { useGetRatingAgencyDashboardQuery } from "@/redux/api/sharedstakeholders";

export default function RatingAgencyDashboard() {
  return (
    <ViewOnlyStakeholderDashboard
      useDashboardQuery={useGetRatingAgencyDashboardQuery}
    />
  );
}
