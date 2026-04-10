# Product Requirements Document (PRD)
# Hostel Management System — Go Backend

**Version:** 1.1  
**Target:** 8 hours (single developer)  
**Backend:** Go (Golang)  
**Future consumer:** Flutter app (REST API)

---

## 1. Overview

### 1.1 Purpose
A REST API backend in Go for a **hostel management system** that supports core CRUD operations for **Rooms** and **Boarders** (people staying in rooms). The API will be consumable by a Flutter frontend later.

### 1.2 Goals
- Implement full CRUD for **Rooms** and **Boarders**
- Use in-memory or SQLite storage (no external DB setup)
- Clean project structure, ready for Flutter integration
- Achievable in **8 hours** with clear milestones

### 1.3 Out of Scope (for this 8-hour slice)
- Authentication/authorization
- Payments or billing
- File uploads (e.g. ID photos)
- Real-time features (WebSockets)
- Production deployment or Docker

---

## 2. Domain Model

### 2.1 Room
| Field          | Type    | Required | Description                              |
|----------------|---------|----------|------------------------------------------|
| `id`           | string  | auto     | UUID                                     |
| `name`         | string  | yes      | e.g. "Room 101"                          |
| `capacity`     | int     | yes      | Max occupants (e.g. 2, 4)                |
| `status`       | string  | yes      | `available`, `occupied`, `maintenance`   |
| `rentalPrice`  | number  | yes      | Non-negative rental amount               |
| `createdAt`    | string  | auto     | ISO 8601                                 |
| `updatedAt`    | string  | auto     | ISO 8601                                 |

*Implementation:* `internal/models/room.go`

### 2.2 Boarder
A **boarder** is a guest assigned to a room. `roomId` is the logical foreign key to `Room.id`.

| Field        | Type   | Required | Description                                      |
|-------------|--------|----------|--------------------------------------------------|
| `id`        | string | auto     | UUID                                             |
| `firstName` | string | yes*     | Given name (at least one of first/last required) |
| `lastName`  | string | yes*     | Family name                                      |
| `phone`     | string | yes      | Contact number                                   |
| `roomId`    | string | yes      | FK to Room (`Room.id`)                           |
| `status`    | string | yes      | `active`, `checked_out`, `pending`               |
| `createdAt` | string | auto     | ISO 8601                                         |
| `updatedAt` | string | auto     | ISO 8601                                         |

\*Validation: `firstName` and `lastName` cannot both be empty.

*Implementation:* `internal/models/boarder.go`

---

## 3. API Specification

### 3.1 Base URL
- Local: `http://localhost:8080/api/v1`

### 3.2 Rooms

| Method | Path              | Description        |
|--------|-------------------|--------------------|
| GET    | `/rooms`          | List all rooms     |
| GET    | `/rooms/:id`      | Get room by ID     |
| POST   | `/rooms`          | Create room        |
| PUT    | `/rooms/:id`      | Update room        |
| DELETE | `/rooms/:id`      | Delete room        |

**Query (GET /rooms):** optional `status`, `limit`, `offset` (extend with more filters as needed).

### 3.3 Boarders

| Method | Path                | Description           |
|--------|---------------------|-----------------------|
| GET    | `/boarders`         | List all boarders     |
| GET    | `/boarders/:id`     | Get boarder by ID     |
| POST   | `/boarders`         | Create boarder        |
| PUT    | `/boarders/:id`     | Update boarder        |
| DELETE | `/boarders/:id`     | Delete boarder        |

**Query (GET /boarders):** optional `roomId`, `status`, `limit`, `offset`.

### 3.4 Response Format
- Success: `{ "data": <entity|array>, "message": "..." }` with HTTP 200/201
- Error: `{ "error": "message", "code": "ERROR_CODE" }` with 4xx/5xx
- Use `Content-Type: application/json`

---

## 4. Milestones (8 Hours)

### Milestone 1 — Project setup & health (≈1 hour)
**Deliverables:**
- Go module initialized (`go mod init`)
- Project layout: `cmd/server`, `internal/handlers`, `internal/models`, `internal/store`, `internal/router`
- HTTP server on port 8080
- Health check: `GET /health` → `{"status":"ok"}`
- README with run instructions

**Acceptance:** `curl http://localhost:8080/health` returns 200 and JSON.

---

