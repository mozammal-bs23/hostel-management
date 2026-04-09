# Hostel Management System (Go Backend)

A lightweight REST API backend built with Go for managing hostel rooms and residents. This project is designed as a modular, 8-hour development slice, serving as a backend for future Flutter or web applications.

## 🚀 Getting Started

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
The server runs on port `8080` by default. You can override this by setting the `PORT` environment variable.

```bash
# Run using Go
go run ./cmd/server/main.go

# Or build and run the binary
go build -o server ./cmd/server/main.go
./server
```

---

## 🛠 API Endpoints

### Health Check
- **GET** `/health`
  - Returns the current status of the API.
  - Response: `{"status": "ok"}`

### Planned Endpoints (Milestones 2-5)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/rooms` | List all rooms (supports filters) |
| `GET` | `/api/v1/rooms/:id` | Get room by ID |
| `POST` | `/api/v1/rooms` | Create a new room |
| `PUT` | `/api/v1/rooms/:id` | Update room details |
| `DELETE` | `/api/v1/rooms/:id` | Delete a room |
| `GET` | `/api/v1/residents` | List all residents |
| `POST` | `/api/v1/residents` | Register a new resident |

---

## 🏗 Project Structure

```text
hostel-management/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Domain entities (Room, Resident)
│   ├── store/           # Data persistence logic (Interfaces & Impls)
│   └── router/          # API route configuration
├── milestones/          # Detailed development tracking
├── PRD.md               # Product Requirements Document
└── go.mod               # Go module definition
```

---

## 📅 Milestones
The project is divided into 6 clear milestones:
1. **M1: Project Setup & Health Check** (Current)
2. **M2: Room Model & Store**
3. **M3: Room API Handlers**
4. **M4: Resident Model & Store**
5. **M5: Resident API Handlers**
6. **M6: Polish & Documentation**

See the `milestones/` directory for detailed progress on each stage.

## 📄 License
This project is for educational/prototype purposes.
