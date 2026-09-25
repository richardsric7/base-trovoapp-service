"use client";

/**
 * Last-resort boundary. Next.js renders this when an error escapes the root
 * layout, replacing the entire document - which is why it must render its own
 * <html> and <body>, and cannot rely on the app's providers or theme.
 *
 * Without it, such a crash produces a blank white page with nothing in the UI
 * and nothing reported.
 */
import { useEffect } from "react";
import * as Sentry from "@sentry/nextjs";

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    Sentry.captureException(error);
  }, [error]);

  return (
    <html lang="en">
      <body
        style={{
          margin: 0,
          minHeight: "100vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          fontFamily:
            "Montserrat, -apple-system, BlinkMacSystemFont, sans-serif",
          background: "#f7f8fa",
          color: "#00225a",
        }}
      >
        <div style={{ maxWidth: 460, padding: 32, textAlign: "center" }}>
          <h1 style={{ fontSize: 20, fontWeight: 600, margin: "0 0 12px" }}>
            Something went wrong
          </h1>
          <p style={{ fontSize: 14, lineHeight: 1.6, color: "#5b6b8c", margin: "0 0 24px" }}>
            The page could not be displayed. The problem has been reported to
            the team.
          </p>
          {error.digest && (
            <p style={{ fontSize: 12, color: "#8a97b1", margin: "0 0 24px" }}>
              Reference: <code>{error.digest}</code>
            </p>
          )}
          <button
            onClick={() => reset()}
            style={{
              padding: "10px 24px",
              fontSize: 14,
              fontWeight: 500,
              color: "#fff",
              background: "#00225a",
              border: "none",
              borderRadius: 8,
              cursor: "pointer",
            }}
          >
            Try again
          </button>
        </div>
      </body>
    </html>
  );
}
