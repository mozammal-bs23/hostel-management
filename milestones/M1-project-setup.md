# Milestone 1 — Project Setup & Health Check

**Duration:** ~1 hour  
**Goal:** A running Go HTTP server with a clean project layout and a health endpoint.

---

## Prerequisites

- Go 1.21+ installed (`go version`)
- A terminal and a REST client (curl, Postman, or VS Code REST extension)

---

## Steps

### 1.1 Initialize the Go module

- Run `go mod init` with your module path (e.g. `github.com/<you>/hostel-management`).
- This creates `go.mod` at the project root.

### 1.2 Create the directory structure

Create the following empty directories. Files will be added as we go.

```
hostel-management/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── handlers/        # HTTP handler functions
│   ├── models/          # Domain structs & validation
│   ├── store/           # Data-access layer (interface + implementations)
│   └── router/          # Route definitions & middleware
├── milestones/          # These planning docs
├── go.mod
├── PRD.md
└── README.md
```

**Why `internal/`?**  
Go's `internal` package convention prevents other Go modules from importing your private code. It forces consumers (the Flutter app) to use the API, not import Go types directly.

### 1.3 Create `cmd/server/main.go`

This is the entry point. It should:

1. Read the `PORT` environment variable (default to `8080` if not set).
2. Create a new router (for now, just `http.NewServeMux` — chi comes in M3).
3. Register a `GET /health` route.
4. Start the HTTP server and log the port it's listening on.
5. Use `log.Fatal` so the process exits if the server fails to start.

### 1.4 Implement the health handler

- File: `internal/handlers/health.go`
- A simple function matching `http.HandlerFunc` signature.
- It should:
  - Set `Content-Type: application/json`
  - Write HTTP 200
  - Return the JSON body `{"status": "ok"}`

### 1.5 Wire the health handler in main

- Import `internal/handlers` in `main.go`.
- Register `GET /health` to the health handler.
- For now, use the standard library's `http.ServeMux`. We'll switch to chi in M3.

### 1.6 Add `.gitignore`

Include at least:
- `/bin/` or any compiled binary directory
- `.env`
- `.DS_Store`
- IDE folders (`.idea/`, `.vscode/` — optional)
- `vendor/` (if you don't vendor)

### 1.7 Create a basic `README.md`

Include:
- Project name and one-line description
- How to run: `go run ./cmd/server`
- Health check: `curl http://localhost:8080/health`
- Note that more endpoints are coming

### 1.8 Install a dependency — `chi` router (prepare for M3)

- Run `go get github.com/go-chi/chi/v5`
- We won't wire chi yet, but having it in `go.mod` now avoids interruption in M3.

---

## Acceptance Criteria

| # | Check | How to verify |
|---|-------|---------------|
| 1 | Server starts without errors | `go run ./cmd/server` logs "listening on :8080" |
| 2 | Health endpoint responds | `curl -i http://localhost:8080/health` returns 200 + `{"status":"ok"}` |
| 3 | Project compiles cleanly | `go build ./...` with zero errors |
| 4 | Directory layout matches PRD | Visual inspection |
| 5 | `.gitignore` present | `cat .gitignore` shows relevant patterns |

---

## Decisions Made

| Decision | Choice | Rationale |
|----------|--------|-----------|
| HTTP framework | chi (installed, wired in M3) | Lightweight, idiomatic, middleware-friendly |
| Config | Environment variables | Simple; no config library needed for just `PORT` |
| Project layout | `cmd/` + `internal/` | Standard Go project layout |

---

## Notes for You (Flutter Dev)

- `cmd/server/main.go` is like `lib/main.dart` — the entry point.
- `internal/` is like having private packages — nothing outside can import them.
- Go doesn't have a `pubspec.yaml`; instead `go.mod` + `go.sum` track dependencies.
- `go run ./cmd/server` compiles and runs in one step (like `flutter run`).

---

*Next → [M2: Room Model & In-Memory Store](./M2-room-model-store.md)*
