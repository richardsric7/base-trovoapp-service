"use client";

/**
 * Route error boundary. React unmounts the crashed subtree and renders this
 * instead, so a failure here no longer blanks the whole app.
 */
import { useEffect } from "react";
import * as Sentry from "@sentry/nextjs";
import { ErrorState } from "@/components";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    Sentry.captureException(error);
  }, [error]);

  return <ErrorState error={error} reset={reset} what="this page" />;
}
