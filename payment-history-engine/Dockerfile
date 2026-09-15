FROM golang:alpine AS builder

# Set necessary environmet  variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN apk add --no-cache ca-certificates build-base runc curl
# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
# APP_VERSION is stamped into the health endpoint so /health reports exactly
# which build is running (CI passes the git sha).
# RUN go build -buildvcs=false -o main .
ARG APP_VERSION=unknown
RUN go build -ldflags "-X trovo-wallet-payment-history-engine/internal/components/health.buildVersion=${APP_VERSION}" -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .


# Build a small image
FROM alpine:latest
RUN apk add --no-cache ca-certificates runc curl

COPY --from=builder /dist/main /



# Command to run
ENTRYPOINT ["/main"]