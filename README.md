# Hostel Management System (Go Backend)

A lightweight REST API backend built with Go for managing hostel **rooms** and **boarders** (guests assigned to rooms). The project follows a modular, 8-hour development slice and is intended as a backend for a future Flutter or web client.

## Getting Started

### Prerequisites
- [Go](https://go.dev/doc/install) (v1.21 or later)

### Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd hostel-management
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```

### Running the Server
The server runs on port `8080` by default. Override with the `PORT` environment variable.

```bash
go run ./cmd/server/main.go
```

```bash
go build -o server ./cmd/server/main.go
./server
```

### Health Check
```bash
curl http://localhost:8080/health
```
Expected: `{"status":"ok"}` (or equivalent JSON from `internal/handlers/health.go`).

---

## API Endpoints

### Implemented
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Liveness / health check |

### Planned (see `PRD.md` and `milestones/`)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/rooms` | List rooms (filters TBD, e.g. `status`) |
| `GET` | `/api/v1/rooms/:id` | Get room by ID |
| `POST` | `/api/v1/rooms` | Create a room |
| `PUT` | `/api/v1/rooms/:id` | Update a room |
| `DELETE` | `/api/v1/rooms/:id` | Delete a room |
| `GET` | `/api/v1/boarders` | List boarders (`roomId`, `status`, pagination) |
| `GET` | `/api/v1/boarders/:id` | Get boarder by ID |
| `POST` | `/api/v1/boarders` | Create a boarder |
| `PUT` | `/api/v1/boarders/:id` | Update a boarder |
| `DELETE` | `/api/v1/boarders/:id` | Delete a boarder |

Domain model details: **`PRD.md`** (`Room` includes `rentalPrice`; `Boarder` uses `firstName`, `lastName`, `phone`, required `roomId` referencing `Room.id`).

---

## Project Structure

```text
hostel-management/
├── cmd/
│   └── server/              # Application entry point
├── internal/
│   ├── handlers/            # HTTP handlers (health; room/boarder in later milestones)
│   ├── models/              # Domain entities (room.go, boarder.go)
│   ├── store/               # Persistence (interfaces + implementations)
│   └── router/              # Routes and middleware
├── milestones/              # Milestone plans (M1–M6)
├── PRD.md                   # Product requirements
└── go.mod
```

---

## Milestones
1. **M1:** Project setup and health check  
2. **M2:** Room model and in-memory store  
3. **M3:** Room HTTP API  
4. **M4:** Boarder model and store (`roomId` → Room)  
5. **M5:** Boarder HTTP API  
6. **M6:** Polish, docs, optional SQLite, seed examples  

See the `milestones/` directory for step-by-step checklists.

## License
This project is for educational or prototype use.
