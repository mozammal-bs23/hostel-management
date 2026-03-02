# Milestone 6 — Polish & Documentation

**Duration:** ~1 hour  
**Goal:** Production-ready README, seed data examples, optional SQLite upgrade, and final cleanup.

---

## Steps

### 6.1 Write a comprehensive README

- File: `README.md` (replace the placeholder from M1)
- Sections to include:

  | Section | Content |
  |---------|---------|
  | **Project Title** | Hostel Management System — Go Backend |
  | **Description** | One paragraph: what it does, who it's for |
  | **Tech Stack** | Go, chi, in-memory store (or SQLite) |
  | **Prerequisites** | Go 1.21+ |
  | **Getting Started** | Clone → `go mod download` → `go run ./cmd/server` |
  | **Environment Variables** | `PORT` (default 8080) |
  | **API Reference** | Table of all endpoints with method, path, description |
  | **Request/Response Examples** | 2–3 curl examples with expected output |
  | **Project Structure** | Directory tree with one-line descriptions |
  | **Running Tests** | `go test ./... -v` |
  | **Future Work** | Auth, Flutter app, Postgres, etc. |

### 6.2 Add curl seed script

- File: `scripts/seed.sh`
- A bash script that creates sample data via curl:
  - 3 rooms (single, double, dormitory)
  - 3 residents (one per room, one unassigned)
- The script should:
  - Print what it's doing at each step
  - Use `jq` to extract IDs from responses (with a fallback note if jq isn't installed)
  - Be executable (`chmod +x scripts/seed.sh`)

### 6.3 Add CORS middleware (Flutter will need this)

- In `internal/router/router.go`, add a simple CORS middleware or use `github.com/go-chi/cors`.
- Allow:
  - Origins: `*` (for development; restrict in production)
  - Methods: `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`
  - Headers: `Content-Type`, `Authorization`
- This prevents Flutter web from hitting CORS errors during development.
- Run `go get github.com/go-chi/cors` if using the chi cors package.

### 6.4 Add graceful shutdown

- In `main.go`, listen for OS signals (`SIGINT`, `SIGTERM`) using `os/signal`.
- On signal: call `server.Shutdown(ctx)` with a 10-second timeout.
- Log "shutting down..." and "server stopped."
- This prevents data loss and is good practice.

### 6.5 Review and clean up code

Checklist:

- [ ] All files have consistent formatting (`go fmt ./...`)
- [ ] No unused imports (`go vet ./...`)
- [ ] No exported functions without purpose (keep handler constructors exported, keep internal helpers unexported)
- [ ] JSON tags use `camelCase` consistently (matches Flutter/Dart conventions)
- [ ] Error messages are lowercase (Go convention)
- [ ] All tests pass: `go test ./... -v`

### 6.6 Optional — SQLite persistence

If you finished early and want persistence across restarts:

1. Run `go get github.com/mattn/go-sqlite3` (requires CGO) or `go get modernc.org/sqlite` (pure Go, no CGO).
2. Create `internal/store/sqlite.go`.
3. Implement `RoomStore` and `ResidentStore` interfaces using SQL.
4. Create tables on startup (auto-migrate):
   - `rooms` table matching the Room struct
   - `residents` table matching the Resident struct with a foreign key to `rooms`
5. Update `main.go` to choose store based on an env var: `STORE=memory` (default) or `STORE=sqlite`.
6. Update README with SQLite instructions.

**Skip this if time is tight.** The in-memory store is perfectly fine for development and demonstrating CRUD.

### 6.7 Update `.gitignore`

Ensure it includes:
```
# Binaries
/bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Go workspace
go.work
go.work.sum

# Environment
.env
.env.local

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# SQLite (if used)
*.db
*.sqlite
```

### 6.8 Final verification

Run through the complete workflow:

1. `go build ./...` — compiles cleanly
2. `go vet ./...` — no issues
3. `go test ./... -v` — all tests pass
4. `go run ./cmd/server` — server starts
5. `bash scripts/seed.sh` — seed data loads
6. Manual curl through every endpoint — all work
7. `Ctrl+C` — graceful shutdown message appears

---

## Acceptance Criteria

| # | Check | How to verify |
|---|-------|---------------|
| 1 | README has all sections | Read the file |
| 2 | Seed script creates rooms + residents | Run `bash scripts/seed.sh` |
| 3 | CORS headers present | `curl -I -X OPTIONS http://localhost:8080/api/v1/rooms` |
| 4 | Graceful shutdown works | Start server → Ctrl+C → see "server stopped" |
| 5 | `go vet` and `go fmt` pass | Run both with zero output |
| 6 | All tests pass | `go test ./... -v` |
| 7 | A new developer can set up in < 5 min | Follow README from scratch |

---

## Final Project Structure

```
hostel-management/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, server startup, graceful shutdown
├── internal/
│   ├── handlers/
│   │   ├── health.go            # GET /health
│   │   ├── response.go          # JSON response helpers
│   │   ├── room.go              # Room CRUD handlers
│   │   └── resident.go          # Resident CRUD handlers
│   ├── models/
│   │   ├── room.go              # Room struct, validation, constants
│   │   └── resident.go          # Resident struct, validation, constants
│   ├── store/
│   │   ├── store.go             # Interfaces: RoomStore, ResidentStore
│   │   ├── memory.go            # In-memory implementation
│   │   └── memory_test.go       # Unit tests (or split into two test files)
│   └── router/
│       └── router.go            # chi router, middleware, route groups
├── scripts/
│   └── seed.sh                  # Sample data loader
├── milestones/                  # These planning documents
│   ├── M1-project-setup.md
│   ├── M2-room-model-store.md
│   ├── M3-room-handlers.md
│   ├── M4-resident-model-store.md
│   ├── M5-resident-handlers.md
│   └── M6-polish-docs.md
├── .gitignore
├── go.mod
├── go.sum
├── PRD.md
└── README.md
```

---

## What You've Built

By completing all 6 milestones, you have:

- A clean, idiomatic Go REST API
- Full CRUD for Rooms and Residents
- Consistent JSON responses with proper HTTP status codes
- Filtered list endpoints ready for Flutter integration
- CORS support for Flutter web
- Graceful shutdown
- Unit tests for the data layer
- Documentation a Flutter developer can follow

**Total endpoints: 11**
- `GET /health`
- `POST/GET/GET/:id/PUT/:id/DELETE/:id` for `/api/v1/rooms`
- `POST/GET/GET/:id/PUT/:id/DELETE/:id` for `/api/v1/residents`

---

## What's Next (Post 8 Hours)

| Priority | Feature | Estimated Time |
|----------|---------|----------------|
| 1 | Flutter app consuming this API | 4–6h |
| 2 | SQLite or Postgres persistence | 2–3h |
| 3 | JWT authentication | 2–3h |
| 4 | Check-in/check-out workflow | 1–2h |
| 5 | Room occupancy tracking | 1–2h |
| 6 | Pagination metadata in list responses | 30min |
| 7 | Docker + docker-compose | 1h |

---

*Previous → [M5: Resident HTTP Handlers](./M5-resident-handlers.md)*  
*Back to → [PRD](../PRD.md)*
