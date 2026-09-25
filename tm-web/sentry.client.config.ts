// Browser-side error reporting. Loaded automatically by @sentry/nextjs.
//
// With no NEXT_PUBLIC_SENTRY_DSN set, Sentry.init is skipped entirely and the
// SDK stays dormant - local development and any deployment without the
// variable behave exactly as before.
import * as Sentry from "@sentry/nextjs";

const dsn = process.env.NEXT_PUBLIC_SENTRY_DSN;

if (dsn) {
  Sentry.init({
    dsn,
    environment: process.env.NEXT_PUBLIC_APP_ENV ?? "development",
    release: process.env.NEXT_PUBLIC_APP_VERSION,

    // Errors only. GlitchTip implements Sentry's error API, not its
    // performance/tracing API - response-time trends come from Prometheus.
    tracesSampleRate: 0,

    // Session replay would ship a recording of the admin dashboard, including
    // customer records on screen, to the error store. Off deliberately.
    replaysSessionSampleRate: 0,
    replaysOnErrorSampleRate: 0,

    // Noise that is not actionable: browser extensions, cancelled navigations
    // and the network blips every SPA produces would otherwise drown the real
    // crashes.
    ignoreErrors: [
      "ResizeObserver loop limit exceeded",
      "ResizeObserver loop completed with undelivered notifications",
      "Non-Error promise rejection captured",
      "AbortError",
      "Network request failed",
      "Failed to fetch",
      "NetworkError when attempting to fetch resource",
      "Load failed",
    ],
    denyUrls: [
      /extensions\//i,
      /^chrome:\/\//i,
      /^chrome-extension:\/\//i,
      /^moz-extension:\/\//i,
    ],

    beforeSend(event) {
      // Strip anything that could carry a credential before it leaves the
      // browser. Query strings on admin URLs can contain tokens and ids.
      if (event.request?.url) {
        try {
          const url = new URL(event.request.url);
          url.search = "";
          event.request.url = url.toString();
        } catch {
          // A malformed URL is not worth losing the report over.
        }
      }
      if (event.request?.headers) {
        delete event.request.headers.Authorization;
        delete event.request.headers.authorization;
        delete event.request.headers.Cookie;
        delete event.request.headers.cookie;
      }
      return event;
    },
  });
}
