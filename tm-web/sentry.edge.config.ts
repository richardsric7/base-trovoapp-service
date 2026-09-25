// Edge runtime (middleware). Kept minimal - the app has no edge middleware
// today, but Next.js expects this file when the Sentry plugin is enabled.
import * as Sentry from "@sentry/nextjs";

const dsn = process.env.NEXT_PUBLIC_SENTRY_DSN;

if (dsn) {
  Sentry.init({
    dsn,
    environment: process.env.NEXT_PUBLIC_APP_ENV ?? "development",
    release: process.env.NEXT_PUBLIC_APP_VERSION,
    tracesSampleRate: 0,
  });
}
