# Milestone 4 — Boarder Model & In-Memory Store

**Duration:** ~1.5 hours  
**Goal:** A Boarder domain model with validation, store interface extension, in-memory implementation, and unit tests. **`roomId` is the foreign key to `Room.id`.**

---

## Steps

### 4.1 Define the Boarder model

- File: `internal/models/boarder.go`
- Create a `Boarder` struct with JSON tags:

  | Field | Type | JSON tag | Notes |
  |-------|------|----------|-------|
  | `ID` | string | `id` | UUID, auto-generated |
  | `FirstName` | string | `firstName` | Required together with last name: not both empty |
  | `LastName` | string | `lastName` | Same as above |
  | `Phone` | string | `phone` | Required |
  | `RoomID` | string | `roomId` | **Required** — references `Room.ID` |
  | `Status` | string | `status` | Required: `active`, `checked_out`, `pending` |
  | `CreatedAt` | time.Time | `createdAt` | Auto-set |
  | `UpdatedAt` | time.Time | `updatedAt` | Auto-set |

- Status constants (example names): `StatusBoarderActive`, `StatusBoarderCheckedOut`, `StatusBoarderPending`.

### 4.2 Add validation to the Boarder model

- Create `Validate() error` method.
- Rules:
  - At least one of `FirstName` / `LastName` must be non-empty.
  - `Phone` must not be empty.
  - `RoomID` must not be empty (room assignment required for this model).
  - `Status` must be one of the three allowed values.

### 4.3 Create a `CreateBoarderRequest` struct (optional)

- Same pattern as Room: only user-supplied fields.
- `ToBoarder()` fills in `ID`, `CreatedAt`, `UpdatedAt`.

### 4.4 Define the `BoarderStore` interface

- File: `internal/store/store.go` (add to the existing file)

  | Method | Signature | Description |
  |--------|-----------|-------------|
  | Create | `Create(boarder models.Boarder) (models.Boarder, error)` | Persist new boarder |
  | GetByID | `GetByID(id string) (models.Boarder, error)` | Fetch one |
  | List | `List(filters BoarderFilters) ([]models.Boarder, error)` | List with filters |
  | Update | `Update(id string, boarder models.Boarder) (models.Boarder, error)` | Update existing |
  | Delete | `Delete(id string) error` | Remove boarder |

- Define `BoarderFilters`:
  - `RoomID *string` — filter by room assignment
  - `Status *string` — filter by status
  - `Limit int`
  - `Offset int`

### 4.5 Implement the in-memory boarder store

- File: `internal/store/memory.go` (extend the existing file, or add `memory_boarder.go`)

**Approach A — Single struct (recommended):**  
Add a `boarders map[string]models.Boarder` field and coordinate locking with the existing room map (one mutex for both maps is fine for this scope).

**Approach B — Separate struct:**  
Create `MemoryBoarderStore` with its own map and mutex.

**Method details:**

| Method | Logic |
|--------|-------|
| `Create` | Validate → store → return |
| `GetByID` | Look up map → `ErrNotFound` if missing |
| `List` | Iterate → filter by `RoomID`, `Status` → apply limit/offset → return slice |
| `Update` | Check exists → merge → update timestamps → store |
| `Delete` | Check exists → delete |

**Optional — phone uniqueness:**  
If the product requires unique phone numbers, iterate the map on create/update and return `ErrDuplicate` when another boarder has the same phone. Omit if not needed.

### 4.6 Write unit tests

- File: `internal/store/memory_boarder_test.go` (or add to `memory_test.go`)

  | Test | What it verifies |
  |------|------------------|
  | `TestCreateBoarder` | Stored correctly; has ID and timestamps |
  | `TestGetBoarderByID_Found` | Correct boarder returned |
  | `TestGetBoarderByID_NotFound` | Returns `ErrNotFound` |
  | `TestListBoarders_All` | Returns all boarders |
  | `TestListBoarders_FilterByRoomID` | Only matching boarders |
  | `TestListBoarders_FilterByStatus` | Only matching boarders |
  | `TestUpdateBoarder` | Fields update; `UpdatedAt` changes |
  | `TestUpdateBoarder_NotFound` | Returns `ErrNotFound` |
  | `TestDeleteBoarder` | Removed; subsequent Get returns `ErrNotFound` |

- Run: `go test ./internal/store/... -v`

### 4.7 Validate room existence on assignment (optional enhancement)

When a boarder is created or updated with a non-empty `RoomID`:

- Check that the room exists in the room store.
- Optionally check capacity vs. current boarders in that room.

**Recommendation:** Implement in the **handler layer (M5)** so the store stays simple, or inject `RoomStore` into a small service. The domain rule is: **`roomId` must reference an existing `Room.id`**.

---

## Acceptance Criteria

| # | Check | How to verify |
|---|-------|---------------|
| 1 | `Boarder` struct matches PRD / `boarder.go` | Code review |
| 2 | `Validate()` enforces name, phone, `roomId`, status | Unit test |
| 3 | `BoarderStore` interface is defined | Compiles |
| 4 | In-memory impl satisfies `BoarderStore` | Compiler check |
| 5 | List filters by `roomId` and `status` | Unit tests |
| 6 | Tests pass | `go test ./internal/store/... -v` |

---

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| `roomId` required | Matches current model | Clear FK to `Room`; API always assigns a room |
| Room existence check | Handler or service (M5) | Store can stay unaware of `RoomStore` |
| Sentinel errors | `ErrNotFound`, optional `ErrDuplicate` | Map cleanly to HTTP 404 / 409 |

---

## Concepts for Flutter Devs

- **`roomId` in JSON** maps to `Room.id` from the rooms API — load rooms first, then create boarders with a valid `roomId`.
- **Foreign key in Go** is just a string field; referential integrity is enforced in handlers or the database.

---

*Previous → [M3: Room HTTP Handlers](./M3-room-handlers.md)*  
*Next → [M5: Boarder HTTP Handlers](./M5-boarder-handlers.md)*
