FROM golang:1.24.4-bullseye AS builder

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    pkg-config

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /backend ./backend/main.go

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
    librdkafka++1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /
COPY --from=builder /backend /backend
EXPOSE 8080
CMD ["/backend"]
