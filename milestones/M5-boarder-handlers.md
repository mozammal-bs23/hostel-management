# Milestone 5 — Boarder HTTP Handlers & Routes

**Duration:** ~1.5 hours  
**Goal:** Full CRUD REST endpoints for Boarders, mounted on the chi router with query filtering and consistent error handling.

---

## Steps

### 5.1 Create the Boarder handler struct

- File: `internal/handlers/boarder.go`
- Create `BoarderHandler` struct holding a `store.BoarderStore`.
- Constructor: `NewBoarderHandler(s store.BoarderStore) *BoarderHandler`
- This mirrors the pattern from `RoomHandler` in M3.

### 5.2 Implement handler methods

| Handler | Route | Request Body | Status Codes | Logic |
|---------|-------|-------------|-------------|-------|
| `Create` | `POST /` | JSON body | 201, 400, 409? | Decode → validate → `store.Create` → respond 201 |
| `GetByID` | `GET /{id}` | — | 200, 404 | Extract ID → `store.GetByID` → 404 if not found |
| `List` | `GET /` | — | 200 | Parse query params → build `BoarderFilters` → `store.List` → respond 200 |
| `Update` | `PUT /{id}` | JSON | 200, 400, 404, 409? | Extract ID → decode → validate → `store.Update` → handle errors |
| `Delete` | `DELETE /{id}` | — | 204, 404 | Extract ID → `store.Delete` → 204 on success |

If you enforce duplicate phone in the store, map `ErrDuplicate` to **409 Conflict**.

### 5.3 Handle query parameters for List

Parse these from the URL query string (`r.URL.Query()`):

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `roomId` | string | — | Filter by room assignment |
| `status` | string | — | Filter by `active`, `checked_out`, `pending` |
| `limit` | int | 50 | Max results per page |
| `offset` | int | 0 | Skip N results |

**Parsing tips:**

- Use `r.URL.Query().Get("roomId")` — empty string if absent.
- Convert limit/offset with `strconv.Atoi()` — use defaults on error.
- Build `BoarderFilters` with `*string` for optional filters: set pointer only if the query param is non-empty.

### 5.4 Error mapping

Reuse the `RespondError` helper from M3. Map store errors to HTTP codes:

| Store Error | HTTP Code | Response |
|-------------|-----------|----------|
| `ErrNotFound` | 404 | e.g. `{"error": "boarder not found"}` |
| `ErrDuplicate` (if used) | 409 Conflict | e.g. `{"error": "phone already exists"}` |
| Validation error | 400 | `{"error": "<validation message>"}` |
| JSON decode error | 400 | `{"error": "invalid request body"}` |

### 5.5 Wire boarder routes in the router

- Update `internal/router/router.go`:
  - Accept `boarderStore store.BoarderStore` (or a combined store facade).
  - Add a route group:

    ```
    r.Route("/api/v1/boarders", func(r chi.Router) {
        r.Post("/", boarderHandler.Create)
        r.Get("/", boarderHandler.List)
        r.Get("/{id}", boarderHandler.GetByID)
        r.Put("/{id}", boarderHandler.Update)
        r.Delete("/{id}", boarderHandler.Delete)
    })
    ```

### 5.6 Update `main.go`

- Instantiate the boarder store (or unified memory store) and pass it into `router.New(...)`.

### 5.7 Manual testing with curl

```
# Create a boarder (use a real room id from POST /api/v1/rooms)
curl -X POST http://localhost:8080/api/v1/boarders \
  -H "Content-Type: application/json" \
  -d '{"firstName":"John","lastName":"Doe","phone":"+1234567890","roomId":"<room-id>","status":"active"}'

# List all boarders
curl http://localhost:8080/api/v1/boarders

# List by status
curl "http://localhost:8080/api/v1/boarders?status=active"

# Get by ID
curl http://localhost:8080/api/v1/boarders/<id>

# Update — change room
curl -X PUT http://localhost:8080/api/v1/boarders/<id> \
  -H "Content-Type: application/json" \
  -d '{"firstName":"John","lastName":"Doe","phone":"+1234567890","roomId":"<other-room-id>","status":"active"}'

# Delete
curl -X DELETE http://localhost:8080/api/v1/boarders/<id>
```

### 5.8 Cross-resource verification (recommended)

In `Create` / `Update`, when `roomId` is set:

- Call `roomStore.GetByID(roomId)` to verify the room exists.
- Return **400** if missing: `{"error": "room not found"}`

This requires `BoarderHandler` to also hold a `RoomStore` reference (or a small application service).

---

## Acceptance Criteria

| # | Check | How to verify |
|---|-------|---------------|
| 1 | `POST /api/v1/boarders` creates boarder | 201 + JSON |
| 2 | `GET /api/v1/boarders` lists all | 200 + array |
| 3 | `GET /api/v1/boarders?status=active` filters | Only active boarders |
| 4 | `GET /api/v1/boarders?roomId=<id>` filters | Only boarders in that room |
| 5 | `GET /api/v1/boarders/:id` returns one | 200 + object |
| 6 | `PUT /api/v1/boarders/:id` updates | 200 + updated fields |
| 7 | `DELETE /api/v1/boarders/:id` removes | 204 + GET → 404 |
| 8 | Room endpoints still work | Regression curl |
| 9 | Health endpoint still works | `GET /health` → 200 |

---

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| FK validation | Handler + `RoomStore` | Enforces `roomId` → existing room without coupling stores |
| Limit default | 50 | Prevents unbounded list responses |

---

## Concepts for Flutter Devs

- **Nested resource:** Boarders belong to a room via `roomId`; fetch rooms first, then boarders for a room using `?roomId=`.

---

*Previous → [M4: Boarder Model & Store](./M4-boarder-model-store.md)*  
*Next → [M6: Polish & Documentation](./M6-polish-docs.md)*
