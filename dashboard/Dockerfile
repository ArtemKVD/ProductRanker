FROM golang:1.24.4-bullseye AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /dashboard ./dashboard/main.go

FROM debian:bullseye-slim

WORKDIR /
COPY --from=builder /dashboard /dashboard
EXPOSE 8082
CMD ["/dashboard"]