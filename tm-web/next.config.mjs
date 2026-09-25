import { withSentryConfig } from "@sentry/nextjs";

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Only generate browser source maps when they will actually be uploaded;
  // otherwise they would be served to browsers for no benefit.
  productionBrowserSourceMaps: Boolean(process.env.SENTRY_AUTH_TOKEN),
  compiler: {
    styledComponents: true,
  },
  images: {
    domains: ["storage.googleapis.com"],
  },

  webpack: (config) => {
    config.module.rules.push({
      test: /\.pdf$/,
      use: [
        {
          loader: "file-loader",
          options: {
            name: "[name].[ext]",
            publicPath: "/_next/static/files/",
            outputPath: "static/files/",
          },
        },
      ],
    });
    return config;
  },
};

// Sentry's plugin does two jobs: it bundles sentry.client.config.ts into the
// app, and it uploads source maps. The first is needed whenever a DSN is set;
// the second only when an auth token is supplied. So wrap whenever either is
// configured, and let the upload options be inert without a token.
//
// With neither set - local development, and any build that has not configured
// reporting - the config is returned untouched and the build is unchanged.
const withSentry = (config) => {
  if (!process.env.NEXT_PUBLIC_SENTRY_DSN && !process.env.SENTRY_AUTH_TOKEN) {
    return config;
  }
  // Source-map upload is opt-in via SENTRY_UPLOAD_SOURCEMAPS=1, deliberately
  // separate from simply having credentials present.
  //
  // The plugin shells out to sentry-cli, which creates a release before
  // uploading. When the org or project does not exist that throws an
  // unhandled rejection which no plugin option can intercept - it killed the
  // deployment twice. Since upload only makes stack traces prettier, it is
  // never worth risking a deploy for: it stays off unless explicitly asked
  // for, and can then be turned off again by unsetting one variable.
  const uploadSourcemaps =
    process.env.SENTRY_UPLOAD_SOURCEMAPS === "1" &&
    Boolean(
      process.env.SENTRY_AUTH_TOKEN &&
        process.env.SENTRY_ORG &&
        process.env.SENTRY_PROJECT &&
        process.env.SENTRY_URL
    );
  // An auth token is what makes the plugin shell out to sentry-cli for
  // release create/finalize, and those calls abort the whole build when the
  // org or project is wrong - regardless of sourcemaps.disable. Omitting the
  // option is not enough: the plugin falls back to process.env.SENTRY_AUTH_TOKEN
  // (bundler-plugin-core reads it directly), so the variables are deleted from
  // the environment instead. Without a token the plugin only bundles the SDK,
  // which is all that error reporting needs.
  if (!uploadSourcemaps) {
    delete process.env.SENTRY_AUTH_TOKEN;
    delete process.env.SENTRY_ORG;
    delete process.env.SENTRY_PROJECT;
    delete process.env.SENTRY_URL;
  }

  return withSentryConfig(config, {
    org: process.env.SENTRY_ORG,
    project: process.env.SENTRY_PROJECT,
    authToken: process.env.SENTRY_AUTH_TOKEN,
    url: process.env.SENTRY_URL,
    silent: true,
    // Source maps are uploaded then deleted from the output, so the readable
    // sources are available to the error store but never served to browsers.
    sourcemaps: {
      disable: !uploadSourcemaps,
      deleteSourcemapsAfterUpload: true,
    },
    // Creating a release calls the API before any upload and throws an
    // unhandled rejection if the org/project does not exist - which would fail
    // the whole deployment over an optional enhancement. Only attempt it when
    // source-map upload is actually configured.
    release: {
      create: uploadSourcemaps,
      inject: uploadSourcemaps,
    },
    // The app has no edge middleware; skip the warning about it.
    disableLogger: true,
    // Source-map upload is an enhancement, not a requirement: it makes stack
    // traces readable, but the app works without it. A misconfigured org or
    // project, or GlitchTip being briefly unreachable, must never stop a
    // deployment - so upload errors are logged and the build continues.
    errorHandler: (err) => {
      // eslint-disable-next-line no-console
      console.warn(
        "[sentry] source-map upload failed; continuing without it:",
        err?.message ?? err
      );
    },
  });
};

export default withSentry(nextConfig);
