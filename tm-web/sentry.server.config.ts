// Server-side (Node) error reporting for the Next.js server: route handlers,
// server components and SSR. Uses the same DSN as the browser bundle.
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
