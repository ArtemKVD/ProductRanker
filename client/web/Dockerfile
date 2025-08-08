FROM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /web ./client/web/main.go

FROM alpine:3.18
WORKDIR /
COPY --from=builder /web /web
COPY --from=builder /app/client/web/templates /templates 

EXPOSE 8081
CMD ["/web"]