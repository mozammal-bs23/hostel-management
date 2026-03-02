# Milestone 4 — Resident Model & In-Memory Store

**Duration:** ~1.5 hours  
**Goal:** A Resident domain model with validation, store interface extension, in-memory implementation, and unit tests.

---

## Steps

### 4.1 Define the Resident model

- File: `internal/models/resident.go`
- Create a `Resident` struct with JSON tags:

  | Field | Type | JSON tag | Notes |
  |-------|------|----------|-------|
  | `ID` | string | `id` | UUID, auto-generated |
  | `Name` | string | `name` | Required |
  | `Email` | string | `email` | Required, must be unique across residents |
  | `Phone` | string | `phone` | Optional |
  | `RoomID` | string | `roomId` | Optional — empty means unassigned |
  | `CheckIn` | *time.Time | `checkIn` | Optional, pointer so `null` in JSON when absent |
  | `CheckOut` | *time.Time | `checkOut` | Optional, pointer |
  | `Status` | string | `status` | Required: `active`, `checked_out`, `pending` |
  | `CreatedAt` | time.Time | `createdAt` | Auto-set |
  | `UpdatedAt` | time.Time | `updatedAt` | Auto-set |

### 4.2 Add validation to the Resident model

- Create `Validate() error` method.
- Rules:
  - `Name` must not be empty.
  - `Email` must not be empty and should contain `@` (basic check — no need for full RFC validation).
  - `Status` must be one of the three allowed values.
  - If `CheckOut` is set, `CheckIn` must also be set.
  - If both are set, `CheckOut` must be after `CheckIn`.
- Define status constants: `ResidentStatusActive`, `ResidentStatusCheckedOut`, `ResidentStatusPending`.

### 4.3 Create a `CreateResidentRequest` struct (optional)

- Same pattern as Room: only user-supplied fields.
- `ToResident()` fills in `ID`, `CreatedAt`, `UpdatedAt`.

### 4.4 Define the `ResidentStore` interface

- File: `internal/store/store.go` (add to the existing file)

  | Method | Signature | Description |
  |--------|-----------|-------------|
  | Create | `Create(resident models.Resident) (models.Resident, error)` | Persist new resident |
  | GetByID | `GetByID(id string) (models.Resident, error)` | Fetch one |
  | List | `List(filters ResidentFilters) ([]models.Resident, error)` | List with filters |
  | Update | `Update(id string, resident models.Resident) (models.Resident, error)` | Update existing |
  | Delete | `Delete(id string) error` | Remove resident |

- Define `ResidentFilters`:
  - `RoomID *string` — filter by room assignment
  - `Status *string` — filter by status
  - `Limit int`
  - `Offset int`

### 4.5 Implement the in-memory resident store

- File: `internal/store/memory.go` (extend the existing file, or create `memory_resident.go`)

**Approach A — Single struct (recommended):**  
Add a `residents map[string]models.Resident` field and a second `sync.RWMutex` (or reuse one mutex for both) to the existing `MemoryStore` struct. Rename it if needed from `MemoryRoomStore` to `MemoryStore`.

**Approach B — Separate struct:**  
Create `MemoryResidentStore` with its own map and mutex. Both approaches are fine; A is simpler for wiring in `main.go`.

**Method details:**

| Method | Logic |
|--------|-------|
| `Create` | Check email uniqueness (iterate map) → validate → store → return |
| `GetByID` | Look up map → `ErrNotFound` if missing |
| `List` | Iterate → filter by `RoomID`, `Status` → apply limit/offset → return slice |
| `Update` | Check exists → optionally check email uniqueness if email changed → merge → update timestamps → store |
| `Delete` | Check exists → delete |

**Important — Email uniqueness:**  
When creating or updating, iterate the map to ensure no other resident has the same email. Return a descriptive error like `"email already exists"` (use a new sentinel error `var ErrDuplicate = errors.New("duplicate")`).

### 4.6 Write unit tests

- File: `internal/store/memory_resident_test.go` (or add to `memory_test.go`)

  | Test | What it verifies |
  |------|------------------|
  | `TestCreateResident` | Stored correctly; has ID and timestamps |
  | `TestCreateResident_DuplicateEmail` | Returns `ErrDuplicate` |
  | `TestGetResidentByID_Found` | Correct resident returned |
  | `TestGetResidentByID_NotFound` | Returns `ErrNotFound` |
  | `TestListResidents_All` | Returns all residents |
  | `TestListResidents_FilterByRoomID` | Only matching residents |
  | `TestListResidents_FilterByStatus` | Only matching residents |
  | `TestUpdateResident` | Fields update; `UpdatedAt` changes |
  | `TestUpdateResident_NotFound` | Returns `ErrNotFound` |
  | `TestDeleteResident` | Removed; subsequent Get returns `ErrNotFound` |

- Run: `go test ./internal/store/... -v`

### 4.7 Validate room existence on assignment (optional enhancement)

If time permits, when a resident is created or updated with a non-empty `RoomID`:
- Check that the room exists in the room store.
- Check that the room's current occupant count < capacity.
- This requires the resident store to have a reference to the room store (or handle this in the handler layer in M5).

**Recommendation:** Skip this for now. Handle it in the handler layer (M5) or defer to a future milestone. Keep the store layer simple.

---

## Acceptance Criteria

| # | Check | How to verify |
|---|-------|---------------|
| 1 | `Resident` struct has correct fields and tags | Code review |
| 2 | `Validate()` rejects empty name, empty email, invalid status | Unit test |
| 3 | `ResidentStore` interface is defined | Compiles |
| 4 | In-memory impl satisfies `ResidentStore` | Compiler check |
| 5 | Email uniqueness enforced on create | Unit test |
| 6 | List filters by `roomId` and `status` | Unit tests |
| 7 | All 10+ tests pass | `go test ./internal/store/... -v` |

---

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Date fields | `*time.Time` (pointer) | Allows `null` in JSON for optional dates |
| Email uniqueness | Iterate map | O(n) but fine for in-memory with small data |
| Room validation | Deferred to handler layer | Keeps store layer simple and single-responsibility |
| Sentinel errors | `ErrNotFound`, `ErrDuplicate` | Handlers can map errors to HTTP codes cleanly |

---

## Concepts for Flutter Devs

- **Pointer types (`*time.Time`)**: Like Dart's `DateTime?`. A `nil` pointer serializes to JSON `null`.
- **Sentinel errors**: Like custom exception classes in Dart (`class NotFoundException extends AppException`). In Go, you compare with `errors.Is(err, store.ErrNotFound)`.
- **No generics needed here**: The `RoomStore` and `ResidentStore` are separate interfaces. Go's approach favors explicit interfaces over generic repositories.

---

*Previous → [M3: Room HTTP Handlers](./M3-room-handlers.md)*  
*Next → [M5: Resident HTTP Handlers](./M5-resident-handlers.md)*
