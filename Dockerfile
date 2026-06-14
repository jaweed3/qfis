# ── Stage 1: Build frontend ────────────────────────────────────
FROM node:22-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json frontend/
RUN cd frontend && npm ci
COPY frontend/ frontend/
RUN cd frontend && npm run build

# ── Stage 2: Build Go backend ─────────────────────────────────
FROM golang:1.26-alpine AS backend-builder
WORKDIR /app
COPY backend/ backend/
COPY --from=frontend-builder /app/frontend/dist backend/frontend/dist
RUN cd backend && CGO_ENABLED=0 go build -o /qfis .

# ── Stage 3: Runtime ──────────────────────────────────────────
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -h /app qfis
WORKDIR /app
COPY --from=backend-builder /qfis .
RUN mkdir -p /app/data && chown -R qfis:qfis /app
USER qfis

EXPOSE 21465

ENV GIN_MODE=release
ENV PORT=21465
ENV QFIS_DB_PATH=/app/data/qfis.db

VOLUME ["/app/data"]

HEALTHCHECK --interval=30s --timeout=5s --retries=3 --start-period=10s \
  CMD wget -qO- http://localhost:21465/api/v1/dashboard/stats || exit 1

CMD ["./qfis"]
