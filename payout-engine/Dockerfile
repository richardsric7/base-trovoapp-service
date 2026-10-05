FROM golang:1.26-alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags "-s -w" -o payout-engine .

# small runtime image; no ports: the engine has no HTTP interface
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata && adduser -D -H engine
COPY --from=builder /build/payout-engine /payout-engine
USER engine
ENTRYPOINT ["/payout-engine"]
