FROM golang:1.24.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/server

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/api-gateway .
COPY --from=builder /app/config ./config

ENV TZ=UTC
ENV CONFIG_PATH=/app/config.yml

EXPOSE 8080

CMD ["./api-gateway"]
