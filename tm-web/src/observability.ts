/**
 * Shared error-reporting helpers for Trovo Manager.
 *
 * Reporting is opt-in: with no DSN configured (local development, or before
 * the variable is set in a deployment) every function here is a no-op and the
 * app behaves exactly as it did before. Instrumentation must never be the
 * reason the UI breaks.
 */
import * as Sentry from "@sentry/nextjs";

export const SENTRY_DSN = process.env.NEXT_PUBLIC_SENTRY_DSN ?? "";
export const APP_ENV = process.env.NEXT_PUBLIC_APP_ENV ?? "development";
export const APP_VERSION = process.env.NEXT_PUBLIC_APP_VERSION ?? "unknown";

/** Whether crash reporting is active. */
export const reportingEnabled = () => SENTRY_DSN !== "";

/**
 * Header the admin API returns on every response (see the backend's observe
 * package). Capturing it lets a browser error be traced to the exact backend
 * request that produced it.
 */
export const REQUEST_ID_HEADER = "x-request-id";

/** The last request ID seen, attached to reports for cross-service tracing. */
let lastRequestId: string | null = null;

export const setLastRequestId = (id: string | null | undefined) => {
  if (!id) return;
  lastRequestId = id;
  if (reportingEnabled()) {
    Sentry.setTag("last_request_id", id);
  }
};

export const getLastRequestId = () => lastRequestId;

/**
 * Report a caught error. Use for failures the app handles but that still need
 * investigating - a crash boundary reports on its own.
 */
export const reportError = (
  error: unknown,
  context?: Record<string, unknown>
) => {
  if (!reportingEnabled()) {
    // Still surface it locally so developers are not left guessing.
    // eslint-disable-next-line no-console
    console.error("[observability]", error, context ?? "");
    return;
  }
  Sentry.withScope((scope) => {
    if (context) scope.setContext("details", context);
    if (lastRequestId) scope.setTag("last_request_id", lastRequestId);
    Sentry.captureException(error);
  });
};

/**
 * Record who is using the app, so an error can be traced to the person who hit
 * it. Only an opaque id and role are sent - never an email address or name,
 * which would put personal data in the error store for no diagnostic gain.
 */
export const identifyUser = (
  user: { id?: string | number; role?: string } | null
) => {
  if (!reportingEnabled()) return;
  if (!user?.id) {
    Sentry.setUser(null);
    return;
  }
  Sentry.setUser({ id: String(user.id), role: user.role });
};
