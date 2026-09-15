# syntax=docker/dockerfile:1.7

FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json* yarn.lock* ./
RUN if [ -f yarn.lock ]; then yarn install --frozen-lockfile; \
    else npm ci; fi
COPY . .
RUN if [ -f yarn.lock ]; then yarn build; else npm run build; fi

FROM nginx:1.27-alpine
RUN apk add --no-cache gettext
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf.template /etc/nginx/templates/default.conf.template
EXPOSE 80
