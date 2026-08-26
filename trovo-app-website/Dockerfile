# =============================================================================
# Trovo wallet marketing website (CRA) served by nginx
# =============================================================================

# ---- Build stage -----------------------------------------------------------
FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
# CI=false so react-scripts doesn't treat warnings as errors during CRA build.
RUN CI=false npm run build

# ---- Serve stage -----------------------------------------------------------
FROM nginx:1.27-alpine AS runner

COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
