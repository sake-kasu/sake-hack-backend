FROM golang:1.26.0-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

FROM alpine:3.23

RUN apk add --no-cache \
    ca-certificates \
    tzdata

ENV TZ=Asia/Tokyo

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/bin/sake-hack-server /app/server
COPY --from=builder /app/db/migrations /app/db/migrations

RUN mkdir -p /app/config && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["/app/server"]
