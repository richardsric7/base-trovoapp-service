# =============================================================================
# Trovotech marketing website — Next.js static export served by nginx
# =============================================================================

# ---- Build stage -----------------------------------------------------------
FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# ---- Serve stage -----------------------------------------------------------
FROM nginx:1.27-alpine AS runner

# Static site: copy the exported build. next.config.ts sets output: 'export',
# so `next build` writes the site to /app/out.
COPY --from=builder /app/out /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
