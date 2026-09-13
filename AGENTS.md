# Agent Instructions

## Project Overview

- This is a Go 1.22 API using Gin for HTTP routing and GORM for persistence.
- The application currently uses a single `main` package. Keep additions consistent with the existing flat layout unless a larger structure is required.
- `main.go` loads `.env`, connects to MySQL with GORM, runs `AutoMigrate(Asset{})`, configures CORS, and starts the Gin server.
- `asset.go` contains asset-related GORM types and handlers. Keep database access behind handler or service boundaries as the API grows.

## Development Commands

- Format Go files with `gofmt -w <files>`.
- Run the full validation suite with `go test ./...`.
- Start the API with `go run .`; this requires a reachable MySQL database and a valid local `.env` file.
- Run `go mod tidy` only when dependency changes are intentional, then review both `go.mod` and `go.sum`.

## Implementation Conventions

- Use Gin handlers for HTTP behavior and return JSON responses with appropriate status codes.
- Use GORM models for database entities, preserving explicit tags for primary keys, indexes, nullability, sizes, JSON names, and soft deletion.
- Propagate or handle errors from database operations and server startup; do not silently ignore migration or request failures.
- Read configuration from environment variables rather than hardcoding deployment-specific values.
- Never commit credentials, private keys, or other secrets from `.env`; use a local untracked environment file for development.

## Change Checklist

1. Keep changes focused on the relevant API or persistence surface.
2. Run `gofmt` on modified Go files.
3. Run `go test ./...` before finishing.
4. If dependencies changed, review the resulting module diff and confirm the application still builds.