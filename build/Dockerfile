# syntax=docker/dockerfile:1
# check=skip=SecretsUsedInArgOrEnv

FROM oven/bun:1-alpine AS frontend-builder

WORKDIR /app

ARG APP_VERSION

COPY front/ ./

RUN echo "export const appVersion = '${APP_VERSION}';" > ./src/assets/version.ts
RUN echo -e "VITE_API_BASE_URL=auto" > ./.env.production

RUN bun install --ignore-scripts
RUN bun run build

FROM golang:1.26-alpine3.24 AS builder

WORKDIR /app

COPY back/ ./

COPY --from=frontend-builder /app/dist ./internal/app/capuchin/web/dist

ARG APP_VERSION

RUN go mod tidy

RUN go build -trimpath -o capuchin -ldflags="-w -X 'main.appVersion=${APP_VERSION}'" ./cmd/capuchin/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.28.0
RUN cp /go/bin/goose .

FROM alpine:3.24

WORKDIR /app

COPY --from=builder /app/capuchin .
COPY --from=builder /app/goose /usr/local/bin/goose
COPY back/migrations ./migrations

RUN apk add --no-cache libc6-compat su-exec

RUN chmod +x /usr/local/bin/goose /app/capuchin

RUN addgroup --system appgroup && adduser --system appuser --ingroup appgroup
RUN mkdir -p /app/data && chown -R appuser:appgroup /app /usr/local/bin/goose

COPY build/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

ENTRYPOINT ["/entrypoint.sh"]
