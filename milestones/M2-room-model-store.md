# Milestone 2 — Room Model & In-Memory Store

**Duration:** ~1.5 hours  
**Goal:** A Room domain model with validation, a store interface, and an in-memory implementation with unit tests.

---

## Steps

### 2.1 Define the Room model

- File: `internal/models/room.go`
- Create a `Room` struct with JSON tags for every field:
  - `ID` (string) — UUID, auto-generated
  - `Name` (string) — required, e.g. "Room 101"
  - `Capacity` (int) — required, must be > 0
  - `RentalPrice` (float64) — required, must be >= 0
  - `Status` (string) — required, one of: `available`, `occupied`, `maintenance`
  - `CreatedAt` (time.Time) — auto-set on creation
  - `UpdatedAt` (time.Time) — auto-set on creation and update

### 2.2 Add validation to the Room model

- Create a `Validate() error` method on Room.
- Rules:
  - `Name` must not be empty.
  - `Capacity` must be >= 1.
  - `RentalPrice` must be >= 0.
  - `Status` must be one of the three allowed values.
- Return a descriptive error message (e.g. `"name is required"`).
- **Tip:** Define the allowed statuses as constants in the same file for reuse.

### 2.3 Add a helper: `CreateRoom` request struct (optional)

- A separate struct `CreateRoomRequest` with only the user-supplied fields (no ID, timestamps).
- A method `ToRoom() Room` that fills in `ID` (uuid.New()), `CreatedAt`, `UpdatedAt`.
- This keeps the handler layer clean (M3).

### 2.4 Install UUID package

- Run `go get github.com/google/uuid`
- Use `uuid.New().String()` to generate room IDs.

### 2.5 Define the store interface

- File: `internal/store/store.go`
- Create a `RoomStore` interface:

  | Method | Signature | Description |
  |--------|-----------|-------------|
  | Create | `Create(room models.Room) (models.Room, error)` | Persist a new room |
  | GetByID | `GetByID(id string) (models.Room, error)` | Fetch one room |
  | List | `List(filters RoomFilters) ([]models.Room, error)` | List with optional filters |
  | Update | `Update(id string, room models.Room) (models.Room, error)` | Update existing room |
  | Delete | `Delete(id string) error` | Remove room |

- Define a `RoomFilters` struct with optional fields: `Status *string`, `Limit int`, `Offset int` (add more filters later if needed, e.g. price range).

### 2.6 Define a sentinel error

- In `store.go`, declare `var ErrNotFound = errors.New("not found")`.
- All store implementations return this when an ID doesn't exist.
- Handlers will check for this to return HTTP 404 (in M3).

### 2.7 Implement the in-memory room store

- File: `internal/store/memory.go`
- Use a struct with a `sync.RWMutex` and a `map[string]models.Room`.
- Constructor: `NewMemoryRoomStore() *MemoryRoomStore` — returns a ready-to-use instance.

**Method details:**

| Method | Logic |
|--------|-------|
| `Create` | Validate → store in map → return the room |
| `GetByID` | Lock read → look up map → return `ErrNotFound` if missing |
| `List` | Lock read → iterate map → apply filters → apply limit/offset → return slice |
| `Update` | Lock write → check exists → merge fields → update `UpdatedAt` → store |
| `Delete` | Lock write → check exists → delete from map |

**Concurrency:** Use `RLock` for reads and `Lock` for writes. This matters because the HTTP server handles requests concurrently.

### 2.8 Write unit tests

- File: `internal/store/memory_test.go`
- Test cases:

  | Test | What it verifies |
  |------|------------------|
  | `TestCreateRoom` | Room is stored; returned room has ID and timestamps |
  | `TestGetByID_Found` | Returns the correct room |
  | `TestGetByID_NotFound` | Returns `ErrNotFound` |
  | `TestListRooms_All` | Returns all rooms |
  | `TestListRooms_FilterByStatus` | Only matching rooms returned |
  | `TestUpdateRoom` | Fields update; `UpdatedAt` changes |
  | `TestUpdateRoom_NotFound` | Returns `ErrNotFound` |
  | `TestDeleteRoom` | Room removed; subsequent Get returns `ErrNotFound` |
  | `TestDeleteRoom_NotFound` | Returns `ErrNotFound` |

- Run with `go test ./internal/store/... -v`

---

## Acceptance Criteria

| # | Check | How to verify |
|---|-------|---------------|
| 1 | Room struct has correct fields and JSON tags | Code review |
| 2 | `Validate()` rejects invalid input | Unit test or manual call |
| 3 | `RoomStore` interface is defined | Compiles without errors |
| 4 | In-memory store implements `RoomStore` | Compiler verifies this |
| 5 | All 9+ unit tests pass | `go test ./internal/store/... -v` |

---

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Storage | `map[string]Room` + mutex | Simplest for an 8-hour scope; swap later |
| ID generation | UUID v4 | Globally unique, no DB sequence needed |
| Filtering | Pointer fields in `RoomFilters` | `nil` means "don't filter" for optional string filters |
| Interface | `RoomStore` | Easy to swap in-memory → SQLite later |

---

## Concepts for Flutter Devs

- **`sync.RWMutex`**: Like a `ReadWriteLock`. Multiple goroutines can read simultaneously, but writes are exclusive. Similar to using `synchronized` blocks.
- **Interface in Go**: Implicit — if your struct has the right methods, it satisfies the interface. No `implements` keyword like in Dart.
- **Pointer fields in filters**: In Dart you'd use nullable types (`int?`). In Go, use `*int` for the same purpose.

---

*Previous → [M1: Project Setup](./M1-project-setup.md)*  
*Next → [M3: Room HTTP Handlers](./M3-room-handlers.md)*  
*Boarders (guests) are introduced in [M4](./M4-boarder-model-store.md); `boarder.roomId` references `room.id`.*
