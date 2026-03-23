FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/bin/flasherp-developer-api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/bin/flasherp-developer-api .

EXPOSE 8081
ENTRYPOINT ["/app/flasherp-developer-api"]
