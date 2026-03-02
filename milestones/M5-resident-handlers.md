# Milestone 5 — Resident HTTP Handlers & Routes

**Duration:** ~1.5 hours  
**Goal:** Full CRUD REST endpoints for Residents, mounted on the chi router with query filtering and consistent error handling.

---

## Steps

### 5.1 Create the Resident handler struct

- File: `internal/handlers/resident.go`
- Create `ResidentHandler` struct holding a `store.ResidentStore`.
- Constructor: `NewResidentHandler(s store.ResidentStore) *ResidentHandler`
- This mirrors the pattern from `RoomHandler` in M3.

### 5.2 Implement handler methods

| Handler | Route | Request Body | Status Codes | Logic |
|---------|-------|-------------|-------------|-------|
| `Create` | `POST /` | `CreateResidentRequest` JSON | 201, 400, 409 | Decode → validate → `store.Create` → if `ErrDuplicate` return 409 → respond 201 |
| `GetByID` | `GET /{id}` | — | 200, 404 | Extract ID → `store.GetByID` → 404 if not found |
| `List` | `GET /` | — | 200 | Parse query params → build `ResidentFilters` → `store.List` → respond 200 |
| `Update` | `PUT /{id}` | Partial/full JSON | 200, 400, 404, 409 | Extract ID → decode → validate → `store.Update` → handle errors |
| `Delete` | `DELETE /{id}` | — | 204, 404 | Extract ID → `store.Delete` → 204 on success |

### 5.3 Handle query parameters for List

Parse these from the URL query string (`r.URL.Query()`):

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `roomId` | string | — | Filter by room assignment |
| `status` | string | — | Filter by `active`, `checked_out`, `pending` |
| `limit` | int | 50 | Max results per page |
| `offset` | int | 0 | Skip N results |

**Parsing tips:**
- Use `r.URL.Query().Get("roomId")` — returns empty string if absent.
- Convert limit/offset with `strconv.Atoi()` — use defaults on error.
- Build a `ResidentFilters` struct with `*string` for optional fields: set pointer only if the query param is non-empty.

### 5.4 Error mapping

Reuse the `RespondError` helper from M3. Map store errors to HTTP codes:

| Store Error | HTTP Code | Response |
|-------------|-----------|----------|
| `ErrNotFound` | 404 | `{"error": "resident not found"}` |
| `ErrDuplicate` | 409 Conflict | `{"error": "email already exists"}` |
| Validation error | 400 | `{"error": "<validation message>"}` |
| JSON decode error | 400 | `{"error": "invalid request body"}` |

### 5.5 Wire resident routes in the router

- Update `internal/router/router.go`:
  - Accept `residentStore store.ResidentStore` as a parameter (or pass a combined store).
  - Add a route group:

    ```
    r.Route("/api/v1/residents", func(r chi.Router) {
        r.Post("/", residentHandler.Create)
        r.Get("/", residentHandler.List)
        r.Get("/{id}", residentHandler.GetByID)
        r.Put("/{id}", residentHandler.Update)
        r.Delete("/{id}", residentHandler.Delete)
    })
    ```

### 5.6 Update `main.go`

- If you used a single `MemoryStore` (Approach A from M4):
  - Already have both stores in one struct. Pass it to the router.
- If you used separate stores:
  - Create `residentStore := store.NewMemoryResidentStore()`
  - Pass both stores to `router.New(roomStore, residentStore)`

### 5.7 Manual testing with curl

```
# Create a resident
curl -X POST http://localhost:8080/api/v1/residents \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","phone":"+1234567890","status":"active"}'

# List all residents
curl http://localhost:8080/api/v1/residents

# List residents by status
curl "http://localhost:8080/api/v1/residents?status=active"

# Get by ID
curl http://localhost:8080/api/v1/residents/<id>

# Update — assign to a room
curl -X PUT http://localhost:8080/api/v1/residents/<id> \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","roomId":"<room-id>","status":"active"}'

# Duplicate email → 409
curl -X POST http://localhost:8080/api/v1/residents \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"john@example.com","status":"active"}'

# Delete
curl -X DELETE http://localhost:8080/api/v1/residents/<id>
```

### 5.8 Cross-resource verification (optional)

If time allows, add a simple check in the `Create`/`Update` handler:
- If `roomId` is provided, call `roomStore.GetByID(roomId)` to verify the room exists.
- Return 400 if the room doesn't exist: `{"error": "room not found"}`
- This requires the `ResidentHandler` to also hold a `RoomStore` reference.

---

## Acceptance Criteria

| # | Check | How to verify |
|---|-------|---------------|
| 1 | `POST /api/v1/residents` creates resident | 201 + resident JSON |
| 2 | `GET /api/v1/residents` lists all | 200 + array |
| 3 | `GET /api/v1/residents?status=active` filters | Only active residents |
| 4 | `GET /api/v1/residents?roomId=<id>` filters | Only residents in that room |
| 5 | `GET /api/v1/residents/:id` returns one | 200 + resident object |
| 6 | `PUT /api/v1/residents/:id` updates | 200 + updated fields |
| 7 | `DELETE /api/v1/residents/:id` removes | 204 + subsequent GET → 404 |
| 8 | Duplicate email returns 409 | curl test |
| 9 | All Room endpoints still work | Regression curl |
| 10 | Health endpoint still works | `GET /health` → 200 |

---

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Duplicate email | 409 Conflict | Standard HTTP code for uniqueness violations |
| Query params | Parsed manually | No extra library needed; chi doesn't auto-parse query params |
| Cross-resource check | Optional / handler-layer | Keeps store layer simple; validates in the layer that has access to both stores |
| Limit default | 50 | Prevents unbounded responses; reasonable for a hostel |

---

## Concepts for Flutter Devs

- **409 Conflict**: Like getting a `DioException` with status 409. Your Flutter app should catch this and show "Email already in use."
- **Query params in Go**: No magic like Dart's `Uri.queryParameters`. You manually call `r.URL.Query().Get("key")`.
- **strconv.Atoi**: Like `int.parse()` in Dart, but returns `(int, error)` instead of throwing.

---

*Previous → [M4: Resident Model & Store](./M4-resident-model-store.md)*  
*Next → [M6: Polish & Documentation](./M6-polish-docs.md)*
