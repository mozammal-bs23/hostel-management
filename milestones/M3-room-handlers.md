# Milestone 3 — Room HTTP Handlers & Router

**Duration:** ~1.5 hours  
**Goal:** Full CRUD REST endpoints for Rooms, wired with chi router, returning proper JSON and status codes.

---

## Steps

### 3.1 Set up the chi router

- File: `internal/router/router.go`
- Create a function `New(roomStore store.RoomStore) http.Handler` that:
  1. Creates a `chi.NewRouter()`
  2. Adds middleware: `chi.middleware.Logger`, `chi.middleware.Recoverer`, `chi.middleware.RequestID`
  3. Adds a JSON content-type middleware (sets `Content-Type: application/json` on all responses)
  4. Mounts `GET /health` (move from main)
  5. Mounts room routes under `/api/v1/rooms` (see 3.3)
  6. Returns the router

### 3.2 Define a standard JSON response helper

- File: `internal/handlers/response.go`
- Create helper functions that handlers will reuse:

  | Helper | Purpose |
  |--------|---------|
  | `RespondJSON(w, statusCode, data)` | Writes `{"data": ...}` with the given status code |
  | `RespondError(w, statusCode, message)` | Writes `{"error": "message"}` with the given status code |

- These ensure every endpoint returns consistent JSON.

### 3.3 Implement Room handlers

- File: `internal/handlers/room.go`
- Create a `RoomHandler` struct that holds a `store.RoomStore` reference.
- Constructor: `NewRoomHandler(store store.RoomStore) *RoomHandler`

**Handler methods:**

| Handler | Route | Status Codes | Logic |
|---------|-------|-------------|-------|
| `Create` | `POST /` | 201, 400 | Decode JSON body → validate → call `store.Create` → respond 201 with room |
| `GetByID` | `GET /{id}` | 200, 404 | Extract `id` from URL param → call `store.GetByID` → 404 if `ErrNotFound` |
| `List` | `GET /` | 200 | Parse query params (`status`, `floor`, `limit`, `offset`) → build `RoomFilters` → call `store.List` |
| `Update` | `PUT /{id}` | 200, 400, 404 | Extract `id` → decode body → validate → call `store.Update` → 404 if not found |
| `Delete` | `DELETE /{id}` | 204, 404 | Extract `id` → call `store.Delete` → 204 (no body) on success |

### 3.4 Wire room routes

Inside `router.New()`, add a route group:

```
r.Route("/api/v1/rooms", func(r chi.Router) {
    r.Post("/", roomHandler.Create)
    r.Get("/", roomHandler.List)
    r.Get("/{id}", roomHandler.GetByID)
    r.Put("/{id}", roomHandler.Update)
    r.Delete("/{id}", roomHandler.Delete)
})
```

Use `chi.URLParam(r, "id")` to extract the `{id}` path parameter inside handlers.

### 3.5 Update `main.go`

- Create the in-memory room store: `roomStore := store.NewMemoryRoomStore()`
- Create the router: `r := router.New(roomStore)`
- Pass `r` to `http.ListenAndServe`
- Remove the old `http.ServeMux` code

### 3.6 Manual testing with curl

Run through every endpoint and verify:

```
# Create
curl -X POST http://localhost:8080/api/v1/rooms \
  -H "Content-Type: application/json" \
  -d '{"name":"Room 101","capacity":2,"floor":1,"status":"available"}'

# List
curl http://localhost:8080/api/v1/rooms

# Get by ID (use an ID from the create response)
curl http://localhost:8080/api/v1/rooms/<id>

# Update
curl -X PUT http://localhost:8080/api/v1/rooms/<id> \
  -H "Content-Type: application/json" \
  -d '{"name":"Room 101 Deluxe","capacity":3,"floor":1,"status":"available"}'

# Delete
curl -X DELETE http://localhost:8080/api/v1/rooms/<id>

# Verify 404
curl http://localhost:8080/api/v1/rooms/<id>
```

### 3.7 Handle edge cases

- **Invalid JSON body:** Return 400 with a clear message.
- **Empty body on POST/PUT:** Return 400.
- **Non-existent ID:** Return 404 (check for `store.ErrNotFound` with `errors.Is`).
- **Validation failure:** Return 400 with the validation error message.

---

## Acceptance Criteria

| # | Check | How to verify |
|---|-------|---------------|
| 1 | `POST /api/v1/rooms` creates a room | curl returns 201 + room with ID |
| 2 | `GET /api/v1/rooms` lists rooms | curl returns 200 + array |
| 3 | `GET /api/v1/rooms/:id` returns room | curl returns 200 + room object |
| 4 | `PUT /api/v1/rooms/:id` updates room | curl returns 200 + updated fields |
| 5 | `DELETE /api/v1/rooms/:id` removes room | curl returns 204 + subsequent GET returns 404 |
| 6 | Invalid body returns 400 | curl with bad JSON returns 400 |
| 7 | Missing room returns 404 | curl with random UUID returns 404 |
| 8 | Health endpoint still works | `GET /health` returns 200 |

---

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Router | chi v5 | Lightweight, stdlib-compatible, good URL params |
| URL params | `{id}` via chi | Clean RESTful paths |
| Response helpers | Centralized `RespondJSON` / `RespondError` | Consistency across all handlers |
| Delete response | 204 No Content | REST convention — nothing to return |
| Dependency injection | Pass store to handler struct | Testable; easy to swap implementations |

---

## Concepts for Flutter Devs

- **chi router** is like Flutter's `GoRouter` — declarative route registration with path parameters.
- **Middleware** in chi is like Dart shelf middleware or Dio interceptors — runs before/after every request.
- **`chi.URLParam(r, "id")`** is like reading `context.pathParameters['id']` in GoRouter.
- **`json.NewDecoder(r.Body).Decode(&req)`** is like `jsonDecode(response.body)` but decodes directly into a typed struct.

---

*Previous → [M2: Room Model & Store](./M2-room-model-store.md)*  
*Next → [M4: Resident Model & Store](./M4-resident-model-store.md)*
