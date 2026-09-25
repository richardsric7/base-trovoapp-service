# syntax=docker/dockerfile:1.7

# ---- deps ----
FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

# ---- build ----
FROM node:20-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
# NEXT_PUBLIC_* vars are inlined into the client bundle at build time,
# so they must be present here — not at runtime.
ARG NEXT_PUBLIC_API_BASE_URL
ENV NEXT_PUBLIC_API_BASE_URL=${NEXT_PUBLIC_API_BASE_URL}
# Error reporting. All optional: with no DSN the SDK stays dormant and the
# build is unchanged.
ARG NEXT_PUBLIC_SENTRY_DSN
ENV NEXT_PUBLIC_SENTRY_DSN=${NEXT_PUBLIC_SENTRY_DSN}
ARG NEXT_PUBLIC_APP_ENV
ENV NEXT_PUBLIC_APP_ENV=${NEXT_PUBLIC_APP_ENV}
ARG NEXT_PUBLIC_APP_VERSION
ENV NEXT_PUBLIC_APP_VERSION=${NEXT_PUBLIC_APP_VERSION}
# Source-map upload. Only when a token is supplied; otherwise the Sentry build
# plugin is skipped entirely (see next.config.mjs).
ARG SENTRY_AUTH_TOKEN
ENV SENTRY_AUTH_TOKEN=${SENTRY_AUTH_TOKEN}
ARG SENTRY_ORG
ENV SENTRY_ORG=${SENTRY_ORG}
ARG SENTRY_PROJECT
ENV SENTRY_PROJECT=${SENTRY_PROJECT}
ARG SENTRY_URL
ENV SENTRY_URL=${SENTRY_URL}
# Explicit opt-in for source-map upload. Credentials alone do not enable it:
# a wrong org/project makes sentry-cli abort the build, and readable stack
# traces are never worth a failed deployment.
ARG SENTRY_UPLOAD_SOURCEMAPS
ENV SENTRY_UPLOAD_SOURCEMAPS=${SENTRY_UPLOAD_SOURCEMAPS}
ENV NEXT_TELEMETRY_DISABLED=1
# Node's default heap (~2 GB on the CI runner) is not enough to build this app
# once the Sentry webpack plugin instruments every module - the build died with
# "JavaScript heap out of memory". 4 GB leaves headroom on a standard runner.
ENV NODE_OPTIONS=--max-old-space-size=4096
RUN yarn build

# ---- runtime ----
FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    PORT=3000

RUN addgroup -g 1001 -S nodejs && adduser -S nextjs -u 1001

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules

# .next/cache holds webpack's build cache - hundreds of MB of .pack files that
# exist only to speed up a rebuild. Copying them into the runtime image bloats
# every layer, every registry push and every pull, and it filled the CI
# runner's disk during image export ("no space left on device"). Next.js
# recreates what it needs at runtime.
RUN rm -rf /app/.next/cache
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/next.config.mjs ./next.config.mjs

# Next.js writes image optimization cache at runtime into .next/cache.
# Files copied above are root-owned; nextjs user needs write access.
RUN mkdir -p /app/.next/cache/images && chown -R nextjs:nodejs /app/.next

USER nextjs
EXPOSE 3000
CMD ["yarn", "start"]
