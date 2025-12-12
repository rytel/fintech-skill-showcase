# Copilot / AI Agent Instructions for fintech-skill-showcase

This file captures the essential, repository-specific knowledge an AI coding agent needs to be productive.

**Big Picture:**
- **Language & layout:**: Small Go web service (single `main.go`) with a Docker multi-stage build (`dockerfile`).
- **Runtime:**: HTTP server exposing a text endpoint `/` and a JSON health endpoint `/api/status` (see `main.go`).
- **Roadmap:**: Database-backed fintech API planned (see `ToDo.md`) — PostgreSQL integration, transactions, and docker-compose are TODOs, not yet implemented.

**Key files to read and edit:**
- **`main.go`**: central HTTP handlers (`simpleTextHandler`, `jsonHandler`) and `StatusResponse` struct. Use this file as the entrypoint for adding routes and wiring services.
- **`dockerfile`**: multi-stage build (builder -> alpine runner). Keep the existing pattern: copy `go.mod`, `go mod download`, then copy sources and `go build -o server .`.
- **`go.mod`**: module name `go-web-server` and Go version `1.21` (verify compatibility with your local toolchain).
- **`ToDo.md`**: authoritative project roadmap and intended design decisions (DB, migrations, tests, docker-compose, iOS client).

**Project-specific conventions & observed patterns:**
- JSON responses use exported struct fields + `json:"name"` tags (example: `StatusResponse` in `main.go`).
- Handlers set `Content-Type` explicitly and commonly use `json.NewEncoder(w).Encode(...)` to write JSON.
- HTTP routing uses `http.HandleFunc` and default `http.ListenAndServe`. New features should follow the same minimal dependency approach unless adding a router is justified.

**Notable inconsistencies to address before large changes:**
- `dockerfile` sets `ENV PORT 8080`, but `main.go` currently hardcodes `const port = ":8080"` and does not read the `PORT` environment variable. If you add runtime configuration, prefer reading `PORT` from `os.Getenv("PORT")` with a fallback.
- `go.mod` specifies Go 1.21 but the Docker builder uses `golang:1.23-alpine`. Keep builds reproducible by aligning versions or explicitly documenting the mismatch.

**Build / run / debug workflows (verified from repo):**
- Local quick run: `go run .` from repository root.
- Build binary locally: `go build -o server .` then `./server` (listens on :8080).
- Docker build (uses `dockerfile`):
  - `docker build -t fintech-skill-showcase .`
  - `docker run -p 8080:8080 fintech-skill-showcase`

**Testing guidance (project-specific examples):**
- No tests exist yet. For handler-level tests, use `net/http/httptest`:
  - Create `_test.go` files.
  - Use `httptest.NewRequest(...)` + `httptest.NewRecorder()` and call handler functions directly.
- Naming: follow `*_test.go` pattern mentioned in `ToDo.md`.

**Integration & external dependencies:**
- Planned DB: PostgreSQL via `github.com/lib/pq` (not yet added). When adding DB access, use environment variables for credentials and prefer transaction usage with `db.Begin()` / `tx.Commit()` / `tx.Rollback()` as described in `ToDo.md`.

**When editing code, prefer small, clear changes:**
- Add routes to `main.go` near existing `http.HandleFunc` calls.
- If adding packages, update `go.mod` and keep `go.sum` tidy; the `dockerfile` relies on `go mod download` early in the build.

**Examples to copy/paste:**
- JSON handler pattern (from `main.go`):
```
type StatusResponse struct { Status string `json:"status"`; Service string `json:"service"`; Version string `json:"version"` }
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(data)
```

**What not to change without coordination:**
- The multi-stage Docker pattern in `dockerfile` — it's intentionally minimal and reproducible.
- The project roadmap in `ToDo.md`: treat it as the source of truth for planned features and tests.

If anything here is unclear or you want me to expand specific sections (example tests, adding env config, or wiring a PostgreSQL client), say which area and I'll iterate. 
