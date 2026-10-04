FROM golang:1.26-alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o quote-service .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates && adduser -D -H quote
COPY --from=builder /build/quote-service /quote-service
USER quote
EXPOSE 8090
ENTRYPOINT ["/quote-service"]
