# syntax=docker/dockerfile:1

# ── Stage 1: Build ────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN make build

# ── Stage 2: Runtime ──────────────────────────────────────────────────
FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata i2c-tools

WORKDIR /app

COPY --from=builder /src/build/clawdbot /app/clawdbot

# Create workspace
RUN mkdir -p /root/.clawdbot/workspace/vault/decisions \
             /root/.clawdbot/workspace/vault/lessons \
             /root/.clawdbot/workspace/vault/trades \
             /root/.clawdbot/workspace/vault/research \
             /root/.clawdbot/workspace/vault/inbox

EXPOSE 18790

ENTRYPOINT ["/app/clawdbot"]
CMD ["agent"]
