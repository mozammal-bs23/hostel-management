# Milestone 2 Evaluation — Eval #1

**Score: 12/20**  
**Status:** ❌ Incomplete (store implementation & tests missing)

### 🏗 Stack-wise Evaluation Breakdown

| Criteria | Score | Finding |
| :--- | :---: | :--- |
| **Domain model** | 4/4 | `internal/models/room.go` defines `Room` with JSON tags on every field, status constants (`available`, `occupied`, `maintenance`), and `Validate()` with clear rules and messages. |
| **Store interface & errors** | 2/4 | `RoomStore` and `RoomFilters` exist in `internal/store/room_store.go` (milestone asked for `store.go`). **Gaps:** no `var ErrNotFound = errors.New("not found")`; `GetByID` / `List` / `Update` use a custom `ErrNotFound` struct as the second return type instead of `error`, which diverges from the milestone and complicates `errors.Is` for HTTP 404 in M3. `List` should return `([]models.Room, error)`, not a not-found-shaped second value. |
| **In-memory store** | 0/4 | **Missing:** no `internal/store/memory.go`, no `sync.RWMutex` + `map[string]models.Room`, no `NewMemoryRoomStore()`. This is the core deliverable for M2 §2.7. |
| **Unit tests** | 0/4 | **Missing:** no `internal/store/memory_test.go`; none of the nine named tests (`TestCreateRoom`, `TestGetByID_*`, `TestListRooms_*`, `TestUpdateRoom*`, `TestDeleteRoom*`) are present. `go test ./internal/store/...` reports no test files. |
| **Dependencies & optional helper** | 2/4 | `github.com/google/uuid` is in `go.mod`. **Gaps:** optional `CreateRoomRequest` + `ToRoom()` (§2.3) not implemented; UUID generation is not wired in model code yet (expected once `ToRoom` or `Create` exists). |

### 📝 Summary

The **Room** type and **validation** meet the milestone and are ready for handlers. The **store layer is only sketched**: the interface file does not match the specified API (`error` returns, sentinel `ErrNotFound`), and the **thread-safe in-memory implementation** plus **unit tests** are not implemented—so M2 acceptance criteria 4 and 5 fail. Close the gap with `store.go` (sentinel), `memory.go` (mutex-backed map), `memory_test.go` (nine cases), and optionally `CreateRoomRequest` / `ToRoom()` in `room.go`.

### 🚀 Next Step

Complete the remaining M2 items, then proceed to **Milestone 3 (Room HTTP Handlers)** per `milestones/M3-room-handlers.md` (wire `NewMemoryRoomStore()`, REST routes, `errors.Is(err, store.ErrNotFound)` → 404).

---

## Acceptance criteria (from `milestones/M2-room-model-store.md`)

| # | Check | Status |
|---|-------|--------|
| 1 | Room struct has correct fields and JSON tags | ✅ |
| 2 | `Validate()` rejects invalid input | ✅ (no dedicated model test in repo) |
| 3 | `RoomStore` interface is defined | ⚠️ Defined but not aligned with spec |
| 4 | In-memory store implements `RoomStore` | ❌ |
| 5 | All 9+ unit tests pass | ❌ |
