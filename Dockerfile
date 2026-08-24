FROM node:24-bookworm-slim AS frontend

WORKDIR /src/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
COPY api/ ../api/
RUN npm run build

FROM golang:1.26-bookworm AS backend

WORKDIR /src/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM debian:bookworm-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app/backend
COPY --from=backend /out/server /app/server
COPY backend/migrations ./migrations
COPY --from=frontend /src/frontend/dist /app/frontend/dist

ENV ADDRESS=0.0.0.0:8080
ENV FRONTEND_DIST=/app/frontend/dist

EXPOSE 8080

CMD ["/app/server", "run"]
