"use client";

import { useEffect } from "react";
import { scheduleTokenRefresh } from "@/redux/scheduleTokenRefresh";

export default function TokenRefresher() {
  useEffect(() => {
    scheduleTokenRefresh();
  }, []);

  return null; // No UI needed
}