### Milestone 2 — Room model & in-memory store (≈1.5 hours)
**Deliverables:**
- `internal/models/room.go` — Room struct (including `rentalPrice`) and validation
- `internal/store/store.go` — interface (e.g. `RoomStore`) with CRUD
- `internal/store/memory.go` — in-memory implementation with map keyed by ID
- Unit tests for store (create, get, list, update, delete)

**Acceptance:** All store tests pass; IDs are UUIDs; list supports optional filters.

---

### Milestone 3 — Room HTTP handlers & router (≈1.5 hours)
**Deliverables:**
- `internal/handlers/room.go` — Create, GetByID, List, Update, Delete
- `internal/router/router.go` — mount `/api/v1/rooms` and wire handlers
- Request/response JSON binding
- 404 on missing room; 400 on invalid body

**Acceptance:** Full Room CRUD via curl/Postman; correct status codes and JSON.

---

### Milestone 4 — Boarder model & store (≈1.5 hours)
**Deliverables:**
- `internal/models/boarder.go` — Boarder struct and validation (`roomId` required; aligns with Room FK)
- Extend store interface and in-memory store for boarders
- Optional: ensure “list boarders by `roomId`” is efficient
- Unit tests for boarder store

**Acceptance:** Boarder CRUD in store works; list by `roomId` works.

---

### Milestone 5 — Boarder HTTP handlers & routes (≈1.5 hours)
**Deliverables:**
- `internal/handlers/boarder.go` — Create, GetByID, List, Update, Delete
- Mount `/api/v1/boarders` in router
- Query params: `roomId`, `status`, `limit`, `offset`
- Consistent error responses with Room API
- Optional: when `roomId` is set, verify room exists via `RoomStore` (400 if invalid)

**Acceptance:** Full Boarder CRUD and filtered list via API.

---

### Milestone 6 — Polish & documentation (≈1 hour)
**Deliverables:**
- Optional: swap to SQLite (e.g. `internal/store/sqlite.go`) or keep in-memory and document it
- README: how to run, env vars (e.g. `PORT`), list of endpoints with examples
- `.gitignore` for Go (binaries, vendor, IDE)
- Seed script or curl examples for 2–3 rooms and 2–3 boarders

**Acceptance:** New developer can clone, run server, and call all endpoints from README.

---

## 5. Technical Stack

| Layer    | Choice              | Notes                          |
|----------|---------------------|--------------------------------|
| Language | Go 1.21+            |                                |
| HTTP     | `net/http` or chi   | Chi recommended for routing    |
| Storage  | In-memory (default) | Optional SQLite in M6          |
| ID       | `github.com/google/uuid` | For `id` fields          |
| Config   | Env (e.g. `PORT`)   | Default port 8080              |

---

## 6. Project Structure (Suggested)

```
hostel-management/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── room.go
│   │   └── boarder.go
│   ├── models/
│   │   ├── room.go
│   │   └── boarder.go
│   ├── store/
│   │   ├── store.go      # interfaces
│   │   ├── memory.go     # in-memory impl
│   │   └── (optional) sqlite.go
│   └── router/
│       └── router.go
├── go.mod
├── go.sum
├── README.md
├── PRD.md
└── .gitignore
```

---

## 7. Time Allocation Summary

| Milestone | Duration | Cumulative |
|-----------|----------|------------|
| M1 — Setup & health | 1h  | 1h  |
| M2 — Room model & store | 1.5h | 2.5h |
| M3 — Room API | 1.5h | 4h  |
| M4 — Boarder model & store | 1.5h | 5.5h |
| M5 — Boarder API | 1.5h | 7h  |
| M6 — Polish & docs | 1h  | 8h  |

---

## 8. Success Criteria

- [ ] Server runs with `go run ./cmd/server` (or binary).
- [ ] All Room CRUD endpoints work and return correct JSON and status codes.
- [ ] All Boarder CRUD endpoints work with optional query filters.
- [ ] README allows a Flutter developer to integrate against the API without reading code.
- [ ] Code is structured so adding SQLite or Postgres later is straightforward.

---

## 9. Future Extensions (Post–8 hours)

- JWT auth and role-based access
- Flutter app consuming this API
- SQLite/Postgres persistence
- Check-in/check-out dates on boarders and room status automation
- Pagination metadata in list responses (`total`, `limit`, `offset`)
- Floor or building metadata on rooms if product needs it

---

*PRD aligned with `Room` and `Boarder` models in `internal/models`. Adjust milestones if you spend more or less time on testing or SQLite.*
